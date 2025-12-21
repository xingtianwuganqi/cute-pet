package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middleware.AdminOnly())
	{
		admin.GET("/categories/common", handler.GetCommonCategories)
		admin.POST("/categories/common", handler.CreateCommonCategory)
		admin.DELETE("/categories/common/:id", handler.DeleteCommonCategory)
		admin.POST("/categories/list", handler.CreateCommonCategoryList)

		admin.GET("user/list", handler.GetUserList)

	}
}
