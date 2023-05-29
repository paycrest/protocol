package accounts

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/user"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/paycrest/paycrest-protocol/utils/token"
)

// AuthController is the controller type for the auth endpoints
type AuthController struct{}

// Register controller validates the payload and creates a new user.
// It hashes the password provided by the user.
// It also sends an email to verify the user's email address.
func (ctrl *AuthController) Register(ctx *gin.Context) {
	var payload RegisterPayload

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

	// TODO: Send email to verify the user's email address

	u.APIResponse(ctx, http.StatusCreated, "success", "User created successfully",
		&RegisterResponse{
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
	var payload LoginPayload

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

	u.APIResponse(ctx, http.StatusOK, "success", "Successfully logged in", &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshJWT controller returns a new access token given a valid refresh token.
func (ctrl *AuthController) RefreshJWT(ctx *gin.Context) {
	var payload RefreshJWTPayload

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
	u.APIResponse(ctx, http.StatusOK, "success", "Successfully refreshed access token", &RefreshResponse{
		AccessToken: accessToken,
	})
}

func mapPayloadScopeToAPIKeyScope(payloadScope string) apikey.Scope {
	switch payloadScope {
	case "provider":
		return apikey.ScopeProvider
	case "sender":
		return apikey.ScopeSender
	case "validator":
		return apikey.ScopeTxValidator
	default:
		// Return a default value or handle the error accordingly
		return apikey.ScopeProvider
	}
}

// GenerateAPIKey controller generates a new API key pair for the user.
func (ctrl *AuthController) GenerateAPIKey(ctx *gin.Context) {
	// Get the user ID from the context
	userIDString, _ := ctx.Get("user_id")

	// Parse the user ID string to uuid.UUID
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		logger.Errorf("error parsing user ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid user ID", nil)
		return
	}

	var payload GenerateAPIKeyPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid request body", u.GetErrorData(err))
		return
	}

	// Generate a new API key pair
	publicKey, secretKey, err := token.GenerateHMACKeys()
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to generate API key", err.Error())
		return
	}

	// Fetch the User entity from the database using the userID value
	user, err := db.Client.User.Get(ctx, userID)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to generate API key", err.Error())
		return
	}

	// Create a new APIKey entity
	apiKey, err := db.Client.APIKey.
		Create().
		SetName(payload.Name).
		SetScope(mapPayloadScopeToAPIKeyScope(payload.Scope)).
		SetPair(publicKey + "::" + secretKey).
		SetUser(user).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to generate API key", err.Error())
		return
	}

	// Return the newly generated API key
	u.APIResponse(ctx, http.StatusOK, "success", "Successfully generated API key", &GenerateAPIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Scope:     apiKey.Scope,
		PublicKey: publicKey,
		SecretKey: secretKey,
		IsActive:  apiKey.IsActive,
		CreatedAt: apiKey.CreatedAt,
	})
}

// GetAPIKeys controller returns all API keys for the user.
func (ctrl *AuthController) GetAPIKeys(ctx *gin.Context) {
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

	// TODO: split API key pair to public and secret keys

	// Return the API keys
	u.APIResponse(ctx, http.StatusOK, "success", "Successfully retrieved API keys", apiKeys)
}

// DeleteAPIKey controller deletes an API key for the user.
func (ctrl *AuthController) DeleteAPIKey(ctx *gin.Context) {
	// Return the new access token
	u.APIResponse(ctx, http.StatusOK, "success", "Successfully deleted API key", nil)
}
