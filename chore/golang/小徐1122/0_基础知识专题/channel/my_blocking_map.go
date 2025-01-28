// 优化MyBlockingMap的性能，可以从以下几个方面入手：
//
// 减少锁竞争：通过更细粒度的锁或读写锁优化并发性能。
// 减少通道创建开销：复用通道或使用更轻量的通知机制。
// 避免惊群效应：优化唤醒机制，避免大量goroutine同时被唤醒。
// 减少内存分配：通过对象池或其他方式复用数据结构。

package main

import (
	"errors"
	"sync"
	"time"
)

func main() {
	bm := NewMyBlockingMap()

	// 示例1：基础使用
	go func() {
		time.Sleep(100 * time.Millisecond)
		bm.Put(1, 42)
	}()

	if v, err := bm.Get(1, 200*time.Millisecond); err == nil {
		println("Get 1:", v) // 应输出42
	}

	// 示例2：超时测试
	if _, err := bm.Get(2, 50*time.Millisecond); err != nil {
		println("Get 2:", err.Error()) // 应输出timeout
	}

	// 示例3：并发测试
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			if v, err := bm.Get(3, time.Duration(i)*time.Second); err != nil {
				println("Concurrent Get 3:", err.Error())
			} else {
				println("Concurrent Get 3:", v)
			}
		}()
	}

	time.Sleep(2 * time.Second)
	bm.Put(3, 100)
	wg.Wait()
}

// MyBlockingMap 支持阻塞获取的并发安全Map
type MyBlockingMap struct {
	mu    sync.Mutex              // 互斥锁，保护data和waits
	data  map[int]int             // 存储键值对
	waits map[int][]chan struct{} // 每个key对应的等待通道列表
}

// NewMyBlockingMap 创建新的MyBlockingMap实例
func NewMyBlockingMap() *MyBlockingMap {
	return &MyBlockingMap{
		data:  make(map[int]int),
		waits: make(map[int][]chan struct{}),
	}
}

// Put 存储键值对，并唤醒所有等待该key的goroutine
func (mp *MyBlockingMap) Put(k, v int) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.data[k] = v // 存储值

	// 唤醒所有等待该key的goroutine
	if chs, ok := mp.waits[k]; ok {
		for _, ch := range chs {
			close(ch) // 关闭通道以通知所有等待者
		}
		delete(mp.waits, k) // 清空等待列表
	}
}

// Get 获取值。若key不存在则阻塞，直到超时或值被设置(防止协程泄漏)
func (mp *MyBlockingMap) Get(k int, maxWaitDuration time.Duration) (int, error) {
	mp.mu.Lock()

	// 快速路径：值已存在
	if v, ok := mp.data[k]; ok {
		mp.mu.Unlock()
		return v, nil
	}

	// 创建等待通道并注册
	ch := make(chan struct{})
	mp.waits[k] = append(mp.waits[k], ch)
	mp.mu.Unlock()

	// 设置超时计时器（使用defer确保停止）
	timer := time.NewTimer(maxWaitDuration)
	defer timer.Stop()

	select {
	case <-ch: // 被Put唤醒
		mp.mu.Lock()
		v, ok := mp.data[k] // 必须重新检查，可能被其他goroutine抢先获取
		mp.mu.Unlock()

		if !ok {
			return 0, errors.New("value disappeared")
		}
		return v, nil

	case <-timer.C: // 超时
		mp.mu.Lock()
		defer mp.mu.Unlock()

		// 双重检查：可能在超时前最后一刻被设置
		if v, ok := mp.data[k]; ok {
			// 清理残留的等待通道（可能已被Put处理）
			if waits, exists := mp.waits[k]; exists {
				// 遍历查找并移除当前通道
				for i, c := range waits {
					if c == ch {
						mp.waits[k] = append(waits[:i], waits[i+1:]...)
						break
					}
				}
				if len(mp.waits[k]) == 0 {
					delete(mp.waits, k)
				}
			}
			return v, nil
		}

		// 正式处理超时：从等待列表移除当前通道
		if waits, exists := mp.waits[k]; exists {
			for i, c := range waits {
				if c == ch {
					mp.waits[k] = append(waits[:i], waits[i+1:]...)
					break
				}
			}
			if len(mp.waits[k]) == 0 {
				delete(mp.waits, k)
			}
		}
		return 0, errors.New("timeout")
	}
}
