package response

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
	"pet-project/service"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code uint        `json:"code" form:"code"`
	Msg  string      `json:"msg" form:"msg"`
	Data interface{} `json:"data" form:"data"`
}

func Response(c *gin.Context, code uint, data interface{}, msg string) {
	if data == nil {
		data = gin.H{}
	}
	res := BaseResponse{}
	res.Code = code
	res.Msg = msg
	res.Data = data

	c.JSON(http.StatusOK, res)
}

// Success 成功
func Success(c *gin.Context, data interface{}) {
	Response(c, 200, data, "success")
}

// Fail 出错
func Fail(c *gin.Context, code uint, msg string) {
	Response(c, code, gin.H{}, msg)
}

func FailMsg(c *gin.Context, code uint, msg string) {
	lang := c.MustGet("lang").(*i18n.Localizer)
	message := service.LocalizeMsg(lang, msg)
	Response(c, code, gin.H{}, message)
}
