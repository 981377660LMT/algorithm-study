package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N := int32(io.NextInt())
	A, B := make([]int32, N), make([]int32, N)
	sum := int32(0)
	groupSum := [3]int32{}
	for i := int32(0); i < N; i++ {
		A[i] = int32(io.NextInt())
		A[i]--
		B[i] = int32(io.NextInt())
		sum += B[i]
		groupSum[A[i]] += B[i]
	}

	if sum%3 != 0 {
		io.Println(-1)
		return
	}

	S := 6
	T := 7
	mcmf := NewMinCostMaxFlow(8, int(S), int(T))
	for i := int32(0); i < 3; i++ {
		mcmf.AddEdge(S, int(i), int(groupSum[i]), 0)
	}
	for i := int32(0); i < N; i++ {
		for t := int32(0); t < 3; t++ {
			if A[i] == t {
				mcmf.AddEdge(int(A[i]), int(t)+3, INF, 0)
			} else {
				mcmf.AddEdge(int(A[i]), int(t)+3, INF, 1)
			}
		}
	}
	for i := 0; i < 3; i++ {
		mcmf.AddEdge(i+3, T, int(sum/3), 0)
	}

	_, cost := mcmf.Flow()
	io.Println(cost)
}

const INF int = 1e18

type Edge struct{ from, to, cap, flow, cost, id int }
type MinCostMaxFlow struct {
	AddEdge        func(from, to, cap, cost int)
	Flow           func() (maxFlow int, minCost int)
	FlowWithLimit  func(flowLimit int) (maxFlow int, minCost int)
	Slope          func() [][2]int
	SlopeWithLimit func(flowLimit int) [][2]int
	Edges          func() []Edge
}

func NewMinCostMaxFlow(n, start, end int) *MinCostMaxFlow {
	type neighbor struct {
		to   int
		rid  int
		cap  int
		cost int
		eid  int
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

	FlowWithLimit := func(flowLimit int) (maxFlow int, minCost int) {
		for maxFlow < flowLimit {
			if !spfa() {
				break
			}

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
