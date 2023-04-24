package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func APIResponse(ctx *gin.Context, httpCode, code int, message string, data interface{}) {
	ctx.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
