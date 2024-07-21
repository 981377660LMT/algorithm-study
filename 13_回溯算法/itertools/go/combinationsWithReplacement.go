package main

import (
	"fmt"
	"time"
)

func main() {
	EnumerateCombinationsWithReplacement(3, 2, func(indicesView []int) bool {
		fmt.Println(indicesView)
		return false
	})

	time1 := time.Now()
	count := 0
	EnumerateCombinationsWithReplacement(20, 10, func(indicesView []int) bool {
		count++
		return false
	})
	fmt.Println(count)
	fmt.Println(time.Since(time1))
}

// 从 n 个元素中选择 k 个元素，允许重复选择同一个元素，按字典序生成所有组合，每个组合用下标表示
// https://docs.python.org/3/library/itertools.html#itertools.combinations_with_replacement
// https://en.wikipedia.org/wiki/Combination#Number_of_combinations_with_repetition
// 方案数 H(n,k)=C(n+k-1,k) https://oeis.org/A059481
// 相当于长度为 k，元素范围在 [0,n-1] 的非降序列的个数
// 2e7 => 100ms
func EnumerateCombinationsWithReplacement(n, r int, f func(indicesView []int) (shouldBreak bool)) {
	ids := make([]int, r)
	if f(ids) {
		return
	}
	for {
		i := r - 1
		for ; i >= 0; i-- {
			if ids[i] != n-1 {
				break
			}
		}
		if i == -1 {
			return
		}
		ids[i]++
		for j := i + 1; j < r; j++ {
			ids[j] = ids[i]
		}
		if f(ids) {
			return
		}
	}
}
