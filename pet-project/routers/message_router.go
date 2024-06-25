package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/handler"
	"pet-project/middleware"
)

func RegisterMessageRouters(r *gin.Engine) {
	router := r.Group("/v1/msg")
	{
		router.POST("/like/action", middleware.JWTTokenMiddleware(), handler.LikeMessageHandler)
		router.POST("/collection/action", middleware.JWTTokenMiddleware(), handler.CollectionMessageHandler)
		router.POST("/list", middleware.JWTTokenMiddleware(), handler.MessageListHandler)
	}
}
