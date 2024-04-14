package routers

import (
	handler "pet-project/handler/api_handler"
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
		actionRouter.GET("/list", handler.GetPetActionList)
		actionRouter.GET("/custom/list", middleware.JWTTokenMiddleware(), handler.GetPetActionList)
		actionRouter.POST("/create/action", middleware.JWTTokenMiddleware(), handler.CreatePetActionType)
		actionRouter.POST("/create/custom", middleware.JWTTokenMiddleware(), handler.CreatePetCustomType)

	}
}
