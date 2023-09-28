package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/lockorderfulfillment"
	"github.com/paycrest/paycrest-protocol/ent/validatorprofile"
	"github.com/paycrest/paycrest-protocol/storage"
	"github.com/paycrest/paycrest-protocol/types"
	u "github.com/paycrest/paycrest-protocol/utils"
	"github.com/paycrest/paycrest-protocol/utils/logger"

	"github.com/gin-gonic/gin"
)

// Controller is the default controller for other endpoints
type Controller struct{}

// GetFiatCurrencies controller fetches the supported fiat currencies
func (ctrl *Controller) GetFiatCurrencies(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// GetInstitutionsByCurrency controller fetches the supported institutions for a given currency
func (ctrl *Controller) GetInstitutionsByCurrency(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// GetTokenRates controller fetches the current market rates for the supported cryptocurrencies
func (ctrl *Controller) GetTokenRates(ctx *gin.Context) {
	u.APIResponse(ctx, http.StatusOK, "success", "OK", nil)
}

// ValidateOrder is a hook to receive validation decisions from validators
func (ctrl *Controller) ValidateOrder(ctx *gin.Context) {
	var payload types.ValidateOrderPayload

	// Parse the payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get lock order fulfillment ID from URL
	fulfillmentID := ctx.Param("fulfillment_id")

	// Parse the order fulfillment ID string into a UUID
	fulfillmentUUID, err := uuid.Parse(fulfillmentID)
	if err != nil {
		logger.Errorf("error parsing fulfillment ID: %v", err)
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid fulfillment ID", nil)
		return
	}

	// Get the api key ID from the context
	apiKey, _ := ctx.Get("api_key")

	// Fetch validator from db
	validator, err := storage.Client.ValidatorProfile.
		Query().
		Where(
			validatorprofile.HasAPIKeyWith(
				apikey.ID(apiKey.(*ent.APIKey).ID),
			),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Could not find validator profile", nil)
		return
	}

	// Fetch order fulfillment from db
	fulfillment, err := storage.Client.LockOrderFulfillment.
		Query().
		Where(
			lockorderfulfillment.IDEQ(fulfillmentUUID),
		).
		Only(ctx)
	if err != nil {
		logger.Errorf("error: %v", err)
		u.APIResponse(ctx, http.StatusNotFound, "error", "Could not find order fulfillment", nil)
		return
	}

	// Update lock order fulfillment status
	if payload.IsValid {
		_, err = fulfillment.Update().
			SetConfirmations(fulfillment.Confirmations + 1).
			AddValidatorIDs(validator.ID).
			Save(ctx)
		if err != nil {
			logger.Errorf("error: %v", err)
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to validate order fulfillment", nil)
			return
		}
	} else {
		fulfillment.Update().
			AppendValidationErrors([]string{payload.ErrorMsg}).
			Save(ctx)
	}

	ctx.JSON(http.StatusOK, "OK")
}
