下面将对这段 `golang.org/x/sync/semaphore` 包中加权信号量（Weighted Semaphore）的代码进行较为详细的解读和分析。该实现允许在同一把信号量中按“重量”（资源量）的方式获取和释放资源，从而限制并发访问的总资源使用量。

### 整体简介

该加权信号量通过以下方式实现：

- `Weighted` 结构中维护了信号量的最大容量 `size` 和当前已占用的容量 `cur`。
- 当 Goroutine 需要获取一定数量（权重）的资源时，如果当前可用资源足够且没有等待队列，就直接占用资源；否则将请求放入等待队列中。
- 等待队列以 FIFO 顺序组织为一个双向链表 `list.List`，队列中的每个元素都是一个 `waiter`，包含请求的数量和一个 `ready` 通道。当资源可用时，信号量会通知队列中排在前面的请求者（通过关闭 `ready` 通道）以获得资源。
- 通过 `context.Context` 支持超时或取消，如果在获取资源过程中 `ctx.Done()` 触发，等待的 Goroutine 将被取消并从队列中移除。

### 结构体定义

```go
type waiter struct {
	n     int64        // 请求的资源数量（权重）
	ready chan<- struct{} // 当可获取资源时关闭该通道通知调用方
}
```

`waiter` 表示排队等待获取资源的请求者，`ready` 通道在资源可用时被关闭，表示该 `waiter` 可以继续执行。

```go
type Weighted struct {
	size    int64       // 信号量最大可用资源量
	cur     int64       // 当前已分配的资源量
	mu      sync.Mutex  // 保护 size, cur, waiters 的互斥锁
	waiters list.List    // 等待队列，元素为 waiter
}
```

`Weighted` 表示加权信号量：

- `size` 是总容量。
- `cur` 表示已经被占用的资源量。
- `waiters` 是一个双向链表存放等待者队列。

### 创建信号量

```go
func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}
```

传入最大容量 `n`，构造 `Weighted` 对象即可。

### Acquire方法

```go
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
	done := ctx.Done()

	s.mu.Lock()
	select {
	case <-done:
		// 如果在获取锁之前或刚获取锁时 ctx 就已经完成(取消或超时),
		// 那么无需尝试获取资源，直接返回错误。
		s.mu.Unlock()
		return ctx.Err()
	default:
	}
```

进入 Acquire 时，先加锁，然后检查 `ctx` 是否已经被取消。如果已取消，直接退出。

```go
	if s.size - s.cur >= n && s.waiters.Len() == 0 {
		// 当前可用资源足够，并且没有等待队列，直接占用资源即可
		s.cur += n
		s.mu.Unlock()
		return nil
	}
```

如果当前可用资源足够而且没有人等待，则直接分配资源给当前请求者。

```go
	if n > s.size {
		// 如果请求资源量大于总容量，不可能满足，因此立即返回错误。
		s.mu.Unlock()
		<-done
		return ctx.Err()
	}
```

如果请求比总容量还大，这个请求永远无法满足，不必进入等待队列，直接等 ctx 完成后返回错误。

```go
	// 走到这里说明可用资源不足或有等待队列，需要排队
	ready := make(chan struct{})
	w := waiter{n: n, ready: ready}
	elem := s.waiters.PushBack(w)
	s.mu.Unlock()
```

创建 `waiter`，把当前请求放入等待队列的末尾，然后解锁。当被通知资源可用时，会关闭 `ready` 通道。

```go
	select {
	case <-done:
		// ctx 在等待过程中被取消
		s.mu.Lock()
		select {
		case <-ready:
			// 如果在 ctx 被取消前已经拿到了资源(ready 已被关闭),
			// 则需要撤销之前的资源分配(s.cur -= n)并通知其他等待者。
			s.cur -= n
			s.notifyWaiters()
		default:
			// 如果还没获取到资源则从队列移除当前 waiter
			isFront := s.waiters.Front() == elem
			s.waiters.Remove(elem)
			// 如果移除的是队首 waiter，且当前有可用资源，需要通知下一个等待者。
			if isFront && s.size > s.cur {
				s.notifyWaiters()
			}
		}
		s.mu.Unlock()
		return ctx.Err()

	case <-ready:
		// 成功获取到资源了 (ready 被关闭意味着此 waiter 已被满足)
		// 但依然需要检查 ctx，因为可能在获取资源同时 ctx 被取消。
		select {
		case <-done:
			// 如果 ctx 已取消，则释放刚刚获取到的资源
			s.Release(n)
			return ctx.Err()
		default:
		}
		return nil
	}
```

在解锁后，会阻塞在 `select`，等待两个事件之一：

1. `<-done`: 上下文被取消/超时

   - 如果在未获取资源前取消了，就从等待队列中移除请求。
   - 如果取消时已经成功获取资源（`ready` 已经关闭），就需要把已分配的资源还回去（`cur -= n`）。
   - 最终返回 `ctx.Err()`。

2. `<-ready`: 资源已获得
   - 此时资源计数已经增加到 `s.cur` 内。检查 `ctx` 是否同时已取消，如果取消则调用 `Release` 还回资源。
   - 若没取消则正常返回 nil。

### TryAcquire 方法

```go
func (s *Weighted) TryAcquire(n int64) bool {
	s.mu.Lock()
	success := s.size - s.cur >= n && s.waiters.Len() == 0
	if success {
		s.cur += n
	}
	s.mu.Unlock()
	return success
}
```

`TryAcquire` 不会等待资源，仅在能立即获取的情况下返回 true，否则返回 false。

- 加锁检查当前是否有足够资源并且无等待者。
- 若满足则直接增加 `cur`；否则返回失败。

### Release 方法

```go
func (s *Weighted) Release(n int64) {
	s.mu.Lock()
	s.cur -= n
	if s.cur < 0 {
		s.mu.Unlock()
		panic("semaphore: released more than held")
	}
	s.notifyWaiters() // 释放资源后可能使得某些等待者可被满足，尝试通知
	s.mu.Unlock()
}
```

释放 `n` 个资源，减少 `cur`。如果释放后 `cur` <0，说明释放的资源比持有的还多，这是逻辑错误，引发 panic。  
然后调用 `s.notifyWaiters()` 尝试为队列中阻塞的请求者分配资源。

### notifyWaiters 方法

```go
func (s *Weighted) notifyWaiters() {
	for {
		next := s.waiters.Front()
		if next == nil {
			break // 没有等待者
		}

		w := next.Value.(waiter)
		if s.size - s.cur < w.n {
			// 下一个等待者需要的资源量仍不足，不能满足。
			// 不继续查找下一个等待者，防止大的请求永远饿死。
			break
		}

		// 可以满足下一个等待者
		s.cur += w.n
		s.waiters.Remove(next)
		close(w.ready) // 通知该等待者已获取资源
	}
}
```

`notifyWaiters` 从队首开始尝试为等待的请求分配资源：

1. 获取队首 `waiter` 的需求 `w.n`。
2. 检查是否有足够的可用资源（`size - cur >= w.n`）。
3. 如果足够：

   - 增加 `cur`。
   - 移除该 `waiter`。
   - 关闭其 `ready` 通道让其获得资源。
   - 继续下一轮，试图唤醒更多的等待者。

4. 如果不够满足当前队首 `waiter`，就停止尝试，即使后面有可能存在要求更少资源的 `waiter`，也不前进队列。这样可以避免饥饿现象，保证队列顺序(FIFO)公正性。

### 设计思路与特点

- **FIFO 语义**：通过队列顺序保证等待顺序，不允许“跳过”前面的请求以满足后面的轻量请求，否则会有可能导致大的请求永远得不到满足。
- **使用 channel 通知**：当资源可用时，信号量通过 `close(w.ready)` 通知等待者，从而实现类似 "锁+条件变量" 的语义，但更简单直接，等待方只需 `<-ready` 即可被唤醒。

- **取消和超时支持**：通过 `context.Context`，如果在等待资源的过程中 `ctx` 被取消，那么该 `waiter` 会被移除队列，不会阻塞后续请求。

- **性能优化**：在资源足够且无等待者的情况下，`Acquire` 不会创建 goroutine 或消耗额外内存，直接快速返回。

### 总结

该加权信号量的实现要点在于：

1. 使用 `size` 表示总资源量，`cur` 表示当前已分配的资源。
2. 使用一个双向链表作为等待队列，每个等待者用 `waiter` 表示并有一个 `ready` 通道。
3. 在可直接满足请求时立即分配资源，否则将请求加入队列等待。
4. 当资源有释放或新的资源可用时，通过 `notifyWaiters` 顺序唤醒等待队列前端的请求。
5. 利用 `context.Context` 支持取消与超时，如果 `ctx` 在等待过程中被取消，尝试安全移除该等待者。
6. 严格的队列化处理，防止饿死问题，保证 FIFO 和公平性。

通过这些机制，实现一个高效并且可取消的加权信号量控制并发访问资源的总量。
