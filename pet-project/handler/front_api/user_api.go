package front_api

import (
	"github.com/gin-gonic/gin"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
)

func UserList(c *gin.Context) {
	// 获取用户列表
	var page models.PageModel
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (page.PageNum - 1) * page.PageSize
	var userList []models.UserInfo
	result := db.DB.Model(models.UserInfo{}).Offset(offset).Limit(page.PageSize).Find(&userList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, userList)
}

func UserSuggestionList(c *gin.Context) {
	var page models.PageModel
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	offset := (page.PageNum - 1) * page.PageSize
	var suggestions []models.SuggestionModel
	result := db.DB.Model(models.SuggestionModel{}).Offset(offset).Limit(page.PageSize).Find(&suggestions)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, suggestions)
}
