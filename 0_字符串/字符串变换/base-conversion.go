// 进制转换

package main

import "fmt"

func main() {
	fmt.Println(ToBase10(3, func(i int) int { return []int{0, 1, 1}[i] }, 2)) // 3
	fmt.Println(ToBaseN(6, 2))                                                // [1 1 0]
}

func ToBase10(n int, f func(int) int, fromBase int) int {
	res := 0
	for i := 0; i < n; i++ {
		res = res*fromBase + f(i)
	}
	return res
}

func ToBaseN(x, toBase int) []int {
	if x == 0 {
		return []int{0}
	}

	res := []int{}
	for x > 0 {
		t := x % toBase
		res = append(res, t)
		x /= toBase
	}

	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
