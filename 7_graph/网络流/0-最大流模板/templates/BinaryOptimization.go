// BinaryOptimization 模型(二元集合划分问题)
// 队伍划分/分组，使得某个值最大/最小
// 燃やす埋める問題, Project Selection Problem

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// yuki_1541()
	// yuki_2320()
	// abc193f()
	abc259g()
}

// https://yukicoder.me/problems/no/1541
// 期末考试
// 有n个考试科目,每学一个科目就能多拿base分
// 对于每个科目i,可以花费cost来学习，学习之后有额外的收益:
// 对于科目subjects[j],如果i和subjects[j]都学习了,那么就能多拿到bonus[j]分
// !最大化(总分-花费)
// n<=100
// !每个科目学习还是不学习 => 燃やす埋める
// 先学习所有科目,然后再割掉不学每个科目的代价(最小割<=>最小代价的划分方案)
func yuki_1541() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	var base int
	fmt.Fscan(in, &n, &base)
	costs := make([]int, n)
	subjects := make([][]int32, n)
	bonuses := make([][]int, n)
	for i := int32(0); i < n; i++ {
		var k, cost int
		fmt.Fscan(in, &k, &cost)
		costs[i] = cost
		subjects[i] = make([]int32, k)
		bonuses[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &subjects[i][j])
			subjects[i][j]--
		}
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &bonuses[i][j])
		}
	}

	B := NewBinaryOptimization(n, false)
	for i := int32(0); i < n; i++ {
		B.Add1(i, 0, base-costs[i])
		for j := 0; j < len(subjects[i]); j++ {
			B.Add2(i, subjects[i][j], 0, 0, 0, bonuses[i][j])
		}
	}
	res, _ := B.Calc()
	fmt.Fprintln(out, res)
}

// No.2320 Game World for PvP
// https://yukicoder.me/problems/no/2320
// 一共有n个人正在进行拔河比赛.
// A数组表示现在红队的人，B数组表示现在蓝队的人.
// 剩下的人加入红队或者蓝队都可以.
// C[i][j] 表示i和j一起在一队的收益.
// 最大化总收益.
func yuki_2320() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, S, T int32
	fmt.Fscan(in, &N, &S, &T)
	A := make([]int, S)
	B := make([]int, T)
	for i := int32(0); i < S; i++ {
		fmt.Fscan(in, &A[i])
		A[i]--
	}
	for i := int32(0); i < T; i++ {
		fmt.Fscan(in, &B[i])
		B[i]--
	}
	C := make([][]int, N)
	for i := int32(0); i < N; i++ {
		C[i] = make([]int, N)
		for j := int32(0); j < N; j++ {
			fmt.Fscan(in, &C[i][j])
		}
	}

	belong := make([]int8, N)
	for i := int32(0); i < N; i++ {
		belong[i] = -1
	}
	for _, v := range A {
		belong[v] = 0
	}
	for _, v := range B {
		belong[v] = 1
	}

	M := NewBinaryOptimization(N, false)
	for i := int32(0); i < N; i++ {
		if belong[i] == 0 {
			M.Add1(i, 0, -INF)
		}
		if belong[i] == 1 {
			M.Add1(i, -INF, 0)
		}
	}
	for i := int32(0); i < N; i++ {
		for j := i + 1; j < N; j++ {
			v := C[i][j]
			M.Add2(i, j, v, 0, 0, v)
		}
	}
	res, _ := M.Calc()
	fmt.Fprintln(out, res)
}

// F - Zebraness
// https://atcoder.jp/contests/abc193/tasks/abc193_f
// 给定一个n*n的矩阵，每个位置有一个值.
// B表示黑色，W表示白色，?表示未知.
// !现在要求填充所有的?为B或者W，最大化相邻位置颜色不同的数量.
// n<=100.
// https://kanpurin.hatenablog.com/entry/2021/02/27/225330
func abc193f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	grid := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &grid[i])
	}

	B := NewBinaryOptimization(n*n, false)
	idx := func(i, j int32) int32 { return i*n + j }
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < n; j++ {
			f := (i + j) & 1
			h := idx(i, j)
			if grid[i][j] == 'W' && f == 0 {
				B.Add1(h, 0, -INF)
			}
			if grid[i][j] == 'W' && f == 1 {
				B.Add1(h, -INF, 0)
			}
			if grid[i][j] == 'B' && f == 0 {
				B.Add1(h, -INF, 0)
			}
			if grid[i][j] == 'B' && f == 1 {
				B.Add1(h, 0, -INF)
			}
		}
	}

	dir2 := [][2]int32{{1, 0}, {0, 1}}
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < n; j++ {
			for _, d := range dir2 {
				ni, nj := i+d[0], j+d[1]
				if ni < 0 || ni >= n || nj < 0 || nj >= n {
					continue
				}
				h := idx(i, j)
				nh := idx(ni, nj)
				B.Add2(h, nh, 1, 0, 0, 1)
			}
		}
	}

	res, _ := B.Calc()
	fmt.Fprintln(out, res)
}

// G - Grid Card Game
// https://atcoder.jp/contests/abc259/tasks/abc259_g
// 给定一个矩阵Anxn (1≤ n ≤ 100)，选择一些行列，可以得到这些行列包含的位置的并的数值和。
// 此外要求任意选中的行列交点处不能是负数。
// !求选择的最大值
// !时间复杂度O(V^2E)
func abc259g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL int32
	fmt.Fscan(in, &ROW, &COL)
	grid := make([][]int, ROW)
	for i := int32(0); i < ROW; i++ {
		grid[i] = make([]int, COL)
		for j := int32(0); j < COL; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	const big int = 1 << 45
	bm := NewBinaryOptimization(ROW+COL, true)

	rowId := func(r int32) int32 { return r }
	colId := func(c int32) int32 { return ROW + c }
	for i := int32(0); i < ROW; i++ {
		for j := int32(0); j < COL; j++ {
			// https://atcoder.jp/contests/abc259/editorial/4286
			x := grid[i][j]
			if x > 0 {
				bm.Add2(rowId(i), colId(j), -x, 0, -x, -x)
			}
			if x < 0 {
				bm.Add2(rowId(i), colId(j), -x, 0, big, -x)
			}
		}
	}

	res, _ := bm.Calc()
	fmt.Fprintln(out, -res)
}

const INF int = 1e18

type BinaryOptimization struct {
	minimize     bool
	n            int32
	source, sink int32
	next         int32
	baseCost     int
	edges        map[[2]int32]int
}

func NewBinaryOptimization(n int32, minimize bool) *BinaryOptimization {
	return &BinaryOptimization{
		minimize: minimize,
		n:        n,
		source:   n,
		sink:     n + 1,
		next:     n + 2,
		edges:    make(map[[2]int32]int),
	}
}

// xi 属于 0/1 时对应的收益.
func (bo *BinaryOptimization) Add1(i int32, x0, x1 int) {
	assert(0 <= i && i < bo.n, "i out of range")
	if !bo.minimize {
		x0, x1 = -x0, -x1
	}
	bo._add_1(i, x0, x1)
}

// (xi,xj) = (00,01,10,11) 时对应的收益.
//
//	!最小化代价时需要满足 x00 + x11 <= x01 + x10.
//	!最大化收益时需要满足 x00 + x11 >= x01 + x10.
func (bo *BinaryOptimization) Add2(i, j int32, x00, x01, x10, x11 int) {
	assert(i != j, "i != j")
	assert(0 <= i && i < bo.n, "i out of range")
	assert(0 <= j && j < bo.n, "j out of range")
	if !bo.minimize {
		x00, x01, x10, x11 = -x00, -x01, -x10, -x11
	}
	bo._add_2(i, j, x00, x01, x10, x11)
}

// (xi,xj,xk) = (000,001,010,011,100,101,110,111) 时对应的收益.
func (bo *BinaryOptimization) Add3(i, j, k int32, x000, x001, x010, x011, x100, x101, x110, x111 int) {
	assert(i != j && i != k && j != k, "i != j && i != k && j != k")
	assert(0 <= i && i < bo.n, "i out of range")
	assert(0 <= j && j < bo.n, "j out of range")
	assert(0 <= k && k < bo.n, "k out of range")
	if !bo.minimize {
		x000, x001, x010, x011, x100, x101, x110, x111 = -x000, -x001, -x010, -x011, -x100, -x101, -x110, -x111
	}
	bo._add_3(i, j, k, x000, x001, x010, x011, x100, x101, x110, x111)
}

// 返回最大收益/最小花费和每个变量的取值0/1.
func (bo *BinaryOptimization) Calc() (int, []bool) {
	flow := NewMaxFlow(bo.next, bo.source, bo.sink)
	for key, cap := range bo.edges {
		from, to := key[0], key[1]
		flow.AddEdge(from, to, cap)
	}
	res, isCut := flow.Cut()
	res += bo.baseCost
	res = min(res, INF)
	if !bo.minimize {
		res = -res
	}
	assign := isCut[:bo.n]
	return res, assign
}

func (bo *BinaryOptimization) Debug() {
	fmt.Println("base_cost", bo.baseCost)
	fmt.Println("source=", bo.source, "sink=", bo.sink)
	for key, cap := range bo.edges {
		fmt.Println(key, cap)
	}
}

func (bo *BinaryOptimization) _add_1(i int32, x0, x1 int) {
	if x0 <= x1 {
		bo.baseCost += x0
		bo._addEdge(bo.source, i, x1-x0)
	} else {
		bo.baseCost += x1
		bo._addEdge(i, bo.sink, x0-x1)
	}
}

// x00 + x11 <= x01 + x10
func (bo *BinaryOptimization) _add_2(i, j int32, x00, x01, x10, x11 int) {
	if x00+x11 > x01+x10 {
		panic("need to satisfy `x00 + x11 <= x01 + x10`.")
	}
	bo._add_1(i, x00, x10)
	bo._add_1(j, 0, x11-x10)
	bo._addEdge(i, j, x01+x10-x00-x11)
}

func (bo *BinaryOptimization) _add_3(i, j, k int32, x000, x001, x010, x011, x100, x101, x110, x111 int) {
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
func (bo *BinaryOptimization) _addEdge(i, j int32, t int) {
	if t == 0 {
		return
	}
	key := [2]int32{i, j}
	bo.edges[key] += t
}

type MaxFlow struct {
	caculated       bool
	n, source, sink int32
	flowRes         int
	prog, level     []int32
	que             []int32
	pos             [][2]int32
	edges           [][]edge
}

func NewMaxFlow(n, source, sink int32) *MaxFlow {
	return &MaxFlow{
		n:      n,
		source: source,
		sink:   sink,
		prog:   make([]int32, n),
		level:  make([]int32, n),
		que:    make([]int32, n),
		edges:  make([][]edge, n),
	}
}

func (mf *MaxFlow) AddEdge(from, to int32, cap int) {
	mf.caculated = false
	if from < 0 || from >= mf.n {
		panic("from out of range")
	}
	if to < 0 || to >= mf.n {
		panic("to out of range")
	}
	if cap < 0 {
		panic("cap must be non-negative")
	}
	a := int32(len(mf.edges[from]))
	var b int32
	if from == to {
		b = a + 1
	} else {
		b = int32(len(mf.edges[to]))
	}
	mf.pos = append(mf.pos, [2]int32{from, a})
	mf.edges[from] = append(mf.edges[from], edge{to: to, rev: b, cap: cap, flow: 0})
	mf.edges[to] = append(mf.edges[to], edge{to: from, rev: a, cap: 0, flow: 0})
}

func (mf *MaxFlow) Flow() int {
	if mf.caculated {
		return mf.flowRes
	}
	mf.caculated = true
	for mf.setLevel() {
		for i := range mf.prog {
			mf.prog[i] = 0
		}
		for {
			f := mf.flowDfs(mf.source, INF)
			if f == 0 {
				break
			}
			mf.flowRes += f
			mf.flowRes = min(mf.flowRes, INF)
			if mf.flowRes == INF {
				return mf.flowRes
			}
		}
	}
	return mf.flowRes
}

// 返回最小割的值和每个点是否属于最小割.
func (mf *MaxFlow) Cut() (int, []bool) {
	mf.Flow()
	isCut := make([]bool, mf.n)
	for i, v := range mf.level {
		isCut[i] = v < 0
	}
	return mf.flowRes, isCut
}

func (mf *MaxFlow) setLevel() bool {
	for i := range mf.level {
		mf.level[i] = -1
	}
	mf.level[mf.source] = 0
	l, r := int32(0), int32(0)
	mf.que[r] = mf.source
	r++
	for l < r {
		v := mf.que[l]
		l++
		for _, e := range mf.edges[v] {
			if e.cap > 0 && mf.level[e.to] == -1 {
				mf.level[e.to] = mf.level[v] + 1
				if e.to == mf.sink {
					return true
				}
				mf.que[r] = e.to
				r++
			}
		}
	}
	return false
}

func (mf *MaxFlow) flowDfs(v int32, lim int) int {
	if v == mf.sink {
		return lim
	}
	res := 0
	for i := &mf.prog[v]; *i < int32(len(mf.edges[v])); *i++ {
		e := &mf.edges[v][*i]
		if e.cap > 0 && mf.level[e.to] == mf.level[v]+1 {
			a := mf.flowDfs(e.to, min(lim, e.cap))
			if a > 0 {
				e.cap -= a
				e.flow += a
				mf.edges[e.to][e.rev].cap += a
				mf.edges[e.to][e.rev].flow -= a
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

type edge = struct {
	to, rev   int32
	cap, flow int
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

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}
