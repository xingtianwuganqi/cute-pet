package routers

import (
	handler "pet-project/handler/api"
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
		// 获取公共的列表
		actionRouter.GET("/list", handler.GetPetActionList)
		// 用户获取自己的列表
		actionRouter.GET("/custom/list", middleware.JWTTokenMiddleware(), handler.GetCustomActionList)
		// 只有管理账号可以创建用户
		actionRouter.POST("/create/action", handler.CreatePetActionType)
		// 用户自己创建
		actionRouter.POST("/create/custom", middleware.JWTTokenMiddleware(), handler.CreatePetCustomType)

	}
}
