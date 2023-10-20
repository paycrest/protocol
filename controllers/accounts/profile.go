package accounts

import (
	"net/http"

	"github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/types"
	u "github.com/paycrest/paycrest-protocol/utils"

	"github.com/gin-gonic/gin"
)

// ProfileController is a controller type for profile settings
type ProfileController struct{}

// UpdateValidatorProfile controller updates the validator profile
func (ctrl *ProfileController) UpdateValidatorProfile(ctx *gin.Context) {
	var payload types.ValidatorProfilePayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get validator profile from the context
	validatorCtx, ok := ctx.Get("validator")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	validator := validatorCtx.(*ent.ValidatorProfile)

	update := validator.Update()

	if payload.WalletAddress != "" {
		update.SetWalletAddress(payload.WalletAddress)
	}

	if payload.HostIdentifier != "" {
		update.SetHostIdentifier(payload.HostIdentifier)
	}

	_, err := update.Save(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update profile", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Profile updated successfully", nil)
}

// GetValidatorProfile retrieves the validator profile
func (ctrl *ProfileController) GetValidatorProfile(ctx *gin.Context) {
	// Get validator profile from the context
	validatorCtx, ok := ctx.Get("validator")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	validator := validatorCtx.(*ent.ValidatorProfile)

	u.APIResponse(ctx, http.StatusOK, "success", "Profile retrieved successfully", &types.ValidatorProfileResponse{
		ID:             validator.ID,
		WalletAddress:  validator.WalletAddress,
		HostIdentifier: validator.HostIdentifier,
	})
}
