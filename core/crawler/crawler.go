package crawler

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"ImageMaster/core/config"
	"ImageMaster/core/crawler/parsers"
	"ImageMaster/core/downloader"
	"ImageMaster/core/request"
)

// Site types
const (
	SiteTypeEHentai   = "ehentai"
	SiteTypeExHentai  = "exhentai"
	SiteTypeTelegraph = "telegraph"
	SiteTypeGeneric   = "generic"
)

// ImageCrawler 图片爬虫接口
type ImageCrawler interface {
	Crawl(url string, saveDir string) (string, error)
	GetDownloader() *downloader.Downloader
}

// CrawlerFactory 爬虫工厂
type CrawlerFactory struct {
	reqClient     *request.Client
	ctx           context.Context
	configManager *config.Manager
}

// NewCrawlerFactory 创建爬虫工厂
func NewCrawlerFactory(ctx context.Context) *CrawlerFactory {
	return &CrawlerFactory{
		reqClient: request.NewClient(),
		ctx:       ctx,
	}
}

// SetConfigManager 设置配置管理器
func (f *CrawlerFactory) SetConfigManager(configManager *config.Manager) {
	f.configManager = configManager

	// 如果配置管理器不为空，设置到请求客户端
	if configManager != nil {
		f.reqClient.SetConfigManager(configManager)
	}
}

// CreateCrawler 创建特定网站的爬虫
func (f *CrawlerFactory) CreateCrawler(siteType string) ImageCrawler {
	fmt.Printf("创建爬虫, 类型: %s\n", siteType)

	// 为每个爬虫创建一个下载器
	dl := downloader.NewDownloader(3, 3, true)

	// 如果有配置管理器，设置到下载器
	if f.configManager != nil {
		dl.SetConfigManager(f.configManager)
	}

	switch siteType {
	case SiteTypeEHentai:
		crawler := &EHentaiCrawler{
			reqClient:  f.reqClient,
			ctx:        f.ctx,
			downloader: dl,
		}
		return crawler
	case SiteTypeExHentai:
		crawler := &EHentaiCrawler{
			reqClient:  f.reqClient,
			ctx:        f.ctx,
			downloader: dl,
		}
		return crawler
	case SiteTypeTelegraph:
		crawler := &TelegraphCrawler{
			reqClient:  f.reqClient,
			downloader: dl,
		}
		return crawler
	default:
		crawler := &GenericCrawler{
			reqClient:  f.reqClient,
			ctx:        f.ctx,
			downloader: dl,
		}
		return crawler
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
	if strings.Contains(host, "e-hentai.org") {
		return SiteTypeEHentai
	}

	// 检测ExHentai
	if strings.Contains(host, "exhentai.org") {
		return SiteTypeExHentai
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
	reqClient  *request.Client
	ctx        context.Context
	downloader *downloader.Downloader
}

// GetDownloader 获取下载器
func (c *GenericCrawler) GetDownloader() *downloader.Downloader {
	return c.downloader
}

// Crawl 执行爬取
func (c *GenericCrawler) Crawl(url string, savePath string) (string, error) {
	// 使用 RemoteGetter 的通用爬虫实现
	// 这里需要实现
	return "", fmt.Errorf("通用爬虫尚未实现")
}

// EHentaiCrawler E-Hentai爬虫
type EHentaiCrawler struct {
	reqClient  *request.Client
	ctx        context.Context
	downloader *downloader.Downloader
}

// GetDownloader 获取下载器
func (c *EHentaiCrawler) GetDownloader() *downloader.Downloader {
	return c.downloader
}

// Crawl 执行爬取
func (c *EHentaiCrawler) Crawl(url string, savePath string) (string, error) {
	// 获取Client用于兼容
	client := c.reqClient.GetHTTPClient()

	// 将下载器传递给解析器
	err := parsers.ParseEHentai(c.ctx, client, url, savePath, c.downloader)
	if err != nil {
		return "", err
	}
	return savePath, nil
}

// TelegraphCrawler Telegraph爬虫
type TelegraphCrawler struct {
	reqClient  *request.Client
	downloader *downloader.Downloader
}

// GetDownloader 获取下载器
func (c *TelegraphCrawler) GetDownloader() *downloader.Downloader {
	return c.downloader
}

// Crawl 执行爬取
func (c *TelegraphCrawler) Crawl(url string, savePath string) (string, error) {
	// 获取Client用于兼容
	client := c.reqClient.GetHTTPClient()

	err := parsers.ParseTelegraph(client, url, savePath, c.downloader)
	if err != nil {
		return "", err
	}
	return savePath, nil
}
