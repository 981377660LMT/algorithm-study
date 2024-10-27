package main

// https://leetcode.cn/problems/maximum-value-of-an-ordered-triplet-ii/description/
func maximumTripletValue(nums []int) int64 {
	P := NewPreSumSuffixSum(nums, func() int { return 0 }, func(a, b int) int { return max(a, b) })
	res := int64(0)
	for j := 1; j < len(nums)-1; j++ {
		preMax := P.PreSum(j)
		suffixMax := P.SuffixSum(j + 1)
		cand := int64((preMax - nums[j]) * (suffixMax))
		if cand > res {
			res = cand
		}
	}
	return res
}

type PreSumSuffixSum[E any] struct {
	preSum    []E
	suffixSum []E
	e         func() E
}

func NewPreSumSuffixSum[E any](arr []E, e func() E, op func(a, b E) E) *PreSumSuffixSum[E] {
	n := len(arr)
	preSum, suffixSum := make([]E, n+1), make([]E, n+1)
	preSum[0] = e()
	suffixSum[n] = e()
	for i := 0; i < n; i++ {
		preSum[i+1] = op(preSum[i], arr[i])
		suffixSum[n-i-1] = op(suffixSum[n-i], arr[n-i-1])
	}
	return &PreSumSuffixSum[E]{preSum: preSum, suffixSum: suffixSum, e: e}
}

// 查询前缀 `[0,end)` 的和.
func (p *PreSumSuffixSum[E]) PreSum(end int) E {
	if end < 0 {
		return p.e()
	}
	if end >= len(p.preSum) {
		return p.preSum[len(p.preSum)-1]
	}
	return p.preSum[end]
}

// 查询后缀 `[start,n)` 的和.
func (p *PreSumSuffixSum[E]) SuffixSum(start int) E {
	if start < 0 {
		return p.suffixSum[0]
	}
	if start >= len(p.suffixSum) {
		return p.e()
	}
	return p.suffixSum[start]
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
