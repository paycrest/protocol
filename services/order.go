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

var CryptoConf = config.CryptoConfig()
var ServerConf = config.ServerConfig()

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
	userOperation, err := utils.InitializeUserOperation(
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
	err = utils.SponsorUserOperation(userOperation, "erc20token")
	if err != nil {
		return fmt.Errorf("failed to sponsor user operation: %w", err)
	}

	// Sign user operation
	_ = utils.SignUserOperation(userOperation)

	// Send user operation
	userOpHash, err := utils.SendUserOperation(userOperation)
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

// RefundOrder refunds sender on canceled lock order
func (s *OrderService) RefundOrder(ctx context.Context, orderID string) error {
	// Fetch lock order from db
	lockOrder, err := db.Client.LockPaymentOrder.
		Query().
		Where(lockpaymentorder.OrderIDEQ(orderID)).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		First(ctx)
	if err != nil {
		return fmt.Errorf("RefundOrder.fetchLockOrder: %w", err)
	}

	// Get default userOperation
	userOperation, err := utils.InitializeUserOperation(
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
	err = utils.SponsorUserOperation(userOperation, "payg")
	if err != nil {
		return fmt.Errorf("RefundOrder.sponsorUserOperation: %w", err)
	}

	// Sign user operation
	_ = utils.SignUserOperation(userOperation)

	// Send user operation
	userOpTxHash, err := utils.SendUserOperation(userOperation)
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

// RevertOrder reverts an initiated payment order on-chain.
func (s *OrderService) RevertOrder(ctx context.Context, order *ent.PaymentOrder, to common.Address) error {
	// Fetch payment order from db
	order, err := db.Client.PaymentOrder.
		Query().
		Where(paymentorder.IDEQ(order.ID)).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithReceiveAddress().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("RevertOrder.fetchOrder: %w", err)
	}

	fees := order.NetworkFeeEstimate.Add(order.SenderFee)
	orderAmountWithFees := order.Amount.Add(fees)

	var amountToRevert decimal.Decimal

	if order.AmountPaid.LessThan(orderAmountWithFees) {
		amountToRevert = order.AmountPaid
	} else if order.AmountPaid.GreaterThan(orderAmountWithFees) {
		amountToRevert = order.AmountPaid.Sub(orderAmountWithFees)
	} else {
		return nil
	}

	// Subtract the network fee from the amount
	amountMinusFee := amountToRevert.Sub(OrderConf.NetworkFee)

	// If amount minus fee is less than zero, return
	if amountMinusFee.LessThan(decimal.Zero) {
		return nil
	}

	// Convert amountMinusFee to big.Int
	amountMinusFeeBigInt := utils.ToSubunit(amountMinusFee, order.Edges.Token.Decimals)

	// Decrypt salt
	saltDecrypted, err := cryptoUtils.DecryptPlain(order.Edges.ReceiveAddress.Salt)
	if err != nil {
		return fmt.Errorf("failed to decrypt salt: %w", err)
	}

	// Get default userOperation
	userOperation, err := utils.InitializeUserOperation(
		ctx, nil, order.Edges.Token.Edges.Network.RPCEndpoint, order.Edges.ReceiveAddress.Address, string(saltDecrypted),
	)
	if err != nil {
		return fmt.Errorf("RevertOrder.initializeUserOperation: %w", err)
	}

	// Create calldata
	calldata, err := s.executeBatchTransferCallData(order, to, amountMinusFeeBigInt)
	if err != nil {
		return fmt.Errorf("RevertOrder.executeBatchTransferCallData: %w", err)
	}
	userOperation.CallData = calldata

	// Sponsor user operation.
	// This will populate the following fields in userOperation: PaymasterAndData, PreVerificationGas, VerificationGasLimit, CallGasLimit
	err = utils.SponsorUserOperation(userOperation, "erc20token")
	if err != nil {
		return fmt.Errorf("RevertOrder.sponsorUserOperation: %w", err)
	}

	// Sign user operation
	_ = utils.SignUserOperation(userOperation)

	// Send user operation
	userOpHash, err := utils.SendUserOperation(userOperation)
	if err != nil {
		return fmt.Errorf("RevertOrder.sendUserOperation: %w", err)
	}

	// Update payment order with userOpHash
	_, err = order.Update().
		SetTxHash(userOpHash).
		SetAmountReturned(amountMinusFee).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("RevertOrder.updateTxHash(%v): %w", userOpHash, err)
	}

	return nil
}

// SettleOrder settles a payment order on-chain.
func (s *OrderService) SettleOrder(ctx context.Context, orderID uuid.UUID) error {
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

	// Get default userOperation
	userOperation, err := utils.InitializeUserOperation(
		ctx, nil, order.Edges.Token.Edges.Network.RPCEndpoint, CryptoConf.AggregatorSmartAccount, CryptoConf.AggregatorSmartAccountSalt,
	)
	if err != nil {
		return fmt.Errorf("SettleOrder.initializeUserOperation: %w", err)
	}

	// Create calldata
	calldata, err := s.executeBatchSettleCallData(ctx, order)
	if err != nil {
		return fmt.Errorf("SettleOrder.settleCallData: %w", err)
	}
	userOperation.CallData = calldata

	// Sponsor user operation.
	// This will populate the following fields in userOperation: PaymasterAndData, PreVerificationGas, VerificationGasLimit, CallGasLimit
	err = utils.SponsorUserOperation(userOperation, "payg")
	if err != nil {
		return fmt.Errorf("SettleOrder.sponsorUserOperation: %w", err)
	}

	// Sign user operation
	_ = utils.SignUserOperation(userOperation)

	// Send user operation
	userOpTxHash, err := utils.SendUserOperation(userOperation)
	if err != nil {
		return fmt.Errorf("SettleOrder.sendUserOperation: %w", err)
	}

	// Update status of lock order
	_, err = order.Update().
		SetTxHash(userOpTxHash).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("SettleOrder.updateTxHash: %w", err)
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

// executeBatchTransferCallData creates the transfer calldata for the execute batch method in the smart account.
func (s *OrderService) executeBatchTransferCallData(order *ent.PaymentOrder, to common.Address, amount *big.Int) ([]byte, error) {
	// Fetch paymaster account
	paymasterAccount, err := s.getPaymasterAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get paymaster account: %w", err)
	}

	if ServerConf.Environment != "staging" && ServerConf.Environment != "production" {
		time.Sleep(5 * time.Second) // TODO: remove in production
	}

	// Create approve data for paymaster contract
	approvePaymasterData, err := s.approveCallData(common.HexToAddress(paymasterAccount), amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create paymaster approve calldata : %w", err)
	}

	// Create approve data for erc20 contract
	approveERC20Data, err := s.approveCallData(to, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create erc20 approve calldata : %w", err)
	}

	// Create transfer data
	transferData, err := s.transferCallData(to, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer calldata: %w", err)
	}

	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse smart account ABI: %w", err)
	}

	executeBatchCallData, err := simpleAccountABI.Pack(
		"executeBatch",
		[]common.Address{
			common.HexToAddress(order.Edges.Token.ContractAddress),
			common.HexToAddress(order.Edges.Token.ContractAddress),
			common.HexToAddress(order.Edges.Token.ContractAddress),
		},
		[][]byte{approvePaymasterData, approveERC20Data, transferData},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack execute ABI: %w", err)
	}

	return executeBatchCallData, nil
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

	if ServerConf.Environment != "staging" && ServerConf.Environment != "production" {
		time.Sleep(5 * time.Second) // TODO: remove in production
	}

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

// transferCallData creates the data for the ERC20 token transfer method
func (s *OrderService) transferCallData(recipient common.Address, amount *big.Int) ([]byte, error) {
	// Create ABI
	erc20ABI, err := abi.JSON(strings.NewReader(contracts.TestTokenMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse erc20 ABI: %w", err)
	}

	// Create calldata
	calldata, err := erc20ABI.Pack("transfer", recipient, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transfer ABI: %w", err)
	}

	return calldata, nil
}

// executeCallData creates the data for the execute method in the smart account.
func (s *OrderService) executeCallData(dest common.Address, value *big.Int, data []byte) ([]byte, error) {
	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse smart account ABI: %w", err)
	}

	executeCallData, err := simpleAccountABI.Pack(
		"execute",
		dest,
		value,
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack execute ABI: %w", err)
	}

	return executeCallData, nil
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
	fee := utils.ToSubunit(OrderConf.NetworkFee, order.Edges.Token.Decimals)
	refundData, err := s.refundCallData(fee, order.OrderID, order.Label)
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData.refundData: %w", err)
	}

	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData.simpleAccountABI: %w", err)
	}

	executeBatchRefundCallData, err := simpleAccountABI.Pack(
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

	return executeBatchRefundCallData, nil
}

// refundCallData creates the data for the refund method
func (s *OrderService) refundCallData(fee *big.Int, orderId, label string) ([]byte, error) {
	paycrestOrderABI, err := abi.JSON(strings.NewReader(contracts.PaycrestOrderMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse PaycrestOrder ABI: %w", err)
	}

	// Generate calldata for refund, orderID, and label should be byte32
	data, err := paycrestOrderABI.Pack(
		"refund",
		fee,
		utils.StringToByte32(orderId),
		utils.StringToByte32(label),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to pack refund ABI: %w", err)
	}

	return data, nil
}

// executeBatchSettleCallData creates the settle calldata for the execute batch method in the smart account.
func (s *OrderService) executeBatchSettleCallData(ctx context.Context, order *ent.LockPaymentOrder) ([]byte, error) {
	// Create approve data for paycrest order contract
	approvePaycrestData, err := s.approveCallData(
		OrderConf.PaycrestOrderContractAddress,
		utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
	)
	if err != nil {
		return nil, fmt.Errorf("executeBatchSettleCallData.approveOrderContract: %w", err)
	}

	// Create settle data
	settleData, err := s.settleCallData(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("executeBatchSettleCallData.refundData: %w", err)
	}

	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("executeBatchSettleCallData.simpleAccountABI: %w", err)
	}

	executeBatchSettleCallData, err := simpleAccountABI.Pack(
		"executeBatch",
		[]common.Address{
			common.HexToAddress(order.Edges.Token.ContractAddress),
			OrderConf.PaycrestOrderContractAddress,
		},
		[][]byte{approvePaycrestData, settleData},
	)
	if err != nil {
		return nil, fmt.Errorf("executeBatchSettledCallData: %w", err)
	}

	return executeBatchSettleCallData, nil
}

// settleCallData creates the data for the settle method in the paycrest order contract
func (s *OrderService) settleCallData(ctx context.Context, order *ent.LockPaymentOrder) ([]byte, error) {
	paycrestOrderABI, err := abi.JSON(strings.NewReader(contracts.PaycrestOrderMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse PaycrestOrder ABI: %w", err)
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
		return nil, fmt.Errorf("failed to fetch provider order token: %w", err)
	}

	var providerAddress string
	for _, addr := range token.Addresses {
		if addr.Network == order.Edges.Token.Edges.Network.Identifier {
			providerAddress = addr.Address
			break
		}
	}

	if providerAddress == "" {
		return nil, fmt.Errorf("failed to fetch provider address: %w", err)
	}

	orderPercent, _ := order.OrderPercent.Float64()

	// Generate calldata for settlement
	data, err := paycrestOrderABI.Pack(
		"settle",
		utils.StringToByte32(order.ID.String()),
		utils.StringToByte32(order.OrderID),
		utils.StringToByte32(order.Label),
		nil, // TODO: remove validators input from contract
		common.HexToAddress(providerAddress),
		uint64(orderPercent),
		order.Edges.Provider.IsPartner,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to pack settle ABI: %w", err)
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

// getPaymasterAccount fetches the paymaster account from stackup
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
