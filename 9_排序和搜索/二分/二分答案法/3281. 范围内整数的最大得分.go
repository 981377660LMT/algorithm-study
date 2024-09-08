// 3281. 范围内整数的最大得分
// https://leetcode.cn/problems/maximize-score-of-numbers-in-ranges/description/
// 给你一个整数数组 start 和一个整数 d，代表 n 个区间 [start[i], start[i] + d]。
// 你需要选择 n 个整数，其中第 i 个整数必须属于第 i 个区间。所选整数的 得分 定义为所选整数两两之间的 最小 绝对差。
// 返回所选整数的 最大可能得分 。

package main

import (
	"slices"
)

func maxPossibleScore(start []int, d int) int {
	n := len(start)
	slices.Sort(start)
	check := func(end int) bool {
		end--
		cur := start[0]
		for i := 1; i < n; i++ {
			left, right := start[i], start[i]+d
			if cur+end > right {
				return false
			}
			cur = max(cur+end, left)
		}
		return true
	}
	right := MaxRight(0, check, (start[n-1]+d-start[0])/(n-1)+1)
	return right - 1
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
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
