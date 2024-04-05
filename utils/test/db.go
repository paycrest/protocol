package test

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
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
	// Deploy ERC20 token contract
	deployedTokenAddress, err := DeployERC20Contract(client)
	if err != nil {
		return nil, err
	}

	// Default payload
	payload := map[string]interface{}{
		"symbol":           "TST",
		"contract_address": deployedTokenAddress.Hex(),
		"decimals":         18,
		"networkRPC":       "ws://localhost:8545",
		"is_enabled":       true,
		"identifier":       "localhost" + uuid.New().String(),
		"chainID":          int64(1337),
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create Network
	network, err := db.Client.Network.
		Create().
		SetIdentifier(payload["identifier"].(string)). // randomize the identifier to avoid conflicts
		SetChainID(payload["chainID"].(int64)).
		SetRPCEndpoint(payload["networkRPC"].(string)).
		SetFee(decimal.NewFromFloat(0.1)).
		SetIsTestnet(true).
		Save(context.Background())
	if err != nil {
		return nil, err
	}

	// Create token
	token, err := db.Client.Token.
		Create().
		SetSymbol(payload["symbol"].(string)).
		SetContractAddress(payload["contract_address"].(string)).
		SetDecimals(int8(payload["decimals"].(int))).
		SetNetwork(network).
		SetIsEnabled(payload["is_enabled"].(bool)).
		Save(context.Background())
	if err != nil {
		return nil, err
	}

	token, err = db.Client.Token.
		Query().
		Where(entToken.IDEQ(token.ID)).
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
		"label":              "thisisatestlabel",
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
	backend, _ := SetUpTestBlockchain()
	token, _ := CreateERC20Token(backend, nil)

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
		"fee_per_token_unit": 0.0,
		"webhook_url":        "https://example.com/hook",
		"domain_whitelist":   []string{"example.com"},
		"fee_address":        "0x1234567890123456789012345678901234567890",
		"refund_address":     "0x0987654321098765432109876543210987654321",
		"user_id":            nil,
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create SenderProfile
	profile, err := db.Client.SenderProfile.
		Create().
		SetWebhookURL(payload["webhook_url"].(string)).
		SetDomainWhitelist(payload["domain_whitelist"].([]string)).
		SetFeePerTokenUnit(decimal.NewFromFloat(payload["fee_per_token_unit"].(float64))).
		SetFeeAddress(payload["fee_address"].(string)).
		SetRefundAddress(payload["refund_address"].(string)).
		SetUserID(payload["user_id"].(uuid.UUID)).
		Save(context.Background())

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

	currency, err := db.Client.FiatCurrency.
		Create().
		SetCode(payload["code"].(string)).
		SetShortName(payload["short_name"].(string)).
		SetDecimals(payload["decimals"].(int)).
		SetSymbol(payload["symbol"].(string)).
		SetName(payload["name"].(string)).
		SetMarketRate(decimal.NewFromFloat(payload["market_rate"].(float64))).
		SetIsEnabled(true).
		Save(context.Background())

	return currency, err

}
