package order

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/services/contracts"
	"github.com/paycrest/tron-wallet/enums"
	"github.com/paycrest/tron-wallet/grpcClient"
	"github.com/paycrest/tron-wallet/grpcClient/proto/api"
	"github.com/paycrest/tron-wallet/grpcClient/proto/core"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/receiveaddress"
	"github.com/paycrest/protocol/ent/senderordertoken"
	"github.com/paycrest/protocol/ent/senderprofile"
	tokenEnt "github.com/paycrest/protocol/ent/token"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	"github.com/paycrest/tron-wallet/util"
	"github.com/shopspring/decimal"

	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	tronWallet "github.com/paycrest/tron-wallet"
)

// OrderTron provides functionality related to on-chain interactions for payment orders
type OrderTron struct{}

// NewOrderTron creates a new instance of OrderTron.
func NewOrderTron() types.OrderService {
	return &OrderTron{}
}

// getNode returns the node to use based on the environment
func (s *OrderTron) getNode() enums.Node {
	if serverConf.Environment == "production" {
		return enums.MAIN_NODE
	} else {
		return enums.SHASTA_NODE
	}
}

// CreateOrder creates a new payment order on-chain.
func (s *OrderTron) CreateOrder(ctx context.Context, client types.RPCClient, orderID uuid.UUID) error {
	var err error
	orderIDPrefix := strings.Split(orderID.String(), "-")[0]

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
		return fmt.Errorf("%s - Tron.CreateOrder.fetchOrder: %w", orderIDPrefix, err)
	}

	// Create wallet
	saltDecrypted, err := cryptoUtils.DecryptPlain(order.Edges.ReceiveAddress.Salt)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.DecryptPlain: %w", orderIDPrefix, err)
	}

	wallet, err := tronWallet.CreateTronWallet(s.getNode(), string(saltDecrypted))
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.CreateTronWallet: %w", orderIDPrefix, err)
	}

	// Transfer TRX from master wallet to receive address for gas
	masterWallet, err := cryptoUtils.GenerateTronAccountFromIndex(0)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.GenerateTronAccountFromIndex: %w", orderIDPrefix, err)
	}

	balance, err := wallet.Balance()
	if err != nil {
		balance = 0
	}

	if balance < 150000000 {
		_, err = masterWallet.Transfer(wallet.AddressBase58, 150000000)
		if err != nil {
			return fmt.Errorf("%s - Tron.CreateOrder.Transfer: %w", orderIDPrefix, err)
		}
		time.Sleep(5 * time.Second) // wait for wallet to be pre-funded with gas
	}

	// Normalize addresses
	gatewayContractAddress, err := util.Base58ToAddress(order.Edges.Token.Edges.Network.GatewayContractAddress)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	senderAddress, err := util.Base58ToAddress(wallet.AddressBase58)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	tokenContractAddress, err := util.Base58ToAddress(order.Edges.Token.ContractAddress)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	// Approve gateway contract to spend token
	calldata, err := s.approveCallData(gatewayContractAddress, utils.ToSubunit(order.Amount.Add(order.ProtocolFee), order.Edges.Token.Decimals))
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.approveCallData: %w", orderIDPrefix, err)
	}

	ct := &core.TriggerSmartContract{
		OwnerAddress:    senderAddress.Bytes(),
		ContractAddress: tokenContractAddress.Bytes(),
		Data:            calldata,
	}
	_, err = s.sendTransaction(wallet, ct, 30000000)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.sendTransaction: %w", orderIDPrefix, err)
	}

	// Create order in gateway contract
	calldata, err = s.createOrderCallData(order)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.createOrderCallData: %w", orderIDPrefix, err)
	}

	ct = &core.TriggerSmartContract{
		OwnerAddress:    senderAddress.Bytes(),
		ContractAddress: gatewayContractAddress.Bytes(),
		Data:            calldata,
	}
	txHash, err := s.sendTransaction(wallet, ct, 150000000)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.sendTransaction: %w", orderIDPrefix, err)
	}

	// Transfer network fee from receive address to master wallet
	_, err = wallet.TransferTRC20(
		&tronWallet.Token{
			ContractAddress: enums.ContractAddress(order.Edges.Token.ContractAddress),
		},
		masterWallet.AddressBase58,
		utils.ToSubunit(order.NetworkFee, order.Edges.Token.Decimals).Int64(),
		30000000,
	)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.TransferTRC20: %w", orderIDPrefix, err)
	}

	// Update payment order
	_, err = order.Update().
		SetBlockNumber(order.Edges.ReceiveAddress.LastIndexedBlock).
		SetTxHash(txHash).
		SetStatus(paymentorder.StatusPending).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.updateTxHash: %w", orderIDPrefix, err)
	}

	paymentOrder, err := db.Client.PaymentOrder.
		Query().
		Where(paymentorder.IDEQ(orderID)).
		WithSenderProfile().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.refetchOrder: %w", orderIDPrefix, err)
	}

	// Send webhook notifcation to sender
	err = utils.SendPaymentOrderWebhook(ctx, paymentOrder)
	if err != nil {
		return fmt.Errorf("%s - Tron.CreateOrder.webhook: %w", orderIDPrefix, err)
	}

	return nil
}

// RefundOrder refunds sender on canceled lock order
func (s *OrderTron) RefundOrder(ctx context.Context, client types.RPCClient, orderID string) error {
	orderIDPrefix := strings.Split(orderID, "-")[0]

	// Fetch lock order from db
	lockOrder, err := db.Client.LockPaymentOrder.
		Query().
		Where(lockpaymentorder.GatewayIDEQ(orderID)).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		First(ctx)
	if err != nil {
		return fmt.Errorf("Tron.RefundOrder.fetchLockOrder: %w", err)
	}

	// Generate master wallet
	wallet, err := cryptoUtils.GenerateTronAccountFromIndex(0)
	if err != nil {
		return fmt.Errorf("%s - Tron.RefundOrder.GenerateTronAccountFromIndex: %w", orderIDPrefix, err)
	}

	// Normalize addresses
	gatewayContractAddress, err := util.Base58ToAddress(lockOrder.Edges.Token.Edges.Network.GatewayContractAddress)
	if err != nil {
		return fmt.Errorf("%s - Tron.RefundOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	senderAddress, err := util.Base58ToAddress(wallet.AddressBase58)
	if err != nil {
		return fmt.Errorf("%s - Tron.RefundOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	tokenContractAddress, err := util.Base58ToAddress(lockOrder.Edges.Token.ContractAddress)
	if err != nil {
		return fmt.Errorf("%s - Tron.RefundOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	// Approve gateway contract to spend token
	sourceOrder, err := db.Client.PaymentOrder.
		Query().
		Where(paymentorder.GatewayIDEQ(lockOrder.GatewayID)).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		Only(context.Background())
	if err != nil {
		return fmt.Errorf("Tron.RefundOrder.fetchPaymentOrder: %w", err)
	}

	calldata, err := s.approveCallData(
		gatewayContractAddress,
		utils.ToSubunit(
			lockOrder.Amount.Add(sourceOrder.SenderFee).Add(sourceOrder.ProtocolFee),
			lockOrder.Edges.Token.Decimals,
		),
	)
	if err != nil {
		return fmt.Errorf("%s - Tron.RefundOrder.approveCallData: %w", orderIDPrefix, err)
	}

	ct := &core.TriggerSmartContract{
		OwnerAddress:    senderAddress.Bytes(),
		ContractAddress: tokenContractAddress.Bytes(),
		Data:            calldata,
	}
	_, err = s.sendTransaction(wallet, ct, 50000000)
	if err != nil {
		return fmt.Errorf("%s - Tron.RefundOrder.sendTransaction: %w", orderIDPrefix, err)
	}

	// Refund order in gateway contract
	fee := utils.ToSubunit(decimal.NewFromInt(0), lockOrder.Edges.Token.Decimals)
	calldata, err = s.refundCallData(fee, lockOrder.GatewayID)
	if err != nil {
		return fmt.Errorf("%s - Tron.RefundOrder.refundCallData: %w", orderIDPrefix, err)
	}

	ct = &core.TriggerSmartContract{
		OwnerAddress:    senderAddress.Bytes(),
		ContractAddress: gatewayContractAddress.Bytes(),
		Data:            calldata,
	}
	_, err = s.sendTransaction(wallet, ct, 50000000)
	if err != nil {
		return fmt.Errorf("%s - Tron.RefundOrder.sendTransaction: %w", orderIDPrefix, err)
	}

	return nil
}

// RevertOrder reverts an initiated payment order on-chain.
func (s *OrderTron) RevertOrder(ctx context.Context, client types.RPCClient, order *ent.PaymentOrder) error {
	if !order.AmountReturned.Equal(decimal.Zero) || strings.HasPrefix(order.Edges.Recipient.Memo, "P#P") {
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
		return fmt.Errorf("%s - Tron.RevertOrder.fetchOrder: %w", orderIDPrefix, err)
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

	// Convert amountMinusFee to big.Int
	amountMinusFeeBigInt := utils.ToSubunit(amountToRevert, order.Edges.Token.Decimals)

	// Create wallet
	saltDecrypted, err := cryptoUtils.DecryptPlain(order.Edges.ReceiveAddress.Salt)
	if err != nil {
		return fmt.Errorf("%s - Tron.RevertOrder.DecryptPlain: %w", orderIDPrefix, err)
	}

	wallet, err := tronWallet.CreateTronWallet(s.getNode(), string(saltDecrypted))
	if err != nil {
		return fmt.Errorf("%s - Tron.RevertOrder.CreateTronWallet: %w", orderIDPrefix, err)
	}

	// Transfer TRX from master wallet to receive address for gas
	masterWallet, err := cryptoUtils.GenerateTronAccountFromIndex(0)
	if err != nil {
		return fmt.Errorf("%s - Tron.RevertOrder.GenerateTronAccountFromIndex: %w", orderIDPrefix, err)
	}

	balance, err := wallet.Balance()
	if err != nil {
		balance = 0
	}

	if balance < 30000000 {
		_, err = masterWallet.Transfer(wallet.AddressBase58, 30000000)
		if err != nil {
			return fmt.Errorf("%s - Tron.RevertOrder.Transfer: %w", orderIDPrefix, err)
		}
		time.Sleep(5 * time.Second) // wait for wallet to be pre-funded with gas
	}

	// Transfer TRC20 token
	token := &tronWallet.Token{
		ContractAddress: enums.ContractAddress(order.Edges.Token.ContractAddress),
	}

	txHash, err := wallet.TransferTRC20(token, order.ReturnAddress, amountMinusFeeBigInt.Int64(), 50000000)
	if err != nil {
		return fmt.Errorf("%s - Tron.RevertOrder.TransferTRC20: %w", orderIDPrefix, err)
	}

	// Update payment order
	_, err = order.Update().
		SetTxHash(txHash).
		SetAmountReturned(amountToRevert).
		SetStatus(paymentorder.StatusReverted).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("%s - Tron.RevertOrder.updateTxHash: %w", orderIDPrefix, err)
	}

	// Send webhook notifcation to sender
	order.Status = paymentorder.StatusReverted
	err = utils.SendPaymentOrderWebhook(ctx, order)
	if err != nil {
		return fmt.Errorf("Tron.RevertOrder.webhook: %v", err)
	}

	return nil
}

// SettleOrder settles a payment order on-chain.
func (s *OrderTron) SettleOrder(ctx context.Context, client types.RPCClient, orderID uuid.UUID) error {
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
		return fmt.Errorf("%s - Tron.SettleOrder.fetchOrder: %w", orderIDPrefix, err)
	}

	// Generate master wallet
	wallet, err := cryptoUtils.GenerateTronAccountFromIndex(0)
	if err != nil {
		return fmt.Errorf("%s - Tron.SettleOrder.GenerateTronAccountFromIndex: %w", orderIDPrefix, err)
	}

	// Normalize addresses
	gatewayContractAddress, err := util.Base58ToAddress(order.Edges.Token.Edges.Network.GatewayContractAddress)
	if err != nil {
		return fmt.Errorf("%s - Tron.SettleOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	senderAddress, err := util.Base58ToAddress(wallet.AddressBase58)
	if err != nil {
		return fmt.Errorf("%s - Tron.SettleOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	tokenContractAddress, err := util.Base58ToAddress(order.Edges.Token.ContractAddress)
	if err != nil {
		return fmt.Errorf("%s - Tron.SettleOrder.Base58ToAddress: %w", orderIDPrefix, err)
	}

	// Approve gateway contract to spend token
	calldata, err := s.approveCallData(
		gatewayContractAddress,
		utils.ToSubunit(order.Amount, order.Edges.Token.Decimals),
	)
	if err != nil {
		return fmt.Errorf("%s - Tron.SettleOrder.approveCallData: %w", orderIDPrefix, err)
	}

	ct := &core.TriggerSmartContract{
		OwnerAddress:    senderAddress.Bytes(),
		ContractAddress: tokenContractAddress.Bytes(),
		Data:            calldata,
	}
	_, err = s.sendTransaction(wallet, ct, 50000000)
	if err != nil {
		return fmt.Errorf("%s - Tron.SettleOrder.sendTransaction: %w", orderIDPrefix, err)
	}

	// Settle order in gateway contract
	calldata, err = s.settleCallData(ctx, order)
	if err != nil {
		return fmt.Errorf("%s - Tron.SettleOrder.settleCallData: %w", orderIDPrefix, err)
	}

	ct = &core.TriggerSmartContract{
		OwnerAddress:    senderAddress.Bytes(),
		ContractAddress: gatewayContractAddress.Bytes(),
		Data:            calldata,
	}
	_, err = s.sendTransaction(wallet, ct, 50000000)
	if err != nil {
		return fmt.Errorf("%s - Tron.SettleOrder.sendTransaction: %w", orderIDPrefix, err)
	}

	return nil
}

// approveCallData creates the data for the ERC20 approve method
func (s *OrderTron) approveCallData(spender util.Address, amount *big.Int) ([]byte, error) {
	// Create ABI
	erc20ABI, err := abi.JSON(strings.NewReader(contracts.ERC20TokenMetaData.ABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse erc20 ABI: %w", err)
	}

	// Create calldata
	calldata, err := erc20ABI.Pack("approve", common.HexToAddress(spender.Hex()[4:]), amount)
	if err != nil {
		return nil, fmt.Errorf("failed to pack approve ABI: %w", err)
	}

	return calldata, nil
}

// createOrderCallData creates the data for the createOrder method
func (s *OrderTron) createOrderCallData(order *ent.PaymentOrder) ([]byte, error) {
	// Encrypt recipient details
	encryptedOrderRecipient, err := s.encryptOrderRecipient(order.Edges.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt recipient details: %w", err)
	}

	// Fetch token configuration
	isTokenConfigured := true
	var networkIdentifier string

	if serverConf.Environment == "production" {
		networkIdentifier = "tron"
	} else {
		networkIdentifier = "tron-shasta"
	}

	token, err := db.Client.SenderOrderToken.
		Query().
		Where(
			senderordertoken.And(
				senderordertoken.HasTokenWith(
					tokenEnt.HasNetworkWith(network.IdentifierEQ(networkIdentifier)),
				),
				senderordertoken.HasSenderWith(
					senderprofile.IDEQ(order.Edges.SenderProfile.ID),
				),
			)).
		Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			isTokenConfigured = false
		} else {
			return nil, fmt.Errorf("failed to fetch order token: %w", err)
		}
	}

	var refundAddress string
	var refundAddressTron util.Address

	if isTokenConfigured {
		refundAddressTron, _ = util.Base58ToAddress(token.RefundAddress)
		refundAddress = refundAddressTron.Hex()[4:]
	} else {
		refundAddressTron, _ = util.Base58ToAddress(order.ReturnAddress)
		refundAddress = refundAddressTron.Hex()[4:]
	}

	amountWithProtocolFee := order.Amount.Add(order.ProtocolFee)
	tokenContractAddressTron, _ := util.Base58ToAddress(order.Edges.Token.ContractAddress)

	senderFeeRecipientTron, _ := util.Base58ToAddress(order.FeeAddress)
	senderFeeRecipient := senderFeeRecipientTron.Hex()[4:]

	// Define params
	params := &types.CreateOrderParams{
		Token:              common.HexToAddress(tokenContractAddressTron.Hex()[4:]),
		Amount:             utils.ToSubunit(amountWithProtocolFee, order.Edges.Token.Decimals),
		Rate:               order.Rate.BigInt(),
		SenderFeeRecipient: common.HexToAddress(senderFeeRecipient),
		SenderFee:          order.SenderFee.BigInt(),
		RefundAddress:      common.HexToAddress(refundAddress),
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

// refundCallData creates the data for the refund method
func (s *OrderTron) refundCallData(fee *big.Int, orderId string) ([]byte, error) {
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

// settleCallData creates the data for the settle method in the gateway contract
func (s *OrderTron) settleCallData(ctx context.Context, order *ent.LockPaymentOrder) ([]byte, error) {
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
			providerAddressTron, _ := util.Base58ToAddress(addr.Address)
			providerAddress = providerAddressTron.Hex()[4:]
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
func (s *OrderTron) encryptOrderRecipient(recipient *ent.PaymentOrderRecipient) (string, error) {
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
	messageCipher, err := cryptoUtils.PublicKeyEncryptJSON(message, cryptoConf.AggregatorPublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt message: %w", err)
	}

	return base64.StdEncoding.EncodeToString(messageCipher), nil
}

// signTransaction signs a transaction with a private key
func (s *OrderTron) signTransaction(transaction *api.TransactionExtention, privateKey *ecdsa.PrivateKey) (*api.TransactionExtention, error) {
	rawData, err := proto.Marshal(transaction.Transaction.GetRawData())
	if err != nil {
		return transaction, fmt.Errorf("proto marshal tx raw data error: %v", err)
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return transaction, fmt.Errorf("sign error: %v", err)
	}

	transaction.Transaction.Signature = append(transaction.Transaction.Signature, signature)
	return transaction, nil
}

// broadcastTransaction broadcasts a transaction to the Tron network
func (s *OrderTron) broadcastTransaction(node enums.Node, transaction *api.TransactionExtention) error {
	c, err := grpcClient.GetGrpcClient(node)
	if err != nil {
		return err
	}

	res, err := c.Broadcast(transaction.Transaction)
	if err != nil {
		return err
	}

	if !res.Result {
		return errors.New(res.Code.String())
	}

	return nil
}

// callMethod reads data from the Tron network
func (s *OrderTron) callMethod(ct *core.TriggerSmartContract) (*api.TransactionExtention, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if len(config.OrderConfig().TronProApiKey) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "TRON-PRO-API-KEY", config.OrderConfig().TronProApiKey)
	}
	defer cancel()

	g, err := grpcClient.GetGrpcClient(s.getNode())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GRPC client: %w", err)
	}

	tx, err := g.Client.TriggerConstantContract(ctx, ct)
	if err != nil {
		return nil, err
	}

	if tx.Result.Code > 0 {
		return tx, fmt.Errorf(string(tx.Result.Message))
	}

	return tx, nil
}

// sendTransaction sends a transaction to the Tron network
func (s *OrderTron) sendTransaction(wallet *tronWallet.TronWallet, ct *core.TriggerSmartContract, feeLimit int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if len(config.OrderConfig().TronProApiKey) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "TRON-PRO-API-KEY", config.OrderConfig().TronProApiKey)
	}
	defer cancel()

	g, err := grpcClient.GetGrpcClient(s.getNode())
	if err != nil {
		return "", fmt.Errorf("failed to connect to GRPC client: %w", err)
	}

	// Trigger smart contract
	tx, err := g.Client.TriggerContract(ctx, ct)
	if err != nil {
		return "", err
	}

	if tx.Result.Code > 0 {
		return "", fmt.Errorf("%s", string(tx.Result.Message))
	}
	if feeLimit > 0 {
		tx.Transaction.RawData.FeeLimit = feeLimit
		g.UpdateHash(tx)
	}

	if tx.Result.Code > 0 {
		return "", fmt.Errorf(string(tx.Result.Message))
	}

	// Sign and broadcast transaction
	privateKey, err := wallet.PrivateKeyRCDSA()
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}
	signedTx, err := s.signTransaction(tx, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = s.broadcastTransaction(wallet.Node, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	return hexutil.Encode(tx.GetTxid())[2:], err
}
