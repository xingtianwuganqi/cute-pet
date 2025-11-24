package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"io"
	"net/http"
	"pet-project/db"
	"pet-project/middleware"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"
	"pet-project/settings"
	"pet-project/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

func GetEmailCode(c *gin.Context) {
	// 验证码
	lang := c.MustGet("lang").(*i18n.Localizer)
	var param models.SendCodeModel
	paramErr := c.ShouldBind(&param)
	if paramErr != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 查询code是否在redis中（是否已经使用过了）
	codeKey := fmt.Sprintf("email_code:%s", param.Email)
	value, err := service.GetCodeFromRedis(c, codeKey)
	fmt.Println("value is ", value)
	fmt.Println("err is ", err)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if len(value) != 0 && value == param.Code {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	fmt.Println("value is ", value)
	fmt.Println("err is ", err)
	// 需要加一个加密信息
	encryptionStr, err := util.Decrypt(param.Code)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	// 判断encryptionStr是否今日日期
	if len(encryptionStr) == 0 || encryptionStr != GetTodayDate() {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	if len(param.Email) == 0 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	email := param.Email
	code := service.GenerateValidationCode(4)

	// 正式环境发验证码
	if settings.Conf.App.Env == "production" {
		smptServer := settings.Conf.EmailService.Host
		smptPort := settings.Conf.EmailService.Port
		username := settings.Conf.EmailService.Username
		password := settings.Conf.EmailService.Password
		// 对方的邮箱
		recipient := email
		subject := service.LocalizeMsg(lang, "VerificationTitle")
		body := service.LocalizeMsgCount(lang, "VerificationDesc", code)

		sendErr := service.SendEmail(recipient, subject, body, smptServer, smptPort, username, password)
		if sendErr != nil {
			response.Fail(c, response.ApiCode.ServerErr, sendErr.Error())
			return
		}
	}
	// 将code保存到redis，设置10分钟失效
	saveErr := service.SaveAccountCodeInRedis(c, email, code, 10*time.Minute)
	if saveErr != nil {
		response.Fail(c, response.ApiCode.ServerErr, saveErr.Error())
		return
	}

	// 保存param.code
	_ = service.SaveAccountCodeInRedis(c, codeKey, param.Code, 24*time.Hour)

	if settings.Conf.App.Env == "production" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": code,
		})
	}

}

// GetPhoneCode 获取手机验证码
func GetPhoneCode(c *gin.Context) {
	// 手机验证码
	var param models.SendCodeModel
	if err := c.ShouldBind(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	// 查询code是否在redis中（是否已经使用过了）
	codeKey := fmt.Sprintf("phone_code:%s", param.Phone)
	value, err := service.GetCodeFromRedis(c, codeKey)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if len(value) != 0 && value == param.Code {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	// 需要加一个加密信息
	encryptionStr, err := util.Decrypt(param.Code)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	// 判断encryptionStr是否今日日期
	if len(encryptionStr) == 0 || encryptionStr != GetTodayDate() {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	phone := param.Phone
	if len(phone) == 0 {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	code := service.GenerateValidationCode(4)

	if settings.Conf.App.Env == "production" {
		url := fmt.Sprintf("https://push.spug.cc/send/gL1QGmWdKWjlRD65?key1=%s&key2=%s&key3=%s&targets=%s",
			"[Pawpal]", code, "10", phone)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			response.Fail(c, response.ApiCode.Fail, response.ApiMsg.Fail)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
				return
			}
		}(resp.Body)
	}

	// 将code保存到redis，设置10分钟失效
	saveErr := service.SaveAccountCodeInRedis(c, phone, code, 10*time.Minute)
	if saveErr != nil {
		response.Fail(c, response.ApiCode.ServerErr, saveErr.Error())
		return
	}

	_ = service.SaveAccountCodeInRedis(c, codeKey, param.Code, 24*time.Hour)

	if settings.Conf.App.Env == "production" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": code,
		})
	}
}

func GetTodayDate() string {
	now := time.Now()
	date := now.Format("2006-01-02")
	return date
}

// GetEncryptionCode 获取今天的加密密钥
// 该函数没有输入参数，但会使用当前请求的上下文 *gin.Context
// 它首先调用 GetTodayDate() 获取今天的日期，然后使用 util.Encrypt() 对日期进行加密
// 如果加密过程中发生错误，它会发送一个失败的 HTTP 响应并返回
// 如果成功，它将返回一个包含加密密钥的 JSON 响应
func GetEncryptionCode(c *gin.Context) {
	// 调用 Encrypt 函数对今天的日期进行加密
	encryptionCode, err := util.Encrypt(GetTodayDate())
	if err != nil {
		// 如果加密过程中出现错误，发送失败的 HTTP 响应
		response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
		return
	}
	// 发送包含加密密钥的 JSON 响应
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": encryptionCode,
	})
}

// CheckRdbCode 校验验证码
func CheckRdbCode(c *gin.Context) {
	var param models.SendCodeModel
	if err := c.ShouldBind(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if len(param.Phone) != 0 {
		code, err := service.GetCodeFromRedis(c, param.Phone)
		if err != nil {
			// Redis 查询确实出错了，非 redis.Nil
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
			return
		}
		if code == param.Code {
			response.Success(c, gin.H{})
		} else {
			response.Fail(c, response.ApiCode.CheckCodeErr, response.ApiMsg.CheckCodeErr)
		}

	} else {
		code, err := service.GetCodeFromRedis(c, param.Email)
		if err != nil {
			// Redis 查询确实出错了，非 redis.Nil
			response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
			return
		}

		if code == param.Code {
			response.Success(c, gin.H{})
		} else {
			response.Fail(c, response.ApiCode.CheckCodeErr, response.ApiMsg.CheckCodeErr)
		}
	}
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
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
			if code != login.Code {
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
			if code != login.Code {
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
		response.Fail(c, response.ApiCode.UserExistsErr, response.ApiMsg.UserExistsErr)
	}

}

// UserPhoneLogin 用户登录
func UserPhoneLogin(c *gin.Context) {
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
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
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
		response.Fail(c, response.ApiCode.PasswordErr, response.ApiMsg.PasswordErr)
	}
}

// UserFindPassword MARK: 找回密码
func UserFindPassword(c *gin.Context) {
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
	} else {
		response.Fail(c, response.ApiCode.ParamLack, response.ApiMsg.ParamLack)
		return
	}
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	} else {
		// 验证验证码
		if len(loginInfo.Phone) > 0 {
			code, err := service.GetCodeFromRedis(c, loginInfo.Phone)
			if err != nil {
				response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
				return
			}
			if code != loginInfo.Code {
				fmt.Println("code error", code)
				response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
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
			if err != nil {
				fmt.Println("err is", err)
				response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
				return
			}

			if code != loginInfo.Code {
				fmt.Println("code err is", err)
				response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
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

// UserUpdatePassword 用户更新密码
func UserUpdatePassword(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var updatePasswordInfo models.UploadPasswordModel
	if err := c.ShouldBind(&updatePasswordInfo); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if updatePasswordInfo.NewPassword != updatePasswordInfo.ConfirmPassword {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	var user models.UserInfo
	result := db.DB.Where("id = ?", userId).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}
	if user.Password != updatePasswordInfo.Password {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	result = db.DB.Model(&user).Where("id = ?", userId).Update("password", updatePasswordInfo.Password)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
		return
	}
	response.Success(c, map[string]interface{}{})
}

func CreateSuggestion(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var suggestion models.SuggestionModel
	if err := c.ShouldBind(&suggestion); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	suggestion.UserId = userId
	result := db.DB.Omit(clause.Associations).Create(&suggestion)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}
	response.Success(c, nil)
}

func GetIpInfo(c *gin.Context) {
	var ipInfo models.IPInfoModel
	if err := c.ShouldBind(&ipInfo); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	url1 := fmt.Sprintf("https://ipapi.co/%s/json/", ipInfo.IP)
	url2 := fmt.Sprintf("https://ipinfo.io/%s/json", ipInfo.IP)
	url3 := fmt.Sprintf("https://ip9.com.cn/get?ip=%s", ipInfo.IP)
	// 获取IP信息
	ipResult, err := GetIPInfoWith(url1, url2, url3)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	response.Success(c, ipResult)
	return
}

// GetIPInfoWith 尝试两个URL获取IP信息
func GetIPInfoWith(url1, url2, url3 string) (*models.IPInfo, error) {
	// 创建HTTP客户端，设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 先尝试第一个URL
	info, err := fetchIPInfo(client, url1)
	if err == nil {
		return info, nil
	}
	fmt.Printf("第一个URL请求失败 (%s): %v\n", url1, err)

	// 第一个失败后尝试第二个URL
	info, err = fetchIPInfo(client, url2)
	if err == nil {
		return info, nil
	}
	fmt.Printf("第二个URL请求失败 (%s): %v\n", url2, err)

	// 第一个失败后尝试第二个URL
	info, err = fetchIPInfo(client, url3)
	if err == nil {
		return info, nil
	}
	fmt.Printf("第三个URL请求失败 (%s): %v\n", url2, err)

	// 两个都失败，返回错误
	return nil, errors.New("三个IP查询URL都请求失败")
}

// fetchIPInfo 从指定URL获取IP信息
func fetchIPInfo(client *http.Client, url string) (*models.IPInfo, error) {
	// 发送HTTP GET请求
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 尝试解析为IPInfo结构体
	var ipInfo models.IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		// 如果解析失败，可能是简单的只返回IP的接口
		// 尝试作为纯文本IP处理
		ip := string(body)
		// 清理可能的空格和换行符
		ip = strings.TrimSpace(ip)

		// 检查是否是纯IP地址（简单验证）
		if len(ip) > 0 {
			ipInfo.IP = ip
			return &ipInfo, nil
		}

		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	return &ipInfo, nil
}

func UploadUserInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var userInfo models.UploadUserInfoModel
	if err := c.ShouldBind(&userInfo); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}
	if len(userInfo.Username) == 0 && len(userInfo.Avatar) == 0 {
		response.Fail(c, response.ApiCode.ParamLack, response.ApiMsg.ParamLack)
		return
	}
	result := db.DB.Model(&models.UserInfo{}).Where("id = ?", userId).
		Update("username", userInfo.Username).
		Update("avatar", userInfo.Avatar)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	response.Success(c, nil)
}

func GetUserInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var userInfo models.UserInfo
	result := db.DB.Where("id = ?", userId).First(&userInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}
	response.Success(c, userInfo)
}

// UserDeactivate 用户退出登录
func UserDeactivate(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	var userInfo models.UserInfo
	result := db.DB.Model(&userInfo).Where("id = ?", userId)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	deactivateResult := result.Delete(&userInfo)
	if deactivateResult.Error != nil {
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}
	response.Success(c, nil)
}

func GetUserList(c *gin.Context) {
	var userList []models.UserInfo
	result := db.DB.Find(&userList)
	if result.Error != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}
	response.Success(c, userList)
}
