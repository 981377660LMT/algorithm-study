// 矩形区间修改单点查询(离线)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func demo() {
	e := func() int32 { return 0 }
	op := func(a, b int32) int32 { return a + b }
	inv := func(a int32) int32 { return -a }
	ps := NewRectangleAddPointSumOffline(e, op, inv, false)
	ps.AddRectangle(0, 2, 0, 2, 1)
	ps.AddRectangle(1, 3, 1, 3, 2)
	ps.AddQuery(1, 1)
	ps.AddQuery(2, 2)
	ps.AddQuery(3, 3)
	res := ps.Calc()
	for _, v := range res {
		fmt.Println(v)
	}
}

func main() {
	yuki2338()
}

// No.2338 Range AtCoder Query
// https://yukicoder.me/problems/no/2338
// 给定长度为n的提交记录，每一项记录形如(pid, status)，表示问题编号和提交状态(AC/WA)。
// 再给定q个查询，每个查询形如(l, r)，表示查询区间[l, r)内的AC数和WA数。
// !AC数：通过的题目数量(每道题最多一次AC)。
// !WA数：第一次AC之前的提交次数。
//
// 对于每一个wa的记录，这个wa对查询[l,r)有贡献的条件为:
// 1. l之前没有ac，r之前有ac。
func yuki2338() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)
	pids := make([]int32, n)
	status := make([]bool, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &pids[i])
		pids[i]--
		var tmp string
		fmt.Fscan(in, &tmp)
		status[i] = tmp == "AC"
	}

	groups := make([][]int32, m)
	for i, v := range pids {
		groups[v] = append(groups[v], int32(i))
	}

	e := func() int32 { return 0 }
	op := func(a, b int32) int32 { return a + b }
	inv := func(a int32) int32 { return -a }
	ac := NewRectangleAddPointSumOffline(e, op, inv, true)
	wa := NewRectangleAddPointSumOffline(e, op, inv, true)

	for _, group := range groups {
		len_ := int32(len(group))
		prev := make([]int32, len_) // pre[i]: i之前的第一个AC
		for i := int32(0); i < len_; i++ {
			prev[i] = -1
		}
		for i, v := range group {
			if i > 0 {
				prev[i] = prev[i-1]
			}
			if status[v] {
				prev[i] = v
			}
		}

		nextAc := n
		for i := len_ - 1; i >= 0; i-- {
			index := group[i]
			if status[index] { // AC
				ac.AddRectangle(0, index+1, index+1, nextAc+1, 1)
				nextAc = index
			} else { // WA
				a, b := prev[i], nextAc
				if b == n {
					continue
				}
				wa.AddRectangle(a+1, index+1, b+1, n+1, 1)
			}
		}
	}

	for i := int32(0); i < q; i++ {
		var l, r int32
		fmt.Fscan(in, &l, &r)
		l--
		ac.AddQuery(l, r)
		wa.AddQuery(l, r)
	}

	res1, res2 := ac.Calc(), wa.Calc()
	for i := int32(0); i < q; i++ {
		fmt.Fprintln(out, res1[i], res2[i])
	}
}

const INF32 int32 = 1e9 + 10

type rect[E any] struct {
	y, x1, x2 int32
	w         E
}

type point struct {
	id   int32
	x, y int32
}

type RectangleAddPointSumOffline[E any] struct {
	smallX bool
	rects  []rect[E]
	points []point
	e      func() E
	op     func(e1, e2 E) E
	inv    func(e E) E
}

func NewRectangleAddPointSumOffline[E any](
	e func() E, op func(e1, e2 E) E, inv func(e E) E,
	smallX bool,
) *RectangleAddPointSumOffline[E] {
	return &RectangleAddPointSumOffline[E]{smallX: smallX, e: e, op: op, inv: inv}
}

func (ps *RectangleAddPointSumOffline[E]) AddRectangle(xl, xr, yl, yr int32, w E) {
	ps.rects = append(ps.rects, rect[E]{y: yl, x1: xl, x2: xr, w: w})
	ps.rects = append(ps.rects, rect[E]{y: yr, x1: xr, x2: xl, w: w})
}

func (ps *RectangleAddPointSumOffline[E]) AddQuery(x, y int32) {
	ps.points = append(ps.points, point{id: int32(len(ps.points)), x: x, y: y})
}

func (ps *RectangleAddPointSumOffline[E]) Calc() []E {
	n, q := int32(len(ps.rects)), int32(len(ps.points))
	if n == 0 || q == 0 {
		res := make([]E, q)
		for i := int32(0); i < q; i++ {
			res[i] = ps.e()
		}
		return res
	}

	rects, points := ps.rects, ps.points
	// compress x
	nx := int32(0)
	if !ps.smallX {
		sort.Slice(points, func(i, j int) bool { return points[i].x < points[j].x })
		keyX := make([]int32, 0, q)
		for i := int32(0); i < q; i++ {
			x := points[i].x
			if len(keyX) == 0 || keyX[len(keyX)-1] != x {
				keyX = append(keyX, x)
			}
			points[i].x = int32(len(keyX) - 1)
		}
		for i := int32(0); i < n; i++ {
			rects[i].x1 = lowerBound32(keyX, rects[i].x1)
			rects[i].x2 = lowerBound32(keyX, rects[i].x2)
		}
		nx = int32(len(keyX))
	} else {
		mx := INF32
		for i := int32(0); i < q; i++ {
			if tmp := points[i].x; tmp < mx {
				mx = tmp
			}
		}
		for i := int32(0); i < q; i++ {
			points[i].x -= mx
			if tmp := points[i].x + 1; tmp > nx {
				nx = tmp
			}
		}
		for i := int32(0); i < n; i++ {
			rects[i].x1 = clamp32(rects[i].x1-mx, 0, nx)
			rects[i].x2 = clamp32(rects[i].x2-mx, 0, nx)
		}
	}

	sort.Slice(points, func(i, j int) bool { return points[i].y < points[j].y })
	sort.Slice(rects, func(i, j int) bool { return rects[i].y < rects[j].y })

	bit := newFenwickTree(ps.e, ps.op, ps.inv)
	bit.Build(nx, func(i int32) E { return ps.e() })
	res := make([]E, q)
	for i := range res {
		res[i] = ps.e()
	}
	j := int32(0)
	for _, point := range points {
		q, x, y := point.id, point.x, point.y
		for j < n && rects[j].y <= y {
			rect := rects[j]
			j++
			bit.Update(rect.x1, rect.w)
			bit.Update(rect.x2, ps.inv(rect.w))
		}
		res[q] = bit.QueryPrefix(x + 1)
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
