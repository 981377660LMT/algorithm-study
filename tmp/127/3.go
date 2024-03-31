package main

type interval = struct{ leftStart, leftEnd, value int }

// 将 nums 的所有非空子数组的元素进行 op 操作，返回所有不同的结果和其出现次数.
//
// nums: 1 <= nums.length <= 1e5.
// op: 与/或/gcd/lcm 中的一种操作，具有单调性.
// f:
// 数组的右端点为right.
// interval 的 leftStart/leftEnd 表示子数组的左端点left的范围.
// interval 的 value 表示该子数组 arr[left,right] 的 op 结果.
func LogTrick(nums []int, op func(int, int) int, f func(left []interval, right int)) map[int]int {
	res := make(map[int]int)

	dp := []interval{}
	for pos, cur := range nums {
		for i, pre := range dp {
			dp[i].value = op(pre.value, cur)
		}
		dp = append(dp, interval{leftStart: pos, leftEnd: pos + 1, value: cur})

		// 去重
		ptr := 0
		for _, v := range dp[1:] {
			if v.value != dp[ptr].value {
				ptr++
				dp[ptr] = v
			} else {
				dp[ptr].leftEnd = v.leftEnd
			}
		}
		dp = dp[:ptr+1]

		// 将区间[0,pos]分成了dp.length个左闭右开区间.
		// 每一段区间的左端点left范围 在 [dp[i].leftStart,dp[i].leftEnd) 中。
		// 对应子数组 arr[left:pos+1] 的 op 值为 dp[i].value.
		for _, v := range dp {
			res[v.value] += v.leftEnd - v.leftStart
		}
		if f != nil {
			f(dp, pos)
		}
	}

	return res
}

const INF int = 4e18

// 给你一个 非负 整数数组 nums 和一个整数 k 。

// 如果一个数组中所有元素的按位或运算 OR 的值 至少 为 k ，那么我们称这个数组是 特别的 。

// 请你返回 nums 中 最短特别非空 子数组的长度，如果特别子数组不存在，那么返回 -1 。
func minimumSubarrayLength(nums []int, k int) int {
	res := INF
	LogTrick(nums, func(a, b int) int { return a | b }, func(left []interval, right int) {
		for _, v := range left {
			if v.value >= k {
				res = min(res, right-v.leftEnd+2)
			}
		}
	})
	if res == INF {
		return -1
	}
	return res
}
