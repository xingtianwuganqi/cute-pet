package handler

import (
	"context"
	"pet-project/response"
	"pet-project/settings"
	"pet-project/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
)

// GetQiNiuToken 获取七牛 tokan
func GetQiNiuToken(c *gin.Context) {
	logger.Logger.Info("GetQiNiuToken handler called", zap.String("clientIP", c.ClientIP()))

	bucket := "petproject"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(settings.Conf.ApiKeys.QiniuAccessKey, settings.Conf.ApiKeys.QiniuSecretKey)
	upToken := putPolicy.UploadToken(mac)
	logger.Logger.Info("Successfully generated Qiniu token")
	response.Success(c, gin.H{
		"token": upToken,
	})
}

func QiNiuDeleteFile(c *gin.Context) {
	logger.Logger.Info("QiNiuDeleteFile handler called", zap.String("clientIP", c.ClientIP()))

	key := c.Param("key")

	creds := credentials.NewCredentials(settings.Conf.ApiKeys.QiniuAccessKey, settings.Conf.ApiKeys.QiniuSecretKey)
	objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: creds},
	})
	bucketName := "petproject"
	bucket := objectsManager.Bucket(bucketName)
	err := bucket.Object(key).Delete().Call(context.Background())
	if err != nil {
		logger.Logger.Error("Failed to delete Qiniu file", zap.Error(err), zap.String("key", key))
		response.Fail(c, response.ApiCode.Fail, response.ApiMsg.Fail)
		return
	}
	logger.Logger.Info("Successfully deleted Qiniu file", zap.String("key", key))
	response.Success(c, nil)
}