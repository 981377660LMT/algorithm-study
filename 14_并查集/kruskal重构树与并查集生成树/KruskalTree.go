// reference: https://nyaannyaan.github.io/library/forest/process-of-merging-forest.hpp
// 表示合并过程的树,按照edges中边的顺序合并顶点.(kruskal重构树)
// process-of-merging-forest/KruskalTree
// 加入第i条边时其所在连通块大小就是权值为i的点所在子树的叶节点数量
//
//
// 例如 0和1合并,边权为2；0和2合并,边权为3,那么返回值为:
// forest: [[] [] [] [{0 1}] [{2 3}]]  (森林的有向图邻接表)
// roots: [4] (新图中的各个根节点)
// values: [0 0 0 2 3] (每个结点的权值)
//
//      4(3)
//     /    \
//    3(2)   2
//   / \
//  0   1
//
// 0-n-1 为原始顶点, n-2n-2 为辅助顶点

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	luogu4197()
}

// https://atcoder.jp/contests/agc002/tasks/agc002_d
// 一张`连通`图，q 次询问从两个点 x 和 y 出发，
// 希望经过的点数量等于 z（每个点可以重复经过，但是重复经过只计算一次）
// 求经过的边最大编号最小是多少。
//
// !根据性质：kruskal重构树中两个点u,v路径上的最大边权等于lca(u,v)的点权; kruskal重构树是一个大根堆.
// 因此可以二分编号，对于询问(x,y,z),我们从x和y分别倍增向上跳到点权大于当前二分值的位置，
// 然后再判断此时跳到节点的子树中的叶子节点数量是否大于等于z.
func StampRally() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = Edge{u, v, i + 1}
	}

	// !注意图是连通的，因此只有一个根节点
	forest, roots, values := KruskalTree(n, edges) // 已经按照边权排序
	tree := _NewTree(forest, roots)
	treeMonoid := NewStaticTreeMonoid(tree, values)

	root := roots[0]
	subTreeLeafCount := make([]int, len(forest)) // 子树中的叶子节点数量
	var dfs func(int) int
	dfs = func(cur int) int {
		if cur < n { // 原始顶点
			subTreeLeafCount[cur] = 1
		}
		for _, to := range forest[cur] {
			subTreeLeafCount[cur] += dfs(to)
		}
		return subTreeLeafCount[cur]
	}
	dfs(root)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x, y, z int
		fmt.Fscan(in, &x, &y, &z)
		x--
		y--

		check := func(mid int) bool {
			to1 := treeMonoid.MaxPath(func(e E) bool { return e <= mid }, x, root)
			to2 := treeMonoid.MaxPath(func(e E) bool { return e <= mid }, y, root)
			if to1 == to2 {
				return subTreeLeafCount[to1] >= z
			} else {
				return subTreeLeafCount[to1]+subTreeLeafCount[to2] >= z
			}
		}
		left, right := 1, m
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		fmt.Fprintln(out, left)
	}
}

// https://www.luogu.com.cn/problem/CF1706E
// Qpwoeirut and Vertices
// 给出 n 个点， m 条边的不带权连通无向图， q 次询问至少要加完编号前多少的边，
// 才能使得 [start,end) 中的所有点两两连通。
//
// 考虑 Kruskal 重构树。
// 由于这题求的是按照边的顺序，
// 所以我们把边的权值赋为边的序号，排序后建出 Kruskal 重构树
// !由kruskal重构树性质，等价于求[start,end)中的所有点的lca.
func CF1706E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	solve := func() {
		var n, m, q int
		fmt.Fscan(in, &n, &m, &q)
		edges := make([]Edge, m)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u, v = u-1, v-1
			edges[i] = Edge{u, v, i + 1}
		}
		forest, roots, values := KruskalTree(n, edges)
		tree := _NewTree(forest, roots)
		order := make([]int, n)
		for i := range order {
			order[i] = i
		}
		rangeLca := RangeLCA(order, tree.LCA)

		for i := 0; i < q; i++ {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			lca := rangeLca(start, end)
			fmt.Fprintln(out, values[lca])
		}
	}

	for t := 0; t < T; t++ {
		solve()
	}
}

// https://www.luogu.com.cn/problem/P4197
// 给定一张无向带权图和q个询问(u，limit，k)，每个点有一个得分，每条边有一个边权。
// !每次查询从u出发，只经过边权小于等于limit的点中第k大的点的得分，如果不存在这样的点，输出-1。
// P4197 Peaks
func luogu4197() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	scores := make([]int, n)
	for i := range scores {
		fmt.Fscan(in, &scores[i])
	}
	edges := make([]Edge, m)
	for i := range edges {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges[i] = Edge{u, v, w}
	}

	sort.Slice(edges, func(i, j int) bool { return edges[i].weight < edges[j].weight })
	forest, roots, values := KruskalTree(n, edges)
	tree := _NewTree(forest, roots)
	treeMonoid := NewStaticTreeMonoid(tree, values)

	belongRoot := make([]int, len(forest))
	lid, rid := make([]int, len(forest)), make([]int, len(forest))
	dfn := 0
	var dfs func(int, int)
	dfs = func(cur, root int) {
		belongRoot[cur] = root
		lid[cur] = dfn
		// !注意kruskal中的dfs序带上了辅助结点，需要忽略这些结点
		if cur < n {
			dfn++
		}
		for _, to := range forest[cur] {
			dfs(to, root)
		}
		rid[cur] = dfn
	}
	for _, root := range roots {
		dfs(root, root)
	}

	newScores := make([]int, n)
	for i, v := range scores {
		newScores[lid[i]] = v
	}
	wm := NewWaveletMatrix(newScores)

	for i := 0; i < q; i++ {
		var u, limit, k int
		fmt.Fscan(in, &u, &limit, &k)
		u--
		k--
		curRoot := belongRoot[u]

		to := treeMonoid.MaxPath(func(x int) bool { return x <= limit }, u, curRoot)
		if to == -1 {
			fmt.Fprintln(out, -1)
			continue
		}
		lid, rid := lid[to], rid[to]
		if k >= rid-lid {
			fmt.Fprintln(out, -1)
			continue
		}
		fmt.Fprintln(out, wm.KthMax(lid, rid, k))
	}

}

// https://www.luogu.com.cn/problem/P1967
// P1967 [NOIP2013 提高组] 货车运输
// 给定一张无向带权图和q个询问(u,v), 每次查询从u到v的路径上的最小边权的最大值.如果不存在这样的路径,输出-1.
// 解：
// !从u到v最大边权最小值即为kruskal重构树中lca(u,v)的值.
func luogu1967() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	edges := make([]Edge, m)
	uf := NewUf(n)
	for i := range edges {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges[i] = Edge{u, v, w}
		uf.Union(u, v, nil)
	}

	// 因为要求最大的最小边权,所以边权从大到小排序
	sort.Slice(edges, func(i, j int) bool { return edges[i].weight > edges[j].weight })
	forest, roots, values := KruskalTree(n, edges)
	tree := _NewTree(forest, roots)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		if uf.Find(u) != uf.Find(v) {
			fmt.Fprintln(out, -1)
			continue
		}
		fmt.Fprintln(out, values[tree.LCA(u, v)])
	}
}

// https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths-ii/description/
// 1724. 检查边长度限制的路径是否存在
// 当存在一条从 p 到 q 的路径，且路径中每条边的距离都严格小于 limit 时，返回 true
type DistanceLimitedPathsExist struct {
	uf         *Uf
	tree       *_Tree
	treeMonoid *StaticTreeMonoid
}

func Constructor(n int, edgeList [][]int) DistanceLimitedPathsExist {
	uf := NewUf(n)
	for _, e := range edgeList {
		uf.Union(e[0], e[1], nil)
	}

	sort.Slice(edgeList, func(i, j int) bool { return edgeList[i][2] < edgeList[j][2] })
	sortedEdges := make([]Edge, len(edgeList))
	for i, e := range edgeList {
		sortedEdges[i] = Edge{e[0], e[1], e[2]}
	}
	forest, roots, values := KruskalTree(n, sortedEdges)
	tree := _NewTree(forest, roots)
	treeMonoid := NewStaticTreeMonoid(tree, values)
	return DistanceLimitedPathsExist{uf: uf, tree: tree, treeMonoid: treeMonoid}
}

func (this *DistanceLimitedPathsExist) Query(p int, q int, limit int) bool {
	if this.uf.Find(p) != this.uf.Find(q) {
		return false
	}
	return this.treeMonoid.QueryPath(p, q) < limit
}

// https://yukicoder.me/problems/no/1451
// 初始时有n个人
// 给定m个操作,每次将i和j所在的班级合并,大小相等时随机选取班长,否则选取较大的班级的班长作为新班级的班长
// 对每个人,问最后成为班长的概率
func yuki1451() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7
	pow := func(base, exp, mod int) int {
		base %= mod
		res := 1 % mod
		for ; exp > 0; exp >>= 1 {
			if exp&1 == 1 {
				res = res * base % mod
			}
			base = base * base % mod
		}
		return res
	}

	var INV2 = pow(2, MOD-2, MOD)

	var n, m int
	fmt.Fscan(in, &n, &m)
	sortedEdges := make([]Edge, 0, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		sortedEdges = append(sortedEdges, Edge{u, v, 1})
	}

	forest, roots, _ := KruskalTree(n, sortedEdges)
	subSize := make([]int, len(forest))
	var getSubSize func(int) int
	getSubSize = func(cur int) int {
		if cur < n { // 原始顶点
			subSize[cur] = 1
		}
		for _, to := range forest[cur] {
			subSize[cur] += getSubSize(to)
		}
		return subSize[cur]
	}
	for _, root := range roots {
		getSubSize(root)
	}

	res := make([]int, n)
	var run func(int, int)
	run = func(cur, p int) {
		if cur < n { // 原始顶点
			res[cur] = p
			return
		}

		if len(forest[cur]) == 1 { // 只有一个子节点
			run(forest[cur][0], p)
			return
		}

		left, right := forest[cur][0], forest[cur][1] // 两个子节点
		if subSize[left] > subSize[right] {
			run(left, p)
		} else if subSize[left] < subSize[right] {
			run(right, p)
		} else {
			run(left, p*INV2%MOD)
			run(right, p*INV2%MOD)
		}
	}

	for _, root := range roots {
		run(root, 1)
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

const INF int = 1e18

type E = int

const IS_COMMUTATIVE = true // 幺半群是否满足交换律
func e() E                  { return 0 }
func op(e1, e2 E) E         { return max(e1, e2) }

type Edge struct{ from, to, weight int }

// 表示合并过程的树,按照edges中边的顺序合并顶点.
//
//	返回:
//		forest: 森林的有向图邻接表
//		roots: 新图中的各个根节点
//		values: 每个辅助结点的权值(即对应边的权值，叶子结点权值为0).
func KruskalTree(n int, edges []Edge) (forest [][]int, roots []int, values []int) {
	parent := make([]int32, 2*n-1)
	for i := range parent {
		parent[i] = int32(i)
	}

	forest = make([][]int, 2*n-1)
	values = make([]int, 2*n-1)
	uf := NewUf(n)
	aux := int32(n)
	for i := range edges {
		e := &edges[i]
		from, to := e.from, e.to
		f := func(big, small int) {
			w, p1, p2 := e.weight, int(parent[big]), int(parent[small])
			forest[aux] = append(forest[aux], p1)
			forest[aux] = append(forest[aux], p2)
			parent[p1], parent[p2] = aux, aux
			parent[big], parent[small] = aux, aux
			values[aux] = w
			aux++
		}
		uf.Union(from, to, f)
	}

	forest = forest[:aux]
	values = values[:aux]
	for i := int32(0); i < aux; i++ {
		if parent[i] == i {
			roots = append(roots, int(i))
		}
	}
	return
}

type Uf struct {
	data []int32
}

func NewUf(n int) *Uf {
	data := make([]int32, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &Uf{data: data}
}

func (ufa *Uf) Union(key1, key2 int, f func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = int32(root1)
	if f != nil {
		f(root1, root2)
	}
	return true
}

func (ufa *Uf) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = int32(ufa.Find(int(ufa.data[key])))
	return int(ufa.data[key])
}

type StaticTreeMonoid struct {
	tree *_Tree
	n    int
	unit E
	seg  *DisjointSparseTable
	segR *DisjointSparseTable
}

// 静态树的路径查询, 维护的量需要满足幺半群的性质.
//
//	data: 顶点的值.
func NewStaticTreeMonoid(tree *_Tree, data []E) *StaticTreeMonoid {
	n := len(tree.Tree)
	res := &StaticTreeMonoid{tree: tree, n: n, unit: e()}
	leaves := make([]E, n)
	for v := range leaves {
		leaves[tree.LID[v]] = data[v]
	}
	res.seg = NewDisjointSparse(leaves, e, op)
	if !IS_COMMUTATIVE {
		res.segR = NewDisjointSparse(leaves, e, func(e1, e2 E) E { return op(e2, e1) }) // opRev
	}
	return res
}

// 查询 start 到 target 的路径上的值.(点权/边权 由 isVertex 决定)
func (st *StaticTreeMonoid) QueryPath(start, target int) E {
	path := st.tree.GetPathDecomposition(start, target)
	val := st.unit
	for _, ab := range path {
		a, b := int(ab[0]), int(ab[1])
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
	if !check(st.QueryPath(start, start)) {
		return -1
	}
	path := st.tree.GetPathDecomposition(start, target)
	val := st.unit
	for _, ab := range path {
		a, b := int(ab[0]), int(ab[1])
		x := st._getProd(a, b)
		if tmp := op(val, x); check(tmp) {
			val = tmp
			start = int(st.tree.IdToNode[b])
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
			return int(st.tree.IdToNode[i-1])
		} else {
			var i int
			if IS_COMMUTATIVE {
				i = st.seg.MinLeft(a+1, checkTmp)
			} else {
				i = st.segR.MinLeft(a+1, checkTmp)
			}
			if i == a+1 {
				return start
			}
			return int(st.tree.IdToNode[i])
		}
	}
	return target
}

func (st *StaticTreeMonoid) QuerySubtree(root int) E {
	l, r := st.tree.LID[root], st.tree.RID[root]
	return st.seg.Query(int(l), int(r))
}

func (st *StaticTreeMonoid) _getProd(a, b int) E {
	if IS_COMMUTATIVE {
		if a <= b {
			return st.seg.Query(a, b+1)
		}
		return st.seg.Query(b, a+1)
	} else {
		if a <= b {
			return st.seg.Query(a, b+1)
		}
		return st.segR.Query(b, a+1)
	}
}

type _Tree struct {
	Tree          [][]int
	Depth         []int32
	Parent        []int32
	LID, RID      []int32 // 欧拉序[in,out)
	IdToNode      []int32
	top, heavySon []int32
	timer         int32
}

func _NewTree(tree [][]int, roots []int) *_Tree {
	n := len(tree)
	lid := make([]int32, n)
	rid := make([]int32, n)
	IdToNode := make([]int32, n)
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	res := &_Tree{
		Tree:     tree,
		Depth:    depth,
		Parent:   parent,
		LID:      lid,
		RID:      rid,
		IdToNode: IdToNode,
		top:      top,
		heavySon: heavySon,
	}
	res._build(roots)
	return res
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *_Tree) Id(root int) (int, int) {
	return int(tree.LID[root]), int(tree.RID[root])
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *_Tree) Eid(u, v int) int {
	if tree.LID[u] > tree.LID[v] {
		return int(tree.LID[u])
	}
	return int(tree.LID[v])
}

func (tree *_Tree) LCA(u, v int) int {
	for {
		if tree.LID[u] > tree.LID[v] {
			u, v = v, u
		}
		if tree.top[u] == tree.top[v] {
			return u
		}
		v = int(tree.Parent[tree.top[v]])
	}
}

func (tree *_Tree) RootedLCA(u, v int, root int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *_Tree) Dist(u, v int) int {
	return int(tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)])
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
func (tree *_Tree) KthAncestor(root, k int) int {
	root32 := int32(root)
	k32 := int32(k)
	if k32 > tree.Depth[root32] {
		return -1
	}
	for {
		u := tree.top[root32]
		if tree.LID[root32]-k32 >= tree.LID[u] {
			return int(tree.IdToNode[tree.LID[root32]-k32])
		}
		k32 -= tree.LID[root32] - tree.LID[u] + 1
		root32 = tree.Parent[u]
	}
}

// 从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed)
//
//	返回跳到的节点,如果不存在这样的节点,返回-1
func (tree *_Tree) Jump(from, to, step int) int {
	step32 := int32(step)
	if step32 == 1 {
		if from == to {
			return -1
		}
		if tree.IsInSubtree(to, from) {
			return tree.KthAncestor(to, int(tree.Depth[to]-tree.Depth[from]-1))
		}
		return int(tree.Parent[from])
	}
	c := tree.LCA(from, to)
	dac := tree.Depth[from] - tree.Depth[c]
	dbc := tree.Depth[to] - tree.Depth[c]
	if step32 > dac+dbc {
		return -1
	}
	if step32 <= dac {
		return tree.KthAncestor(from, int(step32))
	}
	return tree.KthAncestor(to, int(dac+dbc-step32))
}

func (tree *_Tree) CollectChild(root int) []int {
	res := []int{}
	for _, next := range tree.Tree[root] {
		if next != int(tree.Parent[root]) {
			res = append(res, next)
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *_Tree) GetPathDecomposition(u, v int) [][2]int32 {
	u32, v32 := int32(u), int32(v)
	up, down := [][2]int32{}, [][2]int32{}
	for {
		if tree.top[u32] == tree.top[v32] {
			break
		}
		if tree.LID[u32] < tree.LID[v32] {
			down = append(down, [2]int32{tree.LID[tree.top[v32]], tree.LID[v32]})
			v32 = tree.Parent[tree.top[v32]]
		} else {
			up = append(up, [2]int32{tree.LID[u32], tree.LID[tree.top[u32]]})
			u32 = tree.Parent[tree.top[u32]]
		}
	}
	if tree.LID[u32] < tree.LID[v32] {
		down = append(down, [2]int32{tree.LID[u32], tree.LID[v32]})
	} else if tree.LID[v32] <= tree.LID[u32] {
		up = append(up, [2]int32{tree.LID[u32], tree.LID[v32]})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

func (tree *_Tree) EnumeratePathDecomposition(u, v int, f func(a, b int)) {
	u32, v32 := int32(u), int32(v)
	down := [][2]int32{}
	for {
		if tree.top[u32] == tree.top[v32] {
			break
		}
		if tree.LID[u32] < tree.LID[v32] {
			down = append(down, [2]int32{tree.LID[tree.top[v32]], tree.LID[v32]})
			v32 = tree.Parent[tree.top[v32]]
		} else {
			f(int(tree.LID[u32]), int(tree.LID[tree.top[u32]]))
			u32 = tree.Parent[tree.top[u32]]
		}
	}
	if tree.LID[u32] < tree.LID[v32] {
		down = append(down, [2]int32{tree.LID[u32], tree.LID[v32]})
	} else if tree.LID[v32] <= tree.LID[u32] {
		f(int(tree.LID[u32]), int(tree.LID[v32]))
	}
	for i := len(down) - 1; i >= 0; i-- {
		f(int(down[i][0]), int(down[i][1]))
	}
}

func (tree *_Tree) GetPath(u, v int) []int {
	res := []int{}
	composition := tree.GetPathDecomposition(u, v)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, int(tree.IdToNode[i]))
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, int(tree.IdToNode[i]))
			}
		}
	}
	return res
}

// 以root为根时,结点v的子树大小.
func (tree *_Tree) SubtreeSize(v, root int) int {
	if root == -1 {
		return int(tree.RID[v] - tree.LID[v])
	}
	if v == root {
		return len(tree.Tree)
	}
	x := tree.Jump(v, root, 1)
	if tree.IsInSubtree(v, x) {
		return int(tree.RID[v] - tree.LID[v])
	}
	return len(tree.Tree) - int(tree.RID[x]) + int(tree.LID[x])
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *_Tree) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

func (tree *_Tree) ELID(u int) int {
	return int(2*tree.LID[u] - tree.Depth[u])
}

func (tree *_Tree) ERID(u int) int {
	return int(2*tree.RID[u] - tree.Depth[u] - 1)
}

func (tree *_Tree) build(cur, pre, dep int32) int {
	subSize, heavySize, heavySon := 1, 0, int32(-1)
	for _, next := range tree.Tree[cur] {
		next32 := int32(next)
		if next32 != pre {
			nextSize := tree.build(next32, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next32
			}
		}
	}
	tree.Depth[cur] = dep
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *_Tree) markTop(cur, top int32) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	if tree.heavySon[cur] != -1 {
		tree.markTop(tree.heavySon[cur], top)
		for _, next := range tree.Tree[cur] {
			next32 := int32(next)
			if next32 != tree.heavySon[cur] && next32 != tree.Parent[cur] {
				tree.markTop(next32, next32)
			}
		}
	}
	tree.RID[cur] = tree.timer
}

// root:0-based
func (tree *_Tree) _build(roots []int) {
	for _, root := range roots {
		root32 := int32(root)
		tree.build(root32, -1, 0)
		tree.markTop(root32, root32)
	}
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

// 给定非负整数数组 nums 构建一个 WaveletMatrix.
func NewWaveletMatrix(data []int) *WaveletMatrix {
	dataCopy := make([]int, len(data))
	max_ := 0
	for i, v := range data {
		if v > max_ {
			max_ = v
		}
		dataCopy[i] = v
	}
	maxLog := bits.Len(uint(max_)) + 1
	n := len(dataCopy)
	mat := make([]*BitVector, maxLog)
	zs := make([]int, maxLog)
	buff1 := make([]int, maxLog)
	buff2 := make([]int, maxLog)

	ls, rs := make([]int, n), make([]int, n)
	for dep := 0; dep < maxLog; dep++ {
		mat[dep] = NewBitVector(n + 1)
		p, q := 0, 0
		for i := 0; i < n; i++ {
			k := (dataCopy[i] >> uint(maxLog-dep-1)) & 1
			if k == 1 {
				rs[q] = dataCopy[i]
				mat[dep].Set(i)
				q++
			} else {
				ls[p] = dataCopy[i]
				p++
			}
		}

		zs[dep] = p
		mat[dep].Build()
		ls = dataCopy
		for i := 0; i < q; i++ {
			dataCopy[p+i] = rs[i]
		}
	}

	return &WaveletMatrix{
		n:      n,
		maxLog: maxLog,
		mat:    mat,
		zs:     zs,
		buff1:  buff1,
		buff2:  buff2,
	}
}

type WaveletMatrix struct {
	n            int
	maxLog       int
	mat          []*BitVector
	zs           []int
	buff1, buff2 []int
}

// [start, end) 内的 value 的個数.
func (w *WaveletMatrix) Count(start, end, value int) int {
	return w.count(value, end) - w.count(value, start)
}

// [start, end) 内 [lower, upper) 的个数.
func (w *WaveletMatrix) CountRange(start, end, lower, upper int) int {
	return w.freqDfs(0, start, end, 0, lower, upper)
}

// 第k(0-indexed)个value的位置.
func (w *WaveletMatrix) Index(value, k int) int {
	w.count(value, w.n)
	for dep := w.maxLog - 1; dep >= 0; dep-- {
		bit := (value >> uint(w.maxLog-dep-1)) & 1
		k = w.mat[dep].IndexWithStart(bit, k, w.buff1[dep])
		if k < 0 || k >= w.buff2[dep] {
			return -1
		}
		k -= w.buff1[dep]
	}
	return k
}

func (w *WaveletMatrix) IndexWithStart(value, k, start int) int {
	return w.Index(value, k+w.count(value, start))
}

// [start, end) 内第k(0-indexed)小的数(kth).
func (w *WaveletMatrix) Kth(start, end, k int) int {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) 内第k(0-indexed)大的数.
func (w *WaveletMatrix) KthMax(start, end, k int) int {
	if k < 0 || k >= end-start {
		return -1
	}
	res := 0
	for dep := 0; dep < w.maxLog; dep++ {
		p, q := w.mat[dep].Count(1, start), w.mat[dep].Count(1, end)
		if k < q-p {
			start = w.zs[dep] + p
			end = w.zs[dep] + q
			res |= 1 << uint(w.maxLog-dep-1)
		} else {
			k -= q - p
			start -= p
			end -= q
		}
	}
	return res
}

// [start, end) 内第k(0-indexed)小的数(kth).
func (w *WaveletMatrix) KthMin(start, end, k int) int {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) 中比 value 严格小的数, 不存在返回 -INF.
func (w *WaveletMatrix) Lower(start, end, value int) int {
	k := w.lt(start, end, value)
	if k != 0 {
		return w.KthMin(start, end, k-1)
	}
	return -INF
}

// [start, end) 中比 value 严格大的数, 不存在返回 INF.
func (w *WaveletMatrix) Higher(start, end, value int) int {
	k := w.le(start, end, value)
	if k == end-start {
		return INF
	}
	return w.KthMin(start, end, k)
}

// [start, end) 中不超过 value 的最大值, 不存在返回 -INF.
func (w *WaveletMatrix) Floor(start, end, value int) int {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Lower(start, end, value)
}

// [start, end) 中不小于 value 的最小值, 不存在返回 INF.
func (w *WaveletMatrix) Ceiling(start, end, value int) int {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Higher(start, end, value)
}

func (w *WaveletMatrix) access(k int) int {
	res := 0
	for dep := 0; dep < w.maxLog; dep++ {
		bit := w.mat[dep].Get(k)
		res = (res << 1) | bit
		k = w.mat[dep].Count(bit, k) + w.zs[dep]*dep
	}
	return res
}

func (w *WaveletMatrix) count(value, end int) int {
	left, right := 0, end
	for dep := 0; dep < w.maxLog; dep++ {
		w.buff1[dep] = left
		w.buff2[dep] = right
		bit := (value >> uint(w.maxLog-dep-1)) & 1
		left = w.mat[dep].Count(bit, left) + w.zs[dep]*bit
		right = w.mat[dep].Count(bit, right) + w.zs[dep]*bit
	}
	return right - left
}

func (w *WaveletMatrix) freqDfs(d, left, right, val, a, b int) int {
	if left == right {
		return 0
	}
	if d == w.maxLog {
		if a <= val && val < b {
			return right - left
		}
		return 0
	}

	nv := (1 << uint(w.maxLog-d-1)) | val
	nnv := ((1 << uint(w.maxLog-d-1)) - 1) | nv
	if nnv < a || b <= val {
		return 0
	}
	if a <= val && nnv < b {
		return right - left
	}
	lc, rc := w.mat[d].Count(1, left), w.mat[d].Count(1, right)
	return w.freqDfs(d+1, left-lc, right-rc, val, a, b) + w.freqDfs(d+1, lc+w.zs[d], rc+w.zs[d], nv, a, b)
}

func (w *WaveletMatrix) ll(left, right, v int) (int, int) {
	res := 0
	for dep := 0; dep < w.maxLog; dep++ {
		w.buff1[dep] = left
		w.buff2[dep] = right
		bit := v >> uint(w.maxLog-dep-1) & 1
		if bit == 1 {
			res += right - left + w.mat[dep].Count(1, left) - w.mat[dep].Count(1, right)
		}
		left = w.mat[dep].Count(bit, left) + w.zs[dep]*bit
		right = w.mat[dep].Count(bit, right) + w.zs[dep]*bit
	}
	return res, right - left
}

func (w *WaveletMatrix) lt(left, right, v int) int {
	a, _ := w.ll(left, right, v)
	return a
}

func (w *WaveletMatrix) le(left, right, v int) int {
	a, b := w.ll(left, right, v)
	return a + b
}

type BitVector struct {
	n     int
	block []int
	sum   []int
}

func NewBitVector(n int) *BitVector {
	blockCount := (n + 63) >> 6
	return &BitVector{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int, blockCount+1),
	}
}

func (f *BitVector) Set(i int) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *BitVector) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + bits.OnesCount(uint(f.block[i]))
	}
}

func (f *BitVector) Get(i int) int {
	return (f.block[i>>6] >> uint(i&63)) & 1
}

func (f *BitVector) Count(value, end int) int {
	mask := (1 << uint(end&63)) - 1
	res := f.sum[end>>6] + bits.OnesCount(uint(f.block[end>>6]&mask))
	if value == 1 {
		return res
	}
	return end - res
}

func (f *BitVector) Index(value, k int) int {
	if k < 0 || f.Count(value, f.n) <= k {
		return -1
	}

	left, right := 0, f.n
	for right-left > 1 {
		mid := (left + right) >> 1
		if f.Count(value, mid) >= k+1 {
			right = mid
		} else {
			left = mid
		}
	}
	return right - 1
}

func (f *BitVector) IndexWithStart(value, k, start int) int {
	return f.Index(value, k+f.Count(value, start))
}

// 区间LCA.
//
//	points 顶点数组.
//	getLCA LCA实现.
//	返回一个查询 points[start:end) lca 的函数.
func RangeLCA(
	points []int,
	getLCA func(u, v int) int,
) func(start, end int) int {
	n := 1
	for n < len(points) {
		n <<= 1
	}
	seg := make([]int32, n<<1)
	for i := 0; i < len(points); i++ {
		seg[n+i] = int32(points[i])
	}
	for i := n - 1; i >= 0; i-- {
		seg[i] = int32(getLCA(int(seg[i<<1]), int(seg[(i<<1)|1])))
	}
	lca := func(u, v int32) int32 {
		if u == -1 || v == -1 {
			if u == -1 {
				return v
			}
			return u
		}
		return int32(getLCA(int(u), int(v)))
	}
	query := func(start, end int) int {
		res := int32(-1)
		for ; start > 0 && start+(start&-start) <= end; start += start & -start {
			res = lca(res, seg[(n+start)/(start&-start)])
		}
		for ; start < end; end -= end & -end {
			res = lca(res, seg[(n+end)/(end&-end)-1])
		}
		return int(res)
	}
	return query
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
