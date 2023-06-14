package routers

import (
	"github.com/paycrest/paycrest-protocol/config"
	"github.com/paycrest/paycrest-protocol/routers/middleware"

	"github.com/gin-gonic/gin"
)

// Routes function registers all routes
func Routes() *gin.Engine {
	conf := config.ServerConfig()

	environment := conf.Debug
	if environment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.SetTrustedProxies([]string{conf.AllowedHosts})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	RegisterRoutes(router) //routes register

	return router
}
