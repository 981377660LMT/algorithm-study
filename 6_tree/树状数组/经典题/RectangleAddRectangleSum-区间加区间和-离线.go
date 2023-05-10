// RectangleAddRectangleSum-区间加区间和-离线

// RectangleAddRectangleSum<int, mint> rect_sum;
// rect_sum.add_rectangle(xl, xr, yl, yr, mint(w));  // Add w to each of [xl, xr) * [yl, yr)
// rect_sum.add_query(l, r, d, u); // Get sum of [l, r) * [d, u)
// vector<mint> ret = rect_sum.solve();

// https://hitonanode.github.io/cplib-cpp/data_structure/rectangle_add_rectangle_sum.hpp
// O(qlogq)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// https://judge.yosupo.jp/problem/static_rectangle_add_rectangle_sum
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	S := NewRectangleAddRectangleSum()
	for i := 0; i < n; i++ {
		var left, down, right, up, w int
		fmt.Fscan(in, &left, &down, &right, &up, &w)
		S.AddRectangle(left, right, down, up, w)
	}
	for i := 0; i < q; i++ {
		var left, down, right, up int
		fmt.Fscan(in, &left, &down, &right, &up)
		S.AddQuery(left, right, down, up)
	}

	res := S.Solve()
	for _, r := range res {
		fmt.Fprintln(out, r)
	}
}

const MOD int = 998244353 // 如果不需要取模，则为1e18

type RectangleAddRectangleSum struct {
	addQueries []_addQuery
	sumQueries []_sumQuery
}

type _addQuery struct{ x1, x2, y1, y2, val int }
type _sumQuery struct{ x1, x2, y1, y2 int }

func NewRectangleAddRectangleSum() *RectangleAddRectangleSum {
	return &RectangleAddRectangleSum{}
}

// Add val to each of [x1, x2) * [y1, y2)
func (rars *RectangleAddRectangleSum) AddRectangle(x1, x2, y1, y2 int, val int) {
	rars.addQueries = append(rars.addQueries, _addQuery{x1: x1, x2: x2, y1: y1, y2: y2, val: val})
}

// Get sum of [x1, x2) * [y1, y2)
func (rars *RectangleAddRectangleSum) AddQuery(x1, x2, y1, y2 int) {
	rars.sumQueries = append(rars.sumQueries, _sumQuery{x1: x1, x2: x2, y1: y1, y2: y2})
}

func (rars *RectangleAddRectangleSum) Solve() []int {

	ys := make([]int, 0, len(rars.addQueries)*2)
	for _, a := range rars.addQueries {
		ys = append(ys, a.y1, a.y2)
	}
	ys, _ = _sortedSet(ys)

	ops := make([][3]int, 0, (len(rars.sumQueries)+len(rars.addQueries))*2)
	for q := 0; q < len(rars.sumQueries); q++ {
		ops = append(ops, [3]int{rars.sumQueries[q].x1, 0, q})
		ops = append(ops, [3]int{rars.sumQueries[q].x2, 1, q})
	}
	for q := 0; q < len(rars.addQueries); q++ {
		ops = append(ops, [3]int{rars.addQueries[q].x1, 2, q})
		ops = append(ops, [3]int{rars.addQueries[q].x2, 3, q})
	}
	sort.Slice(ops, func(i, j int) bool { return ops[i][0] < ops[j][0] })

	Y := len(ys)
	b00, b01, b10, b11 := _NewBIT(Y), _NewBIT(Y), _NewBIT(Y), _NewBIT(Y)
	res := make([]int, len(rars.sumQueries))

	for _, o := range ops {
		qtype, q := o[1], o[2]
		if qtype >= 2 {
			addQuery := rars.addQueries[q]
			i := sort.SearchInts(ys, addQuery.y1)
			j := sort.SearchInts(ys, addQuery.y2)
			x := o[0]
			yi, yj := addQuery.y1, addQuery.y2
			if qtype&1 == 1 {
				i, j = j, i
				yi, yj = yj, yi
			}

			b00.Add(i, x*yi%MOD*addQuery.val%MOD)
			b01.Add(i, -x*addQuery.val%MOD)
			b10.Add(i, -yi*addQuery.val%MOD)
			b11.Add(i, addQuery.val)
			b00.Add(j, -x*yj%MOD*addQuery.val%MOD)
			b01.Add(j, x*addQuery.val%MOD)
			b10.Add(j, yj*addQuery.val%MOD)
			b11.Add(j, -addQuery.val)
		} else {
			sumQuery := rars.sumQueries[q]
			i := sort.SearchInts(ys, sumQuery.y1)
			j := sort.SearchInts(ys, sumQuery.y2)
			x := o[0]
			yi, yj := sumQuery.y1, sumQuery.y2
			if qtype&1 == 1 {
				i, j = j, i
				yi, yj = yj, yi
			}

			res[q] += b00.Sum(i) + b01.Sum(i)*yi + b10.Sum(i)*x + b11.Sum(i)*x%MOD*yi
			res[q] -= b00.Sum(j) + b01.Sum(j)*yj + b10.Sum(j)*x + b11.Sum(j)*x%MOD*yj
			res[q] = (res[q]%MOD + MOD) % MOD
		}
	}

	return res
}

func _sortedSet(xs []int) (sorted []int, rank map[int]int) {
	set := make(map[int]struct{}, len(xs))
	for _, v := range xs {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	rank = make(map[int]int, len(sorted))
	for i, v := range sorted {
		rank[v] = i
	}
	return
}

type _BIT struct {
	n    int
	data []int
}

func _NewBIT(n int) *_BIT {
	return &_BIT{n: n, data: make([]int, n)}
}

func (bit *_BIT) Add(pos int, v int) {
	pos++
	v = (v + MOD) % MOD
	for pos > 0 && pos <= bit.n {
		bit.data[pos-1] = (bit.data[pos-1] + v) % MOD
		pos += pos & -pos
	}
}

func (bit *_BIT) Sum(k int) int {
	res := 0
	for k > 0 {
		res = (res + bit.data[k-1]) % MOD
		k -= k & -k
	}
	return res
}

func (bit *_BIT) Sum2(l, r int) int {
	return (bit.Sum(r) - bit.Sum(l)) % MOD
}
