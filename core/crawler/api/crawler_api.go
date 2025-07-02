package api

import (
	"context"
	"sort"
	"sync"

	"ImageMaster/core/download/core"
	"ImageMaster/core/download/models"
	"ImageMaster/core/task"
	"ImageMaster/core/types"
)

// CrawlerAPI 爬虫API接口
type CrawlerAPI struct {
	taskManager   *task.TaskManager
	configManager types.ConfigProvider
	storageAPI    interface{}     // 存储API
	ctx           context.Context // Wails上下文
	sync.Mutex
}

// Config API配置
type Config struct {
	TaskManagerConfig task.Config
}

// NewCrawlerAPI 创建爬虫API
func NewCrawlerAPI(configManager types.ConfigProvider) *CrawlerAPI {
	// 默认配置
	config := Config{
		TaskManagerConfig: task.Config{
			DownloaderConfig: core.Config{
				RetryCount:  3,
				RetryDelay:  2,
				ShowProcess: true,
			},
		},
	}

	api := &CrawlerAPI{
		taskManager:   task.NewTaskManager(config.TaskManagerConfig),
		configManager: configManager,
	}

	// 设置配置管理器
	api.taskManager.SetConfigManager(configManager)

	return api
}

// SetStorageAPI 设置存储API
func (api *CrawlerAPI) SetStorageAPI(storageAPI interface{}) {
	api.storageAPI = storageAPI
	// 同时设置到任务管理器
	api.taskManager.SetStorageAPI(storageAPI)
}

// SetContext 设置Wails上下文
func (api *CrawlerAPI) SetContext(ctx context.Context) {
	api.ctx = ctx
	// 将上下文传递给任务管理器
	api.taskManager.SetContext(ctx)
}

// StartCrawl 开始爬取网页图片
func (api *CrawlerAPI) StartCrawl(url string) string {
	// 调用任务管理器创建爬取任务
	return api.taskManager.CrawlWebImages(url)
}

// CancelCrawl 取消爬取任务
func (api *CrawlerAPI) CancelCrawl(taskID string) bool {
	return api.taskManager.CancelTask(taskID)
}

// GetAllTasks 获取所有任务
func (api *CrawlerAPI) GetAllTasks() []*models.DownloadTask {
	return api.taskManager.GetAllTasks()
}

// GetActiveTasks 获取活跃任务
func (api *CrawlerAPI) GetActiveTasks() []*models.DownloadTask {
	return api.taskManager.GetActiveTasks()
}

// GetHistoryTasks 获取历史任务
func (api *CrawlerAPI) GetHistoryTasks() []*models.DownloadTask {
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
func (api *CrawlerAPI) ClearHistory() {
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
func (api *CrawlerAPI) GetTaskByID(taskID string) *models.DownloadTask {
	return api.taskManager.GetTaskByID(taskID)
}

// GetTaskProgress 获取任务进度
func (api *CrawlerAPI) GetTaskProgress(taskID string) map[string]interface{} {
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
