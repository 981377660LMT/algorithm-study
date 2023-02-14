// https://ei1333.github.io/library/other/dynamic-point-add-rectangle-sum.hpp
// Static Rectangle Add Rectangle Sum
// 二维矩形区间计数,支持单点添加
// n<=2e5 xi,yi,wi<=1e9
// AddRectangle(x1,x2,y1,y2) 向 [x1, x2) * [y1, y2) 矩形区间上添加w权重
// AddQuery(x1,x2,y1,y2) 添加查询为区间 [x1, x2) * [y1, y2) 的权重和
// CalculateQueries() 返回所有查询结果

// TODO 有问题
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int = 998244353

func main() {
	// https://judge.yosupo.jp/problem/static_rectangle_add_rectangle_sum
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	srars := NewStaticRectangleAddRectangleSum(n, q)
	for i := 0; i < n; i++ {
		var l, d, r, u, w int
		fmt.Fscan(in, &l, &d, &r, &u, &w)
		srars.AddRectangle(l, r, d, u, w)
	}

	for i := 0; i < q; i++ {
		var l, d, r, u int
		fmt.Fscan(in, &l, &d, &r, &u)
		srars.AddQuery(l, r, d, u)
	}

	res := srars.CalculateQueries()
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type Rectangle struct{ l, d, r, u, w int } // [l, r) * [d, u) 的权重为w.
type Query struct{ l, d, r, u int }
type StaticRectangleAddRectangleSum struct {
	rectangles []Rectangle
	queries    []Query
}

// 根据总点数和查询数初始化.
func NewStaticRectangleAddRectangleSum(n, q int) *StaticRectangleAddRectangleSum {
	return &StaticRectangleAddRectangleSum{
		rectangles: make([]Rectangle, 0, n),
		queries:    make([]Query, 0, q),
	}
}

// 添加矩形 [x1, x2) * [y1, y2) 的权重为w.
func (srars *StaticRectangleAddRectangleSum) AddRectangle(x1, x2, y1, y2, w int) {
	srars.rectangles = append(srars.rectangles, Rectangle{l: x1, d: y1, r: x2, u: y2, w: w})
}

// 添加查询为区间 [x1, x2) * [y1, y2) 的权重和.
func (srars *StaticRectangleAddRectangleSum) AddQuery(x1, x2, y1, y2 int) {
	srars.queries = append(srars.queries, Query{l: x1, d: y1, r: x2, u: y2})
}

// 返回所有查询结果.
func (srars *StaticRectangleAddRectangleSum) CalculateQueries() []int {
	n := len(srars.rectangles)
	q := len(srars.queries)
	res := make([]int, q)
	if n == 0 || q == 0 {
		return res
	}

	ys := make([]int, 0, n+n)
	for i := range srars.rectangles {
		ys = append(ys, srars.rectangles[i].d, srars.rectangles[i].u)
	}
	ys = sortedSet(ys)

	type Q struct {
		x    int
		d, u int
		t    bool
		idx  int
	}
	rs, qs := make([]Q, 0, n+n), make([]Q, 0, q+q)
	for i := range srars.rectangles {
		d := sort.SearchInts(ys, srars.rectangles[i].d)
		u := sort.SearchInts(ys, srars.rectangles[i].u)
		rs = append(rs,
			Q{t: false, x: srars.rectangles[i].l, d: d, u: u, idx: i},
			Q{t: true, x: srars.rectangles[i].r, d: d, u: u, idx: i},
		)
	}

	for i := range srars.queries {
		d := sort.SearchInts(ys, srars.queries[i].d)
		u := sort.SearchInts(ys, srars.queries[i].u)
		qs = append(qs,
			Q{t: false, x: srars.queries[i].l, d: d, u: u, idx: i},
			Q{t: true, x: srars.queries[i].r, d: d, u: u, idx: i},
		)
	}

	sort.Slice(rs, func(i, j int) bool { return rs[i].x < rs[j].x })
	sort.Slice(qs, func(i, j int) bool { return qs[i].x < qs[j].x })
	j := 0
	bit := newBinaryIndexedTree(len(ys))

	// TODO
	for i := range qs {
		for j < n && rs[j].x < qs[i].x {
			p := srars.rectangles[j]
			if rs[j].t {
				bit.Apply(rs[j].d, Hikari{-p.w * p.r * p.d, -p.w, p.d * p.w, p.r * p.w})
				bit.Apply(rs[j].u, Hikari{p.w * p.r * p.u, p.w, -p.u * p.w, -p.r * p.w})
			} else {
				bit.Apply(rs[j].d, Hikari{p.w * p.l * p.d, p.w, -p.d * p.w, -p.l * p.w})
				bit.Apply(rs[j].u, Hikari{-p.w * p.l * p.u, -p.w, p.u * p.w, p.l * p.w})
			}
			j++
		}

		p := srars.queries[qs[i].idx]
		uret := bit.Prod(qs[i].u)
		res[qs[i].idx] += uret[0]
		res[qs[i].idx] += uret[1] * qs[i].x * p.u
		res[qs[i].idx] += uret[2] * qs[i].x
		res[qs[i].idx] += uret[3] * p.u
		dret := bit.Prod(qs[i].d)
		res[qs[i].idx] -= dret[0]
		res[qs[i].idx] -= dret[1] * qs[i].x * p.d
		res[qs[i].idx] -= dret[2] * qs[i].x
		res[qs[i].idx] -= dret[3] * p.d
		if !qs[i].t {
			res[qs[i].idx] *= -1
		}
	}

	return res
}

// sorted(set(nums)) in python.
func sortedSet(nums []int) []int {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	res := make([]int, 0, len(set))
	for v := range set {
		res = append(res, v)
	}
	sort.Ints(res)
	return res
}

type Hikari = []int

func (*binaryIndexedTree) e() Hikari { return []int{0, 0, 0, 0} }
func (*binaryIndexedTree) op(a, b Hikari) Hikari {
	return []int{a[0] + b[0], a[1] + b[1], a[2] + b[2], a[3] + b[3]}
}

type binaryIndexedTree struct {
	n    int
	data []Hikari
}

// 長さ n の 0で初期化された配列で構築する.
func newBinaryIndexedTree(n int) *binaryIndexedTree {
	res := &binaryIndexedTree{n: n}
	data := make([]Hikari, n+1)
	for i := range data {
		data[i] = res.e()
	}
	res.data = data
	return res
}

// 要素 i に値 v を加える.
func (b *binaryIndexedTree) Apply(i int, v Hikari) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] = b.op(b.data[i], v)
	}
}

// [0, r) の要素の総和を求める.
func (b *binaryIndexedTree) Prod(r int) Hikari {
	res := b.e()
	for ; r > 0; r -= r & -r {
		res = b.op(res, b.data[r])
	}
	return res
}
