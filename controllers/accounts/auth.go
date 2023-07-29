package accounts

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/user"
	svc "github.com/paycrest/paycrest-protocol/services"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/paycrest/paycrest-protocol/utils/token"
)

// AuthController is the controller type for the auth endpoints
type AuthController struct {
	// apiKeyService *svc.APIKeyService
}

// NewAuthController creates a new instance of AuthController with injected services
// func NewAuthController() *AuthController {
// 	return &AuthController{
// 		apiKeyService: svc.NewAPIKeyService(db.Client),
// 	}
// }

// Register controller validates the payload and creates a new user.
// It hashes the password provided by the user.
// It also sends an email to verify the user's email address.
func (ctrl *AuthController) Register(ctx *gin.Context) {
	var payload svc.RegisterPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Check if user with email already exists
	userTmp, _ := db.Client.User.
		Query().
		Where(user.EmailEQ(strings.ToLower(payload.Email))).
		Only(ctx)

	if userTmp != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"User with email already exists", nil)
		return
	}

	// Save the user
	user, err := db.Client.User.
		Create().
		SetFirstName(payload.FirstName).
		SetLastName(payload.LastName).
		SetEmail(strings.ToLower(payload.Email)).
		SetPassword(payload.Password).
		Save(ctx)

	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to create new user", err.Error())
		return
	}

	// Create a provider API Key and profile in the background
	// TODO: Replace provider with a UUID environment variable
	if appID := ctx.GetHeader("X-APP-ID"); appID == "provider" {
		apiKeyInput := svc.CreateAPIKeyPayload{
			Name:  payload.TradingName + " API Key",
			Scope: apikey.ScopeProvider,
		}

		// Generate the API key using the service
		apiKeyService := svc.NewAPIKeyService(db.Client)
		apiKey, _, err := apiKeyService.GenerateAPIKey(ctx, user.ID, apiKeyInput)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", err.Error())
			return
		}

		// Create a provider profile
		_, err = db.Client.ProviderProfile.
			Create().
			SetTradingName(payload.TradingName).
			SetCountry(payload.Country).
			SetAPIKey(apiKey).
			Save(ctx)

		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", err.Error())
			return
		}
	}

	// TODO: Send email to verify the user's email address

	u.APIResponse(ctx, http.StatusCreated, "success", "User created successfully",
		&svc.RegisterResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
}

// Login controller validates the payload and creates a new user.
func (ctrl *AuthController) Login(ctx *gin.Context) {
	var payload svc.LoginPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Fetch user by email
	user, emailErr := db.Client.User.
		Query().
		Where(user.EmailEQ(strings.ToLower(payload.Email))).
		Only(ctx)

	// Check if the password is correct
	passwordMatch := crypto.CheckPasswordHash(payload.Password, user.Password)

	if !passwordMatch || emailErr != nil {
		logger.Errorf("error: %v", "Invalid email or password")
		u.APIResponse(ctx, http.StatusUnauthorized, "error",
			"Invalid credentials", "Email and password do not match any user",
		)
		return
	}

	// Generate JWT pair
	accessToken, refreshToken, err := token.GeneratePairJWT(user.ID.String())

	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to create token pair", err.Error(),
		)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Successfully logged in", &svc.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshJWT controller returns a new access token given a valid refresh token.
func (ctrl *AuthController) RefreshJWT(ctx *gin.Context) {
	var payload svc.RefreshJWTPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Validate the refresh token
	claims, err := token.ValidateJWT(payload.RefreshToken)
	userID, ok := claims["sub"].(string)
	if err != nil || !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid or expired refresh token", err.Error())
		return
	}

	// Generate a new access token
	accessToken, err := token.GenerateAccessJWT(userID)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to generate access token", err.Error())
		return
	}

	// Return the new access token
	u.APIResponse(ctx, http.StatusOK, "success", "Successfully refreshed access token", &svc.RefreshResponse{
		AccessToken: accessToken,
	})
}

// CreateAPIKey controller creates a new API key pair for the user.
func (ctrl *AuthController) CreateAPIKey(ctx *gin.Context) {
	// Get the user ID from the context
	userIDString, _ := ctx.Get("user_id")

	// Parse the user ID string to uuid.UUID
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		logger.Errorf("error parsing user ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid user ID", nil)
		return
	}

	var payload svc.CreateAPIKeyPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid request body", u.GetErrorData(err))
		return
	}

	// Generate the API key using the service
	apiKeyService := svc.NewAPIKeyService(db.Client)
	apiKey, secretKey, err := apiKeyService.GenerateAPIKey(ctx, userID, payload)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to generate API key", err.Error())
		return
	}

	// Return the newly generated API key
	u.APIResponse(ctx, http.StatusCreated, "success", "Successfully generated API key", &svc.APIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Scope:     apiKey.Scope,
		Secret:    secretKey,
		IsActive:  apiKey.IsActive,
		CreatedAt: apiKey.CreatedAt,
	})
}

// ListAPIKeys controller returns all API keys for the user.
func (ctrl *AuthController) ListAPIKeys(ctx *gin.Context) {
	// Get the user ID from the context
	userIDString, _ := ctx.Get("user_id")

	// Parse the user ID string to uuid.UUID
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		logger.Errorf("error parsing user ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid user ID", nil)
		return
	}

	// Query the user's API keys
	apiKeys, err := db.Client.User.
		Query().
		Where(user.IDEQ(userID)).
		QueryAPIKeys().
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to retrieve API keys", err.Error())
		return
	}

	// Create APIKeyResponse objects without the Pair field
	apiKeyResponses := make([]svc.APIKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		// Decode the stored secret key to bytes
		decodedSecret, err := base64.StdEncoding.DecodeString(apiKey.Secret)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decode API key", err.Error())
			return
		}

		// Decrypt the decoded secret
		decryptedSecret, err := crypto.Decrypt(decodedSecret)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decrypt API key", err.Error())
			return
		}
		apiKeyResponses[i] = svc.APIKeyResponse{
			ID:        apiKey.ID,
			CreatedAt: apiKey.CreatedAt,
			Name:      apiKey.Name,
			Scope:     apiKey.Scope,
			Secret:    string(decryptedSecret),
			IsActive:  apiKey.IsActive,
		}
	}

	// Return the API keys
	u.APIResponse(ctx, http.StatusOK, "success", "Successfully retrieved API keys", apiKeyResponses)
}

// DeleteAPIKey controller deletes an API key for the user.
func (ctrl *AuthController) DeleteAPIKey(ctx *gin.Context) {
	// Get the API key ID from the request URL
	apiKeyID := ctx.Param("id")

	// Parse the API key ID string to uuid.UUID
	apiKeyUUID, err := uuid.Parse(apiKeyID)
	if err != nil {
		logger.Errorf("error parsing API key ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid API key ID", nil)
		return
	}

	// Get the user ID from the context
	userIDString, _ := ctx.Get("user_id")

	// Parse the user ID string to uuid.UUID
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		logger.Errorf("error parsing user ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid user ID", nil)
		return
	}

	// Check if the API key belongs to the user making the request
	apiKey, err := db.Client.APIKey.
		Query().
		Where(apikey.IDEQ(apiKeyUUID), apikey.HasOwnerWith(user.IDEQ(userID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			u.APIResponse(ctx, http.StatusNotFound, "error", "API key not found", nil)
		} else {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to delete API key", err.Error())
		}
		return
	}

	// Delete the API key
	err = db.Client.APIKey.
		DeleteOne(apiKey).
		Exec(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to delete API key", err.Error())
		return
	}

	// Return a success response
	u.APIResponse(ctx, http.StatusNoContent, "success", "API key deleted successfully", nil)
}
