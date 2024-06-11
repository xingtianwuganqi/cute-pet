package routers

import (
	"pet-project/handler/admin"
	"pet-project/handler/api"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {

	adminUserRouter := r.Group("/admin/user")
	{
		adminUserRouter.GET("/list", admin.UserList)

	}

	userRouter := r.Group("/v1/user")
	{
		userRouter.POST("/register", api.UserRegister)
		userRouter.POST("/login", api.UserPhoneLogin)
		userRouter.POST("/tencent/code", api.GetTencentCode)
		userRouter.POST("/find/password", api.UserFindPassword)
	}

}
