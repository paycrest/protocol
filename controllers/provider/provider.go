package provider

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/lockpaymentorder"
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
	// Get lock order ID from URL
	orderID := ctx.Param("id")

	// Parse the Order ID string into a UUID
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		logger.Errorf("error parsing API key ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid API key ID", nil)
		return
	}

	// Get the user ID from the context
	userID, _ := ctx.Get("user_id")

	// Get Order request from Redis
	result, err := storage.RedisClient.HGetAll(ctx, fmt.Sprintf("order_request_%d", orderUUID)).Result()
	if err != nil {
		logger.Errorf("error getting order request from Redis: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Internal server error", nil)
		return
	}

	if len(result) == 0 {
		logger.Errorf("order request not found in Redis: %d", orderUUID)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Order request not found or is expired", nil)
		return
	}

	// Update lock order status to processing
	order, err := storage.Client.LockPaymentOrder.
		UpdateOneID(orderUUID).
		SetStatus(lockpaymentorder.StatusProcessing).
		SetProviderID(userID.(string)).
		Save(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update lock order status", nil)
		return
	}

	u.APIResponse(ctx, http.StatusCreated, "success", "OK", &types.LockOrderResponse{
		ID:                orderUUID,
		Amount:            order.Amount.Mul(order.Rate),
		Token:             order.Edges.Token.Symbol,
		Institution:       order.Institution,
		AccountIdentifier: order.AccountIdentifier,
		AccountName:       order.AccountName,
		Status:            lockpaymentorder.StatusProcessing,
	})
}

// DeclineOrder controller declines an order
func (ctrl *ProviderController) DeclineOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// FulfillOrder controller fulfills an order
func (ctrl *ProviderController) FulfillOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// CancelOrder controller cancels an order
func (ctrl *ProviderController) CancelOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}
