package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// Downloader 下载器
type Downloader struct {
	client      *http.Client
	retryCount  int
	retryDelay  time.Duration
	showProcess bool
	// 添加进度通知回调
	progressCallback func(current, total int)
}

// NewDownloader 创建新的下载器
func NewDownloader(retryCount int, retryDelay int, showProcess bool) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryCount:  retryCount,
		retryDelay:  time.Duration(retryDelay) * time.Second,
		showProcess: showProcess,
	}
}

// SetProgressCallback 设置进度回调函数
func (d *Downloader) SetProgressCallback(callback func(current, total int)) {
	fmt.Printf("下载器: 设置进度回调函数, callback 是否为nil: %v\n", callback == nil)
	d.progressCallback = callback
}

// DownloadFile 下载文件到指定路径
func (d *Downloader) DownloadFile(url string, filepath string, headers map[string]string) error {
	// 确保目录存在
	dir := path.Dir(filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	fmt.Printf("开始下载文件: %s -> %s, 回调是否设置: %v\n", url, filepath, d.progressCallback != nil)

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

		// 创建请求
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = fmt.Errorf("创建请求失败: %w", err)
			continue
		}

		// 设置默认头部
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		// 设置额外头部
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		// 执行请求
		resp, err := d.client.Do(req)
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

// BatchDownload 批量下载文件并报告进度
func (d *Downloader) BatchDownload(urls []string, filepaths []string, headers map[string]string) (int, error) {
	total := len(urls)
	if total == 0 {
		return 0, nil
	}

	if len(filepaths) != total {
		return 0, fmt.Errorf("URL和文件路径数量不匹配")
	}

	successCount := 0

	// 初始化进度
	fmt.Printf("批量下载开始, 总数: %d, 回调是否设置: %v\n", total, d.progressCallback != nil)
	if d.progressCallback != nil {
		d.progressCallback(successCount, total)
	}

	for i, url := range urls {
		if err := d.DownloadFile(url, filepaths[i], headers); err == nil {
			successCount++
			// 更新进度
			if d.progressCallback != nil {
				fmt.Printf("下载进度: %d/%d\n", successCount, total)
				d.progressCallback(successCount, total)
			} else {
				fmt.Printf("警告: 下载进度回调未设置, 无法通知进度 %d/%d\n", successCount, total)
			}
		}
	}

	return successCount, nil
}
