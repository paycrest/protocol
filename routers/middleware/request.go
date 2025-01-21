package middleware

import (
	"net/http"
	"sync"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/config"
	u "github.com/paycrest/protocol/utils"
)

var (
	unauthenticatedLimiter gin.HandlerFunc
	authenticatedLimiter   gin.HandlerFunc
	initOnce               sync.Once
)

// RateLimitMiddleware initializes and returns the rate limiter middleware
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize limiters once
		initOnce.Do(func() {
			conf := config.ServerConfig()

			// Unauthenticated limiter
			unauthenticatedStore := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
				Rate:  time.Second,
				Limit: uint(conf.RateLimitUnauthenticated),
			})
			unauthenticatedLimiter = ratelimit.RateLimiter(unauthenticatedStore, &ratelimit.Options{
				ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
					u.APIResponse(
						c,
						http.StatusTooManyRequests,
						"error",
						"Too many requests from this IP address",
						map[string]interface{}{
							"retry_after": time.Until(info.ResetTime).Seconds(),
							"limit":       info.Limit,
						},
					)
					c.Abort()
				},
				KeyFunc: func(c *gin.Context) string {
					return "ip:" + c.ClientIP()
				},
			})

			// Authenticated limiter
			authenticatedStore := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
				Rate:  time.Second,
				Limit: uint(conf.RateLimitAuthenticated),
			})
			authenticatedLimiter = ratelimit.RateLimiter(authenticatedStore, &ratelimit.Options{
				ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
					u.APIResponse(
						c,
						http.StatusTooManyRequests,
						"error",
						"Too many requests for this API key",
						map[string]interface{}{
							"retry_after": time.Until(info.ResetTime).Seconds(),
							"limit":       info.Limit,
						},
					)
					c.Abort()
				},
				KeyFunc: func(c *gin.Context) string {
					return "auth:" + c.GetHeader("Authorization")
				},
			})
		})

		// Apply appropriate limiter based on authentication status
		if token := c.GetHeader("Authorization"); token != "" {
			authenticatedLimiter(c)
		} else {
			unauthenticatedLimiter(c)
		}
		
		c.Next()
	}
}
