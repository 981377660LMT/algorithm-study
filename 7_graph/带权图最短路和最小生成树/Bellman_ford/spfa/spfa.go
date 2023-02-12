// 带负边的单源最短路
// v<=1000 e<=2000

package main

import (
	"bufio"
	"fmt"
	"os"
)

type WeightedEdge struct{ to, weighgt int }

const INF int = 1e18

// spfa 求带有负权边的最短路，时间复杂度 O(V*E)，有负环时返回 nil
//   !只是找负环的话，初始时将所有点入队即可
func spfa(n int, graph [][]WeightedEdge, start int) (dist []int) {
	dist = make([]int, n)
	queue := []int{}
	relaxedConut := make([]int, n)
	inQueue := make([]bool, n)
	for i := range dist {
		dist[i] = INF
	}

	queue = append(queue, start)
	inQueue[start] = true
	relaxedConut[start] = 1
	dist[start] = 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		inQueue[cur] = false
		for _, e := range graph[cur] {
			next, weight := e.to, e.weighgt
			cand := dist[cur] + weight
			if cand >= dist[next] {
				continue
			}
			dist[next] = cand
			if !inQueue[next] {
				relaxedConut[next]++
				// 找到一个从 start 出发可达的负环
				if relaxedConut[next] >= n {
					return nil
				}
				inQueue[next] = true
				queue = append(queue, next)
			}
		}
	}

	return
}

func main() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_1_B&lang=jp

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, s int
	fmt.Fscan(in, &n, &m, &s)
	graph := make([][]WeightedEdge, n)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		graph[u] = append(graph[u], WeightedEdge{v, w})
	}

	dist := spfa(n, graph, s)
	if dist == nil {
		fmt.Fprintln(out, "NEGATIVE CYCLE")
		return
	}

	for _, d := range dist {
		if d == INF {
			fmt.Fprintln(out, "INF")
		} else {
			fmt.Fprintln(out, d)
		}
	}
}
