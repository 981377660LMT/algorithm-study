// 三维偏序
// https://www.luogu.com.cn/problem/P3810
// 有n个元素,每个元素有三个属性xi,yi,zi
// !对每个元素i,求出有多少个j满足xj<=xi,yj<=yi,zj<=zi且i!=j.

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"runtime/debug"
	"sort"
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

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, _ := io.NextInt(), io.NextInt()
	xs, ys, zs := make([]int, n), make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		xs[i], ys[i], zs[i] = io.NextInt(), io.NextInt(), io.NextInt()
	}
	res := Solve(xs, ys, zs)

	counter := make([]int, n)
	for _, v := range res {
		counter[v]++
	}
	for i := 0; i < n; i++ {
		fmt.Fprintln(out, counter[i])
	}
}

// !对每个元素i,求出有多少个j满足xj<=xi,yj<=yi,zj<=zi且i!=j.
func Solve(xs, ys, zs []int) []int {
	n := len(xs)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		a, b := order[i], order[j]
		if xs[a] != xs[b] {
			return xs[a] < xs[b]
		}
		if ys[a] != ys[b] {
			return ys[a] < ys[b]
		}
		return zs[a] < zs[b]
	})

	S := NewPointAddRectangleSum(n)

	round := 0
	roundId := make([]int, n)
	for i := 0; i < n; i++ {
		group := []int{order[i]}
		b, c := ys[order[i]], zs[order[i]]
		for i+1 < n && ys[order[i+1]] == b && zs[order[i+1]] == c { // 把第二、三个维度一样的一起拿出来(因为是小于等于)
			i++
			group = append(group, order[i])
		}
		S.AddPoint(b, c, len(group))
		S.AddQuery(0, b+1, 0, c+1)
		for _, v := range group {
			roundId[v] = round
		}
		round++
	}
	queries := S.Run()
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = queries[roundId[i]] - 1 // 减去自己
	}
	return res
}

type Point struct{ x, y, w int }
type Query struct{ l, d, r, u int }
type DynamicPointAddRectangleSum struct {
	queries []interface{} // Point or Query
}

// 根据总点数初始化容量.
func NewPointAddRectangleSum(n int) *DynamicPointAddRectangleSum {
	return &DynamicPointAddRectangleSum{queries: make([]interface{}, 0, n)}
}

// 在(x,y)点上添加w权重.
func (dpars *DynamicPointAddRectangleSum) AddPoint(x, y, w int) {
	dpars.queries = append(dpars.queries, Point{x, y, w})
}

// 添加查询为区间 [x1, x2) * [y1, y2) 的权重和.
func (dpars *DynamicPointAddRectangleSum) AddQuery(x1, x2, y1, y2 int) {
	dpars.queries = append(dpars.queries, Query{x1, y1, x2, y2})
}

// 返回所有查询结果.
func (dpars *DynamicPointAddRectangleSum) Run() []int {
	q := len(dpars.queries)
	rev := make([]int, q)
	sz := 0
	for i := 0; i < q; i++ {
		if _, ok := dpars.queries[i].(Query); ok { // holds_alternative
			rev[i] = sz
			sz++
		}
	}

	res := make([]int, sz)
	rangeQ := [][]int{{0, q}}
	for len(rangeQ) > 0 {
		l, r := rangeQ[0][0], rangeQ[0][1]
		rangeQ = rangeQ[1:]
		m := (l + r) >> 1
		solver := newStaticPointAddRectangleSum(0, 0)
		for k := l; k < m; k++ {
			if p, ok := dpars.queries[k].(Point); ok {
				solver.AddPoint(p.x, p.y, p.w)
			}
		}

		for k := m; k < r; k++ {
			if q, ok := dpars.queries[k].(Query); ok {
				solver.AddQuery(q.l, q.r, q.d, q.u)
			}
		}

		sub := solver.Run()
		for k, t := m, 0; k < r; k++ {
			if _, ok := dpars.queries[k].(Query); ok {
				res[rev[k]] += sub[t]
				t++
			}
		}

		if l+1 < m {
			rangeQ = append(rangeQ, []int{l, m})
		}
		if m+1 < r {
			rangeQ = append(rangeQ, []int{m, r})
		}
	}

	return res
}

type staticPointAddRectangleSum struct {
	points  []Point
	queries []Query
}

// 指定点集和查询个数初始化容量.
func newStaticPointAddRectangleSum(n, q int) *staticPointAddRectangleSum {
	return &staticPointAddRectangleSum{
		points:  make([]Point, 0, n),
		queries: make([]Query, 0, q),
	}
}

// 在(x,y)点上添加w权重.
func (sp *staticPointAddRectangleSum) AddPoint(x, y, w int) {
	sp.points = append(sp.points, Point{x: x, y: y, w: w})
}

// 添加查询为区间 [x1, x2) * [y1, y2) 的权重和.
func (sp *staticPointAddRectangleSum) AddQuery(x1, x2, y1, y2 int) {
	sp.queries = append(sp.queries, Query{l: x1, r: x2, d: y1, u: y2})
}

// 按照添加顺序返回所有查询结果.
func (sp *staticPointAddRectangleSum) Run() []int {
	n := len(sp.points)
	q := len(sp.queries)
	res := make([]int, q)
	if n == 0 || q == 0 {
		return res
	}

	sort.Slice(sp.points, func(i, j int) bool { return sp.points[i].y < sp.points[j].y })
	ys := make([]int, 0, n)
	for i := range sp.points {
		if len(ys) == 0 || ys[len(ys)-1] != sp.points[i].y {
			ys = append(ys, sp.points[i].y)
		}
		sp.points[i].y = len(ys) - 1
	}

	type Q struct {
		x    int
		d, u int
		t    bool
		idx  int
	}

	qs := make([]Q, 0, q+q)
	for i := 0; i < q; i++ {
		query := sp.queries[i]
		d := sort.SearchInts(ys, query.d)
		u := sort.SearchInts(ys, query.u)
		qs = append(qs, Q{x: query.l, d: d, u: u, t: false, idx: i}, Q{x: query.r, d: d, u: u, t: true, idx: i})
	}

	sort.Slice(sp.points, func(i, j int) bool { return sp.points[i].x < sp.points[j].x })
	sort.Slice(qs, func(i, j int) bool { return qs[i].x < qs[j].x })

	j := 0
	bit := NewBit(len(ys))
	for i := range qs {
		for j < n && sp.points[j].x < qs[i].x {
			bit.Add(sp.points[j].y, sp.points[j].w)
			j++
		}
		if qs[i].t {
			res[qs[i].idx] += bit.QueryRange(qs[i].d, qs[i].u)
		} else {
			res[qs[i].idx] -= bit.QueryRange(qs[i].d, qs[i].u)
		}
	}

	return res
}

type BIT struct {
	n    int
	data []int
}

func NewBit(n int) *BIT {
	res := &BIT{n: n, data: make([]int, n)}
	return res
}

func (b *BIT) Add(index int, v int) {
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BIT) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BIT) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}
