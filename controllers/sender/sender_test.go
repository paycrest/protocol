package sender

import (
	"context"
	"encoding/json"
	"fmt"
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

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/utils/test"
	"github.com/paycrest/protocol/utils/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock indexer service
type MockIndexerService struct {
	mock.Mock
}

// IndexERC20Transfer mocks the IndexERC20Transfer method
func (m *MockIndexerService) IndexERC20Transfer(ctx context.Context, client types.RPCClient, receiveAddress *ent.ReceiveAddress, done chan<- bool) error {
	done <- true
	return nil
}

// IndexOrderDeposits mocks the IndexOrderDeposits method
func (m *MockIndexerService) IndexOrderDeposits(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	return nil
}

// IndexOrderSettlements mocks the IndexOrderSettlements method
func (m *MockIndexerService) IndexOrderSettlements(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	return nil
}

// IndexOrderRefunds mocks the IndexOrderRefunds method
func (m *MockIndexerService) IndexOrderRefunds(ctx context.Context, client types.RPCClient, network *ent.Network) error {
	return nil
}

var testCtx = struct {
	user              *ent.User
	token             *ent.Token
	apiKey            *ent.APIKey
	apiKeySecret      string
	networkIdentifier string
	paymentOrderUUID  uuid.UUID
}{}

func createPaymentOrder(t *testing.T, router *gin.Engine) {
	// Fetch network from db
	network, err := db.Client.Network.
		Query().
		Where(network.IdentifierEQ(testCtx.networkIdentifier)).
		Only(context.Background())
	assert.NoError(t, err)

	payload := map[string]interface{}{
		"amount":  100.0,
		"token":   testCtx.token.Symbol,
		"rate":    750.0,
		"network": network.Identifier,
		"recipient": map[string]interface{}{
			"institution":       "First Bank Nigeria PLC",
			"accountIdentifier": "1234567890",
			"accountName":       "John Doe",
			"memo":              "Shola Kehinde - rent for May 2021",
		},
		"label":     "cc334ncldfa",
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
	paymentOrderUUID, err := uuid.Parse(data["id"].(string))
	assert.NoError(t, err)

	// Query the database for the payment order
	paymentOrder, err := db.Client.PaymentOrder.
		Query().
		Where(paymentorder.IDEQ(paymentOrderUUID)).
		WithRecipient().
		Only(context.Background())
	assert.NoError(t, err)

	testCtx.paymentOrderUUID = paymentOrderUUID
	assert.NotNil(t, paymentOrder.Edges.Recipient)
	assert.Equal(t, paymentOrder.Edges.Recipient.AccountIdentifier, payload["recipient"].(map[string]interface{})["accountIdentifier"])
	assert.Equal(t, paymentOrder.Edges.Recipient.Memo, payload["recipient"].(map[string]interface{})["memo"])
	assert.Equal(t, paymentOrder.Edges.Recipient.AccountName, payload["recipient"].(map[string]interface{})["accountName"])
	assert.Equal(t, paymentOrder.Edges.Recipient.Institution, payload["recipient"].(map[string]interface{})["institution"])
	assert.Equal(t, data["senderFee"], "0.2")
	assert.Equal(t, data["networkFee"], "1")
}

func setup() error {
	// Set up test data
	user, err := test.CreateTestUser(nil)
	if err != nil {
		return err
	}
	testCtx.user = user

	senderProfile, err := test.CreateTestSenderProfile(map[string]interface{}{
		"user_id":            user.ID,
		"fee_per_token_unit": 1.5,
	})
	if err != nil {
		return err
	}

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
	backend, err := test.NewSimulatedBlockchain()
	if err != nil {
		return err
	}

	// Create a test token
	testCtx.networkIdentifier = "polygon-mumbai" + uuid.New().String()
	token, err := test.CreateTestToken(backend, map[string]interface{}{
		"identifier": testCtx.networkIdentifier,
	})
	if err != nil {
		return err
	}
	testCtx.token = token

	testCtx.apiKeySecret = secretKey

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

	// Create a mock instance of the IndexerService
	mockIndexerService := &MockIndexerService{}

	// Create a new instance of the SenderController with the mock service
	ctrl := NewSenderController(mockIndexerService)
	router.POST("/orders", ctrl.CreatePaymentOrder)
	router.GET("/orders/:id", ctrl.GetPaymentOrderByID)
	router.GET("/orders/", ctrl.GetPaymentOrders)

	var paymentOrderUUID uuid.UUID

	t.Run("CreatePaymentOrder", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			createPaymentOrder(t, router)
		}
		paymentOrderUUID = testCtx.paymentOrderUUID
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
			assert.Equal(t, len(data["orders"].([]interface{})), pageSize)
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
			assert.Greater(t, firstOrderTimestamp, lastOrderTimestamp)
		})
	})
}
