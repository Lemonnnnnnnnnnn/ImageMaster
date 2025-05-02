package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetterConfig 配置
type GetterConfig struct {
	Concurrency int    // 并发下载数
	MaxRetries  int    // 最大重试次数
	RetryDelay  int    // 重试延迟(秒)
	OutputDir   string // 输出目录
}

// MediaInfo 媒体信息
type MediaInfo struct {
	URL      *url.URL
	Filename string
}

// Getter 结构体
type Getter struct {
	ctx        context.Context
	config     GetterConfig
	limiter    *Limiter
	mediaTypes map[string]bool
	ioManager  *IOManager
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

// NewGetter 创建新的 Getter 实例
func NewGetter(ctx context.Context) *Getter {
	// 默认支持的图片格式
	validExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true, ".bmp": true,
	}

	// 默认配置
	config := GetterConfig{
		Concurrency: 5,
		MaxRetries:  3,
		RetryDelay:  2,
	}

	// 设置默认输出目录
	userDir, err := os.UserHomeDir()
	if err != nil {
		userDir, _ = os.Getwd()
	}
	config.OutputDir = filepath.Join(userDir, "Pictures", "ImageMaster")

	// 确保输出目录存在
	os.MkdirAll(config.OutputDir, 0755)

	return &Getter{
		ctx:        ctx,
		config:     config,
		limiter:    NewLimiter(config.Concurrency),
		mediaTypes: validExts,
		ioManager:  NewIOManager(config.OutputDir),
	}
}

// SetOutputDir 设置输出目录
func (g *Getter) SetOutputDir(dir string) {
	g.config.OutputDir = dir
	g.ioManager = NewIOManager(dir)
}

// LoadMangaLibrary 加载漫画库
func (g *Getter) LoadMangaLibrary(rootPath string, mangas *[]Manga) bool {
	// 递归获取文件夹下的所有子文件夹
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理文件夹
		if !d.IsDir() {
			return nil
		}

		// 跳过根路径
		if path == rootPath {
			return nil
		}

		// 获取文件夹中的图片
		images, err := g.GetImagesInDir(path)
		if err != nil || len(images) == 0 {
			return nil
		}

		// 排序图片
		g.SortImages(images)

		// 创建漫画信息
		manga := Manga{
			Name:        filepath.Base(path),
			Path:        path,
			PreviewImg:  images[0],
			ImagesCount: len(images),
			Images:      nil, // 不预加载所有图片路径
		}

		*mangas = append(*mangas, manga)

		return nil
	})

	return err == nil
}

// GetImagesInDir 获取指定目录中的所有图片
func (g *Getter) GetImagesInDir(dirPath string) ([]string, error) {
	var images []string

	// 读取目录中的所有文件
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if g.mediaTypes[ext] {
			images = append(images, filepath.Join(dirPath, entry.Name()))
		}
	}

	return images, nil
}

// GetMangaImages 获取指定漫画的所有图片
func (g *Getter) GetMangaImages(path string) []string {
	images, _ := g.GetImagesInDir(path)
	g.SortImages(images)
	return images
}

// SortImages 排序图片文件
func (g *Getter) SortImages(images []string) {
	sort.Slice(images, func(i, j int) bool {
		nameI := filepath.Base(images[i])
		nameJ := filepath.Base(images[j])

		// 尝试提取 page_offset 格式
		partsI := strings.Split(strings.TrimSuffix(nameI, filepath.Ext(nameI)), "_")
		partsJ := strings.Split(strings.TrimSuffix(nameJ, filepath.Ext(nameJ)), "_")

		if len(partsI) == 2 && len(partsJ) == 2 {
			pageI, errI1 := strconv.Atoi(partsI[0])
			offsetI, errI2 := strconv.Atoi(partsI[1])
			pageJ, errJ1 := strconv.Atoi(partsJ[0])
			offsetJ, errJ2 := strconv.Atoi(partsJ[1])

			if errI1 == nil && errI2 == nil && errJ1 == nil && errJ2 == nil {
				if pageI != pageJ {
					return pageI < pageJ
				}
				return offsetI < offsetJ
			}
		}

		// 回退到提取数字排序逻辑
		reNum := regexp.MustCompile(`\d+`)
		numsI := reNum.FindAllString(nameI, -1)
		numsJ := reNum.FindAllString(nameJ, -1)

		if len(numsI) > 0 && len(numsJ) > 0 {
			numI, _ := strconv.Atoi(numsI[0])
			numJ, _ := strconv.Atoi(numsJ[0])
			return numI < numJ
		}

		if len(numsI) > 0 {
			return true
		}
		if len(numsJ) > 0 {
			return false
		}

		return nameI < nameJ
	})
}

// GetImageDataUrl 获取图片的DataURL
func (g *Getter) GetImageDataUrl(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	// 获取MIME类型
	ext := strings.ToLower(filepath.Ext(path))
	mimeType := "image/jpeg" // 默认
	switch ext {
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	case ".bmp":
		mimeType = "image/bmp"
	}

	// 构建data URL
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))
}

// CrawlWebImages 从网页爬取图片
func (g *Getter) CrawlWebImages(webUrl string, saveName string) (string, error) {
	Info("开始爬取网页图片: %s", webUrl)

	// 检查URL格式
	parsedURL, err := url.Parse(webUrl)
	if err != nil {
		return "", fmt.Errorf("无效的URL: %w", err)
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	req, err := http.NewRequest("GET", webUrl, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// 执行请求
	resp, err := client.Do(req)
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
	saveDir := filepath.Join(g.config.OutputDir, title)
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
			URL:      imgURL,
			Filename: filename,
		})
	})

	if len(imageUrls) == 0 {
		return "", fmt.Errorf("未找到图片")
	}

	Info("共找到 %d 张图片", len(imageUrls))

	// 下载图片
	downloader := NewDownloader(g.config.MaxRetries, g.config.RetryDelay, true)

	for i, img := range imageUrls {
		mediaInfo := img
		g.limiter.Execute(func() {
			savePath := filepath.Join(saveDir, mediaInfo.Filename)
			Info("下载图片 %d/%d: %s", i+1, len(imageUrls), mediaInfo.URL.String())

			err := downloader.DownloadFile(mediaInfo.URL.String(), savePath, nil)
			if err != nil {
				Error("下载图片失败: %s, 错误: %v", mediaInfo.URL.String(), err)
			} else {
				Info("图片下载成功: %s", savePath)
			}
		})
	}

	g.limiter.Wait()
	return saveDir, nil
}

// normalizeURL 标准化图片URL
func (g *Getter) normalizeURL(baseURL *url.URL, src string) (*url.URL, error) {
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

// 通过对象创建一个本地存储管理器
type IOManager struct {
	baseDir string
}

// 创建一个新的IO管理器
func NewIOManager(baseDir string) *IOManager {
	return &IOManager{
		baseDir: baseDir,
	}
}

// 确保目录存在
func (m *IOManager) EnsureDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

// 写入文件
func (m *IOManager) WriteFile(data interface{}, filename string, subDir string, mainDir string) error {
	// 构建完整路径
	fullPath := filepath.Join(m.baseDir, mainDir, subDir, filename)

	// 确保目录存在
	if err := m.EnsureDir(fullPath); err != nil {
		return err
	}

	// 根据数据类型写入文件
	var err error
	switch v := data.(type) {
	case []byte:
		err = os.WriteFile(fullPath, v, 0644)
	case string:
		err = os.WriteFile(fullPath, []byte(v), 0644)
	default:
		return fmt.Errorf("unsupported data type")
	}

	return err
}
