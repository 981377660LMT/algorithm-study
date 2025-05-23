// Functional-函数图(每个顶点出度为1的有向图,有向的NamoriGraph)
// 定义：Directed graphs in which every vertex has exactly one outgoing edge.
// !每个点的出度为1(如果顶点没有出边，那么它的出边指向自己)
// 连通分量个数=环的个数
//
// !建议使用 基环树找到所有环.go
//
// 1. 如果竞赛图无环，那么竞赛图的拓扑序是唯一确定的
// 2. 竞赛图的强连通分量缩点后呈链状
// 3. 如果竞赛图是强连通的，则一定存在一条哈密顿回路
// 4. 竞赛图存在一条哈密顿路径
// 5. 大小为n的竞赛图如果强连通，则恰好有长度为3,4,…,n的简单环。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// edges = [2,-1,3,1]
	edges := []int{2, -1, 3, 1}
	fmt.Println(longestCycle(edges))
}

func demo() {
	F := NewFunctionalGraph(6)
	F.AddDirectedEdge(0, 1, 1)
	F.AddDirectedEdge(1, 2, 1)
	F.AddDirectedEdge(2, 3, 1)
	F.AddDirectedEdge(3, 4, 1)
	F.AddDirectedEdge(4, 5, 1)
	F.AddDirectedEdge(5, 0, 1)
	F.Build()
	fmt.Println(F.Jump(0, 7))
}

func yuki1242() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &nums[i])
	}
	G := NewFunctionalGraph(64)
	for s := 0; s < 64; s++ {
		t := (2 * s) & 63
		ok := true
		if s&1 == 0 && s&32 == 0 {
			ok = false
		}
		if s&2 == 0 && s&16 == 0 {
			ok = false
		}
		if s&4 == 0 && s&8 == 0 {
			ok = false
		}
		if ok {
			t |= 1
		}
		G.AddDirectedEdge(s, t, 1)
	}
	G.Build()
	x := n
	s := 63
	for i := k - 1; i >= 0; i-- {
		y := nums[i]
		s = G.Jump(s, x-y)
		s &= 62
		x = y
	}
	s = G.Jump(s, x-1)
	if s&1 == 1 {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}

func transitionGame() {
	// https://atcoder.jp/contests/abc296/tasks/abc296_e
	// 给定一个竞赛图,求多少个点在环中
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	F := NewFunctionalGraph(n)
	for i := 0; i < n; i++ {
		F.AddDirectedEdge(i, nums[i]-1, 1)
	}
	F.Build()
	res := 0
	for i := 0; i < n; i++ {
		if F.IsInCycle(i) {
			res++
		}
	}
	fmt.Fprintln(out, res)
}

// 给定一个竞赛图,求最长的环的长度,如果没有环,返回-1.
func longestCycle(edges []int) int {
	n := len(edges)
	F := NewFunctionalGraph(n)
	for i := 0; i < n; i++ {
		if edges[i] != -1 {
			F.AddDirectedEdge(i, edges[i], 1)
		} else {
			F.AddDirectedEdge(i, i, 1) // !如果顶点没有出边,那么它的出边指向自己.
		}
	}
	F.Build()

	res := -1
	cycles := F.CollectAllCycles()
	for _, cycle := range cycles {
		if len(cycle) > 1 {
			res = max(res, len(cycle))
		}
	}
	return res
}

// !每个点的出度为1(如果顶点没有出边，那么它的出边指向自己).
type FunctionalGraph struct {
	G      [][][2]int // (next, weight) 有向图
	Tree   *Tree
	n, m   int
	to     []int
	weight []int
	root   []int
}

func NewFunctionalGraph(n int) *FunctionalGraph {
	fg := &FunctionalGraph{
		n:      n,
		to:     make([]int, n),
		weight: make([]int, n),
		root:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		fg.to[i] = -1
		fg.root[i] = -1
	}
	return fg
}

func (fg *FunctionalGraph) AddDirectedEdge(u, v, w int) {
	if fg.to[u] != -1 {
		panic("FunctionalGraph: u already has an outgoing edge")
	}
	fg.m++
	fg.to[u] = v
	fg.weight[u] = w
}

func (fg *FunctionalGraph) Build() {
	n := fg.n
	uf := _NewUF(n)
	for i := 0; i < n; i++ {
		if !uf.Union(i, fg.to[i]) {
			fg.root[i] = i
		}
	}
	for i := 0; i < n; i++ {
		if fg.root[i] == i {
			fg.root[uf.Find(i)] = i
		}
	}
	for i := 0; i < n; i++ {
		fg.root[i] = fg.root[uf.Find(i)]
	}

	g := make([][][2]int, n+1)
	for i := 0; i < n; i++ {
		if fg.root[i] == i {
			g[n] = append(g[n], [2]int{i, fg.weight[i]})
		} else {
			to := fg.to[i]
			g[to] = append(g[to], [2]int{i, fg.weight[i]})
		}
	}

	tree := _NT(g)
	tree.Build(n)

	fg.G = g
	fg.Tree = tree
}

// 从 v 出发，走 step 步，返回到达的点.
func (fg *FunctionalGraph) Jump(v, step int) int {
	d := fg.Tree.Depth[v]
	if step <= d-1 {
		return fg.Tree.Jump(v, fg.n, step)
	}
	v = fg.root[v]
	step -= d - 1
	bottom := fg.to[v]
	c := fg.Tree.Depth[bottom]
	step %= c
	if step == 0 {
		return v
	}
	return fg.Tree.Jump(bottom, fg.n, step-1)
}

// 给定跳跃步数，返回`每个节点`在该步数内跳跃到的目标节点编号.
func (fg *FunctionalGraph) JumpAll(step int) []int {
	G := fg.Tree.Tree
	res := make([]int, fg.n)
	for i := 0; i < fg.n; i++ {
		res[i] = -1
	}
	query := make([][][2]int, fg.n)
	for v := 0; v < fg.n; v++ {
		d := fg.Tree.Depth[v]
		r := fg.root[v]
		if d-1 > step {
			query[v] = append(query[v], [2]int{v, step})
		} else {
			k := step - (d - 1)
			bottom := fg.to[r]
			c := fg.Tree.Depth[bottom]
			k %= c
			if k == 0 {
				res[v] = r
				continue
			}
			query[bottom] = append(query[bottom], [2]int{v, k - 1})
		}
	}

	path := make([]int, 0)
	var dfs func(v int)
	dfs = func(v int) {
		path = append(path, v)
		for _, q := range query[v] {
			res[q[0]] = path[len(path)-1-q[1]]
		}
		for _, e := range G[v] {
			dfs(e[0])
		}
		path = path[:len(path)-1]
	}
	for _, e := range G[fg.n] {
		dfs(e[0])
	}
	return res
}

// 判断节 v 点是否在 FunctionalGraph 对应的有向图中的环中.
func (fg *FunctionalGraph) IsInCycle(v int) bool {
	return fg.Tree.IsInSubtree(fg.to[fg.root[v]], v)
}

// 给定环的根节点，返回该环上所有节点的编号.
func (fg *FunctionalGraph) CollectCycle(root int) []int {
	if !fg.IsInCycle(root) {
		return nil
	}
	cycle := []int{fg.to[root]}
	for cycle[len(cycle)-1] != root {
		cycle = append(cycle, fg.to[cycle[len(cycle)-1]])
	}
	return cycle
}

// 返回所有环.
func (fg *FunctionalGraph) CollectAllCycles() [][]int {
	res := make([][]int, 0)
	visited := make([]bool, fg.n)
	for i := 0; i < fg.n; i++ {
		if !visited[i] {
			cycle := fg.CollectCycle(i)
			for _, v := range cycle {
				visited[v] = true
			}
			res = append(res, cycle)
		}
	}
	return res
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

type Tree struct {
	Tree          [][][2]int // (next, weight) 有向图
	Depth         []int
	Parent        []int
	LID, RID      []int // 欧拉序[in,out)
	IdToNode      []int
	top, heavySon []int
	timer         int
}

func _NT(graph [][][2]int) *Tree {
	n := len(graph)
	lid := make([]int, n)
	rid := make([]int, n)
	IdToNode := make([]int, n)
	top := make([]int, n)
	depth := make([]int, n)    // 深度
	parent := make([]int, n)   // 父结点
	heavySon := make([]int, n) // 重儿子
	for i := range parent {
		parent[i] = -1
	}

	return &Tree{
		Tree:     graph, // 有向图
		Depth:    depth,
		Parent:   parent,
		LID:      lid,
		RID:      rid,
		IdToNode: IdToNode,
		top:      top,
		heavySon: heavySon,
	}
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

func (tree *Tree) RootedLCA(u, v int, w int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, w) ^ tree.LCA(v, w)
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

// child 是否在 root 的子树中 (child和root不能相等)
func (tree *Tree) IsInSubtree(child, root int) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

func (tree *Tree) build(cur, pre, dep, dist int) int {
	subSize, heavySize, heavySon := 1, 0, -1
	for _, e := range tree.Tree[cur] { // 有向
		next, weight := e[0], e[1]
		nextSize := tree.build(next, cur, dep+1, dist+weight)
		subSize += nextSize
		if nextSize > heavySize {
			heavySize, heavySon = nextSize, next
		}
	}
	tree.Depth[cur] = dep
	tree.heavySon[cur] = heavySon
	tree.Parent[cur] = pre
	return subSize
}

func (tree *Tree) markTop(cur, top int) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	if tree.heavySon[cur] != -1 {
		tree.markTop(tree.heavySon[cur], top)
		for _, e := range tree.Tree[cur] {
			next := e[0]
			if next != tree.heavySon[cur] { // 有向
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
}

func _NewUF(n int) *_UF {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UF{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type _UF struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *_UF) Union(key1, key2 int) bool {
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

func (ufa *_UF) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UF) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UF) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}
