package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
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
	mr, client := setupTestRedis()
	defer mr.Close()

	cacheService := &CacheService{
		client: client,
		metrics: CacheMetrics{
			hits:   prometheus.NewCounter(prometheus.CounterOpts{Name: "cache_hits_total"}),
			misses: prometheus.NewCounter(prometheus.CounterOpts{Name: "cache_misses_total"}),
		},
	}

	router := gin.Default()
	router.GET("/v1/currencies", cacheService.CacheMiddleware(24*time.Hour), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"currencies": []string{"USD", "EUR", "GBP"}})
	})

	// First request should be a cache miss
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/currencies", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "MISS", w.Header().Get("X-Cache"))

	// Second request should be a cache hit
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "HIT", w.Header().Get("X-Cache"))
}

func TestWarmCache(t *testing.T) {
	mr, client := setupTestRedis()
	defer mr.Close()

	cacheService := &CacheService{
		client: client,
		metrics: CacheMetrics{
			hits:   prometheus.NewCounter(prometheus.CounterOpts{Name: "cache_hits_total"}),
			misses: prometheus.NewCounter(prometheus.CounterOpts{Name: "cache_misses_total"}),
		},
	}

	// Mock server to return currencies
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/currencies" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]string{"USD", "EUR", "GBP"})
		}
	}))
	defer mockServer.Close()

	// Override the base URL in the server configuration
	conf := config.ServerConfig()
	conf.HostDomain = mockServer.URL

	ctx := context.Background()
	err := cacheService.WarmCache(ctx)
	assert.NoError(t, err)

	// Verify that the currencies are cached
	for _, currency := range []string{"USD", "EUR", "GBP"} {
		key := fmt.Sprintf("%s:api:institutions:%s", conf.CacheVersion, currency)
		val, err := client.Get(ctx, key).Result()
		assert.NoError(t, err)
		assert.NotEmpty(t, val)
	}
}
