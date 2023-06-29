package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paycrest/paycrest-protocol/controllers"
	"github.com/paycrest/paycrest-protocol/controllers/accounts"
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
	miscRoutes(route)
}

func authRoutes(route *gin.Engine) {
	var ctrl accounts.AuthController

	v1 := route.Group("/v1/auth/")
	v1.POST("register/", ctrl.Register)
	v1.POST("login/", ctrl.Login)
	v1.POST("refresh/", middleware.JWTMiddleware, ctrl.RefreshJWT)
	v1.POST("api-keys/", middleware.JWTMiddleware, ctrl.GenerateAPIKey)
	v1.GET("api-keys/", middleware.JWTMiddleware, ctrl.ListAPIKeys)
	v1.DELETE("api-keys/:id", middleware.JWTMiddleware, ctrl.DeleteAPIKey)
}

func senderRoutes(route *gin.Engine) {
	var ctrl sender.Controller

	v1 := route.Group("/v1/sender/")
	// v1.Use(middleware.HMACVerificationMiddleware)

	v1.POST("orders/", ctrl.CreateOrder)
	v1.GET("orders/:id", ctrl.GetOrderByID)
	v1.DELETE("orders/:id", ctrl.DeleteOrder)
}

func providerRoutes(route *gin.Engine) {
	var ctrl controllers.ProviderController

	v1 := route.Group("/v1/provider/")
	v1.GET("orders/", ctrl.GetOrders)
	v1.POST("orders/:id/accept", ctrl.AcceptOrder)
	v1.POST("orders/:id/decline", ctrl.DeclineOrder)
	v1.POST("orders/:id/fulfill", ctrl.FulfillOrder)
	v1.POST("orders/:id/cancel", ctrl.CancelOrder)
}

func miscRoutes(route *gin.Engine) {
	var ctrl controllers.ProviderController

	v1 := route.Group("/v1/misc/")
	v1.GET("currencies/", ctrl.GetOrders)
	v1.GET("institutions/:currencyCode", ctrl.GetOrders)
	v1.GET("rates/:crypto", ctrl.GetOrders)
}
