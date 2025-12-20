package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(r *gin.Engine) {
	admin := r.Group("/admin/categories")
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/common", handler.GetCommonCategories)
		admin.POST("/common", handler.CreateCommonCategory)
		admin.DELETE("/common/:id", handler.DeleteCommonCategory)
		admin.POST("/list", handler.CreateCommonCategoryList)
	}
}
