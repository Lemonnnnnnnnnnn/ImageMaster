package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"ImageMaster/core/request"
	"ImageMaster/core/types"
)

// DownloadTask 下载任务
type DownloadTask struct {
	ID           string    `json:"id"`           // 任务ID
	URL          string    `json:"url"`          // 下载URL
	Status       string    `json:"status"`       // 状态: pending, downloading, completed, failed
	SavePath     string    `json:"savePath"`     // 保存路径
	StartTime    time.Time `json:"startTime"`    // 开始时间
	CompleteTime time.Time `json:"completeTime"` // 完成时间
	Error        string    `json:"error"`        // 错误信息
	Progress     struct {
		Current int `json:"current"` // 当前已下载项目数
		Total   int `json:"total"`   // 总项目数
	} `json:"progress"` // 下载进度
}

// ProgressCallback 定义进度回调函数类型
type ProgressCallback func(current, total int)

// Downloader 下载器
type Downloader struct {
	reqClient        *request.Client
	retryCount       int
	retryDelay       time.Duration
	showProcess      bool
	configManager    types.ConfigProvider
	progressCallback ProgressCallback // 进度回调函数
}

// NewDownloader 创建新的下载器
func NewDownloader(retryCount int, retryDelay int, showProcess bool) *Downloader {
	return &Downloader{
		reqClient:   request.NewClient(),
		retryCount:  retryCount,
		retryDelay:  time.Duration(retryDelay) * time.Second,
		showProcess: showProcess,
	}
}

// SetConfigManager 设置配置管理器
func (d *Downloader) SetConfigManager(configManager types.ConfigProvider) {
	d.configManager = configManager

	// 将配置管理器传递给请求客户端
	d.reqClient.SetConfigManager(configManager)
}

// SetProxy 设置代理
func (d *Downloader) SetProxy(proxyURL string) error {
	// 使用请求客户端设置代理
	return d.reqClient.SetProxy(proxyURL)
}

// GetProxy 获取当前代理设置
func (d *Downloader) GetProxy() string {
	return d.reqClient.GetProxy()
}

// GetConfigManager 获取配置管理器
func (d *Downloader) GetConfigManager() interface{} {
	return d.configManager
}

// DownloadFile 下载文件到指定路径
func (d *Downloader) DownloadFile(url string, filepath string, headers map[string]string) error {
	// 确保目录存在
	dir := path.Dir(filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 直接创建最终文件
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer out.Close()

	// 执行下载
	success := false
	var lastErr error
	for attempt := 0; attempt <= d.retryCount; attempt++ {
		if attempt > 0 {
			fmt.Printf("重试下载 %s (第 %d 次)\n", url, attempt)
			time.Sleep(d.retryDelay)
		}

		// 设置请求头
		if headers != nil {
			d.reqClient.SetHeaders(headers)
		}

		// 执行请求
		resp, err := d.reqClient.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("请求失败: %w", err)
			continue
		}

		// 检查状态码
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("状态码错误: %d", resp.StatusCode)
			continue
		}

		// 清空文件内容
		if _, err := out.Seek(0, 0); err != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("文件定位失败: %w", err)
			continue
		}
		if err := out.Truncate(0); err != nil {
			resp.Body.Close()
			lastErr = fmt.Errorf("清空文件失败: %w", err)
			continue
		}

		// 复制数据
		_, err = io.Copy(out, resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("数据写入失败: %w", err)
			continue
		}

		success = true
		break
	}

	if !success {
		// 下载失败时删除文件
		os.Remove(filepath)
		return fmt.Errorf("下载失败: %w", lastErr)
	}

	return nil
}

// BatchDownload 批量下载文件
func (d *Downloader) BatchDownload(urls []string, filepaths []string, headers map[string]string) (int, error) {
	total := len(urls)
	if total == 0 {
		return 0, nil
	}

	if len(filepaths) != total {
		return 0, fmt.Errorf("URL和文件路径数量不匹配")
	}

	successCount := 0

	for i, url := range urls {
		if err := d.DownloadFile(url, filepaths[i], headers); err == nil {
			successCount++

			// 下载成功后调用进度回调
			if d.progressCallback != nil {
				d.progressCallback(successCount, total)
			}
		}
	}

	return successCount, nil
}

// SetProgressCallback 设置进度回调函数
func (d *Downloader) SetProgressCallback(callback func(current, total int)) {
	d.progressCallback = callback
}
