package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"pet-project/db"
	"pet-project/middleware"
	"pet-project/models"
	"pet-project/response"
	"pet-project/settings"
	"pet-project/util"
	"strconv"
	"time"
)

func GetTencentCode(c *gin.Context) {
	// 验证码
	email := c.PostForm("email")
	//language := c.PostForm("lan")
	if len(email) == 0 {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}

	smptServer := settings.Conf.EmailService.Host
	smptPort := settings.Conf.EmailService.Port
	username := settings.Conf.EmailService.Username
	password := settings.Conf.EmailService.Password
	code := generateValidationCode(6)
	// 对方的邮箱
	recipient := email
	subject := "【您的验证码】"
	body := fmt.Sprintf("您的验证码为 %s,请在10分钟内使用。", code)

	err := sendEmail(recipient, subject, body, smptServer, smptPort, username, password)
	if err != nil {
		response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
		return
	}
	// 将code保存到redis，设置10分钟失效
	saveErr := saveAccountCodeInRedis(c, email, code)
	if saveErr != nil {
		response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{},
	})
}

// 保存到redis
func saveAccountCodeInRedis(c *gin.Context, email string, code string) error {
	expiration := 10 * time.Minute
	err := db.Rdb.Set(c, email, code, expiration).Err()
	return err
}

// redis中取出code值
func getCodeFromRedis(c *gin.Context, email string) (string, error) {
	value, err := db.Rdb.Get(c, email).Result()
	return value, err
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
	locale := c.MustGet("locale").(*i18n.Localizer)
	userExist := locale.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "UserNotFund",
	})
	var login models.LoginInfo

	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
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
			response.Fail(c, util.ApiCode.EmailError, util.ApiMessage.EmailError)
			return
		}
	} else {
		response.Fail(c, util.ApiCode.ParamLack, util.ApiMessage.ParamLack)
		return
	}

	// 如果查不到，则开始验证验证码
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		// 取出redis中的验证码
		if len(login.Email) > 0 {
			code, err := getCodeFromRedis(c, login.Email)
			if err != nil {
				response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
				return
			}

			// 验证验证码是否正确
			codeValue, _ := strconv.Atoi(code)
			if codeValue != login.Code {
				response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
				return
			}
		} else { // 验证手机验证码
			code, err := getCodeFromRedis(c, login.Phone)
			if err != nil {
				response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
				return
			}

			// 验证验证码是否正确
			codeValue, _ := strconv.Atoi(code)
			if codeValue != login.Code {
				response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
				return
			}
		}

		user := models.UserInfo{
			Phone:    login.Phone,
			Password: login.Password,
			Email:    login.Email,
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
		response.Fail(c, util.ApiCode.UserExistsError, userExist)
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

// 生成指定长度的随机数字验证码
func generateValidationCode(length int) string {
	var code string
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(9))
	}
	return code
}

// 发送验证码
func sendEmail(recipient string,
	subject string,
	body string,
	smtpServier string,
	smtpPort int,
	username string,
	password string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpServier, smtpPort, username, password)
	d.SSL = true //使用ssl连接

	return d.DialAndSend(m)
}
