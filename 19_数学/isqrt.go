// python 里的math.isqrt()函数使用牛顿迭代法,避免了浮点数转整数的误差
// https://leetcode.cn/problems/sqrtx/solution/x-de-ping-fang-gen-by-leetcode-solution/

package main

import (
	"fmt"
	"math"
)

func Isqrt(x int) int {
	x0 := x >> 1
	if x0 == 0 {
		return x
	}
	x1 := (x0 + x/x0) >> 1
	for x1 < x0 {
		x0 = x1
		x1 = (x0 + x/x0) >> 1
	}
	return x0
}

// 返回floor(sqrt(x)).x为负数时返回0.
func Isqrt2(x int) int {
	if x <= 0 {
		return 0
	}
	sqrt := int(math.Sqrt(float64(x)))
	for (sqrt+1)*(sqrt+1) <= x {
		sqrt++
	}
	for sqrt*sqrt > x {
		sqrt--
	}
	return sqrt
}

func main() {
	fmt.Println(Isqrt(5))
	// check 1e18 to 1e18+100
	for i := 0; i < 100; i++ {
		fmt.Println(Isqrt(1000000000000000000 + i))
	}
	fmt.Println(Isqrt2(5))
}
