package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
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
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DSL_2_C

	// !KDTree 先找出k*k矩形内的点, 再逐个检查是否成立

	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	xs, ys := make([]int, n), make([]int, n)
	idx := make([][]int, n)
	for i := 0; i < n; i++ {
		xs[i] = io.NextInt()
		ys[i] = io.NextInt()
		idx[i] = []int{i}
	}
	kd := NewKDTree(xs, ys, idx)

	q := io.NextInt()
	for i := 0; i < q; i++ {
		xl, xr, yl, yr := io.NextInt(), io.NextInt(), io.NextInt(), io.NextInt()
		e := kd.Query(xl, xr+1, yl, yr+1)
		sort.Ints(e)
		for _, v := range e {
			io.Println(v)
		}
		io.Println()
	}
}

func main2() {
	points := [][2]int{{1, 1}, {2, 2}}
	ids := make([][]int, len(points))
	for i := range ids {
		ids[i] = []int{i}
	}
	kd := NewKDTree([]int{1, 2}, []int{1, 2}, ids)
	fmt.Println(kd.Query(0, 1, 0, 1))
}

type P = int // 点的坐标类型
const INF P = 1e18

// !满足交换律的幺半群
type E = []int

func e() E { return []int{} }
func op(a, b E) E {
	return append(a, b...)
}

type KDTree struct {
	n           int
	closedRange [][4]P
	data        []E
}

// 数据集中所有点的横坐标、纵坐标和对应的值.
func NewKDTree(xs, ys []P, vs []E) *KDTree {
	n := len(xs)
	log := 0
	for (1 << uint(log)) < n {
		log++
	}
	data := make([]E, 1<<uint(log+1))
	closedRange := make([][4]P, 1<<uint(log+1))
	res := &KDTree{
		n:           n,
		closedRange: closedRange,
		data:        data,
	}
	if n > 0 {
		res.build(1, xs, ys, vs, true)
	}
	return res
}

func (kd *KDTree) Add(x, y P, v E) {
	kd.addRec(1, x, y, v)
}

// [xl, xr) x [yl, yr)。
func (kd *KDTree) Query(xl, xr, yl, yr P) E {
	if xr <= xl || yr <= yl {
		return e()
	}
	return kd.queryRec(1, xl, xr, yl, yr)
}

func (kd *KDTree) QueryAll() E {
	return kd.data[1]
}

func (kd *KDTree) build(idx int, xs, ys []P, vs []E, divx bool) {
	n := len(xs)
	range4 := &kd.closedRange[idx]
	xmin, xmax, ymin, ymax := &range4[0], &range4[1], &range4[2], &range4[3]
	*xmin, *ymin = INF, INF
	*xmax, *ymax = -INF, -INF

	for i := 0; i < n; i++ {
		x, y := xs[i], ys[i]
		if x < *xmin {
			*xmin = x
		}
		if x > *xmax {
			*xmax = x
		}
		if y < *ymin {
			*ymin = y
		}
		if y > *ymax {
			*ymax = y
		}
	}

	if *xmin == *xmax && *ymin == *ymax {
		x := e()
		for _, v := range vs {
			x = op(x, v)
		}
		kd.data[idx] = x
		return
	}

	m := n / 2
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}

	if divx {
		// nthElement(order, m, func(i, j int) bool {
		// 	return xs[i] < xs[j]
		// })
		sort.Slice(order, func(i, j int) bool {
			return xs[order[i]] < xs[order[j]]
		})
	} else {
		// nthElement(order, m, func(i, j int) bool {
		// 	return ys[i] < ys[j]
		// })
		sort.Slice(order, func(i, j int) bool {
			return ys[order[i]] < ys[order[j]]
		})
	}

	xs, ys, vs = reArrage(xs, order), reArrage(ys, order), reArrage2(vs, order)
	kd.build(2*idx, xs[:m], ys[:m], vs[:m], !divx)
	kd.build(2*idx+1, xs[m:], ys[m:], vs[m:], !divx)
	kd.data[idx] = op(kd.data[2*idx], kd.data[2*idx+1])
}

func (kd *KDTree) isLeaf(idx int) bool {
	range4 := kd.closedRange[idx]
	xmin, xmax, ymin, ymax := range4[0], range4[1], range4[2], range4[3]
	return xmin == xmax && ymin == ymax
}

func (kd *KDTree) isin(x, y P, idx int) bool {
	range4 := kd.closedRange[idx]
	xmin, xmax, ymin, ymax := range4[0], range4[1], range4[2], range4[3]
	return xmin <= x && x <= xmax && ymin <= y && y <= ymax
}

func (kd *KDTree) addRec(idx int, x, y P, v E) bool {
	if !kd.isin(x, y, idx) {
		return false
	}
	if kd.isLeaf(idx) {
		kd.data[idx] = op(kd.data[idx], v)
		return true
	}
	done := false
	if kd.addRec(2*idx, x, y, v) {
		done = true
	}
	if !done && kd.addRec(2*idx+1, x, y, v) {
		done = true
	}
	if done {
		kd.data[idx] = op(kd.data[2*idx], kd.data[2*idx+1])
	}
	return done
}

func (kd *KDTree) queryRec(idx int, x1, x2, y1, y2 P) E {
	range4 := kd.closedRange[idx]
	xmin, xmax, ymin, ymax := range4[0], range4[1], range4[2], range4[3]
	if x2 <= xmin || xmax < x1 || y2 <= ymin || ymax < y1 {
		return e()
	}
	if x1 <= xmin && xmax < x2 && y1 <= ymin && ymax < y2 {
		return kd.data[idx]
	}
	return op(kd.queryRec(2*idx, x1, x2, y1, y2), kd.queryRec(2*idx+1, x1, x2, y1, y2))
}

func reArrage(nums []P, order []int) []P {
	res := make([]P, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}

func reArrage2(nums []E, order []int) []E {
	res := make([]E, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
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
