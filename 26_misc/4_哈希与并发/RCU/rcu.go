// 轻量级的“RCU风格”数据刷新器：
// 读者总能读到一个最新可用的版本，
// 写者通过一个回调在后台周期性地生成新的数据并原子地切换指针，
// 达到了高并发读、低锁开销的目标。

package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type RCU[T any] struct {
	value  *atomic.Pointer[T] // 原子的指向当前数据 T 的指针
	update func(*T) *T        // 更新函数, 用于生成新的 T
	ticker *time.Ticker       // 定时器, 周期性触发 update
	done   chan struct{}      // 用于关闭 goroutine 的信号
	mu     *sync.Mutex        // 在执行更新时用到的互斥锁
}

func NewRCU[T any](updateInterval time.Duration, update func(*T) *T) *RCU[T] {
	self := &RCU[T]{
		value:  &atomic.Pointer[T]{},
		update: update,
		ticker: time.NewTicker(updateInterval),
		done:   make(chan struct{}),
		mu:     &sync.Mutex{},
	}
	self.forceUpdate()
	go self.scheduleUpdate()
	return self
}

func (r *RCU[T]) Load() *T {
	return r.value.Load()
}

func (r *RCU[T]) scheduleUpdate() {
	for {
		select {
		case <-r.done:
			return
		case <-r.ticker.C:
			r.forceUpdate()
		}
	}
}

func (r *RCU[T]) forceUpdate() {
	r.mu.Lock()
	defer r.mu.Unlock()

	newValue := r.update(r.value.Load())
	r.value.Store(newValue)
}

func (r *RCU[T]) Close() {
	r.ticker.Stop()
	r.done <- struct{}{}
}
