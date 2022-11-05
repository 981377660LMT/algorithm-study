package mincostmaxflow

import (
	"fmt"
	"io"
)

// 最小费用流 MCFP
// 最小费用最大流 MCMF（即满流时的费用）
// 将 Edmonds-Karp 中的 BFS 改成 SPFA O(f*VE) 或 Dijkstra O(f*ElogV)
// 要求初始网络中无负权圈
// 性能对比（洛谷 P3381，由于数据不强所以 SPFA 很快）：SPFA 1.05s(max 365ms)   Dijkstra 1.91s(max 688ms)
// https://oi-wiki.org/graph/flow/min-cost/
// !https://github.dev/EndlessCheng/codeforces-go/blob/master/copypasta/graph.go#L3586
func MinCostFlowSPFA(in io.Reader, n, m, start, end int) (int, int64) {
	const inf int = 1e18
	start--
	end--

	type neighbor struct {
		to int
		// rid 为反向边在邻接表中的下标
		rid int
		// 边的残量
		cap  int
		cost int
		// -1表示是反向边
		eid int
	}
	graph := make([][]neighbor, n)
	addEdge := func(from, to, cap, cost, eid int) {
		graph[from] = append(graph[from], neighbor{to, len(graph[to]), cap, cost, eid})
		graph[to] = append(graph[to], neighbor{from, len(graph[from]) - 1, 0, -cost, -1}) // !无向图上 0 换成 cap
	}

	for i := 0; i < m; i++ {
		var v, w, cap, cost int
		fmt.Fscan(in, &v, &w, &cap, &cost)
		v--
		w--
		addEdge(v, w, cap, cost, i)
	}

	dist := make([]int64, len(graph))
	type vi struct{ v, i int }
	pre := make([]vi, len(graph))
	spfa := func() bool {
		const _inf int64 = 1e18
		for i := range dist {
			dist[i] = _inf
		}
		dist[start] = 0
		inQueue := make([]bool, len(graph))
		inQueue[start] = true
		queue := []int{start}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			inQueue[cur] = false
			for i, edge := range graph[cur] {
				if edge.cap == 0 {
					continue
				}
				next := edge.to
				if cand := dist[cur] + int64(edge.cost); cand < dist[next] {
					dist[next] = cand
					pre[next] = vi{cur, i}
					if !inQueue[next] {
						queue = append(queue, next)
						inQueue[next] = true
					}
				}
			}
		}
		return dist[end] < _inf
	}

	ek := func() (maxFlow int, minCost int64) {
		for spfa() {
			// 沿 st-end 的最短路尽量增广
			flow := inf
			for cur := end; cur != start; {
				p := pre[cur]
				if c := graph[p.v][p.i].cap; c < flow {
					flow = c
				}
				cur = p.v
			}
			for cur := end; cur != start; {
				p := pre[cur]
				edge := &graph[p.v][p.i]
				edge.cap -= flow
				graph[cur][edge.rid].cap += flow
				cur = p.v
			}
			maxFlow += flow
			minCost += dist[end] * int64(flow)
		}
		return
	}

	return ek()
}
