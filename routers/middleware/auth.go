package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/predicate"
	"github.com/paycrest/paycrest-protocol/ent/user"
	"github.com/paycrest/paycrest-protocol/storage"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/paycrest/paycrest-protocol/utils/token"
)

// JWTMiddleware is a middleware to handle JWT authentication
func JWTMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		u.APIResponse(c, http.StatusUnauthorized, "error",
			"Authorization header is missing", "Expected: Bearer <token>")
		c.Abort()
		return
	}

	// Split the Authorization header value into two parts: the authentication scheme and the token value
	authParts := strings.SplitN(authHeader, " ", 2)
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		u.APIResponse(c, http.StatusUnauthorized, "error",
			"Invalid Authorization header format", "Expected: Bearer <token>")
		c.Abort()
		return
	}

	// Validate the token and extract the user ID
	claims, err := token.ValidateJWT(authParts[1])
	userID, ok := claims["sub"].(string)
	if err != nil || !ok {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid or expired token", err.Error())
		c.Abort()
		return
	}

	// Set the user_id value in the context of the request
	c.Set("user_id", userID)

	c.Next()
}

// HMACVerificationMiddleware is a middleware for HMAC verification.
// It verifies the HMAC signature in the Authorization header of the request.
func HMACVerificationMiddleware(c *gin.Context) {
	// Get the authorization header value
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Missing Authorization header", nil)
		c.Abort()
		return
	}

	// Parse the authorization header
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "HMAC" {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid Authorization header", "Expected: HMAC <public_key>:<signature>")
		c.Abort()
		return
	}

	// Extract the public key and signature
	parts = strings.SplitN(parts[1], ":", 2)
	publicKey, signature := parts[0], parts[1]
	if publicKey == "" || signature == "" {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid Authorization header format", "Expected: HMAC <public_key>:<signature>")
		c.Abort()
		return
	}

	// Get the request payload
	payload, err := c.GetRawData()
	if err != nil {
		u.APIResponse(c, http.StatusInternalServerError, "error", "Failed to read request payload", err.Error())
		c.Abort()
		return
	}

	// Parse the payload to retrieve timestamp
	var payloadData map[string]interface{}
	err = json.Unmarshal(payload, &payloadData)
	if err != nil {
		u.APIResponse(c, http.StatusBadRequest, "error", "Invalid payload format", err.Error())
		c.Abort()
		return
	}

	// Get the timestamp from the payload
	timestamp, ok := payloadData["timestamp"].(float64) // unix timestamp
	if !ok || timestamp == 0 {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Missing or invalid timestamp in payload", nil)
		c.Abort()
		return
	}

	var conf = config.AuthConfig()

	// Check if the timestamp is within the acceptable window
	if time.Now().Unix()-int64(timestamp) > int64(conf.HmacTimestampAge.Seconds()) {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid timestamp", nil)
		c.Abort()
		return
	}

	// Parse the API key ID string to uuid.UUID
	apiKeyUUID, err := uuid.Parse(publicKey)
	if err != nil {
		logger.Errorf("error parsing API key ID: %v", err)
		u.APIResponse(c, http.StatusBadRequest, "error", "Invalid API key ID", nil)
		c.Abort()
		return
	}

	// Fetch the API key from the database
	apiKey, err := storage.Client.APIKey.
		Query().
		Where(apikey.IDEQ(apiKeyUUID)).
		Only(c)
	if err != nil {
		if ent.IsNotFound(err) {
			u.APIResponse(c, http.StatusNotFound, "error", "API key not found", nil)
		} else {
			logger.Errorf("error: %v", err)
			u.APIResponse(c, http.StatusInternalServerError, "error", "Failed to fetch API key", err.Error())
		}
		c.Abort()
		return
	}

	c.Set("api_key", apiKey)

	// Decode the stored secret key to bytes
	decodedSecret, err := base64.StdEncoding.DecodeString(apiKey.Secret)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(c, http.StatusInternalServerError, "error", "Failed to decode API key", err.Error())
		return
	}

	// Decrypt the decoded secret
	decryptedSecret, err := crypto.DecryptPlain(decodedSecret)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(c, http.StatusInternalServerError, "error", "Failed to decrypt API key", err.Error())
		return
	}

	// Verify the HMAC signature
	valid := token.VerifyHMACSignature(payloadData, string(decryptedSecret), signature)
	if !valid {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid HMAC signature", nil)
		c.Abort()
		return
	}

	// Remove the timestamp key from the payload
	delete(payloadData, "timestamp")

	// Convert the payload data back to JSON
	modifiedPayload, err := json.Marshal(payloadData)
	if err != nil {
		u.APIResponse(c, http.StatusInternalServerError, "error", "Failed to modify payload", err.Error())
		c.Abort()
		return
	}

	// Create a new buffer with the modified payload
	buffer := bytes.NewBuffer(modifiedPayload)

	// Set the modified payload as the request body
	c.Request.Body = io.NopCloser(buffer)

	// Continue to the next middleware
	c.Next()
}

// DynamicAuthMiddleware is a middleware that dynamically selects the authentication method
func DynamicAuthMiddleware(c *gin.Context) {
	// Check the request headers to determine the desired authentication method
	clientType := c.GetHeader("Client-Type")

	// Select the authentication middleware based on the client type
	switch clientType {
	case "frontend":
		JWTMiddleware(c)
	case "backend":
		HMACVerificationMiddleware(c)
	default:
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid client type", nil)
		c.Abort()
		return
	}

	c.Next()
}

// ScopeMiddleware is a middleware that checks the request scope.
func ScopeMiddleware(c *gin.Context) {
	conf := config.ServerConfig()
	var subdomain string

	if conf.Debug {
		subdomain = c.Query("scope")
	} else {
		host := c.Request.Host
		parts := strings.Split(host, ".")
		subdomain = parts[0]
	}

	if subdomain == "validator" {
		subdomain = "tx_validator"
	}

	scope := user.Scope(subdomain)

	if scope == "" {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid scope", nil)
		c.Abort()
		return
	}

	c.Set("scope", scope)
	c.Next()
}

// OnlySenderMiddleware is a middleware that checks if the user scope is sender.
func OnlySenderMiddleware(c *gin.Context) {
	apiKeyCtx, apiKeyOk := c.Get("api_key")
	userID, userOk := c.Get("user_id")

	user, err := storage.Client.User.
		Query().
		Where(
			user.ScopeEQ(user.ScopeSender),
			func() predicate.User {
				if apiKeyOk {
					// HMAC auth was used
					return user.HasAPIKeysWith(
						apikey.IDEQ(apiKeyCtx.(*ent.APIKey).ID),
					)
				} else if userOk {
					// JWT auth was used
					userUUID, _ := uuid.Parse(userID.(string))
					return user.IDEQ(userUUID)
				}
				return nil
			}(),
		).
		WithSenderProfile().
		Only(c)
	if err != nil {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid API key", nil)
		c.Abort()
		return
	}

	c.Set("sender", user.Edges.SenderProfile)
	c.Next()
}

// OnlyProviderMiddleware is a middleware that checks if the user scope is provider.
func OnlyProviderMiddleware(c *gin.Context) {
	apiKeyCtx, apiKeyOk := c.Get("api_key")
	userID, userOk := c.Get("user_id")

	user, err := storage.Client.User.
		Query().
		Where(
			user.ScopeEQ(user.ScopeProvider),
			func() predicate.User {
				if apiKeyOk {
					// HMAC auth was used
					return user.HasAPIKeysWith(
						apikey.IDEQ(apiKeyCtx.(*ent.APIKey).ID),
					)
				} else if userOk {
					// JWT auth was used
					userUUID, _ := uuid.Parse(userID.(string))
					return user.IDEQ(userUUID)
				}
				return nil
			}(),
		).
		WithProviderProfile().
		Only(c)
	if err != nil {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid API key", nil)
		c.Abort()
		return
	}

	c.Set("provider", user.Edges.ProviderProfile)
	c.Next()
}

// OnlyValidatorMiddleware is a middleware that checks if the user scope is validator.
func OnlyValidatorMiddleware(c *gin.Context) {
	apiKeyCtx, apiKeyOk := c.Get("api_key")
	userID, userOk := c.Get("user_id")

	user, err := storage.Client.User.
		Query().
		Where(
			user.ScopeEQ(user.ScopeTxValidator),
			func() predicate.User {
				if apiKeyOk {
					// HMAC auth was used
					return user.HasAPIKeysWith(
						apikey.IDEQ(apiKeyCtx.(*ent.APIKey).ID),
					)
				} else if userOk {
					// JWT auth was used
					userUUID, _ := uuid.Parse(userID.(string))
					return user.IDEQ(userUUID)
				}
				return nil
			}(),
		).
		WithValidatorProfile().
		Only(c)
	if err != nil {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid API keyaa", nil)
		c.Abort()
		return
	}

	c.Set("validator", user.Edges.ValidatorProfile)
	c.Next()
}
