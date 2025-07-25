// 路径查询
// 子树查询
// MaxPath

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	// demo()
	// n = 2, edges = [[0,1,7]], queries = [[1,0],[0,1]]
	fmt.Println(findMedian(2, [][]int{{0, 1, 2}}, [][]int{{1, 1}}))                               // [1]
	fmt.Println(findMedian(4, [][]int{{0, 1, 3}, {1, 2, 3}, {0, 3, 5}}, [][]int{{0, 0}, {3, 2}})) // [0,3]
}

func demo() {
	{
		//    0
		//   / \
		//  1   2
		//     / \
		//    3   4

		tree := NewTree32(5)
		tree.AddEdge(0, 1, 0)
		tree.AddEdge(0, 2, 0)
		tree.AddEdge(2, 3, 0)
		tree.AddEdge(2, 4, 0)
		tree.Build(0)

		S := NewTreeMonoid32Static(tree, false)
		S.Build(func(vidOrEid int32) E { return int(vidOrEid) })
		fmt.Println(S.QuerySubtree(0))          // 10
		fmt.Println(S.QuerySubtree(1))          // 1
		fmt.Println(S.QuerySubtree(2))          // 9
		fmt.Println(S.QuerySubtree(3))          // 3
		fmt.Println(S.QuerySubtree(4))          // 4
		fmt.Println(S.QueryPath(1, 3))          // 6
		fmt.Println(S.QuerySubtreeRooted(0, 3)) // 1

		fmt.Println(S.QuerySubtree(0))          // 10
		fmt.Println(S.QuerySubtreeRooted(4, 3)) // 1
		fmt.Println(S.QuerySubtreeRooted(2, 3)) // 7

		fmt.Println(S.MaxPath(1, 3, func(x E) bool { return x < -1 })) // 2
	}
}

// 3585. 树中找到带权中位节点
// https://leetcode.cn/problems/find-weighted-median-node-in-tree/description/
// 给你一个整数 n，以及一棵 无向带权 树，根节点为节点 0，树中共有 n 个节点，编号从 0 到 n - 1。
// 该树由一个长度为 n - 1 的二维数组 edges 表示，其中 edges[i] = [ui, vi, wi] 表示存在一条从节点 ui 到 vi 的边，权重为 wi。
// 带权中位节点 定义为从 ui 到 vi 路径上的 第一个 节点 x，使得从 ui 到 x 的边权之和 大于等于 该路径总权值和的一半。
// 给你一个二维整数数组 queries。对于每个 queries[j] = [uj, vj]，求出从 uj 到 vj 路径上的带权中位节点。
// 返回一个数组 ans，其中 ans[j] 表示查询 queries[j] 的带权中位节点编号。
func findMedian(n int, edges [][]int, queries [][]int) []int {
	tree := NewTree32(int32(n))
	for _, e := range edges {
		u, v, w := int32(e[0]), int32(e[1]), e[2]
		tree.AddEdge(u, v, w)
	}
	tree.Build(0)

	tm := NewTreeMonoid32Static(tree, true)
	tm.Build(func(vidOrEid int32) E { return edges[vidOrEid][2] })

	query := func(u, v int32) int32 {
		if u == v {
			return u
		}
		pathSum := tm.QueryPath(u, v)
		target := (pathSum + 1) / 2
		check := func(x E) bool { return x < target }
		bound := tm.MaxPath(u, v, check)
		return tree.Jump(bound, v, 1)
	}

	res := make([]int, len(queries))
	for i, q := range queries {
		u, v := int32(q[0]), int32(q[1])
		res[i] = int(query(u, v))
	}
	return res
}

type E = int

const commutative bool = true // E是否可交换，即 op(a,b) == op(b,a)

func e() E        { return 0 }
func op(a, b E) E { return a + b }

type TreeMonoid32Static struct {
	edge      int32
	n         int32
	tree      *Tree32
	seg, segR *DisjointSparseTableFast[E]
}

func NewTreeMonoid32Static(tree *Tree32, edge bool) *TreeMonoid32Static {
	var edgeValue int32
	if edge {
		edgeValue = 1
	}
	return &TreeMonoid32Static{edge: edgeValue, n: tree.n, tree: tree}
}

func (tag *TreeMonoid32Static) Build(f func(vidOrEid int32) E) {
	idToNode := tag.tree.IdToNode
	vToE := tag.tree.vToE
	if tag.edge == 0 {
		fv := func(i int32) E { return f(idToNode[i]) }
		tag.seg = NewDisjointSparseTableFast(tag.n, fv, e, op)
		if !commutative {
			tag.segR = NewDisjointSparseTableFast(tag.n, fv, e, func(a, b E) E { return op(b, a) })
		}
	} else {
		fe := func(i int32) E {
			if i == 0 {
				return e()
			}
			return f(vToE[idToNode[i]])
		}
		tag.seg = NewDisjointSparseTableFast(tag.n, fe, e, op)
		if !commutative {
			tag.segR = NewDisjointSparseTableFast(tag.n, fe, e, func(a, b E) E { return op(b, a) })
		}
	}
}

func (tag *TreeMonoid32Static) QueryPath(from, to int32) E {
	pd := tag.tree.GetPathDecomposition(from, to, tag.edge)
	res := e()
	for i := 0; i < len(pd); i++ {
		res = op(res, tag.getProd(pd[i][0], pd[i][1]))
	}
	return res
}

func (tag *TreeMonoid32Static) QueryAll() E {
	return tag.QuerySubtree(tag.tree.IdToNode[0])
}

func (tag *TreeMonoid32Static) QuerySubtree(u int32) E {
	return tag.QuerySubtreeRooted(u, -1)
}

func (tag *TreeMonoid32Static) QuerySubtreeRooted(u, root int32) E {
	if root == u {
		return tag.QueryAll()
	}
	if root == -1 || tag.tree.InSubtree(u, root) {
		l, r := tag.tree.Lid[u], tag.tree.Rid[u]
		return tag.seg.Query(l+tag.edge, r)
	}
	if tag.edge == 1 {
		panic("not implemented")
	}
	u = tag.tree.Jump(u, root, 1)
	l, r := tag.tree.Lid[u], tag.tree.Rid[u]
	return op(tag.seg.Query(0, l), tag.seg.Query(r, tag.n))
}

func (tag *TreeMonoid32Static) MaxPath(from, to int32, check func(E) bool) int32 {
	if tag.edge != 0 {
		return tag.maxPathEdge(from, to, check)
	}
	if !check(tag.QueryPath(from, from)) {
		return -1
	}
	pd := tag.tree.GetPathDecomposition(from, to, tag.edge)
	val := e()
	idToNode := tag.tree.IdToNode
	for _, e := range pd {
		x := tag.getProd(e[0], e[1])
		if tmp := op(val, x); check(tmp) {
			val = tmp
			from = idToNode[e[1]]
			continue
		}
		checkTmp := func(x E) bool { return check(op(val, x)) }
		if e[0] <= e[1] {
			i := tag.seg.MaxRight(e[0], checkTmp)
			if i == e[0] {
				return from
			}
			return idToNode[i-1]
		} else {
			i := int32(0)
			if commutative {
				i = tag.seg.MinLeft(e[0]+1, checkTmp)
			} else {
				i = tag.segR.MinLeft(e[0]+1, checkTmp)
			}
			if i == e[0]+1 {
				return from
			}
			return idToNode[i]
		}
	}
	return to
}

func (tag *TreeMonoid32Static) maxPathEdge(from, to int32, check func(E) bool) int32 {
	// if !check(e()) {
	// 	return -1
	// }
	lca := tag.tree.Lca(from, to)
	pd := tag.tree.GetPathDecomposition(from, lca, tag.edge)
	val := e()
	parent, idToNode := tag.tree.Parent, tag.tree.IdToNode
	for _, e := range pd {
		x := tag.getProd(e[0], e[1])
		if tmp := op(val, x); check(tmp) {
			val = tmp
			from = parent[idToNode[e[1]]]
			continue
		}
		checkTmp := func(x E) bool { return check(op(val, x)) }
		i := int32(0)
		if commutative {
			i = tag.seg.MinLeft(e[0]+1, checkTmp)
		} else {
			i = tag.segR.MinLeft(e[0]+1, checkTmp)
		}
		if i == e[0]+1 {
			return from
		}
		return parent[idToNode[i]]
	}
	pd = tag.tree.GetPathDecomposition(lca, to, tag.edge)
	for _, e := range pd {
		x := tag.getProd(e[0], e[1])
		if tmp := op(val, x); check(tmp) {
			val = tmp
			from = idToNode[e[1]]
			continue
		}
		checkTmp := func(x E) bool { return check(op(val, x)) }
		i := tag.seg.MaxRight(e[0], checkTmp)
		if i == e[0] {
			return from
		}
		return idToNode[i-1]
	}
	return to
}

func (tag *TreeMonoid32Static) getProd(a, b int32) E {
	if commutative {
		if a <= b {
			return tag.seg.Query(a, b+1)
		}
		return tag.seg.Query(b, a+1)
	} else {
		if a <= b {
			return tag.seg.Query(a, b+1)
		}
		return tag.segR.Query(b, a+1)
	}
}

// Static RMQ, O(n)预处理, O(1)查询.
type DisjointSparseTableFast[E any] struct {
	n        int32
	leaves   []E
	pre, suf []E
	st       *disjointSparseTable[E]
	e        func() E
	op       func(E, E) E
}

func NewDisjointSparseTableFast[E comparable](n int32, f func(int32) E, e func() E, op func(E, E) E) *DisjointSparseTableFast[E] {
	res := &DisjointSparseTableFast[E]{}
	bNum := n >> 4
	leaves := make([]E, 0, n)
	for i := int32(0); i < n; i++ {
		leaves = append(leaves, f(i))
	}
	pre, suf := append(leaves[:0:0], leaves...), append(leaves[:0:0], leaves...)
	for i := int32(1); i < n; i++ {
		if i&15 != 0 {
			pre[i] = op(pre[i-1], leaves[i])
		}
	}
	for i := n - 1; i > 0; i-- {
		if i&15 != 0 {
			suf[i-1] = op(leaves[i-1], suf[i])
		}
	}
	st := newDisjointSparse(bNum, func(i int32) E { return suf[i<<4] }, e, op)
	res.n = n
	res.leaves = leaves
	res.pre, res.suf = pre, suf
	res.st = st
	res.e = e
	res.op = op
	return res
}

func NewDisjointSparseTableFastFrom[E comparable](leaves []E, e func() E, op func(E, E) E) *DisjointSparseTableFast[E] {
	return NewDisjointSparseTableFast(int32(len(leaves)), func(i int32) E { return leaves[i] }, e, op)
}

func (st *DisjointSparseTableFast[E]) Query(start, end int32) E {
	if start >= end {
		return st.e()
	}
	end--
	a, b := start>>4, end>>4
	if a < b {
		x := st.st.Query(a+1, b)
		x = st.op(st.suf[start], x)
		x = st.op(x, st.pre[end])
		return x
	}
	x := st.leaves[start]
	for i := start + 1; i <= end; i++ {
		x = st.op(x, st.leaves[i])
	}
	return x
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *DisjointSparseTableFast[E]) MaxRight(left int32, check func(e E) bool) int32 {
	if left == ds.n {
		return ds.n
	}
	ok, ng := left, ds.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(ds.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (ds *DisjointSparseTableFast[E]) MinLeft(right int32, check func(e E) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(ds.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

type disjointSparseTable[E any] struct {
	n    int32
	e    func() E
	op   func(E, E) E
	data [][]E
}

// DisjointSparseTable 支持幺半群的区间静态查询.
//
//	eg: 区间乘积取模/区间仿射变换...
func newDisjointSparse[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *disjointSparseTable[E] {
	res := &disjointSparseTable[E]{}
	log := int32(1)
	for (1 << log) < n {
		log++
	}
	data := make([][]E, log)
	data[0] = make([]E, 0, n)
	for i := int32(0); i < n; i++ {
		data[0] = append(data[0], f(i))
	}
	for i := int32(1); i < log; i++ {
		data[i] = append(data[i], data[0]...)
		tmp := data[i]
		b := int32(1 << i)
		for m := b; m <= n; m += 2 * b {
			l, r := m-b, min32(m+b, n)
			for j := m - 1; j >= l+1; j-- {
				tmp[j-1] = op(tmp[j-1], tmp[j])
			}
			for j := m; j < r-1; j++ {
				tmp[j+1] = op(tmp[j], tmp[j+1])
			}
		}
	}
	res.n = n
	res.e = e
	res.op = op
	res.data = data
	return res
}

func (ds *disjointSparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return ds.e()
	}
	end--
	if start == end {
		return ds.data[0][start]
	}
	lca := bits.Len32(uint32(start^end)) - 1
	return ds.op(ds.data[lca][start], ds.data[lca][end])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *disjointSparseTable[E]) MaxRight(left int32, check func(e E) bool) int32 {
	if left == ds.n {
		return ds.n
	}
	ok, ng := left, ds.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(ds.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (ds *disjointSparseTable[E]) MinLeft(right int32, check func(e E) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(ds.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func topbit(x int32) int32 {
	if x == 0 {
		return -1
	}
	return int32(bits.Len32(uint32(x)) - 1)
}

func lowbit(x int32) int32 {
	if x == 0 {
		return -1
	}
	return int32(bits.TrailingZeros32(uint32(x)))
}

type neighbor = struct {
	to   int32
	eid  int32
	cost int
}

type Tree32 struct {
	Lid, Rid      []int32
	IdToNode      []int32
	Depth         []int32
	DepthWeighted []int
	Parent        []int32
	Head          []int32 // 重链头
	Tree          [][]neighbor
	Edges         [][2]int32
	vToE          []int32 // 节点v的父边的id
	n             int32
}

func NewTree32(n int32) *Tree32 {
	res := &Tree32{Tree: make([][]neighbor, n), Edges: make([][2]int32, 0, n-1), n: n}
	return res
}

func (t *Tree32) AddEdge(u, v int32, w int) {
	eid := int32(len(t.Edges))
	t.Tree[u] = append(t.Tree[u], neighbor{to: v, eid: eid, cost: w})
	t.Tree[v] = append(t.Tree[v], neighbor{to: u, eid: eid, cost: w})
	t.Edges = append(t.Edges, [2]int32{u, v})
}

func (t *Tree32) AddDirectedEdge(from, to int32, cost int) {
	eid := int32(len(t.Edges))
	t.Tree[from] = append(t.Tree[from], neighbor{to: to, eid: eid, cost: cost})
	t.Edges = append(t.Edges, [2]int32{from, to})
}

func (t *Tree32) Build(root int32) {
	if root != -1 && int32(len(t.Edges)) != t.n-1 {
		panic("edges count != n-1")
	}
	n := t.n
	t.Lid = make([]int32, n)
	t.Rid = make([]int32, n)
	t.IdToNode = make([]int32, n)
	t.Depth = make([]int32, n)
	t.DepthWeighted = make([]int, n)
	t.Parent = make([]int32, n)
	t.Head = make([]int32, n)
	t.vToE = make([]int32, n)
	for i := int32(0); i < n; i++ {
		t.Depth[i] = -1
		t.Head[i] = root
		t.vToE[i] = -1
	}
	if root != -1 {
		t._dfsSize(root, -1)
		time := int32(0)
		t._dfsHld(root, &time)
	} else {
		time := int32(0)
		for i := int32(0); i < n; i++ {
			if t.Depth[i] == -1 {
				t._dfsSize(i, -1)
				t._dfsHld(i, &time)
			}
		}
	}
}

// 从v开始沿着重链向下收集节点.
func (t *Tree32) HeavyPathAt(v int32) []int32 {
	path := []int32{v}
	for {
		a := path[len(path)-1]
		for _, e := range t.Tree[a] {
			if e.to != t.Parent[a] && t.Head[e.to] == v {
				path = append(path, e.to)
				break
			}
		}
		if path[len(path)-1] == a {
			break
		}
	}
	return path
}

// 返回重儿子，如果没有返回 -1.
func (t *Tree32) HeavyChild(v int32) int32 {
	k := t.Lid[v] + 1
	if k == t.n {
		return -1
	}
	w := t.IdToNode[k]
	if t.Parent[w] == v {
		return w
	}
	return -1
}

// 从v开始向上走k步.
func (t *Tree32) KthAncestor(v, k int32) int32 {
	if k > t.Depth[v] {
		return -1
	}
	for {
		u := t.Head[v]
		if t.Lid[v]-k >= t.Lid[u] {
			return t.IdToNode[t.Lid[v]-k]
		}
		k -= t.Lid[v] - t.Lid[u] + 1
		v = t.Parent[u]
	}
}

func (t *Tree32) Lca(u, v int32) int32 {
	for {
		if t.Lid[u] > t.Lid[v] {
			u, v = v, u
		}
		if t.Head[u] == t.Head[v] {
			return u
		}
		v = t.Parent[t.Head[v]]
	}
}

func (t *Tree32) LcaRooted(u, v, root int32) int32 {
	return t.Lca(u, v) ^ t.Lca(u, root) ^ t.Lca(v, root)
}

func (t *Tree32) Dist(a, b int32) int32 {
	c := t.Lca(a, b)
	return t.Depth[a] + t.Depth[b] - 2*t.Depth[c]
}

func (t *Tree32) DistWeighted(a, b int32) int {
	c := t.Lca(a, b)
	return t.DepthWeighted[a] + t.DepthWeighted[b] - 2*t.DepthWeighted[c]
}

// c 是否在 p 的子树中.c和p不能相等.
func (t *Tree32) InSubtree(c, p int32) bool {
	return t.Lid[p] <= t.Lid[c] && t.Lid[c] < t.Rid[p]
}

// 从 a 开始走 k 步到 b.
func (t *Tree32) Jump(a, b, k int32) int32 {
	if k == 1 {
		if a == b {
			return -1
		}
		if t.InSubtree(b, a) {
			return t.KthAncestor(b, t.Depth[b]-t.Depth[a]-1)
		}
		return t.Parent[a]
	}
	c := t.Lca(a, b)
	dac := t.Depth[a] - t.Depth[c]
	dbc := t.Depth[b] - t.Depth[c]
	if k > dac+dbc {
		return -1
	}
	if k <= dac {
		return t.KthAncestor(a, k)
	}
	return t.KthAncestor(b, dac+dbc-k)
}

func (t *Tree32) SubtreeSize(v int32) int32 {
	return t.Rid[v] - t.Lid[v]
}

func (t *Tree32) SubtreeSizeRooted(v, root int32) int32 {
	if v == root {
		return t.n
	}
	x := t.Jump(v, root, 1)
	if t.InSubtree(v, x) {
		return t.Rid[v] - t.Lid[v]
	}
	return t.n - t.Rid[x] + t.Lid[x]
}

func (t *Tree32) CollectChild(v int32) []int32 {
	var res []int32
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			res = append(res, e.to)
		}
	}
	return res
}

// 收集与 v 相邻的轻边.
func (t *Tree32) CollectLight(v int32) []int32 {
	var res []int32
	skip := true
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			if !skip {
				res = append(res, e.to)
			}
			skip = false
		}
	}
	return res
}

func (tree *Tree32) RestorePath(from, to int32) []int32 {
	res := []int32{}
	composition := tree.GetPathDecomposition(from, to, 0)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.IdToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.IdToNode[i])
			}
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree32) GetPathDecomposition(u, v int32, edge int32) [][2]int32 {
	up, down := [][2]int32{}, [][2]int32{}
	lid, head, parent := tree.Lid, tree.Head, tree.Parent
	for {
		if head[u] == head[v] {
			break
		}
		if lid[u] < lid[v] {
			down = append(down, [2]int32{lid[head[v]], lid[v]})
			v = parent[head[v]]
		} else {
			up = append(up, [2]int32{lid[u], lid[head[u]]})
			u = parent[head[u]]
		}
	}
	if lid[u] < lid[v] {
		down = append(down, [2]int32{lid[u] + edge, lid[v]})
	} else if lid[v]+edge <= lid[u] {
		up = append(up, [2]int32{lid[u], lid[v] + edge})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, edge int32, f func(start, end int32)) {
	head, lid, parent := tree.Head, tree.Lid, tree.Parent
	for {
		if head[u] == head[v] {
			break
		}
		if lid[u] < lid[v] {
			a, b := lid[head[v]], lid[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = parent[head[v]]
		} else {
			a, b := lid[u], lid[head[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = parent[head[u]]
		}
	}
	if lid[u] < lid[v] {
		a, b := lid[u]+edge, lid[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if lid[v]+edge <= lid[u] {
		a, b := lid[u], lid[v]+edge
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree32) Id(root int32) (int32, int32) {
	return tree.Lid[root], tree.Rid[root]
}

// 返回返回边 u-v 对应的 边id.
func (tree *Tree32) Eid(u, v int32) int32 {
	if tree.Parent[u] != v {
		u, v = v, u
	}
	return tree.vToE[u]
}

// 点v对应的父边的边id.如果v是根节点则返回-1.
func (tre *Tree32) VToE(v int32) int32 {
	return tre.vToE[v]
}

// 第i条边对应的深度更深的那个节点.
func (tree *Tree32) EToV(i int32) int32 {
	u, v := tree.Edges[i][0], tree.Edges[i][1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (tree *Tree32) ELid(u int32) int32 {
	return 2*tree.Lid[u] - tree.Depth[u]
}

func (tree *Tree32) ERid(u int32) int32 {
	return 2*tree.Rid[u] - tree.Depth[u] - 1
}

func (t *Tree32) _dfsSize(cur, pre int32) {
	size := t.Rid
	t.Parent[cur] = pre
	if pre != -1 {
		t.Depth[cur] = t.Depth[pre] + 1
	} else {
		t.Depth[cur] = 0
	}
	size[cur] = 1
	nexts := t.Tree[cur]
	for i := int32(len(nexts)) - 2; i >= 0; i-- {
		e := nexts[i+1]
		if t.Depth[e.to] == -1 {
			nexts[i], nexts[i+1] = nexts[i+1], nexts[i]
		}
	}
	hldSize := int32(0)
	for i, e := range nexts {
		to := e.to
		if t.Depth[to] == -1 {
			t.DepthWeighted[to] = t.DepthWeighted[cur] + e.cost
			t.vToE[to] = e.eid
			t._dfsSize(to, cur)
			size[cur] += size[to]
			if size[to] > hldSize {
				hldSize = size[to]
				if i != 0 {
					nexts[0], nexts[i] = nexts[i], nexts[0]
				}
			}
		}
	}
}

func (t *Tree32) _dfsHld(cur int32, times *int32) {
	t.Lid[cur] = *times
	*times++
	t.Rid[cur] += t.Lid[cur]
	t.IdToNode[t.Lid[cur]] = cur
	heavy := true
	for _, e := range t.Tree[cur] {
		to := e.to
		if t.Depth[to] > t.Depth[cur] {
			if heavy {
				t.Head[to] = t.Head[cur]
			} else {
				t.Head[to] = to
			}
			heavy = false
			t._dfsHld(to, times)
		}
	}
}

// 路径 [a,b] 与 [c,d] 的交集.
// 如果为空则返回 {-1,-1}，如果只有一个交点则返回 {x,x}，如果有两个交点则返回 {x,y}.
func (t *Tree32) PathIntersection(a, b, c, d int32) (int32, int32) {
	ab := t.Lca(a, b)
	ac := t.Lca(a, c)
	ad := t.Lca(a, d)
	bc := t.Lca(b, c)
	bd := t.Lca(b, d)
	cd := t.Lca(c, d)
	x := ab ^ ac ^ bc // meet(a,b,c)
	y := ab ^ ad ^ bd // meet(a,b,d)
	if x != y {
		return x, y
	}
	z := ac ^ ad ^ cd
	if x != z {
		x = -1
	}
	return x, x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
