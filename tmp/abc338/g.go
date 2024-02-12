package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

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
	defer func() {
		if err := recover(); err != nil {
			io.Println("No")
		}
	}()

	n, m := io.NextInt(), io.NextInt()
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		from, to, cost := io.NextInt(), io.NextInt(), io.NextInt()
		from--
		to--
		edges[i] = Edge{from: from, to: to, cost: cost, ei: i}
	}

	const INF int = 1e18

	res := INF
	for start := 0; start < n; start++ {
		minCost, eis := directedMST(n, edges, start)
		io.Println(minCost, eis)
	}
	io.Println(res)
}

type Edge struct{ from, to, cost, ei int }

// 给定一个连通的有向图，求以root为根节点的最小生成树
//
//	返回值：最小生成树的权值和，最小生成树的边的编号
func directedMST(n int, edges []Edge, root int) (int, []int) {
	for i := 0; i < n; i++ {
		if i != root {
			edges = append(edges, Edge{i, root, 0, -1})
		}
	}

	x := 0

	par := make([]int, 2*n)
	vis := make([]int, 2*n)
	link := make([]int, 2*n)
	for i := range par {
		par[i] = -1
		vis[i] = -1
		link[i] = -1
	}

	heap := NewSkewHeap(true)
	ins := make([]*SkewHeapNode, 2*n)

	for i := range edges {
		e := edges[i]
		ins[e.to] = heap.Push(ins[e.to], e.cost, i)
	}

	st := []int{}

	go_ := func(x int) int {
		x = edges[ins[x].index].from
		for link[x] != -1 {
			st = append(st, x)
			x = link[x]
		}
		for _, p := range st {
			link[p] = x
		}
		st = st[:0]
		return x
	}

	for i := n; ins[x] != nil; i++ {
		for ; vis[x] == -1; x = go_(x) {
			vis[x] = 0
		}
		for ; x != i; x = go_(x) {
			w := ins[x].key
			v := heap.Pop(ins[x])
			v = heap.Add(v, -w)
			ins[i] = heap.Meld(ins[i], v)
			par[x] = i
			link[x] = i
		}
		for ; ins[x] != nil && go_(x) == x; ins[x] = heap.Pop(ins[x]) {
		}
	}

	cost := 0
	res := []int{}
	for i := root; i != -1; i = par[i] {
		vis[i] = 1
	}
	for i := x; i >= 0; i-- {
		if vis[i] == 1 {
			continue
		}
		cost += edges[ins[i].index].cost
		res = append(res, edges[ins[i].index].ei)
		for j := edges[ins[i].index].to; j != -1 && vis[j] == 0; j = par[j] {
			vis[j] = 1
		}
	}

	return cost, res
}

type E = int

type SkewHeapNode struct {
	key, lazy   E
	left, right *SkewHeapNode
	index       int
}

type SkewHeap struct {
	isMin bool
}

func NewSkewHeap(isMin bool) *SkewHeap {
	return &SkewHeap{isMin: isMin}
}

func (sk *SkewHeap) Push(t *SkewHeapNode, key E, index int) *SkewHeapNode {
	return sk.Meld(t, newNode(key, index))
}

func (sk *SkewHeap) Pop(t *SkewHeapNode) *SkewHeapNode {
	return sk.Meld(t.left, t.right)
}

func (sk *SkewHeap) Top(t *SkewHeapNode) E {
	return t.key
}

func (sk *SkewHeap) Meld(x, y *SkewHeapNode) *SkewHeapNode {
	sk.propagate(x)
	sk.propagate(y)
	if x == nil {
		return y
	}
	if y == nil {
		return x
	}
	if (x.key < y.key) != sk.isMin {
		x, y = y, x
	}
	x.right = sk.Meld(y, x.right)
	x.left, x.right = x.right, x.left
	return x
}

func (sk *SkewHeap) Add(t *SkewHeapNode, lazy E) *SkewHeapNode {
	if t == nil {
		return t
	}
	t.lazy += lazy
	sk.propagate(t)
	return t
}

func (sk *SkewHeap) propagate(t *SkewHeapNode) *SkewHeapNode {
	if t != nil && t.lazy != 0 {
		if t.left != nil {
			t.left.lazy += t.lazy
		}
		if t.right != nil {
			t.right.lazy += t.lazy
		}
		t.key += t.lazy
		t.lazy = 0
	}
	return t
}

func newNode(key E, index int) *SkewHeapNode {
	return &SkewHeapNode{key: key, index: index}
}
