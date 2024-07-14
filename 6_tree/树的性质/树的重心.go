package main

import "fmt"

func main() {
	n := 5
	tree := [][]int{
		{1, 2},
		{0, 3},
		{0, 4},
		{1},
		{2},
	}
	root := 0
	centroids := FindCentroids(n, tree, root)
	fmt.Println(centroids)
}

// getCentroids
func FindCentroids(n int, tree [][]int, root int) (centroids []int) {
	weight := make([]int, n)
	subSize := make([]int, n)
	var dfs func(cur, pre int)
	dfs = func(cur, pre int) {
		subSize[cur] = 1
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			dfs(next, cur)
			subSize[cur] += subSize[next]
			weight[cur] = max(weight[cur], subSize[next])
		}
		weight[cur] = max(weight[cur], n-subSize[cur])
		if weight[cur] <= n/2 {
			centroids = append(centroids, cur)
		}
	}

	dfs(root, -1)
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func FindCentroids32(n int32, tree [][]int32, root int32) (centroids []int32) {
	weight := make([]int32, n)
	subSize := make([]int32, n)
	var dfs func(cur, pre int32)
	dfs = func(cur, pre int32) {
		subSize[cur] = 1
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			dfs(next, cur)
			subSize[cur] += subSize[next]
			weight[cur] = max32(weight[cur], subSize[next])
		}
		weight[cur] = max32(weight[cur], n-subSize[cur])
		if weight[cur] <= n/2 {
			centroids = append(centroids, cur)
		}
	}

	dfs(root, -1)
	return
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
