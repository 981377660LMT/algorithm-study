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

	N, M, Q := io.NextInt(), io.NextInt(), io.NextInt()
	edges := make([][3]int, M)
	for i := 0; i < M; i++ {
		u, v, w := io.NextInt()-1, io.NextInt()-1, io.NextInt()
		edges[i] = [3]int{u, v, w}
	}

	queries := make([][3]int, Q)
	alived := make([]bool, M)
	for i := 0; i < M; i++ {
		alived[i] = true
	}
	for i := 0; i < Q; i++ {
		op := io.NextInt()
		if op == 1 {
			x := io.NextInt() - 1
			queries[i] = [3]int{op, x, -1}
			alived[x] = false
		} else {
			x, y := io.NextInt()-1, io.NextInt()-1
			queries[i] = [3]int{op, x, y}
		}
	}

	newEdges := make([][3]int, 0, M)
	for i, e := range edges {
		if alived[i] {
			newEdges = append(newEdges, e)
		}
	}
	floyd := NewFloydFastAdd(N, newEdges)

	var res []int
	for i := Q - 1; i >= 0; i-- {
		query := queries[i]
		if query[0] == 1 {
			x := query[1]
			e := edges[x]
			floyd.AddEdge(e[0], e[1], e[2])
		} else {
			x, y := query[1], query[2]
			res = append(res, floyd.ShortestPath(x, y))
		}
	}

	for i := len(res) - 1; i >= 0; i-- {
		io.Println(res[i])
	}
}

const INF int = 1e18

type FloydFastAdd struct {
	dist [][]int
}

func NewFloydFastAdd(n int, edges [][3]int) *FloydFastAdd {
	dist := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = INF
		}
		dist[i][i] = 0
	}
	for _, edge := range edges {
		u, v, w := edge[0], edge[1], edge[2]
		dist[u][v] = min(w, dist[u][v])
		dist[v][u] = min(w, dist[v][u])
	}
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			if dist[i][k] == INF {
				continue
			}
			for j := 0; j < n; j++ {
				if dist[k][j] == INF {
					continue
				}
				cand := dist[i][k] + dist[k][j]
				if dist[i][j] > cand {
					dist[i][j] = cand
				}
			}
		}
	}
	return &FloydFastAdd{dist: dist}
}

func (ff *FloydFastAdd) AddEdge(u, v, w int) {
	n := len(ff.dist)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			cand := ff.dist[i][u] + w + ff.dist[v][j]
			if ff.dist[i][j] > cand {
				ff.dist[i][j] = cand
			}
			cand = ff.dist[i][v] + w + ff.dist[u][j]
			if ff.dist[i][j] > cand {
				ff.dist[i][j] = cand
			}
		}
	}
}

func (ff *FloydFastAdd) ShortestPath(start, target int) int {
	if ff.dist[start][target] < INF {
		return ff.dist[start][target]
	}
	return -1
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
