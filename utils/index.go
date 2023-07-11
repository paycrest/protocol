package utils

import (
	"math"
	"math/big"

	"github.com/shopspring/decimal"
)

// ToSubUnit converts a decimal amount to the smallest subunit representation.
// It takes the amount and the number of decimal places (decimals) and returns
// the amount in subunits as a *big.Int.
func ToSubUnit(amount decimal.Decimal, decimals int8) *big.Int {
	amountInSubUnit := big.NewInt(0)

	// Multiply the amount by 10^decimals to convert it to subunits
	return amountInSubUnit.Mul(amount.BigInt(), big.NewInt(int64(math.Pow(10, float64(decimals)))))
}

// FromSubUnit converts an amount in subunits represented as a *big.Int back
// to its decimal representation with the given number of decimal places (decimals).
// It returns the amount as a decimal.Decimal.
func FromSubUnit(amountInSubUnit *big.Int, decimals int8) decimal.Decimal {
	return decimal.NewFromBigInt(amountInSubUnit, int32(decimals))
}
