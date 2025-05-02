package viewer

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"ImageMaster/core/crawler"
	"ImageMaster/core/getter"
)

// 应用配置
type Config struct {
	Libraries []string `json:"libraries"`
}

// Viewer 结构体
type Viewer struct {
	ctx            context.Context
	config         Config
	mangas         []getter.Manga
	configPath     string
	localGetter    *getter.LocalGetter
	crawlerFactory *crawler.CrawlerFactory
}

// NewViewer 创建新的 Viewer 实例
func NewViewer() *Viewer {
	return &Viewer{}
}

// Startup 启动应用
func (v *Viewer) Startup(ctx context.Context) {
	v.ctx = ctx

	// 设置配置文件路径
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, _ = os.Getwd()
	}
	v.configPath = filepath.Join(configDir, "manga-viewer-config.json")

	// 初始化 LocalGetter
	// 设置默认输出目录
	userDir, err := os.UserHomeDir()
	if err != nil {
		userDir, _ = os.Getwd()
	}
	outputDir := filepath.Join(userDir, "Pictures", "ImageMaster")
	os.MkdirAll(outputDir, 0755)

	v.localGetter = getter.NewLocalGetter(outputDir)
	v.crawlerFactory = crawler.NewCrawlerFactory(ctx)

	// 加载配置
	v.LoadConfig()

	// 如果配置中有图书馆，自动加载
	if len(v.config.Libraries) > 0 {
		v.LoadAllLibraries()
	}
}

// LoadConfig 加载应用配置
func (v *Viewer) LoadConfig() bool {
	data, err := os.ReadFile(v.configPath)
	if err != nil {
		v.config = Config{Libraries: []string{}}
		return false
	}

	err = json.Unmarshal(data, &v.config)
	if err != nil {
		v.config = Config{Libraries: []string{}}
		return false
	}

	return true
}

// SaveConfig 保存应用配置
func (v *Viewer) SaveConfig() bool {
	data, err := json.Marshal(v.config)
	if err != nil {
		return false
	}

	err = os.WriteFile(v.configPath, data, 0644)
	if err != nil {
		return false
	}

	return true
}

// SelectLibrary 选择漫画库文件夹
func (v *Viewer) SelectLibrary() string {
	dir, err := runtime.OpenDirectoryDialog(v.ctx, runtime.OpenDialogOptions{
		Title: "选择漫画库文件夹",
	})

	if err != nil || dir == "" {
		return ""
	}

	// 检查是否已经添加过该库
	for _, lib := range v.config.Libraries {
		if lib == dir {
			return dir
		}
	}

	// 添加到配置中
	v.config.Libraries = append(v.config.Libraries, dir)
	v.SaveConfig()

	// 加载这个新库
	v.LoadLibrary(dir)

	return dir
}

// GetLibraries 获取所有图书馆路径
func (v *Viewer) GetLibraries() []string {
	return v.config.Libraries
}

// LoadAllLibraries 加载所有图书馆
func (v *Viewer) LoadAllLibraries() {
	v.mangas = []getter.Manga{}
	for _, lib := range v.config.Libraries {
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

// CrawlFromWeb 从网页爬取图片
func (v *Viewer) CrawlFromWeb(url string, saveName string) string {
	// 检测网站类型
	siteType := v.crawlerFactory.DetectSiteType(url)

	// 创建对应爬虫
	crawler := v.crawlerFactory.CreateCrawler(siteType)

	// 执行爬取
	saveDir, err := crawler.Crawl(url, v.GetOutputDir())
	if err != nil {
		return ""
	}

	// 刷新库
	v.LoadAllLibraries()

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
	return dir
}

// GetOutputDir 获取当前输出目录
func (v *Viewer) GetOutputDir() string {
	return v.localGetter.GetOutputDir()
}
