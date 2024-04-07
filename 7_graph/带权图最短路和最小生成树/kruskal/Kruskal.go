package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	CF609E()
	// P3366()
}

// E. Minimum spanning tree for each edge
// https://www.luogu.com.cn/problem/CF609E
// 如果对于一个最小生成树中要求必须包括第 i 条边,那么这个最小生成树的权值是多少?
// n, m <= 2*10^5
//
// !对于第 i 个查询，先加上第 i 条边，
// !接着在构成的这个环中减去一条边权最大的边（第 i 条边除外）。
func CF609E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	mstCost, inMst, _ := Kruskal(n, edges)
	tree := make([][][2]int, n)
	for ei, b := range inMst {
		if b {
			u, v, w := edges[ei][0], edges[ei][1], edges[ei][2]
			tree[u] = append(tree[u], [2]int{v, w})
			tree[v] = append(tree[v], [2]int{u, w})
		}
	}
	lca := NewLCADoubling(tree, []int{0})

	res := make([]int, m)
	for i := 0; i < m; i++ {
		if inMst[i] {
			res[i] = mstCost
			continue
		}
		u, v, w := edges[i][0], edges[i][1], edges[i][2]
		maxWeight := lca.QueryMaxWeight(u, v, true)
		res[i] = mstCost + w - maxWeight
	}

	for _, v := range res {
		fmt.Println(v)
	}
}

// P3366 【模板】最小生成树
func P3366() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	mstCost, _, isTree := Kruskal(n, edges)
	if isTree {
		fmt.Fprintln(out, mstCost)
	} else {
		fmt.Fprintln(out, "orz")
	}
}

// !给定无向图的边，求出一个最小生成树(如果不存在,则求出的是森林中的多个最小生成树)
// mstCost: 最小生成树(森林)的权值之和
// inMst: 是否在最小生成树(森林)中
// isTree: 是否是树
func Kruskal(n int, edges [][3]int) (mstCost int, inMst []bool, isTree bool) {
	order := make([]int, len(edges))
	for i := range edges {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return edges[order[i]][2] < edges[order[j]][2] })

	uf := NewUf(n)
	inMst = make([]bool, len(edges))
	count := 0
	for _, ei := range order {
		u, v, w := edges[ei][0], edges[ei][1], edges[ei][2]
		if uf.Union(u, v) {
			inMst[ei] = true
			mstCost += w
			count++
			if count == n-1 {
				isTree = true
				return
			}
		}
	}
	return
}

type Uf struct {
	n    int
	data []int
}

func NewUf(n int) *Uf {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &Uf{n: n, data: data}
}

func (ufa *Uf) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
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

const INF int = 1e18

type LCADoubling struct {
	Tree          [][][2]int
	Depth         []int32
	DepthWeighted []int
	n             int
	bitLen        int
	dp            [][]int32 // 节点j向上跳2^i步的父节点
	dpWeight1     [][]int   // 节点j向上跳2^i步经过的最大边权
	dpWeight2     [][]int   // 节点j向上跳2^i步经过的最小边权
}

func NewLCADoubling(tree [][][2]int, roots []int) *LCADoubling {
	n := len(tree)
	depth := make([]int32, n)
	lca := &LCADoubling{
		Tree:          tree,
		Depth:         depth,
		DepthWeighted: make([]int, n),
		n:             n,
		bitLen:        bits.Len(uint(n)),
	}
	lca.dp, lca.dpWeight1, lca.dpWeight2 = makeDp(lca)
	for _, root := range roots {
		lca.dfsAndInitDp(int32(root), -1, 0, 0)
	}
	lca.fillDp()
	return lca
}

// 查询树节点两点的最近公共祖先
func (lca *LCADoubling) QueryLCA(root1, root2 int) int {
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	root1 = lca.UpToDepth(root1, int(lca.Depth[root2]))
	if root1 == root2 {
		return root1
	}
	root132, root232 := int32(root1), int32(root2)
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root132] != lca.dp[i][root232] {
			root132 = lca.dp[i][root132]
			root232 = lca.dp[i][root232]
		}
	}
	return int(lca.dp[0][root132])
}

// 查询树节点两点间距离
//
//	weighted: 是否将边权计入距离
func (lca *LCADoubling) QueryDist(root1, root2 int, weighted bool) int {
	if weighted {
		return lca.DepthWeighted[root1] + lca.DepthWeighted[root2] - 2*lca.DepthWeighted[lca.QueryLCA(root1, root2)]
	}
	return int(lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.QueryLCA(root1, root2)])
}

// 查询树节点两点路径上最大边权(倍增的时候维护其他属性)
//
//	isEdge 为true表示查询路径上边权,为false表示查询路径上点权
func (lca *LCADoubling) QueryMaxWeight(root1, root2 int, isEdge bool) int {
	res := -INF
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	toDepth := lca.Depth[root2]
	root132, root232 := int32(root1), int32(root2)
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.Depth[root132]-toDepth)&(1<<i) > 0 {
			res = max(res, lca.dpWeight1[i][root132])
			root132 = lca.dp[i][root132]
		}
	}
	if root132 == root232 {
		return res
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root132] != lca.dp[i][root232] {
			res = max(res, max(lca.dpWeight1[i][root132], lca.dpWeight1[i][root232]))
			root132 = lca.dp[i][root132]
			root232 = lca.dp[i][root232]
		}
	}
	res = max(res, max(lca.dpWeight1[0][root132], lca.dpWeight1[0][root232]))
	if !isEdge {
		lca_ := lca.dp[0][root132]
		res = max(res, lca.dpWeight1[0][lca_])
	}
	return res
}

// 查询树节点两点路径上最小边权(倍增的时候维护其他属性)
//
//	isEdge 为true表示查询路径上边权,为false表示查询路径上点权
func (lca *LCADoubling) QueryMinWeight(root1, root2 int, isEdge bool) int {
	res := INF
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	toDepth := lca.Depth[root2]
	root132, root232 := int32(root1), int32(root2)
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.Depth[root132]-toDepth)&(1<<i) > 0 {
			res = min(res, lca.dpWeight2[i][root132])
			root132 = lca.dp[i][root132]
		}
	}
	if root132 == root232 {
		return res
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root132] != lca.dp[i][root232] {
			res = min(res, min(lca.dpWeight2[i][root132], lca.dpWeight2[i][root232]))
			root132 = lca.dp[i][root132]
			root232 = lca.dp[i][root232]
		}
	}
	res = min(res, min(lca.dpWeight2[0][root132], lca.dpWeight2[0][root232]))
	if !isEdge {
		lca_ := lca.dp[0][root132]
		res = min(res, lca.dpWeight2[0][lca_])
	}
	return res
}

// 查询树节点root的第k个祖先(向上跳k步),如果不存在这样的祖先节点,返回 -1
func (lca *LCADoubling) QueryKthAncestor(root, k int) int {
	root32 := int32(root)
	if k > int(lca.Depth[root32]) {
		return -1
	}
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			root32 = lca.dp[bit][root32]
			if root32 == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return int(root32)
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点
func (lca *LCADoubling) UpToDepth(root, toDepth int) int {
	toDepth32 := int32(toDepth)
	if toDepth32 >= lca.Depth[root] {
		return root
	}
	root32 := int32(root)
	for i := lca.bitLen - 1; i >= 0; i-- {
		if (lca.Depth[root32]-toDepth32)&(1<<i) > 0 {
			root32 = lca.dp[i][root32]
		}
	}
	return int(root32)
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *LCADoubling) Jump(start, target, step int) int {
	lca_ := lca.QueryLCA(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := int(dep1 + dep2 - 2*deplca)
	if step > dist {
		return -1
	}
	if step <= int(dep1-deplca) {
		return lca.QueryKthAncestor(start, step)
	}
	return lca.QueryKthAncestor(target, dist-step)
}

func (lca *LCADoubling) dfsAndInitDp(cur, pre, dep int32, dist int) {
	lca.Depth[cur] = dep
	lca.dp[0][cur] = pre
	lca.DepthWeighted[cur] = dist
	for _, e := range lca.Tree[cur] {
		next, weight := int32(e[0]), e[1]
		if next != pre {
			lca.dpWeight1[0][next] = weight
			lca.dpWeight2[0][next] = weight
			lca.dfsAndInitDp(next, cur, dep+1, dist+weight)
		}
	}
}

func makeDp(lca *LCADoubling) (dp [][]int32, dpWeight1, dpWeight2 [][]int) {
	dp, dpWeight1, dpWeight2 = make([][]int32, lca.bitLen), make([][]int, lca.bitLen), make([][]int, lca.bitLen)
	for i := 0; i < lca.bitLen; i++ {
		dp[i], dpWeight1[i], dpWeight2[i] = make([]int32, lca.n), make([]int, lca.n), make([]int, lca.n)
		for j := 0; j < lca.n; j++ {
			dp[i][j] = -1
			dpWeight1[i][j] = -INF
			dpWeight2[i][j] = INF
		}
	}
	return
}

func (lca *LCADoubling) fillDp() {
	for i := 0; i < lca.bitLen-1; i++ {
		for j := 0; j < lca.n; j++ {
			pre := lca.dp[i][j]
			if pre == -1 {
				lca.dp[i+1][j] = -1
			} else {
				lca.dp[i+1][j] = lca.dp[i][pre]
				lca.dpWeight1[i+1][j] = max(lca.dpWeight1[i][j], lca.dpWeight1[i][lca.dp[i][j]])
				lca.dpWeight2[i+1][j] = min(lca.dpWeight2[i][j], lca.dpWeight2[i][lca.dp[i][j]])
			}
		}
	}

	return
}

func max(a, b int) int {
	if a > b {
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
