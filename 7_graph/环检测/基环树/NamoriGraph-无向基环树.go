// https://maspypy.github.io/library/graph/unicyclic.hpp
// !Namori Graph 无向图基环树
// !一个有 n 个顶点和 n 条边的“连通”无向图。图中只有一个环。
// 这个图被称为“Namori图”，这是根据一位漫画家的名字而命名的，
// 但在学术界中，正确的称呼是“单环图(Unicyclic)”或“伪森林(Pseudoforest)”。

// 例如：下图是一个无向基环树，其中的边权都是1。
//
//	  0
//	  |
//	  1
//	 / \
//	2   3 - 4 - 5
//	 \ /
//	  6
//
// root: 3
// outEdge: 3-6
// to: [1 3 1 6 3 4 2]
// cycle: [6 2 1 3] (bottom到root的路径)
// directedGraph:
//
//	    3(root)
//	   / \
//	  1   4
//	 / \   \
//	0   2   5
//	     \
//	      6(bottom)
//
// !性质1：点u在所在子树的根节点(在环上)为lca(u, bottom).
//
// !这种维护方法的优势是支持动态修改点权或者边权

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// yuki1254()
	// abc266f()
	namoriCut()
}

// https://yukicoder.me/problems/no/1254
// !找基环树中的环上的边.
func yuki1254() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	edges := make([]edge, n)
	for i := int32(0); i < n; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = edge{u: u, v: v, w: 1}
	}

	namori := NewNamoriGraph(n, edges)
	res := []int32{}
	for i, e := range edges {
		if namori.InCycle[e.u] && namori.InCycle[e.v] {
			res = append(res, int32(i+1))
		}
	}

	fmt.Fprintln(out, len(res))
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// https://atcoder.jp/contests/abc266/tasks/abc266_f
// 每次查询基环树中两个点是否由唯一的路径相连.
// !等价于：所在子树的根节点是否相同.
func abc266f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	edges := make([]edge, n)
	for i := int32(0); i < n; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = edge{u, v, 1}
	}

	var q int32
	fmt.Fscan(in, &q)
	queries := make([][2]int32, q)
	for i := int32(0); i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		queries[i] = [2]int32{u, v}
	}

	namori := NewNamoriGraph(n, edges)
	_, tree := namori.BuildTree()
	root := namori.Root
	bottom := namori.To[root]
	for _, q := range queries {
		u, v := q[0], q[1]
		lca1, lca2 := tree.LCA(u, bottom), tree.LCA(v, bottom)
		if lca1 == lca2 {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

func namoriCut() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2891
	// 给定一个基环树 q个询问(x,y) ， x!=y
	// 求使得x和y不连通最少需要切掉多少条边
	// !如果两个点都在环上,则答案为2；否则答案为1(两点之间只有唯一的一条路径)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	edges := make([]edge, n)
	for i := int32(0); i < n; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		edges[i] = edge{a, b, 1}
	}

	G := NewNamoriGraph(n, edges)

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var x, y int32
		fmt.Fscan(in, &x, &y)
		x, y = x-1, y-1
		if G.InCycle[x] && G.InCycle[y] {
			fmt.Fprintln(out, 2)
		} else {
			fmt.Fprintln(out, 1)
		}
	}
}

func demo() {
	edges := []edge{
		{0, 1, 1}, {1, 2, 2}, {1, 3, 3}, {2, 6, 6}, {3, 6, 6}, {3, 4, 4}, {4, 5, 5},
	}
	namori := NewNamoriGraph(7, edges)
	fmt.Println(namori.Root)
	fmt.Println(namori.OutEdgeId)
	fmt.Println(namori.OutCost)
	fmt.Println(namori.To)
	fmt.Println(namori.Cycle)

	directedGraph, tree := namori.BuildTree()
	fmt.Println(directedGraph)
	fmt.Println(namori.Dist(tree, 0, 5))
}

type edge = struct {
	u, v int32
	w    int
}

type neighbor = struct {
	to, eid int32
	weight  int
}

// !无向基环树.
type NamoriGraph struct {
	RawEdges []edge
	RawGraph [][]neighbor

	N         int32
	Root      int32 // 断开outEdge后有向树的根
	OutEdgeId int32 // build后不在树中的边
	OutCost   int
	To        []int32 // !向Root方向移动1步后的结点，Root的To为对应outEdge的另一端
	Cycle     []int32
	InCycle   []bool
}

func NewNamoriGraph(n int32, edges []edge) *NamoriGraph {
	m := int32(len(edges))
	if m != n {
		panic("invalid namori graph")
	}

	graph := make([][]neighbor, n)
	for eid := int32(0); eid < m; eid++ {
		e := &edges[eid]
		u, v, w := e.u, e.v, e.w
		graph[u] = append(graph[u], neighbor{to: v, weight: w, eid: eid})
		graph[v] = append(graph[v], neighbor{to: u, weight: w, eid: eid})
	}

	uf := newUnionFindArraySimple32(n)
	to := make([]int32, n)
	for i := range to {
		to[i] = -1
	}

	root := int32(-1)
	outEdgeId, outCost := int32(-1), -1
	for eid := int32(0); eid < m; eid++ {
		e := &edges[eid]
		u, v, w := e.u, e.v, e.w
		if uf.Union(u, v) {
			continue
		}
		outEdgeId, outCost = eid, w
		root = u
		to[root] = v
		break
	}
	visited := make([]bool, n)
	stack := []int32{root}
	for len(stack) > 0 {
		pre := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		visited[pre] = true
		for _, e := range graph[pre] {
			next, eid := e.to, e.eid
			if visited[next] || eid == outEdgeId {
				continue
			}
			to[next] = pre
			stack = append(stack, next)
		}
	}

	cycle := []int32{to[root]}
	for cycle[len(cycle)-1] != root {
		cycle = append(cycle, to[cycle[len(cycle)-1]])
	}

	inCycle := make([]bool, n)
	for _, v := range cycle {
		inCycle[v] = true
	}

	return &NamoriGraph{
		RawEdges:  edges,
		RawGraph:  graph,
		N:         n,
		Root:      root,
		OutEdgeId: outEdgeId,
		OutCost:   outCost,
		To:        to,
		Cycle:     cycle,
		InCycle:   inCycle,
	}
}

// 断开outEdge, 生成有向树.
func (ng *NamoriGraph) BuildTree() (directedGraph [][]neighbor, tree *Tree) {
	directedGraph = make([][]neighbor, ng.N)
	for eid := int32(0); eid < ng.N; eid++ {
		if eid == ng.OutEdgeId {
			continue
		}
		e := &ng.RawEdges[eid]
		u, v, w := e.u, e.v, e.w
		if ng.To[u] == v {
			u, v = v, u
		}
		directedGraph[u] = append(directedGraph[u], neighbor{to: v, weight: w, eid: eid})
	}
	tree = NewTree(directedGraph)
	tree.Build(ng.Root)
	return
}

// 基环树求距离.
func (ng *NamoriGraph) Dist(tree *Tree, u, v int32) int32 {
	bottom := ng.To[ng.Root]
	// lca为在环上的点
	lca1, lca2 := tree.LCA(u, bottom), tree.LCA(v, bottom)
	distOnCyle := abs32(tree.Depth[lca1] - tree.Depth[lca2])
	distOnCyle = min32(distOnCyle, int32(len(ng.Cycle))-distOnCyle)
	return distOnCyle + tree.Depth[u] + tree.Depth[v] - tree.Depth[lca1] - tree.Depth[lca2]
}

func (ng *NamoriGraph) DistWeighted(tree *Tree, u, v int32) int {
	bottom := ng.To[ng.Root]
	lca1, lca2 := tree.LCA(u, bottom), tree.LCA(v, bottom)
	distOnCycle := abs(tree.DepthWeighted[lca1] - tree.DepthWeighted[lca2])
	distOnCycle = min(distOnCycle, tree.DepthWeighted[bottom]+ng.OutCost-distOnCycle)
	return distOnCycle + tree.DepthWeighted[u] + tree.DepthWeighted[v] - tree.DepthWeighted[lca1] - tree.DepthWeighted[lca2]
}

type unionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func newUnionFindArraySimple32(n int32) *unionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &unionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *unionFindArraySimple32) Union(key1, key2 int32) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *unionFindArraySimple32) Find(key int32) int32 {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = u.Find(u.data[key])
	return u.data[key]
}

func (u *unionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}

type Tree struct {
	Tree          [][]neighbor // (next, weight)
	Depth         []int32
	DepthWeighted []int
	Parent        []int32
	LID, RID      []int32 // 欧拉序[in,out)
	IdToNode      []int32
	top, heavySon []int32
	timer         int32
}

func NewTree(graph [][]neighbor) *Tree {
	n := int32(len(graph))
	tree := graph
	lid := make([]int32, n)
	rid := make([]int32, n)
	IdToNode := make([]int32, n)
	top := make([]int32, n)   // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n) // 深度
	depthWeighted := make([]int, n)
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
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
		IdToNode:      IdToNode,
		top:           top,
		heavySon:      heavySon,
	}
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *Tree) Build(root int32) {
	if root != -1 {
		tree.build(root, -1, 0, 0)
		tree.markTop(root, root)
	} else {
		for i := int32(0); i < int32(len(tree.Tree)); i++ {
			if tree.Parent[i] == -1 {
				tree.build(i, -1, 0, 0)
				tree.markTop(i, i)
			}
		}
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree) Id(root int32) (int32, int32) {
	return tree.LID[root], tree.RID[root]
}

func (tree *Tree) LCA(u, v int32) int32 {
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

func (tree *Tree) RootedLCA(u, v int32, w int32) int32 {
	return tree.LCA(u, v) ^ tree.LCA(u, w) ^ tree.LCA(v, w)
}
func (tree *Tree) RootedParent(u int32, root int32) int32 {
	return tree.Jump(u, root, 1)
}

func (tree *Tree) Dist(u, v int32, weighted bool) int {
	if weighted {
		return tree.DepthWeighted[u] + tree.DepthWeighted[v] - 2*tree.DepthWeighted[tree.LCA(u, v)]
	}
	return int(tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)])
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
//	kthAncestor(root,0) == root
func (tree *Tree) KthAncestor(root, k int32) int32 {
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
func (tree *Tree) Jump(from, to, step int32) int32 {
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

func (tree *Tree) CollectChild(root int32) []int32 {
	res := []int32{}
	for _, e := range tree.Tree[root] {
		next := e.to
		if next != tree.Parent[root] {
			res = append(res, next)
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree) GetPathDecomposition(u, v int32, vertex bool) [][2]int32 {
	up, down := [][2]int32{}, [][2]int32{}
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			down = append(down, [2]int32{tree.LID[tree.top[v]], tree.LID[v]})
			v = tree.Parent[tree.top[v]]
		} else {
			up = append(up, [2]int32{tree.LID[u], tree.LID[tree.top[u]]})
			u = tree.Parent[tree.top[u]]
		}
	}
	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}
	if tree.LID[u] < tree.LID[v] {
		down = append(down, [2]int32{tree.LID[u] + edgeInt, tree.LID[v]})
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		up = append(up, [2]int32{tree.LID[u], tree.LID[v] + edgeInt})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree) EnumeratePathDecomposition(u, v int32, vertex bool, f func(start, end int32)) {
	for {
		if tree.top[u] == tree.top[v] {
			break
		}
		if tree.LID[u] < tree.LID[v] {
			a, b := tree.LID[tree.top[v]], tree.LID[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = tree.Parent[tree.top[v]]
		} else {
			a, b := tree.LID[u], tree.LID[tree.top[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = tree.Parent[tree.top[u]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if tree.LID[u] < tree.LID[v] {
		a, b := tree.LID[u]+edgeInt, tree.LID[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if tree.LID[v]+edgeInt <= tree.LID[u] {
		a, b := tree.LID[u], tree.LID[v]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

func (tree *Tree) GetPath(u, v int32) []int32 {
	res := []int32{}
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
func (tree *Tree) SubSize(v, root int32) int32 {
	if root == -1 {
		return tree.RID[v] - tree.LID[v]
	}
	if v == root {
		return int32(len(tree.Tree))
	}
	x := tree.Jump(v, root, 1)
	if tree.IsInSubtree(v, x) {
		return tree.RID[v] - tree.LID[v]
	}
	return int32(len(tree.Tree)) - tree.RID[x] + tree.LID[x]
}

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *Tree) IsInSubtree(child, root int32) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链底端节点.
func (tree *Tree) GetHeavyPath(start int32) []int32 {
	heavyPath := []int32{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

// 结点v的重儿子.如果没有重儿子,返回-1.
func (tree *Tree) GetHeavyChild(v int32) int32 {
	k := tree.LID[v] + 1
	if k == int32(len(tree.Tree)) {
		return -1
	}
	w := tree.IdToNode[k]
	if tree.Parent[w] == v {
		return w
	}
	return -1
}

func (tree *Tree) ELID(u int32) int32 {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *Tree) ERID(u int32) int32 {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *Tree) build(cur, pre, dep int32, dist int) int32 {
	subSize, heavySize, heavySon := int32(1), int32(0), int32(-1)
	for _, e := range tree.Tree[cur] {
		next, weight := e.to, e.weight
		if next != pre {
			nextSize := tree.build(next, cur, dep+1, dist+int(weight))
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

func (tree *Tree) markTop(cur, top int32) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	heavySon := tree.heavySon[cur]
	if heavySon != -1 {
		tree.markTop(heavySon, top)
		for _, e := range tree.Tree[cur] {
			next := e.to
			if next != heavySon && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
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

func max(a, b int) int {
	if a > b {
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

func abs32(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
