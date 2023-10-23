package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5}
	LogTrick(nums, func(a, b int) int { return a | b }, func(intervals []interval, right int) {
		fmt.Println(right, intervals)
	})
}

type interval = struct{ leftStart, leftEnd, value int }

// 将 nums 的所有非空子数组的元素进行 op 操作，返回所有不同的结果和其出现次数.
//
// nums: 1 <= nums.length <= 1e5.
// op: 与/或/gcd/lcm 中的一种操作，具有单调性.
// f:
// 数组的右端点为right.
// interval 的 leftStart/leftEnd 表示子数组的左端点left的范围.
// interval 的 value 表示该子数组 arr[left,right] 的 op 结果.
func LogTrick(nums []int, op func(int, int) int, f func(intervals []interval, right int)) map[int]int {
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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
