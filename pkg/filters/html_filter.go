package filters

import (
	"regexp"
)

var (
	reMap = map[string]string{
		`</p>`:    ` `,
		`<[^>]+>`: ` `,
	}
)

// HtmlFilter 用于将信息中的 html 标签过滤掉
func HtmlFilter(content []byte) []byte {
	for re, value := range reMap {
		reg := regexp.MustCompile(re)
		content = reg.ReplaceAll(content, []byte(value))
	}
	return content
}
