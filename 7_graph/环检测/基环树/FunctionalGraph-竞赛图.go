// FunctionalGraph-函数图(每个顶点出度为1的有向图,有向的NamoriGraph)
// 定义：Directed graphs in which every vertex has exactly one outgoing edge.
// !每个点的出度为1(如果顶点没有出边，那么它的出边指向自己)
// 连通分量个数=环的个数
// 1. 如果竞赛图无环，那么竞赛图的拓扑序是唯一确定的
// 2. 竞赛图的强连通分量缩点后呈链状
// 3. 如果竞赛图是强连通的，则一定存在一条哈密顿回路
// 4. 竞赛图存在一条哈密顿路径
// 5. 大小为n的竞赛图如果强连通，则恰好有长度为3,4,…,n的简单环。
//
// 例如：下图是一个竞赛图。
//
//	  0
//	  ↓
//	  1
//	 ↙ ↖
//	2   3 ← 4 ← 5
//	 ↘ ↗
//	  6 ← 7
//
// root: 6
// bottom: 3
// to: [1 2 6 1 3 4 3 6]
// directedGraph:
//
//	        8 (虚拟节点)
//	       /
//	      6(分量root)
//	     / \
//	    2   7
//	   /
//	  1
//	 / \
//	0   3 (bottom)
//	     \
//	      4
//	       \
//	        5
//
// !1.n作为树的虚拟根节点，联通各个分量的起点.
// !2.bottom的祖先节点都在环中.
// !3.点u在所在子树的根节点(在环上)为lca(u, bottom).
//
// api:
// 1. Dist(from, to int32, weighted bool) int
// 2. Jump(v int32, step int) int32
// 3. JumpAll(step int) []int32
// 4. InCycle(v int32) bool
// 5. CollectCycle(r int32) []int32
// 6. MeetTime(i, j int32) int32

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// yuki1242()
	// abc296_e()
	demo()
}

func demo() {
	edges := [][]int{{0, 1, 1}, {1, 2, 1}, {3, 1, 1}, {2, 6, 1}, {6, 3, 9}, {4, 3, 1}, {5, 4, 1}, {7, 6, 1}}
	n := int32(8)
	F := NewFunctionalGraph32(n)
	for _, e := range edges {
		F.AddDirectedEdge(int32(e[0]), int32(e[1]), e[2])
	}
	F.Build()
	fmt.Println(F.Dist(7, 1, true))  // 11
	fmt.Println(F.Dist(7, 1, false)) // 3
	fmt.Println(F.Root, F.graph, F.To)

	fmt.Println(F.Jump(7, 12)) // 2
	fmt.Println(F.JumpAll(100))
	for i := int32(0); i < n; i++ {
		fmt.Println(F.InCycle(i))
	}
	for i := int32(0); i < n; i++ {
		if F.Root[i] == i {
			fmt.Println(F.CollectCycle(i))
		}
	}

	fmt.Println(F.MeetTime(0, 1)) // 3
	fmt.Println(F.MeetTime(0, 3)) // 3
	fmt.Println(F.MeetTime(2, 5)) // 3
}

// No.1242 高橋君とすごろく
// https://yukicoder.me/problems/no/1242
func yuki1242() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	var k int32
	fmt.Fscan(in, &n, &k)
	nums := make([]int, k)
	for i := int32(0); i < k; i++ {
		fmt.Fscan(in, &nums[i])
	}
	G := NewFunctionalGraph32(64)
	for s := int32(0); s < 64; s++ {
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
	s := int32(63)
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

func abc296_e() {
	// https://atcoder.jp/contests/abc296/tasks/abc296_e
	// 给定一个竞赛图,求多少个点在环中
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	F := NewFunctionalGraph32(n)
	for i := int32(0); i < n; i++ {
		F.AddDirectedEdge(i, nums[i]-1, 1)
	}
	F.Build()
	res := 0
	for i := int32(0); i < n; i++ {
		if F.InCycle(i) {
			res++
		}
	}
	fmt.Fprintln(out, res)
}

// 给定一个竞赛图,求最长的环的长度,如果没有环,返回-1.
func longestCycle(edges []int) int {
	n := int32(len(edges))
	F := NewFunctionalGraph32(n)
	for i := int32(0); i < n; i++ {
		if edges[i] != -1 {
			F.AddDirectedEdge(i, int32(edges[i]), 1)
		} else {
			F.AddDirectedEdge(i, i, 1) // !如果顶点没有出边,那么它的出边指向自己.
		}
	}
	F.Build()

	res := -1
	for i := int32(0); i < n; i++ {
		if F.Root[i] == i {
			cycle := F.CollectCycle(int32(i))
			if len(cycle) > 1 {
				res = max(res, len(cycle))
			}
		}
	}
	return res
}

type neighbor struct {
	to     int32
	weight int
}

type FunctionalGraph32 struct {
	To     []int32 // 每个顶点的出边指向的顶点
	Weight []int   // 每个顶点的出边的权值
	Root   []int32 // 每个联通分量的起点
	n, m   int32
	graph  [][]neighbor
	tree   *Tree
}

func NewFunctionalGraph32(n int32) *FunctionalGraph32 {
	to, weight, root := make([]int32, n), make([]int, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		to[i] = -1
		root[i] = -1
	}
	return &FunctionalGraph32{n: n, To: to, Weight: weight, Root: root}
}

func (f *FunctionalGraph32) AddDirectedEdge(from, to int32, weight int) {
	if f.To[from] != -1 {
		panic("FunctionalGraph: each vertex must have exactly one outgoing edge")
	}
	f.m++
	f.To[from] = to
	f.Weight[from] = weight
}

func (f *FunctionalGraph32) Build() ([][]neighbor, *Tree) {
	if f.n != f.m {
		panic("FunctionalGraph: vertex count must be equal to edge count")
	}
	n := f.n
	uf := newUnionFindArraySimple32(n)
	for v := int32(0); v < n; v++ {
		if !uf.Union(v, f.To[v]) {
			f.Root[v] = v
		}
	}
	for v := int32(0); v < n; v++ {
		if f.Root[v] == v {
			f.Root[uf.Find(v)] = v
		}
	}
	for v := int32(0); v < n; v++ {
		f.Root[v] = f.Root[uf.Find(v)]
	}

	graph := make([][]neighbor, n+1)
	for v := int32(0); v < n; v++ {
		if f.Root[v] == v {
			graph[n] = append(graph[n], neighbor{to: v, weight: f.Weight[v]})
		} else {
			graph[f.To[v]] = append(graph[f.To[v]], neighbor{to: v, weight: f.Weight[v]})
		}
	}
	f.graph = graph

	tree := NewTree(graph)
	tree.Build(n)
	f.tree = tree
	return graph, tree
}

// 从from到to的距离,不可达返回-1.
func (f *FunctionalGraph32) Dist(from, to int32, weighted bool) int {
	if weighted {
		if f.tree.IsInSubtree(from, to) {
			return f.tree.DepthWeighted[from] - f.tree.DepthWeighted[to]
		}
		root := f.Root[from]
		bottom := f.To[root]
		// from -> root -> bottom -> to
		if f.tree.IsInSubtree(bottom, to) {
			x := f.tree.DepthWeighted[from] - f.tree.DepthWeighted[root]
			x += f.Weight[root]
			x += f.tree.DepthWeighted[bottom] - f.tree.DepthWeighted[to]
			return x
		}
		return -1
	} else {
		if f.tree.IsInSubtree(from, to) {
			return int(f.tree.Depth[from] - f.tree.Depth[to])
		}
		root := f.Root[from]
		bottom := f.To[root]
		// from -> root -> bottom -> to
		if f.tree.IsInSubtree(bottom, to) {
			x := f.tree.Depth[from] - f.tree.Depth[root]
			x++
			x += f.tree.Depth[bottom] - f.tree.Depth[to]
			return int(x)
		}
		return -1
	}
}

// 从v向前跳step步,返回跳到的节点,不可达返回-1.
func (f *FunctionalGraph32) Jump(v int32, step int) int32 {
	d := f.tree.Depth[v]
	if step <= int(d-1) {
		return f.tree.Jump(v, f.n, int32(step))
	}
	v = f.Root[v]
	step -= int(d - 1)
	bottom := f.To[v]
	c := f.tree.Depth[bottom]
	step %= int(c)
	if step == 0 {
		return v
	}
	return f.tree.Jump(bottom, f.n, int32(step-1))
}

func (f *FunctionalGraph32) JumpAll(step int) []int32 {
	n := f.n
	res := make([]int32, n)
	for v := int32(0); v < n; v++ {
		res[v] = -1
	}
	query := make([][][2]int32, n)
	for v := int32(0); v < n; v++ {
		d := int(f.tree.Depth[v])
		r := f.Root[v]
		if d-1 > step {
			query[v] = append(query[v], [2]int32{v, int32(step)})
		}
		if d-1 <= step {
			k := step - (d - 1)
			bottom := f.To[r]
			c := int(f.tree.Depth[bottom])
			k %= c
			if k == 0 {
				res[v] = r
				continue
			}
			query[bottom] = append(query[bottom], [2]int32{v, int32(k - 1)})
		}
	}

	path := []int32{}
	var dfs func(int32)
	dfs = func(v int32) {
		path = append(path, v)
		for _, e := range query[v] {
			w, k := e[0], e[1]
			res[w] = path[int32(len(path))-1-k]
		}
		for _, e := range f.graph[v] {
			dfs(e.to)
		}
		path = path[:len(path)-1]
	}
	for _, e := range f.graph[n] {
		dfs(e.to)
	}
	return res
}

func (f *FunctionalGraph32) InCycle(v int32) bool {
	root := f.Root[v]
	bottom := f.To[root]
	return f.tree.IsInSubtree(bottom, v)
}

func (f *FunctionalGraph32) CollectCycle(r int32) []int32 {
	if r != f.Root[r] {
		panic("FunctionalGraph: r must be root")
	}
	cycle := []int32{f.To[r]}
	for last := cycle[len(cycle)-1]; last != r; last = cycle[len(cycle)-1] {
		cycle = append(cycle, f.To[last])
	}
	return cycle
}

// 点i和点j的最小相遇时间(跳跃次数).
// 无解返回-1.
func (f *FunctionalGraph32) MeetTime(i, j int32) int32 {
	if i == j {
		return 0
	}
	if f.Root[i] != f.Root[j] {
		return -1
	}
	r := f.Root[i]
	b := f.To[r]
	cycleLen := f.tree.Depth[b] - f.tree.Depth[r] + 1
	if (f.tree.Depth[i]-f.tree.Depth[j])%cycleLen != 0 {
		return -1
	}
	if f.tree.Depth[i] == f.tree.Depth[j] {
		lca := f.tree.LCA(i, j)
		return f.tree.Depth[i] - f.tree.Depth[lca]
	}
	ti := f.tree.Depth[i] - f.tree.Depth[f.tree.LCA(b, i)]
	tj := f.tree.Depth[j] - f.tree.Depth[f.tree.LCA(b, j)]
	return max32(ti, tj)
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

func (tree *Tree) RootedLCA(u, v int32, root int32) int32 {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
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
func (tree *Tree) Jump(from, to int32, step int32) int32 {
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
	u.data[root2] = root1
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
