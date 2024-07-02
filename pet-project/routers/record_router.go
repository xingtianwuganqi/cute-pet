package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRecordRouter(r *gin.Engine) {
	recordRouter := r.Group("/v1/record")
	{
		recordRouter.GET("/pet/action/list", handler.GetPetActionList)
		recordRouter.POST("/pet/create", middleware.JWTTokenMiddleware(), handler.PetInfoCreate)
		recordRouter.POST("/pet/create/action", handler.CreatePetActionType)
		recordRouter.GET("/list", middleware.JWTTokenMiddleware(), handler.GetRecordList)
	}

	actionRouter := r.Group("v1/action")
	{
		// 获取日常列表
		actionRouter.GET("/list", handler.GetPetActionList)
		// 用户获取自己的列表
		actionRouter.GET("/custom/list", middleware.JWTTokenMiddleware(), handler.GetCustomActionList)
		// 用户自己创建
		actionRouter.POST("/create/custom", middleware.JWTTokenMiddleware(), handler.CreatePetCustomType)

	}

	consumerRouter := r.Group("v1/consumer")
	{
		// 用户获取花销列表
		consumerRouter.GET("/list", handler.GetPetConsumeList)
		// 用户获取自己创建的花销列表
		consumerRouter.POST("/custom/list", middleware.JWTTokenMiddleware(), handler.GetPetCustomConsumeList)
		// 创建或者更新花销
		consumerRouter.POST("/create", middleware.JWTTokenMiddleware(), handler.CreateConsumeAction)
	}
}
