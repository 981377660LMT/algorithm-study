// 有向连通图最小生成树
// 假定从root可以到达所有点
// 如果没有指定根节点,就添加一个虚拟源点连接所有点,权重为0
// Directed MST
// n,m<=2e5
// O(ELog(V))

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, root int
	fmt.Fscan(in, &n, &m, &root)
	edges := make([]Edge, 0, m)
	for i := 0; i < m; i++ {
		var from, to, cost int
		fmt.Fscan(in, &from, &to, &cost)
		edges = append(edges, Edge{from, to, cost, i})
	}

	minCost, eis := directedMST(n, edges, root)
	fmt.Fprintln(out, minCost)

	// 求出最小生成树的每个点的父节点
	parents := make([]int, n)
	parents[root] = root
	for i := 0; i < len(eis); i++ {
		ei := eis[i]
		parents[edges[ei].to] = edges[ei].from
	}

	for v := 0; v < n; v++ {
		fmt.Fprint(out, parents[v], " ")
	}
}

type Edge struct{ from, to, cost, ei int }

// 给定一个连通的有向图，求以root为根节点的最小生成树
//  返回值：最小生成树的权值和，最小生成树的边的编号
func directedMST(n int, edges []Edge, root int) (int, []int) {
	for i := 0; i < n; i++ {
		if i != root {
			edges = append(edges, Edge{i, root, 0, -1})
		}
	}

	x := 0

	par := make([]int, 2*n)
	vis := make([]int, 2*n)
	link := make([]int, 2*n)
	for i := range par {
		par[i] = -1
		vis[i] = -1
		link[i] = -1
	}

	heap := NewSkewHeap(true)
	ins := make([]*SkewHeapNode, 2*n)

	for i := range edges {
		e := edges[i]
		ins[e.to] = heap.Push(ins[e.to], e.cost, i)
	}

	st := []int{}

	go_ := func(x int) int {
		x = edges[ins[x].index].from
		for link[x] != -1 {
			st = append(st, x)
			x = link[x]
		}
		for _, p := range st {
			link[p] = x
		}
		st = st[:0]
		return x
	}

	for i := n; ins[x] != nil; i++ {
		for ; vis[x] == -1; x = go_(x) {
			vis[x] = 0
		}
		for ; x != i; x = go_(x) {
			w := ins[x].key
			v := heap.Pop(ins[x])
			v = heap.Add(v, -w)
			ins[i] = heap.Meld(ins[i], v)
			par[x] = i
			link[x] = i
		}
		for ; ins[x] != nil && go_(x) == x; ins[x] = heap.Pop(ins[x]) {
		}
	}

	cost := 0
	res := []int{}
	for i := root; i != -1; i = par[i] {
		vis[i] = 1
	}
	for i := x; i >= 0; i-- {
		if vis[i] == 1 {
			continue
		}
		cost += edges[ins[i].index].cost
		res = append(res, edges[ins[i].index].ei)
		for j := edges[ins[i].index].to; j != -1 && vis[j] == 0; j = par[j] {
			vis[j] = 1
		}
	}

	return cost, res
}

type E = int

type SkewHeapNode struct {
	key, lazy   E
	left, right *SkewHeapNode
	index       int
}

type SkewHeap struct {
	isMin bool
}

func NewSkewHeap(isMin bool) *SkewHeap {
	return &SkewHeap{isMin: isMin}
}

func (sk *SkewHeap) Push(t *SkewHeapNode, key E, index int) *SkewHeapNode {
	return sk.Meld(t, newNode(key, index))
}

func (sk *SkewHeap) Pop(t *SkewHeapNode) *SkewHeapNode {
	return sk.Meld(t.left, t.right)
}

func (sk *SkewHeap) Top(t *SkewHeapNode) E {
	return t.key
}

func (sk *SkewHeap) Meld(x, y *SkewHeapNode) *SkewHeapNode {
	sk.propagate(x)
	sk.propagate(y)
	if x == nil {
		return y
	}
	if y == nil {
		return x
	}
	if (x.key < y.key) != sk.isMin {
		x, y = y, x
	}
	x.right = sk.Meld(y, x.right)
	x.left, x.right = x.right, x.left
	return x
}

func (sk *SkewHeap) Add(t *SkewHeapNode, lazy E) *SkewHeapNode {
	if t == nil {
		return t
	}
	t.lazy += lazy
	sk.propagate(t)
	return t
}

func (sk *SkewHeap) propagate(t *SkewHeapNode) *SkewHeapNode {
	if t != nil && t.lazy != 0 {
		if t.left != nil {
			t.left.lazy += t.lazy
		}
		if t.right != nil {
			t.right.lazy += t.lazy
		}
		t.key += t.lazy
		t.lazy = 0
	}
	return t
}

func newNode(key E, index int) *SkewHeapNode {
	return &SkewHeapNode{key: key, index: index}
}
