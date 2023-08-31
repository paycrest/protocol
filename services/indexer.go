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
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/lockpaymentorder"
	networkent "github.com/paycrest/paycrest-protocol/ent/network"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/provisionbucket"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/ent/token"
	"github.com/paycrest/paycrest-protocol/services/contracts"
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
}

// IndexerService performs blockchain to database extract, transform, load (ETL) operations.
type IndexerService struct {
	indexer Indexer
}

// NewIndexerService creates a new instance of IndexerService.
func NewIndexerService(indexer Indexer) *IndexerService {
	return &IndexerService{
		indexer: indexer,
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
			err := s.saveLockPaymentOrder(ctx, network, log)
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
			err := s.saveLockPaymentOrder(ctx, network, iter.Event)
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
func (s *IndexerService) saveLockPaymentOrder(ctx context.Context, network *ent.Network, deposit *contracts.PaycrestOrderDeposit) error {
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
	provisionBucket, err := s.getProvisionBucket(ctx, nil, amountInDecimals, deposit.InstitutionCode)
	if err != nil {
		return fmt.Errorf("failed to fetch provision bucket: %w", err)
	}

	// Create lock payment order in db
	_, err = db.Client.LockPaymentOrder.
		Create().
		SetToken(token).
		SetOrderID(fmt.Sprintf("0x%v", hex.EncodeToString(deposit.OrderId[:]))).
		SetAmount(amountInDecimals).
		SetAmountPaid(decimal.NewFromInt(0)).
		SetRate(decimal.NewFromBigInt(deposit.Rate, 0)).
		SetBlockNumber(int64(deposit.Raw.BlockNumber)).
		SetInstitution(utils.Byte32ToString(deposit.InstitutionCode)).
		SetAccountIdentifier(recipient.AccountIdentifier).
		SetAccountName(recipient.AccountName).
		SetProviderID(recipient.ProviderID).
		SetProvisionBucket(provisionBucket).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create lock payment order: %w", err)
	}

	return nil
}

// getProvisionBucket returns the provision bucket for a lock payment order
func (s *IndexerService) getProvisionBucket(ctx context.Context, client types.RPCClient, amount decimal.Decimal, institutionCode [32]byte) (*ent.ProvisionBucket, error) {
	instance, err := contracts.NewPaycrestOrder(OrderConf.PaycrestOrderContractAddress, client.(bind.ContractBackend))
	if err != nil {
		return nil, err
	}

	institution, err := instance.GetSupportedInstitutionName(nil, institutionCode)
	if err != nil {
		return nil, err
	}

	provisionBucket, err := db.Client.ProvisionBucket.
		Query().
		Where(
			provisionbucket.MaxAmountGTE(amount),
			provisionbucket.MinAmountLTE(amount),
			provisionbucket.CurrencyEQ(utils.Byte32ToString(institution.Currency)),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return provisionBucket, nil
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
