// 3414. 不重叠区间的最大得分
// https://leetcode.cn/problems/maximum-score-of-non-overlapping-intervals/description/
//
// 给你一个二维整数数组 intervals，其中 intervals[i] = [li, ri, weighti]。
// 区间 i 的起点为 li，终点为 ri，权重为 weighti。你最多可以选择 4 个互不重叠 的区间。所选择区间的 得分 定义为这些区间权重的总和。
// 返回一个至多包含 4 个下标且字典序最小的数组，表示从 intervals 中选中的互不重叠且得分最大的区间。
// 如果两个区间没有任何重叠点，则称二者 互不重叠 。特别地，如果两个区间共享左边界或右边界，也认为二者重叠。

package main

import (
	"slices"
	"sort"
)

type Interval struct {
	l, r, w int
	id      int
}

func maximumWeight(intervals [][]int) []int {
	n := len(intervals)
	arr := make([]Interval, n)
	for i := 0; i < n; i++ {
		arr[i] = Interval{
			l:  intervals[i][0],
			r:  intervals[i][1],
			w:  intervals[i][2],
			id: i,
		}
	}
	slices.SortFunc(arr, func(v1, v2 Interval) int {
		if v1.r != v2.r {
			return v1.r - v2.r
		}
		return v1.l - v2.l
	})

	rs := make([]int, n)
	for i := 0; i < n; i++ {
		rs[i] = arr[i].r
	}

	pre := make([]int, n)
	for i := 0; i < n; i++ {
		left := arr[i].l
		x := sort.Search(n, func(mid int) bool {
			return rs[mid] >= left
		})
		pre[i] = x - 1
	}

	maxK := 4
	dp := make([][]uint64, maxK+1)
	bestPath := make([][][]int, maxK+1)
	for k := 0; k <= maxK; k++ {
		dp[k] = make([]uint64, n+1)
		bestPath[k] = make([][]int, n+1)
	}

	for i := 1; i <= n; i++ {
		curId := i - 1
		curW := uint64(arr[curId].w)
		pi := pre[curId]
		for k := 1; k <= maxK; k++ {
			dp[k][i] = dp[k][i-1]
			bestPath[k][i] = make([]int, len(bestPath[k][i-1]))
			copy(bestPath[k][i], bestPath[k][i-1])

			var candRes uint64
			if pi < 0 {
				candRes = curW
			} else {
				candRes = dp[k-1][pi+1] + curW
			}

			if candRes > dp[k][i] {
				dp[k][i] = candRes
				if pi < 0 {
					bestPath[k][i] = []int{arr[curId].id}
				} else {
					tmp := make([]int, len(bestPath[k-1][pi+1]))
					copy(tmp, bestPath[k-1][pi+1])
					tmp = append(tmp, arr[curId].id)
					slices.Sort(tmp)
					bestPath[k][i] = tmp
				}
			} else if candRes == dp[k][i] {
				var candidate []int
				if pi < 0 {
					candidate = []int{arr[curId].id}
				} else {
					tmp := make([]int, len(bestPath[k-1][pi+1]))
					copy(tmp, bestPath[k-1][pi+1])
					tmp = append(tmp, arr[curId].id)
					slices.Sort(tmp)
					candidate = tmp
				}

				if slices.Compare(candidate, bestPath[k][i]) < 0 {
					bestPath[k][i] = candidate
				}
			}
		}
	}

	var resPath []int
	var res uint64
	for k := 1; k <= maxK; k++ {
		val := dp[k][n]
		if val > res {
			res = val
			resPath = bestPath[k][n]
		} else if val == res {
			if slices.Compare(bestPath[k][n], resPath) < 0 {
				resPath = bestPath[k][n]
			}
		}
	}

	return resPath
}
