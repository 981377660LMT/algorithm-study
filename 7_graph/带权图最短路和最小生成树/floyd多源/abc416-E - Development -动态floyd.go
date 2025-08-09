// https://atcoder.jp/contests/abc416/tasks/abc416_e
// APSP 动态插边技巧
// 多源最短路加边加点
//
// 在 AtCoder 国有编号 1…N 的 N 个城市、M 条道路、K 个机场。
// 道路双向，连接 A_i 与 B_i，耗时 C_i；机场位于城市 D_j 之间可双向以时间 T 互通。
// 支持三类操作：
// 1 x y t：新增道路 x–y，时间 t
// 2 x ：在城市 x 建机场
// 3 ：计算当前所有可达对 x,y 的最短时间之和（不可达计 0），输出结果
//
// !N<=500, M<=1e5, Q<=1000
//
// 把机场网络抽象成一个“超级节点” S：
// • 对每个机场城市 v 建两条有向边：v→S 权 0，S→v 权 T。
// • 任意两机场 p→q 的最短路即 p→S→q，耗时 0+T=T。
// !这样一次“建机场”仅插入两条有向边，而不是 O( |机场| ) 条。

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 4e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M int
	fmt.Fscan(in, &N, &M)

	dummy := N // 超级节点
	dist := make([][]int, N+1)
	for i := range dist {
		dist[i] = make([]int, N+1)
		for j := range dist[i] {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = INF
			}
		}
	}

	for i := 0; i < M; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		a--
		b--
		c *= 2
		dist[a][b] = min(c, dist[a][b])
		dist[b][a] = min(c, dist[b][a])
	}

	var K, T int
	fmt.Fscan(in, &K, &T)
	for i := 0; i < K; i++ {
		var x int
		fmt.Fscan(in, &x)
		x--
		dist[x][dummy] = min(T, dist[x][dummy])
		dist[dummy][x] = min(T, dist[dummy][x])
	}

	for k := 0; k < N+1; k++ {
		for i := 0; i < N+1; i++ {
			if dist[i][k] == INF {
				continue
			}
			dik := dist[i][k]
			for j := 0; j < N+1; j++ {
				djk := dist[k][j]
				if djk == INF {
					continue
				}
				d := dik + djk
				if d < dist[i][j] {
					dist[i][j] = d
				}
			}
		}
	}

	updateEdge := func(u, v, w int) {
		if dist[u][v] <= w {
			return
		}

		dist[u][v] = w
		for i := 0; i < N+1; i++ {
			for j := 0; j < N+1; j++ {
				if dist[i][u] == INF || dist[v][j] == INF {
					continue
				}
				cand := dist[i][u] + w + dist[v][j]
				if dist[i][j] > cand {
					dist[i][j] = cand
				}
			}
		}

		dist[v][u] = w
		for i := 0; i < N+1; i++ {
			for j := 0; j < N+1; j++ {
				if dist[i][v] == INF || dist[u][j] == INF {
					continue
				}
				cand := dist[i][v] + w + dist[u][j]
				if dist[i][j] > cand {
					dist[i][j] = cand
				}
			}
		}
	}

	var Q int
	fmt.Fscan(in, &Q)
	for qi := 0; qi < Q; qi++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var x, y, c int
			fmt.Fscan(in, &x, &y, &c)
			x--
			y--
			c *= 2
			updateEdge(x, y, c)
		} else if t == 2 {
			var x int
			fmt.Fscan(in, &x)
			x--
			updateEdge(x, dummy, T)
		} else if t == 3 {
			res := 0
			for i := 0; i < N; i++ {
				for j := i + 1; j < N; j++ {
					if dist[i][j] < INF {
						res += dist[i][j]
					}
				}
			}
			fmt.Fprintln(out, res)
		}
	}
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
