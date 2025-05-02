package request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"ImageMaster/core/config"
)

// Client HTTP客户端封装
type Client struct {
	client        *http.Client
	proxyURL      string
	configManager *config.Manager
	headers       map[string]string
	cookies       []*http.Cookie
}

// NewClient 创建新的请求客户端
func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
		cookies: make([]*http.Cookie, 0),
	}
}

// SetConfigManager 设置配置管理器
func (c *Client) SetConfigManager(configManager *config.Manager) {
	c.configManager = configManager

	// 如果配置管理器不为空，应用其代理设置
	if configManager != nil {
		proxyURL := configManager.GetProxy()
		if proxyURL != "" {
			c.SetProxy(proxyURL)
		}
	}
}

// SetProxy 设置代理
func (c *Client) SetProxy(proxyURL string) error {
	c.proxyURL = proxyURL

	// 如果代理URL为空，使用默认客户端
	if proxyURL == "" {
		c.client = &http.Client{
			Timeout: 30 * time.Second,
		}
		return nil
	}

	// 设置代理
	proxyURLParsed, err := url.Parse(proxyURL)
	if err != nil {
		return fmt.Errorf("解析代理URL失败: %w", err)
	}

	// 创建带代理的Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURLParsed),
	}

	// 更新客户端
	c.client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return nil
}

// GetProxy 获取当前代理设置
func (c *Client) GetProxy() string {
	return c.proxyURL
}

// SetHeader 设置请求头
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// SetHeaders 批量设置请求头
func (c *Client) SetHeaders(headers map[string]string) {
	for key, value := range headers {
		c.headers[key] = value
	}
}

// AddCookie 添加Cookie
func (c *Client) AddCookie(cookie *http.Cookie) {
	c.cookies = append(c.cookies, cookie)
}

// ClearCookies 清除所有Cookie
func (c *Client) ClearCookies() {
	c.cookies = make([]*http.Cookie, 0)
}

// Get 发送GET请求
func (c *Client) Get(url string) (*http.Response, error) {
	return c.DoRequest("GET", url, nil, nil)
}

// GetWithContext 发送带上下文的GET请求
func (c *Client) GetWithContext(ctx context.Context, url string) (*http.Response, error) {
	return c.DoRequestWithContext(ctx, "GET", url, nil, nil)
}

// Post 发送POST请求
func (c *Client) Post(url string, body io.Reader, contentType string) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": contentType,
	}
	return c.DoRequest("POST", url, body, headers)
}

// PostWithContext 发送带上下文的POST请求
func (c *Client) PostWithContext(ctx context.Context, url string, body io.Reader, contentType string) (*http.Response, error) {
	headers := map[string]string{
		"Content-Type": contentType,
	}
	return c.DoRequestWithContext(ctx, "POST", url, body, headers)
}

// DoRequest 执行HTTP请求
func (c *Client) DoRequest(method, url string, body io.Reader, extraHeaders map[string]string) (*http.Response, error) {
	// 尝试从配置中应用代理（如果尚未设置代理且配置管理器存在）
	if c.proxyURL == "" && c.configManager != nil {
		proxyURL := c.configManager.GetProxy()
		if proxyURL != "" {
			c.SetProxy(proxyURL)
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置默认User-Agent
	if _, ok := c.headers["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	}

	// 应用客户端的通用头部
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// 应用额外的请求头
	for key, value := range extraHeaders {
		req.Header.Set(key, value)
	}

	// 应用Cookie
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// 执行请求
	return c.client.Do(req)
}

// DoRequestWithContext 执行带上下文的HTTP请求
func (c *Client) DoRequestWithContext(ctx context.Context, method, url string, body io.Reader, extraHeaders map[string]string) (*http.Response, error) {
	// 尝试从配置中应用代理（如果尚未设置代理且配置管理器存在）
	if c.proxyURL == "" && c.configManager != nil {
		proxyURL := c.configManager.GetProxy()
		if proxyURL != "" {
			c.SetProxy(proxyURL)
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置默认User-Agent
	if _, ok := c.headers["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	}

	// 应用客户端的通用头部
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// 应用额外的请求头
	for key, value := range extraHeaders {
		req.Header.Set(key, value)
	}

	// 应用Cookie
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	// 执行请求
	return c.client.Do(req)
}

// GetHTTPClient 获取内部的HTTP客户端
// 用于兼容旧代码，逐渐迁移到新版请求模块后应删除此方法
func (c *Client) GetHTTPClient() *http.Client {
	return c.client
}
