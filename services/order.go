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
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/paycrest/paycrest-protocol/utils"
	cryptoUtils "github.com/paycrest/paycrest-protocol/utils/crypto"
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

// UserOperation represents an EIP-4337 style transaction for a smart contract account.
type UserOperation struct {
	Sender               common.Address `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce                *big.Int       `json:"nonce"                mapstructure:"nonce"                validate:"required"`
	InitCode             []byte         `json:"initCode"             mapstructure:"initCode"             validate:"required"`
	CallData             []byte         `json:"callData"             mapstructure:"callData"             validate:"required"`
	CallGasLimit         *big.Int       `json:"callGasLimit"         mapstructure:"callGasLimit"         validate:"required"`
	VerificationGasLimit *big.Int       `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   *big.Int       `json:"preVerificationGas"   mapstructure:"preVerificationGas"   validate:"required"`
	MaxFeePerGas         *big.Int       `json:"maxFeePerGas"         mapstructure:"maxFeePerGas"         validate:"required"`
	MaxPriorityFeePerGas *big.Int       `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte         `json:"paymasterAndData"     mapstructure:"paymasterAndData"     validate:"required"`
	Signature            []byte         `json:"signature"            mapstructure:"signature"            validate:"required"`
}

var conf = config.OrderConfig()

// NewOrderService creates a new instance of OrderService.
func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, client types.RPCClient, order *ent.PaymentOrder) error {
	fromAddress, _, _ := utils.GetMasterAccount()

	// Connect to RPC endpoint
	var err error
	if client == nil {
		client, err = types.NewEthClient(order.Edges.Token.Edges.Network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Create approve data
	approveData, err := s.approveCallData(conf.PaycrestOrderContractAddress, order.Amount.BigInt())
	if err != nil {
		return fmt.Errorf("failed to create approve calldata: %w", err)
	}

	// Create createOrder data
	createOrderData, err := s.createOrderCallData(order)
	if err != nil {
		return fmt.Errorf("failed to create createOrder calldata: %w", err)
	}

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
		return fmt.Errorf("failed to create execute ABI: %w", err)
	}

	calldata, err := executeBatchABI.Pack("executeBatch", calls)
	if err != nil {
		return fmt.Errorf("failed to pack execute ABI: %w", err)
	}

	nonce, err := client.PendingNonceAt(ctx, *fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	// TODO: build user operation
	userOperation := &UserOperation{
		Sender:               common.HexToAddress(order.Edges.ReceiveAddress.Address),
		Nonce:                big.NewInt(int64(nonce)),
		InitCode:             common.FromHex("0x"),
		CallData:             calldata,
		CallGasLimit:         big.NewInt(35000),
		VerificationGasLimit: big.NewInt(70000),
		PreVerificationGas:   big.NewInt(21000),
		MaxFeePerGas:         gasPrice,
		MaxPriorityFeePerGas: gasPrice,
		PaymasterAndData:     common.FromHex("0x"),
		Signature:            common.FromHex("0x"),
	}

	_ = userOperation

	// TODO: send user operation

	return nil
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

// getPaymasterAccount returns the paymaster account address
// ref: https://docs.stackup.sh/docs/paymaster-api-rpc-methods#pm_accounts
func (s *OrderService) getPaymasterAccount() (*common.Address, error) {
	client, err := rpc.Dial(conf.PaymasterURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	requestParams := []interface{}{
		conf.EntryPointContractAddress.Hex(),
	}

	var result json.RawMessage
	err = client.Call(&result, "pm_accounts", requestParams...)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %w", err)
	}

	var accounts []string
	json.Unmarshal(result, &accounts)

	address := common.HexToAddress(accounts[0])

	return &address, nil
}

// sponsorUserOperation sponsors the user operation
// ref: https://docs.stackup.sh/docs/paymaster-api-rpc-methods#pm_sponsoruseroperation
func (s *OrderService) sponsorUserOperation(userOp *UserOperation) (*UserOperation, error) {
	client, err := rpc.Dial(conf.PaymasterURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC client: %w", err)
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
		return nil, fmt.Errorf("RPC error: %w", err)
	}

	type Response struct {
		PaymasterAndData     []byte
		PreVerificationGas   *big.Int
		VerificationGasLimit *big.Int
		CallGasLimit         *big.Int
	}

	var response Response

	json.Unmarshal(result, &response)

	userOp.CallGasLimit = response.CallGasLimit
	userOp.VerificationGasLimit = response.VerificationGasLimit
	userOp.PreVerificationGas = response.PreVerificationGas
	userOp.PaymasterAndData = response.PaymasterAndData

	return userOp, nil
}
