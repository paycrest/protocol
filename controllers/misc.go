package controllers

import (
	"net/http"

	u "github.com/paycrest/paycrest-protocol/utils"

	"github.com/gin-gonic/gin"
)

type MiscController struct{}

func (ctrl *MiscController) GetFiatCurrencies(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

func (ctrl *MiscController) GetInstitutionsByCurrency(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

func (ctrl *MiscController) GetCryptoRates(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}
