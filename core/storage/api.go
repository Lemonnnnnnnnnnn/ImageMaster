package storage

import (
	"ImageMaster/core/logger"
	"ImageMaster/core/types/dto"
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

// GetDownloadHistory 获取下载历史
func (api *API) GetDownloadHistory() []*dto.DownloadTaskDTO {
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
