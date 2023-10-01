package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	copy(nums[3:], nums[2:])
	fmt.Println(nums)
}
