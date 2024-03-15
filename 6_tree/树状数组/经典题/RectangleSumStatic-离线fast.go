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
	bit := NewBit(len(ys))
	for _, query := range qs {
		for j < n && points[j].x < query.x {
			bit.Add(points[j].y, points[j].w)
			j++
		}
		if query.t {
			res[query.idx] += bit.QueryRange(query.d, query.u)
		} else {
			res[query.idx] -= bit.QueryRange(query.d, query.u)
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
