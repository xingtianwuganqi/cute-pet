package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
)

func LikeMessageHandler(c *gin.Context) {
	var likeMsg models.LikeMessageModel
	if err := c.ShouldBind(&likeMsg); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	var findLike models.LikeMessageModel
	findResult := db.DB.Model(models.LikeMessageModel{}).Where("like_id", likeMsg.LikeId).First(&findLike)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		createResult := db.DB.Create(&likeMsg)
		if createResult.Error != nil {
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
			return
		}
		response.Success(c, map[string]interface{}{})
	} else {
		insertResult := db.DB.Model(&findLike).Update("like_status", likeMsg.LikeStatus)
		if insertResult.Error != nil {
			response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
			return
		}
		response.Success(c, map[string]interface{}{})
	}
}
