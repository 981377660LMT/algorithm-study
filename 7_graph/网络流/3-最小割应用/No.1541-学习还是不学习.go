package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/1541
	// 期末考试
	// 有n个考试科目,每学一个科目就能多拿base分
	// 对于每个科目i,可以花费cost来学习，学习之后有额外的收益:
	// 对于科目subjects[j],如果i和subjects[j]都学习了,那么就能多拿到bonus[j]分
	// !最大化(总分-花费)
	// n<=100

	// !每个科目学习还是不学习 => 燃やす埋める
	// 先学习所有科目,然后再割掉不学每个科目的代价(最小割)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, base int
	fmt.Fscan(in, &n, &base)

	flow := NewMaxFlowGraph(n + 2 + n*n)
	S, T := n, n+1 // !S:不学 T:学
	ptr := n + 2

	res := 0
	for i := 0; i < n; i++ {
		var k, cost int
		fmt.Fscan(in, &k, &cost)

		subjects := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &subjects[j])
			subjects[j]--
		}

		bonus := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &bonus[j])
		}

		res += base
		flow.AddEdge(S, i, cost) // 学失去cost分
		flow.AddEdge(i, T, base) // 不学失去base分
		for j := 0; j < k; j++ {
			res += bonus[j]
			flow.AddEdge(i, ptr, INF)
			flow.AddEdge(subjects[j], ptr, INF)
			flow.AddEdge(ptr, T, bonus[j]) // 学了i并且学了subjects[j]，此时不学subjects[j]失去bonus[j]分
			ptr++                          // !ptr表示状态:学了i并且学了subjects[j]
		}
	}

	fmt.Fprintln(out, res-flow.Flow(S, T))
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
