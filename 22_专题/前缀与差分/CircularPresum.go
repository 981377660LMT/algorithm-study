package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	ABC331D()
}

const INF int = 1e18

// 环形前缀和.
func CircularPresum(nums []int) func(start int, end int) int {
	n := len(nums)
	preSum := make([]int, n+1)
	for i, v := range nums {
		preSum[i+1] = preSum[i] + v
	}

	cal := func(r int) int {
		return preSum[n]*(r/n) + preSum[r%n]
	}

	query := func(start int, end int) int {
		if start >= end {
			return 0
		}
		return cal(end) - cal(start)
	}

	return query
}

// 二维环形前缀和.
// 0 <= row1 < row2(不包含) <= n
// 0 <= col1 < col2(不包含) <= m
func CircularPreSum2D(grid [][]int) func(row1 int, col1 int, row2 int, col2 int) int {
	n := len(grid)
	m := len(grid[0])
	preSum := make([][]int, n+1)
	for i := range preSum {
		preSum[i] = make([]int, m+1)
	}
	for i, row := range grid {
		tmp1, tmp2 := preSum[i], preSum[i+1]
		for j := range row {
			tmp2[j+1] = tmp2[j] + tmp1[j+1] - tmp1[j] + row[j]
		}
	}

	cal := func(r int, c int) int {
		res1 := preSum[n][m] * (r / n) * (c / m)
		res2 := preSum[r%n][m] * (c / m)
		res3 := preSum[n][c%m] * (r / n)
		res4 := preSum[r%n][c%m]
		return res1 + res2 + res3 + res4
	}

	query := func(row1 int, col1 int, row2 int, col2 int) int {
		if row1 >= row2 || col1 >= col2 {
			return 0
		}
		res1 := cal(row2, col2)
		res2 := cal(row1, col2)
		res3 := cal(row2, col1)
		res4 := cal(row1, col1)
		return res1 - res2 - res3 + res4
	}

	return query
}

// 100076. 无限数组的最短子数组
// https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/
// 求循环数组中和为 target 的最短子数组的长度.不存在则返回 -1.
// 1 <= nums.length <= 1e5
// 1 <= nums[i] <= 1e5
// 1 <= target <= 1e9
func minSizeSubarray(nums []int, target int) int {
	Q := CircularPresum(nums)
	res := INF
	for start := range nums {
		cand := sort.Search(1e9+10, func(mid int) bool {
			return Q(start, start+mid) >= target
		})
		if Q(start, start+cand) == target {
			res = min(res, cand)
		}
	}

	if res == INF {
		return -1
	}
	return res
}

// https://atcoder.jp/contests/abc331/tasks/abc331_d
func ABC331D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		for j, c := range s {
			if c == 'B' {
				grid[i][j] = 1
			}
		}
	}

	query := CircularPreSum2D(grid)
	for i := 0; i < q; i++ {
		var row1, col1, row2, col2 int
		fmt.Fscan(in, &row1, &col1, &row2, &col2)
		fmt.Fprintln(out, query(row1, col1, row2+1, col2+1))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
