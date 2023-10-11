package services

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var (
	mailgunEndpoint = fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", notificationConf.EmailDomain)

	testToken = "test-token"
	testEmail = "test@paycrest.io"
)

func TestEmailService(t *testing.T) {
	// activate httpmock
	httpmock.Activate()
	defer httpmock.Deactivate()

	// register mock response
	httpmock.RegisterResponder("POST", mailgunEndpoint,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewBytesResponse(200, []byte(`{"id": "01", "message": "Sent"}`)), nil
		},
	)

	t.Run("Mailgun", func(t *testing.T) {

		t.Run("SendVerificationEmail should work properly and return a response payload", func(t *testing.T) {
			srv := NewEmailService(MAILGUN_MAIL_PROVIDER)

			response, err := srv.SendVerificationEmail(context.TODO(), testToken, testEmail)

			// error checker.
			assert.NoError(t, err, "unexpected error")

			// assert the test token was sent.
			assert.NotEmpty(t, response.Id, "response ID should not be empty")
		})
	})
}
