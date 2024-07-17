package tasks

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/enttest"
	"github.com/paycrest/protocol/ent/webhookretryattempt"
	db "github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/test"
	"github.com/stretchr/testify/assert"
)

var testCtx = struct {
	user    *ent.SenderProfile
	webhook *ent.WebhookRetryAttempt
}{}

func setup() error {
	// Set up test data
	user, err := test.CreateTestUser(map[string]interface{}{
		"email": "chibie@paycrest.io",
	})
	if err != nil {
		return err
	}

	// Set up test blockchain client
	backend, err := test.SetUpTestBlockchain()
	if err != nil {
		return err
	}

	// Create a test token
	token, err := test.CreateERC20Token(backend, map[string]interface{}{
		"identifier":     "localhost",
		"deployContract": false,
	})
	if err != nil {
		return fmt.Errorf("CreateERC20Token.task_test: %w", err)
	}

	senderProfile, err := test.CreateTestSenderProfile(map[string]interface{}{
		"user_id":            user.ID,
		"fee_per_token_unit": "5",
	})

	if err != nil {
		return fmt.Errorf("CreateTestSenderProfile.task_test: %w", err)
	}
	testCtx.user = senderProfile

	paymentOrder, err := test.CreateTestPaymentOrder(backend, token, map[string]interface{}{
		"sender": senderProfile,
	})

	// Create the payload
	payloadStruct := types.PaymentOrderWebhookPayload{
		Event: "Test_events",
		Data: types.PaymentOrderWebhookData{
			ID:             paymentOrder.ID,
			Amount:         paymentOrder.Amount,
			AmountPaid:     paymentOrder.AmountPaid,
			AmountReturned: paymentOrder.AmountReturned,
			PercentSettled: paymentOrder.PercentSettled,
			SenderFee:      paymentOrder.SenderFee,
			NetworkFee:     paymentOrder.NetworkFee,
			Rate:           paymentOrder.Rate,
			Network:        token.Edges.Network.Identifier,
			GatewayID:      paymentOrder.GatewayID,
			SenderID:       senderProfile.ID,
			Recipient: types.PaymentOrderRecipient{
				Institution:       "",
				AccountIdentifier: "",
				AccountName:       "021",
				ProviderID:        "",
				Memo:              "",
			},
			FromAddress:   paymentOrder.FromAddress,
			ReturnAddress: paymentOrder.ReturnAddress,
			UpdatedAt:     paymentOrder.UpdatedAt,
			CreatedAt:     paymentOrder.CreatedAt,
			TxHash:        paymentOrder.TxHash,
			Status:        paymentOrder.Status,
		},
	}
	payload := utils.StructToMap(payloadStruct)
	hook, err := db.Client.WebhookRetryAttempt.
		Create().
		SetAttemptNumber(3).
		SetNextRetryTime(time.Now().Add(25 * time.Hour)).
		SetPayload(payload).
		SetSignature("").
		SetWebhookURL(senderProfile.WebhookURL).
		SetNextRetryTime(time.Now().Add(-10 * time.Minute)).
		SetCreatedAt(time.Now().Add(-25 * time.Hour)).
		SetStatus(webhookretryattempt.StatusFailed).
		Save(context.Background())

	testCtx.webhook = hook
	if err != nil {
		return fmt.Errorf("CreateTestSenderProfile.WebhookRetryAttempt: %w", err)
	}

	return nil
}
func TestTask(t *testing.T) {

	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	// httpmock.Activate()
	// httpmock.Deactivate()

	// // Register mock response
	// httpmock.RegisterResponder("POST", testCtx.user.WebhookURL,
	// 	func(r *http.Request) (*http.Response, error) {
	// 		return httpmock.NewBytesResponse(400, []byte(`{"id": "01", "message": "Sent"}`)), nil
	// 	},
	// )

	db.Client = client

	// Setup test data
	err := setup()
	assert.NoError(t, err)
	t.Run("RetryFailedWebhookNotifications", func(t *testing.T) {
		err := RetryFailedWebhookNotifications()
		assert.NoError(t, err)
		hook, err := db.Client.WebhookRetryAttempt.
			Query().
			Where(webhookretryattempt.IDEQ(testCtx.webhook.ID)).
			Only(context.Background())
		assert.Equal(t, hook.Status, webhookretryattempt.StatusExpired)
	})
}
