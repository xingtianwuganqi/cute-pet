package handler

import (
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

	if petInfo.UserId != userId {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	result := db.DB.Create(&petInfo)
	if result.Error != nil {
		log.Println(result.Error)
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
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
	if result != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, petActionList)
}

func CreatePetCustomType(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var petCustomType models.PetCustomType
	if err := c.ShouldBind(&petCustomType); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
	}
	log.Println("userId is", userId)
	petCustomType.UserId = userId
	result := db.DB.Create(&petCustomType)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func GetCustomActionList(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var customActionList []models.PetCustomTypeInfo
	// 如果只返回特定字段，
	//var customActionList []models.PetCustomTypeInfo ,
	//查询db.DB.Model(&models.PetCustomType{}).Where("user_id = ?", userId).Find(&customActionList)
	// 定义为UserId的字段，GORM 自动将结构体字段名称转换为 user_id 作为数据库中的列名。
	//result := db.DB.Preload("User").Where("user_id = ?", userId).Find(&customActionList)
	// Select 或 Omit的字段，不会消失，会显示零值
	result := db.DB.Model(&models.PetCustomType{}).Select("id, created_at, updated_at, deleted_at, custom_name, custom_icon").
		Order("created_at DESC").
		Find(&customActionList, "user_id = ?", userId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, customActionList)
}

// GetRecordList 查询记录列表
func GetRecordList(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var pageNum = c.PostForm("pageNum")
	var pageSize = c.PostForm("pageSize")
	num, err := strconv.Atoi(pageNum)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	size, sizeErr := strconv.Atoi(pageSize)
	if sizeErr != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (num - 1) * size
	var recordList []models.RecordList
	result := db.DB.Where("user_id=?", userId).Offset(offset).Limit(size).Find(&recordList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, recordList)
}

// GetPetConsumeList 获取花销列表
func GetPetConsumeList(c *gin.Context) {
	var consumeModel []models.PetConsumeType
	result := db.DB.Model(&models.PetConsumeType{}).Find(&consumeModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, consumeModel)
}

// GetPetCustomConsumeList 获取用户自己创建的花销列表
func GetPetCustomConsumeList(c *gin.Context) {
	var userid = c.MustGet("userId").(uint)
	pageNum := c.PostForm("pageNum")
	pageSize := c.PostForm("pageSize")
	num, _ := strconv.Atoi(pageNum)
	size, _ := strconv.Atoi(pageSize)
	offset := (num - 1) * size
	var customTypes []models.PetCustomConsumeType
	findResult := db.DB.Model(&models.PetCustomConsumeType{}).Where("user_id = ?", userid).Offset(offset).Limit(size).Find(&customTypes)
	if findResult.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, customTypes)
}

func CreateConsumeAction(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var model models.PetCustomConsumeType
	if err := c.ShouldBind(&model); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	model.UserId = userId
	// 如果model包含主键，则更新（update）所有字段，如果不包含主键，则create
	result := db.DB.Save(&model)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, models.EmptyModel{})

}
