// https://leetcode.cn/problems/kth-ancestor-of-a-tree-node/
package main

import (
	"fmt"
	"math/bits"
)

func main() {
	// 	["TreeAncestor","getKthAncestor","getKthAncestor","getKthAncestor"]
	// [[7,[-1,0,0,1,1,2,2]],[3,1],[5,2],[6,3]]

	treeAncestor := Constructor(7, []int{-1, 0, 0, 1, 1, 2, 2})
	fmt.Println(treeAncestor.GetKthAncestor(3, 1))
	fmt.Println(treeAncestor.GetKthAncestor(5, 2))
	fmt.Println(treeAncestor.GetKthAncestor(6, 3))
}

type TreeAncestor struct {
	lca *LCA
}

func Constructor(n int, parent []int) TreeAncestor {
	tree := make([][]WeightedEdge, n)
	for i := 1; i < n; i++ {
		tree[parent[i]] = append(tree[parent[i]], WeightedEdge{i, 1})
	}
	return TreeAncestor{lca: NewLCA(n, tree, 0)}
}

func (this *TreeAncestor) GetKthAncestor(node int, k int) int {
	return this.lca.QueryKthAncestor(node, k)
}

/**
 * Your TreeAncestor object will be instantiated and called as such:
 * obj := Constructor(n, parent);
 * param_1 := obj.GetKthAncestor(node,k);
 */

const INF int = 1e18

type WeightedEdge struct{ to, weight int }

type LCA struct {
	n          int
	bitLen     int
	tree       [][]WeightedEdge
	depth      []int
	distToRoot []int
	// 节点j向上跳2^i步的父节点
	dp [][]int
	// 节点j向上跳2^i步经过的最大边权
	dpWeight1 [][]int
	// 节点j向上跳2^i步经过的最小边权
	dpWeight2 [][]int
}

func NewLCA(n int, tree [][]WeightedEdge, root int) *LCA {
	lca := &LCA{
		n:          n,
		bitLen:     bits.Len(uint(n)),
		tree:       tree,
		depth:      make([]int, n),
		distToRoot: make([]int, n),
	}

	lca.dp, lca.dpWeight1, lca.dpWeight2 = makeDp(lca)
	lca.dfsAndInitDp(root, -1, 0, 0) // !如果传入的是森林,需要visited遍历各个连通分量 (1697. 检查边长度限制的路径是否存在-在线)
	lca.fillDp()
	return lca
}

// 查询树节点两点的最近公共祖先
func (lca *LCA) QueryLCA(root1, root2 int) int {
	if lca.depth[root1] < lca.depth[root2] {
		root1, root2 = root2, root1
	}

	root1 = lca.UpToDepth(root1, lca.depth[root2])
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
	return lca.depth[root1] + lca.depth[root2] - 2*lca.depth[lca.QueryLCA(root1, root2)]
}

// 查询树节点两点路径上最大边权(倍增的时候维护其他属性)
//  isEdge 为true表示查询路径上边权,为false表示查询路径上点权
func (lca *LCA) QueryMaxWeight(root1, root2 int, isEdge bool) int {
	res := -INF
	if lca.depth[root1] < lca.depth[root2] {
		root1, root2 = root2, root1
	}

	toDepth := lca.depth[root2]
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.depth[root1]-toDepth)&(1<<i) > 0 {
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
	if lca.depth[root1] < lca.depth[root2] {
		root1, root2 = root2, root1
	}

	toDepth := lca.depth[root2]
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.depth[root1]-toDepth)&(1<<i) > 0 {
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
	if toDepth >= lca.depth[root] {
		return root
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if (lca.depth[root]-toDepth)&(1<<i) > 0 {
			root = lca.dp[i][root]
		}
	}
	return root
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *LCA) Jump(start, target, step int) int {
	lca_ := lca.QueryLCA(start, target)
	dep1, dep2, deplca := lca.depth[start], lca.depth[target], lca.depth[lca_]
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
	lca.depth[cur] = dep
	lca.dp[0][cur] = pre
	lca.distToRoot[cur] = dist
	for _, e := range lca.tree[cur] {
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
