package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"ImageMaster/core/download/models"
	"ImageMaster/core/logger"
)

// Manager 历史记录管理器
type Manager struct {
	dataDir         string
	mu              sync.RWMutex
	downloadHistory []*models.DownloadTask
}

// NewManager 创建历史记录管理器
func NewManager(appName string) *Manager {
	// 获取数据目录
	userHome, err := os.UserHomeDir()
	if err != nil {
		userHome = "."
	}

	// 数据目录
	dataDir := filepath.Join(userHome, "."+appName)

	// 确保目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logger.Error("Failed to create data directory: %v", err)
	}

	// 创建管理器实例
	m := &Manager{
		dataDir:         dataDir,
		downloadHistory: make([]*models.DownloadTask, 0),
	}

	// 加载历史记录
	m.loadHistory()

	return m
}

// AddRecord 添加下载记录
func (m *Manager) AddRecord(task interface{}) {
	// 类型断言
	downloadTask, ok := task.(*models.DownloadTask)
	if !ok {
		logger.Warn("Invalid task type for download history")
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// 添加到历史记录
	m.downloadHistory = append(m.downloadHistory, downloadTask)
	logger.Debug("Added download record: %s", downloadTask.Name)

	// 保存历史记录
	m.saveHistory()
}

// GetHistory 获取下载历史
func (m *Manager) GetHistory() []*models.DownloadTask {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回历史记录的副本
	history := make([]*models.DownloadTask, len(m.downloadHistory))
	copy(history, m.downloadHistory)

	return history
}

// ClearHistory 清除下载历史
func (m *Manager) ClearHistory() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 清空历史记录
	m.downloadHistory = make([]*models.DownloadTask, 0)
	logger.Info("Cleared download history")

	// 保存历史记录
	m.saveHistory()
}

// saveHistory 保存下载历史到文件
func (m *Manager) saveHistory() {
	historyPath := filepath.Join(m.dataDir, "download_history.json")

	// 将历史记录序列化为JSON
	data, err := json.MarshalIndent(m.downloadHistory, "", "  ")
	if err != nil {
		logger.Error("Failed to serialize download history: %v", err)
		return
	}

	// 保存到文件
	if err := os.WriteFile(historyPath, data, 0644); err != nil {
		logger.Error("Failed to save download history: %v", err)
		return
	}

	logger.Debug("Download history saved")
}

// loadHistory 从文件加载下载历史
func (m *Manager) loadHistory() {
	historyPath := filepath.Join(m.dataDir, "download_history.json")

	// 检查文件是否存在
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		logger.Debug("No existing download history found")
		return
	}

	// 读取文件
	data, err := os.ReadFile(historyPath)
	if err != nil {
		logger.Error("Failed to read download history: %v", err)
		return
	}

	// 反序列化JSON
	if err := json.Unmarshal(data, &m.downloadHistory); err != nil {
		logger.Error("Failed to parse download history: %v", err)
		return
	}

	logger.Info("Loaded %d download history records", len(m.downloadHistory))
}