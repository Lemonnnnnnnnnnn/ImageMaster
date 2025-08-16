package crawler

import (
	"fmt"
	"net/url"
	"strings"

	"ImageMaster/core/crawler/parsers"
	"ImageMaster/core/logger"
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
	SiteTypeComic18   = "comic18"
	SiteTypeHitomi    = "hitomi"
	SiteTypeGeneric   = "generic"
)

// CrawlerFactory 爬虫工厂
type CrawlerFactory struct {
	reqClient     *request.Client
	configManager types.ConfigProvider
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

		// 从配置中获取代理设置，并直接应用到请求客户端
		if proxyURL := configManager.GetProxy(); proxyURL != "" {
			logger.Info("设置代理: %s", proxyURL)
			f.reqClient.SetProxy(proxyURL)
		}
	}
}

// CreateCrawler 创建特定网站的爬虫
func (f *CrawlerFactory) createCrawler(siteType string) types.ImageCrawler {
	logger.Info("创建爬虫, 类型: %s", siteType)

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
	case SiteTypeComic18:
		return parsers.NewComic18Crawler(f.reqClient)
	case SiteTypeHitomi:
		return parsers.NewHitomiCrawler(f.reqClient)
	default:
		logger.Warn("创建爬虫失败, 类型: %s", siteType)
		return nil
	}
}

// DetectSiteType 检测网站类型
func (f *CrawlerFactory) detectSiteType(rawURL string) string {
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

	// 检测18comic
	if strings.Contains(host, "18comic.vip") || strings.Contains(host, "18comic.org") {
		return SiteTypeComic18
	}

	// 检测Hitomi
	if strings.Contains(host, "hitomi.la") {
		return SiteTypeHitomi
	}

	// 默认使用通用爬虫
	return SiteTypeGeneric
}

func (f *CrawlerFactory) Create(rawURL string) (types.ImageCrawler, error) {
	siteType := f.detectSiteType(rawURL)
	crawler := f.createCrawler(siteType)
	if crawler == nil {
		return nil, fmt.Errorf("unsupported site type: %s", siteType)
	}
	return crawler, nil
}
