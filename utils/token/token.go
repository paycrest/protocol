package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/paycrest/paycrest-protocol/config"
)

var conf = config.AuthConfig()

// GenerateAccessJWT generates an access token with a short expiry time ~ 15 minutes
func GenerateAccessJWT(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(conf.JwtAccessHourLifespan).Unix()

	return token.SignedString([]byte(conf.Secret))
}

// GenerateRefreshJWT generates a refresh token with a long expiry time >= 24 hours
func GenerateRefreshJWT(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(conf.JwtRefreshHourLifespan).Unix()

	return token.SignedString([]byte(conf.Secret))
}

// GeneratePairJWT generates a pair of access and refresh tokens
func GeneratePairJWT(userID string) (string, string, error) {
	access, err := GenerateAccessJWT(userID)
	if err != nil {
		return "", "", err
	}

	refresh, err := GenerateRefreshJWT(userID)
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

// GenerateHMACKeys generates a pair of public and private keys for HMAC authentication.
// It returns the public key (Client Key) and private key (Secret Key).
func GenerateHMACKeys() (string, string, error) {
	// Generate random bytes for the keys
	keySize := 32 // 32 bytes = 256 bits -- for HMAC-SHA256 hashing function
	publicKeyBytes := make([]byte, keySize)
	privateKeyBytes := make([]byte, keySize)
	_, err := rand.Read(publicKeyBytes)
	if err != nil {
		return "", "", err
	}
	_, err = rand.Read(privateKeyBytes)
	if err != nil {
		return "", "", err
	}

	// Encode keys to base64 strings
	publicKey := base64.URLEncoding.EncodeToString(publicKeyBytes)
	privateKey := base64.URLEncoding.EncodeToString(privateKeyBytes)

	return publicKey, privateKey, nil
}
