package downloader

import (
	"context"
	"sort"
	"sync"

	"ImageMaster/core/types"
)

// DownloaderAPI 是下载器的API接口
type DownloaderAPI struct {
	manager       *DownloadManager
	configManager types.ConfigProvider
	storageAPI    interface{} // 存储API
	ctx           context.Context // Wails上下文
	sync.Mutex
}

// NewDownloaderAPI 创建下载器API
func NewDownloaderAPI(configManager types.ConfigProvider) *DownloaderAPI {
	api := &DownloaderAPI{
		manager:       NewDownloadManager(),
		configManager: configManager,
	}

	// 设置配置管理器
	api.manager.SetConfigManager(configManager)

	return api
}

// SetStorageAPI 设置存储API
func (api *DownloaderAPI) SetStorageAPI(storageAPI interface{}) {
	api.storageAPI = storageAPI

	// 同时设置到下载管理器
	api.manager.SetStorageAPI(storageAPI)
}

// SetContext 设置Wails上下文
func (api *DownloaderAPI) SetContext(ctx context.Context) {
	api.ctx = ctx
	// 将上下文传递给下载管理器
	api.manager.SetContext(ctx)
}

// StartDownload 开始下载网页图片
func (api *DownloaderAPI) StartDownload(url string) string {
	// 调用下载管理器创建下载任务
	return api.manager.CrawlWebImages(url)
}

// CancelDownload 取消下载任务
func (api *DownloaderAPI) CancelDownload(taskID string) bool {
	return api.manager.CancelTask(taskID)
}

// GetAllTasks 获取所有任务
func (api *DownloaderAPI) GetAllTasks() []*DownloadTask {
	return api.manager.GetAllTasks()
}

// GetActiveTasks 获取活跃任务
func (api *DownloaderAPI) GetActiveTasks() []*DownloadTask {
	return api.manager.GetActiveTasks()
}

// GetHistoryTasks 获取历史任务
func (api *DownloaderAPI) GetHistoryTasks() []*DownloadTask {
	// 如果有存储API，优先从存储获取
	if api.storageAPI != nil {
		if storage, ok := api.storageAPI.(interface{ GetDownloadHistory() []*DownloadTask }); ok {
			tasks := storage.GetDownloadHistory()

			// 确保从存储获取的任务也按时间倒序排序
			if len(tasks) > 0 {
				sort.Slice(tasks, func(i, j int) bool {
					// 如果completeTime为空，使用startTime
					timeI := tasks[i].CompleteTime
					if timeI.IsZero() {
						timeI = tasks[i].StartTime
					}

					timeJ := tasks[j].CompleteTime
					if timeJ.IsZero() {
						timeJ = tasks[j].StartTime
					}

					// 倒序排列
					return timeI.After(timeJ)
				})
			}

			return tasks
		}
	}

	// 否则从下载管理器获取
	return api.manager.GetHistoryTasks()
}

// ClearHistory 清除历史记录
func (api *DownloaderAPI) ClearHistory() {
	// 如果有存储API，同时清除存储中的历史
	if api.storageAPI != nil {
		if storage, ok := api.storageAPI.(interface{ ClearDownloadHistory() }); ok {
			storage.ClearDownloadHistory()
		}
	}

	// 清除下载管理器中的历史
	api.manager.ClearHistory()
}

// GetTaskByID 根据ID获取任务
func (api *DownloaderAPI) GetTaskByID(taskID string) *DownloadTask {
	return api.manager.GetTaskByID(taskID)
}

// GetTaskProgress 获取任务进度
func (api *DownloaderAPI) GetTaskProgress(taskID string) map[string]interface{} {
	task := api.manager.GetTaskByID(taskID)
	if task == nil {
		return nil
	}

	return map[string]interface{}{
		"id":      task.ID,
		"status":  task.Status,
		"current": task.Progress.Current,
		"total":   task.Progress.Total,
		"percent": calculatePercent(task.Progress.Current, task.Progress.Total),
	}
}

// 计算百分比
func calculatePercent(current, total int) int {
	if total <= 0 {
		return 0
	}
	return int((float64(current) / float64(total)) * 100)
}
