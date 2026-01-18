package handler

import (
	"pet-project/internal"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"
	"pet-project/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetCommonCategories 获取宠物分类
func GetCommonCategories(c *gin.Context) {
	logger.Logger.Info("GetCommonCategories handler called", zap.String("clientIP", c.ClientIP()))

	petActionList, err := service.GetCommonCategoriesService()
	if err != nil {
		logger.Logger.Error("Failed to get common categories", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved common categories", zap.Int("count", len(petActionList)))
	response.Success(c, petActionList)
}

// CreateCommonCategory 创建宠物分类
func CreateCommonCategory(c *gin.Context) {
	logger.Logger.Info("CreateCommonCategory handler called", zap.String("clientIP", c.ClientIP()))

	var recordCategory models.RecordCategory
	if err := c.ShouldBind(&recordCategory); err != nil {
		logger.Logger.Warn("Invalid parameters for CreateCommonCategory", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	err := service.CreateCommonCategoryService(recordCategory)
	if err != nil {
		logger.Logger.Error("Failed to create common category", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	logger.Logger.Info("Successfully created common category", zap.Uint("id", recordCategory.ID))
	response.Success(c, nil)
}

// CreateCommonCategoryList 批量添加宠物行为
func CreateCommonCategoryList(c *gin.Context) {
	logger.Logger.Info("CreateCommonCategoryList handler called", zap.String("clientIP", c.ClientIP()))

	var categories []models.RecordCategory
	if err := c.ShouldBind(&categories); err != nil {
		logger.Logger.Warn("Invalid parameters for CreateCommonCategoryList", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	err := service.CreateCommonCategoryListService(categories)
	if err != nil {
		logger.Logger.Error("Failed to create common category list", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	logger.Logger.Info("Successfully created common category list", zap.Int("count", len(categories)))
	response.Success(c, nil)
}

// DeleteCommonCategory 删除宠物行为
func DeleteCommonCategory(c *gin.Context) {
	logger.Logger.Info("DeleteCommonCategory handler called", zap.String("clientIP", c.ClientIP()))

	id := internal.GetUintFromString(c.Param("id"))
	
	err := service.DeleteCommonCategoryService(id)
	if err != nil {
		logger.Logger.Error("Failed to delete common category", zap.Error(err), zap.Uint("id", id))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully deleted common category", zap.Uint("id", id))
	response.Success(c, nil)
}

/*
获取用户列表
*/
func GetUserList(c *gin.Context) {
	logger.Logger.Info("GetUserList handler called", zap.String("clientIP", c.ClientIP()))

	var page = models.PageModel{}
	if err := c.ShouldBind(&page); err != nil {
		logger.Logger.Warn("Invalid parameters for GetUserList", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	userList, err := service.GetUserListService(page)
	if err != nil {
		logger.Logger.Error("Failed to get user list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved user list", zap.Int("count", len(userList)))
	response.Success(c, userList)
}

/*
点赞列表
*/
func GetLikeList(c *gin.Context) {
	logger.Logger.Info("GetLikeList handler called", zap.String("clientIP", c.ClientIP()))

	var page = models.PageModel{}
	if err := c.ShouldBindQuery(&page); err != nil {
		logger.Logger.Warn("Invalid parameters for GetLikeList", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	likeList, err := service.GetLikeListService(page)
	if err != nil {
		logger.Logger.Error("Failed to get like list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved like list", zap.Int("count", len(likeList)))
	response.Success(c, likeList)
}

func GetCollectionList(c *gin.Context) {
	logger.Logger.Info("GetCollectionList handler called", zap.String("clientIP", c.ClientIP()))

	var page = models.PageModel{}
	if err := c.ShouldBindQuery(&page); err != nil {
		logger.Logger.Warn("Invalid parameters for GetCollectionList", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	collectionList, err := service.GetCollectionListService(page)
	if err != nil {
		logger.Logger.Error("Failed to get collection list", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	logger.Logger.Info("Successfully retrieved collection list", zap.Int("count", len(collectionList)))
	response.Success(c, collectionList)
}

// 辅助函数，将字符串转换为uint
// func getUintFromString(s string) uint {
// 	var n uint
// 	fmt.Sscanf(s, "%d", &n)
// 	return n
// }