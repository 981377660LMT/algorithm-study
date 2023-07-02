package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	nums := make([]int, 1e6)
	for i := range nums {
		nums[i] = rand.Intn(1e9)
	}
	time1 := time.Now()
	SortRange(nums, func(i, j int) bool { return nums[i] < nums[j] }, 0, 1e6)
	fmt.Println(time.Since(time1)) // 183.7696ms

	nums = make([]int, 1e6)
	for i := range nums {
		nums[i] = rand.Intn(1e9)
	}
	time1 = time.Now()
	SortRangeStable(nums, func(i, j int) bool { return nums[i] < nums[j] }, 0, 1e6)
	fmt.Println(time.Since(time1)) // 561.8088ms
}

func SortRange(arr []int, less func(i, j int) bool, start, end int) {
	n := len(arr)
	if start < 0 {
		start = 0
	}
	if end > n {
		end = n
	}
	if start >= end {
		return
	}

	sort.Slice(arr[start:end], func(i, j int) bool { return less(i+start, j+start) })
}

func SortRangeStable(arr []int, less func(i, j int) bool, start, end int) {
	n := len(arr)
	if start < 0 {
		start = 0
	}
	if end > n {
		end = n
	}
	if start >= end {
		return
	}

	sort.SliceStable(arr[start:end], func(i, j int) bool { return less(i+start, j+start) })
}
