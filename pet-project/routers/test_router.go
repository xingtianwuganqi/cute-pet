package routers

import (
	"github.com/gin-gonic/gin"
	handler "pet-project/handler/api_handler"
)

func RegisterTestRouter(router *gin.Engine) {
	testRouter := router.Group("v1/test")
	{
		testRouter.GET("/get/test", handler.QueryTestNetworking)
		testRouter.POST("/post/test", handler.FormTestNetworking)
		testRouter.POST("/path/:name", handler.PathTestNetworking)
		testRouter.POST("/binding/test", handler.BindingNetworking)
	}
}
