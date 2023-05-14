package main

// https://leetcode.cn/problems/maximum-or/
func maximumOr(nums []int, k int) int64 {
	mul := ProductWithoutOne(nums, func() int { return 0 }, func(a, b int) int { return a | b })
	res := int64(0)
	for i := 0; i < len(nums); i++ {
		cur := int64(nums[i]<<k | mul[i])
		if cur > res {
			res = cur
		}
	}
	return res
}

type E = int

// 除自身以外数组的乘积.
func ProductWithoutOne(nums []E, e func() E, op func(E, E) E) []E {
	n := len(nums)
	res := make([]E, n)
	for i := 0; i < n-1; i++ {
		res[i+1] = op(res[i], nums[i])
	}
	x := e()
	for i := n - 1; i >= 0; i-- {
		res[i] = op(res[i], x)
		x = op(nums[i], x)
	}
	return res
}
