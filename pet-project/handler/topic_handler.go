// package service

// import (
// 	"errors"
// 	"pet-project/db"
// 	"pet-project/models"
// 	"pet-project/util"

// 	"gorm.io/gorm"
// 	"gorm.io/gorm/clause"
// )

// // GetStatusTopicListService 获取待审核话题列表
// func GetStatusTopicListService(status uint, pageNum, pageSize int) ([]models.TopicModel, error) {
// 	var topicModels []models.TopicModel
// 	offset := (pageNum - 1) * pageSize
// 	result := db.DB.Model(models.TopicModel{}).Preload("User").
// 		Where("topic_status=?", status).
// 		Offset(offset).
// 		Limit(pageSize).
// 		Order("created_at DESC").
// 		Find(&topicModels)
// 	return topicModels, result.Error
// }

// // ChangeTopicStatusService 修改话题状态
// func ChangeTopicStatusService(topicId uint, status int) error {
// 	result := db.DB.Model(models.TopicModel{}).
// 		Preload("User").
// 		Where("id=?", topicId).
// 		Update("topic_status", status)
// 	return result.Error
// }

// // GetTopicListService 获取话题列表
// func GetTopicListService(pageNum, pageSize int) ([]models.TopicModel, error) {
// 	var topicModels []models.TopicModel
// 	offset := (pageNum - 1) * pageSize
// 	result := db.DB.Model(models.TopicModel{}).
// 		Where("topic_status=?", 1).
// 		Preload("User").
// 		Offset(offset).Limit(pageSize).
// 		Order("created_at DESC").
// 		Find(&topicModels)
// 	return topicModels, result.Error
// }

// // UserCreateTopicService 用户创建话题
// func UserCreateTopicService(userId uint, topicModel models.TopicModel) error {
// 	topicModel.UserId = userId
// 	// 过滤敏感词
// 	filter := util.NewWordFilter()
// 	newTitle := filter.Replace(topicModel.Title)
// 	newContent := filter.Replace(topicModel.Desc)
// 	topicModel.Title = newTitle
// 	topicModel.Desc = newContent
// 	result := db.DB.Omit(clause.Associations).Create(&topicModel)
// 	return result.Error
// }

// // DeleteUserTopicService 删除用户话题
// func DeleteUserTopicService(userId, topicId uint) error {
// 	var topicModel models.TopicModel
// 	result := db.DB.Where("id=? and user_id=?", topicId, userId).First(&topicModel)
// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return errors.New("data not exist")
// 	}
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	result = db.DB.Delete(&topicModel)
// 	return result.Error
// }

// // CreatePostService 创建帖子
// func CreatePostService(userId uint, postModel models.PostModel) error {
// 	// 过滤敏感词
// 	filter := util.NewWordFilter()
// 	newContent := filter.Replace(postModel.Content)
// 	postModel.UserId = userId
// 	postModel.Content = newContent
// 	result := db.DB.Omit("User").Create(&postModel)
// 	return result.Error
// }

// // GetPostListService 获取帖子列表
// func GetPostListService(pageNum, pageSize int, userId uint) ([]models.PostModel, error) {
// 	var postModels []models.PostModel
// 	offset := (pageNum - 1) * pageSize
// 	result := db.DB.Model(models.PostModel{}).
// 		Preload("User").
// 		Offset(offset).Limit(pageSize).
// 		Order("created_at desc").
// 		Find(&postModels)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	// 查询状态
// 	if userId != 0 {
// 		var postIds []uint
// 		for _, post := range postModels {
// 			postIds = append(postIds, post.ID)
// 		}

// 		var likedPosts []models.LikeMessageModel
// 		db.DB.Where("from_uid = ? AND like_id IN ? AND like_status = ?", userId, postIds, 1).Find(&likedPosts)
// 		var collectedPosts []models.CollectionMessageModel
// 		db.DB.Where("from_uid = ? AND collection_id IN ? AND like_status = ?", userId, postIds, 1).Find(&collectedPosts)

// 		likedMap := make(map[uint]bool)
// 		for _, l := range likedPosts {
// 			likedMap[l.LikeId] = true
// 		}

// 		collectedMap := make(map[uint]bool)
// 		for _, collect := range collectedPosts {
// 			collectedMap[collect.CollectionId] = true
// 		}

// 		for i := range postModels {
// 			if likedMap[postModels[i].ID] {
// 				postModels[i].LikeStatus = 1
// 			}
// 			if collectedMap[postModels[i].ID] {
// 				postModels[i].CollectionStatus = 1
// 			}
// 		}
// 	}

// 	return postModels, nil
// }

// // DeletePostService 删除帖子
// func DeletePostService(userId, postId uint) error {
// 	result := db.DB.Where("id=? and user_id=?", postId, userId).Delete(&models.PostModel{})
// 	return result.Error
// }

package handler

import (
	"errors"
	"fmt"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Front api

// GetStatusTopicList 获取待审核话题列表
func GetStatusTopicList(c *gin.Context) {
	status := c.Param("status")
	var page models.PageModel
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	topics, err := service.GetStatusTopicListService(uint(getUintFromString(status)), page.PageNum, page.PageSize)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, topics)
}

func ChangeTopicStatus(c *gin.Context) {
	var topicModel models.TopicStatusModel
	if err := c.ShouldBind(&topicModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	err := service.ChangeTopicStatusService(topicModel.TopicId, topicModel.Status)
	if err != nil {
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	response.Success(c, nil)
}

// Topic User api

// GetTopicList 获取话题列表
func GetTopicList(c *gin.Context) {
	var page models.PageModel
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	topics, err := service.GetTopicListService(page.PageNum, page.PageSize)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, topics)
}

func UserCreateTopic(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var topicModel models.TopicModel
	if err := c.ShouldBind(&topicModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	err := service.UserCreateTopicService(userId, topicModel)
	if err != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func DeleteUserTopic(c *gin.Context) {
	userId, _ := c.Get("userId")
	topicId := c.Param("id")
	
	err := service.DeleteUserTopicService(userId.(uint), getUintFromString(topicId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "data not exist" {
			response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		} else {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		}
		return
	}
	response.Success(c, nil)
}

// Post api

func CreatePost(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	postModel := models.PostModel{}
	if err := c.ShouldBind(&postModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	err := service.CreatePostService(userId, postModel)
	if err != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func GetPostList(c *gin.Context) {
	var page models.PageModel
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	// 获取 userId
	userIdInterface, exists := c.Get("userId")
	var userId uint
	if exists {
		userId = userIdInterface.(uint)
	}
	
	posts, err := service.GetPostListService(page.PageNum, page.PageSize, userId)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, posts)
}

func DeletePost(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	postId := c.Param("id")
	
	result := db.DB.Where("id=? and user_id=?", getUintFromString(postId), userId).First(&models.PostModel{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		return
	}
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	err := service.DeletePostService(userId, getUintFromString(postId))
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

// 辅助函数，将字符串转换为uint
func getUintFromString(s string) uint {
	var n uint
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}
