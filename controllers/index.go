package controllers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/lockorderfulfillment"
	svc "github.com/paycrest/paycrest-protocol/services"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"
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
	tokenIsValid := u.ContainsString([]string{"USDT", "USDC"}, token) // TODO: fetch supported tokens from db
	if !tokenIsValid {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Token is not supported", nil)
	}

	fiatSymbol := ctx.Param("fiat")
	fiatIsValid := u.ContainsString([]string{"NGN"}, fiatSymbol) // TODO: fetch supported fiat currencies from db
	if !fiatIsValid {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Fiat currency is not supported", nil)
	}

	tokenAmount, err := decimal.NewFromString(ctx.Param("amount"))
	if err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid amount", nil)
		return
	}

	// Get redis keys for provision buckets
	keys, _, err := storage.RedisClient.Scan(ctx, uint64(0), "bucket_"+fiatSymbol+"_%d_%d", 100).Result()
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
		return
	}

	rateResponse := decimal.NewFromInt(0)

	// Scan through the buckets to find a matching rate
	for _, key := range keys {
		bucketData := strings.Split(key, "_")
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

		// Check if fiat amount is within the bucket range and set the rate
		if fiatAmount.GreaterThanOrEqual(minAmount) && fiatAmount.LessThanOrEqual(maxAmount) {
			rateResponse = rate
		}
	}

	if rateResponse.Equal(decimal.NewFromInt(0)) {
		// No rate found in the regular buckets, return market rate from a provider in the default bucket
		keys, _, err := storage.RedisClient.Scan(ctx, uint64(0), "bucket_"+fiatSymbol+"_default", 1).Result()
		if err != nil {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
			return
		}
		// Get rate of the topmost provider in the priority queue of the default bucket
		providerData, err := storage.RedisClient.LIndex(ctx, keys[0], 0).Result()
		if err != nil {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
			return
		}
		rateResponse, _ = decimal.NewFromString(strings.Split(providerData, ":")[1])
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Rates fetched successfully", rateResponse)
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

	// Get validator profile from the context
	validatorCtx, ok := ctx.Get("validator")
	if !ok {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	validator := validatorCtx.(*ent.ValidatorProfile)

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
