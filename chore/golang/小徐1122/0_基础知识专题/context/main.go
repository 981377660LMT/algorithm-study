// !为什么Context要放在函数参数的第一个位置？
// 并发调用链路（多叉树）上下文，两个作用：父子协程并发协调控制+数据存储介质
// 不能滥用并发：如果你不知道协程什么时候结束，那么你就不要用协程(协程泄漏)；必须要有并发控制的意识。
//
// Context树生命周期，单向传递终止特性(如果父亲终止，儿子也会终止)：
//
// 四种Context：
// 1. emptyCtx：
// 2. cancelCtx：取消
// 3. timeoutCtx：超时
// 4. valueCtx：类似threadLocal

// 包context定义了Context类型，它在API边界和进程之间传递截止时间、取消信号和其他请求范围的值。
//
// 服务器的传入请求应该创建一个Context，对服务器的传出调用应该接受一个Context。
// 它们之间的函数调用链必须传播Context，可以选择使用WithCancel、WithDeadline、WithTimeout或WithValue创建派生的Context来替换它。
// 当一个Context被取消时，所有从它派生的Context也会被取消。
//
// WithCancel、WithDeadline和WithTimeout函数接受一个Context（父级）并返回一个派生的Context（子级）和一个CancelFunc。
// 调用CancelFunc会取消子级及其子级，删除父级对子级的引用，并停止任何关联的定时器。
// 未调用CancelFunc会导致子级及其子级泄漏，直到父级被取消或定时器触发。
// go vet工具会检查所有控制流路径上是否使用了CancelFunc。
//
// WithCancelCause函数返回一个CancelCauseFunc，它接受一个错误并将其记录为取消原因。
// 在取消的上下文或其任何子级上调用Cause会检索原因。如果未指定原因，Cause(ctx)会返回与ctx.Err()相同的值。
//
// 使用Context的程序应遵循以下规则，以保持跨包一致的接口并使静态分析工具检查上下文传播：
//
// !不要在结构类型内部存储Context，而是将Context显式地传递给每个需要它的函数。Context应该是第一个参数，通常命名为ctx：
//
//	func DoSomething(ctx context.Context, arg Arg) error {
//		// ... 使用ctx ...
//	}
//
// 即使函数允许，也不要传递nil的Context。如果不确定要使用哪个Context，请传递context.TODO。
//
// 仅将context Values用于跨进程和API传递请求范围数据，而不是将可选参数传递给函数。
//
// 同一个Context可以传递给在不同goroutine中运行的函数；Context可以安全地同时被多个goroutine使用。
//
// 请参阅https://blog.golang.org/context，了解使用Context的服务器示例代码。

package main

import (
	"errors"
	"internal/reflectlite"
	"sync"
	"sync/atomic"
	"time"
)

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() { close(closedchan) }

// go 1.19// Context携带截止时间、取消信号和其他值跨API边界传递。
//
// Context的方法可以同时被多个goroutine调用。
type Context interface {
	// Deadline返回应该取消代表此上下文的工作的时间。当未设置截止时间时，Deadline返回ok==false。连续调用Deadline返回相同的结果。
	Deadline() (deadline time.Time, ok bool)
	// Done返回一个在应该取消代表此上下文的工作时关闭的通道。如果此上下文永远不会被取消，Done可能返回nil。连续调用Done返回相同的值。Done通道的关闭可能是异步的，在取消函数返回后发生。
	//
	// WithCancel安排在调用cancel时关闭Done；WithDeadline安排在截止时间到期时关闭Done；WithTimeout安排在超时时间过去时关闭Done。
	//
	// Done用于在select语句中使用：
	//
	//  // Stream使用DoSomething生成值并将其发送到out，直到DoSomething返回错误或ctx.Done关闭。
	//  func Stream(ctx context.Context, out chan<- Value) error {
	//  	for {
	//  		v, err := DoSomething(ctx)
	//  		if err != nil {
	//  			return err
	//  		}
	//  		select {
	//  		case <-ctx.Done():
	//  			return ctx.Err()
	//  		case out <- v:
	//  		}
	//  	}
	//  }
	//
	// 有关如何使用Done通道进行取消的更多示例，请参阅https://blog.golang.org/pipelines。
	Done() <-chan struct{}
	// 如果Done尚未关闭，则Err返回nil。
	// 如果Done已关闭，Err返回一个非nil的错误，解释原因：
	// 如果上下文被取消，则返回Canceled；如果上下文的截止时间过去，则返回DeadlineExceeded。
	// 在Err返回非nil错误后，连续调用Err将返回相同的错误。
	Err() error
	// Value返回与此上下文关联的键的值，如果未关联值，则返回nil。对具有相同键的连续调用Value返回相同的结果。
	//
	// 仅将上下文值用于跨进程和API边界传递请求范围数据，而不是将可选参数传递给函数。
	//
	// 键标识上下文中的特定值。希望在上下文中存储值的函数通常会在全局变量中分配一个键，然后将该键用作context.WithValue和Context.Value的参数。键可以是任何支持相等性的类型；包应将键定义为未导出类型，以避免冲突。
	//
	// 定义上下文键的包应为使用该键存储的值提供类型安全的访问器：
	//
	// 	// 包user定义了存储在上下文中的User类型。
	// 	package user
	//
	// 	import "context"
	//
	// 	// User是存储在上下文中的值的类型。
	// 	type User struct {...}
	//
	// 	// key是在此包中定义的键的未导出类型。
	// 	// 这可以防止与其他包中定义的键发生冲突。
	// 	type key int
	//
	// 	// userKey是上下文中user.User值的键。它是未导出的；客户端应该使用user.NewContext和user.FromContext，而不是直接使用此键。
	// 	var userKey key
	//
	// 	// NewContext返回携带值u的新上下文。
	// 	func NewContext(ctx context.Context, u *User) context.Context {
	// 		return context.WithValue(ctx, userKey, u)
	// 	}
	//
	// 	// FromContext返回存储在ctx中的User值（如果有）。
	// 	func FromContext(ctx context.Context) (*User, bool) {
	// 		u, ok := ctx.Value(userKey).(*User)
	// 		return u, ok
	// 	}
	Value(key any) any
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
	done     atomic.Value // 存储 <-chan struct{}，标识 context 是否结束
	children map[canceler]struct{}
	err      error
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

type CancelFunc func()

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	if parent == nil {
		panic("nil parent")
	}
	c := &cancelCtx{}
	c.propagateCancel(parent, c)
	return c, func() { c.cancel(true, Canceled) }
}

// cancelCtx 没有实现 Deadline 方法.

// 返回ctx中的chan.
func (c *cancelCtx) Done() <-chan struct{} {
	// double-checking lazy initialization
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

// propagateCancel 方法顾名思义，用以传递父子 context 之间的 cancel 事件：
// 在 propagateCancel 方法内启动一个守护协程，以保证父 context 终止时，该 cancelCtx 也会被终止
//
// 倘若 parent 是不会被 cancel 的类型（如 emptyCtx），则直接返回；
// 倘若 parent 已经被 cancel，则直接终止子 context，并以 parent 的 err 作为子 context 的 err；
// 假如 parent 是 cancelCtx 的类型，则加锁，并将子 context 添加到 parent 的 children map 当中；
// 假如 parent 不是 cancelCtx 类型，但又存在 cancel 的能力（比如用户自定义实现的 context），则启动一个协程，通过多路复用的方式监控 parent 状态，倘若其终止，则同时终止子 context，并透传 parent 的 err.
func (c *cancelCtx) propagateCancel(parent Context, child canceler) {
	c.Context = parent
	done := parent.Done()
	if done == nil {
		return // parent is never canceled
	}
	select {
	case <-done:
		// parent is already canceled
		child.cancel(false, parent.Err())
	default:
	}

	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil {
			child.cancel(false, p.err)
		} else {
			// 父子都是 cancelCtx，添加到父节点的 children set中
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		goroutines.Add(1)
		// monitor parent context
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}

func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return // already canceled
	}
	c.err = err
	d, _ := c.done.Load().(chan struct{})
	// 若 channel 此前未初始化，则直接注入一个 closedChan，否则关闭该 channel；
	if d == nil {
		c.done.Store(closedchan)
	} else {
		close(d)
	}
	for child := range c.children {
		// NOTE: acquiring the child's lock while holding parent's lock.
		child.cancel(false, err)
	}
	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}
}

func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}

// 判断 parent 是否为 cancelCtx 的类型
// 基于 cancelCtxKey 为 key 取值时返回 cancelCtx 自身，是 cancelCtx 特有的协议.
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	pdone, _ := p.done.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return p, true
}

//
//
//
//
//
//
//
//
//
//
// timerCtx 用于实现超时控制

// timerCtx携带一个定时器和截止时间。
// 它嵌入了一个cancelCtx来实现Done和Err。
// 它通过停止定时器然后委托给cancelCtx.cancel来实现取消。
type timerCtx struct {
	cancelCtx
	timer *time.Timer

	deadline time.Time
}

func WithTimeOut(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	if parent == nil {
		panic("nil parent")
	}
	if cur, ok := parent.Deadline(); ok && deadline.Before(cur) {
		// The current deadline is already sooner than the new one.
		return WithCancel(parent)
	}
	c := &timerCtx{
		deadline: deadline,
	}
	c.cancelCtx.propagateCancel(parent, c)
	dur := time.Until(deadline)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded)
		return c, func() { c.cancel(false, Canceled) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}
	return c, func() { c.cancel(true, Canceled) }
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		removeChild(c.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}

// valueCtx 类似于 threadLocal，通过key-value的方式存储数据。
// valueCtx 不适合视为存储介质，存放`大量`的 kv 数据
// 只适合存放少量作用域较大的全局 meta 数据.
type valueCtx struct {
	Context
	key, val any
}

// WithValue returns a copy of parent in which the value associated with key is
// val.
//
// Use context Values only for request-scoped data that transits processes and
// APIs, not for passing optional parameters to functions.
//
// The provided key must be comparable and should not be of type
// string or any other built-in type to avoid collisions between
// packages using context. Users of WithValue should define their own
// types for keys. To avoid allocating when assigning to an
// interface{}, context keys often have concrete type
// struct{}. Alternatively, exported context key variables' static
// type should be a pointer or interface.
func WithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}

func (c *valueCtx) Value(key any) any {
	if c.key == key {
		return c.val
	}
	return value(c.Context, key)
}

// 向上查找 value.
// 由下而上，由子及父，依次对 key 进行匹配；
// 其中 cancelCtx、timerCtx、emptyCtx 类型会有特殊的处理方式；
// 找到匹配的 key，则将该组 value 进行返回.
func value(c Context, key any) any {
	for {
		switch ctx := c.(type) {
		case *valueCtx:
			if key == ctx.key {
				return ctx.val
			}
			c = ctx.Context
		case *cancelCtx:
			if key == &cancelCtxKey {
				return c
			}
			c = ctx.Context
		case *timerCtx:
			if key == &cancelCtxKey {
				return &ctx.cancelCtx
			}
			c = ctx.Context
		case *emptyCtx:
			return nil
		default:
			return c.Value(key)
		}
	}
}
