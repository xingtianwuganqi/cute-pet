package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// QueryTestNetworking
/*
http://localhost:8082/v1/test/get/test?nickName=张胜男&password=126
*/
func QueryTestNetworking(c *gin.Context) {
	name := c.Query("nickName")
	password := c.DefaultQuery("password", "123")
	data := make(map[string]interface{})
	data["name"] = name
	data["password"] = password
	c.JSON(200, data)
}

func FormTestNetworking(c *gin.Context) {
	name := c.PostForm("nickName")
	password := c.DefaultPostForm("password", "123")
	data := map[string]interface{}{
		"nickName": name,
		"password": password,
	}
	c.JSON(http.StatusOK, data)
}
