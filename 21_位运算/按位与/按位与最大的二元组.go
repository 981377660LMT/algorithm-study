// 按位与最大的二元组
// !给定一个数组，要求找到两个不同的下标i!=j使得A[i]&A[j]最大.

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
)

// https://www.geeksforgeeks.org/maximum-value-pair-array/
func Solve(nums []int) int {
	countValid := func(pattern int) int {
		count := 0
		for _, num := range nums {
			if (pattern & num) == pattern {
				count++
			}
		}
		return count
	}

	res := 0
	max_ := 0
	for _, num := range nums {
		if num > max_ {
			max_ = num
		}
	}
	maxBit := bits.Len(uint(max_))
	for bit := maxBit; bit >= 0; bit-- {
		count := countValid(res | (1 << bit))
		if count >= 2 {
			res |= 1 << bit
		}
	}
	return res
}

func main() {
	solve2 := func(nums []int) int {
		res := 0
		for i := 0; i < len(nums); i++ {
			for j := i + 1; j < len(nums); j++ {
				res = max(res, nums[i]&nums[j])
			}
		}
		return res
	}

	for i := 0; i < 1000; i++ {
		nums := make([]int, 2)
		for j := range nums {
			nums[j] = rand.Intn(10)
		}
		actual := Solve(nums)
		expected := solve2(nums)
		if expected != actual {
			fmt.Println(expected, actual, nums)
			panic("not equal1")
		}
	}

	fmt.Println("ok")
}
