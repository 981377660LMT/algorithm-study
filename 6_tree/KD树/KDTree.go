package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
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
	// https://atcoder.jp/contests/abc234/tasks/abc234_h
	// 给定二维平面上的N个点(i, gt)和一个正整数K。
	// 请列出所有`欧几里得距离`小于等于K的点对。
	// 1<N<2e5，1<K<1.5×1e9。
	// 保证最多4×1e5对点对将被枚举。

	// !KDTree 先找出k*k矩形内的点, 再逐个检查是否成立

	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, k := io.NextInt(), io.NextInt()
	xs, ys := make([]int, n), make([]int, n)
	for i := range xs {
		xs[i], ys[i] = io.NextInt(), io.NextInt()
	}

	calDist := func(x1, y1, x2, y2 int) int { return abs(x1-x2) + abs(y1-y2) }
	kdt := NewKDTree(xs, ys, calDist)
	xMin, xMax, yMin, yMax := INF, -INF, INF, -INF
	for i := 0; i < n; i++ {
		if xs[i] < xMin {
			xMin = xs[i]
		}
		if xs[i] > xMax {
			xMax = xs[i]
		}
		if ys[i] < yMin {
			yMin = ys[i]
		}
		if ys[i] > yMax {
			yMax = ys[i]
		}
	}

	res := make([][2]int, 0)
	for i := 0; i < n; i++ {
		a, b, c, d := xs[i]-k, xs[i]+k+1, ys[i]-k, ys[i]+k+1
		a, b, c, d = max(a, xMin), min(b, xMax+1), max(c, yMin), min(d, yMax+1)
		cand := kdt.CollectInRectangle(a, b, c, d, -1)
		sort.Ints(cand)
		for _, j := range cand {
			if i >= j {
				continue
			}
			dx, dy := xs[i]-xs[j], ys[i]-ys[j]
			if dx*dx+dy*dy <= k*k {
				res = append(res, [2]int{i, j})
			}
		}
	}

	io.Println(len(res))
	for _, p := range res {
		io.Println(p[0]+1, p[1]+1)
	}

}

func circleGame(toys [][]int, circles [][]int, r int) int {
	points := make([]Point, len(circles))
	for i, circle := range circles {
		points[i] = Point{circle[0], circle[1]}
	}

	kdTree := NewKDTree(points, func(p1, p2 Point) float64 {
		dx, dy := float64(p1[0]-p2[0]), float64(p1[1]-p2[1])
		return math.Sqrt(dx*dx + dy*dy)
	})

	res := 0
	for _, toy := range toys {
		minDist, _ := kdTree.FindNearest(Point{toy[0], toy[1]}, float64(r), true)
		if minDist+float64(toy[2]) <= float64(r) {
			res++
		}
	}

	return res
}

const INF int = 1e18

type KDTree struct {
	n           int
	closedRange [][4]int
	data        []int
	calDist     func(x1, y1, x2, y2 int) int
}

func NewKDTree(xs, ys []int, calDist func(x1, y1, x2, y2 int) int) *KDTree {
	n := len(xs)
	log := 0
	for 1<<log < n {
		log++
	}
	data := make([]int, 1<<(log+1))
	for i := range data {
		data[i] = -1
	}
	closedRange := make([][4]int, 1<<(log+1))
	vs := make([]int, n)
	for i := range vs {
		vs[i] = i
	}
	res := &KDTree{
		n:           n,
		closedRange: closedRange,
		data:        data,
		calDist:     calDist,
	}
	if n > 0 {
		res.build(1, xs, ys, vs, true)
	}
	return res
}

// 返回 矩形 [x1, x2) * [y1, y2) 中的点的编号, 最多 maxSize 个.
//  当 maxSize 为 -1 时, 返回所有点.
func (kd *KDTree) CollectInRectangle(x1, x2, y1, y2, maxSize int) []int {
	if x1 > x2 || y1 > y2 || kd.n == 0 {
		return []int{}
	}
	if maxSize == -1 {
		maxSize = kd.n
	}
	res := []int{}
	kd.rectRec(1, x1, x2, y1, y2, &res, maxSize)
	return res
}

// 返回最近邻点的编号, -1 表示不存在最近邻点.
// 不保证计算量的O(logn), 要求点群随机分布.
//  n,q = 1e5 => 1s
//  ban: 禁止的点的编号, -1 表示不禁止.
func (kd *KDTree) SearchNearestNeighbor(x, y int, ban int) int {
	if kd.n == 0 {
		return -1
	}
	res := [2]int{-1, INF}
	kd.nnsRec(1, x, y, &res, ban)
	return res[0]
}

func (kd *KDTree) build(idx int, xs, ys, vs []int, divx bool) {
	n := len(xs)
	xmin, xmax, ymin, ymax := &kd.closedRange[idx][0], &kd.closedRange[idx][1], &kd.closedRange[idx][2], &kd.closedRange[idx][3]
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
	if n == 1 {
		kd.data[idx] = vs[0]
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

	xs, ys, vs = reArrage(xs, order), reArrage(ys, order), reArrage(vs, order)
	kd.build(2*idx, xs[:m], ys[:m], vs[:m], !divx)
	kd.build(2*idx+1, xs[m:], ys[m:], vs[m:], !divx)
}

func (kd *KDTree) rectRec(i, x1, x2, y1, y2 int, res *[]int, ms int) {
	if len(*res) == ms {
		return
	}
	xmin, xmax, ymin, ymax := kd.closedRange[i][0], kd.closedRange[i][1], kd.closedRange[i][2], kd.closedRange[i][3]
	if x2 <= xmin || xmax < x1 {
		return
	}
	if y2 <= ymin || ymax < y1 {
		return
	}
	if kd.data[i] != -1 {
		*res = append(*res, kd.data[i])
		return
	}
	kd.rectRec(2*i, x1, x2, y1, y2, res, ms)
	kd.rectRec(2*i+1, x1, x2, y1, y2, res, ms)
}

func (kd *KDTree) bestDistSquared(i, x, y int) int {
	if i >= len(kd.closedRange) {
		return INF
	}
	xmin, xmax, ymin, ymax := kd.closedRange[i][0], kd.closedRange[i][1], kd.closedRange[i][2], kd.closedRange[i][3]
	// clamp
	clampedX := x
	if clampedX < xmin {
		clampedX = xmin
	} else if clampedX > xmax {
		clampedX = xmax
	}
	dx := x - clampedX
	clampedY := y
	if clampedY < ymin {
		clampedY = ymin
	} else if clampedY > ymax {
		clampedY = ymax
	}
	dy := y - clampedY
	return kd.calDist(0, 0, dx, dy)
}

func (kd *KDTree) nnsRec(i, x, y int, res *[2]int, ban int) {
	d := kd.bestDistSquared(i, x, y)
	if d >= res[1] {
		return
	}
	if kd.data[i] != -1 && kd.data[i] != ban {
		res[0], res[1] = kd.data[i], d
		return
	}
	d0 := kd.bestDistSquared(2*i, x, y)
	d1 := kd.bestDistSquared(2*i+1, x, y)
	if d0 < d1 {
		kd.nnsRec(2*i, x, y, res, ban)
		kd.nnsRec(2*i+1, x, y, res, ban)
	} else {
		kd.nnsRec(2*i+1, x, y, res, ban)
		kd.nnsRec(2*i, x, y, res, ban)
	}
}

func reArrage(nums []int, order []int) []int {
	res := make([]int, len(order))
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
