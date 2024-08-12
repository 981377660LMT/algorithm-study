package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	StablePartition(nums, 0, int32(len(nums)), func(i int32) bool { return nums[i]%2 == 0 })
	fmt.Println(nums) // Output: [4 0 1 1 3]
}

func Partition[T any](arr []T, start, end int32, predicate func(i int32) bool) int32 {
	ptr := start
	for i := start; i < end; i++ {
		if predicate(i) {
			arr[ptr], arr[i] = arr[i], arr[ptr]
			ptr++
		}
	}
	return ptr
}

func StablePartition[T any](arr []T, start, end int32, predicate func(i int32) bool) int32 {
	ptr := start
	var buffer []T
	for i := start; i < end; i++ {
		if predicate(i) {
			arr[ptr] = arr[i]
			ptr++
		} else {
			buffer = append(buffer, arr[i])
		}
	}
	copy(arr[ptr:], buffer)
	return ptr
}
