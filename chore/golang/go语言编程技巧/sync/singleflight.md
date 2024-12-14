下面将详细讲解 `golang.org/x/sync/singleflight` 包中 `Group` 的代码实现。该包的功能是提供一种“单次请求合并”机制：当多个协程同时对同一个 key 发起相同的操作（函数调用）时，`singleflight` 确保只会对这个 key 实际执行一次函数，其他相同请求的协程会等待第一个请求完成并共享其结果，从而避免重复计算与资源浪费。

### 核心思想

- 当多个 goroutine 同时对同一 `key` 调用 `Do(key, fn)` 时，只有第一个 goroutine 的调用会真正执行 `fn()`，其余 goroutine 都会阻塞等待。
- 执行完成后，所有等待 goroutine 都能获得相同的结果和错误。
- 如果执行函数时出现 panic 或 `runtime.Goexit()`，也会被特殊处理，使得等待的 goroutine 不会永久阻塞。
- 通过 `Forget(key)` 可以让 `Group` 忘记对某个 key 的当前执行状态，下次调用 `Do` 时依然会重新执行。

### 数据结构定义

```go
type Group struct {
    mu sync.Mutex
    m  map[string]*call // 用 key 映射到正在执行或已完成的请求
}
```

`Group` 是管理并发请求的核心结构，内部用一个互斥锁 `mu` 和 `map` 来存储当前的 `call` 对象。

```go
type call struct {
    wg sync.WaitGroup
    val interface{}
    err error

    dups  int              // 有多少个重复请求在等待
    chans []chan<- Result  // 等待结果的通道列表
}
```

`call` 表示一个正在进行的或已经完成的单次请求。`wg` 用来实现等待机制：在 `fn()` 执行完毕前，`WaitGroup` 不会 Done，其它等待的协程通过 `wg.Wait()` 来同步等待结果。

- `val` 和 `err` 是该请求最终的返回值和错误，只会设置一次。
- `dups` 统计重复请求的数量（除了第一个调用 `Do` 的 goroutine，其余都是重复）。
- `chans` 列表中保存需要返回结果的通道，这在 `DoChan` 的实现中起作用。

### panic 和 runtime.Goexit 处理

在执行 `fn()` 时可能会发生 panic 或调用 `runtime.Goexit()`。源码中定义了相关错误和处理逻辑：

```go
var errGoexit = errors.New("runtime.Goexit was called")

type panicError struct {
    value interface{}
    stack []byte
}

func newPanicError(v interface{}) error {
    stack := debug.Stack()
    // 移除goroutine信息的首行，留后续可读堆栈信息
    if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
        stack = stack[line+1:]
    }
    return &panicError{value: v, stack: stack}
}
```

`panicError` 用于包装从 `fn()` 中 recover 得到的 panic 值和堆栈信息。  
`errGoexit` 标识 `fn()` 中调用了 `runtime.Goexit()` 的特殊情况。

### Do 方法实现流程

```go
func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
    g.mu.Lock()
    if g.m == nil {
        g.m = make(map[string]*call)
    }

    // 如果 key 已经有正在进行的 call
    if c, ok := g.m[key]; ok {
        c.dups++      // 重复请求计数+1
        g.mu.Unlock()
        c.wg.Wait()   // 等待第一个请求完成

        // 完成后拿到统一的结果
        if e, ok := c.err.(*panicError); ok {
            panic(e) // 如果是panicError，再次panic，让调用方感知到异常
        } else if c.err == errGoexit {
            runtime.Goexit()
        }
        return c.val, c.err, true // shared=true 表示与他人共享结果
    }

    // 若 key 没有正在进行的 call，则新建一个
    c := new(call)
    c.wg.Add(1)
    g.m[key] = c
    g.mu.Unlock()

    // 实际执行 fn() 的方法
    g.doCall(c, key, fn)
    return c.val, c.err, c.dups > 0
}
```

当调用 `Do(key, fn)` 时：

1. 上锁查看 map 中是否已有该 key 对应的 call。

   - 若有，表示已有一个 goroutine 正在执行 `fn()`，则当前 goroutine 只是重复请求(`dups++`)，随后解锁并 `c.wg.Wait()` 等待结果。
   - 若无，说明本次是第一个请求者，将新建 `call`，添加到 `m` 中，然后解锁并继续执行 `doCall`。

2. `doCall` 执行 `fn()`，执行完成后 `WaitGroup` Done，并从 map 中删除 `call`。

3. `Do` 方法返回时可以直接拿到 `c.val`、`c.err`，并根据 `c.dups` 判断结果是否共享给其他调用方。

### doCall 方法

```go
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
    normalReturn := false
    recovered := false

    defer func() {
        // 此defer在 fn() 执行结束或panic后运行
        if !normalReturn && !recovered {
            // 没有正常返回且没有recover成功，说明是 runtime.Goexit 调用
            c.err = errGoexit
        }

        g.mu.Lock()
        defer g.mu.Unlock()
        c.wg.Done() // 唤醒等待的 goroutine

        if g.m[key] == c {
            delete(g.m, key) // 从 map 中清除该 call
        }

        // 根据 c.err 类型进行处理
        if e, ok := c.err.(*panicError); ok {
            // 如果发生panicError, 并且有等待通道，会通过协程再次panic确保不会死锁
            // 没有等待通道则直接panic
            if len(c.chans) > 0 {
                go panic(e)
                select {} // 阻塞以保证panic状态可见
            } else {
                panic(e)
            }
        } else if c.err == errGoexit {
            // runtime.Goexit 直接结束当前goroutine
        } else {
            // 正常完成的情况，通知所有等待通道发送结果
            for _, ch := range c.chans {
                ch <- Result{c.val, c.err, c.dups > 0}
            }
        }
    }()

    func() {
        defer func() {
            if !normalReturn {
                // 如果此时 normalReturn 还是 false，则可能是panic发生了
                if r := recover(); r != nil {
                    // 将panic信息包装为 panicError
                    c.err = newPanicError(r)
                }
            }
        }()
        c.val, c.err = fn()   // 执行用户函数
        normalReturn = true   // 标记函数正常返回
    }()

    if !normalReturn {
        recovered = true
    }
}
```

`doCall` 的核心流程：

1. 使用 `defer` 嵌套，区分正常返回、panic、或 `runtime.Goexit` 的情况。
2. 尝试执行 `fn()` 函数：
   - 如果 `fn()` 正常返回，则 `normalReturn = true`。
   - 如果 `fn()` 内产生 panic，则在 recover 中设置 `c.err = panicError`。
   - 如果 `fn()` 调用了 `runtime.Goexit()`，则函数无正常返回也无panic，`normalReturn` 和 `recovered` 均为false，于是最终会将 `c.err` 设置为 `errGoexit`。
3. 在主 `defer` 中根据执行结果清理状态：
   - `c.wg.Done()` 允许等待在 `c.wg.Wait()` 的其他 goroutine 获得结果。
   - 从 `g.m` 中删除该 key。
   - 若有等待通道 `c.chans`，发送 `Result` 通知调用方。
   - 若发生 panic，则再次 panic 以使上层感知错误。
   - 若是 `runtime.Goexit`，则让当前 goroutine 结束执行。

### DoChan 方法

```go
func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result {
    ch := make(chan Result, 1)
    g.mu.Lock()
    if g.m == nil {
        g.m = make(map[string]*call)
    }
    if c, ok := g.m[key]; ok {
        c.dups++
        c.chans = append(c.chans, ch)
        g.mu.Unlock()
        return ch
    }
    c := &call{chans: []chan<- Result{ch}}
    c.wg.Add(1)
    g.m[key] = c
    g.mu.Unlock()

    go g.doCall(c, key, fn)
    return ch
}
```

`DoChan` 和 `Do` 类似，但是不会阻塞等待结果，而是返回一个通道。等待 `fn()` 执行完毕时，会将 `Result` 发送到通道 `ch`。重复请求也会共享相同的通道切片 `c.chans`，在 `fn()` 完成后逐一通知。

### Forget 方法

```go
func (g *Group) Forget(key string) {
    g.mu.Lock()
    delete(g.m, key)
    g.mu.Unlock()
}
```

`Forget` 允许主动移除该 key 的 call。这样若之后再调用 `Do`，则会重新执行 `fn()`，不会等待之前的结果。

### 小结

`singleflight` 的设计核心是：

1. 使用 `map[string]*call` 来跟踪特定 key 对应的执行状态。
2. 当一个 key 的执行正在进行时，其他请求仅需等待完成，而不再次执行同样的操作。
3. 执行完成后，所有等待方共享相同的结果和错误。
4. 通过特定的 panic 恢复策略，确保不会导致死锁或等待方永久阻塞。
5. 提供 `Do`（同步等待结果）和 `DoChan`（异步返回通道）两种调用模式。

这种机制在需要防止缓存穿透（例如同一时间大量请求落在同一个 key 上，避免重复查询数据库或访问外部API）或需要避免重复性计算的场景中非常有用。
