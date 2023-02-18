// カッコ列をグラフにする。各頂点の範囲を表す配列 LR も作る。
// 全体を表す根ノードも作って、N+1頂点。

// (())() → 得到(n//2)+1 个结点,其中0号结点是虚拟根结点
// graph: [[1 3] [2] [] []] (有向邻接表)
// leftRight: [[0 5] [0 3] [1 2] [4 5]] (每个顶点的括号序)
//
//           0 (0,5)
//          / \
//   (0,3) 1   3 (4,5)
//        /
// (1,2) 2

// 有效的括号序列形成的树，结合LCA

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// fmt.Println(BracketTree("(())()"))
	// https://yukicoder.me/problems/no/1778
	// 给定一个有效的括号序列,每次可以删除一段匹配的括号
	// 给定q个查询,每个查询形如 [start1,start2]
	// 求包含这两段括号的区间中最靠内部的一段区间

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)

	tree, leftRight := BracketTree(s)

	idToRoot := make([]int, n) // 每个起点/终点位置对应的树节点 [0, n)
	for i := 1; i < len(leftRight); i++ {
		left, right := leftRight[i][0], leftRight[i][1]
		idToRoot[left] = i
		idToRoot[right] = i
	}

	lca := NewLCA((n/2)+1, 0)
	for i := 0; i < len(tree); i++ {
		for _, j := range tree[i] {
			lca.AddEdge(i, j, 1)
		}
	}
	lca.Build()

	for i := 0; i < q; i++ {
		var start, end int // 1<=start<end<=n
		fmt.Fscan(in, &start, &end)
		start--
		end--
		root1, root2 := idToRoot[start], idToRoot[end]
		lca_ := lca.QueryLCA(root1, root2)
		if lca_ == 0 { // 不被包含
			fmt.Fprintln(out, -1)
		} else {
			left, right := leftRight[lca_][0], leftRight[lca_][1]
			fmt.Fprintln(out, left+1, right+1)
		}
	}
}

// BracketTree :有效括号序列转换成树.
//  (())() → 得到 `(n/2)+1` 个结点,其中0号结点是虚拟根结点.
//  tree: [[1 3] [2] [] []] (有向邻接表)
//  leftRight: [[0 5] [0 3] [1 2] [4 5]] (每个顶点的括号序/欧拉序)
//
//            0 (0,5)
//           / \
//    (0,3) 1   3 (4,5)
//         /
//  (1,2) 2
func BracketTree(s string) (tree [][]int, leftRight [][2]int) {
	n := len(s) / 2
	tree = make([][]int, n+1)
	leftRight = make([][2]int, n+1)
	now, nxt := 0, 1
	leftRight[0] = [2]int{0, len(s)}
	par := make([]int, n+1)
	for i := range par {
		par[i] = -1
	}

	for i := range s {
		if s[i] == '(' {
			tree[now] = append(tree[now], nxt)
			par[nxt] = now
			leftRight[nxt][0] = i
			now = nxt
			nxt++
		}
		if s[i] == ')' {
			leftRight[now][1] = i
			now = par[now]
		}
	}

	return
}

const INF int = 1e18

type LCA struct {
	Depth      []int
	Tree       [][]edge
	n          int
	root       int
	bitLen     int
	distToRoot []int
	// 节点j向上跳2^i步的父节点
	dp [][]int
	// 节点j向上跳2^i步经过的最大边权
	dpWeight1 [][]int
	// 节点j向上跳2^i步经过的最小边权
	dpWeight2 [][]int
}

type edge struct{ to, weight int }

func NewLCA(n int, root int) *LCA {
	lca := &LCA{
		Tree:       make([][]edge, n),
		Depth:      make([]int, n),
		n:          n,
		root:       root,
		bitLen:     bits.Len(uint(n)),
		distToRoot: make([]int, n),
	}

	return lca
}

// 添加权值为w的无向边(u, v).
func (lca *LCA) AddEdge(u, v, w int) {
	lca.Tree[u] = append(lca.Tree[u], edge{v, w})
	lca.Tree[v] = append(lca.Tree[v], edge{u, w})
}

func (lca *LCA) Build() {
	lca.dp, lca.dpWeight1, lca.dpWeight2 = makeDp(lca)
	lca.dfsAndInitDp(lca.root, -1, 0, 0)
	lca.fillDp()
}

// 查询树节点两点的最近公共祖先
func (lca *LCA) QueryLCA(root1, root2 int) int {
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}

	root1 = lca.UpToDepth(root1, lca.Depth[root2])
	if root1 == root2 {
		return root1
	}

	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root1] != lca.dp[i][root2] {
			root1 = lca.dp[i][root1]
			root2 = lca.dp[i][root2]
		}
	}

	return lca.dp[0][root1]
}

// 查询树节点两点间距离
//  weighted: 是否将边权计入距离
func (lca *LCA) QueryDist(root1, root2 int, weighted bool) int {
	if weighted {
		return lca.distToRoot[root1] + lca.distToRoot[root2] - 2*lca.distToRoot[lca.QueryLCA(root1, root2)]
	}
	return lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.QueryLCA(root1, root2)]
}

// 查询树节点两点路径上最大边权(倍增的时候维护其他属性)
//  isEdge 为true表示查询路径上边权,为false表示查询路径上点权
func (lca *LCA) QueryMaxWeight(root1, root2 int, isEdge bool) int {
	res := -INF
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}

	toDepth := lca.Depth[root2]
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.Depth[root1]-toDepth)&(1<<i) > 0 {
			res = max(res, lca.dpWeight1[i][root1])
			root1 = lca.dp[i][root1]
		}
	}

	if root1 == root2 {
		return res
	}

	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root1] != lca.dp[i][root2] {
			res = max(res, max(lca.dpWeight1[i][root1], lca.dpWeight1[i][root2]))
			root1 = lca.dp[i][root1]
			root2 = lca.dp[i][root2]
		}
	}

	res = max(res, max(lca.dpWeight1[0][root1], lca.dpWeight1[0][root2]))

	if !isEdge {
		lca_ := lca.dp[0][root1]
		res = max(res, lca.dpWeight1[0][lca_])
	}

	return res
}

// 查询树节点两点路径上最小边权(倍增的时候维护其他属性)
//  isEdge 为true表示查询路径上边权,为false表示查询路径上点权
func (lca *LCA) QueryMinWeight(root1, root2 int, isEdge bool) int {
	res := INF
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}

	toDepth := lca.Depth[root2]
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.Depth[root1]-toDepth)&(1<<i) > 0 {
			res = min(res, lca.dpWeight2[i][root1])
			root1 = lca.dp[i][root1]
		}
	}

	if root1 == root2 {
		return res
	}

	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root1] != lca.dp[i][root2] {
			res = min(res, min(lca.dpWeight2[i][root1], lca.dpWeight2[i][root2]))
			root1 = lca.dp[i][root1]
			root2 = lca.dp[i][root2]
		}
	}

	res = min(res, min(lca.dpWeight2[0][root1], lca.dpWeight2[0][root2]))

	if !isEdge {
		lca_ := lca.dp[0][root1]
		res = min(res, lca.dpWeight2[0][lca_])
	}

	return res
}

// 查询树节点root的第k个祖先(向上跳k步),如果不存在这样的祖先节点,返回 -1
func (lca *LCA) QueryKthAncestor(root, k int) int {
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			root = lca.dp[bit][root]
			if root == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return root
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点
func (lca *LCA) UpToDepth(root, toDepth int) int {
	if toDepth >= lca.Depth[root] {
		return root
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if (lca.Depth[root]-toDepth)&(1<<i) > 0 {
			root = lca.dp[i][root]
		}
	}
	return root
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *LCA) Jump(start, target, step int) int {
	lca_ := lca.QueryLCA(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return lca.QueryKthAncestor(start, step)
	}
	return lca.QueryKthAncestor(target, dist-step)
}

func (lca *LCA) dfsAndInitDp(cur, pre, dep, dist int) {
	lca.Depth[cur] = dep
	lca.dp[0][cur] = pre
	lca.distToRoot[cur] = dist
	for _, e := range lca.Tree[cur] {
		if next := e.to; next != pre {
			lca.dpWeight1[0][next] = e.weight
			lca.dpWeight2[0][next] = e.weight
			lca.dfsAndInitDp(next, cur, dep+1, dist+e.weight)
		}
	}
}

func makeDp(lca *LCA) (dp, dpWeight1, dpWeight2 [][]int) {
	dp, dpWeight1, dpWeight2 = make([][]int, lca.bitLen), make([][]int, lca.bitLen), make([][]int, lca.bitLen)
	for i := 0; i < lca.bitLen; i++ {
		dp[i], dpWeight1[i], dpWeight2[i] = make([]int, lca.n), make([]int, lca.n), make([]int, lca.n)
		for j := 0; j < lca.n; j++ {
			dp[i][j] = -1
			dpWeight1[i][j] = -INF
			dpWeight2[i][j] = INF
		}
	}
	return
}

func (lca *LCA) fillDp() {
	for i := 0; i < lca.bitLen-1; i++ {
		for j := 0; j < lca.n; j++ {
			if lca.dp[i][j] == -1 {
				lca.dp[i+1][j] = -1
			} else {
				lca.dp[i+1][j] = lca.dp[i][lca.dp[i][j]]
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

func maxWithKey(key func(x int) int, args ...int) int {
	max := args[0]
	for _, v := range args[1:] {
		if key(max) < key(v) {
			max = v
		}
	}
	return max
}
