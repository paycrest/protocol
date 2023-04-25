package routers

import (
	"github.com/paycrest/paycrest-protocol/sender/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	var ctrl controllers.UserController
	v1 := route.Group("/v1/")
	v1.GET("users/", ctrl.GetUsers)
	v1.POST("users/", ctrl.CreateUser)
}
