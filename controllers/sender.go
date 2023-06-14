package controllers

import (
	"net/http"

	u "github.com/paycrest/paycrest-protocol/utils"

	"github.com/gin-gonic/gin"
)

// SenderController is a controller type for sender endpoints
type SenderController struct{}

// CreateOrder controller creates an order
func (ctrl *SenderController) CreateOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// GetOrderByID controller fetches an order by ID
func (ctrl *SenderController) GetOrderByID(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// DeleteOrder controller deletes an order
func (ctrl *SenderController) DeleteOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}
