//
// Bounded Knapsack Problem
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
	"fmt"
	"math/rand"
)

func main() {
	rand.Seed(20)
	n := 100
	W := 10000
	ps := make([]int, n)
	ws := make([]int, n)
	ms := make([]int, n)
	for i := 0; i < n; i++ {
		ps[i] = rand.Intn(n) + 1
		ws[i] = rand.Intn(n) + 1
		ms[i] = rand.Intn(n) + 1
	}

	res1 := BoundedKnapsackDP(ps, ws, ms, W)
	res2 := BoundedKnapsackDPNaive(ps, ws, ms, W)
	fmt.Println(res1, res2)
}

// 多重背包问题
// 有n种物品，每种物品得分为scores[i]，重量为weights[i]，数量为counts[i]，背包容量为W.
// 选择物品使得得分最大，总重量不超过W.求出最大得分。
//
// dp[i][w]表示前i种物品，总重量为w时的最大得分.
// dp[i][w] = max{dp[i-1][w-k*weights[i]]+k*scores[i] | 0<=k<=counts[i], w-k*weights[i]>=0}
// 滑窗+分组前缀和优化
// !O(nW)时间复杂度
func BoundedKnapsackDP(scores []int, weights []int, counts []int, W int) int {
	n := len(scores)
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, W+1)
	}

	for i := 0; i < n; i++ {
		for s := 0; s < weights[i]; s++ {
			delta := 0
			queue := make([]int, 0)
			top := make([]int, 0)
			for w := s; w <= W; w += weights[i] {
				delta += scores[i]
				a := dp[i][w] - delta
				queue = append(queue, a)
				for len(top) > 0 && top[len(top)-1] < a {
					top = top[:len(top)-1]
				}
				top = append(top, a)
				for len(queue) > counts[i]+1 {
					if queue[0] == top[0] {
						top = top[1:]
					}
					queue = queue[1:]
				}
				dp[i+1][w] = top[0] + delta
			}
		}
	}

	res := 0
	for w := 0; w <= W; w++ {
		res = max(res, dp[n][w])
	}
	return res
}

func BoundedKnapsackDPNaive(scores []int, weights []int, counts []int, W int) int {
	n := len(scores)
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, W+1)
	}

	for i := 0; i < n; i++ {
		for w := 0; w <= W; w++ {
			dp[i+1][w] = dp[i][w]
			for j := 1; j <= counts[i]; j++ {
				if w-j*weights[i] < 0 {
					break
				}
				if dp[i+1][w] < dp[i][w-j*weights[i]]+j*scores[i] {
					dp[i+1][w] = dp[i][w-j*weights[i]] + j*scores[i]
				}
			}
		}
	}

	res := 0
	for w := 0; w <= W; w++ {
		res = max(res, dp[n][w])
	}
	return res
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
