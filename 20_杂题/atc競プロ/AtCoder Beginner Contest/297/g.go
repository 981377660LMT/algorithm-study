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

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, L, R := io.NextInt(), io.NextInt(), io.NextInt()
	stones := make([]int, n)
	for i := 0; i < n; i++ {
		stones[i] = io.NextInt()
	}

	grundy := GrundyNumber(dag)
	xor := 0
	for _, num := range stones {
		xor ^= grundy[num]
	}
	if xor == 0 {
		fmt.Println("Second")
	} else {
		fmt.Println("First")
	}
}

func GrundyNumber(dag [][]int) (grundy []int) {
	order, ok := topoSort(dag)
	if !ok {
		return
	}

	grundy = make([]int, len(dag))
	memo := make([]int, len(dag)+1)
	for j := len(order) - 1; j >= 0; j-- {
		i := order[j]
		if len(dag[i]) == 0 {
			continue
		}
		for _, v := range dag[i] {
			memo[grundy[v]]++
		}
		for memo[grundy[i]] > 0 {
			grundy[i]++
		}
		for _, v := range dag[i] {
			memo[grundy[v]]--
		}
	}

	return
}

func topoSort(dag [][]int) (order []int, ok bool) {
	n := len(dag)
	visited, temp := make([]bool, n), make([]bool, n)
	var dfs func(int) bool
	dfs = func(i int) bool {
		if temp[i] {
			return false
		}
		if !visited[i] {
			temp[i] = true
			for _, v := range dag[i] {
				if !dfs(v) {
					return false
				}
			}
			visited[i] = true
			order = append(order, i)
			temp[i] = false
		}
		return true
	}

	for i := 0; i < n; i++ {
		if !visited[i] {
			if !dfs(i) {
				return nil, false
			}
		}
	}

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order, true
}

type RangeToRangeGraph struct {
	n     int
	nNode int
	edges [][3]int // [from, to, weight]
}

func NewRangeToRangeGraph(n int) *RangeToRangeGraph {
	g := &RangeToRangeGraph{
		n:     n,
		nNode: n * 3,
	}
	for i := 2; i < n+n; i++ {
		g.edges = append(g.edges, [3]int{g.toUpperIdx(i / 2), g.toUpperIdx(i), 0})
	}
	for i := 2; i < n+n; i++ {
		g.edges = append(g.edges, [3]int{g.toLowerIdx(i), g.toLowerIdx(i / 2), 0})
	}
	return g
}

// 添加有向边 from -> to, 权重为 weight.
func (g *RangeToRangeGraph) Add(from, to int, weight int) {
	g.edges = append(g.edges, [3]int{from, to, weight})
}

// 从区间 [fromStart, fromEnd) 中的每个点到 to 都添加一条有向边，权重为 weight.
func (g *RangeToRangeGraph) AddFromRange(fromStart, fromEnd, to int, weight int) {
	l, r := fromStart+g.n, fromEnd+g.n
	for l < r {
		if l&1 == 1 {
			g.Add(g.toLowerIdx(l), to, weight)
			l++
		}
		if r&1 == 1 {
			r--
			g.Add(g.toLowerIdx(r), to, weight)
		}
		l >>= 1
		r >>= 1
	}
}

// 从 from 到区间 [toStart, toEnd) 中的每个点都添加一条有向边，权重为 weight.
func (g *RangeToRangeGraph) AddToRange(from, toStart, toEnd int, weight int) {
	l, r := toStart+g.n, toEnd+g.n
	for l < r {
		if l&1 == 1 {
			g.Add(from, g.toUpperIdx(l), weight)
			l++
		}
		if r&1 == 1 {
			r--
			g.Add(from, g.toUpperIdx(r), weight)
		}
		l >>= 1
		r >>= 1
	}
}

// 从区间 [fromStart, fromEnd) 中的每个点到区间 [toStart, toEnd) 中的每个点都添加一条有向边，权重为 weight.
func (g *RangeToRangeGraph) AddRangeToRange(fromStart, fromEnd, toStart, toEnd int, weight int) {
	newNode := g.nNode
	g.nNode++
	g.AddFromRange(fromStart, fromEnd, newNode, weight)
	g.AddToRange(newNode, toStart, toEnd, 0)
}

// 返回`新图的有向邻接表和新图的节点数`.
func (g *RangeToRangeGraph) Build() (graph [][][2]int, vertex int) {
	graph = make([][][2]int, g.nNode)
	for _, e := range g.edges {
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], [2]int{v, w})
	}
	return graph, g.nNode
}

func (g *RangeToRangeGraph) toUpperIdx(i int) int {
	if i >= g.n {
		return i - g.n
	}
	return g.n + i
}

func (g *RangeToRangeGraph) toLowerIdx(i int) int {
	if i >= g.n {
		return i - g.n
	}
	return g.n + g.n + i
}

//
//
type D = int
type Deque struct{ l, r []D }

func NewDeque(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
