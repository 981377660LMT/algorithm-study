package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

// from https://atcoder.jp/users/ccppjsrb
var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

// いくつかのボールがあります。 各ボールは色
// 1 、色
// 2 、
// … 、色
// N のうちのいずれかであり、
// i=1,2,…,N について、色
// i のボールは全部で
// A
// i
// ​
//   個あります。

// また、
// M 個の箱があります。
// j=1,2,…,M について、
// j 番目の箱には合計で
// B
// j
// ​
//   個までのボールを入れることができます。

// ただし、
// 1≤i≤N と
// 1≤j≤M を満たすすべての整数の組
// (i,j) について、 色
// i のボールを
// j 番目の箱に入れる個数は
// (i×j) 個以下でなければなりません。

// M 個の箱の中に入れることができるボールの合計個数の最大値を出力してください。

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, M := io.NextInt(), io.NextInt()
	ballCount := make([]int, N)
	for i := 0; i < N; i++ {
		ballCount[i] = io.NextInt()
	}
	boxCapacity := make([]int, M)
	for i := 0; i < M; i++ {
		boxCapacity[i] = io.NextInt()
	}

	START := N + M
	END := N + M + 1
	flow := NewPushRelabel(N + M + 2)
	for i := 0; i < N; i++ {
		flow.AddEdge(START, i, ballCount[i], -1)
	}
	for i := 0; i < M; i++ {
		flow.AddEdge(N+i, END, boxCapacity[i], -1)
	}
	for i := N - 1; i >= 0; i-- {
		for j := M - 1; j >= 0; j-- {
			flow.AddEdge(i, N+j, min(min(boxCapacity[j], (i+1)*(j+1)), ballCount[i]), -1)
		}
	}

	fmt.Println(flow.MaxFlow(START, END))
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
