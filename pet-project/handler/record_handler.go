package handler

import (
	"fmt"
	"log"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"

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

	err := service.PetInfoCreateService(userId, petInfo)
	if err != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// GetPetList 获取创建的宠物列表
func GetPetList(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var pageModel models.PageModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	petList, err := service.GetPetListService(userId, pageModel)
	if err != nil {
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

	err := service.UpdatePetInfoService(userId, petInfo)
	if err != nil {
		if err.Error() == "param error" {
			response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		} else if err.Error() == "data not exit" {
			response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		} else if err.Error() == "update error" {
			response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		}
		return
	}
	response.Success(c, nil)
}

// DeletePetInfo 删除创建的宠物详情
func DeletePetInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	petId := getUintFromString(c.Param("id"))

	err := service.DeletePetInfoService(userId, petId)
	if err != nil {
		if err.Error() == "data not exit" {
			response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		} else {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		}
		return
	}
	response.Success(c, nil)
}

// GetRecordCategoryList 获取宠物行为列表
func GetRecordCategoryList(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var pageModel models.CategoryTypeModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	petActionList, err := service.GetRecordCategoryListService(userId, pageModel)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, petActionList)
}

// CreateRecordCategory 添加宠物行为
func CreateRecordCategory(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var recordCategory models.RecordCategory
	if err := c.ShouldBind(&recordCategory); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.CreateRecordCategoryService(userId, recordCategory)
	if err != nil {
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

	err := service.UpdateRecordCategoryService(userId, recordCategory)
	if err != nil {
		if err.Error() == "query error" {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		} else if err.Error() == "update error" {
			response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		}
		return
	}

	response.Success(c, nil)
}

func DeleteRecordCategory(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id := getUintFromString(c.Param("id"))

	err := service.DeleteRecordCategoryService(userId, id)
	if err != nil {
		if err.Error() == "data not exit" {
			response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		} else {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		}
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

	err := service.CreateRecordService(userId, model)
	if err != nil {
		if err.Error() == "query error" {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		} else {
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		}
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

	recordList, err := service.GetRecordListService(userId, pageModel)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, recordList)
}

// DeleteRecordInfo 删除记录
func DeleteRecordInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	id := getUintFromString(c.Param("id"))

	err := service.DeleteRecordInfoService(userId, id)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, map[string]interface{}{})
}

// 查询宠物花费列表
func GetPetCostList(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var pageModel models.RecordListModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	recordList, err := service.GetPetCostListService(userId, pageModel)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, recordList)
}

// // 辅助函数，将字符串转换为uint
// func getUintFromString(s string) uint {
// 	var n uint
// 	fmt.Sscanf(s, "%d", &n)
// 	return n
// }
