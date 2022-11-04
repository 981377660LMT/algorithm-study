package spfa

type Edge struct{ next, weight int }

const inf int = 1e18

// spfa 求带有负权边的最短路，时间复杂度 O(V*E)，有负环时返回 nil
//   https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L1276
//   !只是找负环的话，初始时将所有点入队即可
func SPFA(n int, adjList [][]Edge, start int) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = inf
	}

	dist[start] = 0
	queue := []int{start}
	inQueue := make([]bool, n)
	inQueue[start] = true
	relaxedConut := make([]int, n)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		inQueue[cur] = false
		for _, edge := range adjList[cur] {
			next, weight := edge.next, edge.weight
			if cand := dist[cur] + weight; cand < dist[next] {
				dist[next] = cand
				relaxedConut[next] = relaxedConut[cur] + 1
				// 找到一个从 start 出发可达的负环
				if relaxedConut[next] >= n {
					return nil
				}
				if !inQueue[next] {
					queue = append(queue, next)
					inQueue[next] = true
				}
			}
		}
	}

	return
}
