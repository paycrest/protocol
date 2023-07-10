package utils

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

// ToSubUnit converts a decimal amount to the smallest subunit representation.
// It takes the amount and the number of decimal places (decimals) and returns
// the amount in subunits as a *big.Int.
func ToSubUnit(amount decimal.Decimal, decimals uint8) *big.Int {
	amountInSubUnit := big.NewInt(0)

	// Multiply the amount by 10^decimals to convert it to subunits
	return amountInSubUnit.Mul(amount.BigInt(), big.NewInt(int64(math.Pow(10, float64(decimals)))))
}

// FromSubUnit converts an amount in subunits represented as a *big.Int back
// to its decimal representation with the given number of decimal places (decimals).
// It returns the amount as a decimal.Decimal.
func FromSubUnit(amountInSubUnit *big.Int, decimals uint8) decimal.Decimal {
	return decimal.NewFromBigInt(amountInSubUnit, int32(decimals))
}

// FilterLogsInBatches fetches logs in batches using pagination.
func FilterLogsInBatches(ctx context.Context, client *ethclient.Client, filter ethereum.FilterQuery, batchSize int) ([]types.Log, error) {
	var allLogs []types.Log
	currentBatchSize := batchSize
	currentBlockNumber := filter.FromBlock

	for {
		// Update the filter parameters
		filter.FromBlock = currentBlockNumber
		filter.ToBlock = new(big.Int).Add(currentBlockNumber, big.NewInt(int64(batchSize-1)))

		// Fetch logs for the current batch
		logs, err := client.FilterLogs(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch logs: %v", err)
		}

		// Append the logs to the result
		allLogs = append(allLogs, logs...)

		// Check if we have fetched all logs
		if len(logs) < currentBatchSize {
			break
		}

		// Update the current block number for the next batch
		currentBlockNumber = new(big.Int).Add(big.NewInt(int64(logs[len(logs)-1].BlockNumber)), big.NewInt(1))

		// Sleep for a short duration between batches to avoid overwhelming the RPC endpoint
		time.Sleep(1 * time.Second)
	}

	return allLogs, nil
}
