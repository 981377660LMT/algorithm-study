// https://codeforces.com/contest/1423/problem/H

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"strings"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	// demo()
	// CF1386C()
	CF1423H()
}

func demo() {
	uf := NewBipartiteChecker(5)
	getAll := func() []int32 {
		res := make([]int32, 5)
		for i := range res {
			a, _ := uf.Find(int32(i))
			res[i] = a
		}
		return res
	}
	edges := [][2]int32{{0, 1}, {1, 2}, {2, 3}, {3, 4}}
	Q := NewUndoQueue(int32(len(edges)))

	steps := make([]*FlaggedCommutativeOperation, len(edges))
	for i := range edges {
		u, v := edges[i][0], edges[i][1]
		step := NewFlaggedCommutativeOperation(
			func() { uf.Union(u, v) },
			func() { uf.Undo() },
		)
		steps[i] = step
	}

	Q.Append(steps[1])
	Q.Append(steps[2])
	Q.Append(steps[3])

	fmt.Println("before undo", getAll()) // before undo [0 1 1 1 1]
	Q.PopLeft()
	fmt.Println("after undo", getAll()) // after undo [0 1 3 3 3]
	Q.PopLeft()
	fmt.Println("after undo", getAll()) // after undo [0 1 2 3 3]
	Q.PopLeft()
	fmt.Println("after undo", getAll()) // after undo [0 1 2 3 4]

}

// Joker (删去一些边后这张图是否是二分图)
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
	steps := make([]*FlaggedCommutativeOperation, m)
	for i := int32(0); i < m; i++ {
		u, v := edges[i][0], edges[i][1]
		step := NewFlaggedCommutativeOperation(
			func() { uf.Union(u, v) },
			func() { uf.Undo() },
		)
		steps[i] = step
	}
	for _, step := range steps {
		undoQueue.Append(step)
	}

	// maxRight[i]：删除[i,maxRight[i]]内的边，使得图中存在奇环(不是二分图)的最大右端点。
	maxRight := GetMaxRight(
		m,
		func(_ int32) {
			undoQueue.PopLeft()
		},
		func(left int32) {
			undoQueue.Append(steps[left])
		},
		func(left, right int32) bool {
			return !uf.IsBipartite()
		},
	)

	for i := int32(0); i < q; i++ {
		var left, right int32 // [left, right]
		fmt.Fscan(in, &left, &right)
		left--
		right--
		hasOddCycle := maxRight[left] >= right
		if hasOddCycle {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

// Virus (窗口大小为k的滑动窗口撤销)
// https://www.luogu.com.cn/problem/CF1423H
// 给定一个 n 个点的图，有 q 次操作，每一条边在其连边的第 k 天开始时被删去。
// 有三种操作：
// 1 x y 将 x 与 y 连边。
// 2 z 询问 z 所在的连通块大小。
// 3 进入下一天。
func CF1423H() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q, k int32
	fmt.Fscan(in, &n, &q, &k)

	uf := NewUnionFindArrayWithUndo(n)
	undoQueue := NewUndoQueue(q)

	day := int32(0)
	dayCount := int32(0)
	dayQueue := [][2]int32{} // (day,dayCount)
	union := func(x, y int32) {
		undoQueue.Append(NewFlaggedCommutativeOperation(
			func() { uf.Union(x, y) },
			func() { uf.Undo() },
		))
		dayCount++
	}

	query := func(z int32) int32 {
		return uf.GetSize(z)
	}

	elapse := func() {
		dayQueue = append(dayQueue, [2]int32{day, dayCount})
		day++
		dayCount = 0
		for len(dayQueue) > 0 && day-dayQueue[0][0] >= k {
			count := dayQueue[0][1]
			for i := int32(0); i < count; i++ {
				undoQueue.PopLeft()
			}
			dayQueue = dayQueue[1:]
		}
	}

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var x, y int32
			fmt.Fscan(in, &x, &y)
			x--
			y--
			union(x, y)
		} else if op == 2 {
			var z int32
			fmt.Fscan(in, &z)
			z--
			fmt.Fprintln(out, query(z))
		} else {
			elapse()
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

// 对每个固定的左端点`left(0<=left<n)`，找到最大的右端点`maxRight`，
// 使得滑动窗口内的元素满足`predicate(left,maxRight)`成立.
// 如果不存在，`maxRight`为-1.
func GetMaxRight(
	n int32,
	append func(right int32),
	popLeft func(left int32),
	predicate func(left, right int32) bool,
) []int32 {
	maxRight := make([]int32, n)
	right := int32(0)
	visitedRight := make([]bool, n)
	visitRight := func(right int32) {
		if visitedRight[right] {
			return
		}
		visitedRight[right] = true
		append(right)
	}

	for left := int32(0); left < n; left++ {
		if right < left {
			right = left
		}
		for right < n {
			visitRight(right)
			if predicate(left, right) {
				right++
			} else {
				break
			}
		}

		if right == n {
			for i := left; i < n; i++ {
				maxRight[i] = n - 1
			}
			break
		}

		if tmp := right - 1; tmp >= left {
			maxRight[left] = tmp
		} else {
			maxRight[left] = -1
		}
		popLeft(left)
	}

	return maxRight
}

func NewUnionFindArrayWithUndo(n int32) *UnionFindArrayWithUndo {
	data := make([]int32, n)
	for i := range data {
		data[i] = -1
	}
	return &UnionFindArrayWithUndo{Part: n, n: n, data: data}
}

type historyItem struct{ a, b int32 }

type UnionFindArrayWithUndo struct {
	Part      int32
	n         int32
	innerSnap int32
	data      []int32
	history   []*historyItem // (root,data)
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *UnionFindArrayWithUndo) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}
	small, smallData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	big, bigData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	uf.data[small] = smallData
	uf.data[big] = bigData
	if big != small {
		uf.Part++
	}
	return true
}

// 保存并查集当前的状态.
//
//	!Snapshot() 之后可以调用 Rollback(-1) 回滚到这个状态.
func (uf *UnionFindArrayWithUndo) Snapshot() {
	uf.innerSnap = int32(len(uf.history) >> 1)
}

// 回滚到指定的状态.
//
//	state 为 -1 表示回滚到上一次 `SnapShot` 时保存的状态.
//	其他值表示回滚到状态id为state时的状态.
func (uf *UnionFindArrayWithUndo) Rollback(state int32) bool {
	if state == -1 {
		state = uf.innerSnap
	}
	state <<= 1
	if state < 0 || state > int32(len(uf.history)) {
		return false
	}
	for state < int32(len(uf.history)) {
		uf.Undo()
	}
	return true
}

// 获取当前并查集的状态id.
//
//	也就是当前合并(Union)被调用的次数.
func (uf *UnionFindArrayWithUndo) GetState() int {
	return len(uf.history) >> 1
}

func (uf *UnionFindArrayWithUndo) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndo) Union(x, y int32) bool {
	x, y = uf.Find(x), uf.Find(y)
	uf.history = append(uf.history, &historyItem{x, uf.data[x]})
	uf.history = append(uf.history, &historyItem{y, uf.data[y]})
	if x == y {
		return false
	}
	if uf.data[x] > uf.data[y] {
		x ^= y
		y ^= x
		x ^= y
	}
	uf.data[x] += uf.data[y]
	uf.data[y] = x
	uf.Part--
	return true
}

func (uf *UnionFindArrayWithUndo) Find(x int32) int32 {
	cur := x
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}
func (ufa *UnionFindArrayWithUndo) SetPart(part int32) { ufa.Part = part }

func (uf *UnionFindArrayWithUndo) IsConnected(x, y int32) bool { return uf.Find(x) == uf.Find(y) }

func (uf *UnionFindArrayWithUndo) GetSize(x int32) int32 { return -uf.data[uf.Find(x)] }

func (ufa *UnionFindArrayWithUndo) GetGroups() map[int32][]int32 {
	groups := make(map[int32][]int32)
	for i := int32(0); i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArrayWithUndo) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int32, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
