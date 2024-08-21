package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
	"net/http"
	"pet-project/db"
	"pet-project/middleware"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"
	"pet-project/settings"
	"pet-project/util"
	"strconv"
)

func GetTencentCode(c *gin.Context) {
	// 验证码
	lang := c.MustGet("lang").(*i18n.Localizer)
	email := c.PostForm("email")
	if len(email) == 0 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	smptServer := settings.Conf.EmailService.Host
	smptPort := settings.Conf.EmailService.Port
	username := settings.Conf.EmailService.Username
	password := settings.Conf.EmailService.Password
	code := service.GenerateValidationCode(6)
	// 对方的邮箱
	recipient := email
	subject := service.LocalizeMsg(lang, "VerificationTitle")
	body := service.LocalizeMsgCount(lang, "VerificationDesc", code)

	err := service.SendEmail(recipient, subject, body, smptServer, smptPort, username, password)
	if err != nil {
		response.Fail(c, response.ApiCode.ServerErr, err.Error())
		return
	}
	// 将code保存到redis，设置10分钟失效
	saveErr := service.SaveAccountCodeInRedis(c, email, code)
	if saveErr != nil {
		response.Fail(c, response.ApiCode.ServerErr, saveErr.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{},
	})
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
	lang := c.MustGet("lang").(*i18n.Localizer)
	var login models.RegisterInfo

	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	var findUser models.UserInfo
	var findResult *gorm.DB
	if len(login.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", login.Phone).First(&findUser)
	} else if len(login.Email) > 0 {
		if util.IsValidEmail(login.Email) {
			findResult = db.DB.Where("email = ?", login.Email).First(&findUser)
		} else {
			response.Fail(c, response.ApiCode.EmailErr, response.ApiMsg.EmailErr)
			return
		}
	} else {
		response.Fail(c, response.ApiCode.ParamLack, response.ApiMsg.ParamLack)
		return
	}

	// 如果查不到，则开始验证验证码
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		// 取出redis中的验证码
		if len(login.Email) > 0 {
			code, err := service.GetCodeFromRedis(c, login.Email)
			if err != nil {
				response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
				return
			}

			// 验证验证码是否正确
			codeValue, _ := strconv.Atoi(code)
			if codeValue != login.Code {
				response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
				return
			}
			_ = service.DeleteCodeFromRedis(c, login.Email)
		} else { // 验证手机验证码
			code, err := service.GetCodeFromRedis(c, login.Phone)
			if err != nil {
				response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
				return
			}

			// 验证验证码是否正确
			codeValue, _ := strconv.Atoi(code)
			if codeValue != login.Code {
				response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
				return
			}
			_ = service.DeleteCodeFromRedis(c, login.Email)
		}

		user := models.UserInfo{
			Phone:    login.Phone,
			Password: login.Password,
			Email:    login.Email,
		}
		result := db.DB.Create(&user)
		if result.Error != nil {
			response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
			return
		}
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
			return
		}
		data := models.LoginUserInfo{
			ID:     user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		response.Success(c, data)
	} else {
		msg := service.LocalizeMsg(lang, response.ApiMsg.UserExistsErr)
		response.Fail(c, response.ApiCode.UserExistsErr, msg)
	}

}

// UserPhoneLogin 用户登录
func UserPhoneLogin(c *gin.Context) {
	lang := c.MustGet("lang").(*i18n.Localizer)
	var login models.LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	var findResult *gorm.DB
	var user models.UserInfo
	if len(login.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", login.Phone).First(&user)
	} else if len(login.Email) > 0 {
		if util.IsValidEmail(login.Email) {
			findResult = db.DB.Where("email = ?", login.Email).First(&user)
		} else {
			response.Fail(c, response.ApiCode.EmailErr, response.ApiMsg.EmailErr)
			return
		}
	} else {
		response.Fail(c, response.ApiCode.ParamLack, response.ApiMsg.ParamLack)
		return
	}
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		msg := service.LocalizeMsg(lang, "AccountUnRegister")
		response.Fail(c, response.ApiCode.QueryErr, msg)
		return
	}
	if user.Password == login.Password {
		// 密码正确, 生成token，登录完成
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
			return
		}
		data := models.LoginUserInfo{
			ID:     user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		response.Success(c, data)
	} else {
		msg := service.LocalizeMsg(lang, "PasswordErr")
		response.Fail(c, 300, msg)
	}
}

// UserFindPassword MARK: 找回密码
func UserFindPassword(c *gin.Context) {
	lang := c.MustGet("lang").(*i18n.Localizer)
	var loginInfo models.RegisterInfo
	if err := c.ShouldBind(&loginInfo); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	var findResult *gorm.DB
	var user models.UserInfo
	if len(loginInfo.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", loginInfo.Phone).First(&user)
	} else if len(loginInfo.Email) > 0 {
		if util.IsValidEmail(loginInfo.Email) {
			findResult = db.DB.Where("email = ?", loginInfo.Email).First(&user)
		} else {
			response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
			return
		}
	}
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.UserNotFont, response.ApiMsg.UserNotFound)
		return
	} else {
		// 验证验证码
		if len(loginInfo.Phone) > 0 {
			code, err := service.GetCodeFromRedis(c, loginInfo.Email)
			codeErr := service.LocalizeMsg(lang, "CheckCodeErr")
			if err != nil {
				response.Fail(c, response.ApiCode.ParamErr, codeErr)
				return
			}
			codeValue, _ := strconv.Atoi(code)
			if codeValue != loginInfo.Code {
				response.Fail(c, response.ApiCode.ParamErr, codeErr)
				return
			}
			// 更新密码
			result := db.DB.Model(&user).Where("phone = ?", loginInfo.Phone).Update("password", loginInfo.Password)
			if result.Error != nil {
				response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
				return
			}
			// redis的数据清除
			_ = service.DeleteCodeFromRedis(c, loginInfo.Phone)

			response.Success(c, map[string]interface{}{})
		} else {
			code, err := service.GetCodeFromRedis(c, loginInfo.Email)
			codeErr := service.LocalizeMsg(lang, "CheckCodeErr")
			if err != nil {
				response.Fail(c, response.ApiCode.ParamErr, codeErr)
				return
			}
			codeValue, _ := strconv.Atoi(code)
			if codeValue != loginInfo.Code {
				response.Fail(c, response.ApiCode.ParamErr, codeErr)
				return
			}
			result := db.DB.Model(&user).Where("email = ?", loginInfo.Email).Update("password", loginInfo.Password)
			if result.Error != nil {
				response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
				return
			}

			// 删除redis数据
			_ = service.DeleteCodeFromRedis(c, loginInfo.Email)

			response.Success(c, map[string]interface{}{})
		}
	}
}
