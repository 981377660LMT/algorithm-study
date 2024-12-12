// 原地过滤
// 记住，这种技术最适用于以下情况：
//
// 在过滤后不再需要 numbers切片。
// 性能至关重要，特别是在处理大型数据集的时候。

package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	filtered := nums[:0] // 共享底层数组
	isOdd := func(n int) bool { return n%2 == 1 }
	for _, num := range nums {
		if isOdd(num) {
			filtered = append(filtered, num)
		}
	}

	// [1 3 5 7 9]
	// [1 3 5 7 9 6 7 8 9 10]
	fmt.Println(filtered)
	fmt.Println(nums)
}
