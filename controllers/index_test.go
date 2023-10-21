package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/routers/middleware"
	"github.com/paycrest/paycrest-protocol/services"
	db "github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/ent/enttest"
	"github.com/paycrest/paycrest-protocol/utils/test"
	"github.com/paycrest/paycrest-protocol/utils/token"
	"github.com/stretchr/testify/assert"
)

var testCtx = struct {
	user             *ent.User
	apiKey           *ent.APIKey
	validatorProfile *ent.ValidatorProfile
	apiKeySecret     string
}{}

func setup() error {
	// Set up test data
	user, _ := test.CreateTestUser(map[string]string{
		"scope": "tx_validator",
	})
	testCtx.user = user

	apiKeyService := services.NewAPIKeyService()
	apiKey, secretKey, err := apiKeyService.GenerateAPIKey(
		context.Background(),
		nil,
		user.ID,
	)
	if err != nil {
		return err
	}
	testCtx.apiKey = apiKey

	validator, err := test.CreateTestValidatorProfile(map[string]interface{}{
		"user_id": user.ID,
	})
	if err != nil {
		return err
	}
	testCtx.validatorProfile = validator

	testCtx.apiKeySecret = secretKey

	if _, err := test.CreateTestFiatCurrency(nil); err != nil {
		return err
	}

	return nil
}

func TestIndex(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setup()
	assert.NoError(t, err)

	// Set up test routers
	var ctrl Controller
	router := gin.New()
	router.Use(middleware.ScopeMiddleware)

	router.POST(
		"orders/:fulfillment_id/validate",
		middleware.HMACVerificationMiddleware,
		middleware.OnlyValidatorMiddleware,
		ctrl.ValidateOrder,
	)
	router.GET("currencies", ctrl.GetFiatCurrencies)

	t.Run("ValidateOrderFulfillment", func(t *testing.T) {

		t.Run("order is valid", func(t *testing.T) {
			// Test register with valid payload
			payload := map[string]interface{}{
				"isValid":   true,
				"errorMsg":  "",
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)
			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			// Get test lock order fulfillment
			fulfillment, err := test.CreateTestLockOrderFulfillment(nil)
			assert.NoError(t, err)

			res, err := test.PerformRequest(t, "POST",
				fmt.Sprintf("/orders/%s/validate?scope=validator", fulfillment.ID.String()), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)
			var response string
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "OK", response)
		})

		t.Run("invalid fulfillment ID", func(t *testing.T) {
			payload := map[string]interface{}{
				"isValid":   true,
				"errorMsg":  "",
				"timestamp": time.Now().Unix(),
			}

			fulfillmentID := "invalid"
			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)
			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			res, err := test.PerformRequest(t, "POST", "/orders/"+fulfillmentID+"/validate?scope=validator", payload, headers, router)
			assert.NoError(t, err)

			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.Contains(t, res.Body.String(), "Invalid fulfillment ID")
		})

		// t.Run("order is invalid", func(t *testing.T) {
		// 	payload := map[string]interface{}{
		// 		"isValid":   false,
		// 		"errorMsg":  "Invalid transaction reference",
		// 		"timestamp": time.Now().Unix(),
		// 	}

		// 	signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)
		// 	headers := map[string]string{
		// 		"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
		// 	}

		// 	// Get test lock order fulfillment
		// 	fulfillment, err := test.CreateTestLockOrderFulfillment(nil)
		// 	assert.NoError(t, err)

		// 	res, err := test.PerformRequest(t, "POST",
		// 		fmt.Sprintf("/orders/%s/validate?scope=validator", fulfillment.ID.String()), payload, headers, router)
		// 	assert.NoError(t, err)

		// 	// Assert the response body
		// 	assert.Equal(t, http.StatusOK, res.Code)
		// 	var response string
		// 	err = json.Unmarshal(res.Body.Bytes(), &response)
		// 	assert.NoError(t, err)
		// 	assert.Equal(t, "OK", response)

		// 	fulfillment, err = db.Client.LockOrderFulfillment.
		// 		Query().
		// 		Where(lockorderfulfillment.IDEQ(fulfillment.ID)).
		// 		Only(context.Background())
		// 	assert.NoError(t, err)

		// 	assert.Contains(t, fulfillment.ValidationErrors, "Invalid transaction reference")
		// })
	})

	t.Run("Currencies", func(t *testing.T) {
		t.Run("fetch supported fiat currencies", func(t *testing.T) {
			res, err := test.PerformRequest(t, "GET", "/currencies?scope=sender", nil, nil, router)
			assert.NoError(t, err)

			// Assert the response code.
			assert.Equal(t, http.StatusOK, res.Code)

			var response struct {
				Data    []types.SupportedCurrencies
				Message string
			}
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "OK", response.Message)

			// Assert /currencies response with the seeded Naira currency.
			nairaCurrency := types.SupportedCurrencies{
				Code:      "NGN",
				Name:      "Nigerian Naira",
				ShortName: "Naira",
				Decimals:  2,
				Symbol:    "₦",
			}
			assert.Equal(t, nairaCurrency, response.Data[0])
		})
	})
}
