package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
	"github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"
)

var conf = config.OrderConfig()

// IndexerService performs blockchain to database extract, transform, load (ETL) operations.
type IndexerService struct {
	db *ent.Client
}

// NewIndexerService creates a new instance of IndexerService.
func NewIndexerService(db *ent.Client) *IndexerService {
	return &IndexerService{
		db: db,
	}
}

// IndexERC20Transfer indexes ERC20 token transfers for a specific address.
func (s *IndexerService) IndexERC20Transfer(ctx context.Context, address common.Address) error {

	// Fetch receive address from db
	receiveAddress, err := s.db.ReceiveAddress.
		Query().
		Where(receiveaddress.AddressEQ(address.Hex())).
		WithPaymentOrder(func(pq *ent.PaymentOrderQuery) { // prefetch payment order, token, and network
			pq.WithToken(func(tq *ent.TokenQuery) {
				tq.WithNetwork()
			})
		}).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch receive address: %w", err)
	}

	paymentOrder := receiveAddress.Edges.PaymentOrder
	token := paymentOrder.Edges.Token

	client, err := ethclient.Dial(token.Edges.Network.RPCEndpoint)
	if err != nil {
		return fmt.Errorf("failed to dial RPC: %w", err)
	}

	// Fetch current block header
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch current block number: %w", err)
	}

	// Number of blocks that will be generated within the receive address valid period.
	// Layer 2 networks have the shortest blocktimes e.g Polygon, Tron etc. are < 5 seconds
	// We assume a blocktime of 2s for the largest number of blocks to scan.
	// number of blocks = receive address validity period in seconds / blocktime
	numOfBlocks := int64(conf.ReceiveAddressValidity.Seconds() / 2)

	var fromBlock *big.Int

	if receiveAddress.LastIndexedBlock > 0 {
		// Continue indexing from last indexed block if the last process failed
		fromBlock = big.NewInt(receiveAddress.LastIndexedBlock + 1)
	} else {
		fromBlock = big.NewInt(header.Number.Int64() - numOfBlocks + 1)
	}

	// Query event logs of the token contract starting from the oldest block
	// within the receive address validity period
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
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
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

	logger.Infof("Indexing transfer logs for %s from Block #%s - #%s", address.Hex(), fromBlock.String(), header.Number.String())

	for {
		// Update the filter parameters
		query.FromBlock = currentBlockNumber
		query.ToBlock = new(big.Int).Add(currentBlockNumber, big.NewInt(int64(500-1)))

		// Check if we have reached the final block number
		if query.ToBlock.Cmp(finalBlockNumber) > 0 {
			break
		}

		// Fetch logs for the current batch
		logs, err := client.FilterLogs(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to fetch logs: %w", err)
		}

		for _, vLog := range logs {
			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				var transferEvent ERC20Transfer

				err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
				if err != nil {
					return fmt.Errorf("failed to unpack Transfer event signature: %w", err)
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
							return fmt.Errorf("failed to update receive address status: %w", err)
						}
					} else if comparisonResult < 0 {
						// Transfer value is less than order amount
						_, err = receiveAddress.
							Update().
							SetStatus(receiveaddress.StatusPartial).
							Save(ctx)
						if err != nil {
							return fmt.Errorf("failed to update receive address status: %w", err)
						}
					}

					// Update the payment order with amount paid
					_, err = paymentOrder.
						Update().
						SetAmountPaid(paymentOrder.AmountPaid.Add(utils.FromSubunit(transferEvent.Value, token.Decimals))).
						Save(ctx)
					if err != nil {
						return fmt.Errorf("failed to record amount paid: %w", err)
					}

					if receiveAddress.Status == receiveaddress.StatusPartial {
						// Refresh the receive address with payment order and compare the amount paid with expected amount,
						receiveAddress, err = s.db.ReceiveAddress.
							Query().
							Where(receiveaddress.AddressEQ(address.String())).
							WithPaymentOrder().
							Only(ctx)
						if err != nil {
							return fmt.Errorf("failed to refresh receive address: %w", err)
						}

						// If amount paid meets or exceeds the expected amount, mark receive address as used
						if paymentOrder.AmountPaid.GreaterThanOrEqual(paymentOrder.Amount) {
							_, err = receiveAddress.
								Update().
								SetStatus(receiveaddress.StatusUsed).
								SetLastUsed(time.Now()).
								Save(ctx)
							if err != nil {
								return fmt.Errorf("failed to update receive address status: %w", err)
							}
						}
					}

					return nil
				}

				// Check if receive address validity period has passed
				timeAgo := time.Now().Add(conf.ReceiveAddressValidity)
				amountNotPaidInFull := (receiveAddress.Status == receiveaddress.StatusPartial || receiveAddress.Status == receiveaddress.StatusUnused)

				if receiveAddress.CreatedAt.Before(timeAgo) && amountNotPaidInFull {
					// Receive address hasn't received full payment after validity period, mark status as expired
					_, err = receiveAddress.
						Update().
						SetStatus(receiveaddress.StatusExpired).
						Save(ctx)
					if err != nil {
						return fmt.Errorf("failed to update receive address status: %w", err)
					}
				}
			}
		}

		// Check if we have fetched all logs
		if len(logs) < currentBlockBatchSize {
			break
		}

		// Update the current block number for the next batch
		currentBlockNumber = new(big.Int).Add(big.NewInt(int64(logs[len(logs)-1].BlockNumber)), big.NewInt(1))

		// Sleep for a short duration between batches to avoid overwhelming the RPC endpoint
		time.Sleep(1 * time.Second)
	}

	return nil
}
