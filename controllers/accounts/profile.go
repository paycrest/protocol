package accounts

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/paycrest/protocol/ent"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/network"
	"github.com/paycrest/protocol/ent/provideravailability"
	"github.com/paycrest/protocol/ent/providerordertoken"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/token"
	svc "github.com/paycrest/protocol/services"
	"github.com/paycrest/protocol/storage"
	db "github.com/paycrest/protocol/storage"
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

	if payload.WebhookURL != "" {
		if !isURL(payload.WebhookURL) {
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Invalid webhook url", nil)
			return
		}
	}

	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	update := sender.Update()

	if payload.WebhookURL != "" {
		update.SetWebhookURL(payload.WebhookURL)
	}

	if payload.DomainWhitelist != nil {
		update.SetDomainWhitelist(payload.DomainWhitelist)
	}

	if !payload.FeePerTokenUnit.IsZero() {
		update.SetFeePerTokenUnit(payload.FeePerTokenUnit)
	}

	if payload.FeeAddress != "" {
		update.SetFeeAddress(payload.FeeAddress)
	}

	if payload.RefundAddress != "" {
		update.SetRefundAddress(payload.RefundAddress)
	}

	senderProfile, err := update.Save(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update profile", nil)
		return
	}

	// check if refundAddress, feeAddress and webhookURL are not empty
	if senderProfile.RefundAddress != "" && senderProfile.FeeAddress != "" && senderProfile.WebhookURL != "" {
		db.Client.SenderProfile.
			UpdateOne(senderProfile).
			SetIsActive(true).
			Save(ctx)

	} else {
		db.Client.SenderProfile.
			UpdateOne(senderProfile).
			SetIsActive(false).
			Save(ctx)
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
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	update := provider.Update()

	if payload.TradingName != "" {
		update.SetTradingName(payload.TradingName)
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
			u.APIResponse(ctx, http.StatusBadRequest, "error", "Currency not supported", nil)
			return
		}
		update.SetCurrency(currency)
	}

	if payload.HostIdentifier != "" {
		update.SetHostIdentifier(payload.HostIdentifier)
	}

	if payload.IsPartner {
		update.SetIsPartner(payload.IsPartner)
	}

	// Update availability
	if payload.Availability.Cadence != "" {
		// Get existing availability if it exists
		availability, err := storage.Client.ProviderAvailability.
			Query().
			Where(provideravailability.HasProviderWith(providerprofile.IDEQ(provider.ID))).
			Only(ctx)

		if err == nil {
			// Availability found, update it
			_, err = availability.Update().
				SetCadence(payload.Availability.Cadence).
				SetStartTime(payload.Availability.StartTime).
				SetEndTime(payload.Availability.EndTime).
				Save(ctx)
			if err != nil {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set availability", nil)
				return
			}

		} else if ent.IsNotFound(err) {
			// No existing availability, create new
			_, err = storage.Client.ProviderAvailability.
				Create().
				SetCadence(payload.Availability.Cadence).
				SetStartTime(payload.Availability.StartTime).
				SetEndTime(payload.Availability.EndTime).
				SetProviderID(provider.ID).
				Save(ctx)
			if err != nil {
				u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set availability", nil)
				return
			}

		} else {
			u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to set availability", nil)
			return
		}
	}

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

	providerProfile, err := update.Save(ctx)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to update profile", nil)
		return
	}

	// check if host id and trading name are not empty
	if providerProfile.HostIdentifier != "" && providerProfile.TradingName != "" {
		db.Client.ProviderProfile.
			UpdateOne(providerProfile).
			SetIsActive(true).
			Save(ctx)
	} else {
		db.Client.ProviderProfile.
			UpdateOne(providerProfile).
			SetIsActive(false).
			Save(ctx)
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Profile updated successfully", nil)
}

// GetSenderProfile retrieves the sender profile
func (ctrl *ProfileController) GetSenderProfile(ctx *gin.Context) {
	// Get sender profile from the context
	senderCtx, ok := ctx.Get("sender")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	sender := senderCtx.(*ent.SenderProfile)

	// Get API key
	apiKey, err := ctrl.apiKeyService.GetAPIKey(ctx, sender, nil)
	if err != nil {
		u.APIResponse(ctx, http.StatusInternalServerError, "error", "Failed to retrieve profile", nil)
		return
	}

	u.APIResponse(ctx, http.StatusOK, "success", "Profile retrieved successfully", &types.SenderProfileResponse{
		ID:              sender.ID,
		WebhookURL:      sender.WebhookURL,
		DomainWhitelist: sender.DomainWhitelist,
		FeePerTokenUnit: sender.FeePerTokenUnit,
		FeeAddress:      sender.FeeAddress,
		RefundAddress:   sender.RefundAddress,
		APIKey:          *apiKey,
	})
}

// GetProviderProfile retrieves the provider profile
func (ctrl *ProfileController) GetProviderProfile(ctx *gin.Context) {
	// Get provider profile from the context
	providerCtx, ok := ctx.Get("provider")
	if !ok {
		u.APIResponse(ctx, http.StatusUnauthorized, "error", "Invalid API key", nil)
		return
	}
	provider := providerCtx.(*ent.ProviderProfile)

	// Get availability
	availability, err := provider.QueryAvailability().Only(ctx)
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
		ID:             provider.ID,
		TradingName:    provider.TradingName,
		Currency:       provider.Edges.Currency.Code,
		HostIdentifier: provider.HostIdentifier,
		IsPartner:      provider.IsPartner,
		Availability:   availability,
		Tokens:         tokens,
		APIKey:         *apiKey,
	})
}

func isURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return true
}
