package parsers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"ImageMaster/core/request"
	"ImageMaster/core/types"
)

// NhentaiGallery Nhentai画廊
type NhentaiGallery struct {
	ID     string
	Name   string
	Images []string // 存储所有图片的URL
}

// ParseNhentai 解析Nhentai网站
func ParseNhentai(ctx context.Context, reqClient *request.Client, url string, savePath string, dl types.Downloader) error {
	fmt.Printf("下载 Nhentai 画廊: %s\n", url)

	// 使用下载器的代理配置
	if dl != nil && dl.GetProxy() != "" {
		reqClient.SetProxy(dl.GetProxy())
	}

	// 设置User-Agent
	reqClient.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	nhentaiGallery, err := GetNhentaiGalleryWithClient(reqClient, url)
	if err != nil {
		return fmt.Errorf("获取画廊失败: %w", err)
	}

	// 使用TaskUpdater更新任务名称（如果可用）
	if dl != nil {
		if taskUpdater := dl.GetTaskUpdater(); taskUpdater != nil {
			taskUpdater.UpdateTaskName(nhentaiGallery.Name)
			taskUpdater.UpdateTaskStatus(string(types.StatusParsing), "")
			fmt.Printf("已更新任务名称为: %s\n", nhentaiGallery.Name)
		}
	}

	// 使用传入的下载器
	var localDownloader types.Downloader
	if dl != nil {
		localDownloader = dl
		fmt.Printf("Nhentai解析器使用传入的下载器\n")
	} else {
		// 未提供下载器，返回错误
		return fmt.Errorf("未提供下载器")
	}

	// 保存路径
	galleryPath := savePath + "/" + nhentaiGallery.Name

	// 计算总图片数量
	totalImages := len(nhentaiGallery.Images)
	fmt.Printf("已收集 %d 张图片URL，开始下载...\n", totalImages)

	// 更新任务状态为下载中
	if dl != nil {
		if taskUpdater := dl.GetTaskUpdater(); taskUpdater != nil {
			taskUpdater.UpdateTaskStatus(string(types.StatusDownloading), "")
			taskUpdater.UpdateTaskProgress(0, totalImages)
		}
	}

	// 准备批量下载的URL和路径
	var filePaths []string
	for i := range nhentaiGallery.Images {
		filename := fmt.Sprintf("%03d.webp", i+1)
		fullPath := fmt.Sprintf("%s/%s", galleryPath, filename)
		filePaths = append(filePaths, fullPath)
	}

	// 批量下载所有图片
	successImages, err := localDownloader.BatchDownload(nhentaiGallery.Images, filePaths, nil)
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

// GetNhentaiGalleryWithClient 获取整个画廊信息，包括所有图片URL
func GetNhentaiGalleryWithClient(reqClient *request.Client, galleryURL string) (*NhentaiGallery, error) {
	// 从URL中提取画廊ID
	galleryID, err := extractGalleryID(galleryURL)
	if err != nil {
		return nil, fmt.Errorf("无法提取画廊ID: %w", err)
	}

	// 使用频率限制的请求获取主页面
	resp, err := reqClient.RateLimitedGet(galleryURL)
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

	// 获取画廊标题
	galleryName := strings.TrimSpace(doc.Find("body > div.main_cnt > div > div.gallery_top > div.info > h1").Text())
	if galleryName == "" {
		galleryName = "Unknown Gallery"
	}

	// 收集前十页的图片URL
	var imageURLs []string
	doc.Find("#thumbs_append > div > a > img").Each(func(i int, s *goquery.Selection) {
		if dataSrc, exists := s.Attr("data-src"); exists && dataSrc != "" {
			// 转换缩略图URL为完整图片URL
			fullImageURL := convertThumbnailToFullImage(dataSrc)
			imageURLs = append(imageURLs, fullImageURL)
		}
	})

	fmt.Printf("从主页面获取到 %d 张图片URL\n", len(imageURLs))

	// 获取更多图片（通过AJAX接口）
	moreImages, err := getMoreImagesFromAPI(reqClient, doc, galleryID, len(imageURLs))
	if err != nil {
		fmt.Printf("获取更多图片失败: %v\n", err)
	} else {
		imageURLs = append(imageURLs, moreImages...)
		fmt.Printf("通过API获取到额外 %d 张图片URL\n", len(moreImages))
	}

	if len(imageURLs) == 0 {
		return nil, fmt.Errorf("未找到任何图片")
	}

	return &NhentaiGallery{
		ID:     galleryID,
		Name:   galleryName,
		Images: imageURLs,
	}, nil
}

// extractGalleryID 从URL中提取画廊ID
func extractGalleryID(galleryURL string) (string, error) {
	// 从类似 "https://nhentai.xxx/g/537651/" 的URL中提取 "537651"
	re := regexp.MustCompile(`/g/(\d+)/?`)
	matches := re.FindStringSubmatch(galleryURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法从URL中提取画廊ID")
	}
	return matches[1], nil
}

// convertThumbnailToFullImage 将缩略图URL转换为完整图片URL
func convertThumbnailToFullImage(thumbnailURL string) string {
	// 将结尾的【数字t.jpg】替换为【数字.webp】
	// 例如：http://i4.nhentaimg.com/016/9sazckpugf/11t.jpg -> http://i4.nhentaimg.com/016/9sazckpugf/11.webp
	re := regexp.MustCompile(`(\d+)t\.jpg$`)
	return re.ReplaceAllString(thumbnailURL, "$1.webp")
}

// getMoreImagesFromAPI 通过AJAX API获取更多图片
func getMoreImagesFromAPI(reqClient *request.Client, doc *goquery.Document, galleryID string, visiblePages int) ([]string, error) {
	// 获取CSRF token
	csrfToken, exists := doc.Find(`meta[name="csrf-token"]`).Attr("content")
	if !exists {
		return nil, fmt.Errorf("未找到CSRF token")
	}

	// 获取其他必需的参数
	server := doc.Find("#load_server").AttrOr("value", "")
	uID := doc.Find("#gallery_id").AttrOr("value", "")
	gID := doc.Find("#load_id").AttrOr("value", "")
	imgDir := doc.Find("#load_dir").AttrOr("value", "")
	totalPagesStr := doc.Find("#load_pages").AttrOr("value", "")

	// 如果无法获取必需参数，返回空结果而不是错误
	if server == "" || uID == "" || gID == "" || imgDir == "" || totalPagesStr == "" {
		fmt.Printf("无法获取API参数，跳过API调用\n")
		return []string{}, nil
	}

	totalPages, err := strconv.Atoi(totalPagesStr)
	if err != nil {
		return nil, fmt.Errorf("无法解析总页数: %w", err)
	}

	// 如果可见页面数量已经等于总页数，不需要调用API
	if visiblePages >= totalPages {
		return []string{}, nil
	}

	// 准备POST数据
	formData := url.Values{}
	formData.Set("_token", csrfToken)
	formData.Set("server", server)
	formData.Set("u_id", uID)
	formData.Set("g_id", gID)
	formData.Set("img_dir", imgDir)
	formData.Set("visible_pages", strconv.Itoa(visiblePages))
	formData.Set("total_pages", totalPagesStr)
	formData.Set("type", "2")

	// 设置请求头
	reqClient.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	reqClient.SetHeader("X-Requested-With", "XMLHttpRequest")

	// 发送POST请求
	resp, err := reqClient.Post("https://nhentai.xxx/modules/thumbs_loader.php", strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
	}

	// 解析返回的HTML
	apiDoc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解析API响应失败: %w", err)
	}

	// 从API响应中提取图片URL
	var moreImages []string
	apiDoc.Find("img").Each(func(i int, s *goquery.Selection) {
		if dataSrc, exists := s.Attr("data-src"); exists && dataSrc != "" {
			// 转换缩略图URL为完整图片URL
			fullImageURL := convertThumbnailToFullImage(dataSrc)
			moreImages = append(moreImages, fullImageURL)
		}
	})

	return moreImages, nil
}

// NhentaiCrawler Nhentai爬虫
type NhentaiCrawler struct {
	reqClient  *request.Client
	ctx        context.Context
	downloader types.Downloader
}

// NewNhentaiCrawler 创建新的Nhentai爬虫
func NewNhentaiCrawler(reqClient *request.Client, ctx context.Context) types.ImageCrawler {
	return &NhentaiCrawler{
		reqClient: reqClient,
		ctx:       ctx,
	}
}

// GetDownloader 获取下载器
func (c *NhentaiCrawler) GetDownloader() types.Downloader {
	return c.downloader
}

// SetDownloader 设置下载器
func (c *NhentaiCrawler) SetDownloader(dl types.Downloader) {
	c.downloader = dl
}

// Crawl 执行爬取
func (c *NhentaiCrawler) Crawl(url string, savePath string) (string, error) {
	// 将下载器传递给解析器，解析器会使用downloader获取的代理设置
	err := ParseNhentai(c.ctx, c.reqClient, url, savePath, c.downloader)
	if err != nil {
		return "", err
	}
	return savePath, nil
}

// CrawlAndSave 执行爬取并保存
func (c *NhentaiCrawler) CrawlAndSave(url string, savePath string) string {
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
