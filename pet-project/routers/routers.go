package routers

import (
	"pet-project/settings"

	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	if settings.Conf.App.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	RegisterUserRouter(r)
	RegisterRecordRouter(r)
	return r
}
