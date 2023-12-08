// 无向图经过每个顶点的最小环 O(n^3)
// 如果不存在环，返回INF.
// Minimum Weight Cycle For Each Vertex

package main

import "fmt"

const INF int = 1e18

func main() {
	fmt.Println(MinimumWeightCycleForEachVertex(4, [][3]int{{0, 1, 1}, {1, 2, 2}, {2, 3, 3}, {3, 0, 4}, {0, 2, 3}}))
}

func MinimumWeightCycleForEachVertex(n int, edges [][3]int) []int {
	graph := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		graph[i] = make([]int, n+1)
		for j := 1; j <= n; j++ {
			graph[i][j] = INF
		}
	}
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		u, v = u+1, v+1
		graph[u][v] = w
		graph[v][u] = w
	}
	visited := make([]bool, n+1)
	dist := make([]int, n+1)
	parent := make([]int, n+1)
	res := make([]int, 0, n)
	for u := 1; u <= n; u++ {
		for i := range visited {
			visited[i] = false
		}
		for i := 0; i <= n; i++ {
			dist[i] = INF
			parent[i] = i
		}
		dist[u] = 0
		parent[u] = u
		cand := INF
		for {
			cur := 0
			for i := 1; i <= n; i++ {
				if !visited[i] && dist[i] < dist[cur] {
					cur = i
				}
			}
			if cur == 0 {
				break
			}
			if visited[cur] {
				break
			}
			visited[cur] = true
			if parent[cur] != cur {
				cand = min(cand, dist[cur]+graph[cur][u])
			}
			for i := 1; i <= n; i++ {
				if i != u {
					if visited[i] {
						if parent[i] != parent[cur] {
							cand = min(cand, dist[cur]+dist[i]+graph[cur][i])
						}
					} else if dist[cur]+graph[cur][i] < dist[i] {
						dist[i] = dist[cur] + graph[cur][i]
						if cur == u {
							parent[i] = i
						} else {
							parent[i] = parent[cur]
						}
					}
				}
			}
		}
		res = append(res, cand)
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
