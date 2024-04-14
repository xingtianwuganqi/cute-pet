package routers

import (
	handler "pet-project/handler/api_handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	userRouter := r.Group("/v1/user")
	{
		userRouter.POST("/register", handler.UserRegister)
		userRouter.POST("/login", handler.UserPhoneLogin)
		userRouter.POST("/test", handler.TestNetworking)
	}

}
