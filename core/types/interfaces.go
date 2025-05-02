package types

// Downloader 下载器接口
type Downloader interface {
	DownloadFile(url string, filepath string, headers map[string]string) error
	BatchDownload(urls []string, filepaths []string, headers map[string]string) (int, error)
	SetProgressCallback(callback func(current, total int))
	GetProxy() string
	SetProxy(proxyURL string) error
}

// ProgressReporter 进度报告接口
type ProgressReporter interface {
	ReportProgress(current, total int)
}

// ConfigProvider 配置提供者接口
type ConfigProvider interface {
	// GetOutputDir 获取输出目录
	GetOutputDir() string

	// GetProxy 获取代理设置
	GetProxy() string
}

// ConfigManager 配置管理接口
type ConfigManager interface {
	GetOutputDir() string
	SetOutputDir(dir string) bool
	GetProxy() string
	SetProxy(proxyURL string) bool
}

// StorageProvider 存储提供者接口
type StorageProvider interface {
	GetConfigManager() ConfigManager
	AddDownloadRecord(task interface{})
	GetDownloadHistory() []interface{}
	ClearDownloadHistory()
}
