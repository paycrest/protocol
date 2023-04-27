package controllers

import (
	"net/http"

	u "github.com/paycrest/paycrest-protocol/utils"

	"github.com/gin-gonic/gin"
)

type ProviderController struct{}

func (ctrl *ProviderController) GetOrders(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

func (ctrl *ProviderController) AcceptOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

func (ctrl *ProviderController) DeclineOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

func (ctrl *ProviderController) FulfillOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

func (ctrl *ProviderController) CancelOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}
