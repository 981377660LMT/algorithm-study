// 区间dp模板
// !dp[i][j]=min{dp[i][k]+dp[k+1][j]+f(i,j,k)} (0<=i<=k<j<n)
// TODO282. 石子合并

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

// !f: 区间[i,j]的代价函数 (0<=i<k<j<n)
//  注意有时候需要特殊处理包含两个点的的区间 (j-i=1)
func rangeDpMin(n int, f func(i, k, j int) int) int {
	dp := make([][]int, n)
	memo := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
		memo[i] = make([]int, n)
		for j := 0; j < n; j++ {
			memo[i][j] = -1
		}
	}

	var dfs func(left, right int) int
	dfs = func(left, right int) int {
		if left >= right {
			return 0
		}
		if right-left == 1 {
			return f(left, left, right)
		}
		if memo[left][right] != -1 {
			return memo[left][right]
		}

		res := INF
		// !左右子区间`不能`共用中间点 k
		for k := left; k < right; k++ {
			cand := dfs(left, k) + dfs(k+1, right) + f(left, k, right)
			if cand < res {
				res = cand
			}
		}

		// !左右子区间`能`共用中间点 k
		// for k := left + 1; k < right; k++ {
		// 	cand := dfs(left, k) + dfs(k, right) + f(left, k, right)
		// 	if cand < res {
		// 		res = cand
		// 	}
		// }

		memo[left][right] = res
		return res
	}

	return dfs(0, n-1)
}

// 1039. 多边形三角剖分的最低得分
func minScoreTriangulation(values []int) int {
	n := len(values)
	f := func(i, k, j int) int {
		if j-i <= 1 {
			return 0
		}
		return values[i] * values[k] * values[j]
	}
	return rangeDpMin(n, f)
}

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	res := rangeDpMin(n, func(i, k, j int) int {
		return preSum[j+1] - preSum[i]
	})
	fmt.Fprintln(out, res)
}
