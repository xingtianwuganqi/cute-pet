package handler

import (
	"errors"
	"net/http"
	"pet-project/db"
	"pet-project/middleware"
	"pet-project/models"
	"pet-project/response"
	"pet-project/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTencentCode(c *gin.Context) {
	// 验证码
	code := c.PostForm("code")
	num, _ := strconv.Atoi(code)
	if num != 200 {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": 1008,
	})
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
	var login models.LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	var findUser models.UserInfo
	findResult := db.DB.Where("phone = ?", login.Phone).First(&findUser)
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		user := models.UserInfo{
			Phone:    login.Phone,
			Password: login.Password,
		}
		result := db.DB.Create(&user)
		if result.Error != nil {
			response.Fail(c, util.ApiCode.CreateErr, util.ApiMessage.CreateErr)
			return
		}
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
			return
		}
		data := models.LoginUserInfo{
			UserId: user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		response.Success(c, data)
	} else {
		response.Fail(c, util.ApiCode.UserExistsError, util.ApiMessage.UserExistsError)
	}

}

// UserPhoneLogin 用户登录
func UserPhoneLogin(c *gin.Context) {
	var login models.LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}
	var user models.UserInfo
	result := db.DB.Where("phone = ?", login.Phone).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, util.ApiCode.QueryError, "该手机号未注册")
		return
	}
	if user.Password == login.Password {
		// 密码正确, 生成token，登录完成
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
			return
		}
		data := models.LoginUserInfo{
			UserId: user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		response.Success(c, data)
	} else {
		response.Fail(c, 300, "密码错误")
	}
}
