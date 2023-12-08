package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	IndexERC20Transfer(ctx context.Context, client types.RPCClient, receiveAddress *ent.ReceiveAddress, done chan<- bool) error
	IndexOrderDeposits(ctx context.Context, client types.RPCClient, network *ent.Network) error
	IndexOrderSettlements(ctx context.Context, client types.RPCClient, network *ent.Network) error
	IndexOrderRefunds(ctx context.Context, client types.RPCClient, network *ent.Network) error
}

// IndexerService performs blockchain to database extract, transform, load (ETL) operations.
type IndexerService struct {
	indexer       Indexer
	priorityQueue *PriorityQueueService
	order         *OrderService
}

// NewIndexerService creates a new instance of IndexerService.
func NewIndexerService(indexer Indexer) *IndexerService {
	priorityQueue := NewPriorityQueueService()
	order := NewOrderService()

	return &IndexerService{
		indexer:       indexer,
		priorityQueue: priorityQueue,
		order:         order,
	}
}

// IndexERC20Transfer indexes ERC20 token transfers for a specific receive address.
func (s *IndexerService) IndexERC20Transfer(ctx context.Context, client types.RPCClient, receiveAddress *ent.ReceiveAddress) (ok bool, err error) {

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
		return false, fmt.Errorf("failed to fetch payment order: %w", err)
	}

	token := paymentOrder.Edges.Token

	if client == nil {
		client, err = types.NewEthClient(token.Edges.Network.RPCEndpoint)
		if err != nil {
			return false, fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Fetch current block header
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to fetch current block number: %w", err)
	}

	// Number of blocks that will be generated within the receive address valid period.
	// Layer 2 networks have the shortest blocktimes e.g Polygon, Tron etc. are < 5 seconds
	// We assume a blocktime of 2s for the largest number of blocks to scan.
	// number of blocks = receive address validity period in seconds / blocktime
	numOfBlocks := int64(OrderConf.ReceiveAddressValidity.Seconds() / 2)

	var fromBlock int64

	if receiveAddress.LastIndexedBlock > 0 {
		// Continue indexing from last indexed block if the last process failed
		fromBlock = receiveAddress.LastIndexedBlock + 1
	} else {
		fromBlock = int64(math.Max(float64(header.Number.Int64()-numOfBlocks+1), 0))
	}

	fromBlockBig := big.NewInt(fromBlock)

	// Query event logs of the token contract starting from the oldest block
	// within the receive address validity period
	query := ethereum.FilterQuery{
		FromBlock: fromBlockBig,
		ToBlock:   header.Number,
		Addresses: []common.Address{
			common.HexToAddress(paymentOrder.Edges.Token.ContractAddress),
		},
	}

	// Fetch logs in block batches.
	// This is important because client.FilterLogs function has a limit of 10k results
	// TODO: explain why batchsize of 200 was chosen
	currentBlockBatchSize := 200
	currentBlockNumber := query.FromBlock
	finalBlockNumber := query.ToBlock

	// Get Transfer event signature hash
	contractAbi, err := abi.JSON(strings.NewReader(ERC20ABI))
	if err != nil {
		return false, fmt.Errorf("failed to parse ABI: %w", err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	logger.Infof(
		"Indexing transfer logs for %s from Block #%s - #%s",
		receiveAddress.Address,
		fromBlockBig.String(),
		finalBlockNumber.String(),
	)

	for {
		// Update the filter parameters
		query.FromBlock = currentBlockNumber

		batchEnd := big.NewInt(currentBlockNumber.Int64() + int64(currentBlockBatchSize-1))
		query.ToBlock = utils.BigMin(batchEnd, finalBlockNumber)

		// Fetch logs for the current batch
		logs, err := client.FilterLogs(ctx, query)
		if err != nil {
			logger.Errorf("IndexERC20Transfer.fetchLogs: %v", err)
			return false, nil
		}

		for _, vLog := range logs {
			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				var transferEvent types.ERC20Transfer

				err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
				if err != nil {
					return false, fmt.Errorf("IndexERC20Transfer.unpackSignature: %w", err)
				}

				transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

				if transferEvent.To.Hex() == receiveAddress.Address {
					// This is a transfer to the receive address to create an order on-chain
					// Compare the transferred value with the expected order amount
					comparisonResult := transferEvent.Value.Cmp(paymentOrder.Amount.BigInt())

					if comparisonResult == 0 {
						// Transfer value equals order amount
						_, err = receiveAddress.
							Update().
							SetStatus(receiveaddress.StatusUsed).
							SetLastUsed(time.Now()).
							SetLastIndexedBlock(query.ToBlock.Int64()).
							Save(ctx)
						if err != nil {
							logger.Errorf("IndexERC20Transfer.db: %v", err)
							return false, nil
						}
						return true, nil
					} else if comparisonResult < 0 {
						// Transfer value is less than order amount
						_, err = receiveAddress.
							Update().
							SetStatus(receiveaddress.StatusPartial).
							Save(ctx)
						if err != nil {
							logger.Errorf("IndexERC20Transfer.db: %v", err)
							return false, nil
						}

						amountPaid := paymentOrder.AmountPaid.Add(utils.FromSubunit(transferEvent.Value, token.Decimals))

						// If amount paid meets or exceeds the expected amount, mark receive address as used
						if amountPaid.GreaterThanOrEqual(paymentOrder.Amount) {
							_, err = receiveAddress.
								Update().
								SetStatus(receiveaddress.StatusUsed).
								SetLastUsed(time.Now()).
								SetLastIndexedBlock(query.ToBlock.Int64()).
								Save(ctx)
							if err != nil {
								logger.Errorf("IndexERC20Transfer.db: %v", err)
								return false, nil
							}

							// Revert excess amount to the from address
							err := s.order.RevertOrder(ctx, paymentOrder, transferEvent.From)
							if err != nil {
								logger.Errorf("IndexERC20Transfer.revertOrder: %v", err)
								return false, nil
							}

							return true, nil
						} else {
							// Update the payment order with amount paid
							_, err = paymentOrder.
								Update().
								SetAmountPaid(amountPaid).
								SetTxHash(vLog.TxHash.Hex()).
								Save(ctx)
							if err != nil {
								logger.Errorf("IndexERC20Transfer.db: %v", err)
								return false, nil
							}
						}

					} else if comparisonResult > 0 {
						// Transfer value is greater than order amount
						_, err = receiveAddress.
							Update().
							SetStatus(receiveaddress.StatusUsed).
							SetLastUsed(time.Now()).
							SetLastIndexedBlock(query.ToBlock.Int64()).
							Save(ctx)
						if err != nil {
							logger.Errorf("IndexERC20Transfer.db: %v", err)
							return false, nil
						}

						// Revert excess amount to the from address
						err := s.order.RevertOrder(ctx, paymentOrder, transferEvent.From)
						if err != nil {
							logger.Errorf("IndexERC20Transfer.revertOrder: %v", err)
							return false, nil
						}

						return true, nil
					}

					return false, nil

				} else if transferEvent.From.Hex() == receiveAddress.Address {
					// This is a revert transfer
					// Compare the transferred value with the expected order amount returned
					orderAmountReturnedInSubunit := utils.ToSubunit(paymentOrder.AmountReturned, token.Decimals)

					if transferEvent.Value.Cmp(orderAmountReturnedInSubunit) == 0 {
						if paymentOrder.AmountPaid.LessThan(paymentOrder.Amount) {
							_, err = receiveAddress.
								Update().
								SetLastIndexedBlock(query.ToBlock.Int64()).
								Save(ctx)
							if err != nil {
								logger.Errorf("IndexERC20Transfer.db: %v", err)
								return false, nil
							}

							_, err := paymentOrder.
								Update().
								SetStatus(paymentorder.StatusReverted).
								SetTxHash(vLog.TxHash.Hex()).
								Save(ctx)
							if err != nil {
								logger.Errorf("IndexERC20Transfer.db: %v", err)
								return false, nil
							}

							// Send webhook notifcation to sender
							paymentOrder.Status = paymentorder.StatusReverted

							err = utils.SendPaymentOrderWebhook(ctx, paymentOrder)
							if err != nil {
								logger.Errorf("IndexERC20Transfer.webhook: %v", err)
								return false, nil
							}
						}

					}
				}

				// Handle receive address validity checks
				if receiveAddress.Status != receiveaddress.StatusUsed {
					amountNotPaidInFull := receiveAddress.Status == receiveaddress.StatusPartial || receiveAddress.Status == receiveaddress.StatusUnused
					validUntilIsFarGone := receiveAddress.ValidUntil.Before(time.Now().Add(-(5 * time.Minute)))
					isExpired := receiveAddress.ValidUntil.Before(time.Now())

					if validUntilIsFarGone {
						_, err = receiveAddress.
							Update().
							SetValidUntil(time.Now().Add(OrderConf.ReceiveAddressValidity)).
							Save(ctx)
						if err != nil {
							logger.Errorf("IndexERC20Transfer.db: %v", err)
							return false, nil
						}
					} else if isExpired && amountNotPaidInFull {
						// Receive address hasn't received full payment after validity period, mark status as expired
						_, err = receiveAddress.
							Update().
							SetStatus(receiveaddress.StatusExpired).
							Save(ctx)
						if err != nil {
							logger.Errorf("IndexERC20Transfer.db: %v", err)
							return false, nil
						}

						// Expire payment order
						_, err = paymentOrder.
							Update().
							SetStatus(paymentorder.StatusExpired).
							Save(ctx)
						if err != nil {
							logger.Errorf("IndexERC20Transfer.db: %v", err)
							return false, nil
						}

						// Revert amount to the from address
						err := s.order.RevertOrder(ctx, paymentOrder, transferEvent.From)
						if err != nil {
							logger.Errorf("IndexERC20Transfer.revertOrder: %v", err)
							return false, nil
						}

						return true, nil
					}
				}
			}
		}

		// Update last indexed block
		_, err = receiveAddress.
			Update().
			SetLastIndexedBlock(query.ToBlock.Int64()).
			Save(ctx)
		if err != nil {
			logger.Errorf("IndexERC20Transfer.db: %v", err)
			return false, nil
		}

		// Check if we have reached the final block number
		if batchEnd.Cmp(finalBlockNumber) >= 0 {
			break
		}

		// Update the current block number for the next batch
		currentBlockNumber = big.NewInt(query.ToBlock.Int64() + 1)

		// Sleep for a short duration between batches to avoid overwhelming the RPC endpoint
		time.Sleep(1 * time.Second)
	}

	return false, nil
}

// RunIndexERC20Transfer runs the indexer service for a receive address
// it loops indefinitely until the address expires or a transfer is found
func (s *IndexerService) RunIndexERC20Transfer(ctx context.Context, receiveAddress *ent.ReceiveAddress) {
	for {
		time.Sleep(2 * time.Minute) // add 2 minutes delay between each indexing operation

		ok, err := s.IndexERC20Transfer(ctx, nil, receiveAddress)
		if err != nil {
			logger.Errorf("failed to index erc20 transfer: %v", err)
			return
		}

		// Refresh the receive address with payment order
		receiveAddress, err = db.Client.ReceiveAddress.
			Query().
			Where(receiveaddress.AddressEQ(receiveAddress.Address)).
			WithPaymentOrder().
			Only(ctx)
		if err != nil {
			logger.Errorf("failed to refresh receive address: %v", err)
			return
		}

		if ok {
			if receiveAddress.Status == receiveaddress.StatusUsed {
				// Create order on-chain
				orderService := NewOrderService()
				err = orderService.CreateOrder(ctx, receiveAddress.Edges.PaymentOrder.ID)
				if err != nil {
					logger.Errorf("failed to create order on-chain: %v", err)
				}
			}

			// Address is expired, stop indexing
			return
		}
	}
}

// IndexOrderDeposit indexes deposits to the order contract for a specific network.
func (s *IndexerService) IndexOrderDeposits(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Initialize contract filterer
	filterer, err := contracts.NewPaycrestOrderFilterer(OrderConf.PaycrestOrderContractAddress, client)
	if err != nil {
		return fmt.Errorf("failed to create filterer: %w", err)
	}

	// Index missed blocks
	go func() {
		// Filter logs from the last lock payment order block number
		opts := s.getMissedBlocksOpts(ctx, client, filterer, network, lockpaymentorder.StatusPending)

		// Fetch logs
		iter, err := filterer.FilterDeposit(opts, nil, nil, nil)
		if err != nil {
			logger.Errorf("IndexOrderDeposits.fetchLogs: %v", err)
		}

		// Iterate over logs
		for iter.Next() {
			err := s.createLockPaymentOrder(ctx, client, network, iter.Event)
			if err != nil {
				logger.Errorf("IndexOrderDeposits.createOrder: %v", err)
				continue
			}
		}
	}()

	// Start listening for deposit events
	logs := make(chan *contracts.PaycrestOrderDeposit)

	sub, err := filterer.WatchDeposit(&bind.WatchOpts{
		Start: nil,
	}, logs, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to watch deposit events: %w", err)
	}

	defer sub.Unsubscribe()

	logger.Infof("Listening for Deposit events...\n")

	for {
		select {
		case log := <-logs:
			err := s.createLockPaymentOrder(ctx, client, network, log)
			if err != nil {
				logger.Errorf("failed to create lock payment order: %v", err)
				continue
			}
		case err := <-sub.Err():
			logger.Errorf("failed to parse deposit event: %v", err)
			continue
		}
	}
}

// IndexOrderSettlements indexes order settlements for a specific network.
func (s *IndexerService) IndexOrderSettlements(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Initialize contract filterer
	filterer, err := contracts.NewPaycrestOrderFilterer(OrderConf.PaycrestOrderContractAddress, client)
	if err != nil {
		return fmt.Errorf("failed to create filterer: %w", err)
	}

	// Index missed blocks
	go func() {
		// Filter logs from the last lock payment order block number
		opts := s.getMissedBlocksOpts(ctx, client, filterer, network, lockpaymentorder.StatusSettled)

		// Fetch logs
		iter, err := filterer.FilterSettled(opts, nil, nil)
		if err != nil {
			logger.Errorf("IndexOrderSettlements.fetchLogs: %v", err)
		}

		// Iterate over logs
		for iter.Next() {
			log := iter.Event
			err := s.updateOrderStatusSettled(ctx, log)
			if err != nil {
				logger.Errorf("IndexOrderSettlements.update: %v", err)
				continue
			}
		}
	}()

	// Start listening for settlement events
	logs := make(chan *contracts.PaycrestOrderSettled)

	sub, err := filterer.WatchSettled(&bind.WatchOpts{
		Start: nil,
	}, logs, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to watch settlement events: %w", err)
	}

	defer sub.Unsubscribe()

	logger.Infof("Listening for Settlement events...\n")

	for {
		select {
		case log := <-logs:
			err := s.updateOrderStatusSettled(ctx, log)
			if err != nil {
				logger.Errorf("IndexOrderSettlements.update: %v", err)
				continue
			}
		case err := <-sub.Err():
			logger.Errorf("failed to parse settlement event: %v", err)
			continue
		}
	}
}

// IndexOrderRefunds indexes order refunds for a specific network.
func (s *IndexerService) IndexOrderRefunds(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	var err error

	// Connect to RPC endpoint
	if client == nil {
		client, err = types.NewEthClient(network.RPCEndpoint)
		if err != nil {
			return fmt.Errorf("failed to connect to RPC client: %w", err)
		}
	}

	// Initialize contract filterer
	filterer, err := contracts.NewPaycrestOrderFilterer(OrderConf.PaycrestOrderContractAddress, client)
	if err != nil {
		return fmt.Errorf("failed to create filterer: %w", err)
	}

	// Index missed blocks
	go func() {
		// Filter logs from the last lock payment order block number
		opts := s.getMissedBlocksOpts(ctx, client, filterer, network, lockpaymentorder.StatusRefunded)

		// Fetch logs
		iter, err := filterer.FilterRefunded(opts, nil)
		if err != nil {
			logger.Errorf("IndexOrderRefunds.fetchLogs: %v", err)
		}

		// Iterate over logs
		for iter.Next() {
			log := iter.Event
			err := s.updateOrderStatusRefunded(ctx, log)
			if err != nil {
				logger.Errorf("IndexOrderRefunds.update: %v", err)
				continue
			}
		}
	}()

	// Start listening for refund events
	logs := make(chan *contracts.PaycrestOrderRefunded)

	sub, err := filterer.WatchRefunded(&bind.WatchOpts{
		Start: nil,
	}, logs, nil)
	if err != nil {
		return fmt.Errorf("failed to watch refund events: %w", err)
	}

	defer sub.Unsubscribe()

	logger.Infof("Listening for Refund events...\n")

	for {
		select {
		case log := <-logs:
			err := s.updateOrderStatusRefunded(ctx, log)
			if err != nil {
				logger.Errorf("IndexOrderRefunds.update: %v", err)
				continue
			}
		case err := <-sub.Err():
			logger.Errorf("IndexOrderRefunds.parseEvent: %v", err)
			continue
		}
	}
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
			lockpaymentorder.OrderIDEQ(utils.Byte32ToString(log.OrderId)),
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
func (s *IndexerService) updateOrderStatusSettled(ctx context.Context, log *contracts.PaycrestOrderSettled) error {
	// Aggregator side status update
	splitOrderId, _ := uuid.Parse(utils.Byte32ToString(log.SplitOrderId))
	_, err := db.Client.LockPaymentOrder.
		Update().
		Where(
			lockpaymentorder.IDEQ(splitOrderId),
		).
		SetBlockNumber(int64(log.Raw.BlockNumber)).
		SetTxHash(log.Raw.TxHash.Hex()).
		SetStatus(lockpaymentorder.StatusSettled).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusSettled.aggregator: %v", err)
	}

	// Sender side status update
	_, err = db.Client.PaymentOrder.
		Update().
		Where(
			paymentorder.LabelEQ(utils.Byte32ToString(log.Label)),
		).
		SetTxHash(log.Raw.TxHash.Hex()).
		SetStatus(paymentorder.StatusSettled).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("updateOrderStatusSettled.sender: %v", err)
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
func (s *IndexerService) createLockPaymentOrder(ctx context.Context, client types.RPCClient, network *ent.Network, deposit *contracts.PaycrestOrderDeposit) error {
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
	institution, err := s.getInstitutionByCode(ctx, client, deposit.InstitutionCode)
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

	provisionBucket, err := s.getProvisionBucket(
		ctx, nil, amountInDecimals.Mul(rate), currency,
	)
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

// getProvisionBucket returns the provision bucket for a lock payment order
func (s *IndexerService) getProvisionBucket(ctx context.Context, client types.RPCClient, amount decimal.Decimal, currency *ent.FiatCurrency) (*ent.ProvisionBucket, error) {

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
func (s *IndexerService) getInstitutionByCode(ctx context.Context, client types.RPCClient, institutionCode [32]byte) (*contracts.PaycrestSettingManagerInstitutionByCode, error) {
	instance, err := contracts.NewPaycrestOrder(OrderConf.PaycrestOrderContractAddress, client.(bind.ContractBackend))
	if err != nil {
		return nil, err
	}

	institution, err := instance.GetSupportedInstitutionName(nil, institutionCode)
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
		numberOfAllocations := amountToSplit.Div(bucket.MaxAmount).IntPart() // e.g 100000 / 10000 = 10, TODO: verify integer conversion

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
		bucket, err := s.getProvisionBucket(ctx, nil, amountToSplit, currency)
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

// getMissedBlocksOpts returns the filter options for fetching missed blocks based on lock payment order status
func (s *IndexerService) getMissedBlocksOpts(
	ctx context.Context, client types.RPCClient, filterer *contracts.PaycrestOrderFilterer, network *ent.Network, status lockpaymentorder.Status,
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
		startBlockNumber = int64(result[0].BlockNumber)
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
