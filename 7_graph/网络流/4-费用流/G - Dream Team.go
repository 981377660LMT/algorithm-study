// https://zhuanlan.zhihu.com/p/496282947
// 题意
// 有n个人来自不同的学校ai,擅长不同的学科bi,每个人有一个能力值ci
// 要求组建一支i个人的梦之队最大化队员的能力值,并且满足队伍中所有人来自的学校和擅长的学科都不同.
// 输出匹配数为1,2,…k个匹配时的最优匹配.
// n<=3e4
// !ai,bi<=150 (暗示作为顶点的数据量)
// ci<=1e9

// 分析
// 把学校和学科看作点,把一个人看成匹配边,能力值看作边权,
// 其实就是求有i条匹配边的最优匹配.可以用费用流解决.
// !对每个学生,从ai到bi连一条容量为1,费用为ci的边.
// 在spfa费用流算法中一次spfa只会找到一条费用最小的增广流,
// 所以每次增广之后得到的费用就对应匹配数为1,2,…k个匹配时的答案.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	schools := make([]int, n)
	types := make([]int, n)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &schools[i], &types[i], &scores[i])
	}

	res := dreamTeam(schools, types, scores)
	fmt.Fprintln(out, len(res))
	for _, r := range res {
		fmt.Fprintln(out, r)
	}
}

func dreamTeam(schools []int, types []int, scores []int) (res []int) {
	n := len(schools)
	START, END, OFFSET := 320, 321, 160
	flow := NewMinCostMaxFlow(330, START, END)
	for i := 0; i < n; i++ {
		flow.AddEdge(schools[i], OFFSET+types[i], 1, -scores[i])
	}
	for i := 0; i < OFFSET; i++ {
		flow.AddEdge(START, i, 1, 0)
		flow.AddEdge(i+OFFSET, END, 1, 0)
	}

	curFlow := 0
	curCost := 0
	for {
		f, c := flow.FlowWithLimit(1) // !每次消耗1流量,每个回合的答案就是容量为1,2,...,时的最小费用
		if f == 0 {
			break
		}
		curFlow += 1
		curCost += -c
		res = append(res, curCost)
	}

	return
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
