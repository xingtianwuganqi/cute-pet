package admin

import (
	"github.com/gin-gonic/gin"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
)

func UserList(c *gin.Context) {
	var pageModel models.PageModel
	if err := c.ShouldBindQuery(&pageModel); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.AMsg.ParamErr)
		return
	}
	pageNum := pageModel.PageNum
	pageSize := pageModel.PageSize
	offset := (pageNum - 1) * pageSize

	/*
		Limit 指定每页返回的记录数，而 Offset 指定跳过的记录数。
	*/
	var userList []models.LoginUserInfo
	result := db.DB.Model(&models.UserInfo{}).Limit(pageNum).Offset(offset).
		Omit("token").Order("created_at asc").Find(&userList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.AMsg.QueryErr)
		return
	}
	response.Success(c, userList)
}
