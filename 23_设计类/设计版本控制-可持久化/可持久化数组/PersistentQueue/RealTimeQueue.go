// https://scrapbox.io/data-structures/Realtime_Queue

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

func init() {
	debug.SetGCPercent(-1)
}

func main2() {
	assert := func(a bool) {
		if !a {
			panic("assert fail")
		}
	}
	queue := NewRealTimeQueue()
	assert(queue.Empty() == true)
	queue1 := queue.Push(1)
	assert(queue1.Empty() == false)
	assert(queue1.Front() == 1)
	queue2 := queue1.Push(2)
	assert(queue2.Empty() == false)
	assert(queue2.Front() == 1)
	queue3 := queue2.Shift()
	assert(queue3.Empty() == false)
	assert(queue3.Front() == 2)
	queue4 := queue3.Shift()
	assert(queue4.Empty() == true)

	time1 := time.Now()
	for i := 0; i < 1e6; i++ {
		queue = queue.Push(i)
		queue = queue.Shift()
		queue.Empty()
	}

	fmt.Println(time.Since(time1)) // 343ms
}

func main() {
	// https://judge.yosupo.jp/problem/persistent_queue
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	gits := make([]*RealTimeQueue, q+1)
	gits[0] = NewRealTimeQueue()

	for i := 0; i < q; i++ {
		var op, time, x int
		fmt.Fscan(in, &op, &time)
		time++
		if op == 0 {
			fmt.Fscan(in, &x)
			gits[i+1] = gits[time].Push(x)
		} else {
			res := gits[time].Front()
			fmt.Fprintln(out, res)
			gits[i+1] = gits[time].Shift()
		}
	}
}

type S = int

// 可持久化队列.
//  https://scrapbox.io/data-structures/Realtime_Queue
type RealTimeQueue struct {
	front    *_Stream
	back     *_PersistentStack
	schedule *_Stream
}

func NewRealTimeQueue() *RealTimeQueue {
	return &RealTimeQueue{
		front:    _NewStream(),
		back:     _NewPersistentStack(),
		schedule: _NewStream(),
	}
}

func (q *RealTimeQueue) Empty() bool {
	return q.front.Empty()
}

func (q *RealTimeQueue) Front() S {
	if q.Empty() {
		panic("queue is empty")
	}
	return q.front.Top()
}

func (q *RealTimeQueue) Push(x S) *RealTimeQueue {
	return q._makeQueue(q.front, q.back.Push(x), q.schedule)
}

func (q *RealTimeQueue) Shift() *RealTimeQueue {
	if q.Empty() {
		return NewRealTimeQueue()
	}
	return q._makeQueue(q.front.Pop(), q.back, q.schedule)
}

func (q *RealTimeQueue) _rotate(f *_Stream, b *_PersistentStack, s *_Stream) *_Stream {
	return &_Stream{_Suspension: _NewSuspensionWith(func() interface{} {
		if f.Empty() {
			return &_Cell{resolved: b.Top(), next: s}
		}
		return &_Cell{resolved: f.Top(), next: q._rotate(f.Pop(), b.Pop(), s.Push(b.Top()))}
	})}
}

func (q *RealTimeQueue) _makeQueue(f *_Stream, b *_PersistentStack, s *_Stream) *RealTimeQueue {
	if !s.Empty() {
		return &RealTimeQueue{front: f, back: b, schedule: s.Pop()}
	}
	tmp := q._rotate(f, b, _NewStream())
	return &RealTimeQueue{front: tmp, back: _NewPersistentStack(), schedule: tmp}
}

type _Cell struct {
	resolved S
	next     *_Stream
}

type _Stream struct {
	*_Suspension
}

// 惰性求值的流.
func _NewStream() *_Stream {
	return &_Stream{_Suspension: _NewSuspension()}
}

func (s *_Stream) Empty() bool {
	return s.Resolve() == nil
}

func (s *_Stream) Top() S {
	return s.Resolve().(*_Cell).resolved
}

func (s *_Stream) Pop() *_Stream {
	return s.Resolve().(*_Cell).next
}

func (s *_Stream) Push(x S) *_Stream {
	return &_Stream{_Suspension: _NewSuspensionWith(&_Cell{resolved: x, next: s})}
}

func (s *_Stream) String() string {
	x := s
	res := []S{}
	for !x.Empty() {
		res = append(res, x.Top())
		x = x.Pop()
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return fmt.Sprintf("Stream%v", res)
}

// 惰性求值.
type _Suspension struct {
	x        interface{}
	resolved interface{}
}

func _NewSuspension() *_Suspension {
	return &_Suspension{}
}

func _NewSuspensionWith(x interface{}) *_Suspension {
	return &_Suspension{x: x}
}

func (s *_Suspension) Resolve() interface{} {
	if s.resolved == nil {
		if f, ok := s.x.(func() interface{}); ok {
			s.resolved = f()
		} else {
			s.resolved = s.x
		}
	}
	return s.resolved
}

type _PersistentStack struct {
	root *_StackNode
}

type _StackNode struct {
	value S
	pre   *_StackNode
}

// 创建一个新的可持久化栈.
func _NewPersistentStack() *_PersistentStack {
	return &_PersistentStack{}
}

func (s *_PersistentStack) Push(value S) *_PersistentStack {
	return &_PersistentStack{root: &_StackNode{value: value, pre: s.root}}
}

func (s *_PersistentStack) Pop() *_PersistentStack {
	if s.root == nil {
		panic("stack is empty")
	}
	return &_PersistentStack{root: s.root.pre}
}

func (s *_PersistentStack) Top() S {
	if s.root == nil {
		panic("stack is empty")
	}
	return s.root.value
}

func (s *_PersistentStack) Empty() bool {
	return s.root == nil
}

func (s *_PersistentStack) Reverse() *_PersistentStack {
	res := _NewPersistentStack()
	x := s
	for !x.Empty() {
		res = res.Push(x.Top())
		x = x.Pop()
	}
	return res
}

func (s *_PersistentStack) String() string {
	sb := []string{}
	x := s
	for !x.Empty() {
		sb = append(sb, fmt.Sprintf("%v", x.Top()))
		x = x.Pop()
	}
	for i, j := 0, len(sb)-1; i < j; i, j = i+1, j-1 {
		sb[i], sb[j] = sb[j], sb[i]
	}
	return fmt.Sprintf("Stack%v", sb)
}
