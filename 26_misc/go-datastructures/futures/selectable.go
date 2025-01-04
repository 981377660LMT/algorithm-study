// **多播一次性事件**：
// 执行异步操作，填充一个结果，后续多个 goroutine 同时想要该结果时，
// !通过 `WaitChan()` 可以与其他事件/超时并行等待，通过 `GetResult()` 可以阻塞等待或立即获取已完成结果。

package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// 1. 创建一个 Selectable future
	f := NewSelectable()

	// !2. 启动一个异步任务，模拟1秒后得到结果
	go func() {
		time.Sleep(1 * time.Second)
		// f.SetValue("Hello Future") or f.SetError(err)
		f.SetValue("Hello Future")
	}()

	// !3. 用 WaitChan() + select 方式等待结果
	go func() {
		select {
		case <-f.WaitChan():
			val, err := f.GetResult()
			fmt.Println("[Goroutine1] got result:", val, "error:", err)
		case <-time.After(2 * time.Second):
			fmt.Println("[Goroutine1] Timeout!")
		}
	}()

	// !4. 另一goroutine直接调用 GetResult()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := f.GetResult()
		fmt.Println("[Goroutine2] got result:", val, "error:", err)
	}()

	wg.Wait()
	fmt.Println("Done.")
}

// ErrFutureCanceled signals that futures in canceled by a call to `f.Cancel()`
var ErrFutureCanceled = errors.New("future canceled")

var closed = make(chan struct{})

func init() {
	close(closed)
}

// Selectable is a future with channel exposed for external `select`.
// Many simultaneous listeners may wait for result either with `f.Value()`
// or by selecting/fetching from `f.WaitChan()`, which is closed when future
// fulfilled.
// Selectable contains sync.Mutex, so it is not movable/copyable.
type Selectable struct {
	m    sync.Mutex
	val  any
	err  error
	wait chan struct{}

	filled uint32
}

// NewSelectable returns new selectable future.
// Note: this method is for backward compatibility.
// You may allocate it directly on stack or embedding into larger structure
func NewSelectable() *Selectable {
	return &Selectable{}
}

// WaitChan returns channel, which is closed when future is fullfilled.
func (f *Selectable) WaitChan() <-chan struct{} {
	if atomic.LoadUint32(&f.filled) == 1 {
		return closed
	}
	return f.wchan()
}

// GetResult waits for future to be fullfilled and returns value or error,
// whatever is set first
func (f *Selectable) GetResult() (any, error) {
	if atomic.LoadUint32(&f.filled) == 0 {
		<-f.wchan() // 等待 channel 关闭
	}
	return f.val, f.err
}

// SetValue is alias for Fill(v, nil)
func (f *Selectable) SetValue(v any) error {
	return f.fill(v, nil)
}

// SetError is alias for Fill(nil, e)
func (f *Selectable) SetError(e error) {
	f.fill(nil, e)
}

// Cancel is alias for SetError(ErrFutureCanceled)
func (f *Selectable) Cancel() {
	f.SetError(ErrFutureCanceled)
}

func (f *Selectable) wchan() <-chan struct{} {
	f.m.Lock()
	if f.wait == nil {
		f.wait = make(chan struct{})
	}
	ch := f.wait
	f.m.Unlock()
	return ch
}

// fill sets value for future, if it were not already fullfilled
// Returns error, if it were already set to future.
func (f *Selectable) fill(v any, e error) error {
	f.m.Lock()
	if f.filled == 0 {
		f.val = v
		f.err = e
		atomic.StoreUint32(&f.filled, 1)
		w := f.wait
		f.wait = closed
		if w != nil {
			close(w)
		}
	}
	f.m.Unlock()
	return f.err
}
