package utils

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestHttp(t *testing.T) {
	t.Run("Make JSON Request with retries on failure", func(t *testing.T) {

		// setup httpmock
		httpmock.Activate()
		defer httpmock.Deactivate()

		// register mock response
		httpmock.RegisterResponder("POST", "https://example.com/", func(r *http.Request) (*http.Response, error) {
			return httpmock.NewBytesResponse(500, nil), nil
		},
		)
		httpmock.RegisterResponder("PUT", "https://example.com/", func(r *http.Request) (*http.Response, error) {
			return httpmock.NewBytesResponse(400, nil), nil
		},
		)
		httpmock.RegisterResponder("PATCH", "https://example.com/", func(r *http.Request) (*http.Response, error) {
			return httpmock.NewBytesResponse(200, nil), errors.New("429 Too Many Requests")
		},
		)

		httpmock.RegisterResponder("GET", "https://example.com/", func(r *http.Request) (*http.Response, error) {
			return httpmock.NewBytesResponse(200, []byte(`{"success": "success"}`)), nil
		},
		)

		t.Run("with 500 error", func(t *testing.T) {
			timeStart := time.Now()
			_, err := MakeJSONRequest(context.Background(), "POST", "https://example.com/", nil, nil)
			assert.Error(t, err)
			assert.True(t, time.Since(timeStart) > 10*time.Second)
		})

		t.Run("with 429 error", func(t *testing.T) {
			timeStart := time.Now()
			_, err := MakeJSONRequest(context.Background(), "PATCH", "https://example.com/", nil, nil)
			assert.Error(t, err)
			assert.True(t, time.Since(timeStart) > 10*time.Second)
		})

		t.Run("with 400 error", func(t *testing.T) {
			timeStart := time.Now()
			_, err := MakeJSONRequest(context.Background(), "PUT", "https://example.com/", nil, nil)
			assert.Error(t, err)
			assert.True(t, time.Since(timeStart) < 1*time.Second)
		})

		t.Run("with 200 sucess", func(t *testing.T) {
			res, err := MakeJSONRequest(context.Background(), "GET", "https://example.com/", nil, nil)
			assert.NoError(t, err)
			assert.Equal(t, map[string]interface{}{"success": "success"}, res)
		})

	})

}
