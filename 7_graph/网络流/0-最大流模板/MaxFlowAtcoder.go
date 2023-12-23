package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	demo()
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_6_A
func demo() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	F := NewMaxFlowAtcoder(n)
	for i := 0; i < m; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		F.AddEdge(a, b, c)
	}
	fmt.Println(F.Flow(0, n-1))
}

const INF int = 1e18

type Edge struct{ from, to, cap, flow int }

type _edge struct {
	to, rev int32
	cap     int
}

type MaxFlowAtcoder struct {
	n   int32
	pos [][2]int32
	g   [][]_edge
}

// https://github.com/atcoder/ac-library/blob/master/atcoder/maxflow.hpp
func NewMaxFlowAtcoder(n int) *MaxFlowAtcoder {
	return &MaxFlowAtcoder{
		n: int32(n),
		g: make([][]_edge, n),
	}
}

// 添加一条从from到to的容量为cap的边，返回边的编号.
func (mf *MaxFlowAtcoder) AddEdge(from, to int, cap int) int {
	m := len(mf.pos)
	mf.pos = append(mf.pos, [2]int32{int32(from), int32(len(mf.g[from]))})
	fromId := int32(len(mf.g[from]))
	toId := int32(len(mf.g[to]))
	if from == to {
		toId++
	}
	mf.g[from] = append(mf.g[from], _edge{int32(to), toId, cap})
	mf.g[to] = append(mf.g[to], _edge{int32(from), fromId, 0})
	return m
}

func (mf *MaxFlowAtcoder) GetEdge(i int) Edge {
	first, second := mf.pos[i][0], mf.pos[i][1]
	e := mf.g[first][second]
	re := mf.g[e.to][e.rev]
	return Edge{int(first), int(e.to), e.cap + re.cap, re.cap}
}

func (mf *MaxFlowAtcoder) GetEdges() []Edge {
	m := len(mf.pos)
	res := make([]Edge, 0, m)
	for i := 0; i < m; i++ {
		res = append(res, mf.GetEdge(i))
	}
	return res
}

func (mf *MaxFlowAtcoder) ChangeEdge(i int, newCap, newFlow int) {
	e := &mf.g[mf.pos[i][0]][mf.pos[i][1]]
	re := &mf.g[e.to][e.rev]
	e.cap = newCap - newFlow
	re.cap = newFlow
}

func (mf *MaxFlowAtcoder) Flow(s, t int) int {
	return mf.FlowWithLimit(s, t, INF)
}

func (mf *MaxFlowAtcoder) FlowWithLimit(s, t int, flowLimit int) int {
	level := make([]int32, mf.n)
	iter := make([]int32, mf.n)
	queue := make([]int32, 0, mf.n)
	s32, t32 := int32(s), int32(t)

	bfs := func() {
		for i := range level {
			level[i] = -1
		}
		level[s] = 0
		queue = queue[:0]
		queue = append(queue, s32)
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for _, e := range mf.g[v] {
				if e.cap == 0 || level[e.to] >= 0 {
					continue
				}
				level[e.to] = level[v] + 1
				if e.to == t32 {
					return
				}
				queue = append(queue, e.to)
			}
		}
	}

	var dfs func(int32, int) int
	dfs = func(v int32, up int) int {
		if v == s32 {
			return up
		}
		res := 0
		levelV := level[v]
		for i := &iter[v]; *i < int32(len(mf.g[v])); *i++ {
			e := &mf.g[v][*i]
			if levelV <= level[e.to] || mf.g[e.to][e.rev].cap == 0 {
				continue
			}
			d := dfs(e.to, min(up-res, mf.g[e.to][e.rev].cap))
			if d <= 0 {
				continue
			}
			e.cap += d
			mf.g[e.to][e.rev].cap -= d
			res += d
			if res == up {
				return res
			}
		}
		level[v] = mf.n
		return res
	}

	flow := 0
	for flow < flowLimit {
		bfs()
		if level[t] == -1 {
			break
		}
		for i := range iter {
			iter[i] = 0
		}
		for {
			f := dfs(t32, flowLimit-flow)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

// 返回剩余网络中从源点s可以到达的所有节点的集合。
func (mf *MaxFlowAtcoder) MinCut(s int) []bool {
	visited := make([]bool, mf.n)
	q := []int32{int32(s)}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		visited[p] = true
		for _, e := range mf.g[p] {
			if e.cap > 0 && !visited[e.to] {
				visited[e.to] = true
				q = append(q, e.to)
			}
		}
	}
	return visited
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
