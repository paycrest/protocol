package accounts

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

// Register controller validates the payload and creates a new user.
// It hashes the password provided by the user.
// It also sends an email to verify the user's email address.
func (ctrl *AuthController) Register(ctx *gin.Context) {
	var payload RegisterPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Check if user with email already exists
	userTmp, _ := db.Client.User.
		Query().
		Where(user.EmailEQ(strings.ToLower(payload.Email))).
		Only(ctx)

	if userTmp != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"User with email already exists", nil)
		return
	}

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
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to create new user", err.Error())
		return
	}

	// TODO: Send email to verify the user's email address

	u.APIResponse(ctx, http.StatusCreated, "success", "User created successfully",
		&RegisterResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
}

// Login controller validates the payload and creates a new user.
func (ctrl *AuthController) Login(ctx *gin.Context) {
	var payload LoginPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
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
		u.APIResponse(ctx, http.StatusUnauthorized, "error",
			"Invalid credentials", "Email and password do not match any user",
		)
		return
	}

	// Generate JWT pair
	accessToken, refreshToken, err := token.GeneratePairJWT(user.ID.String())

	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to create token pair", err.Error(),
		)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Successfully logged in", &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshJWT controller returns a new access token given a valid refresh token.
func (ctrl *AuthController) RefreshJWT(ctx *gin.Context) {
	var payload RefreshJWTPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Validate the refresh token
	claims, err := token.ValidateJWT(payload.RefreshToken)
	userID, ok := claims["sub"].(string)
	if err != nil || !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid or expired refresh token", err.Error())
		return
	}

	// Generate a new access token
	accessToken, err := token.GenerateAccessJWT(userID)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to generate access token", err.Error())
		return
	}

	// Return the new access token
	u.APIResponse(ctx, http.StatusOK, "success", "Successfully refreshed access token", &RefreshResponse{
		AccessToken: accessToken,
	})
}
