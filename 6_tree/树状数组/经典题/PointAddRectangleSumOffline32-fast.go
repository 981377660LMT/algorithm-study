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
	"os"
	"sort"
)

func main() {
	// https://judge.yosupo.jp/problem/point_add_rectangle_sum
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	dpars := NewPointAddRectangleSum(n + q)
	for i := int32(0); i < n; i++ {
		var x, y, w int32
		fmt.Fscan(in, &x, &y, &w)
		dpars.AddPoint(x, y, w)
	}

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var x, y, w int32
			fmt.Fscan(in, &x, &y, &w)
			dpars.AddPoint(x, y, w)
		} else {
			var l, d, r, u int32
			fmt.Fscan(in, &l, &d, &r, &u)
			dpars.AddQuery(l, r, d, u)
		}
	}

	res := dpars.Work()
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type Point struct{ x, y, w int32 }
type Query struct{ l, d, r, u int32 }
type DynamicPointAddRectangleSum struct {
	queries []interface{} // Point or Query
}

// 根据总点数初始化容量.
func NewPointAddRectangleSum(n int32) *DynamicPointAddRectangleSum {
	return &DynamicPointAddRectangleSum{queries: make([]interface{}, 0, n)}
}

// 在(x,y)点上添加w权重.
func (dpars *DynamicPointAddRectangleSum) AddPoint(x, y, w int32) {
	dpars.queries = append(dpars.queries, Point{x, y, w})
}

// 添加查询为区间 [x1, x2) * [y1, y2) 的权重和.
func (dpars *DynamicPointAddRectangleSum) AddQuery(x1, x2, y1, y2 int32) {
	dpars.queries = append(dpars.queries, Query{x1, y1, x2, y2})
}

// 按照添加顺序返回所有查询结果..
func (dpars *DynamicPointAddRectangleSum) Work() []int {
	q := int32(len(dpars.queries))
	rev := make([]int32, q)
	sz := int32(0)
	for i := int32(0); i < q; i++ {
		if _, ok := dpars.queries[i].(Query); ok { // holds_alternative
			rev[i] = sz
			sz++
		}
	}

	res := make([]int, sz)
	rangeQ := [][]int32{{0, q}}
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
			rangeQ = append(rangeQ, []int32{l, m})
		}
		if m+1 < r {
			rangeQ = append(rangeQ, []int32{m, r})
		}
	}

	return res
}

type staticPointAddRectangleSum struct {
	points  []Point
	queries []Query
}

// 指定点集和查询个数初始化容量.
func newStaticPointAddRectangleSum(n, q int32) *staticPointAddRectangleSum {
	return &staticPointAddRectangleSum{
		points:  make([]Point, 0, n),
		queries: make([]Query, 0, q),
	}
}

// 在(x,y)点上添加w权重.
func (sp *staticPointAddRectangleSum) AddPoint(x, y, w int32) {
	sp.points = append(sp.points, Point{x: x, y: y, w: w})
}

// 添加查询为区间 [x1, x2) * [y1, y2) 的权重和.
func (sp *staticPointAddRectangleSum) AddQuery(x1, x2, y1, y2 int32) {
	sp.queries = append(sp.queries, Query{l: x1, r: x2, d: y1, u: y2})
}

// 按照添加顺序返回所有查询结果..
func (sp *staticPointAddRectangleSum) Work() []int {
	n := int32(len(sp.points))
	q := int32(len(sp.queries))
	res := make([]int, q)
	if n == 0 || q == 0 {
		return res
	}

	sort.Slice(sp.points, func(i, j int) bool { return sp.points[i].y < sp.points[j].y })
	ys := make([]int32, 0, n)
	for i := range sp.points {
		if len(ys) == 0 || ys[len(ys)-1] != sp.points[i].y {
			ys = append(ys, sp.points[i].y)
		}
		sp.points[i].y = int32(len(ys) - 1)
	}

	type Q struct {
		x    int32
		d, u int32
		t    bool
		idx  int32
	}

	qs := make([]Q, 0, q+q)
	for i := int32(0); i < q; i++ {
		query := sp.queries[i]
		d := int32(sort.Search(len(ys), func(j int) bool { return ys[j] >= query.d }))
		u := int32(sort.Search(len(ys), func(j int) bool { return ys[j] >= query.u }))
		qs = append(qs, Q{x: query.l, d: d, u: u, t: false, idx: i}, Q{x: query.r, d: d, u: u, t: true, idx: i})
	}

	sort.Slice(sp.points, func(i, j int) bool { return sp.points[i].x < sp.points[j].x })
	sort.Slice(qs, func(i, j int) bool { return qs[i].x < qs[j].x })

	j := int32(0)
	bit := newBinaryIndexedTree(int32(len(ys)))
	for i := range qs {
		for j < n && sp.points[j].x < qs[i].x {
			bit.Apply(sp.points[j].y, int(sp.points[j].w))
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
	n    int32
	data []int
}

func newBinaryIndexedTree(n int32) *binaryIndexedTree {
	return &binaryIndexedTree{n: n, data: make([]int, n+1)}
}

func (b *binaryIndexedTree) Apply(i int32, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

func (b *binaryIndexedTree) Prod(r int32) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

func (b *binaryIndexedTree) ProdRange(l, r int32) int {
	if l == 0 {
		return b.Prod(r)
	}
	pos, neg := 0, 0
	for r > l {
		pos += b.data[r]
		r &= r - 1
	}
	for l > r {
		neg += b.data[l]
		l &= l - 1
	}
	return pos - neg
}
