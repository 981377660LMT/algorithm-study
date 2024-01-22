// 给定一棵n个顶点的树
// 求包含根节点0、顶点数为k的树的个数
// !结点编号为0~n-1
// https://snuke.hatenablog.com/entry/2019/01/15/211812
// TODO, Not Verified

package main

import "fmt"

const MOD int = 998244353

func main() {
	tree := make([][]int, 5)
	for i := 0; i < 4; i++ {
		tree[i] = append(tree[i], i+1)
		tree[i+1] = append(tree[i+1], i)
	}
	fmt.Println(Solve1(tree, 3))
	fmt.Println(Solve2(tree, 3))
}

// 1<=k<=n<=3000
// !dp[i][v]表示以i为根的子树中包含v个顶点的树的个数
// O(n^2)
// 咋一看这样的复杂度会是O(n^3)，
// 但如果每次枚举的范围都是儿子树的大小，可以证明这样的树型 dp的复杂度是 O(n^2)的。 (完全图)
func Solve1(tree [][]int, k int) int {
	n := len(tree)
	subSize := make([]int, n)
	dp := make([][]int, n)

	var f func(int, int)
	f = func(cur, pre int) {
		subSize[cur] = 1
		dp[cur] = make([]int, 2)
		dp[cur][0] = 0
		dp[cur][1] = 1
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			f(next, cur)
			ndp := make([]int, subSize[cur]+subSize[next]+1)
			dp1, dp2 := dp[cur], dp[next]
			for i := 0; i <= subSize[cur]; i++ {
				for j := 0; j <= subSize[next]; j++ {
					ndp[i+j] = (ndp[i+j] + dp1[i]*dp2[j]) % MOD
				}
			}
			subSize[cur] += subSize[next]
			dp[cur] = ndp
		}

		dp[cur][0] = 1 // !0个顶点的树的个数
	}

	f(0, -1)
	return dp[0][k]
}

// 1<=n<=1e5, 1<=k<=n
// 原理:
// n个集合合并,一开始每个集合只有一个元素
// !如果合并两个集合的代价为min(|A|, K) * min(|B|, K)，
// !则合并成一个集合的最大代价为O(n*K)
func Solve2(tree [][]int, k int) int {
	n := len(tree)
	subSize := make([]int, n)
	dp := make([][]int, n)
	res := 0

	var f func(int, int)
	f = func(cur, pre int) {
		subSize[cur] = 1
		dp[cur] = make([]int, 2)
		dp[cur][0] = 0
		dp[cur][1] = 1

		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			f(next, cur)
			ndp := make([]int, subSize[cur]+subSize[next]+1)
			dp1, dp2 := dp[cur], dp[next]
			for i := 0; i <= subSize[cur]; i++ {
				for j := 0; j <= subSize[next]; j++ {
					ndp[i+j] = (ndp[i+j] + dp1[i]*dp2[j]) % MOD
				}
			}
			subSize[cur] += subSize[next]
			dp[cur] = ndp
			if subSize[cur] > k {
				subSize[cur] = k
				dp[cur] = dp[cur][:k+1]
			}
		}

		if subSize[cur] >= k {
			res += dp[cur][k]
			res %= MOD
		}

		dp[cur][0] = 1 // !0个顶点的树的个数
	}

	f(0, -1)
	return res
}
