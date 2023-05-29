package accounts

import (
	"time"

	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
)

// RegisterPayload is the payload for the register endpoint
type RegisterPayload struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6,max=20"`
}

// RegisterResponse is the response for the register endpoint
type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

// LoginPayload is the payload for the login endpoint
type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response for the login endpoint
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshJWTPayload is the payload for the refresh endpoint
type RefreshJWTPayload struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshResponse is the response for the refresh endpoint
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}

// GnerateAPIKeyPayload is the payload for the generate api key endpoint
type GenerateAPIKeyPayload struct {
	Name  string `json:"name" binding:"required"`
	Scope string `json:"type" binding:"required,oneof=sender provider tx_validator"`
}

// GenerateAPIKeyResponse is the response for the generate api key endpoint
type GenerateAPIKeyResponse struct {
	ID        int          `json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	Name      string       `json:"name"`
	Scope     apikey.Scope `json:"scope"`
	PublicKey string       `json:"publicKey"`
	SecretKey string       `json:"secretKey"`
	IsActive  bool         `json:"isActive"`
}
