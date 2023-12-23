package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_6_A
func demo() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	F := NewMaxFlow(n, 0, n-1)
	for i := 0; i < m; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		F.Add(a, b, c, 0)
	}
	fmt.Println(F.Flow())
}

const INF int = 1e18

type Edge struct {
	to, rev   int32
	cap, flow int
}

type MaxFlow struct {
	n, source, sink int32
	edges           [][]Edge
	level, iter     []int32
	queue           []int32
	calculated      bool
	flowRes         int
}

func NewMaxFlow(n, source, sink int) *MaxFlow {
	return &MaxFlow{
		n:      int32(n),
		source: int32(source),
		sink:   int32(sink),
		edges:  make([][]Edge, n),
	}
}

func (mf *MaxFlow) Add(from, to int, cap, revCap int) {
	if from == to {
		return
	}
	mf.calculated = false
	a := len(mf.edges[from])
	b := len(mf.edges[to])
	mf.edges[from] = append(mf.edges[from], Edge{to: int32(to), rev: int32(b), cap: cap})
	mf.edges[to] = append(mf.edges[to], Edge{to: int32(from), rev: int32(a), cap: revCap})
}

// from, to, flow
func (mf *MaxFlow) GetFlowEdges() [][3]int {
	edges := make([][3]int, 0)
	for from := int32(0); from < mf.n; from++ {
		nexts := mf.edges[from]
		for i := 0; i < len(nexts); i++ {
			e := &nexts[i]
			if e.flow <= 0 {
				continue
			}
			edges = append(edges, [3]int{int(from), int(e.to), e.flow})
		}
	}
	return edges
}

func (mf *MaxFlow) Flow() int {
	if mf.calculated {
		return mf.flowRes
	}
	mf.calculated = true
	for mf.setLevel() {
		mf.iter = make([]int32, mf.n)
		for {
			x := mf.flowDfs(mf.source, INF)
			if x == 0 {
				break
			}
			mf.flowRes += x
			mf.flowRes = min(mf.flowRes, INF)
			if mf.flowRes == INF {
				return mf.flowRes
			}
		}
	}
	return mf.flowRes
}

// 最小割(边).返回最小割的值和每个点是否在最小割中.
func (mf *MaxFlow) MinCut() (res int, isCut []bool) {
	mf.Flow()
	isCut = make([]bool, mf.n)
	for i := int32(0); i < mf.n; i++ {
		isCut[i] = mf.level[i] < 0
	}
	return mf.flowRes, isCut
}

// O(流量*(边数+点数)) 还原出所有的简单路径.
func (mf *MaxFlow) PathDecomposition() [][]int {
	mf.Flow()
	edges := mf.GetFlowEdges()
	tos := make([][]int32, mf.n)
	for i := 0; i < len(edges); i++ {
		e := &edges[i]
		from, to, flow := e[0], int32(e[1]), e[2]
		for j := 0; j < flow; j++ {
			tos[from] = append(tos[from], to)
		}
	}

	res := make([][]int, 0)
	visited := make([]bool, mf.n)
	sinkInt := int(mf.sink)
	for i := 0; i < mf.flowRes; i++ {
		path := []int{int(mf.source)}
		visited[mf.source] = true
		for path[len(path)-1] != sinkInt {
			last := path[len(path)-1]
			to := tos[last][len(tos[last])-1]
			tos[last] = tos[last][:len(tos[last])-1]
			for visited[to] {
				last := path[len(path)-1]
				path = path[:len(path)-1]
				visited[last] = false
			}
			path = append(path, int(to))
			visited[to] = true
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
	fmt.Println("edges (from, to, cap, flow)")
	for v := int32(0); v < mf.n; v++ {
		for _, e := range mf.edges[v] {
			if e.cap == 0 && e.flow == 0 {
				continue
			}
			fmt.Println(v, e.to, e.cap, e.flow)
		}
	}
}

func (mf *MaxFlow) setLevel() bool {
	mf.queue = []int32{}
	mf.level = make([]int32, mf.n)
	for i := int32(0); i < mf.n; i++ {
		mf.level[i] = -1
	}
	mf.level[mf.source] = 0
	mf.queue = append(mf.queue, int32(mf.source))
	for len(mf.queue) > 0 {
		cur := mf.queue[0]
		mf.queue = mf.queue[1:]
		nexts := mf.edges[cur]
		for i := 0; i < len(nexts); i++ {
			e := &nexts[i]
			if e.cap > 0 && mf.level[e.to] == -1 {
				mf.level[e.to] = mf.level[cur] + 1
				if e.to == mf.sink {
					return true
				}
				mf.queue = append(mf.queue, e.to)
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
	for i := &mf.iter[v]; *i < int32(len(mf.edges[v])); *i++ {
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

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
