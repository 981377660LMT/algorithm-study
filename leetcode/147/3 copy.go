package main

import "time"

func longestSubsequence(nums []int) int {
	n := len(nums)
	if n <= 1 {
		return n
	}

	maxDiff := maxs(nums...) - mins(nums...)
	dp := make([][]uint16, n)
	prefix := make([][]uint16, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]uint16, maxDiff+1)
		prefix[i] = make([]uint16, maxDiff+1)
	}

	res := uint16(1)
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			d := abs(nums[i] - nums[j])
			dp[i][d] = max16(dp[i][d], 1+prefix[j][d])
		}
		prefix[i][maxDiff] = dp[i][maxDiff]
		for d := maxDiff - 1; d >= 0; d-- {
			prefix[i][d] = max16(dp[i][d], prefix[i][d+1])
		}
		res = max16(res, prefix[i][0])
	}

	return int(res) + 1
}

func max16(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func main() {
	n := int(1e4)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = i
	}
	cur := time.Now()
	longestSubsequence(nums)
	println(time.Since(cur).Milliseconds())
}
