package controllers

import (
	"net/http"

	u "github.com/paycrest/paycrest-protocol/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (ctrl *AuthController) Register(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, 200, "OK", nil)
}

func (ctrl *AuthController) Login(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, 200, "OK", nil)
}

func (ctrl *AuthController) RefreshJWT(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, 200, "OK", nil)
}
