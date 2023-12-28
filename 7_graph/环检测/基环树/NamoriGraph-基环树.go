// https://maspypy.github.io/library/graph/unicyclic.hpp
// Namori Graph 无向图基环树
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
	abc266f()
}

// https://yukicoder.me/problems/no/1254
func yuki1254() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	edges := make([]Edge, n)
	for i := 0; i < n; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = Edge{u, v, 1}
	}

	namori := NewNamoriGraph(n, edges)
	res := []int{}
	for i, e := range edges {
		if namori.InCycle[e[0]] && namori.InCycle[e[1]] {
			res = append(res, i+1)
		}
	}

	fmt.Fprintln(out, len(res))
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// https://atcoder.jp/contests/abc266/tasks/abc266_f
// 每次查询基环树中两个点是否由唯一的路径相连.
// 等价于：所在子树的根节点是否相同.
func abc266f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	edges := make([]Edge, n)
	for i := 0; i < n; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = Edge{u, v, 1}
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		queries[i] = [2]int{u, v}
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

func demo() {
	edges := []Edge{
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

type Edge = [3]int // (u,v,w)
type Neighbor = struct{ to, weight, eid int }

// !无向基环树.
type NamoriGraph struct {
	RawEdges []Edge
	RawGraph [][]Neighbor

	N         int
	Root      int // 断开outEdge后有向树的根
	OutEdgeId int // build后不在树中的边
	OutCost   int
	To        []int // !向Root方向移动1步后的结点，Root的To为对应outEdge的另一端
	Cycle     []int
	InCycle   []bool
}

func NewNamoriGraph(n int, edges []Edge) *NamoriGraph {
	m := len(edges)
	if m != n {
		panic("invalid namori graph")
	}

	graph := make([][]Neighbor, n)
	for eid := 0; eid < m; eid++ {
		e := &edges[eid]
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], Neighbor{to: v, weight: w, eid: eid})
		graph[v] = append(graph[v], Neighbor{to: u, weight: w, eid: eid})
	}

	uf := NewUf(n)
	to := make([]int, n)
	for i := range to {
		to[i] = -1
	}

	root := -1
	outEdgeId, outCost := -1, -1
	for eid := 0; eid < m; eid++ {
		e := &edges[eid]
		u, v, w := e[0], e[1], e[2]
		if uf.Union(u, v) {
			continue
		}
		outEdgeId, outCost = eid, w
		root = u
		to[root] = v
		break
	}
	visited := make([]bool, n)
	stack := []int{root}
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

	cycle := []int{to[root]}
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
func (ng *NamoriGraph) BuildTree() (directedGraph [][]Neighbor, tree *Tree) {
	directedGraph = make([][]Neighbor, ng.N)
	for eid := 0; eid < ng.N; eid++ {
		if eid == ng.OutEdgeId {
			continue
		}
		e := &ng.RawEdges[eid]
		u, v, w := e[0], e[1], e[2]
		if ng.To[u] == v {
			u, v = v, u
		}
		directedGraph[u] = append(directedGraph[u], Neighbor{to: v, weight: w, eid: eid})
	}
	tree = NewTree(directedGraph, ng.Root)
	return
}

func (ng *NamoriGraph) Dist(tree *Tree, u, v int) int {
	bottom := ng.To[ng.Root]
	// lca为在环上的点
	lca1, lca2 := tree.LCA(u, bottom), tree.LCA(v, bottom)
	distOnCyle := abs(tree.Depth[lca1] - tree.Depth[lca2])
	distOnCyle = min(distOnCyle, len(ng.Cycle)-distOnCyle)
	return distOnCyle + tree.Depth[u] + tree.Depth[v] - tree.Depth[lca1] - tree.Depth[lca2]
}

func (ng *NamoriGraph) DistWeighted(tree *Tree, u, v int) int {
	bottom := ng.To[ng.Root]
	lca1, lca2 := tree.LCA(u, bottom), tree.LCA(v, bottom)
	distOnCycle := abs(tree.DepthWeighted[lca1] - tree.DepthWeighted[lca2])
	distOnCycle = min(distOnCycle, tree.DepthWeighted[bottom]+ng.OutCost-distOnCycle)
	return distOnCycle + tree.DepthWeighted[u] + tree.DepthWeighted[v] - tree.DepthWeighted[lca1] - tree.DepthWeighted[lca2]
}

type Uf struct {
	data []int
}

func NewUf(n int) *Uf {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &Uf{data: data}
}

func (ufa *Uf) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *Uf) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

type Tree struct {
	Tree                 [][]Neighbor
	Depth, DepthWeighted []int
	Parent               []int
	LID, RID             []int // 欧拉序[in,out)
	IdToNode             []int
	top, heavySon        []int
	timer                int
}

func NewTree(adjList [][]Neighbor, root int) *Tree {
	n := len(adjList)
	lid := make([]int, n)
	rid := make([]int, n)
	idToNode := make([]int, n)
	top := make([]int, n)   // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int, n) // 深度
	depthWeighted := make([]int, n)
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	res := &Tree{
		Tree:          adjList,
		Depth:         depth,
		DepthWeighted: depthWeighted,
		Parent:        parent,
		LID:           lid,
		RID:           rid,
		IdToNode:      idToNode,
		top:           top,
		heavySon:      heavySon,
	}
	res.Build(root)
	return res
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

func (tree *Tree) RootedParent(u int, root int) int {
	return tree.Jump(u, root, 1)
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
//	kthAncestor(root,0) == root
func (tree *Tree) KthAncestor(root, k int) int {
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

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree) EnumeratePathDecomposition(u, v int, vertex bool, f func(start, end int)) {
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

	edgeInt := 1
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

func (tree *Tree) GetPath(u, v int) []int {
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
func (tree *Tree) SubSize(v, root int) int {
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

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链底端节点.
func (tree *Tree) GetHeavyPath(start int) []int {
	heavyPath := []int{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

// 结点v的重儿子.如果没有重儿子,返回-1.
func (tree *Tree) GetHeavyChild(v int) int {
	k := tree.LID[v] + 1
	if k == len(tree.Tree) {
		return -1
	}
	w := tree.IdToNode[k]
	if tree.Parent[w] == v {
		return w
	}
	return -1
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
		next, weight := e.to, e.weight
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

func max(a, b int) int {
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
