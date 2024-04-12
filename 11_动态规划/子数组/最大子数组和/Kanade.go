// MaximumContinuousSubsequence
package main

import "fmt"

func main() {
	fmt.Println(Kanade(func(i int) int { return i }, 0, 10))
}

const INF int = 2e18

// O(n)找到最大的一个子数组和（不能为空）.
func Kanade(f func(int) int, start, end int) int {
	prefix := 0
	max_ := -INF
	for i := start; i < end; i++ {
		prefix += f(i)
		max_ = max(max_, prefix)
		prefix = max(0, prefix)
		prefix = min(prefix, INF)
	}
	return max_
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
