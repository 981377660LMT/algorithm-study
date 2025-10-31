// https://maspypy.github.io/library/ds/node_pool.hpp
// 线程安全：可在 NodePool 方法外层加 sync.Mutex 包裹，或改为每线程私有池。
// 深拷贝：若 T 包含指针并要求深拷贝，可把 Clone 改成接收自定义 copier：CloneWith(r, func(dst *T, src *T) { ... })。

package main

import "fmt"

func main() {
	type Node struct{ A, B int }
	pool := NewNodePool[Node](0) // 使用默认 1<<16

	// Create / 写入
	r1 := pool.Create()
	n1 := r1.Ptr()
	n1.A, n1.B = 1, 2

	// CreateWith
	r2 := pool.CreateWith(Node{A: 3, B: 4})

	// Clone
	r3 := pool.Clone(r2)

	fmt.Println(*r1.Ptr(), *r2.Ptr(), *r3.Ptr()) // {1 2} {3 4} {3 4}

	// Destroy
	pool.Destroy(r1)
	pool.Destroy(r2)
	pool.Destroy(r3)

	// Reset（不释放内存，复用第一块）
	pool.Reset()
}

// NodePool 是一个通用对象池，按固定大小的 chunk 分配 Slot。
// 每个 Slot 保存一个值和一个空闲链表指针，Destroy 后可 O(1) 回收。
// 非并发安全；如需并发，请在外层加锁。
type NodePool[T any] struct {
	constChunkSize int

	chunks   [][]nodeSlot[T] // 若干块，每块固定容量
	cur      []nodeSlot[T]   // 当前块
	curUsed  int             // 当前块已用数量
	freeHead *nodeSlot[T]    // 空闲链表头
}

type nodeSlot[T any] struct {
	next *nodeSlot[T]
	val  T
}

// Ref 是对池中对象的句柄。通过 Ref.Ptr() 可拿到 *T 直接读写。
type Ref[T any] struct{ s *nodeSlot[T] }

// NewNodePool 创建一个对象池；chunkSize<=0 时使用默认 1<<16。
func NewNodePool[T any](chunkSize int) *NodePool[T] {
	if chunkSize <= 0 {
		chunkSize = 1 << 16
	}
	p := &NodePool[T]{constChunkSize: chunkSize}
	p.allocChunk()
	return p
}

// Create 分配一个新对象（零值初始化），返回句柄。
func (p *NodePool[T]) Create() Ref[T] {
	s := p.newSlot()
	var zero T
	s.val = zero
	return Ref[T]{s: s}
}

// CreateWith 用给定值初始化，返回句柄。
func (p *NodePool[T]) CreateWith(v T) Ref[T] {
	s := p.newSlot()
	s.val = v
	return Ref[T]{s: s}
}

// Clone 复制一个对象（浅拷贝）。
func (p *NodePool[T]) Clone(r Ref[T]) Ref[T] {
	if r.s == nil {
		panic("Clone: nil Ref")
	}
	s := p.newSlot()
	s.val = r.s.val
	return Ref[T]{s: s}
}

// Destroy 归还到对象池。
func (p *NodePool[T]) Destroy(r Ref[T]) {
	if r.s == nil {
		return
	}
	// 可选：清零值，帮助 GC 更快回收引用对象
	var zero T
	r.s.val = zero

	r.s.next = p.freeHead
	p.freeHead = r.s
}

// Reset 保留已分配内存，仅把可用位置复位到第一块开头，并清空空闲链表。
func (p *NodePool[T]) Reset() {
	p.freeHead = nil
	if len(p.chunks) > 0 {
		p.cur = p.chunks[0]
		p.curUsed = 0
	}
}

// Ptr 返回对象地址，可直接读写对象内容。
func (r Ref[T]) Ptr() *T {
	if r.s == nil {
		return nil
	}
	return &r.s.val
}

func (p *NodePool[T]) allocChunk() {
	ch := make([]nodeSlot[T], p.constChunkSize)
	p.chunks = append(p.chunks, ch)
	p.cur = ch
	p.curUsed = 0
}

func (p *NodePool[T]) newSlot() *nodeSlot[T] {
	if p.freeHead != nil {
		s := p.freeHead
		p.freeHead = p.freeHead.next
		s.next = nil
		return s
	}
	if p.curUsed == len(p.cur) {
		p.allocChunk()
	}
	s := &p.cur[p.curUsed]
	p.curUsed++
	return s
}
