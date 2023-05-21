package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"math/bits"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
)

func init() {
	debug.SetGCPercent(-1)
}

// from https://atcoder.jp/users/ccppjsrb
var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()
	n := io.NextInt()
	values1 := make([]int, n+10) // 顶点权值
	values2 := make([]int, n+10)
	for i := 0; i < n; i++ {
		a, b := io.NextInt(), io.NextInt()
		values1[i] = a
		values2[i] = b
	}
	lca := NewTree(n)
	tree := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		u, v := io.NextInt()-1, io.NextInt()-1
		tree[u] = append(tree[u], Edge{v, 1})
		tree[v] = append(tree[v], Edge{u, 1})
		lca.AddEdge(u, v, 1)
	}
	lca.Build(0)
	queries := make([][2]int, n-1) // u, v
	for i := 0; i < n-1; i++ {
		u := 0
		v := i + 1
		queries[i] = [2]int{u, v}
	}

	mo := NewMoOnTree(tree, 0)
	for _, q := range queries {
		mo.AddQuery(q[0], q[1])
	}

	res := make([]int, n-1)
	counter, kind := make(map[int]int, n), 0
	overlap := 0
	add := func(i int) {
		counter[values1[i]]++
		counter[values2[i]]++
		if counter[values1[i]] == 1 {
			kind++
		}
		if counter[values2[i]] == 1 {
			kind++
		}
		if values1[i] == values2[i] {
			overlap++
		}
	}
	remove := func(i int) {
		counter[values1[i]]--
		if counter[values1[i]] == 0 {
			kind--
		}
		counter[values2[i]]--
		if counter[values2[i]] == 0 {
			kind--
		}
		if values1[i] == values2[i] {
			overlap--
		}
	}
	query := func(qid int) {
		u, v := queries[qid][0], queries[qid][1]
		res[qid] = mins(kind-overlap+1, kind, lca.Dist(u, v, false)+1)
	}

	mo.Run(add, remove, query)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

type MoOnTree struct {
	tree    [][]Edge
	root    int
	in, vs  []int
	queries [][2]int
}

type Edge struct{ to, weight int }

func NewMoOnTree(tree [][]Edge, root int) *MoOnTree {
	return &MoOnTree{tree: tree, root: root}
}

// 添加从顶点u到顶点v的查询.
func (mo *MoOnTree) AddQuery(u, v int) { mo.queries = append(mo.queries, [2]int{u, v}) }

// 处理每个查询.
//  add: 将数据添加到窗口.
//  remove: 将数据从窗口移除.
//  query: 查询窗口内的数据.
func (mo *MoOnTree) Run(add func(rootId int), remove func(rootId int), query func(qid int)) {
	n := len(mo.tree)
	q := len(mo.queries)

	vs := make([]int, 0, 2*n)
	tin := make([]int, n)
	tout := make([]int, n)

	var initTime func(v, fa int)
	initTime = func(v, fa int) {
		tin[v] = len(vs)
		vs = append(vs, v)
		for _, w := range mo.tree[v] {
			to := w.to
			if to != fa {
				initTime(to, v)
			}
		}
		tout[v] = len(vs)
		vs = append(vs, v)
	}
	initTime(mo.root, -1)
	lca := offlineLCA(mo.tree, mo.queries, mo.root)
	blockSize := int(math.Ceil(float64(2*n) / math.Sqrt(float64(q)))) // int(math.Round(math.Pow(float64(2*n), 2.0/3)))
	type Q struct{ lb, l, r, lca, qid int }
	qs := make([]Q, len(mo.queries))
	for i := range qs {
		v, w := mo.queries[i][0], mo.queries[i][1]
		if tin[v] > tin[w] {
			v, w = w, v
		}
		if lca_ := lca[i]; lca_ != v {
			qs[i] = Q{tout[v] / blockSize, tout[v], tin[w] + 1, lca_, i}
		} else {
			qs[i] = Q{tin[v] / blockSize, tin[v], tin[w] + 1, -1, i}
		}
	}

	sort.Slice(qs, func(i, j int) bool {
		a, b := qs[i], qs[j]
		if a.lb != b.lb {
			return a.lb < b.lb
		}
		if a.lb&1 == 0 {
			return a.r < b.r
		}
		return a.r > b.r
	})

	flip := make([]bool, n)
	f := func(u int) {
		flip[u] = !flip[u]
		if flip[u] {
			add(u)
		} else {
			remove(u)
		}
	}

	l, r := 0, 0
	for _, q := range qs {
		for ; r < q.r; r++ {
			f(vs[r])
		}
		for ; l < q.l; l++ {
			f(vs[l])
		}
		for l > q.l {
			l--
			f(vs[l])
		}
		for r > q.r {
			r--
			f(vs[r])
		}
		if q.lca >= 0 {
			f(q.lca)
		}
		query(q.qid)
		if q.lca >= 0 {
			f(q.lca)
		}
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

// LCA离线.
func offlineLCA(graph [][]Edge, queries [][2]int, root int) []int {
	n := len(graph)
	ufa := newUnionFindArray(n)
	st, mark, ptr, res := make([]int, n), make([]int, n), make([]int, n), make([]int, len(queries))
	for i := 0; i < len(queries); i++ {
		res[i] = -1
	}
	top := 0
	st[top] = root
	for _, q := range queries {
		mark[q[0]]++
		mark[q[1]]++
	}
	q := make([][][2]int, n)
	for i := 0; i < n; i++ {
		q[i] = make([][2]int, 0, mark[i])
		mark[i] = -1
		ptr[i] = len(graph[i])
	}
	for i := range queries {
		q[queries[i][0]] = append(q[queries[i][0]], [2]int{queries[i][1], i})
		q[queries[i][1]] = append(q[queries[i][1]], [2]int{queries[i][0], i})
	}

	run := func(u int) bool {
		for ptr[u] != 0 {
			v := graph[u][ptr[u]-1].to
			ptr[u]--
			if mark[v] == -1 {
				top++
				st[top] = v
				return true
			}
		}
		return false
	}

	for top != -1 {
		u := st[top]
		if mark[u] == -1 {
			mark[u] = u
		} else {
			ufa.Union(u, graph[u][ptr[u]].to)
			mark[ufa.Find(u)] = u
		}

		if !run(u) {
			for _, v := range q[u] {
				if mark[v[0]] != -1 && res[v[1]] == -1 {
					res[v[1]] = mark[ufa.Find(v[0])]
				}
			}
			top--
		}
	}

	return res
}

func newUnionFindArray(n int) *unionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &unionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type unionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *unionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *unionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *unionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *unionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *unionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *unionFindArray) String() string {
	sb := []string{"unionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}

type BitArray struct {
	n    int
	log  int
	data []int
}

// 長さ n の 0で初期化された配列で構築する.
func NewBitArray(n int) *BitArray {
	return &BitArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

// 配列で構築する.
func NewBitArrayFrom(arr []int) *BitArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BitArray) Build(arr []int) {
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

// 要素 i に値 v を加える.
func (b *BitArray) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r) の要素の総和を求める.
func (b *BitArray) Prod(r int) int {
	res := int(0)
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BitArray) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) LowerBound(x int) int {
	i := 0
	for k := 1 << uint(b.log); k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) UpperBound(x int) int {
	i := 0
	for k := 1 << uint(b.log); k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

func (b *BitArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.ProdRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
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
