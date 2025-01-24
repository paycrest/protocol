package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/utils/test"
	"github.com/stretchr/testify/assert"
)

var testCtx struct {
	router *gin.Engine
}

func setup() error {
	// Set Gin mode for testing
	gin.SetMode(gin.TestMode)

	// Initialize router
	router := gin.New()
	router.Use(RateLimitMiddleware())

	// Add test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Assign router to the test context
	testCtx.router = router
	return nil
}

func TestMain(m *testing.M) {
	// Perform setup before running tests
	if err := setup(); err != nil {
		panic(err)
	}

	// Run all tests
	m.Run()
}

// Helper function to decode JSON responses
func decodeResponseBody(t *testing.T, body *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	err := json.NewDecoder(body.Body).Decode(&response)
	assert.NoError(t, err)
	return response
}

func TestRateLimitMiddleware(t *testing.T) {
	router := testCtx.router

	tests := []struct {
		name           string
		authenticated  bool
		numRequests    int
		expectedStatus int
		token          string
	}{
		{
			name:           "Unauthenticated Under Limit",
			authenticated:  false,
			numRequests:    4,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthenticated Over Limit",
			authenticated:  false,
			numRequests:    6,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "Authenticated Under Limit",
			authenticated:  true,
			numRequests:    45,
			expectedStatus: http.StatusOK,
			token:          "test-token",
		},
		{
			name:           "Authenticated Over Limit",
			authenticated:  true,
			numRequests:    55,
			expectedStatus: http.StatusTooManyRequests,
			token:          "test-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wait for the rate limit to reset
			time.Sleep(1 * time.Second)

			headers := map[string]string{}
			if tt.authenticated {
				headers["Authorization"] = tt.token
			}

			var lastStatus int
			for i := 0; i < tt.numRequests; i++ {
				w, _ := test.PerformRequest(t, "GET", "/test", nil, headers, router)
				lastStatus = w.Code
			}

			assert.Equal(t, tt.expectedStatus, lastStatus)
		})
	}
}

func TestRateLimitErrorResponse(t *testing.T) {
	router := testCtx.router

	// Make enough requests to trigger the rate limit
	headers := map[string]string{}
	for i := 0; i < 6; i++ {
		w, _ := test.PerformRequest(t, "GET", "/test", nil, headers, router)

		// Check if the rate limit was triggered
		if w.Code == http.StatusTooManyRequests {
			response := decodeResponseBody(t, w)

			// Verify the top-level structure of the response
			assert.Equal(t, "error", response["status"])
			assert.Equal(t, "Too many requests from this IP address", response["message"])

			// Check nested data fields
			data, ok := response["data"].(map[string]interface{})
			assert.True(t, ok, "data field should be a map")
			assert.Contains(t, data, "retry_after")
			assert.Contains(t, data, "limit")

			// Verify values
			assert.Greater(t, data["retry_after"].(float64), 0.0)
			assert.Equal(t, 5, int(data["limit"].(float64)))

			break
		}
	}
}
