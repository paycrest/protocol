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
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/lockpaymentorder"
	networkent "github.com/paycrest/paycrest-protocol/ent/network"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/provisionbucket"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/ent/token"
	"github.com/paycrest/paycrest-protocol/services/contracts"
	db "github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/paycrest/paycrest-protocol/utils"
	cryptoUtils "github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/shopspring/decimal"
)

type orderRecipient struct {
	AccountIdentifier string
	AccountName       string
	Institution       string
	ProviderID        string
}

// Indexer is an interface for indexing blockchain data to the database.
type Indexer interface {
	IndexERC20Transfer(ctx context.Context, client types.RPCClient, receiveAddress *ent.ReceiveAddress, done chan<- bool) error
	IndexOrderDeposits(ctx context.Context, client types.RPCClient, network *ent.Network) error
}

// IndexerService performs blockchain to database extract, transform, load (ETL) operations.
type IndexerService struct {
	indexer       Indexer
	priorityQueue *PriorityQueueService
}

// NewIndexerService creates a new instance of IndexerService.
func NewIndexerService(indexer Indexer) *IndexerService {
	priorityQueue := NewPriorityQueueService()

	return &IndexerService{
		indexer:       indexer,
		priorityQueue: priorityQueue,
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
			return false, fmt.Errorf("failed to fetch logs: %w", err)
		}

		for _, vLog := range logs {
			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				var transferEvent types.ERC20Transfer

				err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
				if err != nil {
					return false, fmt.Errorf("failed to unpack Transfer event signature: %w", err)
				}

				transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

				if transferEvent.To.Hex() == receiveAddress.Address {
					// Compare the transferred value with the expected order amount
					orderAmountInSubunit := utils.ToSubunit(paymentOrder.Amount, token.Decimals)
					var comparisonResult = transferEvent.Value.Cmp(orderAmountInSubunit)

					if comparisonResult == 0 {
						// Transfer value equals order amount
						_, err = receiveAddress.
							Update().
							SetStatus(receiveaddress.StatusUsed).
							SetLastUsed(time.Now()).
							Save(ctx)
						if err != nil {
							return false, fmt.Errorf("failed to update receive address status: %w", err)
						}
						return true, nil
					} else if comparisonResult < 0 {
						// Transfer value is less than order amount
						_, err = receiveAddress.
							Update().
							SetStatus(receiveaddress.StatusPartial).
							Save(ctx)
						if err != nil {
							return false, fmt.Errorf("failed to update receive address status: %w", err)
						}
					}

					// Update the payment order with amount paid
					_, err = paymentOrder.
						Update().
						SetAmountPaid(paymentOrder.AmountPaid.Add(utils.FromSubunit(transferEvent.Value, token.Decimals))).
						Save(ctx)
					if err != nil {
						return false, fmt.Errorf("failed to record amount paid: %w", err)
					}

					if receiveAddress.Status == receiveaddress.StatusPartial {
						// Refresh the receive address with payment order and compare the amount paid with expected amount,
						receiveAddress, err = db.Client.ReceiveAddress.
							Query().
							Where(receiveaddress.AddressEQ(receiveAddress.Address)).
							WithPaymentOrder().
							Only(ctx)
						if err != nil {
							return false, fmt.Errorf("failed to refresh receive address: %w", err)
						}

						// If amount paid meets or exceeds the expected amount, mark receive address as used
						if paymentOrder.AmountPaid.GreaterThanOrEqual(paymentOrder.Amount) {
							_, err = receiveAddress.
								Update().
								SetStatus(receiveaddress.StatusUsed).
								SetLastUsed(time.Now()).
								Save(ctx)
							if err != nil {
								return false, fmt.Errorf("failed to update receive address status: %w", err)
							}
							return true, nil
						}
					}

					return false, nil
				}

				// Handle receive address validity checks
				amountNotPaidInFull := receiveAddress.Status == receiveaddress.StatusPartial || receiveAddress.Status == receiveaddress.StatusUnused
				validUntilIsFarGone := receiveAddress.ValidUntil.Before(time.Now().Add(-(5 * time.Minute)))
				isExpired := receiveAddress.ValidUntil.Before(time.Now())

				if validUntilIsFarGone {
					_, err = receiveAddress.
						Update().
						SetValidUntil(time.Now().Add(OrderConf.ReceiveAddressValidity)).
						Save(ctx)
					if err != nil {
						return false, fmt.Errorf("failed to update receive address valid until: %w", err)
					}
				} else if isExpired && amountNotPaidInFull {
					// Receive address hasn't received full payment after validity period, mark status as expired
					_, err = receiveAddress.
						Update().
						SetStatus(receiveaddress.StatusExpired).
						Save(ctx)
					if err != nil {
						return false, fmt.Errorf("failed to update receive address status: %w", err)
					}
					return true, nil
				}
			}
		}

		// Update last indexed block
		_, err = receiveAddress.
			Update().
			SetLastIndexedBlock(query.ToBlock.Int64()).
			Save(ctx)

		if err != nil {
			return false, fmt.Errorf("failed to update receive address last indexed block: %w", err)
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
				err = orderService.CreateOrder(ctx, nil, receiveAddress.Edges.PaymentOrder.ID)
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

	go s.indexMissingBlocks(ctx, client, filterer, network)

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
			err := s.saveLockPaymentOrder(ctx, client, network, log)
			if err != nil {
				logger.Errorf("failed to save lock payment order: %v", err)
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

	// TODO: index missing blocks here
	// e.g go s.indexMissingBlocks(ctx, client, filterer, network)

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
			splitOrderId, _ := uuid.Parse(utils.Byte32ToString(log.SplitOrderId))

			_, err := db.Client.LockPaymentOrder.
				Update().
				Where(
					lockpaymentorder.IDEQ(splitOrderId),
				).
				SetBlockNumber(int64(log.Raw.BlockNumber)).
				SetStatus(lockpaymentorder.StatusSettled).
				Save(ctx)
			if err != nil {
				logger.Errorf("failed to fetch lock payment order: %v", err)
				continue
			}

			// TODO: Settle payment order on sender side
			// orderId, _ := uuid.Parse(utils.Byte32ToString(log.OrderId))
			// order, err = db.Client.PaymentOrder.
			// 	Query().
			// 	Where(
			// 		paymentorder.OrderIDEQ(orderId),
			// 	).
			// 	Only(ctx)

			// if err != nil {
			// 	logger.Errorf("failed to fetch lock payment order: %v", err)
			// 	continue
			// }
		case err := <-sub.Err():
			logger.Errorf("failed to parse settlement event: %v", err)
			continue
		}
	}
}

// indexMissingBlocks indexes missing blocks from the last lock payment order block number
func (s *IndexerService) indexMissingBlocks(ctx context.Context, client types.RPCClient, filterer *contracts.PaycrestOrderFilterer, network *ent.Network) {
	// Get the last lock payment order from db
	result, err := db.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.HasTokenWith(
				token.HasNetworkWith(
					networkent.IDEQ(network.ID),
				),
			),
			lockpaymentorder.StatusEQ(lockpaymentorder.StatusPending),
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

	if len(result) > 0 {
		lockPaymentOrder := result[0]

		// Filter logs from the last lock payment order block number
		toBlock := header.Number.Uint64()
		opts := &bind.FilterOpts{
			Start: uint64(lockPaymentOrder.BlockNumber),
			End:   &toBlock,
		}

		// Fetch logs
		iter, err := filterer.FilterDeposit(opts, nil, nil, nil)
		if err != nil {
			logger.Errorf("failed to fetch logs: %v", err)
		}

		// Iterate over logs
		for iter.Next() {
			err := s.saveLockPaymentOrder(ctx, client, network, iter.Event)
			if err != nil {
				logger.Errorf("failed to save lock payment order: %v", err)
				continue
			}
		}
	}
}

// getOrderRecipientFromMessageHash decrypts the message hash and returns the order recipient
func (s *IndexerService) getOrderRecipientFromMessageHash(messageHash string) (*orderRecipient, error) {
	messageCipher, err := hex.DecodeString(strings.TrimPrefix(messageHash, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode message hash: %w", err)
	}

	message, err := cryptoUtils.DecryptJSON(messageCipher)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message hash: %w", err)
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	var recipient *orderRecipient
	if err := json.Unmarshal(messageBytes, &recipient); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return recipient, nil
}

// saveLockPaymentOrder saves a lock payment order in the database
func (s *IndexerService) saveLockPaymentOrder(ctx context.Context, client types.RPCClient, network *ent.Network, deposit *contracts.PaycrestOrderDeposit) error {
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

	provisionBucket, err := s.getProvisionBucket(
		ctx, nil, amountInDecimals, utils.Byte32ToString(institution.Currency),
	)
	if err != nil {
		return fmt.Errorf("failed to fetch provision bucket: %w", err)
	}

	// Create lock payment order fields
	lockPaymentOrder := types.LockPaymentOrderFields{
		Token:             token,
		OrderID:           fmt.Sprintf("0x%v", hex.EncodeToString(deposit.OrderId[:])),
		Amount:            amountInDecimals,
		Rate:              decimal.NewFromBigInt(deposit.Rate, 0),
		BlockNumber:       int64(deposit.Raw.BlockNumber),
		Institution:       utils.Byte32ToString(deposit.InstitutionCode),
		AccountIdentifier: recipient.AccountIdentifier,
		AccountName:       recipient.AccountName,
		ProviderID:        recipient.ProviderID,
		ProvisionBucket:   provisionBucket,
	}

	if provisionBucket == nil {
		// Split lock payment order into multiple orders
		err := s.splitLockPaymentOrder(
			ctx, lockPaymentOrder, utils.Byte32ToString(institution.Currency),
		)
		if err != nil {
			return fmt.Errorf("failed to split lock payment order: %w", err)
		}
	} else {
		// Create lock payment order in db
		orderCreated, err := db.Client.LockPaymentOrder.
			Create().
			SetToken(lockPaymentOrder.Token).
			SetOrderID(lockPaymentOrder.OrderID).
			SetAmount(lockPaymentOrder.Amount).
			SetRate(lockPaymentOrder.Rate).
			SetBlockNumber(lockPaymentOrder.BlockNumber).
			SetInstitution(lockPaymentOrder.Institution).
			SetAccountIdentifier(lockPaymentOrder.AccountIdentifier).
			SetAccountName(lockPaymentOrder.AccountName).
			SetProviderID(lockPaymentOrder.ProviderID).
			SetProvisionBucket(lockPaymentOrder.ProvisionBucket).
			Save(ctx)
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
func (s *IndexerService) getProvisionBucket(ctx context.Context, client types.RPCClient, amount decimal.Decimal, currency string) (*ent.ProvisionBucket, error) {
	provisionBucket, err := db.Client.ProvisionBucket.
		Query().
		Where(
			provisionbucket.MaxAmountGTE(amount),
			provisionbucket.MinAmountLTE(amount),
			provisionbucket.CurrencyEQ(currency),
		).
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
func (s *IndexerService) splitLockPaymentOrder(ctx context.Context, lockPaymentOrder types.LockPaymentOrderFields, currency string) error {
	buckets, err := db.Client.ProvisionBucket.
		Query().
		Where(provisionbucket.CurrencyEQ(currency)).
		WithProviderProfiles().
		Order(ent.Desc(provisionbucket.FieldMaxAmount)).
		All(ctx)
	if err != nil {
		logger.Errorf("failed to fetch provision buckets: %v", err)
		return err
	}

	amountToSplit := lockPaymentOrder.Amount // e.g 100,000

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
				SetAmount(bucket.MaxAmount).
				SetRate(lockPaymentOrder.Rate).
				SetOrderPercent(orderPercent).
				SetBlockNumber(lockPaymentOrder.BlockNumber).
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
			SetAmount(amountToSplit).
			SetRate(lockPaymentOrder.Rate).
			SetBlockNumber(lockPaymentOrder.BlockNumber).
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
		lockPaymentOrder.Amount = amountToSplit
		err := s.splitLockPaymentOrder(ctx, lockPaymentOrder, currency)
		if err != nil {
			return err
		}
	}

	return nil
}
