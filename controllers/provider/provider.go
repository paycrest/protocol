package provider

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/token"
	"github.com/paycrest/protocol/services"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	u "github.com/paycrest/protocol/utils"
	cryptoUtils "github.com/paycrest/protocol/utils/crypto"
	"github.com/paycrest/protocol/utils/logger"
	tokenUtils "github.com/paycrest/protocol/utils/token"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// ProviderController is a controller type for provider endpoints
type ProviderController struct {
	orderService services.Order
}

// NewProviderController creates a new instance of ProviderController with injected services
func NewProviderController() *ProviderController {
	return &ProviderController{
		orderService: services.NewOrderService(),
	}
}

// GetLockPaymentOrders controller fetches all assigned orders
func (ctrl *ProviderController) GetLockPaymentOrders(ctx *gin.Context) {
	// get page and pageSize query params
	page, offset, pageSize := u.Paginate(ctx)

	// Set ordering
	ordering := ctx.Query("ordering")
	order := ent.Desc(lockpaymentorder.FieldCreatedAt)
	if ordering == "asc" {
		order = ent.Asc(lockpaymentorder.FieldCreatedAt)
	}

	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	lockPaymentOrderQuery := storage.Client.LockPaymentOrder.Query()

	// Filter by status
	statusMap := map[string]lockpaymentorder.Status{
		"pending":    lockpaymentorder.StatusPending,
		"validated":  lockpaymentorder.StatusValidated,
		"fulfilled":  lockpaymentorder.StatusFulfilled,
		"cancelled":  lockpaymentorder.StatusCancelled,
		"processing": lockpaymentorder.StatusProcessing,
		"settled":    lockpaymentorder.StatusSettled,
	}

	statusQueryParam := ctx.Query("status")

	if status, ok := statusMap[statusQueryParam]; ok {
		lockPaymentOrderQuery = lockPaymentOrderQuery.Where(
			lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)),
			lockpaymentorder.StatusEQ(status),
		)
	} else {
		lockPaymentOrderQuery = lockPaymentOrderQuery.Where(
			lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)),
		)
	}

	count, err := lockPaymentOrderQuery.Count(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch orders", nil)
		return
	}

	// Fetch all orders assigned to the provider
	lockPaymentOrders, err := lockPaymentOrderQuery.
		Limit(pageSize).
		Offset(offset).
		Order(order).
		WithProvider().
		WithToken(
			func(query *ent.TokenQuery) {
				query.WithNetwork()
			},
		).
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch orders", nil)
		return
	}

	var orders []types.LockPaymentOrderResponse
	for _, order := range lockPaymentOrders {
		orders = append(orders, types.LockPaymentOrderResponse{
			ID:                order.ID,
			Token:             order.Edges.Token.Symbol,
			GatewayID:         order.GatewayID,
			Amount:            order.Amount,
			Rate:              order.Rate,
			Institution:       order.Institution,
			AccountIdentifier: order.AccountIdentifier,
			AccountName:       order.AccountName,
			TxHash:            order.TxHash,
			Status:            order.Status.String(),
			Memo:              order.Memo,
			Network:           order.Edges.Token.Edges.Network.Identifier,
			UpdatedAt:         order.UpdatedAt,
			CreatedAt:         order.CreatedAt,
		})
	}

	// return paginated orders
	u.APIResponse(ctx, http.StatusOK, "success", "Orders successfully retrieved", types.ProviderLockOrderList{
		Page:         page,
		PageSize:     pageSize,
		TotalRecords: count,
		Orders:       orders,
	})
}

// AcceptOrder controller accepts an order
func (ctrl *ProviderController) AcceptOrder(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the Order ID string into a UUID
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return
	}

	// Get Order request from Redis
	result, err := storage.RedisClient.HGetAll(ctx, fmt.Sprintf("order_request_%s", orderID)).Result()
	if err != nil {
		logger.Errorf("error getting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to accept order request", nil)
		return
	}

	if result["providerId"] != provider.ID || len(result) == 0 {
		logger.Errorf("order request not found in Redis: %v", orderID)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order request not found or is expired", nil)
		return
	}

	// Delete order request from Redis
	_, err = storage.RedisClient.Del(ctx, fmt.Sprintf("order_request_%s", orderID)).Result()
	if err != nil {
		logger.Errorf("error deleting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to accept order request", nil)
		return
	}

	// Update lock order status to processing
	order, err := storage.Client.LockPaymentOrder.
		UpdateOneID(orderID).
		SetStatus(lockpaymentorder.StatusProcessing).
		SetProviderID(provider.ID).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		if ent.IsNotFound(err) {
			u.APIResponse(ctx, http.StatusNotFound, "error", "Order not found", nil)
		} else {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		}
		return
	}

	u.APIResponse(ctx, http.StatusCreated, "success", "Order request accepted successfully", &types.AcceptOrderResponse{
		ID:                orderID,
		Amount:            order.Amount.Mul(order.Rate).RoundBank(0),
		Institution:       order.Institution,
		AccountIdentifier: order.AccountIdentifier,
		AccountName:       order.AccountName,
		Memo:              order.Memo,
	})
}

// DeclineOrder controller declines an order
func (ctrl *ProviderController) DeclineOrder(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the Order ID string into a UUID
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return
	}

	// Get Order request from Redis
	result, err := storage.RedisClient.HGetAll(ctx, fmt.Sprintf("order_request_%s", orderID)).Result()
	if err != nil {
		logger.Errorf("error getting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	if result["providerId"] != provider.ID || len(result) == 0 {
		logger.Errorf("order request not found in Redis: %v", orderID)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order request not found or is expired", nil)
		return
	}

	// Delete order request from Redis
	_, err = storage.RedisClient.Del(ctx, fmt.Sprintf("order_request_%s", orderID)).Result()
	if err != nil {
		logger.Errorf("error deleting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	// Push provider ID to order exclude list
	orderKey := fmt.Sprintf("order_exclude_list_%s", orderID)
	_, err = storage.RedisClient.RPush(ctx, orderKey, provider.ID).Result()
	if err != nil {
		logger.Errorf("error pushing provider %s to order %s exclude_list on Redis: %v", provider.ID, orderID, err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Order request declined successfully", nil)
}

// FulfillOrder controller fulfills an order
func (ctrl *ProviderController) FulfillOrder(ctx *gin.Context) {
	var payload types.FulfillLockOrderPayload

	// Parse the order payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get provider profile from the context
	_, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}

	// Parse the Order ID string into a UUID
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return
	}

	updateLockOrder := storage.Client.LockPaymentOrder.UpdateOneID(orderID)

	// Query or create lock order fulfillment
	fulfillment, err := storage.Client.LockOrderFulfillment.
		Query().
		Where(lockorderfulfillment.TxIDEQ(payload.TxID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			fulfillment, err = storage.Client.LockOrderFulfillment.
				Create().
				SetOrderID(orderID).
				SetTxID(payload.TxID).
				Save(ctx)
			if err != nil {
				logger.Errorf("error: %v", err)
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
				return
			}
		} else {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}
	}

	if payload.ValidationStatus == lockorderfulfillment.ValidationStatusSuccess {
		_, err := fulfillment.Update().
			SetValidationStatus(lockorderfulfillment.ValidationStatusSuccess).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

		_, err = updateLockOrder.
			SetStatus(lockpaymentorder.StatusValidated).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

		// Settle order or fail silently
		err = ctrl.orderService.SettleOrder(ctx, orderID)
		if err != nil {
			logger.Errorf("FulfillOrder.SettleOrder: %v", err)
		}

	} else if payload.ValidationStatus == lockorderfulfillment.ValidationStatusFailed {
		_, err := fulfillment.Update().
			SetValidationStatus(lockorderfulfillment.ValidationStatusFailed).
			SetValidationError(payload.ValidationError).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

		_, err = updateLockOrder.
			SetStatus(lockpaymentorder.StatusFulfilled).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}

	} else {
		_, err := updateLockOrder.
			SetStatus(lockpaymentorder.StatusFulfilled).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			return
		}
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Order fulfilled successfully", nil)
}

// CancelOrder controller cancels an order
func (ctrl *ProviderController) CancelOrder(ctx *gin.Context) {
	var payload types.CancelLockOrderPayload

	// Parse the order payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the Order ID string into a UUID
	orderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return
	}

	// Fetch lock payment order from db
	order, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.IDEQ(orderID),
			lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)),
		).
		WithProvider().
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Could not find payment order", nil)
		return
	}

	// Get new cancellation count based on cancel reason
	orderConf := config.OrderConfig()
	orderUpdate := storage.Client.LockPaymentOrder.UpdateOneID(orderID)
	CancellationCount := order.CancellationCount
	if payload.Reason == "Invalid recipient bank details" || provider.VisibilityMode == providerprofile.VisibilityModePrivate {
		CancellationCount += orderConf.RefundCancellationCount // Allows us refund immediately for invalid recipient
		orderUpdate.AppendCancellationReasons([]string{payload.Reason})
	} else if payload.Reason != "Insufficient funds" {
		CancellationCount += 1
		orderUpdate.AppendCancellationReasons([]string{payload.Reason})
	}
	// Update lock order status to cancelled
	orderUpdate.
		SetStatus(lockpaymentorder.StatusCancelled).
		SetCancellationCount(CancellationCount)
	order, err = orderUpdate.Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to cancel order", nil)
		return
	}

	// Check if order cancellation count is equal or greater than RefundCancellationCount in config,
	// and the order has not been refunded, then trigger refund
	if order.CancellationCount >= orderConf.RefundCancellationCount && order.Status == lockpaymentorder.StatusCancelled {
		err = ctrl.orderService.RefundOrder(ctx, order.GatewayID)
		if err != nil {
			logger.Errorf("CancelOrder.RefundOrder(%v): %v", orderID, err)
		}
	}

	// Push provider ID to order exclude list
	orderKey := fmt.Sprintf("order_exclude_list_%s", orderID)
	_, err = storage.RedisClient.RPush(ctx, orderKey, provider.ID).Result()
	if err != nil {
		logger.Errorf("error pushing provider %s to order %s exclude_list on Redis: %v", provider.ID, orderID, err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Order cancelled successfully", nil)
}

// GetMarketRate controller fetches the median rate of the cryptocurrency token against the fiat currency
func (ctrl *ProviderController) GetMarketRate(ctx *gin.Context) {
	// Parse path parameters
	tokenExists, err := storage.Client.Token.
		Query().
		Where(token.Symbol(strings.ToUpper(ctx.Param("token")))).
		Exist(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to get market rate", nil)
		return
	}

	if !tokenExists {
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

	orderConf := config.OrderConfig()
	deviation := currency.MarketRate.Mul(orderConf.PercentDeviationFromMarketRate.Div(decimal.NewFromInt(100)))

	u.APIResponse(ctx, http.StatusOK, "success", "Rate fetched successfully", &types.MarketRateResponse{
		MarketRate:  currency.MarketRate,
		MinimumRate: currency.MarketRate.Sub(deviation),
		MaximumRate: currency.MarketRate.Add(deviation),
	})
}

// Stats controller fetches provider stats
func (ctrl *ProviderController) Stats(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Fetch provider stats
	var v []struct {
		Count int
		Sum   decimal.Decimal
	}
	var totalFiatVolume decimal.Decimal

	query := storage.Client.LockPaymentOrder.
		Query().
		Where(lockpaymentorder.HasProviderWith(providerprofile.IDEQ(provider.ID)))

	err := query.
		Aggregate(
			ent.Count(),
			ent.Sum(lockpaymentorder.FieldAmount),
		).
		Scan(ctx, &v)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider stats", nil)
		return
	}

	orders, err := query.All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch provider stats", nil)
		return
	}

	for _, order := range orders {
		totalFiatVolume = totalFiatVolume.Add(order.Amount.Mul(order.Rate).RoundBank(0))
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Provider stats fetched successfully", &types.ProviderStatsResponse{
		TotalOrders:       v[0].Count,
		TotalFiatVolume:   totalFiatVolume,
		TotalCryptoVolume: v[0].Sum,
	})
}

// NodeInfo controller fetches the provision node info
func (ctrl *ProviderController) NodeInfo(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}

	provider, err := storage.Client.ProviderProfile.
		Query().
		Where(providerprofile.IDEQ(providerCtx.(*ent.ProviderProfile).ID)).
		WithAPIKey().
		WithCurrency().
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch node info", nil)
		return
	}

	// Compute HMAC
	decodedSecret, err := base64.StdEncoding.DecodeString(provider.Edges.APIKey.Secret)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch node info", nil)
		return
	}
	decryptedSecret, err := cryptoUtils.DecryptPlain(decodedSecret)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch node info", nil)
		return
	}
	signature := tokenUtils.GenerateHMACSignature(map[string]interface{}{}, string(decryptedSecret))

	res, err := fastshot.NewClient(provider.HostIdentifier).
		Config().SetTimeout(30*time.Second).
		Header().Add("X-Request-Signature", signature).
		Build().GET("/health").
		Retry().Set(3, 5*time.Second).
		Send()
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to fetch node info", nil)
		return
	}

	data, err := utils.ParseJSONResponse(res.RawResponse)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to fetch node info", nil)
		return
	}

	currency := data["data"].(map[string]interface{})["currency"].(string)
	if currency != provider.Edges.Currency.Code {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusServiceUnavailable, "error", "Failed to fetch node info", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Node info fetched successfully", data)
}
