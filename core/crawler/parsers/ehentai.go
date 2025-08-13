package parsers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"

	"ImageMaster/core/request"
	"ImageMaster/core/types"
)

// EHentaiAlbum EH专辑
type EHentaiAlbum struct {
	Name  string
	Pages []string
}

// EHentaiParser EHentai解析器实现
type EHentaiParser struct{}

// GetName 获取解析器名称
func (p *EHentaiParser) GetName() string {
	return "eHentai"
}

// Parse 解析URL获取图片信息
func (p *EHentaiParser) Parse(reqClient *request.Client, url string) (*ParseResult, error) {
	// 设置ehentai特殊配置
	err := SetupEHentaiClient(reqClient, nil)
	if err != nil {
		return nil, fmt.Errorf("设置EHentai客户端失败: %w", err)
	}

	eHentaiAlbum, err := GetAlbumWithClient(reqClient, url)
	if err != nil {
		return nil, fmt.Errorf("获取专辑失败: %w", err)
	}

	// 批量下载URL和路径
	var imgURLs []string
	var filePaths []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 遍历每一页
	for pageIndex, page := range eHentaiAlbum.Pages {
		links := ParseLinks(page)

		// 并发处理每个链接
		for linkIndex, link := range links {
			wg.Add(1)
			go func(pageIdx, linkIdx int, linkURL string) {
				defer wg.Done()

				// 解析页面获取真实图片URL
				imgURL, err := ParsePageWithClient(reqClient, linkURL)
				if err != nil {
					fmt.Printf("解析页面失败 %s: %v\n", linkURL, err)
					return
				}
				fmt.Printf("解析到图片：%s\n", imgURL)

				// 构建保存文件名
				filename := fmt.Sprintf("%d_%d.jpg", pageIdx, linkIdx)

				// 线程安全地添加到结果中
				mu.Lock()
				imgURLs = append(imgURLs, imgURL)
				filePaths = append(filePaths, filename)
				mu.Unlock()
			}(pageIndex, linkIndex, link)
		}
	}

	// 等待所有并发任务完成
	wg.Wait()

	return &ParseResult{
		Name:      eHentaiAlbum.Name,
		ImageURLs: imgURLs,
		FilePaths: filePaths,
	}, nil
}

// ParsePageWithClient 解析EH页面获取真实图片URL，使用request客户端
func ParsePageWithClient(reqClient *request.Client, link string) (string, error) {
	realURL, err := GetRealURLWithClient(reqClient, link)
	if err != nil {
		return "", fmt.Errorf("获取真实URL失败: %w", err)
	}

	realPage, err := ParseRealPageWithClient(reqClient, realURL)
	if err != nil {
		return "", fmt.Errorf("解析真实页面失败: %w", err)
	}

	return realPage, nil
}

// GetRealURLWithClient 获取真实图片URL，使用request客户端
func GetRealURLWithClient(reqClient *request.Client, link string) (string, error) {
	resp, err := reqClient.RateLimitedGet(link)
	fmt.Printf("获取真实URL成功...: %s\n", link)

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

// ParseRealPageWithClient 解析真实页面获取图片URL，使用request客户端
func ParseRealPageWithClient(reqClient *request.Client, realURL string) (string, error) {
	resp, err := reqClient.RateLimitedGet(realURL)
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

// GetAlbumWithClient 获取整个专辑信息，使用request客户端
func GetAlbumWithClient(reqClient *request.Client, url string) (*EHentaiAlbum, error) {
	// 首先访问第一页获取专辑名称和分页信息
	resp, err := reqClient.RateLimitedGet(url)
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
	resp.Body.Close()

	// 获取专辑名称
	albumName := ""
	doc.Find("#gn").Each(func(i int, s *goquery.Selection) {
		albumName = s.Text()
	})

	if albumName == "" {
		return nil, fmt.Errorf("无法获取专辑名称")
	}

	// 获取所有页面URL
	pageURLs := []string{url} // 包含当前页面

	// 获取第一个 .gtb 元素中的所有 td，排除第一个和最后一个
	gtbElement := doc.Find("body > .gtb").First()
	if gtbElement.Length() > 0 {
		tds := gtbElement.Find("td")
		totalTds := tds.Length()

		tds.Each(func(i int, s *goquery.Selection) {
			if i == 0 || i == 1 || i == totalTds-1 {
				return
			}

			// 从td中的a标签获取href
			s.Find("a").Each(func(j int, a *goquery.Selection) {
				if href, exists := a.Attr("href"); exists {
					pageURLs = append(pageURLs, href)
				}
			})
		})
	}

	// 如果只有一页，直接返回当前页面内容
	if len(pageURLs) == 1 {
		html, err := doc.Html()
		if err != nil {
			return nil, err
		}
		return &EHentaiAlbum{
			Name:  albumName,
			Pages: []string{html},
		}, nil
	}

	// 并发访问所有页面
	var pages []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 预分配pages切片
	pages = make([]string, len(pageURLs))

	for index, pageURL := range pageURLs {
		fmt.Printf("访问第%d页, %s\n", index+1, pageURL)
		wg.Add(1)
		go func(idx int, pURL string) {
			defer wg.Done()

			pageResp, err := reqClient.RateLimitedGet(pURL)
			if err != nil {
				fmt.Printf("访问页面失败 %s: %v\n", pURL, err)
				return
			}
			defer pageResp.Body.Close()

			if pageResp.StatusCode != http.StatusOK {
				fmt.Printf("页面HTTP状态码错误 %s: %d\n", pURL, pageResp.StatusCode)
				return
			}

			pageDoc, err := goquery.NewDocumentFromReader(pageResp.Body)
			if err != nil {
				fmt.Printf("解析页面失败 %s: %v\n", pURL, err)
				return
			}

			html, err := pageDoc.Html()
			if err != nil {
				fmt.Printf("获取页面HTML失败 %s: %v\n", pURL, err)
				return
			}

			mu.Lock()
			pages[idx] = html
			mu.Unlock()
		}(index, pageURL)
	}

	// 等待所有页面访问完成
	wg.Wait()

	// 过滤掉空的页面
	var validPages []string
	for _, page := range pages {
		if page != "" {
			validPages = append(validPages, page)
		}
	}

	return &EHentaiAlbum{
		Name:  albumName,
		Pages: validPages,
	}, nil
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

// EHentaiCrawler E-Hentai爬虫
type EHentaiCrawler struct {
	*BaseCrawler
}

// NewEHentaiCrawler 创建新的E-Hentai爬虫
func NewEHentaiCrawler(reqClient *request.Client) types.ImageCrawler {
	parser := &EHentaiParser{}
	baseCrawler := NewBaseCrawler(reqClient, parser)
	return &EHentaiCrawler{
		BaseCrawler: baseCrawler,
	}
}

// SetupEHentaiClient 设置EHentai特殊的客户端配置
func SetupEHentaiClient(reqClient *request.Client, downloader types.Downloader) error {
	// 先执行通用设置
	if err := SetupRequestClient(reqClient, downloader); err != nil {
		return err
	}

	// 设置ehentai需要的cookie
	reqClient.AddCookie(&http.Cookie{
		Name:  "nw",
		Value: "1",
	})

	return nil
}
