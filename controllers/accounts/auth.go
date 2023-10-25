package accounts

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/ent/fiatcurrency"
	"github.com/paycrest/paycrest-protocol/ent/providerprofile"
	userEnt "github.com/paycrest/paycrest-protocol/ent/user"
	"github.com/paycrest/paycrest-protocol/ent/verificationtoken"
	svc "github.com/paycrest/paycrest-protocol/services"
	db "github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/crypto"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/paycrest/paycrest-protocol/utils/token"
)

// AuthController is the controller type for the auth endpoints
type AuthController struct {
	apiKeyService *svc.APIKeyService
	emailService  *svc.EmailService
}

// NewAuthController creates a new instance of AuthController with injected services
func NewAuthController() *AuthController {
	return &AuthController{
		apiKeyService: svc.NewAPIKeyService(),
		emailService:  svc.NewEmailService(svc.MAILGUN_MAIL_PROVIDER),
	}
}

// Register controller validates the payload and creates a new user.
// It hashes the password provided by the user.
// It also sends an email to verify the user's email address.
func (ctrl *AuthController) Register(ctx *gin.Context) {
	var payload types.RegisterPayload
	scope := ctx.MustGet("scope").(userEnt.Scope)

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	tx, err := db.Client.Tx(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to create new user", nil)
		return
	}

	// Check if user with email already exists
	userTmp, _ := tx.User.
		Query().
		Where(
			userEnt.EmailEQ(strings.ToLower(payload.Email)),
			userEnt.ScopeEQ(userEnt.Scope(scope)),
		).
		Only(ctx)

	if userTmp != nil {
		_ = tx.Rollback()
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"User with email already exists", nil)
		return
	}

	// Save the user
	user, err := tx.User.
		Create().
		SetFirstName(payload.FirstName).
		SetLastName(payload.LastName).
		SetEmail(strings.ToLower(payload.Email)).
		SetPassword(payload.Password).
		SetScope(scope).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to create new user", nil)
		return
	}

	// Send verification email
	verificationToken, err := tx.VerificationToken.
		Create().
		SetOwner(user).
		SetScope(verificationtoken.ScopeVerification).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
	}

	if verificationToken != nil {
		if _, err := ctrl.emailService.SendVerificationEmail(ctx, verificationToken.Token, user.Email); err != nil {
			logger.Errorf("error: %v", err)
		}
	}

	// Create a provider profile
	clientType := ctx.GetHeader("Client-Type")
	if scope == userEnt.ScopeProvider && (clientType == "web" || clientType == "mobile") {
		// Fetch currency
		currency, err := tx.FiatCurrency.
			Query().
			Where(fiatcurrency.CodeEQ(payload.Currency)).
			Only(ctx)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", nil)
			return
		}

		provider, err := tx.ProviderProfile.
			Create().
			SetTradingName(payload.TradingName).
			SetCurrency(currency).
			SetUser(user).
			SetProvisionMode(providerprofile.ProvisionModeManual).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", nil)
			return
		}

		// Generate the API key using the service
		_, _, err = ctrl.apiKeyService.GenerateAPIKey(ctx, tx, nil, provider, nil)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", nil)
			return
		}
	}

	// Create a sender profile
	if scope == userEnt.ScopeSender {
		sender, err := tx.SenderProfile.
			Create().
			SetUser(user).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", nil)
			return
		}

		// Generate the API key using the service
		_, _, err = ctrl.apiKeyService.GenerateAPIKey(ctx, tx, sender, nil, nil)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", nil)
			return
		}
	}

	// Create a validator profile
	if scope == userEnt.ScopeTxValidator {
		validator, err := tx.ValidatorProfile.
			Create().
			SetUser(user).
			Save(ctx)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", nil)
			return
		}

		// Generate the API key using the service
		_, _, err = ctrl.apiKeyService.GenerateAPIKey(ctx, tx, nil, nil, validator)
		if err != nil {
			_ = tx.Rollback()
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error",
				"Failed to create new user", nil)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error",
			"Failed to create new user", nil)
		return
	}

	u.APIResponse(ctx, http.StatusCreated, "success", "User created successfully",
		&types.RegisterResponse{
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
	var payload types.LoginPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Fetch user by email
	user, emailErr := db.Client.User.
		Query().
		Where(userEnt.EmailEQ(strings.ToLower(payload.Email))).
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
			"Failed to create token pair", nil,
		)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Successfully logged in", &types.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshJWT controller returns a new access token given a valid refresh token.
func (ctrl *AuthController) RefreshJWT(ctx *gin.Context) {
	var payload types.RefreshJWTPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Validate the refresh token
	claims, err := token.ValidateJWT(payload.RefreshToken)
	userID, ok := claims["sub"].(string)
	if err != nil || !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid or expired refresh token", nil)
		return
	}

	// Generate a new access token
	accessToken, err := token.GenerateAccessJWT(userID)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to generate access token", nil)
		return
	}

	// Return the new access token
	u.APIResponse(ctx, http.StatusOK, "success", "Successfully refreshed access token", &types.RefreshResponse{
		AccessToken: accessToken,
	})
}

// ConfirmEmail controller validates the payload and confirm the users email.
func (ctrl *AuthController) ConfirmEmail(ctx *gin.Context) {
	var payload types.ConfirmEmailPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Fetch verificationtoken
	verificationToken, vtErr := db.Client.VerificationToken.
		Query().
		Where(verificationtoken.TokenEQ(payload.Token)).
		WithOwner().
		Only(ctx)
	if vtErr != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid verification token", vtErr.Error())
		return
	}

	// Update User IsVerified to true
	_, setIfVerifiedErr := verificationToken.Edges.Owner.
		Update().
		SetIsVerified(true).
		Save(ctx)
	if setIfVerifiedErr != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to verify user email", setIfVerifiedErr.Error())
		return
	}

	// Return a success response
	u.APIResponse(ctx, http.StatusOK, "success", "User email verified successfully", nil)
}

// ResendVerificationToken controller resends the verification token to the users email.
func (ctrl *AuthController) ResendVerificationToken(ctx *gin.Context) {
	var payload types.ResendTokenPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Fetch User account.
	user, userErr := db.Client.User.Query().Where(userEnt.EmailEQ(payload.Email)).Only(ctx)
	if userErr != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid credential", userErr.Error())
	}

	// Generate VerificationToken.
	verificationtoken, vtErr := db.Client.VerificationToken.Create().SetOwner(user).SetScope(verificationtoken.Scope(payload.Scope)).Save(ctx)
	if vtErr != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to generate verification token", vtErr.Error())
		return
	}

	if _, err := ctrl.emailService.SendVerificationEmail(ctx, verificationtoken.Token, user.Email); err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to send verification email", vtErr.Error())
		return
	}

	// Return a success response
	u.APIResponse(ctx, http.StatusOK, "success", "Verification token has been sent to your email", nil)
}
