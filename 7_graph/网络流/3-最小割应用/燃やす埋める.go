// https://maspypy.github.io/library/flow/binary_optimization.hpp

package main

import (
	"bufio"
	"fmt"
	"os"
)

func 期末考试() {
	// https://yukicoder.me/problems/no/1541
	// 期末考试
	// 有n个考试科目,每学一个科目就能多拿base分
	// 对于每个科目i,可以花费cost来学习，学习之后有额外的收益:
	// 对于科目subjects[j],如果i和subjects[j]都学习了,那么就能多拿到bonus[j]分
	// !最大化(总分-花费)
	// n<=100

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, base int
	fmt.Fscan(in, &n, &base)
	bo := NewBinaryOptimization(n, false)
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

		bo.Add1(i, 0, base-cost) // 学习的收益
		for j := 0; j < k; j++ {
			bo.Add2(i, subjects[j], 0, 0, 0, bonus[j]) // 一起学习的收益
		}
	}

	res, _ := bo.Run()
	fmt.Fprintln(out, res)
}

func 选择卡牌() {
	// https://atcoder.jp/contests/abc259/tasks/abc259_g
	// 给定一个矩阵Anxn (1≤ n ≤ 100)，选择一些行列，可以得到这些行列包含的位置的并的数值和。
	// 此外要求任意选中的行列交点处不能是负数。
	// !求选择的最大值
	// !时间复杂度O(V^2E)

	// 0:不选 1:选
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)
	grid := make([][]int, ROW)
	for i := 0; i < ROW; i++ {
		grid[i] = make([]int, COL)
		for j := 0; j < COL; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	bm := NewBinaryOptimization(ROW+COL, true)

	rowId := func(r int) int { return r }
	colId := func(c int) int { return ROW + c }
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			// !https://atcoder.jp/contests/abc259/editorial/4286
			// f(x,y)=>-f(x,^y)之后就满足劣モジュラ函数了
			// 00 => 01 , 01 => 00 , 10 => 11 , 11 => 10
			x := grid[i][j]
			if x > 0 {
				bm.Add2(rowId(i), colId(j), -x, 0, -x, -x)
			}
			if x < 0 {
				bm.Add2(rowId(i), colId(j), -x, 0, 1<<45, -x)
			}
		}
	}

	res, _ := bm.Run()
	fmt.Fprintln(out, -res)
}

func main() {
	选择卡牌()
}

const INF int = 1e18

type BinaryOptimization struct {
	n            int
	source, sink int
	next         int
	baseCost     int
	edges        map[[2]int]int
	minimize     bool
}

func NewBinaryOptimization(n int, minimize bool) *BinaryOptimization {
	return &BinaryOptimization{
		n:        n,
		source:   n,
		sink:     n + 1,
		next:     n + 2,
		edges:    make(map[[2]int]int),
		minimize: minimize,
	}
}

// xi 属于 0, 1 时对应的收益.
func (bo *BinaryOptimization) Add1(i, x0, x1 int) {
	if !bo.minimize {
		x0, x1 = -x0, -x1
	}
	bo._add_1(i, x0, x1)
}

// (xi,xj) = (00,01,10,11) 时对应的收益.
//  !最小化代价时需要满足 x00 + x11 <= x01 + x10.
//  !最大化收益时需要满足 x00 + x11 >= x01 + x10.
func (bo *BinaryOptimization) Add2(i, j, x00, x01, x10, x11 int) {
	if !bo.minimize {
		x00, x01, x10, x11 = -x00, -x01, -x10, -x11
	}
	bo._add_2(i, j, x00, x01, x10, x11)
}

// (xi,xj,xk) = (000,001,010,011,100,101,110,111) 时对应的收益.
func (bo *BinaryOptimization) Add3(i, j, k, x000, x001, x010, x011, x100, x101, x110, x111 int) {
	if !bo.minimize {
		x000, x001, x010, x011, x100, x101, x110, x111 = -x000, -x001, -x010, -x011, -x100, -x101, -x110, -x111
	}
	bo._add_3(i, j, k, x000, x001, x010, x011, x100, x101, x110, x111)
}

// 返回最大收益/最小花费和每个变量的取值0/1.
func (bo *BinaryOptimization) Run() (res int, assign []int) {
	flow := NewMaxFlowGraph(bo.next)
	for key, cap := range bo.edges {
		from, to := key[0], key[1]
		flow.AddEdge(from, to, cap)
	}

	res, isCut := flow.Cut(bo.source, bo.sink)
	res += bo.baseCost
	res = min(res, INF)
	assign = make([]int, bo.n)
	for i := 0; i < bo.n; i++ {
		if isCut[i] {
			assign[i] = 1
		} else {
			assign[i] = 0
		}
	}
	if !bo.minimize {
		res = -res
	}
	return
}

func (bo *BinaryOptimization) Debug() {
	fmt.Println("base_cost", bo.baseCost)
	fmt.Println("source=", bo.source, "sink=", bo.sink)
	for key, cap := range bo.edges {
		fmt.Println(key, cap)
	}
}

func (bo *BinaryOptimization) _add_1(i, x0, x1 int) {
	if x0 <= x1 {
		bo.baseCost += x0
		bo._addEdge(bo.source, i, x1-x0)
	} else {
		bo.baseCost += x1
		bo._addEdge(i, bo.sink, x0-x1)
	}
}

// x00 + x11 <= x01 + x10
func (bo *BinaryOptimization) _add_2(i, j, x00, x01, x10, x11 int) {
	if x00+x11 > x01+x10 {
		panic("need to satisfy `x00 + x11 <= x01 + x10`.")
	}
	bo._add_1(i, x00, x10)
	bo._add_1(j, 0, x11-x10)
	bo._addEdge(i, j, x01+x10-x00-x11)
}

func (bo *BinaryOptimization) _add_3(i, j, k, x000, x001, x010, x011, x100, x101, x110, x111 int) {
	p := x000 - x100 - x010 - x001 + x110 + x101 + x011 - x111
	if p > 0 {
		bo.baseCost += x000
		bo._add_1(i, 0, x100-x000)
		bo._add_1(j, 0, x010-x000)
		bo._add_1(k, 0, x001-x000)
		bo._add_2(i, j, 0, 0, 0, x000+x110-x100-x010)
		bo._add_2(i, k, 0, 0, 0, x000+x101-x100-x001)
		bo._add_2(j, k, 0, 0, 0, x000+x011-x010-x001)
		bo.baseCost -= p
		bo._addEdge(i, bo.next, p)
		bo._addEdge(j, bo.next, p)
		bo._addEdge(k, bo.next, p)
		bo._addEdge(bo.next, bo.sink, p)
		bo.next++
	} else {
		p = -p
		bo.baseCost += x111
		bo._add_1(i, x011-x111, 0)
		bo._add_1(i, x101-x111, 0)
		bo._add_1(i, x110-x111, 0)
		bo._add_2(i, j, x111+x001-x011-x101, 0, 0, 0)
		bo._add_2(i, k, x111+x010-x011-x110, 0, 0, 0)
		bo._add_2(j, k, x111+x100-x101-x110, 0, 0, 0)
		bo.baseCost -= p
		bo._addEdge(bo.next, i, p)
		bo._addEdge(bo.next, j, p)
		bo._addEdge(bo.next, k, p)
		bo._addEdge(bo.source, bo.next, p)
		bo.next++
	}
}

// t>=0
func (bo *BinaryOptimization) _addEdge(i, j, t int) {
	if t == 0 {
		return
	}
	key := [2]int{i, j}
	bo.edges[key] += t
	bo.edges[key] = min(bo.edges[key], INF)
}

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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
