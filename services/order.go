package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/services/contracts"
	db "github.com/paycrest/protocol/storage"
	"github.com/shopspring/decimal"

	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// OrderService provides functionality related to on-chain interactions for payment orders
type OrderService struct{}

type CreateOrderParams struct {
	Token              common.Address
	Amount             *big.Int
	InstitutionCode    [32]byte
	Label              [32]byte
	Rate               *big.Int
	SenderFeeRecipient common.Address
	SenderFee          *big.Int
	RefundAddress      common.Address
	MessageHash        string
}

var (
	fromAddress, privateKey, _ = cryptoUtils.GenerateAccountFromIndex(0)
	CryptoConf                 = config.CryptoConfig()
)

// NewOrderService creates a new instance of OrderService.
func NewOrderService() *OrderService {
	return &OrderService{}
}

// CreateOrder creates a new payment order on-chain.
func (s *OrderService) CreateOrder(ctx context.Context, orderID uuid.UUID) error {
	var err error

	// Fetch payment order from db
	order, err := db.Client.PaymentOrder.
		Query().
		Where(paymentorder.IDEQ(orderID)).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithSenderProfile().
		WithRecipient().
		WithReceiveAddress().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch payment order: %w", err)
	}

	saltDecrypted, err := cryptoUtils.DecryptPlain(order.Edges.ReceiveAddress.Salt)
	if err != nil {
		return fmt.Errorf("failed to decrypt salt: %w", err)
	}

	// Initialize user operation with defaults
	userOperation, err := s.initializeUserOperation(
		ctx, nil, order.Edges.Token.Edges.Network.RPCEndpoint, order.Edges.ReceiveAddress.Address, string(saltDecrypted),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize user operation: %w", err)
	}

	// Create calldata
	calldata, err := s.executeBatchCreateOrderCallData(order)
	if err != nil {
		return fmt.Errorf("failed to create calldata: %w", err)
	}
	userOperation.CallData = calldata

	// Sponsor user operation.
	// This will populate the following fields in userOperation: PaymasterAndData, PreVerificationGas, VerificationGasLimit, CallGasLimit
	err = s.sponsorUserOperation(userOperation, "erc20token")
	if err != nil {
		return fmt.Errorf("failed to sponsor user operation: %w", err)
	}

	// Sign user operation
	_ = s.signUserOperation(userOperation)

	// Send user operation
	userOpHash, err := s.sendUserOperation(userOperation)
	if err != nil {
		return fmt.Errorf("failed to send user operation: %w", err)
	}

	// Update payment order with userOpHash
	_, err = order.Update().
		SetTxHash(userOpHash).
		SetStatus(paymentorder.StatusPending).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update payment order: %w", err)
	}

	paymentOrder, err := db.Client.PaymentOrder.
		Query().
		Where(paymentorder.IDEQ(orderID)).
		WithSenderProfile().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch payment order: %w", err)
	}

	// Send webhook notifcation to sender
	err = utils.SendPaymentOrderWebhook(ctx, paymentOrder)
	if err != nil {
		return fmt.Errorf("CreateOrder.webhook: %w", err)
	}

	return nil
}

// RefundOrder refunds sender on canceled order
func (s *OrderService) RefundOrder(ctx context.Context, lockOrder *ent.LockPaymentOrder) error {
	// Get default userOperation
	userOperation, err := s.initializeUserOperation(
		ctx, nil, lockOrder.Edges.Token.Edges.Network.RPCEndpoint, CryptoConf.AggregatorSmartAccount, CryptoConf.AggregatorSmartAccountSalt,
	)
	if err != nil {
		return fmt.Errorf("RefundOrder.initializeUserOperation: %w", err)
	}

	// Create calldata
	calldata, err := s.executeBatchRefundCallData(lockOrder)
	if err != nil {
		return fmt.Errorf("RefundOrder.refundCallData: %w", err)
	}
	userOperation.CallData = calldata

	// Sponsor user operation.
	// This will populate the following fields in userOperation: PaymasterAndData, PreVerificationGas, VerificationGasLimit, CallGasLimit
	err = s.sponsorUserOperation(userOperation, "payg")
	if err != nil {
		return fmt.Errorf("RefundOrder.sponsorUserOperation: %w", err)
	}

	// Sign user operation
	_ = s.signUserOperation(userOperation)

	// Send user operation
	userOpTxHash, err := s.sendUserOperation(userOperation)
	if err != nil {
		return fmt.Errorf("RefundOrder.sendUserOperation: %w", err)
	}

	// Update status of all lock orders with same order_id
	_, err = db.Client.LockPaymentOrder.
		Update().
		Where(lockpaymentorder.OrderIDEQ(lockOrder.OrderID)).
		SetTxHash(userOpTxHash).
		SetStatus(lockpaymentorder.StatusRefunding).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("RefundOrder.updateTxHash(%v): %w", userOpTxHash, err)
	}

	return nil
}

// SettleOrder settles a payment order on-chain.
func (s *OrderService) SettleOrder(ctx context.Context, client types.RPCClient, orderID uuid.UUID) error {
	var err error

	// Fetch payment order from db
	order, err := db.Client.LockPaymentOrder.
		Query().
		Where(lockpaymentorder.IDEQ(orderID)).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithFulfillment().
		WithProvider().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch lock order: %w", err)
	}

	// Fetch provider address from db
	token, err := db.Client.ProviderOrderToken.
		Query().
		Where(
			providerordertoken.SymbolEQ(order.Edges.Token.Symbol),
			providerordertoken.HasProviderWith(
				providerprofile.IDEQ(order.Edges.Provider.ID),
			),
		).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch provider order token: %w", err)
	}

	var providerAddress string
	for _, addr := range token.Addresses {
		if addr.Network == order.Edges.Token.Edges.Network.Identifier {
			providerAddress = addr.Address
			break
		}
	}

	if providerAddress == "" {
		return fmt.Errorf("failed to fetch provider address: %w", err)
	}

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(order.Edges.Token.Edges.Network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Initialize paycrest order contract
	orderContract, err := contracts.NewPaycrestOrder(OrderConf.PaycrestOrderContractAddress, client.(bind.ContractBackend))
	if err != nil {
		return fmt.Errorf("failed to initialize paycrest order contract: %w", err)
	}

	orderPercent, _ := order.OrderPercent.Float64()

	// Settle order
	tx, err := orderContract.Settle(
		nil,
		utils.StringToByte32(order.ID.String()),
		utils.StringToByte32(order.OrderID),
		utils.StringToByte32(order.Label),
		nil, // TODO: remove validators input from contract
		common.HexToAddress(providerAddress),
		uint64(orderPercent),
		order.Edges.Provider.IsPartner,
	)
	if err != nil {
		return fmt.Errorf("failed to settle order: %w", err)
	}

	_, err = order.Update().
		SetTxHash(tx.Hash().Hex()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update payment order: %w", err)
	}

	return nil
}

// GetSupportedInstitutions fetches the supported institutions by currencyCode.
func (s *OrderService) GetSupportedInstitutions(ctx context.Context, client types.RPCClient, currencyCode string) ([]types.Institution, error) {
	// Connect to RPC endpoint
	var err error
	if client == nil {
		// NOTE: RPCEndpoint defaults to polygon-mumbai until contract is deployed to polygon mainnet.
		client, err = types.NewEthClient("https://polygon-mumbai.g.alchemy.com/v2/zfXjaatj2o5xKkqe0iSvnU9JkKZoiS54")
		if err != nil {
			return nil, fmt.Errorf("GetSupportedInstitutions.NewEthClient: %w", err)
		}
	}

	currency := utils.StringToByte32(currencyCode)

	// Initialize contract filterer
	instance, err := contracts.NewPaycrestOrder(OrderConf.PaycrestOrderContractAddress, client.(bind.ContractBackend))
	if err != nil {
		return nil, fmt.Errorf("GetSupportedInstitutions.NewPaycrestOrder: %w", err)
	}

	institutions, err := instance.GetSupportedInstitutions(nil, currency)
	if err != nil {
		return nil, fmt.Errorf("GetSupportedInstitutions: %w", err)
	}

	supportedInstitution := make([]types.Institution, len(institutions))
	for i, v := range institutions {
		institution := types.Institution{
			Name: utils.Byte32ToString(v.Name),
			Code: utils.Byte32ToString(v.Code),
			Type: "BANK", // NOTE: defaults to bank.
		}
		supportedInstitution[i] = institution
	}

	return supportedInstitution, nil
}

// executeBatchCreateOrderCallData creates the calldata for the execute batch method in the smart account.
func (s *OrderService) executeBatchCreateOrderCallData(order *ent.PaymentOrder) ([]byte, error) {
	// Create approve data for paycrest order contract
	approvePaycrestData, err := s.approveCallData(
		OrderConf.PaycrestOrderContractAddress,
		utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create paycrest approve calldata: %w", err)
	}

	// Fetch paymaster account
	paymasterAccount, err := s.getPaymasterAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get paymaster account: %w", err)
	}

	time.Sleep(5 * time.Second) // TODO: remove in production

	// Create approve data for paymaster contract
	approvePaymasterData, err := s.approveCallData(
		common.HexToAddress(paymasterAccount),
		utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create paymaster approve calldata : %w", err)
	}

	// Create createOrder data
	createOrderData, err := s.createOrderCallData(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create createOrder calldata: %w", err)
	}

	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse smart account ABI: %w", err)
	}

	executeBatchCreateOrderCallData, err := simpleAccountABI.Pack(
		"executeBatch",
		[]common.Address{
			common.HexToAddress(order.Edges.Token.ContractAddress),
			common.HexToAddress(order.Edges.Token.ContractAddress),
			OrderConf.PaycrestOrderContractAddress,
		},
		[][]byte{approvePaymasterData, approvePaycrestData, createOrderData},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack execute ABI: %w", err)
	}

	return executeBatchCreateOrderCallData, nil
}

// approveCallData creates the data for the ERC20 approve method
func (s *OrderService) approveCallData(spender common.Address, amount *big.Int) ([]byte, error) {
	// Create ABI
	erc20ABI, err := abi.JSON(strings.NewReader(contracts.TestTokenMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse erc20 ABI: %w", err)
	}

	// Create calldata
	calldata, err := erc20ABI.Pack("approve", spender, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack approve ABI: %w", err)
	}

	return calldata, nil
}

// createOrderCallData creates the data for the createOrder method
func (s *OrderService) createOrderCallData(order *ent.PaymentOrder) ([]byte, error) {
	// Encrypt recipient details
	encryptedOrderRecipient, err := s.encryptOrderRecipient(order.Edges.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt recipient details: %w", err)
	}

	// Calculate sender fee
	feePerTokenUnit := order.Edges.SenderProfile.FeePerTokenUnit
	senderFee := feePerTokenUnit.Mul(order.Amount)

	// Define params
	params := &CreateOrderParams{
		Token:              common.HexToAddress(order.Edges.Token.ContractAddress),
		Amount:             utils.ToSubunit(order.Amount.Sub(decimal.NewFromFloat(3)), order.Edges.Token.Decimals),
		InstitutionCode:    utils.StringToByte32(order.Edges.Recipient.Institution),
		Label:              utils.StringToByte32(order.Label),
		Rate:               order.Rate.BigInt(),
		SenderFeeRecipient: common.HexToAddress("0x3870419Ba2BBf0127060bCB37f69A1b1C090992B"),
		SenderFee:          senderFee.BigInt(),
		RefundAddress:      common.HexToAddress("0x3870419Ba2BBf0127060bCB37f69A1b1C090992B"),
		MessageHash:        encryptedOrderRecipient,
	}

	// Create ABI
	paycrestOrderABI, err := abi.JSON(strings.NewReader(contracts.PaycrestOrderMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse PaycrestOrder ABI: %w", err)
	}

	// Generate call data
	data, err := paycrestOrderABI.Pack(
		"createOrder",
		params.Token,
		params.Amount,
		params.InstitutionCode,
		params.Label,
		params.Rate,
		params.SenderFeeRecipient,
		params.SenderFee,
		params.RefundAddress,
		params.MessageHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack createOrder ABI: %w", err)
	}

	return data, nil
}

// createAccountCallData creates the data for the createAccount method
func (s *OrderService) createAccountCallData(owner common.Address, salt *big.Int) ([]byte, error) {
	// Create ABI
	accountFactoryABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountFactoryMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse account factory ABI: %w", err)
	}

	// Create calldata
	calldata, err := accountFactoryABI.Pack("createAccount", owner, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to pack createAccount ABI: %w", err)
	}

	return calldata, nil
}

// executeBatchRefundCallData creates the refund calldata for the execute batch method in the smart account.
func (s *OrderService) executeBatchRefundCallData(order *ent.LockPaymentOrder) ([]byte, error) {
	// Create approve data for paycrest order contract
	approvePaycrestData, err := s.approveCallData(
		OrderConf.PaycrestOrderContractAddress,
		utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
	)
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData.approveOrderContract: %w", err)
	}

	// Create refund data
	refundData, err := s.refundCallData(order.OrderID, order.Label)
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData.refundData: %w", err)
	}

	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData.simpleAccountABI: %w", err)
	}

	executeBatchCreateOrderCallData, err := simpleAccountABI.Pack(
		"executeBatch",
		[]common.Address{
			common.HexToAddress(order.Edges.Token.ContractAddress),
			OrderConf.PaycrestOrderContractAddress,
		},
		[][]byte{approvePaycrestData, refundData},
	)
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData: %w", err)
	}

	return executeBatchCreateOrderCallData, nil
}

// refundCallData creates the data for the refund method
func (s *OrderService) refundCallData(orderId, label string) ([]byte, error) {
	// Refund ABI
	paycrestOrderABI, err := abi.JSON(strings.NewReader(contracts.PaycrestOrderMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse PaycrestOrder ABI: %w", err)
	}

	// Generate call data for refund, orderID, and label should be byte32
	data, err := paycrestOrderABI.Pack(
		"refund",
		utils.StringToByte32(orderId),
		utils.StringToByte32(label),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to pack refund ABI: %w", err)
	}

	return data, nil
}

// encryptOrderRecipient encrypts the recipient details
func (s *OrderService) encryptOrderRecipient(recipient *ent.PaymentOrderRecipient) (string, error) {
	message := struct {
		AccountIdentifier string
		AccountName       string
		Institution       string
		ProviderID        string
	}{
		recipient.AccountIdentifier, recipient.AccountName, recipient.Institution, recipient.ProviderID,
	}

	// Encrypt with the public key of the aggregator
	messageCipher, err := cryptoUtils.PublicKeyEncryptJSON(message, CryptoConf.AggregatorPublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt message: %w", err)
	}

	return fmt.Sprintf("0x%x", messageCipher), nil
}

func (s *OrderService) getPaymasterAccount() (string, error) {
	client, err := rpc.Dial(OrderConf.PaymasterURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	requestParams := []interface{}{
		OrderConf.EntryPointContractAddress.Hex(),
	}

	var result json.RawMessage
	err = client.Call(&result, "pm_accounts", requestParams...)
	if err != nil {
		return "", fmt.Errorf("RPC error: %w", err)
	}

	var response []string
	err = json.Unmarshal(result, &response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response[0], nil
}

// sponsorUserOperation sponsors the user operation
// ref: https://docs.stackup.sh/docs/paymaster-api-rpc-methods#pm_sponsoruseroperation
func (s *OrderService) sponsorUserOperation(userOp *userop.UserOperation, mode string) error {
	client, err := rpc.Dial(OrderConf.PaymasterURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	var payload map[string]interface{}

	switch mode {
	case "payg":
		payload = map[string]interface{}{
			"type": "payg",
		}
	case "erc20token":
		payload = map[string]interface{}{
			"type":  "erc20token",
			"token": "0x3870419Ba2BBf0127060bCB37f69A1b1C090992B",
		}
	default:
		return fmt.Errorf("invalid mode")
	}

	requestParams := []interface{}{
		userOp,
		OrderConf.EntryPointContractAddress.Hex(),
		payload,
	}

	// op, _ := userOp.MarshalJSON()
	// fmt.Println(string(op))

	var result json.RawMessage
	err = client.Call(&result, "pm_sponsorUserOperation", requestParams...)
	if err != nil {
		return fmt.Errorf("RPC error: %w", err)
	}

	type Response struct {
		PaymasterAndData     string `json:"paymasterAndData"     mapstructure:"paymasterAndData"`
		PreVerificationGas   string `json:"preVerificationGas"   mapstructure:"preVerificationGas"`
		VerificationGasLimit string `json:"verificationGasLimit" mapstructure:"verificationGasLimit"`
		CallGasLimit         string `json:"callGasLimit"         mapstructure:"callGasLimit"`
	}

	var response Response
	err = json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	userOp.CallGasLimit, _ = new(big.Int).SetString(response.CallGasLimit, 0)
	userOp.VerificationGasLimit, _ = new(big.Int).SetString(response.VerificationGasLimit, 0)
	userOp.PreVerificationGas, _ = new(big.Int).SetString(response.PreVerificationGas, 0)
	userOp.PaymasterAndData = common.FromHex(response.PaymasterAndData)

	return nil
}

// sendUserOperation sends the user operation
func (s *OrderService) sendUserOperation(userOp *userop.UserOperation) (string, error) {
	client, err := rpc.Dial(OrderConf.BundlerRPCURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	requestParams := []interface{}{
		userOp,
		OrderConf.EntryPointContractAddress.Hex(),
	}

	op, _ := userOp.MarshalJSON()
	fmt.Println(string(op))

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

// // getUserOperationStatus returns the status of the user operation
// func (s *OrderService) getUserOperationStatus(userOpHash string) (bool, error) {
// 	client, err := rpc.Dial(OrderConf.BundlerRPCURL)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to connect to RPC client: %w", err)
// 	}

// 	requestParams := []interface{}{
// 		userOpHash,
// 	}

// 	var result json.RawMessage
// 	err = client.Call(&result, "eth_getUserOperationReceipt", requestParams)
// 	if err != nil {
// 		return false, fmt.Errorf("RPC error: %w", err)
// 	}

// 	var userOpStatus map[string]interface{}
// 	err = json.Unmarshal(result, &userOpStatus)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to unmarshal response: %w", err)
// 	}

// 	return userOpStatus["success"].(bool), nil
// }

// eip1559GasPrice computes the EIP1559 gas price
func (s *OrderService) eip1559GasPrice(ctx context.Context, client types.RPCClient) (maxFeePerGas, maxPriorityFeePerGas *big.Int) {
	tip, _ := client.SuggestGasTipCap(ctx)
	latestHeader, _ := client.HeaderByNumber(ctx, nil)

	buffer := new(big.Int).Mul(tip, big.NewInt(13)).Div(tip, big.NewInt(100))
	maxPriorityFeePerGas = new(big.Int).Add(tip, buffer)

	if latestHeader.BaseFee != nil {
		maxFeePerGas = new(big.Int).
			Mul(latestHeader.BaseFee, big.NewInt(2)).
			Add(latestHeader.BaseFee, maxPriorityFeePerGas)
	} else {
		maxFeePerGas = maxPriorityFeePerGas
	}

	return maxFeePerGas, maxPriorityFeePerGas
}

// Initialize user operation with defaults
func (s *OrderService) initializeUserOperation(ctx context.Context, client types.RPCClient, rpcUrl, sender, salt string) (*userop.UserOperation, error) {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(rpcUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Build user operation
	userOperation := &userop.UserOperation{
		Sender:               common.HexToAddress(sender),
		Nonce:                big.NewInt(0),
		InitCode:             common.FromHex("0x"),
		CallData:             common.FromHex("0x"),
		CallGasLimit:         big.NewInt(350000),
		VerificationGasLimit: big.NewInt(300000),
		PreVerificationGas:   big.NewInt(100000),
		MaxFeePerGas:         big.NewInt(50000),
		MaxPriorityFeePerGas: big.NewInt(1000),
		PaymasterAndData:     common.FromHex("0x"),
		Signature:            common.FromHex("0xa925dcc5e5131636e244d4405334c25f034ebdd85c0cb12e8cdb13c15249c2d466d0bade18e2cafd3513497f7f968dcbb63e519acd9b76dcae7acd61f11aa8421b"),
	}

	// Get nonce
	nonce, err := client.PendingNonceAt(ctx, *fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}
	userOperation.Nonce = big.NewInt(int64(nonce))

	// Create initcode
	code, err := client.CodeAt(ctx, userOperation.Sender, nil)
	if err != nil {
		return nil, err
	}

	if len(code) == 0 {
		// address does not exist yet
		salt, _ := new(big.Int).SetString(salt, 10)

		createAccountCallData, err := s.createAccountCallData(*fromAddress, salt)
		if err != nil {
			return nil, fmt.Errorf("failed to create init code: %w", err)
		}

		var factoryAddress [20]byte
		copy(factoryAddress[:], common.HexToAddress("0x9406Cc6185a346906296840746125a0E44976454").Bytes())

		userOperation.InitCode = append(factoryAddress[:], createAccountCallData...)
	}

	// Set gas fees
	maxFeePerGas, maxPriorityFeePerGas := s.eip1559GasPrice(ctx, client)
	userOperation.MaxFeePerGas = maxFeePerGas
	userOperation.MaxPriorityFeePerGas = maxPriorityFeePerGas

	return userOperation, nil
}

// signUserOperation signs the user operation
func (s *OrderService) signUserOperation(userOperation *userop.UserOperation) error {
	// Sign user operation
	userOpHash := userOperation.GetUserOpHash(
		OrderConf.EntryPointContractAddress,
		big.NewInt(137),
	)

	signature, err := utils.PersonalSign(string(userOpHash[:]), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign user operation: %w", err)
	}
	userOperation.Signature = signature

	return nil
}
