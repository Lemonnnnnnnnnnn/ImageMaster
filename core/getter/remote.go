package getter

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"ImageMaster/core/downloader"
)

// RemoteGetter 远程图片获取器
type RemoteGetter struct {
	ctx       context.Context
	limiter   *Limiter
	outputDir string
	client    *http.Client
}

// Limiter 简单并发限制器
type Limiter struct {
	wg    sync.WaitGroup
	ch    chan struct{}
	limit int
}

// NewLimiter 创建并发限制器
func NewLimiter(limit int) *Limiter {
	if limit <= 0 {
		limit = 5 // 默认并发数
	}
	return &Limiter{
		ch:    make(chan struct{}, limit),
		limit: limit,
	}
}

// Execute 执行并发任务
func (l *Limiter) Execute(fn func()) {
	l.wg.Add(1)
	l.ch <- struct{}{}

	go func() {
		defer func() {
			<-l.ch
			l.wg.Done()
		}()
		fn()
	}()
}

// Wait 等待所有任务完成
func (l *Limiter) Wait() {
	l.wg.Wait()
}

// NewRemoteGetter 创建远程图片获取器
func NewRemoteGetter(ctx context.Context, outputDir string, concurrency int) *RemoteGetter {
	return &RemoteGetter{
		ctx:       ctx,
		limiter:   NewLimiter(concurrency),
		outputDir: outputDir,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetOutputDir 设置输出目录
func (g *RemoteGetter) SetOutputDir(dir string) {
	g.outputDir = dir
}

// GetOutputDir 获取输出目录
func (g *RemoteGetter) GetOutputDir() string {
	return g.outputDir
}

// CrawlWebImages 从网页爬取图片 - 通用网页爬虫
func (g *RemoteGetter) CrawlWebImages(webUrl string, saveName string) (string, error) {
	fmt.Printf("开始爬取网页图片: %s\n", webUrl)

	// 检查URL格式
	parsedURL, err := url.Parse(webUrl)
	if err != nil {
		return "", fmt.Errorf("无效的URL: %w", err)
	}

	// 发送请求
	req, err := http.NewRequest("GET", webUrl, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 执行请求
	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("解析HTML失败: %w", err)
	}

	// 获取标题
	title := doc.Find("title").Text()
	if title == "" {
		title = doc.Find("h1").First().Text()
	}
	if title == "" {
		title = saveName
	}
	if saveName != "" {
		title = saveName
	}

	// 创建保存目录
	saveDir := filepath.Join(g.outputDir, title)
	err = os.MkdirAll(saveDir, 0755)
	if err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 查找所有图片
	var imageUrls []MediaInfo
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		// 处理相对URL
		imgURL, err := g.normalizeURL(parsedURL, src)
		if err != nil {
			return
		}

		// 生成文件名 (格式: 001.jpg, 002.jpg, ...)
		ext := filepath.Ext(imgURL.Path)
		if ext == "" {
			ext = ".jpg" // 默认扩展名
		}
		filename := fmt.Sprintf("%03d%s", i+1, ext)

		imageUrls = append(imageUrls, MediaInfo{
			URL:      imgURL.String(),
			Filename: filename,
		})
	})

	if len(imageUrls) == 0 {
		return "", fmt.Errorf("未找到图片")
	}

	fmt.Printf("共找到 %d 张图片\n", len(imageUrls))

	// 下载图片
	downloader := downloader.NewDownloader(3, 2, true)

	for i, img := range imageUrls {
		mediaInfo := img
		g.limiter.Execute(func() {
			savePath := filepath.Join(saveDir, mediaInfo.Filename)
			fmt.Printf("下载图片 %d/%d: %s\n", i+1, len(imageUrls), mediaInfo.URL)

			err := downloader.DownloadFile(mediaInfo.URL, savePath, nil)
			if err != nil {
				fmt.Printf("下载图片失败: %s, 错误: %v\n", mediaInfo.URL, err)
			} else {
				fmt.Printf("图片下载成功: %s\n", savePath)
			}
		})
	}

	g.limiter.Wait()
	return saveDir, nil
}

// normalizeURL 标准化图片URL
func (g *RemoteGetter) normalizeURL(baseURL *url.URL, src string) (*url.URL, error) {
	// 如果URL不是以http开头，添加协议和主机
	if !strings.HasPrefix(src, "http") {
		if strings.HasPrefix(src, "//") {
			src = baseURL.Scheme + ":" + src
		} else if strings.HasPrefix(src, "/") {
			src = baseURL.Scheme + "://" + baseURL.Host + src
		} else {
			base := baseURL.Scheme + "://" + baseURL.Host
			if path := filepath.Dir(baseURL.Path); path != "/" && path != "." {
				base += path
			}
			if !strings.HasSuffix(base, "/") {
				base += "/"
			}
			src = base + src
		}
	}

	imgURL, err := url.Parse(src)
	if err != nil {
		return nil, err
	}

	return imgURL, nil
}
