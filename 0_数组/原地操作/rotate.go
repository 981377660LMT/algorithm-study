package main

import (
	"fmt"
)

func SwapRange(arr []int, start, end int) {
	if start >= end {
		return
	}
	end--
	for start < end {
		arr[start], arr[end] = arr[end], arr[start]
		start++
		end--
	}
}

func RotateLeft(arr []int, start, end, step int) {
	n := end - start
	if n <= 1 || step == 0 {
		return
	}
	if step >= n {
		step %= n
	}
	SwapRange(arr, start, start+step)
	SwapRange(arr, start+step, end)
	SwapRange(arr, start, end)
}

func RotateRight(arr []int, start, end, step int) {
	n := end - start
	if n <= 1 || step == 0 {
		return
	}
	if step >= n {
		step %= n
	}
	SwapRange(arr, start, end-step)
	SwapRange(arr, end-step, end)
	SwapRange(arr, start, end)
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	RotateRight(nums, 0, len(nums), 3)
	fmt.Println(nums)
	RotateLeft(nums, 0, len(nums), 3)
	fmt.Println(nums)
	RotateRight(nums, 0, 1, 2)
	fmt.Println(nums)
}
