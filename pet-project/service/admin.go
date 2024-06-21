package service

import (
	_ "github.com/GoAdminGroup/go-admin/adapter/gin" // 引入适配器，必须引入，如若不引入，则需要自己定义
	"github.com/GoAdminGroup/go-admin/engine"
	_ "github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/mysql" // 引入对应数据库引擎
	"github.com/GoAdminGroup/go-admin/modules/language"
	_ "github.com/GoAdminGroup/go-admin/template"
	_ "github.com/GoAdminGroup/go-admin/template/chartjs"
	_ "github.com/GoAdminGroup/themes/adminlte" // 引入主题，必须引入，不然报错
	"github.com/gin-gonic/gin"
	"pet-project/generate"
	"pet-project/settings"
	"strconv"
)

func AdminConfig(r *gin.Engine) {
	eng := engine.Default()
	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				Host:         settings.Conf.Database.Host,
				Port:         strconv.Itoa(settings.Conf.Database.Port),
				User:         settings.Conf.Database.Username,
				Pwd:          settings.Conf.Database.Password,
				Name:         settings.Conf.Database.DataBase,
				MaxIdleConns: 50,
				MaxOpenConns: 150,
				Driver:       config.DriverMysql,
			},
		},
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Theme:    "adminlte",
		Debug:    true,
		Language: language.CN,
	}
	r.Static("/uploads", "./uploads")
	err := eng.AddConfig(&cfg).AddGenerators(generate.Generators).Use(r)
	if err != nil {
		panic(err)
	}
}
