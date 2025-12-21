package handler

import (
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"slices"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// GetCommonCategories 获取宠物分类
func GetCommonCategories(c *gin.Context) {
	var petActionList []models.RecordCategory
	result := db.DB.Model(&models.RecordCategory{}).
		Where("user_id IS NULL").
		Find(&petActionList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, petActionList)
}

// CreateCommonCategory 创建宠物分类
func CreateCommonCategory(c *gin.Context) {
	var recordCategory models.RecordCategory
	if err := c.ShouldBind(&recordCategory); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	recordCategory.UserId = nil
	result := db.DB.Create(&recordCategory)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// CreateCommonCategoryList 批量添加宠物行为
func CreateCommonCategoryList(c *gin.Context) {
	var categories []models.RecordCategory
	if err := c.ShouldBind(&categories); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	var values []bool
	for i := range categories {
		var item = categories[i]
		item.UserId = nil
		result := db.DB.Omit(clause.Associations).Create(&categories[i])
		if result.Error != nil {
			values = append(values, false)
		} else {
			values = append(values, true)
		}
	}
	if slices.Contains(values, false) {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// DeleteCommonCategory 删除宠物行为
func DeleteCommonCategory(c *gin.Context) {
	id := c.Param("id")
	result := db.DB.Delete(&models.RecordCategory{}, "id = ? AND user_id IS NULL", id)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

func GetUserList(c *gin.Context) {
	var userList []models.UserInfo
	result := db.DB.Find(&userList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, userList)
}
