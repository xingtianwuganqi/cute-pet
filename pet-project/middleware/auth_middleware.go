package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"pet-project/db"
	"pet-project/internal"
	"pet-project/models"
	"pet-project/response"
	"pet-project/settings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MyClaims struct {
	UserId uint `json:"userId"`
	jwt.StandardClaims
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

var mySecret = []byte("伍c七Alz1θVx2ψLHNpfωv九nξ捌τD六053λwGμrMνRuegsη八γ陆jOBX8ρ三E9πFS零bδοmkχ7K6PβϵϕoZ五iυU一Jq柒ydYt四QhW4玖κCIαζTaι二σ")

//创建token

func GenToken(userId uint) (string, error) {
	claims := jwt.MapClaims{}
    claims["userId"] = userId
    claims["exp"] = time.Now().AddDate(30, 0, 0).Unix()
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if Claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return Claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTTokenMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if settings.Conf.App.Env != "production" {
			fmt.Printf("token: %s\n", token)
		}
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  internal.LocalizeMsg(c.MustGet("lang").(*i18n.Localizer), response.ApiMsg.AuthErr),
				"data": map[string]interface{}{},
			})
			c.Abort()
			return
		}
		mc, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  internal.LocalizeMsg(c.MustGet("lang").(*i18n.Localizer), response.ApiMsg.AuthErr),
				"data": map[string]interface{}{},
			})
			c.Abort()
			return
		}

		// 查询这个user是不是空
		var user models.UserInfo
		error := db.DB.Where("id = ?", mc.UserId).First(&user).Error
		if errors.Is(error, gorm.ErrRecordNotFound) {
			response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
			return
		}

		// 将当前请求的userId信息保存到请求的上下文c上
		c.Set("userId", mc.UserId)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if settings.Conf.App.Env != "production" {
			fmt.Printf("token: %s\n", token)
		}
		if len(token) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  internal.LocalizeMsg(c.MustGet("lang").(*i18n.Localizer), response.ApiMsg.AuthErr),
				"data": map[string]interface{}{},
			})
			c.Abort()
			return
		}
		mc, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  internal.LocalizeMsg(c.MustGet("lang").(*i18n.Localizer), response.ApiMsg.AuthErr),
				"data": map[string]interface{}{},
			})
			c.Abort()
			return
		}

		// 查询这个user是不是空
		var user models.UserInfo
		userResult := db.DB.Where("ID = ?", mc.UserId).Find(&user)
		if errors.Is(userResult.Error, gorm.ErrRecordNotFound) {
			response.Fail(c, response.ApiCode.UserNotFound, response.ApiMsg.UserNotFound)
			return
		}

		if user.Role != RoleAdmin {
			c.AbortWithStatusJSON(403, gin.H{
				"msg": "permission denied",
			})
			return
		}

		// 将当前请求的userId信息保存到请求的上下文c上
		// c.Set("userId", mc.UserId)
		c.Next()
	}
}

func OptionalJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 默认未登录
		c.Set("userId", uint(0))

		token := c.GetHeader("token")
		if token == "" {
			c.Next()
			return
		}

		// 解析 token
		mc, err := ParseToken(token)
		if err != nil {
			// token 不合法，当未登录处理
			c.Next()
			return
		}

		// 查询用户是否存在
		var user models.UserInfo
		err = db.DB.
			Select("id").
			Where("id = ?", mc.UserId).
			First(&user).Error

		if err != nil {
			// 用户不存在 / 被删除
			c.Next()
			return
		}

		// 登录态有效
		c.Set("userId", mc.UserId)
		c.Next()
	}
}
