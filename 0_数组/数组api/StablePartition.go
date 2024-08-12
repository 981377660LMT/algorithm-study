package main

import "fmt"

func main() {
	arr := []int32{1, 4, 0, 1, 3}
	Partition(arr, 0, int32(len(arr)), func(x int32) bool { return x > 1 })
	fmt.Println(arr)
}

func Partition[T any](arr []T, start, end int32, predicate func(T) bool) int32 {
	ptr := start
	for i := start; i < end; i++ {
		if predicate(arr[i]) {
			arr[ptr], arr[i] = arr[i], arr[ptr]
			ptr++
		}
	}
	return ptr
}

// 参考 24_高级数据结构/waveletmatrix/dynamic/WaveletMatrixActivable.go
func StablePartition[T any](arr []T, start, end int32, predicate func(T) bool) int32 {
	if start == end {
		return start
	}
	if start+1 == end {
		if predicate(arr[start]) {
			return start + 1
		} else {
			return start
		}
	}
	mid := (start + end) / 2
	mid = StablePartition(arr, start, mid, predicate)
	mid = StablePartition(arr, mid, end, predicate)
	return Partition(arr, start, end, predicate)
}
