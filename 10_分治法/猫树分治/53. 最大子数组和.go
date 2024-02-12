// https://leetcode.cn/problems/maximum-subarray/
// 进阶：如果你已经实现复杂度为 O(n) 的解法，尝试使用更为精妙的 分治法 求解。

// f(start,end) = max(f(start,mid),f(mid,end),w(start,end))
// !其中w(start,end)表示包含mid的最大子数组和，前后缀求解即可.类似猫树分治?

package main

const INF int = 1e18

func maxSubArray(nums []int) int {
	n := len(nums)
	res := -INF

	var f func(start, end int)
	f = func(start, end int) {
		if start >= end {
			return
		}
		if start+1 == end {
			res = max(res, nums[start])
			return
		}
		mid := (start + end) >> 1
		f(start, mid) // !不包含mid
		f(mid, end)   // !不包含mid

		// !包含mid
		curSum, preMax := 0, 0
		for i := mid - 1; i >= start; i-- {
			curSum += nums[i]
			preMax = max(preMax, curSum)
		}
		curSum, sufMax := 0, 0
		for i := mid + 1; i < end; i++ {
			curSum += nums[i]
			sufMax = max(sufMax, curSum)
		}
		res = max(res, nums[mid]+preMax+sufMax)
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
