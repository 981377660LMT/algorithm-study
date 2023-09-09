// https://ei1333.github.io/library/graph/flow/push-relabel.hpp
// 最大流预流推进

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_6_A
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	mf := NewPushRelabel(n)
	for i := 0; i < m; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		mf.AddEdge(u, v, c, -1)
	}

	fmt.Fprintln(out, mf.MaxFlow(0, n-1))
}

// 2123. 使矩阵中的 1 互不相邻的最小操作数
// https://leetcode.cn/problems/minimum-operations-to-remove-adjacent-ones-in-matrix/
func minimumOperations(grid [][]int) int {
	DIR4 := [][]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	ROW, COL := len(grid), len(grid[0])
	N := ROW * COL
	START := N
	END := START + 1
	maxFlow := NewPushRelabel(END + 1)

	for r := 0; r < ROW; r++ {
		for c := 0; c < COL; c++ {
			if grid[r][c] == 0 || (r+c)&1 == 1 {
				continue
			}
			cur := r*COL + c
			for _, d := range DIR4 {
				nr, nc := r+d[0], c+d[1]
				if nr >= 0 && nr < ROW && nc >= 0 && nc < COL && grid[nr][nc] == 1 {
					next := nr*COL + nc
					maxFlow.AddEdge(START, cur, 1, -1)
					maxFlow.AddEdge(cur, next, 1, -1)
					maxFlow.AddEdge(next, END, 1, -1)
				}
			}
		}
	}

	return maxFlow.MaxFlow(START, END)
}

// LCP 38. 守卫城堡
// https://leetcode.cn/problems/7rLGCR/description/
func guardCastle(grid []string) int {
	DIR2 := [][]int{{0, 1}, {1, 0}}
	ROW, COL := len(grid), len(grid[0])
	OFFSET := ROW * COL
	TELEPORT := 2 * OFFSET
	START := TELEPORT + 1
	END := START + 1
	maxFlow := NewPushRelabel(END + 1)

	getNext := func(cur int) []int {
		curRow, curCol := cur/COL, cur%COL
		res := []int{}
		for _, d := range DIR2 {
			nextRow, nextCol := curRow+d[0], curCol+d[1]
			if nextRow >= 0 && nextRow < ROW && nextCol >= 0 && nextCol < COL && grid[nextRow][nextCol] != '#' {
				res = append(res, nextRow*COL+nextCol)
			}
		}
		return res
	}

	for r := 0; r < ROW; r++ {
		for c := 0; c < COL; c++ {
			if grid[r][c] == '#' {
				continue
			}
			cur := r*COL + c

			// 所有点拆成 入点 和 出点 两个点
			if grid[r][c] == '.' {
				maxFlow.AddEdge(cur, cur+OFFSET, 1, -1)
			} else {
				maxFlow.AddEdge(cur, cur+OFFSET, INF, -1)
			}

			// 源点连接恶魔出生点
			if grid[r][c] == 'S' {
				maxFlow.AddEdge(START, cur, INF, -1)
			}

			// 城堡连接汇点
			if grid[r][c] == 'C' {
				maxFlow.AddEdge(cur, END, INF, -1)
			}

			// 虚拟点连通所有传送门
			if grid[r][c] == 'P' {
				maxFlow.AddEdge(cur+OFFSET, TELEPORT, INF, -1)
				maxFlow.AddEdge(TELEPORT, cur, INF, -1)
			}

			// 所有出点连通周围的入点
			for _, next := range getNext(cur) {
				maxFlow.AddEdge(cur+OFFSET, next, INF, -1)
				maxFlow.AddEdge(next+OFFSET, cur, INF, -1)
			}
		}
	}

	minCut := maxFlow.MaxFlow(START, END)
	if minCut < INF {
		return minCut
	}
	return -1
}

const INF int = 2e9

type edge struct {
	to    int32
	cap   int
	rev   int32
	isrev bool
	idx   int32
}

// 最高标号预流推进算法.
type MaxFlowPushRelabel struct {
	exFlow       []int
	graph        [][]edge
	potential    []int32
	curEdge      []int32
	allVertex    *_List
	activeVertex *_Stack
	visitedEdge  map[int]struct{}
	n            int32
	height       int32
	relabels     int32
}

func NewPushRelabel(n int) *MaxFlowPushRelabel {
	return &MaxFlowPushRelabel{
		exFlow:       make([]int, n),
		graph:        make([][]edge, n),
		potential:    make([]int32, n),
		curEdge:      make([]int32, n),
		allVertex:    _NewList(int32(n), int32(n)),
		activeVertex: _NewStack(int32(n), int32(n)),
		visitedEdge:  make(map[int]struct{}),
		n:            int32(n),
		height:       -1,
	}
}

// 内部会对边去重.
func (pr *MaxFlowPushRelabel) AddEdge(from, to, cap int, index int) {
	hash := from*int(pr.n) + to
	if _, ok := pr.visitedEdge[hash]; ok {
		return
	}
	pr.visitedEdge[hash] = struct{}{}
	pr.graph[from] = append(pr.graph[from], edge{int32(to), cap, int32(len(pr.graph[to])), false, int32(index)})
	pr.graph[to] = append(pr.graph[to], edge{int32(from), 0, int32(len(pr.graph[from]) - 1), true, int32(index)})
}

func (pr *MaxFlowPushRelabel) MaxFlow(start, target int) int {
	startInt32, targetInt32 := int32(start), int32(target)
	level := pr._init(startInt32, targetInt32)
	for level >= 0 {
		if pr.activeVertex.Empty(level) {
			level--
			continue
		}
		u := pr.activeVertex.Top(level)
		pr.activeVertex.Pop(level)
		level = pr._discharge(u, targetInt32)
		if pr.relabels*2 >= pr.n {
			level = pr._globalRelabel(targetInt32)
			pr.relabels = 0
		}
	}
	return pr.exFlow[target]
}

func (pr *MaxFlowPushRelabel) _calcActive(t int32) int32 {
	pr.height = -1
	for i := int32(0); i < pr.n; i++ {
		if pr.potential[i] < pr.n {
			pr.curEdge[i] = 0
			pr.height = maxInt32(pr.height, pr.potential[i])
			pr.allVertex.Insert(pr.potential[i], i)
			if pr.exFlow[i] > 0 && i != t {
				pr.activeVertex.Push(pr.potential[i], i)
			}
		} else {
			pr.potential[i] = pr.n + 1
		}
	}
	return pr.height
}

func (pr *MaxFlowPushRelabel) _bfs(t int32) {
	for i := int32(0); i < pr.n; i++ {
		pr.potential[i] = maxInt32(pr.potential[i], pr.n)
	}
	pr.potential[t] = 0
	que := []int32{t}
	for len(que) > 0 {
		p := que[0]
		que = que[1:]
		for _, e := range pr.graph[p] {
			if pr.potential[e.to] == pr.n && pr.graph[e.to][e.rev].cap > 0 {
				pr.potential[e.to] = pr.potential[p] + 1
				que = append(que, e.to)
			}
		}
	}
}

func (pr *MaxFlowPushRelabel) _init(s, t int32) int32 {
	pr.potential[s] = pr.n + 1
	pr._bfs(t)
	for i := range pr.graph[s] {
		e := &pr.graph[s][i]
		if pr.potential[e.to] < pr.n {
			pr.graph[e.to][e.rev].cap = e.cap
			pr.exFlow[s] -= e.cap
			pr.exFlow[e.to] += e.cap
		}
		e.cap = 0
	}
	return pr._calcActive(t)
}

func (pr *MaxFlowPushRelabel) _push(u, t int32, e *edge) bool {
	f := min(int(e.cap), pr.exFlow[u])
	v := e.to
	e.cap -= f
	pr.exFlow[u] -= f
	pr.graph[v][e.rev].cap += f
	pr.exFlow[v] += f
	if pr.exFlow[v] == f && v != t {
		pr.activeVertex.Push(pr.potential[v], v)
	}
	return pr.exFlow[u] == 0
}

func (pr *MaxFlowPushRelabel) _discharge(u, t int32) int32 {
	for i := &pr.curEdge[u]; *i < int32(len(pr.graph[u])); *i++ {
		e := &pr.graph[u][*i]
		if pr.potential[u] == pr.potential[e.to]+1 && e.cap > 0 {
			if pr._push(u, t, e) {
				return pr.potential[u]
			}
		}
	}
	return pr._relabel(u)
}

func (pr *MaxFlowPushRelabel) _globalRelabel(t int32) int32 {
	pr._bfs(t)
	pr.allVertex.Clear()
	pr.activeVertex.Clear()
	return pr._calcActive(t)
}

func (pr *MaxFlowPushRelabel) _gapRelabel(u int32) {
	for i := pr.potential[u]; i <= pr.height; i++ {
		for id := pr.allVertex.data[pr.n+i].next; id < pr.n; id = pr.allVertex.data[id].next {
			pr.potential[id] = pr.n + 1
		}
		pr.allVertex.data[pr.n+i].next = pr.n + i
		pr.allVertex.data[pr.n+i].prev = pr.n + i
	}
}

func (pr *MaxFlowPushRelabel) _relabel(u int32) int32 {
	pr.relabels++
	prv, cur := pr.potential[u], pr.n
	for i := int32(0); i < int32(len(pr.graph[u])); i++ {
		e := pr.graph[u][i]
		if cur > pr.potential[e.to]+1 && e.cap > 0 {
			pr.curEdge[u] = i
			cur = pr.potential[e.to] + 1
		}
	}
	if pr.allVertex.MoreOne(prv) {
		pr.allVertex.Erase(u)
		if pr.potential[u] = cur; pr.potential[u] == pr.n {
			pr.potential[u] = pr.n + 1
			return prv
		}
		pr.activeVertex.Push(cur, u)
		pr.allVertex.Insert(cur, u)
		pr.height = maxInt32(pr.height, cur)
	} else {
		pr._gapRelabel(u)
		pr.height = prv - 1
		return pr.height
	}
	return cur
}

type _Stack struct {
	n, h int32
	node []int32
}

func _NewStack(n, h int32) *_Stack {
	res := &_Stack{n: n, h: h, node: make([]int32, n+h)}
	res.Clear()
	return res
}

func (list *_Stack) Empty(h int32) bool {
	return list.node[list.n+h] == list.n+h
}

func (list *_Stack) Top(h int32) int32 {
	return list.node[list.n+h]
}

func (list *_Stack) Pop(h int32) {
	list.node[list.n+h] = list.node[list.node[list.n+h]]
}

func (list *_Stack) Push(h, u int32) {
	list.node[u] = list.node[list.n+h]
	list.node[list.n+h] = u
}

func (list *_Stack) Clear() {
	for i := list.n; i < list.n+list.h; i++ {
		list.node[i] = i
	}
}

type _node struct{ prev, next int32 }

type _List struct {
	n, h int32
	data []_node
}

func _NewList(n, h int32) *_List {
	res := &_List{n: n, h: h, data: make([]_node, n+h)}
	res.Clear()
	return res
}

func (list *_List) Empty(h int32) bool {
	return list.data[list.n+h].next == list.n+h
}

func (list *_List) MoreOne(h int32) bool {
	return list.data[list.n+h].prev != list.data[list.n+h].next
}

func (list *_List) Insert(h, u int32) {
	list.data[u].prev = list.data[list.n+h].prev
	list.data[u].next = list.n + h
	list.data[list.data[list.n+h].prev].next = u
	list.data[list.n+h].prev = u
}

func (list *_List) Erase(u int32) {
	list.data[list.data[u].prev].next = list.data[u].next
	list.data[list.data[u].next].prev = list.data[u].prev
}

func (list *_List) Clear() {
	for i := list.n; i < list.n+list.h; i++ {
		list.data[i].next = i
		list.data[i].prev = i
	}
}

func maxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
