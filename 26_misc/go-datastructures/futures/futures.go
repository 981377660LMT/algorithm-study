/*
Package futures is useful for broadcasting an identical message to a multitude
of listeners as opposed to channels which will choose a listener at random
if multiple listeners are listening to the same channel.  The future will
also cache the result so any future interest will be immediately returned
to the consumer.
*/

// **futures 包的目标**是：
// 只需将结果写入一次，多个消费者都能得到该结果（或超时错误），
// 像一种“只写一次，多读多播”的模式
// spmc (single producer, multiple consumers)

package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	// 1. 建一个只写一次的 channel, 供 Future 监听
	completer := make(chan any)

	// 2. 创建 Future, 超时为 5 秒
	future := NewFuture(completer, 5*time.Second)

	// 3. 异步发起网络请求，成功后通过 completer 发送结果
	go func() {
		resp, err := http.Get("https://www.example.com")
		if err != nil {
			// 如果出错, 可以直接塞 err 到 channel
			// 但由于 channel 是 <-chan any, 这里要转换一下
			completer <- err
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			completer <- err
			return
		}

		// 把结果发给 Future
		completer <- string(body)
	}()

	// 4. 任何想要获取结果的人，都可以调用 future.GetResult()
	//    这里演示两个 goroutine 都来获取结果
	go func() {
		res, err := future.GetResult()
		fmt.Println("Goroutine 1 get:", res, err)
	}()

	go func() {
		res, err := future.GetResult()
		fmt.Println("Goroutine 2 get:", res, err)
	}()

	// main 等一会儿，防止退出
	time.Sleep(2 * time.Second)
}

// Completer is a channel that the future expects to receive
// a result on.  The future only receives on this channel.
type Completer <-chan any

// Future represents an object that can be used to perform asynchronous
// tasks.  The constructor of the future will complete it, and listeners
// will block on getresult until a result is received.  This is different
// from a channel in that the future is only completed once, and anyone
// listening on the future will get the result, regardless of the number
// of listeners.
type Future struct {
	triggered bool       // 标记是否已经填充结果
	item      any        // 填充的结果
	err       error      // 可能的错误（如超时或其它异常）
	lock      sync.Mutex // 保护上述字段的并发安全
	wg        sync.WaitGroup
}

// GetResult will immediately fetch the result if it exists
// or wait on the result until it is ready.
func (f *Future) GetResult() (any, error) {
	f.lock.Lock()
	if f.triggered {
		f.lock.Unlock()
		return f.item, f.err
	}
	f.lock.Unlock()

	f.wg.Wait() // 等待结果就绪
	return f.item, f.err
}

// HasResult will return true iff the result exists
func (f *Future) HasResult() bool {
	f.lock.Lock()
	hasResult := f.triggered
	f.lock.Unlock()
	return hasResult
}

func (f *Future) setItem(item any, err error) {
	f.lock.Lock()
	f.triggered = true
	f.item = item
	f.err = err
	f.lock.Unlock()
	f.wg.Done()
}

func listenForResult(f *Future, ch Completer, timeout time.Duration, wg *sync.WaitGroup) {
	wg.Done() // 通知 New(...) 里的协程：该监听协程已启动
	t := time.NewTimer(timeout)
	select {
	case item := <-ch:
		f.setItem(item, nil)
		t.Stop() // we want to trigger GC of this timer as soon as it's no longer needed
	case <-t.C:
		f.setItem(nil, fmt.Errorf(`timeout after %f seconds`, timeout.Seconds()))
	}
}

// NewFuture is the constructor to generate a new future.  Pass the completed
// item to the toComplete channel and any listeners will get
// notified.  If timeout is hit before toComplete is called,
// any listeners will get passed an error.
func NewFuture(completer Completer, timeout time.Duration) *Future {
	f := &Future{}
	f.wg.Add(1)
	var wg sync.WaitGroup
	wg.Add(1)
	go listenForResult(f, completer, timeout, &wg)
	wg.Wait()
	return f
}
