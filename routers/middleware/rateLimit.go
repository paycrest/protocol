package middleware

import (
	"os"
	"strconv"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

// getEnvInt gets an integer from environment variable with fallback
func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

// getLimits gets rate limits from environment
func getLimits() (unauthenticated, authenticated int) {
	return getEnvInt("RATE_LIMIT_UNAUTHENTICATED", 5),
		getEnvInt("RATE_LIMIT_AUTHENTICATED", 50)
}

func keyFunc(c *gin.Context) string {
	if token := c.GetHeader("Authorization"); token != "" {
		return "auth:" + token
	}
	return "ip:" + c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	message := "Too many requests from this IP address"
	if c.GetHeader("Authorization") != "" {
		message = "Too many requests for this API key"
	}

	c.JSON(429, gin.H{
		"error":       message,
		"retry_after": time.Until(info.ResetTime).Seconds(),
		"limit":       info.Limit,
	})
}

func RateLimitMiddleware() gin.HandlerFunc {
	unauthenticatedLimit, authenticatedLimit := getLimits()

	// Store for unauthenticated requests
	unauthenticatedStore := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: uint(unauthenticatedLimit),
	})

	// Store for authenticated requests
	authenticatedStore := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: uint(authenticatedLimit),
	})

	// Limiter for unauthenticated requests
	unauthenticatedLimiter := ratelimit.RateLimiter(unauthenticatedStore, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	// Limiter for authenticated requests
	authenticatedLimiter := ratelimit.RateLimiter(authenticatedStore, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	return func(c *gin.Context) {
		if token := c.GetHeader("Authorization"); token != "" {
			authenticatedLimiter(c)
		} else {
			unauthenticatedLimiter(c)
		}
		c.Next()
	}
}
