package main

// kthSuperset 返回 mask 的第k个超集, k 从 0 开始.
// 即将k填入mask的二进制中为0的部分, bmi2指令集中有pdep指令可以完成这个操作.
func kthSuperset(mask uint, k uint) uint {
	i, j := 0, 0
	for k>>j > 0 {
		if mask>>i&1 == 0 {
			mask |= k >> j & 1 << i
			j++
		}
		i++
	}
	return mask
}

// https://leetcode.cn/problems/minimum-array-end/description/
func minEnd(n int, x int) int64 {
	return int64(kthSuperset(uint(x), uint(n-1)))
}
