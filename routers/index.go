package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/controllers"
	"github.com/paycrest/paycrest-protocol/controllers/accounts"
	"github.com/paycrest/paycrest-protocol/controllers/provider"
	"github.com/paycrest/paycrest-protocol/controllers/sender"
	"github.com/paycrest/paycrest-protocol/routers/middleware"
	u "github.com/paycrest/paycrest-protocol/utils"
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

	v1.POST(
		"orders/:fulfillment_id/validate",
		middleware.HMACVerificationMiddleware,
		middleware.OnlyValidatorMiddleware,
		ctrl.ValidateOrder,
	)
	v1.GET("currencies/", ctrl.GetFiatCurrencies)
	v1.GET("institutions/:currencyCode", ctrl.GetInstitutionsByCurrency)
	v1.GET("rates/:token", ctrl.GetTokenRates)
}

func authRoutes(route *gin.Engine) {
	var ctrl accounts.AuthController

	v1 := route.Group("/v1/auth/")
	v1.POST("register/", ctrl.Register)
	v1.POST("login/", ctrl.Login)
	v1.POST("confirm-account/", ctrl.ConfirmEmail)
	v1.POST("resend-token/", ctrl.ResendVerificationToken)
	v1.POST("refresh/", middleware.JWTMiddleware, ctrl.RefreshJWT)
	v1.POST("api-keys/", middleware.JWTMiddleware, ctrl.CreateAPIKey)
	v1.GET("api-keys/", middleware.JWTMiddleware, ctrl.ListAPIKeys)
	v1.DELETE("api-keys/:id", middleware.JWTMiddleware, ctrl.DeleteAPIKey)
}

func senderRoutes(route *gin.Engine) {
	var ctrl sender.SenderController

	v1 := route.Group("/v1/sender/")
	v1.Use(middleware.HMACVerificationMiddleware)
	v1.Use(middleware.OnlySenderMiddleware)

	v1.POST("orders/", ctrl.CreatePaymentOrder)
	v1.GET("orders/:id", ctrl.GetPaymentOrderByID)
}

func providerRoutes(route *gin.Engine) {
	var ctrl provider.ProviderController

	v1 := route.Group("/v1/provider/")
	v1.Use(middleware.HMACVerificationMiddleware)
	v1.Use(middleware.OnlyProviderMiddleware)

	v1.GET("orders/", ctrl.GetOrders)
	v1.POST("orders/:id/accept", ctrl.AcceptOrder)
	v1.POST("orders/:id/decline", ctrl.DeclineOrder)
	v1.POST("orders/:id/fulfill", ctrl.FulfillOrder)
	v1.POST("orders/:id/cancel", ctrl.CancelOrder)
}
