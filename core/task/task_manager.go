package task

import (
	"context"
	"sort"
	"sync"
	"time"

	"ImageMaster/core/crawler"
	"ImageMaster/core/download/core"
	"ImageMaster/core/download/models"
	"ImageMaster/core/types"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TaskManager 任务管理器
type TaskManager struct {
	tasks         map[string]*models.DownloadTask // 所有任务，包括活跃和历史
	activeTasks   map[string]bool                 // 活跃任务集合
	taskCancelMap map[string]chan struct{}        // 任务取消通道
	downloaders   map[string]*core.Downloader     // 每个任务对应的下载器实例
	defaultConfig core.Config                     // 默认下载器配置
	mu            sync.RWMutex                    // 并发控制锁
	history       []*models.DownloadTask          // 历史任务记录
	storageAPI    interface{}                     // 存储API
	ctx           context.Context                 // Wails上下文
	configManager types.ConfigProvider            // 配置管理器
}

// Config 任务管理器配置
type Config struct {
	DownloaderConfig core.Config
}

// NewTaskManager 创建任务管理器
func NewTaskManager(config Config) *TaskManager {
	return &TaskManager{
		tasks:         make(map[string]*models.DownloadTask),
		activeTasks:   make(map[string]bool),
		taskCancelMap: make(map[string]chan struct{}),
		downloaders:   make(map[string]*core.Downloader),
		defaultConfig: config.DownloaderConfig,
		history:       make([]*models.DownloadTask, 0),
	}
}

// SetConfigManager 设置配置管理器
func (tm *TaskManager) SetConfigManager(configManager types.ConfigProvider) {
	tm.configManager = configManager
}

// SetStorageAPI 设置存储API
func (tm *TaskManager) SetStorageAPI(storageAPI interface{}) {
	tm.storageAPI = storageAPI
}

// SetContext 设置Wails上下文
func (tm *TaskManager) SetContext(ctx context.Context) {
	tm.ctx = ctx
}

// AddTask 添加下载任务并立即开始下载
func (tm *TaskManager) AddTask(url string) *models.DownloadTask {
	tm.mu.Lock()

	// 创建新任务
	now := time.Now()
	task := &models.DownloadTask{
		ID:        uuid.New().String(),
		URL:       url,
		Status:    string(models.StatusPending),
		StartTime: now,
		UpdatedAt: now,
	}

	// 初始化进度
	task.Progress.Current = 0
	task.Progress.Total = 0

	// 添加到任务列表
	tm.tasks[task.ID] = task
	tm.activeTasks[task.ID] = true

	// 创建取消通道
	cancelChan := make(chan struct{})
	tm.taskCancelMap[task.ID] = cancelChan

	tm.mu.Unlock()

	// 异步执行下载任务
	go tm.executeTask(task.ID, cancelChan)

	return task
}

// CrawlWebImages 从网页下载图片，返回任务ID
func (tm *TaskManager) CrawlWebImages(url string) string {
	// 添加下载任务
	task := tm.AddTask(url)
	return task.ID
}

// createDownloaderForTask 为任务创建专用的下载器实例
func (tm *TaskManager) createDownloaderForTask(taskID string) *core.Downloader {
	// 创建新的下载器实例
	newDownloader := core.NewDownloader(tm.defaultConfig)

	// 复制配置
	if tm.configManager != nil {
		newDownloader.SetConfigManager(tm.configManager)
	}

	// 创建并设置TaskUpdater
	taskUpdater := NewTaskUpdater(taskID, tm)
	newDownloader.SetTaskUpdater(taskUpdater)

	// 保存到下载器映射
	tm.downloaders[taskID] = newDownloader

	return newDownloader
}

// executeTask 执行下载任务
func (tm *TaskManager) executeTask(taskID string, cancelChan chan struct{}) {
	defer func() {
		tm.mu.Lock()
		delete(tm.activeTasks, taskID)
		delete(tm.taskCancelMap, taskID)
		delete(tm.downloaders, taskID)
		tm.mu.Unlock()
	}()

	// 获取任务
	tm.mu.RLock()
	task, exists := tm.tasks[taskID]
	if !exists {
		tm.mu.RUnlock()
		return
	}
	tm.mu.RUnlock()

	// 更新任务状态为下载中
	tm.UpdateTask(taskID, func(task *models.DownloadTask) {
		task.Status = string(models.StatusDownloading)
		task.UpdatedAt = time.Now()
	})

	// 创建下载器
	downloader := tm.createDownloaderForTask(taskID)

	// 创建爬虫工厂
	crawlerFactory := crawler.NewCrawlerFactory(tm.ctx)
	if tm.configManager != nil {
		crawlerFactory.SetConfigManager(tm.configManager)
	}

	// 检测网站类型并创建对应的爬虫
	siteType := crawlerFactory.DetectSiteType(task.URL)
	crawlerInstance := crawlerFactory.CreateCrawler(siteType)
	crawlerInstance.SetDownloader(downloader)

	// 设置输出目录
	var outputDir string
	if tm.configManager != nil {
		outputDir = tm.configManager.GetOutputDir()
	} else {
		outputDir = "downloads"
	}

	// 执行爬取
	savePath, err := crawlerInstance.Crawl(task.URL, outputDir)
	if err != nil {
		// 下载失败
		tm.UpdateTask(taskID, func(task *models.DownloadTask) {
			task.Status = string(models.StatusFailed)
			task.Error = err.Error()
			task.CompleteTime = time.Now()
			task.UpdatedAt = time.Now()
		})
	} else {
		// 下载成功
		tm.UpdateTask(taskID, func(task *models.DownloadTask) {
			task.Status = string(models.StatusCompleted)
			task.SavePath = savePath
			task.CompleteTime = time.Now()
			task.UpdatedAt = time.Now()
		})
	}

	// 将任务移动到历史记录
	tm.moveTaskToHistory(taskID)

	// 发送完成事件到前端
	if tm.ctx != nil {
		runtime.EventsEmit(tm.ctx, "download:completed", map[string]interface{}{
			"taskId": taskID,
			"name":   task.Name,
			"status": task.Status,
		})
	}
}

// moveTaskToHistory 将任务移动到历史记录
func (tm *TaskManager) moveTaskToHistory(taskID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if task, exists := tm.tasks[taskID]; exists {
		// 添加到历史记录
		tm.history = append(tm.history, task)

		// 如果有存储API，保存到存储
		if tm.storageAPI != nil {
			if storage, ok := tm.storageAPI.(interface{ AddDownloadRecord(task interface{}) }); ok {
				storage.AddDownloadRecord(task)
			}
		}
	}
}

// UpdateTask 更新任务
func (tm *TaskManager) UpdateTask(taskID string, updateFunc func(task *models.DownloadTask)) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if task, exists := tm.tasks[taskID]; exists {
		updateFunc(task)
		task.UpdatedAt = time.Now()
	}
}

// UpdateTaskProgress 更新任务进度
func (tm *TaskManager) UpdateTaskProgress(taskID string, current, total int) {
	tm.UpdateTask(taskID, func(task *models.DownloadTask) {
		task.Progress.Current = current
		task.Progress.Total = total
	})
}

// CancelTask 取消任务
func (tm *TaskManager) CancelTask(taskID string) bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if cancelChan, exists := tm.taskCancelMap[taskID]; exists {
		close(cancelChan)
		// 更新任务状态
		if task, exists := tm.tasks[taskID]; exists {
			task.Status = string(models.StatusCancelled)
			task.UpdatedAt = time.Now()
		}
		return true
	}
	return false
}

// GetTaskByID 根据ID获取任务
func (tm *TaskManager) GetTaskByID(taskID string) *models.DownloadTask {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.tasks[taskID]
}

// GetAllTasks 获取所有任务
func (tm *TaskManager) GetAllTasks() []*models.DownloadTask {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tasks := make([]*models.DownloadTask, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}

	// 按时间倒序排序
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartTime.After(tasks[j].StartTime)
	})

	return tasks
}

// GetActiveTasks 获取活跃任务
func (tm *TaskManager) GetActiveTasks() []*models.DownloadTask {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tasks := make([]*models.DownloadTask, 0)
	for taskID := range tm.activeTasks {
		if task, exists := tm.tasks[taskID]; exists {
			tasks = append(tasks, task)
		}
	}

	// 按时间倒序排序
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartTime.After(tasks[j].StartTime)
	})

	return tasks
}

// GetHistoryTasks 获取历史任务
func (tm *TaskManager) GetHistoryTasks() []*models.DownloadTask {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	// 复制历史记录
	history := make([]*models.DownloadTask, len(tm.history))
	copy(history, tm.history)

	// 按完成时间倒序排序
	sort.Slice(history, func(i, j int) bool {
		// 如果completeTime为空，使用startTime
		timeI := history[i].CompleteTime
		if timeI.IsZero() {
			timeI = history[i].StartTime
		}

		timeJ := history[j].CompleteTime
		if timeJ.IsZero() {
			timeJ = history[j].StartTime
		}

		// 倒序排列
		return timeI.After(timeJ)
	})

	return history
}

// ClearHistory 清除历史记录
func (tm *TaskManager) ClearHistory() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 清空历史记录
	tm.history = make([]*models.DownloadTask, 0)

	// 清除非活跃任务
	for id := range tm.tasks {
		if _, active := tm.activeTasks[id]; !active {
			delete(tm.tasks, id)
		}
	}
}
