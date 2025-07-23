package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
)

func CreateCurrentTopic(c *gin.Context) {
	var topicModel models.TopicModel
	if err := c.ShouldBind(&topicModel); err != nil {
	}
	result := db.DB.Omit(clause.Associations).Create(&topicModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func UserCreateTopic(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var topicModel models.TopicModel
	if err := c.ShouldBind(&topicModel); err != nil {
	}
	topicModel.UserId = userId
	result := db.DB.Omit(clause.Associations).Create(&topicModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}
