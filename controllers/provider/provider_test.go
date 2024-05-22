package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/routers/middleware"
	"github.com/paycrest/protocol/services"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/utils/test"
	"github.com/paycrest/protocol/utils/token"
	"github.com/stretchr/testify/assert"
)

var testCtx = struct {
	user         *ent.User
	provider     *ent.ProviderProfile
	apiKey       *ent.APIKey
	apiKeySecret string
	lockOrder    *ent.LockPaymentOrder
}{}

func setup() error {
	// Set up test data
	user, err := test.CreateTestUser(map[string]interface{}{
		"scope": "provider"})
	if err != nil {
		return err
	}
	testCtx.user = user

	currency, err := test.CreateTestFiatCurrency(nil)
	if err != nil {
		return err
	}

	providerProfile, err := test.CreateTestProviderProfile(map[string]interface{}{
		"user_id":     testCtx.user.ID,
		"currency_id": currency.ID,
	})
	if err != nil {
		return err
	}
	testCtx.provider = providerProfile

	for i := 0; i < 10; i++ {
		_, err := test.CreateTestLockPaymentOrder(map[string]interface{}{
			"gateway_id": uuid.New().String(),
			"provider":   providerProfile,
		})
		if err != nil {
			return err
		}
	}

	apiKeyService := services.NewAPIKeyService()
	apiKey, secretKey, err := apiKeyService.GenerateAPIKey(
		context.Background(),
		nil,
		nil,
		providerProfile,
	)
	if err != nil {
		return err
	}

	testCtx.apiKey = apiKey
	testCtx.apiKeySecret = secretKey

	return nil
}

func TestProvider(t *testing.T) {

	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setup()
	assert.NoError(t, err)

	// Set up test routers
	router := gin.New()
	router.Use(middleware.DynamicAuthMiddleware)
	router.Use(middleware.OnlyProviderMiddleware)

	// Create a new instance of the SenderController with the mock service
	ctrl := NewProviderController()
	router.GET("/orders", ctrl.GetLockPaymentOrders)
	router.GET("/stats", ctrl.Stats)
	router.GET("/node-info", ctrl.NodeInfo)

	t.Run("GetLockPaymentOrders", func(t *testing.T) {
		t.Run("fetch default list", func(t *testing.T) {
			// Test default params
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
				"Client-Type":   "backend",
			}

			res, err := test.PerformRequest(t, "GET", "/orders", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Orders successfully retrieved", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, int(data["page"].(float64)), 1)
			assert.Equal(t, int(data["pageSize"].(float64)), 10) // default pageSize
			assert.NotNil(t, data["total"])
			assert.NotEmpty(t, data["orders"])
			assert.Greater(t, len(data["orders"].([]interface{})), 0)

		})

		t.Run("when filtering is applied", func(t *testing.T) {
			// Test different status filters
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
				"Client-Type":   "backend",
			}

			//query params
			status := "pending"

			res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders?status=%s", status), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Orders successfully retrieved", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, int(data["page"].(float64)), 1)
			assert.Equal(t, int(data["pageSize"].(float64)), 10) // default pageSize
			assert.NotNil(t, data["total"])
			assert.NotEmpty(t, data["orders"])
			assert.Greater(t, len(data["orders"].([]interface{})), 0)

		})

		t.Run("with custom page and pageSize", func(t *testing.T) {
			// Test different page and pageSize values
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
				"Client-Type":   "backend",
			}

			//query params
			page := 1
			pageSize := 5

			res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders?page=%s&pageSize=%s", strconv.Itoa(page), strconv.Itoa(pageSize)), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Orders successfully retrieved", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, int(data["page"].(float64)), page)
			assert.Equal(t, int(data["pageSize"].(float64)), pageSize)
			assert.NotNil(t, data["total"])
			assert.NotEmpty(t, data["orders"])
			assert.Equal(t, len(data["orders"].([]interface{})), pageSize)
			assert.Greater(t, len(data["orders"].([]interface{})), 0)

		})

		t.Run("with ordering", func(t *testing.T) {
			// Test ascending and descending ordering
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
				"Client-Type":   "backend",
			}

			//query params
			ordering := "desc"

			res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders?ordering=%s", ordering), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Orders successfully retrieved", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Try to parse the first and last order time strings using a set of predefined layouts
			firstOrderTimestamp, err := time.Parse(
				time.RFC3339Nano,
				data["orders"].([]interface{})[0].(map[string]interface{})["createdAt"].(string),
			)
			if err != nil {
				return
			}

			lastOrderTimestamp, err := time.Parse(
				time.RFC3339Nano,
				data["orders"].([]interface{})[len(data["orders"].([]interface{}))-1].(map[string]interface{})["createdAt"].(string),
			)
			if err != nil {
				return
			}

			assert.Equal(t, int(data["page"].(float64)), 1)
			assert.Equal(t, int(data["pageSize"].(float64)), 10) // default pageSize
			assert.NotNil(t, data["total"])
			assert.NotEmpty(t, data["orders"])
			assert.Greater(t, len(data["orders"].([]interface{})), 0)
			assert.Greater(t, firstOrderTimestamp, lastOrderTimestamp)
		})

	})

	t.Run("GetStats", func(t *testing.T) {
		t.Run("when no orders have been initiated", func(t *testing.T) {
			// Create a new user with no orders
			user, err := test.CreateTestUser(map[string]interface{}{
				"email": "no_order_user@test.com",
			})
			if err != nil {
				return
			}

			currency, err := test.CreateTestFiatCurrency(nil)
			if err != nil {
				return
			}

			providerProfile, err := test.CreateTestProviderProfile(map[string]interface{}{
				"user_id":     user.ID,
				"currency_id": currency.ID,
			})
			if err != nil {
				return
			}

			apiKeyService := services.NewAPIKeyService()
			apiKey, secretKey, err := apiKeyService.GenerateAPIKey(
				context.Background(),
				nil,
				nil,
				providerProfile,
			)
			if err != nil {
				return
			}

			// Test default params
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, secretKey)

			headers := map[string]string{
				"Authorization": "HMAC " + apiKey.ID.String() + ":" + signature,
				"Client-Type":   "backend",
			}

			res, err := test.PerformRequest(t, "GET", "/stats", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Provider stats fetched successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, int(data["totalOrders"].(float64)), 0)

			totalFiatVolumeStr, ok := data["totalFiatVolume"].(string)
			assert.True(t, ok, "totalFiatVolume is not of type string")
			totalFiatVolume, err := decimal.NewFromString(totalFiatVolumeStr)
			assert.NoError(t, err, "Failed to convert totalFiatVolume to decimal")
			assert.Equal(t, totalFiatVolume, decimal.NewFromInt(0))

			totalCryptoVolumeStr, ok := data["totalCryptoVolume"].(string)
			assert.True(t, ok, "totalCryptoVolume is not of type string")
			totalCryptoVolume, err := decimal.NewFromString(totalCryptoVolumeStr)
			assert.NoError(t, err, "Failed to convert totalCryptoVolume to decimal")
			assert.Equal(t, totalCryptoVolume, decimal.NewFromInt(0))
		})

		t.Run("when orders have been initiated", func(t *testing.T) {
			// Test default params
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
				"Client-Type":   "backend",
			}

			res, err := test.PerformRequest(t, "GET", "/stats", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Provider stats fetched successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type *types.ProviderStatsResponse")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the totalOrders value
			totalOrders, ok := data["totalOrders"].(float64)
			assert.True(t, ok, "totalOrders is not of type float64")
			assert.Equal(t, 10, int(totalOrders))

			// Assert the totalFiatVolume value
			totalFiatVolumeStr, ok := data["totalFiatVolume"].(string)
			assert.True(t, ok, "totalFiatVolume is not of type string")
			totalFiatVolume, err := decimal.NewFromString(totalFiatVolumeStr)
			assert.NoError(t, err, "Failed to convert totalFiatVolume to decimal")
			assert.Equal(t, totalFiatVolume.Cmp(decimal.NewFromInt(0)), 0)

			// Assert the totalCryptoVolume value
			totalCryptoVolumeStr, ok := data["totalCryptoVolume"].(string)
			assert.True(t, ok, "totalCryptoVolume is not of type string")
			totalCryptoVolume, err := decimal.NewFromString(totalCryptoVolumeStr)
			assert.NoError(t, err, "Failed to convert totalCryptoVolume to decimal")
			assert.Equal(t, totalCryptoVolume.Cmp(decimal.NewFromInt(0)), 0)
		})

		t.Run("when orders have been settled", func(t *testing.T) {
			amount:=200.5
			_, err := test.CreateTestLockPaymentOrder(map[string]interface{}{
				"amount":    amount,
				"gateway_id": uuid.New().String(),
				"provider":   testCtx.provider,
				"status":     "settled",
			})
			assert.NoError(t, err)

			// Test default params
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
				"Client-Type":   "backend",
			}

			res, err := test.PerformRequest(t, "GET", "/stats", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Provider stats fetched successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type *types.ProviderStatsResponse")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the totalOrders value
			totalOrders, ok := data["totalOrders"].(float64)
			assert.True(t, ok, "totalOrders is not of type float64")
			assert.Equal(t, 11, int(totalOrders))

			// Assert the totalFiatVolume value
			totalFiatVolumeStr, ok := data["totalFiatVolume"].(string)
			assert.True(t, ok, "totalFiatVolume is not of type string")
			totalFiatVolume, err := decimal.NewFromString(totalFiatVolumeStr)
			assert.NoError(t, err, "Failed to convert totalFiatVolume to decimal")
			assert.Equal(t, totalFiatVolume.Cmp(decimal.NewFromInt(150375)), 0)

			// Assert the totalCryptoVolume value
			totalCryptoVolumeStr, ok := data["totalCryptoVolume"].(string)
			assert.True(t, ok, "totalCryptoVolume is not of type string")
			totalCryptoVolume, err := decimal.NewFromString(totalCryptoVolumeStr)
			assert.NoError(t, err, "Failed to convert totalCryptoVolume to decimal")
			assert.Equal(t, totalCryptoVolume.Cmp(decimal.NewFromFloat(amount)), 0)
		})
	})

	t.Run("NodeInfo", func(t *testing.T) {

		t.Run("when node is healthy", func(t *testing.T) {
			// Activate httpmock
			httpmock.Activate()
			defer httpmock.Deactivate()

			// Register mock response
			httpmock.RegisterResponder("GET", "https://example.com/health",
				func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(200, map[string]interface{}{
						"status":  "success",
						"message": "Node is live",
						"data": map[string]interface{}{
							"currency": "NGN",
						},
					})
				},
			)

			// Test default params
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			res, err := test.PerformRequest(t, "GET", "/node-info", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Node info fetched successfully", response.Message)
		})

		t.Run("when node is unhealthy", func(t *testing.T) {
			// Activate httpmock
			httpmock.Activate()
			defer httpmock.Deactivate()

			// Register mock response
			httpmock.RegisterResponder("GET", "https://example.com/health",
				func(r *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(503, nil)
				},
			)

			// Test default params
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			res, err := test.PerformRequest(t, "GET", "/node-info", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusServiceUnavailable, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to fetch node info", response.Message)
		})
	})
}
