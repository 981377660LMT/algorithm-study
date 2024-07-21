package main

import (
	"fmt"
	"time"
)

func main() {

	time1 := time.Now()
	count := 0

	EnumerateCombinations(30, 10, func(indicesView []int) bool {
		count++
		return false
	})
	fmt.Println(count)
	fmt.Println(time.Since(time1))

}

// 从 n 个元素中选择 r 个元素，按字典序生成所有组合，每个组合用下标表示  r <= n
// https://docs.python.org/3/library/itertools.html#itertools.combinations
// https://stackoverflow.com/questions/41694722/algorithm-for-itertools-combinations-in-python
// C(30,10)(3e7) => 150ms.
func EnumerateCombinations(n, r int, f func(indicesView []int) (shouldBreak bool)) {
	ids := make([]int, r)
	for i := range ids {
		ids[i] = i
	}
	if f(ids) {
		return
	}
	for {
		i := r - 1
		for ; i >= 0; i-- {
			if ids[i] != i+n-r {
				break
			}
		}
		if i == -1 {
			return
		}
		ids[i]++
		for j := i + 1; j < r; j++ {
			ids[j] = ids[j-1] + 1
		}
		if f(ids) {
			return
		}
	}
}
