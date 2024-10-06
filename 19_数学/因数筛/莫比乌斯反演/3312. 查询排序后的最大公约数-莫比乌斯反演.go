// 3312. 查询排序后的最大公约数
// https://leetcode.cn/problems/sorted-gcd-pair-queries/description/
// !求第k大的公约数对gcd(nums[i], nums[j]), kthGcdPair
//
// 莫比乌斯反演(莫反)
// !本质上是容斥原理

package main

import "sort"

func gcdValues(nums []int, queries []int64) []int {
	upper := maxs(nums...) + 1
	c1, c2 := make([]int, upper), make([]int, upper)
	for _, v := range nums {
		c1[v]++
	}
	for f := 1; f < upper; f++ {
		for m := f; m < upper; m += f {
			c2[f] += c1[m]
		}
	}
	for i := 1; i < upper; i++ {
		c2[i] = c2[i] * (c2[i] - 1) / 2
	}
	for i := upper - 1; i > 0; i-- {
		for j := 2 * i; j < upper; j += i {
			c2[i] -= c2[j]
		}
	}

	presum := make([]int, len(c2))
	presum[0] = c2[0]
	for i := 1; i < len(c2); i++ {
		presum[i] = presum[i-1] + c2[i]
	}
	res := make([]int, len(queries))
	for i, kth := range queries {
		res[i] = sort.SearchInts(presum, int(kth)+1)
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
