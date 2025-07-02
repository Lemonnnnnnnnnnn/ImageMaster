package storage

import (
	"ImageMaster/core/download/models"
	"fmt"
)

// StorageAPI 存储管理API
type StorageAPI struct {
	storage *Storage
}

// NewStorageAPI 创建存储API
func NewStorageAPI(configName string) *StorageAPI {
	fmt.Println("NewStorageAPI")
	return &StorageAPI{
		storage: NewStorage(configName),
	}
}

// GetOutputDir 获取输出目录
func (api *StorageAPI) GetOutputDir() string {
	return api.storage.GetConfigManager().GetOutputDir()
}

// SetOutputDir 设置输出目录
func (api *StorageAPI) SetOutputDir(dir string) bool {
	return api.storage.GetConfigManager().SetOutputDir(dir)
}

// GetProxy 获取代理设置
func (api *StorageAPI) GetProxy() string {
	return api.storage.GetConfigManager().GetProxy()
}

// SetProxy 设置代理
func (api *StorageAPI) SetProxy(proxyURL string) bool {
	return api.storage.GetConfigManager().SetProxy(proxyURL)
}

// GetDownloadHistory 获取下载历史
func (api *StorageAPI) GetDownloadHistory() []*models.DownloadTask {
	return api.storage.GetDownloadHistory()
}

// AddDownloadRecord 添加下载记录
func (api *StorageAPI) AddDownloadRecord(task interface{}) {
	api.storage.AddDownloadRecord(task)
}

// ClearDownloadHistory 清除下载历史
func (api *StorageAPI) ClearDownloadHistory() {
	api.storage.ClearDownloadHistory()
}

// GetStorage 获取存储管理器
func (api *StorageAPI) GetStorage() *Storage {
	return api.storage
}

// GetLibraries 获取图书馆列表
func (api *StorageAPI) GetLibraries() []string {
	return api.storage.GetConfigManager().GetLibraries()
}

// AddLibrary 添加图书馆
func (api *StorageAPI) AddLibrary(path string) bool {
	return api.storage.GetConfigManager().AddLibrary(path)
}
