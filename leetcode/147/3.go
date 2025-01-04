package main

func longestSubsequence(nums []int) int {
	n := len(nums)
	if n <= 1 {
		return n
	}

	maxDiff := maxs(nums...) - mins(nums...)
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, maxDiff+1)
		for d := 0; d <= maxDiff; d++ {
			dp[i][d] = 1
		}
	}
	dpPreMax := make([][]int, n)
	for j := 0; j < n; j++ {
		dpPreMax[j] = make([]int, maxDiff+1)
	}

	update := func(i int) {
		dpPreMax[i][maxDiff] = dp[i][maxDiff]
		for d := maxDiff - 1; d >= 0; d-- {
			if dp[i][d] > dpPreMax[i][d+1] {
				dpPreMax[i][d] = dp[i][d]
			} else {
				dpPreMax[i][d] = dpPreMax[i][d+1]
			}
		}
	}
	for i := 0; i < n; i++ {
		update(i)
	}

	res := 1
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			d := abs(nums[i] - nums[j])
			tmp := dpPreMax[j][d]
			dp[i][d] = max(dp[i][d], tmp+1)
			if dp[i][d] > res {
				res = dp[i][d]
			}
		}
		update(i)
	}
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
