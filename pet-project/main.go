package main

import (
	"pet-project/db"
	"pet-project/routers"
)

func main() {
	db.LinkInit()
	r := routers.RegisterRouter()
	r.Run()
}
