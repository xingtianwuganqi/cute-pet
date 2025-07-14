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

		petRouter.POST("/create/category", handler.CreateRecordCategory)
		petRouter.GET("/category/list", handler.GetRecordCategoryList)
		petRouter.DELETE("/category/delete/:id", handler.DeleteRecordCategory)
		petRouter.POST("/create/categoryList", handler.CreateCategoryList)

		petRouter.POST("/custom/category/create", middleware.JWTTokenMiddleware(), handler.CreateCustomCategory)
		petRouter.GET("/custom/category/list", middleware.JWTTokenMiddleware(), handler.GetCustomCategoryList)
		petRouter.DELETE("/custom/category/delete/:id", middleware.JWTTokenMiddleware(), handler.DeleteCustomCategory)
	}
}
