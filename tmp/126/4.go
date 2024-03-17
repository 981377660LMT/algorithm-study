package main

import "fmt"

// 给你一个长度为 n 的整数数组 nums 和一个 正 整数 k 。
// 一个整数数组的 能量 定义为和 等于 k 的子序列的数目。
// 请你返回 nums 中所有子序列的 能量和 。
// 由于答案可能很大，请你将它对 109 + 7 取余 后返回。
// 先求有多少种组成k的子集，再看这些子集属于多少个子序列
const MOD int = 1e9 + 7

// nums = [1,2,3], k = 3
func main() {
	fmt.Println(sumOfPower([]int{1, 2, 3}, 3))
}

func sumOfPower(nums []int, k int) int {
	res, ok := SubsetSumTargetDp1(nums, k)
	fmt.Println(res, ok)
	return 0
}

// 能否用nums中的若干个数凑出和为target.
//
//	O(n*max(nums)))
func SubsetSumTargetDp1(nums []int, target int) (res []int, ok bool) {
	if target <= 0 {
		return
	}

	n := len(nums)
	max_ := 0
	for _, v := range nums {
		max_ = max(max_, v)
	}
	right, curSum := 0, 0
	for right < n && curSum+nums[right] <= target {
		curSum += nums[right]
		right++
	}
	if right == n && curSum != target {
		return
	}

	offset := target - max_ + 1
	dp := make([]int, 2*max_)
	for i := range dp {
		dp[i] = -1
	}
	pre := make([][]int, n)
	for i := range pre {
		pre[i] = make([]int, 2*max_)
		for j := range pre[i] {
			pre[i][j] = -1
		}
	}

	dp[curSum-offset] = right
	for i := right; i < n; i++ {
		ndp := make([]int, len(dp))
		copy(ndp, dp)
		p := pre[i]
		a := nums[i]
		for j := 0; j < max_; j++ {
			if ndp[j+a] < dp[j] {
				ndp[j+a] = dp[j]
				p[j+a] = -2
			}
		}
		for j := 2*max_ - 1; j >= max_; j-- {
			for k := ndp[j] - 1; k >= max(dp[j], 0); k-- {
				if ndp[j-nums[k]] < k {
					ndp[j-nums[k]] = k
					p[j-nums[k]] = k
				}
			}
		}
		dp = ndp
	}

	if dp[max_-1] == -1 {
		return
	}

	used := make([]bool, n)
	i, j := n-1, max_-1
	for i >= right {
		p := pre[i][j]
		if p == -2 {
			used[i] = !used[i]
			j -= nums[i]
			i--
		} else if p == -1 {
			i--
		} else {
			used[p] = !used[p]
			j += nums[p]
		}
	}

	for i >= 0 {
		used[i] = !used[i]
		i--
	}

	for i := 0; i < n; i++ {
		if used[i] {
			res = append(res, i)
		}
	}

	ok = true
	return
}
