package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pet-project/db"
	"pet-project/internal"
	"pet-project/middleware"
	"pet-project/models"
	"pet-project/settings"
	"pet-project/util"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

// GetEmailCodeService 处理邮箱验证码发送的业务逻辑
func GetEmailCodeService(email string, code string, lang *i18n.Localizer) (interface{}, error) {
	// 查询code是否在redis中（是否已经使用过了）
	codeKey := fmt.Sprintf("email_code:%s", email)
	value, err := internal.GetCodeFromRedis(nil, codeKey)
	if err != nil {
		return nil, errors.New("invalid param")
	}
	if len(value) != 0 && value == code {
		return nil, errors.New("invalid param")
	}

	// 需要加一个加密信息
	encryptionStr, err := util.Decrypt(code)
	if err != nil {
		return nil, errors.New("invalid param")
	}
	// 判断encryptionStr是否今日日期
	if len(encryptionStr) == 0 || encryptionStr != GetTodayDate() {
		return nil, errors.New("invalid param")
	}

	if len(email) == 0 {
		return nil, errors.New("invalid param")
	}

	sendCode := internal.GenerateValidationCode(4)

	// 正式环境发验证码
	if settings.Conf.App.Env == "production" {
		smptServer := settings.Conf.EmailService.Host
		smptPort := settings.Conf.EmailService.Port
		username := settings.Conf.EmailService.Username
		password := settings.Conf.EmailService.Password
		recipient := email
		subject := internal.LocalizeMsg(lang, "VerificationTitle")
		body := internal.LocalizeMsgCount(lang, "VerificationDesc", sendCode)

		sendErr := internal.SendEmail(recipient, subject, body, smptServer, smptPort, username, password)
		if sendErr != nil {
			return nil, sendErr
		}
	}

	// 将code保存到redis，设置10分钟失效
	saveErr := internal.SaveAccountCodeInRedis(nil, email, sendCode, 10*time.Minute)
	if saveErr != nil {
		return nil, saveErr
	}

	// 保存param.code
	_ = internal.SaveAccountCodeInRedis(nil, codeKey, code, 24*time.Hour)

	if settings.Conf.App.Env == "production" {
		return map[string]interface{}{}, nil
	} else {
		return sendCode, nil
	}
}

// GetPhoneCodeService 处理手机验证码发送的业务逻辑
func GetPhoneCodeService(phone string, code string) (interface{}, error) {
	// 查询code是否在redis中（是否已经使用过了）
	codeKey := fmt.Sprintf("phone_code:%s", phone)
	value, err := internal.GetCodeFromRedis(nil, codeKey)
	if err != nil {
		return nil, errors.New("invalid param")
	}
	if len(value) != 0 && value == code {
		return nil, errors.New("invalid param")
	}

	// 需要加一个加密信息
	encryptionStr, err := util.Decrypt(code)
	if err != nil {
		return nil, errors.New("invalid param")
	}
	// 判断encryptionStr是否今日日期
	if len(encryptionStr) == 0 || encryptionStr != GetTodayDate() {
		return nil, errors.New("invalid param")
	}

	if len(phone) == 0 {
		return nil, errors.New("invalid param")
	}

	sendCode := internal.GenerateValidationCode(4)

	if settings.Conf.App.Env == "production" {
		url := fmt.Sprintf("https://push.spug.cc/send/gL1QGmWdKWjlRD65?key1=%s&key2=%s&key3=%s&targets=%s",
			"[Pawpal]", sendCode, "10", phone)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, errors.New("failed to send verification code")
		}
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	// 将code保存到redis，设置10分钟失效
	saveErr := internal.SaveAccountCodeInRedis(nil, phone, sendCode, 10*time.Minute)
	if saveErr != nil {
		return nil, saveErr
	}

	_ = internal.SaveAccountCodeInRedis(nil, codeKey, code, 24*time.Hour)

	if settings.Conf.App.Env == "production" {
		return map[string]interface{}{}, nil
	} else {
		return sendCode, nil
	}
}

// GetTodayDate 返回今天日期字符串
func GetTodayDate() string {
	now := time.Now()
	date := now.Format("2006-01-02")
	return date
}

// CheckRdbCodeService 验证验证码
func CheckRdbCodeService(phone string, email string, code string) error {
	if len(phone) != 0 {
		storedCode, err := internal.GetCodeFromRedis(nil, phone)
		if err != nil {
			return errors.New("query error")
		}
		if storedCode == code {
			return nil
		} else {
			return errors.New("verification code error")
		}
	} else {
		storedCode, err := internal.GetCodeFromRedis(nil, email)
		if err != nil {
			return errors.New("query error")
		}

		if storedCode == code {
			return nil
		} else {
			return errors.New("verification code error")
		}
	}
}

// UserRegisterService 用户注册业务逻辑
func UserRegisterService(registerInfo models.RegisterInfo) (models.LoginUserInfo, error) {
	var findUser models.UserInfo
	var findResult *gorm.DB
	if len(registerInfo.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", registerInfo.Phone).First(&findUser)
	} else if len(registerInfo.Email) > 0 {
		if util.IsValidEmail(registerInfo.Email) {
			findResult = db.DB.Where("email = ?", registerInfo.Email).First(&findUser)
		} else {
			return models.LoginUserInfo{}, errors.New("invalid email format")
		}
	} else {
		return models.LoginUserInfo{}, errors.New("phone or email required")
	}

	// 如果查不到，则开始验证验证码
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		// 取出redis中的验证码
		if len(registerInfo.Email) > 0 {
			storedCode, err := internal.GetCodeFromRedis(nil, registerInfo.Email)
			if err != nil {
				return models.LoginUserInfo{}, errors.New("server error")
			}
			if storedCode != registerInfo.Code {
				return models.LoginUserInfo{}, errors.New("invalid verification code")
			}
			_ = internal.DeleteCodeFromRedis(nil, registerInfo.Email)
		} else { // 验证手机验证码
			storedCode, err := internal.GetCodeFromRedis(nil, registerInfo.Phone)
			if err != nil {
				return models.LoginUserInfo{}, errors.New("server error")
			}

			// 验证验证码是否正确
			if storedCode != registerInfo.Code {
				return models.LoginUserInfo{}, errors.New("invalid verification code")
			}
			_ = internal.DeleteCodeFromRedis(nil, registerInfo.Email)
		}

		user := models.UserInfo{
			Phone:    registerInfo.Phone,
			Password: registerInfo.Password,
			Email:    registerInfo.Email,
		}
		result := db.DB.Create(&user)
		if result.Error != nil {
			return models.LoginUserInfo{}, errors.New("create user failed")
		}
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			return models.LoginUserInfo{}, errors.New("generate token failed")
		}
		loginUserInfo := models.LoginUserInfo{
			ID:     user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		return loginUserInfo, nil
	} else {
		return models.LoginUserInfo{}, errors.New("user already exists")
	}
}

// UserPhoneLoginService 用户登录业务逻辑
func UserPhoneLoginService(loginInfo models.LoginInfo) (models.LoginUserInfo, error) {
	var findResult *gorm.DB
	var user models.UserInfo
	if len(loginInfo.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", loginInfo.Phone).First(&user)
	} else if len(loginInfo.Email) > 0 {
		if util.IsValidEmail(loginInfo.Email) {
			findResult = db.DB.Where("email = ?", loginInfo.Email).First(&user)
		} else {
			return models.LoginUserInfo{}, errors.New("invalid email format")
		}
	} else {
		return models.LoginUserInfo{}, errors.New("phone or email required")
	}
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return models.LoginUserInfo{}, errors.New("user not found")
	}
	if user.Password == loginInfo.Password {
		// 密码正确, 生成token，登录完成
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			return models.LoginUserInfo{}, errors.New("generate token failed")
		}
		loginUserInfo := models.LoginUserInfo{
			ID:     user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		return loginUserInfo, nil
	} else {
		return models.LoginUserInfo{}, errors.New("password incorrect")
	}
}

// UserFindPasswordService 用户找回密码业务逻辑
func UserFindPasswordService(loginInfo models.RegisterInfo) error {
	var findResult *gorm.DB
	var user models.UserInfo
	if len(loginInfo.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", loginInfo.Phone).First(&user)
	} else if len(loginInfo.Email) > 0 {
		if util.IsValidEmail(loginInfo.Email) {
			findResult = db.DB.Where("email = ?", loginInfo.Email).First(&user)
		} else {
			return errors.New("invalid param")
		}
	} else {
		return errors.New("phone or email required")
	}
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	} else {
		// 验证验证码
		if len(loginInfo.Phone) > 0 {
			storedCode, err := internal.GetCodeFromRedis(nil, loginInfo.Phone)
			if err != nil {
				return errors.New("invalid param")
			}
			if storedCode != loginInfo.Code {
				fmt.Println("code error", storedCode)
				return errors.New("invalid param")
			}
			// 更新密码
			result := db.DB.Model(&user).Where("phone = ?", loginInfo.Phone).Update("password", loginInfo.Password)
			if result.Error != nil {
				return errors.New("server error")
			}
			// redis的数据清除
			_ = internal.DeleteCodeFromRedis(nil, loginInfo.Phone)

			return nil
		} else {
			storedCode, err := internal.GetCodeFromRedis(nil, loginInfo.Email)
			if err != nil {
				fmt.Println("err is", err)
				return errors.New("invalid param")
			}

			if storedCode != loginInfo.Code {
				fmt.Println("code err is", err)
				return errors.New("invalid param")
			}
			result := db.DB.Model(&user).Where("email = ?", loginInfo.Email).Update("password", loginInfo.Password)
			if result.Error != nil {
				return errors.New("server error")
			}

			// 删除redis数据
			_ = internal.DeleteCodeFromRedis(nil, loginInfo.Email)

			return nil
		}
	}
}

// UserUpdatePasswordService 用户更新密码业务逻辑
func UserUpdatePasswordService(userId uint, updatePasswordInfo models.UploadPasswordModel) error {
	if updatePasswordInfo.NewPassword != updatePasswordInfo.ConfirmPassword {
		return errors.New("password mismatch")
	}
	var user models.UserInfo
	result := db.DB.Where("id = ?", userId).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}
	if user.Password != updatePasswordInfo.Password {
		return errors.New("current password incorrect")
	}
	result = db.DB.Model(&user).Where("id = ?", userId).Update("password", updatePasswordInfo.NewPassword)
	if result.Error != nil {
		return errors.New("update failed")
	}
	return nil
}

// CreateSuggestionService 创建建议
func CreateSuggestionService(userId uint, suggestion models.SuggestionModel) error {
	suggestion.UserId = userId
	result := db.DB.Create(&suggestion)
	if result.Error != nil {
		return errors.New("create suggestion failed")
	}
	return nil
}

// GetUserInfoService 获取用户信息
func GetUserInfoService(userId uint) (models.UserInfo, error) {
	var userInfo models.UserInfo
	result := db.DB.Where("id = ?", userId).First(&userInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.UserInfo{}, errors.New("user not found")
	}
	return userInfo, nil
}

// UploadUserInfoService 更新用户信息
func UploadUserInfoService(userId uint, userInfo models.UploadUserInfoModel) error {
	if len(userInfo.Username) == 0 && len(userInfo.Avatar) == 0 {
		return errors.New("username or avatar required")
	}
	result := db.DB.Model(&models.UserInfo{}).Where("id = ?", userId).
		Update("username", userInfo.Username).
		Update("avatar", userInfo.Avatar)
	if result.Error != nil {
		return errors.New("update failed")
	}
	return nil
}

// UserDeactivateService 用户注销
func UserDeactivateService(userId uint) error {
	var userInfo models.UserInfo
	result := db.DB.Model(&userInfo).Where("id = ?", userId)
	if result.Error != nil {
		return errors.New("query error")
	}
	// 删除所有发布的信息
	recordResult := db.DB.Where("user_id = ?", userId).Delete(&models.RecordList{})
	if recordResult.Error != nil {
		return errors.New("update failed")
	}
	// 删除所有发布的帖子
	postResult := db.DB.Where("user_id = ?", userId).Delete(&models.PostModel{})
	if postResult.Error != nil {
		return errors.New("update failed")
	}
	deactivateResult := result.Delete(&userInfo)
	if deactivateResult.Error != nil {
		return errors.New("update failed")
	}
	return nil
}

// GetIpInfoService 获取IP信息
func GetIpInfoService(ip string) (*models.IPInfo, error) {
	url1 := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
	url2 := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	url3 := fmt.Sprintf("https://ip9.com.cn/get?ip=%s", ip)
	// 获取IP信息
	ipResult, err := GetIPInfoWith(url1, url2, url3)
	if err != nil {
		return nil, errors.New("query error")
	}

	return ipResult, nil
}

// GetIPInfoWith 尝试多个URL获取IP信息
func GetIPInfoWith(url1, url2, url3 string) (*models.IPInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 先尝试第一个URL
	info, err := fetchIPInfo(client, url1)
	if err == nil {
		return info, nil
	}
	fmt.Printf("第一个URL请求失败 (%s): %v\n", url1, err)

	// 第二个URL
	info, err = fetchIPInfo(client, url2)
	if err == nil {
		return info, nil
	}
	fmt.Printf("第二个URL请求失败 (%s): %v\n", url2, err)

	// 第三个URL
	info, err = fetchIPInfo(client, url3)
	if err == nil {
		return info, nil
	}
	fmt.Printf("第三个URL请求失败 (%s): %v\n", url3, err)

	// 都失败，返回错误
	return nil, errors.New("all IP query URLs failed")
}

// fetchIPInfo 从指定URL获取IP信息
func fetchIPInfo(client *http.Client, url string) (*models.IPInfo, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	var ipInfo models.IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		// 如果解析失败，可能是简单的只返回IP的接口
		ip := string(body)
		ip = strings.TrimSpace(ip)

		if len(ip) > 0 {
			ipInfo.IP = ip
			return &ipInfo, nil
		}

		return nil, fmt.Errorf("JSON parse failed: %w", err)
	}

	return &ipInfo, nil
}