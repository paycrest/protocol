package controllers

import (
	"net/http"

	u "github.com/paycrest/paycrest-protocol/utils"

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
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
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
