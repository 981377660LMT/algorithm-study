// 最小费用流 MCFP
// 最小费用最大流 MCMF（即满流时的费用）
// 将 Edmonds-Karp 中的 BFS 改成 SPFA O(f*VE) 或 Dijkstra O(f*ElogV)
// 要求初始网络中无负权圈
// 性能对比（洛谷 P3381，由于数据不强所以 SPFA 很快）：SPFA 1.05s(max 365ms)   Dijkstra 1.91s(max 688ms)
// https://oi-wiki.org/graph/flow/min-cost/
// !https://github.dev/EndlessCheng/codeforces-go/blob/master/copypasta/graph.go#L3586

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, s, t int
	fmt.Fscan(in, &n, &m, &s, &t)
	s--
	t--
	mcf := NewMinCostFlow(n, s, t)
	for i := 0; i < m; i++ {
		var u, v, c, w int
		fmt.Fscan(in, &u, &v, &c, &w)
		u--
		v--
		mcf.AddEdge(u, v, c, w)
	}
	maxFlow, minCost := mcf.Work()
	fmt.Fprintln(out, maxFlow, minCost)
}

type MinCostFlow struct {
	AddEdge func(from, to, cap, cost int)
	Work    func() (maxFlow int, minCost int)
}

func NewMinCostFlow(n, start, end int) *MinCostFlow {
	type neighbor struct {
		to   int
		rid  int // rid 为反向边在邻接表中的下标
		cap  int // 边的残量
		cost int
		eid  int // -1表示是反向边
	}

	graph := make([][]neighbor, n)
	ei := 0
	addEdge := func(from, to, cap, cost, eid int) {
		graph[from] = append(graph[from], neighbor{to, len(graph[to]), cap, cost, eid})
		graph[to] = append(graph[to], neighbor{from, len(graph[from]) - 1, 0, -cost, -1})
	}

	dist := make([]int, len(graph))
	type vi struct{ v, i int }
	pre := make([]vi, len(graph))
	spfa := func() bool {
		for i := range dist {
			dist[i] = INF
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
				if cand := dist[cur] + int(edge.cost); cand < dist[next] {
					dist[next] = cand
					pre[next] = vi{cur, i}
					if !inQueue[next] {
						queue = append(queue, next)
						inQueue[next] = true
					}
				}
			}
		}
		return dist[end] < INF
	}

	ek := func() (maxFlow int, minCost int) {
		for spfa() {
			// 沿 st-end 的最短路尽量增广
			flow := INF
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
			minCost += dist[end] * flow
		}
		return
	}

	AddEdge := func(from, to, cap, cost int) {
		addEdge(from, to, cap, cost, ei)
		ei++
	}

	return &MinCostFlow{
		AddEdge: AddEdge,
		Work:    ek,
	}
}
