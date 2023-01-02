// !sort.Search 返回单调数组的前缀切片[0:n]中第一个满足 f(i) == true 的 i

package main

import (
	"fmt"
	"sort"
)

func main() {
	bisectLeft := sort.SearchInts
	// func SearchInts(a []int, x int) int {
	// 	return Search(len(a), func(i int) bool { return a[i] >= x })
	// }

	bisectRight := func(a []int, x int) int {
		return sort.Search(len(a), func(i int) bool { return a[i] > x })
	}

	nums := []int{1, 2, 3, 3, 3, 4, 5, 6}
	start := bisectLeft(nums, 2)    // 2
	end := bisectRight(nums, 3) - 1 // 4
	fmt.Println(start, end)
}

// func Search(n int, f func(int) bool) int {
// Define f(-1) == false and f(n) == true.
// Invariant: f(i-1) == false, f(j) == true.
// 	i, j := 0, n
// 	for i < j {
// 		h := int(uint(i+j) >> 1) // avoid overflow when computing h
// 		if !f(h) {
// 			i = h + 1 // preserves f(i-1) == false
// 		} else {
// 			j = h // preserves f(j) == true
// 		}
// 	}
// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
// 	return i
// }
