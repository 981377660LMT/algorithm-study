package bfs

const inf int = 1e18

func Bfs(n int, adjList [][]int, start int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = inf
	}

	dist[start] = 0
	queue := []int{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adjList[cur] {
			if dist[next] > dist[cur]+1 {
				dist[next] = dist[cur] + 1
				queue = append(queue, next)
			}
		}
	}

	return
}
