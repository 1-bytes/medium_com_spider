package utils

import "strings"

// GetBetweenStr 截取中间字符串
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		return ""
	}
	n += len(start)
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		return ""
	}
	m += len(end)
	return string([]byte(str)[:m])
}
