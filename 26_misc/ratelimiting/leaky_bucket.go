package main

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity    int       // 漏桶容量（最大可同时排队的请求数）
	leakRate    float64   // 漏出速率（每秒能处理多少个请求，或每毫秒能处理多少等）
	currentSize float64   // 当前排队中还未处理完的请求数量（可用浮点表示）
	lastTime    time.Time // 上次更新的时间
	mu          sync.Mutex
}

// NewLeakyBucket 创建一个 LeakyBucket
//   - capacity: 漏桶最大可排队的请求数量
//   - leakRate: 每秒可处理的请求数，例如 5 表示 1 秒能处理 5 个请求
func NewLeakyBucket(capacity int, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		lastTime: time.Now(),
	}
}

// Allow 尝试向漏桶中加入一个请求
//   - 如果在更新后队列未满，则允许（返回 true）
//   - 如果已满，则丢弃（返回 false）
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastTime).Seconds() // 距离上次更新的秒数，float64

	// 1. 计算在 elapsed 时间里，可漏出多少请求
	//    漏出量 = leakRate * elapsed
	leaked := lb.leakRate * elapsed

	// 2. 更新当前的排队请求数量
	if leaked >= lb.currentSize {
		// 如果可漏出的量 >= 当前队列占用，则队列清空
		lb.currentSize = 0
	} else {
		// 否则，减少 leaked 数量
		lb.currentSize -= leaked
	}

	// 3. 更新 lastTime
	lb.lastTime = now

	// 4. 尝试加入新请求
	if lb.currentSize < float64(lb.capacity) {
		lb.currentSize += 1
		return true
	}
	// 队列已满，丢弃请求
	return false
}

// --------------------- 测试示例 ---------------------
func main() {
	// capacity=5，leakRate=2 表示：队列最大可排队 5 个请求，每秒可处理 2 个请求
	lb := NewLeakyBucket(5, 2.0)

	// 总共发 10 个请求试试
	for i := 1; i <= 10; i++ {
		allowed := lb.Allow()
		if allowed {
			fmt.Printf("Request #%d allowed\n", i)
		} else {
			fmt.Printf("Request #%d dropped\n", i)
		}
		time.Sleep(200 * time.Millisecond) // 0.2 秒发一个请求
	}

	// 等待 3 秒，让漏桶继续“漏”一些请求
	time.Sleep(3 * time.Second)
	fmt.Println("After waiting 3 seconds...")

	// 再次发送一批请求
	for i := 11; i <= 15; i++ {
		allowed := lb.Allow()
		if allowed {
			fmt.Printf("Request #%d allowed\n", i)
		} else {
			fmt.Printf("Request #%d dropped\n", i)
		}
		time.Sleep(300 * time.Millisecond) // 0.3 秒发一个请求
	}
}
