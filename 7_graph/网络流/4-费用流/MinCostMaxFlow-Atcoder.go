// !如果带负权边,需要用Big添加位移或者spfa版的MinCostMaxFlow

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
	mcmf1 := NewMinCostMaxFlow(2*n + 5)
	mcmf2 := NewMinCostMaxFlow(2*n + 5)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var cost int
			fmt.Fscan(in, &cost)
			mcmf1.AddEdge(i, j+OFFSET, 1, cost)
			mcmf2.AddEdge(i, j+OFFSET, 1, BIG-cost) // !BIG 消除负权边
		}
	}

	for i := 0; i < n; i++ {
		mcmf1.AddEdge(START, i, 1, 0)
		mcmf2.AddEdge(START, i, 1, 0)
		mcmf1.AddEdge(i+OFFSET, END, 1, 0)
		mcmf2.AddEdge(i+OFFSET, END, 1, 0)
	}

	_, cost1 := mcmf1.Flow(START, END)
	flow, cost2 := mcmf2.Flow(START, END)
	cost2 = BIG*flow - cost2 // !最后把多加的费用补回来

	fmt.Fprintln(out, cost1)
	fmt.Fprintln(out, cost2)
}

const (
	// Remove minus cost when solving max cost max flow problem.
	// !BIG needs to be larger than the maximum cost.
	BIG int = 1e9
	INF int = 1e18
)

type Edge struct {
	From, To        int
	Cap, FLow, Cost int
}

type MinCostMaxFlow struct {
	n   int
	pos [][2]int
	g   [][]*neighbor
}

// NewMinCostMaxFlow creates a graph of n vertices and 0 edges.
func NewMinCostMaxFlow(n int) *MinCostMaxFlow {
	return &MinCostMaxFlow{
		n: n,
		g: make([][]*neighbor, n),
	}
}

// AddEdge adds an edge oriented from from to to with capacity cap and cost cost.
// It returns an integer k such that this is the k-th edge that is added.
func (m *MinCostMaxFlow) AddEdge(from, to int, cap, cost int) int {
	ps := len(m.pos)
	m.pos = append(m.pos, [2]int{from, len(m.g[from])})
	m.g[from] = append(m.g[from], &neighbor{
		To:   to,
		Rev:  len(m.g[to]),
		Cap:  cap,
		Cost: cost,
	})
	m.g[to] = append(m.g[to], &neighbor{
		To:   from,
		Rev:  len(m.g[from]) - 1,
		Cap:  0,
		Cost: -cost,
	})
	return ps
}

// GetEdge returns the current internal state of the edges.
// The edges are ordered in the same order as added by add_edge.
func (m *MinCostMaxFlow) GetEdge(i int) *Edge {
	e := m.g[m.pos[i][0]][m.pos[i][1]]
	re := m.g[e.To][e.Rev]
	return &Edge{
		From: m.pos[i][0],
		To:   e.To,
		Cap:  e.Cap + re.Cap,
		FLow: re.Cap,
		Cost: e.Cost,
	}
}

// Edges returns the current internal state of the edges.
// The edges are ordered in the same order as added by add_edge.
func (m *MinCostMaxFlow) Edges() []*Edge {
	ps := len(m.pos)
	result := make([]*Edge, ps)
	for i := 0; i < ps; i++ {
		result[i] = m.GetEdge(i)
	}
	return result
}

// Flow augments the flow from s to t as much as possible.
// It returns the amount of the flow and the cost.
// It augments the s-t flow as much as possible.
func (m *MinCostMaxFlow) Flow(s, t int) (flow, cost int) {
	return m.FlowWithLimit(s, t, INF)
}

// FlowWithLimit augments the flow from s to t as much as possible.
// It returns the amount of the flow and the cost.
// !It augments the s-t flow as much as possible, until reaching the amount of flow_limit.
func (m *MinCostMaxFlow) FlowWithLimit(s, t int, flowLimit int) (flow, cost int) {
	sl := m.SlopeWithLimit(s, t, flowLimit)
	last := sl[len(sl)-1]
	return last[0], last[1]
}

// Slope returns the list of the changepoints
// Let g be a function such that g(x) is the cost of the minimum cost s-t flow when the amount of the flow is exactly x.
// g is known to be piecewise linear.
// It returns g as the list of the changepoints, that satisfies the followings.
// The first element of the list is (0, 0).
// Both of [0] and [1] are strictly increasing.
// No three changepoints are on the same line.
// The last element of the list is (x, g(x)), where x is the maximum amount of the s-t flow.
func (m *MinCostMaxFlow) Slope(s, t int) [][2]int {
	return m.SlopeWithLimit(s, t, INF)
}

// SlopeWithLimit returns the list of the changepoints
// Let g be a function such that g(x) is the cost of the minimum cost s-t flow when the amount of the flow is exactly x.
// g is known to be piecewise linear.
// It returns g as the list of the changepoints, that satisfies the followings.
// The first element of the list is (0, 0).
// Both of [0] and [1] are strictly increasing.
// No three changepoints are on the same line.
// !The last element of the list is (y, g(y)), where y = min(x, flow_limit).
func (m *MinCostMaxFlow) SlopeWithLimit(s, t int, flowLimit int) [][2]int {
	// variants (C = maxcost):
	// -(n-1)C <= dual[s] <= dual[i] <= dual[t] = 0
	// reduced cost (= e.cost + dual[e.from] - dual[e.to]) >= 0 for all edge
	n := m.n
	dual := make([]int, n)
	dist := make([]int, n)
	pv := make([]int, n)
	pe := make([]int, n)
	vis := make([]bool, n)
	dualRef := func() bool {
		for i := 0; i < n; i++ {
			dist[i] = INF
			pv[i] = -1
			pe[i] = -1
			vis[i] = false
		}
		dist[s] = 0
		que := make(PriorityQueue, 0)
		que.Push(&Q{Key: 0, To: s})
		for len(que) > 0 {
			v := que[0].To
			que.Pop()
			if vis[v] {
				continue
			}
			vis[v] = true
			if v == t {
				break
			}
			// dist[v] = shortest(s, v) + dual[s] - dual[v]
			// dist[v] >= 0 (all reduced cost are positive)
			// dist[v] <= (n-1)C
			for i := 0; i < len(m.g[v]); i++ {
				e := m.g[v][i]
				if vis[e.To] || e.Cap == 0 {
					continue
				}
				// |-dual[e.to] + dual[v]| <= (n-1)C
				// cost <= C - -(n-1)C + 0 = nC
				cost := e.Cost - dual[e.To] + dual[v]
				if dist[e.To]-dist[v] > cost {
					dist[e.To] = dist[v] + cost
					pv[e.To] = v
					pe[e.To] = i
					que.Push(&Q{Key: dist[e.To], To: e.To})
				}
			}
		}
		if !vis[t] {
			return false
		}
		for v := 0; v < n; v++ {
			if !vis[v] {
				continue
			}
			// dual[v] = dual[v] - dist[t] + dist[v]
			//         = dual[v] - (shortest(s, t) + dual[s] - dual[t]) + (shortest(s, v) + dual[s] - dual[v])
			//         = - shortest(s, t) + dual[t] + shortest(s, v)
			//         = shortest(s, v) - shortest(s, t) >= 0 - (n-1)C
			dual[v] -= dist[t] - dist[v]
		}
		return true
	}

	var flow, cost, prevCost int
	prevCost = -1
	result := make([][2]int, 0)
	result = append(result, [2]int{flow, cost})
	for flow < flowLimit {
		if !dualRef() {
			break
		}
		c := flowLimit - flow
		for v := t; v != s; v = pv[v] {
			if c > m.g[pv[v]][pe[v]].Cap {
				c = m.g[pv[v]][pe[v]].Cap
			}
		}
		for v := t; v != s; v = pv[v] {
			e := m.g[pv[v]][pe[v]]
			e.Cap -= c
			m.g[v][e.Rev].Cap += c
		}
		d := -dual[s]
		flow += c
		cost += c * d
		if prevCost == d {
			result = result[:len(result)-1]
		}
		result = append(result, [2]int{flow, cost})
	}

	return result
}

type neighbor struct {
	To, Rev   int
	Cap, Cost int
}

type Q struct {
	Key int
	To  int
}

type PriorityQueue []*Q

func (p PriorityQueue) less(i, j int) bool {
	return p[i].Key < p[j].Key
}

func (p PriorityQueue) swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *PriorityQueue) Push(q *Q) {
	n := len(*p)
	*p = append(*p, q)
	for n > 0 && p.less(n, (n-1)>>1) {
		p.swap(n, (n-1)>>1)
		n = (n - 1) >> 1
	}
}

func (p *PriorityQueue) Pop() *Q {
	old := *p
	x := old[0]
	n := len(old) - 1
	old.swap(0, n)
	*p = (*p)[:n]
	loop := true
	for i := 0; i < n && loop; {
		l := i<<1 + 1
		r := i<<1 + 2
		switch {
		case r < n && p.less(r, l) && p.less(r, i):
			p.swap(i, r)
			i = r
		case l < n && p.less(l, i):
			p.swap(i, l)
			i = l
		default:
			loop = false
		}
	}
	return x
}

func (p *PriorityQueue) Peek() *Q {
	return (*p)[0]
}
