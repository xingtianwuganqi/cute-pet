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
	userId := c.MustGet("userId").(uint)
	var statusModel models.LikeMessageModel
	if err := c.ShouldBind(&statusModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	// 查询到这条帖子
	var postInfo models.PostModel
	if statusModel.LikeType == 1 {
		result := db.DB.Where("id = ?", statusModel.LikeId).First(&postInfo)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
			return
		}
	}

	var likeStatus models.LikeMessageModel
	likeResult := db.DB.Model(&models.LikeMessageModel{}).Where("like_id = ?", statusModel.LikeId).First(&likeStatus)
	if errors.Is(likeResult.Error, gorm.ErrRecordNotFound) {
		likeStatus.LikeType = statusModel.LikeType
		likeStatus.LikeId = statusModel.LikeId
		likeStatus.LikeStatus = statusModel.LikeStatus
		likeStatus.FromUid = userId
		likeStatus.ToUid = statusModel.ToUid
		db.DB.Create(&likeStatus)
	} else {
		db.DB.Model(&likeStatus).Update("like_status", statusModel.LikeStatus)
	}

	// 更新帖子点赞数
	if statusModel.LikeStatus == 1 {
		if statusModel.LikeType == 1 {
			var num = postInfo.LikeNum + 1
			db.DB.Model(&postInfo).Update("like_num", num)
		}
		// 新增一条消息
		msgInfo := models.MessageModel{
			MessageType: 1,
			MessageId:   statusModel.LikeId,
			FromUid:     statusModel.FromUid,
			ToUid:       statusModel.ToUid,
		}
		db.DB.Model(models.MessageModel{}).Create(&msgInfo)
	} else {
		if statusModel.LikeType == 1 {
			var num uint
			if postInfo.LikeNum > 0 {
				num = postInfo.LikeNum - 1
			} else {
				num = 0
			}
			db.DB.Model(&postInfo).Update("like_num", num)
		}
	}
	// 查询到这条帖子
	response.Success(c, statusModel)
}

func CollectionMessageHandler(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	collectionModel := models.CollectionMessageModel{}
	if err := c.ShouldBind(&collectionModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 查询帖子是否存在
	postInfo := models.PostModel{}
	if collectionModel.CollectionType == 1 {
		postResult := db.DB.Where("id = ?", collectionModel.CollectionId).First(&postInfo)
		if errors.Is(postResult.Error, gorm.ErrRecordNotFound) {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
			return
		}
	}

	var collectionStatus models.CollectionMessageModel
	findResult := db.DB.Model(models.CollectionMessageModel{}).Where("collection_id = ?", collectionModel.CollectionId).First(&collectionStatus)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		collectionStatus.CollectionType = collectionModel.CollectionType
		collectionStatus.CollectionId = collectionModel.CollectionId
		collectionStatus.CollectionStatus = collectionModel.CollectionStatus
		collectionStatus.FromUid = userId
		collectionStatus.ToUid = collectionModel.ToUid
		db.DB.Create(&collectionStatus)
	} else {
		db.DB.Model(models.CollectionMessageModel{}).Where("collection_id = ?", collectionModel.CollectionId).Update("collection_status", collectionModel.CollectionStatus)
	}

	if collectionModel.CollectionStatus == 1 {

		if collectionModel.CollectionType == 1 {
			num := postInfo.CollectionNum + 1
			db.DB.Model(&postInfo).Update("collection_num", num)
		}

		// 新增一条消息
		msgInfo := models.MessageModel{
			MessageType: 2,
			MessageId:   collectionModel.CollectionId,
			FromUid:     collectionModel.FromUid,
			ToUid:       collectionModel.ToUid,
		}
		db.DB.Create(&msgInfo)
	} else {
		if collectionModel.CollectionType == 1 {
			var num uint
			if postInfo.CollectionNum > 0 {
				num = postInfo.CollectionNum - 1
			} else {
				num = 0
			}
			db.DB.Model(&postInfo).Update("collection_num", num)
		}
	}

	response.Success(c, collectionStatus)
}

// MessageListHandler 消息列表
func MessageListHandler(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var typeModel models.MessageListType
	if err := c.ShouldBind(&typeModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offer := (typeModel.PageNum - 1) * typeModel.PageSize
	var msgList []models.MessageModel

	var result *gorm.DB
	switch typeModel.MessageType {
	case 0:
		result = db.DB.Model(models.MessageModel{}).
			Where("to_uid = ?", userId).
			Offset(offer).Limit(typeModel.PageSize).
			Order("created_at DESC").
			Find(&msgList)
	case 1, 2:
		result = db.DB.Model(models.MessageModel{}).
			Where("to_uid = ? AND message_type = ?", userId, typeModel.MessageType).
			Offset(offer).Limit(typeModel.PageSize).
			Order("created_at DESC").
			Find(&msgList)
	default:
		result = db.DB.Model(models.MessageModel{}).
			Where("to_uid = ? AND message_type IN ?", userId, []int{3, 4}).
			Offset(offer).Limit(typeModel.PageSize).
			Order("created_at DESC").
			Find(&msgList)
	}

	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	// 提取所有 ID
	var ids []uint
	for _, msg := range msgList {
		if msg.IsRead == false { // 只更新未读的
			ids = append(ids, msg.ID)
		}
	}

	// 批量更新 is_read 字段为 1
	if len(ids) > 0 {
		if err := db.DB.Model(&models.MessageModel{}).
			Where("id IN ?", ids).
			Update("is_read", true).Error; err != nil {
			// 更新失败不终止，但你可以记录日志或返回警告
		} else {
			// 更新内存中的 msgList 中的 isRead 字段，返回给前端时一致
			for i := range msgList {
				msgList[i].IsRead = true
			}
		}
	}
	response.Success(c, msgList)
}

// UnreadNumberHandler 未读消息数量
func UnreadNumberHandler(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)

	var likeNum int64 = 0
	var collectionNum int64 = 0
	var commentNum int64 = 0

	db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type = 1 AND is_read = false", userId).
		Count(&likeNum)
	db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type = 2 AND is_read = false", userId).
		Count(&collectionNum)
	db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type IN ? AND is_read = false", userId, []int{3, 4}).
		Count(&commentNum)

	response.Success(c, gin.H{
		"likeNum":       likeNum,
		"collectionNum": collectionNum,
		"commentNum":    commentNum,
	})
}

// CommentHandler 评论
func CommentHandler(c *gin.Context) {
	var commentModel models.CommentModel
	if err := c.ShouldBind(&commentModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	result := db.DB.Model(models.CommentModel{}).Create(&commentModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, commentModel)
}

func DeleteCommentHandler(c *gin.Context) {
	var commentId = c.Param("commentId")
	result := db.DB.Model(models.CommentModel{}).Delete(&models.CommentModel{}, commentId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

// ReplyHandler 回复
func ReplyHandler(c *gin.Context) {
	var replyModel models.ReplyModel
	if err := c.ShouldBind(&replyModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	result := db.DB.Model(models.ReplyModel{}).Create(&replyModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, replyModel)
}

func DeleteReplyHandler(c *gin.Context) {
	var replyId = c.Param("replyId")
	result := db.DB.Model(models.ReplyModel{}).Delete(&models.ReplyModel{}, replyId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

func GetCommentListHandler(c *gin.Context) {
	param := models.CommentListModel{}
	if err := c.ShouldBind(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (param.PageNum - 1) * param.PageSize
	var commentList []models.CommentModel
	result := db.DB.
		Preload("ReplyList", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC").Limit(5)
		}).
		Where("topic_id = ?", param.TopicId).
		Offset(offset).
		Limit(param.PageSize).
		Order("created_at DESC").
		Find(&commentList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, commentList)
}

func GetReplyListHandler(c *gin.Context) {
	param := models.ReplyListModel{}
	if err := c.ShouldBind(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (param.PageNum - 1) * param.PageSize
	var replyList []models.ReplyModel
	result := db.DB.Model(&models.ReplyModel{}).Where("comment_id = ?", param.CommentId).
		Offset(offset).
		Limit(param.PageSize).
		Order("created_at DESC").
		Find(&replyList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, replyList)
}
