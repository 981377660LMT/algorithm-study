package main

import (
	"fmt"
	"time"
)

func main() {
	time1 := time.Now()
	n := 11
	count := 0
	lens := make([]int, n)
	for i := 0; i < n; i++ {
		lens[i] = i + 1
	}
	EnumerateProduct(lens, func(indicesView []int) bool {
		count++
		return false
	})
	fmt.Println(count)
	fmt.Println(time.Since(time1))
}

// 遍历多个类数组对象的笛卡尔积.
// 11!(4e7) => 170ms.
func EnumerateProduct(lens []int, f func(indicesView []int) (shouldBreak bool)) {
	var dfs func(int, []int) bool
	dfs = func(index int, iters []int) bool {
		if index == len(iters) {
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
