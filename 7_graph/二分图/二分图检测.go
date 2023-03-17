// 二分图检测

package main

func isBipartite(n int, graph [][]int) (colors []int, ok bool) {
	colors = make([]int, n)
	for i := range colors {
		colors[i] = -1
	}
	bfs := func(start int) bool {
		colors[start] = 0
		queue := []int{start}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, next := range graph[cur] {
				if colors[next] == -1 {
					colors[next] = colors[cur] ^ 1
					queue = append(queue, next)
				} else if colors[next] == colors[cur] {
					return false
				}
			}
		}
		return true
	}

	for i := range colors {
		if colors[i] == -1 && !bfs(i) {
			return nil, false
		}
	}
	return colors, true
}
