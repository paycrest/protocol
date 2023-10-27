package accounts

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/routers/middleware"
	"github.com/paycrest/protocol/services"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/paycrest/protocol/ent/user"
	"github.com/paycrest/protocol/ent/validatorprofile"
	"github.com/paycrest/protocol/utils/test"
	"github.com/paycrest/protocol/utils/token"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Set up test routers
	router := gin.New()
	ctrl := &ProfileController{}

	router.Use(middleware.ScopeMiddleware)
	router.GET(
		"settings/validator",
		middleware.JWTMiddleware,
		middleware.OnlyValidatorMiddleware,
		ctrl.GetValidatorProfile,
	)
	router.PATCH(
		"settings/validator",
		middleware.JWTMiddleware,
		middleware.OnlyValidatorMiddleware,
		ctrl.UpdateValidatorProfile,
	)
	router.GET(
		"settings/sender",
		middleware.JWTMiddleware,
		middleware.OnlySenderMiddleware,
		ctrl.GetSenderProfile,
	)
	router.PATCH(
		"settings/sender",
		middleware.JWTMiddleware,
		middleware.OnlySenderMiddleware,
		ctrl.UpdateSenderProfile,
	)

	t.Run("UpdateSenderProfile", func(t *testing.T) {
		testUser, err := test.CreateTestUser(map[string]string{"scope": "sender"})
		assert.NoError(t, err)

		_, err = test.CreateTestSenderProfile(map[string]interface{}{
			"domain_whitelist": []string{"example.com"},
			"user_id":          testUser.ID,
		})
		assert.NoError(t, err)

		// Test partial update
		accessToken, _ := token.GenerateAccessJWT(testUser.ID.String())
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		payload := types.SenderProfilePayload{
			DomainWhitelist: []string{"example.com", "mydomain.com"},
		}

		res, err := test.PerformRequest(t, "PATCH", "/settings/sender?scope=sender", payload, headers, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusOK, res.Code)

		var response types.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Profile updated successfully", response.Message)
		assert.Nil(t, response.Data, "response.Data is not nil")

		senderProfile, err := db.Client.SenderProfile.
			Query().
			Where(senderprofile.HasUserWith(user.ID(testUser.ID))).
			Only(context.Background())
		assert.NoError(t, err)

		assert.Contains(t, senderProfile.DomainWhitelist, "mydomain.com")
	})

	t.Run("UpdateValidatorProfile", func(t *testing.T) {
		testUser, err := test.CreateTestUser(map[string]string{"scope": "tx_validator"})
		assert.NoError(t, err)

		_, err = test.CreateTestValidatorProfile(map[string]interface{}{
			"wallet_address":  "0x000000000000000000000000000000000000dEaD",
			"host_identifier": "example.com",
			"user_id":         testUser.ID,
		})
		assert.NoError(t, err)

		// Test partial update
		accessToken, _ := token.GenerateAccessJWT(testUser.ID.String())
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		payload := types.ValidatorProfilePayload{
			HostIdentifier: "127.0.0.1:8080",
		}

		res, err := test.PerformRequest(t, "PATCH", "/settings/validator?scope=validator", payload, headers, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusOK, res.Code)

		var response types.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Profile updated successfully", response.Message)
		assert.Nil(t, response.Data, "response.Data is not nil")

		validatorProfile, err := db.Client.ValidatorProfile.
			Query().
			Where(validatorprofile.HasUserWith(user.ID(testUser.ID))).
			Only(context.Background())
		assert.NoError(t, err)

		assert.Equal(t, "127.0.0.1:8080", validatorProfile.HostIdentifier)
	})

	t.Run("GetSenderProfile", func(t *testing.T) {
		testUser, err := test.CreateTestUser(map[string]string{
			"email": "hello@test.com",
			"scope": "sender",
		})
		assert.NoError(t, err)

		sender, err := test.CreateTestSenderProfile(map[string]interface{}{
			"domain_whitelist": []string{"mydomain.com"},
			"user_id":          testUser.ID,
		})
		assert.NoError(t, err)

		apiKeyService := services.NewAPIKeyService()
		_, _, err = apiKeyService.GenerateAPIKey(
			context.Background(),
			nil,
			sender,
			nil,
			nil,
		)
		assert.NoError(t, err)

		accessToken, _ := token.GenerateAccessJWT(testUser.ID.String())
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		res, err := test.PerformRequest(t, "GET", "/settings/sender?scope=sender", nil, headers, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusOK, res.Code)
		var response types.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Profile retrieved successfully", response.Message)
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "response.Data is not of type map[string]interface{}")
		assert.NotNil(t, data, "response.Data is nil")

		assert.Contains(t, data["domainWhitelist"], "mydomain.com")
	})

	t.Run("GetValidatorProfile", func(t *testing.T) {
		testUser, err := test.CreateTestUser(map[string]string{
			"email": "hello@test.com",
			"scope": "tx_validator",
		})
		assert.NoError(t, err)

		validator, err := test.CreateTestValidatorProfile(map[string]interface{}{
			"wallet_address": "0x0000000000",
			"user_id":        testUser.ID,
		})
		assert.NoError(t, err)

		apiKeyService := services.NewAPIKeyService()
		_, _, err = apiKeyService.GenerateAPIKey(
			context.Background(),
			nil,
			nil,
			nil,
			validator,
		)
		assert.NoError(t, err)

		accessToken, _ := token.GenerateAccessJWT(testUser.ID.String())
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		res, err := test.PerformRequest(t, "GET", "/settings/validator?scope=validator", nil, headers, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusOK, res.Code)
		var response types.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Profile retrieved successfully", response.Message)
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "response.Data is not of type map[string]interface{}")
		assert.NotNil(t, data, "response.Data is nil")

		assert.Equal(t, "0x0000000000", data["walletAddress"])
	})
}
