package storage

import (
	"ImageMaster/core/logger"
	"ImageMaster/core/task"
)

// API 存储管理API - 对外提供的统一接口
type API struct {
	manager *Manager
}

// NewAPI 创建存储API
func NewAPI(appName string) *API {
	logger.Info("Creating storage API for app: %s", appName)
	return &API{
		manager: NewManager(appName),
	}
}

// GetOutputDir 获取输出目录
func (api *API) GetOutputDir() string {
	return api.manager.GetOutputDir()
}

// SetOutputDir 设置输出目录
func (api *API) SetOutputDir(dir string) bool {
	return api.manager.SetOutputDir(dir)
}

// GetProxy 获取代理设置
func (api *API) GetProxy() string {
	return api.manager.GetProxy()
}

// SetProxy 设置代理
func (api *API) SetProxy(proxyURL string) bool {
	return api.manager.SetProxy(proxyURL)
}

// GetDownloadHistory 获取下载历史
func (api *API) GetDownloadHistory() []*task.DownloadTask {
	return api.manager.GetDownloadHistoryTyped()
}

// AddDownloadRecord 添加下载记录
func (api *API) AddDownloadRecord(task interface{}) {
	api.manager.AddDownloadRecord(task)
}

// ClearDownloadHistory 清除下载历史
func (api *API) ClearDownloadHistory() {
	api.manager.ClearDownloadHistory()
}

// GetStorage 获取存储管理器（为了兼容性保留）
func (api *API) GetStorage() *Manager {
	return api.manager
}

// GetLibraries 获取图书馆列表
func (api *API) GetLibraries() []string {
	return api.manager.GetLibraries()
}

// AddLibrary 添加图书馆
func (api *API) AddLibrary(path string) bool {
	return api.manager.AddLibrary(path)
}
