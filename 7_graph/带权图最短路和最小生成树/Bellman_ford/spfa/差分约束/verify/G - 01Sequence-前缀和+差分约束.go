// https://atcoder.jp/contests/abc216/tasks/abc216_g
// 一个长度为n的序列，只由0和1组成，给出m个约束条件l, r, c，表示l 到r中至少有c个1，
// 问满足条件的序列是什么，如果有多种，则输出1的数量最小的那种
// n<=2e5, m<=2e5

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

	n, m := io.NextInt(), io.NextInt()
	D := NewDualShortestPath(n+5, true)
	for i := 0; i < m; i++ {
		l, r, c := io.NextInt(), io.NextInt(), io.NextInt() // !l,r从1开始
		D.AddEdge(l-1, r, -c)                               // Sr-Sl-1>=c
	}
	for i := 1; i <= n; i++ { // !前缀和约束
		D.AddEdge(i-1, i, 0) // Si-Si-1>=0
		D.AddEdge(i, i-1, 1) // Si-Si-1<=1
	}

	res, _ := D.Run()
	sb := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		sb = append(sb, res[i]-res[i-1]) // !每个位置取0还是1
	}

	for _, v := range sb {
		io.Print(v, " ")
	}
}

const INF int = 1e18

type DualShortestPath struct {
	n   int
	g   [][][2]int
	min bool
}

func NewDualShortestPath(n int, min bool) *DualShortestPath {
	return &DualShortestPath{
		n:   n,
		g:   make([][][2]int, n),
		min: min,
	}
}

// f(j) <= f(i) + w
func (d *DualShortestPath) AddEdge(i, j, w int) {
	if d.min {
		d.g[i] = append(d.g[i], [2]int{j, w})
	} else {
		d.g[j] = append(d.g[j], [2]int{i, w})
	}
}

// 求 `f(i) - f(0)` 的最小值/最大值, 并检测是否有负环/正环
func (d *DualShortestPath) Run() (dist []int, ok bool) {
	if d.min {
		return d.spfaMin()
	}
	return d.spfaMax()
}

func (d *DualShortestPath) spfaMin() (dist []int, ok bool) {
	dist = make([]int, d.n)
	queue := NewDeque2(d.n)
	count := make([]int, d.n)
	inQueue := make([]bool, d.n)
	for i := 0; i < d.n; i++ {
		queue.Append(i)
		inQueue[i] = true
		count[i] = 1
	}
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		inQueue[cur] = false
		for _, e := range d.g[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				if !inQueue[next] {
					count[next]++
					if count[next] >= d.n+1 {
						return nil, false
					}
					inQueue[next] = true
					queue.AppendLeft(next)
				}
			}
		}
	}

	for i := 0; i < d.n; i++ {
		dist[i] = -dist[i]
	}
	ok = true
	return
}

func (d *DualShortestPath) spfaMax() (dist []int, ok bool) {
	dist = make([]int, d.n)
	inQueue := make([]bool, d.n)
	count := make([]int, d.n)
	for i := 0; i < d.n; i++ {
		dist[i] = INF
	}

	queue := []int{0}
	dist[0] = 0
	inQueue[0] = true
	count[0] = 1
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		inQueue[cur] = false
		for _, e := range d.g[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				if !inQueue[next] {
					count[next]++
					if count[next] >= d.n+1 {
						return nil, false
					}
					inQueue[next] = true
					queue = append(queue, next)
				}
			}
		}
	}

	ok = true
	return
}

type D = int
type Deque struct{ l, r []D }

func NewDeque2(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

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
