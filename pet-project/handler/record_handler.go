package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"slices"
)

// PetInfoCreate 提交宠物详情
func PetInfoCreate(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var petInfo models.PetInfo
	if err := c.ShouldBind(&petInfo); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	fmt.Println(petInfo.User)
	// 忽略User是因为ShouldBind会创建一个User默认值，导致插入一条新的用户数据
	//result := db.DB.Omit("User").Create(&petInfo)
	//result := db.DB.Omit(clause.Associations).Create(&petInfo)
	// 现在不会创建关联对象了
	petInfo.UserId = userId
	result := db.DB.Omit(clause.Associations).Create(&petInfo)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// GetPetList 获取创建的宠物列表
func GetPetList(c *gin.Context) {
	var uerId = c.MustGet("userId").(uint)
	var petList []models.PetInfo
	var pageModel models.PageModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (pageModel.PageNum - 1) * pageModel.PageSize
	result := db.DB.
		Preload("User").
		Model(&models.PetInfo{}).
		Where("user_id = ?", uerId).
		Offset(offset).
		Limit(pageModel.PageSize).
		Order("created_at DESC").
		Find(&petList)

	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, petList)
}

// UpdatePetInfo 更新宠物信息
func UpdatePetInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var petInfo models.PetInfo
	if err := c.ShouldBind(&petInfo); err != nil {
		log.Println(err.Error())
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if petInfo.ID == 0 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	// 忽略User是因为ShouldBind会创建一个User默认值，导致插入一条新的用户数据
	petInfo.UserId = userId
	result := db.DB.Model(&models.PetInfo{}).Where("id = ? AND user_id = ?", petInfo.ID, petInfo.UserId).
		Updates(models.PetInfo{
			UserId:   petInfo.UserId,
			PetType:  petInfo.PetType,
			Avatar:   petInfo.Avatar,
			Name:     petInfo.Name,
			Gender:   petInfo.Gender,
			BirthDay: petInfo.BirthDay,
			HomeDay:  petInfo.HomeDay,
			Weight:   petInfo.Weight,
			Unit:     petInfo.Unit,
			Desc:     petInfo.Desc,
		})
	if result.Error != nil {
		log.Println(result.Error)
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	response.Success(c, nil)
}

// DeletePetInfo 删除创建的宠物详情
func DeletePetInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	petId := c.Param("id")
	petInfo := models.PetInfo{}
	findResult := db.DB.Model(&models.PetInfo{}).
		Where("id = ? AND user_id = ?", petId, userId).
		First(&petInfo)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		return
	}
	result := db.DB.Delete(&petInfo, "id = ?", petId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

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

func CreateCommonCategory(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	if userId == 1 {
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
	} else {
		response.Fail(c, response.ApiCode.RejectErr, response.ApiMsg.RejectErr)
	}
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
	userId := c.MustGet("userId").(uint)
	if userId == 1 {
		id := c.Param("id")
		result := db.DB.Delete(&models.RecordCategory{}, "id = ? AND user_id IS NULL", id)
		if result.Error != nil {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
			return
		}
		response.Success(c, nil)
	} else {
		response.Fail(c, response.ApiCode.RejectErr, response.ApiMsg.RejectErr)
	}

}

// GetRecordCategoryList 获取宠物行为列表
func GetRecordCategoryList(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var pageModel models.CategoryTypeModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (pageModel.PageNum - 1) * pageModel.PageSize
	var petActionList []models.RecordCategory
	if pageModel.CategoryType == 1 {
		result := db.DB.Model(&models.RecordCategory{}).
			Where("user_id = ?", userId).
			Offset(offset).
			Limit(pageModel.PageSize).
			Find(&petActionList)
		if result.Error != nil {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
			return
		}
		response.Success(c, petActionList)
	} else {
		result := db.DB.Model(&models.RecordCategory{}).
			Where("user_id IS NULL").Or("user_id = ?", userId).
			Offset(offset).
			Limit(pageModel.PageSize).
			Find(&petActionList)
		if result.Error != nil {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
			return
		}
		response.Success(c, petActionList)
	}

}

// CreateRecordCategory 添加宠物行为
func CreateRecordCategory(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var recordCategory models.RecordCategory
	if err := c.ShouldBind(&recordCategory); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	recordCategory.UserId = &userId
	result := db.DB.Create(&recordCategory)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func UpdateRecordCategory(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	// 解析 JSON
	var recordCategory models.RecordCategory
	if err := c.ShouldBindJSON(&recordCategory); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 确认该分类属于当前用户
	var old models.RecordCategory
	if err := db.DB.Where("id = ? AND user_id = ?", recordCategory.ID, userId).First(&old).Error; err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	// 执行更新（只更新传入的字段）
	if err := db.DB.Model(&old).Updates(recordCategory).Error; err != nil {
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}

	response.Success(c, nil)
}

func DeleteRecordCategory(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id := c.Param("id")
	// 先查询存不存在
	recordCategory := models.RecordCategory{}
	findResult := db.DB.Model(&models.RecordCategory{}).Where("id = ? AND user_id = ?", id, userId).First(&recordCategory)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		return
	}
	result := db.DB.Delete(&models.RecordCategory{}, "id = ? AND user_id = ?", id, userId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

// CreateRecord 创建记录
func CreateRecord(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var model models.RecordList
	if err := c.ShouldBind(&model); err != nil {
		fmt.Println(err)
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	// 查一下宠物是否存在
	var petInfo models.PetInfo
	petResult := db.DB.Model(&models.PetInfo{}).Where("id = ?", model.PetInfoId).First(&petInfo)
	if errors.Is(petResult.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	model.UserId = userId
	fmt.Println(model)
	result := db.DB.Omit(clause.Associations).Create(&model)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, model)
}

// GetRecordList 查询记录列表
func GetRecordList(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var pageModel models.RecordListModel
	if err := c.ShouldBind(&pageModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (pageModel.PageNum - 1) * pageModel.PageSize
	var recordList []models.RecordList

	// 查询的参数
	queryParam := models.RecordList{}
	if pageModel.CategoryId != nil {
		queryParam.RecordCategoryId = pageModel.CategoryId
	}
	if pageModel.PetInfoId != 0 {
		queryParam.PetInfoId = pageModel.PetInfoId
	}
	queryParam.UserId = userId

	result := db.DB.Preload("User").
		Model(&models.RecordList{}).
		Where(&queryParam).
		Offset(offset).
		Limit(pageModel.PageSize).
		Order("record_time DESC").
		Find(&recordList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, recordList)
}

// DeleteRecordInfo 删除记录
func DeleteRecordInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id := c.Param("id")
	result := db.DB.Where("id=? and user_id=?", id, userId).Delete(&models.RecordList{})
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, map[string]interface{}{})
}
