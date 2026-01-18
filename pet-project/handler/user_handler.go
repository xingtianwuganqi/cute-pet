package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pet-project/logger"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"
	"pet-project/settings"
	"pet-project/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

func GetEmailCode(c *gin.Context) {
	logger.Logger.Info("Received request for email verification code", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	// 验证码
	lang := c.MustGet("lang").(*i18n.Localizer)
	var param models.SendCodeModel
	paramErr := c.ShouldBind(&param)
	if paramErr != nil {
		logger.Logger.Warn("Invalid parameters for email code request", zap.Error(paramErr))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 调用service层的业务逻辑
	result, err := service.GetEmailCodeService(param.Email, param.Code, lang)
	if err != nil {
		logger.Logger.Error("Failed to send email verification code", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	logger.Logger.Info("Email verification code request processed successfully", zap.String("email", param.Email))

	if settings.Conf.App.Env == "production" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": result,
		})
	}
}

// GetPhoneCode 获取手机验证码
func GetPhoneCode(c *gin.Context) {
	logger.Logger.Info("Received request for phone verification code", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var param models.SendCodeModel
	if err := c.ShouldBind(&param); err != nil {
		logger.Logger.Warn("Invalid parameters for phone code request", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 调用service层的业务逻辑
	result, err := service.GetPhoneCodeService(param.Phone, param.Code)
	if err != nil {
		logger.Logger.Error("Failed to send phone verification code", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	logger.Logger.Info("Phone verification code request processed successfully", zap.String("phone", param.Phone))

	if settings.Conf.App.Env == "production" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": gin.H{},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"data": result,
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
	logger.Logger.Info("Received request for encryption code", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	// 调用 Encrypt 函数对今天的日期进行加密
	encryptionCode, err := util.Encrypt(GetTodayDate())
	if err != nil {
		logger.Logger.Error("Failed to encrypt date", zap.Error(err))
		// 如果加密过程中出现错误，发送失败的 HTTP 响应
		response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
		return
	}
	// 发送包含加密密钥的 JSON 响应
	logger.Logger.Info("Encryption code generated successfully")
	c.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": encryptionCode,
	})
}

// CheckRdbCode 校验验证码
func CheckRdbCode(c *gin.Context) {
	logger.Logger.Info("Received request to check verification code", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var param models.SendCodeModel
	if err := c.ShouldBind(&param); err != nil {
		logger.Logger.Warn("Invalid parameters for verification code check", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.CheckRdbCodeService(param.Phone, param.Email, param.Code)
	if err != nil {
		logger.Logger.Warn("Verification code check failed", zap.Error(err))
		response.Fail(c, response.ApiCode.CheckCodeErr, response.ApiMsg.CheckCodeErr)
	} else {
		logger.Logger.Info("Verification code checked successfully")
		response.Success(c, gin.H{})
	}
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
	logger.Logger.Info("Received user registration request", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var login models.RegisterInfo
	if err := c.ShouldBind(&login); err != nil {
		logger.Logger.Warn("Invalid parameters for user registration", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 调用service层的业务逻辑
	result, err := service.UserRegisterService(login)
	if err != nil {
		logger.Logger.Error("User registration failed", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}

	logger.Logger.Info("User registration successful", zap.Uint("userID", result.ID))
	response.Success(c, result)
}

// UserPhoneLogin 用户登录
func UserPhoneLogin(c *gin.Context) {
	logger.Logger.Info("Received user login request", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var login models.LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		logger.Logger.Warn("Invalid parameters for user login", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.UserPhoneLoginService(login)
	if err != nil {
		logger.Logger.Error("User login failed", zap.Error(err))
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}

	logger.Logger.Info("User login successful", zap.Uint("userID", result.ID))
	response.Success(c, result)
}

// UserFindPassword MARK: 找回密码
func UserFindPassword(c *gin.Context) {
	logger.Logger.Info("Received user password recovery request", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var loginInfo models.RegisterInfo
	if err := c.ShouldBind(&loginInfo); err != nil {
		logger.Logger.Warn("Invalid parameters for password recovery", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.UserFindPasswordService(loginInfo)
	if err != nil {
		logger.Logger.Error("User password recovery failed", zap.Error(err))
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}

	logger.Logger.Info("User password recovery successful")
	response.Success(c, map[string]interface{}{})
}

// UserUpdatePassword 用户更新密码
func UserUpdatePassword(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	logger.Logger.Info("Received user password update request", 
		zap.Uint("userID", userId),
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var updatePasswordInfo models.UploadPasswordModel
	if err := c.ShouldBind(&updatePasswordInfo); err != nil {
		logger.Logger.Warn("Invalid parameters for password update", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.UserUpdatePasswordService(userId, updatePasswordInfo)
	if err != nil {
		logger.Logger.Error("User password update failed", zap.Error(err))
		response.Fail(c, response.ApiCode.ServerErr, response.ApiMsg.ServerErr)
		return
	}

	logger.Logger.Info("User password updated successfully", zap.Uint("userID", userId))
	response.Success(c, map[string]interface{}{})
}

func CreateSuggestion(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	logger.Logger.Info("Received user suggestion creation request", 
		zap.Uint("userID", userId),
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var suggestion models.SuggestionModel
	if err := c.ShouldBind(&suggestion); err != nil {
		logger.Logger.Warn("Invalid parameters for suggestion creation", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.CreateSuggestionService(userId, suggestion)
	if err != nil {
		logger.Logger.Error("User suggestion creation failed", zap.Error(err))
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}

	logger.Logger.Info("User suggestion created successfully", zap.Uint("userID", userId), zap.Uint("suggestionID", suggestion.ID))
	response.Success(c, nil)
}

func GetIpInfo(c *gin.Context) {
	logger.Logger.Info("Received IP info request", 
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var ipInfo models.IPInfoModel
	if err := c.ShouldBind(&ipInfo); err != nil {
		logger.Logger.Warn("Invalid parameters for IP info request", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.GetIpInfoService(ipInfo.IP)
	if err != nil {
		logger.Logger.Error("Failed to get IP info", zap.Error(err))
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

	logger.Logger.Info("IP info retrieved successfully", zap.String("IP", ipInfo.IP))
	response.Success(c, result)
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
	logger.Logger.Debug("First URL request failed", zap.String("url", url1), zap.Error(err))

	// 第一个失败后尝试第二个URL
	info, err = fetchIPInfo(client, url2)
	if err == nil {
		return info, nil
	}
	logger.Logger.Debug("Second URL request failed", zap.String("url", url2), zap.Error(err))

	// 第一个失败后尝试第二个URL
	info, err = fetchIPInfo(client, url3)
	if err == nil {
		return info, nil
	}
	logger.Logger.Debug("Third URL request failed", zap.String("url", url3), zap.Error(err))

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
	logger.Logger.Info("Received user info update request", 
		zap.Uint("userID", userId),
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))
	
	var userInfo models.UploadUserInfoModel
	if err := c.ShouldBind(&userInfo); err != nil {
		logger.Logger.Warn("Invalid parameters for user info update", zap.Error(err))
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.UploadUserInfoService(userId, userInfo)
	if err != nil {
		logger.Logger.Error("User info update failed", zap.Error(err))
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}

	logger.Logger.Info("User info updated successfully", zap.Uint("userID", userId))
	response.Success(c, nil)
}

func GetUserInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	logger.Logger.Info("Received user info retrieval request", 
		zap.Uint("userID", userId),
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))

	result, err := service.GetUserInfoService(userId)
	if err != nil {
		logger.Logger.Error("Failed to retrieve user info", zap.Error(err))
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}

	logger.Logger.Info("User info retrieved successfully", zap.Uint("userID", userId))
	response.Success(c, result)
}

// UserDeactivate 用户注销
func UserDeactivate(c *gin.Context) {
	userId := c.MustGet("userId").(uint)
	logger.Logger.Info("Received user deactivation request", 
		zap.Uint("userID", userId),
		zap.String("clientIP", c.ClientIP()), 
		zap.String("method", c.Request.Method))

	err := service.UserDeactivateService(userId)
	if err != nil {
		logger.Logger.Error("User deactivation failed", zap.Error(err))
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}

	logger.Logger.Info("User deactivated successfully", zap.Uint("userID", userId))
	response.Success(c, nil)
}