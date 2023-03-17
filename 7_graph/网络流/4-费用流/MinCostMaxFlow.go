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

func main() {
	// https://www.acwing.com/problem/content/description/2195/
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	START, END, OFFSET := 2*n+2, 2*n+3, n
	mcmf1 := NewMinCostMaxFlow(2*n+5, START, END)
	mcmf2 := NewMinCostMaxFlow(2*n+5, START, END)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var cost int
			fmt.Fscan(in, &cost)
			mcmf1.AddEdge(i, j+OFFSET, 1, cost)
			mcmf2.AddEdge(i, j+OFFSET, 1, -cost)
		}
	}

	for i := 0; i < n; i++ {
		mcmf1.AddEdge(START, i, 1, 0)
		mcmf2.AddEdge(START, i, 1, 0)
		mcmf1.AddEdge(i+OFFSET, END, 1, 0)
		mcmf2.AddEdge(i+OFFSET, END, 1, 0)
	}

	_, cost1 := mcmf1.Flow()
	_, cost2 := mcmf2.Flow()
	fmt.Fprintln(out, cost1)
	fmt.Fprintln(out, -cost2)
}

const INF int = 1e18

type Edge struct{ from, to, cap, flow, cost, id int }
type MinCostMaxFlow struct {
	AddEdge        func(from, to, cap, cost int)
	Flow           func() (maxFlow int, minCost int)
	FlowWithLimit  func(flowLimit int) (maxFlow int, minCost int)
	Slope          func() [][2]int // (flow, cost) 的每个转折点
	SlopeWithLimit func(flowLimit int) [][2]int
	Edges          func() []Edge // 注意根据from,to排除虚拟源点汇点; `flow>0` 才是流经的边
}

func NewMinCostMaxFlow(n, start, end int) *MinCostMaxFlow {
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

	// ek
	FlowWithLimit := func(flowLimit int) (maxFlow int, minCost int) {
		for maxFlow < flowLimit {
			if !spfa() {
				break
			}

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

	SlopeWithLimit := func(flowLimit int) (slope [][2]int) {
		maxFlow, minCost := 0, 0
		for maxFlow < flowLimit {
			if !spfa() {
				break
			}
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
			slope = append(slope, [2]int{maxFlow, minCost})
		}
		return
	}

	AddEdge := func(from, to, cap, cost int) {
		addEdge(from, to, cap, cost, ei)
		ei++
	}

	GetEdges := func() (res []Edge) {
		for from, edges := range graph {
			for _, e := range edges {
				if e.eid == -1 {
					continue
				}
				res = append(res, Edge{from, e.to, e.cap + graph[e.to][e.rid].cap, graph[e.to][e.rid].cap, e.cost, e.eid})
			}
		}
		return
	}

	return &MinCostMaxFlow{
		AddEdge:        AddEdge,
		Flow:           func() (maxFlow int, minCost int) { return FlowWithLimit(INF) },
		FlowWithLimit:  FlowWithLimit,
		Slope:          func() [][2]int { return SlopeWithLimit(INF) },
		SlopeWithLimit: SlopeWithLimit,
		Edges:          GetEdges,
	}
}
