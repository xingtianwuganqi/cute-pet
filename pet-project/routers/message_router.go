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
		router.POST("/comment", middleware.JWTTokenMiddleware(), handler.CommentHandler)
		router.DELETE("/comment/:commentId", middleware.JWTTokenMiddleware(), handler.DeleteCommentHandler)
		router.POST("/reply", middleware.JWTTokenMiddleware(), handler.ReplyHandler)
		router.DELETE("/reply/:replyId", middleware.JWTTokenMiddleware(), handler.DeleteReplyHandler)
		router.POST("/comment/list", handler.GetCommentListHandler)
		router.POST("/reply/list", handler.GetReplyListHandler)
	}
}
