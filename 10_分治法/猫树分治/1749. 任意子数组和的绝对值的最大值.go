package main

const INF int = 1e18

// 1749. 任意子数组和的绝对值的最大值
// https://leetcode.cn/problems/maximum-absolute-sum-of-any-subarray/description/
func maxAbsoluteSum(nums []int) int {
	n := len(nums)
	res := -INF

	var f func(start, end int)
	f = func(start, end int) {
		if start >= end {
			return
		}
		if start+1 == end {
			res = max(res, abs(nums[start]))
			return
		}
		mid := (start + end) >> 1
		f(start, mid) // !不包含mid
		f(mid, end)   // !不包含mid

		// !包含mid
		curSum, preMax, preMin := 0, 0, 0
		for i := mid - 1; i >= start; i-- {
			curSum += nums[i]
			preMax = max(preMax, curSum)
			preMin = min(preMin, curSum)
		}
		curSum, sufMax, sufMin := 0, 0, 0
		for i := mid + 1; i < end; i++ {
			curSum += nums[i]
			sufMax = max(sufMax, curSum)
			sufMin = min(sufMin, curSum)
		}
		res = max(res, abs(nums[mid]+preMax+sufMax))
		res = max(res, abs(nums[mid]+preMin+sufMin))
	}

	f(0, n)
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
