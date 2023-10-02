package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pet-project/db"
	"pet-project/models"
)

type LoginInfo struct {
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginUserInfo struct {
	user  models.UserInfo
	token string
}

// UserRegister 注册
func UserRegister(c *gin.Context) {

}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var login LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err,
		})
		return
	}
	var user models.UserInfo
	result := db.DB.Where("phone = ?", login.Phone).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusGone, gin.H{
			"code": 400,
			"msg":  "not find",
		})
		return
	}
	if user.Password == login.Password {
		// 密码正确, 生成token，登录完成
		token := "300200"
		tokenResult := db.DB.Where("userId = ?", user.ID).Update("token", token)
		if tokenResult.Error != nil {
			c.JSON(http.StatusGone, gin.H{
				"code": 400,
				"msg":  "查询失败",
			})
			return
		}
		info := LoginUserInfo{
			user:  user,
			token: token,
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": info,
			"msg":  "success",
		})
	} else {
		c.JSON(http.StatusGone, gin.H{
			"code": 400,
			"msg":  "密码错误",
		})
		return
	}
}
