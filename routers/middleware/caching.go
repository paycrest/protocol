package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

// CacheMetrics holds the metrics for cache hits and misses
type CacheMetrics struct {
	hits   prometheus.Counter
	misses prometheus.Counter
}

// CacheService handles Redis operations
type CacheService struct {
	client  *redis.Client
	metrics CacheMetrics
}

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewCacheService creates a new Redis cache service
func NewCacheService(config config.RedisConfiguration) (*CacheService, error) {
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
	conf := config.RedisConfig()
	path := c.Request.URL.Path
	switch {
	case path == "/v1/currencies":
		return fmt.Sprintf("%s:api:currencies:list", conf.CacheVersion)
	case path == "/v1/pubkey":
		return fmt.Sprintf("%s:api:aggregator:pubkey", conf.CacheVersion)
	case len(c.Param("currency_code")) > 0:
		return fmt.Sprintf("%s:api:institutions:%s", conf.CacheVersion, c.Param("currency_code"))
	default:
		return fmt.Sprintf("%s:api:%s", conf.CacheVersion, path)
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
	conf := config.ServerConfig()
	baseURL := conf.HostDomain
	if baseURL == "" {
		return fmt.Errorf("host domain is not set in the server configuration")
	}

	// Fetch the list of supported currencies with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	currenciesURL := fmt.Sprintf("%s/v1/currencies", baseURL)
	resp, err := client.Get(currenciesURL)
	if err != nil {
		return fmt.Errorf("failed to fetch currencies from %s: %v", currenciesURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch currencies: status code %d", resp.StatusCode)
	}

	var currencies []string
	if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		return fmt.Errorf("failed to decode currencies response: %v", err)
	}

	if len(currencies) == 0 {
		fmt.Println("No currencies found. Using default currencies [USD, EUR, GBP].")
		currencies = []string{"USD", "EUR", "GBP"}
	}

	// Define static and dynamic endpoints
	endpoints := map[string]time.Duration{
		"currencies": time.Duration(conf.CurrenciesCacheDuration) * time.Hour,
		"pubkey":     time.Duration(conf.PubKeyCacheDuration) * time.Hour,
	}

	for _, currency := range currencies {
		endpoints[fmt.Sprintf("institutions/%s", currency)] = time.Duration(conf.InstitutionsCacheDuration) * time.Hour
	}

	// Warm up cache for each endpoint
	for path, duration := range endpoints {
		urlStr := fmt.Sprintf("%s/v1/%s", baseURL, path)
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			fmt.Printf("Failed to parse URL %s: %v\n", urlStr, err)
			continue
		}
		key := generateCacheKey(&gin.Context{Request: &http.Request{URL: parsedURL}})
		if err := s.cacheEndpoint(ctx, urlStr, key, duration); err != nil {
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
