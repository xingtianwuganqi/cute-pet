package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	userRouter := r.Group("/v1/user")
	{
		userRouter.POST("/register", handler.UserRegister)
		userRouter.POST("/login", handler.UserPhoneLogin)
		userRouter.POST("/email/code", handler.GetEmailCode)
		userRouter.POST("/phone/code", handler.GetPhoneCode)
		userRouter.POST("/check/code", handler.CheckRdbCode)
		userRouter.POST("/find/password", handler.UserFindPassword)
		userRouter.POST("/qiniu/upload", middleware.JWTTokenMiddleware(), handler.GetQiNiuToken)
		userRouter.POST("/pwd/upload", middleware.JWTTokenMiddleware(), handler.UserUpdatePassword)
		userRouter.POST("/info/upload", middleware.JWTTokenMiddleware(), handler.UploadUserInfo)
		userRouter.GET("/encryption/code", handler.GetEncryptionCode)
		userRouter.GET("/info", middleware.JWTTokenMiddleware(), handler.GetUserInfo)
		userRouter.POST("/deactivate", middleware.JWTTokenMiddleware(), handler.UserDeactivate)

		userRouter.POST("/suggestion", middleware.JWTTokenMiddleware(), handler.CreateSuggestion)
		userRouter.POST("/ip", handler.GetIpInfo)
	}

}
