package routers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	r := gin.Default()
	RegisterUserRouter(r)
	RegisterRecordRouter(r)
	return r
}
