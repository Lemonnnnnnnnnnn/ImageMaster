package parsers

import (
	"fmt"
	"net/http"
	"net/url"
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

// NhentaiParser Nhentai解析器实现
type NhentaiParser struct{}

// GetName 获取解析器名称
func (p *NhentaiParser) GetName() string {
	return "Nhentai"
}

// Parse 解析URL获取图片信息
func (p *NhentaiParser) Parse(reqClient *request.Client, url string) (*ParseResult, error) {
	nhentaiGallery, err := GetNhentaiGalleryWithClient(reqClient, url)
	if err != nil {
		return nil, fmt.Errorf("获取画廊失败: %w", err)
	}

	// 准备文件路径
	var filePaths []string
	for i := range nhentaiGallery.Images {
		filename := fmt.Sprintf("%03d.webp", i+1)
		filePaths = append(filePaths, filename)
	}

	return &ParseResult{
		Name:      nhentaiGallery.Name,
		ImageURLs: nhentaiGallery.Images,
		FilePaths: filePaths,
	}, nil
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
	*BaseCrawler
}

// NewNhentaiCrawler 创建新的Nhentai爬虫
func NewNhentaiCrawler(reqClient *request.Client) types.ImageCrawler {
	parser := &NhentaiParser{}
	baseCrawler := NewBaseCrawler(reqClient, parser)
	return &NhentaiCrawler{
		BaseCrawler: baseCrawler,
	}
}
