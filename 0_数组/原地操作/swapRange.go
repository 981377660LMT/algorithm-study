package main

import "fmt"

// ReverseRange
func SwapRange(s []int, i, j int) {
	for i < j {
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	SwapRange(nums, 2, 6)
	fmt.Println(nums)
}
