package main

const INF int = 1e18

// bfs求最无权图短路.
func Bfs(start int, adjList [][]int) []int {
	n := len(adjList)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := []int{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adjList[cur] {
			cand := dist[cur] + 1
			if cand < dist[next] {
				dist[next] = cand
				queue = append(queue, next)
			}
		}
	}
	return dist
}

// 多源bfs.
func BfsMultiStart(starts []int, adjList [][]int) []int {
	n := len(adjList)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	queue := append(starts[:0:0], starts...)
	for _, start := range starts {
		dist[start] = 0
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adjList[cur] {
			cand := dist[cur] + 1
			if cand < dist[next] {
				dist[next] = cand
				queue = append(queue, next)
			}
		}
	}
	return dist
}

// bfs求起点到终点的最短距离和路径.
func BfsPath(n int, adjList [][]int, start int, end int) (res int, path []int) {
	dist := make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	queue := []int{start}
	pre := make([]int, n)
	for i := range pre {
		pre[i] = -1
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adjList[cur] {
			cand := dist[cur] + 1
			if cand < dist[next] {
				dist[next] = cand
				pre[next] = cur
				queue = append(queue, next)
			}
		}
	}
	if dist[end] == INF {
		return INF, []int{}
	}
	cur := end
	for pre[cur] != -1 {
		path = append(path, cur)
		cur = pre[cur]
	}
	path = append(path, start)
	for i := 0; i < len(path)/2; i++ {
		path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
	}
	return dist[end], path
}

// 返回距离start为dist的结点.
func BfsDepth(n int, adjList [][]int, start int, dist int) []int {
	if dist < 0 {
		return []int{}
	}
	if dist == 0 {
		return []int{start}
	}

	queue := []int{start}
	visited := make([]bool, n)
	todo := dist
	for len(queue) > 0 && todo > 0 {
		len_ := len(queue)
		for i := 0; i < len_; i++ {
			cur := queue[0]
			queue = queue[1:]
			for _, next := range adjList[cur] {
				if !visited[next] {
					visited[next] = true
					queue = append(queue, next)
				}
			}
		}
		todo--
	}

	return queue
}

// 网格图bfs, 返回每个格子到起点的最短距离.
func BfsGrid(row int, col int, starts [][]int) [][]int {
	DIR4 := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	dist := make([][]int, row)
	for i := range dist {
		dist[i] = make([]int, col)
		for j := range dist[i] {
			dist[i][j] = INF
		}
	}

	queue := append(starts[:0:0], starts...)
	for _, start := range starts {
		dist[start[0]][start[1]] = 0
	}

	for len(queue) > 0 {
		len_ := len(queue)
		for i := 0; i < len_; i++ {
			curX, curY := queue[0][0], queue[0][1]
			queue = queue[1:]
			for _, dir := range DIR4 {
				nextX, nextY := curX+dir[0], curY+dir[1]
				cand := dist[curX][curY] + 1
				if 0 <= nextX && nextX < row && 0 <= nextY && nextY < col && cand < dist[nextX][nextY] {
					dist[nextX][nextY] = cand
					queue = append(queue, []int{nextX, nextY})
				}
			}
		}
	}

	return dist
}
