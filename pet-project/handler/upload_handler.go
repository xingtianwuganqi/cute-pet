package handler

import (
	"context"
	"pet-project/response"
	"pet-project/settings"

	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
)

// GetQiNiuToken 获取七牛 tokan
func GetQiNiuToken(c *gin.Context) {
	bucket := "petproject"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(settings.Conf.ApiKeys.QiniuAccessKey, settings.Conf.ApiKeys.QiniuSecretKey)
	upToken := putPolicy.UploadToken(mac)
	response.Success(c, gin.H{
		"token": upToken,
	})
}

func QiNiuDeleteFile(c *gin.Context) {
	key := c.Param("key")

	creds := credentials.NewCredentials(settings.Conf.ApiKeys.QiniuAccessKey, settings.Conf.ApiKeys.QiniuSecretKey)
	objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: creds},
	})
	bucketName := "petproject"
	bucket := objectsManager.Bucket(bucketName)
	err := bucket.Object(key).Delete().Call(context.Background())
	if err != nil {
		response.Fail(c, response.ApiCode.Fail, response.ApiMsg.Fail)
		return
	}
	response.Success(c, nil)
}

func DeleteQiNiuFile(key string) error {
	creds := credentials.NewCredentials(settings.Conf.ApiKeys.QiniuAccessKey, settings.Conf.ApiKeys.QiniuSecretKey)
	objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: creds},
	})
	bucketName := "petproject"
	bucket := objectsManager.Bucket(bucketName)
	err := bucket.Object(key).Delete().Call(context.Background())
	return err
}
