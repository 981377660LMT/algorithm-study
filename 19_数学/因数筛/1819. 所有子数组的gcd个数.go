// 统计数组所有子数组的 gcd 的不同个数，复杂度 O(n*log^2max)
// https://github.com/981377660LMT/codeforces-go/blob/50efc27004a0864fad32b2070b0f01c87b67b7c1/copypasta/bits.go#L522
// logTrick/bitOpTrick

package main

import "fmt"

func main() {
	fmt.Println(CountGcdOfAllSubarray([]int{6, 10, 15}))
}

// 统计数组所有子数组的gcd的不同个数，复杂度 O(n*log^2max)
func CountGcdOfAllSubarray(nums []int) int {
	return len(LogTrick(nums, gcd, nil))
}

// 将 nums 的所有非空子数组的元素进行 op 操作，返回所有不同的结果和其出现次数.
//  nums: 1 <= nums.length <= 1e5.
//  op: 与/或/gcd/lcm 中的一种操作，具有单调性.
//  f: nums[:end] 中所有子数组的结果为 preCounter.
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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
