package service

import (
	"errors"
	"pet-project/db"
	"pet-project/models"
	"pet-project/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetStatusTopicListService 获取待审核话题列表的业务逻辑
func GetStatusTopicListService(status uint, pageNum, pageSize int) ([]models.TopicModel, error) {
	var topicModel []models.TopicModel
	offset := (pageNum - 1) * pageSize
	result := db.DB.Model(models.TopicModel{}).Preload("User").
		Where("topic_status=?", status).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&topicModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return topicModel, nil
}

// ChangeTopicStatusService 更改话题状态的业务逻辑
func ChangeTopicStatusService(topicId uint, status uint) error {
	result := db.DB.Model(models.TopicModel{}).
		Preload("User").
		Where("id=?", topicId).
		Update("topic_status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetTopicListService 获取话题列表的业务逻辑
func GetTopicListService(pageNum, pageSize int) ([]models.TopicModel, error) {
	var topicModel []models.TopicModel
	offset := (pageNum - 1) * pageSize
	result := db.DB.Model(models.TopicModel{}).
		Where("topic_status=?", 1).
		Preload("User").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&topicModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return topicModel, nil
}

// UserCreateTopicService 用户创建话题的业务逻辑
func UserCreateTopicService(userId uint, topicModel models.TopicModel) error {
	topicModel.UserId = userId
	// 过滤敏感词
	filter := util.NewWordFilter()
	newTitle := filter.Replace(topicModel.Title)
	newContent := filter.Replace(topicModel.Desc)
	topicModel.Title = newTitle
	topicModel.Desc = newContent
	result := db.DB.Omit(clause.Associations).Create(&topicModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteUserTopicService 删除用户话题的业务逻辑
func DeleteUserTopicService(userId uint, topicId uint) error {
	topicModel := models.TopicModel{}
	result := db.DB.Where("id=? and user_id=?", topicId, userId).First(&topicModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("data not exist")
	}
	if result.Error != nil {
		return result.Error
	}
	result = db.DB.Delete(&topicModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreatePostService 创建帖子的业务逻辑
func CreatePostService(userId uint, postModel models.PostModel) error {
	// 过滤敏感词
	filter := util.NewWordFilter()
	newContent := filter.Replace(postModel.Content)
	postModel.UserId = userId
	postModel.Content = newContent
	result := db.DB.Omit("User").Create(&postModel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetPostListService 获取帖子列表的业务逻辑
func GetPostListService(pageNum, pageSize int, userId uint) ([]models.PostModel, error) {
	var postModels []models.PostModel
	offset := (pageNum - 1) * pageSize
	result := db.DB.Model(models.PostModel{}).
		Preload("User").
		Offset(offset).Limit(pageSize).
		Order("created_at desc").
		Find(&postModels)
	if result.Error != nil {
		return nil, result.Error
	}

	// 查询状态
	if userId != 0 {
		var postIds []uint
		for _, post := range postModels {
			postIds = append(postIds, post.ID)
		}

		var likedPosts []models.LikeMessageModel
		db.DB.Where("from_uid = ? AND like_id IN ? AND like_status = ?", userId, postIds, 1).Find(&likedPosts)
		var collectedPosts []models.CollectionMessageModel
		db.DB.Where("from_uid = ? AND collection_id IN ? AND collection_status = ?", userId, postIds, 1).Find(&collectedPosts)

		likedMap := map[uint]bool{}
		for _, l := range likedPosts {
			likedMap[l.LikeId] = true
		}

		collectedMap := map[uint]bool{}
		for _, collect := range collectedPosts {
			collectedMap[collect.CollectionId] = true
		}

		for i := range postModels {
			if likedMap[postModels[i].ID] {
				postModels[i].LikeStatus = 1
			}
			if collectedMap[postModels[i].ID] {
				postModels[i].CollectionStatus = 1
			}
		}
	}

	return postModels, nil
}

// DeletePostService 删除帖子的业务逻辑
func DeletePostService(userId uint, postId uint) error {
	result := db.DB.Where("id=? and user_id=?", postId, userId).Delete(&models.PostModel{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}