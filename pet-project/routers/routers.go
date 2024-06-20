package routers

import (
	"pet-project/middleware"
	"pet-project/service"
	"pet-project/settings"

	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	if settings.Conf.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	bundle := service.ReloadLocalBundle()
	r := gin.Default()
	// 添加admin配置
	service.AdminConfig(r)
	r.Use(middleware.LocaleMiddleware(bundle))
	RegisterUserRouter(r)
	RegisterRecordRouter(r)
	RegisterTestRouter(r)
	return r
}
