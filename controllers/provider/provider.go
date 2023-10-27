package provider

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paycrest/protocol/config"
	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	"github.com/paycrest/protocol/utils"
	u "github.com/paycrest/protocol/utils"
	"github.com/paycrest/protocol/utils/logger"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// ProviderController is a controller type for provider endpoints
type ProviderController struct{}

// GetOrders controller fetches all assigned orders
func (ctrl *ProviderController) GetOrders(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// AcceptOrder controller accepts an order
func (ctrl *ProviderController) AcceptOrder(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the order payload
	orderID, err := parseOrderPayload(ctx, provider)
	if err != nil {
		return
	}

	// Delete order request from Redis
	_, err = storage.RedisClient.Del(ctx, fmt.Sprintf("order_request_%d", orderID)).Result()
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
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		return
	}

	u.APIResponse(ctx, http.StatusCreated, "success", "Order request accepted successfully", &types.LockOrderResponse{
		ID:                orderID,
		Amount:            order.Amount.Mul(order.Rate),
		Token:             order.Edges.Token.Symbol,
		Institution:       order.Institution,
		AccountIdentifier: order.AccountIdentifier,
		AccountName:       order.AccountName,
		Status:            lockpaymentorder.StatusProcessing,
		UpdatedAt:         order.UpdatedAt,
	})
}

// DeclineOrder controller declines an order
func (ctrl *ProviderController) DeclineOrder(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the order payload
	orderID, err := parseOrderPayload(ctx, provider)
	if err != nil {
		return
	}

	// Delete order request from Redis
	_, err = storage.RedisClient.Del(ctx, fmt.Sprintf("order_request_%d", orderID)).Result()
	if err != nil {
		logger.Errorf("error deleting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	// Push provider ID to order exclude list
	orderKey := fmt.Sprintf("order_exclude_list_%d", orderID)
	_, err = storage.RedisClient.RPush(ctx, orderKey, provider.ID).Result()
	if err != nil {
		logger.Errorf("error pushing provider %s to order %d exclude_list on Redis: %v", provider.ID, orderID, err)
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
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the order payload
	orderID, err := parseOrderPayload(ctx, provider)
	if err != nil {
		return
	}

	tx, err := storage.Client.Tx(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		return
	}

	// Save lock order fulfillment
	_, err = tx.LockOrderFulfillment.
		Create().
		SetOrderID(orderID).
		SetTxID(payload.TxID).
		SetTxReceiptImage(payload.TxReceiptImage).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		_ = tx.Rollback()
		return
	}

	// Update lock order status to fulfilled
	_, err = tx.LockPaymentOrder.
		UpdateOneID(orderID).
		SetStatus(lockpaymentorder.StatusFulfilled).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		_ = tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		return
	}

	// Delete the order exclude list
	orderKey := fmt.Sprintf("order_exclude_list_%d", orderID)
	_, err = storage.RedisClient.Del(ctx, orderKey).Result()
	if err != nil {
		logger.Errorf("error deleting order exclude list from Redis: %v", err)
	}

	// Publish fulfillment to validators
	// TODO: figure out if we need to wait for XX minutes before posting to avoid "pending" responses
	conf := config.OrderConfig()
	message := types.FulfillmentMessage{
		TxID:                    payload.TxID,
		TxReceiptImage:          payload.TxReceiptImage,
		Institution:             payload.Institution,
		MaxConcurrentValidators: conf.MaxConcurrentValidators,
	}
	err = storage.RedisClient.Publish(ctx, "order_fulfillments", message).Err()
	if err != nil {
		logger.Errorf("error publishing to validators: %v", err)
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
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Parse the order payload
	orderID, err := parseOrderPayload(ctx, provider)
	if err != nil {
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

	// Update lock order status to cancelled
	_, err = storage.Client.LockPaymentOrder.
		UpdateOneID(orderID).
		SetStatus(lockpaymentorder.StatusCancelled).
		SetCancellationCount(order.CancellationCount + 1).
		AppendCancellationReasons([]string{payload.Reason}).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to cancel order", nil)
		return
	}

	// Push provider ID to order exclude list
	orderKey := fmt.Sprintf("order_exclude_list_%d", orderID)
	_, err = storage.RedisClient.RPush(ctx, orderKey, order.Edges.Provider.ID).Result()
	if err != nil {
		logger.Errorf("error pushing provider %s to order %d exclude_list on Redis: %v", order.Edges.Provider.ID, orderID, err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to decline order request", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Order cancelled successfully", nil)
}

// GetMarketRate controller fetches the median rate of the cryptocurrency token against the fiat currency
func (ctrl *ProviderController) GetMarketRate(ctx *gin.Context) {
	// Parse path parameters
	token := ctx.Param("token")
	tokenIsValid := utils.ContainsString([]string{"USDT", "USDC"}, token) // TODO: fetch supported tokens from db
	if !tokenIsValid {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Token is not supported", nil)
	}

	fiatSymbol := ctx.Param("fiat")
	fiatIsValid := utils.ContainsString([]string{"NGN"}, fiatSymbol) // TODO: fetch supported fiat currencies from db
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

	deviation := decimal.NewFromFloat(0.01) // 1%

	u.APIResponse(ctx, http.StatusOK, "success", "Rate fetched successfully", &types.MarketRateResponse{
		MarketRate:  marketRate,
		MinimumRate: marketRate.Mul(decimal.NewFromFloat(1).Sub(deviation)), // market rate - 1%
		MaximumRate: marketRate.Mul(decimal.NewFromFloat(1).Add(deviation)), // market rate + 1%
	})
}

// parseOrderPayload parses the order payload
func parseOrderPayload(ctx *gin.Context, provider *ent.ProviderProfile) (uuid.UUID, error) {
	// Get lock order ID from URL
	orderID := ctx.Param("id")

	// Parse the Order ID string into a UUID
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		logger.Errorf("error parsing order ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid Order ID", nil)
		return uuid.UUID{}, err
	}

	// Get Order request from Redis
	result, err := storage.RedisClient.HGetAll(ctx, fmt.Sprintf("order_request_%d", orderUUID)).Result()
	if err != nil {
		logger.Errorf("error getting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Internal server error", nil)
		return uuid.UUID{}, err
	}

	if result["provider_id"] != provider.ID || len(result) == 0 {
		logger.Errorf("order request not found in Redis: %d", orderUUID)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order request not found or is expired", nil)
		return uuid.UUID{}, err
	}

	return orderUUID, nil
}
