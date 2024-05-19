// MaxAntiChain
// OnAntiChain

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// abc237ex()
	P4298祭祀()

}

// P4298 [CTSC2008] 祭祀
// https://www.luogu.com.cn/problem/P4298
// 给定一个 n 个点，m 条边的简单有向无环图（DAG).
// 0.求最大独立集(反链)的大小。
// 1.求出它的最长反链(最大独立集)，并构造方案。
// 2.判断每个点是否在最长反链上.
func P4298祭祀() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = [2]int{u, v}
	}

	S := NewMaxAntiChainSolver(n, edges)
	maxAntiChain := S.MaxAntiChain()
	fmt.Fprintln(out, len(maxAntiChain))

	bit := make([]byte, n)
	for _, v := range maxAntiChain {
		bit[v] = 1
	}
	for i := 0; i < n; i++ {
		fmt.Fprint(out, bit[i])
	}
	fmt.Fprintln(out)

	onAntiChain := S.OnAntiChain()
	for i := 0; i < n; i++ {
		if onAntiChain[i] {
			fmt.Fprint(out, 1)
		} else {
			fmt.Fprint(out, 0)
		}
	}
	fmt.Fprintln(out)
}

// Ex - Hakata
// https://atcoder.jp/contests/abc237/tasks/abc237_h
// 给定一个字符串, 你需要从中选出若干回文子串, 并且使得选出的串不存在某一个是另一个的子串, 问最多能选出多少子串.
// n<=200
//
// 给定一些偏序包含关系,求最大独立集(互相无法到达).
// !遍历所有回文子串，如果j是i的子串，则连边 i->j，求dag最长反链即可.
func abc237ex() {
	zAlgo := func(s string) []int {
		n := len(s)
		if n == 0 {
			return nil
		}
		z := make([]int, n)
		j := 0
		for i := 1; i < n; i++ {
			var k int
			if j+z[j] <= i {
				k = 0
			} else {
				k = min(j+z[j]-i, z[i-j])
			}
			for i+k < n && s[k] == s[i+k] {
				k++
			}
			if j+z[j] < i+z[i] {
				j = i
			}
			z[i] = k
		}
		z[0] = n
		return z
	}

	// O(n+m)判断`shorter`是否是`longer`的子串.
	isSubstring := func(longer, shorter string) bool {
		if len(shorter) > len(longer) {
			return false
		}
		if len(shorter) == 0 {
			return true
		}
		n, m := len(longer), len(shorter)
		z := zAlgo(shorter + longer)
		for i := m; i < n+m; i++ {
			if z[i] >= m {
				return true
			}
		}
		return false
	}

	isPalindrome := func(s string) bool {
		n := len(s)
		for i := 0; i < n>>1; i++ {
			if s[i] != s[n-1-i] {
				return false
			}
		}
		return true
	}

	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	set := make(map[string]struct{})
	n := len(s)
	for start := 0; start < n; start++ {
		for end := start + 1; end <= n; end++ {
			cur := s[start:end]
			if isPalindrome(cur) {
				set[cur] = struct{}{}
			}
		}
	}
	allPalindromes := make([]string, 0, len(set))
	for k := range set {
		allPalindromes = append(allPalindromes, k)
	}

	edges := make([][2]int, 0)
	for i := 0; i < len(allPalindromes); i++ {
		for j := 0; j < len(allPalindromes); j++ {
			if i == j {
				continue
			}
			if isSubstring(allPalindromes[i], allPalindromes[j]) {
				edges = append(edges, [2]int{i, j})
			}
		}
	}

	S := NewMaxAntiChainSolver(len(allPalindromes), edges)
	res := S.MaxAntiChain()
	fmt.Fprintln(out, len(res))
}

// TODO: 增加权重
// dag上的最大带权独立集
// https://atcoder.jp/contests/abc354/tasks/abc354_g

const INF int = 1e18

type MaxAntiChainSolver struct {
	n            int
	source, sink int
	maxFlow      *MaxFlowAtcoder
	calculated   bool
}

func NewMaxAntiChainSolver(n int, edges [][2]int) *MaxAntiChainSolver {
	source, sink := 2*n, 2*n+1
	mf := NewMaxFlowAtcoder(2*n + 2)
	for _, e := range edges {
		u, v := e[0], e[1]
		mf.AddEdge(u+n, v, INF)
	}
	for i := 0; i < n; i++ {
		mf.AddEdge(source, i+n, 1)
		mf.AddEdge(i, sink, 1)
		mf.AddEdge(i, i+n, INF)
	}
	return &MaxAntiChainSolver{n: n, source: source, sink: sink, maxFlow: mf}
}

// dag最长反链.
func (solver *MaxAntiChainSolver) MaxAntiChain() []int {
	solver.calFlow()
	isCut := solver.maxFlow.MinCut(solver.source)
	res := make([]int, 0)
	for i := 0; i < solver.n; i++ {
		if !isCut[i] && isCut[i+solver.n] {
			res = append(res, i)
		}
	}
	return res
}

// 每个点能否在最长反链上.
func (solver *MaxAntiChainSolver) OnAntiChain() []bool {
	solver.calFlow()
	adjList := make([][]int, solver.n*2+2)
	for _, e := range solver.maxFlow.GetEdges() {
		if e.cap != e.flow {
			adjList[e.from] = append(adjList[e.from], e.to)
		}
		if e.flow != 0 {
			adjList[e.to] = append(adjList[e.to], e.from)
		}
	}

	reach := make([][]bool, solver.n*2+2)
	for i := 0; i < solver.n*2+2; i++ {
		reach[i] = make([]bool, solver.n*2+2)
	}
	for i := 0; i < solver.n*2+2; i++ {
		queue := []int{i}
		reach[i][i] = true
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, to := range adjList[cur] {
				if !reach[i][to] {
					reach[i][to] = true
					queue = append(queue, to)
				}
			}
		}
	}

	res := make([]bool, solver.n)
	for i := 0; i < solver.n; i++ {
		res[i] = !reach[i+solver.n][solver.sink] && !reach[i+solver.n][i] && !reach[solver.source][i]
	}
	return res
}

func (solver *MaxAntiChainSolver) calFlow() {
	if solver.calculated {
		return
	}
	solver.calculated = true
	solver.maxFlow.Flow(solver.source, solver.sink)
}

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
func (mf *MaxFlowAtcoder) MinCut(s int) (visited []bool) {
	visited = make([]bool, mf.n)
	q := []int32{int32(s)}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		visited[p] = true
		for _, e := range mf.g[p] {
			if e.cap != 0 && !visited[e.to] {
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
