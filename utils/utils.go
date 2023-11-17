package utils

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	tokenUtils "github.com/paycrest/protocol/utils/token"
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
// func Byte32ToString(b [32]byte) string {

// 	// Copy byte array into slice
// 	buf := make([]byte, 32)
// 	copy(buf, b[:])

// 	// Truncate trailing zeros
// 	buf = bytes.TrimRight(buf, "\x00")

//		return string(buf)
//	}
func Byte32ToString(b [32]byte) string {

	// Find first null index if any
	nullIndex := -1
	for i, x := range b {
		if x == 0 {
			nullIndex = i
			break
		}
	}

	// Slice at first null or return full 32 bytes
	if nullIndex >= 0 {
		return string(b[:nullIndex])
	} else {
		return string(b[:])
	}
}

// BigMin returns the minimum value between two big numbers
func BigMin(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return x
	}
	return y
}

// PersonalSign is an equivalent of ethers.personal_sign for signing ethereum messages
// Ref: https://github.com/etaaa/Golang-Ethereum-Personal-Sign/blob/main/main.go
func PersonalSign(message string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	fullMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(fullMessage))
	signatureBytes, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	signatureBytes[64] += 27
	return signatureBytes, nil
}

// Difference returns the elements in `a` that aren't in `b`.
func Difference(a, b []string) []string {
	setB := make(map[string]struct{})
	for _, x := range b {
		setB[x] = struct{}{}
	}

	var diff []string
	for _, x := range a {
		if _, found := setB[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// ContainsString returns true if the slice contains the given string
func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Median returns the median value of a decimal slice
func Median(data []decimal.Decimal) decimal.Decimal {
	l := len(data)
	if l == 0 {
		return decimal.Zero
	}

	// Sort data in ascending order
	sort.Slice(data, func(i, j int) bool {
		return data[i].LessThan(data[j])
	})

	middle := l / 2
	result := data[middle]

	// Handle even length slices
	if l%2 == 0 {
		result = result.Add(data[middle-1])
		result = result.Div(decimal.NewFromInt(2))
	}

	return result
}

// SendPaymentOrderWebhook notifies a sender when the status of a payment order changes
func SendPaymentOrderWebhook(ctx context.Context, paymentOrder *ent.PaymentOrder) error {
	profile := paymentOrder.Edges.SenderProfile

	// If webhook URL is empty, return
	if profile.WebhookURL == "" {
		return nil
	}

	// Determine the event
	var event string

	switch paymentOrder.Status {
	case paymentorder.StatusPending:
		event = "payment_order.pending"
	case paymentorder.StatusReverted:
		event = "payment_order.reverted"
	case paymentorder.StatusExpired:
		event = "payment_order.expired"
	case paymentorder.StatusSettled:
		event = "payment_order.settled"
	case paymentorder.StatusRefunded:
		event = "payment_order.refunded"
	default:
		return nil
	}

	// Fetch the recipient
	recipient, err := paymentOrder.QueryRecipient().Only(ctx)
	if err != nil {
		return err
	}

	// Fetch the token
	token, err := paymentOrder.
		QueryToken().
		WithNetwork().
		Only(ctx)
	if err != nil {
		return err
	}

	// Create the payload
	payloadStruct := types.PaymentOrderWebhookPayload{
		Event: event,
		Data: types.PaymentOrderWebhookData{
			ID:         paymentOrder.ID,
			Amount:     paymentOrder.Amount,
			AmountPaid: paymentOrder.AmountPaid,
			Rate:       paymentOrder.Rate,
			Network:    token.Edges.Network.Identifier,
			Label:      paymentOrder.Label,
			SenderID:   profile.ID,
			Recipient: types.PaymentOrderRecipient{
				Institution:       recipient.Institution,
				AccountIdentifier: recipient.AccountIdentifier,
				AccountName:       recipient.AccountName,
				ProviderID:        recipient.ProviderID,
				Memo:              recipient.Memo,
			},
			UpdatedAt: paymentOrder.UpdatedAt,
			CreatedAt: paymentOrder.CreatedAt,
			TxHash:    paymentOrder.TxHash,
			Status:    paymentOrder.Status,
		},
	}

	payload := StructToMap(payloadStruct)

	// Compute HMAC signature
	apiKey, err := profile.QueryAPIKey().Only(ctx)
	if err != nil {
		return err
	}

	decodedSecret, err := base64.StdEncoding.DecodeString(apiKey.Secret)
	if err != nil {
		return err
	}

	decryptedSecret, err := cryptoUtils.DecryptPlain(decodedSecret)
	if err != nil {
		return err
	}

	signature := tokenUtils.GenerateHMACSignature(payload, string(decryptedSecret))

	// Send the webhook
	_, err = MakeJSONRequest(
		ctx,
		"POST",
		profile.WebhookURL,
		payload,
		map[string]string{
			"X-Paycrest-Signature": signature,
		},
	)
	if err != nil {
		// Log retry attempt
		_, err := storage.Client.WebhookRetryAttempt.
			Create().
			SetAttemptNumber(1).
			SetNextRetryTime(time.Now().Add(2 * time.Minute)).
			SetPayload(payload).
			SetSignature(signature).
			SetWebhookURL(profile.WebhookURL).
			SetStatus("failed").
			Save(ctx)
		return err
	}

	return nil
}

// StructToMap converts a struct to a map[string]interface{}
func StructToMap(input interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Use reflection to iterate over the struct fields
	valueOf := reflect.ValueOf(input)
	typeOf := valueOf.Type()

	for i := 0; i < valueOf.NumField(); i++ {
		field := valueOf.Field(i)
		fieldName := typeOf.Field(i).Name

		// Convert the field value to interface{}
		result[fieldName] = field.Interface()
	}

	return result
}

// StringSliceContains returns true if the slice contains the given string
func StringSliceContains(slice []string, str string) bool {
  for _, s := range slice {
    if s == str {
      return true
    }
  }
  return false
}