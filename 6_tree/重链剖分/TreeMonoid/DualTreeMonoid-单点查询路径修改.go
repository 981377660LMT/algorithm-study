// Get
// UpdatePath
// UpdateSubtree

package main

import "fmt"

func main() {
	tree := _NT(5)
	tree.AddEdge(0, 1, 1)
	tree.AddEdge(0, 2, 1)
	tree.AddEdge(1, 3, 1)
	tree.AddEdge(1, 4, 1)
	tree.Build(0)

	DM := NewDualTreeMonoid(tree, true)
	DM.UpdatePath(0, 4, 1)
	fmt.Println(DM.Get(0))
	DM.UpdateSubtree(1, 1)
	fmt.Println(DM.Get(0))
}

const INF int = 1e18

type Id = int

func id() Id                 { return 0 }
func composition(f, g Id) Id { return f + g }

type DualTreeMonoid struct {
	tree     *_T
	n        int
	isVertex bool
	seg      *_DST
}

// 树的单点查询 + 路径修改, 维护的量需要满足幺半群的性质.
//  isVertex: data是否为顶点的值以及查询的时候是否是顶点权值.
func NewDualTreeMonoid(tree *_T, isVertex bool) *DualTreeMonoid {
	res := &DualTreeMonoid{tree: tree, isVertex: isVertex, n: len(tree.Tree)}
	res.seg = _NewDST(res.n, id, composition)
	return res
}

func (tm *DualTreeMonoid) Get(i int) Id {
	v := i
	if !tm.isVertex {
		e := tm.tree.Edges[i]
		from, to := e[0], e[1]
		if tm.tree.Parent[from] == to {
			v = from
		} else {
			v = to
		}
	}
	return tm.seg.Get(tm.tree.LID[v])
}

func (tm *DualTreeMonoid) UpdatePath(start, target int, lazy Id) {
	path := tm.tree.GetPathDecomposition(start, target, tm.isVertex)
	for _, p := range path {
		a, b := p[0], p[1]
		if a <= b {
			tm.seg.Update(a, b+1, lazy)
		} else {
			tm.seg.Update(b, a+1, lazy)
		}
	}
}

func (tm *DualTreeMonoid) UpdateSubtree(root int, lazy Id) {
	l, r := tm.tree.LID[root], tm.tree.RID[root]
	offset := 1
	if tm.isVertex {
		offset = 0
	}
	tm.seg.Update(l+offset, r, lazy)
}

type _T struct {
	Tree                 [][][2]int // (next, weight)
	Edges                [][3]int   // (u, v, w)
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	IdToNode             []int
	top, heavySon        []int
	timer                int
}

func _NT(n int) *_T {
	tree := make([][][2]int, n)
	lid := make([]int, n)
	rid := make([]int, n)
	IdToNode := make([]int, n)
	top := make([]int, n)   // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n) // 深度
	depthWeighted := make([]int, n)
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	edges := make([][3]int, 0, n-1)
	for i := range parent {
		parent[i] = -1
	}

	return &_T{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: depthWeighted,
		Parent:        parent,
		LID:           lid,
		RID:           rid,
		IdToNode:      IdToNode,
		top:           top,
		heavySon:      heavySon,
		Edges:         edges,
	}
}

// 添加无向边 u-v, 边权为w.
func (tree *_T) AddEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Tree[v] = append(tree.Tree[v], [2]int{u, w})
	tree.Edges = append(tree.Edges, [3]int{u, v, w})
}

// 添加有向边 u->v, 边权为w.
func (tree *_T) AddDirectedEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Edges = append(tree.Edges, [3]int{u, v, w})
}

// root:0-based
//  当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *_T) Build(root int) {
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
func (tree *_T) Id(root int) (int, int) {
	return tree.LID[root], tree.RID[root]
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *_T) Eid(u, v int) int {
	if tree.LID[u] > tree.LID[v] {
		return tree.LID[u]
	}
	return tree.LID[v]
}

// 较深的那个点作为边的编号.
func (tree *_T) EidtoV(eid int) int {
	e := tree.Edges[eid]
	u, v := e[0], e[1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (tree *_T) LCA(u, v int) int {
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

func (tree *_T) RootedLCA(u, v int, root int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *_T) Dist(u, v int, weighted bool) int {
	if weighted {
		return tree.DepthWeighted[u] + tree.DepthWeighted[v] - 2*tree.DepthWeighted[tree.LCA(u, v)]
	}
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//  如果不存在第k个祖先，返回-1
func (tree *_T) KthAncestor(root, k int) int {
	if k > tree.Depth[root] {
		return -1
	}
	for {
		u := tree.top[root]
		if tree.LID[root]-k >= tree.LID[u] {
			return tree.IdToNode[tree.LID[root]-k]
		}
		k -= tree.LID[root] - tree.LID[u] + 1
		root = tree.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//  返回跳到的节点,如果不存在这样的节点,返回-1
func (tree *_T) Jump(from, to, step int) int {
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

func (tree *_T) CollectChild(root int) []int {
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
//  !eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *_T) GetPathDecomposition(u, v int, vertex bool) [][2]int {
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

func (tree *_T) GetPath(u, v int) []int {
	res := []int{}
	composition := tree.GetPathDecomposition(u, v, true)
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

// 以root为根时,结点v的子树大小.
func (tree *_T) SubtreeSize(v, root int) int {
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
func (tree *_T) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

func (tree *_T) ELID(u int) int {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *_T) ERID(u int) int {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *_T) build(cur, pre, dep, dist int) int {
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

func (tree *_T) markTop(cur, top int) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
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

//
//

type _DST struct {
	size, height int
	lazy         []Id

	id          func() Id
	composition func(f, g Id) Id
	lazyUnit    Id
}

func _NewDST(n int, id func() Id, composition func(f, g Id) Id) *_DST {
	res := &_DST{id: id, composition: composition, lazyUnit: id()}
	size := 1
	height := 0
	for size < n {
		size <<= 1
		height++
	}
	lazy := make([]Id, 2*size)
	for i := 0; i < 2*size; i++ {
		lazy[i] = res.lazyUnit
	}
	res.size = size
	res.height = height
	res.lazy = lazy
	return res
}

func (seg *_DST) Get(index int) Id {
	index += seg.size
	seg.thrust(index)
	return seg.lazy[index]
}

func (seg *_DST) Update(left, right int, value Id) {
	left += seg.size
	right += seg.size - 1
	seg.thrust(left)
	seg.thrust(right)
	l, r := left, right+1
	for l < r {
		if l&1 == 1 {
			seg.lazy[l] = seg.composition(seg.lazy[l], value)
			l++
		}
		if r&1 == 1 {
			r--
			seg.lazy[r] = seg.composition(seg.lazy[r], value)
		}
		l >>= 1
		r >>= 1
	}
}

func (seg *_DST) thrust(k int) {
	for i := seg.height; i > 0; i-- {
		seg.propagate(k >> i)
	}
}

func (seg *_DST) propagate(k int) {
	if seg.lazy[k] != seg.lazyUnit {
		seg.lazy[k<<1] = seg.composition(seg.lazy[k<<1], seg.lazy[k])
		seg.lazy[k<<1|1] = seg.composition(seg.lazy[k<<1|1], seg.lazy[k])
		seg.lazy[k] = seg.lazyUnit
	}
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
