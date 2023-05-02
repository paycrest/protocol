package token

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/paycrest/paycrest-protocol/config"
)

var conf = config.AuthConfig()

// GenerateAccessJWT generates an access token with a short expiry time ~ 15 minutes
func GenerateAccessJWT(user_id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user_id
	claims["exp"] = time.Now().Add(conf.JwtAccessHourLifespan).Unix()

	return token.SignedString([]byte(conf.Secret))

}

// GenerateRefreshJWT generates a refresh token with a long expiry time >= 24 hours
func GenerateRefreshJWT(user_id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user_id
	claims["exp"] = time.Now().Add(conf.JwtRefreshHourLifespan).Unix()

	return token.SignedString([]byte(conf.Secret))
}

// GeneratePairJWT generates a pair of access and refresh tokens
func GeneratePairJWT(user_id, name string) (string, string, error) {
	access, err := GenerateAccessJWT(user_id)
	if err != nil {
		return "", "", err
	}

	refresh, err := GenerateRefreshJWT(user_id)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// ValidateJWT validates the JWT token string and returns the claims if valid.
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.Secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if token is expired
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
