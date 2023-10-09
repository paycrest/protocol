package utils

import (
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {

	t.Run("ToSubunit", func(t *testing.T) {
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

		for _, tc := range testCases {
			actualVal := ToSubunit(tc.amount, tc.decimals)
			assert.Equal(t, tc.expectVal, actualVal)
		}
	})

	t.Run("FromSubunit", func(t *testing.T) {
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

		for _, tc := range testCases {
			actualVal := FromSubunit(tc.amountInSubunit, tc.decimals)
			assert.Equal(t, tc.expectVal, actualVal)
		}
	})

	t.Run("TestMedian", func(t *testing.T) {
		data := []decimal.Decimal{
			decimal.NewFromInt(9),
			decimal.NewFromInt(1),
			decimal.NewFromInt(5),
			decimal.NewFromInt(6),
			decimal.NewFromInt(2),
			decimal.NewFromInt(1),
			decimal.NewFromInt(3),
			decimal.NewFromInt(1),
			decimal.NewFromInt(1),
			decimal.NewFromInt(2),
		}

		median := Median(data)

		assert := assert.New(t)
		assert.True(median.Equal(decimal.NewFromInt(2)), "Median calculation is incorrect")
	})
}
