package accounts

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/paycrest-protocol/routers/middleware"
	svc "github.com/paycrest/paycrest-protocol/services"
	db "github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/ent/enttest"
	"github.com/paycrest/paycrest-protocol/ent/providerprofile"
	"github.com/paycrest/paycrest-protocol/ent/senderprofile"
	"github.com/paycrest/paycrest-protocol/ent/user"
	"github.com/paycrest/paycrest-protocol/ent/verificationtoken"
	"github.com/paycrest/paycrest-protocol/utils/test"
	"github.com/paycrest/paycrest-protocol/utils/token"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	// setup httpmock
	httpmock.Activate()
	defer httpmock.Deactivate()

	// register mock response
	httpmock.RegisterResponder("POST", "https://api.mailgun.net/v3/sandbox9c66b379b78d43d2b1533bf2a09a5325.mailgun.org/messages",
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewBytesResponse(200, []byte(`{"id": "01", "message": "Sent"}`)), nil
		},
	)

	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Set up test routers
	router := gin.New()
	ctrl := &AuthController{
		apiKeyService: svc.NewAPIKeyService(),
		emailService:  svc.NewEmailService(svc.SENDGRID_MAIL_PROVIDER),
	}

	router.Use(middleware.ScopeMiddleware)
	router.POST("/register", ctrl.Register)
	router.POST("/login", ctrl.Login)
	router.POST("/confirm-account", ctrl.ConfirmEmail)
	router.POST("/resend-token", ctrl.ResendVerificationToken)
	router.POST("/refresh", middleware.JWTMiddleware, ctrl.RefreshJWT)

	var userID string

	t.Run("Register", func(t *testing.T) {
		t.Run("with valid payload", func(t *testing.T) {
			// Test register with valid payload
			payload := types.RegisterPayload{
				FirstName: "Ike",
				LastName:  "Ayo",
				Email:     "ikeayo@example.com",
				Password:  "password",
			}

			res, err := test.PerformRequest(t, "POST", "/register?scope=sender", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusCreated, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User created successfully", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response data
			assert.Contains(t, data, "id")
			match, _ := regexp.MatchString(
				`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`,
				data["id"].(string),
			)
			if !match {
				t.Errorf("Expected '%s' to be a valid UUID", data["id"].(string))
			}

			userID = data["id"].(string)

			// Parse the user ID string to uuid.UUID
			userUUID, err := uuid.Parse(userID)
			assert.NoError(t, err)
			assert.Equal(t, payload.Email, data["email"].(string))
			assert.Equal(t, payload.FirstName, data["firstName"].(string))
			assert.Equal(t, payload.LastName, data["lastName"].(string))

			// Query the database to check if API key and profile were created for the sender
			senderProfile, err := db.Client.SenderProfile.
				Query().
				Where(senderprofile.HasUserWith(user.ID(userUUID))).
				WithAPIKey().
				Only(context.Background())
			assert.NoError(t, err)

			assert.NotNil(t, senderProfile.Edges.APIKey)
			assert.NotNil(t, senderProfile)
		})

		t.Run("from the provider app", func(t *testing.T) {
			// Test register with valid payload
			payload := types.RegisterPayload{
				FirstName:   "Ike",
				LastName:    "Ayo",
				Email:       "ikeayoprovider@example.com",
				Password:    "password",
				TradingName: "Africana LP",
				Currency:    "NGN",
			}

			_, err := test.CreateTestFiatCurrency(nil)
			assert.NoError(t, err)

			headers := map[string]string{
				"Client-Type": "mobile",
			}

			res, err := test.PerformRequest(t, "POST", "/register?scope=provider", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusCreated, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Parse the user ID string to uuid.UUID
			userUUID, err := uuid.Parse(data["id"].(string))
			assert.NoError(t, err)

			// Query the database to check if API key and profile were created for the provider
			providerProfile, err := db.Client.ProviderProfile.
				Query().
				Where(providerprofile.HasUserWith(user.ID(userUUID))).
				WithAPIKey().
				Only(context.Background())
			assert.NoError(t, err)

			assert.NotNil(t, providerProfile.Edges.APIKey)
			assert.NotNil(t, providerProfile)
		})

		t.Run("with existing user", func(t *testing.T) {
			// Test register with existing user
			payload := types.RegisterPayload{
				FirstName: "Ike",
				LastName:  "Ayo",
				Email:     "ikeayo@example.com",
				Password:  "password",
			}

			res, err := test.PerformRequest(t, "POST", "/register?scope=sender", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusBadRequest, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "User with email already exists", response.Message)
			assert.Nil(t, response.Data)
		})

		t.Run("with invalid email", func(t *testing.T) {
			// Test register with invalid email
			payload := types.RegisterPayload{
				FirstName: "Ike",
				LastName:  "Ayo",
				Email:     "invalid-email",
				Password:  "password",
			}

			res, err := test.PerformRequest(t, "POST", "/register?scope=sender", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusBadRequest, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to validate payload", response.Message)
			assert.Equal(t, "error", response.Status)
			data, ok := response.Data.([]interface{})
			assert.True(t, ok, "response.Data is not of type []interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response errors in data
			assert.Len(t, data, 1)
			errorMap, ok := data[0].(map[string]interface{})
			assert.True(t, ok, "error is not of type map[string]interface{}")
			assert.NotNil(t, errorMap, "error is nil")
			assert.Contains(t, errorMap, "field")
			assert.Equal(t, "Email", errorMap["field"].(string))
			assert.Contains(t, errorMap, "message")
			assert.Equal(t, "Must be a valid email address", errorMap["message"].(string))
		})
	})

	t.Run("ConfirmEmail", func(t *testing.T) {
		// fetch user
		userUUID, err := uuid.Parse(userID)
		assert.NoError(t, err)

		user, fetchUserErr := db.Client.User.
			Query().
			Where(user.IDEQ(userUUID)).
			Only(context.Background())
		assert.NoError(t, fetchUserErr, "failed to fetch user by userID")

		// generate verificationToken
		verificationtoken, vtErr := user.QueryVerificationToken().Only(context.Background())
		assert.NoError(t, vtErr)

		t.Run("confirm user email", func(t *testing.T) {
			// Test user email confirmation-token
			payload := types.ConfirmEmailPayload{
				Token: verificationtoken.Token,
			}

			res, err := test.PerformRequest(t, "POST", "/confirm-account?scope=sender", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Nil(t, response.Data)

			updateUser, uErr := verificationtoken.QueryOwner().Only(context.Background())
			assert.NoError(t, uErr)
			assert.Equal(t, true, updateUser.IsVerified)
		})
	})

	t.Run("Login", func(t *testing.T) {
		t.Run("with valid credentials", func(t *testing.T) {
			// Test login with valid credentials
			payload := types.LoginPayload{
				Email:    "ikeayo@example.com",
				Password: "password",
			}

			res, err := test.PerformRequest(t, "POST", "/login?scope=sender", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Successfully logged in", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response data
			assert.Contains(t, data, "accessToken")
			assert.NotEmpty(t, data["accessToken"].(string))
			assert.Contains(t, data, "refreshToken")
			assert.NotEmpty(t, data["refreshToken"].(string))
		})

		t.Run("with invalid credentials", func(t *testing.T) {
			// Test login with invalid credentials
			payload := types.LoginPayload{
				Email:    "ikeayo@example.com",
				Password: "wrong-password",
			}

			res, err := test.PerformRequest(t, "POST", "/login?scope=sender", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusUnauthorized, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid credentials", response.Message)
			assert.Equal(t, "Email and password do not match any user", response.Data)
		})
	})

	t.Run("RefreshJWT", func(t *testing.T) {
		t.Run("with a valid refresh token", func(t *testing.T) {
			refreshToken, err := token.GenerateRefreshJWT(userID)
			assert.NoError(t, err, "failed to generate refresh token")

			// Test refresh token with valid refresh token
			payload := types.RefreshJWTPayload{
				RefreshToken: refreshToken,
			}

			headers := map[string]string{
				"Authorization": "Bearer " + refreshToken,
			}

			res, err := test.PerformRequest(t, "POST", "/refresh?scope=sender", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Successfully refreshed access token", response.Message)
			data, ok := response.Data.(map[string]interface{})
			assert.True(t, ok, "response.Data is not of type map[string]interface{}")
			assert.NotNil(t, data, "response.Data is nil")

			// Assert the response data
			assert.Contains(t, data, "accessToken")
			assert.NotEmpty(t, data["accessToken"].(string))
			assert.NotContains(t, data, "refreshToken")
		})

		t.Run("with an invalid refresh token", func(t *testing.T) {
			refreshToken := "invalid-refresh-token"

			// Test refresh token with invalid refresh token
			payload := types.RefreshJWTPayload{
				RefreshToken: refreshToken,
			}

			refreshTokenForHeader, err := token.GenerateRefreshJWT(userID)
			assert.NoError(t, err, "failed to generate refresh token")

			headers := map[string]string{
				"Authorization": "Bearer " + refreshTokenForHeader,
			}
			res, err := test.PerformRequest(t, "POST", "/refresh?scope=sender", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusUnauthorized, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid or expired refresh token", response.Message)
		})
	})

	t.Run("ResendVerificationToken", func(t *testing.T) {
		// fetch user
		user, fetchUserErr := db.Client.User.Query().Where(user.IDEQ(uuid.MustParse(userID))).Only(context.Background())
		assert.NoError(t, fetchUserErr, "failed to fetch user by userID")

		_, err := user.Update().SetIsVerified(false).Save(context.Background())
		assert.NoError(t, err, "failed to set isVerified to false")

		t.Run("verification token should be resent", func(t *testing.T) {
			// construct resend verification token payload
			payload := types.ResendTokenPayload{
				Scope: verificationtoken.ScopeVerification.String(),
				Email: user.Email,
			}

			res, err := test.PerformRequest(t, "POST", "/resend-token?scope=sender", payload, nil, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			// verificationtokens should be two
			amount := user.QueryVerificationToken().CountX(context.Background())
			assert.Equal(t, 2, amount)
		})
	})
}
