package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	cryptoUtils "github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/shopspring/decimal"
)

// ToSubunit converts a decimal amount to the smallest subunit representation.
// It takes the amount and the number of decimal places (decimals) and returns
// the amount in subunits as a *big.Int.
func ToSubunit(amount decimal.Decimal, decimals int8) *big.Int {
	// Compute the multiplier: 10^decimals
	multiplier := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))

	// Multiply the amount by the multiplier to convert it to subunits
	subunitInDecimal := amount.Mul(multiplier)

	// Create a new big.Int from the string representation of the subunit amount
	subunit := new(big.Int)
	subunit.SetString(subunitInDecimal.String(), 10)

	return subunit
}

// FromSubunit converts an amount in subunits represented as a *big.Int back
// to its decimal representation with the given number of decimal places (decimals).
// It returns the amount as a decimal.Decimal.
func FromSubunit(amountInSubunit *big.Int, decimals int8) decimal.Decimal {
	// Compute the divisor: 10^decimals
	divisor := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals))).BigFloat()

	// Create a new big.Float with the desired precision and rounding mode
	f := new(big.Float).SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)

	// Create a new big.Float for the subunit amount with the desired precision and rounding mode
	fSubunit := new(big.Float).SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fSubunit.SetMode(big.ToNearestEven)

	// Divide the subunit amount by the divisor and convert it to a float64
	result, _ := f.Quo(fSubunit.SetInt(amountInSubunit), divisor).Float64()

	return decimal.NewFromFloat(result)
}

// StringToByte32 converts string to [32]byte
func StringTo32Byte(s string) [32]byte {

	buf := []byte(s)

	// Pad or truncate
	if len(buf) < 32 {
		padded := make([]byte, 32)
		copy(padded, buf)
		buf = padded
	} else if len(buf) > 32 {
		buf = buf[:32]
	}

	// Hex encode and decode
	hexStr := hex.EncodeToString(buf)
	bytesBytes, _ := hex.DecodeString(hexStr)

	var result [32]byte
	copy(result[:], bytesBytes)

	return result
}

// Byte32ToString converts [32]byte to string
func Byte32ToString(b [32]byte) string {

	// Copy byte array into slice
	buf := make([]byte, 32)
	copy(buf, b[:])

	// Truncate trailing zeros
	buf = bytes.TrimRight(buf, "\x00")

	// Hex encode and decode
	hexStr := hex.EncodeToString(buf)
	strBytes, _ := hex.DecodeString(hexStr)

	return string(strBytes)
}

// GetMasterAccount returns the master account address and private key.
func GetMasterAccount() (*common.Address, *ecdsa.PrivateKey, error) {
	fromAddress, privateKeyHex, err := cryptoUtils.GenerateAccountFromIndex(0)
	if err != nil {
		return nil, nil, err
	}

	address := common.HexToAddress(fromAddress)

	privateKeyBytes, err := hexutil.Decode(privateKeyHex)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	return &address, privateKey, nil
}
