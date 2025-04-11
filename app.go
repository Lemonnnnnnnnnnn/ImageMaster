package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 漫画信息结构
type Manga struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	PreviewImg  string   `json:"previewImg"`
	ImagesCount int      `json:"imagesCount"`
	Images      []string `json:"images,omitempty"`
}

// 应用配置
type Config struct {
	Libraries []string `json:"libraries"`
}

// App struct
type App struct {
	ctx        context.Context
	config     Config
	mangas     []Manga
	configPath string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 设置配置文件路径
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, _ = os.Getwd()
	}
	a.configPath = filepath.Join(configDir, "manga-viewer-config.json")

	// 加载配置
	a.LoadConfig()

	// 如果配置中有图书馆，自动加载
	if len(a.config.Libraries) > 0 {
		a.LoadAllLibraries()
	}
}

// LoadConfig 加载应用配置
func (a *App) LoadConfig() bool {
	data, err := os.ReadFile(a.configPath)
	if err != nil {
		a.config = Config{Libraries: []string{}}
		return false
	}

	err = json.Unmarshal(data, &a.config)
	if err != nil {
		a.config = Config{Libraries: []string{}}
		return false
	}

	return true
}

// SaveConfig 保存应用配置
func (a *App) SaveConfig() bool {
	data, err := json.Marshal(a.config)
	if err != nil {
		return false
	}

	err = os.WriteFile(a.configPath, data, 0644)
	if err != nil {
		return false
	}

	return true
}

// SelectLibrary 选择漫画库文件夹
func (a *App) SelectLibrary() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择漫画库文件夹",
	})

	if err != nil || dir == "" {
		return ""
	}

	// 检查是否已经添加过该库
	for _, lib := range a.config.Libraries {
		if lib == dir {
			return dir
		}
	}

	// 添加到配置中
	a.config.Libraries = append(a.config.Libraries, dir)
	a.SaveConfig()

	// 加载这个新库
	a.LoadLibrary(dir)

	return dir
}

// GetLibraries 获取所有图书馆路径
func (a *App) GetLibraries() []string {
	return a.config.Libraries
}

// LoadAllLibraries 加载所有图书馆
func (a *App) LoadAllLibraries() {
	a.mangas = []Manga{}
	for _, lib := range a.config.Libraries {
		a.LoadLibrary(lib)
	}
}

// LoadLibrary 加载指定图书馆
func (a *App) LoadLibrary(path string) bool {
	// 递归获取文件夹下的所有子文件夹
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只处理文件夹
		if !d.IsDir() {
			return nil
		}

		// 跳过根路径
		if path == a.config.Libraries[len(a.config.Libraries)-1] {
			return nil
		}

		// 获取文件夹中的图片
		images, err := a.GetImagesInDir(path)
		if err != nil || len(images) == 0 {
			return nil
		}

		// 排序图片
		a.SortImages(images)

		// 创建漫画信息
		manga := Manga{
			Name:        filepath.Base(path),
			Path:        path,
			PreviewImg:  images[0],
			ImagesCount: len(images),
			Images:      nil, // 不预加载所有图片路径
		}

		a.mangas = append(a.mangas, manga)

		return nil
	})

	return err == nil
}

// GetAllMangas 获取所有漫画
func (a *App) GetAllMangas() []Manga {
	return a.mangas
}

// GetMangaImages 获取指定漫画的所有图片
func (a *App) GetMangaImages(path string) []string {
	images, _ := a.GetImagesInDir(path)
	a.SortImages(images)
	return images
}

// GetImagesInDir 获取指定目录中的所有图片
func (a *App) GetImagesInDir(dirPath string) ([]string, error) {
	var images []string

	// 读取目录中的所有文件
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// 图片扩展名
	validExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true, ".bmp": true,
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if validExts[ext] {
			images = append(images, filepath.Join(dirPath, entry.Name()))
		}
	}

	return images, nil
}

// DeleteManga 删除漫画（删除文件夹）
func (a *App) DeleteManga(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		return false
	}

	// 从manga列表中移除
	for i, manga := range a.mangas {
		if manga.Path == path {
			a.mangas = append(a.mangas[:i], a.mangas[i+1:]...)
			break
		}
	}

	return true
}

// SortImages 排序图片文件
func (a *App) SortImages(images []string) {
	sort.Slice(images, func(i, j int) bool {
		nameI := filepath.Base(images[i])
		nameJ := filepath.Base(images[j])

		// 尝试提取 page_offset 格式
		partsI := strings.Split(strings.TrimSuffix(nameI, ".jpg"), "_")
		partsJ := strings.Split(strings.TrimSuffix(nameJ, ".jpg"), "_")

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

		// 回退到原来的排序逻辑
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
func (a *App) GetImageDataUrl(path string) string {
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
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encodeBytesToBase64(data))
}

// 辅助函数 - 编码为base64
func encodeBytesToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// 辅助函数 - 获取系统分隔符
func (a *App) GetPathSeparator() string {
	return string(filepath.Separator)
}

// 辅助函数 - 获取操作系统类型
func (a *App) GetOSType() string {
	// 获取操作系统类型
	osType := "unknown"

	if runtime.Environment(a.ctx).Platform == "windows" {
		osType = "windows"
	} else if runtime.Environment(a.ctx).Platform == "darwin" {
		osType = "darwin"
	} else if runtime.Environment(a.ctx).Platform == "linux" {
		osType = "linux"
	}

	return osType
}
