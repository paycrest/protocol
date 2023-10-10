package controllers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/lockorderfulfillment"
	"github.com/paycrest/paycrest-protocol/ent/validatorprofile"
	svc "github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"
	"github.com/paycrest/paycrest-protocol/utils"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// Controller is the default controller for other endpoints
type Controller struct {
	orderService *svc.OrderService
}

// NewController creates a new instance of AuthController with injected services
func NewController() *Controller {
	return &Controller{
		orderService: svc.NewOrderService(),
	}
}

// GetFiatCurrencies controller fetches the supported fiat currencies
func (ctrl *Controller) GetFiatCurrencies(ctx *gin.Context) {
	// fetch stored fiat currencies.
	fiatcurrencies, err := storage.Client.FiatCurrency.Query().All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to fetch FiatCurrencies", err.Error())
		return
	}

	currencies := make([]types.SupportedCurrencies, 0, len(fiatcurrencies))
	for _, currency := range fiatcurrencies {
		currencies = append(currencies, types.SupportedCurrencies{
			Code:      currency.Code,
			Name:      currency.Name,
			ShortName: currency.ShortName,
			Decimals:  int8(currency.Decimals),
			Symbol:    currency.Symbol,
		})
	}

	u.APIResponse(ctx, http.StatusOK, "success", "OK", currencies)
}

// GetInstitutionsByCurrency controller fetches the supported institutions for a given currency
func (ctrl *Controller) GetInstitutionsByCurrency(ctx *gin.Context) {
	// Get currency code from the URL
	currencyCode := ctx.Param("currency_code")

	institutions, err := ctrl.orderService.GetSupportedInstitution(ctx, nil, currencyCode)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to fetch institutions", err.Error())
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "OK", institutions)
}

// GetRates controller fetches the current market rates for the supported cryptocurrencies
func (ctrl *Controller) GetTokenRates(ctx *gin.Context) {
	// Parse path parameters
	token := ctx.Param("token")
	tokenIsValid := utils.ContainsString([]string{"USDT", "USDC"}, token)
	if !tokenIsValid {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Token is not supported", nil)
	}

	tokenAmount, err := decimal.NewFromString(ctx.Param("amount"))
	if err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid amount", nil)
		return
	}

	rates := map[string]decimal.Decimal{}
	cursor := 0

scanLoop:
	for {
		// Get redis keys for provision buckets
		keys, c, err := storage.RedisClient.Scan(ctx, uint64(cursor), "bucket_%s_%d_%d", 100).Result()
		if err != nil {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
			return
		}

		for _, key := range keys {
			bucketData := strings.Split(key, "_")
			fiatCurrency := bucketData[1]
			minAmount, _ := decimal.NewFromString(bucketData[2])
			maxAmount, _ := decimal.NewFromString(bucketData[3])

			// Get the topmost provider in the priority queue of the bucket
			providerData, err := storage.RedisClient.LIndex(ctx, key, 0).Result()
			if err != nil {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
				return
			}

			// Get fiat equivalent of the token amount
			rate, _ := decimal.NewFromString(strings.Split(providerData, ":")[1])
			fiatAmount := tokenAmount.Mul(rate)

			// Check if fiat amount is within the bucket range
			if fiatAmount.GreaterThanOrEqual(minAmount) && fiatAmount.LessThanOrEqual(maxAmount) {
				// Assign rate to the fiat currency
				rates[fiatCurrency] = rate

				break scanLoop
			}
		}

		if len(rates) == 0 {
			// No rate found in the regular buckets, return market rate from a provider in the default bucket
			keys, _, err := storage.RedisClient.Scan(ctx, uint64(cursor), "bucket_%s_default", 100).Result()
			if err != nil {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
				return
			}

			for _, key := range keys {
				// Get rate of the topmost provider in the priority queue of the default bucket
				providerData, err := storage.RedisClient.LIndex(ctx, key, 0).Result()
				if err != nil {
					u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
					return
				}
				rate, _ := decimal.NewFromString(strings.Split(providerData, ":")[1])

				// Assign rate to the fiat currency
				fiatCurrency := strings.Split(key, "_")[1]
				rates[fiatCurrency] = rate

				break scanLoop
			}
		}

		// Break when cursor is back to 0
		cursor = int(c)
		if cursor == 0 {
			break
		}
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Rates fetched successfully", rates)
}

// ValidateOrder is a hook to receive validation decisions from validators
func (ctrl *Controller) ValidateOrder(ctx *gin.Context) {
	var payload types.ValidateOrderPayload

	// Parse the payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get lock order fulfillment ID from URL
	fulfillmentID := ctx.Param("fulfillment_id")

	// Parse the order fulfillment ID string into a UUID
	fulfillmentUUID, err := uuid.Parse(fulfillmentID)
	if err != nil {
		logger.Errorf("error parsing fulfillment ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid fulfillment ID", nil)
		return
	}

	// Get the api key ID from the context
	apiKey, _ := ctx.Get("api_key")

	// Fetch validator from db
	validator, err := storage.Client.ValidatorProfile.
		Query().
		Where(
			validatorprofile.HasAPIKeyWith(
				apikey.ID(apiKey.(*ent.APIKey).ID),
			),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Could not find validator profile", nil)
		return
	}

	// Fetch order fulfillment from db
	fulfillment, err := storage.Client.LockOrderFulfillment.
		Query().
		Where(
			lockorderfulfillment.IDEQ(fulfillmentUUID),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Could not find order fulfillment", nil)
		return
	}

	// Update lock order fulfillment status
	if payload.IsValid {
		_, err = fulfillment.Update().
			SetConfirmations(fulfillment.Confirmations + 1).
			AddValidatorIDs(validator.ID).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to validate order fulfillment", nil)
			return
		}
	} else {
		_, err = fulfillment.Update().
			AppendValidationErrors([]string{payload.ErrorMsg}).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to validate order fulfillment", nil)
			return
		}
	}

	ctx.JSON(http.StatusOK, "OK")
}
