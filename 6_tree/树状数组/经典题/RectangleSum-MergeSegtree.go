// https://kopricky.github.io/code/SegmentTrees/merge_segtree.html
// 静态二维矩形区间计数(RectangleSum-MergeSegtree)
// !慢

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// https://judge.yosupo.jp/problem/rectangle_sum
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	points := make([][2]int, n)
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		var x, y, w int
		fmt.Fscan(in, &x, &y, &w)
		points[i] = [2]int{x, y}
		weights[i] = w
	}

	S := NewRectangleSum(points, weights)
	for i := 0; i < q; i++ {
		// left<=x<right, down<=y<up
		var left, down, right, up int
		fmt.Fscan(in, &left, &down, &right, &up)
		fmt.Fprintln(out, S.Query(left, down, right, up))
	}
}

// RectangleSumMergeSegtree
type RectangleSum struct {
	n   int
	xs  []int
	ys  [][]int
	sum [][]int
}

func NewRectangleSum(points [][2]int, weights []int) *RectangleSum {
	sz := len(points)
	xs := make([]int, sz)
	n := 1
	for n < sz {
		n <<= 1
	}
	hoge := make([][][2]int, 2*n)
	sorted := make([][2]int, sz)
	for i := 0; i < sz; i++ {
		sorted[i] = [2]int{points[i][0], i}
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][0] < sorted[j][0]
	})
	ys := make([][]int, 2*n)
	sum := make([][]int, 2*n)
	for i := 0; i < sz; i++ {
		a := sorted[i][0]
		b := sorted[i][1]
		xs[i] = a
		ys[i+n] = []int{points[b][1]}
		hoge[i+n] = [][2]int{{points[b][1], b}}
		sum[i+n] = []int{0, weights[b]}
	}

	for i := n - 1; i >= 1; i-- {
		child1 := hoge[i<<1]
		child2 := hoge[(i<<1)|1]
		hoge[i] = make([][2]int, len(child1)+len(child2))
		cur := hoge[i]
		p := 0
		q := 0
		r := 0
		for p < len(child1) && q < len(child2) {
			if child1[p][0] < child2[q][0] {
				cur[r] = child1[p]
				p++
			} else {
				cur[r] = child2[q]
				q++
			}
			r++
		}
		for p < len(child1) {
			cur[r] = child1[p]
			p++
			r++
		}
		for q < len(child2) {
			cur[r] = child2[q]
			q++
			r++
		}
		ys[i] = make([]int, len(cur))
		sum[i] = make([]int, len(cur)+1)
		sum[i][0] = 0
		for j := 0; j < len(cur); j++ {
			ys[i][j] = cur[j][0]
			sum[i][j+1] = sum[i][j] + weights[cur[j][1]]
		}
	}

	return &RectangleSum{
		n:   n,
		xs:  xs,
		ys:  ys,
		sum: sum,
	}
}

func (rs *RectangleSum) Query(x1, y1, x2, y2 int) int {
	lxid := sort.SearchInts(rs.xs, x1)
	rxid := sort.SearchInts(rs.xs, x2)
	if lxid >= rxid {
		return 0
	}
	return rs._query(lxid, rxid, y1, y2)
}

func (rs *RectangleSum) _query(a, b, ly, ry int) int {
	res1 := 0
	res2 := 0
	a += rs.n
	b += rs.n
	sum := rs.sum
	for a^b > 0 {
		if a&1 > 0 {
			c := sort.SearchInts(rs.ys[a], ly)
			d := sort.SearchInts(rs.ys[a], ry)
			res1 += sum[a][d] - sum[a][c]
			a++
		}
		if b&1 > 0 {
			b--
			c := sort.SearchInts(rs.ys[b], ly)
			d := sort.SearchInts(rs.ys[b], ry)
			res2 += sum[b][d] - sum[b][c]
		}
		a >>= 1
		b >>= 1
	}
	return res1 + res2
}
