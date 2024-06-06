package handler

import (
	"log"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PetInfoCreate 提交宠物详情
func PetInfoCreate(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var petInfo models.PetInfo
	if err := c.ShouldBind(&petInfo); err != nil {
		log.Println(err)
		response.Fail(c, response.ApiCode.ParamErr, response.AMsg.ParamErr)
		return
	}

	// 如果token的userId和参数的userId不一样，说明不是同一个人
	log.Println("userId is", userId)
	log.Println("+++++++", reflect.TypeOf(petInfo.UserId), petInfo.UserId, petInfo.PetType)
	log.Println("0000", c.PostForm("pet_type"))
	if petInfo.UserId != userId {
		response.Fail(c, response.ApiCode.QueryErr, response.AMsg.QueryErr)
		return
	}

	result := db.DB.Create(&petInfo)
	if result.Error != nil {
		log.Println(result.Error)
		response.Fail(c, response.ApiCode.CreateErr, response.AMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// CreatePetActionType 添加宠物行为
func CreatePetActionType(c *gin.Context) {
	var actionModel models.PetActionType
	if err := c.ShouldBind(&actionModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.AMsg.ParamErr)
		return
	}

	result := db.DB.Create(&actionModel)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.AMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// GetPetActionList 获取宠物行为
func GetPetActionList(c *gin.Context) {
	var petActionList []models.PetActionType
	result := db.DB.Model(&models.PetActionType{}).Find(&petActionList)
	if result != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.AMsg.QueryErr)
		return
	}
	response.Success(c, petActionList)
}

func CreatePetCustomType(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var petCustomType models.PetCustomType
	if err := c.ShouldBind(&petCustomType); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.AMsg.ParamErr)
	}
	log.Println("userId is", userId)
	petCustomType.UserId = userId
	result := db.DB.Create(&petCustomType)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.AMsg.CreateErr)
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
		response.Fail(c, response.ApiCode.QueryErr, response.AMsg.QueryErr)
		return
	}
	response.Success(c, customActionList)
}

// GetRecordList 查询记录列表
func GetRecordList(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var pageNum = c.PostForm("pageNum")
	var pageSize = c.Query("pageSize")
	num, err := strconv.Atoi(pageNum)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.AMsg.ParamErr)
		return
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.AMsg.ParamErr)
		return
	}
	offset := (num - 1) * size
	var recordList []models.RecordList
	result := db.DB.Where("userId=?", userId).Offset(offset).Limit(size).Find(&recordList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.AMsg.QueryErr)
		return
	}
	response.Success(c, recordList)
}
