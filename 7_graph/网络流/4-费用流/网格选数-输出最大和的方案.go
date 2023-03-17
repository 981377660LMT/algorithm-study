// https://atcoder.jp/contests/practice2/tasks/practice2_e
// 网格选数-输出最大和的方案
// 给定一个全是非负数的网格,从中选取一些数使得和最大.
// 注意任意一行不能超过k个数, 任意一列不能超过k个数.
// 选择用'X'表示，不选用'.'表示

package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(grid [][]int, rowLimit, colLimit int) (maxSum int, selected [][2]int) {
	ROW, COL := len(grid), len(grid[0])
	START, END := ROW+COL, ROW+COL+1

	getFlow := func() *MinCostMaxFlow {
		flow := NewMinCostMaxFlow(ROW+COL+2, START, END)
		for i := 0; i < ROW; i++ {
			flow.AddEdge(START, i, rowLimit, 0)
		}
		for j := 0; j < COL; j++ {
			flow.AddEdge(j+ROW, END, colLimit, 0)
		}
		for i := 0; i < ROW; i++ {
			for j := 0; j < COL; j++ {
				flow.AddEdge(i, j+ROW, 1, -grid[i][j]) // max value
			}
		}
		return flow
	}

	// !流量最大时费用不一定最大
	// !所以要根据slope判断,求出此时的最佳流量限制
	flow1 := getFlow()
	slope := flow1.Slope()
	bestFlowLimit := INF
	for _, s := range slope {
		f, c := s[0], -s[1]
		if c > maxSum {
			maxSum = c
			bestFlowLimit = f
		}
	}

	// 再求一次最大流,流量限制为bestFlowLimit
	flow2 := getFlow()
	flow2.FlowWithLimit(bestFlowLimit)
	edges := flow2.Edges()
	for _, e := range edges {
		if e.from == START || e.to == END || e.flow <= 0 { // invalid flow edge
			continue
		}
		from, to := e.from, e.to-ROW
		selected = append(selected, [2]int{from, to})
	}

	return
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, n)
		for j := range grid[i] {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	solution := make([][]byte, n)
	for i := range solution {
		solution[i] = make([]byte, n)
		for j := range solution[i] {
			solution[i][j] = '.'
		}
	}

	res, selected := solve(grid, k, k)
	for _, p := range selected {
		solution[p[0]][p[1]] = 'X'
	}

	fmt.Fprintln(out, res)
	for _, row := range solution {
		fmt.Fprintln(out, string(row))
	}
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
