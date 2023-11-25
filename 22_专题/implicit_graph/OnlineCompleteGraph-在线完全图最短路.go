// OnlineCompleteGraph-在线完全图补图bfs
// https://maspypy.github.io/library/graph/implicit_graph/complement_graph_bfs.hpp

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18
const MOD int = 998244353

func demo() {
	fmt.Println(ComplementGraphBfs(5, 0, [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}}))
}

// G - Counting Shortest Paths
// https://atcoder.jp/contests/abc319/tasks/abc319_g
// 在一个n个顶点的无向无权完全图中删除m条边,求从START到TARGET的最路径数模998244353.
//
// bfs最短路计数分解成两个问题:
//  1. 在线bfs求出START到其他点的最短路;
//  2. START到TARGET的路径上所有的边 u->v 满足 dist[u]+1 == dist[v].(最短路的充要条件)
//     !可以按照距离从小到大dp.用总数减去不合法的路径数即可.
//
// 参考:
// E - Safety Journey
// https://atcoder.jp/contests/abc212/tasks/abc212_e
// 转移边数很多的dp问题 => 正难则反.
func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	banEdges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &banEdges[i][0], &banEdges[i][1])
		banEdges[i][0]--
		banEdges[i][1]--
	}
	START, TARGET := 0, n-1

	dist, _ := ComplementGraphBfs(n, START, banEdges)
	if dist[TARGET] == INF {
		fmt.Fprintln(out, -1)
		return
	}

	ban := make([][]int, n) // 每个点的禁止转移的邻居
	for _, e := range banEdges {
		u, v := e[0], e[1]
		ban[u] = append(ban[u], v)
		ban[v] = append(ban[v], u)
	}

	groupByDist := make([][]int, n) // 距离START为0,1,2,...,n-1的点
	for i := 0; i < n; i++ {
		if d := dist[i]; d != INF {
			groupByDist[d] = append(groupByDist[d], i)
		}
	}

	dp := make([]int, n) // 到达i点的路径数
	dp[START] = 1
	for d := 0; d < n-1; d++ {
		preCount := 0
		for _, pre := range groupByDist[d] {
			preCount += dp[pre]
			preCount %= MOD
		}
		for _, cur := range groupByDist[d+1] {
			dp[cur] = preCount
			for _, pre := range ban[cur] {
				if dist[pre]+1 == dist[cur] {
					dp[cur] -= dp[pre]
					dp[cur] %= MOD
					if dp[cur] < 0 {
						dp[cur] += MOD
					}
				}
			}
		}
	}

	fmt.Fprintln(out, dp[TARGET])
}

// 完全图最短路.
//
//	给定一个无向无权的完全图，求出完全图上从start到其他点的最短路.不可达的点距离为INF.
//	banEdges是禁止通行的边.
//	O(V+len(banEdges)).
func ComplementGraphBfs(n int, start int, banEdges [][2]int) (dist []int, pre []int) {
	banGraph := make([][]int, n)
	for _, e := range banEdges {
		u, v := e[0], e[1]
		banGraph[u] = append(banGraph[u], v)
		banGraph[v] = append(banGraph[v], u)
	}

	dist = make([]int, n)
	pre = make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
		pre[i] = -1
	}
	dist[start] = 0
	queue := make([]int, 0, n)
	queue = append(queue, start)

	notNeightBor := make([]bool, n)
	unVisited := make([]int, 0, n-1)
	for i := 0; i < n; i++ {
		if i != start {
			unVisited = append(unVisited, i)
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, u := range banGraph[cur] {
			notNeightBor[u] = true
		}

		nextUnvisited := []int{}
		for _, next := range unVisited {
			if notNeightBor[next] {
				nextUnvisited = append(nextUnvisited, next) // findUnvisited
			} else {
				// setVisited
				dist[next] = dist[cur] + 1
				pre[next] = cur
				queue = append(queue, next)
			}
		}
		unVisited = nextUnvisited

		for _, u := range banGraph[cur] {
			notNeightBor[u] = false
		}
	}

	return
}
