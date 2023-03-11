// 最小斯坦纳树
// !一个带权的无向图上有k个关键点，求联通k个关键点最小的代价(边权之和)。
// n≤100,m≤500,k≤10。

// O(n*3^k+mlogm*2^k)

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

// 一个联通的无向带权图上有k个关键点 criticals，求联通k个关键点最小的代价(边权之和)。
func MinimumSteinerTree(n int, edges [][]int, criticals []int) int {
	k := len(criticals)
	visited := make([]bool, n)
	graph := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], [2]int{v, w})
		graph[v] = append(graph[v], [2]int{u, w})
	}

	dp := make([][]int, 1<<k)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}

	for i := 0; i < k; i++ {
		dp[1<<i][criticals[i]] = 0
	}

	spfa := func(s int) {
		q := []int{}
		for i := 0; i < n; i++ {
			if dp[s][i] != INF {
				q = append(q, i)
				visited[i] = true
			}
		}

		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			visited[u] = false
			for _, e := range graph[u] {
				v, w := e[0], e[1]
				if dp[s][u]+w < dp[s][v] {
					dp[s][v] = dp[s][u] + w
					if !visited[v] {
						q = append(q, v)
						visited[v] = true
					}
				}
			}
		}
	}

	for s := 1; s < (1 << k); s++ {
		for t := s & (s - 1); t > 0; t = (t - 1) & s {
			if t < (s ^ t) {
				break
			}
			for i := 0; i < n; i++ {
				dp[s][i] = min(dp[s][i], dp[t][i]+dp[t^s][i])
			}
		}
		spfa(s)
	}

	res := INF
	for i := 0; i < n; i++ {
		res = min(res, dp[(1<<k)-1][i])
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// https://www.luogu.com.cn/problem/P6192
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v, w})
	}

	criticals := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &criticals[i])
		criticals[i]--
	}

	fmt.Fprintln(out, MinimumSteinerTree(n, edges, criticals))
}
