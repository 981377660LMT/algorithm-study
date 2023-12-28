package main

// 还原路径/dp复原.
func RestorePath(target int, pre []int) []int {
	path := []int{target}
	for pre[path[len(path)-1]] != -1 {
		path = append(path, pre[path[len(path)-1]])
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
