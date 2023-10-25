package services

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"
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
func (s *APIKeyService) GenerateAPIKey(
	ctx context.Context,
	tx *ent.Tx,
	sender *ent.SenderProfile,
	provider *ent.ProviderProfile,
	validator *ent.ValidatorProfile,
) (*ent.APIKey, string, error) {
	// Generate a new secret key
	secretKey, err := token.GeneratePrivateKey()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate API key: %w", err)
	}

	// Encrypt the secret key
	encryptedSecret, _ := crypto.EncryptPlain([]byte(secretKey))
	encodedSecret := base64.StdEncoding.EncodeToString(encryptedSecret)

	var apiKey *ent.APIKey

	if sender != nil {
		if tx != nil {
			apiKey, err = tx.APIKey.
				Create().
				SetSecret(encodedSecret).
				SetSenderProfile(sender).
				Save(ctx)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create API key: %w", err)
			}
		} else {
			apiKey, err = storage.Client.APIKey.
				Create().
				SetSecret(encodedSecret).
				SetSenderProfile(sender).
				Save(ctx)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create API key: %w", err)
			}
		}
	} else if provider != nil {
		if tx != nil {
			apiKey, err = tx.APIKey.
				Create().
				SetSecret(encodedSecret).
				SetProviderProfile(provider).
				Save(ctx)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create API key: %w", err)
			}
		} else {
			apiKey, err = storage.Client.APIKey.
				Create().
				SetSecret(encodedSecret).
				SetProviderProfile(provider).
				Save(ctx)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create API key: %w", err)
			}
		}
	} else if validator != nil {
		if tx != nil {
			apiKey, err = tx.APIKey.
				Create().
				SetSecret(encodedSecret).
				SetValidatorProfile(validator).
				Save(ctx)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create API key: %w", err)
			}
		} else {
			apiKey, err = storage.Client.APIKey.
				Create().
				SetSecret(encodedSecret).
				SetValidatorProfile(validator).
				Save(ctx)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create API key: %w", err)
			}
		}
	} else {
		return nil, "", fmt.Errorf("profile not provided")
	}

	return apiKey, secretKey, nil
}

// GetAPIKey gets the API key for a user profile.
func (s *APIKeyService) GetAPIKey(
	ctx context.Context,
	sender *ent.SenderProfile,
	provider *ent.ProviderProfile,
	validator *ent.ValidatorProfile,
) (*types.APIKeyResponse, error) {
	var apiKey *ent.APIKey

	if sender != nil {
		apiKey, _ = sender.QueryAPIKey().Only(ctx)
	} else if provider != nil {
		apiKey, _ = provider.QueryAPIKey().Only(ctx)
	} else if validator != nil {
		apiKey, _ = validator.QueryAPIKey().Only(ctx)
	} else {
		return nil, fmt.Errorf("profile not provided")
	}

	// Decrypt the secret key
	decodedSecret, _ := base64.StdEncoding.DecodeString(apiKey.Secret)
	decryptedSecret, _ := crypto.DecryptPlain(decodedSecret)

	return &types.APIKeyResponse{
		ID:     apiKey.ID,
		Secret: string(decryptedSecret),
	}, nil
}
