package front_api

import (
	"github.com/gin-gonic/gin"
	"pet-project/db"
	"pet-project/models"
	"pet-project/response"
	"strconv"
)

func UserList(c *gin.Context) {
	// 获取用户列表
	var users []models.UserInfo
	pageNumber := c.Query("page")
	pageSize := 20
	number, _ := strconv.Atoi(pageNumber)
	offset := (number - 1) * pageSize
	var userList []models.UserInfo
	result := db.DB.Model(models.UserInfo{}).Offset(offset).Limit(pageSize).Find(&userList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, users)
}
