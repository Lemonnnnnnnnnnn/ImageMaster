package getter

// Manga 漫画信息结构
type Manga struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	PreviewImg  string   `json:"previewImg"`
	ImagesCount int      `json:"imagesCount"`
	Images      []string `json:"images,omitempty"`
}

// MediaInfo 媒体信息
type MediaInfo struct {
	URL      string
	Filename string
}

// IOManager 文件存储管理器
type IOManager struct {
	baseDir string
}

// NewIOManager 创建新的IO管理器
func NewIOManager(baseDir string) *IOManager {
	return &IOManager{
		baseDir: baseDir,
	}
}
