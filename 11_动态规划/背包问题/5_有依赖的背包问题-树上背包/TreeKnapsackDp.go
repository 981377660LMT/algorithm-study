// 树上背包/树形背包/依赖背包
// 子树合并背包的复杂度证明 https://blog.csdn.net/lyd_7_29/article/details/79854245
// 复杂度(二乘木dp) https://leetcode.cn/circle/discuss/t7l62c/
// !通过子树大小优化背包循环上界

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// P2014选课()
	// CF815C()
	// abc287F()
	// ABC207F()
	// yuki196()
	// ABC207F()
	P2015二叉苹果树()
}

// P2014 [CTSC1997] 选课, O(nk) dp
// https://www.luogu.com.cn/problem/U53204#submit
// https://www.luogu.com.cn/problem/P2014
// https://zhuanlan.zhihu.com/p/103813542
func P2014选课() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)

	// 添加一个虚拟根节点，转换为选择k+1个节点的树形背包问题
	DUMMY := n
	graph := make([][]int, n+1)
	scores := make([]int, n+1)

	for i := 0; i < n; i++ {
		var parent, score int
		fmt.Fscan(in, &parent, &score)
		parent--
		if parent == -1 {
			parent = DUMMY
		}
		graph[parent] = append(graph[parent], i)
		scores[i] = score
	}

	dp := TreeKnapsackDpSquare2(graph, DUMMY, scores, k+1)
	fmt.Fprintln(out, dp[k+1])
}

// https://www.luogu.com.cn/problem/CF815C
// CF815C Karen and Supermarket
// 有n件商品,每件有价格ci,优惠券di,
// 对于i>=2,使用di的条件为:xi的优惠券需要被使用
// 问初始金钱为money时 最多能买多少件商品? n<=5000,ci,di,money<=1e9
// !dp[u][j][0/1]是以u为根子树中,u用或者不用优惠券时,选j件物品需要的最小代价
//
// 因为使用优惠劵就必须买这件商品，我们可以将优惠劵的关系看成一棵树。
func CF815C() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int = 1e9 + 10

	var n, money int
	fmt.Fscan(in, &n, &money)

	tree := make([][]int, n)
	prices := make([]int, n)
	coupons := make([]int, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			var price, coupon int
			fmt.Fscan(in, &price, &coupon)
			prices[i] = price
			coupons[i] = coupon
		} else {
			var price, coupon, parent int
			fmt.Fscan(in, &price, &coupon, &parent)
			parent--
			tree[parent] = append(tree[parent], i)
			prices[i] = price
			coupons[i] = coupon
		}
	}

	subSize := make([]int, n)

	var f func(int, int) [][2]int
	f = func(cur, pre int) [][2]int {
		subSize[cur] = 1
		dp := make([][2]int, n+1)
		for i := 0; i <= n; i++ {
			dp[i][0] = INF
			dp[i][1] = INF
		}
		dp[0][0] = 0
		dp[1][0] = prices[cur]
		dp[1][1] = prices[cur] - coupons[cur]

		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			nextDp := f(next, cur)
			// 当前节点可以不选, 所以是 j>=0
			for j := subSize[cur]; j >= 0; j-- {
				for w := 0; w <= subSize[next]; w++ {
					dp[j+w][0] = min(dp[j+w][0], dp[j][0]+nextDp[w][0])
					dp[j+w][1] = min(dp[j+w][1], dp[j][1]+nextDp[w][0])
					dp[j+w][1] = min(dp[j+w][1], dp[j][1]+nextDp[w][1])
				}
			}
			subSize[cur] += subSize[next]
		}
		return dp
	}
	rootDp := f(0, -1)

	for i := n; i >= 0; i-- {
		if rootDp[i][0] <= money || rootDp[i][1] <= money {
			fmt.Fprintln(out, i)
			return
		}
	}
}

// F - Components-连通块个数为i的导出子图数
// https://atcoder.jp/contests/abc287/tasks/abc287_f
// 给定一棵有n个点的树，在所有2^n−1的非空点集中，回答下列问题：
// !对于i∈[1,n]，有多少个导出子图所形成的连通块个数恰好是i。
// 数量对 998244353取模。
// 1<=n<=5000
// 二乗の木DP(二乘木dp)
// https://snuke.hatenablog.com/entry/2019/01/15/211812
// 解:
// 合并子树的时候，如果点u和子树节点都选择的时候，连通块个数会减一，其余情况都不会
// 因此还要加上是否选择点 u的状态。
// !dp[i][k][0/1] 表示当前节点为i，当前连通块个数为k，当前节点是否被选中的状态下的方案数
func abc287F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n int
	fmt.Fscan(in, &n)
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}

	components := func(n int, edges [][2]int) []int {
		tree := make([][]int, n)
		for i := 0; i < n-1; i++ {
			u, v := edges[i][0], edges[i][1]
			tree[u] = append(tree[u], v)
			tree[v] = append(tree[v], u)
		}

		subSize := make([]int, n)
		dp := make([][][2]int, n)

		var dfs func(cur, pre int)
		dfs = func(cur, pre int) {
			subSize[cur] = 1
			dp[cur] = make([][2]int, 2)
			dp[cur][0][0] = 1
			dp[cur][1][1] = 1

			for _, next := range tree[cur] {
				if next == pre {
					continue
				}

				dfs(next, cur)
				ndp := make([][2]int, subSize[cur]+subSize[next]+1) // 当前不选/当前选
				dp1, dp2 := dp[cur], dp[next]
				for i := 0; i <= subSize[cur]; i++ {
					for j := 0; j <= subSize[next]; j++ {
						ndp[i+j][0] += dp1[i][0] * (dp2[j][0] + dp2[j][1])
						ndp[i+j][0] %= MOD

						ndp[i+j][1] += dp1[i][1] * dp2[j][0]
						ndp[i+j][1] %= MOD

						if i+j-1 >= 0 {
							ndp[i+j-1][1] += dp1[i][1] * dp2[j][1]
							ndp[i+j-1][1] %= MOD
						}
					}
				}

				subSize[cur] += subSize[next]
				dp[cur] = ndp
			}
		}
		dfs(0, -1)

		res := make([]int, n+1)
		for i := 1; i <= n; i++ {
			res[i] = (dp[0][i][0] + dp[0][i][1]) % MOD
		}
		return res
	}

	res := components(n, edges)
	for i := 1; i <= n; i++ {
		fmt.Fprintln(out, res[i])
	}
}

// https://www.luogu.com.cn/remoteJudgeRedirect/atcoder/abc207_f
// [ABC207F] Tree Patrolling
// 类似二叉树监控
// 给出一棵有 n 个节点的树，每个点可能有一个警卫，每个警卫控制当前节点以及相邻节点。
// 对每个 k=0,1,2,⋯n 求出正好有 k 个节点被控制的方案数。n≤2000
//
// !dp[i][j][k] 表示以 i 为根的子树中，有 j 个节点被控制，此时 i 节点状态为 k
// k=0: 未被覆盖
// k=1: 自身有警卫
// k=2: 被某个儿子覆盖
func ABC207F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	var n int
	fmt.Fscan(in, &n)
	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	subSize := make([]int, n)
	dp := make([][][3]int, n)

	var f func(int, int)
	f = func(cur int, pre int) {
		subSize[cur] = 1
		dp[cur] = make([][3]int, 2)
		dp[cur][0][0] = 1
		dp[cur][1][1] = 1

		for _, next := range tree[cur] {
			if next == pre {
				continue
			}

			f(next, cur)
			ndp := make([][3]int, subSize[cur]+subSize[next]+2)
			dp1, dp2 := dp[cur], dp[next]
			for i := 0; i <= subSize[cur]; i++ {
				for j := 0; j <= subSize[next]; j++ {
					// 1.当前节点没被覆盖，儿子没被覆盖，转移后一定没被覆盖，被覆盖节点数不变
					ndp[i+j][0] = (ndp[i+j][0] + dp1[i][0]*dp2[j][0]) % MOD
					// 2.当前节点没被覆盖，儿子有警卫，转移后一定被覆盖，被覆盖节点数加1（自己被覆盖了）
					ndp[i+j+1][2] = (ndp[i+j+1][2] + dp1[i][0]*dp2[j][1]) % MOD
					// 3.当前节点没被覆盖，儿子被覆盖，转移后一定没被覆盖，被覆盖节点数不变（孙子有警卫对自己没有影响）
					ndp[i+j][0] = (ndp[i+j][0] + dp1[i][0]*dp2[j][2]) % MOD

					// 4.当前节点有警卫，儿子没被覆盖，转移后一定有警卫，被覆盖节点数加1（儿子被覆盖了）
					ndp[i+j+1][1] = (ndp[i+j+1][1] + dp1[i][1]*dp2[j][0]) % MOD
					// 5.当前节点有警卫，儿子有警卫，转移后一定有警卫，被覆盖节点数不变
					ndp[i+j][1] = (ndp[i+j][1] + dp1[i][1]*dp2[j][1]) % MOD
					// 6.当前节点有警卫，儿子被覆盖，转移后一定有警卫，被覆盖节点数不变
					ndp[i+j][1] = (ndp[i+j][1] + dp1[i][1]*dp2[j][2]) % MOD

					// 7.当前节点被覆盖，儿子没被覆盖，转移后一定被覆盖，被覆盖节点数不变
					ndp[i+j][2] = (ndp[i+j][2] + dp1[i][2]*dp2[j][0]) % MOD
					// 8.当前节点被覆盖，儿子有警卫，转移后一定被覆盖，被覆盖节点数不变
					ndp[i+j][2] = (ndp[i+j][2] + dp1[i][2]*dp2[j][1]) % MOD
					// 9.当前节点被覆盖，儿子被覆盖，转移后一定被覆盖，被覆盖节点数不变
					ndp[i+j][2] = (ndp[i+j][2] + dp1[i][2]*dp2[j][2]) % MOD
				}
			}

			subSize[cur] += subSize[next]
			dp[cur] = ndp
		}

	}

	f(0, -1)
	rootDp := dp[0]
	for i := 0; i <= n; i++ {
		count := (rootDp[i][0] + rootDp[i][1] + rootDp[i][2]) % MOD
		fmt.Fprintln(out, count)
	}
}

// No.196 典型DP-涂黑k个顶点的方案数(树染色)
// https://yukicoder.me/problems/no/196
// 给定一棵n个点的树,每个顶点染黑或者白
// 求涂黑k个顶点的方案数模1e9+7,且满足黑色顶点的子树也是黑色
// n<=2000 0<=k<=n
// !dp[i][j] 表示以i为根的子树中,涂黑j个顶点的方案数
func yuki196() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	var n, k int
	fmt.Fscan(in, &n, &k)

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	dp := make([][]int, n)
	subSize := make([]int, n)

	var f func(int, int)
	f = func(cur int, pre int) {
		subSize[cur] = 1
		dp[cur] = make([]int, 2)
		dp[cur][0] = 1
		dp[cur][1] = 0

		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			f(next, cur)
			ndp := make([]int, subSize[cur]+subSize[next]+1)
			for i := 0; i <= subSize[cur]; i++ {
				for j := 0; j <= subSize[next]; j++ {
					ndp[i+j] += dp[cur][i] * dp[next][j]
					ndp[i+j] %= MOD
				}
			}
			subSize[cur] += subSize[next]
			dp[cur] = ndp
		}

		dp[cur][len(dp[cur])-1] = 1 // !涂黑当前节点
	}
	f(0, -1)

	fmt.Fprintln(out, dp[0][k])
}

// P2015二叉苹果树
// https://www.luogu.com.cn/problem/P2015
// 给出一个 n 个节点的二叉树，每条边上有一个权值，
// 求最多选 m 条边（与根节点连通）的最大权值和是多少。
// dp[i][j] 表示以 i 为根的子树中保留 j 条边的最大权值和.
func P2015二叉苹果树() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)

	tree := make([][][2]int, n)
	weights := make([]map[int]int, n)
	for i := range weights {
		weights[i] = make(map[int]int)
	}
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		tree[u] = append(tree[u], [2]int{v, w})
		tree[v] = append(tree[v], [2]int{u, w})
	}

	subSize := make([]int, n)
	var f func(int, int) []int
	f = func(cur int, pre int) []int {
		subSize[cur] = 1
		dp := make([]int, k+1)
		for _, e := range tree[cur] {
			next, weight := e[0], e[1]
			if next == pre {
				continue
			}
			nextDp := f(next, cur)
			for i := subSize[cur]; i >= 1; i-- { // 当前节点必须选，所以是 i>=1
				for j := 0; j <= subSize[next]; j++ {
					if i+j > k {
						break
					}
					dp[i+j] = max(dp[i+j], dp[i-1]+nextDp[j]+weight) // !注意这里的i-1
				}
			}

			subSize[cur] += subSize[next]
		}

		return dp
	}
	rootDp := f(0, -1)
	fmt.Fprintln(out, rootDp[k])
}

// O(n*k) dp, 二乘木dp
// !滚动dp的写法(推荐)
// 树形背包/依赖背包，选取maxSelect个节点，使得选取的节点的权值和最大.
//
//	tree: 树的邻接表表示
//	root: 根节点
//	scores: 节点的分数
func TreeKnapsackDpSquare1(tree [][]int, root int, scores []int) []int {
	n := len(tree)
	dp := make([][]int, n) // dp[i][j] 表示以 i 为根的子树中选择 j 个节点的最大权值和
	subSize := make([]int, n)

	var f func(int, int)
	f = func(cur int, pre int) {
		subSize[cur] = 1
		dp[cur] = make([]int, 2)
		dp[cur][0] = 0
		dp[cur][1] = scores[cur]

		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			f(next, cur)
			ndp := make([]int, subSize[cur]+subSize[next]+1)
			// 当前节点必须选，所以是 i>=1
			for i := 1; i <= subSize[cur]; i++ {
				for j := 0; j <= subSize[next]; j++ {
					ndp[i+j] = max(ndp[i+j], dp[cur][i]+dp[next][j])
				}
			}
			subSize[cur] += subSize[next]
			dp[cur] = ndp
		}
	}

	f(root, -1)
	return dp[root]
}

// O(n*k) dp, 二乘木dp
// !原地更新dp数组的写法
// 树形背包/依赖背包，选取maxSelect个节点，使得选取的节点的权值和最大.
//
//	tree: 树的邻接表表示
//	root: 根节点
//	scores: 节点的分数
//	maxSelect: 选择的节点数，0<=maxSelect<=n
func TreeKnapsackDpSquare2(tree [][]int, root int, scores []int, maxSelect int) []int {
	n := len(tree)
	subSize := make([]int, n)

	// dp[i][j] 表示以 i 为根的子树中选择 j 个节点的最大权值和
	var f func(int, int) []int
	f = func(cur int, pre int) []int {
		subSize[cur] = 1
		dp := make([]int, maxSelect+1)
		dp[0] = 0
		dp[1] = scores[cur]

		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			nextDp := f(next, cur)
			for j := subSize[cur]; j >= 1; j-- { // 当前节点必须选，所以是 j>=1；原地更新背包
				for w := 0; w <= subSize[next]; w++ {
					if j+w > maxSelect {
						break
					}
					dp[j+w] = max(dp[j+w], dp[j]+nextDp[w])
				}
			}
			subSize[cur] += subSize[next]
		}

		return dp
	}

	dp := f(root, -1)
	return dp
}

type Node = struct{ weight, value int }

// O(n*k^2) dp
// 树形背包/依赖背包，选取maxSelect个节点，使得选取的节点的权值和最大.
// nodes[i] 表示第 i 个节点的权值和分数.
func TreeKnapsackDp(tree [][]int, root int, nodes []Node, maxWeight int) []int {
	var f func(int, int) []int
	f = func(cur int, pre int) []int {
		node := nodes[cur]
		dp := make([]int, maxWeight+1) // dp[i] 表示选择容量为 i 时的最大价值
		for i := node.weight; i <= maxWeight; i++ {
			dp[i] = node.value // 根节点必须选
		}
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			ndp := f(next, cur)
			for j := maxWeight; j >= node.weight; j-- {
				// 类似分组背包，枚举分给子树的容量 w，对应的子树的最大价值为 ndp[w]
				// w 不可超过 j-node.weight，否则无法选择根节点
				for w := 0; w <= j-node.weight; w++ {
					dp[j] = max(dp[j], dp[j-w]+ndp[w])
				}
			}
		}
		return dp
	}

	return f(root, -1)
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
