package routers

import (
	"github.com/paycrest/paycrest-services/sender/controllers"

	"github.com/gin-gonic/gin"
)

func ExamplesRoutes(route *gin.Engine) {
	var ctrl controllers.ExampleController
	v1 := route.Group("/v1/")
	v1.GET("test/", ctrl.GetExampleData)
	v1.POST("test/", ctrl.CreateExample)
}
