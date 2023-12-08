// 无向图经过每条边的最小环 O(n^3)
// 如果不存在环，返回INF.
// Minimum Weight Cycle For Each Edge

package main

import "fmt"

const INF int = 1e18

func main() {
	fmt.Println(MinimumWeightCycleForEachEdge(4, [][3]int{{0, 1, 1}, {1, 2, 2}, {2, 3, 3}, {3, 0, 4}, {0, 2, 3}}))
}

func MinimumWeightCycleForEachEdge(n int, edges [][3]int) []int {
	graph := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		graph[i] = make([]int, n+1)
		for j := 1; j <= n; j++ {
			graph[i][j] = INF
		}
	}
	edgeId := make([][]int, n+1)
	for i := range edgeId {
		edgeId[i] = make([]int, n+1)
	}
	edgeWeight := make([]int, len(edges)+1)
	res := make([]int, len(edges)+1)
	for i := 1; i <= len(edges); i++ {
		e := edges[i-1]
		u, v, w := e[0], e[1], e[2]
		u, v = u+1, v+1
		graph[u][v] = w
		graph[v][u] = w
		edgeId[u][v] = i
		edgeId[v][u] = i
		res[i] = INF
		edgeWeight[i] = w
	}

	visited := make([]bool, n+1)
	dist := make([]int, n+1)
	dp := make([]int, len(edges)+1)
	parent := make([]int, n+1)
	last := make([]int, n+1)
	for u := 1; u <= n; u++ {
		for i := range visited {
			visited[i] = false
		}
		for i := 0; i <= n; i++ {
			dist[i] = INF
			parent[i] = i
			last[i] = -1
		}
		for i := 1; i <= len(edges); i++ {
			dp[i] = INF
		}
		dist[u] = 0
		parent[u] = u
		last[u] = 0
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
				dp[edgeId[last[cur]][cur]] = min(dp[edgeId[last[cur]][cur]], dist[cur]+graph[cur][u])
				res[edgeId[cur][u]] = min(res[edgeId[cur][u]], dist[cur]+graph[cur][u])
				res[edgeId[u][cur]] = min(res[edgeId[u][cur]], dist[cur]+graph[cur][u])
			}
			for i := 1; i <= n; i++ {
				if i != u {
					if visited[i] {
						if parent[i] != parent[cur] {
							dp[edgeId[last[cur]][cur]] = min(dp[edgeId[last[cur]][cur]], dist[cur]+dist[i]+graph[cur][i])
							dp[edgeId[last[i]][i]] = min(dp[edgeId[last[i]][i]], dist[cur]+dist[i]+graph[cur][i])
							res[edgeId[cur][i]] = min(res[edgeId[cur][i]], dist[cur]+dist[i]+graph[cur][i])
							res[edgeId[i][cur]] = min(res[edgeId[i][cur]], dist[cur]+dist[i]+graph[cur][i])
						}
					} else if dist[cur]+graph[cur][i] < dist[i] {
						dist[i] = dist[cur] + graph[cur][i]
						last[i] = cur
						if cur == u {
							parent[i] = i
						} else {
							parent[i] = parent[cur]
						}
					}
				}
			}

		}

		for i := 1; i <= n; i++ {
			if last[i] == -1 {
				continue
			}
			x := dp[edgeId[last[i]][i]]
			for cur := i; cur != 0; cur = last[cur] {
				if last[cur] != 0 {
					dp[edgeId[last[cur]][cur]] = min(dp[edgeId[last[cur]][cur]], x)
				}
			}
		}

		for i := 1; i <= len(edges); i++ {
			res[i] = min(res[i], dp[i])
		}

	}

	return res[1:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
