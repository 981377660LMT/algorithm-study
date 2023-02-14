// https://ei1333.github.io/library/other/static-rectangle-add-rectangle-sum.hpp
// 静态二维矩形区间计数
// n<=2e5 xi,yi,wi<=1e9
// AddPoint(x,y,w) 向(x,y)点上添加w权重
// AddQuery(x1,x2,y1,y2) 添加查询为区间 [x1, x2) * [y1, y2) 的权重和
// CalculateQueries() 返回所有查询结果

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	// https://judge.yosupo.jp/problem/rectangle_sum
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	spars := NewStaticPointAddRectangleSum(n, q)
	for i := 0; i < n; i++ {
		var x, y, w int
		fmt.Fscan(in, &x, &y, &w)
		spars.AddPoint(x, y, w)
	}

	for i := 0; i < q; i++ {
		var l, d, r, u int
		fmt.Fscan(in, &l, &d, &r, &u)
		spars.AddQuery(l, r, d, u)
	}

	res := spars.CalculateQueries()
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type Point struct{ x, y, w int }
type Query struct{ l, d, r, u int }
type StaticPointAddRectangleSum struct {
	points  []Point
	queries []Query
}

// 指定点集和查询个数初始化容量.
func NewStaticPointAddRectangleSum(n, q int) *StaticPointAddRectangleSum {
	return &StaticPointAddRectangleSum{
		points:  make([]Point, 0, n),
		queries: make([]Query, 0, q),
	}
}

// 在(x,y)点上添加w权重.
func (sp *StaticPointAddRectangleSum) AddPoint(x, y, w int) {
	sp.points = append(sp.points, Point{x: x, y: y, w: w})
}

// 添加查询为区间 [x1, x2) * [y1, y2) 的权重和.
func (sp *StaticPointAddRectangleSum) AddQuery(x1, x2, y1, y2 int) {
	sp.queries = append(sp.queries, Query{l: x1, r: x2, d: y1, u: y2})
}

// 返回所有查询结果.
func (sp *StaticPointAddRectangleSum) CalculateQueries() []int {
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
		qs = append(qs, Q{x: query.l, d: d, u: u, t: false, idx: i})
		qs = append(qs, Q{x: query.r, d: d, u: u, t: true, idx: i})
	}

	sort.Slice(sp.points, func(i, j int) bool { return sp.points[i].x < sp.points[j].x })
	sort.Slice(qs, func(i, j int) bool { return qs[i].x < qs[j].x })

	j := 0
	bit := newBinaryIndexedTree(len(ys))
	for _, query := range qs {
		for j < n && sp.points[j].x < query.x {
			bit.Apply(sp.points[j].y, sp.points[j].w)
			j++
		}
		if query.t {
			res[query.idx] += bit.ProdRange(query.d, query.u)
		} else {
			res[query.idx] -= bit.ProdRange(query.d, query.u)
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
	res := int(0)
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
