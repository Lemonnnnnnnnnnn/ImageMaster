package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ImageMaster/core/crawler/parsers"
)

const (
	SiteTypeGeneric   = "generic"
	SiteTypeEHentai   = "ehentai"
	SiteTypeTelegraph = "telegraph"
)

// ImageCrawler 图片爬虫接口
type ImageCrawler interface {
	Crawl(url string, savePath string) (string, error)
}

// CrawlerFactory 爬虫工厂
type CrawlerFactory struct {
	client *http.Client
	ctx    context.Context
}

// NewCrawlerFactory 创建爬虫工厂
func NewCrawlerFactory(ctx context.Context) *CrawlerFactory {
	return &CrawlerFactory{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		ctx: ctx,
	}
}

// CreateCrawler 创建特定网站的爬虫
func (f *CrawlerFactory) CreateCrawler(siteType string) ImageCrawler {
	switch siteType {
	case SiteTypeEHentai:
		return &EHentaiCrawler{
			client: f.client,
			ctx:    f.ctx,
		}
	case SiteTypeTelegraph:
		return &TelegraphCrawler{
			client: f.client,
		}
	default:
		return &GenericCrawler{
			client: f.client,
			ctx:    f.ctx,
		}
	}
}

// DetectSiteType 检测网站类型
func (f *CrawlerFactory) DetectSiteType(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return SiteTypeGeneric
	}

	host := parsedURL.Host

	// 检测E-Hentai
	if strings.Contains(host, "e-hentai.org") || strings.Contains(host, "exhentai.org") {
		return SiteTypeEHentai
	}

	// 检测Telegraph
	if strings.Contains(host, "telegra.ph") || strings.Contains(host, "telegraph.com") {
		return SiteTypeTelegraph
	}

	// 默认使用通用爬虫
	return SiteTypeGeneric
}

// GenericCrawler 通用网页爬虫
type GenericCrawler struct {
	client *http.Client
	ctx    context.Context
}

// Crawl 执行爬取
func (c *GenericCrawler) Crawl(url string, savePath string) (string, error) {
	// 使用 RemoteGetter 的通用爬虫实现
	// 这里需要实现
	return "", fmt.Errorf("通用爬虫尚未实现")
}

// EHentaiCrawler E-Hentai爬虫
type EHentaiCrawler struct {
	client *http.Client
	ctx    context.Context
}

// Crawl 执行爬取
func (c *EHentaiCrawler) Crawl(url string, savePath string) (string, error) {
	err := parsers.ParseEHentai(c.ctx, c.client, url, savePath)
	if err != nil {
		return "", err
	}
	return savePath, nil
}

// TelegraphCrawler Telegraph爬虫
type TelegraphCrawler struct {
	client *http.Client
}

// Crawl 执行爬取
func (c *TelegraphCrawler) Crawl(url string, savePath string) (string, error) {
	err := parsers.ParseTelegraph(c.client, url, savePath)
	if err != nil {
		return "", err
	}
	return savePath, nil
}
