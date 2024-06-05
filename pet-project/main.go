package main

import (
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
	if settings.Conf.App.Env == "production" {
		err := r.Run(":8082")
		if err != nil {
			return
		}
	} else {
		err := r.Run(":8086")
		if err != nil {
			return
		}
	}

}
