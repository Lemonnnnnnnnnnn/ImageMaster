package request

import (
	"net/http"
	"sync"
	"time"
)

// TokenBucket 令牌桶算法实现请求频率控制
type TokenBucket struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration // 令牌补充间隔
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewTokenBucket 创建令牌桶
func NewTokenBucket(maxTokens int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// TakeToken 获取令牌，如果没有令牌则等待
func (tb *TokenBucket) TakeToken() {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	for {
		// 补充令牌
		now := time.Now()
		elapsed := now.Sub(tb.lastRefill)
		tokensToAdd := int(elapsed / tb.refillRate)

		if tokensToAdd > 0 {
			tb.tokens += tokensToAdd
			if tb.tokens > tb.maxTokens {
				tb.tokens = tb.maxTokens
			}
			tb.lastRefill = now
		}

		// 如果有令牌，消费一个令牌并返回
		if tb.tokens > 0 {
			tb.tokens--
			return
		}

		// 如果没有令牌，等待下一个令牌补充时间
		waitTime := tb.refillRate
		tb.mutex.Unlock()
		time.Sleep(waitTime)
		tb.mutex.Lock()
		// 循环继续，重新检查令牌状态
	}
}

// GetTokenCount 获取当前令牌数量（仅用于调试）
func (tb *TokenBucket) GetTokenCount() int {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	return tb.tokens
}

// 默认令牌桶配置
const (
	DefaultMaxTokens  = 5                      // 默认最大令牌数
	DefaultRefillRate = 500 * time.Millisecond // 默认令牌补充间隔
)

// 全局默认令牌桶实例
var DefaultTokenBucket = NewTokenBucket(DefaultMaxTokens, DefaultRefillRate)

// RateLimitedGet 使用默认令牌桶的频率限制GET请求
func (c *Client) RateLimitedGet(url string) (*http.Response, error) {
	// 获取令牌（会自动等待如果需要）
	DefaultTokenBucket.TakeToken()

	// 执行请求
	return c.Get(url)
}

// RateLimitedGetWithBucket 使用指定令牌桶的频率限制GET请求
func (c *Client) RateLimitedGetWithBucket(url string, bucket *TokenBucket) (*http.Response, error) {
	// 获取令牌（会自动等待如果需要）
	bucket.TakeToken()

	// 执行请求
	return c.Get(url)
}
