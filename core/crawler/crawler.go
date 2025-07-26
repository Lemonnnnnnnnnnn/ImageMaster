package crawler

import (
	"fmt"
	"net/url"
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
	SiteTypeWnacg     = "wnacg"
	SiteTypeNhentai   = "nhentai"
	SiteTypeGeneric   = "generic"
)

// CrawlerFactory 爬虫工厂
type CrawlerFactory struct {
	reqClient     *request.Client
	configManager types.ConfigProvider
	proxyManager  *proxy.ProxyManager
}

// NewCrawlerFactory 创建爬虫工厂
func NewCrawlerFactory() *CrawlerFactory {
	return &CrawlerFactory{
		reqClient: request.NewClient(),
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
func (f *CrawlerFactory) CreateCrawler(siteType string) types.ImageCrawler {
	fmt.Printf("创建爬虫, 类型: %s\n", siteType)

	// 所有爬虫共用同一个配置好的请求客户端
	switch siteType {
	case SiteTypeEHentai:
		return parsers.NewEHentaiCrawler(f.reqClient)
	case SiteTypeExHentai:
		return parsers.NewEHentaiCrawler(f.reqClient)
	case SiteTypeTelegraph:
		return parsers.NewTelegraphCrawler(f.reqClient)
	case SiteTypeWnacg:
		return parsers.NewWnacgCrawler(f.reqClient)
	case SiteTypeNhentai:
		return parsers.NewNhentaiCrawler(f.reqClient)
	default:
		crawler := &GenericCrawler{
			reqClient: f.reqClient,
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

	// 检测Wnacg
	if strings.Contains(host, "wnacg.com") {
		return SiteTypeWnacg
	}

	// 检测Nhentai
	if strings.Contains(host, "nhentai.xxx") {
		return SiteTypeNhentai
	}

	// 默认使用通用爬虫
	return SiteTypeGeneric
}

// GenericCrawler 通用网页爬虫
type GenericCrawler struct {
	reqClient  *request.Client
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
