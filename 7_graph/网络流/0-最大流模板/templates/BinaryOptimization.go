package main

import "fmt"

func main() {

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
func (bo *BinaryOptimization) Run() (int, []bool) {
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
	mf.edges[from] = append(mf.edges[from], edge{to, b, cap, 0})
	mf.edges[to] = append(mf.edges[to], edge{from, a, 0, 0})
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
			f := mf.flowDfs(mf.sink, INF)
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
