package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/paycrest/paycrest-protocol/database"
	"github.com/paycrest/paycrest-protocol/ent/user"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/paycrest/paycrest-protocol/utils/token"
)

type AuthController struct{}

// RegisterPayload is the payload for the register endpoint
type RegisterPayload struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// Register controller validates the payload and creates a new user.
// It hashes the password provided by the user.
// It also sends an email to verify the user's email address.
func (ctrl *AuthController) Register(ctx *gin.Context) {
	var payload RegisterPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "error", err.Error())
		return
	}

	// TODO: Check if user with email already exists

	// Save the user
	user, err := db.Client.User.
		Create().
		SetFirstName(payload.FirstName).
		SetLastName(payload.LastName).
		SetEmail(strings.ToLower(payload.Email)).
		SetPassword(payload.Password).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "error", err.Error())
		return
	}

	// TODO: Send email to verify the user's email address

	u.APIResponse(ctx, http.StatusCreated, "success", "User created successfully", &user)
}

// LoginPayload is the payload for the login endpoint
type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login controller validates the payload and creates a new user.
func (ctrl *AuthController) Login(ctx *gin.Context) {
	var payload LoginPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "error", err.Error())
		return
	}

	// Fetch user by email
	user, emailErr := db.Client.User.
		Query().
		Where(user.EmailEQ(strings.ToLower(payload.Email))).
		Only(ctx)

	// Check if the password is correct
	passwordMatch := crypto.CheckPasswordHash(payload.Password, user.Password)

	if !passwordMatch || emailErr != nil {
		logger.Errorf("error: %v", "Invalid email or password")
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Email and password do not match any user", emailErr.Error(),
		)
		return
	}

	// Generate JWT pair
	accessToken, refreshToken, err := token.GeneratePairJWT(user.ID.String(), user.FirstName+" "+user.LastName)

	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to create token pair", err.Error(),
		)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Successfully logged in", gin.H{
		"access":  accessToken,
		"refresh": refreshToken,
	})
}

func (ctrl *AuthController) RefreshJWT(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}
