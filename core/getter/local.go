package getter

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// LocalGetter 本地图片获取器
type LocalGetter struct {
	mediaTypes map[string]bool
	outputDir  string
}

// NewLocalGetter 创建本地图片获取器
func NewLocalGetter(outputDir string) *LocalGetter {
	// 默认支持的图片格式
	validExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
		".gif": true, ".webp": true, ".bmp": true,
	}

	return &LocalGetter{
		mediaTypes: validExts,
		outputDir:  outputDir,
	}
}

// SetOutputDir 设置输出目录
func (g *LocalGetter) SetOutputDir(dir string) {
	g.outputDir = dir
}

// GetOutputDir 获取输出目录
func (g *LocalGetter) GetOutputDir() string {
	return g.outputDir
}

// LoadMangaLibrary 加载漫画库
func (g *LocalGetter) LoadMangaLibrary(rootPath string, mangas *[]Manga) bool {
	// 递归获取文件夹下的所有子文件夹
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
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
func (g *LocalGetter) GetImagesInDir(dirPath string) ([]string, error) {
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
func (g *LocalGetter) GetMangaImages(path string) []string {
	images, _ := g.GetImagesInDir(path)
	g.SortImages(images)
	return images
}

// SortImages 排序图片文件
func (g *LocalGetter) SortImages(images []string) {
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
func (g *LocalGetter) GetImageDataUrl(path string) string {
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
