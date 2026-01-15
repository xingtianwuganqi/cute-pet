package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"pet-project/models"
	"pet-project/response"
	"pet-project/service"
	"pet-project/settings"
	"pet-project/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

	// 调用service层的业务逻辑
	result, err := service.GetEmailCodeService(param.Email, param.Code, lang)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

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
	var param models.SendCodeModel
	if err := c.ShouldBind(&param); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 调用service层的业务逻辑
	result, err := service.GetPhoneCodeService(param.Phone, param.Code)
	if err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

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

	err := service.CheckRdbCodeService(param.Phone, param.Email, param.Code)
	if err != nil {
		response.Fail(c, response.ApiCode.CheckCodeErr, response.ApiMsg.CheckCodeErr)
	} else {
		response.Success(c, gin.H{})
	}
}

// UserRegister 注册
func UserRegister(c *gin.Context) {
	var login models.RegisterInfo
	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	// 调用service层的业务逻辑
	result, err := service.UserRegisterService(login)
	if err != nil {
		response.Fail(c, response.ApiCode.CreateErr, response.ApiMsg.CreateErr)
		return
	}

	response.Success(c, result)
}

// UserPhoneLogin 用户登录
func UserPhoneLogin(c *gin.Context) {
	var login models.LoginInfo
	if err := c.ShouldBind(&login); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	result, err := service.UserPhoneLoginService(login)
	if err != nil {
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}

	response.Success(c, result)
}

// UserFindPassword MARK: 找回密码
func UserFindPassword(c *gin.Context) {
	var loginInfo models.RegisterInfo
	if err := c.ShouldBind(&loginInfo); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.UserFindPasswordService(loginInfo)
	if err != nil {
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}

	response.Success(c, map[string]interface{}{})
}

// UserUpdatePassword 用户更新密码
func UserUpdatePassword(c *gin.Context) {
	var userId = c.MustGet("userId").(uint)
	var updatePasswordInfo models.UploadPasswordModel
	if err := c.ShouldBind(&updatePasswordInfo); err != nil {
		response.Fail(c, response.ApiCode.ParamErr, response.ApiMsg.ParamErr)
		return
	}

	err := service.UserUpdatePasswordService(userId, updatePasswordInfo)
	if err != nil {
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

	err := service.CreateSuggestionService(userId, suggestion)
	if err != nil {
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

	result, err := service.GetIpInfoService(ipInfo.IP)
	if err != nil {
		response.Fail(c, response.ApiCode.QueryErr, response.ApiMsg.QueryErr)
		return
	}

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

	err := service.UploadUserInfoService(userId, userInfo)
	if err != nil {
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}

	response.Success(c, nil)
}

func GetUserInfo(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	result, err := service.GetUserInfoService(userId)
	if err != nil {
		response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
		return
	}

	response.Success(c, result)
}

// UserDeactivate 用户注销
func UserDeactivate(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	err := service.UserDeactivateService(userId)
	if err != nil {
		response.Fail(c, response.ApiCode.UpdateErr, response.ApiMsg.UpdateErr)
		return
	}

	response.Success(c, nil)
}
