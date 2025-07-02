package crawler

import (
	"context"
	"net/url"
	"strings"

	"ImageMaster/core/crawler"
	"ImageMaster/core/types"
)

// SimpleCrawler 简化的爬虫接口
type SimpleCrawler struct {
	factory       *crawler.CrawlerFactory
	configManager types.ConfigProvider
	downloader    types.Downloader
}

// CrawlResult 爬取结果
type CrawlResult struct {
	Title    string
	SavePath string
	Error    error
}

// NewSimpleCrawler 创建简化爬虫
func NewSimpleCrawler() *SimpleCrawler {
	return &SimpleCrawler{
		factory: crawler.NewCrawlerFactory(context.Background()),
	}
}

// SetDownloader 设置下载器
func (sc *SimpleCrawler) SetDownloader(downloader types.Downloader) {
	// 保存下载器引用
	sc.downloader = downloader

	// 从下载器获取配置管理器并设置到爬虫工厂
	if configProvider, ok := downloader.(interface{ GetConfigManager() interface{} }); ok {
		if configManager, ok := configProvider.GetConfigManager().(types.ConfigProvider); ok {
			sc.configManager = configManager
			sc.factory.SetConfigManager(configManager)
		}
	}
}

// CrawlImages 爬取图片
func (sc *SimpleCrawler) CrawlImages(targetURL, outputDir string) (*CrawlResult, error) {
	// 检测网站类型
	siteType := sc.detectSiteType(targetURL)

	// 创建对应的爬虫
	crawlerInstance := sc.factory.CreateCrawler(siteType)

	// 将下载器设置给爬虫实例
	if sc.downloader != nil {
		crawlerInstance.SetDownloader(sc.downloader)
	}

	// 执行爬取
	savePath, err := crawlerInstance.Crawl(targetURL, outputDir)
	if err != nil {
		return &CrawlResult{
			Error: err,
		}, err
	}

	return &CrawlResult{
		Title:    targetURL,
		SavePath: savePath,
	}, nil
}

// detectSiteType 检测网站类型
func (sc *SimpleCrawler) detectSiteType(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return crawler.SiteTypeGeneric
	}

	host := parsedURL.Host

	// 检测E-Hentai
	if strings.Contains(host, "e-hentai.org") {
		return crawler.SiteTypeEHentai
	}

	// 检测ExHentai
	if strings.Contains(host, "exhentai.org") {
		return crawler.SiteTypeExHentai
	}

	// 检测Telegraph
	if strings.Contains(host, "telegra.ph") || strings.Contains(host, "telegraph.com") {
		return crawler.SiteTypeTelegraph
	}

	// 默认使用通用爬虫
	return crawler.SiteTypeGeneric
}
