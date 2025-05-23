package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"ImageMaster/core/types"
)

// 确保ConfigManager实现ConfigProvider接口
var _ types.ConfigProvider = (*ConfigManager)(nil)

// Config 应用配置结构体
type Config struct {
	Libraries []string `json:"libraries"`
	OutputDir string   `json:"output_dir"`
	ProxyURL  string   `json:"proxy_url"`
}

// ConfigManager 配置管理器
type ConfigManager struct {
	config     Config
	configPath string
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager(configName string) *ConfigManager {
	// 设置配置文件路径
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, _ = os.Getwd()
	}
	configPath := filepath.Join(configDir, configName)

	return &ConfigManager{
		config:     Config{Libraries: []string{}, OutputDir: "", ProxyURL: ""},
		configPath: configPath,
	}
}

// LoadConfig 加载应用配置
func (m *ConfigManager) LoadConfig() bool {
	data, err := os.ReadFile(m.configPath)
	fmt.Println("LoadConfig", m.configPath)
	if err != nil {
		fmt.Println("LoadConfig", err)
		// 加载失败，使用默认配置
		m.config = Config{Libraries: []string{}, OutputDir: "", ProxyURL: ""}
		return false
	}

	err = json.Unmarshal(data, &m.config)
	if err != nil {
		fmt.Println("LoadConfig", err)
		m.config = Config{Libraries: []string{}, OutputDir: "", ProxyURL: ""}
		return false
	}

	return true
}

// SaveConfig 保存应用配置
func (m *ConfigManager) SaveConfig() bool {
	data, err := json.Marshal(m.config)
	if err != nil {
		return false
	}

	err = os.WriteFile(m.configPath, data, 0644)
	if err != nil {
		return false
	}

	return true
}

// GetConfig 获取配置
func (m *ConfigManager) GetConfig() Config {
	return m.config
}

// SetConfig 设置配置
func (m *ConfigManager) SetConfig(config Config) {
	m.config = config
}

// GetLibraries 获取图书馆列表
func (m *ConfigManager) GetLibraries() []string {
	return m.config.Libraries
}

// AddLibrary 添加图书馆
func (m *ConfigManager) AddLibrary(path string) bool {
	// 检查是否已经添加过该库
	for _, lib := range m.config.Libraries {
		if lib == path {
			return false
		}
	}

	// 添加到配置中
	m.config.Libraries = append(m.config.Libraries, path)
	return m.SaveConfig()
}

// SetOutputDir 设置输出目录
func (m *ConfigManager) SetOutputDir(path string) bool {
	m.config.OutputDir = path
	return m.SaveConfig()
}

// GetOutputDir 获取输出目录
func (m *ConfigManager) GetOutputDir() string {
	return m.config.OutputDir
}

// SetProxy 设置代理
func (m *ConfigManager) SetProxy(proxyURL string) bool {
	m.config.ProxyURL = proxyURL
	return m.SaveConfig()
}

// GetProxy 获取代理设置
func (m *ConfigManager) GetProxy() string {
	return m.config.ProxyURL
}
