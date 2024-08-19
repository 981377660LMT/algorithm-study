// 最大流.
// api:
//  NewMaxFlow(n, source, sink int32) *MaxFlow
//  AddEdge(from, to int32, cap int)
//  Flow() int
//  Cut() (int, []bool)
//  PathDecomposition() [][]int32
//  ChangeCapacity(i int32, after int)
//  Debug()
//  GetFlowEdges() []flowEdge

package main

import "fmt"

func main() {

}

const INF int = 1e18

// 支持动态加边，支持修改边容量.
// 记变化的容量为 F，时间复杂度为 O((N+M)|F|).
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

// 修改边的容量.
func (mf *MaxFlow) ChangeCapacity(i int32, after int) {
	from, idx := mf.pos[i][0], mf.pos[i][1]
	e := &mf.edges[from][idx]
	before := e.cap + e.flow
	if before < after {
		mf.caculated = e.cap > 0
		e.cap += after - before
		return
	}
	e.cap = after - e.flow
	if e.cap < 0 {
		mf.FlowPushBack(e)
	}
}

func (mf *MaxFlow) FlowPushBack(e0 *edge) {
	re0 := &mf.edges[e0.to][e0.rev]
	a, b := re0.to, e0.to
	findPath := func(s, t int32, lim int) int {
		vis := make([]bool, mf.n)
		for i := range mf.prog {
			mf.prog[i] = 0
		}
		var dfs func(int32, int) int
		dfs = func(v int32, f int) int {
			if v == t {
				return f
			}
			for i := &mf.prog[v]; *i < int32(len(mf.edges[v])); *i++ {
				e := &mf.edges[v][*i]
				if vis[e.to] || e.cap <= 0 {
					continue
				}
				vis[e.to] = true
				a := dfs(e.to, min(f, e.cap))
				if a == 0 {
					continue
				}
				e.cap -= a
				e.flow += a
				mf.edges[e.to][e.rev].cap += a
				mf.edges[e.to][e.rev].flow -= a
				return a
			}
			return 0
		}
		return dfs(s, lim)
	}

	for e0.cap < 0 {
		c := findPath(a, b, -e0.cap)
		if c == 0 {
			break
		}
		e0.cap += c
		e0.flow -= c
		re0.cap -= c
		re0.flow += c
	}
	c := -e0.cap
	for c > 0 && a != mf.source {
		x := findPath(a, mf.source, c)
		c -= x
	}
	c = -e0.cap
	for c > 0 && b != mf.sink {
		x := findPath(mf.sink, b, c)
		c -= x
	}
	c = -e0.cap
	e0.cap += c
	e0.flow -= c
	re0.cap -= c
	re0.flow += c
	mf.flowRes -= c
}

func (mf *MaxFlow) GetFlowEdges() []flowEdge {
	res := make([]flowEdge, 0)
	for frm, edges := range mf.edges {
		for _, e := range edges {
			if e.flow <= 0 {
				continue
			}
			res = append(res, flowEdge{int32(frm), e.to, e.flow})
		}
	}
	return res
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

// O(F(N+M)) 还原路径.
func (mf *MaxFlow) PathDecomposition() [][]int32 {
	mf.Flow()
	tos := make([][]int32, mf.n)
	edges := mf.GetFlowEdges()
	for i := 0; i < len(edges); i++ {
		from, to, flow := edges[i].from, edges[i].to, edges[i].flow
		for j := 0; j < flow; j++ {
			tos[from] = append(tos[from], to)
		}
	}

	var res [][]int32
	visited := make([]bool, mf.n)
	for i := 0; i < mf.flowRes; i++ {
		path := []int32{mf.source}
		visited[mf.source] = true
		for path[len(path)-1] != mf.sink {
			last := &tos[path[len(path)-1]]
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

func (mf *MaxFlow) Debug() {
	fmt.Println("source", mf.source)
	fmt.Println("sink", mf.sink)
	fmt.Println("edges (frm, to, cap, flow)")
	for i, edges := range mf.edges {
		for _, e := range edges {
			if e.cap == 0 && e.flow == 0 {
				continue
			}
			fmt.Println(i, e.to, e.cap, e.flow)
		}
	}
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

type flowEdge = struct {
	from, to int32
	flow     int
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
