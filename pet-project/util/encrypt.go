package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var key = []byte("thisis16bytekey!") // 刚好16字节

func Md5String(s string) string {
	timeStr := string(time.Now().UnixNano())
	timeStr += s
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}

func IsValidEmail(email string) bool {
	// 正则表达式匹配邮箱格式
	// 匹配常见的邮箱格式，不保证能匹配所有有效的邮箱
	// 匹配a@b.com这样的格式，但不匹配a@b这样的格式
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegexp.MatchString(email)
}

// GenerateDynamicParam 生成包含时间戳的动态参数
func GenerateDynamicParam() string {
	timestamp := time.Now().Unix()
	randomString := "random_string" // 可以使用更多的随机生成逻辑
	return fmt.Sprintf("%d|%s", timestamp, randomString)
}

// Encrypt 加密参数
func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	b := base64.StdEncoding.EncodeToString([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密参数
func Decrypt(cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	data, _ := base64.StdEncoding.DecodeString(string(ciphertext))
	return string(data), nil
}

// IsValidParam isValidParam 验证解密后的参数是否有效
func IsValidParam(decryptedParam string, maxAge time.Duration) bool {
	parts := strings.Split(decryptedParam, "|")
	if len(parts) != 2 {
		return false
	}

	timestampStr, randomString := parts[0], parts[1]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	if time.Since(time.Unix(timestamp, 0)) > maxAge {
		return false
	}

	// 检查 randomString 是否已被使用，可以维护一个已使用的随机字符串的列表
	if randomString != "random_string" {
		return false
	}
	return true
}
