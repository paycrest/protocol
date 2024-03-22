package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	networkent "github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/provisionbucket"
	"github.com/paycrest/protocol/ent/receiveaddress"
	"github.com/paycrest/protocol/ent/token"
	"github.com/paycrest/protocol/services/contracts"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/shopspring/decimal"
)

var OrderConf = config.OrderConfig()

// Indexer is an interface for indexing blockchain data to the database.
type Indexer interface {
	IndexERC20Transfer(ctx context.Context, client types.RPCClient, receiveAddress *ent.ReceiveAddress) error
	IndexOrderCreated(ctx context.Context, client types.RPCClient, network *ent.Network) error
	IndexOrderSettled(ctx context.Context, client types.RPCClient, network *ent.Network) error
	IndexOrderRefunded(ctx context.Context, client types.RPCClient, network *ent.Network) error
	HandleReceiveAddressValidity(ctx context.Context, receiveAddress *ent.ReceiveAddress, paymentOrder *ent.PaymentOrder) error
}

// IndexerService performs blockchain to database extract, transform, load (ETL) operations.
type IndexerService struct {
	priorityQueue *PriorityQueueService
	order         Order
}

// NewIndexerService creates a new instance of IndexerService.
func NewIndexerService(order Order) Indexer {
	priorityQueue := NewPriorityQueueService()

	return &IndexerService{
		priorityQueue: priorityQueue,
		order:         order,
	}
}

// IndexERC20Transfer indexes deposits to the order contract for a specific network.
func (s *IndexerService) IndexERC20Transfer(ctx context.Context, client types.RPCClient, receiveAddress *ent.ReceiveAddress) error {
	var err error

	// Fetch payment order from db
	paymentOrder, err := db.Client.PaymentOrder.
		Query().
		Where(
			paymentorder.HasReceiveAddressWith(
				receiveaddress.AddressEQ(receiveAddress.Address),
			),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		Only(ctx)
	if err != nil {
		logger.Errorf("IndexERC20Transfer.db: %v", err)
		return err
	}

	token := paymentOrder.Edges.Token

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(token.Edges.Network.RPCEndpoint)
		if err != nil {
			logger.Errorf("IndexERC20Transfer.NewEthClient: %v", err)
			return err
		}
	}

	// Initialize contract filterer
	filterer, err := contracts.NewERC20TokenFilterer(common.HexToAddress(token.ContractAddress), client)
	if err != nil {
		logger.Errorf("IndexERC20Transfer.NewERC20TokenFilterer: %v", err)
		return err
	}

	// Index missed blocks
	go func() {
		// Filter logs from the last lock payment order block number
		opts := s.getMissedERC20BlocksOpts(ctx, client, token.Edges.Network)

		// Fetch logs
		var iter *contracts.ERC20TokenTransferIterator
		retryErr := utils.Retry(3, 5*time.Second, func() error {
			var err error
			iter, err = filterer.FilterTransfer(opts, nil, nil)
			return err
		})
		if retryErr != nil {
			logger.Errorf("IndexERC20Transfer.FilterTransfer: %v", retryErr)
			return
		}

		// Iterate over logs
		for iter.Next() {
			ok, err := s.updateReceiveAddressStatus(ctx, receiveAddress, paymentOrder, iter.Event)
			if err != nil {
				logger.Errorf("IndexERC20Transfer.updateReceiveAddressStatus: %v", err)
				continue
			}
			if ok {
				return
			}
		}
	}()

	if ServerConf.Environment != "test" {
		// Start listening for ERC20 transfer events
		logs := make(chan *contracts.ERC20TokenTransfer)

		sub, err := filterer.WatchTransfer(&bind.WatchOpts{
			Start: nil,
		}, logs, nil, nil)
		if err != nil {
			logger.Errorf("IndexERC20Transfer.WatchTransfer: %v", err)
			return err
		}

		defer sub.Unsubscribe()

		logger.Infof(fmt.Sprintf("Listening for ERC20 Transfer event: %s\n", receiveAddress.Address))

		for {
			select {
			case log := <-logs:
				ok, err := s.updateReceiveAddressStatus(ctx, receiveAddress, paymentOrder, log)
				if err != nil {
					logger.Errorf("IndexERC20Transfer.updateReceiveAddressStatus: %v", err)
					continue
				}
				if ok {
					close(logs)
					return nil
				}
			case err := <-sub.Err():
				if err == nil {
					sub.Unsubscribe()

					// Retry the subscription
					retryErr := utils.Retry(3, 5*time.Second, func() error {
						sub, err = filterer.WatchTransfer(&bind.WatchOpts{
							Start: nil,
						}, logs, nil, nil)
						return err
					})
					if retryErr != nil {
						logger.Errorf("IndexERC20Transfer.WatchTransfer: %v", retryErr)
						return retryErr
					}
				} else {
					logger.Errorf("IndexERC20Transfer.logError: %v", err)
					continue
				}
			}
		}
	}

	return nil
}

// IndexOrderCreated indexes deposits to the order contract for a specific network.
func (s *IndexerService) IndexOrderCreated(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("IndexOrderCreated.NewEthClient: %w", err)
		}
	}

	// Initialize contract filterer
	filterer, err := contracts.NewPaycrestFilterer(OrderConf.PaycrestOrderContractAddress, client)
	if err != nil {
		return fmt.Errorf("IndexOrderCreated.NewPaycrestFilterer: %w", err)
	}

	// Index missed blocks
	go func() {
		// Filter logs from the last lock payment order block number
		opts := s.getMissedOrderBlocksOpts(ctx, client, network, lockpaymentorder.StatusPending)

		// Fetch logs
		var iter *contracts.PaycrestOrderCreatedIterator
		retryErr := utils.Retry(3, 5*time.Second, func() error {
			var err error
			iter, err = filterer.FilterOrderCreated(opts, nil, nil, nil)
			return err
		})
		if retryErr != nil {
			logger.Errorf("IndexOrderCreated.FilterOrderCreated: %v", retryErr)
			return
		}

		// Iterate over logs
		for iter.Next() {
			err := s.createLockPaymentOrder(ctx, client, network, iter.Event)
			if err != nil {
				logger.Errorf("IndexOrderCreated.createOrder: %v", err)
				continue
			}
		}
	}()

	if ServerConf.Environment != "test" {
		// Start listening for deposit events
		logs := make(chan *contracts.PaycrestOrderCreated)

		sub, err := filterer.WatchOrderCreated(&bind.WatchOpts{
			Start: nil,
		}, logs, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("IndexOrderCreated.WatchOrderCreated: %w", err)
		}

		defer sub.Unsubscribe()

		logger.Infof("Listening for OrderCreated events...\n")

		for {
			select {
			case log := <-logs:
				err := s.createLockPaymentOrder(ctx, client, network, log)
				if err != nil {
					logger.Errorf("IndexOrderCreated.createLockPaymentOrder: %v", err)
					continue
				}
			case err := <-sub.Err():
				if err == nil {
					sub.Unsubscribe()

					// Retry the subscription
					retryErr := utils.Retry(3, 5*time.Second, func() error {
						sub, err = filterer.WatchOrderCreated(&bind.WatchOpts{
							Start: nil,
						}, logs, nil, nil, nil)
						return err
					})
					if retryErr != nil {
						logger.Errorf("IndexOrderCreated.WatchOrderCreated: %v", retryErr)
						return retryErr
					}
				} else {
					logger.Errorf("IndexOrderCreated.logError: %v", err)
					continue
				}
			}
		}
	}

	return nil
}

// IndexOrderSettled indexes order settlements for a specific network.
func (s *IndexerService) IndexOrderSettled(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("IndexOrderSettled.NewEthClient: %w", err)
		}
	}

	// Initialize contract filterer
	filterer, err := contracts.NewPaycrestFilterer(OrderConf.PaycrestOrderContractAddress, client)
	if err != nil {
		return fmt.Errorf("IndexOrderSettled.NewPaycrestFilterer: %w", err)
	}

	// Index missed blocks
	go func() {
		// Filter logs from the last lock payment order block number
		opts := s.getMissedOrderBlocksOpts(ctx, client, network, lockpaymentorder.StatusSettling)

		// Fetch logs
		var iter *contracts.PaycrestOrderSettledIterator
		retryErr := utils.Retry(3, 5*time.Second, func() error {
			var err error
			iter, err = filterer.FilterOrderSettled(opts, nil, nil)
			return err
		})
		if retryErr != nil {
			logger.Errorf("IndexOrderSettled.FilterOrderSettled: %v", retryErr)
			return
		}

		// Iterate over logs
		for iter.Next() {
			log := iter.Event
			err := s.updateOrderStatusSettled(ctx, log)
			if err != nil {
				logger.Errorf("IndexOrderSettled.update: %v", err)
				continue
			}
		}
	}()

	if ServerConf.Environment != "test" {
		// Start listening for settlement events
		logs := make(chan *contracts.PaycrestOrderSettled)

		sub, err := filterer.WatchOrderSettled(&bind.WatchOpts{
			Start: nil,
		}, logs, nil, nil)
		if err != nil {
			return fmt.Errorf("IndexOrderSettled.WatchOrderSettled: %w", err)
		}

		defer sub.Unsubscribe()

		logger.Infof("Listening for OrderSettled events...\n")

		for {
			select {
			case log := <-logs:
				err := s.updateOrderStatusSettled(ctx, log)
				if err != nil {
					logger.Errorf("IndexOrderSettled.update: %v", err)
					continue
				}
			case err := <-sub.Err():
				if err == nil {
					sub.Unsubscribe()

					// Retry the subscription
					retryErr := utils.Retry(3, 5*time.Second, func() error {
						sub, err = filterer.WatchOrderSettled(&bind.WatchOpts{
							Start: nil,
						}, logs, nil, nil)
						return err
					})
					if retryErr != nil {
						logger.Errorf("IndexOrderSettled.WatchOrderSettled: %v", retryErr)
						return retryErr
					}
				} else {
					logger.Errorf("IndexOrderSettled.logError: %v", err)
					continue
				}
			}
		}
	}
	return nil
}

// IndexOrderRefunded indexes order refunds for a specific network.
func (s *IndexerService) IndexOrderRefunded(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("IndexOrderRefunded.NewEthClient: %w", err)
		}
	}

	// Initialize contract filterer
	filterer, err := contracts.NewPaycrestFilterer(OrderConf.PaycrestOrderContractAddress, client)
	if err != nil {
		return fmt.Errorf("IndexOrderRefunded.NewPaycrestFilterer: %w", err)
	}

	// Index missed blocks
	go func() {
		// Filter logs from the last lock payment order block number
		opts := s.getMissedOrderBlocksOpts(ctx, client, network, lockpaymentorder.StatusRefunding)

		// Fetch logs
		var iter *contracts.PaycrestOrderRefundedIterator
		retryErr := utils.Retry(3, 5*time.Second, func() error {
			var err error
			iter, err = filterer.FilterOrderRefunded(opts, nil)
			return err
		})
		if retryErr != nil {
			logger.Errorf("IndexOrderRefunded.FilterOrderRefunded: %v", retryErr)
			return
		}

		// Iterate over logs
		for iter.Next() {
			log := iter.Event
			err := s.updateOrderStatusRefunded(ctx, log)
			if err != nil {
				logger.Errorf("IndexOrderRefunded.update: %v", err)
				continue
			}
		}
	}()

	if ServerConf.Environment != "test" {
		// Start listening for refund events
		logs := make(chan *contracts.PaycrestOrderRefunded)

		sub, err := filterer.WatchOrderRefunded(&bind.WatchOpts{
			Start: nil,
		}, logs, nil)
		if err != nil {
			return fmt.Errorf("IndexOrderRefunded.WatchOrderRefunded: %w", err)
		}

		defer sub.Unsubscribe()

		logger.Infof("Listening for OrderRefunded events...\n")

		for {
			select {
			case log := <-logs:
				err := s.updateOrderStatusRefunded(ctx, log)
				if err != nil {
					logger.Errorf("IndexOrderRefunded.update: %v", err)
					continue
				}
			case err := <-sub.Err():
				if err == nil {
					sub.Unsubscribe()

					// Retry the subscription
					retryErr := utils.Retry(3, 5*time.Second, func() error {
						sub, err = filterer.WatchOrderRefunded(&bind.WatchOpts{
							Start: nil,
						}, logs, nil)
						return err
					})
					if retryErr != nil {
						logger.Errorf("IndexOrderRefunded.WatchOrderRefunded: %v", retryErr)
						return retryErr
					}
				} else {
					logger.Errorf("IndexOrderRefunded.logError: %v", err)
					continue
				}
			}
		}
	}

	return nil
}

// HandleReceiveAddressValidity checks the validity of a receive address
func (s *IndexerService) HandleReceiveAddressValidity(ctx context.Context, receiveAddress *ent.ReceiveAddress, paymentOrder *ent.PaymentOrder) error {
	if receiveAddress.Status != receiveaddress.StatusUsed {
		amountNotPaidInFull := receiveAddress.Status == receiveaddress.StatusPartial || receiveAddress.Status == receiveaddress.StatusUnused
		validUntilIsFarGone := receiveAddress.ValidUntil.Before(time.Now().Add(-(5 * time.Minute)))
		isExpired := receiveAddress.ValidUntil.Before(time.Now())

		if validUntilIsFarGone {
			_, err := receiveAddress.
				Update().
				SetValidUntil(time.Now().Add(OrderConf.ReceiveAddressValidity)).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("HandleReceiveAddressValidity.db: %v", err)
			}
		} else if isExpired && amountNotPaidInFull {
			// Receive address hasn't received full payment after validity period, mark status as expired
			_, err := receiveAddress.
				Update().
				SetStatus(receiveaddress.StatusExpired).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("HandleReceiveAddressValidity.db: %v", err)
			}

			// Expire payment order
			_, err = paymentOrder.
				Update().
				SetStatus(paymentorder.StatusExpired).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("HandleReceiveAddressValidity.db: %v", err)
			}

			// Revert amount to the from address
			err = s.order.RevertOrder(ctx, paymentOrder)
			if err != nil {
				return fmt.Errorf("HandleReceiveAddressValidity.RevertOrder: %v", err)
			}
		}
	} else {
		// Revert excess amount to the from address
		err := s.order.RevertOrder(ctx, paymentOrder)
		if err != nil {
			return fmt.Errorf("HandleReceiveAddressValidity.RevertOrder: %v", err)
		}
	}

	return nil
}

// getOrderRecipientFromMessageHash decrypts the message hash and returns the order recipient
func (s *IndexerService) getOrderRecipientFromMessageHash(messageHash string) (*types.PaymentOrderRecipient, error) {
	messageCipher, err := hex.DecodeString(strings.TrimPrefix(messageHash, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode message hash: %w", err)
	}

	// Decrypt with the private key of the aggregator
	message, err := cryptoUtils.PublicKeyDecryptJSON(messageCipher, CryptoConf.AggregatorPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message hash: %w", err)
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	var recipient *types.PaymentOrderRecipient
	if err := json.Unmarshal(messageBytes, &recipient); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return recipient, nil
}

// updateOrderStatusRefunded updates the status of a payment order to refunded
func (s *IndexerService) updateOrderStatusRefunded(ctx context.Context, log *contracts.PaycrestOrderRefunded) error {
	// Aggregator side status update
	_, err := db.Client.LockPaymentOrder.
		Update().
		Where(
			lockpaymentorder.OrderIDEQ(fmt.Sprintf("0x%v", hex.EncodeToString(log.OrderId[:]))),
		).
		SetBlockNumber(int64(log.Raw.BlockNumber)).
		SetTxHash(log.Raw.TxHash.Hex()).
		SetStatus(lockpaymentorder.StatusRefunded).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusRefunded.aggregator: %v", err)
	}

	// Sender side status update
	_, err = db.Client.PaymentOrder.
		Update().
		Where(
			paymentorder.LabelEQ(utils.Byte32ToString(log.Label)),
		).
		SetTxHash(log.Raw.TxHash.Hex()).
		SetStatus(paymentorder.StatusRefunded).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusRefunded.sender: %v", err)
	}

	// Fetch payment order
	paymentOrder, err := db.Client.PaymentOrder.
		Query().
		Where(
			paymentorder.LabelEQ(utils.Byte32ToString(log.Label)),
		).
		WithSenderProfile().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusRefunded.fetchOrder: %v", err)
	}

	// Send webhook notifcation to sender
	err = utils.SendPaymentOrderWebhook(ctx, paymentOrder)
	if err != nil {
		return fmt.Errorf("updateOrderStatusRefunded.webhook: %v", err)
	}

	return nil
}

// updateOrderStatusSettled updates the status of a payment order to settled
func (s *IndexerService) updateOrderStatusSettled(ctx context.Context, event *contracts.PaycrestOrderSettled) error {
	// Check for existing lock order with txHash
	orderCount, err := db.Client.LockPaymentOrder.
		Query().
		Where(lockpaymentorder.TxHashEQ(event.Raw.TxHash.Hex())).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusSettled.db: %v", err)
	}

	if orderCount > 0 {
		// This log has already been indexed
		return nil
	}

	// Aggregator side status update
	splitOrderId, _ := uuid.Parse(utils.Byte32ToString(event.SplitOrderId))
	_, err = db.Client.LockPaymentOrder.
		Update().
		Where(
			lockpaymentorder.IDEQ(splitOrderId),
		).
		SetBlockNumber(int64(event.Raw.BlockNumber)).
		SetTxHash(event.Raw.TxHash.Hex()).
		SetStatus(lockpaymentorder.StatusSettled).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusSettled.aggregator: %v", err)
	}

	// Sender side status update
	paymentOrderUpdate := db.Client.PaymentOrder.
		Update().
		Where(
			paymentorder.LabelEQ(utils.Byte32ToString(event.Label)),
		)

	// Convert settled percent to BPS
	settledPercent := decimal.NewFromBigInt(event.SettlePercent, 0).Div(decimal.NewFromInt(1000))

	// If settled percent is 100%, mark order as settled
	if settledPercent.Equal(decimal.NewFromInt(100)) {
		paymentOrderUpdate = paymentOrderUpdate.SetStatus(paymentorder.StatusSettled)
	}

	_, err = paymentOrderUpdate.
		SetTxHash(event.Raw.TxHash.Hex()).
		SetPercentSettled(settledPercent).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusSettled.sender: %v", err)
	}

	// Fetch payment order
	paymentOrder, err := db.Client.PaymentOrder.
		Query().
		Where(
			paymentorder.LabelEQ(utils.Byte32ToString(event.Label)),
		).
		WithSenderProfile().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusSettled.fetchOrder: %v", err)
	}

	// Send webhook notifcation to sender
	err = utils.SendPaymentOrderWebhook(ctx, paymentOrder)
	if err != nil {
		return fmt.Errorf("updateOrderStatusSettled.webhook: %v", err)
	}

	return nil
}

// createLockPaymentOrder saves a lock payment order in the database
func (s *IndexerService) createLockPaymentOrder(ctx context.Context, client types.RPCClient, network *ent.Network, deposit *contracts.PaycrestOrderCreated) error {
	// Check for existing address with txHash
	orderCount, err := db.Client.LockPaymentOrder.
		Query().
		Where(lockpaymentorder.TxHashEQ(deposit.Raw.TxHash.Hex())).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("createLockPaymentOrder.db: %v", err)
	}

	if orderCount > 0 {
		// This transfer has already been indexed
		return nil
	}

	// Get token from db
	token, err := db.Client.Token.
		Query().
		Where(
			token.ContractAddressEQ(deposit.Token.Hex()),
			token.HasNetworkWith(
				networkent.IDEQ(network.ID),
			),
		).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch token: %w", err)
	}

	// Get order recipient from message hash
	recipient, err := s.getOrderRecipientFromMessageHash(deposit.MessageHash)
	if err != nil {
		return fmt.Errorf("failed to decrypt message hash: %w", err)
	}

	// Get provision bucket
	amountInDecimals := utils.FromSubunit(deposit.Amount, token.Decimals)
	institution, err := s.getInstitutionByCode(client, deposit.InstitutionCode)
	if err != nil {
		return fmt.Errorf("failed to fetch institution: %w", err)
	}

	currency, err := db.Client.FiatCurrency.
		Query().
		Where(
			fiatcurrency.IsEnabledEQ(true),
			fiatcurrency.CodeEQ(utils.Byte32ToString(institution.Currency)),
		).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch fiat currency: %w", err)
	}

	rate := decimal.NewFromBigInt(deposit.Rate, 0)

	provisionBucket, err := s.getProvisionBucket(ctx, amountInDecimals.Mul(rate), currency)
	if err != nil {
		return fmt.Errorf("failed to fetch provision bucket: %w", err)
	}

	// Create lock payment order fields
	lockPaymentOrder := types.LockPaymentOrderFields{
		Token:             token,
		OrderID:           fmt.Sprintf("0x%v", hex.EncodeToString(deposit.OrderId[:])),
		Amount:            amountInDecimals,
		Rate:              rate,
		Label:             utils.Byte32ToString(deposit.Label),
		BlockNumber:       int64(deposit.Raw.BlockNumber),
		TxHash:            deposit.Raw.TxHash.Hex(),
		Institution:       utils.Byte32ToString(deposit.InstitutionCode),
		AccountIdentifier: recipient.AccountIdentifier,
		AccountName:       recipient.AccountName,
		ProviderID:        recipient.ProviderID,
		Memo:              recipient.Memo,
		ProvisionBucket:   provisionBucket,
	}

	if provisionBucket == nil {
		currency, err := db.Client.FiatCurrency.
			Query().
			Where(
				fiatcurrency.IsEnabledEQ(true),
				fiatcurrency.CodeEQ(utils.Byte32ToString(institution.Currency)),
			).
			Only(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch fiat currency: %w", err)
		}

		// Split lock payment order into multiple orders
		err = s.splitLockPaymentOrder(
			ctx, lockPaymentOrder, currency,
		)
		if err != nil {
			return fmt.Errorf("failed to split lock payment order: %w", err)
		}
	} else {
		// Create lock payment order in db
		orderBuilder := db.Client.LockPaymentOrder.
			Create().
			SetToken(lockPaymentOrder.Token).
			SetOrderID(lockPaymentOrder.OrderID).
			SetAmount(lockPaymentOrder.Amount).
			SetRate(lockPaymentOrder.Rate).
			SetLabel(lockPaymentOrder.Label).
			SetOrderPercent(decimal.NewFromInt(100)).
			SetBlockNumber(lockPaymentOrder.BlockNumber).
			SetTxHash(lockPaymentOrder.TxHash).
			SetInstitution(lockPaymentOrder.Institution).
			SetAccountIdentifier(lockPaymentOrder.AccountIdentifier).
			SetAccountName(lockPaymentOrder.AccountName).
			SetMemo(lockPaymentOrder.Memo).
			SetProvisionBucket(lockPaymentOrder.ProvisionBucket)

		if lockPaymentOrder.ProviderID != "" {
			orderBuilder = orderBuilder.SetProviderID(lockPaymentOrder.ProviderID)
		}

		orderCreated, err := orderBuilder.Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create lock payment order: %w", err)
		}

		// Assign the lock payment order to a provider
		lockPaymentOrder.ID = orderCreated.ID
		_ = s.priorityQueue.AssignLockPaymentOrder(ctx, lockPaymentOrder)
	}

	return nil
}

// updateReceiveAddressStatus updates the status of a receive address. if `done` is true, the indexing process is complete for the given receive address
func (s *IndexerService) updateReceiveAddressStatus(
	ctx context.Context, receiveAddress *ent.ReceiveAddress, paymentOrder *ent.PaymentOrder, event *contracts.ERC20TokenTransfer,
) (done bool, err error) {

	if event.To.Hex() == receiveAddress.Address {
		// Check for existing address with txHash
		count, err := db.Client.ReceiveAddress.
			Query().
			Where(receiveaddress.TxHashEQ(event.Raw.TxHash.Hex())).
			Count(ctx)
		if err != nil {
			return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
		}

		if count > 0 && receiveAddress.Status != receiveaddress.StatusUnused {
			// This transfer has already been indexed
			return false, nil
		}

		// This is a transfer to the receive address to create an order on-chain
		// Compare the transferred value with the expected order amount + fees
		fees := paymentOrder.NetworkFee.Add(paymentOrder.SenderFee).Add(paymentOrder.ProtocolFee)
		orderAmountWithFees := paymentOrder.Amount.Add(fees).Round(int32(paymentOrder.Edges.Token.Decimals))
		orderAmountWithFeesInSubunit := utils.ToSubunit(orderAmountWithFees, paymentOrder.Edges.Token.Decimals)
		comparisonResult := event.Value.Cmp(orderAmountWithFeesInSubunit)

		paymentOrder, err = paymentOrder.
			Update().
			SetFromAddress(event.From.Hex()).
			SetTxHash(event.Raw.TxHash.Hex()).
			AddAmountPaid(utils.FromSubunit(event.Value, paymentOrder.Edges.Token.Decimals)).
			Save(ctx)
		if err != nil {
			return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
		}

		if comparisonResult == 0 {
			// Transfer value equals order amount with fees
			_, err = receiveAddress.
				Update().
				SetStatus(receiveaddress.StatusUsed).
				SetLastUsed(time.Now()).
				SetTxHash(event.Raw.TxHash.Hex()).
				SetLastIndexedBlock(int64(event.Raw.BlockNumber)).
				Save(ctx)
			if err != nil {
				return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
			}

			err = s.order.CreateOrder(ctx, paymentOrder.ID)
			if err != nil {
				return true, fmt.Errorf("updateReceiveAddressStatus.CreateOrder: %v", err)
			}

			return true, nil

		} else if comparisonResult < 0 {
			// Transfer value is less than order amount with fees

			// If amount paid meets or exceeds the order amount with fees, mark receive address as used
			if paymentOrder.AmountPaid.GreaterThanOrEqual(orderAmountWithFees) {
				_, err = receiveAddress.
					Update().
					SetStatus(receiveaddress.StatusUsed).
					SetLastUsed(time.Now()).
					SetTxHash(event.Raw.TxHash.Hex()).
					SetLastIndexedBlock(int64(event.Raw.BlockNumber)).
					Save(ctx)
				if err != nil {
					return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
				}
			} else {
				_, err = receiveAddress.
					Update().
					SetStatus(receiveaddress.StatusPartial).
					SetTxHash(event.Raw.TxHash.Hex()).
					SetLastIndexedBlock(int64(event.Raw.BlockNumber)).
					Save(ctx)
				if err != nil {
					return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
				}
			}

		} else if comparisonResult > 0 {
			// Transfer value is greater than order amount with fees
			_, err = receiveAddress.
				Update().
				SetStatus(receiveaddress.StatusUsed).
				SetLastUsed(time.Now()).
				SetLastIndexedBlock(int64(event.Raw.BlockNumber)).
				Save(ctx)
			if err != nil {
				return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
			}
		}

		err = s.HandleReceiveAddressValidity(ctx, receiveAddress, paymentOrder)
		if err != nil {
			return true, fmt.Errorf("updateReceiveAddressStatus.HandleReceiveAddressValidity: %v", err)
		}

	} else if event.From.Hex() == receiveAddress.Address {
		// This is a revert transfer from the receive address

		// Check for existing address with txHash
		count, err := db.Client.ReceiveAddress.
			Query().
			Where(receiveaddress.TxHashEQ(event.Raw.TxHash.Hex())).
			Count(ctx)
		if err != nil {
			return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
		}

		if count > 0 && paymentOrder.Status == paymentorder.StatusReverted {
			// This transfer has already been indexed
			return false, nil
		}

		// Compare the transferred value with the expected order amount returned
		indexedValue := utils.FromSubunit(event.Value, paymentOrder.Edges.Token.Decimals)

		if indexedValue.Equal(paymentOrder.AmountReturned) {
			_, err := receiveAddress.
				Update().
				SetTxHash(event.Raw.TxHash.Hex()).
				SetLastIndexedBlock(int64(event.Raw.BlockNumber)).
				Save(ctx)
			if err != nil {
				return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
			}

			_, err = paymentOrder.
				Update().
				SetStatus(paymentorder.StatusReverted).
				SetTxHash(event.Raw.TxHash.Hex()).
				Save(ctx)
			if err != nil {
				return true, fmt.Errorf("updateReceiveAddressStatus.db: %v", err)
			}

			// Send webhook notifcation to sender
			paymentOrder.Status = paymentorder.StatusReverted

			err = utils.SendPaymentOrderWebhook(ctx, paymentOrder)
			if err != nil {
				return true, fmt.Errorf("updateReceiveAddressStatus.webhook: %v", err)
			}

			return true, nil
		}
	}

	return false, nil
}

// getProvisionBucket returns the provision bucket for a lock payment order
func (s *IndexerService) getProvisionBucket(ctx context.Context, amount decimal.Decimal, currency *ent.FiatCurrency) (*ent.ProvisionBucket, error) {

	provisionBucket, err := db.Client.ProvisionBucket.
		Query().
		Where(
			provisionbucket.MaxAmountGTE(amount),
			provisionbucket.MinAmountLTE(amount),
			provisionbucket.HasCurrencyWith(
				fiatcurrency.IDEQ(currency.ID),
			),
		).
		WithCurrency().
		Only(ctx)
	if err != nil {
		logger.Errorf("failed to fetch provision bucket: %v", err)
		return nil, err
	}

	return provisionBucket, nil
}

// getInstitutionByCode returns the institution for a given institution code
func (s *IndexerService) getInstitutionByCode(client types.RPCClient, institutionCode [32]byte) (*contracts.SharedStructsInstitutionByCode, error) {
	instance, err := contracts.NewPaycrest(OrderConf.PaycrestOrderContractAddress, client.(bind.ContractBackend))
	if err != nil {
		return nil, err
	}

	institution, err := instance.GetSupportedInstitutionByCode(nil, institutionCode)
	if err != nil {
		return nil, err
	}

	return &institution, nil
}

// splitLockPaymentOrder splits a lock payment order into multiple orders
func (s *IndexerService) splitLockPaymentOrder(ctx context.Context, lockPaymentOrder types.LockPaymentOrderFields, currency *ent.FiatCurrency) error {
	buckets, err := db.Client.ProvisionBucket.
		Query().
		Where(provisionbucket.HasCurrencyWith(
			fiatcurrency.IDEQ(currency.ID),
		)).
		WithProviderProfiles().
		Order(ent.Desc(provisionbucket.FieldMaxAmount)).
		All(ctx)
	if err != nil {
		logger.Errorf("failed to fetch provision buckets: %v", err)
		return err
	}

	amountToSplit := lockPaymentOrder.Amount.Mul(lockPaymentOrder.Rate) // e.g 100,000

	for _, bucket := range buckets {
		// Get the number of providers in the bucket
		bucketSize := int64(len(bucket.Edges.ProviderProfiles))

		// Get the number of allocations to make in the bucket
		numberOfAllocations := amountToSplit.Div(bucket.MaxAmount).IntPart() // e.g 100000 / 10000 = 10

		var trips int64

		if bucketSize >= numberOfAllocations {
			trips = numberOfAllocations // e.g 10
		} else {
			trips = bucketSize // e.g 2
		}

		// Create a slice to hold the LockPaymentOrder entities for this bucket
		lockOrders := make([]*ent.LockPaymentOrderCreate, 0, trips)

		tx, err := db.Client.Tx(ctx)
		if err != nil {
			return err
		}

		for i := int64(0); i < trips; i++ {
			ratio := bucket.MaxAmount.Div(amountToSplit)
			orderPercent := ratio.Mul(decimal.NewFromInt(100))
			lockOrder := tx.LockPaymentOrder.
				Create().
				SetToken(lockPaymentOrder.Token).
				SetOrderID(lockPaymentOrder.OrderID).
				SetAmount(bucket.MaxAmount.Div(lockPaymentOrder.Rate)).
				SetRate(lockPaymentOrder.Rate).
				SetOrderPercent(orderPercent).
				SetBlockNumber(lockPaymentOrder.BlockNumber).
				SetTxHash(lockPaymentOrder.TxHash).
				SetInstitution(lockPaymentOrder.Institution).
				SetAccountIdentifier(lockPaymentOrder.AccountIdentifier).
				SetAccountName(lockPaymentOrder.AccountName).
				SetProviderID(lockPaymentOrder.ProviderID).
				SetProvisionBucket(bucket)
			lockOrders = append(lockOrders, lockOrder)
		}

		// Batch insert all LockPaymentOrder entities for this bucket in a single transaction
		ordersCreated, err := tx.LockPaymentOrder.
			CreateBulk(lockOrders...).
			Save(ctx)
		if err != nil {
			logger.Errorf("failed to create lock payment orders in bulk: %v", err)
			_ = tx.Rollback()
			return err
		}

		// Commit the transaction if everything succeeded
		if err := tx.Commit(); err != nil {
			logger.Errorf("failed to split lock payment order: %v", err)
			return err
		}

		// Assign the lock payment orders to providers
		for _, order := range ordersCreated {
			lockPaymentOrder.ID = order.ID
			_ = s.priorityQueue.AssignLockPaymentOrder(ctx, lockPaymentOrder)
		}

		amountToSplit = amountToSplit.Sub(bucket.MaxAmount)
	}

	largestBucket := buckets[0]

	if amountToSplit.LessThan(largestBucket.MaxAmount) {
		bucket, err := s.getProvisionBucket(ctx, amountToSplit, currency)
		if err != nil {
			return err
		}

		orderCreated, err := db.Client.LockPaymentOrder.
			Create().
			SetToken(lockPaymentOrder.Token).
			SetOrderID(lockPaymentOrder.OrderID).
			SetAmount(amountToSplit.Div(lockPaymentOrder.Rate)).
			SetRate(lockPaymentOrder.Rate).
			SetBlockNumber(lockPaymentOrder.BlockNumber).
			SetTxHash(lockPaymentOrder.TxHash).
			SetInstitution(lockPaymentOrder.Institution).
			SetAccountIdentifier(lockPaymentOrder.AccountIdentifier).
			SetAccountName(lockPaymentOrder.AccountName).
			SetProviderID(lockPaymentOrder.ProviderID).
			SetProvisionBucket(bucket).
			Save(ctx)
		if err != nil {
			logger.Errorf("failed to create lock payment order: %v", err)
			return err
		}

		// Assign the lock payment order to a provider
		lockPaymentOrder.ID = orderCreated.ID
		_ = s.priorityQueue.AssignLockPaymentOrder(ctx, lockPaymentOrder)

	} else {
		// TODO: figure out how to handle this case, currently it recursively splits the amount
		lockPaymentOrder.Amount = amountToSplit.Div(lockPaymentOrder.Rate)
		err := s.splitLockPaymentOrder(ctx, lockPaymentOrder, currency)
		if err != nil {
			return err
		}
	}

	return nil
}

// getMissedOrderBlocksOpts returns the filter options for fetching missed blocks based on lock payment order status
func (s *IndexerService) getMissedOrderBlocksOpts(
	ctx context.Context, client types.RPCClient, network *ent.Network, status lockpaymentorder.Status,
) *bind.FilterOpts {

	// Get the last lock payment order from db
	result, err := db.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.HasTokenWith(
				token.HasNetworkWith(
					networkent.IDEQ(network.ID),
				),
			),
			lockpaymentorder.StatusEQ(status),
			lockpaymentorder.BlockNumberNEQ(0),
		).
		Order(ent.Desc(lockpaymentorder.FieldBlockNumber)).
		Limit(1).
		All(ctx)
	if err != nil {
		logger.Errorf("failed to fetch lock payment order: %v", err)
	}

	// Fetch current block header
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		logger.Errorf("failed to fetch current block number: %v", err)
	}

	var startBlockNumber int64
	toBlock := header.Number.Uint64()

	if len(result) > 0 {
		startBlockNumber = int64(result[0].BlockNumber) + 1
	} else {
		startBlockNumber = int64(toBlock) - 500
	}

	// Filter logs from the last lock payment order block number
	opts := &bind.FilterOpts{
		Start: uint64(startBlockNumber),
		End:   &toBlock,
	}

	return opts
}

// getMissedERC20BlocksOpts returns the filter options for fetching missed blocks based on receive address status
func (s *IndexerService) getMissedERC20BlocksOpts(ctx context.Context, client types.RPCClient, network *ent.Network) *bind.FilterOpts {

	// Get receive address with most recent indexed block from db
	result, err := db.Client.ReceiveAddress.
		Query().
		Where(
			receiveaddress.HasPaymentOrderWith(
				paymentorder.HasTokenWith(
					token.HasNetworkWith(
						networkent.IDEQ(network.ID),
					),
				),
			),
			receiveaddress.LastIndexedBlockNEQ(0),
		).
		Order(ent.Desc(receiveaddress.FieldLastIndexedBlock)).
		Limit(1).
		All(ctx)
	if err != nil {
		logger.Errorf("getMissedERC20BlocksOpts.db: %v", err)
	}

	// Fetch current block header
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		logger.Errorf("getMissedERC20BlocksOpts.HeaderByNumber: %v", err)
	}

	var startBlockNumber int64
	toBlock := header.Number.Uint64()

	if len(result) > 0 {
		startBlockNumber = int64(result[0].LastIndexedBlock)
	} else {
		startBlockNumber = int64(toBlock) - 500
	}

	// Filter logs from the last indexed block number
	opts := &bind.FilterOpts{
		Start: uint64(startBlockNumber),
		End:   &toBlock,
	}

	return opts
}
