package main

import "fmt"

func main() {
	nums := []int{1, 2, 3}
	EnumerateSetPartition(nums, func(groups [][]E) bool {
		fmt.Println(groups)
		return false
	})

	// [[1] [2] [3]]
	// [[1 3] [2]]
	// [[1] [2 3]]
	// [[1 2] [3]]
	// [[1 2 3]]
}

type E = int

// 遍历所有的集合划分.
// f: 返回 true 表示停止遍历.
func EnumerateSetPartition(arr []E, f func(groups [][]E) bool) {
	n := len(arr)
	groups := [][]E{}

	var dfs func(int) bool
	dfs = func(pos int) bool {
		if pos == n {
			return f(groups)
		}
		groups = append(groups, []E{arr[pos]})
		dfs(pos + 1)
		groups = groups[:len(groups)-1]
		for i := range groups {
			groups[i] = append(groups[i], arr[pos])
			if dfs(pos + 1) {
				return true
			}
			groups[i] = groups[i][:len(groups[i])-1]
		}
		return false
	}

	dfs(0)
}
