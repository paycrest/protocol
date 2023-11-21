package token

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/paycrest/protocol/config"
)

var conf = config.AuthConfig()

// GenerateAccessJWT generates an access token with a short expiry time ~ 15 minutes
func GenerateAccessJWT(userID string, scope string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["scope"] = scope
	claims["exp"] = time.Now().Add(conf.JwtAccessHourLifespan).Unix()

	return token.SignedString([]byte(conf.Secret))
}

// GenerateRefreshJWT generates a refresh token with a long expiry time >= 24 hours
func GenerateRefreshJWT(userID string, scope string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["scope"] = scope
	claims["exp"] = time.Now().Add(conf.JwtRefreshHourLifespan).Unix()

	return token.SignedString([]byte(conf.Secret))
}

// GeneratePairJWT generates a pair of access and refresh tokens
func GeneratePairJWT(userID string, scope string) (string, string, error) {
	access, err := GenerateAccessJWT(userID, scope)
	if err != nil {
		return "", "", err
	}

	refresh, err := GenerateRefreshJWT(userID, scope)
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
		return nil, fmt.Errorf("invalid token")
	}

	// Check if token is expired
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// GeneratePrivateKey generates a private key (Secret Key).
func GeneratePrivateKey() (string, error) {
	// Generate random bytes for the key
	keySize := 32 // 32 bytes = 256 bits -- for HMAC-SHA256 hashing function
	privateKeyBytes := make([]byte, keySize)
	_, err := rand.Read(privateKeyBytes)
	if err != nil {
		return "", err
	}

	// Encode private key to base64 string
	privateKey := base64.URLEncoding.EncodeToString(privateKeyBytes)

	return privateKey, nil
}

// VerifyHMACSignature verifies the HMAC signature for the given payload using the private key
// and returns true if the signature is valid.
func VerifyHMACSignature(payload map[string]interface{}, privateKey string, signature string) bool {
	expectedSignature := []byte(GenerateHMACSignature(payload, privateKey))
	computedSignature := []byte(signature)
	return hmac.Equal(expectedSignature, computedSignature)
}

// GenerateHMACSignature generates the HMAC signature for the given payload using the private key.
// The signature is returned as a hex-encoded string.
func GenerateHMACSignature(payload map[string]interface{}, privateKey string) string {
	key := []byte(privateKey)
	h := hmac.New(sha256.New, key)
	payloadBytes, _ := json.Marshal(payload)

	h.Write(payloadBytes)
	return hex.EncodeToString(h.Sum(nil))
}
