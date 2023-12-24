// 统计区间内的个数
// DivCount, DivCount2, DivSum, DivSum2

package main

import "fmt"

func main() {
	fmt.Println(DivCount1(1, 10, 2))
	fmt.Println(DivCount2(1, 10, 2, 1))
	fmt.Println(DivSum(1, 10, 2))
	fmt.Println(DivSum2(1, 10, 2, 1))
}

// 区间[lower,upper]内k的倍数个数.
func DivCount1(lower, upper int, k int) int {
	return upper/k - (lower-1)/k
}

// 区间[lower,upper]内形如k*x+b的个数.
func DivCount2(lower, upper int, k int, b int) int {
	return (upper-b)/k - (lower-1-b)/k
}

// 区间[lower,upper]内k的倍数和.
func DivSum(lower, upper int, k int) int {
	if lower > upper {
		return 0
	}

	var f func(right int) int
	f = func(right int) int {
		if right < k {
			return 0
		}
		first := k
		last := right / k * k
		count := (last-first)/k + 1
		return (first + last) * count / 2
	}

	return f(upper) - f(lower-1)
}

// 区间[lower,upper]内形如k*x+b的和.
func DivSum2(lower, upper int, k int, b int) int {
	//
}
