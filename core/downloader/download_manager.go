package downloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"ImageMaster/core/crawler"
	"ImageMaster/core/types"

	"github.com/google/uuid"
)

// DownloadStatus 表示下载任务状态
type DownloadStatus string

const (
	StatusPending     DownloadStatus = "pending"     // 等待下载
	StatusDownloading DownloadStatus = "downloading" // 下载中
	StatusCompleted   DownloadStatus = "completed"   // 下载完成
	StatusFailed      DownloadStatus = "failed"      // 下载失败
	StatusCancelled   DownloadStatus = "cancelled"   // 已取消
)

// DownloadManager 下载管理器
type DownloadManager struct {
	tasks         map[string]*DownloadTask // 所有任务，包括活跃和历史
	activeTasks   map[string]bool          // 活跃任务集合
	taskCancelMap map[string]chan struct{} // 任务取消通道
	downloader    *Downloader              // 下载器实例
	mu            sync.RWMutex             // 并发控制锁
	history       []*DownloadTask          // 历史任务记录
	storageAPI    interface{}              // 存储API
}

// NewDownloadManager 创建下载管理器
func NewDownloadManager() *DownloadManager {
	return &DownloadManager{
		tasks:         make(map[string]*DownloadTask),
		activeTasks:   make(map[string]bool),
		taskCancelMap: make(map[string]chan struct{}),
		downloader:    NewDownloader(3, 2, true),
		history:       make([]*DownloadTask, 0),
	}
}

// SetConfigManager 设置配置管理器
func (dm *DownloadManager) SetConfigManager(configManager interface{}) {
	if cm, ok := configManager.(types.ConfigProvider); ok {
		dm.downloader.SetConfigManager(cm)
	}
}

// SetStorageAPI 设置存储API
func (dm *DownloadManager) SetStorageAPI(storageAPI interface{}) {
	dm.storageAPI = storageAPI
}

// AddTask 添加下载任务并立即开始下载
func (dm *DownloadManager) AddTask(url string) *DownloadTask {
	dm.mu.Lock()

	// 创建新任务
	task := &DownloadTask{
		ID:        uuid.New().String(),
		URL:       url,
		Status:    string(StatusPending),
		StartTime: time.Now(),
	}

	// 初始化进度
	task.Progress.Current = 0
	task.Progress.Total = 0

	// 添加到任务列表
	dm.tasks[task.ID] = task
	dm.activeTasks[task.ID] = true

	// 创建取消通道
	cancelChan := make(chan struct{})
	dm.taskCancelMap[task.ID] = cancelChan

	dm.mu.Unlock()

	// 异步执行下载任务
	go dm.executeTask(task.ID, cancelChan)

	return task
}

// CrawlWebImages 从网页下载图片，返回任务ID
func (dm *DownloadManager) CrawlWebImages(url string) string {
	// 添加下载任务
	task := dm.AddTask(url)
	return task.ID
}

// executeTask 执行下载任务
func (dm *DownloadManager) executeTask(taskID string, cancelChan chan struct{}) {
	dm.mu.Lock()
	task := dm.tasks[taskID]
	dm.mu.Unlock()

	if task == nil {
		return
	}

	// 更新任务状态为下载中
	dm.updateTaskStatus(taskID, StatusDownloading, "")

	// 异步执行实际下载
	go func() {
		var savedDir string
		var err error

		// 创建一个监听取消信号的通道
		done := make(chan struct{})

		// 在单独的goroutine中执行实际下载
		go func() {
			defer close(done)

			// 使用crawler包处理不同网站的抓取

			// 1. 创建爬虫工厂
			factory := crawler.NewCrawlerFactory(context.Background())

			// 2. 如果有配置管理器，设置到爬虫工厂
			if cm, ok := dm.downloader.GetConfigManager().(types.ConfigProvider); ok {
				factory.SetConfigManager(cm)
			}

			// 3. 检测网站类型并创建相应的爬虫
			siteType := factory.DetectSiteType(task.URL)
			imageCrawler := factory.CreateCrawler(siteType)

			// 4. 设置下载器到爬虫
			imageCrawler.SetDownloader(dm.downloader)

			// 5. 设置下载进度回调函数
			dm.downloader.SetProgressCallback(func(current, total int) {
				dm.updateTaskProgress(taskID, current, total)
			})

			// 获取配置的输出目录
			var outputDir string
			if cm, ok := dm.downloader.GetConfigManager().(types.ConfigProvider); ok {
				outputDir = cm.GetOutputDir()
				if outputDir == "" {
					// 如果配置中未设置输出目录，使用用户home目录下的默认位置
					homeDir, err := os.UserHomeDir()
					if err != nil {
						homeDir = "."
					}
					outputDir = filepath.Join(homeDir, "ImageMaster/downloads")
				}
			} else {
				// 未配置时的默认值
				homeDir, err := os.UserHomeDir()
				if err != nil {
					homeDir = "."
				}
				outputDir = filepath.Join(homeDir, "ImageMaster/downloads")
			}

			// 构建完整保存路径
			savePath := filepath.Join(outputDir)

			// 确保目录存在
			if err = os.MkdirAll(savePath, 0755); err != nil {
				dm.updateTaskStatus(taskID, StatusFailed, fmt.Sprintf("创建目录失败: %v", err))
				return
			}

			// 执行抓取
			savedDir, err = imageCrawler.Crawl(task.URL, savePath)
			if err != nil {
				dm.updateTaskStatus(taskID, StatusFailed, err.Error())
				return
			}

			// 获取实际保存目录
			if savedDir == "" {
				savedDir = savePath // 使用构建的路径作为默认值
			}
		}()

		// 等待下载完成或接收取消信号
		select {
		case <-cancelChan:
			// 任务被取消
			dm.updateTaskStatus(taskID, StatusCancelled, "任务被用户取消")
		case <-done:
			// 下载完成
			if err != nil {
				dm.updateTaskStatus(taskID, StatusFailed, err.Error())
			} else {
				// 更新任务保存路径
				dm.mu.Lock()
				if t, exists := dm.tasks[taskID]; exists {
					t.SavePath = savedDir
				}
				dm.mu.Unlock()

				dm.updateTaskStatus(taskID, StatusCompleted, "")
			}
		}

		// 标记任务完成
		dm.markTaskComplete(taskID)
	}()
}

// updateTaskStatus 更新任务状态
func (dm *DownloadManager) updateTaskStatus(taskID string, status DownloadStatus, errorMsg string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if task, exists := dm.tasks[taskID]; exists {
		task.Status = string(status)

		if errorMsg != "" {
			task.Error = errorMsg
		}

		if status == StatusCompleted || status == StatusFailed || status == StatusCancelled {
			task.CompleteTime = time.Now()
		}
	}
}

// updateTaskProgress 更新任务进度
func (dm *DownloadManager) updateTaskProgress(taskID string, current, total int) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if task, exists := dm.tasks[taskID]; exists {
		task.Progress.Current = current
		task.Progress.Total = total
	}

	fmt.Println("updateTaskProgress", taskID, current, total)
}

// markTaskComplete 标记任务完成，移出活跃任务列表
func (dm *DownloadManager) markTaskComplete(taskID string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// 从活跃任务中移除
	delete(dm.activeTasks, taskID)

	// 关闭取消通道
	if cancelChan, exists := dm.taskCancelMap[taskID]; exists {
		close(cancelChan)
		delete(dm.taskCancelMap, taskID)
	}

	// 添加到历史记录
	if task, exists := dm.tasks[taskID]; exists {
		dm.history = append(dm.history, task)

		// 如果有存储API，保存到永久存储
		if dm.storageAPI != nil {
			if api, ok := dm.storageAPI.(interface{ AddDownloadRecord(*DownloadTask) }); ok {
				api.AddDownloadRecord(task)
			}
		}
	}
}

// CancelTask 取消正在进行的下载任务
func (dm *DownloadManager) CancelTask(taskID string) bool {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// 检查任务是否存在且处于活跃状态
	if _, active := dm.activeTasks[taskID]; !active {
		return false
	}

	// 发送取消信号
	if cancelChan, exists := dm.taskCancelMap[taskID]; exists {
		select {
		case cancelChan <- struct{}{}:
			// 取消信号已发送
		default:
			// 通道已满或已关闭，直接返回
		}
	}

	return true
}

// GetAllTasks 获取所有任务，包括历史任务
func (dm *DownloadManager) GetAllTasks() []*DownloadTask {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	tasks := make([]*DownloadTask, 0, len(dm.tasks))
	for _, task := range dm.tasks {
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartTime.After(tasks[j].StartTime)
	})

	return tasks
}

// GetActiveTasks 获取当前活跃的任务
func (dm *DownloadManager) GetActiveTasks() []*DownloadTask {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	tasks := make([]*DownloadTask, 0, len(dm.activeTasks))
	for id := range dm.activeTasks {
		if task, exists := dm.tasks[id]; exists {
			tasks = append(tasks, task)
		}
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].StartTime.After(tasks[j].StartTime)
	})

	fmt.Println("GetActiveTasks", tasks)
	return tasks
}

// GetHistoryTasks 获取历史任务
func (dm *DownloadManager) GetHistoryTasks() []*DownloadTask {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	// 返回历史任务的副本
	history := make([]*DownloadTask, len(dm.history))
	copy(history, dm.history)

	// 按完成时间倒序排序（最新的任务在前面）
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

		// 倒序排列，所以是timeI在后返回true
		return timeI.After(timeJ)
	})

	return history
}

// ClearHistory 清除历史记录
func (dm *DownloadManager) ClearHistory() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	// 清空历史记录
	dm.history = make([]*DownloadTask, 0)

	// 清除非活跃任务
	for id := range dm.tasks {
		if _, active := dm.activeTasks[id]; !active {
			delete(dm.tasks, id)
		}
	}
}

// GetTaskByID 根据ID获取任务
func (dm *DownloadManager) GetTaskByID(taskID string) *DownloadTask {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if task, exists := dm.tasks[taskID]; exists {
		return task
	}

	return nil
}
