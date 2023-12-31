// https://37zigen.com/bankers-queue/
// Okasaki C. Purely functional data structures[M]. Cambridge University Press, 1999.

package main

import (
	"fmt"
	"time"
)

type V = int32

func main() {
	// demo()
	timeit()
}

func yosupo() {
	// https://judge.yosupo.jp/problem/persistent_queue
	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// var q int
	// fmt.Fscan(in, &q)

	// gits := make([]*BankerQueue, q+1)
	// gits[0] = NewRealTimeQueue()

	// for i := 0; i < q; i++ {
	// 	var op, time, x int
	// 	fmt.Fscan(in, &op, &time)
	// 	time++
	// 	if op == 0 {
	// 		fmt.Fscan(in, &x)
	// 		gits[i+1] = gits[time].Push(x)
	// 	} else {
	// 		res := gits[time].Front()
	// 		fmt.Fprintln(out, res)
	// 		gits[i+1] = gits[time].Shift()
	// 	}
	// }
}

func demo() {
	queue0 := NewPersistentQueue()
	queue1 := queue0.Push(1)
	queue2 := queue1.Push(2)
	queue3 := queue2.Pop()
	fmt.Println(queue1.Top(), queue2.Top(), queue3.Top())
}

func timeit() {
	assert := func(a bool) {
		if !a {
			panic("assert fail")
		}
	}
	queue := NewPersistentQueue()
	assert(queue.IsEmpty() == true)
	queue1 := queue.Push(1)
	assert(queue1.IsEmpty() == false)
	assert(queue1.Top() == 1)
	queue2 := queue1.Push(2)
	assert(queue2.IsEmpty() == false)
	assert(queue2.Top() == 1)
	queue3 := queue2.Pop()
	assert(queue3.IsEmpty() == false)
	assert(queue3.Top() == 2)
	queue4 := queue3.Pop()
	assert(queue4.IsEmpty() == true)

	time1 := time.Now()
	for i := int32(0); i < 1e7; i++ {
		queue = queue.Push(i)
		queue = queue.Pop()
		queue.IsEmpty()
	}

	fmt.Println(time.Since(time1)) // 1.666934792s
}

// 完全可持久化队列.均摊时间复杂度O(1)，最坏时间复杂度O(N).
// https://www.kmonos.net/pub/Presen/PFDS.pdf
// https://37zigen.com/bankers-queue/
type PersistentQueue struct {
	frontSize, rearSize int32
	front, rear         *PersistentStack
}

func NewPersistentQueue() *PersistentQueue {
	return &PersistentQueue{}
}

func _newPersistentQueue(f *PersistentStack, fsize int32, r *PersistentStack, rsize int32) *PersistentQueue {
	return &PersistentQueue{frontSize: fsize, rearSize: rsize, front: f, rear: r}
}

func (q *PersistentQueue) IsEmpty() bool {
	return q.frontSize == 0
}

func (q *PersistentQueue) Top() V {
	return q.front.Top()
}

func (q *PersistentQueue) Push(x V) *PersistentQueue {
	return _newPersistentQueue(q.front, q.frontSize, Push(q.rear, x), q.rearSize+1).normalize()
}

func (q *PersistentQueue) Pop() *PersistentQueue {
	return _newPersistentQueue(q.front.Pop(), q.frontSize-1, q.rear, q.rearSize).normalize()
}

func (q *PersistentQueue) Len() int {
	return int(q.frontSize + q.rearSize)
}

func (q *PersistentQueue) normalize() *PersistentQueue {
	if q.frontSize >= q.rearSize {
		return q
	} else {
		return _newPersistentQueue(Concat(q.front, Reverse(q.rear)), q.frontSize+q.rearSize, nil, 0)
	}
}

type PersistentStack struct {
	next     *PersistentStack
	value    V
	evaluate func() *PersistentStack // 惰性求值
}

func NewPersistentStack() *PersistentStack {
	res := &PersistentStack{}
	res.evaluate = func() *PersistentStack { return res }
	return res
}

func _newPersistentStack(value V, next *PersistentStack) *PersistentStack {
	res := &PersistentStack{value: value, next: next}
	res.evaluate = func() *PersistentStack { return res }
	return res
}

func (s *PersistentStack) IsEmpty() bool {
	return s.next == nil
}

func (s *PersistentStack) Top() V {
	return s.value
}

func (s *PersistentStack) Pop() *PersistentStack {
	if s.next != nil {
		s.next = s.next.evaluate()
	}
	return s.next
}

func Push(x *PersistentStack, v V) *PersistentStack {
	return _newPersistentStack(v, x)
}

func Concat(x, y *PersistentStack) *PersistentStack {
	if x == nil {
		return y.evaluate()
	} else {
		next := &PersistentStack{}
		next.evaluate = func() *PersistentStack { return Concat(x.Pop(), y) }
		return _newPersistentStack(x.value, next)
	}
}

func Reverse(head *PersistentStack) *PersistentStack {
	res := &PersistentStack{}
	res.evaluate = func() *PersistentStack {
		var tmp *PersistentStack
		for x := head; x != nil; x = x.Pop() {
			tmp = Push(tmp, x.Top())
		}
		return tmp
	}
	return res
}
