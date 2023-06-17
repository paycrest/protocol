package sender

import (
	"net/http"

	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"

	"github.com/gin-gonic/gin"
)

// Controller is a controller type for sender endpoints
type Controller struct{}

// CreateOrder controller creates an order
func (ctrl *Controller) CreateOrder(ctx *gin.Context) {
	var payload interface{}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}
	u.APIResponse(ctx, http.StatusOK, "success", "OK", &payload)
}

// GetOrderByID controller fetches an order by ID
func (ctrl *Controller) GetOrderByID(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// DeleteOrder controller deletes an order
func (ctrl *Controller) DeleteOrder(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}
