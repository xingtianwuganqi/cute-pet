package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRecordRouter(r *gin.Engine) {
	recordRouter := r.Group("/v1/record")
	{
		recordRouter.POST("/create", middleware.JWTTokenMiddleware(), handler.CreateRecord)
		recordRouter.GET("/list", middleware.JWTTokenMiddleware(), handler.GetRecordList)
		recordRouter.DELETE("/delete/:id", middleware.JWTTokenMiddleware(), handler.DeleteRecordInfo)
	}

	petRouter := r.Group("/v1/pet")
	{
		petRouter.GET("/list", middleware.JWTTokenMiddleware(), handler.GetPetList)
		petRouter.POST("/create", middleware.JWTTokenMiddleware(), handler.PetInfoCreate)
		petRouter.PUT("/update", middleware.JWTTokenMiddleware(), handler.UpdatePetInfo)
		petRouter.DELETE("/delete/:id", middleware.JWTTokenMiddleware(), handler.DeletePetInfo)

		petRouter.POST("/create/action", handler.CreatePetActionType)
		petRouter.GET("/action/list", handler.GetPetActionList)
		petRouter.DELETE("/action/delete/:id", handler.DeletePetAction)
		petRouter.POST("/custom/action/create", middleware.JWTTokenMiddleware(), handler.CreatePetCustomAction)
		petRouter.GET("/custom/action/list", middleware.JWTTokenMiddleware(), handler.GetCustomActionList)
		petRouter.DELETE("/custom/action/delete/:id", middleware.JWTTokenMiddleware(), handler.DeletePetCustomAction)

		petRouter.POST("/consume/create", handler.CreateConsumeAction)
		petRouter.GET("/consume/list", handler.GetPetConsumeList)
		petRouter.DELETE("/consume/delete/:id", handler.DeletePetConsumeAction)
		petRouter.POST("/custom/consume/create", middleware.JWTTokenMiddleware(), handler.CreateCustomConsumeAction)
		petRouter.GET("/custom/consume/list", middleware.JWTTokenMiddleware(), handler.GetPetCustomConsumeList)
		petRouter.DELETE("/custom/consume/delete/:id", middleware.JWTTokenMiddleware(), handler.DeleteCustomConsumeAction)
	}
}
