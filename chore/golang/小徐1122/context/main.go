// !为什么Context要放在函数参数的第一个位置？
// 并发调用链路（多叉树）上下文，两个作用：父子协程并发协调控制+数据存储介质
// 不能滥用并发：如果你不知道协程什么时候结束，那么你就不要用协程(协程泄漏)
// Context树生命周期，单向传递终止特性
//
// 四种Context：
// 1. emptyCtx：
// 2. cancelCtx：取消
// 3. timeoutCtx：超时
// 4. valueCtx：类似threadLocal

package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
}

// go 1.19
type Context interface {
	Deadline() (deadline time.Time, ok bool) // 返回 context 的截止时间, 可选
	Done() <-chan struct{}                   // 返回一个只读 channel，当 context 被取消或者超时时，该 channel 会被关闭
	Err() error                              // 返回 context 结束的原因，它只会在 Done channel 被关闭后才会返回非 nil 值
	Value(key any) any                       // 返回 context 中 key 对应的 value
}

var Canceled = errors.New("context canceled")
var DeadlineExceeded error = deadlineExceededError{}
var goroutines atomic.Int32
var cancelCtxKey int // symbol

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }

type emptyCtx int // 根节点，儿子通过WithCancel、WithTimeout、WithValue增加能力

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) { return }
func (*emptyCtx) Done() <-chan struct{}                   { return nil } // 读取会阻塞
func (*emptyCtx) Err() error                              { return nil }
func (*emptyCtx) Value(key any) any                       { return nil }

// 仅语义上的区别
var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func Background() Context { return background }
func TODO() Context       { return todo }

type cancelCtx struct {
	Context // 父节点

	mu       sync.Mutex
	done     atomic.Value // 存储 <-chan struct{}
	children map[canceler]struct{}
	err      error
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

func WithCancel(parent Context) (ctx Context, cancel func()) {
	if parent == nil {
		panic("nil parent")
	}
	c := &cancelCtx{}
	c.propagateCancel(parent, c)
	return c, func() { c.cancel(true, Canceled) }
}

// cancelCtx 没有实现 Deadline 方法.

func (c *cancelCtx) Done() <-chan struct{} {
	// double-check initialization
	d := c.done.Load()
	if d != nil {
		return d.(chan struct{})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	d = c.done.Load()
	if d == nil {
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *cancelCtx) Value(key any) any {
	// 若 key 特定值 &cancelCtxKey，则返回 cancelCtx自身.
	// 仅内部使用.
	if key == &cancelCtxKey {
		return c
	}
	return value(c.Context, key)
}

// 向上查找 value.
func value(c Context, key any) any {}

// propagateCancel 方法顾名思义，用以传递父子 context 之间的 cancel 事件：
// 在 propagateCancel 方法内启动一个守护协程，以保证父 context 终止时，该 cancelCtx 也会被终止
func (c *cancelCtx) propagateCancel(parent Context, child canceler) {
	c.Context = parent
}

func (c *cancelCtx) cancel(removeFromParent bool, err error) {}
