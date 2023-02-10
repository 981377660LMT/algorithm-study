// 全局最小割
// Stoer-Wagner 算法
// http://blog.kongfy.com/2015/02/kargermincut/
// https://www.luogu.com.cn/problem/P5632
// 全局最小割不指定源点和汇点，而是要求将图中的所有点分成两部分(两个连通分量)
// 且删除的边权之和最小
// 如果图不连通,返回0
// n<=500

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func globalMinCut(n int, edges [][]int) int {
	s, t := 0, 0
	dap := make([]bool, n+1)
	ord := make([]int, n+1)
	var visited []bool
	var w []int
	dist := make([][]int, n+1)
	for i := 0; i < n+1; i++ {
		dist[i] = make([]int, n+1)
	}

	for _, edge := range edges {
		u, v, w := edge[0], edge[1], edge[2]
		u, v = u+1, v+1
		dist[u][v] += w
		dist[v][u] += w
	}

	proc := func(x int) int {
		visited = make([]bool, n+1)
		w = make([]int, n+1)
		w[0] = -1
		for i := 1; i <= n-x+1; i++ {
			mx := 0
			for j := 1; j <= n; j++ {
				if !dap[j] && !visited[j] && w[j] > w[mx] {
					mx = j
				}
			}

			visited[mx] = true
			ord[i] = mx
			for j := 1; j <= n; j++ {
				if !dap[j] && !visited[j] {
					w[j] += dist[mx][j]
				}
			}
		}

		s = ord[n-x]
		t = ord[n-x+1]
		return w[t]
	}

	sw := func() int {
		res := INF
		for i := 1; i < n; i++ {
			res = min(res, proc(i))
			dap[t] = true
			for j := 1; j <= n; j++ {
				dist[s][j] += dist[t][j]
				dist[j][s] += dist[j][t]
			}
		}

		return res
	}

	return sw()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	const INF int = int(1e18)
	const MOD int = 998244353

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v, w})
	}

	fmt.Fprintln(out, globalMinCut(n, edges))
}
