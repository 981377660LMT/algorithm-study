// 动态全局最小割(星图最小割)
// https://ei1333.github.io/library/graph/flow/global-minimum-cut-of-dynamic-star-augmented-graph.hpp
// 给定一张n个点的无向图带权图.
// 添加一个虚拟节点n，对于每个节点i，添加一条边i<->n，权值为weights[i].
// 之后进行q次操作，每次修改i<->n的权值为w，求全局最小割的权值.
//
// - NewGlobalMinimumCutofDynamicStarAugmentedGraph(n, edges) 构造函数，edges为边集.
//   时间复杂度O(VE+V^2logV).
// - Update(v, cost) 修改v<->n的权值为cost，返回全局最小割的权值.
//   时间复杂度O(log^2V).

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// https://judge.yosupo.jp/problem/global_minimum_cut_of_dynamic_star_augmented_graph
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)
	weights := make([]int32, n)
	for i := range weights {
		fmt.Fscan(in, &weights[i])
	}
	edges := make([]Edge, 0, m)
	for i := int32(0); i < m; i++ {
		var u, v int32
		var w int
		fmt.Fscan(in, &u, &v, &w)
		edges = append(edges, Edge{from: u, to: v, cost: w})
	}

	G := NewGlobalMinimumCutofDynamicStarAugmentedGraph(n, edges)
	for i := int32(0); i < n; i++ {
		G.Update(i, int(weights[i]))
	}

	for i := int32(0); i < q; i++ {
		var u int32
		var w int
		fmt.Fscan(in, &u, &w)
		fmt.Fprintln(out, G.Update(u, w))
	}
}

// 无向边
type Edge = struct {
	from, to int32
	cost     int
}

type Neighbor = struct {
	to   int32
	cost int
}

type GlobalMinimumCutofDynamicStarAugmentedGraph struct {
	n   int32
	hld *HldSpecified
	cur []int
	seg *LazySegTree32
}

func NewGlobalMinimumCutofDynamicStarAugmentedGraph(n int32, edges []Edge) *GlobalMinimumCutofDynamicStarAugmentedGraph {
	res := &GlobalMinimumCutofDynamicStarAugmentedGraph{n: n}
	hldGraph := ExtremeVertexSet(n, edges)
	res.hld = NewHldSpecified(hldGraph)
	res.cur = make([]int, n)
	res.hld.Build(int32(len(hldGraph)) - 1)
	vs := make([]int, 2*n-1)
	ins := res.hld.in
	for i := int32(0); i < 2*n-1; i++ {
		for _, e := range hldGraph[i] {
			vs[ins[e.to]] = e.cost
		}
	}
	res.seg = NewLazySegTree32(2*n-1, func(i int32) int { return vs[i] })
	return res
}

func (g *GlobalMinimumCutofDynamicStarAugmentedGraph) Update(v int32, cost int) int {
	g.hld.Add(v, int32(len(g.hld.g))-1, func(l, r int32) { g.seg.Update(l, r, cost-g.cur[v]) }, false)
	g.cur[v] = cost
	return g.seg.QueryAll()
}

// O(VE+V^2logV)
func ExtremeVertexSet(n int32, edges []Edge) [][]Neighbor {
	res := make([][]Neighbor, 2*n-1)
	uf := make([]int32, n)
	for i := int32(0); i < n; i++ {
		uf[i] = i
	}
	cur := make([]int, 2*n-1)
	leaf := make([]bool, 2*n-1)
	for i := int32(0); i < n; i++ {
		leaf[i] = true
	}

	type pair = struct {
		first  int
		second int32
	}
	pq := NewHeap[pair](func(a, b pair) bool { return a.first < b.first })
	for phase := int32(0); phase < n-1; phase++ {
		g := make([][]Neighbor, 2*n-1)
		cost := make([]int, 2*n-1)
		for _, e := range edges {
			e.from = uf[e.from]
			e.to = uf[e.to]
			if e.from != e.to {
				cost[e.from] += e.cost
				cost[e.to] += e.cost
				g[e.from] = append(g[e.from], Neighbor{to: e.to, cost: e.cost})
				g[e.to] = append(g[e.to], Neighbor{to: e.from, cost: e.cost})
			}
		}
		for i := int32(0); i < 2*n-1; i++ {
			if leaf[i] {
				cur[i] = cost[i]
				pq.Push(pair{first: cost[i], second: i})
			}
		}
		x, y := int32(-1), int32(-1)
		for pq.Len() > 0 {
			v := pq.Pop().second
			if cur[v] == -1 {
				continue
			}
			cur[v] = -1
			y = x
			x = v
			for _, e := range g[v] {
				if cur[e.to] != -1 {
					cur[e.to] -= e.cost
					pq.Push(pair{first: cur[e.to], second: e.to})
				}
			}
		}
		z := n + phase
		res[z] = append(res[z], Neighbor{to: x, cost: cost[x]})
		res[z] = append(res[z], Neighbor{to: y, cost: cost[y]})
		for i := int32(0); i < n; i++ {
			if uf[i] == x || uf[i] == y {
				uf[i] = z
			}
		}
		leaf[x] = false
		leaf[y] = false
		leaf[z] = true
	}
	return res
}

func NewHeap[H any](less func(a, b H) bool, nums ...H) *Heap[H] {
	res := &Heap[H]{less: less, data: append(nums[:0:0], nums...)}
	if len(nums) > 1 {
		res.heapify()
	}
	return res
}

func NewHeapWithCapacity[H any](less func(a, b H) bool, capacity int32, nums ...H) *Heap[H] {
	if n := int32(len(nums)); capacity < n {
		capacity = n
	}
	res := &Heap[H]{less: less, data: make([]H, 0, capacity)}
	res.data = append(res.data, nums...)
	if len(nums) > 1 {
		res.heapify()
	}
	return res
}

type Heap[H any] struct {
	less func(a, b H) bool
	data []H
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.up(h.Len() - 1)
}

func (h *Heap[H]) Pop() H {
	n := h.Len() - 1
	h.data[0], h.data[n] = h.data[n], h.data[0]
	h.down(0, n)
	res := h.data[n]
	h.data = h.data[:n]
	return res
}

func (h *Heap[H]) Top() H {
	return h.data[0]
}

// replace 弹出并返回堆顶，同时将 v 入堆.
// 需保证 h 非空.
func (h *Heap[H]) Replace(v H) H {
	top := h.Top()
	h.data[0] = v
	h.fix(0)
	return top
}

// pushPop 先将 v 入堆，然后弹出并返回堆顶.
func (h *Heap[H]) PushPop(v H) H {
	data, less := h.data, h.less
	if len(data) > 0 && less(data[0], v) {
		v, data[0] = data[0], v
		h.fix(0)
	}
	return v
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.down(i, n)
	}
}

func (h *Heap[H]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(h.data[j], h.data[i]) {
			break
		}
		h.data[i], h.data[j] = h.data[j], h.data[i]
		j = i
	}
}

func (h *Heap[H]) down(i0, n int) bool {
	i := i0
	for {
		j1 := (i << 1) | 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(h.data[j2], h.data[j1]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.less(h.data[j], h.data[i]) {
			break
		}
		h.data[i], h.data[j] = h.data[j], h.data[i]
		i = j
	}
	return i > i0
}

func (h *Heap[H]) fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

const INF int = 1e18

// RangeAddRangeMin

type E = int
type Id = int

func (*LazySegTree32) e() E                   { return INF }
func (*LazySegTree32) id() Id                 { return 0 }
func (*LazySegTree32) op(left, right E) E     { return min(left, right) }
func (*LazySegTree32) mapping(f Id, g E) E    { return f + g }
func (*LazySegTree32) composition(f, g Id) Id { return f + g }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// !template
type LazySegTree32 struct {
	n    int32
	size int32
	log  int32
	data []E
	lazy []Id
}

func NewLazySegTree32(n int32, f func(int32) E) *LazySegTree32 {
	tree := &LazySegTree32{}
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32) Query(left, right int32) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *LazySegTree32) QueryAll() E {
	return tree.data[1]
}
func (tree *LazySegTree32) GetAll() []E {
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32) Update(left, right int32, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 单点查询(不需要 pushUp/op 操作时使用)
func (tree *LazySegTree32) Get(index int32) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *LazySegTree32) Set(index int32, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree32) pushUp(root int32) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *LazySegTree32) pushDown(root int32) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *LazySegTree32) propagate(root int32, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

type HldSpecified struct {
	g    [][]Neighbor
	size []int32
	in   []int32
	head []int32
	rev  []int32
	par  []int32
}

func NewHldSpecified(g [][]Neighbor) *HldSpecified {
	return &HldSpecified{g: g}
}

func (hld *HldSpecified) Build(root int32) {
	n := int32(len(hld.g))
	hld.size = make([]int32, n)
	hld.in = make([]int32, n)
	hld.head = make([]int32, n)
	hld.rev = make([]int32, n)
	hld.par = make([]int32, n)
	hld.dfsSize(root, -1)
	hld.head[root] = root
	t := int32(0)
	hld.dfsHld(root, -1, &t)
}

func (hld *HldSpecified) Add(u, v int32, f func(l, r int32), edge bool) {
	in := hld.in
	head, par := hld.head, hld.par
	for {
		if in[u] > in[v] {
			u, v = v, u
		}
		if head[u] == head[v] {
			break
		}
		f(in[head[v]], in[v]+1)
		v = par[head[v]]
	}
	if edge {
		f(in[u]+1, in[v]+1)
	} else {
		f(in[u], in[v]+1)
	}
}

func (hld *HldSpecified) dfsSize(idx, p int32) {
	hld.par[idx] = p
	hld.size[idx] = 1
	nexts := hld.g[idx]
	if len(nexts) > 0 && nexts[0].to == p {
		nexts[0], nexts[len(nexts)-1] = nexts[len(nexts)-1], hld.g[idx][0]
	}
	for i, e := range nexts {
		if e.to == p {
			continue
		}
		hld.dfsSize(e.to, idx)
		hld.size[idx] += hld.size[e.to]
		if hld.size[nexts[0].to] < hld.size[e.to] {
			nexts[0], nexts[i] = nexts[i], nexts[0]
		}
	}
}

func (hld *HldSpecified) dfsHld(idx, p int32, times *int32) {
	hld.in[idx] = *times
	*times++
	hld.rev[hld.in[idx]] = idx
	nexts := hld.g[idx]
	for _, e := range nexts {
		if e.to == p {
			continue
		}
		if nexts[0].to == e.to {
			hld.head[e.to] = hld.head[idx]
		} else {
			hld.head[e.to] = e.to
		}
		hld.dfsHld(e.to, idx, times)
	}
}
