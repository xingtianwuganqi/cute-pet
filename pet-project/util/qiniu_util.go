package util

import (
	"context"
	"pet-project/settings"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
)

// DeleteQiNiuFile 删除七牛云文件
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