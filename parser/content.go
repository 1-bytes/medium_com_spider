package parser

import (
	"strings"
	"unicode"
)

type content struct {
	url    string
	title  string
	author string
	date   string
	text   string
}

type paragraph map[string]string

// Content 解析正文
func Content(items Paragraphs) []string {
	var paragraphs []string
	for _, item := range items {
		if item.Type != "P" && item.Type != "p" {
			continue
		}
		paragraphs = append(paragraphs, strings.TrimSpace(item.Text))
	}
	return paragraphs
}

// hasChinese 判断字符串中是否包含中文
func hasChinese(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) {
			return true
		}
	}
	return false
}
