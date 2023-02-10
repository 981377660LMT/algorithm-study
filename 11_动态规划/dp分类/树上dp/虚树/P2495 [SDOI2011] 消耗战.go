// https://www.luogu.com.cn/problem/P2495
// 在一场战争中，战场由n个岛屿和n-1个桥梁组成，保证每两个岛屿间有且仅有一条路径可达。
// 现在，我军已经侦查到敌军的总部在编号为0的岛屿，而且他们已经没有足够多的能源维系战斗，
// 我军胜利在望。已知在其他k个岛屿上有丰富能源，
// 为了防止敌军获取能源，我军的任务是炸毁一些桥梁，使得敌军不能到达任何能源丰富的岛屿。
// 由于不同桥梁的材质和结构不同，所以炸毁不同的桥梁有不同的代价，
// 我军希望在满足目标的同时使得总代价最小。
// 侦查部门还发现，敌军有一台神秘机器。即使我军切断所有能源之后，他们也可以用那台机器。
// 机器产生的效果不仅仅会修复所有我军炸毁的桥梁，
// 而且会重新随机资源分布（但可以保证的是，资源不会分布到0号岛屿上)。
// 不过侦查部门还发现了这台机器只能够使用q次，所以我们只需要把每次任务完成即可。

// 虚树上dp
// 第一行一个整数n，表示岛屿数量。
// 接下来n-1行，每行三个整数u, o, w，表示u号岛屿和v号岛屿由一条代价为w的桥梁直接相连。
// 第n＋1行，一个整数q，代表敌方机器能使用的次数。
// 接下来q行，第i行一个整数k;，代表第i次后，有k个岛屿资源丰富。
// 接下来k个整数h1, h2,..., hk;，表示资源丰富岛屿的编号。
// n<=2e5 q<=5e5 ∑ki<=5e5

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func solve(n int, edges [][3]int, q int, groups [][]int) []int {
	adjList := make([][]WeightedEdge, n)
	edges_ := make([][]int, 0, len(edges))
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		edges_ = append(edges_, []int{u, v})
		adjList[u] = append(adjList[u], WeightedEdge{v, w})
		adjList[v] = append(adjList[v], WeightedEdge{u, w})
	}

	lca := NewLCA(n, adjList, 0)
	A := NewAuxiliaryTree(n, edges_)
	visited := make([]bool, n) // 标记是否已经被分组
	res := make([]int, q)
	for i := 0; i < q; i++ {
		points := groups[i]
		for _, p := range points {
			visited[p] = true
		}

		points = append(points, 0)    // 加上根节点组成虚树
		tree, root := A.Query(points) // 虚树有向邻接表, 根节点
		// dp[i]表示i和以i为根的子树中的关键点都不相连的最小代价
		// 如果子节点是关键点，dp[i] += minWeight[i][child]
		// 如果子节点不是关键点，dp[i] += min(dp[child], minWeight[i][child])
		var dfs2 func(int, int) int
		dfs2 = func(cur, pre int) int {
			res := 0
			for _, next := range tree[cur] {
				if next == pre {
					continue
				}
				nextRes := dfs2(next, cur)
				minWeight := lca.QueryMinWeight(cur, next, true)
				if visited[next] {
					res += minWeight
				} else {
					res += min(nextRes, minWeight)
				}
			}
			return res
		}

		res[i] = dfs2(root, -1)
		for _, p := range points {
			visited[p] = false
		}
	}

	return res
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	edges := make([][3]int, 0, n-1)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges = append(edges, [3]int{u, v, w})
	}

	var q int
	fmt.Fscan(in, &q)
	points := make([][]int, 0, q)
	for i := 0; i < q; i++ {
		var count int
		fmt.Fscan(in, &count)
		curPoints := make([]int, 0, count)
		for j := 0; j < count; j++ {
			var p int
			fmt.Fscan(in, &p)
			curPoints = append(curPoints, p-1)
		}
		points = append(points, curPoints)
	}

	res := solve(n, edges, q, points)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type AuxiliaryTree struct {
	g             [][]int // 原图邻接表
	g0            [][]int // 虚树邻接表
	s             []int
	fs, ls, depth []int

	lg []int
	st [][]int
}

// 给定顶点个数n和无向边集(u,v)构建.
//  O(nlogn)
func NewAuxiliaryTree(n int, edges [][]int) *AuxiliaryTree {
	g := make([][]int, n)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
		g[e[1]] = append(g[e[1]], e[0])
	}
	res := &AuxiliaryTree{
		g:     g,
		g0:    make([][]int, n),
		s:     []int{},
		fs:    make([]int, n),
		ls:    make([]int, n),
		depth: make([]int, n),
	}

	res.dfs(0, -1, 0)
	res.buildSt()
	return res
}

// 指定点集,返回虚树的(有向图邻接表,虚树的根).
//  如果虚树不存在`(len(points<=1))`,返回空邻接表和-1.
//  O(klogk) 构建虚树.
func (t *AuxiliaryTree) Query(points []int) ([][]int, int) {
	k := len(points)
	if k <= 1 {
		return nil, -1
	}

	points = append(points[:0:0], points...)
	sort.Slice(points, func(i, j int) bool {
		return t.fs[points[i]] < t.fs[points[j]]
	})

	for i := 0; i < k-1; i++ {
		x, y := t.fs[points[i]], t.fs[points[i+1]]
		l := t.lg[y-x+1]
		w := t.st[l][x]
		if t.depth[t.st[l][y-(1<<l)+1]] < t.depth[t.st[l][x]] {
			w = t.st[l][y-(1<<l)+1]
		}
		points = append(points, w)
	}

	sort.Slice(points, func(i, j int) bool {
		return t.fs[points[i]] < t.fs[points[j]]
	})

	stk := []int{}
	pre := -1
	root := -1
	for _, v := range points {
		if pre == v {
			continue
		}
		for len(stk) > 0 && t.ls[stk[len(stk)-1]] < t.fs[v] {
			stk = stk[:len(stk)-1]
		}
		if len(stk) > 0 {
			parent := stk[len(stk)-1]
			t.g0[parent] = append(t.g0[parent], v)
			if root == -1 {
				root = parent
			}
		}

		t.g0[v] = t.g0[v][:0]
		stk = append(stk, v)
		pre = v
	}

	return t.g0, root
}

func (t *AuxiliaryTree) dfs(v, p, d int) {
	t.depth[v] = d
	t.fs[v] = len(t.s)
	t.s = append(t.s, v)
	for _, w := range t.g[v] {
		if w == p {
			continue
		}
		t.dfs(w, v, d+1)
		t.s = append(t.s, v)
	}
	t.ls[v] = len(t.s)
	t.s = append(t.s, v)
}

func (t *AuxiliaryTree) buildSt() {
	l := len(t.s)
	lg := make([]int, l+1)
	for i := 2; i <= l; i++ {
		lg[i] = lg[i>>1] + 1
	}
	st := make([][]int, lg[l]+1)
	for i := range st {
		st[i] = make([]int, l-(1<<i)+1)
		for j := range st[i] {
			st[i][j] = l
		}
	}

	copy(st[0], t.s)
	b := 1
	for i := 0; i < lg[l]; i++ {
		st0, st1 := st[i], st[i+1]
		for j := 0; j < l-(b<<1)+1; j++ {
			st1[j] = st0[j]
			if t.depth[st0[j+b]] < t.depth[st0[j]] {
				st1[j] = st0[j+b]
			}
		}
		b <<= 1
	}

	t.lg = lg
	t.st = st
}

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

func maxWithKey(key func(x int) int, args ...int) int {
	max := args[0]
	for _, v := range args[1:] {
		if key(max) < key(v) {
			max = v
		}
	}
	return max
}
