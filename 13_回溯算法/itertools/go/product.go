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

	EnumerateProduct32ByLex([]int32{3, 4}, func(indicesView []int32, num int32) bool {
		count++
		fmt.Println(indicesView, num)
		return false
	})

	EnumeratePrefix([]int32{3, 3}, func(indicesView []int32, mul int32) bool {
		count++
		fmt.Println(indicesView, mul)
		return false
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
func EnumerateProduct32ForNum(sizes []int32, f func(digits []int32, num int32) (shouldBreak bool)) {
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

// 遍历高维前缀和.
func EnumeratePrefix(sizes []int32, f func(digits []int32, num int32) bool) {
	n := int32(len(sizes))
	var dfs func(int32, []int32, int32, int32) bool
	dfs = func(pos int32, digits []int32, num int32, base int32) bool {
		if pos == -1 {
			return f(digits, num)
		}
		for digits[pos] = 0; digits[pos] < sizes[pos]; digits[pos]++ {
			if dfs(pos-1, digits, num+(digits[pos]+1)*base, base*(sizes[pos]+1)) {
				return true
			}
		}
		return false
	}
	dfs(n-1, make([]int32, len(sizes)), 0, 1)
}

// 按照字典序遍历多个类数组对象的笛卡尔积.
func EnumerateProduct32ByLex(sizes []int32, f func(indicesView []int32, num int32) bool) {

	var dfs func(int32, []int32, int32, int32) bool
	dfs = func(index int32, iters []int32, num int32, base int32) bool {
		if index == -1 {
			return f(iters, num)
		}
		for iters[index] = 0; iters[index] < sizes[index]; iters[index]++ {
			if dfs(index-1, iters, num+(iters[index])*base, base*(sizes[index])) {
				return true
			}
		}
		return false
	}
	dfs(int32(len(sizes))-1, make([]int32, len(sizes)), 0, 1)
}
