package main

import (
	"strconv"
	"pet-project/db"
	"pet-project/logger"
	"pet-project/routers"
	"pet-project/settings"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	if err := settings.LoadConfig(); err != nil {
		panic(err)
	}

	// 初始化日志系统
	logger.InitLogger()
	defer logger.Sync()

	// 记录启动日志
	logger.Logger.Info("Starting application",
		zap.String("env", settings.Conf.App.Env),
		zap.Int("port", settings.Conf.App.Port))

	// 初始化数据库连接
	db.LinkDataBase()

	// 启动路由
	r := routers.RegisterRouter()
	r.Run(":" + strconv.Itoa(settings.Conf.App.Port))
}

/*
windows启动redis命令

redis-server.exe redis.windows.conf


mac 启动redis

cd /usr/local/bin
redis-server
// 终止
redis-cli shutdown
*/
