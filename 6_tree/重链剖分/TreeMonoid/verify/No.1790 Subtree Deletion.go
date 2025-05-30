// 子树异或和/子树删除

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https: //yukicoder.me/problems/no/1790
	// 查询子树边权异或和/子树删除
	// !1. 子树删除可以用 removed 标记, 包括和父亲的那条边(如果有的话)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	tree := _NT(n)
	data := make([]E, n-1)
	for i := 0; i < n-1; i++ {
		var u, v, xor int
		fmt.Fscan(in, &u, &v, &xor)
		u--
		v--
		tree.AddEdge(u, v, xor)
		data[i] = E{xor: xor}
	}
	tree.Build(0)

	LM := NewLazyTreeMonoid(tree, data, false)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op, root int
		fmt.Fscan(in, &op, &root)
		root--
		if op == 1 { // !删除以root为根的子树(打标记)
			LM.UpdateSubtree(root, E{removed: true})
			if root != 0 {
				LM.UpdatePath(root, tree.Parent[root], E{removed: true}) // 注意这里,包括和父亲的边
			}
		} else { // 查询以root为根的子树的异或和
			fmt.Fprintln(out, LM.QuerySubtree(root).xor)
		}
	}
}

const INF = 1e18

type E = struct {
	xor     int
	removed bool
}
type Id = E

func e() E   { return E{} }
func id() Id { return E{} }
func op(e1, e2 E) E {
	if e1.removed {
		return e2
	}
	if e2.removed {
		return e1
	}
	return E{xor: e1.xor ^ e2.xor}
}
func mapping(f Id, g E) E {
	if f.removed || g.removed {
		return E{removed: true}
	}
	return E{xor: f.xor ^ g.xor}
}
func composition(f, g Id) Id {
	if f.removed || g.removed {
		return E{removed: true}
	}
	return E{xor: f.xor ^ g.xor}
}

type LazyTreeMonoid struct {
	tree     *_T
	n        int
	unit     E
	isVertex bool
	seg      *_LST
}

// !树的路径查询 + 区间修改, 维护的量需要满足幺半群的性质，且必须要满足交换律.
//
//	data: 顶点的值, 或者边的值.(边的编号为两个定点中较深的那个点的编号)
//	isVertex: data是否为顶点的值以及查询的时候是否是顶点权值.
func NewLazyTreeMonoid(tree *_T, data []E, isVertex bool) *LazyTreeMonoid {
	n := len(tree.Tree)
	res := &LazyTreeMonoid{tree: tree, n: n, unit: e(), isVertex: isVertex}
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
	res.seg = _NewLST(leaves, e, id, op, mapping, composition)

	return res
}

// 第i个顶点或者第i条边的值修改为e.
func (tm *LazyTreeMonoid) Set(i int, e E) {
	if !tm.isVertex {
		i = tm.tree.EidtoV(i)
	}
	i = tm.tree.LID[i]
	tm.seg.Set(i, e)
}

func (tm *LazyTreeMonoid) Get(i int) E {
	return tm.seg.Get(tm.tree.LID[i])
}

func (tm *LazyTreeMonoid) UpdatePath(start, target int, lazy Id) {
	path := tm.tree.GetPathDecomposition(start, target, tm.isVertex)
	for _, ab := range path {
		left, right := ab[0], ab[1]
		if left > right {
			left, right = right, left
		}
		tm.seg.Update(left, right+1, lazy)
	}
}

// 查询 start 到 target 的路径上的值.(点权/边权 由 isVertex 决定)
func (tm *LazyTreeMonoid) QueryPath(start, target int) E {
	path := tm.tree.GetPathDecomposition(start, target, tm.isVertex)
	val := tm.unit
	for _, ab := range path {
		a, b := ab[0], ab[1]
		var x E
		if a <= b {
			x = tm.seg.Query(a, b+1)
		} else {
			x = tm.seg.Query(b, a+1)
		}
		val = op(val, x)
	}
	return val
}

// 找到路径上最后一个 x 使得 QueryPath(start,x) 满足check函数.不存在返回-1.
func (tm *LazyTreeMonoid) MaxPath(check func(E) bool, start, target int) int {
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
		} else {
			x = tm.seg.Query(b, a+1)
		}
		if tmp := op(val, x); check(tmp) {
			val = tmp
			start = tm.tree.IdToNode[b]
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
			return tm.tree.IdToNode[i-1]
		} else {
			i := tm.seg.MinLeft(a+1, checkTmp)
			if i == a+1 {
				return start
			}
			return tm.tree.IdToNode[i]
		}
	}
	return target
}

func (tm *LazyTreeMonoid) QuerySubtree(root int) E {
	l, r := tm.tree.LID[root], tm.tree.RID[root]
	offset := 1
	if tm.isVertex {
		offset = 0
	}
	return tm.seg.Query(l+offset, r)
}

func (tm *LazyTreeMonoid) UpdateSubtree(root int, lazy Id) {
	l, r := tm.tree.LID[root], tm.tree.RID[root]
	offset := 1
	if tm.isVertex {
		offset = 0
	}
	tm.seg.Update(l+offset, r, lazy)
}

func (tm *LazyTreeMonoid) QueryAll() E { return tm.seg.QueryAll() }

func (tm *LazyTreeMonoid) _maxPathEdge(check func(E) bool, u, v int) int {
	lca := tm.tree.LCA(u, v)
	path := tm.tree.GetPathDecomposition(u, lca, tm.isVertex)
	val := tm.unit
	// climb
	for _, ab := range path {
		a, b := ab[0], ab[1]
		x := tm.seg.Query(b, a+1)
		if tmp := op(val, x); check(tmp) {
			val = tmp
			u = tm.tree.Parent[tm.tree.IdToNode[b]]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		i := tm.seg.MinLeft(a+1, checkTmp)
		if i == a+1 {
			return u
		}
		return tm.tree.Parent[tm.tree.IdToNode[i]]
	}

	// down
	path = tm.tree.GetPathDecomposition(lca, v, tm.isVertex)
	for _, ab := range path {
		a, b := ab[0], ab[1]
		x := tm.seg.Query(a, b+1)
		if tmp := op(val, x); check(tmp) {
			val = tmp
			u = tm.tree.IdToNode[b]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		i := tm.seg.MaxRight(a, checkTmp)
		if i == a {
			return u
		}
		return tm.tree.IdToNode[i-1]
	}
	return v
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
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
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

func (tree *_T) RootedLCA(u, v int, w int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, w) ^ tree.LCA(v, w)
}

func (tree *_T) Dist(u, v int, weighted bool) int {
	if weighted {
		return tree.DepthWeighted[u] + tree.DepthWeighted[v] - 2*tree.DepthWeighted[tree.LCA(u, v)]
	}
	return tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)]
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
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
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
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
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
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

type _LST struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id

	e           func() E
	id          func() Id
	op          func(E, E) E
	mapping     func(Id, E) E
	composition func(Id, Id) Id
	dataUint    E
	lazyUint    Id
}

func _NewLST(
	leaves []E,
	e func() E, id func() Id, op func(E, E) E,
	mapping func(Id, E) E, composition func(Id, Id) Id,
) *_LST {
	tree := &_LST{
		e: e, id: id, op: op, mapping: mapping, composition: composition,
		dataUint: e(), lazyUint: id(),
	}
	n := len(leaves)
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.dataUint
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.lazyUint
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = leaves[i]
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *_LST) Query(left, right int) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.dataUint
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
	sml, smr := tree.dataUint, tree.dataUint
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
func (tree *_LST) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *_LST) Update(left, right int, f Id) {
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
	for i := 1; i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *_LST) MinLeft(right int, predicate func(data E) bool) int {
	if right == 0 {
		return 0
	}
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}
	res := tree.dataUint
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}
		if !predicate(tree.op(tree.data[right], res)) {
			for right < tree.size {
				tree.pushDown(right)
				right = right<<1 | 1
				if predicate(tree.op(tree.data[right], res)) {
					res = tree.op(tree.data[right], res)
					right--
				}
			}
			return right + 1 - tree.size
		}
		res = tree.op(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}
	return 0
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (tree *_LST) MaxRight(left int, predicate func(data E) bool) int {
	if left == tree.n {
		return tree.n
	}
	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}
	res := tree.dataUint
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left <<= 1
				if predicate(tree.op(res, tree.data[left])) {
					res = tree.op(res, tree.data[left])
					left++
				}
			}
			return left - tree.size
		}
		res = tree.op(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return tree.n
}

// 单点查询(不需要 pushUp/op 操作时使用)
func (tree *_LST) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *_LST) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *_LST) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *_LST) pushDown(root int) {
	if tree.lazy[root] != tree.lazyUint {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.lazyUint
	}
}
func (tree *_LST) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
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
