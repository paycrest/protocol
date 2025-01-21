package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type CacheMetrics struct {
	hits   prometheus.Counter
	misses prometheus.Counter
}

type CacheService struct {
	client  *redis.Client
	metrics CacheMetrics
}

func NewCacheService(config CacheConfig) (*CacheService, error) {
	metrics := CacheMetrics{
		hits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		}),
		misses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		}),
	}
	prometheus.MustRegister(metrics.hits, metrics.misses)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &CacheService{client: client, metrics: metrics}, nil
}

func generateCacheKey(c *gin.Context) string {
	path := c.Request.URL.Path
	switch {
	case path == "/v1/currencies":
		return "api:currencies:list"
	case path == "/v1/pubkey":
		return "api:aggregator:pubkey"
	case len(c.Param("currency_code")) > 0:
		return fmt.Sprintf("api:institutions:%s", c.Param("currency_code"))
	default:
		return fmt.Sprintf("api:v1:%s:%s", c.Request.Method, path)
	}
}

func generateETag(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *CacheService) CacheMiddleware(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := generateCacheKey(c)
		ctx := context.Background()

		// Check ETag
		if etag := c.GetHeader("If-None-Match"); etag != "" {
			if storedETag, _ := s.client.Get(ctx, key+":etag").Result(); etag == storedETag {
				c.Status(http.StatusNotModified)
				return
			}
		}

		// Try to get from cache
		val, err := s.client.Get(ctx, key).Result()
		if err == nil {
			s.metrics.hits.Inc()
			etag, _ := s.client.Get(ctx, key+":etag").Result()

			c.Header("X-Cache", "HIT")
			c.Header("Cache-Control", fmt.Sprintf("max-age=%d, stale-while-revalidate=60", int(duration.Seconds())))
			c.Header("ETag", etag)
			c.String(200, val)

			// Background revalidation if approaching expiry
			if ttl, _ := s.client.TTL(ctx, key).Result(); ttl < time.Minute {
				go s.revalidateCache(c.Copy(), key, duration)
			}
			return
		}

		s.metrics.misses.Inc()
		c.Writer = &cacheWriter{ResponseWriter: c.Writer, body: make([]byte, 0)}
		c.Next()

		if c.Writer.Status() == 200 {
			response := c.Writer.(*cacheWriter).body
			etag := generateETag(response)

			s.client.Set(ctx, key, string(response), duration)
			s.client.Set(ctx, key+":etag", etag, duration)

			c.Header("ETag", etag)
			c.Header("Cache-Control", fmt.Sprintf("max-age=%d, stale-while-revalidate=60", int(duration.Seconds())))
		}
	}
}

type cacheWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *cacheWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func (s *CacheService) WarmCache(ctx context.Context) error {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	endpoints := map[string]time.Duration{
		"currencies":       24 * time.Hour,
		"pubkey":           365 * time.Hour,
		"institutions/USD": 24 * time.Hour,
		"institutions/EUR": 24 * time.Hour,
		"institutions/GBP": 24 * time.Hour,
	}

	for path, duration := range endpoints {
		url := fmt.Sprintf("%s/v1/%s", baseURL, path)
		key := generateCacheKey(&gin.Context{Request: &http.Request{URL: &url.URL{Path: "/v1/" + path}}})
		if err := s.cacheEndpoint(ctx, url, key, duration); err != nil {
			fmt.Printf("Failed to cache %s: %v\n", path, err)
		}
	}

	return nil
}

func (s *CacheService) cacheEndpoint(ctx context.Context, url, key string, duration time.Duration) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data from %s: status code %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	etag := generateETag(body)
	s.client.Set(ctx, key, string(body), duration)
	s.client.Set(ctx, key+":etag", etag, duration)

	return nil
}

func (s *CacheService) revalidateCache(c *gin.Context, key string, duration time.Duration) {
	ctx := context.Background()
	req, _ := http.NewRequest(c.Request.Method, c.Request.URL.String(), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		etag := generateETag(body)
		s.client.Set(ctx, key, string(body), duration)
		s.client.Set(ctx, key+":etag", etag, duration)
	}
}
