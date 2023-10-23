package main

import "fmt"

func main() {
	nums := []int{1, 2, 3}
	EnumerateSetPartition(nums, func(groups [][]E) {
		fmt.Println(groups)
	})

	// [[1] [2] [3]]
	// [[1 3] [2]]
	// [[1] [2 3]]
	// [[1 2] [3]]
	// [[1 2 3]]
}

type E = int

// 遍历所有的集合划分.
func EnumerateSetPartition(arr []E, f func(groups [][]E)) {
	n := len(arr)
	groups := [][]E{}

	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			f(groups)
			return
		}
		groups = append(groups, []E{arr[pos]})
		dfs(pos + 1)
		groups = groups[:len(groups)-1]
		for i := range groups {
			groups[i] = append(groups[i], arr[pos])
			dfs(pos + 1)
			groups[i] = groups[i][:len(groups[i])-1]
		}
	}

	dfs(0)
}
