下面给出对这个 **优先队列（PriorityQueue）** 实现的详细分析，以及一个简要用法示例。该优先队列与普通 `Queue` 类似，同样支持阻塞获取、`Dispose` 语义等，但底层采用 **最小堆**（min-heap）来存储并管理元素的优先级顺序。

---

## 一、功能与特性

1. **优先级顺序**：

   - 队列中的元素类型必须实现 `Item` 接口，即 `Compare(other Item) int`，返回 `-1/0/1` 表示小于/等于/大于 `other`。
   - 内部使用最小堆保存元素，“最小”的（Compare(...) 返回负的）元素在堆顶，会被最先取出。

2. **阻塞读取**：

   - 若队列中无元素，则 `Get(...)` 会阻塞，直到其他协程 `Put` 了新元素或队列被 `Dispose`。
   - 被唤醒后即可读取到当前最小（优先级最高）的元素列表。

3. **支持或禁止重复元素**：

   - 初始化时可指定 `allowDuplicates`。
   - 若不允许重复，会用 `itemMap` 跟踪已存在的 `Item`，重复 `Put` 时直接忽略。

4. **Dispose 语义**：

   - 一旦队列被 `Dispose()`，则标记 `disposed = true`；后续任何 `Put/Get` 均返回 `ErrDisposed`，且所有阻塞的 `Get` 也会立即被唤醒并返回错误。

5. **并发安全**：
   - 同 `queue.Queue`，内部通过一个 `sync.Mutex` 来保护堆 (`items`) 与等待队列 (`waiters`)；不使用更底层的 CAS/无锁结构。

---

## 二、核心数据结构

### 1. `PriorityQueue`

```go
type PriorityQueue struct {
    waiters         waiters
    items           priorityItems
    itemMap         map[Item]struct{}
    lock            sync.Mutex
    disposeLock     sync.Mutex
    disposed        bool
    allowDuplicates bool
}
```

- **items**：`priorityItems` 类型（`[]Item`），内部维护一个 **最小堆** (min-heap)。
- **itemMap**：用于跟踪已存在的 `Item`，在 `allowDuplicates=false` 时保证去重。
- **waiters**：当队列为空有协程想要 `Get` 时，会阻塞在这里；一旦 `Put` 添加了新元素，就唤醒一个阻塞的协程。
- **disposed**：标记是否已被 `Dispose`。若是，则所有操作返回 `ErrDisposed`。
- **lock**：读写队列与 waiters 时加锁，避免并发冲突。
- **disposeLock**：在 `Dispose()` 时也会加锁，保证原子性。

### 2. `Item` 接口

```go
type Item interface {
    Compare(other Item) int
}
```

- 要插入优先队列的元素必须实现此接口。
- `Compare`：若 `this.Compare(other)` < 0，表示 `this` 优先级**更高**（本实现中 `this` 被当成“更小”的元素）。
- 优先队列存储时，以**最小元素**在堆顶的方式做排序。

### 3. `priorityItems` — 小顶堆

```go
type priorityItems []Item
func (items *priorityItems) push(item Item) { ... }   // 上滤
func (items *priorityItems) pop() Item       { ... }  // 下滤
func (items *priorityItems) swap(i, j int)   { ... }
func (items *priorityItems) get(number int) []Item    { ... }
```

- `push(item Item)`：将元素放到末尾，然后向上冒泡 (bubble up)，保持堆的最小性质。
- `pop()`：将堆首元素（最小）弹出，用末尾元素放到堆首，再向下筛 (bubble down) 恢复堆性质。
- 通过 `(*items)[i].Compare((*items)[j]) < 0` 判断优先级高低。

**示意**：

- **小顶堆**结构：最小的 `Item` 保持在 `index=0`。
- 访问最小元素只需看 `items[0]`。取出时做 heapify。

---

## 三、主要方法

### 1. `Put(items ...Item) error`

- 加锁 `pq.lock.Lock()` -> 若 `disposed` 则返回 `ErrDisposed`。
- 对传入的每个 `Item`：
  - 如果 `allowDuplicates==true` 或 `itemMap[item]` 不存在：
    - 插入堆 `pq.items.push(item)`
    - 若 `!allowDuplicates` 则 `pq.itemMap[item] = struct{}{}` 记入 map
- 唤醒等待队列 `pq.waiters`：
  - 循环取 `sema := pq.waiters.get()`（可能多个协程在等待）
  - 对每个 sema：先 `sema.response.Add(1)`，再 `sema.ready <- true` 发信号，然后 `sema.response.Wait()` 等待 `Get` 端 `Done()`
  - 如果堆已空就停止唤醒。
- 释放锁并返回。

### 2. `Get(number int) ([]Item, error)`

- 若 `number < 1`，直接返回空。
- 加锁：若 `disposed` 返回 `ErrDisposed`。
- 若堆不为空：
  - 从堆顶开始 `pop` 最多 `number` 个 (`pq.items.get(number)`)，并从 `pq.itemMap` 删除这些元素（若不允许重复）。
  - 返回这些元素。
- 否则（堆为空），创建一个新的信号量 `sema := newSema()` 放进 `pq.waiters`，解锁。
- 等待 `<-sema.ready` 信号，如果队列被 `Dispose`（或者**被 Put**的线程写信号），当前阻塞会结束。
  - 再检查 `Disposed()`：若处于释放状态返回 `ErrDisposed`。
  - 否则再次从堆取出 `number` 个，清除 `itemMap`，`sema.response.Done()` 并返回。

### 3. `Peek() Item`

- 加锁：若 `items` 非空，返回 `items[0]`（堆顶即最小元素）。若为空则返回 `nil`。

### 4. `Empty() / Len()`

- 分别查看 `len(pq.items)` 是否为 0，或直接返回长度。

### 5. `Dispose()`

- 加锁 -> 同时上 `disposeLock` -> `pq.disposed = true`。
- 唤醒所有在 `waiters` 里的 `sema`（`waiter.ready <- true`），让阻塞的协程返回 `ErrDisposed`。
- 置 `pq.items = nil`, `pq.waiters = nil`，释放资源。

---

## 四、示例用法

假设我们要实现一个**最小优先级**结构（数值越小优先级越高），定义一个 `MyItem`：

```go
package main

import (
    "fmt"
    "workiva/queue"
)

type MyItem int
func (a MyItem) Compare(b queue.Item) int {
    x, y := int(a), int(b.(MyItem))
    switch {
    case x < y:
        return -1
    case x > y:
        return 1
    default:
        return 0
    }
}

func main() {
    pq := queue.NewPriorityQueue(10, false) // 初始容量10, 不允许重复

    // 1. Put
    pq.Put(MyItem(5), MyItem(2), MyItem(8), MyItem(2)) // 2 will be inserted once if allowDuplicates=false

    // 2. Get
    items, err := pq.Get(2) // get 2 smallest
    if err != nil {
        fmt.Println("Get error:", err)
        return
    }
    fmt.Println("Got items:", items) // Might be [2,5] in ascending order

    // 3. Peek
    top := pq.Peek()
    fmt.Println("Peek top:", top) // Could be 8 if only 8 left

    // 4. Dispose
    pq.Dispose()
    fmt.Println("Is disposed?", pq.Disposed())
}
```

**输出示例**：

```
Got items: [2 5]
Peek top: 8
Is disposed? true
```

---

## 五、关键点总结

1. **最小堆**实现：`priorityItems` 用 heapify (“bubble up/down”) 维持堆的最小性质，`pop()` 总是弹出**最小** `Item`。
2. **阻塞获取**：若队列为空，`Get()` 会把调用方挂到 `waiters`，等待下一次 `Put()` 或 `Dispose()` 发信号。
3. **可选去重**：`allowDuplicates=false` 时使用 `itemMap` 跟踪已插入 `Item`；再度插入相同 `Item` 将被忽略。
4. **Dispose**：在 `Dispose()` 过程中，会设置 `dispsed=true` 并唤醒所有等待 `Get` 的协程让它们返回 `ErrDisposed`。后续 `Put/Get` 都直接报错。

这使得 `PriorityQueue` 成为一个**线程安全**、**支持阻塞**、**可选去重**、**最小优先**队列，适合在多协程环境中处理带有优先级的任务调度、事件处理等场景。
