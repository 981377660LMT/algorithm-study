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
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var pos, startTime, alive int
			fmt.Fscan(in, &pos, &startTime, &alive)
			pos--
			check := func(e int) bool { return e <= alive } // 不消失可以达到的最远点
			to := S.MaxPath(check, pos, 0)
			p := tree.Parent[to]
			query = append(query, [3]int{1, lid[pos], startTime}) // 灯笼开始流动
			if p != -1 {                                          // 还没进入河流就消失了,在p处减一个灯笼
				query = append(query, [3]int{-1, lid[p], startTime + dist[pos]})
			}
		} else {
			var pos, endTime, null int
			fmt.Fscan(in, &pos, &endTime, &null)
			pos--
			query = append(query, [3]int{0, pos, endTime}) // 查询在pos处,<=endTime前可以看到的灯笼
		}
	}

	R := NewStaticPointAddRectangleSum(q, q)
	// fmt.Fprintln(out, query)
	for _, q := range query {
		op, pos, time := q[0], q[1], q[2]
		if op == 0 {
			l, r := lid[pos], tree.RID[pos]
			time += dist[pos]
			// fmt.Fprintln(out, op, pos, time, l, r, time+1)
			R.AddQuery(l, r, 0, time+1)
		} else if op == 1 {
			R.AddPoint(pos, time, 1)
		} else {
			R.AddPoint(pos, time, -1)
		}
	}

	res := R.Work()
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

const INF int = 1e18

type E = int

const IS_COMMUTATIVE = true // 幺半群是否满足交换律
func e() E                  { return 0 }
func op(e1, e2 E) E         { return e1 + e2 }

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
			start = st.tree.idToNode[b]
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
			return st.tree.idToNode[i-1]
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
				return st.tree.idToNode[i]
			}
			return st.tree.Parent[st.tree.idToNode[i]]
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
			u = st.tree.Parent[st.tree.idToNode[b]]
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
		return st.tree.Parent[st.tree.idToNode[i]]
	}

	// down
	path = st.tree.GetPathDecomposition(lca, v, st.isVertex)
	for _, ab := range path {
		a, b := ab[0], ab[1]
		x := st.seg.Query(a, b+1)
		if tmp := op(val, x); check(tmp) {
			val = tmp
			u = st.tree.idToNode[b]
			continue
		}
		checkTmp := func(x E) bool {
			return check(op(val, x))
		}
		i := st.seg.MaxRight(a, checkTmp)
		if i == a {
			return u
		}
		return st.tree.idToNode[i-1]
	}
	return v
}

type _T struct {
	Tree                 [][][2]int // (next, weight)
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	idToNode             []int
	top, heavySon        []int
	edges                [][3]int
	timer                int
}

func _NT(n int) *_T {
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

	return &_T{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: depthWeighted,
		Parent:        parent,
		LID:           lid,
		RID:           rid,
		idToNode:      idToNode,
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

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *_T) Eid(u, v int) int {
	if tree.LID[u] > tree.LID[v] {
		return tree.LID[u]
	}
	return tree.LID[v]
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
//
//	如果不存在第k个祖先，返回-1
func (tree *_T) KthAncestor(root, k int) int {
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

type Point struct{ x, y, w int }
type Query struct{ l, d, r, u int }
type DynamicPointAddRectangleSum struct {
	queries []interface{} // Point or Query
}

// 根据总点数初始化容量.
func NewPointAddRectangleSum(n int) *DynamicPointAddRectangleSum {
	return &DynamicPointAddRectangleSum{queries: make([]interface{}, 0, n)}
}

// 在(x,y)点上添加w权重.
func (dpars *DynamicPointAddRectangleSum) AddPoint(x, y, w int) {
	dpars.queries = append(dpars.queries, Point{x, y, w})
}

// 添加查询为区间 [x1, x2) * [y1, y2) 的权重和.
func (dpars *DynamicPointAddRectangleSum) AddQuery(x1, x2, y1, y2 int) {
	dpars.queries = append(dpars.queries, Query{x1, y1, x2, y2})
}

// 按照添加顺序返回所有查询结果..
func (dpars *DynamicPointAddRectangleSum) Work() []int {
	q := len(dpars.queries)
	rev := make([]int, q)
	sz := 0
	for i := 0; i < q; i++ {
		if _, ok := dpars.queries[i].(Query); ok { // holds_alternative
			rev[i] = sz
			sz++
		}
	}

	res := make([]int, sz)
	rangeQ := [][]int{{0, q}}
	for len(rangeQ) > 0 {
		l, r := rangeQ[0][0], rangeQ[0][1]
		rangeQ = rangeQ[1:]
		m := (l + r) >> 1
		solver := NewStaticPointAddRectangleSum(0, 0)
		for k := l; k < m; k++ {
			if p, ok := dpars.queries[k].(Point); ok {
				solver.AddPoint(p.x, p.y, p.w)
			}
		}

		for k := m; k < r; k++ {
			if q, ok := dpars.queries[k].(Query); ok {
				solver.AddQuery(q.l, q.r, q.d, q.u)
			}
		}

		sub := solver.Work()
		for k, t := m, 0; k < r; k++ {
			if _, ok := dpars.queries[k].(Query); ok {
				res[rev[k]] += sub[t]
				t++
			}
		}

		if l+1 < m {
			rangeQ = append(rangeQ, []int{l, m})
		}
		if m+1 < r {
			rangeQ = append(rangeQ, []int{m, r})
		}
	}

	return res
}

type StaticPointAddRectangleSum struct {
	points  []Point
	queries []Query
}

// 指定点集和查询个数初始化容量.
func NewStaticPointAddRectangleSum(n, q int) *StaticPointAddRectangleSum {
	return &StaticPointAddRectangleSum{
		points:  make([]Point, 0, n),
		queries: make([]Query, 0, q),
	}
}

// 在(x,y)点上添加w权重.
func (sp *StaticPointAddRectangleSum) AddPoint(x, y, w int) {
	sp.points = append(sp.points, Point{x: x, y: y, w: w})
}

// 添加查询为区间 [x1, x2) * [y1, y2) 的权重和.
func (sp *StaticPointAddRectangleSum) AddQuery(x1, x2, y1, y2 int) {
	sp.queries = append(sp.queries, Query{l: x1, r: x2, d: y1, u: y2})
}

// 按照添加顺序返回所有查询结果..
func (sp *StaticPointAddRectangleSum) Work() []int {
	n := len(sp.points)
	q := len(sp.queries)
	res := make([]int, q)
	if n == 0 || q == 0 {
		return res
	}

	sort.Slice(sp.points, func(i, j int) bool { return sp.points[i].y < sp.points[j].y })
	ys := make([]int, 0, n)
	for i := range sp.points {
		if len(ys) == 0 || ys[len(ys)-1] != sp.points[i].y {
			ys = append(ys, sp.points[i].y)
		}
		sp.points[i].y = len(ys) - 1
	}

	type Q struct {
		x    int
		d, u int
		t    bool
		idx  int
	}

	qs := make([]Q, 0, q+q)
	for i := 0; i < q; i++ {
		query := sp.queries[i]
		d := sort.SearchInts(ys, query.d)
		u := sort.SearchInts(ys, query.u)
		qs = append(qs, Q{x: query.l, d: d, u: u, t: false, idx: i}, Q{x: query.r, d: d, u: u, t: true, idx: i})
	}

	sort.Slice(sp.points, func(i, j int) bool { return sp.points[i].x < sp.points[j].x })
	sort.Slice(qs, func(i, j int) bool { return qs[i].x < qs[j].x })

	j := 0
	bit := newBinaryIndexedTree(len(ys))
	for i := range qs {
		for j < n && sp.points[j].x < qs[i].x {
			bit.Apply(sp.points[j].y, sp.points[j].w)
			j++
		}
		if qs[i].t {
			res[qs[i].idx] += bit.ProdRange(qs[i].d, qs[i].u)
		} else {
			res[qs[i].idx] -= bit.ProdRange(qs[i].d, qs[i].u)
		}
	}

	return res
}

type binaryIndexedTree struct {
	n    int
	log  int
	data []int
}

// 長さ n の 0で初期化された配列で構築する.
func newBinaryIndexedTree(n int) *binaryIndexedTree {
	return &binaryIndexedTree{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

// 配列で構築する.
func newBinaryIndexedTreeFrom(arr []int) *binaryIndexedTree {
	res := newBinaryIndexedTree(len(arr))
	res.build(arr)
	return res
}

// 要素 i に値 v を加える.
func (b *binaryIndexedTree) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r) の要素の総和を求める.
func (b *binaryIndexedTree) Prod(r int) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *binaryIndexedTree) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *binaryIndexedTree) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.
func (b *binaryIndexedTree) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

func (b *binaryIndexedTree) build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}
