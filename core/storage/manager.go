package storage

import (
	"ImageMaster/core/logger"
	"ImageMaster/core/task"
	"ImageMaster/core/types"
)

// 确保Manager实现StorageProvider接口
var _ types.StorageProvider = (*Manager)(nil)

// Manager 存储管理器 - 统一的存储接口
type Manager struct {
	historyManager *HistoryManager
}

// NewManager 创建存储管理器
func NewManager(appName string) *Manager {
	logger.Info("Initializing storage manager for app: %s", appName)

	return &Manager{
		historyManager: NewHistoryManager(appName),
	}
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
