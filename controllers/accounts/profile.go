package accounts

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/token"
	svc "github.com/paycrest/protocol/services"
	"github.com/paycrest/protocol/storage"
	"github.com/paycrest/protocol/types"
	u "github.com/paycrest/protocol/utils"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

// ProfileController is a controller type for profile settings
type ProfileController struct {
	apiKeyService *svc.APIKeyService
}

// NewProfileController creates a new instance of ProfileController
func NewProfileController() *ProfileController {
	return &ProfileController{
		apiKeyService: svc.NewAPIKeyService(),
	}
}

// UpdateSenderProfile controller updates the sender profile
func (ctrl *ProfileController) UpdateSenderProfile(ctx *gin.Context) {
	var payload types.SenderProfilePayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	if payload.WebhookURL != "" && !u.IsURL(payload.WebhookURL) {
		u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", []types.ErrorData{{
			Field:   "WebhookURL",
			Message: "Invalid URL",
		}})
		return
	}

	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	update := sender.Update()

	if payload.WebhookURL != "" || (payload.WebhookURL == "" && sender.WebhookURL != "") {
		update.SetWebhookURL(payload.WebhookURL)
	}

	if payload.DomainWhitelist != nil || (payload.DomainWhitelist == nil && sender.DomainWhitelist != nil) {
		update.SetDomainWhitelist(payload.DomainWhitelist)
	}

	if !payload.FeePerTokenUnit.IsZero() && payload.FeeAddress != "" {
		update.SetFeePerTokenUnit(payload.FeePerTokenUnit).SetFeeAddress(payload.FeeAddress)
	} else {
		if !payload.FeePerTokenUnit.IsZero() && payload.FeeAddress == "" {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", types.ErrorData{
				Field:   "FeeAddress",
				Message: "This field is required",
			})
			return
		}
		if payload.FeePerTokenUnit.IsZero() && payload.FeeAddress != "" {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", types.ErrorData{
				Field:   "FeePerTokenUnit",
				Message: "This field is required",
			})
			return
		}
	}

	if payload.RefundAddress != "" {
		update.SetRefundAddress(payload.RefundAddress).SetIsActive(true)
	} else if payload.RefundAddress == "" && sender.RefundAddress != "" {
		update.SetRefundAddress(payload.RefundAddress).SetIsActive(false)
	}

	_, err := update.Save(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update profile", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Profile updated successfully", nil)
}

// UpdateProviderProfile controller updates the provider profile
func (ctrl *ProfileController) UpdateProviderProfile(ctx *gin.Context) {
	var payload types.ProviderProfilePayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		u.APIResponse(ctx, http.StatusBadRequest, "error",
			"Failed to validate payload", u.GetErrorData(err))
		return
	}

	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	update := provider.Update()

	if payload.TradingName != "" || (payload.TradingName == "" && provider.TradingName != "") {
		update.SetTradingName(payload.TradingName)
	}

	if payload.HostIdentifier != "" || (payload.HostIdentifier == "" && provider.HostIdentifier != "") {
		update.SetHostIdentifier(payload.HostIdentifier)
	}

	if payload.Currency != "" {
		currency, err := storage.Client.FiatCurrency.
			Query().
			Where(
				fiatcurrency.IsEnabledEQ(true),
				fiatcurrency.CodeEQ(payload.Currency),
			).
			Only(ctx)
		if err != nil {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Failed to validate payload", types.ErrorData{
				Field:   "FiatCurrency",
				Message: "This field is required",
			})
			return
		}
		update.SetCurrency(currency)
	}

	if payload.IsPartner {
		update.SetIsPartner(true)
	} else {
		update.SetIsPartner(false)
	}

	if payload.VisibilityMode != "" {
		update.SetVisibilityMode(providerprofile.VisibilityMode(payload.VisibilityMode))
	}

	if payload.Address != "" {
		update.SetAddress(payload.Address)
	}

	if payload.MobileNumber != "" {
		if !u.IsValidMobileNumber(payload.MobileNumber) {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid mobile number", nil)
			return
		}
		update.SetMobileNumber(payload.MobileNumber)
	}

	if !payload.DateOfBirth.IsZero() {
		update.SetDateOfBirth(payload.DateOfBirth)
	}

	if payload.BusinessName != "" {
		update.SetBusinessName(payload.BusinessName)
	}

	if payload.IdentityDocumentType != "" {
		if providerprofile.IdentityDocumentType(payload.IdentityDocumentType) != providerprofile.IdentityDocumentTypePassport &&
			providerprofile.IdentityDocumentType(payload.IdentityDocumentType) != providerprofile.IdentityDocumentTypeDriversLicense &&
			providerprofile.IdentityDocumentType(payload.IdentityDocumentType) != providerprofile.IdentityDocumentTypeNationalID {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid identity document type", nil)
			return
		}
		update.SetIdentityDocumentType(providerprofile.IdentityDocumentType(payload.IdentityDocumentType))
	}

	if payload.IdentityDocument != "" {
		if !u.IsValidFileURL(payload.IdentityDocument) {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid identity document URL", nil)
			return
		}
		update.SetIdentityDocument(payload.IdentityDocument)
	}

	if payload.BusinessDocument != "" {
		if !u.IsValidFileURL(payload.BusinessDocument) {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid business document URL", nil)
			return
		}
		update.SetBusinessDocument(payload.BusinessDocument)
	}

	// // Update availability
	// if payload.Availability != (types.ProviderAvailabilityPayload{}) && payload.Availability.Cadence != "" {
	// 	// Get existing availability if it exists
	// 	availability, err := storage.Client.ProviderAvailability.
	// 		Query().
	// 		Where(provideravailability.HasProviderWith(providerprofile.IDEQ(provider.ID))).
	// 		Only(ctx)

	// 	if err == nil {
	// 		// Availability found, update it
	// 		_, err = availability.Update().
	// 			SetCadence(payload.Availability.Cadence).
	// 			SetStartTime(payload.Availability.StartTime).
	// 			SetEndTime(payload.Availability.EndTime).
	// 			Save(ctx)
	// 		if err != nil {
	// 			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set availability", nil)
	// 			return
	// 		}

	// 	} else if ent.IsNotFound(err) {
	// 		// No existing availability, create new
	// 		_, err = storage.Client.ProviderAvailability.
	// 			Create().
	// 			SetCadence(payload.Availability.Cadence).
	// 			SetStartTime(payload.Availability.StartTime).
	// 			SetEndTime(payload.Availability.EndTime).
	// 			SetProviderID(provider.ID).
	// 			Save(ctx)
	// 		if err != nil {
	// 			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set availability", nil)
	// 			return
	// 		}

	// 	} else {
	// 		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set availability", nil)
	// 		return
	// 	}
	// }

	// Update tokens
	for _, tokenPayload := range payload.Tokens {
		if len(tokenPayload.Addresses) == 0 {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "No addresses provided", nil)
			return
		}

		// Check if token is supported
		_, err := storage.Client.Token.
			Query().
			Where(token.Symbol(tokenPayload.Symbol)).
			First(ctx)
		if err != nil {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Token not supported", nil)
			return
		}

		// Check if network is supported
		for _, addressPayload := range tokenPayload.Addresses {
			_, err = storage.Client.Network.
				Query().
				Where(network.IdentifierEQ(addressPayload.Network)).
				First(ctx)
			if err != nil {
				u.APIResponse(
					ctx,
					http.StatusBadRequest,
					"error", "Network not supported - "+addressPayload.Network,
					nil,
				)
				return
			}
		}

		// Ensure rate is within allowed deviation from the market rate
		partnerProviderData, _ := storage.RedisClient.LIndex(ctx, fmt.Sprintf("bucket_%s_default", payload.Currency), 0).Result()
		marketRate, _ := decimal.NewFromString(strings.Split(partnerProviderData, ":")[1])

		var rate decimal.Decimal

		if tokenPayload.ConversionRateType == providerordertoken.ConversionRateTypeFixed {
			rate = tokenPayload.FixedConversionRate
		} else {
			floatingRate := tokenPayload.FloatingConversionRate // in percentage
			rate = marketRate.Mul(floatingRate.Div(decimal.NewFromInt(100)))
		}

		allowedDeviation := decimal.NewFromFloat(0.01) // 1%

		if marketRate.Cmp(decimal.Zero) != 0 {
			if rate.LessThan(marketRate.Mul(decimal.NewFromFloat(1).Sub(allowedDeviation))) ||
				rate.GreaterThan(marketRate.Mul(decimal.NewFromFloat(1).Add(allowedDeviation))) {
				u.APIResponse(ctx, http.StatusBadRequest, "error", "Rate is too far from market rate", nil)
				return
			}
		}

		// See if token already exists for provider
		orderToken, err := storage.Client.ProviderOrderToken.
			Query().
			Where(
				providerordertoken.SymbolEQ(tokenPayload.Symbol),
				providerordertoken.HasProviderWith(providerprofile.IDEQ(provider.ID)),
			).
			Only(ctx)

		if err != nil {
			if ent.IsNotFound(err) {
				// Token doesn't exist, create it
				_, err = storage.Client.ProviderOrderToken.
					Create().
					SetSymbol(tokenPayload.Symbol).
					SetConversionRateType(tokenPayload.ConversionRateType).
					SetFixedConversionRate(tokenPayload.FixedConversionRate).
					SetFloatingConversionRate(tokenPayload.FloatingConversionRate).
					SetMaxOrderAmount(tokenPayload.MaxOrderAmount).
					SetMinOrderAmount(tokenPayload.MinOrderAmount).
					SetAddresses(tokenPayload.Addresses).
					SetProviderID(provider.ID).
					Save(ctx)
				if err != nil {
					u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set token - "+tokenPayload.Symbol, nil)
				}
			} else {
				if err != nil {
					u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set token - "+tokenPayload.Symbol, nil)
				}
			}
		} else {
			// Token exists, update it
			_, err = orderToken.Update().
				SetSymbol(tokenPayload.Symbol).
				SetConversionRateType(tokenPayload.ConversionRateType).
				SetFixedConversionRate(tokenPayload.FixedConversionRate).
				SetFloatingConversionRate(tokenPayload.FloatingConversionRate).
				SetMaxOrderAmount(tokenPayload.MaxOrderAmount).
				SetMinOrderAmount(tokenPayload.MinOrderAmount).
				SetAddresses(tokenPayload.Addresses).
				SetProviderID(provider.ID).
				Save(ctx)
			if err != nil {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set token - "+tokenPayload.Symbol, nil)
			}
		}
	}

	// Activate profile
	if payload.HostIdentifier != "" && payload.TradingName != "" {
		update.SetIsActive(true)
	} else {
		update.SetIsActive(false)
	}

	_, err := update.Save(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update profile", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Profile updated successfully", nil)
}

// GetSenderProfile retrieves the sender profile
func (ctrl *ProfileController) GetSenderProfile(ctx *gin.Context) {
	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	user, err := sender.QueryUser().Only(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to retrieve profile", nil)
		return
	}

	// Get API key
	apiKey, err := ctrl.apiKeyService.GetAPIKey(ctx, sender, nil)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to retrieve profile", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Profile retrieved successfully", &types.SenderProfileResponse{
		ID:              sender.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		WebhookURL:      sender.WebhookURL,
		DomainWhitelist: sender.DomainWhitelist,
		FeePerTokenUnit: sender.FeePerTokenUnit,
		FeeAddress:      sender.FeeAddress,
		RefundAddress:   sender.RefundAddress,
		APIKey:          *apiKey,
		IsActive:        sender.IsActive,
	})
}

// GetProviderProfile retrieves the provider profile
func (ctrl *ProfileController) GetProviderProfile(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key or token", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	user, err := provider.QueryUser().Only(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to retrieve profile", nil)
		return
	}

	// Get tokens
	tokens, err := provider.QueryOrderTokens().All(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to retrieve profile", nil)
		return
	}

	// Get API key
	apiKey, err := ctrl.apiKeyService.GetAPIKey(ctx, nil, provider)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to retrieve profile", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Profile retrieved successfully", &types.ProviderProfileResponse{
		ID:                   provider.ID,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Email:                user.Email,
		TradingName:          provider.TradingName,
		Currency:             provider.Edges.Currency.Code,
		HostIdentifier:       provider.HostIdentifier,
		IsPartner:            provider.IsPartner,
		IsAvailable:          provider.IsAvailable,
		Tokens:               tokens,
		APIKey:               *apiKey,
		IsActive:             provider.IsActive,
		Address:              provider.Address,
		MobileNumber:         provider.MobileNumber,
		DateOfBirth:          provider.DateOfBirth,
		BusinessName:         provider.BusinessName,
		IdentityDocumentType: string(provider.IdentityDocumentType),
		IdentityDocument:     provider.IdentityDocument,
		BusinessDocument:     provider.BusinessDocument,
	})
}
