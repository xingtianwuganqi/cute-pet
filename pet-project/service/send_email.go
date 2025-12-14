package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/gomail.v2"
	"math/rand"
	"pet-project/db"
	"time"
)

// GenerateValidationCode 生成指定长度的随机数字验证码
func GenerateValidationCode(length int) string {
	var code string
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(9))
	}
	return code
}

// SendEmail 发送验证码
func SendEmail(recipient string,
	subject string,
	body string,
	smtpServier string,
	smtpPort int,
	username string,
	password string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpServier, smtpPort, username, password)
	d.SSL = true //使用ssl连接

	return d.DialAndSend(m)
}

// SaveAccountCodeInRedis 保存到redis
func SaveAccountCodeInRedis(c *gin.Context, email, code string, expiration time.Duration) error {
	key := fmt.Sprintf("verify_code:%s", email)
	if err := db.Rdb.Set(c, key, code, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func GetCodeFromRedis(c *gin.Context, email string) (string, error) {
	key := fmt.Sprintf("verify_code:%s", email)

	value, err := db.Rdb.Get(c, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

// DeleteCodeFromRedis redis删除数据
func DeleteCodeFromRedis(c *gin.Context, key string) error {
	err := db.Rdb.Del(c, key).Err()
	return err
}
