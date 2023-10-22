package main

import (
	"pet-project/config/settings"
	"pet-project/db"
	"pet-project/routers"
)

func main() {
	settings.ConfigEnvironment()
	db.LinkInit()
	r := routers.RegisterRouter()
	r.Run(":8082")
}
