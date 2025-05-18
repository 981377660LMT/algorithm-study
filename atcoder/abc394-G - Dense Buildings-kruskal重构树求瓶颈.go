// abc394-G - Dense Buildings-kruskal重构树
// https://atcoder.jp/contests/abc394/tasks/abc394_g
//
// 在一个H×W的城市网格中，每个格子有一座F_i,j层高的建筑。
// 高桥君有两种移动方式：
//
// 使用楼梯在同一建筑内上下移动（每次上下一层计为1次使用楼梯）
// 使用空中通道从当前建筑的X层横向移动到相邻建筑的X层（前提是目标建筑至少有X层）
// !计算从起点建筑的某一层到终点建筑的某一层所需的最少楼梯使用次数。
//
// !本质上是找到一条从起点到终点的路径，使得路径上所有建筑的最低高度尽可能高（因为这决定了我们可以在多高的楼层使用空中通道）

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
	"strconv"
)

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

	h, w := io.NextInt(), io.NextInt()
	grid := make([][]int, h)
	for i := range grid {
		grid[i] = make([]int, w)
		for j := range grid[i] {
			grid[i][j] = io.NextInt()
		}
	}

	edges := make([]Edge, 0, h*w*2)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if i+1 < h {
				// !边权为两个格子的最小高度
				from, to, weight := i*w+j, (i+1)*w+j, min(grid[i][j], grid[i+1][j])
				edges = append(edges, Edge{from: from, to: to, weight: weight})
			}
			if j+1 < w {
				from, to, weight := i*w+j, i*w+j+1, min(grid[i][j], grid[i][j+1])
				edges = append(edges, Edge{from: from, to: to, weight: weight})
			}
		}
	}

	// !最大生成树 => 最大瓶颈路
	sort.Slice(edges, func(i, j int) bool { return edges[i].weight > edges[j].weight })
	forest, roots, values := KruskalTree(h*w, edges)
	tree := _NewTree(forest, roots)

	// !从(x1,y1,h1)到(x2,y2,h2)，爬楼梯的最小代价
	query := func(x1, y1, h1, x2, y2, h2 int) int {
		u, v := x1*w+y1, x2*w+y2
		var minH int
		if u == v {
			minH = INF // 边权不存在
		} else {
			minH = values[tree.LCA(u, v)]
		}

		minH = min(minH, min(h1, h2))
		return (h1 - minH) + (h2 - minH)
	}

	q := io.NextInt()
	for i := 0; i < q; i++ {
		x1, y1, h1 := io.NextInt(), io.NextInt(), io.NextInt()
		x1, y1 = x1-1, y1-1
		x2, y2, h2 := io.NextInt(), io.NextInt(), io.NextInt()
		x2, y2 = x2-1, y2-1
		io.Println(query(x1, y1, h1, x2, y2, h2))
	}
}

const INF int = 1e18

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

func (tree *_Tree) RootedLCA(u, v int, w int) int {
	return tree.LCA(u, v) ^ tree.LCA(u, w) ^ tree.LCA(v, w)
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
