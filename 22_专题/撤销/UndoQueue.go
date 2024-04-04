// https://codeforces.com/contest/1423/problem/H
// 给定n条顶点，m条无向边，之后提供q个请求，第i个请求为li,ri，要求回答仅考虑编号在[li,ri]之间的边，判断所有顶点是否连通。其中1≤n,m≤105,1≤q≤106

// https://codeforces.com/contest/1386/problem/C
// 提供n个商品，第i件商品的价格为wi，价值为vi。
// 我们总共有m单位金钱，希望能买到总价值最大的货物。
// 换言之，我们希望选择一些商品，这些商品的总价格不超过m的前提下保证总价值最大。
// 同时，如果我们同时购买k件物品，那么这k件物品中的一件我们可以免费以五折优化买走（如果价格不能整除2,则上取整），且上述优惠最多只能发生一次。
// 问最大能取走的货物的总价值。其中1≤n,m≤3000，1≤wi,vi≤109，1≤k≤n。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	demo()
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

// int main() {
// 	ios::sync_with_stdio(0), cin.tie(0);
// 	cin >> n >> m >> q;
// 	dsuqueue D(n);
// 	e.resize(m);
// 	for (auto &i : e) {
// 		cin >> i.first >> i.second;
// 		--i.first, --i.second;
// 	}

// 	R.resize(m);
// 	for (int i = 0; i < m; i++)
// 		D.unite(e[i].first, e[i].second);
// 	int nxt = 0;
// 	for (int i = 0; i < m; i++) {
// 		while (!D.is_bipartite() && nxt < m) {
// 			D.undo();
// 			nxt++;
// 		}
// 		if (D.is_bipartite())
// 			R[i] = nxt - 1;
// 		else
// 			R[i] = md;

// 		D.unite(e[i].first, e[i].second);
// 	}

//		while (q--) {
//			int l, r;
//			cin >> l >> r;
//			--l, --r;
//			if (R[l] <= r) cout << "NO\n";
//			else cout << "YES\n";
//		}
//	}
//
// Joker
// https://www.luogu.com.cn/problem/CF1386C
// 给定一张 n 个点 m 条边的图，q 次询问，每次询问删掉 [li,ri] 内的边，问这张图是否存在奇环。
// 存在奇环<=>不是二分图
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
	queries := make([][2]int32, q)
	for i := int32(0); i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
		queries[i][0]--
		queries[i][1]--
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
	dq, bufA, bufB []*FlaggedCommutativeOperation
}

func NewUndoQueue(capacity int) *UndoQueue {
	return &UndoQueue{
		dq:   make([]*FlaggedCommutativeOperation, 0, capacity),
		bufA: make([]*FlaggedCommutativeOperation, 0, capacity),
		bufB: make([]*FlaggedCommutativeOperation, 0, capacity),
	}
}

func (uq *UndoQueue) Append(op *FlaggedCommutativeOperation) {
	op.flag = false
	uq._pushAndDo(op)
}

func (uq *UndoQueue) PopLeft() *FlaggedCommutativeOperation {
	if !uq.dq[len(uq.dq)-1].flag {
		uq._popAndUndo()
		for len(uq.dq) > 0 && len(uq.bufB) != len(uq.bufA) {
			uq._popAndUndo()
		}
		if len(uq.dq) == 0 {
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

	res := uq.dq[len(uq.dq)-1]
	uq.dq = uq.dq[:len(uq.dq)-1]
	res.Undo()
	return res
}

func (uq *UndoQueue) Len() int {
	return len(uq.dq)
}

func (uq *UndoQueue) Empty() bool {
	return len(uq.dq) == 0
}

func (uq *UndoQueue) Clear() {
	n := len(uq.dq)
	for i := 0; i < n; i++ {
		uq._popAndUndo()
	}
	uq.bufA = uq.bufA[:0]
	uq.bufB = uq.bufB[:0]
}

func (uq *UndoQueue) _pushAndDo(op *FlaggedCommutativeOperation) {
	uq.dq = append(uq.dq, op)
	op.Apply()
}

func (uq *UndoQueue) _popAndUndo() {
	res := uq.dq[len(uq.dq)-1]
	uq.dq = uq.dq[:len(uq.dq)-1]
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
