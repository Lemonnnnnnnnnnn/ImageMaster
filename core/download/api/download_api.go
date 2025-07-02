package api

import (
	"context"
	"sort"
	"sync"

	"ImageMaster/core/download/core"
	"ImageMaster/core/download/manager"
	"ImageMaster/core/download/models"
	"ImageMaster/core/types"
)

// DownloadAPI 下载API接口
type DownloadAPI struct {
	taskManager   *manager.TaskManager
	configManager types.ConfigProvider
	storageAPI    interface{} // 存储API
	ctx           context.Context // Wails上下文
	sync.Mutex
}

// Config API配置
type Config struct {
	TaskManagerConfig manager.Config
}

// NewDownloadAPI 创建下载API
func NewDownloadAPI(configManager types.ConfigProvider) *DownloadAPI {
	// 默认配置
	config := Config{
		TaskManagerConfig: manager.Config{
			DownloaderConfig: core.Config{
				RetryCount:  3,
				RetryDelay:  2,
				ShowProcess: true,
			},
		},
	}

	api := &DownloadAPI{
		taskManager:   manager.NewTaskManager(config.TaskManagerConfig),
		configManager: configManager,
	}

	// 设置配置管理器
	api.taskManager.SetConfigManager(configManager)

	return api
}

// SetStorageAPI 设置存储API
func (api *DownloadAPI) SetStorageAPI(storageAPI interface{}) {
	api.storageAPI = storageAPI
	// 同时设置到任务管理器
	api.taskManager.SetStorageAPI(storageAPI)
}

// SetContext 设置Wails上下文
func (api *DownloadAPI) SetContext(ctx context.Context) {
	api.ctx = ctx
	// 将上下文传递给任务管理器
	api.taskManager.SetContext(ctx)
}

// StartDownload 开始下载网页图片
func (api *DownloadAPI) StartDownload(url string) string {
	// 调用任务管理器创建下载任务
	return api.taskManager.CrawlWebImages(url)
}

// CancelDownload 取消下载任务
func (api *DownloadAPI) CancelDownload(taskID string) bool {
	return api.taskManager.CancelTask(taskID)
}

// GetAllTasks 获取所有任务
func (api *DownloadAPI) GetAllTasks() []*models.DownloadTask {
	return api.taskManager.GetAllTasks()
}

// GetActiveTasks 获取活跃任务
func (api *DownloadAPI) GetActiveTasks() []*models.DownloadTask {
	return api.taskManager.GetActiveTasks()
}

// GetHistoryTasks 获取历史任务
func (api *DownloadAPI) GetHistoryTasks() []*models.DownloadTask {
	// 如果有存储API，优先从存储获取
	if api.storageAPI != nil {
		if storage, ok := api.storageAPI.(interface{ GetDownloadHistory() []*models.DownloadTask }); ok {
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

	// 否则从任务管理器获取
	return api.taskManager.GetHistoryTasks()
}

// ClearHistory 清除历史记录
func (api *DownloadAPI) ClearHistory() {
	// 如果有存储API，同时清除存储中的历史
	if api.storageAPI != nil {
		if storage, ok := api.storageAPI.(interface{ ClearDownloadHistory() }); ok {
			storage.ClearDownloadHistory()
		}
	}

	// 清除任务管理器中的历史
	api.taskManager.ClearHistory()
}

// GetTaskByID 根据ID获取任务
func (api *DownloadAPI) GetTaskByID(taskID string) *models.DownloadTask {
	return api.taskManager.GetTaskByID(taskID)
}

// GetTaskProgress 获取任务进度
func (api *DownloadAPI) GetTaskProgress(taskID string) map[string]interface{} {
	task := api.taskManager.GetTaskByID(taskID)
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