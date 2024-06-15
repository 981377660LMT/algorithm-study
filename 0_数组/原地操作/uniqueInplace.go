package main

import "fmt"

// 有序数组原地去重.
// Compact/dedup.
func UniqueInplace(sorted *([]int)) {
	if len(*sorted) == 0 {
		return
	}
	nums := *sorted
	size := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != nums[size] {
			size++
			nums[size] = nums[fast]
		}
	}
	size++
	*sorted = nums[:size]
}

func main() {
	nums := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	UniqueInplace(&nums)
	fmt.Println(nums)
}
