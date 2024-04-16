package services

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/services/contracts"
	db "github.com/paycrest/protocol/storage"
	"github.com/shopspring/decimal"

	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/receiveaddress"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	"github.com/paycrest/protocol/utils/logger"
)

type CreateOrderParams struct {
	Token              common.Address
	Amount             *big.Int
	InstitutionCode    [32]byte
	Rate               *big.Int
	SenderFeeRecipient common.Address
	SenderFee          *big.Int
	RefundAddress      common.Address
	MessageHash        string
}

var CryptoConf = config.CryptoConfig()
var ServerConf = config.ServerConfig()

// Order provides an interface for the OrderService
type Order interface {
	CreateOrder(ctx context.Context, orderID uuid.UUID) error
	RefundOrder(ctx context.Context, orderID string) error
	RevertOrder(ctx context.Context, order *ent.PaymentOrder) error
	SettleOrder(ctx context.Context, orderID uuid.UUID) error
	GetSupportedInstitutions(ctx context.Context, client types.RPCClient, currencyCode string) ([]types.Institution, error)
}

// OrderService provides functionality related to on-chain interactions for payment orders
type OrderService struct{}

// NewOrderService creates a new instance of OrderService.
func NewOrderService() Order {
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
	if ServerConf.Environment != "production" {
		err = utils.SponsorUserOperation(userOperation, "erc20token", order.Edges.Token.ContractAddress, order.Edges.Token.Edges.Network.ChainID)
	} else {
		err = utils.SponsorUserOperation(userOperation, "payg", "", order.Edges.Token.Edges.Network.ChainID)
	}
	if err != nil {
		return fmt.Errorf("failed to sponsor user operation: %w", err)
	}

	// Sign user operation
	err = utils.SignUserOperation(userOperation, order.Edges.Token.Edges.Network.ChainID)
	if err != nil {
		return fmt.Errorf("failed to sign user operation: %w", err)
	}

	// Send user operation
	txHash, blockNumber, err := utils.SendUserOperation(userOperation, order.Edges.Token.Edges.Network.ChainID)
	if err != nil {
		return fmt.Errorf("failed to send user operation: %w", err)
	}

	// Update payment order with userOpHash
	_, err = order.Update().
		SetTxHash(txHash).
		SetBlockNumber(blockNumber).
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
		Where(lockpaymentorder.GatewayIDEQ(orderID)).
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
	if ServerConf.Environment != "production" {
		err = utils.SponsorUserOperation(userOperation, "erc20token", lockOrder.Edges.Token.ContractAddress, lockOrder.Edges.Token.Edges.Network.ChainID)
	} else {
		err = utils.SponsorUserOperation(userOperation, "payg", "", lockOrder.Edges.Token.Edges.Network.ChainID)
	}
	if err != nil {
		return fmt.Errorf("RefundOrder.sponsorUserOperation: %w", err)
	}

	// Sign user operation
	_ = utils.SignUserOperation(userOperation, lockOrder.Edges.Token.Edges.Network.ChainID)

	// Send user operation
	txHash, blockNumber, err := utils.SendUserOperation(userOperation, lockOrder.Edges.Token.Edges.Network.ChainID)
	if err != nil {
		return fmt.Errorf("RefundOrder.sendUserOperation: %w", err)
	}

	// Update status of all lock orders with same order_id
	_, err = db.Client.LockPaymentOrder.
		Update().
		Where(lockpaymentorder.GatewayIDEQ(lockOrder.GatewayID)).
		SetTxHash(txHash).
		SetBlockNumber(blockNumber).
		SetStatus(lockpaymentorder.StatusRefunded).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("RefundOrder.updateTxHash(%v): %w", txHash, err)
	}

	return nil
}

// RevertOrder reverts an initiated payment order on-chain.
func (s *OrderService) RevertOrder(ctx context.Context, order *ent.PaymentOrder) error {
	if !order.AmountReturned.Equal(decimal.Zero) {
		return nil
	}

	orderIDPrefix := strings.Split(order.ID.String(), "-")[0]

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
		return fmt.Errorf("%s - RevertOrder.fetchOrder: %w", orderIDPrefix, err)
	}

	fees := order.NetworkFee.Add(order.SenderFee).Add(order.ProtocolFee)
	orderAmountWithFees := order.Amount.Add(fees)

	var amountToRevert decimal.Decimal

	if order.AmountPaid.LessThan(orderAmountWithFees) {
		amountToRevert = order.AmountPaid
	} else if order.AmountPaid.GreaterThan(orderAmountWithFees) {
		amountToRevert = order.AmountPaid.Sub(orderAmountWithFees)
	} else if order.Status == paymentorder.StatusInitiated && order.Edges.ReceiveAddress.Status == receiveaddress.StatusUsed && order.UpdatedAt.Before(time.Now().Add(-5*time.Minute)) {
		amountToRevert = order.AmountPaid
	} else {
		return nil
	}

	if amountToRevert.Equal(decimal.Zero) {
		return nil
	}

	// Subtract the network fee from the amount
	amountMinusFee := amountToRevert.Sub(order.NetworkFee)

	// If amount minus fee is less than zero, return
	if amountMinusFee.LessThan(decimal.Zero) {
		return nil
	}

	logger.Infof(orderIDPrefix, amountMinusFee.String())

	// Convert amountMinusFee to big.Int
	amountMinusFeeBigInt := utils.ToSubunit(amountMinusFee, order.Edges.Token.Decimals)

	// Decrypt salt
	saltDecrypted, err := cryptoUtils.DecryptPlain(order.Edges.ReceiveAddress.Salt)
	if err != nil {
		return fmt.Errorf("%s - RevertOrder.DecryptPlain: %w", orderIDPrefix, err)
	}

	// Get default userOperation
	userOperation, err := utils.InitializeUserOperation(
		ctx, nil, order.Edges.Token.Edges.Network.RPCEndpoint, order.Edges.ReceiveAddress.Address, string(saltDecrypted),
	)
	if err != nil {
		return fmt.Errorf("%s - RevertOrder.initializeUserOperation: %w", orderIDPrefix, err)
	}

	// Create calldata
	calldata, err := s.executeBatchTransferCallData(order, common.HexToAddress(order.FromAddress), amountMinusFeeBigInt)
	if err != nil {
		return fmt.Errorf("RevertOrder.executeBatchTransferCallData: %w", err)
	}
	userOperation.CallData = calldata

	// Sponsor user operation.
	// This will populate the following fields in userOperation: PaymasterAndData, PreVerificationGas, VerificationGasLimit, CallGasLimit
	if ServerConf.Environment != "production" {
		err = utils.SponsorUserOperation(userOperation, "erc20token", order.Edges.Token.ContractAddress, order.Edges.Token.Edges.Network.ChainID)
	} else {
		err = utils.SponsorUserOperation(userOperation, "payg", "", order.Edges.Token.Edges.Network.ChainID)
	}
	if err != nil {
		return fmt.Errorf("%s - RevertOrder.sponsorUserOperation: %w", orderIDPrefix, err)
	}

	// Sign user operation
	_ = utils.SignUserOperation(userOperation, order.Edges.Token.Edges.Network.ChainID)

	// Send user operation
	txHash, blockNumber, err := utils.SendUserOperation(userOperation, order.Edges.Token.Edges.Network.ChainID)
	if err != nil {
		return fmt.Errorf("%s - RevertOrder.sendUserOperation: %w", orderIDPrefix, err)
	}

	// Update payment order
	_, err = order.Update().
		SetTxHash(txHash).
		SetBlockNumber(blockNumber).
		SetAmountReturned(amountMinusFee).
		SetStatus(paymentorder.StatusReverted).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("RevertOrder.updateTxHash(%v): %w", txHash, err)
	}

	// Send webhook notifcation to sender
	order.Status = paymentorder.StatusReverted
	err = utils.SendPaymentOrderWebhook(ctx, order)
	if err != nil {
		return fmt.Errorf("RevertOrder.webhook: %v", err)
	}

	return nil
}

// SettleOrder settles a payment order on-chain.
func (s *OrderService) SettleOrder(ctx context.Context, orderID uuid.UUID) error {
	var err error

	orderIDPrefix := strings.Split(orderID.String(), "-")[0]

	// Fetch payment order from db
	order, err := db.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.IDEQ(orderID),
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusValidated),
			lockpaymentorder.HasFulfillmentWith(
				lockorderfulfillment.ValidationStatusEQ(lockorderfulfillment.ValidationStatusSuccess),
			),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithProvider().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("%s - SettleOrder.fetchOrder: %w", orderIDPrefix, err)
	}

	// Get default userOperation
	userOperation, err := utils.InitializeUserOperation(
		ctx, nil, order.Edges.Token.Edges.Network.RPCEndpoint, CryptoConf.AggregatorSmartAccount, CryptoConf.AggregatorSmartAccountSalt,
	)
	if err != nil {
		return fmt.Errorf("%s - SettleOrder.initializeUserOperation: %w", orderIDPrefix, err)
	}

	// Create calldata
	calldata, err := s.executeBatchSettleCallData(ctx, order)
	if err != nil {
		return fmt.Errorf("%s - SettleOrder.executeBatchSettleCallData: %w", orderIDPrefix, err)
	}
	userOperation.CallData = calldata

	// Sponsor user operation.
	// This will populate the following fields in userOperation: PaymasterAndData, PreVerificationGas, VerificationGasLimit, CallGasLimit
	if ServerConf.Environment != "production" {
		err = utils.SponsorUserOperation(userOperation, "erc20token", order.Edges.Token.ContractAddress, order.Edges.Token.Edges.Network.ChainID)
	} else {
		err = utils.SponsorUserOperation(userOperation, "payg", "", order.Edges.Token.Edges.Network.ChainID)
	}
	if err != nil {
		return fmt.Errorf("%s - SettleOrder.sponsorUserOperation: %w", orderIDPrefix, err)
	}

	// Sign user operation
	_ = utils.SignUserOperation(userOperation, order.Edges.Token.Edges.Network.ChainID)

	// Send user operation
	txHash, blockNumber, err := utils.SendUserOperation(userOperation, order.Edges.Token.Edges.Network.ChainID)
	if err != nil {
		return fmt.Errorf("%s - SettleOrder.sendUserOperation: %w", orderIDPrefix, err)
	}

	// Update status of lock order
	_, err = order.Update().
		SetTxHash(txHash).
		SetBlockNumber(blockNumber).
		SetStatus(lockpaymentorder.StatusSettled).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("%s - SettleOrder.updateTxHash: %w", orderIDPrefix, err)
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
	instance, err := contracts.NewGateway(OrderConf.GatewayContractAddress, client.(bind.ContractBackend))
	if err != nil {
		return nil, fmt.Errorf("GetSupportedInstitutions.NewGatewayOrder: %w", err)
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
	paymasterAccount, err := utils.GetPaymasterAccount(order.Edges.Token.Edges.Network.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get paymaster account: %w", err)
	}

	if ServerConf.Environment != "staging" && ServerConf.Environment != "production" {
		time.Sleep(5 * time.Second)
	}

	// Create approve data for paymaster contract
	approvePaymasterData, err := s.approveCallData(common.HexToAddress(paymasterAccount), amount)
	if err != nil {
		return nil, fmt.Errorf("failed to create paymaster approve calldata : %w", err)
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
		},
		[][]byte{approvePaymasterData, transferData},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack execute ABI: %w", err)
	}

	return executeBatchCallData, nil
}

// executeBatchCreateOrderCallData creates the calldata for the execute batch method in the smart account.
func (s *OrderService) executeBatchCreateOrderCallData(order *ent.PaymentOrder) ([]byte, error) {
	orderAmountWithFees := order.Amount.Add(order.ProtocolFee).Add(order.SenderFee)

	// Create approve data for gateway contract
	approveGatewayData, err := s.approveCallData(
		OrderConf.GatewayContractAddress,
		utils.ToSubunit(orderAmountWithFees, order.Edges.Token.Decimals),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway approve calldata: %w", err)
	}

	// Fetch paymaster account
	paymasterAccount, err := utils.GetPaymasterAccount(order.Edges.Token.Edges.Network.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get paymaster account: %w", err)
	}

	if ServerConf.Environment != "production" {
		time.Sleep(5 * time.Second)
	}

	// Create approve data for paymaster contract
	approvePaymasterData, err := s.approveCallData(
		common.HexToAddress(paymasterAccount),
		utils.ToSubunit(orderAmountWithFees, order.Edges.Token.Decimals),
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
			OrderConf.GatewayContractAddress,
		},
		[][]byte{approvePaymasterData, approveGatewayData, createOrderData},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack execute ABI: %w", err)
	}

	return executeBatchCreateOrderCallData, nil
}

// approveCallData creates the data for the ERC20 approve method
func (s *OrderService) approveCallData(spender common.Address, amount *big.Int) ([]byte, error) {
	// Create ABI
	erc20ABI, err := abi.JSON(strings.NewReader(contracts.ERC20TokenMetaData.ABI))
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
	erc20ABI, err := abi.JSON(strings.NewReader(contracts.ERC20TokenMetaData.ABI))
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

// // executeCallData creates the data for the execute method in the smart account.
// func (s *OrderService) executeCallData(dest common.Address, value *big.Int, data []byte) ([]byte, error) {
// 	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse smart account ABI: %w", err)
// 	}

// 	executeCallData, err := simpleAccountABI.Pack(
// 		"execute",
// 		dest,
// 		value,
// 		data,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to pack execute ABI: %w", err)
// 	}

// 	return executeCallData, nil
// }

// createOrderCallData creates the data for the createOrder method
func (s *OrderService) createOrderCallData(order *ent.PaymentOrder) ([]byte, error) {
	// Encrypt recipient details
	encryptedOrderRecipient, err := s.encryptOrderRecipient(order.Edges.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt recipient details: %w", err)
	}

	var refundAddress common.Address
	if order.Edges.SenderProfile.RefundAddress == "" {
		refundAddress = common.HexToAddress(order.FromAddress)
	} else {
		refundAddress = common.HexToAddress(order.Edges.SenderProfile.RefundAddress)
	}

	amountWithProtocolFee := order.Amount.Add(order.ProtocolFee)

	// Define params
	params := &CreateOrderParams{
		Token:              common.HexToAddress(order.Edges.Token.ContractAddress),
		Amount:             utils.ToSubunit(amountWithProtocolFee, order.Edges.Token.Decimals),
		InstitutionCode:    utils.StringToByte32(order.Edges.Recipient.Institution),
		Rate:               order.Rate.BigInt(),
		SenderFeeRecipient: common.HexToAddress(order.FeeAddress),
		SenderFee:          order.SenderFee.BigInt(),
		RefundAddress:      refundAddress,
		MessageHash:        encryptedOrderRecipient,
	}

	// Create ABI
	gatewayABI, err := abi.JSON(strings.NewReader(contracts.GatewayMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse GatewayOrder ABI: %w", err)
	}

	// Generate call data
	data, err := gatewayABI.Pack(
		"createOrder",
		params.Token,
		params.Amount,
		params.InstitutionCode,
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
	// Create approve data for gateway contract
	approveGatewayData, err := s.approveCallData(
		OrderConf.GatewayContractAddress,
		utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
	)
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData.approveOrderContract: %w", err)
	}

	// Create refund data
	fee := utils.ToSubunit(order.Edges.Token.Edges.Network.Fee, order.Edges.Token.Decimals)
	refundData, err := s.refundCallData(fee, order.GatewayID)
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData.refundData: %w", err)
	}

	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData.simpleAccountABI: %w", err)
	}

	contractAddresses := []common.Address{
		common.HexToAddress(order.Edges.Token.ContractAddress),
		OrderConf.GatewayContractAddress,
	}

	data := [][]byte{approveGatewayData, refundData}

	if ServerConf.Environment != "production" {
		paymasterAccount, err := utils.GetPaymasterAccount(order.Edges.Token.Edges.Network.ChainID)
		if err != nil {
			return nil, fmt.Errorf("failed to get paymaster account: %w", err)
		}
		time.Sleep(5 * time.Second)

		// Create approve data for paymaster contract
		approvePaymasterData, err := s.approveCallData(
			common.HexToAddress(paymasterAccount),
			utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create paymaster approve calldata : %w", err)
		}

		contractAddresses = append(
			[]common.Address{common.HexToAddress(order.Edges.Token.ContractAddress)},
			contractAddresses...,
		)
		data = append(
			[][]byte{approvePaymasterData},
			data...,
		)
	}

	executeBatchRefundCallData, err := simpleAccountABI.Pack(
		"executeBatch",
		contractAddresses,
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("executeBatchRefundCallData: %w", err)
	}

	return executeBatchRefundCallData, nil
}

// refundCallData creates the data for the refund method
func (s *OrderService) refundCallData(fee *big.Int, orderId string) ([]byte, error) {
	gatewayABI, err := abi.JSON(strings.NewReader(contracts.GatewayMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse GatewayOrder ABI: %w", err)
	}

	decodedOrderID, err := hex.DecodeString(orderId[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode orderId: %w", err)
	}

	// Generate calldata for refund, orderID, and label should be byte32
	data, err := gatewayABI.Pack(
		"refund",
		fee,
		utils.StringToByte32(string(decodedOrderID)),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to pack refund ABI: %w", err)
	}

	return data, nil
}

// executeBatchSettleCallData creates the settle calldata for the execute batch method in the smart account.
func (s *OrderService) executeBatchSettleCallData(ctx context.Context, order *ent.LockPaymentOrder) ([]byte, error) {
	// Create approve data for gateway contract
	approveGatewayData, err := s.approveCallData(
		OrderConf.GatewayContractAddress,
		utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
	)
	if err != nil {
		return nil, fmt.Errorf("approveOrderContract: %w", err)
	}

	contractAddresses := []common.Address{
		common.HexToAddress(order.Edges.Token.ContractAddress),
	}

	data := [][]byte{approveGatewayData}

	if ServerConf.Environment != "production" {
		// Fetch paymaster account
		paymasterAccount, err := utils.GetPaymasterAccount(order.Edges.Token.Edges.Network.ChainID)
		if err != nil {
			return nil, fmt.Errorf("failed to get paymaster account: %w", err)
		}
		time.Sleep(5 * time.Second)

		// Create approve data for paymaster contract
		approvePaymasterData, err := s.approveCallData(
			common.HexToAddress(paymasterAccount),
			utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create paymaster approve calldata : %w", err)
		}

		contractAddresses = append(
			contractAddresses,
			common.HexToAddress(order.Edges.Token.ContractAddress),
		)
		data = append(data, approvePaymasterData)
	}

	// Create settle data
	settleData, err := s.settleCallData(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("settleData: %w", err)
	}

	simpleAccountABI, err := abi.JSON(strings.NewReader(contracts.SimpleAccountMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("simpleAccountABI: %w", err)
	}

	contractAddresses = append(
		contractAddresses,
		OrderConf.GatewayContractAddress,
	)
	data = append(data, settleData)

	executeBatchSettleCallData, err := simpleAccountABI.Pack(
		"executeBatch",
		contractAddresses,
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("executeBatchSettledCallData: %w", err)
	}

	return executeBatchSettleCallData, nil
}

// settleCallData creates the data for the settle method in the gateway contract
func (s *OrderService) settleCallData(ctx context.Context, order *ent.LockPaymentOrder) ([]byte, error) {
	gatewayABI, err := abi.JSON(strings.NewReader(contracts.GatewayMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse GatewayOrder ABI: %w", err)
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

	orderPercent, _ := order.OrderPercent.
		Mul(decimal.NewFromInt(1000)). // convert percent to BPS
		Float64()

	orderID, err := hex.DecodeString(order.GatewayID[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode orderID: %w", err)
	}

	splitOrderID := strings.ReplaceAll(order.ID.String(), "-", "")

	// Generate calldata for settlement
	data, err := gatewayABI.Pack(
		"settle",
		utils.StringToByte32(splitOrderID),
		utils.StringToByte32(string(orderID)),
		common.HexToAddress(providerAddress),
		uint64(orderPercent),
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
		Memo              string
	}{
		recipient.AccountIdentifier, recipient.AccountName, recipient.Institution, recipient.ProviderID, recipient.Memo,
	}

	// Encrypt with the public key of the aggregator
	messageCipher, err := cryptoUtils.PublicKeyEncryptJSON(message, CryptoConf.AggregatorPublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt message: %w", err)
	}

	return base64.StdEncoding.EncodeToString(messageCipher), nil
}
