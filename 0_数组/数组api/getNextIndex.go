package main

import "fmt"

func main() {
	fmt.Println(GetPrevIndex([]int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}))
	fmt.Println(GetNextIndex([]int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}))
	fmt.Println(GetPrevAndNextIndex([]int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}))
}

// GetPrevIndex 获取数组中相同元素的前一个元素的位置.不存在则返回-1.
func GetPrevIndex[T comparable](arr []T) []int32 {
	n := int32(len(arr))
	left := make([]int32, n)
	last := map[T]int32{}
	for i := int32(0); i < n; i++ {
		cur := arr[i]
		if v, ok := last[cur]; ok {
			left[i] = v
		} else {
			left[i] = -1
		}
		last[cur] = i
	}
	return left
}

// GetNextIndex 获取数组中相同元素的后一个元素的位置.不存在则返回-1.
func GetNextIndex[T comparable](arr []T) []int32 {
	n := int32(len(arr))
	right := make([]int32, n)
	last := map[T]int32{}
	for i := n - 1; i >= 0; i-- {
		cur := arr[i]
		if v, ok := last[cur]; ok {
			right[i] = v
		} else {
			right[i] = -1
		}
		last[cur] = i
	}
	return right
}

func GetPrevAndNextIndex[T comparable](arr []T) ([]int32, []int32) {
	n := int32(len(arr))
	left, right := make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		left[i], right[i] = -1, -1
	}
	last := map[T]int32{}
	for i := int32(0); i < n; i++ {
		cur := arr[i]
		if v, ok := last[cur]; ok {
			left[i] = v
			right[v] = i
		}
		last[cur] = i
	}
	return left, right
}
