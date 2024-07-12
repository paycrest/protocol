package services

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/paycrest/protocol/types"
	"github.com/stretchr/testify/assert"
)

var (
	mailgunEndpoint  = fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", notificationConf.EmailDomain)
	sendGridEndpoint = "https://api.sendgrid.com/v3/mail/send"

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

	httpmock.RegisterResponder("POST", sendGridEndpoint,
		func(r *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(202, nil)
			resp.Header.Set("X-Message-Id", "thisisatestid")
			return resp, nil
		},
	)

	// t.Run("Mailgun", func(t *testing.T) {

	// 	t.Run("SendVerificationEmail should work properly and return a response payload", func(t *testing.T) {
	// 		srv := NewEmailService(MAILGUN_MAIL_PROVIDER)

	// 		response, err := srv.SendVerificationEmail(context.Background(), testToken, testEmail)

	// 		// error checker.
	// 		assert.NoError(t, err, "unexpected error")

	// 		// assert the test token was sent.
	// 		assert.NotEmpty(t, response.Id, "response ID should not be empty")
	// 	})
	// })

	t.Run("SendGrid", func(t *testing.T) {

		t.Run("SendVerificationEmail should work properly and return a response payload", func(t *testing.T) {
			srv := NewEmailService(SENDGRID_MAIL_PROVIDER)

			response, err := srv.SendVerificationEmail(context.Background(), testToken, testEmail)

			// error checker.
			assert.NoError(t, err, "unexpected error")

			// assert the test token was sent.
			assert.NotEmpty(t, response.Id, "response ID should not be empty")
			assert.Equal(t, "thisisatestid", response.Id, "response ID should be equal to thisisatestid")
		})

		t.Run("testMail service",
			func(t *testing.T) {
				// srv := NewEmailService(SENDGRID_MAIL_PROVIDER)
				err := SendTemplateEmail(types.SendEmailPayload{
					FromAddress: "oayobami15@gmail.com",
					ToAddress:   "xlassixxx@gmail.com",
					Subject:     "Main test",
					DynamicData: map[string]interface{}{
						"code":           "654321",
						"recipient_name": "",
						"subject":        "trying",
					},
				}, "d-be6f6e592fd242ca9e14db589f21c1b1")

				assert.NoError(t, err)
				// t.Run("Fetch Templates", func(t *testing.T) {
				// 	SendVerificationEmailV2(types.SendEmailPayload{
				// 		FromAddress: "oayobami15@gmail.com",
				// 		ToAddress:   "xlassixxx@gmail.com",
				// 		Subject:     "Main test",
				// 		DynamicData: map[string]interface{}{
				// 			"code":           "84321",
				// 			"recipient_name": "",
				// 			"subject":        "trying",
				// 		},
				// 	})
				// })
			})
	})
}
