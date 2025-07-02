package storage

import (
	"ImageMaster/core/config"
	"ImageMaster/core/logger"
	"ImageMaster/core/task"
	"ImageMaster/core/types"
)

// 确保Manager实现StorageProvider接口
var _ types.StorageProvider = (*Manager)(nil)

// Manager 存储管理器 - 统一的存储接口
type Manager struct {
	configManager  *config.Manager
	historyManager *HistoryManager
}

// NewManager 创建存储管理器
func NewManager(appName string) *Manager {
	logger.Info("Initializing storage manager for app: %s", appName)

	return &Manager{
		configManager:  config.NewManager(appName),
		historyManager: NewHistoryManager(appName),
	}
}

// GetConfigManager 获取配置管理器
func (m *Manager) GetConfigManager() types.ConfigManager {
	return m.configManager
}

// AddDownloadRecord 添加下载记录
func (m *Manager) AddDownloadRecord(task interface{}) {
	m.historyManager.AddRecord(task)
}

// GetDownloadHistory 获取下载历史
func (m *Manager) GetDownloadHistory() []interface{} {
	history := m.historyManager.GetHistory()
	// 转换为interface{}切片
	result := make([]interface{}, len(history))
	for i, task := range history {
		result[i] = task
	}
	return result
}

// GetDownloadHistoryTyped 获取类型化的下载历史
func (m *Manager) GetDownloadHistoryTyped() []*task.DownloadTask {
	return m.historyManager.GetHistory()
}

// ClearDownloadHistory 清除下载历史
func (m *Manager) ClearDownloadHistory() {
	m.historyManager.ClearHistory()
}

// GetOutputDir 获取输出目录
func (m *Manager) GetOutputDir() string {
	return m.configManager.GetOutputDir()
}

// SetOutputDir 设置输出目录
func (m *Manager) SetOutputDir(dir string) bool {
	return m.configManager.SetOutputDir(dir)
}

// GetProxy 获取代理设置
func (m *Manager) GetProxy() string {
	return m.configManager.GetProxy()
}

// SetProxy 设置代理
func (m *Manager) SetProxy(proxyURL string) bool {
	return m.configManager.SetProxy(proxyURL)
}

// GetLibraries 获取图书馆列表
func (m *Manager) GetLibraries() []string {
	return m.configManager.GetLibraries()
}

// AddLibrary 添加图书馆
func (m *Manager) AddLibrary(path string) bool {
	return m.configManager.AddLibrary(path)
}
