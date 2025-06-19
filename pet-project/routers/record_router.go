package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRecordRouter(r *gin.Engine) {
	recordRouter := r.Group("/v1/record")
	{
		recordRouter.GET("/list", middleware.JWTTokenMiddleware(), handler.GetRecordList)
	}

	petRouter := r.Group("/v1/pet")
	{
		petRouter.GET("/list", middleware.JWTTokenMiddleware(), handler.GetPetList)
		petRouter.POST("/create", middleware.JWTTokenMiddleware(), handler.PetInfoCreate)
		petRouter.PUT("/update", middleware.JWTTokenMiddleware(), handler.UpdatePetInfo)
		petRouter.DELETE("/delete/:id", middleware.JWTTokenMiddleware(), handler.DeletePetInfo)

		petRouter.POST("/create/action", handler.CreatePetActionType)
		petRouter.GET("/action/list", handler.GetPetActionList)
		petRouter.POST("/custom/action/create", middleware.JWTTokenMiddleware(), handler.CreatePetCustomAction)
		petRouter.GET("/custom/action/list", middleware.JWTTokenMiddleware(), handler.GetCustomActionList)

		petRouter.POST("/consume/create", handler.CreateConsumeAction)
		petRouter.GET("/consume/list", handler.GetPetConsumeList)
		petRouter.POST("/custom/consume/create", middleware.JWTTokenMiddleware(), handler.CreateCustomConsumeAction)
		petRouter.GET("/custom/consume/list", middleware.JWTTokenMiddleware(), handler.GetPetCustomConsumeList)
	}

	actionRouter := r.Group("v1/action")
	{
		// 获取日常列表
		actionRouter.GET("/list", handler.GetPetActionList)
		// 用户获取自己的列表
		actionRouter.GET("/custom/list", middleware.JWTTokenMiddleware(), handler.GetCustomActionList)
		// 用户自己创建
		actionRouter.POST("/create/custom", middleware.JWTTokenMiddleware(), handler.CreateConsumeAction)

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
