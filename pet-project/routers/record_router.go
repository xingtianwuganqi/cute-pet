package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/handler"
)

func RegisterRecordRouter(r *gin.Engine) {
	r.Group("v1/record")
	{
		r.POST("/create", handler.RecordCreate)
	}
}
