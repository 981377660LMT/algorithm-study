package main

// 给你一个数组 nums 和一个整数 k 。你需要找到 nums 的一个 子数组 ，满足子数组中所有元素按位与运算 AND 的值与 k 的 绝对差 尽可能 小 。
// 换言之，你需要选择一个子数组 nums[l..r] 满足 |k - (nums[l] AND nums[l + 1] ... AND nums[r])| 最小。

// 请你返回 最小 的绝对差值。

// 子数组是数组中连续的 非空 元素序列。

const INF int = 1e18

func minimumDifference(nums []int, k int) int {
	andRes := LogTrick(nums, func(a, b int) int { return a & b }, nil)
	res := INF
	for and := range andRes {
		diff := abs(and - k)
		if diff < res {
			res = diff
		}
	}
	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 将 nums 的所有非空子数组的元素进行 op 操作，返回所有不同的结果和其出现次数.
//
//	nums: 1 <= nums.length <= 1e5.
//	op: 与/或/gcd/lcm 中的一种操作，具有单调性.
//	f: nums[:end] 中所有子数组的结果为 preCounter.
func LogTrick(nums []int, op func(int, int) int, f func(end int, preCounter map[int]int)) map[int]int {
	res := make(map[int]int)
	dp := []int{}
	for pos, cur := range nums {
		for i, pre := range dp {
			dp[i] = op(pre, cur)
		}
		dp = append(dp, cur)

		// 去重
		ptr := 0
		for _, v := range dp[1:] {
			if v != dp[ptr] {
				ptr++
				dp[ptr] = v
			}
		}

		dp = dp[:ptr+1]
		for _, v := range dp {
			res[v]++
		}
		if f != nil {
			f(pos+1, res)
		}
	}

	return res
}
