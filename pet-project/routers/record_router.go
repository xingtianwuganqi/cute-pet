package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRecordRouter(r *gin.Engine) {
	recordRouter := r.Group("/v1/records")
	{
		recordRouter.POST("", middleware.JWTTokenMiddleware(), handler.CreateRecord)
		recordRouter.POST("/list", middleware.JWTTokenMiddleware(), handler.GetRecordList)
		recordRouter.DELETE("/:id", middleware.JWTTokenMiddleware(), handler.DeleteRecordInfo)
	}

	petRouter := r.Group("/v1/pets")
	{
		petRouter.GET("", middleware.JWTTokenMiddleware(), handler.GetPetList)
		petRouter.POST("", middleware.JWTTokenMiddleware(), handler.PetInfoCreate)
		petRouter.PATCH("", middleware.JWTTokenMiddleware(), handler.UpdatePetInfo)
		petRouter.DELETE("/:id", middleware.JWTTokenMiddleware(), handler.DeletePetInfo)

	}

	categoryRouter := r.Group("/v1/categories")
	{
		categoryRouter.GET("/common", handler.GetCommonCategories)
		categoryRouter.POST("/common", handler.CreateCommonCategory)
		categoryRouter.DELETE("/common/:id", handler.DeleteCommonCategory)

		categoryRouter.GET("", middleware.JWTTokenMiddleware(), handler.GetRecordCategoryList)
		categoryRouter.POST("", middleware.JWTTokenMiddleware(), handler.CreateRecordCategory)
		categoryRouter.DELETE("/:id", middleware.JWTTokenMiddleware(), handler.DeleteRecordCategory)
		categoryRouter.POST("/list", middleware.JWTTokenMiddleware(), handler.CreateCategoryList)
	}
}
