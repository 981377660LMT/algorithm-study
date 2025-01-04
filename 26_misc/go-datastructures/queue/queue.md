下面给出对这个队列（`Queue`）实现的详细分析和简要用法示例。本队列是一个 **可无限扩容** 的普通先进先出（FIFO）队列（并非环形/无锁），通过 **锁 + 信号量**（或更准确地说一个“waiters”队列）来在“空”等情况下阻塞，或在“dispose”时立即唤醒所有阻塞操作。

## 一、功能概述

1. **多生产者、多消费者**：通过一个 `sync.Mutex` 来保护队列的 `items` 切片，实现并发安全。
2. **阻塞获取**：`Get` / `Poll` 在队列为空时，会阻塞等待 `Put` 插入新元素。
3. **无上限**：相对于 Go 原生 `chan` 可以设置“缓冲区上限”，本队列可以无限扩张 `items` 切片存储。
4. **支持超时**：`Poll(number, timeout)` 会在 `timeout` 后返回 `ErrTimeout`（若未获取到数据）。
5. **Dispose 机制**：一旦调用 `Dispose`，队列被标记为 `disposed`，后续所有 `Put/Get/Poll` 调用都会返回 `ErrDisposed`；并且阻塞在 `Get/Poll` 的 goroutine 会被立即唤醒并返回错误。

## 二、核心数据结构

### 1. `Queue`

```go
type Queue struct {
    waiters  waiters    // 等待获取的 goroutine 列表
    items    items      // 存储队列元素的切片
    lock     sync.Mutex // 互斥锁保护
    disposed bool       // 是否已 Dispose
}
```

- **waiters**：一组 `sema`，表示当前阻塞等待数据的请求。
- **items**：实际存储的数据，类型是 `[]interface{}`。
- **lock**：统一用 `sync.Mutex` 在 `Put` / `Get` / `Poll` / `Dispose` / ... 等操作时保护 `Queue` 状态。
- **disposed**：标记队列是否已被释放，任何阻塞或后续操作都要立即返回错误。

### 2. `waiters` (type `waiters []*sema`)

这是一个辅助切片，用来存放等待“有数据可读”的 goroutine 的信号对象 `sema`。

- 每当队列为空时，`Get/Poll` 会创建一个新的 `sema`，放进 `waiters`。
- 当队列有数据写入时，会从 `waiters` 取出第一个 `sema` 并发信号告诉它可以读取了。
- 如果 `Dispose`，就给所有 `sema` 发信号，让它们立刻返回错误。

### 3. `sema` 信号量

```go
type sema struct {
    ready    chan bool
    response *sync.WaitGroup
}
```

- `ready`：一个容量为 1 的通道，用于发信号“可以读取了”或“队列被 dispose”。
- `response`：`WaitGroup` 用于配合 `ready` 通知完成。主要在 `Poll` 超时处理或被 `Put` 时需要对应 `Done()`。

当 `Poll` 超时或 `Dispose`，需要把已注册的 `sema` 从队列中移除或手动“完成”以解除阻塞。

---

## 三、主要方法

### 1. `Put(items ...interface{}) error`

1. `lock.Lock()` -> 检查 `disposed`。若已 dispose, 返回 `ErrDisposed`。
2. 把 `items` 追加到 `q.items`。
3. **唤醒等待队列**：
   - 循环从 `q.waiters.get()` 取出一个 `sema`，如果为空就 break；否则发信号 `sema.ready <- true` 让对应的 `Poll/Get` 可以获取数据。
   - 那边拿到信号后会 `sema.response.Add(1)` 并 `Wait()`，然后 `sema.response.Done()` 以完成一次交互。
   - 如果 `q.items` 被一次获取后又变空，就 break，不再发下一个信号。
4. `lock.Unlock()`。

**要点**：一次 `Put` 会尽量唤醒尽可能多的阻塞读者，但若队列数据瞬间又读完，就不再继续唤醒更多读者，以免空读。

### 2. `Get(number int64) ([]interface{}, error)`

== `Poll(number, 0)`

### 3. `Poll(number int64, timeout time.Duration)`

1. 若 `number < 1`，直接返回空。
2. `lock.Lock()` 并检查 `disposed`；若已 dispose，返回 `ErrDisposed`。
3. 如果 `q.items` 有数据 (`len(q.items) != 0`)，就直接获取最多 `number` 条数据 (`q.items.get(number)`)，`lock.Unlock()` 并返回。
4. 否则 (队列为空)：
   - 创建一个新的 `sema`，放进 `q.waiters`。
   - `lock.Unlock()`。
   - `select { case <-sema.ready: ... case <-time.After(timeout): ... }`
     - 若 `timeout > 0`，等待指定时间；若 `timeout == 0`，则会一直阻塞直到 `Put` 发信号或 `Dispose`。
     - 若收到 `sema.ready`, 说明 `Put` 中给了信号。再次进入互斥区检查 `disposed`；然后实际获取 `q.items.get(number)`。
     - 若超时，尝试把 `sema` 从 `q.waiters` 移除，如果成功则返回 `ErrTimeout`，如果已经被 `Put` 拿走了，就 `sema.response.Done()` 以解锁 `Put`。
5. 返回获取到的数据或错误。

### 4. `Peek() (interface{}, error)`

- 在 `lock` 下安全地看 `q.items` 的头部元素，不取出，不修改队列。如果空则 `ErrEmptyQueue`，若 dispose 则 `ErrDisposed`。

### 5. `TakeUntil(checker func(item interface{}) bool)`

- 仅在 `lock` 下，把队列前面连续满足 `checker(...) == true` 的元素取出，并返回。若队列空或第一个不满足，则直接返回空数组。
- 不阻塞。如果队列空，不会等数据到达。

### 6. `Dispose() []interface{}`

- `lock.Lock()` -> `q.disposed = true`
- 给所有 `waiters` 发送信号 `waiter.ready <- true`，让它们从阻塞中醒来并返回 `ErrDisposed`。
- 保存当前 `q.items` 到 `disposedItems`，然后 `q.items = nil, q.waiters = nil`。
- `return disposedItems` 作为被丢弃的剩余数据。

### 7. 其它

- `Empty()`：判断队列是否为空（`len(q.items) == 0`）。
- `Len()`：返回队列中元素数量。
- `Disposed()`：是否已标记为 disposed。

---

## 四、执行流程简图

![](https://user-images.githubusercontent.com/5938217/229349065-414459ed-2df9-4763-a2a6-60513bc6ebdc.png)

1. **Put**：加锁 -> 入队 + 唤醒阻塞读者 -> 解锁。
2. **Get** / **Poll**：若队列非空，直接取；否则注册 `sema` 并阻塞，直到被 `Put` 或超时或 `Dispose`。

---

## 五、使用示例

```go
package main

import (
    "fmt"
    "log"
    "time"

    "your_module/queue"
)

func main() {
    q := queue.New(10) // 初始容量10，队列可动态扩容

    // Producer
    go func() {
        for i := 0; i < 5; i++ {
            err := q.Put(fmt.Sprintf("Msg%d", i))
            if err != nil {
                log.Println("Put err:", err)
                return
            }
            time.Sleep(100 * time.Millisecond)
        }
        // Dispose after produce 5 messages
        disposed := q.Dispose()
        fmt.Println("Disposed with leftover items:", disposed)
    }()

    // Consumer
    for {
        items, err := q.Poll(1, 2*time.Second)
        if err == queue.ErrDisposed {
            fmt.Println("Queue is disposed, consumer end.")
            break
        } else if err == queue.ErrTimeout {
            fmt.Println("Poll timed out, try again or exit.")
            break
        } else if err != nil {
            log.Println("Poll error:", err)
            break
        } else {
            fmt.Println("Got item:", items[0])
        }
    }
}
```

**运行流程**：

1. Producer 每 100ms 放一个字符串 (`Msg0 ~ Msg4`) 到队列，然后 `Dispose()`。
2. Consumer 用 `Poll(1, 2s)` 获取 1 条数据；若没有，会阻塞直到：
   - 新数据进入（续执行），或
   - 2s 超时返回 `ErrTimeout`，或
   - `Dispose()` 导致 `ErrDisposed` 并退出。

输出可能类似：

```
Got item: Msg0
Got item: Msg1
Got item: Msg2
Got item: Msg3
Got item: Msg4
Disposed with leftover items: []
Queue is disposed, consumer end.
```

---

## 六、注意事项

1. **阻塞模式**：当队列为空、调用 `Get/Poll` 会阻塞。若要非阻塞检查，请使用 `Peek()` 或 `Empty()` 先判断，或通过超时 `Poll(..., someDuration)`.
2. **超时处理**：`Poll(number, timeout)` 用 `time.After(timeout)` 判断超时，可能受 Go 的计时器小数精度影响。
3. **锁开销**：队列使用 `sync.Mutex`，高并发情况下可能成为瓶颈。如果需要更高性能且固定容量，可以考虑本包中的无锁环形缓冲区 `RingBuffer` 。
4. **Dispose 语义**：
   - 一旦调用 `Dispose`, 该队列无法再使用；所有后续 `Put/Get` 都会返回 `ErrDisposed`。
   - 所有阻塞在 `Get/Poll` 的 goroutine 也会立刻被唤醒并返回 `ErrDisposed`。
   - `Dispose()` 返回最后未取出的数据，供调用者自行处理。

---

## 七、总结

- **queue.Queue** 提供了一个可无限扩容、可阻塞、可超时的 FIFO 队列，支持 `Dispose` 语义来强制唤醒并让调用方结束阻塞。
- 关键原理：使用 `sync.Mutex` + “waiters 信号量” 模式来协调阻塞 `Poll/Get`，`Put` 时依次唤醒阻塞读者。
- 与原生 `chan` 相比：
  - 不限容量，任何时刻 `Put` 不会阻塞（除了加锁时短暂等待），而 `chan` 若缓冲满或 0 缓冲时就会阻塞生产者。
  - `Dispose()` 允许立刻唤醒所有等待 goroutine 并返回错误，这在 Go 原生 `chan` 里无法优雅处理（可能用 close(chan) 但语义略不同）。
- 适用于需要**有界或无界**队列，且希望**可自定义超时**及**可关闭**的场景，具备一定的可扩展性和易用性。
