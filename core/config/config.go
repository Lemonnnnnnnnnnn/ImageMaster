package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"ImageMaster/core/logger"
	"ImageMaster/core/types"
)

// 确保Manager实现ConfigProvider和ConfigManager接口
var _ types.ConfigProvider = (*Manager)(nil)
var _ types.ConfigManager = (*Manager)(nil)

// Config 应用配置结构体
type Config struct {
	Libraries []string `json:"libraries"`
	OutputDir string   `json:"output_dir"`
	ProxyURL  string   `json:"proxy_url"`
}

// Manager 配置管理器
type Manager struct {
	config     Config
	configPath string
}

// NewManager 创建新的配置管理器
func NewManager(configName string) *Manager {
	// 设置配置文件路径
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, _ = os.Getwd()
	}
	configPath := filepath.Join(configDir, configName)

	m := &Manager{
		config:     Config{Libraries: []string{}, OutputDir: "", ProxyURL: ""},
		configPath: configPath,
	}

	// 自动加载配置
	m.LoadConfig()
	return m
}

// LoadConfig 加载应用配置
func (m *Manager) LoadConfig() bool {
	data, err := os.ReadFile(m.configPath)
	logger.Debug("Loading config from: %s", m.configPath)
	if err != nil {
		logger.Warn("Failed to load config: %v, using default config", err)
		// 加载失败，使用默认配置
		m.config = Config{Libraries: []string{}, OutputDir: "", ProxyURL: ""}
		return false
	}

	err = json.Unmarshal(data, &m.config)
	if err != nil {
		logger.Error("Failed to parse config: %v, using default config", err)
		m.config = Config{Libraries: []string{}, OutputDir: "", ProxyURL: ""}
		return false
	}

	logger.Info("Config loaded successfully")
	return true
}

// SaveConfig 保存应用配置
func (m *Manager) SaveConfig() bool {
	data, err := json.Marshal(m.config)
	if err != nil {
		logger.Error("Failed to marshal config: %v", err)
		return false
	}

	err = os.WriteFile(m.configPath, data, 0644)
	if err != nil {
		logger.Error("Failed to save config: %v", err)
		return false
	}

	logger.Debug("Config saved successfully")
	return true
}

// GetConfig 获取配置
func (m *Manager) GetConfig() Config {
	return m.config
}

// SetConfig 设置配置
func (m *Manager) SetConfig(config Config) {
	m.config = config
}

// GetLibraries 获取图书馆列表
func (m *Manager) GetLibraries() []string {
	return m.config.Libraries
}

// AddLibrary 添加图书馆
func (m *Manager) AddLibrary(path string) bool {
	// 检查是否已经添加过该库
	for _, lib := range m.config.Libraries {
		if lib == path {
			logger.Warn("Library already exists: %s", path)
			return false
		}
	}

	// 添加到配置中
	m.config.Libraries = append(m.config.Libraries, path)
	logger.Info("Added library: %s", path)
	return m.SaveConfig()
}

// SetOutputDir 设置输出目录
func (m *Manager) SetOutputDir(path string) bool {
	m.config.OutputDir = path
	logger.Info("Set output directory: %s", path)
	return m.SaveConfig()
}

// GetOutputDir 获取输出目录
func (m *Manager) GetOutputDir() string {
	return m.config.OutputDir
}

// SetProxy 设置代理
func (m *Manager) SetProxy(proxyURL string) bool {
	m.config.ProxyURL = proxyURL
	logger.Info("Set proxy: %s", proxyURL)
	return m.SaveConfig()
}

// GetProxy 获取代理设置
func (m *Manager) GetProxy() string {
	return m.config.ProxyURL
}