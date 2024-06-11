package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// 还有一个问题，ShouldBind是动态获取userInfo的类型，并找到这个结构体里的值，要想被找到值，u结构内的值必须大写
type user struct {
	NickName string `json:"nickName" form:"nickName" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// QueryTestNetworking
/*
http://localhost:8082/v1/test/get/test?nickName=张胜男&password=126
*/
func QueryTestNetworking(c *gin.Context) {
	name := c.Query("nickName")
	password := c.DefaultQuery("password", "123")
	code, _ := c.GetQuery("code")
	data := make(map[string]interface{})
	data["name"] = name
	data["password"] = password
	data["code"] = code
	c.JSON(200, data)
}

func FormTestNetworking(c *gin.Context) {
	name := c.PostForm("nickName")
	password := c.DefaultPostForm("password", "123")
	code, _ := c.GetPostForm("code")
	data := map[string]interface{}{
		"nickName": name,
		"password": password,
		"code":     code,
	}
	c.JSON(http.StatusOK, data)
}

func PathTestNetworking(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{"name": name})
}

func BindingNetworking(c *gin.Context) {
	ip := c.ClientIP()
	log.Println("ip is", ip)
	header := c.Request.Header
	headerV := c.GetHeader("User-Agent")
	log.Println("header is ", header)
	log.Println("headerV is ", headerV)
	nickName := c.PostForm("nickName")
	password := c.DefaultPostForm("password", "123")
	log.Println("nickName is ", nickName, password)
	// 这里需要传的是值，所以加&，因为要改变的是userInfo根的值。而不是像拷贝过来一样，更改现值，原值不变
	// 还有一个问题，ShouldBind是动态获取userInfo的类型，并找到这个结构体里的值，要想被找到值，u结构内的值必须大写
	var userInfo user
	if err := c.ShouldBind(&userInfo); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("user is ", userInfo)
	c.JSON(http.StatusOK, userInfo)
}
