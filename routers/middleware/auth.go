package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/token"
)

// JWTMiddleware is a middleware to handle JWT authentication
func JWTMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Authorization header is missing", nil)
		c.Abort()
		return
	}

	// Split the Authorization header value into two parts: the authentication scheme and the token value
	authParts := strings.SplitN(authHeader, " ", 2)
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid Authorization header format", nil)
		c.Abort()
		return
	}

	// Validate the token and extract the user ID
	claims, err := token.ValidateJWT(authParts[1])
	userID, ok := claims["sub"].(string)
	if err != nil || !ok {
		u.APIResponse(c, http.StatusUnauthorized, "error", "Invalid token", nil)
		c.Abort()
		return
	}

	// Set the user_id value in the context of the request
	c.Set("user_id", userID)

	c.Next()
}
