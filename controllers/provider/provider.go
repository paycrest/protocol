package provider

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/services"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	u "github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
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
	page, pageSize := u.Paginate(ctx)

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

	// Fetch all orders assigned to the provider
	lockPaymentOrders, err := lockPaymentOrderQuery.
		Limit(pageSize).
		Offset(page).
		Order(order).
		WithProvider().
		All(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch orders", nil)
		return
	}

	var orders []types.LockPaymentOrderFields

	for _, order := range lockPaymentOrders {
		orders = append(orders, types.LockPaymentOrderFields{
			ID:                order.ID,
			Token:             order.Edges.Token,
			OrderID:           order.OrderID,
			Amount:            order.Amount.Mul(order.Rate),
			Rate:              order.Rate,
			Label:             order.Label,
			Institution:       order.Institution,
			AccountIdentifier: order.AccountIdentifier,
			AccountName:       order.AccountName,
			TxHash:            order.TxHash,
			UpdatedAt:         order.UpdatedAt,
			CreatedAt:         order.CreatedAt,
		})
	}
	// return paginated orders
	u.APIResponse(ctx, http.StatusOK, "success", "Orders successfully retrieved", types.ProviderLockOrderList{
		Page:         page + 1,
		PageSize:     pageSize,
		TotalRecords: len(orders),
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
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Internal server error", nil)
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
		Amount:            order.Amount.Mul(order.Rate),
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
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Internal server error", nil)
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

	tx, err := storage.Client.Tx(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		return
	}

	updateLockOrder := tx.LockPaymentOrder.UpdateOneID(orderID)

	// Query or create lock order fulfillment
	fulfillment, err := tx.LockOrderFulfillment.
		Query().
		Where(lockorderfulfillment.TxIDEQ(payload.TxID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			fulfillment, err = tx.LockOrderFulfillment.
				Create().
				SetOrderID(orderID).
				SetTxID(payload.TxID).
				Save(ctx)
			if err != nil {
				logger.Errorf("error: %v", err)
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
				_ = tx.Rollback()
				return
			}
		} else {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			_ = tx.Rollback()
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
			_ = tx.Rollback()
			return
		}

		_, err = updateLockOrder.
			SetStatus(lockpaymentorder.StatusValidated).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			_ = tx.Rollback()
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
			_ = tx.Rollback()
			return
		}

		_, err = updateLockOrder.
			SetStatus(lockpaymentorder.StatusFulfilled).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			_ = tx.Rollback()
			return
		}

	} else {
		_, err := updateLockOrder.
			SetStatus(lockpaymentorder.StatusFulfilled).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
			_ = tx.Rollback()
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		return
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
	if payload.Reason == "Invalid recipient bank details" {
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
		err = ctrl.orderService.RefundOrder(ctx, order.OrderID)
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

	// Get rate of the topmost provider in the priority queue of the default bucket
	keys, _, err := storage.RedisClient.Scan(ctx, uint64(0), "bucket_"+fiatSymbol+"_default", 1).Result()
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
		return
	}
	providerData, err := storage.RedisClient.LIndex(ctx, keys[0], 0).Result()
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to fetch rates", nil)
		return
	}
	marketRate, _ := decimal.NewFromString(strings.Split(providerData, ":")[1])

	deviation := config.OrderConfig().PercentDeviationFromMarketRate // 10%

	u.APIResponse(ctx, http.StatusOK, "success", "Rate fetched successfully", &types.MarketRateResponse{
		MarketRate:  marketRate,
		MinimumRate: marketRate.Mul(decimal.NewFromFloat(1).Sub(deviation)), // market rate - 10%
		MaximumRate: marketRate.Mul(decimal.NewFromFloat(1).Add(deviation)), // market rate + 10%
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
		totalFiatVolume = totalFiatVolume.Add(order.Amount.Mul(order.Rate))
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Provider stats fetched successfully", &types.ProviderStatsResponse{
		TotalOrders:       v[0].Count,
		TotalFiatVolume:   totalFiatVolume,
		TotalCryptoVolume: v[0].Sum,
	})
}
