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
	// # N 頂点
	// # M 辺の有向グラフがあります。 頂点には
	// # 1 から
	// # N までの番号が付けられていて、
	// # i 番目の辺は頂点
	// # U
	// # i
	// # ​
	// #   から頂点
	// # V
	// # i
	// # ​
	// #   に向かって伸びています。

	// # あなたは今頂点
	// # 1 にいます。 以下の行動をちょうど
	// # 10
	// # 10
	// # 100

	// #   回繰り返して頂点
	// # 1 に戻ってくることが可能かどうか判定してください。

	// # 今いる頂点から伸びている辺を
	// # 1 つ選び、その辺が伸びている先の頂点に移動する。
	// # T 個のテストケースが与えられるので、それぞれについて解いてください。

	T := io.NextInt()
	for t := 0; t < T; t++ {
		n, m := io.NextInt(), io.NextInt()
		S := NewStronglyConnectedComponents(n)
		g := make([][]int, n)
		indeg := make([]int, n)
		for i := 0; i < m; i++ {
			u, v := io.NextInt()-1, io.NextInt()-1
			S.AddEdge(u, v, 1)
			g[u] = append(g[u], v)
			indeg[v]++
		}

		if indeg[0] == 0 {
			io.Println("No")
			continue
		}
		S.Build()
		groupOf0 := S.CompId[0]
		// fmt.Println(g, groupOf0, S.Group[groupOf0])

		cycle := len(S.Group[groupOf0])
		// mod 2/5
		for cycle%2 == 0 {
			cycle /= 2
		}
		for cycle%5 == 0 {
			cycle /= 5
		}
		if cycle == 1 {
			io.Println("Yes")
		} else {
			io.Println("No")
		}
	}
}

type WeightedEdge struct{ from, to, cost, index int }
type StronglyConnectedComponents struct {
	G      [][]WeightedEdge // 原图
	Dag    [][]int          // 强连通分量缩点后的DAG(有向图邻接表)
	CompId []int            // 每个顶点所属的强连通分量的编号
	Group  [][]int          // 每个强连通分量所包含的顶点
	rg     [][]WeightedEdge
	order  []int
	used   []bool
	eid    int
}

func NewStronglyConnectedComponents(n int) *StronglyConnectedComponents {
	return &StronglyConnectedComponents{G: make([][]WeightedEdge, n)}
}

func (scc *StronglyConnectedComponents) AddEdge(from, to, cost int) {
	scc.G[from] = append(scc.G[from], WeightedEdge{from, to, cost, scc.eid})
	scc.eid++
}

func (scc *StronglyConnectedComponents) Build() {
	scc.rg = make([][]WeightedEdge, len(scc.G))
	for i := range scc.G {
		for _, e := range scc.G[i] {
			scc.rg[e.to] = append(scc.rg[e.to], WeightedEdge{e.to, e.from, e.cost, e.index})
		}
	}

	scc.CompId = make([]int, len(scc.G))
	for i := range scc.CompId {
		scc.CompId[i] = -1
	}
	scc.used = make([]bool, len(scc.G))
	for i := range scc.G {
		scc.dfs(i)
	}
	for i, j := 0, len(scc.order)-1; i < j; i, j = i+1, j-1 {
		scc.order[i], scc.order[j] = scc.order[j], scc.order[i]
	}

	ptr := 0
	for _, v := range scc.order {
		if scc.CompId[v] == -1 {
			scc.rdfs(v, ptr)
			ptr++
		}
	}

	dag := make([][]int, ptr)
	visited := make(map[int]struct{}) // 边去重
	for i := range scc.G {
		for _, e := range scc.G[i] {
			x, y := scc.CompId[e.from], scc.CompId[e.to]
			if x == y {
				continue // 原来的边 x->y 的顶点在同一个强连通分量内,可以汇合同一个 SCC 的权值
			}
			hash := x*len(scc.G) + y
			if _, ok := visited[hash]; !ok {
				dag[x] = append(dag[x], y)
				visited[hash] = struct{}{}
			}
		}
	}
	scc.Dag = dag

	scc.Group = make([][]int, ptr)
	for i := range scc.G {
		scc.Group[scc.CompId[i]] = append(scc.Group[scc.CompId[i]], i)
	}
}

// 获取顶点k所属的强连通分量的编号
func (scc *StronglyConnectedComponents) Get(k int) int {
	return scc.CompId[k]
}

func (scc *StronglyConnectedComponents) dfs(idx int) {
	tmp := scc.used[idx]
	scc.used[idx] = true
	if tmp {
		return
	}
	for _, e := range scc.G[idx] {
		scc.dfs(e.to)
	}
	scc.order = append(scc.order, idx)
}

func (scc *StronglyConnectedComponents) rdfs(idx int, cnt int) {
	if scc.CompId[idx] != -1 {
		return
	}
	scc.CompId[idx] = cnt
	for _, e := range scc.rg[idx] {
		scc.rdfs(e.to, cnt)
	}
}
