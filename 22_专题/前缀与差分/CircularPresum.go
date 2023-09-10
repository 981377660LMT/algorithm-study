package main

// 环形数组前缀和
func CircularPresum(nums []int) func(start int, end int) int {
	n := len(nums)
	preSum := make([]int, n+1)
	for i, v := range nums {
		preSum[i+1] = preSum[i] + v
	}

	cal := func(r int) int {
		return preSum[n]*(r/n) + preSum[r%n]
	}

	query := func(start int, end int) int {
		if start >= end {
			return 0
		}
		return cal(end) - cal(start)
	}

	return query
}
