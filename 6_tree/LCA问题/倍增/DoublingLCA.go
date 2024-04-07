// https://github.dev/EndlessCheng/codeforces-go/blob/39cf48ce938640075a7b4885326f3d77929be143/copypasta/graph_tree.go#L668
// 倍增法求LCA

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	jump()
}

func jump() {
	// https://judge.yosupo.jp/problem/jump_on_tree
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	tree := make([][][2]int, n)
	for i := int(0); i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], [2]int{v, 0})
		tree[v] = append(tree[v], [2]int{u, 0})
	}
	D := NewLCADoubling(tree, []int{0})

	for i := int(0); i < q; i++ {
		var from, to, k int
		fmt.Fscan(in, &from, &to, &k)
		fmt.Fprintln(out, D.Jump(from, to, k))
	}
}

// https://leetcode.cn/problems/closest-node-to-path-in-tree/
// edges = [[0,1],[0,2],[0,3],[1,4],[2,5],[2,6]]
func closestNode(n int, edges [][]int, query [][]int) []int {
	tree := make([][][2]int, n)
	for _, edge := range edges {
		tree[edge[0]] = append(tree[edge[0]], [2]int{edge[1], 0})
		tree[edge[1]] = append(tree[edge[1]], [2]int{edge[0], 0})
	}

	lca := NewLCADoubling(tree, []int{0})
	res := make([]int, len(query))
	for i, q := range query {
		// lca里最深的那个
		res[i] = maxWithKey(
			func(x int) int { return int(lca.Depth[x]) },
			lca.QueryLCA(q[0], q[1]),
			lca.QueryLCA(q[0], q[2]),
			lca.QueryLCA(q[1], q[2]),
		)
	}

	return res
}

const INF int = 1e18

type LCADoubling struct {
	Tree          [][][2]int
	Depth         []int32
	DepthWeighted []int
	n             int
	log           int
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
		log:           bits.Len32(uint32(n)) - 1,
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
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.dp[i][root132], lca.dp[i][root232]; a != b {
			root132, root232 = a, b
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
	for i := lca.log; i >= 0; i-- { // upToDepth
		if (lca.Depth[root132]-toDepth)&(1<<i) > 0 {
			res = max(res, lca.dpWeight1[i][root132])
			root132 = lca.dp[i][root132]
		}
	}
	if root132 == root232 {
		return res
	}
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.dp[i][root132], lca.dp[i][root232]; a != b {
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
	for i := lca.log; i >= 0; i-- { // upToDepth
		if (lca.Depth[root132]-toDepth)&(1<<i) > 0 {
			res = min(res, lca.dpWeight2[i][root132])
			root132 = lca.dp[i][root132]
		}
	}
	if root132 == root232 {
		return res
	}
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.dp[i][root132], lca.dp[i][root232]; a != b {
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
	for i := lca.log; i >= 0; i-- {
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
	log := lca.log
	dp, dpWeight1, dpWeight2 = make([][]int32, log+1), make([][]int, log+1), make([][]int, log+1)
	for i := 0; i < log+1; i++ {
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
	for i := 0; i < lca.log; i++ {
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

func maxWithKey(key func(x int) int, args ...int) int {
	max := args[0]
	for _, v := range args[1:] {
		if key(max) < key(v) {
			max = v
		}
	}
	return max
}
