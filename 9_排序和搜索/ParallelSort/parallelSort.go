package main

import "fmt"

func main() {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	SortDemo(nums, func(i, j int) bool { return nums[i] < nums[j] })
	fmt.Println(nums)
}

func SortDemo(nums []int, less func(i, j int) bool) {
	ops := ParallelSort(len(nums))
	cas := func(i, j int) {
		if less(j, i) {
			nums[i], nums[j] = nums[j], nums[i]
		}
	}
	for _, op := range ops {
		for _, p := range op { // can be parallelized
			cas(p[0], p[1])
		}
	}
}

// 双调排序Bitonic Sort，适合并行计算的排序算法
//
// Batchers algorithm is a parallel sorting algorithm that is based on the merge sort algorithm.
// It is a comparison-based algorithm that is used to sort a list of elements in parallel.
//
// Each operation consists of disjoint pairs (compare and swap).
// Batcher's algorithm
// (1 + 2 + ... + ceil(log_2(n))) operations
func ParallelSort(n int) (ops [][][2]int) {
	if n < 0 {
		panic("invalid input")
	}
	for m := 1; m < n; m <<= 1 { // sorted blocks of m elements
		for d := m; d >= 1; d >>= 1 { // operate on pairs of distance d
			ops = append(ops, [][2]int{})
			for i := 0; i < n; i += (m << 1) {
				for j := i + d%m; j+d < i+(m<<1); j += (d << 1) {
					for k := j; k < j+d && k+d < n; k++ {
						ops[len(ops)-1] = append(ops[len(ops)-1], [2]int{k, k + d})
					}
				}
			}
		}
	}
	return ops
}
