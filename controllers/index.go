package controllers

import (
	"net/http"
	"strings"

	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/token"
	svc "github.com/paycrest/protocol/services"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	u "github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// Controller is the default controller for other endpoints
type Controller struct {
	orderService         svc.Order
	priorityQueueService *svc.PriorityQueueService
}

// NewController creates a new instance of AuthController with injected services
func NewController() *Controller {
	return &Controller{
		orderService:         svc.NewOrderService(),
		priorityQueueService: svc.NewPriorityQueueService(),
	}
}

// GetFiatCurrencies controller fetches the supported fiat currencies
func (ctrl *Controller) GetFiatCurrencies(ctx *gin.Context) {
	// fetch stored fiat currencies.
	fiatcurrencies, err := storage.Client.FiatCurrency.
		Query().
		Where(fiatcurrency.IsEnabledEQ(true)).
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to fetch FiatCurrencies", err.Error())
		return
	}

	currencies := make([]types.SupportedCurrencies, 0, len(fiatcurrencies))
	for _, currency := range fiatcurrencies {
		currencies = append(currencies, types.SupportedCurrencies{
			Code:       currency.Code,
			Name:       currency.Name,
			ShortName:  currency.ShortName,
			Decimals:   int8(currency.Decimals),
			Symbol:     currency.Symbol,
			MarketRate: currency.MarketRate,
		})
	}

	u.APIResponse(ctx, http.StatusOK, "success", "OK", currencies)
}

// GetInstitutionsByCurrency controller fetches the supported institutions for a given currency
func (ctrl *Controller) GetInstitutionsByCurrency(ctx *gin.Context) {
	// Get currency code from the URL
	currencyCode := ctx.Param("currency_code")

	institutions, err := ctrl.orderService.GetSupportedInstitutions(ctx, nil, strings.ToUpper(currencyCode))
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to fetch institutions", err.Error())
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "OK", institutions)
}

// GetTokenRate controller fetches the current rate of the cryptocurrency token against the fiat currency
func (ctrl *Controller) GetTokenRate(ctx *gin.Context) {
	// Parse path parameters
	_, err := storage.Client.Token.
		Query().
		Where(token.Symbol(strings.ToUpper(ctx.Param("token")))).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Token is not supported", nil)
		return
	}

	currency, err := storage.Client.FiatCurrency.
		Query().
		Where(
			fiatcurrency.IsEnabledEQ(true),
			fiatcurrency.CodeEQ(strings.ToUpper(ctx.Param("fiat"))),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Fiat currency is not supported", nil)
		return
	}

	tokenAmount, err := decimal.NewFromString(ctx.Param("amount"))
	if err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid amount", nil)
		return
	}

	rateResponse := decimal.NewFromInt(0)

	// get providerID from query params
	providerID := ctx.Query("provider_id")
	if providerID != "" {
		// get the provider from the bucket
		provider, err := storage.Client.ProviderProfile.
			Query().
			Where(providerprofile.IDEQ(providerID)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				u.APIResponse(ctx, http.StatusBadRequest, "error", "Provider not found", nil)
				return
			} else {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider profile", nil)
				return
			}
		}

		rateResponse, err = ctrl.priorityQueueService.GetProviderRate(ctx, provider)
		if err != nil {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider rate", nil)
			return
		}

	} else {
		// Get redis keys for provision buckets
		keys, _, err := storage.RedisClient.Scan(ctx, uint64(0), "bucket_"+currency.Code+"_*_*", 100).Result()
		if err != nil {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
			return
		}

		highestMaxAmount := decimal.NewFromInt(0)

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
				break
			} else {
				// Get the highest max amount
				if maxAmount.GreaterThan(highestMaxAmount) {
					highestMaxAmount = maxAmount
					rateResponse = rate
				}
			}
		}
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Rate fetched successfully", rateResponse)
}

// GetAggregatorPublicKey controller expose Aggregator Public Key
func (ctrl *Controller) GetAggregatorPublicKey(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", map[string]interface{}{
		"aggregatorPublicKey": config.CryptoConfig().AggregatorPublicKey,
	})
}
