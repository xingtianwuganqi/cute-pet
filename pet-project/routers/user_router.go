package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/handler"
	"pet-project/middleware"
)

func RegisterUserRouter(r *gin.Engine) {
	userRouter := r.Group("/v1/user")
	{
		userRouter.POST("/register", handler.UserRegister)
		userRouter.POST("/login", handler.UserPhoneLogin)
		userRouter.POST("/tencent/code", handler.GetTencentCode)
		userRouter.POST("/find/password", handler.UserFindPassword)
		userRouter.POST("/upload", middleware.JWTTokenMiddleware(), handler.GetQiNiuToken)
		userRouter.POST("/pwd/upload", middleware.JWTTokenMiddleware(), handler.UserUpdatePassword)
	}

}
