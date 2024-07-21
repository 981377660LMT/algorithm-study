package main

import (
	"fmt"
	"sort"
)

func main() {
	// 记录排序过程中交换元素的下标
	r := swapRecorder{sort.IntSlice{3, 1, 2}, &[][2]int{}}
	sort.Sort(r)
	fmt.Println(r.swaps)
}

// 记录排序过程中交换元素的下标
// r := swapRecorder{a, &[][2]int{}}
// sort.Sort(r)
// 相关题目 https://codeforces.com/problemset/problem/266/C
type swapRecorder struct {
	sort.IntSlice
	swaps *[][2]int
}

func (r swapRecorder) Swap(i, j int) {
	// 快排时可能会有 i 和 j 相等的情况
	if i == j {
		return
	}
	*r.swaps = append(*r.swaps, [2]int{i, j})
	r.IntSlice[i], r.IntSlice[j] = r.IntSlice[j], r.IntSlice[i]
}
