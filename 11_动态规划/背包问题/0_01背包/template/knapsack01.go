// https://nyaannyaan.github.io/library/verify/verify-aoj-dpl/aoj-dpl-1-f.test.cpp
// 1. dp[i][j] 表示前 i 个物品，取得重量为 j 时的最大价值
// 2. dp[i][j] 表示前 i 个物品，达成总价值为 j 时的最小重量
// 3. meet in the middle
// 4. knapsack01-branch-and-bound

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, limit int
	fmt.Fscan(in, &n, &limit)
	values := make([]int, n)
	weights := make([]int, n)
	for i := range values {
		fmt.Fscan(in, &values[i], &weights[i])
	}
	fmt.Fprintln(out, knapsack01(values, weights, limit))
}

const INF int = 1e18

func knapsack01(values []int, weights []int, limit int) int {
	vSum := 0
	for _, v := range values {
		vSum += v
	}
	cost1 := bits.Len(uint(limit) * uint((len(values) + 1)))
	cost2 := bits.Len(uint(vSum) * uint((len(values) + 1)))
	cost3 := (len(values)+1)/2 + bits.Len(uint(len(values)+1))
	minCost := min(cost1, min(cost2, cost3))
	if minCost > 32 {
		return solver4(values, weights, limit)
	}
	if minCost == cost1 {
		return solver1(values, weights, limit)
	}
	if minCost == cost2 {
		return solver2(values, weights, limit)
	}
	return solver3(values, weights, limit)
}

// limit*n 很小的情况下，dp[i][j] 表示前 i 个物品，取得重量为 j 时的最大价值
func solver1(values []int, weights []int, limit int) int {
	dp := make([]int, limit+1)
	for i := 0; i < len(values); i++ {
		v, w := values[i], weights[i]
		for j := limit - w; j >= 0; j-- {
			dp[j+w] = max(dp[j+w], dp[j]+v)
		}
	}
	res := 0
	for _, v := range dp {
		res = max(res, v)
	}
	return res
}

// ∑vi*n 很小的情况下，dp[i][j] 表示前 i 个物品，达成总价值为 j 时的最小重量
// 把重量看成价值，价值看成重量，求同等价值下能得到的最小重量，若该最小重量不超过背包容量，则该价值合法。所有合法价值的最大值即为答案
func solver2(values []int, weights []int, limit int) int {
	vSum := 0
	for _, v := range values {
		vSum += v
	}
	dp := make([]int, vSum+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for i := 0; i < len(values); i++ {
		v, w := values[i], weights[i]
		for j := vSum - v; j >= 0; j-- {
			dp[j+v] = min(dp[j+v], dp[j]+w)
		}
	}
	res := 0
	for j := vSum; j >= 0; j-- {
		if dp[j] <= limit {
			res = j
			break
		}
	}
	return res
}

// n<=40 的情况下，可以 meet in the middle
//
//	dp处理出子集的(价值，重量)对，排序后双指针求解
func solver3(values []int, weights []int, limit int) int {
	getDp := func(vs, ws []int) [][2]int {
		dp := make([][2]int, 0, 1<<uint(len(vs)))
		dp = append(dp, [2]int{0, 0})
		for i := range vs {
			n := len(dp)
			for j := 0; j < n; j++ {
				if dp[j][1]+ws[i] <= limit {
					dp = append(dp, [2]int{dp[j][0] + vs[i], dp[j][1] + ws[i]})
				}
			}
		}

		sort.Slice(dp, func(i, j int) bool {
			return dp[i][1] < dp[j][1]
		})
		for i := 1; i < len(dp); i++ {
			dp[i][0] = max(dp[i][0], dp[i-1][0])
		}
		return dp
	}

	n := len(values)
	leftValues, leftWeights := values[:n/2], weights[:n/2]
	rightValues, rightWeights := values[n/2:], weights[n/2:]
	dp1, dp2 := getDp(leftValues, leftWeights), getDp(rightValues, rightWeights)
	res := 0
	right := len(dp2) - 1
	for left := 0; left < len(dp1); left++ {
		for right >= 0 && (dp1[left][1]+dp2[right][1] > limit) {
			right--
		}
		if right == -1 {
			break
		}
		res = max(res, dp1[left][0]+dp2[right][0])
	}

	return res
}

// knapsack01-branch-and-bound
// n<=2000,vi<=1e9,wi<=1e9,limit<=1e9
func solver4(values []int, weights []int, limit int) int {
	n := len(values)
	goods := make([][]int, n)
	for i := range goods {
		goods[i] = []int{values[i], weights[i]}
	}
	sort.Slice(goods, func(i, j int) bool {
		return goods[i][0]*goods[j][1] > goods[j][0]*goods[i][1]
	})

	best := 0
	relax := func(i, v, w int) (int, bool) {
		res := v
		flag := true
		for i < n {
			if w == 0 {
				break
			}
			if w >= goods[i][1] {
				w -= goods[i][1]
				res += goods[i][0]
				i++
				continue
			}
			flag = false
			res += goods[i][0] * w / goods[i][1]
			break
		}
		return res, flag
	}

	var dfs func(i int, v int, w int)
	dfs = func(i int, v int, w int) {
		if i == n {
			if v > best {
				best = v
			}
			return
		}

		rel, flag := relax(i, v, w)
		if flag {
			if rel > best {
				best = rel
			}
			return
		}

		if rel < best {
			return
		}

		if w >= goods[i][1] {
			dfs(i+1, v+goods[i][0], w-goods[i][1])
		}

		dfs(i+1, v, w)
	}

	dfs(0, 0, limit)
	return best
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
