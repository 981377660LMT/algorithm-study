// 1819. 序列中不同最大公约数的数目
// 统计数组中所有子序列的gcd个数，复杂度 O(maxlog^2max)
// https://leetcode.cn/problems/number-of-different-subsequences-gcds/solution/ji-bai-100mei-ju-gcdxun-huan-you-hua-pyt-get7/
//
// 1<=nums.length<=1e5
// 1<=nums[i]<=2e5

package main

func countDifferentSubsequenceGCDs(nums []int) (res int) {
	max := 0
	for _, v := range nums {
		if v > max {
			max = v
		}
	}

	has := make([]bool, max+1)
	for _, v := range nums {
		has[v] = true
	}

	// !枚举答案
	for i := 1; i <= max; i++ {
		gcd_ := 0
		for j := i; j <= max; j += i { // 枚举 i 的倍数 j
			if has[j] { // 如果 j 在 nums 中
				gcd_ = gcd(gcd_, j) // 更新最大公约数  TODO: 可以用光速 gcd
				if gcd_ == i {      // 找到一个答案
					res++
					break
				}
			}
		}
	}

	return
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
