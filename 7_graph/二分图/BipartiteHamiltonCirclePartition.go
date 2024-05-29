package main

import "fmt"

func main() {
	P := NewBipartiteHamiltonCirclePartition(3, 3)
	P.AddEdge(0, 0)
	P.AddEdge(1, 1)
	P.AddEdge(2, 2)
	P.Solve()
	fmt.Println(P.GetAllCircle())
}

// 二分图哈密尔顿回路分解.
type BipartiteHamiltonCirclePartition struct {
	g       [][]*FlowEdge
	l, r    int32
	visited []bool
}

func NewBipartiteHamiltonCirclePartition(l, r int32) *BipartiteHamiltonCirclePartition {
	res := &BipartiteHamiltonCirclePartition{l: l, r: r}
	res.g = make([][]*FlowEdge, res.sink()+1)
	for i := int32(0); i < l; i++ {
		res.addFlowEdge(res.g, res.src(), res.left(i), 2)
	}
	for i := int32(0); i < r; i++ {
		res.addFlowEdge(res.g, res.right(i), res.sink(), 2)
	}
	return res
}

func (b *BipartiteHamiltonCirclePartition) AddEdge(u, v int32) {
	b.addFlowEdge(b.g, b.left(u), b.right(v), 1)
}

func (b *BipartiteHamiltonCirclePartition) Solve() bool {
	if b.l != b.r {
		return false
	}
	dinic := NewDinic()
	flow := dinic.Apply(b.g, b.src(), b.sink(), 2*int(b.l))
	return flow == 2*int(b.l)
}

// 左边的点从 0 开始, 右边的点从 L 开始.
func (b *BipartiteHamiltonCirclePartition) GetAllCircle() [][]int32 {
	b.visited = make([]bool, b.sink()+1)
	res := make([][]int32, 0)
	seq := make([]int32, 0, b.sink()+1)
	for i := int32(0); i < b.l; i++ {
		if b.visited[i] {
			continue
		}
		seq = seq[:0]
		b.dfsL(i, &seq)
		res = append(res, append(seq[:0:0], seq...))
	}
	return res
}

func (b *BipartiteHamiltonCirclePartition) left(i int32) int32  { return i }
func (b *BipartiteHamiltonCirclePartition) right(i int32) int32 { return i + b.l }
func (b *BipartiteHamiltonCirclePartition) src() int32          { return b.l + b.r }
func (b *BipartiteHamiltonCirclePartition) sink() int32         { return b.src() + 1 }

func (b *BipartiteHamiltonCirclePartition) dfsL(root int32, seq *[]int32) {
	b.visited[root] = true
	*seq = append(*seq, root)
	for _, e := range b.g[root] {
		if e.real && e.flow == 1 && !b.visited[e.to] {
			b.dfsR(e.to, seq)
			return
		}
	}
}

func (b *BipartiteHamiltonCirclePartition) dfsR(root int32, seq *[]int32) {
	b.visited[root] = true
	*seq = append(*seq, root)
	for _, e := range b.g[root] {
		if !e.real && e.rev.flow == 1 && !b.visited[e.to] {
			b.dfsL(e.to, seq)
			return
		}
	}
}

func (b *BipartiteHamiltonCirclePartition) addFlowEdge(g [][]*FlowEdge, s, t int32, cap int) {
	real := NewFlowEdge(t, 0, true)
	virtual := NewFlowEdge(s, cap, false)
	real.rev = virtual
	virtual.rev = real
	g[s] = append(g[s], real)
	g[t] = append(g[t], virtual)
}

const INF int = 1e18
const INF32 int32 = 1e9 + 10

type Dinic struct {
	g     [][]*FlowEdge
	s, t  int32
	queue []int32
	dists []int32
	iters []int32
}

func NewDinic() *Dinic {
	return &Dinic{}
}

func (d *Dinic) Apply(g [][]*FlowEdge, s, t int32, send int) int {
	d.prepare(int32(len(g)))
	d.s = s
	d.t = t
	d.g = g
	flow := 0
	for flow < send {
		// bfs for flow
		for i := 0; i < len(g); i++ {
			d.dists[i] = INF32
		}
		d.dists[t] = 0
		d.queue = d.queue[:0]
		d.queue = append(d.queue, t)
		for len(d.queue) > 0 {
			head := d.queue[0]
			d.queue = d.queue[1:]
			for _, e := range g[head] {
				if e.flow != 0 && d.dists[e.to] == INF32 {
					d.dists[e.to] = d.dists[head] + 1
					d.queue = append(d.queue, e.to)
				}
			}
		}
		if d.dists[s] == INF32 {
			break
		}
		for i := 0; i < len(g); i++ {
			d.iters[i] = int32(len(g[i])) - 1
		}
		flow += d.send(s, send-flow)
	}
	return flow
}

func (d *Dinic) send(root int32, flow int) int {
	if root == d.t {
		return flow
	}
	snapshot := flow
	for d.iters[root] >= 0 && flow > 0 {
		e := d.g[root][d.iters[root]]
		if d.dists[e.to]+1 == d.dists[root] && e.rev.flow != 0 {
			sent := d.send(e.to, min(flow, e.rev.flow))
			if sent > 0 {
				flow -= sent
				e.flow += sent
				e.rev.flow -= sent
				continue
			}
		}
		d.iters[root]--
	}
	return snapshot - flow
}

func (d *Dinic) prepare(vertexNum int32) {
	if d.dists != nil && int32(len(d.dists)) >= vertexNum {
		return
	}
	d.queue = make([]int32, 0, vertexNum)
	d.dists = make([]int32, vertexNum)
	d.iters = make([]int32, vertexNum)
}

type FlowEdge struct {
	to   int32
	flow int
	real bool
	rev  *FlowEdge
}

func NewFlowEdge(to int32, flow int, real bool) *FlowEdge {
	return &FlowEdge{to: to, flow: flow, real: real}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
