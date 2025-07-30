package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"pet-project/util"
)

// Front api

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

// GetStatusTopicList 获取待审核话题列表
func GetStatusTopicList(c *gin.Context) {
	status := c.Param("status")
	var page models.PageModel
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	var topicModel []models.TopicModel
	offset := (page.PageNum - 1) * page.PageSize
	result := db.DB.Model(models.TopicModel{}).Preload("User").
		Where("topic_status=?", status).Offset(offset).Limit(page.PageSize).Find(&topicModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, topicModel)
}

func ChangeTopicStatus(c *gin.Context) {
	var topicModel models.TopicStatusModel
	if err := c.ShouldBind(&topicModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	result := db.DB.Model(models.TopicModel{}).
		Preload("User").
		Where("id=?", topicModel.TopicId).
		Update("topic_status", topicModel.Status)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.UploadErr, response.ApiMsg.UploadErr)
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
	var topicModel []models.TopicModel
	offset := (page.PageNum - 1) * page.PageSize
	result := db.DB.Model(models.TopicModel{}).
		Where("topic_status=?", 1).
		Preload("User").
		Offset(offset).Limit(page.PageSize).Find(&topicModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, topicModel)

}

func UserCreateTopic(c *gin.Context) {
	userId, _ := c.Get("userId")
	var topicModel models.TopicModel
	if err := c.ShouldBind(&topicModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	value, _ := userId.(int)
	topicModel.UserId = uint(value)
	// 过滤敏感词
	filter := util.NewWordFilter()
	newTitle := filter.Replace(topicModel.Title)
	newContent := filter.Replace(topicModel.Desc)
	topicModel.Title = newTitle
	topicModel.Desc = newContent
	result := db.DB.Omit(clause.Associations).Create(&topicModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func DeleteUserTopic(c *gin.Context) {
	userId, _ := c.Get("userId")
	topicId := c.Param("id")
	topicModel := models.TopicModel{}
	result := db.DB.Where("id=? and user_id=?", topicId, userId).First(&topicModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		return
	}
	result = db.DB.Delete(&topicModel, "id = ?", topicId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
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
	// 过滤敏感词
	filter := util.NewWordFilter()
	newContent := filter.Replace(postModel.Content)
	postModel.UserId = userId
	postModel.Content = newContent
	result := db.DB.Omit("User").Create(&postModel)
	if result.Error != nil {
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
	var postModels []models.PostModel
	offset := (page.PageNum - 1) * page.PageSize
	result := db.DB.Model(models.PostModel{}).
		Preload("User").
		Offset(offset).Limit(page.PageSize).
		Order("created_at desc").
		Find(&postModels)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	// 获取 userId
	userIdInterface, exists := c.Get("userId")
	var userId uint
	if exists {
		userId = userIdInterface.(uint)
	}

	// 查询状态
	if userId != 0 {
		var postIds []uint
		for _, post := range postModels {
			postIds = append(postIds, post.ID)
		}

		var likedPosts []models.LikeMessageModel
		db.DB.Where("user_id = ? AND like_id IN ?", userId, postIds).Find(&likedPosts)
		var collectedPosts []models.CollectionMessageModel
		db.DB.Where("user_id = ? AND collection_id IN ?", userId, postIds).Find(&collectedPosts)

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

	response.Success(c, postModels)
}

func DeletePost(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	postId := c.Param("id")
	result := db.DB.Where("id=? and user_id=?", postId, userId).Delete(&models.PostModel{})
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}
