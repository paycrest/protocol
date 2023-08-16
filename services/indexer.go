package services

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/paycrest/paycrest-protocol/config"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"
)

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
	var conf = config.OrderConfig()

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
	numOfBlocks := int64(conf.ReceiveAddressValidity.Seconds() / 2)

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
	// TODO: explain why batchsize of 500 was chosen
	currentBlockBatchSize := 500
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
		"\nIndexing transfer logs for %s from Block #%s - #%s\n\n",
		receiveAddress.Address,
		fromBlockBig.String(),
		header.Number.String(),
	)

	for {
		// Update the filter parameters
		query.FromBlock = currentBlockNumber
		query.ToBlock = big.NewInt(currentBlockNumber.Int64() + int64(currentBlockBatchSize-1))

		// Check if we have reached the final block number
		if query.ToBlock.Cmp(finalBlockNumber) > 0 {
			break
		}

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
						SetValidUntil(time.Now().Add(conf.ReceiveAddressValidity)).
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

		// Update last indexed block number
		_, err = receiveAddress.
			Update().
			SetLastIndexedBlock(query.ToBlock.Int64()).
			Save(ctx)
		if err != nil {
			return false, fmt.Errorf("failed to update receive address last indexed block: %w", err)
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
		// time.Sleep(2 * time.Minute) // add 2 minutes delay between each indexing operation

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
