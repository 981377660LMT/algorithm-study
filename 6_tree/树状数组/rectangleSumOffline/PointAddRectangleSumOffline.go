// 矩形单点修改区间查询(离线)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func demo() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPointAddRectangleSumOffline(e, op, inv, false)
	S.AddPoint(1, 1, 1)
	S.AddPoint(2, 2, 2)
	S.AddQuery(0, 3, 0, 3)
	S.AddQuery(1, 2, 1, 2)
	fmt.Println(S.Calc())
}

// https://judge.yosupo.jp/problem/rectangle_sum
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPointAddRectangleSumOffline(e, op, inv, false)

	for i := int32(0); i < n; i++ {
		var x, y int32
		var w int
		fmt.Fscan(in, &x, &y, &w)
		S.AddPoint(x, y, w)
	}

	for i := int32(0); i < q; i++ {
		var l, d, r, u int32
		fmt.Fscan(in, &l, &d, &r, &u)
		S.AddQuery(l, r, d, u)
	}

	res := S.Calc()
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

const INF32 int32 = 1e9 + 10

type point[E any] struct {
	x, y int32
	w    E
}

type query struct{ xl, xr, yl, yr int32 }

type PointAddRectangleSumOffline[E any] struct {
	smallX  bool
	points  []point[E]
	queires []query
	e       func() E
	op      func(e1, e2 E) E
	inv     func(e E) E
}

func NewPointAddRectangleSumOffline[E any](
	e func() E, op func(e1, e2 E) E, inv func(e E) E,
	smallX bool,
) *PointAddRectangleSumOffline[E] {
	return &PointAddRectangleSumOffline[E]{smallX: smallX, e: e, op: op, inv: inv}
}

func (ps *PointAddRectangleSumOffline[E]) AddPoint(x, y int32, w E) {
	ps.points = append(ps.points, point[E]{x: x, y: y, w: w})
}

func (ps *PointAddRectangleSumOffline[E]) AddQuery(xl, xr, yl, yr int32) {
	ps.queires = append(ps.queires, query{xl: xl, xr: xr, yl: yl, yr: yr})
}

func (ps *PointAddRectangleSumOffline[E]) Calc() []E {
	n, q := int32(len(ps.points)), int32(len(ps.queires))
	if n == 0 || q == 0 {
		res := make([]E, q)
		for i := int32(0); i < q; i++ {
			res[i] = ps.e()
		}
		return res
	}

	points, queires := ps.points, ps.queires
	// compress x
	nx := int32(0)
	if !ps.smallX {
		sort.Slice(points, func(i, j int) bool { return points[i].x < points[j].x })
		keyX := make([]int32, 0, n)
		for i := int32(0); i < n; i++ {
			x := points[i].x
			if len(keyX) == 0 || keyX[len(keyX)-1] != x {
				keyX = append(keyX, x)
			}
			points[i].x = int32(len(keyX) - 1)
		}
		for i := int32(0); i < q; i++ {
			queires[i].xl = lowerBound32(keyX, queires[i].xl)
			queires[i].xr = lowerBound32(keyX, queires[i].xr)
		}
		nx = int32(len(keyX))
	} else {
		mx := INF32
		for i := int32(0); i < n; i++ {
			if tmp := points[i].x; tmp < mx {
				mx = tmp
			}
		}
		for i := int32(0); i < n; i++ {
			points[i].x -= mx
			if tmp := points[i].x + 1; tmp > nx {
				nx = tmp
			}
		}
		for i := int32(0); i < q; i++ {
			queires[i].xl = clamp32(queires[i].xl-mx, 0, nx)
			queires[i].xr = clamp32(queires[i].xr-mx, 0, nx)
		}
	}

	events := make([][4]int32, 0, q+q)
	for _, query := range queires {
		events = append(events, [4]int32{query.yl, query.xl, query.xr, int32(len(events))})
		events = append(events, [4]int32{query.yr, query.xl, query.xr, int32(len(events))})
	}
	sort.Slice(points, func(i, j int) bool { return points[i].y < points[j].y })
	sort.Slice(events, func(i, j int) bool { return events[i][0] < events[j][0] })

	bit := newFenwickTree(ps.e, ps.op, ps.inv)
	bit.Build(nx, func(i int32) E { return ps.e() })
	res := make([]E, q)
	for i := range res {
		res[i] = ps.e()
	}
	j := int32(0)
	for _, event := range events {
		y, xl, xr, idx := event[0], event[1], event[2], event[3]
		for j < n && points[j].y < y {
			bit.Update(points[j].x, points[j].w)
			j++
		}
		w := bit.QueryRange(xl, xr)
		q := idx >> 1
		if idx&1 == 0 {
			res[q] = ps.op(res[q], ps.inv(w))
		} else {
			res[q] = ps.op(res[q], w)
		}
	}
	return res
}

type fenwickTree[E any] struct {
	n     int32
	total E
	data  []E
	e     func() E
	op    func(e1, e2 E) E
	inv   func(e E) E
}

func newFenwickTree[E any](e func() E, op func(e1, e2 E) E, inv func(e E) E) *fenwickTree[E] {
	return &fenwickTree[E]{e: e, op: op, inv: inv}
}

func (fw *fenwickTree[E]) Build(n int32, f func(i int32) E) {
	data := make([]E, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	for i := int32(1); i <= n; i++ {
		if j := i + (i & -i); j <= n {
			data[j-1] = fw.op(data[i-1], data[j-1])
		}
	}
	fw.n = n
	fw.data = data
	fw.total = fw.QueryPrefix(n)
}

func (fw *fenwickTree[E]) QueryAll() E { return fw.total }

// [0, end)
func (fw *fenwickTree[E]) QueryPrefix(end int32) E {
	if end > fw.n {
		end = fw.n
	}
	res := fw.e()
	for ; end > 0; end &= end - 1 {
		res = fw.op(res, fw.data[end-1])
	}
	return res
}

// [start, end)
func (fw *fenwickTree[E]) QueryRange(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start > end {
		return fw.e()
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	pos, neg := fw.e(), fw.e()
	for end > start {
		pos = fw.op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = fw.op(neg, fw.data[start-1])
		start &= start - 1
	}
	return fw.op(pos, fw.inv(neg))
}

// 要求op满足交换律(commute).
func (fw *fenwickTree[E]) Update(i int32, x E) {
	fw.total = fw.op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = fw.op(fw.data[i-1], x)
	}
}

func (fw *fenwickTree[E]) GetAll() []E {
	res := make([]E, fw.n)
	for i := int32(0); i < fw.n; i++ {
		res[i] = fw.QueryRange(i, i+1)
	}
	return res
}

func lowerBound32(nums []int32, x int32) int32 {
	left, right := int32(0), int32(len(nums)-1)
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func clamp32(x, l, r int32) int32 {
	if x < l {
		return l
	}
	if x > r {
		return r
	}
	return x
}
