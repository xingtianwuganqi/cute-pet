package main

import (
	"fmt"
	"pet-project/db"
	"pet-project/routers"
	"pet-project/settings"
)

func main() {
	if err := settings.LoadConfig(); err != nil {
		panic(err)
	}
	db.LinkDataBase()
	r := routers.RegisterRouter()
	port := fmt.Sprintf(":%d", settings.Conf.App.Port)
	err := r.Run(port)
	if err != nil {
		return
	}
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
