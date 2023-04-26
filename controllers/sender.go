package controllers

import (
	"net/http"

	u "github.com/paycrest/paycrest-protocol/utils"

	"github.com/gin-gonic/gin"
)

type SenderController struct{}

func (ctrl *SenderController) CreateOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, 200, "OK", nil)
}

func (ctrl *SenderController) GetOrderByID(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, 200, "OK", nil)
}

func (ctrl *SenderController) DeleteOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, 200, "OK", nil)
}
