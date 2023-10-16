package services

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/user"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/paycrest/paycrest-protocol/utils/token"
)

// APIKeyService provides functionality related to API keys.
type APIKeyService struct{}

// NewAPIKeyService creates a new instance of APIKeyService.
func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{}
}

// GenerateAPIKey generates a new API key for the user.
func (s *APIKeyService) GenerateAPIKey(ctx context.Context, tx *ent.Tx, userID uuid.UUID) (*ent.APIKey, string, error) {
	// Generate a new secret key
	secretKey, err := token.GeneratePrivateKey()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate API key: %w", err)
	}

	// Encrypt the secret key
	encryptedSecret, err := crypto.EncryptPlain([]byte(secretKey))
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt API key: %w", err)
	}

	// Encode the encrypted secret to base64
	encodedSecret := base64.StdEncoding.EncodeToString(encryptedSecret)

	var apiKey *ent.APIKey

	if tx == nil {
		// Fetch the User entity from the database using the userID value
		user, err := storage.Client.User.
			Query().
			Where(user.IDEQ(userID)).
			Only(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("failed to fetch user: %w", err)
		}

		// Create a new APIKey entity
		apiKey, err = storage.Client.APIKey.
			Create().
			SetSecret(encodedSecret).
			SetOwner(user).
			Save(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create API key: %w", err)
		}
	} else {
		// Fetch the User entity from the database using the userID value
		user, err := tx.User.
			Query().
			Where(user.IDEQ(userID)).
			Only(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("failed to fetch user: %w", err)
		}

		// Create a new APIKey entity
		apiKey, err = tx.APIKey.
			Create().
			SetSecret(encodedSecret).
			SetOwner(user).
			Save(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create API key: %w", err)
		}
	}

	return apiKey, secretKey, nil
}
