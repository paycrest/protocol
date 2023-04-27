package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func APIResponse(ctx *gin.Context, httpCode int, status string, message string, data interface{}) {
	ctx.JSON(httpCode, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
