// 基础：
//   EnumerateProduct/EnumerateProduct32
// 带num:
//   EnumerateProductWithNum
// 初始化高维前缀和：
//   EnumeratePrefix

package main

import (
	"fmt"
	"time"
)

func main() {
	time1 := time.Now()
	n := int32(11)
	count := 0
	lens := make([]int32, n)
	for i := int32(0); i < n; i++ {
		lens[i] = i + 1
	}

	EnumerateProduct32(lens, func(indicesView []int32) bool {
		count++
		return false
	})
	fmt.Println(count)
	fmt.Println(time.Since(time1))

	fmt.Println("=====")
	EnumerateProductWithNum([]int32{3, 3}, func(digits []int32, num int32) bool {
		count++
		fmt.Println(digits, num)
		return false
	})

	fmt.Println("=====")
	EnumeratePrefix([]int32{3, 3}, func(indices []int32, mul int32) {
		count++
		fmt.Println(indices, mul)
	})
}

// 遍历多个类数组对象的笛卡尔积.
// 11!(4e7) => 170ms.
func EnumerateProduct(lens []int, f func(indicesView []int) (shouldBreak bool)) {
	n := len(lens)
	var dfs func(int, []int) bool
	dfs = func(index int, iters []int) bool {
		if index == n {
			return f(iters)
		}
		for iters[index] = 0; iters[index] < lens[index]; iters[index]++ {
			if dfs(index+1, iters) {
				return true
			}
		}
		return false
	}
	dfs(0, make([]int, len(lens)))
}

func EnumerateProduct32(lens []int32, f func(indicesView []int32) (shouldBreak bool)) {
	n := int32(len(lens))
	var dfs func(int32, []int32) bool
	dfs = func(index int32, iters []int32) bool {
		if index == n {
			return f(iters)
		}
		for iters[index] = 0; iters[index] < lens[index]; iters[index]++ {
			if dfs(index+1, iters) {
				return true
			}
		}
		return false
	}
	dfs(0, make([]int32, len(lens)))
}

// EnumerateDigits.
func EnumerateProductWithNum(sizes []int32, f func(digits []int32, num int32) (shouldBreak bool)) {
	n := int32(len(sizes))
	var dfs func(int32, []int32, int32, int32) bool
	dfs = func(pos int32, digits []int32, num int32, base int32) bool {
		if pos == -1 {
			return f(digits, num)
		}
		for digits[pos] = 0; digits[pos] < sizes[pos]; digits[pos]++ {
			if dfs(pos-1, digits, num+digits[pos]*base, base*sizes[pos]) {
				return true
			}
		}
		return false
	}
	dfs(n-1, make([]int32, len(sizes)), 0, 1)
}

// 用于初始化高维前缀和.
//
//	用坐标为digits的点的值填充presum[num].
func EnumeratePrefix(sizes []int32, f func(indices []int32, num int32)) {
	n := int32(len(sizes))
	var dfs func(int32, []int32, int32, int32)
	dfs = func(pos int32, indicies []int32, num int32, base int32) {
		if pos == -1 {
			f(indicies, num)
			return
		}
		for indicies[pos] = 0; indicies[pos] < sizes[pos]; indicies[pos]++ {
			dfs(pos-1, indicies, num+(indicies[pos]+1)*base, base*(sizes[pos]+1))
		}
	}
	dfs(n-1, make([]int32, len(sizes)), 0, 1)
}

func EnumerateDigits(mins []int32, maxs []int32, bases []int32, f func(num int32)) {
	n := int32(len(bases))
	var dfs func(int32, int32)
	dfs = func(pos int32, num int32) {
		if pos == n {
			f(num)
			return
		}
		for i := mins[pos]; i <= maxs[pos]; i++ {
			dfs(pos+1, num+bases[pos]*i)
		}
	}
	dfs(0, 0)
}
