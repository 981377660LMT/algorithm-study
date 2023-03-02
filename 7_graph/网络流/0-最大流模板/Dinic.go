// !最大流 Dinic's algorithm O(V^2 * E)  二分图上为 O(E√V)
// 如果容量是浮点数，下面代码中 > 0 的判断要改成 > eps
// https://ei1333.github.io/library/graph/flow/dinic.hpp
// https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph.go#L3207

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

	var n, m int
	fmt.Fscan(in, &n, &m)
	dinic := NewDinic(n)
	for i := 0; i < m; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		dinic.AddEdge(u, v, c)
	}

	fmt.Fprintln(out, dinic.MaxFlow(0, n-1))
}

const INF int = 1e18

type DinicEdge struct {
	to    int
	cap   int // 剩余容量
	rev   int // 逆向边在 G[to] 中的序号
	isRev bool
	index int
}

type Dinic struct {
	graph   [][]DinicEdge
	minCost []int
	iter    []int
}

func NewDinic(n int) *Dinic {
	return &Dinic{
		graph: make([][]DinicEdge, n),
	}
}

func (d *Dinic) AddEdge(from, to, cap int) {
	d.AddEdgeWithIndex(from, to, cap, -1)
}

func (d *Dinic) AddEdgeWithIndex(from, to, cap, index int) {
	d.graph[from] = append(d.graph[from], DinicEdge{to, cap, len(d.graph[to]), false, index})
	d.graph[to] = append(d.graph[to], DinicEdge{from, 0, len(d.graph[from]) - 1, true, index})
}

func (d *Dinic) MaxFlow(s, t int) int {
	flow := 0
	for d.buildAugmentingPath(s, t) {
		d.iter = make([]int, len(d.graph))
		f := 0
		for {
			f = d.findMinDistAugmentPath(s, t, INF)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

// (from,to,流量,容量)
func (d *Dinic) GetEdges() [][4]int {
	res := make([][4]int, 0)
	for i, edges := range d.graph {
		for _, e := range edges {
			if e.isRev {
				continue
			}
			revEdge := d.graph[e.to][e.rev]
			res = append(res, [4]int{i, e.to, revEdge.cap, e.cap + revEdge.cap})
		}
	}
	return res
}

func (d *Dinic) findMinDistAugmentPath(idx, t, flow int) int {
	if idx == t {
		return flow
	}

	i := d.iter[idx]
	for i < len(d.graph[idx]) {
		e := d.graph[idx][i]
		if e.cap > 0 && d.minCost[idx] < d.minCost[e.to] {
			f := d.findMinDistAugmentPath(e.to, t, min(flow, e.cap))
			if f > 0 {
				d.graph[idx][i].cap -= f
				d.graph[e.to][e.rev].cap += f
				return f
			}
		}
		i++
		d.iter[idx]++
	}
	return 0
}

func (d *Dinic) buildAugmentingPath(s, t int) bool {
	d.minCost = make([]int, len(d.graph))
	for i := range d.minCost {
		d.minCost[i] = -1
	}
	d.minCost[s] = 0
	queue := []int{s}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, e := range d.graph[v] {
			if e.cap > 0 && d.minCost[e.to] == -1 {
				d.minCost[e.to] = d.minCost[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return d.minCost[t] != -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
