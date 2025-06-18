package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var key = []byte("thisis16bytekey!") // 刚好16字节

// UuidString 生成uuid
func UuidString() string {
	uuidWithHyphen := uuid.New().String()
	uuidWithoutHyphen := strings.ReplaceAll(uuidWithHyphen, "-", "")
	return uuidWithoutHyphen
}

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

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, pad...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:(length - unpadding)], nil
}

// Encrypt encrypts plain text using AES-CBC with PKCS7 padding
func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plainData := pkcs7Padding([]byte(plainText), aes.BlockSize)

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(plainData))
	mode.CryptBlocks(cipherText, plainData)

	final := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(final), nil
}

// Decrypt decrypts base64-encoded AES-CBC encrypted text
func Decrypt(cipherTextBase64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	cipherData := data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherData, cipherData)

	plainData, err := pkcs7UnPadding(cipherData)
	if err != nil {
		return "", err
	}

	return string(plainData), nil
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
