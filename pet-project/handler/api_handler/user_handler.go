package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"pet-project/db"
	"pet-project/middleware"
	"pet-project/models"
	"pet-project/response"
	"pet-project/util"
)

func GetTencentCode(c *gin.Context) {
	// 验证码
	email := c.PostForm("email")
	//language := c.PostForm("lan")
	if len(email) == 0 {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}

	smptServer := "smtp.163.com"
	smptPort := 465
	username := "xingtianwuganqi123@163.com"
	password := "DQHEPDTSYPAVKEFB"
	code := generateValidationCode(6)
	// 对方的邮箱
	recipient := email
	subject := "【您的验证码】"
	body := fmt.Sprintf("您的验证码为 %s,请在10分钟内使用。", code)

	// 将code保存到redis，设置10分钟失效

	err := sendEmail(recipient, subject, body, smptServer, smptPort, username, password)
	if err != nil {
		response.Fail(c, util.ApiCode.ServerError, util.ApiMessage.ServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{},
	})
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
	var login models.LoginInfo

	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, util.ApiCode.ParamError, util.ApiMessage.ParamError)
		return
	}

	// 验证验证码
	if login.Code != 2024 {
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

	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
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
		println("用户已存在", util.ApiCode.UserExistsError)
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

// 获取验证码
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
