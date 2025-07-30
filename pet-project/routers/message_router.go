package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/handler"
	"pet-project/middleware"
)

func RegisterMessageRouters(r *gin.Engine) {
	router := r.Group("/v1/msg")
	{
		router.POST("/like", middleware.JWTTokenMiddleware(), handler.LikeMessageHandler)
		router.POST("/collection", middleware.JWTTokenMiddleware(), handler.CollectionMessageHandler)
		router.POST("/list", middleware.JWTTokenMiddleware(), handler.MessageListHandler)
		router.GET("/unread", middleware.JWTTokenMiddleware(), handler.UnreadNumberHandler)
	}
}
