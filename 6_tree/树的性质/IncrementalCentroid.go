// 在树的节点权重增加1的情况下，动态地维护树的重心。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	demo()
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2636
// !对每个i=0,1,...,n-1 求 i到点(0-i)的距离和
func judge() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	tree := NewTree(n)
	for i := 1; i < n; i++ {
		var parent, dist int
		fmt.Fscan(in, &parent, &dist)
		tree.AddEdge(parent-1, i, dist)
	}
	tree.Build(0)

	icc := NewIncrementalCentroid(tree, make([]Abel, n))
	distSum := 0
	for v := 0; v < n; v++ {
		preCent := icc.Centroid
		icc.AddOne(v)
		curCent := icc.Centroid
		distSum += tree.Dist(preCent, v, true)
		distSum -= tree.Dist(preCent, curCent, true)
		fmt.Fprintln(out, distSum)
	}
}

func demo() {
	tree := NewTree(5)
	tree.AddEdge(0, 1, 1)
	tree.AddEdge(0, 2, 1)
	tree.AddEdge(0, 3, 1)
	tree.AddEdge(1, 4, 1)
	tree.Build(0)
	icc := NewIncrementalCentroid(tree, make([]Abel, 5))
	fmt.Println(icc._getSubTreeWeight(1))
	fmt.Println(icc.Centroid) // 0
	icc.AddOne(1)
	fmt.Println(icc.Centroid) // 1
	icc.AddOne(2)
	icc.AddOne(2)
	fmt.Println(icc.Centroid) // 2
	icc.AddOne(0)
	fmt.Println(icc.MaxSubtree)
}

type Abel = int

func e() Abel             { return 0 }
func op(e1, e2 Abel) Abel { return e1 + e2 }
func inv(e Abel) Abel     { return -e }

type IncrementalCentroid struct {
	Centroid     int    // 重心
	MaxSubtree   [2]int // (maxChild, weight)
	Tree         *Tree
	n            int
	weightSum    int
	vertexWeight *_TreeAbelGroup // 维护顶点权值
	fs           *_fastSet
}

func NewIncrementalCentroid(tree *Tree, vertexWeight []Abel) *IncrementalCentroid {
	vertexWeight = append(vertexWeight[:0:0], vertexWeight...)
	res := &IncrementalCentroid{
		Tree: tree, n: len(tree.Tree),
		vertexWeight: NewTreeAbelGroup(tree, vertexWeight, true, false, true),
		fs:           _newFastSet(len(tree.Tree)),
	}
	return res
}

// 结点v的权值加1.
func (icc *IncrementalCentroid) AddOne(v int) {
	icc.fs.Insert(icc.Tree.LID[v])
	icc.vertexWeight.Add(v, 1)
	icc.weightSum++
	if v == icc.Centroid {
		return
	}
	wt := icc._getSubTreeWeight(v)
	if icc.MaxSubtree[1] < wt {
		icc.MaxSubtree = [2]int{icc.Tree.Jump(icc.Centroid, v, 1), wt}
	}
	if 2*wt <= icc.weightSum {
		return
	}
	k := wt
	if icc.weightSum != 2*k-1 {
		panic("icc.wtSm != 2*k-1")
	}
	to := icc._moveTo(v)
	icc.MaxSubtree = [2]int{icc.Tree.Jump(to, icc.Centroid, 1), k - 1}
	icc.Centroid = to
}

func (icc *IncrementalCentroid) GetCentroid() int {
	return icc.Centroid
}

func (icc *IncrementalCentroid) _moveTo(v int) int {
	cent := icc.Centroid
	if icc.Tree.IsInSubtree(v, cent) {
		a := icc.Tree.Jump(cent, v, 1)
		left, right := icc.Tree.LID[a], icc.Tree.RID[a]
		left, right = icc.fs.Next(left), icc.fs.Prev(right-1)
		x, y := icc.Tree.idToNode[left], icc.Tree.idToNode[right]
		return icc.Tree.LCA(x, y)
	}
	left, right := icc.Tree.LID[cent], icc.Tree.RID[cent]
	x := v
	ids := []int{}
	ids = append(ids, icc.fs.Next(0), icc.fs.Prev(left-1))
	ids = append(ids, icc.fs.Next(right), icc.fs.Prev(icc.n-1))
	for _, idx := range ids {
		if idx == -1 || idx == icc.n {
			continue
		}
		if left <= idx && idx < right {
			continue
		}
		y := icc.Tree.idToNode[idx]
		x = icc.Tree.RootedLCA(x, y, cent)
	}
	return x
}

// 获取结点v的子树大小，结点v不能是重心.
func (icc *IncrementalCentroid) _getSubTreeWeight(v int) int {
	cent := icc.Centroid
	if v == cent {
		panic("v == cent")
	}
	if icc.Tree.IsInSubtree(v, cent) {
		return icc.vertexWeight.QuerySubtree(icc.Tree.Jump(cent, v, 1))
	}
	return icc.weightSum - icc.vertexWeight.QuerySubtree(cent)
}

//
//
//

type _TreeAbelGroup struct {
	tree         *Tree
	n            int
	isVertex     bool
	pathQuery    bool
	subtreeQuery bool
	bit          *_FTree
	bitSubtree   *_FTree
	unit         Abel
}

// 静态树的路径查询，维护的量需要满足阿贝尔群的性质.
//
//	data: 顶点的值, 或者边的值.(边的编号为两个定点中较深的那个点的编号)
//	isVertex: data是否为顶点的值以及查询的时候是否是顶点权值.
func NewTreeAbelGroup(tree *Tree, data []Abel, isVertex, pathQuery, subtreeQuery bool) *_TreeAbelGroup {
	res := &_TreeAbelGroup{
		tree: tree, n: len(tree.Tree),
		isVertex:  isVertex,
		pathQuery: pathQuery, subtreeQuery: subtreeQuery,
		unit: e(),
	}

	if pathQuery {
		bitRaw := make([]Abel, 2*res.n)
		if isVertex {
			for v := 0; v < res.n; v++ {
				bitRaw[tree.ELID(v)] = data[v]
				bitRaw[tree.ERID(v)] = inv(data[v])
			}
		} else {
			for i := range bitRaw {
				bitRaw[i] = res.unit
			}
			for e := 0; e < res.n-1; e++ {
				v := tree.EidtoV(e)
				bitRaw[tree.ELID(v)] = data[e]
				bitRaw[tree.ERID(v)] = inv(data[e])
			}
		}
		res.bit = _NewFTree(len(bitRaw), bitRaw...)
	}

	if subtreeQuery {
		bitRaw := make([]Abel, res.n)
		if isVertex {
			for v := 0; v < res.n; v++ {
				bitRaw[tree.LID[v]] = data[v]
			}
		} else {
			for i := range bitRaw {
				bitRaw[i] = res.unit
			}
			for e := 0; e < res.n-1; e++ {
				v := tree.EidtoV(e)
				bitRaw[tree.LID[v]] = data[e]
			}
		}
		res.bitSubtree = _NewFTree(len(bitRaw), bitRaw...)
	}

	return res
}

// 第i个点或者第i条边的值加上x.
func (ta *_TreeAbelGroup) Add(i int, x Abel) {
	v := i
	if !ta.isVertex {
		v = ta.tree.EidtoV(i)
	}
	if ta.pathQuery {
		inv_x := inv(x)
		ta.bit.Update(ta.tree.ELID(v), x)
		ta.bit.Update(ta.tree.ERID(v), inv_x)
	}
	if ta.subtreeQuery {
		ta.bitSubtree.Update(ta.tree.LID[v], x)
	}
}

func (ta *_TreeAbelGroup) QueryPath(from, to int) Abel {
	lca := ta.tree.LCA(from, to)
	x1 := ta.bit.QueryRange(ta.tree.ELID(lca)+1, ta.tree.ELID(from)+1)
	offset := 1
	if ta.isVertex {
		offset = 0
	}
	x2 := ta.bit.QueryRange(ta.tree.ELID(lca)+offset, ta.tree.ELID(to)+1)
	return op(x1, x2)
}

func (ta *_TreeAbelGroup) QuerySubtree(root int) Abel {
	l, r := ta.tree.LID[root], ta.tree.RID[root]
	offset := 1
	if ta.isVertex {
		offset = 0
	}
	return ta.bitSubtree.QueryRange(l+offset, r)
}

type Tree struct {
	Tree                 [][][2]int // (next, weight)
	Edges                [][3]int   // (u, v, w)
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	idToNode             []int
	top, heavySon        []int
	timer                int
}

func NewTree(n int) *Tree {
	tree := make([][][2]int, n)
	lid := make([]int, n)
	rid := make([]int, n)
	idToNode := make([]int, n)
	top := make([]int, n)   // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n) // 深度
	depthWeighted := make([]int, n)
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	edges := make([][3]int, 0, n-1)
	for i := range parent {
		parent[i] = -1
	}

	return &Tree{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: depthWeighted,
		Parent:        parent,
		LID:           lid,
		RID:           rid,
		idToNode:      idToNode,
		top:           top,
		heavySon:      heavySon,
		Edges:         edges,
	}
}

// 添加无向边 u-v, 边权为w.
func (tree *Tree) AddEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Tree[v] = append(tree.Tree[v], [2]int{u, w})
	tree.Edges = append(tree.Edges, [3]int{u, v, w})
}

// 添加有向边 u->v, 边权为w.
func (tree *Tree) AddDirectedEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Edges = append(tree.Edges, [3]int{u, v, w})
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *Tree) Build(root int) {
	if root != -1 {
		tree.build(root, -1, 0, 0)
		tree.markTop(root, root)
	} else {
		for i := 0; i < len(tree.Tree); i++ {
			if tree.Parent[i] == -1 {
				tree.build(i, -1, 0, 0)
				tree.markTop(i, i)
			}
		}
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree) Id(root int) (int, int) {
	return tree.LID[root], tree.RID[root]
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *Tree) Eid(u, v int) int {
	if tree.LID[u] > tree.LID[v] {
		return tree.LID[u]
	}
	return tree.LID[v]
}

// 较深的那个点作为边的编号.
func (tree *Tree) EidtoV(eid int) int {
	e := tree.Edges[eid]
	u, v := e[0], e[1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (tree *Tree) LCA(u, v int) int {
	for {
		if tree.LID[u] > tree.LID[v] {
			u, v = v, u
		}
		if tree.top[u] == tree.top[v] {
			return u
		}
		v = tree.Parent[tree.top[v]]
	}
}

func (tree *Tree) RootedLCA(u, v int, root int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *Tree) Dist(u, v int, weighted bool) int {
	if weighted {
		return tree.DepthWeighted[u] + tree.DepthWeighted[v] - 2*tree.DepthWeighted[tree.LCA(u, v)]
	}
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
func (tree *Tree) KthAncestor(root, k int) int {
	if k > tree.Depth[root] {
		return -1
	}
	for {
		u := tree.top[root]
		if tree.LID[root]-k >= tree.LID[u] {
			return tree.idToNode[tree.LID[root]-k]
		}
		k -= tree.LID[root] - tree.LID[u] + 1
		root = tree.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (tree *Tree) Jump(from, to, step int) int {
	if step == 1 {
		if from == to {
			return -1
		}
		if tree.IsInSubtree(to, from) {
			return tree.KthAncestor(to, tree.Depth[to]-tree.Depth[from]-1)
		}
		return tree.Parent[from]
	}
	c := tree.LCA(from, to)
	dac := tree.Depth[from] - tree.Depth[c]
	dbc := tree.Depth[to] - tree.Depth[c]
	if step > dac+dbc {
		return -1
	}
	if step <= dac {
		return tree.KthAncestor(from, step)
	}
	return tree.KthAncestor(to, dac+dbc-step)
}

func (tree *Tree) CollectChild(root int) []int {
	res := []int{}
	for _, e := range tree.Tree[root] {
		next := e[0]
		if next != tree.Parent[root] {
			res = append(res, next)
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree) GetPathDecomposition(u, v int, vertex bool) [][2]int {
	up, down := [][2]int{}, [][2]int{}
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			down = append(down, [2]int{tree.LID[tree.top[v]], tree.LID[v]})
			v = tree.Parent[tree.top[v]]
		} else {
			up = append(up, [2]int{tree.LID[u], tree.LID[tree.top[u]]})
			u = tree.Parent[tree.top[u]]
		}
	}
	edgeInt := 1
	if vertex {
		edgeInt = 0
	}
	if tree.LID[u] < tree.LID[v] {
		down = append(down, [2]int{tree.LID[u] + edgeInt, tree.LID[v]})
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		up = append(up, [2]int{tree.LID[u], tree.LID[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

func (tree *Tree) GetPath(u, v int) []int {
	res := []int{}
	composition := tree.GetPathDecomposition(u, v, true)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.idToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.idToNode[i])
			}
		}
	}
	return res
}

// 以root为根时,结点v的子树大小.
func (tree *Tree) SubtreeSize(v, root int) int {
	if root == -1 {
		return tree.RID[v] - tree.LID[v]
	}
	if v == root {
		return len(tree.Tree)
	}
	x := tree.Jump(v, root, 1)
	if tree.IsInSubtree(v, x) {
		return tree.RID[v] - tree.LID[v]
	}
	return len(tree.Tree) - tree.RID[x] + tree.LID[x]
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *Tree) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

func (tree *Tree) ELID(u int) int {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *Tree) ERID(u int) int {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *Tree) build(cur, pre, dep, dist int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, e := range tree.Tree[cur] {
		next, weight := e[0], e[1]
		if next != pre {
			nextSize := tree.build(next, cur, dep+1, dist+weight)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next
			}
		}
	}
	tree.Depth[cur] = dep
	tree.DepthWeighted[cur] = dist
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *Tree) markTop(cur, top int) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.idToNode[tree.timer] = cur
	tree.timer++
	if tree.heavySon[cur] != -1 {
		tree.markTop(tree.heavySon[cur], top)
		for _, e := range tree.Tree[cur] {
			next := e[0]
			if next != tree.heavySon[cur] && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
}

type _FTree struct {
	n    int
	data []Abel
	unit Abel
}

func _NewFTree(n int, nums ...Abel) *_FTree {
	fw := &_FTree{n: n, unit: e()}
	if len(nums) == 0 {
		nums = make([]Abel, n)
		for i := range nums {
			nums[i] = fw.unit
		}
	}
	fw.build(nums)
	return fw
}

// [0, right)
func (fw *_FTree) QueryPrefix(right int) Abel {
	if right > fw.n {
		right = fw.n
	}
	res := fw.unit
	for right > 0 {
		res = op(res, fw.data[right-1])
		right &= right - 1
	}
	return res
}

// [left, right)
func (fw *_FTree) QueryRange(left, right int) Abel {
	if left < 0 {
		left = 0
	}
	if right > fw.n {
		right = fw.n
	}
	if left == 0 {
		return fw.QueryPrefix(right)
	}
	if left > right {
		return fw.unit
	}
	pos, neg := fw.unit, fw.unit
	for right > left {
		pos = op(pos, fw.data[right-1])
		right &= right - 1
	}
	for left > right {
		neg = op(neg, fw.data[left-1])
		left &= left - 1
	}
	return op(pos, inv(neg))
}

func (fw *_FTree) Update(i int, x Abel) {
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = op(fw.data[i-1], x)
	}
}

func (fw *_FTree) build(nums []Abel) {
	n := fw.n
	fw.data = append(fw.data, nums...)
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			fw.data[j-1] = op(fw.data[i-1], fw.data[j-1])
		}
	}
}

type _fastSet struct {
	n, lg int
	seg   [][]int
}

func _newFastSet(n int) *_fastSet {
	res := &_fastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func (fs *_fastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>uint((i&63)))&1 != 0
}

func (fs *_fastSet) Insert(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << uint((i & 63))
		i >>= 6
	}
}

func (fs *_fastSet) Erase(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] &= ^(1 << uint((i & 63)))
		if fs.seg[h][i>>6] != 0 {
			break
		}
		i >>= 6
	}
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *_fastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		if i>>6 == len(fs.seg[h]) {
			break
		}
		d := fs.seg[h][i>>6] >> uint((i & 63))
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *_fastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << uint((63 - i&63))
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *_fastSet) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *_fastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("_fastSet{%v}", strings.Join(res, ", "))
}

func (*_fastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*_fastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
