package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/handler"
	"pet-project/middleware"
)

func RegisterRecordRouter(r *gin.Engine) {
	r.Group("v1/record")
	{
		r.GET("/petAction/list", handler.GetPetActionList)
		r.POST("/pet/create", middleware.JWTTokenMiddleware(), handler.PetInfoCreate)
	}
}
