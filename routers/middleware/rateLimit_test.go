
package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimitMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

func TestRateLimitMiddleware(t *testing.T) {
	router := setupTestRouter()

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
			// Wait for rate limit to reset
			time.Sleep(1 * time.Second)

			var lastStatus int
			for i := 0; i < tt.numRequests; i++ {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/test", nil)

				if tt.authenticated {
					req.Header.Set("Authorization", tt.token)
				}

				router.ServeHTTP(w, req)
				lastStatus = w.Code
			}

			assert.Equal(t, tt.expectedStatus, lastStatus)
		})
	}
}

func TestRateLimitErrorResponse(t *testing.T) {
	router := setupTestRouter()

	// Make enough requests to trigger rate limit
	for i := 0; i < 6; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code == http.StatusTooManyRequests {
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Check response structure
			assert.Contains(t, response, "error")
			assert.Contains(t, response, "retry_after")
			assert.Contains(t, response, "limit")
			assert.Equal(t, "Too many requests from this IP address", response["error"])
			break
		}
	}
}