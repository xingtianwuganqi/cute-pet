package routers

import (
	"github.com/gin-gonic/gin"
	"pet-project/handler"
	"pet-project/middleware"
)

func RegisterMessageRouters(r *gin.Engine) {
	router := r.Group("/v1/messages")
	{
		router.POST("/like", middleware.JWTTokenMiddleware(), handler.LikeMessageHandler)
		router.POST("/collection", middleware.JWTTokenMiddleware(), handler.CollectionMessageHandler)
		router.GET("", middleware.JWTTokenMiddleware(), handler.MessageListHandler)
		router.GET("/unreads", middleware.JWTTokenMiddleware(), handler.UnreadNumberHandler)
		router.POST("/comments", middleware.JWTTokenMiddleware(), handler.CommentHandler)
		router.DELETE("/comments/:commentId", middleware.JWTTokenMiddleware(), handler.DeleteCommentHandler)
		router.POST("/replies", middleware.JWTTokenMiddleware(), handler.ReplyHandler)
		router.DELETE("/replies/:replyId", middleware.JWTTokenMiddleware(), handler.DeleteReplyHandler)
		router.GET("/comments", handler.GetCommentListHandler)
		router.GET("/replies", handler.GetReplyListHandler)
	}
}
