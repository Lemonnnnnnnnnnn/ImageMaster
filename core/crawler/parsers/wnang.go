package parsers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"

	"ImageMaster/core/request"
	"ImageMaster/core/types"
)

// WnacgAlbum Wnacg专辑
type WnacgAlbum struct {
	Name  string
	Pages []string // 存储所有分页的URL
}

// ParseWnacg 解析Wnacg网站
func ParseWnacg(ctx context.Context, reqClient *request.Client, url string, savePath string, dl types.Downloader) error {
	fmt.Printf("下载 Wnacg 专辑: %s\n", url)

	// 使用下载器的代理配置
	if dl != nil && dl.GetProxy() != "" {
		reqClient.SetProxy(dl.GetProxy())
	}

	// 设置User-Agent
	reqClient.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	wnacgAlbum, err := GetWnacgAlbumWithClient(reqClient, url)
	if err != nil {
		return fmt.Errorf("获取专辑失败: %w", err)
	}

	// 使用TaskUpdater更新任务名称（如果可用）
	if dl != nil {
		if taskUpdater := dl.GetTaskUpdater(); taskUpdater != nil {
			taskUpdater.UpdateTaskName(wnacgAlbum.Name)
			taskUpdater.UpdateTaskStatus(string(types.StatusParsing), "")
			fmt.Printf("已更新任务名称为: %s\n", wnacgAlbum.Name)
		}
	}

	// 使用传入的下载器
	var localDownloader types.Downloader
	if dl != nil {
		localDownloader = dl
		fmt.Printf("Wnacg解析器使用传入的下载器\n")
	} else {
		// 未提供下载器，返回错误
		return fmt.Errorf("未提供下载器")
	}

	// 保存路径
	albumPath := savePath + "/" + wnacgAlbum.Name

	// 批量下载URL和路径
	var imgURLs []string
	var filePaths []string
	var urlMutex sync.Mutex

	// 收集所有漫画页面链接
	var allMangaLinks []string
	for pageIndex, pageURL := range wnacgAlbum.Pages {
		links, err := GetMangaLinksFromPage(reqClient, pageURL)
		if err != nil {
			fmt.Printf("获取分页 %s 的漫画链接失败: %v\n", pageURL, err)
			continue
		}

		// 为每个链接添加页面索引信息
		for linkIndex, link := range links {
			allMangaLinks = append(allMangaLinks, fmt.Sprintf("%d_%d|%s", pageIndex, linkIndex, link))
		}
	}

	totalMangaLinks := len(allMangaLinks)
	fmt.Printf("总共需要处理 %d 个漫画页面\n", totalMangaLinks)

	// 更新解析进度
	if dl != nil {
		if taskUpdater := dl.GetTaskUpdater(); taskUpdater != nil {
			taskUpdater.UpdateTaskProgress(0, totalMangaLinks)
		}
	}

	// 并发处理控制
	var wg sync.WaitGroup
	processedCount := 0
	var progressMutex sync.Mutex

	// 并发处理所有漫画页面链接
	for _, mangaData := range allMangaLinks {
		// 解析页面索引和URL
		parts := strings.SplitN(mangaData, "|", 2)
		if len(parts) != 2 {
			continue
		}
		indexPart := parts[0]
		mangaURL := parts[1]

		wg.Add(1)

		go func(indexPart, mangaURL string) {
			defer wg.Done()

			// 解析漫画页面获取真实图片URL
			imgURL, err := ParseWnacgPageWithClient(reqClient, mangaURL)
			if err != nil {
				fmt.Printf("解析漫画页面失败 %s: %v\n", mangaURL, err)
				return
			}

			// 构建保存文件名
			filename := fmt.Sprintf("%s.jpg", indexPart)
			fullPath := fmt.Sprintf("%s/%s", albumPath, filename)

			urlMutex.Lock()
			imgURLs = append(imgURLs, imgURL)
			filePaths = append(filePaths, fullPath)
			urlMutex.Unlock()

			// 更新进度
			progressMutex.Lock()
			processedCount++
			if dl != nil {
				if taskUpdater := dl.GetTaskUpdater(); taskUpdater != nil {
					taskUpdater.UpdateTaskProgress(processedCount, totalMangaLinks)
				}
			}
			progressMutex.Unlock()

			fmt.Printf("解析完成 %s\n", filename)
		}(indexPart, mangaURL)
	}

	// 等待所有URL收集任务完成
	wg.Wait()

	// 计算总图片数量
	totalImages := len(imgURLs)
	fmt.Printf("已收集 %d 张图片URL，开始下载...\n", totalImages)

	// 更新任务状态为下载中
	if dl != nil {
		if taskUpdater := dl.GetTaskUpdater(); taskUpdater != nil {
			taskUpdater.UpdateTaskStatus(string(types.StatusDownloading), "")
			taskUpdater.UpdateTaskProgress(0, totalImages)
		}
	}

	// 批量下载所有图片
	successImages, err := localDownloader.BatchDownload(imgURLs, filePaths, nil)
	if err != nil {
		fmt.Printf("批量下载出错: %v\n", err)
		return fmt.Errorf("批量下载出错: %w", err)
	}

	fmt.Printf("下载完成，总共 %d 张图片，成功 %d 张\n", totalImages, successImages)

	// 更新最终状态
	if dl != nil {
		if taskUpdater := dl.GetTaskUpdater(); taskUpdater != nil {
			if successImages == totalImages {
				taskUpdater.UpdateTaskStatus(string(types.StatusCompleted), "")
			} else {
				failedCount := totalImages - successImages
				taskUpdater.UpdateTaskStatus(string(types.StatusFailed), fmt.Sprintf("成功 %d 张，失败 %d 张", successImages, failedCount))
			}
		}
	}

	// 如果有图片下载失败，返回错误
	if successImages < totalImages {
		failedCount := totalImages - successImages
		return fmt.Errorf("下载未完全成功，总共 %d 张图片，成功 %d 张，失败 %d 张", totalImages, successImages, failedCount)
	}

	return nil
}

// GetWnacgAlbumWithClient 获取整个专辑信息，包括所有分页URL
func GetWnacgAlbumWithClient(reqClient *request.Client, url string) (*WnacgAlbum, error) {
	var pageURLs []string
	currentURL := url

	albumName := ""

	// 添加第一页
	pageURLs = append(pageURLs, currentURL)

	// 使用频率限制的请求获取第一页
	resp, err := reqClient.RateLimitedGet(currentURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 读取响应
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}

	// 获取专辑名称（从页面标题获取）
	albumName = strings.TrimSpace(doc.Find("#bodywrap > h2").Text())
	if albumName == "" {
		albumName = "Unknown Album"
	}

	// 获取所有分页链接
	doc.Find(".paginator a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			// 拼接完整URL
			fullURL := href
			if !strings.HasPrefix(href, "http") {
				fullURL = "https://www.wnacg.com" + href
			}

			// 避免重复添加当前页
			if fullURL != currentURL {
				pageURLs = append(pageURLs, fullURL)
			}
		}
	})

	resp.Body.Close()

	if albumName == "" {
		return nil, fmt.Errorf("无法获取专辑名称")
	}

	// 去重分页URL
	uniqueURLs := make([]string, 0, len(pageURLs))
	seen := make(map[string]bool)
	for _, pageURL := range pageURLs {
		if !seen[pageURL] {
			seen[pageURL] = true
			uniqueURLs = append(uniqueURLs, pageURL)
		}
	}

	return &WnacgAlbum{
		Name:  albumName,
		Pages: uniqueURLs,
	}, nil
}

// GetMangaLinksFromPage 从分页中获取所有漫画页面的链接
func GetMangaLinksFromPage(reqClient *request.Client, pageURL string) ([]string, error) {
	resp, err := reqClient.RateLimitedGet(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	// 获取 class = cc 的 ul 元素，从 li 的 a 标签中获取每一页漫画的网址
	doc.Find("#bodywrap ul li a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			// 拼接完整URL
			fullURL := href
			if !strings.HasPrefix(href, "http") {
				fullURL = "https://www.wnacg.com" + href
			}
			links = append(links, fullURL)
		}
	})

	return links, nil
}

// ParseWnacgPageWithClient 解析Wnacg漫画页面获取真实图片URL
func ParseWnacgPageWithClient(reqClient *request.Client, link string) (string, error) {
	resp, err := reqClient.RateLimitedGet(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// 获取真实图片URL - 访问每一页漫画的网址，将每一页 id = picarea 的 img 作为图片结果
	imgURL := ""
	doc.Find("#picarea").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			imgURL = src
		}
	})

	if imgURL == "" {
		return "", fmt.Errorf("找不到图片URL")
	}

	// 如果URL是相对路径，转换为绝对路径
	if strings.HasPrefix(imgURL, "//") {
		imgURL = "https:" + imgURL
	} else if strings.HasPrefix(imgURL, "/") {
		imgURL = "https://www.wnacg.com" + imgURL
	}

	return imgURL, nil
}

// WnacgCrawler Wnacg爬虫
type WnacgCrawler struct {
	reqClient  *request.Client
	ctx        context.Context
	downloader types.Downloader
}

// NewWnacgCrawler 创建新的Wnacg爬虫
func NewWnacgCrawler(reqClient *request.Client, ctx context.Context) types.ImageCrawler {
	return &WnacgCrawler{
		reqClient: reqClient,
		ctx:       ctx,
	}
}

// GetDownloader 获取下载器
func (c *WnacgCrawler) GetDownloader() types.Downloader {
	return c.downloader
}

// SetDownloader 设置下载器
func (c *WnacgCrawler) SetDownloader(dl types.Downloader) {
	c.downloader = dl
}

// Crawl 执行爬取
func (c *WnacgCrawler) Crawl(url string, savePath string) (string, error) {
	// 将下载器传递给解析器，解析器会使用downloader获取的代理设置
	err := ParseWnacg(c.ctx, c.reqClient, url, savePath, c.downloader)
	if err != nil {
		return "", err
	}
	return savePath, nil
}

// CrawlAndSave 执行爬取并保存
func (c *WnacgCrawler) CrawlAndSave(url string, savePath string) string {
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
