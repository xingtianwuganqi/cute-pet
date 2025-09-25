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
		recordRouter.GET("", middleware.JWTTokenMiddleware(), handler.GetRecordList)
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
		categoryRouter.GET("/common", middleware.JWTTokenMiddleware(), handler.GetCommonCategories)
		categoryRouter.POST("/common", middleware.JWTTokenMiddleware(), handler.CreateCommonCategory)
		categoryRouter.DELETE("/common/:id", middleware.JWTTokenMiddleware(), handler.DeleteCommonCategory)
		categoryRouter.POST("/list", middleware.JWTTokenMiddleware(), handler.CreateCommonCategoryList)

		categoryRouter.GET("", middleware.JWTTokenMiddleware(), handler.GetRecordCategoryList)
		categoryRouter.POST("", middleware.JWTTokenMiddleware(), handler.CreateRecordCategory)
		categoryRouter.PATCH("", middleware.JWTTokenMiddleware(), handler.UpdateRecordCategory)
		categoryRouter.DELETE("/:id", middleware.JWTTokenMiddleware(), handler.DeleteRecordCategory)
	}
}
