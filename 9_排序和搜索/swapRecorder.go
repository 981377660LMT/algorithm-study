// SortRecorder/SwapRecorder
// 记录排序过程中交换元素的下标.
// 相关题目 https://codeforces.com/problemset/problem/266/C
//
// Sort
// 用于记录排序过程中的交换.

package main

import (
	"fmt"
	"sort"
)

func main() {
	{
		nums := []int{1, 2, 3, 5, 4}
		beforeSwap := func(i, j int) {
			fmt.Println(fmt.Sprintf("swap %d %d", i, j))
		}
		R := NewSwapRecorder(sort.IntSlice(nums), beforeSwap)
		sort.Sort(R)
	}

	{
		nums := []int{1, 2, 3, 5, 4}
		less := func(i, j int) bool {
			return nums[i] < nums[j]
		}
		swap := func(i, j int) {
			fmt.Println(fmt.Sprintf("swap %d %d", i, j))
			nums[i], nums[j] = nums[j], nums[i]
		}
		Sort(len(nums), less, swap)
	}
}

type SwapRecorder struct {
	sort.Interface
	beforeSwap func(i, j int)
}

func NewSwapRecorder(sortInterface sort.Interface, beforeSwap func(i, j int)) *SwapRecorder {
	return &SwapRecorder{Interface: sortInterface, beforeSwap: beforeSwap}
}

func (r *SwapRecorder) Swap(i, j int) {
	if i == j {
		return
	}
	r.beforeSwap(i, j)
	r.Interface.Swap(i, j)
}

// /////////////////////////
// ///////////////////////////
// ///////////////////////////
// 用于记录排序过程中的交换.
func Sort(n int, less func(i, j int) bool, swap func(i, j int)) {
	sort.Sort(&sorter{n: n, less: less, swap: swap})
}

type sorter struct {
	n    int
	less func(i, j int) bool
	swap func(i, j int)
}

func (s *sorter) Len() int           { return s.n }
func (s *sorter) Less(i, j int) bool { return s.less(i, j) }
func (s *sorter) Swap(i, j int) {
	if i == j {
		return
	}
	s.swap(i, j)
}
