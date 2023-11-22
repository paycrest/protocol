package test

import (
	"context"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/shopspring/decimal"
)

// CreateTestUser creates a test user with default or custom values
func CreateTestUser(overrides map[string]interface{}) (*ent.User, error) {

	// Default payload
	payload := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "johndoe@test.com",
		"password":  "password",
		"scope":     []string{"sender"},
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	var scopes []string

	if scopes, ok := payload["scope"].([]string); ok {

		// Check if slice only contains "provider"
		isProviderOnly := true
		for _, s := range scopes {
			if s != "provider" {
				isProviderOnly = false
				break
			}
		}

		if isProviderOnly {
			scopes = []string{"sender", "provider"}
		} else {
			scopes = payload["scope"].([]string)
		}
	}

	// Create user
	user, err := db.Client.User.
		Create().
		SetFirstName(payload["firstName"].(string)).
		SetLastName(payload["lastName"].(string)).
		SetEmail(payload["email"].(string)).
		SetPassword(payload["password"].(string)).
		SetScopes(scopes).
		Save(context.Background())

	return user, err
}

// CreateTestToken creates a test token with default or custom values
func CreateTestToken(client types.RPCClient, overrides map[string]interface{}) (*ent.Token, error) {

	// Deploy ERC20 token contract
	tokenAddress, err := DeployERC20Contract(client)
	if err != nil {
		return nil, err
	}

	// Default payload
	payload := map[string]interface{}{
		"symbol":           "TST",
		"contract_address": tokenAddress.Hex(),
		"decimals":         18,
		"networkRPC":       "http://localhost:8545",
		"is_enabled":       true,
		"identifier":       "polygon-mumbai" + uuid.New().String(),
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create Network
	network, err := db.Client.Network.
		Create().
		SetIdentifier(payload["identifier"].(string)). // randomize the identifier to avoid conflicts
		SetChainID(1337).
		SetRPCEndpoint(payload["networkRPC"].(string)).
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

	return token, err
}

// CreateTestLockPaymentOrder creates a test LockPaymentOrder with default or custom values
func CreateTestLockPaymentOrder(overrides map[string]interface{}) (*ent.LockPaymentOrder, error) {

	// Default payload
	payload := map[string]interface{}{
		"order_id":           "order-123",
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
	backend, _ := NewSimulatedBlockchain()
	token, _ := CreateTestToken(backend, nil)

	// Create LockPaymentOrder
	order, err := db.Client.LockPaymentOrder.
		Create().
		SetOrderID(payload["order_id"].(string)).
		SetAmount(decimal.NewFromFloat(payload["amount"].(float64))).
		SetRate(decimal.NewFromFloat(payload["rate"].(float64))).
		SetStatus(lockpaymentorder.Status(payload["status"].(string))).
		SetLabel(payload["label"].(string)).
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
		"tx_receipt_image":  "https://picsum.photos/200",
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
		SetTxReceiptImage(payload["tx_receipt_image"].(string)).
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
		"fee_address":        "0x123...",
		"refund_address":     "0x123...",
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
func CreateTestProviderProfile(overrides map[string]interface{}, user *ent.User, currency *ent.FiatCurrency) (*ent.ProviderProfile, error) {

	// Default payload
	payload := map[string]interface{}{
		"user_id":         uuid.New(),
		"trading_name":    "Elon Musk Trading Co.",
		"host_identifier": "https://example.com/hook",
		"provision_mode":  "auto",
		"is_partner":      false,
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
		SetUser(user).
		SetCurrency(currency).
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
