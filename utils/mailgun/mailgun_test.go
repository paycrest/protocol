package mailgun

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var (
	mailgunEndpoint = "https://api.mailgun.net/v3/sandbox9c66b379b78d43d2b1533bf2a09a5325.mailgun.org/messages"

	testToken = "test-token"
	testEmail = "test@paycrest.io"
)

func TestMailGun(t *testing.T) {
	// activate httpmock
	httpmock.Activate()
	defer httpmock.Deactivate()

	// register mock response
	httpmock.RegisterResponder("POST", mailgunEndpoint,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewBytesResponse(200, []byte(`{"id": "01", "message": "Sent"}`)), nil
		},
	)

	response, err := SendVerificationEmail(testToken, testEmail)

	// error checker.
	assert.NoError(t, err, "unexpected error")

	// assert the test token was sent.
	assert.NotEmpty(t, response.Id, "response ID should not be empty")
}
