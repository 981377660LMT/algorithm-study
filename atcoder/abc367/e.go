package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

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

	N, K := io.NextInt(), io.NextInt()
	X := make([]int, N)
	A := make([]int, N)
	for i := 0; i < N; i++ {
		X[i] = io.NextInt() - 1
	}
	for i := 0; i < N; i++ {
		A[i] = io.NextInt() - 1
	}

	F := NewFunctionalGraph(int32(N))
	for i := 0; i < N; i++ {
		F.AddDirectedEdge(int32(i), int32(X[i]), 0)
	}

	F.Build()

	res := F.JumpAll(K)
	for i := 0; i < N; i++ {
		io.Print(A[res[i]]+1, " ")
	}

}

type neighbor struct {
	to     int32
	weight int
}

type FunctionalGraph struct {
	To     []int32 // 每个顶点的出边指向的顶点
	Weight []int   // 每个顶点的出边的权值
	Root   []int32 // 每个联通分量的起点
	n, m   int32
	graph  [][]neighbor
	tree   *Tree
}

func NewFunctionalGraph(n int32) *FunctionalGraph {
	to, weight, root := make([]int32, n), make([]int, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		to[i] = -1
		root[i] = -1
	}
	return &FunctionalGraph{n: n, To: to, Weight: weight, Root: root}
}

func (f *FunctionalGraph) AddDirectedEdge(from, to int32, weight int) {
	if f.To[from] != -1 {
		panic("FunctionalGraph: each vertex must have exactly one outgoing edge")
	}
	f.m++
	f.To[from] = to
	f.Weight[from] = weight
}

func (f *FunctionalGraph) Build() ([][]neighbor, *Tree) {
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
func (f *FunctionalGraph) Dist(from, to int32, weighted bool) int {
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
func (f *FunctionalGraph) Jump(v int32, step int) int32 {
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

func (f *FunctionalGraph) JumpAll(step int) []int32 {
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

func (f *FunctionalGraph) InCycle(v int32) bool {
	root := f.Root[v]
	bottom := f.To[root]
	return f.tree.IsInSubtree(bottom, v)
}

func (f *FunctionalGraph) CollectCycle(r int32) []int32 {
	if r != f.Root[r] {
		panic("FunctionalGraph: r must be root")
	}
	cycle := []int32{f.To[r]}
	for last := cycle[len(cycle)-1]; last != r; last = cycle[len(cycle)-1] {
		cycle = append(cycle, f.To[last])
	}
	return cycle
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
