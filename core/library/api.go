package library

import (
	"context"
	"os"
	"path/filepath"

	"ImageMaster/core/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// API 图书馆API - 提供前端调用的接口
type API struct {
	ctx            context.Context
	configManager  types.ConfigManager
	libraryManager *Manager
}

// NewAPI 创建新的图书馆API实例
func NewAPI(configManager types.ConfigManager) *API {
	return &API{
		configManager: configManager,
	}
}

// SetContext 设置上下文
func (a *API) SetContext(ctx context.Context) {
	a.ctx = ctx
}

// InitializeLibraryManager 初始化图书馆管理器
func (a *API) InitializeLibraryManager() {
	// 设置默认输出目录
	userDir, err := os.UserHomeDir()
	if err != nil {
		userDir, _ = os.Getwd()
	}
	defaultOutputDir := filepath.Join(userDir, "Pictures", "ImageMaster")

	// 如果配置中有指定的输出目录，使用配置中的目录，否则使用默认目录
	outputDir := defaultOutputDir
	if configOutputDir := a.configManager.GetOutputDir(); configOutputDir != "" {
		outputDir = configOutputDir
	} else {
		// 如果是第一次使用，将默认目录保存到配置中
		a.configManager.SetOutputDir(defaultOutputDir)
	}

	// 初始化图书馆管理器
	a.libraryManager = NewManager(a.configManager, outputDir)
	a.libraryManager.SetContext(a.ctx)
	a.libraryManager.EnsureDir(outputDir)

	// 如果配置中有图书馆，自动加载
	if len(a.configManager.GetLibraries()) > 0 {
		a.libraryManager.LoadAllLibraries()
	}
}

// SelectLibrary 选择漫画库文件夹
func (a *API) SelectLibrary() string {
	return a.libraryManager.SelectLibrary()
}

// LoadAllLibraries 加载所有图书馆
func (a *API) LoadAllLibraries() {
	a.libraryManager.LoadAllLibraries()
}

// GetAllMangas 获取所有漫画
func (a *API) GetAllMangas() []Manga {
	return a.libraryManager.GetAllMangas()
}

// GetMangaImages 获取指定漫画的所有图片
func (a *API) GetMangaImages(path string) []string {
	return a.libraryManager.GetMangaImages(path)
}

// DeleteManga 删除漫画（删除文件夹）
func (a *API) DeleteManga(path string) bool {
	return a.libraryManager.DeleteManga(path)
}

// GetImageDataUrl 获取图片的DataURL
func (a *API) GetImageDataUrl(path string) string {
	return a.libraryManager.GetImageDataUrl(path)
}

// SetOutputDir 设置输出目录
func (a *API) SetOutputDir() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择保存目录",
	})

	if err != nil || dir == "" {
		return a.configManager.GetOutputDir()
	}

	a.configManager.SetOutputDir(dir)

	// 更新图书馆管理器的输出目录
	a.libraryManager.SetOutputDir(dir)

	return dir
}