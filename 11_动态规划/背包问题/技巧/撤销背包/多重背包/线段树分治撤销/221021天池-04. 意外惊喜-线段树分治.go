// 阿里云天池专场
// 有n个单调不减数组，求从这n个数组中选择k个数的最大和，（只能拿数组的一个前缀）
// 221021天池-04. 意外惊喜 https://leetcode.cn/contest/tianchi2022/problems/tRZfIV/
// 比 2218. 从栈中取出K个硬币的最大面值和 https://leetcode.cn/problems/maximum-value-of-k-coins-from-piles/
// !多了每个礼物包中的礼物价值 非严格递增 的条件，去除了背包总数的限制
// n<=3000 k<=3000 O(nklogn)
//
// 解：
// !因为数组是单调的，所以能得到一个结论就是我们最多只会拿不全一个数组
// 因为当有两个数组都拿一部分时，那么肯定可以通过减少某个数组拿的数字去拿另一个数组来增大这个和
// 换句话说就是，只要选了一个序列的第一个数，那么最优方案一定是尽量把这个序列的所有数都选上。
// 即：给定 n 个物品，对 i=1,2,⋯,n，你需要求出：去掉第 i 个物品后，对其他物品做背包的结果(可撤销01背包)。
// !所以我们就可以通过枚举哪个数组不全选，采用分治删点(线段树分治/可撤销背包)的方式dp

package main

import (
	"bufio"
	"fmt"
	"os"
)

func brilliantSurprise(present [][]int, limit int) int {
	n := len(present)
	groupSum := make([]int, n)
	for i, group := range present {
		for _, v := range group {
			groupSum[i] += v
		}
	}

	res := 0
	initState := make([]int, limit+1)
	MutateWithoutOne(
		&initState,
		0, n,
		func(state *S) *S {
			newState := append([]int{}, *state...)
			return &newState
		},
		func(state *S, index int) {
			m := len(present[index])
			dp := *state
			for i := limit; i >= m; i-- {
				dp[i] = max(dp[i], dp[i-m]+groupSum[index])
			}
		},
		func(state *S, index int) {
			dp := *state
			res = max(res, dp[limit])
			group, curSum := present[index], 0
			for i := 0; i < min(len(group), limit); i++ {
				curSum += group[i]
				res = max(res, dp[limit-i-1]+curSum)
			}
		},
	)

	return res
}

type S = []int

// 线段树分治的特殊情形.
func MutateWithoutOne(
	initState *S,
	start, end int,
	copy func(state *S) *S,
	mutate func(state *S, index int),
	query func(state *S, index int),
) {
	var dfs func(state *S, curStart, curEnd int)
	dfs = func(state *S, curStart, curEnd int) {
		if curEnd == curStart+1 {
			query(state, curStart)
			return
		}

		mid := (curStart + curEnd) >> 1
		leftCopy := copy(state)
		for i := curStart; i < mid; i++ {
			mutate(leftCopy, i)
		}
		dfs(leftCopy, mid, curEnd)

		rightCopy := copy(state)
		for i := mid; i < curEnd; i++ {
			mutate(rightCopy, i)
		}
		dfs(rightCopy, curStart, mid)
	}

	dfs(initState, start, end)
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

// https://www.luogu.com.cn/problem/CF1442D
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, limit int
	fmt.Fscan(in, &n, &limit)
	present := make([][]int, n)
	for i := 0; i < n; i++ {
		var m int
		fmt.Fscan(in, &m)
		group := make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &group[j])
		}
		present[i] = group
	}

	res := brilliantSurprise(present, limit)
	fmt.Fprintln(out, res)
}
