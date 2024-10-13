package main

const INF int = 1e18

// O(n^2) 加边, O(1) 查询最短路.
type FloydQueryFast struct {
	dist [][]int
}

func NewFloydQueryFast(n int, directedEdges [][3]int) *FloydQueryFast {
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = INF
		}
		dist[i][i] = 0
	}
	for _, edge := range directedEdges {
		u, v, w := edge[0], edge[1], edge[2]
		dist[u][v] = min(w, dist[u][v])
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if dist[i][k] == INF {
				continue
			}
			for j := 0; j < n; j++ {
				if dist[k][j] == INF {
					continue
				}
				cand := dist[i][k] + dist[k][j]
				if dist[i][j] > cand {
					dist[i][j] = cand
				}
			}
		}
	}
	return &FloydQueryFast{dist: dist}
}

func (floyd *FloydQueryFast) AddDirectedEdge(u, v, w int) {
	n := len(floyd.dist)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cand := floyd.dist[i][u] + w + floyd.dist[v][j]
			if floyd.dist[i][j] > cand {
				floyd.dist[i][j] = cand
			}
		}
	}
}

// 求u到v的最短距离.
// 如果不存在从u到v的路径, 返回-1.
func (floyd *FloydQueryFast) Dist(u, v int) int {
	if floyd.dist[u][v] == INF {
		return -1
	}
	return floyd.dist[u][v]
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
