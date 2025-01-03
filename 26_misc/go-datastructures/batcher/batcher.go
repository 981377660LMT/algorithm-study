// 批处理器 (Batcher)。
// 用于累积数据并在满足指定条件（最大条目数、最大字节数、或调用 Flush() / 等待超时）时将数据打包产出

package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {

	// ExampleCalculateBytes: 计算 item 的字节大小，这里简单假设 item 是字符串
	calculateBytes := func(i interface{}) uint {
		s, ok := i.(string)
		if !ok {
			return 0
		}
		return uint(len(s))
	}

	// 1. 创建一个 Batcher:
	//  - 最大等候时间: 2秒
	//  - 每批最多5个item
	//  - 每批最多10字节
	//  - batchChan队列可容纳3个
	//  - calculate函数为ExampleCalculateBytes
	b, err := NewBatcher(2*time.Second, 5, 10, 3, calculateBytes)
	if err != nil {
		panic(err)
	}

	// 2. 启动一个协程, 不断Get批次
	go func() {
		for {
			batch, err := b.Get()
			if err == ErrDisposed {
				fmt.Println("Batcher disposed, stop consumer.")
				return
			}
			// 其它错误处理...
			if err != nil {
				fmt.Println("Get error:", err)
				return
			}
			fmt.Println("Got a batch:", batch)
		}
	}()

	// 3. 放数据
	b.Put("hello") // size=5
	time.Sleep(500 * time.Millisecond)
	b.Put("world") // size=5
	// 此时 total=10 bytes => 等一下, size=10 => 触发 flush?
	// Yes, flush => "hello","world" 立即形成一个batch

	// 4. Sleep后再Put
	time.Sleep(1 * time.Second)
	b.Put("abc")    // size=3
	b.Put("123456") // size=6 => total=9, 未超过10 => 还不flush
	// 2秒到时 => 这时 get() 超时 => 取走 [abc,123456]

	time.Sleep(3 * time.Second)

	// 5. Dispose
	b.Dispose()
	// 之后, Put() => ErrDisposed, Get() => 可能取到空, eventually ErrDisposed
}

// 可重入锁.
type mutex struct {
	// This is really more of a semaphore design, but eh
	// Full -> locked, empty -> unlocked
	lock chan struct{}
}

func newMutex() *mutex {
	return &mutex{lock: make(chan struct{}, 1)}
}

func (m *mutex) Lock() {
	m.lock <- struct{}{}
}

func (m *mutex) Unlock() {
	<-m.lock
}

func (m *mutex) TryLock() bool {
	select {
	case m.lock <- struct{}{}:
		return true
	default:
		return false
	}
}

// Batcher provides an API for accumulating items into a batch for processing.
type Batcher interface {
	// Put adds items to the batcher.
	Put(any) error

	// Get retrieves a batch from the batcher. This call will block until
	// one of the conditions for a "complete" batch is reached.
	Get() ([]any, error)

	// Flush forcibly completes the batch currently being built
	Flush() error

	// Dispose will dispose of the batcher. Any calls to Put or Flush
	// will return ErrDisposed, calls to Get will return an error iff
	// there are no more ready batches.
	Dispose()

	// IsDisposed will determine if the batcher is disposed
	IsDisposed() bool
}

// ErrDisposed is the error returned for a disposed Batcher
var ErrDisposed = errors.New("batcher: disposed")

// CalculateBytes evaluates the number of bytes in an item added to a Batcher.
type CalculateBytes func(any) uint

type basicBatcher struct {
	maxTime  time.Duration // 批次最大等待时间
	maxItems uint          // 批次最大条目数
	maxBytes uint          // 批次最大字节数

	calculateBytes CalculateBytes // 用来计算每个item的字节大小

	disposed       bool       // 是否已销毁
	items          []any      // 当前正在积累的item列表
	batchChan      chan []any // 用于**完成的批次**在此通道中排队等待被 `Get()` 取走
	availableBytes uint       // 当前批次已累计的字节数
	lock           *mutex     // 自定义锁,支持TryLock
}

// NewBatcher creates a new Batcher using the provided arguments.
// Batch readiness can be determined in three ways:
//   - Maximum number of bytes per batch
//   - Maximum number of items per batch
//   - Maximum amount of time waiting for a batch
//
// Values of zero for one of these fields indicate they should not be
// taken into account when evaluating the readiness of a batch.
// This provides an ordering guarantee for any given thread such that if a
// thread places two items in the batcher, Get will guarantee the first
// item is returned before the second, whether before the second in the same
// batch, or in an earlier batch.
func NewBatcher(maxTime time.Duration, maxItems, maxBytes, queueLen uint, calculate CalculateBytes) (Batcher, error) {
	if maxBytes > 0 && calculate == nil {
		return nil, errors.New("batcher: must provide CalculateBytes function")
	}

	return &basicBatcher{
		maxTime:        maxTime,
		maxItems:       maxItems,
		maxBytes:       maxBytes,
		calculateBytes: calculate,
		items:          make([]any, 0, maxItems),
		batchChan:      make(chan []any, queueLen),
		lock:           newMutex(),
	}, nil
}

// Put adds items to the batcher.
func (b *basicBatcher) Put(item any) error {
	b.lock.Lock()
	if b.disposed {
		b.lock.Unlock()
		return ErrDisposed
	}

	b.items = append(b.items, item)
	if b.calculateBytes != nil {
		b.availableBytes += b.calculateBytes(item)
	}
	if b.ready() {
		// To guarantee ordering this MUST be in the lock, otherwise multiple
		// flush calls could be blocked at the same time, in which case
		// there's no guarantee each batch is placed into the channel in
		// the proper order
		b.flush()
	}

	b.lock.Unlock()
	return nil
}

// Get retrieves a batch from the batcher. This call will block until
// one of the conditions for a "complete" batch is reached.
func (b *basicBatcher) Get() ([]any, error) {
	// Don't check disposed yet so any items remaining in the queue
	// will be returned properly.

	var timeout <-chan time.Time
	if b.maxTime > 0 {
		timeout = time.After(b.maxTime)
	}

	select {
	case items, ok := <-b.batchChan:
		// If there's something on the batch channel, we definitely want that.
		if !ok {
			return nil, ErrDisposed
		}
		return items, nil
	case <-timeout:
		// 时间到了，但**我们要注意**可能在超时时刻 `batchChan` 又有了数据，所以需要**循环**判断
		for {
			if b.lock.TryLock() {
				// locked => safe to read b.items.
				// 先检查 batchChan 里是否又塞进来东西
				// 如果空，就直接把当前 items 列表当作批次
				select {
				case items, ok := <-b.batchChan:
					b.lock.Unlock()
					if !ok {
						return nil, ErrDisposed
					}
					return items, nil
				default:
				}

				// If that is unsuccessful, nothing was added to the channel,
				// and the temp buffer can't have changed because of the lock,
				// so grab that
				items := b.items
				b.items = make([]any, 0, b.maxItems)
				b.availableBytes = 0
				b.lock.Unlock()
				return items, nil
			} else {
				// lock被占用 => 说明flush/put正在进行
				// 先尝试 batchChan
				select {
				case items, ok := <-b.batchChan:
					if !ok {
						return nil, ErrDisposed
					}
					return items, nil
				default:
				}
			}
		}
	}
}

// Flush forcibly completes the batch currently being built
func (b *basicBatcher) Flush() error {
	// This is the same pattern as a Put
	b.lock.Lock()
	if b.disposed {
		b.lock.Unlock()
		return ErrDisposed
	}
	b.flush()
	b.lock.Unlock()
	return nil
}

// Dispose will dispose of the batcher. Any calls to Put or Flush
// will return ErrDisposed, calls to Get will return an error iff
// there are no more ready batches. Any items not flushed and retrieved
// by a Get may or may not be retrievable after calling this.
func (b *basicBatcher) Dispose() {
	for {
		if b.lock.TryLock() {
			// We've got a lock
			if b.disposed {
				b.lock.Unlock()
				return
			}

			b.disposed = true
			b.items = nil
			b.drainBatchChan()
			close(b.batchChan)
			b.lock.Unlock()
		} else {
			// Two cases here:
			// 1) Something is blocked and holding onto the lock
			// 2) Something temporarily has a lock
			// For case 1, we have to clear at least some space so the blocked
			// Put/Flush can release the lock. For case 2, nothing bad
			// will happen here
			b.drainBatchChan()
		}

	}
}

// IsDisposed will determine if the batcher is disposed
func (b *basicBatcher) IsDisposed() bool {
	b.lock.Lock()
	disposed := b.disposed
	b.lock.Unlock()
	return disposed
}

// flush adds the batch currently being built to the queue of completed batches.
// flush is not threadsafe, so should be synchronized externally.
func (b *basicBatcher) flush() {
	b.batchChan <- b.items
	b.items = make([]any, 0, b.maxItems)
	b.availableBytes = 0
}

func (b *basicBatcher) ready() bool {
	if b.maxItems != 0 && uint(len(b.items)) >= b.maxItems {
		return true
	}
	if b.maxBytes != 0 && b.availableBytes >= b.maxBytes {
		return true
	}
	return false
}

func (b *basicBatcher) drainBatchChan() {
	for {
		select {
		case <-b.batchChan:
		default:
			return
		}
	}
}
