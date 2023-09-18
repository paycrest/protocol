package utils

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/paycrest/paycrest-protocol/types"
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
func StringToByte32(s string) [32]byte {
	var result [32]byte

	// Convert the input string to bytes
	inputBytes := []byte(s)

	// Copy the input bytes into the result array, limiting to 32 bytes
	copy(result[:], inputBytes)

	return result
}

// Byte32ToString converts [32]byte to string
func Byte32ToString(b [32]byte) string {

	// Copy byte array into slice
	buf := make([]byte, 32)
	copy(buf, b[:])

	// Truncate trailing zeros
	buf = bytes.TrimRight(buf, "\x00")

	return string(buf)
}

// BigMin returns the minimum value between two big numbers
func BigMin(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return x
	}
	return y
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

// PersonalSign is an equivalent of ethers.personal_sign for signing ethereum messages
// Ref: https://github.com/etaaa/Golang-Ethereum-Personal-Sign/blob/main/main.go
func PersonalSign(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(fullMessage))
	signatureBytes, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	signatureBytes[64] += 27
	return signatureBytes, nil
}

// EIP1559GasPrice returns the maxFeePerGas and maxPriorityFeePerGas for EIP-1559
func EIP1559GasPrice(ctx context.Context, client types.RPCClient) (maxFeePerGas, maxPriorityFeePerGas *big.Int) {
	tip, _ := client.SuggestGasTipCap(ctx)
	latestHeader, _ := client.HeaderByNumber(ctx, nil)

	buffer := new(big.Int).Mul(tip, big.NewInt(13)).Div(tip, big.NewInt(100))
	maxPriorityFeePerGas = new(big.Int).Add(tip, buffer)

	if latestHeader.BaseFee != nil {
		maxFeePerGas = new(big.Int).
			Mul(latestHeader.BaseFee, big.NewInt(2)).
			Add(latestHeader.BaseFee, maxPriorityFeePerGas)
	} else {
		maxFeePerGas = maxPriorityFeePerGas
	}

	return maxFeePerGas, maxPriorityFeePerGas
}
