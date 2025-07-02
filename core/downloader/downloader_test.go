package downloader

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDownloaderWithProxy(t *testing.T) {
	// 创建下载器实例
	d := NewDownloader(3, 2, true)

	// 设置代理
	proxyURL := "http://127.0.0.1:7890"
	err := d.SetProxy(proxyURL)
	if err != nil {
		t.Fatalf("设置代理失败: %v", err)
	}

	// 检查代理是否正确设置
	if d.GetProxy() != proxyURL {
		t.Errorf("代理设置错误, 期望 %s, 实际 %s", proxyURL, d.GetProxy())
	}

	// 设置临时下载目录
	tempDir, err := os.MkdirTemp("", "downloader_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试谷歌网站访问
	t.Run("TestGoogleAccess", func(t *testing.T) {
		url := "https://www.google.com/favicon.ico"
		filePath := filepath.Join(tempDir, "google_favicon.ico")

		err := d.DownloadFile(url, filePath, nil)
		if err != nil {
			t.Errorf("通过代理下载Google文件失败: %v", err)
		} else {
			// 检查文件是否存在且大小大于0
			fileInfo, err := os.Stat(filePath)
			if err != nil || fileInfo.Size() == 0 {
				t.Errorf("下载的文件无效: %v", err)
			}
		}
	})

	// 测试e-hentai网站访问
	t.Run("TestEHentaiAccess", func(t *testing.T) {
		url := "https://e-hentai.org/g/3064128/dc3f9ad7a7/"
		filePath := filepath.Join(tempDir, "ehentai.html")

		err := d.DownloadFile(url, filePath, nil)
		if err != nil {
			t.Errorf("通过代理下载E-Hentai文件失败: %v", err)
		} else {
			// 检查文件是否存在且大小大于0
			fileInfo, err := os.Stat(filePath)
			if err != nil || fileInfo.Size() == 0 {
				t.Errorf("下载的文件无效: %v", err)
			}
		}
	})
}

func TestBatchDownloadWithProxy(t *testing.T) {
	// 创建下载器实例
	d := NewDownloader(3, 2, true)

	// 设置代理
	err := d.SetProxy("http://127.0.0.1:7890")
	if err != nil {
		t.Fatalf("设置代理失败: %v", err)
	}

	// 设置临时下载目录
	tempDir, err := os.MkdirTemp("", "downloader_batch_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)



	// 准备批量下载的URL和路径
	urls := []string{
		"https://www.google.com/favicon.ico",
		"https://e-hentai.org/favicon.ico",
	}

	filepaths := []string{
		filepath.Join(tempDir, "google_favicon.ico"),
		filepath.Join(tempDir, "ehentai_favicon.ico"),
	}

	// 执行批量下载
	successCount, err := d.BatchDownload(urls, filepaths, nil)
	if err != nil {
		t.Errorf("批量下载出错: %v", err)
	}

	// 检查下载结果
	if successCount != len(urls) {
		t.Errorf("批量下载不完整, 成功: %d, 总数: %d", successCount, len(urls))
	}



	// 检查所有文件是否存在
	for i, filepath := range filepaths {
		fileInfo, err := os.Stat(filepath)
		if err != nil {
			t.Errorf("下载的文件 %s 不存在: %v", urls[i], err)
		} else if fileInfo.Size() == 0 {
			t.Errorf("下载的文件 %s 大小为0", urls[i])
		}
	}
}

func TestNoProxyDownload(t *testing.T) {
	// 创建下载器实例（不使用代理）
	d := NewDownloader(3, 2, true)

	// 设置临时下载目录
	tempDir, err := os.MkdirTemp("", "downloader_noproxy_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试国内可访问网站
	t.Run("TestBaiduAccess", func(t *testing.T) {
		url := "https://www.baidu.com/favicon.ico"
		filePath := filepath.Join(tempDir, "baidu_favicon.ico")

		err := d.DownloadFile(url, filePath, nil)
		if err != nil {
			t.Errorf("不使用代理下载百度文件失败: %v", err)
		} else {
			// 检查文件是否存在且大小大于0
			fileInfo, err := os.Stat(filePath)
			if err != nil || fileInfo.Size() == 0 {
				t.Errorf("下载的文件无效: %v", err)
			}
		}
	})
}

func TestDownloadWithHeaders(t *testing.T) {
	// 创建下载器实例
	d := NewDownloader(3, 2, true)

	// 设置代理
	err := d.SetProxy("http://127.0.0.1:7890")
	if err != nil {
		t.Fatalf("设置代理失败: %v", err)
	}

	// 设置临时下载目录
	tempDir, err := os.MkdirTemp("", "downloader_headers_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 准备自定义请求头
	headers := map[string]string{
		"Referer":         "https://www.google.com/",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
	}

	// 使用自定义请求头下载文件
	url := "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"
	filePath := filepath.Join(tempDir, "google_logo.png")

	err = d.DownloadFile(url, filePath, headers)
	if err != nil {
		t.Errorf("使用自定义头下载文件失败: %v", err)
	} else {
		// 检查文件是否存在且大小大于0
		fileInfo, err := os.Stat(filePath)
		if err != nil || fileInfo.Size() == 0 {
			t.Errorf("下载的文件无效: %v", err)
		}
	}
}

func TestRetryBehavior(t *testing.T) {
	// 创建下载器实例，设置重试次数为3，延迟为1秒
	d := NewDownloader(3, 1, true)

	// 设置代理
	err := d.SetProxy("http://127.0.0.1:7890")
	if err != nil {
		t.Fatalf("设置代理失败: %v", err)
	}

	// 设置临时下载目录
	tempDir, err := os.MkdirTemp("", "downloader_retry_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 测试下载不存在的URL，应该尝试重试3次
	startTime := time.Now()
	url := "https://www.google.com/nonexistent_file_123456789.xyz"
	filePath := filepath.Join(tempDir, "should_not_exist.txt")

	err = d.DownloadFile(url, filePath, nil)
	duration := time.Since(startTime)

	// 应该失败，且文件不应存在
	if err == nil {
		t.Errorf("预期下载失败，但成功了")
		if _, err := os.Stat(filePath); err == nil {
			t.Errorf("不应该创建文件，但文件存在")
		}
	}

	// 重试3次，每次间隔1秒，总耗时应该至少3秒左右
	// 考虑到请求和处理耗时，预期总耗时应该在3秒以上
	if duration < 3*time.Second {
		t.Errorf("重试时间过短，可能没有正确重试: %v", duration)
	}
}
