package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"strconv"
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
		msgInfo := models.MessageModel{
			MessageType: 1,
			MessageId:   likeMsg.ID,
			FromUid:     likeMsg.FromUid,
			ToUid:       likeMsg.ToUid,
		}
		db.DB.Create(&msgInfo)
		response.Success(c, map[string]interface{}{})
	} else {
		insertResult := db.DB.Model(&findLike).Update("like_status", likeMsg.LikeStatus)
		if insertResult.Error != nil {
			response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
			return
		}
		if likeMsg.LikeStatus == 1 {
			// 创建一条消息
			msgInfo := models.MessageModel{
				MessageType: 1,
				MessageId:   likeMsg.ID,
				FromUid:     likeMsg.FromUid,
				ToUid:       likeMsg.ToUid,
			}
			db.DB.Create(&msgInfo)
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
		// 创建一条消息
		msgInfo := models.MessageModel{
			MessageType: 2,
			MessageId:   collection.ID,
			FromUid:     collection.FromUid,
			ToUid:       collection.ToUid,
		}
		db.DB.Create(&msgInfo)
		response.Success(c, map[string]interface{}{})
	} else {
		insertResult := db.DB.Model(&findCollection).Update("collection_status", collection.CollectionStatus)
		if insertResult.Error != nil {
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
			return
		}
		if collection.CollectionStatus == 1 {
			// 创建一条消息
			msgInfo := models.MessageModel{
				MessageType: 2,
				MessageId:   collection.ID,
				FromUid:     collection.FromUid,
				ToUid:       collection.ToUid,
			}
			db.DB.Create(&msgInfo)
		}
		response.Success(c, map[string]interface{}{})
	}
}

func MessageListHandler(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var pageSize = c.PostForm("pageSize")
	var pageNum = c.PostForm("pageNum")
	size, _ := strconv.Atoi(pageSize)
	num, _ := strconv.Atoi(pageNum)
	offer := (num - 1) * size
	var msgList []models.MessageModel
	result := db.DB.Model(models.MessageModel{}).Where("userId=?", userId).Offset(offer).Limit(size).Find(&msgList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, msgList)
}
