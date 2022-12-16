package main

import "math/bits"

func closestNode(n int, edges [][]int, query [][]int) []int {
	tree := make([][]Edge, n)
	for _, edge := range edges {
		tree[edge[0]] = append(tree[edge[0]], Edge{edge[1], 0})
		tree[edge[1]] = append(tree[edge[1]], Edge{edge[0], 0})
	}

	lca := NewLCA(n, tree, 0)
	res := make([]int, len(query))
	for i, q := range query {
		// lca里最深的那个
		res[i] = maxWithKey(
			func(x int) int { return lca.depth[x] },
			lca.QueryLCA(q[0], q[1]),
			lca.QueryLCA(q[0], q[2]),
			lca.QueryLCA(q[1], q[2]),
		)
	}

	return res
}

const INF int = 1e18

type Edge struct{ to, weight int }

type LCA struct {
	n      int
	bitLen int
	tree   [][]Edge
	depth  []int

	// dp[i][j] 表示节点j向上跳2^i步的父节点
	dp [][]int
}

func NewLCA(n int, tree [][]Edge, root int) *LCA {
	lca := &LCA{
		n:      n,
		bitLen: bits.Len(uint(n)),
		tree:   tree,
		depth:  make([]int, n),
	}

	lca.dp = makeDp(lca)
	lca.dfsAndInitDp(root, -1, 0)
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
func (lca *LCA) QueryDist(root1, root2 int) int {
	return lca.depth[root1] + lca.depth[root2] - 2*lca.depth[lca.QueryLCA(root1, root2)]
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

func (lca *LCA) dfsAndInitDp(cur, pre, dep int) {
	lca.depth[cur] = dep
	lca.dp[0][cur] = pre
	for _, e := range lca.tree[cur] {
		if next := e.to; next != pre {
			lca.dfsAndInitDp(next, cur, dep+1)
		}
	}
}

func makeDp(lca *LCA) (dp [][]int) {
	dp = make([][]int, lca.bitLen)
	for i := 0; i < lca.bitLen; i++ {
		dp[i] = make([]int, lca.n)
		for j := 0; j < lca.n; j++ {
			dp[i][j] = -1
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

func maxWithKey(key func(x int) int, args ...int) int {
	if len(args) == 0 {
		panic("max: empty args")
	}
	max := args[0]
	for _, v := range args[1:] {
		if key(max) < key(v) {
			max = v
		}
	}
	return max
}
