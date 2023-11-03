package main

import "fmt"

// 原地删除数组中满足条件的元素.
func RemoveInplace(arr *([]int), shouldRemove func(index int) bool) {
	nums := *arr
	ptr := 0
	for i := 0; i < len(nums); i++ {
		if !shouldRemove(i) {
			nums[ptr] = nums[i]
			ptr++
		}
	}
	*arr = nums[:ptr]
}

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	RemoveInplace(&arr, func(index int) bool { return arr[index]%2 == 0 })
	fmt.Println(arr)
}
