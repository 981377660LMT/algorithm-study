// https://maspypy.github.io/library/flow/maxflow.hpp
// Mincut
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
	g := NewMaxFlowGraph(n)
	for i := 0; i < m; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		g.AddEdge(a, b, c)
	}
	fmt.Fprintln(out, g.Flow(0, n-1))

}

const INF int = 1e18

type Edge struct{ to, rev, cap int }
type MaxFlowGraph struct {
	N           int
	G           [][]Edge
	prog, level []int
	flowRes     int
	calculated  bool
}

func NewMaxFlowGraph(n int) *MaxFlowGraph {
	return &MaxFlowGraph{N: n, G: make([][]Edge, n)}
}

func (g *MaxFlowGraph) AddEdge(from, to, cap int) {
	g.G[from] = append(g.G[from], Edge{to, len(g.G[to]), cap})
	g.G[to] = append(g.G[to], Edge{from, len(g.G[from]) - 1, 0})
}

func (g *MaxFlowGraph) Flow(source, sink int) int {
	if g.calculated {
		return g.flowRes
	}
	g.calculated = true
	for g.setLevel(source, sink) {
		g.prog = make([]int, g.N)
		for {
			f := g.flowDfs(source, sink, INF)
			if f == 0 {
				break
			}
			g.flowRes += f
			g.flowRes = min(g.flowRes, INF)
			if g.flowRes == INF {
				return g.flowRes
			}
		}
	}
	return g.flowRes
}

// 返回最小割的值和每个点是否属于最小割
func (g *MaxFlowGraph) Cut(source, sink int) (minCut int, isCut []bool) {
	minCut = g.Flow(source, sink)
	isCut = make([]bool, g.N)
	for i := 0; i < g.N; i++ {
		isCut[i] = g.level[i] < 0
	}
	return
}

// 残量图的边(from,to,remainCap)
func (g *MaxFlowGraph) GetEdges() (edges [][3]int) {
	for v := 0; v < g.N; v++ {
		for _, e := range g.G[v] {
			edges = append(edges, [3]int{v, e.to, e.cap})
		}
	}
	return
}

func (g *MaxFlowGraph) setLevel(source, sink int) bool {
	g.level = make([]int, g.N)
	for i := range g.level {
		g.level[i] = -1
	}
	g.level[source] = 0
	q := []int{source}
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for _, e := range g.G[v] {
			if e.cap > 0 && g.level[e.to] == -1 {
				g.level[e.to] = g.level[v] + 1
				if e.to == sink {
					return true
				}
				q = append(q, e.to)
			}
		}
	}
	return false
}

func (g *MaxFlowGraph) flowDfs(v, sink, lim int) int {
	if v == sink {
		return lim
	}
	res := 0
	for i := &g.prog[v]; *i < len(g.G[v]); *i++ {
		e := &g.G[v][*i]
		if e.cap > 0 && g.level[e.to] == g.level[v]+1 {
			a := g.flowDfs(e.to, sink, min(lim, e.cap))
			if a > 0 {
				e.cap -= a
				g.G[e.to][e.rev].cap += a
				res += a
				lim -= a
				if lim == 0 {
					break
				}
			}
		}
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
