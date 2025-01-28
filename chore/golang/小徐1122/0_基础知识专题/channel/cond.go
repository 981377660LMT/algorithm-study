// 核心流程：
// 1. 使用 `sync.NewCond` 创建条件变量，绑定一个互斥锁。
// 2. 需要等待的 Goroutine 执行 `cond.Wait()` 前需先 `Lock()`，等待时会临时解锁。
// 3. 另一个 Goroutine 使用 `cond.Signal()` 或 `cond.Broadcast()` 发送通知，唤醒一个或全部等待的 Goroutine。
// 4. 被唤醒后，Goroutine 会重新获取锁并继续执行。

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	ready := false

	// 等待方
	go func() {
		cond.L.Lock()
		for !ready {
			cond.Wait() // 等待状态，直到被 Signal 或 Broadcast
		}
		fmt.Println("Goroutine 1 收到通知，继续执行")
		cond.L.Unlock()
	}()

	// 通知方
	go func() {
		time.Sleep(time.Second)
		cond.L.Lock()
		ready = true
		cond.Signal() // 唤醒一个等待 Goroutine
		cond.L.Unlock()
	}()

	time.Sleep(2 * time.Second)
}
