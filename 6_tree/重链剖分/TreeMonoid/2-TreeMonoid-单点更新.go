// !- NewTreeMonoid(tree *Tree, data []E, isVertex bool) *TreeMonoid:
//   需要传入树、节点（或边）的初始值，以及一个布尔值表示给定的值是否是节点的值。
// !- Set(i int, x E):
//   将顶点i或者边i的值修改为x。
// !- Update(i int, x E):
//   将顶点i或者边i的值修改为原值与x的和。
// !- QueryPath(start, target int) E:
//   用于查询两个节点之间路径的值。(取决于isVertex)
// !- MaxPath(check func(E) bool, start, target int) int:
//   在树上查找最后一个满足条件的节点，使得从 start 到该节点的路径上的值满足某个条件。
//   若不存在这样的节点则返回 -1。
// !- QuerySubtree(root int) E:
//   用于查询以给定节点为根的子树上的所有节点值的代数和。

package main

import (
	"fmt"
)

func demo() {
	tree := NewTree(5)
	tree.AddEdge(0, 1, 1)
	tree.AddEdge(0, 2, 1)
	tree.AddEdge(1, 3, 1)
	tree.AddEdge(1, 4, 1)
	tree.Build(0)

	S := NewTreeMonoid(tree, []int{1, 2, 3, 4, 5}, false)
	fmt.Println(S.QueryPath(3, 4))
	fmt.Println(S.QueryPath(3, 1))
	fmt.Println(S.QuerySubtree(1))
	fmt.Println(S.MaxPath(func(x int) bool { return x <= 4 }, 3, 2))
}

const INF int = 1e18

type E = int

const IS_COMMUTATIVE = true // 幺半群是否满足交换律
func e() E                  { return 0 }
func op(e1, e2 E) E         { return e1 + e2 }

type TreeMonoid struct {
	tree     *Tree
	n        int
	unit     E
	isVertex bool
	seg      *_ST
	segR     *_ST
}

// 树的路径查询 + 单点修改, 维护的量需要满足幺半群的性质.
//  data: 顶点的值, 或者边的值.(边的编号为两个定点中较深的那个点的编号)
//  isVertex: data是否为顶点的值以及查询的时候是否是顶点权值.
func NewTreeMonoid(tree *Tree, data []E, isVertex bool) *TreeMonoid {
	n := len(tree.Tree)
	res := &TreeMonoid{tree: tree, n: n, unit: e(), isVertex: isVertex}
	leaves := make([]E, n)
	if isVertex {
		for v := range leaves {
			leaves[tree.LID[v]] = data[v]
		}
	} else {
		for i := range leaves {
			leaves[i] = res.unit
		}
		for e := 0; e < n-1; e++ {
			v := tree.EidtoV(e)
			leaves[tree.LID[v]] = data[e]
		}
	}
	res.seg = _NewS(leaves, e, op)
	if !IS_COMMUTATIVE {
		res.segR = _NewS(leaves, e, func(e1, e2 E) E { return op(e2, e1) }) // opRev
	}
	return res
}

// 第i个顶点或者第i条边的值修改为e.
func (tm *TreeMonoid) Set(i int, e E) {
	if !tm.isVertex {
		i = tm.tree.EidtoV(i)
	}
	i = tm.tree.LID[i]
	tm.seg.Set(i, e)
	if !IS_COMMUTATIVE {
		tm.segR.Set(i, e)
	}
}

// 第i个顶点或者第i条边的值与delta进行运算.
func (tm *TreeMonoid) Update(i int, delta E) {
	if !tm.isVertex {
		i = tm.tree.EidtoV(i)
	}
	i = tm.tree.LID[i]
	tm.seg.Update(i, delta)
	if !IS_COMMUTATIVE {
		tm.segR.Update(i, delta)
	}
}

// 查询 start 到 target 的路径上的值.(点权/边权 由 isVertex 决定)
func (tm *TreeMonoid) QueryPath(start, target int) E {
	path := tm.tree.GetPathDecomposition(start, target, tm.isVertex)
	val := tm.unit
	for _, ab := range path {
		a, b := ab[0], ab[1]
		var x E
		if a <= b {
			x = tm.seg.Query(a, b+1)
		} else if IS_COMMUTATIVE {
			x = tm.seg.Query(b, a+1)
		} else {
			x = tm.segR.Query(b, a+1)
		}
		val = op(val, x)
	}
	return val
}

// 找到路径上最后一个 x 使得 QueryPath(start,x) 满足check函数.不存在返回-1.
func (tm *TreeMonoid) MaxPath(check func(E) bool, start, target int) int {
	if !tm.isVertex {
		return tm._maxPathEdge(check, start, target)
	}
	if !check(tm.QueryPath(start, start)) {
		return -1
	}
	path := tm.tree.GetPathDecomposition(start, target, tm.isVertex)
	val := tm.unit
	for _, ab := range path {
		a, b := ab[0], ab[1]
		var x E
		if a <= b {
			x = tm.seg.Query(a, b+1)
		} else if IS_COMMUTATIVE {
			x = tm.seg.Query(b, a+1)
		} else {
			x = tm.segR.Query(b, a+1)
		}
		if tmp := op(val, x); check(tmp) {
			val = tmp
			start = tm.tree.idToNode[b]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		if a <= b {
			i := tm.seg.MaxRight(a, checkTmp)
			if i == a {
				return start
			}
			return tm.tree.idToNode[i-1]
		} else {
			var i int
			if IS_COMMUTATIVE {
				i = tm.seg.MinLeft(a+1, checkTmp)
			} else {
				i = tm.segR.MinLeft(a+1, checkTmp)
			}
			if i == a+1 {
				return start
			}
			return tm.tree.idToNode[i]
		}
	}
	return target
}

func (tm *TreeMonoid) QuerySubtree(root int) E {
	l, r := tm.tree.LID[root], tm.tree.RID[root]
	offset := 1
	if tm.isVertex {
		offset = 0
	}
	return tm.seg.Query(l+offset, r)
}

func (tm *TreeMonoid) _maxPathEdge(check func(E) bool, u, v int) int {
	lca := tm.tree.LCA(u, v)
	path := tm.tree.GetPathDecomposition(u, lca, tm.isVertex)
	val := tm.unit
	// climb
	for _, ab := range path {
		a, b := ab[0], ab[1]
		var x E
		if IS_COMMUTATIVE {
			x = tm.seg.Query(b, a+1)
		} else {
			x = tm.segR.Query(b, a+1)
		}
		if tmp := op(val, x); check(tmp) {
			val = tmp
			u = tm.tree.Parent[tm.tree.idToNode[b]]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		var i int
		if IS_COMMUTATIVE {
			i = tm.seg.MinLeft(a+1, checkTmp)
		} else {
			i = tm.segR.MinLeft(a+1, checkTmp)
		}
		if i == a+1 {
			return u
		}
		return tm.tree.Parent[tm.tree.idToNode[i]]
	}

	// down
	path = tm.tree.GetPathDecomposition(lca, v, tm.isVertex)
	for _, ab := range path {
		a, b := ab[0], ab[1]
		x := tm.seg.Query(a, b+1)
		if tmp := op(val, x); check(tmp) {
			val = tmp
			u = tm.tree.idToNode[b]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		i := tm.seg.MaxRight(a, checkTmp)
		if i == a {
			return u
		}
		return tm.tree.idToNode[i-1]
	}
	return v
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
//  当root设为-1时，会从0开始遍历未访问过的连通分量
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

// 返回边 u-v 对应的 欧拉序起点编号, 0-indexed.
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
//  如果不存在第k个祖先，返回-1
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
//  返回跳到的节点,如果不存在这样的节点,返回-1
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
//  !eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
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

//
//
type _ST struct {
	n, size int
	seg     []E

	e    func() E
	op   func(E, E) E
	unit E
}

func _NewS(leaves []E, e func() E, op func(E, E) E) *_ST {
	res := &_ST{e: e, op: op, unit: e()}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}

func (st *_ST) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.unit
	}
	return st.seg[index+st.size]
}

func (st *_ST) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

func (st *_ST) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *_ST) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.unit
	}
	leftRes, rightRes := st.unit, st.unit
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}

func (st *_ST) QueryAll() E { return st.seg[1] }

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *_ST) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := st.unit
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(st.op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if predicate(st.op(res, st.seg[left])) {
					res = st.op(res, st.seg[left])
					left++
				}
			}
			return left - st.size
		}
		res = st.op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *_ST) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	right += st.size
	res := st.unit
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(st.op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if predicate(st.op(st.seg[right], res)) {
					res = st.op(st.seg[right], res)
					right--
				}
			}
			return right + 1 - st.size
		}
		res = st.op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
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
