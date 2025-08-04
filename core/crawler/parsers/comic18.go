package parsers

import (
	"fmt"
	"net/http"
	"path"

	"github.com/PuerkitoBio/goquery"

	"ImageMaster/core/request"
	"ImageMaster/core/types"
)

// Comic18Parser 18comic解析器实现
type Comic18Parser struct{}

// GetName 获取解析器名称
func (p *Comic18Parser) GetName() string {
	return "18Comic"
}

// Parse 解析URL获取图片信息
func (p *Comic18Parser) Parse(reqClient *request.Client, url string) (*ParseResult, error) {
	resp, err := reqClient.Get(url)
	if err != nil {
		fmt.Println("18comic解析器获取URL失败", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("18comic解析器解析HTML失败", err)
		return nil, err
	}

	// 获取专辑名称
	albumName := ""
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		albumName = s.Text()
	})

	if albumName == "" {
		albumName = "18Comic Album" // 默认名称
	}

	// 获取所有图片
	var imgURLs []string
	var filePaths []string
	doc.Find(".scramble-page > img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			imgURLs = append(imgURLs, src)

			// 从 src 中提取文件扩展名
			ext := path.Ext(src)
			if ext == "" {
				ext = ".webp" // 默认扩展名
			}
			filePaths = append(filePaths, fmt.Sprintf("%d%s", i, ext))
		}
	})

	return &ParseResult{
		Name:      albumName,
		ImageURLs: imgURLs,
		FilePaths: filePaths,
	}, nil
}

// Comic18Crawler 18comic爬虫
type Comic18Crawler struct {
	*BaseCrawler
}

// NewComic18Crawler 创建新的18comic爬虫
func NewComic18Crawler(reqClient *request.Client) types.ImageCrawler {
	parser := &Comic18Parser{}
	baseCrawler := NewBaseCrawler(reqClient, parser)
	return &Comic18Crawler{
		BaseCrawler: baseCrawler,
	}
}
