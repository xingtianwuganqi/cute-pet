package handler

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"strconv"

	"github.com/gin-gonic/gin"
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
	result := db.DB.Create(&petInfo)
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

// CreatePetActionType 添加宠物行为
func CreatePetActionType(c *gin.Context) {
	var actionModel models.PetActionType
	if err := c.ShouldBind(&actionModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result := db.DB.Create(&actionModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// GetPetActionList 获取宠物行为
func GetPetActionList(c *gin.Context) {
	var petActionList []models.PetActionType
	result := db.DB.Model(&models.PetActionType{}).Find(&petActionList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, petActionList)
}

func DeletePetAction(c *gin.Context) {
	id := c.Param("id")
	result := db.DB.Delete(&models.PetActionType{}, "id = ?", id)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

// CreatePetCustomAction 创建自定义日常
func CreatePetCustomAction(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var petCustomAction models.PetCustomAction
	if err := c.ShouldBind(&petCustomAction); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	petCustomAction.UserId = userId
	result := db.DB.Omit("User").Create(&petCustomAction)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func GetCustomActionList(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	param := models.PageModel{}
	if err := c.ShouldBindQuery(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (param.PageNum - 1) * param.PageSize

	var customActionList []models.PetCustomAction
	// 如果只返回特定字段，
	//var customActionList []models.PetCustomTypeInfo ,
	//查询db.DB.Model(&models.PetCustomType{}).Where("user_id = ?", userId).Find(&customActionList)
	// 定义为UserId的字段，GORM 自动将结构体字段名称转换为 user_id 作为数据库中的列名。
	//result := db.DB.Preload("User").Where("user_id = ?", userId).Find(&customActionList)
	// Select 或 Omit的字段，不会消失，会显示零值
	result := db.DB.Preload("User").Model(&models.PetCustomAction{}).
		Offset(offset).Limit(param.PageSize).
		Order("created_at DESC").
		Find(&customActionList, "user_id = ?", userId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, customActionList)
}

func DeletePetCustomAction(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	petCustomActionId := c.Param("id")
	result := db.DB.Delete(&models.PetCustomAction{}, "id = ? AND user_id = ?", petCustomActionId, userId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

// CreateConsumeAction 创建公共花销
func CreateConsumeAction(c *gin.Context) {
	var consumeModel models.PetConsumeType
	if err := c.ShouldBind(&consumeModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	result := db.DB.Create(&consumeModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, consumeModel)
}

// GetPetConsumeList 获取花销列表
func GetPetConsumeList(c *gin.Context) {
	var consumeModels []models.PetConsumeType
	findResult := db.DB.Model(&models.PetConsumeType{}).Find(&consumeModels)
	if findResult.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, consumeModels)
}

func DeletePetConsumeAction(c *gin.Context) {
	id := c.Param("id")
	result := db.DB.Delete(&models.PetConsumeType{}, "id = ?", id)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

// CreateCustomConsumeAction 创建用户自定义花销
func CreateCustomConsumeAction(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var model models.PetCustomConsume
	if err := c.ShouldBind(&model); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	model.UserId = userId
	// 如果model包含主键，则更新（update）所有字段，如果不包含主键，则create
	result := db.DB.Create(&model)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, models.EmptyModel{})

}

// GetPetCustomConsumeList 获取用户自己创建的花销列表
func GetPetCustomConsumeList(c *gin.Context) {
	var userid = c.MustGet("userId").(uint)
	pageNum := c.Query("pageNum")
	pageSize := c.Query("pageSize")
	num, _ := strconv.Atoi(pageNum)
	size, _ := strconv.Atoi(pageSize)
	offset := (num - 1) * size
	var customTypes []models.PetCustomConsume
	findResult := db.DB.Model(&models.PetCustomConsume{}).Where("user_id = ?", userid).Offset(offset).Limit(size).Find(&customTypes)
	if findResult.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, customTypes)
}

// DeleteCustomConsumeAction 删除用户自己创建的花销
func DeleteCustomConsumeAction(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id := c.Param("id")
	result := db.DB.Delete(&models.PetCustomConsume{}, "id = ? AND user_id = ?", id, userId)
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
	model.UserId = userId
	result := db.DB.Create(&model)
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
