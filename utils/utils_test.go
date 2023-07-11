package utils_test

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/paycrest/paycrest-protocol/utils"
)

func TestToSubunit(t *testing.T) {
	// Test cases
	testCases := []struct {
		amount    decimal.Decimal
		decimals  int8
		expectVal *big.Int
	}{
		{
			amount:    decimal.NewFromFloat(1.23),
			decimals:  2,
			expectVal: big.NewInt(123),
		},
		{
			amount:    decimal.NewFromFloat(0.001),
			decimals:  8,
			expectVal: big.NewInt(100000),
		},
		{
			amount:    decimal.NewFromFloat(0.005),
			decimals:  18,
			expectVal: big.NewInt(5000000000000000),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		actualVal := utils.ToSubunit(tc.amount, tc.decimals)
		assert.Equal(t, tc.expectVal, actualVal)
	}
}

func TestFromSubunit(t *testing.T) {
	// Test cases
	testCases := []struct {
		amountInSubunit *big.Int
		decimals        int8
		expectVal       decimal.Decimal
	}{
		{
			amountInSubunit: big.NewInt(123),
			decimals:        2,
			expectVal:       decimal.NewFromFloat(1.23),
		},
		{
			amountInSubunit: big.NewInt(1),
			decimals:        8,
			expectVal:       decimal.NewFromFloat(0.00000001),
		},
		{
			amountInSubunit: big.NewInt(5000000000000000),
			decimals:        18,
			expectVal:       decimal.NewFromFloat(0.005),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		actualVal := utils.FromSubunit(tc.amountInSubunit, tc.decimals)
		assert.Equal(t, tc.expectVal, actualVal)
	}
}
