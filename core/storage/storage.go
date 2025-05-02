package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ImageMaster/core/downloader"
)

// Storage 存储管理器
type Storage struct {
	configManager   *ConfigManager
	dataDir         string
	mu              sync.RWMutex
	downloadHistory []*downloader.DownloadTask
}

// NewStorage 创建存储管理器
func NewStorage(configName string) *Storage {
	fmt.Println("NewStorage")
	// 创建配置管理器
	configManager := NewConfigManager(configName)
	configManager.LoadConfig()

	// 获取数据目录
	userHome, err := os.UserHomeDir()
	if err != nil {
		userHome = "."
	}

	// 数据目录
	dataDir := filepath.Join(userHome, "."+configName)

	// 确保目录存在
	os.MkdirAll(dataDir, 0755)

	// 创建存储实例
	storage := &Storage{
		configManager:   configManager,
		dataDir:         dataDir,
		downloadHistory: make([]*downloader.DownloadTask, 0),
	}

	// 加载历史记录
	storage.loadDownloadHistory()

	return storage
}

// GetConfigManager 获取配置管理器
func (s *Storage) GetConfigManager() *ConfigManager {
	return s.configManager
}

// AddDownloadRecord 添加下载记录
func (s *Storage) AddDownloadRecord(task *downloader.DownloadTask) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 添加到历史记录
	s.downloadHistory = append(s.downloadHistory, task)

	// 保存历史记录
	s.saveDownloadHistory()
}

// GetDownloadHistory 获取下载历史
func (s *Storage) GetDownloadHistory() []*downloader.DownloadTask {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 返回历史记录的副本
	history := make([]*downloader.DownloadTask, len(s.downloadHistory))
	copy(history, s.downloadHistory)

	return history
}

// ClearDownloadHistory 清除下载历史
func (s *Storage) ClearDownloadHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 清空历史记录
	s.downloadHistory = make([]*downloader.DownloadTask, 0)

	// 保存历史记录
	s.saveDownloadHistory()
}

// saveDownloadHistory 保存下载历史到文件
func (s *Storage) saveDownloadHistory() {
	historyPath := filepath.Join(s.dataDir, "download_history.json")

	// 将历史记录序列化为JSON
	data, err := json.MarshalIndent(s.downloadHistory, "", "  ")
	if err != nil {
		fmt.Printf("序列化下载历史失败: %v\n", err)
		return
	}

	// 保存到文件
	if err := os.WriteFile(historyPath, data, 0644); err != nil {
		fmt.Printf("保存下载历史失败: %v\n", err)
	}
}

// loadDownloadHistory 从文件加载下载历史
func (s *Storage) loadDownloadHistory() {
	historyPath := filepath.Join(s.dataDir, "download_history.json")

	// 检查文件是否存在
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return
	}

	// 读取文件
	data, err := os.ReadFile(historyPath)
	if err != nil {
		fmt.Printf("读取下载历史失败: %v\n", err)
		return
	}

	// 反序列化JSON
	if err := json.Unmarshal(data, &s.downloadHistory); err != nil {
		fmt.Printf("解析下载历史失败: %v\n", err)
	}
}

// 序列化辅助函数
type serializedTask struct {
	ID           string    `json:"id"`
	URL          string    `json:"url"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	SavePath     string    `json:"savePath"`
	StartTime    time.Time `json:"startTime"`
	CompleteTime time.Time `json:"completeTime"`
	Error        string    `json:"error"`
	Progress     struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progress"`
}
