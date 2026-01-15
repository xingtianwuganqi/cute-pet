package handler

import (
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"

	"github.com/gin-gonic/gin"
)

// GetCommonCategories 获取宠物分类
func GetCommonCategories(c *gin.Context) {
	petActionList, err := service.GetCommonCategoriesService()
	if err != nil {
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
	
	err := service.CreateCommonCategoryService(recordCategory)
	if err != nil {
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
	
	err := service.CreateCommonCategoryListService(categories)
	if err != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

// DeleteCommonCategory 删除宠物行为
func DeleteCommonCategory(c *gin.Context) {
	id := getUintFromString(c.Param("id"))
	
	err := service.DeleteCommonCategoryService(id)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, nil)
}

/*
获取用户列表
*/
func GetUserList(c *gin.Context) {
	var page = models.PageModel{}
	if err := c.ShouldBind(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	userList, err := service.GetUserListService(page)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, userList)
}

/*
点赞列表
*/
func GetLikeList(c *gin.Context) {
	var page = models.PageModel{}
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	
	likeList, err := service.GetLikeListService(page)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, likeList)
}

func GetCollectionList(c *gin.Context) {
	var page = models.PageModel{}
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	collectionList, err := service.GetCollectionListService(page)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, collectionList)
}

// 辅助函数，将字符串转换为uint
// func getUintFromString(s string) uint {
// 	var n uint
// 	fmt.Sscanf(s, "%d", &n)
// 	return n
// }