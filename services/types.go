package services

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
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
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

// GenerateAPIKeyPayload is the payload for the generate API key endpoint
type GenerateAPIKeyPayload struct {
	Name  string       `json:"name" binding:"required"`
	Scope apikey.Scope `json:"scope" binding:"required,oneof=sender provider tx_validator"`
}

// APIKeyResponse is the response type for an API key
type APIKeyResponse struct {
	ID        uuid.UUID    `json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	Name      string       `json:"name"`
	Scope     apikey.Scope `json:"scope"`
	Secret    string       `json:"secret"`
	IsActive  bool         `json:"isActive"`
}
