package routers

import (
	"pet-project/handler"
	"pet-project/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTopicRouter(r *gin.Engine) {

	topicRouter := r.Group("/v1/topic")
	{
		topicRouter.POST("/create", middleware.JWTTokenMiddleware(), handler.UserCreateTopic)
		topicRouter.GET("/list", handler.GetTopicList)
		topicRouter.DELETE("/delete/:id", middleware.JWTTokenMiddleware(), handler.DeleteUserTopic)
	}

	frontTopicRouter := r.Group("/v1/front/topic")
	{
		frontTopicRouter.GET("/status/:status", handler.GetStatusTopicList)
		frontTopicRouter.POST("/status/change", handler.ChangeTopicStatus)
	}

	postRouter := r.Group("/v1/post")
	{
		postRouter.POST("/create", middleware.JWTTokenMiddleware(), handler.CreatePost)
		postRouter.GET("/list", handler.GetPostList)
		postRouter.DELETE("/delete/:id", middleware.JWTTokenMiddleware(), handler.DeletePost)

	}
}
