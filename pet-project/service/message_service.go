package service

import (
	"errors"
	"pet-project/db"
	"pet-project/models"
	"pet-project/logger"

	"gorm.io/gorm"
	"go.uber.org/zap"
)

// LikeMessageService 点赞消息的业务逻辑
func LikeMessageService(userId uint, statusModel models.LikeMessageModel) (models.LikeMessageModel, error) {
	logger.Logger.Info("LikeMessageService called", zap.Uint("userId", userId), zap.Uint("likeId", statusModel.LikeId))

	// 查询到这条帖子
	var postInfo models.PostModel
	if statusModel.LikeType == 1 {
		result := db.DB.Where("id = ?", statusModel.LikeId).First(&postInfo)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				logger.Logger.Warn("Post not found for like", zap.Uint("postId", statusModel.LikeId))
				return models.LikeMessageModel{}, errors.New("query error")
			}
			logger.Logger.Error("Database error in LikeMessageService", zap.Error(result.Error))
			return models.LikeMessageModel{}, result.Error
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
		createErr := db.DB.Create(&likeStatus).Error
		if createErr != nil {
			logger.Logger.Error("Database error creating like status", zap.Error(createErr))
			return models.LikeMessageModel{}, createErr
		}
	} else {
		updateErr := db.DB.Model(&likeStatus).Update("like_status", statusModel.LikeStatus).Error
		if updateErr != nil {
			logger.Logger.Error("Database error updating like status", zap.Error(updateErr))
			return models.LikeMessageModel{}, updateErr
		}
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

	logger.Logger.Info("Successfully processed like message", zap.Uint("userId", userId), zap.Uint("likeId", statusModel.LikeId))
	return statusModel, nil
}

// CollectionMessageService 收藏消息的业务逻辑
func CollectionMessageService(userId uint, collectionModel models.CollectionMessageModel) (models.CollectionMessageModel, error) {
	logger.Logger.Info("CollectionMessageService called", zap.Uint("userId", userId), zap.Uint("collectionId", collectionModel.CollectionId))

	// 查询帖子是否存在
	postInfo := models.PostModel{}
	if collectionModel.CollectionType == 1 {
		postResult := db.DB.Where("id = ?", collectionModel.CollectionId).First(&postInfo)
		if errors.Is(postResult.Error, gorm.ErrRecordNotFound) {
			logger.Logger.Warn("Post not found for collection", zap.Uint("postId", collectionModel.CollectionId))
			return models.CollectionMessageModel{}, errors.New("query error")
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
		err := db.DB.Create(&collectionStatus).Error
		if err != nil {
			logger.Logger.Error("Database error creating collection status", zap.Error(err))
			return models.CollectionMessageModel{}, err
		}
	} else {
		err := db.DB.Model(models.CollectionMessageModel{}).Where("collection_id = ?", collectionModel.CollectionId).Update("collection_status", collectionModel.CollectionStatus).Error
		if err != nil {
			logger.Logger.Error("Database error updating collection status", zap.Error(err))
			return models.CollectionMessageModel{}, err
		}
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

	logger.Logger.Info("Successfully processed collection message", zap.Uint("userId", userId), zap.Uint("collectionId", collectionModel.CollectionId))
	return collectionStatus, nil
}

// MessageListService 消息列表的业务逻辑
func MessageListService(userId uint, typeModel models.MessageListType) ([]models.MessageModel, error) {
	logger.Logger.Info("MessageListService called", zap.Uint("userId", userId), zap.Uint("messageType", typeModel.MessageType))

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
		logger.Logger.Error("Database error in MessageListService", zap.Error(result.Error))
		return nil, result.Error
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
			logger.Logger.Warn("Failed to update message read status", zap.Error(err))
		} else {
			// 更新内存中的 msgList 中的 isRead 字段，返回给前端时一致
			for i := range msgList {
				msgList[i].IsRead = true
			}
		}
	}
	logger.Logger.Info("Successfully retrieved message list", zap.Int("count", len(msgList)), zap.Uint("userId", userId))
	return msgList, nil
}

// UnreadNumberService 未读消息数量的业务逻辑
func UnreadNumberService(userId uint) (map[string]int64, error) {
	logger.Logger.Info("UnreadNumberService called", zap.Uint("userId", userId))

	var likeNum int64 = 0
	var collectionNum int64 = 0
	var commentNum int64 = 0

	err := db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type = 1 AND is_read = false", userId).
		Count(&likeNum).Error
	if err != nil {
		logger.Logger.Error("Database error counting like messages", zap.Error(err))
		return nil, err
	}
	err = db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type = 2 AND is_read = false", userId).
		Count(&collectionNum).Error
	if err != nil {
		logger.Logger.Error("Database error counting collection messages", zap.Error(err))
		return nil, err
	}
	err = db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type IN ? AND is_read = false", userId, []int{3, 4}).
		Count(&commentNum).Error
	if err != nil {
		logger.Logger.Error("Database error counting comment messages", zap.Error(err))
		return nil, err
	}

	result := map[string]int64{
		"likeNum":       likeNum,
		"collectionNum": collectionNum,
		"commentNum":    commentNum,
	}

	logger.Logger.Info("Successfully retrieved unread numbers", zap.Uint64("likeNum", uint64(likeNum)), zap.Uint64("collectionNum", uint64(collectionNum)), zap.Uint64("commentNum", uint64(commentNum)))
	return result, nil
}

// CommentService 评论的业务逻辑
func CommentService(commentModel models.CommentModel) (models.CommentModel, error) {
	logger.Logger.Info("CommentService called", zap.Uint("topicId", commentModel.TopicId), zap.Uint("fromUid", commentModel.FromUid))

	result := db.DB.Model(models.CommentModel{}).Create(&commentModel)
	if result.Error != nil {
		logger.Logger.Error("Database error in CommentService", zap.Error(result.Error))
		return models.CommentModel{}, result.Error
	}
	// 查询fromUser和toUser
	fromUser := models.UserInfo{}
	db.DB.Model(models.UserInfo{}).Where("id = ?", commentModel.FromUid).First(&fromUser)
	toUser := models.UserInfo{}
	db.DB.Model(models.UserInfo{}).Where("id = ?", commentModel.ToUid).First(&toUser)
	commentModel.FromUser = &fromUser
	commentModel.ToUser = &toUser
	// 查询帖子
	postInfo := models.PostModel{}
	db.DB.Model(models.PostModel{}).Where("id = ?", commentModel.TopicId).First(&postInfo)
	db.DB.Model(&postInfo).Update("comment_num", postInfo.CommentNum+1)
	logger.Logger.Info("Successfully processed comment", zap.Uint("commentId", commentModel.ID))
	return commentModel, nil
}

// DeleteCommentService 删除评论的业务逻辑
func DeleteCommentService(commentId uint) error {
	logger.Logger.Info("DeleteCommentService called", zap.Uint("commentId", commentId))

	var commentModel models.CommentModel
	result := db.DB.Model(models.CommentModel{}).Where("id = ?", commentId).First(&commentModel)
	if result.Error != nil {
		logger.Logger.Error("Database error finding comment to delete", zap.Error(result.Error))
		return result.Error
	}
	// 查询帖子
	postInfo := models.PostModel{}
	db.DB.Model(models.PostModel{}).Where("id = ?", commentModel.TopicId).First(&postInfo)
	if postInfo.CommentNum > 0 {
		num := postInfo.CommentNum - 1
		db.DB.Model(&postInfo).Update("comment_num", num)
	}
	db.DB.Model(models.CommentModel{}).Delete(&models.CommentModel{}, commentId)
	logger.Logger.Info("Successfully deleted comment", zap.Uint("commentId", commentId))
	return nil
}

// ReplyService 回复的业务逻辑
func ReplyService(replyModel models.ReplyModel) (models.ReplyModel, error) {
	logger.Logger.Info("ReplyService called", zap.Uint("commentId", replyModel.CommentId), zap.Uint("fromUid", replyModel.FromUid))

	result := db.DB.Model(models.ReplyModel{}).Create(&replyModel)
	if result.Error != nil {
		logger.Logger.Error("Database error in ReplyService", zap.Error(result.Error))
		return models.ReplyModel{}, result.Error
	}
	// 查询fromUser和toUser
	fromUser := models.UserInfo{}
	db.DB.Model(models.UserInfo{}).Where("id = ?", replyModel.FromUid).First(&fromUser)
	toUser := models.UserInfo{}
	db.DB.Model(models.UserInfo{}).Where("id = ?", replyModel.ToUid).First(&toUser)
	replyModel.FromUser = &fromUser
	replyModel.ToUser = &toUser

	// 查询到post
	commentModel := models.CommentModel{}
	db.DB.Model(models.PostModel{}).Where("id = ?", replyModel.CommentId).First(&commentModel)
	postModel := models.PostModel{}
	db.DB.Model(models.PostModel{}).Where("id = ?", commentModel.TopicId).First(&postModel)
	num := postModel.CommentNum + 1
	db.DB.Model(&postModel).Update("comment_num", num)

	logger.Logger.Info("Successfully processed reply", zap.Uint("replyId", replyModel.ID))
	return replyModel, nil
}

// DeleteReplyService 删除回复的业务逻辑
func DeleteReplyService(replyId uint) error {
	logger.Logger.Info("DeleteReplyService called", zap.Uint("replyId", replyId))

	result := db.DB.Model(models.ReplyModel{}).Delete(&models.ReplyModel{}, replyId)
	if result.Error != nil {
		logger.Logger.Error("Database error deleting reply", zap.Error(result.Error))
		return result.Error
	}
	replyModel := models.ReplyModel{}
	db.DB.Model(models.ReplyModel{}).Where("id = ?", replyId).First(&replyModel)
	commentModel := models.CommentModel{}
	db.DB.Model(models.CommentModel{}).Where("id = ?", replyModel.CommentId).First(&commentModel)
	postModel := models.PostModel{}
	db.DB.Model(models.PostModel{}).Where("id = ?", commentModel.TopicId).First(&postModel)
	if postModel.CommentNum > 0 {
		num := postModel.CommentNum - 1
		db.DB.Model(&postModel).Update("comment_num", num)
	}
	logger.Logger.Info("Successfully deleted reply", zap.Uint("replyId", replyId))
	return nil
}

// GetCommentListService 获取评论列表的业务逻辑
func GetCommentListService(param models.CommentListModel) ([]models.CommentModel, error) {
	logger.Logger.Info("GetCommentListService called", zap.Uint("topicId", param.TopicId))

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
		logger.Logger.Error("Database error in GetCommentListService", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.Logger.Info("Successfully retrieved comment list", zap.Int("count", len(commentList)))
	return commentList, nil
}

// GetReplyListService 获取回复列表的业务逻辑
func GetReplyListService(param models.ReplyListModel) ([]models.ReplyModel, error) {
	logger.Logger.Info("GetReplyListService called", zap.Uint("commentId", param.CommentId))

	offset := (param.PageNum - 1) * param.PageSize
	var replyList []models.ReplyModel
	result := db.DB.Model(&models.ReplyModel{}).Where("comment_id = ?", param.CommentId).
		Offset(offset).
		Limit(param.PageSize).
		Order("created_at DESC").
		Find(&replyList)
	if result.Error != nil {
		logger.Logger.Error("Database error in GetReplyListService", zap.Error(result.Error))
		return nil, result.Error
	}
	logger.Logger.Info("Successfully retrieved reply list", zap.Int("count", len(replyList)))
	return replyList, nil
}