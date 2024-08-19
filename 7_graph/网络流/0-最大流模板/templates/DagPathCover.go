package main

import (
	"fmt"
)

func demo() {
	// 0
	// | \
	// 1  2
	// |  |
	// 3  4

	dag := [][]int32{{1, 2}, {3}, {4}, {}, {}}
	colors := DagPathCover(dag)
	fmt.Println(colors) // [0 0 1 0 1]
}

const INF int = 1e18

// dag路径覆盖(最小不相交路径覆盖).
// 给每个点染色，使得每种颜色的点构成一条路径.
// !dag的边i->j要求i<j.
func DagPathCover(dag [][]int32) (colors []int32) {
	n := int32(len(dag))
	for from := int32(0); from < n; from++ {
		for _, to := range dag[from] {
			if from >= to {
				panic("edge i->j must satisfy i<j")
			}
		}
	}

	source, sink := 2*n, 2*n+1
	mf := newMaxFlow(2*n+2, source, sink)
	for v := int32(0); v < n; v++ {
		mf.AddEdge(source, 2*v+1, 1)
		mf.AddEdge(2*v+0, sink, 1)
		mf.AddEdge(2*v+0, 2*v+1, INF)
	}
	for from := int32(0); from < n; from++ {
		for _, to := range dag[from] {
			mf.AddEdge(2*from+1, 2*to+0, INF)
		}
	}

	mf.Flow()
	paths := mf.PathDecomposition()

	uf := newUnionFindArraySimple32(n)
	for _, path := range paths {
		a, b := path[1], path[len(path)-2]
		uf.Union(a/2, b/2)
	}

	colors = make([]int32, n)
	for v := int32(0); v < n; v++ {
		colors[v] = -1
	}
	p := int32(0)
	for v := int32(0); v < n; v++ {
		if uf.Find(v) == v {
			colors[v] = p
			p++
		}
	}
	for v := int32(0); v < n; v++ {
		if root := uf.Find(v); root != v {
			colors[v] = colors[root]
		}
	}
	return
}

type unionFindArraySimple32 struct {
	n    int32
	data []int32
}

func newUnionFindArraySimple32(n int32) *unionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &unionFindArraySimple32{n: n, data: data}
}

func (u *unionFindArraySimple32) Union(key1, key2 int32) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	return true
}

func (u *unionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *unionFindArraySimple32) Size(key int32) int32 {
	return -u.data[u.Find(key)]
}

// 支持动态加边，支持修改边容量.
// 记变化的容量为 F，时间复杂度为 O((N+M)|F|).
type maxFlow struct {
	caculated       bool
	n, source, sink int32
	flowRes         int
	prog, level     []int32
	que             []int32
	pos             [][2]int32
	edges           [][]edge
}

func newMaxFlow(n, source, sink int32) *maxFlow {
	return &maxFlow{
		n:      n,
		source: source,
		sink:   sink,
		prog:   make([]int32, n),
		level:  make([]int32, n),
		que:    make([]int32, n),
		edges:  make([][]edge, n),
	}
}

func (mf *maxFlow) AddEdge(from, to int32, cap int) {
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

func (mf *maxFlow) GetFlowEdges() []flowEdge {
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

func (mf *maxFlow) Flow() int {
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

// O(F(N+M)) 还原路径.
func (mf *maxFlow) PathDecomposition() [][]int32 {
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

func (mf *maxFlow) setLevel() bool {
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

func (mf *maxFlow) flowDfs(v int32, lim int) int {
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
