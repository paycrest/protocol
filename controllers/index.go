package controllers

import (
	"net/http"
	"slices"
	"strings"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/institution"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/token"
	svc "github.com/paycrest/protocol/services"
	orderSvc "github.com/paycrest/protocol/services/order"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	u "github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// Controller is the default controller for other endpoints
type Controller struct {
	orderService         types.OrderService
	priorityQueueService *svc.PriorityQueueService
}

// NewController creates a new instance of AuthController with injected services
func NewController() *Controller {
	return &Controller{
		orderService:         orderSvc.NewOrderEVM(),
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

	institutions, err := storage.Client.Institution.
		Query().
		Where(institution.HasFiatCurrencyWith(
			fiatcurrency.CodeEQ(currencyCode),
		)).
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to fetch institutions", nil)
		return
	}

	response := make([]types.SupportedInstitutions, 0, len(institutions))
	for _, institution := range institutions {
		response = append(response, types.SupportedInstitutions{
			Code: institution.Code,
			Name: institution.Name,
			Type: institution.Type,
		})
	}

	u.APIResponse(ctx, http.StatusOK, "success", "OK", response)
}

// GetTokenRate controller fetches the current rate of the cryptocurrency token against the fiat currency
func (ctrl *Controller) GetTokenRate(ctx *gin.Context) {
	// Parse path parameters
	token, err := storage.Client.Token.
		Query().
		Where(token.Symbol(strings.ToUpper(ctx.Param("token")))).
		First(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch token rate", nil)
		return
	}

	if token == nil {
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
	u.APIResponse(ctx, http.StatusOK, "success", "OK", config.CryptoConfig().AggregatorPublicKey)
}

// VerifyAccount controller verifies an account of a given institution
func (ctrl *Controller) VerifyAccount(ctx *gin.Context) {
	var payload types.VerifyAccountRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	institution, err := storage.Client.Institution.
		Query().
		Where(institution.CodeEQ(payload.Institution)).
		WithFiatCurrency().
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", []types.ErrorData{{
			Field:   "Institution",
			Message: "Institution is not supported",
		}})
		return
	}

	provider, err := storage.Client.ProviderProfile.
		Query().
		Where(
			providerprofile.HasCurrencyWith(
				fiatcurrency.CodeEQ(institution.Edges.FiatCurrency.Code),
			),
			providerprofile.HostIdentifierNotNil(),
			providerprofile.IsActiveEQ(true),
			providerprofile.IsAvailableEQ(true),
		).
		First(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to verify account", err.Error())
		return
	}

	res, err := fastshot.NewClient(provider.HostIdentifier).
		Config().SetTimeout(30 * time.Second).
		Build().POST("/verify_account").
		Body().AsJSON(payload).
		Send()
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to verify account", nil)
		return
	}

	data, err := u.ParseJSONResponse(res.RawResponse)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to verify account", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Account name was fetched successfully", data["data"].(string))
}

// GetLockPaymentOrderStatus controller fetches a payment order status by ID
func (ctrl *Controller) GetLockPaymentOrderStatus(ctx *gin.Context) {
	// Get order ID from the URL
	orderID := ctx.Param("id")

	// Fetch related payment orders from the database
	orders, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.GatewayIDEQ(orderID),
		).
		WithToken(func(tq *ent.TokenQuery) {
			tq.WithNetwork()
		}).
		WithTransactions().
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch order status", nil)
		return
	}

	var settlements []types.LockPaymentOrderSplitOrder
	var receipts []types.LockPaymentOrderTxReceipt
	var settlePercent decimal.Decimal
	var totalAmount decimal.Decimal

	for _, order := range orders {
		for _, transaction := range order.Edges.Transactions {
			if u.ContainsString([]string{"order_settled", "order_created", "order_refunded"}, transaction.Status.String()) {
				var status lockpaymentorder.Status
				if transaction.Status.String() == "order_created" {
					status = lockpaymentorder.StatusPending
				} else {
					status = lockpaymentorder.Status(strings.TrimPrefix(transaction.Status.String(), "order_"))
				}
				receipts = append(receipts, types.LockPaymentOrderTxReceipt{
					Status:    status,
					TxHash:    transaction.TxHash,
					Timestamp: transaction.CreatedAt,
				})
			}
		}

		settlements = append(settlements, types.LockPaymentOrderSplitOrder{
			SplitOrderID: order.ID,
			Amount:       order.Amount,
			Rate:         order.Rate,
			OrderPercent: order.OrderPercent,
		})

		settlePercent = settlePercent.Add(order.OrderPercent)
		totalAmount = totalAmount.Add(order.Amount)
	}

	// Sort receipts by latest timestamp
	slices.SortStableFunc(receipts, func(a, b types.LockPaymentOrderTxReceipt) int {
		return b.Timestamp.Compare(a.Timestamp)
	})

	if (len(orders) == 0) || (len(receipts) == 0) {
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order not found", nil)
		return
	}

	response := &types.LockPaymentOrderStatusResponse{
		OrderID:       orders[0].GatewayID,
		Amount:        totalAmount,
		Token:         orders[0].Edges.Token.Symbol,
		Network:       orders[0].Edges.Token.Edges.Network.Identifier,
		SettlePercent: settlePercent,
		Status:        orders[0].Status,
		TxHash:        receipts[0].TxHash,
		Settlements:   settlements,
		TxReceipts:    receipts,
		UpdatedAt:     orders[0].UpdatedAt,
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Order status fetched successfully", response)
}
