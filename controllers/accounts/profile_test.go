package accounts

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/routers/middleware"
	"github.com/paycrest/protocol/services"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/paycrest/protocol/ent/user"
	"github.com/paycrest/protocol/utils/test"
	"github.com/paycrest/protocol/utils/token"
	"github.com/stretchr/testify/assert"
)

var testCtx = struct {
	user            *ent.User
	providerProfile *ent.ProviderProfile
}{}

func setup() error {
	// Set up test data
	user, err := test.CreateTestUser(map[string]string{
		"scope": "provider",
		"email": "providerjohndoe@test.com",
	})
	if err != nil {
		return err
	}
	testCtx.user = user

	currency, err := test.CreateTestFiatCurrency(map[string]interface{}{
		"code":        "KES",
		"short_name":  "Shilling",
		"decimals":    2,
		"symbol":      "KSh",
		"name":        "Kenyan Shilling",
		"market_rate": 550.0,
	})
	if err != nil {
		return err
	}

	provderProfile, err := test.CreateTestProviderProfile(map[string]interface{}{
		"user_id":     testCtx.user.ID,
		"currency_id": currency.ID,
	})
	if err != nil {
		return err
	}
	testCtx.providerProfile = provderProfile

	return nil
}

func TestProfile(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setup()
	assert.NoError(t, err)

	// Set up test routers
	router := gin.New()
	ctrl := &ProfileController{}

	router.GET(
		"settings/sender",
		middleware.JWTMiddleware,
		middleware.OnlySenderMiddleware,
		ctrl.GetSenderProfile,
	)
	router.GET(
		"settings/provider",
		middleware.JWTMiddleware,
		middleware.OnlyProviderMiddleware,
		ctrl.GetProviderProfile,
	)
	router.PATCH(
		"settings/sender",
		middleware.JWTMiddleware,
		middleware.OnlySenderMiddleware,
		ctrl.UpdateSenderProfile,
	)
	router.PATCH(
		"settings/provider",
		middleware.JWTMiddleware,
		middleware.OnlyProviderMiddleware,
		ctrl.UpdateProviderProfile,
	)

	t.Run("UpdateSenderProfile", func(t *testing.T) {
		t.Run("with all fields", func(t *testing.T) {
			testUser, err := test.CreateTestUser(map[string]string{"scope": "sender"})
			assert.NoError(t, err)

			_, err = test.CreateTestSenderProfile(map[string]interface{}{
				"domain_whitelist": []string{"example.com"},
				"user_id":          testUser.ID,
			})
			assert.NoError(t, err)

			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testUser.ID.String(), "sender")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.SenderProfilePayload{
				DomainWhitelist: []string{"example.com", "mydomain.com"},
				RefundAddress:   "0x1234567890",
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/sender", payload, headers, router)
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

		t.Run("with an invalid webhook", func(t *testing.T) {
			testUser, err := test.CreateTestUser(map[string]string{
				"scope": "sender",
				"email": "johndoe2@test.com",
			})
			assert.NoError(t, err)

			_, err = test.CreateTestSenderProfile(map[string]interface{}{
				"domain_whitelist": []string{"example.com"},
				"user_id":          testUser.ID,
			})
			assert.NoError(t, err)

			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testUser.ID.String(), "sender")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.SenderProfilePayload{
				WebhookURL:      "examplecom",
				DomainWhitelist: []string{"example.com", "mydomain.com"},
				RefundAddress:   "0x1234567890",
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/sender", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusBadRequest, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid webhook url", response.Message)
			assert.Nil(t, response.Data, "response.Data is not nil")
		})

		t.Run("with all fields and check if it is active", func(t *testing.T) {
			testUser, err := test.CreateTestUser(map[string]string{
				"scope": "sender",
				"email": "johndoe3@test.com",
			})
			assert.NoError(t, err)

			_, err = test.CreateTestSenderProfile(map[string]interface{}{
				"domain_whitelist": []string{"example.com"},
				"user_id":          testUser.ID,
			})
			assert.NoError(t, err)

			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testUser.ID.String(), "sender")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.SenderProfilePayload{
				DomainWhitelist: []string{"example.com", "mydomain.com"},
				RefundAddress:   "0x1234567890",
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/sender", payload, headers, router)
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
			assert.True(t, senderProfile.IsActive)
		})

	})

	t.Run("UpdateProviderProfile", func(t *testing.T) {
		t.Run("with all fields complete and check if it is active", func(t *testing.T) {
			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testCtx.user.ID.String(), "provider")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.ProviderProfilePayload{
				TradingName:    "My Trading Name",
				Currency:       "KES",
				HostIdentifier: "example.com",
				Availability: types.ProviderAvailabilityPayload{
					Cadence:   "weekdays",
					StartTime: time.Now(),
					EndTime:   time.Now().Add(time.Hour * 24),
				},
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Profile updated successfully", response.Message)
			assert.Nil(t, response.Data, "response.Data is not nil")

			providerProfile, err := db.Client.ProviderProfile.
				Query().
				Where(providerprofile.HasUserWith(user.ID(testCtx.user.ID))).
				WithCurrency().
				Only(context.Background())
			assert.NoError(t, err)

			assert.Contains(t, providerProfile.TradingName, payload.TradingName)
			assert.Contains(t, providerProfile.HostIdentifier, payload.HostIdentifier)
			assert.Contains(t, providerProfile.Edges.Currency.Code, payload.Currency)
			assert.True(t, providerProfile.IsActive)
		})
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
		)
		assert.NoError(t, err)

		accessToken, _ := token.GenerateAccessJWT(testUser.ID.String(), "sender")
		headers := map[string]string{
			"Authorization": "Bearer " + accessToken,
		}
		res, err := test.PerformRequest(t, "GET", "/settings/sender", nil, headers, router)
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

	t.Run("UpdateProviderProfile", func(t *testing.T) {
		t.Run("with visibility", func(t *testing.T) {
			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testCtx.user.ID.String(), "provider")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.ProviderProfilePayload{
				VisibilityMode: "private",
				TradingName:    testCtx.providerProfile.TradingName,
				HostIdentifier: testCtx.providerProfile.HostIdentifier,
				Currency:       "KES",
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Profile updated successfully", response.Message)
			assert.Nil(t, response.Data, "response.Data is not nil")

			providerProfile, err := db.Client.ProviderProfile.Query().
				Where(providerprofile.VisibilityModeEQ(providerprofile.VisibilityModePrivate)).
				Count(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, 1, providerProfile)
		})
	})

}
