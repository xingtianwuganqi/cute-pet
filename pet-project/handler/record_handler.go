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
		Offset(offset).Limit(pageModel.PageSize).
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
		response.Fail(c, response.ApiCode.UploadErr, response.ApiMsg.UploadErr)
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

// CreateRecordCategory 添加宠物行为
func CreateRecordCategory(c *gin.Context) {
	var recordCategory models.RecordCategory
	if err := c.ShouldBind(&recordCategory); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result := db.DB.Create(&recordCategory)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// GetRecordCategoryList 获取宠物行为
func GetRecordCategoryList(c *gin.Context) {
	var petActionList []models.RecordCategory
	result := db.DB.Model(&models.RecordCategory{}).Find(&petActionList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, petActionList)
}

func DeleteRecordCategory(c *gin.Context) {
	id := c.Param("id")
	result := db.DB.Delete(&models.RecordCategory{}, "id = ?", id)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

// CreateCustomCategory 创建自定义分类
func CreateCustomCategory(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var customCategory models.CustomCategory
	if err := c.ShouldBind(&customCategory); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	customCategory.UserId = userId
	result := db.DB.Omit(clause.Associations).Create(&customCategory)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func GetCustomCategoryList(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	param := models.PageModel{}
	if err := c.ShouldBindQuery(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (param.PageNum - 1) * param.PageSize

	var customActionList []models.CustomCategory
	// 如果只返回特定字段，
	//var customActionList []models.PetCustomTypeInfo ,
	//查询db.DB.Model(&models.PetCustomType{}).Where("user_id = ?", userId).Find(&customActionList)
	// 定义为UserId的字段，GORM 自动将结构体字段名称转换为 user_id 作为数据库中的列名。
	//result := db.DB.Preload("User").Where("user_id = ?", userId).Find(&customActionList)
	// Select 或 Omit的字段，不会消失，会显示零值
	result := db.DB.Preload("User").Model(&models.CustomCategory{}).
		Offset(offset).Limit(param.PageSize).
		Order("created_at DESC").
		Find(&customActionList, "user_id = ?", userId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, customActionList)
}

func DeleteCustomCategory(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	petCustomActionId := c.Param("id")
	result := db.DB.Delete(&models.CustomCategory{}, "id = ? AND user_id = ?", petCustomActionId, userId)
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
	var pageModel models.PageModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (pageModel.PageNum - 1) * pageModel.PageSize
	var recordList []models.RecordList
	result := db.DB.Preload("User").
		Model(&models.RecordList{}).Where("user_id = ?", userId).
		Offset(offset).
		Limit(pageModel.PageSize).
		Order("created_at DESC").
		Find(&recordList)
	if result.Error != nil {
		fmt.Println(result.Error)
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
