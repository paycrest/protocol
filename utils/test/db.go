package test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/institution"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/senderordertoken"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/paycrest/protocol/ent/token"
	entToken "github.com/paycrest/protocol/ent/token"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/shopspring/decimal"
)

// CreateTestUser creates a test user with default or custom values
func CreateTestUser(overrides map[string]interface{}) (*ent.User, error) {

	// Default payload
	payload := map[string]interface{}{
		"firstName":       "John",
		"lastName":        "Doe",
		"email":           "johndoe@test.com",
		"password":        "password",
		"scope":           "sender",
		"isEmailVerified": false,
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create user
	user, err := db.Client.User.
		Create().
		SetFirstName(payload["firstName"].(string)).
		SetLastName(payload["lastName"].(string)).
		SetEmail(strings.ToLower(payload["email"].(string))).
		SetPassword(payload["password"].(string)).
		SetScope(payload["scope"].(string)).
		SetIsEmailVerified(payload["isEmailVerified"].(bool)).
		Save(context.Background())

	return user, err
}

// CreateERC20Token creates a test token with default or custom values
func CreateERC20Token(client types.RPCClient, overrides map[string]interface{}) (*ent.Token, error) {

	// Default payload
	payload := map[string]interface{}{
		"symbol":         "TST",
		"decimals":       6,
		"networkRPC":     "http://localhost:8545",
		"is_enabled":     true,
		"identifier":     "localhost",
		"chainID":        int64(1337),
		"deployContract": true,
	}

	var contractAddress string

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	if payload["deployContract"].(bool) {
		// Deploy ERC20 token contract
		deployedTokenAddress, err := DeployERC20Contract(client)
		if err != nil {
			return nil, err
		}
		contractAddress = deployedTokenAddress.Hex()
	} else {
		contractAddress = "0xd4E96eF8eee8678dBFf4d535E033Ed1a4F7605b7"
	}

	// Create Network
	networkId, err := db.Client.Network.
		Create().
		SetIdentifier(payload["identifier"].(string)).
		SetChainID(payload["chainID"].(int64)).
		SetRPCEndpoint(payload["networkRPC"].(string)).
		SetFee(decimal.NewFromFloat(0.1)).
		SetIsTestnet(true).
		OnConflict().
		UpdateNewValues().
		UpdateRPCEndpoint().
		UpdateChainID().
		UpdateIdentifier().
		ID(context.Background())

	if err != nil {
		return nil, fmt.Errorf("CreateERC20Token.networkId: %w", err)
	}
	// Create token
	tokenId := db.Client.Token.
		Create().
		SetSymbol(payload["symbol"].(string)).
		SetContractAddress(contractAddress).
		SetDecimals(int8(payload["decimals"].(int))).
		SetNetworkID(networkId).
		SetIsEnabled(payload["is_enabled"].(bool)).
		OnConflict().
		// Use the new values that were set on create.
		UpdateNewValues().
		IDX(context.Background())

	token, err := db.Client.Token.
		Query().
		Where(entToken.IDEQ(tokenId)).
		WithNetwork().
		Only(context.Background())

	return token, err
}

// CreateERC20Token creates a test token with default or custom values
func CreateTRC20Token(client types.RPCClient, overrides map[string]interface{}) (*ent.Token, error) {

	// Default payload
	payload := map[string]interface{}{
		"symbol":     "TRON_ST",
		"decimals":   6,
		"networkRPC": "ws://localhost:8544",
		"is_enabled": true,
		"identifier": "tron",
		"chainID":    int64(13378),
	}

	contractAddress := "TFRKiHrHCeSyWL67CEwydFvUMYJ6CbYYX6"

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create Network
	networkId, err := db.Client.Network.
		Create().
		SetIdentifier(payload["identifier"].(string)).
		SetChainID(payload["chainID"].(int64)).
		SetRPCEndpoint(payload["networkRPC"].(string)).
		SetFee(decimal.NewFromFloat(0.1)).
		SetIsTestnet(true).
		OnConflict().
		UpdateNewValues().
		ID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("CreateERC20Token.networkId: %w", err)
	}

	// Create token
	tokenId := db.Client.Token.
		Create().
		SetSymbol(payload["symbol"].(string)).
		SetContractAddress(contractAddress).
		SetDecimals(int8(payload["decimals"].(int))).
		SetNetworkID(networkId).
		SetIsEnabled(payload["is_enabled"].(bool)).
		OnConflict().
		UpdateNewValues().
		IDX(context.Background())

	token, err := db.Client.Token.
		Query().
		Where(entToken.IDEQ(tokenId)).
		WithNetwork().
		Only(context.Background())

	return token, err
}

// CreateTestLockPaymentOrder creates a test LockPaymentOrder with default or custom values
func CreateTestLockPaymentOrder(overrides map[string]interface{}) (*ent.LockPaymentOrder, error) {

	// Default payload
	payload := map[string]interface{}{
		"gateway_id":         "order-123",
		"amount":             100.50,
		"rate":               750.0,
		"status":             "pending",
		"block_number":       12345,
		"institution":        "Test Bank",
		"account_identifier": "1234567890",
		"account_name":       "Test Account",
	}

	// Create provider profile
	var providerProfile *ent.ProviderProfile
	if overrides["provider"] == nil {
		providerProfile = nil
	} else {
		providerProfile = overrides["provider"].(*ent.ProviderProfile)
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create test token
	backend, _ := SetUpTestBlockchain(nil)
	token, err := CreateERC20Token(backend, map[string]interface{}{
		"deployContract": false,
	})
	if err != nil {
		return nil, err
	}

	// Create LockPaymentOrder
	order, err := db.Client.LockPaymentOrder.
		Create().
		SetGatewayID(payload["gateway_id"].(string)).
		SetAmount(decimal.NewFromFloat(payload["amount"].(float64))).
		SetRate(decimal.NewFromFloat(payload["rate"].(float64))).
		SetStatus(lockpaymentorder.Status(payload["status"].(string))).
		SetOrderPercent(decimal.NewFromFloat(100.0)).
		SetBlockNumber(int64(payload["block_number"].(int))).
		SetInstitution(payload["institution"].(string)).
		SetAccountIdentifier(payload["account_identifier"].(string)).
		SetAccountName(payload["account_name"].(string)).
		SetTokenID(token.ID).
		SetProvider(providerProfile).
		Save(context.Background())

	return order, err
}

// CreateTestPaymentOrder creates a test PaymentOrder with default or custom values for sender
func CreateTestPaymentOrder(client types.RPCClient, token *ent.Token, overrides map[string]interface{}) (*ent.PaymentOrder, error) {
	// Default payload
	payload := map[string]interface{}{
		"amount":             100.50,
		"rate":               750.0,
		"status":             "pending",
		"fee_per_token_unit": 0.0,
		"fee_address":        "0x1234567890123456789012345678901234567890",
		"return_address":     "0x0987654321098765432109876543210987654321",
		"institution":        "Test Bank",
		"account_identifier": "1234567890",
		"account_name":       "Test Account",
		"memo":               "Shola Kehinde - rent for May 2021",
		"providerId":         "",
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create smart wallet
	receiveAddress, err := CreateSmartAddress(
		context.Background(), client)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)

	// Create payment order
	paymentOrder, err := db.Client.PaymentOrder.
		Create().
		SetSenderProfile(overrides["sender"].(*ent.SenderProfile)).
		SetAmount(decimal.NewFromFloat(payload["amount"].(float64))).
		SetAmountPaid(decimal.NewFromInt(0)).
		SetAmountReturned(decimal.NewFromInt(0)).
		SetPercentSettled(decimal.NewFromInt(0)).
		SetNetworkFee(token.Edges.Network.Fee).
		SetProtocolFee(decimal.NewFromFloat(payload["amount"].(float64)).Mul(decimal.NewFromFloat(0))).
		SetSenderFee(decimal.NewFromFloat(payload["fee_per_token_unit"].(float64)).Mul(decimal.NewFromFloat(payload["amount"].(float64))).Div(decimal.NewFromFloat(payload["rate"].(float64))).Round(int32(token.Decimals))).
		SetToken(token).
		SetRate(decimal.NewFromFloat(payload["rate"].(float64))).
		SetReceiveAddress(receiveAddress).
		SetReceiveAddressText(receiveAddress.Address).
		SetFeePerTokenUnit(decimal.NewFromFloat(payload["fee_per_token_unit"].(float64))).
		SetFeeAddress(payload["fee_address"].(string)).
		SetReturnAddress(payload["return_address"].(string)).
		SetStatus(paymentorder.Status(payload["status"].(string))).
		Save(context.Background())
	if err != nil {
		return nil, err
	}

	// Create payment order recipient
	_, err = db.Client.PaymentOrderRecipient.
		Create().
		SetInstitution(payload["institution"].(string)).
		SetAccountIdentifier(payload["account_identifier"].(string)).
		SetAccountName(payload["account_name"].(string)).
		SetProviderID(payload["providerId"].(string)).
		SetMemo(payload["memo"].(string)).
		SetPaymentOrder(paymentOrder).
		Save(context.Background())

	return paymentOrder, err
}

// CreateTestLockOrderFulfillment creates a test LockOrderFulfillment with defaults or custom values
func CreateTestLockOrderFulfillment(overrides map[string]interface{}) (*ent.LockOrderFulfillment, error) {

	// Default payload
	payload := map[string]interface{}{
		"tx_id":             "0x123...",
		"validation_errors": []string{},
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create lock order
	order, _ := CreateTestLockPaymentOrder(nil)

	// Create LockOrderFulfillment
	fulfillment, err := db.Client.LockOrderFulfillment.
		Create().
		SetTxID(payload["tx_id"].(string)).
		SetOrderID(order.ID).
		Save(context.Background())

	return fulfillment, err
}

// CreateTestSenderProfile creates a test SenderProfile with defaults or custom values
func CreateTestSenderProfile(overrides map[string]interface{}) (*ent.SenderProfile, error) {

	// Default payload
	payload := map[string]interface{}{
		"fee_per_token_unit": "0.0",
		"webhook_url":        "https://example.com/hook",
		"domain_whitelist":   []string{"example.com"},
		"fee_address":        "0x1234567890123456789012345678901234567890",
		"refund_address":     "0x0987654321098765432109876543210987654321",
		"user_id":            nil,
		"token":              "TST",
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	_token, err := db.Client.Token.
		Query().
		Where(
			token.SymbolEQ(payload["token"].(string)),
		).
		Only(context.Background())
	if err != nil {
		return nil, err
	}

	feePerTokenUnit, _ := decimal.NewFromString(payload["fee_per_token_unit"].(string))

	// Create SenderProfile
	profile, err := db.Client.SenderProfile.
		Create().
		SetWebhookURL(payload["webhook_url"].(string)).
		SetDomainWhitelist(payload["domain_whitelist"].([]string)).
		SetUserID(payload["user_id"].(uuid.UUID)).
		Save(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = db.Client.SenderOrderToken.
		Query().
		Where(
			senderordertoken.And(
				senderordertoken.HasTokenWith(token.IDEQ(_token.ID)),
				senderordertoken.HasSenderWith(senderprofile.IDEQ(profile.ID)),
			),
		).Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			_, err := db.Client.SenderOrderToken.
				Create().
				SetSenderID(profile.ID).
				SetTokenID(_token.ID).
				SetRefundAddress(payload["refund_address"].(string)).
				SetFeePerTokenUnit(feePerTokenUnit).
				SetFeeAddress(payload["fee_address"].(string)).
				Save(context.Background())
			if err != nil {
				return nil, fmt.Errorf("CreateTestSenderProfile: %w", err)
			}
			return profile, nil
		} else {
			return nil, fmt.Errorf("CreateTestSenderProfile: %w", err)
		}
	}
	return profile, err
}

// CreateTestProviderProfile creates a test ProviderProfile with defaults or custom values
func CreateTestProviderProfile(overrides map[string]interface{}) (*ent.ProviderProfile, error) {

	// Default payload
	payload := map[string]interface{}{
		"user_id":         uuid.New(),
		"trading_name":    "Elon Musk Trading Co.",
		"currency_id":     uuid.New(),
		"host_identifier": "https://example.com",
		"provision_mode":  "auto",
		"is_partner":      false,
		"visibility_mode": "public",
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create ProviderProfile
	profile, err := db.Client.ProviderProfile.
		Create().
		SetID(payload["user_id"].(uuid.UUID).String()).
		SetTradingName(payload["trading_name"].(string)).
		SetHostIdentifier(payload["host_identifier"].(string)).
		SetProvisionMode(providerprofile.ProvisionMode(payload["provision_mode"].(string))).
		SetIsPartner(payload["is_partner"].(bool)).
		SetUserID(payload["user_id"].(uuid.UUID)).
		SetCurrencyID(payload["currency_id"].(uuid.UUID)).
		SetVisibilityMode(providerprofile.VisibilityMode(payload["visibility_mode"].(string))).
		Save(context.Background())

	return profile, err
}

// CreateTestFiatCurrency creates a test FiatCurrency with defaults or custom values
func CreateTestFiatCurrency(overrides map[string]interface{}) (*ent.FiatCurrency, error) {

	// Default payload.
	payload := map[string]interface{}{
		"code":        "NGN",
		"short_name":  "Naira",
		"decimals":    2,
		"symbol":      "₦",
		"name":        "Nigerian Naira",
		"market_rate": 950.0,
	}

	// Apply overrides.
	for key, value := range overrides {
		payload[key] = value
	}

	institutions, err := db.Client.Institution.CreateBulk(
		db.Client.Institution.
			Create().
			SetName("Kuda Microfinance Bank").
			SetCode("KUDANGN").
			SetType(institution.TypeMobileMoney),
		db.Client.Institution.
			Create().
			SetName("FirstBank Bank").
			SetCode("FBNNGN"),
	).Save(context.Background())

	if err != nil {
		return nil, err
	}

	currency, err := db.Client.FiatCurrency.
		Create().
		SetCode(payload["code"].(string)).
		SetShortName(payload["short_name"].(string)).
		SetDecimals(payload["decimals"].(int)).
		SetSymbol(payload["symbol"].(string)).
		SetName(payload["name"].(string)).
		SetMarketRate(decimal.NewFromFloat(payload["market_rate"].(float64))).
		SetIsEnabled(true).
		AddInstitutions(institutions...).
		Save(context.Background())

	return currency, err

}

// CreateEnvFile creates a new file with Key=Value format.
func CreateEnvFile(filePath string, data map[string]string) (string, error) {
	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Iterate over the map entries and write each key-value pair to the file
	for key, value := range data {
		_, err := writer.WriteString(fmt.Sprintf("%s='%s'\n", key, value))
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}
