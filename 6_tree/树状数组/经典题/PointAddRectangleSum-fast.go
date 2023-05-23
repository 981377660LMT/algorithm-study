// https://ei1333.github.io/library/other/dynamic-point-add-rectangle-sum.hpp
// Point Add Rectangle Sum
// 二维矩形区间计数,支持单点添加
// n<=2e5 xi,yi,wi<=1e9
// AddPoint(x,y,w) 向(x,y)点上添加w权重
// AddQuery(x1,x2,y1,y2) 添加查询为区间 [x1, x2) * [y1, y2) 的权重和
// Work() 按照添加顺序返回所有查询结果.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	// https://judge.yosupo.jp/problem/point_add_rectangle_sum
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	dpars := NewPointAddRectangleSum(n + q)
	for i := 0; i < n; i++ {
		var x, y, w int
		fmt.Fscan(in, &x, &y, &w)
		dpars.AddPoint(x, y, w)
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var x, y, w int
			fmt.Fscan(in, &x, &y, &w)
			dpars.AddPoint(x, y, w)
		} else {
			var l, d, r, u int
			fmt.Fscan(in, &l, &d, &r, &u)
			dpars.AddQuery(l, r, d, u)
		}
	}

	res := dpars.Work()
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
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

// 按照添加顺序返回所有查询结果..
func (dpars *DynamicPointAddRectangleSum) Work() []int {
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

		sub := solver.Work()
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

// 按照添加顺序返回所有查询结果..
func (sp *staticPointAddRectangleSum) Work() []int {
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
	bit := newBinaryIndexedTree(len(ys))
	for i := range qs {
		for j < n && sp.points[j].x < qs[i].x {
			bit.Apply(sp.points[j].y, sp.points[j].w)
			j++
		}
		if qs[i].t {
			res[qs[i].idx] += bit.ProdRange(qs[i].d, qs[i].u)
		} else {
			res[qs[i].idx] -= bit.ProdRange(qs[i].d, qs[i].u)
		}
	}

	return res
}

type binaryIndexedTree struct {
	n    int
	log  int
	data []int
}

// 長さ n の 0で初期化された配列で構築する.
func newBinaryIndexedTree(n int) *binaryIndexedTree {
	return &binaryIndexedTree{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

// 配列で構築する.
func newBinaryIndexedTreeFrom(arr []int) *binaryIndexedTree {
	res := newBinaryIndexedTree(len(arr))
	res.build(arr)
	return res
}

// 要素 i に値 v を加える.
func (b *binaryIndexedTree) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r) の要素の総和を求める.
func (b *binaryIndexedTree) Prod(r int) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *binaryIndexedTree) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *binaryIndexedTree) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.
func (b *binaryIndexedTree) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

func (b *binaryIndexedTree) build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}
