package routers

import (
	"pet-project/handler"
	"pet-project/internal"
	"pet-project/middleware"
	"pet-project/settings"

	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	if settings.Conf.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	bundle := internal.ReloadLocalBundle()
	r := gin.Default()
	
	// 添加健康检查路由
	r.GET("/health", handler.HealthCheck)
	
	// 添加admin配置
	//internal.AdminConfig(r)
	r.Use(middleware.LocaleMiddleware(bundle))
	RegisterUserRouter(r)
	RegisterRecordRouter(r)
	RegisterTestRouter(r)
	RegisterMessageRouters(r)
	RegisterTopicRouter(r)
	RegisterAdminRouter(r)
	return r
}