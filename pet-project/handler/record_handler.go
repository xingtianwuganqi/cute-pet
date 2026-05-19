package handler

import (
	"fmt"
	"log"
	"pet-project/internal"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"
	"pet-project/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PetInfoCreate 提交宠物详情
func PetInfoCreate(c *gin.Context) {
	logger.Logger.Info("PetInfoCreate handler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)
	var petInfo models.PetInfo
	if err := c.ShouldBind(&petInfo); err != nil {
		logger.Logger.Warn("Invalid parameters for PetInfoCreate", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.PetInfoCreateService(userId, petInfo)
	if err != nil {
		logger.Logger.Error("Failed to create pet info", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	logger.Logger.Info("Successfully created pet info", zap.Uint("userId", userId), zap.Uint("petId", petInfo.ID))
	response.Success(c, nil)
}

// GetPetList 获取创建的宠物列表
func GetPetList(c *gin.Context) {
	logger.Logger.Info("GetPetList handler called", zap.String("clientIP", c.ClientIP()))

	var userId = c.MustGet("userId").(uint)
	var pageModel models.PageModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		logger.Logger.Warn("Invalid parameters for GetPetList", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	petList, err := service.GetPetListService(userId, pageModel)
	if err != nil {
		logger.Logger.Error("Failed to get pet list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved pet list", zap.Int("count", len(petList)), zap.Uint("userId", userId))
	response.Success(c, petList)
}

// UpdatePetInfo 更新宠物信息
func UpdatePetInfo(c *gin.Context) {
	logger.Logger.Info("UpdatePetInfo handler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)
	var petInfo models.PetInfo
	if err := c.ShouldBind(&petInfo); err != nil {
		log.Println(err.Error())
		logger.Logger.Warn("Invalid parameters for UpdatePetInfo", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.UpdatePetInfoService(userId, petInfo)
	if err != nil {
		logger.Logger.Error("Failed to update pet info", zap.Error(err))
		if err.Error() == "param error" {
			response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		} else if err.Error() == "data not exit" {
			response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		} else if err.Error() == "update error" {
			response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		}
		return
	}
	logger.Logger.Info("Successfully updated pet info", zap.Uint("userId", userId), zap.Uint("petId", petInfo.ID))
	response.Success(c, nil)
}

// DeletePetInfo 删除创建的宠物详情
func DeletePetInfo(c *gin.Context) {
	logger.Logger.Info("DeletePetInfo handler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)
	petId := internal.GetUintFromString(c.Param("id"))

	err := service.DeletePetInfoService(userId, petId)
	if err != nil {
		logger.Logger.Error("Failed to delete pet info", zap.Error(err), zap.Uint("userId", userId), zap.Uint("petId", petId))
		if err.Error() == "data not exit" {
			response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		} else {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		}
		return
	}
	logger.Logger.Info("Successfully deleted pet info", zap.Uint("userId", userId), zap.Uint("petId", petId))
	response.Success(c, nil)
}

// GetRecordCategoryList 获取宠物行为列表
func GetRecordCategoryList(c *gin.Context) {
	logger.Logger.Info("GetRecordCategoryList handler called", zap.String("clientIP", c.ClientIP()))

	var userId = c.MustGet("userId").(uint)
	var pageModel models.CategoryTypeModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		logger.Logger.Warn("Invalid parameters for GetRecordCategoryList", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	petActionList, err := service.GetRecordCategoryListService(userId, pageModel)
	if err != nil {
		logger.Logger.Error("Failed to get record category list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved record category list", zap.Int("count", len(petActionList)), zap.Uint("userId", userId))
	response.Success(c, petActionList)
}

// CreateRecordCategory 添加宠物行为
func CreateRecordCategory(c *gin.Context) {
	logger.Logger.Info("CreateRecordCategory handler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)
	var recordCategory models.RecordCategory
	if err := c.ShouldBind(&recordCategory); err != nil {
		logger.Logger.Warn("Invalid parameters for CreateRecordCategory", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.CreateRecordCategoryService(userId, recordCategory)
	if err != nil {
		logger.Logger.Error("Failed to create record category", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	logger.Logger.Info("Successfully created record category", zap.Uint("userId", userId), zap.Uint("categoryId", recordCategory.ID))
	response.Success(c, nil)
}

func UpdateRecordCategory(c *gin.Context) {
	logger.Logger.Info("UpdateRecordCategory handler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)

	// 解析 JSON
	var recordCategory models.RecordCategory
	if err := c.ShouldBindJSON(&recordCategory); err != nil {
		logger.Logger.Warn("Invalid parameters for UpdateRecordCategory", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.UpdateRecordCategoryService(userId, recordCategory)
	if err != nil {
		logger.Logger.Error("Failed to update record category", zap.Error(err))
		if err.Error() == "query error" {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		} else if err.Error() == "update error" {
			response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		}
		return
	}

	logger.Logger.Info("Successfully updated record category", zap.Uint("userId", userId), zap.Uint("categoryId", recordCategory.ID))
	response.Success(c, nil)
}

func DeleteRecordCategory(c *gin.Context) {
	logger.Logger.Info("DeleteRecordCategory handler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)
	id := internal.GetUintFromString(c.Param("id"))

	err := service.DeleteRecordCategoryService(userId, id)
	if err != nil {
		logger.Logger.Error("Failed to delete record category", zap.Error(err), zap.Uint("userId", userId), zap.Uint("id", id))
		if err.Error() == "data not exit" {
			response.Fail(c, response.ApiCode.DataNotExit, response.ApiMsg.DataNotExit)
		} else {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		}
		return
	}
	logger.Logger.Info("Successfully deleted record category", zap.Uint("userId", userId), zap.Uint("id", id))
	response.Success(c, nil)
}

// CreateRecord 创建记录
func CreateRecord(c *gin.Context) {
	logger.Logger.Info("CreateRecord handler called", zap.String("clientIP", c.ClientIP()))

	var userId = c.MustGet("userId").(uint)
	var model models.RecordList
	if err := c.ShouldBind(&model); err != nil {
		fmt.Println(err)
		logger.Logger.Warn("Invalid parameters for CreateRecord", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.CreateRecordService(userId, model)
	if err != nil {
		logger.Logger.Error("Failed to create record", zap.Error(err))
		if err.Error() == "query error" {
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		} else {
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		}
		return
	}
	logger.Logger.Info("Successfully created record", zap.Uint("userId", userId), zap.Uint("recordId", model.ID))
	response.Success(c, model)
}

// GetRecordList 查询记录列表
func GetRecordList(c *gin.Context) {
	logger.Logger.Info("GetRecordList handler called", zap.String("clientIP", c.ClientIP()))

	var userId = c.MustGet("userId").(uint)
	var pageModel models.RecordListModel
	if err := c.ShouldBind(&pageModel); err != nil {
		logger.Logger.Warn("Invalid parameters for GetRecordList", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	recordList, err := service.GetRecordListService(userId, pageModel)
	if err != nil {
		logger.Logger.Error("Failed to get record list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	logger.Logger.Info("Successfully retrieved record list", zap.Int("count", len(recordList)), zap.Uint("userId", userId))
	response.Success(c, recordList)
}

// DeleteRecordInfo 删除记录
func DeleteRecordInfo(c *gin.Context) {
	logger.Logger.Info("DeleteRecordInfo handler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)
	id := internal.GetUintFromString(c.Param("id"))

	err := service.DeleteRecordInfoService(userId, id)
	if err != nil {
		logger.Logger.Error("Failed to delete record info", zap.Error(err), zap.Uint("userId", userId), zap.Uint("id", id))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully deleted record info", zap.Uint("userId", userId), zap.Uint("id", id))
	response.Success(c, map[string]interface{}{})
}

// 查询宠物花费列表
func GetPetCostList(c *gin.Context) {
	logger.Logger.Info("GetPetCostList handler called", zap.String("clientIP", c.ClientIP()))

	userId := c.MustGet("userId").(uint)

	var pageModel models.RecordListModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		logger.Logger.Warn("Invalid parameters for GetPetCostList", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	recordList, err := service.GetPetCostListService(userId, pageModel)
	if err != nil {
		logger.Logger.Error("Failed to get pet cost list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	logger.Logger.Info("Successfully retrieved pet cost list", zap.Int("count", len(recordList)), zap.Uint("userId", userId))
	response.Success(c, recordList)
}

// // 辅助函数，将字符串转换为uint
// func getUintFromString(s string) uint {
// 	var n uint
// 	fmt.Sscanf(s, "%d", &n)
// 	return n
// }