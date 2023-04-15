//
// Fully Retroactive Queue
//
// Description:
//   It maintains a list of actions ("push" and "pop").
//   We can insert/erase the actions and ask status ("front") in
//   any position of the list.
//
//   To implement this structure, we maintain the list by a BST
//   and keep track the number of pushs/pops in each position
//   by using the range addition and range minimum.
//
// Complexity:
//   O(n log n).
//
// References:
//   E. Demaine, J. Iacono, and S. Langerman (2007):
//   "Retroactive data structures". ACM Transactions of Algorithms,
//   vol.3, no.2, pp.1--21.
//

package main

import (
	"fmt"
	"time"
)

func main() {
	queue := NewRetroactiveQueue()
	time1 := queue.InsertPush(queue.Head, 10)
	time2 := queue.InsertPush(time1, 20)
	time3 := queue.InsertPush(time2, 30)
	time4 := queue.InsertPop(time3)
	fmt.Println(queue.Front(time1)) // [10]
	fmt.Println(queue.Front(time2)) // [10, 20]
	fmt.Println(queue.Front(time3)) // [10, 20, 30]
	fmt.Println(queue.Front(time4)) // [20]
	queue.Erase(time1)              // time1之后, InsertPush10这个操作被删除
	fmt.Println(queue.Front(time1)) // [10]
	fmt.Println(queue.Front(time2)) // [20]
	fmt.Println(queue.Front(time3)) // [20, 30]
	fmt.Println(queue.Front(time4)) // [30]

	t := time.Now()
	// 1e5
	for i := 0; i < 100000; i++ {
		queue.InsertPush(queue.Head, 10)
	}
	fmt.Println(time.Since(t))
}

type T int

// 完全可追溯队列.
type RetroactiveQueue struct {
	Head *Time
}

type Time struct {
	kind                                int // 0:nil, 1:push, 2:pop
	value                               T
	child                               [2]*Time
	parent                              *Time
	aPush, dPush, aPop, dPop, minRemain int
}

func NewRetroactiveQueue() *RetroactiveQueue {
	return &RetroactiveQueue{Head: &Time{}}
}

func (rq *RetroactiveQueue) InsertPush(x *Time, value T) *Time {
	return rq._insert(x, &Time{kind: 1, value: value})
}

func (rq *RetroactiveQueue) InsertPop(x *Time) *Time {
	return rq._insert(x, &Time{kind: 2})
}

func (rq *RetroactiveQueue) Erase(x *Time) *Time {
	rq._splay(x)
	y := x.child[1]
	if y == nil {
		x = x.child[0]
		x.parent = nil
		return x
	}
	if x.kind == 1 {
		y.dPush--
	} else if x.kind == 2 {
		y.dPop--
	}
	y.parent = nil
	rq._update(y)
	for y.child[0] != nil {
		y = y.child[0]
	}
	rq._splay(y)
	y.child[0] = x.child[0]
	if y.child[0] != nil {
		y.child[0].parent = y
	}
	rq._update(y.child[0])
	rq._update(y)
	return y
}

func (rq *RetroactiveQueue) Front(x *Time) T {
	rq._splay(x)
	k := x.aPop + x.dPop
	for y := x; y != nil; {
		rq._pushDown(y)
		if y.aPush > k {
			x = y
			y = y.child[0]
		} else {
			y = y.child[1]
		}
	}
	return x.value
}

func (rq *RetroactiveQueue) _remain(x *Time) int {
	if x == nil {
		return 0
	}
	return x.aPush + x.dPush - x.aPop - x.dPop
}

func (rq *RetroactiveQueue) _update(x *Time) *Time {
	if x == nil {
		return x
	}
	x.minRemain = rq._remain(x)
	for i := 0; i < 2; i++ {
		if x.child[i] != nil {
			x.minRemain = min(x.minRemain, x.child[i].minRemain)
		}
	}
	return x
}

func (rq *RetroactiveQueue) _pushDown(x *Time) *Time {
	if x == nil {
		return x
	}
	for i := 0; i < 2; i++ {
		if x.child[i] != nil {
			x.child[i].dPush += x.dPush
			x.child[i].dPop += x.dPop
		}
	}
	x.aPush += x.dPush
	x.aPop += x.dPop
	x.dPush = 0
	x.dPop = 0
	return x
}

func (rq *RetroactiveQueue) _dir(x *Time) int {
	if x.parent != nil && x.parent.child[1] == x {
		return 1
	}
	return 0
}

func (rq *RetroactiveQueue) _rot(x *Time) {
	d := rq._dir(x)
	p := x.parent
	rq._pushDown(p.parent)
	rq._pushDown(p)
	rq._pushDown(x)
	rq._link(p.parent, x, rq._dir(p))
	rq._link(p, x.child[1^d], d)
	rq._link(x, p, 1^d)
	rq._update(p)
	rq._update(x)
}

func (rq *RetroactiveQueue) _link(x, y *Time, d int) {
	if x != nil {
		x.child[d] = y
	}
	if y != nil {
		y.parent = x
	}
}

func (rq *RetroactiveQueue) _splay(x *Time) {
	if x == nil {
		return
	}
	for x.parent != nil {
		if x.parent.parent != nil {
			if rq._dir(x) == rq._dir(x.parent) {
				rq._rot(x.parent)
			} else {
				rq._rot(x)
			}
		}

		rq._rot(x)
	}
	rq._pushDown(x)
}

func (rq *RetroactiveQueue) _insert(x *Time, y *Time) *Time {
	rq._splay(x)
	y.child[0] = x
	x.parent = y
	y.child[1] = x.child[1]
	x.child[1] = nil
	if y.child[1] != nil {
		y.child[1].parent = y
		if y.kind == 1 {
			y.child[1].dPush += 1
		} else if y.kind == 2 {
			y.child[1].dPop += 1
		}
	}
	y.aPush = x.aPush + x.dPush + (y.kind & 1)
	y.aPop = x.aPop + x.dPop + (y.kind >> 1)
	rq._update(y.child[0])
	rq._update(y.child[1])
	rq._update(y)
	return y
}

func (rq *RetroactiveQueue) _valid(x *Time) bool {
	rq._splay(x)
	return x.minRemain >= 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
