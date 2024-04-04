// https://codeforces.com/contest/1386/problem/C
// 给定n条顶点，m条无向边，之后提供q个请求，第i个请求为li,ri，要求回答仅考虑编号在[li,ri]之间的边，判断所有顶点是否连通。其中1≤n,m≤105,1≤q≤106

// https://codeforces.com/contest/1423/problem/H

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()
	CF1386C()
}

func demo() {
	sum := 0
	Q := NewUndoQueue(10)
	Q.Append(NewFlaggedCommutativeOperation(
		func() {
			sum++
		},
		func() {
			sum--
		},
	))
	fmt.Println(sum)
	Q.PopLeft()
	fmt.Println(sum)
}

// Joker
// https://www.luogu.com.cn/problem/CF1386C
// 给定一张 n 个点 m 条边的图，q 次询问，每次询问删掉 [li,ri] 内的边，问这张图是否存在奇环。
// 存在奇环<=>不是二分图
//
// !注意到对每个固定的左端点 i，向右删除边的过程是单调的，即删除的边越多，越容易使得图为二分图(不存在奇环)。
// !因此可以滑动窗口处理出每个左端点固定时的最大右端点maxRight，然后对每个询问判断是否在这个区间内即可。
func CF1386C() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)
	edges := make([][2]int32, m)
	for i := int32(0); i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}

	uf := NewBipartiteChecker(n)
	undoQueue := NewUndoQueue(m)
	maxRight := make([]int32, m) // 左端点固定为i时，使得图为二分图(不存在奇环)的最大右端点(包含)

	for i := int32(0); i < q; i++ {
		var left, right int32 // [left, right]
		fmt.Fscan(in, &left, &right)
		left--
		right--
		isBipartite := maxRight[left] <= right
		if isBipartite {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
		}
	}
}

// 满足交换律的操作，即(op1, op2) 和 (op2, op1) 作用效果相同.
type FlaggedCommutativeOperation struct {
	apply, undo func()
	flag        bool
}

func NewFlaggedCommutativeOperation(apply, undo func()) *FlaggedCommutativeOperation {
	return &FlaggedCommutativeOperation{apply: apply, undo: undo}
}
func (op *FlaggedCommutativeOperation) Apply() { op.apply() }
func (op *FlaggedCommutativeOperation) Undo()  { op.undo() }

// 支持撤销操作的队列，每个操作会被应用和撤销 O(logn) 次.
type UndoQueue struct {
	stack, bufA, bufB []*FlaggedCommutativeOperation
}

func NewUndoQueue(capacity int32) *UndoQueue {
	return &UndoQueue{
		stack: make([]*FlaggedCommutativeOperation, 0, capacity),
		bufA:  make([]*FlaggedCommutativeOperation, 0, capacity),
		bufB:  make([]*FlaggedCommutativeOperation, 0, capacity),
	}
}

func (uq *UndoQueue) Append(op *FlaggedCommutativeOperation) {
	op.flag = false
	uq._pushAndDo(op)
}

func (uq *UndoQueue) PopLeft() *FlaggedCommutativeOperation {
	if !uq.stack[len(uq.stack)-1].flag {
		uq._popAndUndo()
		for len(uq.stack) > 0 && len(uq.bufB) != len(uq.bufA) {
			uq._popAndUndo()
		}
		if len(uq.stack) == 0 {
			for len(uq.bufB) > 0 {
				res := uq.bufB[0]
				uq.bufB = uq.bufB[1:]
				res.flag = true
				uq._pushAndDo(res)
			}
		} else {
			for len(uq.bufB) > 0 {
				res := uq.bufB[len(uq.bufB)-1]
				uq.bufB = uq.bufB[:len(uq.bufB)-1]
				uq._pushAndDo(res)
			}
		}
		for len(uq.bufA) > 0 {
			uq._pushAndDo(uq.bufA[len(uq.bufA)-1])
			uq.bufA = uq.bufA[:len(uq.bufA)-1]
		}
	}

	res := uq.stack[len(uq.stack)-1]
	uq.stack = uq.stack[:len(uq.stack)-1]
	res.Undo()
	return res
}

func (uq *UndoQueue) Len() int32 {
	return int32(len(uq.stack))
}

func (uq *UndoQueue) Empty() bool {
	return len(uq.stack) == 0
}

func (uq *UndoQueue) Clear() {
	n := len(uq.stack)
	for i := 0; i < n; i++ {
		uq._popAndUndo()
	}
	uq.bufA = uq.bufA[:0]
	uq.bufB = uq.bufB[:0]
}

func (uq *UndoQueue) _pushAndDo(op *FlaggedCommutativeOperation) {
	uq.stack = append(uq.stack, op)
	op.Apply()
}

func (uq *UndoQueue) _popAndUndo() {
	res := uq.stack[len(uq.stack)-1]
	uq.stack = uq.stack[:len(uq.stack)-1]
	res.Undo()
	if res.flag {
		uq.bufA = append(uq.bufA, res)
	} else {
		uq.bufB = append(uq.bufB, res)
	}
}

// 在线二分图检测.
type BipartiteChecker struct {
	n            int32
	parent       []int32
	rank         []int32
	color        []int32
	version      int32
	firstViolate int32
	history      []*setValueStep // plugin
}

func NewBipartiteChecker(n int32) *BipartiteChecker {
	res := &BipartiteChecker{
		n:            n,
		parent:       make([]int32, n),
		rank:         make([]int32, n),
		color:        make([]int32, n),
		firstViolate: -1,
	}
	for i := int32(0); i < n; i++ {
		res.parent[i] = i
	}
	return res
}

func (b *BipartiteChecker) IsBipartite() bool {
	return b.firstViolate == -1
}

// (leader, color)
func (b *BipartiteChecker) Find(x int32) (int32, int32) {
	if x == b.parent[x] {
		return x, 0
	}
	leader, color := b.Find(b.parent[x])
	color ^= b.color[x]
	return leader, color
}

func (b *BipartiteChecker) Union(x, y int32) {
	b.version++
	color := int32(1)
	leaderX, distX := b.Find(x)
	x, color = leaderX, color^distX
	leaderY, distY := b.Find(y)
	y, color = leaderY, color^distY
	if x == y {
		if color == 1 && b.firstViolate == -1 {
			b.firstViolate = b.version
		}
		b.setValue(&b.parent[0], b.parent[0])
		return
	}
	if b.rank[x] < b.rank[y] {
		b.setValue(&b.parent[x], y)
		b.setValue(&b.color[x], color)
	} else {
		b.setValue(&b.parent[y], x)
		b.setValue(&b.color[y], color)
		if b.rank[x] == b.rank[y] {
			b.setValue(&b.rank[x], b.rank[x]+1)
		}
	}
}

func (b *BipartiteChecker) Undo() {
	if len(b.history) == 0 {
		return
	}
	v := b.history[len(b.history)-1].version
	if b.firstViolate == v {
		b.firstViolate = -1
	}
	for len(b.history) > 0 && b.history[len(b.history)-1].version == v {
		b.history[len(b.history)-1].Revert()
		b.history = b.history[:len(b.history)-1]
	}
}

func (b *BipartiteChecker) setValue(cell *int32, newValue int32) {
	step := newSetValueStep(cell, *cell, b.version)
	*cell = newValue // apply
	b.history = append(b.history, step)
}

type setValueStep struct {
	cell     *int32
	oldValue int32
	version  int32
}

func newSetValueStep(cell *int32, oldValue int32, version int32) *setValueStep {
	return &setValueStep{cell: cell, oldValue: oldValue, version: version}
}

// func (s *SetValueStep) Apply()  { *s.cell = s.newValue }
func (s *setValueStep) Revert() { *s.cell = s.oldValue }
