package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/maphash"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/routers/middleware"
	"github.com/paycrest/protocol/services"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/utils/test"
	"github.com/paycrest/protocol/utils/token"
	"github.com/stretchr/testify/assert"
)

var testCtx = struct {
	user              *ent.SenderProfile
	token             *ent.Token
	apiKey            *ent.APIKey
	apiKeySecret      string
	client            types.RPCClient
	networkIdentifier string
}{}

// func createPaymentOrder(t *testing.T, router *gin.Engine) {
// 	// Fetch network from db
// 	network, err := db.Client.Network.
// 		Query().
// 		Where(network.IdentifierEQ(testCtx.networkIdentifier)).
// 		Only(context.Background())
// 	assert.NoError(t, err)

// 	r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))

// 	payload := map[string]interface{}{
// 		"amount":  "100",
// 		"token":   testCtx.token.Symbol,
// 		"rate":    "750",
// 		"network": network.Identifier,
// 		"recipient": map[string]interface{}{
// 			"institution":       "ABNGNGLA",
// 			"accountIdentifier": "1234567890",
// 			"accountName":       "John Doe",
// 			"memo":              "Shola Kehinde - rent for May 2021",
// 		},
// 		"label":     fmt.Sprintf("%d", r.Intn(100000)),
// 		"timestamp": time.Now().Unix(),
// 	}

// 	signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

// 	headers := map[string]string{
// 		"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
// 	}

// 	res, err := test.PerformRequest(t, "POST", "/orders", payload, headers, router)
// 	assert.NoError(t, err)

// 	// Assert the response body
// 	assert.Equal(t, http.StatusCreated, res.Code)

// 	var response types.Response
// 	err = json.Unmarshal(res.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Payment order initiated successfully", response.Message)
// 	data, ok := response.Data.(map[string]interface{})
// 	assert.True(t, ok, "response.Data is not of type map[string]interface{}")
// 	assert.NotNil(t, data, "response.Data is nil")

// 	assert.Equal(t, data["amount"], payload["amount"])
// 	assert.Equal(t, data["network"], payload["network"])
// 	assert.NotEmpty(t, data["validUntil"])

// 	// Parse the payment order ID string to uuid.UUID
// 	paymentOrderUUID, err := uuid.Parse(data["id"].(string))
// 	assert.NoError(t, err)

// 	// Query the database for the payment order
// 	paymentOrder, err := db.Client.PaymentOrder.
// 		Query().
// 		Where(paymentorder.IDEQ(paymentOrderUUID)).
// 		WithRecipient().
// 		Only(context.Background())
// 	assert.NoError(t, err)

// 	assert.NotNil(t, paymentOrder.Edges.Recipient)
// 	assert.Equal(t, paymentOrder.Edges.Recipient.AccountIdentifier, payload["recipient"].(map[string]interface{})["accountIdentifier"])
// 	assert.Equal(t, paymentOrder.Edges.Recipient.Memo, payload["recipient"].(map[string]interface{})["memo"])
// 	assert.Equal(t, paymentOrder.Edges.Recipient.AccountName, payload["recipient"].(map[string]interface{})["accountName"])
// 	assert.Equal(t, paymentOrder.Edges.Recipient.Institution, payload["recipient"].(map[string]interface{})["institution"])
// 	assert.Equal(t, data["senderFee"], "0.666667")
// 	assert.Equal(t, data["transactionFee"], network.Fee.Add(paymentOrder.Amount.Mul(decimal.NewFromFloat(0.001))).String()) // 0.1% protocol fee
// }

func setup() error {
	// Set up test data
	user, err := test.CreateTestUser(nil)
	if err != nil {
		return err
	}

	senderProfile, err := test.CreateTestSenderProfile(map[string]interface{}{
		"user_id":            user.ID,
		"fee_per_token_unit": "5",
	})
	if err != nil {
		return err
	}
	testCtx.user = senderProfile

	apiKeyService := services.NewAPIKeyService()
	apiKey, secretKey, err := apiKeyService.GenerateAPIKey(
		context.Background(),
		nil,
		senderProfile,
		nil,
	)
	if err != nil {
		return err
	}
	testCtx.apiKey = apiKey

	// Set up test blockchain client
	backend, err := test.SetUpTestBlockchain()
	if err != nil {
		return err
	}

	// Create a test token
	testCtx.networkIdentifier = "localhost" + uuid.New().String()
	token, err := test.CreateERC20Token(backend, map[string]interface{}{
		"identifier": testCtx.networkIdentifier,
	})
	if err != nil {
		return err
	}
	testCtx.token = token
	testCtx.client = backend

	testCtx.apiKeySecret = secretKey

	for i := 0; i < 9; i++ {
		receiveAddress, err := test.CreateSmartAccount(
			context.Background(), backend)
		if err != nil {
			return err
		}
		test.CreateTestPaymentOrder(backend, token, map[string]interface{}{
			"sender":         senderProfile,
			"receiveAddress": receiveAddress,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestSender(t *testing.T) {

	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setup()
	assert.NoError(t, err)

	// Set up test routers
	router := gin.New()
	router.Use(middleware.HMACVerificationMiddleware)
	router.Use(middleware.OnlySenderMiddleware)

	// Create a mock instance of the OrderService
	mockOrderService := &test.MockOrderService{}

	// Create a mock instance of the IndexerService
	mockIndexerService := &test.MockIndexerService{}

	// Create a new instance of the SenderController with the mock service
	ctrl := NewSenderController(mockIndexerService, mockOrderService)
	router.POST("/orders", ctrl.InitiatePaymentOrder)
	router.GET("/orders/:id", ctrl.GetPaymentOrderByID)
	router.GET("/orders/", ctrl.GetPaymentOrders)
	router.GET("/stats", ctrl.Stats)

	var paymentOrderUUID uuid.UUID

	t.Run("InitiatePaymentOrder", func(t *testing.T) {

		// Fetch network from db
		network, err := db.Client.Network.
			Query().
			Where(network.IdentifierEQ(testCtx.networkIdentifier)).
			Only(context.Background())
		assert.NoError(t, err)

		r := rand.New(rand.NewSource(int64(new(maphash.Hash).Sum64())))

		payload := map[string]interface{}{
			"amount":  "100",
			"token":   testCtx.token.Symbol,
			"rate":    "750",
			"network": network.Identifier,
			"recipient": map[string]interface{}{
				"institution":       "ABNGNGLA",
				"accountIdentifier": "1234567890",
				"accountName":       "John Doe",
				"memo":              "Shola Kehinde - rent for May 2021",
			},
			"label":     fmt.Sprintf("%d", r.Intn(100000)),
			"timestamp": time.Now().Unix(),
		}

		signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

		headers := map[string]string{
			"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
		}

		res, err := test.PerformRequest(t, "POST", "/orders", payload, headers, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusCreated, res.Code)

		var response types.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Payment order initiated successfully", response.Message)
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "response.Data is not of type map[string]interface{}")
		assert.NotNil(t, data, "response.Data is nil")

		assert.Equal(t, data["amount"], payload["amount"])
		assert.Equal(t, data["network"], payload["network"])
		assert.NotEmpty(t, data["validUntil"])

		// Parse the payment order ID string to uuid.UUID
		paymentOrderUUID, err = uuid.Parse(data["id"].(string))
		assert.NoError(t, err)

		// Query the database for the payment order
		paymentOrder, err := db.Client.PaymentOrder.
			Query().
			Where(paymentorder.IDEQ(paymentOrderUUID)).
			WithRecipient().
			Only(context.Background())
		assert.NoError(t, err)

		assert.NotNil(t, paymentOrder.Edges.Recipient)
		assert.Equal(t, paymentOrder.Edges.Recipient.AccountIdentifier, payload["recipient"].(map[string]interface{})["accountIdentifier"])
		assert.Equal(t, paymentOrder.Edges.Recipient.Memo, payload["recipient"].(map[string]interface{})["memo"])
		assert.Equal(t, paymentOrder.Edges.Recipient.AccountName, payload["recipient"].(map[string]interface{})["accountName"])
		assert.Equal(t, paymentOrder.Edges.Recipient.Institution, payload["recipient"].(map[string]interface{})["institution"])
		assert.Equal(t, data["senderFee"], "0.666667")
		assert.Equal(t, data["transactionFee"], network.Fee.Add(paymentOrder.Amount.Mul(decimal.NewFromFloat(0.001))).String()) // 0.1% protocol fee
	})

	t.Run("GetPaymentOrderByID", func(t *testing.T) {
		var payload = map[string]interface{}{
			"timestamp": time.Now().Unix(),
		}

		signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

		headers := map[string]string{
			"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
		}

		res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders/%s", paymentOrderUUID), payload, headers, router)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusOK, res.Code)

		var response types.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "The order has been successfully retrieved", response.Message)
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "response.Data is of not type map[string]interface{}")
		assert.NotNil(t, data, "response.Data is nil")
	})

	t.Run("GetPaymentOrders", func(t *testing.T) {
		t.Run("fetch default list", func(t *testing.T) {
			// Test default params
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			res, err := test.PerformRequest(t, "GET", "/orders/", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Payment orders retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, int(data["page"].(float64)), 1)
			assert.Equal(t, int(data["pageSize"].(float64)), 10) // default pageSize
			assert.NotEmpty(t, data["total"])
			assert.NotEmpty(t, data["orders"])
		})

		t.Run("when filtering is applied", func(t *testing.T) {
			// Test different status filters
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			//query params
			status := "initiated"

			res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders/?status=%s", status), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Payment orders retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, int(data["page"].(float64)), 1)
			assert.Equal(t, int(data["pageSize"].(float64)), 10) // default pageSize
			assert.NotEmpty(t, data["total"])
			assert.NotEmpty(t, data["orders"])
		})

		t.Run("with custom page and pageSize", func(t *testing.T) {
			// Test different page and pageSize values
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			//query params
			page := 1
			pageSize := 10

			res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders/?page=%s&pageSize=%s", strconv.Itoa(page), strconv.Itoa(pageSize)), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Payment orders retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, int(data["page"].(float64)), page)
			assert.Equal(t, int(data["pageSize"].(float64)), pageSize)
			assert.Equal(t, 10, len(data["orders"].([]interface{})))
			assert.NotEmpty(t, data["total"])
			assert.NotEmpty(t, data["orders"])
		})

		t.Run("with ordering", func(t *testing.T) {
			// Test ascending and descending ordering
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			//query params
			ordering := "desc"

			res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders/?ordering=%s", ordering), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Payment orders retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Try to parse the first and last order time strings using a set of predefined layouts
			firstOrderTimestamp, err := time.Parse(time.RFC3339Nano, data["orders"].([]interface{})[0].(map[string]interface{})["createdAt"].(string))
			if err != nil {
				return
			}

			lastOrderTimestamp, err := time.Parse(time.RFC3339Nano, data["orders"].([]interface{})[len(data["orders"].([]interface{}))-1].(map[string]interface{})["createdAt"].(string))
			if err != nil {
				return
			}

			assert.Equal(t, int(data["page"].(float64)), 1)
			assert.Equal(t, int(data["pageSize"].(float64)), 10) // default pageSize
			assert.NotEmpty(t, data["total"])
			assert.NotEmpty(t, data["orders"])
			assert.Greater(t, len(data["orders"].([]interface{})), 0)
			assert.GreaterOrEqual(t, firstOrderTimestamp, lastOrderTimestamp)
		})

		t.Run("with filtering by network", func(t *testing.T) {
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			//query params
			network := testCtx.networkIdentifier

			res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders/?network=%s", network), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Payment orders retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.NotEmpty(t, data["total"])
			assert.NotEmpty(t, data["orders"])
			assert.Greater(t, len(data["orders"].([]interface{})), 0)

			for _, order := range data["orders"].([]interface{}) {
				assert.Equal(t, order.(map[string]interface{})["network"], network)
			}
		})

		t.Run("with filtering by token", func(t *testing.T) {
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			//query params
			token := testCtx.token.Symbol

			res, err := test.PerformRequest(t, "GET", fmt.Sprintf("/orders/?token=%s", token), payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Payment orders retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.NotEmpty(t, data["total"])
			assert.NotEmpty(t, data["orders"])
			assert.Greater(t, len(data["orders"].([]interface{})), 0)

			for _, order := range data["orders"].([]interface{}) {
				assert.Equal(t, order.(map[string]interface{})["token"], token)
			}
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

			senderProfile, err := test.CreateTestSenderProfile(map[string]interface{}{
				"user_id":            user.ID,
				"fee_per_token_unit": "5",
			})
			if err != nil {
				return
			}

			apiKeyService := services.NewAPIKeyService()
			apiKey, secretKey, err := apiKeyService.GenerateAPIKey(
				context.Background(),
				nil,
				senderProfile,
				nil,
			)
			if err != nil {
				return
			}

			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, secretKey)

			headers := map[string]string{
				"Authorization": "HMAC " + apiKey.ID.String() + ":" + signature,
			}

			res, err := test.PerformRequest(t, "GET", "/stats", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Sender stats retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			assert.Equal(t, int(data["totalOrders"].(float64)), 0)

			totalOrderVolumeStr, ok := data["totalOrderVolume"].(string)
			assert.True(t, ok, "totalOrderVolume is not of type string")
			totalOrderVolume, err := decimal.NewFromString(totalOrderVolumeStr)
			assert.NoError(t, err, "Failed to convert totalOrderVolume to decimal")
			assert.Equal(t, totalOrderVolume, decimal.NewFromInt(0))

			totalFeeEarningsStr, ok := data["totalFeeEarnings"].(string)
			assert.True(t, ok, "totalFeeEarnings is not of type string")
			totalFeeEarnings, err := decimal.NewFromString(totalFeeEarningsStr)
			assert.NoError(t, err, "Failed to convert totalFeeEarnings to decimal")
			assert.Equal(t, totalFeeEarnings, decimal.NewFromInt(0))
		})

		t.Run("when orders have been initiated", func(t *testing.T) {
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			res, err := test.PerformRequest(t, "GET", "/stats", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Sender stats retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the totalOrders value
			totalOrders, ok := data["totalOrders"].(float64)
			assert.True(t, ok, "totalOrders is not of type float64")
			assert.Equal(t, 10, int(totalOrders))

			// Assert the totalOrderVolume value
			totalOrderVolumeStr, ok := data["totalOrderVolume"].(string)
			assert.True(t, ok, "totalOrderVolume is not of type string")
			totalOrderVolume, err := decimal.NewFromString(totalOrderVolumeStr)
			assert.NoError(t, err, "Failed to convert totalOrderVolume to decimal")
			assert.Equal(t, 0, totalOrderVolume.Cmp(decimal.NewFromInt(0)))

			// Assert the totalFeeEarnings value
			totalFeeEarningsStr, ok := data["totalFeeEarnings"].(string)
			assert.True(t, ok, "totalFeeEarnings is not of type string")
			totalFeeEarnings, err := decimal.NewFromString(totalFeeEarningsStr)
			assert.NoError(t, err, "Failed to convert totalFeeEarnings to decimal")
			assert.Equal(t, 0, totalFeeEarnings.Cmp(decimal.NewFromInt(0)))
		})

		t.Run("should only calculate volumes of settled orders", func(t *testing.T) {

			receiveAddress, err := test.CreateSmartAccount(
				context.Background(), testCtx.client)
			assert.NoError(t, err)

			// create settled Order
			_, err = test.CreateTestPaymentOrder(testCtx.client, testCtx.token, map[string]interface{}{
				"sender":             testCtx.user,
				"amount":             100.0,
				"token":              testCtx.token.Symbol,
				"rate":               750.0,
				"status":             "settled",
				"fee_per_token_unit": 5.0,
				"receiveAddress":     receiveAddress,
			})
			assert.NoError(t, err)
			var payload = map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}

			signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

			headers := map[string]string{
				"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
			}

			res, err := test.PerformRequest(t, "GET", "/stats", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Sender stats retrieved successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is of not type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the totalOrders value
			totalOrders, ok := data["totalOrders"].(float64)
			assert.True(t, ok, "totalOrders is not of type float64")
			assert.Equal(t, 11, int(totalOrders))

			// Assert the totalOrderVolume value
			totalOrderVolumeStr, ok := data["totalOrderVolume"].(string)
			assert.True(t, ok, "totalOrderVolume is not of type string")
			totalOrderVolume, err := decimal.NewFromString(totalOrderVolumeStr)
			assert.NoError(t, err, "Failed to convert totalOrderVolume to decimal")
			assert.Equal(t, 0, totalOrderVolume.Cmp(decimal.NewFromInt(100)))

			// Assert the totalFeeEarnings value
			totalFeeEarningsStr, ok := data["totalFeeEarnings"].(string)
			assert.True(t, ok, "totalFeeEarnings is not of type string")
			totalFeeEarnings, err := decimal.NewFromString(totalFeeEarningsStr)
			assert.NoError(t, err, "Failed to convert totalFeeEarnings to decimal")
			assert.Equal(t, 0, totalFeeEarnings.Cmp(decimal.NewFromFloat(0.666667)))
		})
	})
}
