package main

const INF int = 1e18

// Floyd 求多源最短路.
func Floyd(n int, edges [][]int) [][]int {
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = INF
		}
	}
	for i := 0; i < n; i++ {
		dist[i][i] = 0
	}

	for _, road := range edges {
		u, v, w := road[0], road[1], road[2]
		dist[u][v] = min(w, dist[u][v]) // 有重边，取最小值
		dist[v][u] = min(w, dist[v][u])
	}

	// dis[k][i][j] 表示「经过若干个编号不超过 k 的节点」时，从 i 到 j 的最短路长度
	for k := 0; k < n; k++ { // 经过的中转点
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				// 松弛：如果一条边可以被松弛了，说明这条边就没有必要留下了
				cand := dist[i][k] + dist[k][j]
				dist[i][j] = min(cand, dist[i][j])
			}
		}
	}
	return dist
}

type FloydStatic struct {
}

// 动态Floyd算法, 支持拷贝、增加删除点、增加删除边.
type FloydDynamic struct {
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
