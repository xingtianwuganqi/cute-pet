package handler

import (
	"github.com/gin-gonic/gin"
	"pet-project/db"
	"pet-project/models"
	"pet-project/util"
)

// PetInfoCreate 提交宠物详情
func PetInfoCreate(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var petInfo models.PetInfo
	if err := c.ShouldBind(&petInfo); err != nil {
		util.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}

	// 如果token的userId和参数的userId不一样，说明不是同一个人
	if petInfo.UserId != userId {
		util.Fail(c, util.ApiCode.QueryError, util.ApiMessage.QueryError)
		return
	}

	result := db.DB.Create(&petInfo)
	if result.Error != nil {
		util.Fail(c, util.ApiCode.CreateErr, util.ApiMessage.CreateErr)
		return
	}
	util.Success(c, nil)
}

// GetPetActionList 获取宠物行为
func GetPetActionList(c *gin.Context) {
	var petActionList []models.PetActionType
	result := db.DB.Model(models.PetActionType{}).Find(&petActionList)
	if result != nil {
		util.Fail(c, util.ApiCode.QueryError, util.ApiMessage.QueryError)
		return
	}
	util.Success(c, petActionList)
}
