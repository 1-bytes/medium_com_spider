package parser

import (
	"medium_com/pkg/utils"
	"regexp"
	"strconv"
	"strings"
)

// ID 解析 ID
func ID(u string) int {
	idRe := regexp.MustCompile(`/(\d+)\.html`)
	match := idRe.FindAllStringSubmatch(u, -1)
	if len(match) > 0 {
		id, _ := strconv.Atoi(match[0][1])
		return id
	}
	return -1
}

// Author 作者
func Author() string {
	return "admin"
}

// Category 解析分类
func Category() string {
	return "blog"
}

// ReleaseDate 发布日期
func ReleaseDate(body []byte) string {
	info := utils.GetBetweenStr(string(body), `<span class="info_l">`, `</span>`)
	infoSplit := strings.Split(info, "|")
	if len(infoSplit) >= 3 {
		date := strings.TrimSpace(infoSplit[2])
		date = strings.TrimSpace(utils.GetBetweenStr(date, "Updated: ", "\n"))
		return date
	}
	return ""
}
