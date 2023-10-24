// 和等于乘积的子数组个数

package main

import "fmt"

func main() {
	fmt.Println(CountSubarrayWithSumEqualToProduct([]int{1, 3, 2, 2}))
}

// 给一数组 a，元素均为正整数，求区间和等于区间积的区间个数
// n<=2e5,nums[i]<=2e5
// https://www.dotcpp.com/oj/problem2622.html
//
// 考虑对每个区间右端点，有多少个合法的左端点，使得区间和等于区间积。
// 由于乘积至少要乘 2 才会变化，所以对于一个固定的区间右端点，不同的区间积至多有 O(log(sum(a))) 个
// 只需要在加入一个新的数后，去重并去掉区间积超过 sum(a) 的区间，就可以暴力做出此题
func CountSubarrayWithSumEqualToProduct(nums []int) int {
	type interval = struct{ leftStart, leftEnd, value int }
	allSum := 0
	for _, v := range nums {
		allSum += v
	}

	res := 0
	preSum := map[int]int{0: 0} // 前缀和互不相同
	curSum := 0
	dp := []interval{} // value记录乘积

	for pos, cur := range nums {
		curSum += cur
		for i := range dp {
			dp[i].value *= cur
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

		// 去掉超过 allSum 的，从而保证 dp 中至多有 O(log(allSum)) 个元素
		for len(dp) > 0 && dp[0].value > allSum {
			dp = dp[1:]
		}

		// 将区间[0,pos]分成了dp.length个左闭右开区间.
		// 每一段区间的左端点left范围 在 [dp[i].leftStart,dp[i].leftEnd) 中。
		// 对应子数组 arr[left:pos+1] 的 op 值为 dp[i].value.
		for _, v := range dp {
			if pos, ok := preSum[curSum-v.value]; ok && v.leftStart <= pos && pos < v.leftEnd {
				res++
			}
		}
		preSum[curSum] = pos + 1

	}

	return res
}
