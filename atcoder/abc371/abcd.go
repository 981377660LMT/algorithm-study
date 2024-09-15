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

	sab, sac, sbc := io.Text(), io.Text(), io.Text()
	graph := make([][]int, 3)

	if sab == "<" {
		graph[0] = append(graph[0], 1)
	} else {
		graph[1] = append(graph[1], 0)
	}
	if sac == "<" {
		graph[0] = append(graph[0], 2)
	} else {
		graph[2] = append(graph[2], 0)
	}
	if sbc == "<" {
		graph[1] = append(graph[1], 2)
	} else {
		graph[2] = append(graph[2], 1)
	}

	order, ok := TopoSort(3, graph, true)
	if !ok {
		io.Println("Impossible")
		return
	}

	res := order[1]
	if res == 0 {
		io.Println("A")
	} else if res == 1 {
		io.Println("B")
	} else {
		io.Println("C")
	}

}

func TopoSort(n int, adjList [][]int, directed bool) (order []int, ok bool) {
	deg := make([]int, n)
	startDeg := 0
	if directed {
		for _, adj := range adjList {
			for _, j := range adj {
				deg[j]++
			}
		}
	} else {
		for i, adj := range adjList {
			deg[i] = len(adj)
		}
		startDeg = 1
	}

	queue := []int{}
	for v := 0; v < n; v++ {
		if deg[v] == startDeg {
			queue = append(queue, v)
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)
		for _, next := range adjList[cur] {
			deg[next]--
			if deg[next] == startDeg {
				queue = append(queue, next)
			}
		}
	}

	if len(order) < n {
		return nil, false
	}
	return order, true
}
