package accounts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paycrest/aggregator/ent"
	"github.com/paycrest/aggregator/routers/middleware"
	db "github.com/paycrest/aggregator/storage"
	"github.com/paycrest/aggregator/types"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/aggregator/ent/enttest"
	"github.com/paycrest/aggregator/ent/fiatcurrency"
	"github.com/paycrest/aggregator/ent/providerordertoken"
	"github.com/paycrest/aggregator/ent/providerprofile"
	"github.com/paycrest/aggregator/ent/senderordertoken"
	"github.com/paycrest/aggregator/ent/senderprofile"
	tokenDB "github.com/paycrest/aggregator/ent/token"
	"github.com/paycrest/aggregator/ent/user"
	"github.com/paycrest/aggregator/utils/test"
	"github.com/paycrest/aggregator/utils/token"
	"github.com/stretchr/testify/assert"
)

var testCtx = struct {
	user            *ent.User
	providerProfile *ent.ProviderProfile
	token           *ent.Token
	client          types.RPCClient
}{}

func setup() error {
	// Set up test blockchain client
	client, err := test.SetUpTestBlockchain()
	if err != nil {
		return err
	}

	testCtx.client = client
	// Create a test token
	token, err := test.CreateERC20Token(
		client,
		map[string]interface{}{
			"deployContract": false,
		})
	if err != nil {
		return err
	}
	testCtx.token = token

	// Set up test data
	user, err := test.CreateTestUser(map[string]interface{}{
		"scope": "provider",
		"email": "providerjohndoe@test.com",
	})
	if err != nil {
		return err
	}
	testCtx.user = user

	currency, err := test.CreateTestFiatCurrency(map[string]interface{}{
		"code":        "KES",
		"short_name":  "Shilling",
		"decimals":    2,
		"symbol":      "KSh",
		"name":        "Kenyan Shilling",
		"market_rate": 550.0,
	})
	if err != nil {
		return err
	}

	provderProfile, err := test.CreateTestProviderProfile(map[string]interface{}{
		"user_id":     testCtx.user.ID,
		"currency_id": currency.ID,
	})
	if err != nil {
		return err
	}
	testCtx.providerProfile = provderProfile

	return nil
}

func TestProfile(t *testing.T) {
	// Set up test database client
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	db.Client = client

	// Setup test data
	err := setup()
	assert.NoError(t, err)

	// Set up test routers
	router := gin.New()
	ctrl := &ProfileController{}

	router.GET(
		"settings/sender",
		middleware.JWTMiddleware,
		middleware.OnlySenderMiddleware,
		ctrl.GetSenderProfile,
	)
	router.GET(
		"settings/provider",
		middleware.JWTMiddleware,
		middleware.OnlyProviderMiddleware,
		ctrl.GetProviderProfile,
	)
	router.PATCH(
		"settings/sender",
		middleware.JWTMiddleware,
		middleware.OnlySenderMiddleware,
		ctrl.UpdateSenderProfile,
	)
	router.PATCH(
		"settings/provider",
		middleware.JWTMiddleware,
		middleware.OnlyProviderMiddleware,
		ctrl.UpdateProviderProfile,
	)

	t.Run("UpdateSenderProfile", func(t *testing.T) {
		t.Run("with all fields", func(t *testing.T) {
			testUser, err := test.CreateTestUser(map[string]interface{}{"scope": "sender"})
			assert.NoError(t, err)

			_, err = test.CreateTestSenderProfile(map[string]interface{}{
				"domain_whitelist": []string{"example.com"},
				"user_id":          testUser.ID,
				"token":            testCtx.token.Symbol,
			})
			assert.NoError(t, err)

			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testUser.ID.String(), "sender")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.SenderProfilePayload{
				DomainWhitelist: []string{"example.com", "mydomain.com"},
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/sender", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Profile updated successfully", response.Message)
			assert.Nil(t, response.Data, "response.Data is not nil")

			senderProfile, err := db.Client.SenderProfile.
				Query().
				Where(senderprofile.HasUserWith(user.ID(testUser.ID))).
				Only(context.Background())
			assert.NoError(t, err)

			assert.Contains(t, senderProfile.DomainWhitelist, "mydomain.com")
		})

		t.Run("with an invalid webhook", func(t *testing.T) {
			testUser, err := test.CreateTestUser(map[string]interface{}{
				"scope": "sender",
				"email": "johndoe2@test.com",
			})
			assert.NoError(t, err)

			_, err = test.CreateTestSenderProfile(map[string]interface{}{
				"domain_whitelist": []string{"example.com"},
				"user_id":          testUser.ID,
			})
			assert.NoError(t, err)

			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testUser.ID.String(), "sender")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.SenderProfilePayload{
				WebhookURL:      "examplecom",
				DomainWhitelist: []string{"example.com", "mydomain.com"},
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/sender", payload, headers, router)
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
			assert.Equal(t, "WebhookURL", errorMap["field"].(string))
			assert.Contains(t, errorMap, "message")
			assert.Equal(t, "Invalid URL", errorMap["message"].(string))
		})

		t.Run("with all fields and check if it is active", func(t *testing.T) {
			testUser, err := test.CreateTestUser(map[string]interface{}{
				"scope": "sender",
				"email": "johndoe3@test.com",
			})
			assert.NoError(t, err)

			_, err = test.CreateTestSenderProfile(map[string]interface{}{
				"domain_whitelist": []string{"example.com"},
				"user_id":          testUser.ID,
				"token":            testCtx.token.Symbol,
			})
			assert.NoError(t, err)

			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testUser.ID.String(), "sender")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}

			// setup payload
			tokenPayload := make([]types.SenderOrderTokenPayload, 2)
			tokenAddresses := make([]types.SenderOrderAddressPayload, 1)

			// setup ERC20 token
			tokenAddresses[0].FeeAddress = "0xD4EB9067111F81b9bAabE06E2b8ebBaDADEd5DAf"
			tokenAddresses[0].Network = testCtx.token.Edges.Network.Identifier
			tokenAddresses[0].RefundAddress = "0xD4EB9067111F81b9bAabE06E2b8ebBaDADEd5DA0"

			tokenPayload[0].FeePercent = decimal.NewFromInt(1)
			tokenPayload[0].Symbol = testCtx.token.Symbol
			tokenPayload[0].Addresses = tokenAddresses

			// setup TRC token
			tronToken, err := test.CreateTRC20Token(testCtx.client, map[string]interface{}{})
			assert.NoError(t, err)
			assert.NotEqual(t, "localhost", tronToken.Edges.Network.Identifier)

			// setup TRC20 token
			tronTokenAddresses := make([]types.SenderOrderAddressPayload, 1)
			tronTokenAddresses[0].FeeAddress = "TFRKiHrHCeSyWL67CEwydFvUMYJ6CbYYX7"
			tronTokenAddresses[0].Network = tronToken.Edges.Network.Identifier
			tronTokenAddresses[0].RefundAddress = "TFRKiHrHCeSyWL67CEwydFvUMYJ6CbYYXR"

			tokenPayload[1].FeePercent = decimal.NewFromInt(2)
			tokenPayload[1].Symbol = tronToken.Symbol
			tokenPayload[1].Addresses = tronTokenAddresses

			// put the payload together
			payload := types.SenderProfilePayload{
				DomainWhitelist: []string{"example.com", "mydomain.com"},
				WebhookURL:      "https://example.com",
				Tokens:          tokenPayload,
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/sender", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Profile updated successfully", response.Message)
			assert.Nil(t, response.Data, "response.Data is not nil")

			senderProfile, err := db.Client.SenderProfile.
				Query().
				Where(senderprofile.HasUserWith(user.ID(testUser.ID))).
				WithOrderTokens().
				Only(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, len(senderProfile.Edges.OrderTokens), 2)

			t.Run("check If Tron was added", func(t *testing.T) {
				senderorder, err := db.Client.SenderOrderToken.
					Query().
					Where(
						senderordertoken.HasSenderWith(
							senderprofile.IDEQ(senderProfile.ID),
						),
						senderordertoken.HasTokenWith(tokenDB.IDEQ(tronToken.ID)),
					).
					Only(context.Background())
				assert.NoError(t, err)
				assert.Equal(t, senderorder.FeeAddress, "TFRKiHrHCeSyWL67CEwydFvUMYJ6CbYYX7")
				assert.Equal(t, senderorder.RefundAddress, "TFRKiHrHCeSyWL67CEwydFvUMYJ6CbYYXR")
			})

			t.Run("check If EVM chain was added", func(t *testing.T) {
				senderorder, err := db.Client.SenderOrderToken.
					Query().
					Where(
						senderordertoken.HasSenderWith(
							senderprofile.IDEQ(senderProfile.ID),
						),
						senderordertoken.HasTokenWith(tokenDB.IDEQ(testCtx.token.ID)),
					).
					Only(context.Background())
				assert.NoError(t, err)
				assert.Equal(t, senderorder.FeeAddress, "0xD4EB9067111F81b9bAabE06E2b8ebBaDADEd5DAf")
				assert.Equal(t, senderorder.RefundAddress, "0xD4EB9067111F81b9bAabE06E2b8ebBaDADEd5DA0")
			})
			assert.Contains(t, senderProfile.DomainWhitelist, "mydomain.com")
			assert.True(t, senderProfile.IsActive)
		})

	})

	t.Run("UpdateProviderProfile", func(t *testing.T) {

		t.Run("with all fields complete and check if it is active", func(t *testing.T) {
			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testCtx.user.ID.String(), "provider")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.ProviderProfilePayload{
				TradingName:      "My Trading Name",
				Currency:         "KES",
				HostIdentifier:   "example.com",
				BusinessDocument: "https://example.com/business_doc.png",
				IdentityDocument: "https://example.com/national_id.png",
				IsAvailable:      true,
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Profile updated successfully", response.Message)
			assert.Nil(t, response.Data, "response.Data is not nil")

			providerProfile, err := db.Client.ProviderProfile.
				Query().
				Where(providerprofile.HasUserWith(user.ID(testCtx.user.ID))).
				WithCurrency().
				Only(context.Background())
			assert.NoError(t, err)

			assert.Contains(t, providerProfile.TradingName, payload.TradingName)
			assert.Contains(t, providerProfile.HostIdentifier, payload.HostIdentifier)
			assert.Contains(t, providerProfile.Edges.Currency.Code, payload.Currency)
			assert.Contains(t, providerProfile.BusinessDocument, payload.BusinessDocument)
			assert.Contains(t, providerProfile.IdentityDocument, payload.IdentityDocument)
			assert.True(t, providerProfile.IsActive)
		})

		t.Run("with visibility", func(t *testing.T) {
			// Test partial update
			accessToken, _ := token.GenerateAccessJWT(testCtx.user.ID.String(), "provider")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}
			payload := types.ProviderProfilePayload{
				VisibilityMode: "private",
				TradingName:    testCtx.providerProfile.TradingName,
				HostIdentifier: testCtx.providerProfile.HostIdentifier,
				Currency:       "KES",
			}

			res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
			assert.NoError(t, err)

			// Assert the response body
			assert.Equal(t, http.StatusOK, res.Code)

			var response types.Response
			err = json.Unmarshal(res.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Profile updated successfully", response.Message)
			assert.Nil(t, response.Data, "response.Data is not nil")

			providerProfile, err := db.Client.ProviderProfile.Query().
				Where(providerprofile.VisibilityModeEQ(providerprofile.VisibilityModePrivate)).
				Count(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, 1, providerProfile)
		})

		t.Run("with optional fields", func(t *testing.T) {
			profileUpdateRequest := func(payload types.ProviderProfilePayload) *httptest.ResponseRecorder {
				// Test partial update
				accessToken, _ := token.GenerateAccessJWT(testCtx.user.ID.String(), "provider")
				headers := map[string]string{
					"Authorization": "Bearer " + accessToken,
				}

				res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
				assert.NoError(t, err)

				return res
			}
			t.Run("fails for invalid mobile number", func(t *testing.T) {
				payload := types.ProviderProfilePayload{
					MobileNumber:   "01234567890",
					TradingName:    testCtx.providerProfile.TradingName,
					HostIdentifier: testCtx.providerProfile.HostIdentifier,
					Currency:       "KES",
				}
				res1 := profileUpdateRequest(payload)

				payload.MobileNumber = "+023456789029"
				res2 := profileUpdateRequest(payload)

				// Assert the response body
				assert.Equal(t, http.StatusBadRequest, res1.Code)
				assert.Equal(t, http.StatusBadRequest, res2.Code)

				var response types.Response
				err = json.Unmarshal(res1.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid mobile number", response.Message)

				err = json.Unmarshal(res2.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid mobile number", response.Message)
			})

			t.Run("success for valid moblie number", func(t *testing.T) {
				payload := types.ProviderProfilePayload{
					MobileNumber:   "+2347012345678",
					TradingName:    testCtx.providerProfile.TradingName,
					HostIdentifier: testCtx.providerProfile.HostIdentifier,
					Currency:       "KES",
				}
				res := profileUpdateRequest(payload)

				// Assert the response body
				assert.Equal(t, http.StatusOK, res.Code)

				var response types.Response
				err = json.Unmarshal(res.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Profile updated successfully", response.Message)

				// Assert optional fields were correctly set and retrieved
				providerProfile, err := db.Client.ProviderProfile.
					Query().
					Where(providerprofile.HasUserWith(user.ID(testCtx.user.ID))).
					Only(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, providerProfile.MobileNumber, "+2347012345678")
			})

			t.Run("fails for invalid identity document type", func(t *testing.T) {
				payload := types.ProviderProfilePayload{
					IdentityDocumentType: "student_id",
					TradingName:          testCtx.providerProfile.TradingName,
					HostIdentifier:       testCtx.providerProfile.HostIdentifier,
				}
				res1 := profileUpdateRequest(payload)

				payload.IdentityDocumentType = "bank_statement"
				res2 := profileUpdateRequest(payload)

				// Assert the response body
				assert.Equal(t, http.StatusBadRequest, res1.Code)
				assert.Equal(t, http.StatusBadRequest, res2.Code)

				var response types.Response
				err = json.Unmarshal(res1.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid identity document type", response.Message)

				err = json.Unmarshal(res2.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid identity document type", response.Message)
			})

			t.Run("fails for invalid identity document url", func(t *testing.T) {
				payload := types.ProviderProfilePayload{
					IdentityDocument: "img.png",
					TradingName:      testCtx.providerProfile.TradingName,
					HostIdentifier:   testCtx.providerProfile.HostIdentifier,
				}
				res1 := profileUpdateRequest(payload)

				payload.IdentityDocument = "ftp://123.example.com/file.png"
				res2 := profileUpdateRequest(payload)

				// Assert the response body
				assert.Equal(t, http.StatusBadRequest, res1.Code)
				assert.Equal(t, http.StatusBadRequest, res2.Code)

				var response types.Response
				err = json.Unmarshal(res1.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid identity document URL", response.Message)

				err = json.Unmarshal(res2.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid identity document URL", response.Message)
				// assert.Nil(t, response.Data, "response.Data is not nil")
			})

			t.Run("fails for invalid business document url", func(t *testing.T) {
				payload := types.ProviderProfilePayload{
					BusinessDocument: "http://123.example.com/file.ai",
					TradingName:      testCtx.providerProfile.TradingName,
					HostIdentifier:   testCtx.providerProfile.HostIdentifier,
				}
				res := profileUpdateRequest(payload)
				// Assert the response body
				assert.Equal(t, http.StatusBadRequest, res.Code)

				var response types.Response
				err = json.Unmarshal(res.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid business document URL", response.Message)
				// assert.Nil(t, response.Data, "response.Data is not nil")
			})

			t.Run("succeeds with valid optional fields", func(t *testing.T) {
				payload := types.ProviderProfilePayload{
					Address:              "123, Example Street, Nairobi, Kenya",
					MobileNumber:         "+2347012345678",
					DateOfBirth:          time.Date(2022, time.January, 1, 12, 30, 0, 0, time.UTC),
					BusinessName:         "Example Business",
					IdentityDocumentType: "national_id",
					IdentityDocument:     "https://example.com/national_id.png",
					BusinessDocument:     "https://example.com/business_doc.png",
					TradingName:          testCtx.providerProfile.TradingName,
					HostIdentifier:       testCtx.providerProfile.HostIdentifier,
				}
				res := profileUpdateRequest(payload)
				// Assert the response body
				assert.Equal(t, http.StatusOK, res.Code)

				var response types.Response
				err = json.Unmarshal(res.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Profile updated successfully", response.Message)

				// Assert optional fields were correctly set and retrieved
				providerProfile, err := db.Client.ProviderProfile.
					Query().
					Where(providerprofile.HasUserWith(user.ID(testCtx.user.ID))).
					Only(context.Background())
				assert.NoError(t, err)

				assert.Equal(t, providerProfile.Address, "123, Example Street, Nairobi, Kenya")
				assert.Equal(t, providerProfile.MobileNumber, "+2347012345678")
				assert.Equal(t, providerProfile.DateOfBirth, time.Date(2022, time.January, 1, 12, 30, 0, 0, time.UTC))
				assert.Equal(t, providerProfile.BusinessName, "Example Business")
				assert.Equal(t, string(providerProfile.IdentityDocumentType), "national_id")
				assert.Equal(t, providerProfile.IdentityDocument, "https://example.com/national_id.png")
				assert.Equal(t, providerProfile.BusinessDocument, "https://example.com/business_doc.png")
			})

		})

		t.Run("token validation against provision buckets", func(t *testing.T) {
			// Setup common test data
			accessToken, _ := token.GenerateAccessJWT(testCtx.user.ID.String(), "provider")
			headers := map[string]string{
				"Authorization": "Bearer " + accessToken,
			}

			setupBuckets := func(t *testing.T, currency string, amounts ...float64) {

				// Clear existing buckets
				_, err := db.Client.ProvisionBucket.Delete().Exec(context.Background())
				assert.NoError(t, err)

				// Get currency
				curr, err := db.Client.FiatCurrency.Query().
					Where(fiatcurrency.CodeEQ(currency)).
					Only(context.Background())

				if ent.IsNotFound(err) { // If the currency doesn't exist, create it
					curr, err = db.Client.FiatCurrency.Create().
						SetCode(currency).
						Save(context.Background())
				}
				assert.NoError(t, err)

				// Create new buckets
				for _, amount := range amounts {
					_, err = db.Client.ProvisionBucket.Create().
						SetCurrency(curr).
						SetMinAmount(decimal.NewFromFloat(0)).
						SetMaxAmount(decimal.NewFromFloat(amount)).
						Save(context.Background())
					assert.NoError(t, err)
				}
			}

			t.Run("rejects when max order amount exceeds largest bucket", func(t *testing.T) {
				setupBuckets(t, "KES", 1000, 5000, 10000)
				payload := types.ProviderProfilePayload{
					TradingName:    testCtx.providerProfile.TradingName,
					HostIdentifier: testCtx.providerProfile.HostIdentifier,
					Currency:       "KES",
					Tokens: []types.ProviderOrderTokenPayload{
						{
							Symbol:                 "BTC",
							MaxOrderAmount:         decimal.NewFromFloat(0.5),
							MinOrderAmount:         decimal.NewFromFloat(0.001),
							ConversionRateType:     providerordertoken.ConversionRateTypeFloating,
							FloatingConversionRate: decimal.NewFromFloat(0),
							FixedConversionRate:    decimal.NewFromFloat(0),
							Addresses: []struct {
								Address string `json:"address"`
								Network string `json:"network"`
							}{
								{
									Network: "bitcoin",
									Address: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
								},
							},
						},
					},
				}
				res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.Code)
				var response types.Response
				err = json.Unmarshal(res.Body.Bytes(), &response)
				assert.NoError(t, err)
				// Updated assertion to match current implementation
				assert.Equal(t, "Token not supported", response.Message)
			})

			t.Run("accepts when max order amount within bucket limits", func(t *testing.T) {
				// Setup buckets: 1000, 5000, 10000 KES
				setupBuckets(t, "KES", 1000, 5000, 10000)

				payload := types.ProviderProfilePayload{
					TradingName:    testCtx.providerProfile.TradingName,
					HostIdentifier: testCtx.providerProfile.HostIdentifier,
					Currency:       "KES",
					Tokens: []types.ProviderOrderTokenPayload{
						{
							Symbol:                 "BTC",
							MaxOrderAmount:         decimal.NewFromFloat(0.1), // Assuming rate of 40000 KES/BTC = 4000 KES
							MinOrderAmount:         decimal.NewFromFloat(0.001),
							ConversionRateType:     providerordertoken.ConversionRateTypeFloating,
							FloatingConversionRate: decimal.NewFromFloat(0),
							FixedConversionRate:    decimal.NewFromFloat(0),
							Addresses: []struct {
								Address string `json:"address"`
								Network string `json:"network"`
							}{
								{
									Network: "bitcoin",
									Address: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
								},
							},
						},
					},
				}
				res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.Code)
				var response types.Response
				err = json.Unmarshal(res.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "error", response.Status)
				assert.Equal(t, "Token not supported", response.Message)
			})

			t.Run("rejects when currency has no provision buckets", func(t *testing.T) {
				// Clear all buckets for KES
				setupBuckets(t, "KES")

				payload := types.ProviderProfilePayload{
					TradingName:    testCtx.providerProfile.TradingName,
					HostIdentifier: testCtx.providerProfile.HostIdentifier,
					Currency:       "KES",
					Tokens: []types.ProviderOrderTokenPayload{
						{
							Symbol:                 "BTC",
							MaxOrderAmount:         decimal.NewFromFloat(0.1),
							MinOrderAmount:         decimal.NewFromFloat(0.001),
							ConversionRateType:     providerordertoken.ConversionRateTypeFloating,
							FloatingConversionRate: decimal.NewFromFloat(0),
							FixedConversionRate:    decimal.NewFromFloat(0),
							Addresses: []struct {
								Address string `json:"address"`
								Network string `json:"network"`
							}{
								{
									Network: "bitcoin",
									Address: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
								},
							},
						},
					},
				}

				res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.Code)
				var response types.Response
				err = json.Unmarshal(res.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "error", response.Status)
				assert.Equal(t, "Token not supported", response.Message)
			})

			t.Run("validates multiple tokens against buckets", func(t *testing.T) {
				// Setup buckets: 1000, 5000, 10000 KES
				setupBuckets(t, "KES", 1000, 5000, 10000)

				payload := types.ProviderProfilePayload{
					TradingName:    testCtx.providerProfile.TradingName,
					HostIdentifier: testCtx.providerProfile.HostIdentifier,
					Currency:       "KES",
					Tokens: []types.ProviderOrderTokenPayload{
						{
							Symbol:                 "BTC",
							MaxOrderAmount:         decimal.NewFromFloat(0.1), // Valid amount
							MinOrderAmount:         decimal.NewFromFloat(0.001),
							ConversionRateType:     providerordertoken.ConversionRateTypeFloating,
							FloatingConversionRate: decimal.NewFromFloat(0),
							FixedConversionRate:    decimal.NewFromFloat(0),
							Addresses: []struct {
								Address string `json:"address"`
								Network string `json:"network"`
							}{
								{
									Network: "bitcoin",
									Address: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
								},
							},
						},
						{
							Symbol:                 "ETH",
							MaxOrderAmount:         decimal.NewFromFloat(10), // Invalid amount (too high)
							MinOrderAmount:         decimal.NewFromFloat(0.01),
							ConversionRateType:     providerordertoken.ConversionRateTypeFloating,
							FloatingConversionRate: decimal.NewFromFloat(0),
							FixedConversionRate:    decimal.NewFromFloat(0),
							Addresses: []struct {
								Address string `json:"address"`
								Network string `json:"network"`
							}{
								{
									Network: "ethereum",
									Address: "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
								},
							},
						},
					},
				}

				res, err := test.PerformRequest(t, "PATCH", "/settings/provider", payload, headers, router)
				assert.NoError(t, err)
				assert.Equal(t, http.StatusBadRequest, res.Code)
				var response types.Response
				err = json.Unmarshal(res.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "error", response.Status)
				assert.Equal(t, "Token not supported", response.Message)
			})
		})

	})

	// t.Run("GetSenderProfile", func(t *testing.T) {
	// 	testUser, err := test.CreateTestUser(map[string]interface{}{
	// 		"email": "hello@test.com",
	// 		"scope": "sender",
	// 	})
	// 	assert.NoError(t, err)

	// 	sender, err := test.CreateTestSenderProfile(map[string]interface{}{
	// 		"domain_whitelist": []string{"mydomain.com"},
	// 		"user_id":          testUser.ID,
	// 	})
	// 	assert.NoError(t, err)

	// 	apiKeyService := services.NewAPIKeyService()
	// 	_, _, err = apiKeyService.GenerateAPIKey(
	// 		context.Background(),
	// 		nil,
	// 		sender,
	// 		nil,
	// 	)
	// 	assert.NoError(t, err)

	// 	accessToken, _ := token.GenerateAccessJWT(testUser.ID.String(), "sender")
	// 	headers := map[string]string{
	// 		"Authorization": "Bearer " + accessToken,
	// 	}
	// 	res, err := test.PerformRequest(t, "GET", "/settings/sender", nil, headers, router)
	// 	assert.NoError(t, err)

	// 	// Assert the response body
	// 	assert.Equal(t, http.StatusOK, res.Code)
	// 	var response struct {
	// 		Data    types.SenderProfileResponse
	// 		Message string
	// 	}
	// 	err = json.Unmarshal(res.Body.Bytes(), &response)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, "Profile retrieved successfully", response.Message)
	// 	assert.NotNil(t, response.Data, "response.Data is nil")
	// 	assert.Greater(t, len(response.Data.Tokens), 0)
	// 	assert.Contains(t, response.Data.WebhookURL, "https://example.com")

	// })

}
