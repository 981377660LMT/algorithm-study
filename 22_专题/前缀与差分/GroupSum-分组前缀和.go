package main

// https://leetcode.cn/problems/ways-to-make-a-fair-array/
// LC2902 https://leetcode.cn/problems/count-of-sub-multisets-with-bounded-sum/

// 模分组前缀和/同余前缀和.
//  返回: 求下标在[start,end)范围内, 模mod为key的元素的和.
func GroupPresum(arr []int, mod int) func(start, end, key int) int {
	preSum := make([]int, len(arr)+mod)
	for i, v := range arr {
		preSum[i+mod] = preSum[i] + v
	}
	cal := func(r, k int) int {
		if r%mod <= k {
			return preSum[r/mod*mod+k]
		}
		return preSum[(r+mod-1)/mod*mod+k]
	}
	query := func(start, end, key int) int {
		if start >= end {
			return 0
		}
		key %= mod
		return cal(end, key) - cal(start, key)
	}
	return query
}
