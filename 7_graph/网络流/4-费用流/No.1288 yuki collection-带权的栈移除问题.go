// No.1288 yuki collection-带权的消消乐问题
// 给定一个由yuki组成的长为n的字符串(n<=2000)
// 消除每个字母可以获得分数为s[i]
// 你可以消除连续的yuki
// !求你可以获得的最大分数
// https://maspypy.github.io/library/test/yukicoder/1288.test.cpp
// O(n^2logn)

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

// 拆点
func removeYuki(s string, scores []int) int {
	n := len(s)
	A := func(i int) int { return i }
	B := func(i int) int { return n + 1 + i }
	C := func(i int) int { return 2*(n+1) + i }
	D := func(i int) int { return 3*(n+1) + i }
	E := func(i int) int { return 4*(n+1) + i }
	mcmf := NewMinCostFlow(5*n+5, A(0), E(n))
	for i := 0; i < n; i++ {
		mcmf.AddEdge(A(i), A(i+1), n, 0)
		mcmf.AddEdge(B(i), B(i+1), n, 0)
		mcmf.AddEdge(C(i), C(i+1), n, 0)
		mcmf.AddEdge(D(i), D(i+1), n, 0)
		mcmf.AddEdge(E(i), E(i+1), n, 0)
	}
	for i := 0; i < n; i++ {
		if s[i] == 'y' {
			mcmf.AddEdge(A(i), B(i+1), 1, -scores[i])
		} else if s[i] == 'u' {
			mcmf.AddEdge(B(i), C(i+1), 1, -scores[i])
		} else if s[i] == 'k' {
			mcmf.AddEdge(C(i), D(i+1), 1, -scores[i])
		} else if s[i] == 'i' {
			mcmf.AddEdge(D(i), E(i+1), 1, -scores[i])
		}
	}

	// 流量最大时费用不一定最大,所以要根据slope判断
	slope := mcmf.Slope()
	res := 0
	for _, s := range slope {
		res = max(res, -s[1])
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type MinCostFlow struct {
	AddEdge func(from, to, cap, cost int)
	Flow    func() (maxFlow int, minCost int)
	Slope   func() [][2]int // (flow, cost)
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

	// ek
	Work := func() (maxFlow int, minCost int) {
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

	Slope := func() (slope [][2]int) {
		maxFlow, minCost := 0, 0
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
			slope = append(slope, [2]int{maxFlow, minCost})
		}
		return
	}

	AddEdge := func(from, to, cap, cost int) {
		addEdge(from, to, cap, cost, ei)
		ei++
	}

	return &MinCostFlow{
		AddEdge: AddEdge,
		Flow:    Work,
		Slope:   Slope,
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &scores[i])
	}
	fmt.Fprintln(out, removeYuki(s, scores))
}
