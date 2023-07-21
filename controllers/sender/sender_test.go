package sender

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/routers/middleware"
	"github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/utils"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/ent/enttest"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/utils/test"
	"github.com/paycrest/paycrest-protocol/utils/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock indexer service
type MockIndexerService struct {
	mock.Mock
}

// IndexERC20Transfer mocks the IndexERC20Transfer method
func (m *MockIndexerService) IndexERC20Transfer(ctx context.Context, receiveAddress *ent.ReceiveAddress, done chan<- bool) error {
	// Call through to mock object's AssertCalled
	// args := m.Called(ctx, receiveAddress, done)
	done <- true
	return nil
}

var testCtx = struct {
	user         *ent.User
	apiKey       *ent.APIKey
	apiKeySecret string
}{}

func setup(client *ent.Client) error {
	// Set up test data
	user, err := test.CreateTestUser(db.Client, nil)
	if err != nil {
		return err
	}
	testCtx.user = user

	apiKeyService := services.NewAPIKeyService(db.Client)
	apiKey, secretKey, err := apiKeyService.GenerateAPIKey(
		context.Background(),
		user.ID,
		services.CreateAPIKeyPayload{
			Name:  "name",
			Scope: "sender",
		})
	testCtx.apiKey = apiKey
	if err != nil {
		return err
	}

	testCtx.apiKeySecret = secretKey

	return nil
}

func TestSender(t *testing.T) {

	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setup(db.Client)
	assert.NoError(t, err)

	// Set up test routers
	router := gin.New()
	router.Use(middleware.HMACVerificationMiddleware)

	// Create a mock instance of the IndexerService
	mockIndexerService := &MockIndexerService{}

	// Create a new instance of the SenderController with the mock service
	ctrl := NewSenderController(mockIndexerService)
	router.POST("/orders", ctrl.CreatePaymentOrder)

	t.Run("CreatePaymentOrder", func(t *testing.T) {
		payload := map[string]interface{}{
			"amount":  100.0,
			"token":   "USDT",
			"network": "bnb-smart-chain",
			"recipient": map[string]interface{}{
				"institution":       "First Bank Nigeria PLC",
				"accountIdentifier": "1234567890",
				"accountName":       "John Doe",
			},
			"timestamp": time.Now().Unix(),
		}

		signature := token.GenerateHMACSignature(payload, testCtx.apiKeySecret)

		headers := map[string]string{
			"Authorization": "HMAC " + testCtx.apiKey.ID.String() + ":" + signature,
		}

		res, err := test.PerformRequest(t, "POST", "/orders", payload, headers, router)
		time.Sleep(1 * time.Second)
		assert.NoError(t, err)

		// Assert the response body
		assert.Equal(t, http.StatusCreated, res.Code)

		var response utils.Response
		err = json.Unmarshal(res.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Payment order initiated successfully", response.Message)
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok, "response.Data is not of type map[string]interface{}")
		assert.NotNil(t, data, "response.Data is nil")

		assert.Equal(t, data["amount"], payload["amount"])
		assert.Equal(t, data["network"], payload["network"])

		// Parse the payment order ID string to uuid.UUID
		paymentOrderUUID, err := uuid.Parse(data["id"].(string))
		assert.NoError(t, err)

		// Query the database for the payment order
		_, err = db.Client.PaymentOrder.
			Query().
			Where(paymentorder.IDEQ(paymentOrderUUID)).
			Only(context.Background())
		assert.NoError(t, err)
	})
}
