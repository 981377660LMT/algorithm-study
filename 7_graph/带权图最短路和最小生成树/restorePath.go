package main

// 还原路径/dp复原.
func RestorePath(target int, pre []int) []int {
	path := []int{target}
	cur := target
	for pre[cur] != -1 {
		cur = pre[cur]
		path = append(path, cur)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// 还原路径/dp复原.
func RestorePath2(start, target int, pre []int) []int {
	path := []int{target}
	cur := target
	for cur != start {
		cur = pre[cur]
		path = append(path, cur)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}
