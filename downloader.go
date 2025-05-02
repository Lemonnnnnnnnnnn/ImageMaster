package main

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

// DownloadFile 下载文件到指定路径
func (d *Downloader) DownloadFile(url string, filepath string, headers map[string]string) error {
	// 确保目录存在
	dir := path.Dir(filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建临时文件
	tmpFile := filepath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %w", err)
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
		os.Remove(tmpFile)
		return fmt.Errorf("下载失败: %w", lastErr)
	}

	// 重命名临时文件
	if err := os.Rename(tmpFile, filepath); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("重命名文件失败: %w", err)
	}

	return nil
}

// 更新 Getter 中的下载方法
func (g *Getter) DownloadImage(url string, savePath string) error {
	downloader := NewDownloader(
		g.config.MaxRetries,
		g.config.RetryDelay,
		false,
	)
	return downloader.DownloadFile(url, savePath, nil)
}
