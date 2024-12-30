package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int64         // 令牌桶容量
	tokens       int64         // 当前可用令牌数
	fillInterval time.Duration // 生成一个令牌的时间间隔
	lastFill     time.Time     // 上次填充令牌的时间
	mu           sync.Mutex    // 互斥锁，保证并发安全
}

// NewTokenBucket 创建并初始化一个令牌桶
//   - capacity: 令牌桶容量
//   - fillInterval: 每隔多长时间生成一个令牌
func NewTokenBucket(capacity int64, fillInterval time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		fillInterval: fillInterval,
		lastFill:     time.Now(),
	}
}

// 请求获取一个令牌，获取成功返回 true，否则返回 false.
func (tb *TokenBucket) Acquire() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 计算距离上次填充已经过去了多长时间
	elapsed := time.Now().Sub(tb.lastFill)

	// 根据 elapsed 计算需要填充的令牌数量
	if elapsed > 0 {
		newTokens := elapsed / tb.fillInterval
		if newTokens > 0 {
			// 计算填充后的令牌数量，不能超过容量
			tb.tokens = min64(tb.capacity, tb.tokens+int64(newTokens))
			// 更新 lastFill, 只往后挪 newTokens * fillInterval
			tb.lastFill = tb.lastFill.Add(tb.fillInterval * newTokens)
		}
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	// 创建一个令牌桶，容量为 5，每 500ms 生成 1 个令牌
	tb := NewTokenBucket(5, 500*time.Millisecond)

	// 模拟请求
	for i := 1; i <= 10; i++ {
		if tb.Acquire() {
			fmt.Printf("Request #%d allowed\n", i)
		} else {
			fmt.Printf("Request #%d blocked\n", i)
		}
		// 睡眠一会儿模拟请求间隔
		time.Sleep(200 * time.Millisecond)
	}

	// 可以尝试再过一段时间后（让令牌桶重新填充）继续请求
	time.Sleep(2 * time.Second)
	fmt.Println("After waiting 2 seconds...")

	for i := 11; i <= 15; i++ {
		if tb.Acquire() {
			fmt.Printf("Request #%d allowed\n", i)
		} else {
			fmt.Printf("Request #%d blocked\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
