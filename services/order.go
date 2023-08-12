package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/config"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"

	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/paycrest/paycrest-protocol/utils"
	cryptoUtils "github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// OrderService provides functionality related to on-chain interactions for payment orders
type OrderService struct{}

type CreateOrderParams struct {
	Token              common.Address
	Amount             *big.Int
	RefundAddress      common.Address
	SenderFeeRecipient common.Address
	SenderFee          *big.Int
	Rate               *big.Int
	InstitutionCode    [32]byte
	MessageHash        string
}

var conf = config.OrderConfig()

// NewOrderService creates a new instance of OrderService.
func NewOrderService() *OrderService {
	return &OrderService{}
}

// CreateOrder creates a new payment order on-chain.
func (s *OrderService) CreateOrder(ctx context.Context, client types.RPCClient, orderID uuid.UUID) error {
	var err error

	// Fetch payment order from db
	order, err := db.Client.PaymentOrder.
		Query().
		Where(paymentorder.IDEQ(orderID)).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithRecipient().
		WithReceiveAddress().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch payment order: %w", err)
	}

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(order.Edges.Token.Edges.Network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	fromAddress, privateKey, _ := utils.GetMasterAccount()

	// Initialize user operation with defaults
	userOperation := &userop.UserOperation{
		Sender:               common.HexToAddress(order.Edges.ReceiveAddress.Address),
		Nonce:                big.NewInt(0),
		InitCode:             common.FromHex("0x"),
		CallData:             common.FromHex("0x"),
		CallGasLimit:         big.NewInt(0),
		VerificationGasLimit: big.NewInt(0),
		PreVerificationGas:   big.NewInt(0),
		MaxFeePerGas:         big.NewInt(0),
		MaxPriorityFeePerGas: big.NewInt(0),
		PaymasterAndData:     common.FromHex("0x"),
		Signature:            common.FromHex("0x"),
	}

	// Get nonce
	nonce, err := client.PendingNonceAt(ctx, *fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}
	userOperation.Nonce = big.NewInt(int64(nonce))

	// Create initcode
	code, err := client.CodeAt(ctx, userOperation.Sender, nil)
	if err != nil {
		return err
	}

	if len(code) == 0 {
		// address does not exist yet
		createAccountCallData, err := s.createAccountCallData(*fromAddress, big.NewInt(0))
		if err != nil {
			return fmt.Errorf("failed to create init code: %w", err)
		}

		var factoryAddress [20]byte
		copy(factoryAddress[:], common.HexToAddress("0x9406Cc6185a346906296840746125a0E44976454").Bytes())

		userOperation.InitCode = append(factoryAddress[:], createAccountCallData...)
	}

	// Create calldata
	calldata, err := s.executeBatchCallData(order)
	if err != nil {
		return fmt.Errorf("failed to create calldata: %w", err)
	}
	userOperation.CallData = calldata

	// Sponsor user operation.
	// This will populate the following fields in userOperation: PaymasterAndData, PreVerificationGas, VerificationGasLimit, CallGasLimit
	err = s.sponsorUserOperation(userOperation)
	if err != nil {
		return fmt.Errorf("failed to sponsor user operation: %w", err)
	}

	// Set gas fees
	gasPrice, _ := client.SuggestGasPrice(ctx)
	userOperation.MaxFeePerGas = big.NewInt(0).Mul(gasPrice, userOperation.CallGasLimit)
	userOperation.MaxPriorityFeePerGas = big.NewInt(0).Mul(gasPrice, big.NewInt(110)) // 110%

	// Sign user operation
	userOpHash := userOperation.GetUserOpHash(conf.EntryPointContractAddress, big.NewInt(order.Edges.Token.Edges.Network.ChainID)).Hex()

	signature, err := crypto.Sign([]byte(userOpHash), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign user operation: %w", err)
	}
	userOperation.Signature = signature

	// Send user operation
	_, err = s.sendUserOperation(userOperation)
	if err != nil {
		return fmt.Errorf("failed to send user operation: %w", err)
	}

	// update payment order with userOpHash
	_, err = order.Update().
		SetTxHash(userOpHash).
		SetStatus(paymentorder.StatusPending).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update payment order: %w", err)
	}

	return nil
}

// executeBatchCallData creates the calldata for the execute batch method in the smart account.
func (s *OrderService) executeBatchCallData(order *ent.PaymentOrder) ([]byte, error) {
	// Create approve data
	approveData, err := s.approveCallData(conf.PaycrestOrderContractAddress, order.Amount.BigInt())
	if err != nil {
		return nil, fmt.Errorf("failed to create approve calldata: %w", err)
	}

	// Create createOrder data
	createOrderData, err := s.createOrderCallData(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create createOrder calldata: %w", err)
	}

	// Create executeBatch data
	calls := []ethereum.CallMsg{
		{
			To:    &conf.PaycrestOrderContractAddress,
			Value: big.NewInt(0),
			Data:  approveData,
		},
		{
			To:    &conf.PaycrestOrderContractAddress,
			Value: big.NewInt(0),
			Data:  createOrderData,
		},
	}

	executeBatchABI, err := abi.JSON(strings.NewReader("executeBatch(address[],uint256[],bytes[])"))
	if err != nil {
		return nil, fmt.Errorf("failed to create execute ABI: %w", err)
	}

	executeBatchCallData, err := executeBatchABI.Pack("executeBatch", calls)
	if err != nil {
		return nil, fmt.Errorf("failed to pack execute ABI: %w", err)
	}

	return executeBatchCallData, nil
}

// approveCallData creates the data for the ERC20 approve method
func (s *OrderService) approveCallData(spender common.Address, amount *big.Int) ([]byte, error) {
	// Define params
	params := struct {
		Spender common.Address
		Amount  *big.Int
	}{spender, amount}

	// Create ABI
	approveABI, err := abi.JSON(strings.NewReader("approve(address,uint256)"))
	if err != nil {
		return nil, fmt.Errorf("failed to create approve ABI: %w", err)
	}

	// Create calldata
	calldata, err := approveABI.Pack("approve", params)
	if err != nil {
		return nil, fmt.Errorf("failed to pack approve ABI: %w", err)
	}

	return calldata, nil
}

// createOrderCallData creates the data for the createOrder method
func (s *OrderService) createOrderCallData(order *ent.PaymentOrder) ([]byte, error) {
	fromAddress, _, _ := utils.GetMasterAccount()

	// Encrypt recipient details
	encryptedOrderRecipient, err := s.encryptOrderRecipient(order.Edges.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt recipient details: %w", err)
	}

	// Define params
	params := &CreateOrderParams{
		Token:              common.HexToAddress(order.Edges.Token.ContractAddress),
		Amount:             order.Amount.BigInt(),
		RefundAddress:      *fromAddress,
		SenderFeeRecipient: *fromAddress,
		SenderFee:          big.NewInt(0),
		Rate:               big.NewInt(0),
		InstitutionCode:    utils.StringTo32Byte(order.Edges.Recipient.Institution),
		MessageHash:        encryptedOrderRecipient,
	}

	// Create ABI
	createOrderABI, err := abi.JSON(strings.NewReader("createOrder(address,uint256,address,address,uint256,uint96,bytes32,string)"))
	if err != nil {
		return nil, fmt.Errorf("failed to create createOrder ABI: %w", err)
	}

	// Generate call data
	data, err := createOrderABI.Pack("createOrder", params)
	if err != nil {
		return nil, fmt.Errorf("failed to pack createOrder ABI: %w", err)
	}

	return data, nil
}

// createAccountCallData creates the data for the createAccount method
func (s *OrderService) createAccountCallData(address common.Address, salt *big.Int) ([]byte, error) {
	// Define params
	params := struct {
		Address common.Address
		Salt    *big.Int
	}{address, salt}

	// Create ABI
	createAccountABI, err := abi.JSON(strings.NewReader("createAccount(address,uint256)"))
	if err != nil {
		return nil, fmt.Errorf("failed to create createAccount ABI: %w", err)
	}

	// Create calldata
	calldata, err := createAccountABI.Pack("createAccount", params)
	if err != nil {
		return nil, fmt.Errorf("failed to pack createAccount ABI: %w", err)
	}

	return calldata, nil
}

// encryptOrderRecipient encrypts the recipient details
func (s *OrderService) encryptOrderRecipient(recipient *ent.PaymentOrderRecipient) (string, error) {
	message := struct {
		AccountIdentifier string
		AccountName       string
		Institution       string
	}{
		recipient.AccountIdentifier, recipient.AccountName, recipient.Institution,
	}

	messageCipher, err := cryptoUtils.EncryptJSON(message)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt message: %w", err)
	}

	return fmt.Sprintf("0x%x", messageCipher), nil
}

// sponsorUserOperation sponsors the user operation
// ref: https://docs.stackup.sh/docs/paymaster-api-rpc-methods#pm_sponsoruseroperation
func (s *OrderService) sponsorUserOperation(userOp *userop.UserOperation) error {
	client, err := rpc.Dial(conf.PaymasterURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	requestParams := []interface{}{
		userOp,
		conf.EntryPointContractAddress.Hex(),
		map[string]interface{}{
			"type": "payg",
		},
	}

	var result json.RawMessage
	err = client.Call(&result, "pm_sponsorUserOperation", requestParams...)
	if err != nil {
		return fmt.Errorf("RPC error: %w", err)
	}

	type Response struct {
		PaymasterAndData     []byte
		PreVerificationGas   *big.Int
		VerificationGasLimit *big.Int
		CallGasLimit         *big.Int
	}

	var response Response
	err = json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	userOp.CallGasLimit = response.CallGasLimit
	userOp.VerificationGasLimit = response.VerificationGasLimit
	userOp.PreVerificationGas = response.PreVerificationGas
	userOp.PaymasterAndData = response.PaymasterAndData

	return nil
}

// sendUserOperation sends the user operation
func (s *OrderService) sendUserOperation(userOp *userop.UserOperation) (string, error) {
	client, err := rpc.Dial(conf.BundlerRPCURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	requestParams := []interface{}{
		userOp,
		conf.EntryPointContractAddress.Hex(),
	}

	var result json.RawMessage
	err = client.Call(&result, "eth_sendUserOperation", requestParams...)
	if err != nil {
		return "", fmt.Errorf("RPC error: %w", err)
	}

	var userOpHash string
	err = json.Unmarshal(result, &userOpHash)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return userOpHash, nil
}
