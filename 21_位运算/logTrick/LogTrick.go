package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	// var q int
	// fmt.Fscan(in, &q)
	// queries := make([]int, q)
	// for i := range queries {
	// 	fmt.Fscan(in, &queries[i])
	// }
	// res := CGCDSSQ(nums, queries)
	// for _, v := range res {
	// 	fmt.Fprintln(out, v)
	// }

	res := NewYearConcert(nums)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// 给定一个数组和一些查询，每次查询给出一个值x.
// !查询数组中有多少个子数组的gcd等于x.
// https://www.luogu.com.cn/problem/CF475D
func CGCDSSQ(nums []int, queries []int) []int {
	counter := LogTrick(nums, gcd, nil)
	res := make([]int, len(queries))
	for i, v := range queries {
		res[i] = counter[v]
	}
	return res
}

// https://www.luogu.com.cn/problem/CF1632D
// !一个数组不合法当且仅当数组的gcd等于其长度.
// !对nums的每个非空前缀,求最少修改次数，使得该前缀的所有子数组都合法.
// nums.length<=2e5, nums[i]<=1e9.
//
// !结论：如果某个前缀需要修改，那么最多改其中一个数即可(改成一个大质数)
func NewYearConcert(nums []int) []int {
	n := len(nums)
	res := make([]int, len(nums))

	dp := []Interval{}
	for pos, cur := range nums {
		for i, pre := range dp {
			dp[i].value = gcd(pre.value, cur)
		}
		dp = append(dp, Interval{leftStart: pos, leftEnd: pos + 1, value: cur})

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

		shouldModify := false
		for _, interval := range dp {
			minLen := pos - interval.leftEnd + 2
			maxLen := pos - interval.leftStart + 1
			value := interval.value
			if minLen <= value && value <= maxLen {
				shouldModify = true
				break
			}
		}
		if shouldModify {
			res[pos] = 1
			dp = dp[:0] // !如果需要修改，那么dp直接清空
		}
	}

	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + res[i-1]
	}
	return preSum[1:]
}

func maxGcdSum(nums []int, k int) int64 {
	res := 0
	preSum := make([]int, len(nums)+1)
	for i, v := range nums {
		preSum[i+1] = preSum[i] + v
	}
	LogTrick(nums, gcd, func(left []Interval, right int) {
		for _, v := range left {
			if right-v.leftStart+1 >= k {
				res = max(res, v.value*(preSum[right+1]-preSum[v.leftStart]))
			}
		}
	})
	return int64(res)
}

// TODO
// 3574. 最大子数组 GCD 分数
// https://leetcode.cn/problems/maximize-subarray-gcd-score/description/

type Interval = struct{ leftStart, leftEnd, value int }

// 将 nums 的所有非空子数组的元素进行 op 操作，返回所有不同的结果和其出现次数.
//
// nums: 1 <= nums.length <= 1e5.
// op: 与/或/gcd/lcm 中的一种操作，具有单调性.
// f:
// 数组的右端点为right.
// Interval 的 leftStart/leftEnd 表示子数组的左端点left的范围.
// Interval 的 value 表示该子数组 arr[left,right] 的 op 结果.
func LogTrick(nums []int, op func(int, int) int, f func(left []Interval, right int)) map[int]int {
	res := make(map[int]int)

	dp := []Interval{}
	for pos, cur := range nums {
		for i, pre := range dp {
			dp[i].value = op(pre.value, cur)
		}
		dp = append(dp, Interval{leftStart: pos, leftEnd: pos + 1, value: cur})

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
