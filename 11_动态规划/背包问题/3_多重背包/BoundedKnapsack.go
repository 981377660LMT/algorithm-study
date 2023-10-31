//
// Bounded Knapsack Problem （多重背包问题）
//
// Description:
//
//   There are n kinds of items of profit pi, weight bi, and
//   amount mi (i in {0...n-1}). For a given W, we want to select items
//   to maximize the total profit subject to the total weight
//   is at most W and the number of each item is at most mi.
//
//   This problem can be solved by the following DP.
//     E[j][w] = max {E[j-1][w], E[j-1][w-bj]+pj, E[j-1][w-2bj]+2pj. ...}
//   A naive implementation requires O(nmW) time. However, we can reduce
//   this to O(nW) as follows.
//
//   We compute E[j][s], E[j][s+bj]. E[s+2bj] ... for each s in [0,bj).
//   For simplicity, we consider s = 0. Then, we have
//     E[j][w+bj]  = max {E[j-1][w+bj], E[j-1][w]+pj, E[j-1][w-bj]+2pj. ...}
//   By comparing this with the original formula,
//   - E[j][w+bj] contains E[j-1][w+bj] term
//   - E[j][w+bj] does not contain E[j-1][w-mjbj] term
//   - The all terms have been added pj
//   Thus, by using a data structure that supports these operations,
//   we can perform the DP efficiently. The data structure is implemented
//   by a maximum queue with one accumulation parameter.
//
// Complexity:
//
//   O(n W)
//

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Luogu1776()
	Luogu1077()
}

// const MOD int = 1e9 + 7
const MOD int = 1e6 + 7

func Luogu1776() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, maxCapacity int
	fmt.Fscan(in, &n, &maxCapacity)
	values := make([]int, n)
	weights := make([]int, n)
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i], &weights[i], &counts[i])
	}

	res := BoundedKnapsackDP(values, weights, counts, maxCapacity)
	fmt.Fprintln(out, res)
}

// P1077 [NOIP2012 普及组] 摆花
// https://www.luogu.com.cn/problem/P1077
func Luogu1077() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, target int
	fmt.Fscan(in, &n, &target)
	counts := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &counts[i])
	}

	values := make([]int, n)
	for i := 0; i < n; i++ {
		values[i] = 1
	}

	dp := BoundedKnapsackDPCountWays(values, counts)
	fmt.Fprintln(out, dp[target])

}

func waysToReachTarget(target int, types [][]int) int {
	counts := make([]int, len(types))
	for i := range types {
		counts[i] = types[i][0]
	}
	values := make([]int, len(types))
	for i := range types {
		values[i] = types[i][1]
	}

	dp := BoundedKnapsackDPCountWays(values, counts)
	return dp[target]
}

// 多重背包问题
// 有n种物品，每种物品得分为values[i]，重量为weights[i]，数量为counts[i]，背包容量为W.
// 选择物品使得得分最大，总重量不超过maxW.求出最大得分。
//
// dp[i][w]表示前i种物品，总重量为w时的最大得分.
// dp[i][w] = max{dp[i-1][w-k*weights[i]]+k*values[i] | 0<=k<=counts[i], w-k*weights[i]>=0}
// 分组转移，单调队列维护滑动窗口最大值
// !时间复杂度 O(nW)
func BoundedKnapsackDP(values []int, weights []int, counts []int, maxCapacity int) int {
	type item struct{ max, j int }
	dp := make([]int, maxCapacity+1)
	for i, count := range counts {
		v, w := values[i], weights[i]
		for rem := 0; rem < w; rem++ { // 按照 j%w 的结果，`分组转移`，rem 表示 remainder
			queue := []item{}
			for j := 0; j*w+rem <= maxCapacity; j++ {
				cand := dp[j*w+rem] - j*v
				for len(queue) > 0 && queue[len(queue)-1].max <= cand {
					queue = queue[:len(queue)-1] // 及时去掉无用数据
				}
				queue = append(queue, item{cand, j})
				// 本质是查表法，q[0].max 就表示 f[(j-1)*w+r]-(j-1)*v, f[(j-2)*w+r]-(j-2)*v, …… 这些转移来源的最大值
				dp[j*w+rem] = queue[0].max + j*v // 物品个数为两个 j 的差（前缀和思想）
				if j-queue[0].j == count {       // 至多选 num 个物品
					queue = queue[1:] // 及时去掉无用数据
				}
			}
		}
	}
	return dp[maxCapacity]
}

// 多重背包二进制优化.
func BoundedKnapsackDPBinary(values []int, weights []int, counts []int, maxCapacity int) int {
	dp := make([]int, maxCapacity+1)
	for i, num := range counts {
		v, w := values[i], weights[i]
		for k1 := 1; num > 0; k1 <<= 1 {
			k := min(k1, num)
			for j := maxCapacity; j >= k*w; j-- {
				dp[j] = max(dp[j], dp[j-k*w]+k*v)
			}
			num -= k
		}
	}
	return dp[maxCapacity]
}

func BoundedKnapsackDPNaive(values []int, weights []int, counts []int, maxCapacity int) int {
	n := len(counts)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, maxCapacity+1)
	}
	for i, count := range counts {
		v, w := values[i], weights[i]
		for j := range dp[i] {
			// 枚举选了 k=0,1,2,...num 个第 i 种物品
			for k := 0; k <= count && k*w <= j; k++ {
				dp[i+1][j] = max(dp[i+1][j], dp[i][j-k*w]+k*v)
			}
		}
	}
	return dp[n][maxCapacity]
}

// 多重背包求方案数(分组前缀和优化).
//
//		values: 物品价值
//		counts: 物品数量
//		返回值: dp[i] 表示总价值为 i 的方案数 % MOD.
//	 时间复杂度 O(n*sum(values[i]*counts[i]))
func BoundedKnapsackDPCountWays(values []int, counts []int) []int {
	allSum := 0
	count0 := 0
	for i, count := range counts {
		v := values[i]
		if v == 0 {
			count0 += count
			continue
		}
		allSum += count * values[i]
	}
	dp := make([]int, allSum+1)
	dp[0] = count0 + 1

	maxJ := 0
	for i, count := range counts {
		v := values[i]
		if v == 0 {
			continue
		}
		maxJ += v * count
		for j := v; j <= maxJ; j++ {
			dp[j] = (dp[j] + dp[j-v]) % MOD // 同余前缀和
		}
		for j := maxJ; j >= v*(count+1); j-- { // 超过个数
			dp[j] = (dp[j] - dp[j-v*(count+1)]) % MOD
		}
	}

	for i := range dp {
		if dp[i] < 0 {
			dp[i] += MOD
		}
	}
	return dp
}

// O(n*upper).
func BoundedKnapsackDpCountWaysWithUpper(values []int, counts []int, upper int) []int {
	count0 := 0
	for i, count := range counts {
		v := values[i]
		if v == 0 {
			count0 += count
			continue
		}
	}
	dp := make([]int, upper+1)
	dp[0] = count0 + 1

	maxJ := 0
	for i, count := range counts {
		v := values[i]
		if v == 0 {
			continue
		}
		maxJ += v * count
		if maxJ > upper {
			maxJ = upper
		}
		for j := v; j <= maxJ; j++ {
			dp[j] = (dp[j] + dp[j-v]) % MOD // 同余前缀和
		}
		for j := maxJ; j >= v*(count+1); j-- {
			dp[j] = (dp[j] - dp[j-v*(count+1)]) % MOD
		}
	}

	for i := range dp {
		if dp[i] < 0 {
			dp[i] += MOD
		}
	}

	return dp
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
