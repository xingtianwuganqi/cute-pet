package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"pet-project/response"
	"pet-project/settings"
)

// GetQiNiuToken 获取七牛 tokan
func GetQiNiuToken(c *gin.Context) {
	bucket := ""
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(settings.Conf.ApiKeys.QiniuAccessKey, settings.Conf.ApiKeys.QiniuSecretKey)
	upToken := putPolicy.UploadToken(mac)
	response.Success(c, gin.H{
		"token": upToken,
	})
}
