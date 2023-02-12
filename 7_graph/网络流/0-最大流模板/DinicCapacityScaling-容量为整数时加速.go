// https://ei1333.github.io/library/graph/flow/dinic-capacity-scaling.hpp
// 如果边的容量为整数，可以使用容量缩放加速Dinic算法。
// O(EVlogU) // U为边的最大容量

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, s, t int
	fmt.Fscan(in, &n, &m, &s, &t)
	s, t = s-1, t-1
	dinic := NewDinicCapacityScaling(n)
	for i := 0; i < m; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		u, v = u-1, v-1
		dinic.AddEdge(u, v, c)
	}

	fmt.Fprintln(out, dinic.MaxFlow(s, t))
}

const INF int = 1e18

type edge struct {
	to    int
	cap   int
	rev   int
	isrev bool
	idx   int
}

type DinicCapacityScaling struct {
	graph   [][]edge
	minCost []int
	iter    []int
	maxCap  int
}

func NewDinicCapacityScaling(n int) *DinicCapacityScaling {
	return &DinicCapacityScaling{
		graph: make([][]edge, n),
	}
}

func (d *DinicCapacityScaling) AddEdge(from, to, cap int) {
	d.AddEdgeWithIndex(from, to, cap, -1)
}

func (d *DinicCapacityScaling) AddEdgeWithIndex(from, to, cap, index int) {
	d.maxCap = max(d.maxCap, cap)
	d.graph[from] = append(d.graph[from], edge{to, cap, len(d.graph[to]), false, index})
	d.graph[to] = append(d.graph[to], edge{from, 0, len(d.graph[from]) - 1, true, index})
}

func (d *DinicCapacityScaling) MaxFlow(s, t int) int {
	if d.maxCap == 0 {
		return 0
	}

	flow := 0
	tmp := 63 - bits.LeadingZeros(uint(d.maxCap))
	for i := tmp; i >= 0; i-- {
		now := 1 << i
		for d.buildAugmentPath(s, t, now) {
			d.iter = make([]int, len(d.graph))
			flow += d.findAugmentPath(s, t, now, INF)
		}
	}
	return flow
}

// (from, to, flow, capacity, eid)
func (d *DinicCapacityScaling) GetEdges() [][5]int {
	res := make([][5]int, 0)
	for i := range d.graph {
		for _, e := range d.graph[i] {
			if e.isrev {
				continue
			}
			revE := d.graph[e.to][e.rev]
			res = append(res, [5]int{i, e.to, revE.cap, e.cap + revE.cap, e.idx})
		}
	}
	return res
}

func (d *DinicCapacityScaling) buildAugmentPath(s, t, base int) bool {
	d.minCost = make([]int, len(d.graph))
	for i := range d.minCost {
		d.minCost[i] = -1
	}
	q := []int{s}
	d.minCost[s] = 0
	for len(q) > 0 && d.minCost[t] == -1 {
		p := q[0]
		q = q[1:]
		for _, e := range d.graph[p] {
			if e.cap >= base && d.minCost[e.to] == -1 {
				d.minCost[e.to] = d.minCost[p] + 1
				q = append(q, e.to)
			}
		}
	}
	return d.minCost[t] != -1
}

func (ds *DinicCapacityScaling) findAugmentPath(idx, t, base, flow int) int {
	if idx == t {
		return flow
	}
	sum := 0
	cur := ds.iter[idx]
	for cur < len(ds.graph[idx]) {
		e := &ds.graph[idx][cur]
		if e.cap >= base && ds.minCost[idx] < ds.minCost[e.to] {
			d := ds.findAugmentPath(e.to, t, base, min(flow-sum, e.cap))
			if d > 0 {
				e.cap -= d
				ds.graph[e.to][e.rev].cap += d
				sum += d
				if flow-sum < base {
					break
				}
			}
		}
		cur++
		ds.iter[idx]++
	}

	return sum
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
