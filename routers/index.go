package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/protocol/controllers"
	"github.com/paycrest/protocol/controllers/accounts"
	"github.com/paycrest/protocol/controllers/provider"
	"github.com/paycrest/protocol/controllers/sender"
	"github.com/paycrest/protocol/routers/middleware"
	u "github.com/paycrest/protocol/utils"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		u.APIResponse(ctx, http.StatusNotFound, "error", "Route Not Found", nil)
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

	// Add all routes
	authRoutes(route)
	senderRoutes(route)
	providerRoutes(route)

	var ctrl controllers.Controller

	v1 := route.Group("/v1/")

	v1.GET(
		"currencies/",
		middleware.DynamicAuthMiddleware,
		ctrl.GetFiatCurrencies,
	)
	v1.GET(
		"institutions/:currency_code",
		middleware.DynamicAuthMiddleware,
		ctrl.GetInstitutionsByCurrency,
	)
	v1.GET("rates/:token/:amount/:fiat_symbol", ctrl.GetTokenRate)
}

func authRoutes(route *gin.Engine) {
	authCtrl := accounts.NewAuthController()
	var profileCtrl accounts.ProfileController

	v1 := route.Group("/v1/")
	v1.POST("auth/register/", authCtrl.Register)
	v1.POST("auth/login/", authCtrl.Login)
	v1.POST("auth/confirm-account/", authCtrl.ConfirmEmail)
	v1.POST("auth/resend-token/", authCtrl.ResendVerificationToken)
	v1.POST("auth/refresh/", middleware.JWTMiddleware, authCtrl.RefreshJWT)
	v1.POST("auth/reset-password-token", middleware.JWTMiddleware, authCtrl.ResetPasswordToken)
	v1.PATCH("auth/reset-password", middleware.JWTMiddleware, authCtrl.ResetPassword)

	v1.GET(
		"settings/provider",
		middleware.JWTMiddleware,
		middleware.OnlyProviderMiddleware,
		profileCtrl.GetProviderProfile,
	)
	v1.PATCH(
		"settings/provider",
		middleware.JWTMiddleware,
		middleware.OnlyProviderMiddleware,
		profileCtrl.UpdateProviderProfile,
	)

	v1.GET(
		"settings/sender",
		middleware.JWTMiddleware,
		middleware.OnlySenderMiddleware,
		profileCtrl.GetSenderProfile,
	)
	v1.PATCH(
		"settings/sender",
		middleware.JWTMiddleware,
		middleware.OnlySenderMiddleware,
		profileCtrl.UpdateSenderProfile,
	)
}

func senderRoutes(route *gin.Engine) {
	var senderCtrl sender.SenderController

	v1 := route.Group("/v1/sender/")
	v1.Use(middleware.HMACVerificationMiddleware)
	v1.Use(middleware.OnlySenderMiddleware)

	v1.POST("orders/", senderCtrl.CreatePaymentOrder)
	v1.GET("orders/:id", senderCtrl.GetPaymentOrderByID)
}

func providerRoutes(route *gin.Engine) {
	var providerCtrl provider.ProviderController

	v1 := route.Group("/v1/provider/")
	v1.Use(middleware.DynamicAuthMiddleware)
	v1.Use(middleware.OnlyProviderMiddleware)

	v1.GET("orders/", providerCtrl.GetOrders)
	v1.POST("orders/:id/accept", providerCtrl.AcceptOrder)
	v1.POST("orders/:id/decline", providerCtrl.DeclineOrder)
	v1.POST("orders/:id/fulfill", providerCtrl.FulfillOrder)
	v1.POST("orders/:id/cancel", providerCtrl.CancelOrder)
	v1.GET("rates/:token/:fiat_symbol", providerCtrl.GetMarketRate)
}
