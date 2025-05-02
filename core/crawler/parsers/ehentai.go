package parsers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"ImageMaster/core/downloader"
)

const PARALLEL = 5

// EHentaiAlbum EH专辑
type EHentaiAlbum struct {
	Name  string
	Pages []string
}

// ParseEHentai 解析EH网站
func ParseEHentai(ctx context.Context, client *http.Client, url string, savePath string, dl *downloader.Downloader) error {
	fmt.Printf("下载 eHentai 专辑: %s\n", url)

	eHentaiAlbum, err := GetAlbum(client, url)
	if err != nil {
		return fmt.Errorf("获取专辑失败: %w", err)
	}

	// 创建信号量来控制并发
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, PARALLEL)

	// 统计结果
	var successMutex sync.Mutex
	successImages := 0

	// 使用传入的下载器或创建新的下载器
	localDownloader := dl
	if localDownloader == nil {
		localDownloader = downloader.NewDownloader(3, 3, true)
		fmt.Printf("EHentai解析器创建了新的下载器: %p\n", localDownloader)
	} else {
		fmt.Printf("EHentai解析器使用传入的下载器: %p\n", localDownloader)
	}

	// 保存路径
	albumPath := savePath + "/" + eHentaiAlbum.Name

	// 批量下载URL和路径
	var imgURLs []string
	var filePaths []string

	// 遍历每一页
	for pageIndex, page := range eHentaiAlbum.Pages {
		links := ParseLinks(page)

		for linkIndex, link := range links {
			link := link // 防止闭包问题

			wg.Add(1)
			semaphore <- struct{}{} // 获取信号量

			go func(pageIndex, linkIndex int, linkURL string) {
				defer func() {
					<-semaphore // 释放信号量
					wg.Done()
				}()

				// 解析页面获取真实图片URL
				imgURL, err := ParsePage(client, linkURL)
				if err != nil {
					fmt.Printf("解析页面失败 %s: %v\n", linkURL, err)
					return
				}

				// 构建保存文件名
				filename := fmt.Sprintf("%d_%d.jpg", pageIndex, linkIndex)
				fullPath := fmt.Sprintf("%s/%s", albumPath, filename)

				successMutex.Lock()
				imgURLs = append(imgURLs, imgURL)
				filePaths = append(filePaths, fullPath)
				successMutex.Unlock()

				// 随机休眠1到3秒防止被ban
				sleepDuration := time.Duration(1+rand.Intn(3)) * time.Second
				time.Sleep(sleepDuration)
			}(pageIndex, linkIndex, link)
		}
	}

	// 等待所有URL收集任务完成
	wg.Wait()

	// 批量下载所有图片并跟踪进度
	successImages, err = localDownloader.BatchDownload(imgURLs, filePaths, nil)
	if err != nil {
		fmt.Printf("批量下载出错: %v\n", err)
	}

	fmt.Printf("下载完成，总共 %d 张图片，成功 %d 张\n", len(imgURLs), successImages)
	return nil
}

// ParsePage 解析EH页面获取真实图片URL
func ParsePage(client *http.Client, link string) (string, error) {
	realURL, err := GetRealURL(client, link)
	if err != nil {
		return "", fmt.Errorf("获取真实URL失败: %w", err)
	}

	realPage, err := ParseRealPage(client, realURL)
	if err != nil {
		return "", fmt.Errorf("解析真实页面失败: %w", err)
	}

	return realPage, nil
}

// GetRealURL 获取真实图片URL
func GetRealURL(client *http.Client, link string) (string, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}

	// 设置Cookie
	req.AddCookie(&http.Cookie{
		Name:  "nw",
		Value: "1",
	})

	// 设置User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// 获取img标签的onerror属性
	imgOnError := ""
	doc.Find("#img").Each(func(i int, s *goquery.Selection) {
		if onError, exists := s.Attr("onerror"); exists {
			imgOnError = onError
		}
	})

	if imgOnError == "" {
		return "", fmt.Errorf("找不到图片onerror属性")
	}

	// 提取nl参数
	re := regexp.MustCompile(`nl\('(.+)'\)`)
	matches := re.FindStringSubmatch(imgOnError)
	if len(matches) < 2 {
		return "", fmt.Errorf("无法解析nl参数")
	}

	nl := matches[1]
	realURL := fmt.Sprintf("%s?nl=%s", link, nl)
	return realURL, nil
}

// ParseRealPage 解析真实页面获取图片URL
func ParseRealPage(client *http.Client, realURL string) (string, error) {
	resp, err := client.Get(realURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// 获取真实图片URL
	imgURL := ""
	doc.Find("#img").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			imgURL = src
		}
	})

	if imgURL == "" {
		return "", fmt.Errorf("找不到图片URL")
	}

	return imgURL, nil
}

// ParseLinks 解析页面中的图片链接
func ParseLinks(body string) []string {
	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil
	}

	var links []string
	doc.Find("#gdt > a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			links = append(links, href)
		}
	})

	return links
}

// GetAlbum 获取整个专辑信息
func GetAlbum(client *http.Client, url string) (*EHentaiAlbum, error) {
	var pages []string
	currentURL := url

	albumName := ""

	for {
		req, err := http.NewRequest("GET", currentURL, nil)
		if err != nil {
			return nil, err
		}

		// 设置Cookie
		req.AddCookie(&http.Cookie{
			Name:  "nw",
			Value: "1",
		})

		// 设置User-Agent
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
		}

		// 读取响应
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}

		// 获取页面内容
		html, err := doc.Html()
		if err != nil {
			resp.Body.Close()
			return nil, err
		}

		pages = append(pages, html)

		// 获取专辑名称
		if albumName == "" {
			doc.Find("#gn").Each(func(i int, s *goquery.Selection) {
				albumName = s.Text()
			})
		}

		// 查找下一页
		nextPageHref := ""
		doc.Find("body > .gtb td:last-child > a").Each(func(i int, s *goquery.Selection) {
			if href, exists := s.Attr("href"); exists {
				nextPageHref = href
			}
		})

		resp.Body.Close()

		if nextPageHref != "" {
			currentURL = nextPageHref
		} else {
			break
		}

		// 延迟一下，防止被ban
		time.Sleep(1 * time.Second)
	}

	if albumName == "" {
		return nil, fmt.Errorf("无法获取专辑名称")
	}

	return &EHentaiAlbum{
		Name:  albumName,
		Pages: pages,
	}, nil
}
