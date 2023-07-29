package test

import (
	"context"
	"strings"

	"github.com/paycrest/paycrest-protocol/ent"
)

// CreateTestUser creates a test user with default or custom values
func CreateTestUser(client *ent.Client, overrides map[string]string) (*ent.User, error) {

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
	user, err := client.User.
		Create().
		SetFirstName(payload["firstName"]).
		SetLastName(payload["lastName"]).
		SetEmail(strings.ToLower(payload["email"])).
		SetPassword(payload["password"]).
		Save(context.Background())

	return user, err
}
