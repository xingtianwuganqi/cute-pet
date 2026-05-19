package internal

import "fmt"

// 辅助函数，将字符串转换为uint
func GetUintFromString(s string) uint {
	var n uint
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}