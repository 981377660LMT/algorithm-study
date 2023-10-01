package main

import "fmt"

func main() {
	fmt.Println(countInvMergeSort([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}))
}

// 归并排序求逆序对.
func countInvMergeSort(arr []int) int {
	if len(arr) < 2 {
		return 0
	}
	if len(arr) == 2 {
		if arr[0] > arr[1] {
			return 1
		}
		return 0
	}
	res := 0
	midCount := 0
	upper := make([]int, 0)
	lower := make([]int, 0)
	mid := arr[len(arr)/2]
	for i := 0; i < len(arr); i++ {
		num := arr[i]
		if num < mid {
			lower = append(lower, num)
			res += len(upper)
			res += midCount
		} else if num > mid {
			upper = append(upper, num)
		} else {
			midCount++
			res += len(upper)
		}
	}
	res += countInvMergeSort(lower)
	res += countInvMergeSort(upper)
	return res
}
