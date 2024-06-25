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
	findResult := db.DB.Model(models.LikeMessageModel{}).Where("like_id = ?", likeMsg.LikeId).First(&findLike)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		createResult := db.DB.Create(&likeMsg)
		if createResult.Error != nil {
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
			return
		}
		// 创建一条消息
		//MsgInfo := models.MessageModel{
		//	MessageId: createResult.id
		//}
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

func CollectionMessageHandler(c *gin.Context) {
	var collection models.CollectionMessageModel
	if err := c.ShouldBind(&collection); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	var findCollection models.CollectionMessageModel
	findResult := db.DB.Model(models.CollectionMessageModel{}).Where("collection_id = ?", collection.CollectionId).First(&findCollection)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		createResult := db.DB.Create(&collection)
		if createResult.Error != nil {
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
			return
		}
		response.Success(c, map[string]interface{}{})
	} else {
		insertResult := db.DB.Model(&findCollection).Update("collection_status", collection.CollectionStatus)
		if insertResult.Error != nil {
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
			return
		}
	}
}
