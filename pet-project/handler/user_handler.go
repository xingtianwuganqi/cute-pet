package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"pet-project/db"
	"pet-project/models"
	"pet-project/util"
)

type LoginInfo struct {
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Password string `form:"password" json:"password"`
	Code     string `form:"code" json:"code"`
}

type LoginUserInfo struct {
	UserId uint   `json:"userId"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Token  string `json:"token"`
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
	var login LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		util.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	var findUser models.UserInfo
	findResult := db.DB.Where("phone = ?", login.Phone).First(&findUser)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		token := util.Md5String(login.Phone)
		user := models.UserInfo{
			Phone:    login.Phone,
			Password: login.Password,
			UserToken: models.UserToken{
				Token: token,
			},
		}
		result := db.DB.Create(&user)
		if result.Error != nil {
			util.Fail(c, util.ApiCode.CreateErr, util.ApiMessage.CreateErr)
			return
		}
		data := LoginUserInfo{
			UserId: user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		util.Success(c, data)
	} else {
		util.Fail(c, 400, "fail")
	}

}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var login LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		util.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	var user models.UserInfo
	result := db.DB.Where("phone = ?", login.Phone).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		util.Fail(c, util.ApiCode.QueryError, "该手机号未注册")
		return
	}
	if user.Password == login.Password {
		// 密码正确, 生成token，登录完成
		token := util.Md5String(user.Phone)
		var tokenInfo models.UserToken
		tokenResult := db.DB.Where("id = ?", user.ID).First(&tokenInfo)
		if tokenResult.Error != nil {
			c.JSON(http.StatusGone, gin.H{
				"code": 400,
				"msg":  "查询失败",
			})
			return
		}
		tokenInfo.Token = token
		db.DB.Save(&tokenInfo)
		data := LoginUserInfo{
			UserId: user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		util.Success(c, data)
	} else {
		util.Fail(c, 300, "密码错误")
	}
}
