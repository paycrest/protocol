package token

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/paycrest/paycrest-protocol/config"
)

var conf = config.AuthConfig()

// GenerateAccessJWT generates an access token with a short expiry time ~ 15 minutes
func GenerateAccessJWT(user_id, name string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user_id
	claims["name"] = name
	claims["exp"] = time.Now().Add(conf.JwtAccessHourLifespan).Unix()

	return token.SignedString([]byte(conf.Secret))

}

// GenerateRefreshJWT generates a refresh token with a long expiry time >= 24 hours
func GenerateRefreshJWT(user_id, name string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user_id
	claims["name"] = name
	claims["exp"] = time.Now().Add(conf.JwtRefreshHourLifespan).Unix()

	return token.SignedString([]byte(conf.Secret))
}

// GeneratePairJWT generates a pair of access and refresh tokens
func GeneratePairJWT(user_id, name string) (string, string, error) {
	access, err := GenerateAccessJWT(user_id, name)
	if err != nil {
		return "", "", err
	}

	refresh, err := GenerateRefreshJWT(user_id, name)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
