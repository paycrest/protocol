package accounts

import (
	"time"

	"github.com/google/uuid"
)

// RegisterPayload is the payload for the register endpoint
type RegisterPayload struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// RegisterResponse is the response for the register endpoint
type RegisterResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

// LoginPayload is the payload for the login endpoint
type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response for the login endpoint
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshJWTPayload is the payload for the refresh endpoint
type RefreshJWTPayload struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponse is the response for the refresh endpoint
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
