package parsers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"ImageMaster/core/request"
	"ImageMaster/core/types"
)

// TelegraphAlbum Telegraph专辑
type TelegraphAlbum struct {
	Name   string
	Images []TelegraphImage
}

// TelegraphImage Telegraph图片
type TelegraphImage struct {
	Name string
	URL  string
}

// ParseTelegraph 解析Telegraph网站
func ParseTelegraph(client *http.Client, url string, savePath string, dl types.Downloader) error {
	fmt.Printf("下载 Telegraph 专辑: %s\n", url)

	// 创建请求客户端
	reqClient := request.NewClient()

	// 使用下载器的代理配置
	if dl != nil && dl.GetProxy() != "" {
		reqClient.SetProxy(dl.GetProxy())
	}

	// 设置User-Agent
	reqClient.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	album, err := GetTelegraphAlbum(reqClient, url)
	if err != nil {
		return fmt.Errorf("获取专辑失败: %w", err)
	}

	// 使用传入的下载器
	var localDownloader types.Downloader
	if dl != nil {
		localDownloader = dl
		fmt.Printf("Telegraph解析器使用传入的下载器\n")
	} else {
		// 未提供下载器，返回错误
		return fmt.Errorf("未提供下载器")
	}

	// 保存路径
	albumPath := savePath + "/" + album.Name

	// 准备批量下载
	var imgURLs []string
	var filePaths []string

	for _, image := range album.Images {
		fullPath := fmt.Sprintf("%s/%s", albumPath, image.Name)
		imgURLs = append(imgURLs, image.URL)
		filePaths = append(filePaths, fullPath)
	}

	// 批量下载所有图片
	totalImages := len(imgURLs)
	successCount, err := localDownloader.BatchDownload(imgURLs, filePaths, nil)
	if err != nil {
		fmt.Printf("批量下载出错: %v\n", err)
		return fmt.Errorf("批量下载出错: %w", err)
	}

	fmt.Printf("下载完成，总共 %d 张图片，成功 %d 张\n", totalImages, successCount)
	
	// 如果有图片下载失败，返回错误
	if successCount < totalImages {
		failedCount := totalImages - successCount
		return fmt.Errorf("下载未完全成功，总共 %d 张图片，成功 %d 张，失败 %d 张", totalImages, successCount, failedCount)
	}
	
	return nil
}

// GetTelegraphAlbum 获取Telegraph专辑
func GetTelegraphAlbum(reqClient *request.Client, url string) (*TelegraphAlbum, error) {
	resp, err := reqClient.Get(url)
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

	// 获取专辑名称
	albumName := ""
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		albumName = s.Text()
	})

	if albumName == "" {
		albumName = "Telegraph Album" // 默认名称
	}

	// 获取所有图片
	var images []TelegraphImage
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			url := FormatTelegraphURL(src)

			// 创建图片信息
			image := TelegraphImage{
				Name: fmt.Sprintf("%d.jpg", i),
				URL:  url,
			}

			images = append(images, image)
		}
	})

	return &TelegraphAlbum{
		Name:   albumName,
		Images: images,
	}, nil
}

// FormatTelegraphURL 格式化Telegraph URL
func FormatTelegraphURL(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}
	return "https://telegra.ph" + url
}
