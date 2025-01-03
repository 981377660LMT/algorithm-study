下面给出对这个可选（Selectable）Future 实现的详细分析，并提供一个简单示例，便于理解其使用方式。此实现提供了一个“可通过 channel 来等待结果或错误”的 Future 模型，同时支持多个 goroutine 并发等待结果。

---

## 一、功能概述

这个包实现了一个可被外部通过 channel 等待的 Future——称为 **Selectable**。其特点如下：

1. **只写一次，多读多播**：一旦通过 `Fill` / `SetValue` / `SetError` 等方法设置了结果，就会关闭一个专门的 `wait` channel；从而任何 goroutine 都能通过 `WaitChan()` 或 `GetResult()` 获得结果或错误，并且**不会再阻塞**。
2. **并发安全**：用 `sync.Mutex` 和原子变量 `filled` 来保证只有一次成功的填充，后续填充只会返回已经设置的错误值（或结果）。
3. **可取消**：调用 `Cancel()` 会将该 Future 标记为 `ErrFutureCanceled`，所有等待的 goroutine 都能感知到这个“取消”结果。
4. **带独立 channel**：`WaitChan()` 返回一个只读的 `chan struct{}`，在 Future 被填充后会关闭，因此可以把它放在 `select` 里，与其他事件一起等待。

---

## 二、主要结构与字段

```go
type Selectable struct {
	m      sync.Mutex
	val    any
	err    error
	wait   chan struct{}
	filled uint32
}
```

- **m**：互斥锁，用于保护 `val`, `err`, `wait` 等字段。
- **val**：填充的结果值，可以是任意类型（`any`）。
- **err**：填充的错误，如果为非 `nil`，表示此 Future 以错误结束。
- **wait**：一个可选的通道，初始可能是 `nil`，在实际需要时才分配。一旦 Future 被填充，就会关闭这个通道。
- **filled**：原子变量，值为 `0` 表示还未填充，`1` 表示已填充。  
  通过 `atomic.LoadUint32` / `atomic.StoreUint32` 来减少锁争用。

---

## 三、关键方法解析

### 1. `WaitChan()`

```go
func (f *Selectable) WaitChan() <-chan struct{} {
    if atomic.LoadUint32(&f.filled) == 1 {
        return closed
    }
    return f.wchan()
}
```

- 如果 `filled` 为 1，表示结果已准备好，直接返回一个预先关闭好的全局 `closed` 通道（这样等待时立即返回）。
- 否则返回 `f.wchan()`，后者会安全地创建（或返回已有）通道。

```go
func (f *Selectable) wchan() <-chan struct{} {
    f.m.Lock()
    if f.wait == nil {
        f.wait = make(chan struct{})
    }
    ch := f.wait
    f.m.Unlock()
    return ch
}
```

- 在需要时创建一个新的 `wait` channel，并返回给调用方。当后面调用 `Fill` 时会 `close` 这个 channel。

### 2. `GetResult() (any, error)`

```go
func (f *Selectable) GetResult() (any, error) {
    if atomic.LoadUint32(&f.filled) == 0 {
        <-f.wchan() // 等待 channel 关闭
    }
    return f.val, f.err
}
```

- 如果尚未填充（`filled=0`），则阻塞在 `<-f.wchan()`，直到 `Fill`（或 `Cancel` 等）把通道关闭。
- 返回填充好的 `val` 和 `err`。对于任何后续调用（或其他 goroutine 并发调用），若 `filled=1`，则不会阻塞，会立刻返回。

### 3. `Fill(v any, e error)`

```go
func (f *Selectable) Fill(v any, e error) error {
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
```

- 如果 `filled` 还没被设置（即 0），就把 `val` 和 `err` 记录下来，并将 `filled` 置为 1。
- 如果之前创建了 `f.wait` 通道，则 `close(w)` 通知所有等待方结果已就绪。然后把 `f.wait` 重置为 `closed`（一个全局已关闭的通道），这样后续 `WaitChan()` 不再需要新的分配。
- 后续任何对 `Fill` 的调用发现 `filled != 0`，就不会再改动 `val` 和 `err`。
- 返回 `f.err` 供调用方参考，一般表示最终确定的错误是什么。

### 4. `SetValue`, `SetError`, `Cancel`

- `SetValue(v any)` 相当于 `Fill(v, nil)`。
- `SetError(e error)` 相当于 `Fill(nil, e)`。
- `Cancel()` 直接置为一个全局的 `ErrFutureCanceled`。

这些方法都基于 `Fill` 内部逻辑，只是封装了“只填充值”或“只填充错误”。

### 5. 全局变量 `closed` 和 `init()`

```go
var closed = make(chan struct{})

func init() {
	close(closed)
}
```

- 一个已提前关闭好的通道，便于在已经填充的场合下，快速返回一个不会再阻塞的通道。

---

## 四、典型使用场景

1. **多播一次性事件**：执行异步操作，填充一个结果，后续多个 goroutine 同时想要该结果时，通过 `GetResult()` 或 `WaitChan()` 获取相同的值或错误。
2. **取消逻辑**：可以在某些超时或特定条件触发时，调用 `f.Cancel()`，让所有等待结果的 goroutine 都收到 `ErrFutureCanceled`。
3. **在 `select` 中等待**：因为有 `WaitChan()`，可以直接把它放在 `select` 里与其他事件竞争：

   ```go
   select {
   case <-f.WaitChan():
       val, err := f.GetResult()
       // do something
   case <-otherSignal:
       // handle other logic
   }
   ```

---

## 五、使用示例

下面以一个示例来演示如何创建并使用 `Selectable` future。在该示例中，我们模拟一个异步操作（例如网络请求、文件读取等），完成后填充到 future；并由两个 goroutine 以不同方式获取结果。

```go
package main

import (
    "fmt"
    "time"
    "sync"
    "example.com/futures" // 假设你将以上代码放在此包路径
)

func main() {
    // 1. 创建一个 Selectable future
    f := futures.NewSelectable()

    // 2. 启动一个异步任务，模拟1秒后得到结果
    go func() {
        time.Sleep(1 * time.Second)
        // f.SetValue("Hello Future") or f.SetError(err)
        f.SetValue("Hello Future")
    }()

    // 3. 用 WaitChan() + select 方式等待结果
    go func() {
        select {
        case <-f.WaitChan():
            val, err := f.GetResult()
            fmt.Println("[Goroutine1] got result:", val, "error:", err)
        case <-time.After(2 * time.Second):
            fmt.Println("[Goroutine1] Timeout!")
        }
    }()

    // 4. 另一goroutine直接调用 GetResult()
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
```

**流程**：

1. `f := futures.NewSelectable()`：创建一个全新的 future，尚未填充。
2. 异步任务在 1 秒后调用 `f.SetValue(...)`，填充结果 `"Hello Future"`。
3. Goroutine1 使用 `select { case <-f.WaitChan(): ... }` 进行等待，一旦填充后 `WaitChan()` 就会被关闭，立刻执行并拿到结果。
4. Goroutine2 直接调用 `GetResult()`，该方法也会在结果填充时返回正确值，不再阻塞。
5. 两者都能拿到相同的结果字符串和 `nil` 错误，且只填充了一次。

---

## 六、实现细节与注意事项

1. **效率**：
   - 通过 `atomic.LoadUint32(&f.filled)` 来减少对互斥锁的使用，但在需要修改状态时仍会加 `f.m.Lock()` 来保障互斥安全。
   - `wait` channel 只会创建一次，如果 Future 最终都没有被任何 goroutine 等待，这个 channel 可能一直是 `nil`。
2. **只填充一次**：
   - 不管是成功 (`SetValue`) 还是失败 (`SetError`), 还是取消 (`Cancel`)，一旦 `filled=1`，就不会再更改。
3. **已关闭的通道**：
   - 当 Future 完成后，会把 `f.wait` 替换为一个全局的已经关闭的通道 `closed`，这样后续不管 `WaitChan()` 调用几次，都返回已关闭通道，确保立即可读。

---

## 七、小结

- **Selectable** 这个 Future 实现适用于“仅一次填充结果、多次读取”的场景，尤其适用于需要 `select` 等待的需求。
- 通过 `WaitChan()` 可以与其他事件/超时并行等待，通过 `GetResult()` 可以阻塞等待或立即获取已完成结果。
- 仅支持“一次完成”，填充后便不会再更新，符合大多数 Future 场景的需求，并且在并发条件下能安全地广播结果给所有等待方。
