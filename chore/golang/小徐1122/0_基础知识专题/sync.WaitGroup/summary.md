# Golang并发等待组sync.WaitGroup深度解析

## 一、前言

Golang作为一门高并发语言，其核心并发组件包括：

- **goroutine**：GMP模型下的轻量级协程
- **channel**：goroutine间通信机制
- **context**：跨goroutine生命周期控制
- **sync.WaitGroup**：实现goroutine的等待聚合模式

本文重点解析`sync.WaitGroup`的工作原理、使用场景及底层实现。

---

## 二、场景驱动：等待聚合模式

### 1. 典型场景

**需求**：主goroutine需等待所有子goroutine完成工作后再继续执行。

**传统channel方案**：

```go
func TestChannelWait() {
    const N = 10
    ch := make(chan struct{}, N)

    for i := 0; i < N; i++ {
        go func() {
            defer func() { ch <- struct{}{} }()
            time.Sleep(time.Second)
        }()
    }

    for i := 0; i < N; i++ { <-ch }
}
```

**缺陷**：

- 需预先确定goroutine数量
- Channel信号无法复用
- 多级等待需复杂Channel管理

---

### 2. WaitGroup解决方案

```go
func TestWaitGroup() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            time.Sleep(time.Second)
        }()
    }
    wg.Wait()
}
```

**核心方法**：
| 方法 | 作用 | 调用位置 |
|------------|------------------------------|---------------|
| Add(n) | 计数器+n（登记goroutine） | 主goroutine |
| Done() | 计数器-1（任务完成上报） | 子goroutine |
| Wait() | 阻塞直到计数器清零 | 主goroutine |

---

## 三、高级使用技巧

### 1. 数据聚合模式

**需求**：收集所有子goroutine的结果

**版本3.0优化实现**：

```go
func DataAggregation() {
    const tasks = 10
    dataCh := make(chan int)

    go func() { // 写协程
        var wg sync.WaitGroup
        for i := 0; i < tasks; i++ {
            wg.Add(1)
            go func(i int) {
                defer wg.Done()
                dataCh <- i*2
            }(i)
        }
        wg.Wait()
        close(dataCh)
    }()

    var results []int // 主协程同时作为读协程
    for v := range dataCh {
        results = append(results, v)
    }
    fmt.Println(results)
}
```

**优势**：

- 无需额外同步机制
- 自动保证数据完整性
- 资源高效利用

---

### 2. 工程实践案例（etcd）

**服务发现模块中的优雅关闭**：

```go
type resolver struct {
    wg sync.WaitGroup
}

func (r *resolver) Build() {
    r.wg.Add(1)
    go r.watch()
}

func (r *resolver) watch() {
    defer r.wg.Done()
    // ...监控逻辑...
}

func (r *resolver) Close() {
    r.wg.Wait() // 确保监控协程退出
}
```

---

## 四、源码实现剖析

### 1. 数据结构

```go
type WaitGroup struct {
    noCopy  // 防拷贝标识
    state1 uint64 // 高32位：计数器，低32位：等待者数量
    state2 uint32 // 信号量
}
```

内存布局：

```
|----- state1 (64 bits) -----|-- state2 (32) --|
[ 计数器 (32) | 等待者数 (32) ]   [  信号量      ]
```

---

### 2. 核心方法实现

#### (1) Add(delta int)

```go
func (wg *WaitGroup) Add(delta int) {
    state := atomic.AddUint64(&wg.state1, uint64(delta)<<32)
    v := int32(state >> 32) // 当前计数器值

    if v < 0 { panic("negative counter") }
    if v > 0 || w == 0 { return }

    // 唤醒所有等待者
    for ; w != 0; w-- {
        runtime_Semrelease(&wg.state2, false)
    }
}
```

#### (2) Wait()

```go
func (wg *WaitGroup) Wait() {
    for {
        state := atomic.LoadUint64(&wg.state1)
        if (state>>32) == 0 { return }

        if atomic.CompareAndSwapUint64(&wg.state1, state, state+1) {
            runtime_Semacquire(&wg.state2) // 阻塞等待
            if atomic.LoadUint64(&wg.state1) != 0 {
                panic("reuse before wait return")
            }
            return
        }
    }
}
```

**关键机制**：

- 原子操作保证计数器修改的线程安全
- 信号量实现goroutine的阻塞/唤醒
- 状态值的高效位操作

---

## 五、使用注意事项

| 错误类型       | 示例                  | 解决方案                    |
| -------------- | --------------------- | --------------------------- |
| 过早调用Wait() | Add在子协程中调用     | Add必须在Wait前的主协程调用 |
| 计数器负数     | Done()调用多于Add()   | 保证Add/Done成对出现        |
| 跨轮次重用     | 前一轮未结束就重用    | 等待前一轮Wait完成          |
| 值拷贝导致失效 | 函数传值传递WaitGroup | 始终使用指针传递            |

---

## 六、性能对比

| 特性         | sync.WaitGroup  | Channel方案           |
| ------------ | --------------- | --------------------- |
| 内存开销     | 16字节          | 每个任务需单独Channel |
| 调度开销     | 原子操作+信号量 | Channel发送/接收      |
| 扩展性       | 动态增减计数器  | 需预分配Channel容量   |
| 多级等待支持 | 单次使用        | 可通过多个Channel实现 |

---

## 七、总结

**适用场景**：

- 主协程等待多个子协程完成
- 需要动态增减任务数量
- 无数据传递的纯同步场景

**最佳实践**：

1. Add()调用应在启动子协程前完成
2. 推荐使用defer wg.Done()
3. 避免在子协程中调用Add()
4. 需要数据传递时结合Channel使用
5. 复杂场景可组合context使用

**设计哲学**：

- 最小化同步原语
- 显式优于隐式
- 组合而非继承

通过深入理解WaitGroup的机制，开发者可以构建更高效、可靠的并发系统，在GMP模型中实现精准的协程协同控制。
