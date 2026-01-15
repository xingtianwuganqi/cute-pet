package service

import (
	"errors"
	"pet-project/db"
	"pet-project/models"

	"gorm.io/gorm"
)

// LikeMessageService 点赞消息的业务逻辑
func LikeMessageService(userId uint, statusModel models.LikeMessageModel) (models.LikeMessageModel, error) {
	// 查询到这条帖子
	var postInfo models.PostModel
	if statusModel.LikeType == 1 {
		result := db.DB.Where("id = ?", statusModel.LikeId).First(&postInfo)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return models.LikeMessageModel{}, errors.New("query error")
			}
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
			return models.LikeMessageModel{}, createErr
		}
	} else {
		updateErr := db.DB.Model(&likeStatus).Update("like_status", statusModel.LikeStatus).Error
		if updateErr != nil {
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

	return statusModel, nil
}

// CollectionMessageService 收藏消息的业务逻辑
func CollectionMessageService(userId uint, collectionModel models.CollectionMessageModel) (models.CollectionMessageModel, error) {
	// 查询帖子是否存在
	postInfo := models.PostModel{}
	if collectionModel.CollectionType == 1 {
		postResult := db.DB.Where("id = ?", collectionModel.CollectionId).First(&postInfo)
		if errors.Is(postResult.Error, gorm.ErrRecordNotFound) {
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
			return models.CollectionMessageModel{}, err
		}
	} else {
		err := db.DB.Model(models.CollectionMessageModel{}).Where("collection_id = ?", collectionModel.CollectionId).Update("collection_status", collectionModel.CollectionStatus).Error
		if err != nil {
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

	return collectionStatus, nil
}

// MessageListService 消息列表的业务逻辑
func MessageListService(userId uint, typeModel models.MessageListType) ([]models.MessageModel, error) {
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
			// 更新失败不终止，但你可以记录日志或返回警告
		} else {
			// 更新内存中的 msgList 中的 isRead 字段，返回给前端时一致
			for i := range msgList {
				msgList[i].IsRead = true
			}
		}
	}
	return msgList, nil
}

// UnreadNumberService 未读消息数量的业务逻辑
func UnreadNumberService(userId uint) (map[string]int64, error) {
	var likeNum int64 = 0
	var collectionNum int64 = 0
	var commentNum int64 = 0

	err := db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type = 1 AND is_read = false", userId).
		Count(&likeNum).Error
	if err != nil {
		return nil, err
	}
	err = db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type = 2 AND is_read = false", userId).
		Count(&collectionNum).Error
	if err != nil {
		return nil, err
	}
	err = db.DB.Model(&models.MessageModel{}).
		Where("to_uid = ? AND message_type IN ? AND is_read = false", userId, []int{3, 4}).
		Count(&commentNum).Error
	if err != nil {
		return nil, err
	}

	return map[string]int64{
		"likeNum":       likeNum,
		"collectionNum": collectionNum,
		"commentNum":    commentNum,
	}, nil
}

// CommentService 评论的业务逻辑
func CommentService(commentModel models.CommentModel) (models.CommentModel, error) {
	result := db.DB.Model(models.CommentModel{}).Create(&commentModel)
	if result.Error != nil {
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
	return commentModel, nil
}

// DeleteCommentService 删除评论的业务逻辑
func DeleteCommentService(commentId uint) error {
	var commentModel models.CommentModel
	result := db.DB.Model(models.CommentModel{}).Where("id = ?", commentId).First(&commentModel)
	if result.Error != nil {
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
	return nil
}

// ReplyService 回复的业务逻辑
func ReplyService(replyModel models.ReplyModel) (models.ReplyModel, error) {
	result := db.DB.Model(models.ReplyModel{}).Create(&replyModel)
	if result.Error != nil {
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

	return replyModel, nil
}

// DeleteReplyService 删除回复的业务逻辑
func DeleteReplyService(replyId uint) error {
	result := db.DB.Model(models.ReplyModel{}).Delete(&models.ReplyModel{}, replyId)
	if result.Error != nil {
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
	return nil
}

// GetCommentListService 获取评论列表的业务逻辑
func GetCommentListService(param models.CommentListModel) ([]models.CommentModel, error) {
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
		return nil, result.Error
	}
	return commentList, nil
}

// GetReplyListService 获取回复列表的业务逻辑
func GetReplyListService(param models.ReplyListModel) ([]models.ReplyModel, error) {
	offset := (param.PageNum - 1) * param.PageSize
	var replyList []models.ReplyModel
	result := db.DB.Model(&models.ReplyModel{}).Where("comment_id = ?", param.CommentId).
		Offset(offset).
		Limit(param.PageSize).
		Order("created_at DESC").
		Find(&replyList)
	if result.Error != nil {
		return nil, result.Error
	}
	return replyList, nil
}