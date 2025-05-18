package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	// 放花灯
	// https://yukicoder.me/problems/no/1216
	// Paken River 拥有 n 个“检查点”，每个检查点被编号为 1 到 n。
	// 检查点 i 和检查点 j 直接连接，灯笼将沿着河流从 i 到 j，用时间 t 流动。
	// 河流具有树形结构，其节点是检查点，其根是出口。
	// 检查点 1 是 Paken River 的出口，河流从这里流入。
	// 灯笼被流动到某个检查点后，它将沿着河流顺流而下，一段时间后将关闭。
	// 回答以下 q 个查询。
	// !0:添加查询：在时间 t 将灯笼从检查点 i 流到河流上。灯笼在 alive 秒后消失。
	// !1:回答查询：在检查点 i，输出在时刻 t 之前亮起并可见的灯笼的总数(子树里灯笼树)。 包括在时刻 t 正好可见、刚开始流动、刚好在到达时关闭的灯笼。
	// n<=5e4,q<=1e5
	// 1<=i<=n,0<=t<=1e12,0<=alive<=1e12

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	tree := _NT(n)
	for i := 0; i < n-1; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		a--
		b--
		tree.AddEdge(a, b, c)
	}
	tree.Build(0)

	edgeW := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		edgeW[i] = tree.edges[i][2] // weight
	}
	S := NewStaticTreeMonoid(tree, edgeW, false)

	dist := tree.DepthWeighted
	lid := tree.LID
	query := make([][3]int, 0, q)
	xs, ys := []int{}, []int{}
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var pos, startTime, alive int
			fmt.Fscan(in, &pos, &startTime, &alive)
			pos--
			check := func(e int) bool { return e <= alive } // 不消失可以达到的最远点
			to := S.MaxPath(check, pos, 0)
			parent := tree.Parent[to]
			query = append(query, [3]int{1, lid[pos], startTime + dist[pos]}) // 灯笼开始流动
			xs = append(xs, lid[pos])
			ys = append(ys, startTime+dist[pos])
			if parent != -1 { // 还没进入河流就消失了,在p处减一个灯笼
				xs = append(xs, lid[parent])
				ys = append(ys, startTime+dist[pos])
				query = append(query, [3]int{-1, lid[parent], startTime + dist[pos]})
			}
		} else {
			var pos, endTime, null int
			fmt.Fscan(in, &pos, &endTime, &null)
			pos--
			query = append(query, [3]int{0, pos, endTime}) // 查询在pos处,<=endTime前可以看到的灯笼
		}
	}

	R := NewBIT2DSparse(xs, ys, false)

	for _, q := range query {
		op, pos, time := q[0], q[1], q[2]
		if op == 0 {
			l, r := lid[pos], tree.RID[pos]
			time += dist[pos]

			fmt.Fprintln(out, R.QueryRange(l, r, 0, time+1))
		} else if op == 1 {
			R.Update(pos, time, 1)
		} else {
			R.Update(pos, time, -1)
		}
	}

}

const INF int = 1e18

type E = int

const IS_COMMUTATIVE = true // 幺半群是否满足交换律
type Able = int

// 需要是阿贝尔群(满足交换律)
func e() Able           { return 0 }
func op(a, b Able) Able { return a + b }
func inv(a Able) Able   { return -a }

type StaticTreeMonoid struct {
	tree     *_T
	n        int
	unit     E
	isVertex bool
	seg      *DisjointSparseTable
	segR     *DisjointSparseTable
}

// data: 顶点的值, 或者边的值.(边的编号为两个定点中较深的那个点的编号)
// isVertex: data是否为顶点的值以及查询的时候是否是顶点权值.
func NewStaticTreeMonoid(tree *_T, data []E, isVertex bool) *StaticTreeMonoid {
	n := len(tree.Tree)
	res := &StaticTreeMonoid{tree: tree, n: n, unit: e(), isVertex: isVertex}
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
	res.seg = NewDisjointSparse(leaves, e, op)
	if !IS_COMMUTATIVE {
		res.segR = NewDisjointSparse(leaves, e, func(e1, e2 E) E { return op(e2, e1) }) // opRev
	}
	return res
}

// 查询 start 到 target 的路径上的值.(点权/边权 由 isVertex 决定)
func (st *StaticTreeMonoid) QueryPath(start, target int) E {
	path := st.tree.GetPathDecomposition(start, target, st.isVertex)
	val := st.unit
	for _, ab := range path {
		a, b := ab[0], ab[1]
		var x E
		if a <= b {
			x = st.seg.Query(a, b+1)
		} else if IS_COMMUTATIVE {
			x = st.seg.Query(b, a+1)
		} else {
			x = st.segR.Query(b, a+1)
		}
		val = op(val, x)
	}
	return val
}

// 找到路径上最后一个 x 使得 QueryPath(start,x) 满足check函数.不存在返回-1.
func (st *StaticTreeMonoid) MaxPath(check func(E) bool, start, target int) int {
	if !st.isVertex {
		return st._maxPathEdge(check, start, target)
	}
	if !check(st.QueryPath(start, start)) {
		return -1
	}
	path := st.tree.GetPathDecomposition(start, target, st.isVertex)
	val := st.unit
	for _, ab := range path {
		a, b := ab[0], ab[1]
		var x E
		if a <= b {
			x = st.seg.Query(a, b+1)
		} else if IS_COMMUTATIVE {
			x = st.seg.Query(b, a+1)
		} else {
			x = st.segR.Query(b, a+1)
		}
		if tmp := op(val, x); check(tmp) {
			val = tmp
			start = st.tree.IdToNode[b]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		if a <= b {
			i := st.seg.MaxRight(a, checkTmp)
			if i == a {
				return start
			}
			return st.tree.IdToNode[i-1]
		} else {
			var i E
			if IS_COMMUTATIVE {
				i = st.seg.MinLeft(a+1, checkTmp)
			} else {
				i = st.segR.MinLeft(a+1, checkTmp)
			}
			if i == a+1 {
				return start
			}
			if st.isVertex {
				return st.tree.IdToNode[i]
			}
			return st.tree.Parent[st.tree.IdToNode[i]]
		}
	}
	return target
}

func (st *StaticTreeMonoid) QuerySubtree(root int) E {
	l, r := st.tree.LID[root], st.tree.RID[root]
	offset := 1
	if st.isVertex {
		offset = 0
	}
	return st.seg.Query(l+offset, r)
}

func (st *StaticTreeMonoid) _maxPathEdge(check func(E) bool, u, v int) int {
	lca := st.tree.LCA(u, v)
	path := st.tree.GetPathDecomposition(u, lca, st.isVertex)
	val := st.unit
	// climb
	for _, ab := range path {
		a, b := ab[0], ab[1]
		var x E
		if IS_COMMUTATIVE {
			x = st.seg.Query(b, a+1)
		} else {
			x = st.segR.Query(b, a+1)
		}
		if tmp := op(val, x); check(tmp) {
			val = tmp
			u = st.tree.Parent[st.tree.IdToNode[b]]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		var i E
		if IS_COMMUTATIVE {
			i = st.seg.MinLeft(a+1, checkTmp)
		} else {
			i = st.segR.MinLeft(a+1, checkTmp)
		}
		if i == a+1 {
			return u
		}
		return st.tree.Parent[st.tree.IdToNode[i]]
	}

	// down
	path = st.tree.GetPathDecomposition(lca, v, st.isVertex)
	for _, ab := range path {
		a, b := ab[0], ab[1]
		x := st.seg.Query(a, b+1)
		if tmp := op(val, x); check(tmp) {
			val = tmp
			u = st.tree.IdToNode[b]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		i := st.seg.MaxRight(a, checkTmp)
		if i == a {
			return u
		}
		return st.tree.IdToNode[i-1]
	}
	return v
}

type _T struct {
	Tree                 [][][2]int // (next, weight)
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	IdToNode             []int
	top, heavySon        []int
	edges                [][3]int
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
		edges:         edges,
	}
}

// 添加无向边 u-v, 边权为w.
func (tree *_T) AddEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.Tree[v] = append(tree.Tree[v], [2]int{u, w})
	tree.edges = append(tree.edges, [3]int{u, v, w})
}

// 添加有向边 u->v, 边权为w.
func (tree *_T) AddDirectedEdge(u, v, w int) {
	tree.Tree[u] = append(tree.Tree[u], [2]int{v, w})
	tree.edges = append(tree.edges, [3]int{u, v, w})
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
	e := tree.edges[eid]
	from, to := e[0], e[1]
	if tree.Parent[from] == to {
		return from
	}
	return to
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

type DisjointSparseTable struct {
	n, log int
	data   [][]E
	unit   E
	op     func(E, E) E
}

// DisjointSparseTable 支持幺半群的区间静态查询.
//
//	eg: 区间乘积取模/区间仿射变换...
func NewDisjointSparse(leaves []E, e func() E, op func(E, E) E) *DisjointSparseTable {
	res := &DisjointSparseTable{}
	n := len(leaves)
	log := 1
	for (1 << log) < n {
		log++
	}
	data := make([][]E, log)
	data[0] = append(data[0], leaves...)
	for i := 1; i < log; i++ {
		data[i] = append(data[i], data[0]...)
		v := data[i]
		b := 1 << i
		for m := b; m <= n; m += 2 * b {
			l, r := m-b, min(m+b, n)
			for j := m - 1; j >= l+1; j-- {
				v[j-1] = op(v[j-1], v[j])
			}
			for j := m; j < r-1; j++ {
				v[j+1] = op(v[j], v[j+1])
			}
		}
	}
	res.n = n
	res.log = log
	res.data = data
	res.unit = e()
	res.op = op
	return res
}

func (ds *DisjointSparseTable) Query(start, end int) E {
	if start == end {
		return ds.unit
	}
	end--
	if start == end {
		return ds.data[0][start]
	}
	k := 31 - bits.LeadingZeros32(uint32(start^end))
	return ds.op(ds.data[k][start], ds.data[k][end])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *DisjointSparseTable) MaxRight(left int, check func(e E) bool) int {
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
func (ds *DisjointSparseTable) MinLeft(right int, check func(e E) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
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

type BIT2DSparse struct {
	n          int
	keyX       []int
	keyY       []int
	minX       int
	indptr     []int
	data       []Able
	discretize bool
	unit       Able
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewBIT2DSparse(xs, ys []int, discretize bool) *BIT2DSparse {
	res := &BIT2DSparse{discretize: discretize, unit: e()}
	ws := make([]Able, len(xs))
	for i := range ws {
		ws[i] = res.unit
	}
	res._build(xs, ys, ws)
	return res
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewBIT2DSparseWithWeights(xs, ys []int, ws []Able, discretize bool) *BIT2DSparse {
	res := &BIT2DSparse{discretize: discretize, unit: e()}
	res._build(xs, ys, ws)
	return res
}

// 点 (x,y) 的值加上 val.
func (fwt *BIT2DSparse) Update(x, y int, val Able) {
	i := fwt._xtoi(x)
	for ; i < fwt.n; i += ((i + 1) & -(i + 1)) {
		fwt._add(i, y, val)
	}
}

// [lx,rx) * [ly,ry)
func (t *BIT2DSparse) QueryRange(lx, rx, ly, ry int) Able {
	pos, neg := t.unit, t.unit
	l, r := t._xtoi(lx)-1, t._xtoi(rx)-1
	for l < r {
		pos = op(pos, t._prodI(r, ly, ry))
		r -= ((r + 1) & -(r + 1))
	}
	for r < l {
		neg = op(neg, t._prodI(l, ly, ry))
		l -= ((l + 1) & -(l + 1))
	}
	return op(pos, inv(neg))
}

// [0,rx) * [0,ry)
func (t *BIT2DSparse) QueryPrefix(rx, ry int) Able {
	pos := t.unit
	r := t._xtoi(rx) - 1
	for r >= 0 {
		pos = op(pos, t._prefixProdI(r, ry))
		r -= ((r + 1) & -(r + 1))
	}
	return pos
}

func (t *BIT2DSparse) _add(i int, y int, val Able) {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	j := bisectLeft(t.keyY, y, lid, lid+n-1) - lid
	for j < n {
		t.data[lid+j] = op(t.data[lid+j], val)
		j += ((j + 1) & -(j + 1))
	}
}

func (t *BIT2DSparse) _prodI(i int, ly, ry int) Able {
	pos, neg := t.unit, t.unit
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	left := bisectLeft(t.keyY, ly, lid, lid+n-1) - lid - 1
	right := bisectLeft(t.keyY, ry, lid, lid+n-1) - lid - 1
	for left < right {
		pos = op(pos, t.data[lid+right])
		right -= ((right + 1) & -(right + 1))
	}
	for right < left {
		neg = op(neg, t.data[lid+left])
		left -= ((left + 1) & -(left + 1))
	}
	return op(pos, inv(neg))
}

func (t *BIT2DSparse) _prefixProdI(i int, ry int) Able {
	pos := t.unit
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	R := bisectLeft(t.keyY, ry, lid, lid+n-1) - lid - 1
	for R >= 0 {
		pos = op(pos, t.data[lid+R])
		R -= ((R + 1) & -(R + 1))
	}
	return pos
}

func (ft *BIT2DSparse) _build(X, Y []int, wt []Able) {
	if len(X) != len(Y) || len(X) != len(wt) {
		panic("Lengths of X, Y, and wt must be equal.")
	}

	if ft.discretize {
		ft.keyX = unique(X)
		ft.n = len(ft.keyX)
	} else {
		if len(X) > 0 {
			min_, max_ := 0, 0
			for _, x := range X {
				if x < min_ {
					min_ = x
				}
				if x > max_ {
					max_ = x
				}
			}
			ft.minX = min_
			ft.n = max_ - min_ + 1
		}
		ft.keyX = make([]int, ft.n)
		for i := 0; i < ft.n; i++ {
			ft.keyX[i] = ft.minX + i
		}
	}

	N := ft.n
	keyYRaw := make([][]int, N)
	datRaw := make([][]Able, N)
	indices := argSort(Y)

	for _, i := range indices {
		ix := ft._xtoi(X[i])
		y := Y[i]
		for ix < N {
			kY := keyYRaw[ix]
			if len(kY) == 0 || kY[len(kY)-1] < y {
				keyYRaw[ix] = append(keyYRaw[ix], y)
				datRaw[ix] = append(datRaw[ix], wt[i])
			} else {
				datRaw[ix][len(datRaw[ix])-1] = op(datRaw[ix][len(datRaw[ix])-1], wt[i])
			}
			ix += ((ix + 1) & -(ix + 1))
		}
	}

	ft.indptr = make([]int, N+1)
	for i := 0; i < N; i++ {
		ft.indptr[i+1] = ft.indptr[i] + len(keyYRaw[i])
	}
	ft.keyY = make([]int, ft.indptr[N])
	ft.data = make([]Able, ft.indptr[N])

	for i := 0; i < N; i++ {
		for j := 0; j < ft.indptr[i+1]-ft.indptr[i]; j++ {
			ft.keyY[ft.indptr[i]+j] = keyYRaw[i][j]
			ft.data[ft.indptr[i]+j] = datRaw[i][j]
		}
	}

	for i := 0; i < N; i++ {
		n := ft.indptr[i+1] - ft.indptr[i]
		for j := 0; j < n-1; j++ {
			k := j + ((j + 1) & -(j + 1))
			if k < n {
				ft.data[ft.indptr[i]+k] = op(ft.data[ft.indptr[i]+k], ft.data[ft.indptr[i]+j])
			}
		}
	}
}

func (ft *BIT2DSparse) _xtoi(x int) int {
	if ft.discretize {
		return bisectLeft(ft.keyX, x, 0, len(ft.keyX)-1)
	}
	tmp := x - ft.minX
	if tmp < 0 {
		tmp = 0
	} else if tmp > ft.n {
		tmp = ft.n
	}
	return tmp
}

func bisectLeft(nums []int, x int, left, right int) int {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func unique(nums []int) (sorted []int) {
	set := make(map[int]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	return
}

func argSort(nums []int) []int {
	order := make([]int, len(nums))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
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
