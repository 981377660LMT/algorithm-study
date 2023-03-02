// https://nyaannyaan.github.io/library/set-function/enumerate-set.hpp
// 枚举子集/超集

package main

import "fmt"

func main() {
	fmt.Println(EnumerateSubset(5))      // 101 100 001 000
	fmt.Println(EnumerateSuperset(5, 3)) // 101 111
}

// 枚举子集
func EnumerateSubset(b int) []int {
	var res []int
	for s := b; s >= 0; s-- {
		s &= b
		res = append(res, s)
	}
	return res
}

// 枚举超集
func EnumerateSuperset(b, n int) []int {
	var res []int
	for s := b; s < (1 << n); s = (s + 1) | b {
		res = append(res, s)
	}
	return res
}
