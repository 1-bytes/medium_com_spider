package parser

// Title 用于解析页面中的标题
func Title(items Paragraphs) string {
	for _, item := range items {
		if item.Type != "H3" {
			continue
		}
		return item.Text
	}
	return ""
}
