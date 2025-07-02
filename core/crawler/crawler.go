package crawler

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"ImageMaster/core/crawler/parsers"
	"ImageMaster/core/proxy"
	"ImageMaster/core/request"
	"ImageMaster/core/types"
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
	CrawlAndSave(url string, savePath string) string
	GetDownloader() types.Downloader
	SetDownloader(dl types.Downloader)
}

// CrawlerFactory 爬虫工厂
type CrawlerFactory struct {
	reqClient     *request.Client
	ctx           context.Context
	configManager types.ConfigProvider
	proxyManager  *proxy.ProxyManager
}

// NewCrawlerFactory 创建爬虫工厂
func NewCrawlerFactory(ctx context.Context) *CrawlerFactory {
	return &CrawlerFactory{
		reqClient: request.NewClient(),
		ctx:       ctx,
	}
}

// SetConfigManager 设置配置管理器
func (f *CrawlerFactory) SetConfigManager(configManager types.ConfigProvider) {
	f.configManager = configManager

	// 如果配置管理器不为空，设置到请求客户端
	if configManager != nil {
		f.reqClient.SetConfigManager(configManager)

		// 创建代理管理器
		f.proxyManager = proxy.NewProxyManager(configManager)

		// 从配置中获取代理设置，并直接应用到请求客户端
		if proxyURL := configManager.GetProxy(); proxyURL != "" {
			fmt.Printf("设置代理: %s\n", proxyURL)
			f.reqClient.SetProxy(proxyURL)
		}
	}
}

// CreateCrawler 创建特定网站的爬虫
func (f *CrawlerFactory) CreateCrawler(siteType string) ImageCrawler {
	fmt.Printf("创建爬虫, 类型: %s\n", siteType)

	// 所有爬虫共用同一个配置好的请求客户端
	switch siteType {
	case SiteTypeEHentai:
		crawler := &EHentaiCrawler{
			reqClient: f.reqClient,
			ctx:       f.ctx,
		}
		return crawler
	case SiteTypeExHentai:
		crawler := &EHentaiCrawler{
			reqClient: f.reqClient,
			ctx:       f.ctx,
		}
		return crawler
	case SiteTypeTelegraph:
		crawler := &TelegraphCrawler{
			reqClient: f.reqClient,
		}
		return crawler
	default:
		crawler := &GenericCrawler{
			reqClient: f.reqClient,
			ctx:       f.ctx,
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
	downloader types.Downloader
}

// GetDownloader 获取下载器
func (c *GenericCrawler) GetDownloader() types.Downloader {
	return c.downloader
}

// SetDownloader 设置下载器
func (c *GenericCrawler) SetDownloader(dl types.Downloader) {
	c.downloader = dl
}

// Crawl 执行爬取
func (c *GenericCrawler) Crawl(url string, savePath string) (string, error) {
	// 通用爬虫暂未实现
	return "", fmt.Errorf("通用爬虫尚未实现")
}

// CrawlAndSave 执行爬取并保存
func (c *GenericCrawler) CrawlAndSave(url string, savePath string) string {
	_, err := c.Crawl(url, savePath)
	if err != nil {
		return ""
	}
	return savePath
}

// EHentaiCrawler E-Hentai爬虫
type EHentaiCrawler struct {
	reqClient  *request.Client
	ctx        context.Context
	downloader types.Downloader
}

// GetDownloader 获取下载器
func (c *EHentaiCrawler) GetDownloader() types.Downloader {
	return c.downloader
}

// SetDownloader 设置下载器
func (c *EHentaiCrawler) SetDownloader(dl types.Downloader) {
	c.downloader = dl
}

// Crawl 执行爬取
func (c *EHentaiCrawler) Crawl(url string, savePath string) (string, error) {
	// 将下载器传递给解析器，解析器会使用downloader获取的代理设置
	err := parsers.ParseEHentai(c.ctx, c.reqClient, url, savePath, c.downloader)
	if err != nil {
		return "", err
	}
	return savePath, nil
}

// CrawlAndSave 执行爬取并保存
func (c *EHentaiCrawler) CrawlAndSave(url string, savePath string) string {
	// 从URL中提取标题作为文件夹名
	name := filepath.Base(savePath)
	if name == "" || name == "." {
		// 如果无法从路径中提取有效的名称，使用"download"作为默认名称
		name = "download"
	}

	result, err := c.Crawl(url, savePath)
	if err != nil {
		fmt.Printf("爬取失败: %v\n", err)
		return ""
	}

	return result
}

// TelegraphCrawler Telegraph爬虫
type TelegraphCrawler struct {
	reqClient  *request.Client
	downloader types.Downloader
}

// GetDownloader 获取下载器
func (c *TelegraphCrawler) GetDownloader() types.Downloader {
	return c.downloader
}

// SetDownloader 设置下载器
func (c *TelegraphCrawler) SetDownloader(dl types.Downloader) {
	c.downloader = dl
}

// Crawl 执行爬取
func (c *TelegraphCrawler) Crawl(url string, savePath string) (string, error) {
	// 将下载器传递给解析器，解析器会使用downloader获取的代理设置
	err := parsers.ParseTelegraph(nil, url, savePath, c.downloader)
	if err != nil {
		return "", err
	}
	return savePath, nil
}

// CrawlAndSave 执行爬取并保存
func (c *TelegraphCrawler) CrawlAndSave(url string, savePath string) string {
	result, err := c.Crawl(url, savePath)
	if err != nil {
		return ""
	}
	return result
}
