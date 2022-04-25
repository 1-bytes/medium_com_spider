package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/olivere/elastic/v7"
	"medium_com/bootstrap"
	"medium_com/parser"
	elasticsearch "medium_com/pkg/elastic"
	"medium_com/pkg/queued"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NewCollector 传入配置信息，创建并返回一个 colly 的 collector 实例
func NewCollector(options ...colly.CollectorOption) *colly.Collector {
	c := colly.NewCollector(options...)
	//// 代理设置
	//fmt.Println(config.GetString("spider.socks5"))
	//rp, err := proxy.RoundRobinProxySwitcher(config.GetString("spider.socks5"))
	//if err != nil {
	//	log.Println("attempt to use Socks5 proxy failed.")
	//	panic(err)
	//}
	// 爬虫速度以及响应时间等参数的控制
	c.WithTransport(&http.Transport{
		//Proxy: rp,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		//DisableKeepAlives:     true,
	})
	// 初始化 Redis Storage，将其用作爬虫的持久化队列
	if err := c.SetStorage(bootstrap.Storage); err != nil {
		panic(err)
	}
	return c
}

// SpiderCallbacks colly 的回调函数
func SpiderCallbacks(c *colly.Collector) {
	// 请求发起之前要处理的一些事件
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("Referer", "https://medium.com")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	})

	// 抓取新的页面
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		// 链接中存在某些关键字的直接跳过
		skipURLKeywordsMap := []string{
			"#",
		}
		for _, keyword := range skipURLKeywordsMap {
			if strings.Index(url, keyword) != -1 {
				return
			}
		}
		//_ = e.Request.Visit(url)
		_ = queued.Queued.AddURL(url)
	})

	// 处理请求结果
	c.OnResponse(func(r *colly.Response) {
		url := r.Request.URL.String()
		domain := r.Request.URL.Host
		if strings.Index(url, ".htm") == -1 {
			return
		}
		if strings.Index(string(r.Body), "pagestyle") != -1 {
			return
		}

		body := r.Body
		//id := parser.ID(url)
		//title := parser.Title(body)
		author := parser.Author()
		category := parser.Category()
		releaseDate := parser.ReleaseDate(body)
		paragraphs, err := parser.Content(body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		if len(paragraphs) == 0 {
			//fmt.Printf("Error: request failed, title or content is nil\n")
			//	_ = r.Request.Retry()
			return
		}
		model := parser.DictArticleModel{
			//ID:                  id,
			Type: parser.TypeMap[category],
			//Title:               title,
			Author:              author,
			ReleaseDate:         releaseDate,
			MostRecentlyUpdated: "",
			SourceDomain:        domain,
			SourceUrl:           url,
		}

		err = SaveDataToMySQL("dict_article", &model)
		if err != nil {
			fmt.Printf("SaveData error: %v\n", err)
			return
		}

		for _, paragraph := range paragraphs {
			//fmt.Printf("ID: %d\n", id)
			//fmt.Printf("Title: %s\n", title)
			//fmt.Printf("Author: %s\n", author)
			//fmt.Printf("Category: %s\n", category)
			//fmt.Printf("ReleaseDate: %s\n", releaseDate)
			//fmt.Printf("EN: %s\n", paragraph["EN"])
			//fmt.Printf("CN: %s\n", paragraph["CN"])
			//fmt.Println()

			data := parser.JsonData{
				ID:           strconv.Itoa(model.ID),
				SourceDomain: domain,
				SourceURL:    url,
				Paragraph:    paragraph,
			}
			if err = SaveDataToElastic("dict_article", "", data); err != nil {
				fmt.Printf("SaveData error: %v\n", err)
			}
		}
	})

	// 错误处理
	c.OnError(func(resp *colly.Response, err error) {
		//err = resp.Request.Retry()
		err = queued.Queued.AddRequest(resp.Request)
		if err != nil {
			fmt.Println("Request URL:", resp.Request.URL, "failed with response:", resp, "\nError:", err)
		}
	})
}

// SaveDataToElastic 存储数据至 ES
func SaveDataToElastic(index string, id string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var e *elastic.IndexService
	e = elasticsearch.GetInstance().Index()
	if id != "" {
		e.Id(id)
	}
	do, err := e.Index(index).BodyJson(string(j)).Do(context.Background())
	if do != nil {
		fmt.Printf("%+v: %+v\n", do.Result, do.Id)
	}
	return err
}

// SaveDataToMySQL 存储数据至 mysql
func SaveDataToMySQL(tables string, data *parser.DictArticleModel) error {
	db := bootstrap.DB
	tx := db.Table(tables).Create(data)
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}
