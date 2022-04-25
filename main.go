package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"medium_com/bootstrap"
	"medium_com/cmd"
	"medium_com/parser"
	"medium_com/pkg/config"
	"medium_com/pkg/fetcher"
	"medium_com/pkg/queued"
)

func main() {
	urls := []string{
		"https://medium.com/jobs-at-medium/work-at-medium-959d1a85284e",
		"https://medium.com/@thebullzone/how-to-trade-using-the-concept-of-bollinger-bands-93e55fe4261d",
		"https://medium.com/illumination/best-data-science-tools-automation-analytics-and-visualisation-c5fef67d6140",
		"https://medium.com/the-straight-dope/a-small-request-of-my-loyal-audience-4fac9f8a6ddb",
	}
	bootstrap.Setup()
	c := cmd.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.Async(config.GetBool("spider.async", false)),
		colly.AllowedDomains("medium.com"),
		colly.DetectCharset(),
	)
	cmd.SpiderCallbacks(c)

	for _, url := range urls {
		_ = queued.Queued.AddURL(url)
	}
	_ = queued.Queued.Run(c)
	//testCase()
}

// testCase 测试用例
func testCase() {
	bootstrap.Setup()
	ArticleParagraphsJson, err := fetcher.GetArticleParagraphs("959d1a85284e")
	if err != nil {
		panic(err)
	}
	var ArticleParagraphs []parser.ArticleParagraphs
	err = json.Unmarshal(ArticleParagraphsJson, &ArticleParagraphs)
	if err != nil {
		panic(err)
	}

	items := ArticleParagraphs[0].Data.Post.ViewerEdge.FullContent.BodyModel.Paragraphs
	title := parser.Title(items)
	fmt.Println(title)
	author := parser.Author()
	category := parser.Category()
	//releaseDate := parser.ReleaseDate(bytes)
	paragraphs := parser.Content(items)

	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Author: %s\n", author)
	fmt.Printf("Category: %s\n", category)
	for _, paragraph := range paragraphs {
		//fmt.Printf("ReleaseDate: %s\n", releaseDate)
		fmt.Printf("paragraph: %s\n", paragraph)

		//data := parser.JsonData{
		//	//ID:        strconv.Itoa(id),
		//	SourceURL: url,
		//	Paragraph: paragraph,
		//}
		//if err = cmd.SaveDataToElastic("medium_com", "", data); err != nil {
		//	fmt.Printf("SaveData error: %v\n", err)
		//}
	}
}
