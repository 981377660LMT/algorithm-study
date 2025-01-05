package main

import (
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
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].r == arr[j].r {
			return arr[i].l < arr[j].l
		}
		return arr[i].r < arr[j].r
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

	lexLess := func(a, b []int) bool {
		la, lb := len(a), len(b)
		for i := 0; i < la && i < lb; i++ {
			if a[i] < b[i] {
				return true
			} else if a[i] > b[i] {
				return false
			}
		}
		return la < lb
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
					sort.Ints(tmp)
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
					sort.Ints(tmp)
					candidate = tmp
				}
				if lexLess(candidate, bestPath[k][i]) {
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
			if lexLess(bestPath[k][n], resPath) {
				resPath = bestPath[k][n]
			}
		}
	}

	return resPath
}
