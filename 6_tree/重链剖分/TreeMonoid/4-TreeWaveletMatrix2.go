package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	AnUnusualKinginPakenKingdom()
	// L番目の数字()
}

// https://atcoder.jp/contests/pakencamp-2022-day1/tasks/pakencamp_2022_day1_j
// 路径上的中位数
func AnUnusualKinginPakenKingdom() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	tree := NewTree32(n)
	edgeValues := make([]int, n-1)
	for i := int32(0); i < n-1; i++ {
		var a, b int32
		var c int
		fmt.Fscan(in, &a, &b, &c)
		tree.AddEdge(a-1, b-1, 1)
		edgeValues[i] = c
	}
	tree.Build(0)
	twm := NewTreeWaveletMatrix(tree, edgeValues, false, nil, -1, true)

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		x := twm.MedianPath(a, b, false, 0)
		y := twm.MedianPath(a, b, true, 0)
		fmt.Fprintln(out, (x+y)/2)
	}
}

func L番目の数字() {
	// // https://atcoder.jp/contests/utpc2011/tasks/utpc2011_12
	// // 路径上第k小的数

	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// var n, q int
	// fmt.Fscan(in, &n, &q)
	// tree := _NT(n)
	// values := make([]int, n)
	// for i := 0; i < n; i++ {
	// 	fmt.Fscan(in, &values[i])
	// }
	// for i := 0; i < n-1; i++ {
	// 	var a, b int
	// 	fmt.Fscan(in, &a, &b)
	// 	tree.AddEdge(a-1, b-1, 1)
	// }
	// tree.Build(0)
	// twm := NewTreeWaveletMatrix(tree, values, true, -1, nil)
	// for i := 0; i < q; i++ {
	// 	var a, b, k int
	// 	fmt.Fscan(in, &a, &b, &k)
	// 	a, b, k = a-1, b-1, k-1
	// 	fmt.Fprintln(out, twm.KthPath(a, b, k, 0))
	// }
}

const INF WmValue = 1e18

type WmValue = int
type WmSum = int

func e() WmSum            { return 0 }
func op(a, b WmSum) WmSum { return a + b }
func inv(a WmSum) WmSum   { return -a }

type TreeWaveletMatrix struct {
	tree     *Tree32
	n        int32
	wm       *WaveletMatrixWithSum
	isVertex bool
	compress bool
}

// data: 第i个顶点的值或者第i条边的权值.
// isVertex: data是否为顶点的值以及查询的时候是否是顶点权值.
// sumData: 和数据，nil表示不需要和数据.
// log: 如果需要支持异或查询则需要传入log，-1表示默认.
// compress: 是否对data进行离散化(值域较大(1e9)时可以离散化加速).
func NewTreeWaveletMatrix(tree *Tree32, data []WmValue, isVertex bool, sumData []WmSum, log int32, compress bool) *TreeWaveletMatrix {
	res := &TreeWaveletMatrix{
		tree: tree, n: int32(len(tree.Tree)),
		wm: &WaveletMatrixWithSum{}, isVertex: isVertex, compress: compress,
	}
	n := res.n
	A1 := make([]WmValue, n)
	var S1 []WmSum
	if len(sumData) > 0 {
		S1 = make([]WmSum, n)
	}

	if isVertex {
		if !(len(data) == int(n) && (len(sumData) == 0 || len(sumData) == int(n))) {
			panic("invalid input")
		}
		for v := int32(0); v < n; v++ {
			A1[tree.LID[v]] = data[v]
		}
		if int32(len(sumData)) == n {
			for v := int32(0); v < n; v++ {
				S1[tree.LID[v]] = sumData[v]
			}
		}
		res.wm.build(A1, S1, log, compress)
	} else {
		if !(len(data) == int(n-1) && (len(sumData) == 0 || len(sumData) == int(n-1))) {
			panic("invalid input")
		}
		if len(sumData) != 0 {
			for e := int32(0); e < n-1; e++ {
				S1[tree.LID[tree.EidtoV(e)]] = sumData[e]
			}
		}
		for e := int32(0); e < n-1; e++ {
			A1[tree.LID[tree.EidtoV(e)]] = data[e]
		}
		res.wm.build(A1, S1, log, compress)
	}

	return res
}

// 路径s-t上值在[a,b)范围内的结点个数.
func (twm *TreeWaveletMatrix) CountPath(s, t int32, a, b, xorVal WmValue) int32 {
	return twm.wm.CountSegments(twm.getSegments(s, t), a, b, xorVal)
}

// 子树中值在[a,b)范围内的结点个数.
func (twm *TreeWaveletMatrix) CountSubtree(node int32, a, b, xorVal WmValue) int32 {
	l, r := twm.tree.LID[node], twm.tree.RID[node]
	offset := int32(1)
	if twm.isVertex {
		offset = 0
	}
	return twm.wm.CountRange(l+offset, r, a, b, xorVal)
}

// 路径s-t上第k小的值和前k小的结点和(k>=0).
func (twm *TreeWaveletMatrix) KthValueAndSumPath(s, t, k int32, xorVal WmValue) (WmValue, WmSum) {
	return twm.wm.KthValueAndSumSegments(twm.getSegments(s, t), k, xorVal)
}

// 子树中第k小的值和前k小的结点和(k>=0).
func (twm *TreeWaveletMatrix) KthValueAndSumSubtree(node, k int32, xorVal WmValue) (WmValue, WmSum) {
	l, r := twm.tree.LID[node], twm.tree.RID[node]
	offset := int32(1)
	if twm.isVertex {
		offset = 0
	}
	return twm.wm.KthValueAndSum(l+offset, r, k, xorVal)
}

// 路径s-t上第k小的值(k>=0).
func (twm *TreeWaveletMatrix) KthPath(s, t, k int32, xorVal WmValue) WmValue {
	return twm.wm.KthSegments(twm.getSegments(s, t), k, xorVal)
}

// 子树中第k小的值(k>=0).
func (twm *TreeWaveletMatrix) KthSubtree(node, k int32, xorVal WmValue) WmValue {
	l, r := twm.tree.LID[node], twm.tree.RID[node]
	offset := int32(1)
	if twm.isVertex {
		offset = 0
	}
	return twm.wm.Kth(l+offset, r, k, xorVal)
}

// 路径s-t上中位数.
func (twm *TreeWaveletMatrix) MedianPath(s, t int32, upper bool, xorVal WmValue) WmValue {
	return twm.wm.MedianSegments(twm.getSegments(s, t), upper, xorVal)
}

// 子树中中位数.
func (twm *TreeWaveletMatrix) MedianSubtree(node int32, upper bool, xorVal WmValue) WmValue {
	l, r := twm.tree.LID[node], twm.tree.RID[node]
	offset := int32(1)
	if twm.isVertex {
		offset = 0
	}
	return twm.wm.Median(l+offset, r, upper, xorVal)
}

// 路径s-t上第k1到k2小的和.
func (twm *TreeWaveletMatrix) SumPath(s, t, k1, k2 int32, xorVal WmValue) WmSum {
	return twm.wm.SumSegments(twm.getSegments(s, t), k1, k2, xorVal)
}

// 子树中第k1到k2小的和.
func (twm *TreeWaveletMatrix) SumSubtree(node, k1, k2 int32, xorVal WmValue) WmSum {
	l, r := twm.tree.LID[node], twm.tree.RID[node]
	offset := int32(1)
	if twm.isVertex {
		offset = 0
	}
	return twm.wm.SumRange(l+offset, r, k1, k2, xorVal)
}

// 路径s-t上所有结点和.
func (twm *TreeWaveletMatrix) SumAllPath(s, t int32) WmSum {
	return twm.wm.SumAllSegments(twm.getSegments(s, t))
}

// 子树中所有结点和.
func (twm *TreeWaveletMatrix) SumAllSubtree(node int32) WmSum {
	l, r := twm.tree.LID[node], twm.tree.RID[node]
	offset := int32(1)
	if twm.isVertex {
		offset = 0
	}
	return twm.wm.SumAll(l+offset, r)
}

func (twm *TreeWaveletMatrix) getSegments(s, t int32) [][2]int32 {
	segments := twm.tree.GetPathDecomposition(s, t, twm.isVertex)
	for i := 0; i < len(segments); i++ {
		seg := &segments[i]
		if seg[0] >= seg[1] {
			seg[0], seg[1] = seg[1], seg[0]
		}
		seg[1]++
	}
	return segments
}

type Tree32 struct {
	Tree          [][][2]int32 // (next, weight)
	Edges         [][3]int32   // (u,v,w)
	Depth         []int32
	DepthWeighted []int
	Parent        []int32
	LID, RID      []int32 // 欧拉序[in,out)
	IdToNode      []int32
	top, heavySon []int32
	timer         int32
}

func NewTree32(n int32) *Tree32 {
	tree := make([][][2]int32, n)
	lid := make([]int32, n)
	rid := make([]int32, n)
	IdToNode := make([]int32, n)
	top := make([]int32, n)   // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n) // 深度
	depthWeighted := make([]int, n)
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	edges := make([][3]int32, 0, n-1)
	for i := range parent {
		parent[i] = -1
	}

	return &Tree32{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: depthWeighted,
		Parent:        parent,
		LID:           lid,
		RID:           rid,
		IdToNode:      IdToNode,
		top:           top,
		heavySon:      heavySon,
		Edges:         edges,
	}
}

// 添加无向边 u-v, 边权为w.
func (tree *Tree32) AddEdge(u, v, w int32) {
	tree.Tree[u] = append(tree.Tree[u], [2]int32{v, w})
	tree.Tree[v] = append(tree.Tree[v], [2]int32{u, w})
	tree.Edges = append(tree.Edges, [3]int32{u, v, w})
}

// 添加有向边 u->v, 边权为w.
func (tree *Tree32) AddDirectedEdge(u, v, w int32) {
	tree.Tree[u] = append(tree.Tree[u], [2]int32{v, w})
	tree.Edges = append(tree.Edges, [3]int32{u, v, w})
}

// root:0-based
//
//	当root设为-1时，会从0开始遍历未访问过的连通分量
func (tree *Tree32) Build(root int32) {
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
func (tree *Tree32) Id(root int32) (int32, int32) {
	return tree.LID[root], tree.RID[root]
}

// 返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed.
func (tree *Tree32) Eid(u, v int32) int32 {
	if tree.LID[u] > tree.LID[v] {
		return tree.LID[u]
	}
	return tree.LID[v]
}

// 较深的那个点作为边的编号.
func (tree *Tree32) EidtoV(eid int32) int32 {
	e := tree.Edges[eid]
	u, v := e[0], e[1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (tree *Tree32) LCA(u, v int32) int32 {
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

func (tree *Tree32) RootedLCA(u, v int32, root int32) int32 {
	return tree.LCA(u, v) ^ tree.LCA(u, root) ^ tree.LCA(v, root)
}

func (tree *Tree32) RootedParent(u int32, root int32) int32 {
	return tree.Jump(u, root, 1)
}

func (tree *Tree32) Dist(u, v int32, weighted bool) int {
	if weighted {
		return tree.DepthWeighted[u] + tree.DepthWeighted[v] - 2*tree.DepthWeighted[tree.LCA(u, v)]
	}
	return int(tree.Depth[u] + tree.Depth[v] - 2*tree.Depth[tree.LCA(u, v)])
}

// k: 0-based
//
//	如果不存在第k个祖先，返回-1
//	kthAncestor(root,0) == root
func (tree *Tree32) KthAncestor(root, k int32) int32 {
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
func (tree *Tree32) Jump(from, to, step int32) int32 {
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

func (tree *Tree32) CollectChild(root int32) []int32 {
	res := []int32{}
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
func (tree *Tree32) GetPathDecomposition(u, v int32, vertex bool) [][2]int32 {
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
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, vertex bool, f func(start, end int32)) {
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

func (tree *Tree32) GetPath(u, v int32) []int32 {
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
func (tree *Tree32) SubSize(v, root int32) int32 {
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
func (tree *Tree32) IsInSubtree(child, root int32) bool {
	return tree.LID[root] <= tree.LID[child] && tree.LID[child] < tree.RID[root]
}

// 寻找以 start 为 top 的重链 ,heavyPath[-1] 即为重链底端节点.
func (tree *Tree32) GetHeavyPath(start int32) []int32 {
	heavyPath := []int32{start}
	cur := start
	for tree.heavySon[cur] != -1 {
		cur = tree.heavySon[cur]
		heavyPath = append(heavyPath, cur)
	}
	return heavyPath
}

// 结点v的重儿子.如果没有重儿子,返回-1.
func (tree *Tree32) GetHeavyChild(v int32) int32 {
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

func (tree *Tree32) ELID(u int32) int32 {
	return 2*tree.LID[u] - tree.Depth[u]
}

func (tree *Tree32) ERID(u int32) int32 {
	return 2*tree.RID[u] - tree.Depth[u] - 1
}

func (tree *Tree32) build(cur, pre, dep int32, dist int) int32 {
	subSize, heavySize, heavySon := int32(1), int32(0), int32(-1)
	for _, e := range tree.Tree[cur] {
		next, weight := e[0], e[1]
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

func (tree *Tree32) markTop(cur, top int32) {
	tree.top[cur] = top
	tree.LID[cur] = tree.timer
	tree.IdToNode[tree.timer] = cur
	tree.timer++
	heavySon := tree.heavySon[cur]
	if heavySon != -1 {
		tree.markTop(heavySon, top)
		for _, e := range tree.Tree[cur] {
			next := e[0]
			if next != heavySon && next != tree.Parent[cur] {
				tree.markTop(next, next)
			}
		}
	}
	tree.RID[cur] = tree.timer
}

type WaveletMatrixWithSum struct {
	n, log   int32
	mid      []int32
	bv       []*BitVector
	key      []WmValue
	setLog   bool
	presum   [][]WmSum
	compress bool
}

// nums: 数组元素.
// sumData: 和数据，nil表示不需要和数据.
// log: 如果需要支持异或查询则需要传入log，-1表示默认.
// compress: 是否对nums进行离散化(值域较大(1e9)时可以离散化加速).
func NewWaveletMatrixWithSum(nums []WmValue, sumData []WmSum, log int32, compress bool) *WaveletMatrixWithSum {
	wm := &WaveletMatrixWithSum{}
	wm.build(nums, sumData, log, compress)
	return wm
}

func (wm *WaveletMatrixWithSum) build(nums []WmValue, sumData []WmSum, log int32, compress bool) {
	numsCopy := append(nums[:0:0], nums...)
	sumDataCopy := append(sumData[:0:0], sumData...)

	wm.n = int32(len(numsCopy))
	wm.log = log
	wm.compress = compress
	wm.setLog = log != -1
	if wm.n == 0 {
		wm.log = 0
		return
	}
	makeSum := len(sumData) > 0
	if compress {
		if wm.setLog {
			panic("compress and log should not be set at the same time")
		}
		wm.key = make([]WmValue, 0, wm.n)
		order := wm._argSort(numsCopy)
		for _, i := range order {
			if len(wm.key) == 0 || wm.key[len(wm.key)-1] != numsCopy[i] {
				wm.key = append(wm.key, numsCopy[i])
			}
			numsCopy[i] = WmValue(len(wm.key) - 1)
		}
	}
	if wm.log == -1 {
		tmp := wm._maxs(numsCopy)
		if tmp < 1 {
			tmp = 1
		}
		wm.log = int32(bits.Len(uint(tmp)))
	}
	wm.mid = make([]int32, wm.log)
	wm.bv = make([]*BitVector, wm.log)
	for i := range wm.bv {
		wm.bv[i] = NewBitVector(wm.n)
	}
	if makeSum {
		wm.presum = make([][]WmSum, 1+wm.log)
		for i := range wm.presum {
			sums := make([]WmSum, wm.n+1)
			for j := range sums {
				sums[j] = e()
			}
			wm.presum[i] = sums
		}
	}
	if len(sumDataCopy) == 0 {
		sumDataCopy = make([]WmSum, len(numsCopy))
	}

	A, S := numsCopy, sumDataCopy
	A0, A1 := make([]WmValue, wm.n), make([]WmValue, wm.n)
	S0, S1 := make([]WmSum, wm.n), make([]WmSum, wm.n)
	for d := wm.log - 1; d >= -1; d-- {
		p0, p1 := int32(0), int32(0)
		if makeSum {
			tmp := wm.presum[d+1]
			for i := int32(0); i < wm.n; i++ {
				tmp[i+1] = op(tmp[i], S[i])
			}
		}
		if d == -1 {
			break
		}
		for i := int32(0); i < wm.n; i++ {
			f := (A[i] >> d & 1) == 1
			if !f {
				if makeSum {
					S0[p0] = S[i]
				}
				A0[p0] = A[i]
				p0++
			} else {
				if makeSum {
					S1[p1] = S[i]
				}
				wm.bv[d].Set(i)
				A1[p1] = A[i]
				p1++
			}
		}
		wm.mid[d] = p0
		wm.bv[d].Build()
		A, A0 = A0, A
		S, S0 = S0, S
		for i := int32(0); i < p1; i++ {
			A[p0+i] = A1[i]
			S[p0+i] = S1[i]
		}
	}
}

// 区间 [start, end) 中范围在 [a, b) 中的元素的个数.
func (wm *WaveletMatrixWithSum) CountRange(start, end int32, a, b WmValue, xorVal WmValue) int32 {
	return wm.CountPrefix(start, end, b, xorVal) - wm.CountPrefix(start, end, a, xorVal)
}

func (wm *WaveletMatrixWithSum) SumRange(start, end, k1, k2 int32, xorVal WmValue) WmSum {
	if k1 < 0 {
		k1 = 0
	}
	if k2 > end-start {
		k2 = end - start
	}
	if k1 >= k2 {
		return e()
	}
	add := wm.SumPrefix(start, end, k2, xorVal)
	sub := wm.SumPrefix(start, end, k1, xorVal)
	return op(add, inv(sub))
}

// 返回区间 [start, end) 中 范围在 [0, x) 中的元素的个数.
func (wm *WaveletMatrixWithSum) CountPrefix(start, end int32, x WmValue, xor WmValue) int32 {
	if xor != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if wm.compress {
		x = wm._lowerBound(wm.key, x)
	}
	if x == 0 {
		return 0
	}
	if x >= 1<<wm.log {
		return end - start
	}
	count := int32(0)
	for d := wm.log - 1; d >= 0; d-- {
		add := (x>>d)&1 == 1
		f := (xor>>d)&1 == 1
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		kf := int32(0)
		if f {
			kf = (end - start - r0 + l0)
		} else {
			kf = (r0 - l0)
		}
		if add {
			count += kf
			if f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		} else {
			if !f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		}
	}
	return count
}

// 返回区间 [start, end) 中 [0, k) 的和.
func (wm *WaveletMatrixWithSum) SumPrefix(start, end, k int32, xor WmValue) WmSum {
	_, sum := wm.KthValueAndSum(start, end, k, xor)
	return sum
}

// [start, end)区间内第k(k>=0)小的元素.
func (wm *WaveletMatrixWithSum) Kth(start, end, k int32, xorVal WmValue) WmValue {
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set")
		}
	}
	count := int32(0)
	res := WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		var c int32
		if f {
			c = (end - start) - (r0 - l0)
		} else {
			c = r0 - l0
		}
		if count+c > k {
			if !f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		} else {
			count += c
			res |= 1 << d
			if !f {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			} else {
				start, end = l0, r0
			}
		}
	}
	if wm.compress {
		res = wm.key[res]
	}
	return res
}

// 返回区间 [start, end) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果).
// 如果k >= end-start, 返回 (-1, 区间 op 的结果).
func (wm *WaveletMatrixWithSum) KthValueAndSum(start, end, k int32, xorVal WmValue) (WmValue, WmSum) {
	if start >= end {
		return -1, e()
	}
	if k >= end-start {
		return -1, wm.SumAll(start, end)
	}
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if len(wm.presum) == 0 {
		panic("sumData is not provided")
	}
	count := int32(0)
	sum := e()
	res := WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		c := int32(0)
		if f {
			c = (end - start) - (r0 - l0)
		} else {
			c = r0 - l0
		}
		if count+c > k {
			if !f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		} else {
			var s WmSum
			if f {
				s = wm._get(d, start+wm.mid[d]-l0, end+wm.mid[d]-r0)
			} else {
				s = wm._get(d, l0, r0)
			}
			count += c
			sum = op(sum, s)
			res |= 1 << d
			if !f {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			} else {
				start, end = l0, r0
			}
		}
	}
	sum = op(sum, wm._get(0, start, start+k-count))
	if wm.compress {
		res = wm.key[res]
	}
	return res, sum
}

// upper: 向上取中位数还是向下取中位数.
func (wm *WaveletMatrixWithSum) Median(start, end int32, upper bool, xorVal WmValue) WmValue {
	n := end - start
	var k int32
	if upper {
		k = n / 2
	} else {
		k = (n - 1) / 2
	}
	return wm.Kth(start, end, k, xorVal)
}

func (wm *WaveletMatrixWithSum) SumAll(start, end int32) WmSum {
	return wm._get(wm.log, start, end)
}

// 使得predicate(count, sum)为true的最大的(count, sum).
func (wm *WaveletMatrixWithSum) MaxRight(predicate func(int32, WmSum) bool, start, end int32, xorVal WmValue) (int32, WmSum) {
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if start == end {
		return end - start, e()
	}
	if s := wm._get(wm.log, start, end); predicate(end-start, s) {
		return end - start, s
	}
	count := int32(0)
	sum := e()
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		c := int32(0)
		if f {
			c = (end - start) - (r0 - l0)
		} else {
			c = (r0 - l0)
		}
		var s WmSum
		if f {
			s = wm._get(d, start+wm.mid[d]-l0, end+wm.mid[d]-r0)
		} else {
			s = wm._get(d, l0, r0)
		}
		if tmp := op(sum, s); predicate(count+c, tmp) {
			count += c
			sum = tmp
			if f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		} else {
			if !f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		}
	}
	k := wm._binarySearch(func(k int32) bool {
		return predicate(count+k, op(sum, wm._get(0, start, start+k)))
	}, 0, end-start)
	count += k
	sum = op(sum, wm._get(0, start, start+k))
	return count, sum
}

// [start, end) 中小于等于 x 的数中最大的数
//
//	如果不存在则返回-INF
func (wm *WaveletMatrixWithSum) Floor(start, end int32, x WmValue, xor WmValue) WmValue {
	less := wm.CountPrefix(start, end, x, xor)
	if less == 0 {
		return -INF
	}
	res := wm.Kth(start, end, less-1, xor)
	return res
}

// [start, end) 中大于等于 x 的数中最小的数
//
//	如果不存在则返回INF
func (wm *WaveletMatrixWithSum) Ceil(start, end int32, x WmValue, xor WmValue) int {
	less := wm.CountPrefix(start, end, x, xor)
	if less == end-start {
		return INF
	}
	res := wm.Kth(start, end, less, xor)
	return res
}

func (wm *WaveletMatrixWithSum) CountSegments(segments [][2]int32, a, b WmValue, xorVal WmValue) int32 {
	res := int32(0)
	for _, seg := range segments {
		res += wm.CountRange(seg[0], seg[1], a, b, xorVal)
	}
	return res
}

func (wm *WaveletMatrixWithSum) KthSegments(segments [][2]int32, k int32, xorVal WmValue) WmValue {
	totalLen := int32(0)
	for _, seg := range segments {
		totalLen += seg[1] - seg[0]
	}
	count := int32(0)
	res := WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		c := int32(0)
		for _, seg := range segments {
			l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
			if f {
				c += (seg[1] - seg[0]) - (r0 - l0)
			} else {
				c += r0 - l0
			}
		}
		if count+c > k {
			for _, seg := range segments {
				l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
				if !f {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		} else {
			count += c
			res |= 1 << d
			for _, seg := range segments {
				l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
				if f {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		}
	}
	if wm.compress {
		res = wm.key[res]
	}
	return res
}

func (wm *WaveletMatrixWithSum) KthValueAndSumSegments(segments [][2]int32, k int32, xorVal WmValue) (WmValue, WmSum) {
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	totalLen := int32(0)
	for _, seg := range segments {
		totalLen += seg[1] - seg[0]
	}
	if k >= totalLen {
		return -1, wm.SumAllSegments(segments)
	}
	count := int32(0)
	sum := e()
	res := WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		c := int32(0)
		for _, seg := range segments {
			l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
			if f {
				c += (seg[1] - seg[0]) - (r0 - l0)
			} else {
				c += r0 - l0
			}
		}
		if count+c > k {
			for _, seg := range segments {
				l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
				if !f {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		} else {
			count += c
			res |= 1 << d
			for _, seg := range segments {
				l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
				var s WmSum
				if f {
					s = wm._get(d, seg[0]+wm.mid[d]-l0, seg[1]+wm.mid[d]-r0)
				} else {
					s = wm._get(d, l0, r0)
				}
				sum = op(sum, s)
				if !f {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				} else {
					seg[0], seg[1] = l0, r0
				}
			}
		}
	}
	for _, seg := range segments {
		t := min32(seg[1]-seg[0], k-count)
		sum = op(sum, wm._get(0, seg[0], seg[0]+t))
		count += t
	}
	if wm.compress {
		res = wm.key[res]
	}
	return res, sum
}

func (wm *WaveletMatrixWithSum) MedianSegments(segments [][2]int32, upper bool, xorVal WmValue) WmValue {
	n := int32(0)
	for _, seg := range segments {
		n += seg[1] - seg[0]
	}
	var k int32
	if upper {
		k = n / 2
	} else {
		k = (n - 1) / 2
	}
	return wm.KthSegments(segments, k, xorVal)
}

func (wm *WaveletMatrixWithSum) SumAllSegments(segments [][2]int32) WmSum {
	sum := e()
	for _, seg := range segments {
		sum = op(sum, wm._get(wm.log, seg[0], seg[1]))
	}
	return sum
}

func (wm *WaveletMatrixWithSum) SumSegments(segments [][2]int32, k1, k2 int32, xorVal WmValue) WmSum {
	if k1 >= k2 {
		return e()
	}
	add := wm._prefixSumSegments(segments, k2, xorVal)
	sub := wm._prefixSumSegments(segments, k1, xorVal)
	return op(add, inv(sub))
}

func (wm *WaveletMatrixWithSum) _prefixSumSegments(segments [][2]int32, k int32, xor WmValue) WmSum {
	_, sum := wm.KthValueAndSumSegments(segments, k, xor)
	return sum
}

func (wm *WaveletMatrixWithSum) _get(d, l, r int32) WmSum {
	return op(inv(wm.presum[d][l]), wm.presum[d][r])
}

func (wm *WaveletMatrixWithSum) _argSort(nums []WmValue) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func (wm *WaveletMatrixWithSum) _maxs(nums []WmValue) WmValue {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}

func (wm *WaveletMatrixWithSum) _lowerBound(nums []WmValue, target WmValue) WmValue {
	left, right := int32(0), int32(len(nums)-1)
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return WmValue(left)
}

func (wm *WaveletMatrixWithSum) _binarySearch(f func(int32) bool, ok, ng int32) int32 {
	for abs32(ok-ng) > 1 {
		x := (ok + ng) >> 1
		if f(x) {
			ok = x
		} else {
			ng = x
		}
	}
	return ok
}

type BitVector struct {
	bits   []uint64
	preSum []int32
}

func NewBitVector(n int32) *BitVector {
	return &BitVector{bits: make([]uint64, n>>6+1), preSum: make([]int32, n>>6+1)}
}

func (bv *BitVector) Set(i int32) {
	bv.bits[i>>6] |= 1 << (i & 63)
}

func (bv *BitVector) Build() {
	for i := 0; i < len(bv.bits)-1; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bits[i]))
	}
}

func (bv *BitVector) Rank(k int32, f bool) int32 {
	m, s := bv.bits[k>>6], bv.preSum[k>>6]
	res := s + int32(bits.OnesCount64(m&((1<<(k&63))-1)))
	if f {
		return res
	}
	return k - res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
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
