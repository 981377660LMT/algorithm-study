// https://ei1333.github.io/library/other/static-rectangle-add-rectangle-sum.hpp
// 静态二维矩形区间计数
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
	// https://judge.yosupo.jp/problem/rectangle_sum
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	points, queries := make([]Point, 0, n), make([]Query, 0, q)
	for i := 0; i < n; i++ {
		var x, y, w int
		fmt.Fscan(in, &x, &y, &w)
		points = append(points, Point{x, y, w})
	}

	for i := 0; i < q; i++ {
		var l, d, r, u int
		fmt.Fscan(in, &l, &d, &r, &u)
		queries = append(queries, Query{l, r, d, u})
	}

	res := RectangleSum(points, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type Point = struct{ x, y, w int }
type Query = struct{ lx, rx, ly, ry int }

// Point: (x,y,w)
// Query: [lx,rx) * [ly,ry)
func RectangleSum(points []Point, queries []Query) []int {
	n := len(points)
	q := len(queries)
	res := make([]int, q)
	if n == 0 || q == 0 {
		return res
	}

	sort.Slice(points, func(i, j int) bool { return points[i].y < points[j].y })
	ys := make([]int, 0, n)
	for i := range points {
		if len(ys) == 0 || ys[len(ys)-1] != points[i].y {
			ys = append(ys, points[i].y)
		}
		points[i].y = len(ys) - 1
	}

	type Q struct {
		x    int
		d, u int
		t    bool
		idx  int
	}

	qs := make([]Q, 0, q+q)
	for i := 0; i < q; i++ {
		query := queries[i]
		d := sort.SearchInts(ys, query.ly)
		u := sort.SearchInts(ys, query.ry)
		qs = append(qs, Q{x: query.lx, d: d, u: u, t: false, idx: i})
		qs = append(qs, Q{x: query.rx, d: d, u: u, t: true, idx: i})
	}

	sort.Slice(points, func(i, j int) bool { return points[i].x < points[j].x })
	sort.Slice(qs, func(i, j int) bool { return qs[i].x < qs[j].x })

	j := 0
	bit := newBinaryIndexedTree(len(ys))
	for _, query := range qs {
		for j < n && points[j].x < query.x {
			bit.Apply(points[j].y, points[j].w)
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

func newBinaryIndexedTree(n int) *binaryIndexedTree {
	return &binaryIndexedTree{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

func newBinaryIndexedTreeFrom(arr []int) *binaryIndexedTree {
	res := newBinaryIndexedTree(len(arr))
	res.build(arr)
	return res
}

func (b *binaryIndexedTree) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

func (b *binaryIndexedTree) Prod(r int) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

func (b *binaryIndexedTree) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
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
