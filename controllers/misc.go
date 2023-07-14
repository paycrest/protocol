package controllers

import (
	"net/http"

	u "github.com/paycrest/paycrest-protocol/utils"

	"github.com/gin-gonic/gin"
)

// MiscController is a controller type for misc endpoints
type MiscController struct{}

// GetFiatCurrencies controller fetches the supported fiat currencies
func (ctrl *MiscController) GetFiatCurrencies(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// GetInstitutionsByCurrency controller fetches the supported institutions for a given currency
func (ctrl *MiscController) GetInstitutionsByCurrency(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// GetTokenRates controller fetches the current market rates for the supported cryptocurrencies
func (ctrl *MiscController) GetTokenRates(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}
