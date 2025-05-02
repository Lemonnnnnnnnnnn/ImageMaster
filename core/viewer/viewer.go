package viewer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"ImageMaster/core/config"
	"ImageMaster/core/crawler"
	"ImageMaster/core/getter"
)

// 下载进度回调函数类型
type DownloadProgressCallback func(current int, total int)

// Viewer 结构体
type Viewer struct {
	ctx                  context.Context
	configManager        *config.Manager
	mangas               []getter.Manga
	localGetter          *getter.LocalGetter
	crawlerFactory       *crawler.CrawlerFactory
	progressCallbackLock sync.Mutex
	// progressCallback     DownloadProgressCallback
}

// NewViewer 创建新的 Viewer 实例
func NewViewer() *Viewer {
	return &Viewer{}
}

// Startup 启动应用
func (v *Viewer) Startup(ctx context.Context) {
	v.ctx = ctx

	// 创建配置管理器
	v.configManager = config.NewManager("manga-viewer-config.json")
	v.configManager.LoadConfig()

	// 设置默认输出目录
	userDir, err := os.UserHomeDir()
	if err != nil {
		userDir, _ = os.Getwd()
	}
	defaultOutputDir := filepath.Join(userDir, "Pictures", "ImageMaster")

	// 如果配置中有指定的输出目录，使用配置中的目录，否则使用默认目录
	outputDir := defaultOutputDir
	if cfg := v.configManager.GetConfig(); cfg.OutputDir != "" {
		outputDir = cfg.OutputDir
	} else {
		// 如果是第一次使用，将默认目录保存到配置中
		v.configManager.SetOutputDir(defaultOutputDir)
	}

	// 确保输出目录存在
	os.MkdirAll(outputDir, 0755)

	v.localGetter = getter.NewLocalGetter(outputDir)
	v.crawlerFactory = crawler.NewCrawlerFactory(ctx)

	// 如果配置中有图书馆，自动加载
	if len(v.configManager.GetLibraries()) > 0 {
		v.LoadAllLibraries()
	}
}

// SelectLibrary 选择漫画库文件夹
func (v *Viewer) SelectLibrary() string {
	dir, err := runtime.OpenDirectoryDialog(v.ctx, runtime.OpenDialogOptions{
		Title: "选择漫画库文件夹",
	})

	if err != nil || dir == "" {
		return ""
	}

	// 添加到配置中
	v.configManager.AddLibrary(dir)

	// 加载这个新库
	v.LoadLibrary(dir)

	return dir
}

// GetLibraries 获取所有图书馆路径
func (v *Viewer) GetLibraries() []string {
	return v.configManager.GetLibraries()
}

// LoadAllLibraries 加载所有图书馆
func (v *Viewer) LoadAllLibraries() {
	v.mangas = []getter.Manga{}
	for _, lib := range v.configManager.GetLibraries() {
		v.LoadLibrary(lib)
	}
}

// LoadLibrary 加载指定图书馆
func (v *Viewer) LoadLibrary(path string) bool {
	return v.localGetter.LoadMangaLibrary(path, &v.mangas)
}

// GetAllMangas 获取所有漫画
func (v *Viewer) GetAllMangas() []getter.Manga {
	return v.mangas
}

// GetMangaImages 获取指定漫画的所有图片
func (v *Viewer) GetMangaImages(path string) []string {
	return v.localGetter.GetMangaImages(path)
}

// DeleteManga 删除漫画（删除文件夹）
func (v *Viewer) DeleteManga(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		return false
	}

	// 从manga列表中移除
	for i, manga := range v.mangas {
		if manga.Path == path {
			v.mangas = append(v.mangas[:i], v.mangas[i+1:]...)
			break
		}
	}

	return true
}

// GetImageDataUrl 获取图片的DataURL
func (v *Viewer) GetImageDataUrl(path string) string {
	return v.localGetter.GetImageDataUrl(path)
}

// 辅助函数 - 获取系统分隔符
func (v *Viewer) GetPathSeparator() string {
	return string(filepath.Separator)
}

// 辅助函数 - 获取操作系统类型
func (v *Viewer) GetOSType() string {
	// 获取操作系统类型
	osType := "unknown"

	if runtime.Environment(v.ctx).Platform == "windows" {
		osType = "windows"
	} else if runtime.Environment(v.ctx).Platform == "darwin" {
		osType = "darwin"
	} else if runtime.Environment(v.ctx).Platform == "linux" {
		osType = "linux"
	}

	return osType
}

// NotifyDownloadProgress 通知下载进度
func (v *Viewer) NotifyDownloadProgress(current, total int) {
	fmt.Printf("通知下载进度: %d/%d\n", current, total)

	// 关键修改：直接通过事件通知前端，不使用回调机制
	runtime.EventsEmit(v.ctx, "download:progress", current, total)

	v.progressCallbackLock.Lock()
	defer v.progressCallbackLock.Unlock()

	// if v.progressCallback != nil {
	// 	// 仍然执行老的回调，保持兼容性
	// 	v.progressCallback(current, total)
	// }
}

// CrawlFromWeb 从网页爬取图片
func (v *Viewer) CrawlFromWeb(url string, saveName string) string {
	// 设置爬虫工厂的配置管理器
	v.crawlerFactory.SetConfigManager(v.configManager)

	// 检测网站类型
	siteType := v.crawlerFactory.DetectSiteType(url)
	fmt.Printf("检测到网站类型: %s\n", siteType)

	// 创建对应爬虫
	crawler := v.crawlerFactory.CreateCrawler(siteType)

	// 设置进度回调
	downloader := crawler.GetDownloader()
	if downloader != nil {
		fmt.Printf("设置进度回调\n")

		// 直接使用匿名函数包装，确保调用是正确的
		downloader.SetProgressCallback(func(current, total int) {
			fmt.Printf("进度回调被触发: %d/%d\n", current, total)
			v.NotifyDownloadProgress(current, total)
		})
	} else {
		fmt.Printf("警告: 下载器为nil，无法设置进度回调\n")
	}

	// 执行爬取
	saveDir, err := crawler.Crawl(url, v.GetOutputDir())
	if err != nil {
		fmt.Printf("爬取失败: %v\n", err)
		return ""
	}

	// 刷新库
	v.LoadAllLibraries()

	fmt.Printf("爬取完成，保存到: %s\n", saveDir)
	return saveDir
}

// SetOutputDir 设置输出目录
func (v *Viewer) SetOutputDir() string {
	dir, err := runtime.OpenDirectoryDialog(v.ctx, runtime.OpenDialogOptions{
		Title: "选择图片保存目录",
	})

	if err != nil || dir == "" {
		return ""
	}

	v.localGetter.SetOutputDir(dir)

	// 更新配置并保存
	v.configManager.SetOutputDir(dir)

	return dir
}

// GetOutputDir 获取当前输出目录
func (v *Viewer) GetOutputDir() string {
	return v.localGetter.GetOutputDir()
}

// SetProxy 设置代理
func (v *Viewer) SetProxy(proxyURL string) bool {
	// 更新配置
	return v.configManager.SetProxy(proxyURL)
}

// GetProxy 获取当前代理设置
func (v *Viewer) GetProxy() string {
	return v.configManager.GetProxy()
}
