package util

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"time"
)

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
