package main

import (
	"fmt"
	"sync"
	"time"
)

// SlidingWindowRateLimiter 滑动窗口限流器
type SlidingWindowRateLimiter struct {
	windowSize  time.Duration // 时间窗口大小，例如 1 秒
	maxRequests int           // 窗口内允许的最大请求数

	mu         sync.Mutex
	timestamps []time.Time // 存放最近的请求时间戳
}

// NewSlidingWindowRateLimiter 创建并初始化一个滑动窗口限流器
//   - windowSize: 窗口大小，例如 1s、5s
//   - maxRequests: 在一个 windowSize 内允许的最大请求数
func NewSlidingWindowRateLimiter(windowSize time.Duration, maxRequests int) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		windowSize:  windowSize,
		maxRequests: maxRequests,
	}
}

// Allow 尝试获取一次请求配额，成功返回 true，超过限流返回 false
func (sw *SlidingWindowRateLimiter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	// 1. 移除窗口区间外的旧请求记录
	cutoff := now.Add(-sw.windowSize)
	ptr := 0
	for ptr < len(sw.timestamps) && sw.timestamps[ptr].Before(cutoff) {
		ptr++
	}
	// idx 之前的时间戳都在窗口外，直接截断
	sw.timestamps = sw.timestamps[ptr:]

	// 2. 判断当前窗口内的请求数是否已达上限
	if len(sw.timestamps) >= sw.maxRequests {
		return false
	}

	// 3. 记录当前请求
	sw.timestamps = append(sw.timestamps, now)
	return true
}

func main() {
	// 创建一个滑动窗口限流器：窗口大小 1 秒，最多允许 5 个请求
	limiter := NewSlidingWindowRateLimiter(1*time.Second, 5)

	// 模拟 10 次请求
	for i := 1; i <= 10; i++ {
		allowed := limiter.Allow()
		if allowed {
			fmt.Printf("Request #%d: allowed\n", i)
		} else {
			fmt.Printf("Request #%d: blocked\n", i)
		}

		// 为了测试效果，这里每 100ms 来一个请求
		time.Sleep(100 * time.Millisecond)
	}

	// 等待 1 秒再发送一批请求，观察结果
	time.Sleep(1 * time.Second)
	fmt.Println("After waiting 1 seconds...")

	// 再发送一批请求
	for i := 11; i <= 20; i++ {
		allowed := limiter.Allow()
		if allowed {
			fmt.Printf("Request #%d: allowed\n", i)
		} else {
			fmt.Printf("Request #%d: blocked\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
