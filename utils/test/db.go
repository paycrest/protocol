package test

import (
	"context"
	"strings"

	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/types"
)

// CreateTestUser creates a test user with default or custom values
func CreateTestUser(overrides map[string]string) (*ent.User, error) {

	// Default payload
	payload := map[string]string{
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "johndoe@test.com",
		"password":  "password",
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create user
	user, err := db.Client.User.
		Create().
		SetFirstName(payload["firstName"]).
		SetLastName(payload["lastName"]).
		SetEmail(strings.ToLower(payload["email"])).
		SetPassword(payload["password"]).
		Save(context.Background())

	return user, err
}

// CreateTestToken creates a test token with default or custom values
func CreateTestToken(client types.RPCClient, overrides map[string]string) (*ent.Token, error) {

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
		"is_enabled":       true,
	}

	// Apply overrides
	for key, value := range overrides {
		payload[key] = value
	}

	// Create Network
	network, err := db.Client.Network.
		Create().
		SetIdentifier("polygon-mumbai").
		SetChainID(1337).
		SetRPCEndpoint("http://localhost:8545").
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
