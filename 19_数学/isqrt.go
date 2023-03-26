// python 里的math.isqrt()函数使用牛顿迭代法,避免了浮点数转整数的误差
// https://leetcode.cn/problems/sqrtx/solution/x-de-ping-fang-gen-by-leetcode-solution/

package main

import "math"

// 求 x 的平方根的整数部分.
//  0 <= x < 1<<31
func ISqrt(x int) int {
	if x == 0 {
		return 0
	}
	C, x0 := float64(x), float64(x)
	for {
		xi := 0.5 * (x0 + C/x0)
		if math.Abs(x0-xi) < 1e-7 {
			break
		}
		x0 = xi
	}
	return int(x0)
}
