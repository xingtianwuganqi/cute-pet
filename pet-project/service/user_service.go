package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pet-project/db"
	"pet-project/internal"
	"pet-project/logger"
	"pet-project/middleware"
	"pet-project/models"
	"pet-project/settings"
	"pet-project/util"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetEmailCodeService 处理邮箱验证码发送的业务逻辑
func GetEmailCodeService(email string, code string, lang *i18n.Localizer) (interface{}, error) {
	logger.Logger.Info("Attempting to send email verification code", zap.String("email", email))
	
	// 查询code是否在redis中（是否已经使用过了）
	codeKey := fmt.Sprintf("email_code:%s", email)
	value, err := internal.GetCodeFromRedis(nil, codeKey)
	if err != nil {
		logger.Logger.Error("Error retrieving code from redis", zap.Error(err), zap.String("codeKey", codeKey))
		return nil, errors.New("invalid param")
	}
	if len(value) != 0 && value == code {
		logger.Logger.Warn("Verification code already used", zap.String("email", email))
		return nil, errors.New("invalid param")
	}

	// 需要加一个加密信息
	encryptionStr, err := util.Decrypt(code)
	if err != nil {
		logger.Logger.Error("Error decrypting code", zap.Error(err))
		return nil, errors.New("invalid param")
	}
	// 判断encryptionStr是否今日日期
	if len(encryptionStr) == 0 || encryptionStr != GetTodayDate() {
		logger.Logger.Warn("Verification code is expired", zap.String("decrypted", encryptionStr))
		return nil, errors.New("invalid param")
	}

	if len(email) == 0 {
		logger.Logger.Warn("Email is empty")
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
			logger.Logger.Error("Failed to send email", zap.Error(sendErr), zap.String("email", email))
			return nil, sendErr
		}
		logger.Logger.Info("Email sent successfully", zap.String("email", email))
	}

	// 将code保存到redis，设置10分钟失效
	saveErr := internal.SaveAccountCodeInRedis(nil, email, sendCode, 10*time.Minute)
	if saveErr != nil {
		logger.Logger.Error("Failed to save code to redis", zap.Error(saveErr), zap.String("email", email))
		return nil, saveErr
	}

	// 保存param.code
	_ = internal.SaveAccountCodeInRedis(nil, codeKey, code, 24*time.Hour)

	logger.Logger.Info("Email verification code processed successfully", zap.String("email", email))
	if settings.Conf.App.Env == "production" {
		return map[string]interface{}{}, nil
	} else {
		return sendCode, nil
	}
}

// GetPhoneCodeService 处理手机验证码发送的业务逻辑
func GetPhoneCodeService(phone string, code string) (interface{}, error) {
	logger.Logger.Info("Attempting to send phone verification code", zap.String("phone", phone))
	
	// 查询code是否在redis中（是否已经使用过了）
	codeKey := fmt.Sprintf("phone_code:%s", phone)
	value, err := internal.GetCodeFromRedis(nil, codeKey)
	if err != nil {
		logger.Logger.Error("Error retrieving code from redis", zap.Error(err), zap.String("codeKey", codeKey))
		return nil, errors.New("invalid param")
	}
	if len(value) != 0 && value == code {
		logger.Logger.Warn("Phone verification code already used", zap.String("phone", phone))
		return nil, errors.New("invalid param")
	}

	// 需要加一个加密信息
	encryptionStr, err := util.Decrypt(code)
	if err != nil {
		logger.Logger.Error("Error decrypting code", zap.Error(err))
		return nil, errors.New("invalid param")
	}
	// 判断encryptionStr是否今日日期
	if len(encryptionStr) == 0 || encryptionStr != GetTodayDate() {
		logger.Logger.Warn("Phone verification code is expired", zap.String("decrypted", encryptionStr))
		return nil, errors.New("invalid param")
	}

	if len(phone) == 0 {
		logger.Logger.Warn("Phone number is empty")
		return nil, errors.New("invalid param")
	}

	sendCode := internal.GenerateValidationCode(4)

	if settings.Conf.App.Env == "production" {
		url := fmt.Sprintf("https://push.spug.cc/send/gL1QGmWdKWjlRD65?key1=%s&key2=%s&key3=%s&targets=%s",
			"[Pawpal]", sendCode, "10", phone)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			logger.Logger.Error("Failed to send phone verification code", zap.Error(err), zap.Int("statusCode", resp.StatusCode))
			return nil, errors.New("failed to send verification code")
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		logger.Logger.Info("Phone verification code sent successfully", zap.String("phone", phone))
	}

	// 将code保存到redis，设置10分钟失效
	saveErr := internal.SaveAccountCodeInRedis(nil, phone, sendCode, 10*time.Minute)
	if saveErr != nil {
		logger.Logger.Error("Failed to save code to redis", zap.Error(saveErr), zap.String("phone", phone))
		return nil, saveErr
	}

	_ = internal.SaveAccountCodeInRedis(nil, codeKey, code, 24*time.Hour)

	logger.Logger.Info("Phone verification code processed successfully", zap.String("phone", phone))
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
	logger.Logger.Info("Checking verification code", zap.String("phone", phone), zap.String("email", email))
	
	if len(phone) != 0 {
		storedCode, err := internal.GetCodeFromRedis(nil, phone)
		if err != nil {
			logger.Logger.Error("Error retrieving phone code from redis", zap.Error(err), zap.String("phone", phone))
			return errors.New("query error")
		}
		if storedCode == code {
			logger.Logger.Info("Phone verification code verified successfully", zap.String("phone", phone))
			return nil
		} else {
			logger.Logger.Warn("Incorrect phone verification code", zap.String("phone", phone))
			return errors.New("verification code error")
		}
	} else {
		storedCode, err := internal.GetCodeFromRedis(nil, email)
		if err != nil {
			logger.Logger.Error("Error retrieving email code from redis", zap.Error(err), zap.String("email", email))
			return errors.New("query error")
		}

		if storedCode == code {
			logger.Logger.Info("Email verification code verified successfully", zap.String("email", email))
			return nil
		} else {
			logger.Logger.Warn("Incorrect email verification code", zap.String("email", email))
			return errors.New("verification code error")
		}
	}
}

// UserRegisterService 用户注册业务逻辑
func UserRegisterService(registerInfo models.RegisterInfo) (models.LoginUserInfo, error) {
	logger.Logger.Info("User registration attempt", zap.String("email", registerInfo.Email), zap.String("phone", registerInfo.Phone))
	
	var findUser models.UserInfo
	var findResult *gorm.DB
	if len(registerInfo.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", registerInfo.Phone).First(&findUser)
	} else if len(registerInfo.Email) > 0 {
		if util.IsValidEmail(registerInfo.Email) {
			findResult = db.DB.Where("email = ?", registerInfo.Email).First(&findUser)
		} else {
			logger.Logger.Error("Invalid email format", zap.String("email", registerInfo.Email))
			return models.LoginUserInfo{}, errors.New("invalid email format")
		}
	} else {
		logger.Logger.Error("Phone or email required")
		return models.LoginUserInfo{}, errors.New("phone or email required")
	}

	// 如果查不到，则开始验证验证码
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		// 取出redis中的验证码
		if len(registerInfo.Email) > 0 {
			storedCode, err := internal.GetCodeFromRedis(nil, registerInfo.Email)
			if err != nil {
				logger.Logger.Error("Error retrieving email code from redis", zap.Error(err), zap.String("email", registerInfo.Email))
				return models.LoginUserInfo{}, errors.New("server error")
			}
			if storedCode != registerInfo.Code {
				logger.Logger.Warn("Invalid email verification code", zap.String("email", registerInfo.Email))
				return models.LoginUserInfo{}, errors.New("invalid verification code")
			}
			_ = internal.DeleteCodeFromRedis(nil, registerInfo.Email)
		} else { // 验证手机验证码
			storedCode, err := internal.GetCodeFromRedis(nil, registerInfo.Phone)
			if err != nil {
				logger.Logger.Error("Error retrieving phone code from redis", zap.Error(err), zap.String("phone", registerInfo.Phone))
				return models.LoginUserInfo{}, errors.New("server error")
			}

			// 验证验证码是否正确
			if storedCode != registerInfo.Code {
				logger.Logger.Warn("Invalid phone verification code", zap.String("phone", registerInfo.Phone))
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
			logger.Logger.Error("Failed to create user", zap.Error(result.Error))
			return models.LoginUserInfo{}, errors.New("create user failed")
		}
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			logger.Logger.Error("Failed to generate token", zap.Error(err))
			return models.LoginUserInfo{}, errors.New("generate token failed")
		}
		loginUserInfo := models.LoginUserInfo{
			ID:     user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		logger.Logger.Info("User registered successfully", zap.Uint("userID", user.ID))
		return loginUserInfo, nil
	} else {
		logger.Logger.Warn("Registration attempt for existing user", zap.String("email", registerInfo.Email), zap.String("phone", registerInfo.Phone))
		return models.LoginUserInfo{}, errors.New("user already exists")
	}
}

// UserPhoneLoginService 用户登录业务逻辑
func UserPhoneLoginService(loginInfo models.LoginInfo) (models.LoginUserInfo, error) {
	logger.Logger.Info("User login attempt", zap.String("email", loginInfo.Email), zap.String("phone", loginInfo.Phone))
	
	var findResult *gorm.DB
	var user models.UserInfo
	if len(loginInfo.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", loginInfo.Phone).First(&user)
	} else if len(loginInfo.Email) > 0 {
		if util.IsValidEmail(loginInfo.Email) {
			findResult = db.DB.Where("email = ?", loginInfo.Email).First(&user)
		} else {
			logger.Logger.Error("Invalid email format", zap.String("email", loginInfo.Email))
			return models.LoginUserInfo{}, errors.New("invalid email format")
		}
	} else {
		logger.Logger.Error("Phone or email required")
		return models.LoginUserInfo{}, errors.New("phone or email required")
	}
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		logger.Logger.Warn("Login attempt for non-existent user", zap.String("email", loginInfo.Email), zap.String("phone", loginInfo.Phone))
		return models.LoginUserInfo{}, errors.New("user not found")
	}
	if user.Password == loginInfo.Password {
		// 密码正确, 生成token，登录完成
		userId := user.ID
		token, err := middleware.GenToken(userId)
		if err != nil {
			logger.Logger.Error("Failed to generate token", zap.Error(err))
			return models.LoginUserInfo{}, errors.New("generate token failed")
		}
		loginUserInfo := models.LoginUserInfo{
			ID:     user.ID,
			Phone:  user.Phone,
			Avatar: user.Avatar,
			Email:  user.Email,
			Token:  token,
		}
		logger.Logger.Info("User logged in successfully", zap.Uint("userID", user.ID))
		return loginUserInfo, nil
	} else {
		logger.Logger.Warn("Login attempt with incorrect password", zap.String("email", loginInfo.Email), zap.String("phone", loginInfo.Phone))
		return models.LoginUserInfo{}, errors.New("password incorrect")
	}
}

// UserFindPasswordService 用户找回密码业务逻辑
func UserFindPasswordService(loginInfo models.RegisterInfo) error {
	logger.Logger.Info("User password recovery attempt", zap.String("email", loginInfo.Email), zap.String("phone", loginInfo.Phone))
	
	var findResult *gorm.DB
	var user models.UserInfo
	if len(loginInfo.Phone) > 0 {
		findResult = db.DB.Where("phone = ?", loginInfo.Phone).First(&user)
	} else if len(loginInfo.Email) > 0 {
		if util.IsValidEmail(loginInfo.Email) {
			findResult = db.DB.Where("email = ?", loginInfo.Email).First(&user)
		} else {
			logger.Logger.Error("Invalid email format", zap.String("email", loginInfo.Email))
			return errors.New("invalid param")
		}
	} else {
		logger.Logger.Error("Phone or email required")
		return errors.New("phone or email required")
	}
	if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
		logger.Logger.Warn("Password recovery attempt for non-existent user", zap.String("email", loginInfo.Email), zap.String("phone", loginInfo.Phone))
		return errors.New("user not found")
	} else {
		// 验证验证码
		if len(loginInfo.Phone) > 0 {
			storedCode, err := internal.GetCodeFromRedis(nil, loginInfo.Phone)
			if err != nil {
				logger.Logger.Error("Error retrieving phone code from redis", zap.Error(err), zap.String("phone", loginInfo.Phone))
				return errors.New("invalid param")
			}
			if storedCode != loginInfo.Code {
				logger.Logger.Warn("Invalid phone verification code during password recovery", zap.String("phone", loginInfo.Phone))
				return errors.New("invalid param")
			}
			// 更新密码
			result := db.DB.Model(&user).Where("phone = ?", loginInfo.Phone).Update("password", loginInfo.Password)
			if result.Error != nil {
				logger.Logger.Error("Failed to update password", zap.Error(result.Error), zap.String("phone", loginInfo.Phone))
				return errors.New("server error")
			}
			// redis的数据清除
			_ = internal.DeleteCodeFromRedis(nil, loginInfo.Phone)

			logger.Logger.Info("Password updated successfully via phone", zap.String("phone", loginInfo.Phone))
			return nil
		} else {
			storedCode, err := internal.GetCodeFromRedis(nil, loginInfo.Email)
			if err != nil {
				logger.Logger.Error("Error retrieving email code from redis", zap.Error(err), zap.String("email", loginInfo.Email))
				return errors.New("invalid param")
			}

			if storedCode != loginInfo.Code {
				logger.Logger.Warn("Invalid email verification code during password recovery", zap.String("email", loginInfo.Email))
				return errors.New("invalid param")
			}
			result := db.DB.Model(&user).Where("email = ?", loginInfo.Email).Update("password", loginInfo.Password)
			if result.Error != nil {
				logger.Logger.Error("Failed to update password", zap.Error(result.Error), zap.String("email", loginInfo.Email))
				return errors.New("server error")
			}

			// 删除redis数据
			_ = internal.DeleteCodeFromRedis(nil, loginInfo.Email)

			logger.Logger.Info("Password updated successfully via email", zap.String("email", loginInfo.Email))
			return nil
		}
	}
}

// UserUpdatePasswordService 用户更新密码业务逻辑
func UserUpdatePasswordService(userId uint, updatePasswordInfo models.UploadPasswordModel) error {
	logger.Logger.Info("User attempting to update password", zap.Uint("userID", userId))
	
	if updatePasswordInfo.NewPassword != updatePasswordInfo.ConfirmPassword {
		logger.Logger.Warn("New password does not match confirmation", zap.Uint("userID", userId))
		return errors.New("password mismatch")
	}
	var user models.UserInfo
	result := db.DB.Where("id = ?", userId).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Logger.Warn("Attempt to update password for non-existent user", zap.Uint("userID", userId))
		return errors.New("user not found")
	}
	if user.Password != updatePasswordInfo.Password {
		logger.Logger.Warn("Current password incorrect", zap.Uint("userID", userId))
		return errors.New("current password incorrect")
	}
	result = db.DB.Model(&user).Where("id = ?", userId).Update("password", updatePasswordInfo.NewPassword)
	if result.Error != nil {
		logger.Logger.Error("Failed to update password", zap.Error(result.Error), zap.Uint("userID", userId))
		return errors.New("update failed")
	}
	logger.Logger.Info("Password updated successfully", zap.Uint("userID", userId))
	return nil
}

// CreateSuggestionService 创建建议
func CreateSuggestionService(userId uint, suggestion models.SuggestionModel) error {
	logger.Logger.Info("Creating user suggestion", zap.Uint("userID", userId))
	
	suggestion.UserId = userId
	result := db.DB.Create(&suggestion)
	if result.Error != nil {
		logger.Logger.Error("Failed to create suggestion", zap.Error(result.Error), zap.Uint("userID", userId))
		return errors.New("create suggestion failed")
	}
	logger.Logger.Info("Suggestion created successfully", zap.Uint("userID", userId), zap.Uint("suggestionID", suggestion.ID))
	return nil
}

// GetUserInfoService 获取用户信息
func GetUserInfoService(userId uint) (models.UserInfo, error) {
	logger.Logger.Info("Retrieving user info", zap.Uint("userID", userId))
	
	var userInfo models.UserInfo
	result := db.DB.Where("id = ?", userId).First(&userInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Logger.Warn("Attempt to retrieve non-existent user info", zap.Uint("userID", userId))
		return models.UserInfo{}, errors.New("user not found")
	}
	logger.Logger.Info("User info retrieved successfully", zap.Uint("userID", userId))
	return userInfo, nil
}

// UploadUserInfoService 更新用户信息
func UploadUserInfoService(userId uint, userInfo models.UploadUserInfoModel) error {
	logger.Logger.Info("Updating user info", zap.Uint("userID", userId))
	
	if len(userInfo.Username) == 0 && len(userInfo.Avatar) == 0 {
		logger.Logger.Warn("Username or avatar required", zap.Uint("userID", userId))
		return errors.New("username or avatar required")
	}
	result := db.DB.Model(&models.UserInfo{}).Where("id = ?", userId).
		Update("username", userInfo.Username).
		Update("avatar", userInfo.Avatar)
	if result.Error != nil {
		logger.Logger.Error("Failed to update user info", zap.Error(result.Error), zap.Uint("userID", userId))
		return errors.New("update failed")
	}
	logger.Logger.Info("User info updated successfully", zap.Uint("userID", userId))
	return nil
}

// UserDeactivateService 用户注销
func UserDeactivateService(userId uint) error {
	logger.Logger.Info("Deactivating user account", zap.Uint("userID", userId))
	
	var userInfo models.UserInfo
	result := db.DB.Model(&userInfo).Where("id = ?", userId)
	if result.Error != nil {
		logger.Logger.Error("Query error during user deactivation", zap.Error(result.Error), zap.Uint("userID", userId))
		return errors.New("query error")
	}
	// 删除所有发布的信息
	recordResult := db.DB.Where("user_id = ?", userId).Delete(&models.RecordList{})
	if recordResult.Error != nil {
		logger.Logger.Error("Failed to delete user records", zap.Error(recordResult.Error), zap.Uint("userID", userId))
		return errors.New("update failed")
	}
	// 删除所有发布的帖子
	postResult := db.DB.Where("user_id = ?", userId).Delete(&models.PostModel{})
	if postResult.Error != nil {
		logger.Logger.Error("Failed to delete user posts", zap.Error(postResult.Error), zap.Uint("userID", userId))
		return errors.New("update failed")
	}
	deactivateResult := result.Delete(&userInfo)
	if deactivateResult.Error != nil {
		logger.Logger.Error("Failed to deactivate user", zap.Error(deactivateResult.Error), zap.Uint("userID", userId))
		return errors.New("update failed")
	}
	logger.Logger.Info("User deactivated successfully", zap.Uint("userID", userId))
	return nil
}

// GetIpInfoService 获取IP信息
func GetIpInfoService(ip string) (*models.IPInfo, error) {
	logger.Logger.Info("Getting IP info", zap.String("ip", ip))
	
	url1 := fmt.Sprintf("https://ipapi.co/%s/json/", ip)
	url2 := fmt.Sprintf("https://ipinfo.io/%s/json", ip)
	url3 := fmt.Sprintf("https://ip9.com.cn/get?ip=%s", ip)
	// 获取IP信息
	ipResult, err := GetIPInfoWith(url1, url2, url3)
	if err != nil {
		logger.Logger.Error("Failed to get IP info", zap.Error(err), zap.String("ip", ip))
		return nil, errors.New("query error")
	}

	logger.Logger.Info("IP info retrieved successfully", zap.String("ip", ip))
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
	logger.Logger.Debug("First URL request failed", zap.String("url", url1), zap.Error(err))

	// 第二个URL
	info, err = fetchIPInfo(client, url2)
	if err == nil {
		return info, nil
	}
	logger.Logger.Debug("Second URL request failed", zap.String("url", url2), zap.Error(err))

	// 第三个URL
	info, err = fetchIPInfo(client, url3)
	if err == nil {
		return info, nil
	}
	logger.Logger.Debug("Third URL request failed", zap.String("url", url3), zap.Error(err))

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