// 最小费用流/最小费用最大流
// "边权要求非负"，或者"给出的图为dag且顶点编号为拓扑序".
// https://kopricky.github.io/code/NetworkFlow/min_cost_flow_DAG.html
// https://atcoder.jp/contests/tdpc/tasks/tdpc_graph
// 时间复杂度O(|f|mlogn), |f|为流量
//
// api:
// NewMinCostFlow(n, source, sink int32) *MinCostFlow
// NewMinCostFlowFromDag(n, source, sink int32) *MinCostFlow
// AddEdge(from, to int32, cap, cost int) int32
// Flow() (flow, cost int)
// FlowWithLimit(limit int) (flow, cost int)
// Slope() [][2]int
// SlopeWithLimit(limit int) [][2]int
// PathDecomposition() [][]int32
// GetEdge(i int32) edge
// Edges() []edge
// Debug()

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// assignment()
	// judge()
	// abc214h()
	abcGraph()
	// abcMinCostFlow()
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_6_B
func judge() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	var f int
	fmt.Fscan(in, &n, &m, &f)
	pd := NewMinCostFlow(n, 0, n-1)
	for i := int32(0); i < m; i++ {
		var from, to int32
		var cap, cost int
		fmt.Fscan(in, &from, &to, &cap, &cost)
		pd.AddEdge(from, to, cap, cost)
	}
	flow, cost := pd.FlowWithLimit(f)
	if flow < f {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, cost)
	}
}

// https://judge.yosupo.jp/problem/assignment
// 给定一个n*n的矩阵，每个元素表示从i到j的费用.
// 选择n个元素，使得每行每列只有一个元素，求最小费用.
// 输出每行选择的列.
func assignment() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	grid := make([][]int, n)
	for i := int32(0); i < n; i++ {
		grid[i] = make([]int, n)
		for j := int32(0); j < n; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	source := int32(0)
	left := func(i int32) int32 { return 1 + i }
	right := func(i int32) int32 { return 1 + n + i }
	sink := right(n)
	M := NewMinCostFlowFromDag(n+n+2, source, sink)
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < n; j++ {
			M.AddEdge(left(i), right(j), 1, grid[i][j])
		}
	}
	for i := int32(0); i < n; i++ {
		M.AddEdge(source, left(i), 1, 0)
		M.AddEdge(right(i), sink, 1, 0)
	}

	_, minCost := M.Flow()
	edges := M.Edges()
	res := make([]int32, n)
	for _, e := range edges {
		if e.flow > 0 && 1 <= e.from && e.from <= n {
			res[e.from-1] = e.to - right(0)
		}
	}
	fmt.Fprintln(out, minCost)
	for i := int32(0); i < n; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}

// H - Collecting
// https://atcoder.jp/contests/abc214h/tasks/abc214_h
// !有一张N个点M条边的有向图，每个点有一个点权ai.
// !现在要找出k条经过点0的路径，使得这些路径的并集的点权和尽量大。
// dag路径覆盖最大点权和
// n,m<=2e5,k<=10.
func abc214h() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	var k int
	fmt.Fscan(in, &n, &m, &k)
	graph := make([][]int32, n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		graph[u] = append(graph[u], v)
	}
	weights := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &weights[i])
	}

	count, belong := StronglyConnectedComponent(graph)
	dag := SccDag(graph, count, belong)
	groupSum := make([]int, count)
	for i := int32(0); i < n; i++ {
		groupSum[belong[i]] += weights[i]
	}

	source := int32(0)
	left := func(i int32) int32 { return 1 + 2*i + 0 }
	right := func(i int32) int32 { return 1 + 2*i + 1 }
	sink := 1 + count + count
	M := NewMinCostFlowFromDag(count+count+2, source, sink)
	M.AddEdge(source, left(belong[0]), k, 0) // 经过0
	for i := int32(0); i < count; i++ {
		M.AddEdge(left(i), right(i), 1, -groupSum[i])
		M.AddEdge(left(i), right(i), k, 0)
	}
	for i := int32(0); i < count; i++ {
		M.AddEdge(right(i), sink, k, 0)
	}
	for from := int32(0); from < int32(len(dag)); from++ {
		nexts := dag[from]
		for _, to := range nexts {
			M.AddEdge(right(from), left(to), k, 0)
		}
	}

	_, minCost := M.Flow()
	fmt.Fprintln(out, -minCost)
}

// https://atcoder.jp/contests/tdpc/tasks/tdpc_graph
// !有一张N个点M条边的有向图.
// !现在要找出两条路径，使得这两条路径的并集的点的个数最多.
// n<=300.
func abcGraph() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	adjMatrix := make([][]bool, n)
	for i := int32(0); i < n; i++ {
		adjMatrix[i] = make([]bool, n)
	}
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < n; j++ {
			var v int8
			fmt.Fscan(in, &v)
			adjMatrix[i][j] = v == 1
		}
	}

	k := 2
	graph := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < n; j++ {
			if adjMatrix[i][j] {
				graph[i] = append(graph[i], j)
			}
		}
	}
	weights := make([]int, n)
	for i := int32(0); i < n; i++ {
		weights[i] = 1
	}

	count, belong := StronglyConnectedComponent(graph)
	dag := SccDag(graph, count, belong)
	groupSum := make([]int, count)
	for i := int32(0); i < n; i++ {
		groupSum[belong[i]] += weights[i]
	}
	source := int32(0)
	source2 := int32(1)
	left := func(i int32) int32 { return 2 + 2*i + 0 }
	right := func(i int32) int32 { return 2 + 2*i + 1 }
	sink := 2 + count + count
	M := NewMinCostFlowFromDag(count+count+3, source, sink)
	M.AddEdge(source, source2, k, 0)
	for i := int32(0); i < count; i++ {
		M.AddEdge(source2, left(i), k, 0)
	}
	for i := int32(0); i < count; i++ {
		M.AddEdge(left(i), right(i), 1, -groupSum[i])
		M.AddEdge(left(i), right(i), k, 0)
	}
	for i := int32(0); i < count; i++ {
		M.AddEdge(right(i), sink, k, 0)
	}
	for from := int32(0); from < int32(len(dag)); from++ {
		nexts := dag[from]
		for _, to := range nexts {
			M.AddEdge(right(from), left(to), k, 0)
		}
	}

	_, minCost := M.Flow()
	fmt.Fprintln(out, -minCost)
}

// https://atcoder.jp/contests/practice2/tasks/practice2_e
// 给定一个n*n的矩阵，选择若干个单元格使得和最大.
// 但是不能同一行或者同一列选择个数不能超过k.
// 输出方案.
func abcMinCostFlow() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	var k int
	fmt.Fscan(in, &n, &k)
	grid := make([][]int, n)
	for i := int32(0); i < n; i++ {
		grid[i] = make([]int, n)
		for j := int32(0); j < n; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	ROW, COL := n, n
	S, T := ROW+COL, ROW+COL+1
	BIG := int(1e9 + 10)
	/**
	 * generate (s -> row -> column -> t) graph
	 * i-th row correspond to vertex i
	 * i-th col correspond to vertex n + i
	 **/
	M := NewMinCostFlow(ROW+COL+2, S, T)

	// we can "waste" the flow
	M.AddEdge(S, T, int(n)*k, BIG)
	for i := int32(0); i < n; i++ {
		M.AddEdge(S, i, k, 0)
		M.AddEdge(ROW+i, T, k, 0)
	}
	for i := int32(0); i < ROW; i++ {
		for j := int32(0); j < COL; j++ {
			M.AddEdge(i, ROW+j, 1, BIG-grid[i][j])
		}
	}

	_, cost := M.FlowWithLimit(int(n) * k)
	fmt.Fprintln(out, -cost+BIG*int(n)*k)

	visited := make([][]bool, ROW)
	for i := int32(0); i < ROW; i++ {
		visited[i] = make([]bool, COL)
	}
	path := M.PathDecomposition()
	for _, p := range path {
		if len(p) == 4 {
			x, y := p[1], p[2]-ROW
			visited[x][y] = true
		}
	}

	for i := int32(0); i < ROW; i++ {
		for j := int32(0); j < COL; j++ {
			if visited[i][j] {
				fmt.Fprint(out, "X")
			} else {
				fmt.Fprint(out, ".")
			}
		}
		fmt.Fprintln(out)
	}
}

// 100401. 放三个车的价值之和最大 II
// https://leetcode.cn/problems/maximum-value-sum-by-placing-three-rooks-ii/
func maximumValueSum(board [][]int) int64 {
	ROW, COL := int32(len(board)), int32(len(board[0]))
	S, T := ROW+COL, ROW+COL+1
	BIG := int(1e9 + 10)
	M := NewMinCostFlow(ROW+COL+2, S, T)
	for r := int32(0); r < ROW; r++ {
		M.AddEdge(S, r, 1, 0)
	}
	for r := int32(0); r < ROW; r++ {
		for c := int32(0); c < COL; c++ {
			M.AddEdge(r, ROW+c, 1, BIG-board[r][c]) // 求最大值，取负数
		}
	}
	for c := int32(0); c < COL; c++ {
		M.AddEdge(ROW+c, T, 1, 0)
	}
	_, cost := M.FlowWithLimit(3) // !放三个车，流量为3
	return int64(-cost + 3*BIG)
}

const INF int = 1e18

type MinCostFlow struct {
	dag             bool
	n, source, sink int32
	edges           []edge
}

type edge struct {
	from, to  int32
	cap, flow int
	cost      int
}

type _edge struct {
	to, rev   int32
	cap, cost int
}

func NewMinCostFlow(n, source, sink int32) *MinCostFlow {
	checkArguments(n, source, sink)
	return &MinCostFlow{n: n, source: source, sink: sink}
}

func NewMinCostFlowFromDag(n, source, sink int32) *MinCostFlow {
	checkArguments(n, source, sink)
	return &MinCostFlow{dag: true, n: n, source: source, sink: sink}
}

func (mcf *MinCostFlow) AddEdge(from, to int32, cap, cost int) int32 {
	if from < 0 || from >= mcf.n {
		panic("from out of range")
	}
	if to < 0 || to >= mcf.n {
		panic("to out of range")
	}
	if cap < 0 {
		panic("cap is negative")
	}
	if !mcf.dag && cost < 0 {
		panic("cost is negative in non-dag")
	}
	if mcf.dag && from >= to {
		panic("from >= to in dag")
	}
	m := int32(len(mcf.edges))
	mcf.edges = append(mcf.edges, edge{from: from, to: to, cap: cap, cost: cost})
	return m
}

func (mcf *MinCostFlow) Flow() (flow, cost int) {
	return mcf.FlowWithLimit(INF)
}
func (mcf *MinCostFlow) FlowWithLimit(limit int) (flow, cost int) {
	res := mcf.SlopeWithLimit(limit)
	return res[len(res)-1][0], res[len(res)-1][1]
}

func (mcf *MinCostFlow) Slope() [][2]int { return mcf.SlopeWithLimit(INF) }
func (mcf *MinCostFlow) SlopeWithLimit(limit int) [][2]int {
	m := int32(len(mcf.edges))
	edgeIndex := make([]int32, m)

	g := func() *csr {
		degree, redgeIndex := make([]int32, mcf.n), make([]int32, m)
		elist := make([]elistPair, 0, 2*m)
		for i := int32(0); i < m; i++ {
			e := mcf.edges[i]
			edgeIndex[i] = degree[e.from]
			degree[e.from]++
			redgeIndex[i] = degree[e.to]
			degree[e.to]++
			elist = append(elist, elistPair{first: e.from, second: _edge{to: e.to, rev: -1, cap: e.cap - e.flow, cost: e.cost}})
			elist = append(elist, elistPair{first: e.to, second: _edge{to: e.from, rev: -1, cap: e.flow, cost: -e.cost}})
		}
		csr := newCsr(mcf.n, elist)
		for i := int32(0); i < m; i++ {
			e := mcf.edges[i]
			edgeIndex[i] += csr.start[e.from]
			redgeIndex[i] += csr.start[e.to]
			csr.elist[edgeIndex[i]].rev = redgeIndex[i]
			csr.elist[redgeIndex[i]].rev = edgeIndex[i]
		}
		return csr
	}()

	res := mcf._slope(g, limit)
	for i := int32(0); i < m; i++ {
		e := g.elist[edgeIndex[i]]
		mcf.edges[i].flow = mcf.edges[i].cap - e.cap
	}
	return res
}

// 路径还原, O(f*(n+m)).
func (mcf *MinCostFlow) PathDecomposition() [][]int32 {
	to := make([][]int32, mcf.n)
	for _, e := range mcf.edges {
		for i := 0; i < e.flow; i++ {
			to[e.from] = append(to[e.from], e.to)
		}
	}
	var res [][]int32
	visited := make([]bool, mcf.n)
	for len(to[mcf.source]) > 0 {
		path := []int32{mcf.source}
		visited[mcf.source] = true
		for path[len(path)-1] != mcf.sink {
			last := &to[path[len(path)-1]]
			tmp := (*last)[len(*last)-1]
			*last = (*last)[:len(*last)-1]
			for visited[tmp] {
				visited[path[len(path)-1]] = false
				path = path[:len(path)-1]
			}
			path = append(path, tmp)
			visited[tmp] = true
		}
		for _, v := range path {
			visited[v] = false
		}
		res = append(res, path)
	}
	return res
}

func (mcf *MinCostFlow) _slope(g *csr, flowLimit int) [][2]int {
	if mcf.dag {
		if mcf.source != 0 || mcf.sink != mcf.n-1 {
			panic("source and sink must be 0 and n-1 in dag")
		}
	}
	dualDist := make([][2]int, mcf.n)
	prevE := make([]int32, mcf.n)
	visited := make([]bool, mcf.n)

	queMin := make([]int32, 0)
	pq := make([]pqPair, 0)
	pqLess := func(i, j int32) bool { return pq[i].key < pq[j].key }

	dualRef := func() bool {
		for i := int32(0); i < mcf.n; i++ {
			dualDist[i][1] = INF
		}
		for i := int32(0); i < mcf.n; i++ {
			visited[i] = false
		}
		queMin = queMin[:0]
		pq = pq[:0]
		heapR := int32(0)

		dualDist[mcf.source][1] = 0
		queMin = append(queMin, mcf.source)
		for len(queMin) > 0 || len(pq) > 0 {
			var v int32
			if len(queMin) > 0 {
				v = queMin[len(queMin)-1]
				queMin = queMin[:len(queMin)-1]
			} else {
				for heapR < int32(len(pq)) {
					heapR++
					heapUp(pq, heapR-1, pqLess)
				}
				v = pq[0].to
				pq[0], pq[heapR-1] = pq[heapR-1], pq[0]
				heapDown(pq, 0, heapR-1, pqLess)
				pq = pq[:len(pq)-1]
				heapR--
			}
			if visited[v] {
				continue
			}
			visited[v] = true
			if v == mcf.sink {
				break
			}
			dualV, distV := dualDist[v][0], dualDist[v][1]
			for i := g.start[v]; i < g.start[v+1]; i++ {
				e := &g.elist[i]
				if e.cap == 0 {
					continue
				}
				cost := e.cost - dualDist[e.to][0] + dualV
				if dualDist[e.to][1] > distV+cost {
					distTo := distV + cost
					dualDist[e.to][1] = distTo
					prevE[e.to] = e.rev
					if distTo == distV {
						queMin = append(queMin, e.to)
					} else {
						pq = append(pq, pqPair{to: e.to, key: distTo})
					}
				}
			}
		}
		if !visited[mcf.sink] {
			return false
		}
		for i := int32(0); i < mcf.n; i++ {
			if !visited[i] {
				continue
			}
			dualDist[i][0] -= dualDist[mcf.sink][1] - dualDist[i][1]
		}
		return true
	}

	dualRefDag := func() bool {
		for i := int32(0); i < mcf.n; i++ {
			dualDist[i][1] = INF
		}
		dualDist[mcf.source][1] = 0
		for i := int32(0); i < mcf.n; i++ {
			visited[i] = false
		}
		visited[mcf.source] = true
		for v := int32(0); v < mcf.n; v++ {
			if !visited[v] {
				continue
			}
			dualV, distV := dualDist[v][0], dualDist[v][1]
			for i := g.start[v]; i < g.start[v+1]; i++ {
				e := &g.elist[i]
				if e.cap == 0 {
					continue
				}
				cost := e.cost - dualDist[e.to][0] + dualV
				if dualDist[e.to][1] > distV+cost {
					visited[e.to] = true
					distTo := distV + cost
					dualDist[e.to][1] = distTo
					prevE[e.to] = e.rev
				}
			}
		}
		if !visited[mcf.sink] {
			return false
		}
		for i := int32(0); i < mcf.n; i++ {
			if !visited[i] {
				continue
			}
			dualDist[i][0] -= dualDist[mcf.sink][1] - dualDist[i][1]
		}
		return true
	}

	flow, cost := 0, 0
	prevCostPerFlow := -1
	res := [][2]int{{0, 0}}
	for flow < flowLimit {
		if mcf.dag && flow == 0 {
			if !dualRefDag() {
				break
			}
		} else {
			if !dualRef() {
				break
			}
		}
		c := flowLimit - flow
		for v := mcf.sink; v != mcf.source; v = g.elist[prevE[v]].to {
			c = min(c, g.elist[g.elist[prevE[v]].rev].cap)
		}
		for v := mcf.sink; v != mcf.source; v = g.elist[prevE[v]].to {
			e := &g.elist[prevE[v]]
			e.cap += c
			g.elist[e.rev].cap -= c
		}
		d := -dualDist[mcf.source][0]
		flow += c
		cost += c * d
		if prevCostPerFlow == d {
			res = res[:len(res)-1]
		}
		res = append(res, [2]int{flow, cost})
		prevCostPerFlow = d
	}
	return res
}

func (mcf *MinCostFlow) GetEdge(i int32) edge {
	return mcf.edges[i]
}
func (mcf *MinCostFlow) Edges() []edge {
	return mcf.edges
}

func (mcf *MinCostFlow) Debug() {
	fmt.Println("flow graph")
	fmt.Println("from, to, cap, cost")
	for _, e := range mcf.edges {
		fmt.Println(e.from, e.to, e.cap, e.cost)
	}
}

func checkArguments(n int32, source int32, sink int32) {
	if source < 0 || source >= n {
		panic("source out of range")
	}
	if sink < 0 || sink >= n {
		panic("sink out of range")
	}
	if source == sink {
		panic("source equals to sink")
	}
}

type elistPair struct {
	first  int32
	second _edge
}
type csr struct {
	start []int32
	elist []_edge
}

func newCsr(n int32, edges []elistPair) *csr {
	start := make([]int32, n+1)
	elist := make([]_edge, len(edges))
	for i := int32(0); i < int32(len(edges)); i++ {
		start[edges[i].first+1]++
	}
	for i := int32(1); i <= n; i++ {
		start[i] += start[i-1]
	}
	counter := append(start[:0:0], start...)
	for _, e := range edges {
		elist[counter[e.first]] = e.second
		counter[e.first]++
	}
	return &csr{start: start, elist: elist}
}

// heapUtils
type pqPair = struct {
	to  int32
	key int
}

func heapUp(data []pqPair, i0 int32, less func(a, b int32) bool) {
	for {
		i := (i0 - 1) / 2
		if i == i0 || !less(i0, i) {
			break
		}
		data[i], data[i0] = data[i0], data[i]
		i0 = i
	}
}
func heapDown(data []pqPair, i0, n int32, less func(a, b int32) bool) {
	i := i0
	for {
		j1 := (i << 1) | 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !less(j, i) {
			break
		}
		data[i], data[j] = data[j], data[i]
		i = j
	}
}

// 有向图强连通分量分解.
func StronglyConnectedComponent(graph [][]int32) (count int32, belong []int32) {
	n := int32(len(graph))
	belong = make([]int32, n)
	low := make([]int32, n)
	order := make([]int32, n)
	for i := range order {
		order[i] = -1
	}
	now := int32(0)
	path := []int32{}

	var dfs func(int32)
	dfs = func(v int32) {
		low[v] = now
		order[v] = now
		now++
		path = append(path, v)
		for _, to := range graph[v] {
			if order[to] == -1 {
				dfs(to)
				low[v] = min32(low[v], low[to])
			} else {
				low[v] = min32(low[v], order[to])
			}
		}
		if low[v] == order[v] {
			for {
				u := path[len(path)-1]
				path = path[:len(path)-1]
				order[u] = n
				belong[u] = count
				if u == v {
					break
				}
			}
			count++
		}
	}

	for i := int32(0); i < n; i++ {
		if order[i] == -1 {
			dfs(i)
		}
	}
	for i := int32(0); i < n; i++ {
		belong[i] = count - 1 - belong[i]
	}
	return
}

// 有向图的强连通分量缩点.
func SccDag(graph [][]int32, count int32, belong []int32) (dag [][]int32) {
	dag = make([][]int32, count)
	adjSet := make([]map[int32]struct{}, count)
	for i := int32(0); i < count; i++ {
		adjSet[i] = make(map[int32]struct{})
	}
	for cur, nexts := range graph {
		for _, next := range nexts {
			if bid1, bid2 := belong[cur], belong[next]; bid1 != bid2 {
				adjSet[bid1][bid2] = struct{}{}
			}
		}
	}
	for i := int32(0); i < count; i++ {
		for next := range adjSet[i] {
			dag[i] = append(dag[i], next)
		}
	}
	return
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
