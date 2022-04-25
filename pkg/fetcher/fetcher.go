package fetcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// GetArticleParagraphs 根据文章 ID 获取文章的段落详情
func GetArticleParagraphs(articleID string) ([]byte, error) {
	//type Graphql
	graphql := []struct {
		OperationName string `json:"operationName"`
		Variables     struct {
			PostID string `json:"postId"`
		} `json:"variables"`
		Query string `json:"query"`
	}{
		{
			OperationName: "PostViewerEdgeContentQuery",
			Query: `query PostViewerEdgeContentQuery($postId: ID!$postMeteringOptions:PostMeteringOptions)
{post(id: $postId) {...on Post {id viewerEdge {id fullContent(postMeteringOptions: $postMeteringOptions) 
{bodyModel {...PostBody_bodyModel}}}}}}fragment PostBody_bodyModel on RichText {paragraphs {id name type text}}`,
			Variables: struct {
				PostID string `json:"postId"`
			}{
				PostID: articleID,
			},
		},
	}

	body, err := json.Marshal(graphql)
	if err != nil {
		return nil, err
	}

	header := &http.Header{}
	header.Set("content-type", "application/json")
	respJson, err := Fetch(http.MethodPost, "https://cryptowhale.medium.com/_/graphql",
		strings.NewReader(string(body)), true, header)
	if err != nil {
		return nil, err
	}
	return respJson, nil
}

// Fetch 用于获取网页内容
func Fetch(method, u string, body io.Reader, useProxy bool, header *http.Header) ([]byte, error) {
	// 创建客户端
	client := &http.Client{}
	// 设置代理
	if useProxy {
		proxy := "socks5://127.0.0.1:1080"
		p := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
		client.Transport = &http.Transport{
			Proxy: p,
		}
	}
	// 创建 request 请求
	request, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	// 增加 header
	header.Set("Referer", "https://medium.com/")
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	request.Header = *header
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("return status code is not expected (%d)", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return ioutil.ReadAll(utf8Reader)
}

// determineEncoding 自动判断网页编码.
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
