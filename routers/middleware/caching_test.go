package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis() (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestCacheMiddleware(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mr, client := setupTestRedis()
	defer mr.Close()

	cacheService := &CacheService{
		client: client,
		metrics: CacheMetrics{
			hits:   prometheus.NewCounter(prometheus.CounterOpts{Name: "test_cache_hits_total"}),
			misses: prometheus.NewCounter(prometheus.CounterOpts{Name: "test_cache_misses_total"}),
		},
	}

	// Create test router
	router := gin.New() // Use gin.New() instead of Default() to avoid extra middleware
	router.GET("/v1/currencies", cacheService.CacheMiddleware(24*time.Hour), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"currencies": []string{"USD", "EUR", "GBP"}})
	})

	// First request (should miss)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/v1/currencies", nil)
	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "MISS", w1.Header().Get("X-Cache"))

	// Verify data was cached
	key := "v1:api:currencies:list"
	_, err := client.Get(context.Background(), key).Result()
	assert.NoError(t, err)

	// Second request (should hit)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/v1/currencies", nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "HIT", w2.Header().Get("X-Cache"))
}

func TestWarmCache(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mr, client := setupTestRedis()
	defer mr.Close()

	// Verify Redis connection
	ctx := context.Background()
	err := client.Ping(ctx).Err()
	require.NoError(t, err, "Redis connection failed")

	// Create mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var response interface{}

		switch r.URL.Path {
		case "/v1/currencies":
			response = []string{"USD", "EUR", "GBP"}
		case "/v1/pubkey":
			response = map[string]string{"key": "test-key"}
		default:
			if matched, _ := path.Match("/v1/institutions/*", r.URL.Path); matched {
				currency := path.Base(r.URL.Path)
				response = map[string]interface{}{
					"institutions": []map[string]string{
						{"id": "bank1", "name": "Bank 1", "currency": currency},
						{"id": "bank2", "name": "Bank 2", "currency": currency},
					},
				}
			} else {
				http.NotFound(w, r)
				return
			}
		}

		err := json.NewEncoder(w).Encode(response)
		require.NoError(t, err, "Failed to encode response")
	}))
	defer mockServer.Close()

	// Create cache service
	cacheService := &CacheService{
		client: client,
		metrics: CacheMetrics{
			hits: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "test_cache_hits_total",
				Help: "Test cache hits",
			}),
			misses: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "test_cache_misses_total",
				Help: "Test cache misses",
			}),
		},
	}

	// Set required configuration
	viper.Reset()
	viper.Set("HOST_DOMAIN", mockServer.URL)
	viper.Set("CURRENCIES_CACHE_DURATION", 24)
	viper.Set("INSTITUTIONS_CACHE_DURATION", 24)
	viper.Set("PUBKEY_CACHE_DURATION", 365)

	// Execute warm cache
	err = cacheService.WarmCache(ctx)
	require.NoError(t, err, "WarmCache failed")

	// Verify caches with explicit error checking
	keys := []string{
		"v1:api:currencies:list",
		"v1:api:aggregator:pubkey",
		"v1:api:institutions:USD",
		"v1:api:institutions:EUR",
		"v1:api:institutions:GBP",
	}

	for _, key := range keys {
		t.Run(fmt.Sprintf("Verify cache for %s", key), func(t *testing.T) {
			// Verify data cache
			val, err := client.Get(ctx, key).Result()
			if err != nil {
				t.Fatalf("Failed to get key %s from cache: %v", key, err)
			}
			if val == "" {
				t.Fatalf("Empty value for key %s", key)
			}

			// Verify JSON validity
			var jsonCheck interface{}
			if err = json.Unmarshal([]byte(val), &jsonCheck); err != nil {
				t.Fatalf("Invalid JSON for key %s: %v", key, err)
			}

			// Log the cached data for debugging
			t.Logf("Cached data for %s: %s", key, val)

			// Check if ETag exists
			etagKey := key + ":etag"
			etag, err := client.Get(ctx, etagKey).Result()
			if err != nil {
				t.Logf("Error getting ETag for key %s: %v", key, err)
				return // Skip ETag verification if not present
			}

			// Only verify ETag if one was retrieved
			if etag != "" {
				t.Logf("Found ETag for %s: %s", key, etag)
				expectedETag := generateETag([]byte(val))
				t.Logf("Expected ETag: %s", expectedETag)

				if etag != expectedETag {
					t.Errorf("ETag mismatch for key %s\nGot:      %s\nExpected: %s",
						key, etag, expectedETag)
				}
			}
		})
	}
}
