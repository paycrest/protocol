package accounts

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/routers/middleware"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/enttest"
	"github.com/paycrest/paycrest-protocol/ent/providerprofile"
	"github.com/paycrest/paycrest-protocol/ent/user"
	svc "github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/test"
	"github.com/paycrest/paycrest-protocol/utils/token"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Set up test routers
	router := gin.New()
	ctrl := &AuthController{}
	router.POST("/register", ctrl.Register)
	router.POST("/login", ctrl.Login)
	router.POST("/refresh", middleware.JWTMiddleware, ctrl.RefreshJWT)
	router.POST("/api-keys", middleware.JWTMiddleware, ctrl.GenerateAPIKey)
	router.GET("/api-keys", middleware.JWTMiddleware, ctrl.ListAPIKeys)
	router.DELETE("/api-keys/:id", middleware.JWTMiddleware, ctrl.DeleteAPIKey)

	var userID string

	t.Run("Register", func(t *testing.T) {
		t.Run("with valid payload", func(t *testing.T) {
			// Test register with valid payload
			payload := svc.RegisterPayload{
				FirstName: "Ike",
				LastName:  "Ayo",
				Email:     "ikeayo@example.com",
				Password:  "password",
			}

			res, err := test.PerformRequest(t, "POST", "/register", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusCreated, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User created successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response data
			assert.Contains(t, data, "id")
			match, _ := regexp.MatchString(
				`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`,
				data["id"].(string),
			)
			if !match {
				t.Errorf("Expected '%s' to be a valid UUID", data["id"].(string))
			}
			userID = data["id"].(string)
			assert.Equal(t, payload.Email, data["email"].(string))
			assert.Equal(t, payload.FirstName, data["firstName"].(string))
			assert.Equal(t, payload.LastName, data["lastName"].(string))
		})

		t.Run("from the provider app", func(t *testing.T) {
			// Test register with valid payload
			payload := svc.RegisterPayload{
				FirstName:   "Ike",
				LastName:    "Ayo",
				Email:       "ikeayoprovider@example.com",
				Password:    "password",
				TradingName: "Africana LP",
				Country:     "Nigeria",
			}

			headers := map[string]string{
				"X-App-ID": "provider",
			}

			res, err := test.PerformRequest(t, "POST", "/register", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusCreated, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Parse the user ID string to uuid.UUID
			userUUID, err := uuid.Parse(data["id"].(string))
			assert.NoError(t, err)

			// Query the database to check if API key and profile were created for the provider
			apiKey, err := db.Client.APIKey.
				Query().
				Where(apikey.HasOwnerWith(user.ID(userUUID))).
				Only(context.Background())
			assert.NoError(t, err)

			providerProfile, err := db.Client.ProviderProfile.
				Query().
				Where(providerprofile.HasAPIKeyWith(apikey.ID(apiKey.ID))).
				Only(context.Background())
			assert.NoError(t, err)

			assert.NotNil(t, apiKey)
			assert.NotNil(t, providerProfile)
		})

		t.Run("with existing user", func(t *testing.T) {
			// Test register with existing user
			payload := svc.RegisterPayload{
				FirstName: "Ike",
				LastName:  "Ayo",
				Email:     "ikeayo@example.com",
				Password:  "password",
			}

			res, err := test.PerformRequest(t, "POST", "/register", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusBadRequest, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User with email already exists", response.Message)
			assert.Nil(t, response.Data)
		})

		t.Run("with invalid email", func(t *testing.T) {
			// Test register with invalid email
			payload := svc.RegisterPayload{
				FirstName: "Ike",
				LastName:  "Ayo",
				Email:     "invalid-email",
				Password:  "password",
			}

			res, err := test.PerformRequest(t, "POST", "/register", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusBadRequest, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to validate payload", response.Message)
			assert.Equal(t, "error", response.Status)
			data, ok := response.Data.([]interface{})
			assert.True(t, ok, "response.Data is not of type []interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response errors in data
			assert.Len(t, data, 1)
			errorMap, ok := data[0].(map[string]interface{})
			assert.True(t, ok, "error is not of type map[string]interface{}")
			assert.NotNil(t, errorMap, "error is nil")
			assert.Contains(t, errorMap, "field")
			assert.Equal(t, "Email", errorMap["field"].(string))
			assert.Contains(t, errorMap, "message")
			assert.Equal(t, "Must be a valid email address", errorMap["message"].(string))
		})
	})

	t.Run("Login", func(t *testing.T) {
		t.Run("with valid credentials", func(t *testing.T) {
			// Test login with valid credentials
			payload := svc.LoginPayload{
				Email:    "ikeayo@example.com",
				Password: "password",
			}

			res, err := test.PerformRequest(t, "POST", "/login", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Successfully logged in", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response data
			assert.Contains(t, data, "accessToken")
			assert.NotEmpty(t, data["accessToken"].(string))
			assert.Contains(t, data, "refreshToken")
			assert.NotEmpty(t, data["refreshToken"].(string))
		})

		t.Run("with invalid credentials", func(t *testing.T) {
			// Test login with invalid credentials
			payload := svc.LoginPayload{
				Email:    "ikeayo@example.com",
				Password: "wrong-password",
			}

			res, err := test.PerformRequest(t, "POST", "/login", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusUnauthorized, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid credentials", response.Message)
			assert.Equal(t, "Email and password do not match any user", response.Data)
		})
	})

	t.Run("RefreshJWT", func(t *testing.T) {
		t.Run("with a valid refresh token", func(t *testing.T) {
			refreshToken, err := token.GenerateRefreshJWT(userID)
			assert.NoError(t, err, "failed to generate refresh token")

			// Test refresh token with valid refresh token
			payload := svc.RefreshJWTPayload{
				RefreshToken: refreshToken,
			}

			headers := map[string]string{
				"Authorization": "Bearer " + refreshToken,
			}

			res, err := test.PerformRequest(t, "POST", "/refresh", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Successfully refreshed access token", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response data
			assert.Contains(t, data, "accessToken")
			assert.NotEmpty(t, data["accessToken"].(string))
			assert.NotContains(t, data, "refreshToken")
		})

		t.Run("with an invalid refresh token", func(t *testing.T) {
			refreshToken := "invalid-refresh-token"

			// Test refresh token with invalid refresh token
			payload := svc.RefreshJWTPayload{
				RefreshToken: refreshToken,
			}

			refreshTokenForHeader, err := token.GenerateRefreshJWT(userID)
			assert.NoError(t, err, "failed to generate refresh token")

			headers := map[string]string{
				"Authorization": "Bearer " + refreshTokenForHeader,
			}
			res, err := test.PerformRequest(t, "POST", "/refresh", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusUnauthorized, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid or expired refresh token", response.Message)
			assert.NotNil(t, response.Data, "response.Data is nil")
		})
	})

	t.Run("GenerateAPIKey", func(t *testing.T) {
		accessToken, _ := token.GenerateAccessJWT(userID)

		t.Run("with a valid scope", func(t *testing.T) {
			payload := svc.GenerateAPIKeyPayload{
				Name:  "Test API Key",
				Scope: "sender",
			}

			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}

			res, err := test.PerformRequest(t, "POST", "/api-keys", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusCreated, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Successfully generated API key", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, "sender", data["scope"])
			assert.True(t, data["isActive"].(bool))

			// Assert that id is a valid UUID
			assert.Contains(t, data, "id")
			id, err := uuid.Parse(data["id"].(string))
			assert.NoError(t, err)
			_ = id

			// Assert that the secret is a base64 encoded string
			secret, err := base64.URLEncoding.DecodeString(data["secret"].(string))
			assert.NoError(t, err)
			_ = secret
		})

		t.Run("with an invalid scope", func(t *testing.T) {
			payload := svc.GenerateAPIKeyPayload{
				Name:  "Test API Key",
				Scope: "bad-scope",
			}

			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}

			res, err := test.PerformRequest(t, "POST", "/api-keys", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusBadRequest, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid request body", response.Message)
			assert.Equal(t, "error", response.Status)
			data, ok := response.Data.([]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response errors in data
			assert.Len(t, data, 1)
			errorMap, ok := data[0].(map[string]interface{})
			assert.True(t, ok, "error is not of type map[string]interface{}")
			assert.NotNil(t, errorMap, "error is nil")
			assert.Contains(t, errorMap, "field")
			assert.Equal(t, "Scope", errorMap["field"].(string))
			assert.Contains(t, errorMap, "message")
			assert.Equal(t, "Must be one of sender provider tx_validator", errorMap["message"].(string))
		})
	})

	t.Run("ListAPIKeys", func(t *testing.T) {
		accessToken, _ := token.GenerateAccessJWT(userID)
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		res, err := test.PerformRequest(t, "GET", "/api-keys", nil, headers, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusOK, res.Code)

		var response utils.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Successfully retrieved API keys", response.Message)
		assert.Equal(t, "success", response.Status)
		data, ok := response.Data.([]interface{})
		assert.True(t, ok, "response.Data is not of type []interface{}")
		assert.NotNil(t, data, "response.Data is nil")

		assert.Len(t, data, 1)
		dataMap, ok := data[0].(map[string]interface{})
		assert.True(t, ok, "error is not of type map[string]interface{}")
		assert.Contains(t, dataMap, "secret")
		assert.Contains(t, dataMap, "id")
	})

	t.Run("DeleteAPIKey", func(t *testing.T) {
		accessToken, _ := token.GenerateAccessJWT(userID)

		t.Run("with a valid API key", func(t *testing.T) {

			userUUID, err := uuid.Parse(userID)
			assert.NoError(t, err)

			// Create a context with a timeout of 5 seconds
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// BEFORE: Query the user's API keys
			apiKeys, err := client.User.
				Query().
				Where(user.IDEQ(userUUID)).
				QueryAPIKeys().
				All(ctx)

			assert.NoError(t, err)
			assert.Len(t, apiKeys, 1)

			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}

			res, err := test.PerformRequest(t, "DELETE", "/api-keys/"+apiKeys[0].ID.String(), nil, headers, router)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusNoContent, res.Code)

			// AFTER: Query the user's API keys again
			apiKeysAgain, err := client.User.
				Query().
				Where(user.IDEQ(userUUID)).
				QueryAPIKeys().
				All(ctx)

			assert.NoError(t, err)
			assert.Len(t, apiKeysAgain, 0)
		})

		t.Run("with an invalid API key", func(t *testing.T) {

			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			res, err := test.PerformRequest(t, "DELETE", "/api-keys/invalid-api-key", nil, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusBadRequest, res.Code)

			var response utils.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, "Invalid API key ID", response.Message)
		})
	})
}
