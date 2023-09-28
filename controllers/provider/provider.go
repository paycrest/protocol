package provider

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/lockpaymentorder"
	"github.com/paycrest/paycrest-protocol/ent/providerprofile"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"

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
	// Parse the order payload
	orderID, err := parseOrderPayload(ctx)
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

	// Fetch provider from db
	apiKey, _ := ctx.Get("api_key")
	provider, err := storage.Client.ProviderProfile.
		Query().
		Where(
			providerprofile.HasAPIKeyWith(
				apikey.IDEQ(apiKey.(*ent.APIKey).ID),
			),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Could not find provider profile", nil)
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
	// Parse the order payload
	orderID, err := parseOrderPayload(ctx)
	if err != nil {
		return
	}

	// Fetch provider from db
	apiKey, _ := ctx.Get("api_key")
	provider, err := storage.Client.ProviderProfile.
		Query().
		Where(
			providerprofile.HasAPIKeyWith(
				apikey.IDEQ(apiKey.(*ent.APIKey).ID),
			),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Could not find provider profile", nil)
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

	orderID, err := parseOrderPayload(ctx)
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

	orderID, err := parseOrderPayload(ctx)
	if err != nil {
		return
	}

	// Fetch lock payment order from db
	apiKey, _ := ctx.Get("api_key")
	order, err := storage.Client.LockPaymentOrder.
		Query().
		Where(
			lockpaymentorder.IDEQ(orderID),
			lockpaymentorder.HasProviderWith(
				providerprofile.HasAPIKeyWith(
					apikey.IDEQ(apiKey.(*ent.APIKey).ID),
				),
			),
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

// parseOrderPayload parses the order payload
func parseOrderPayload(ctx *gin.Context) (uuid.UUID, error) {
	// Get lock order ID from URL
	orderID := ctx.Param("id")

	// Parse the Order ID string into a UUID
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		logger.Errorf("error parsing API key ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid API key ID", nil)
		return uuid.UUID{}, err
	}

	// Get Order request from Redis
	result, err := storage.RedisClient.HGetAll(ctx, fmt.Sprintf("order_request_%d", orderUUID)).Result()
	if err != nil {
		logger.Errorf("error getting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Internal server error", nil)
		return uuid.UUID{}, err
	}

	if len(result) == 0 {
		logger.Errorf("order request not found in Redis: %d", orderUUID)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order request not found or is expired", nil)
		return uuid.UUID{}, err
	}

	return orderUUID, nil
}
