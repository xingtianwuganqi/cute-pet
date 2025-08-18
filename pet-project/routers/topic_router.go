package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTopicRouter(r *gin.Engine) {

	topicRouter := r.Group("/v1/topics")
	{
		topicRouter.POST("", middleware.JWTTokenMiddleware(), handler.UserCreateTopic)
		topicRouter.GET("", handler.GetTopicList)
		topicRouter.DELETE("/:id", middleware.JWTTokenMiddleware(), handler.DeleteUserTopic)
	}

	postRouter := r.Group("/v1/posts")
	{
		postRouter.POST("", middleware.JWTTokenMiddleware(), handler.CreatePost)
		postRouter.GET("", handler.GetPostList)
		postRouter.DELETE("/:id", middleware.JWTTokenMiddleware(), handler.DeletePost)

	}
}
