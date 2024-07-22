package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"unsafe"
)

func main() {
	// yosupo()
	arc130_f()
}

// https://leetcode.cn/problems/erect-the-fence/description/
func outerTrees(trees [][]int) [][]int {
	points := make([][2]int, len(trees))
	for i, p := range trees {
		points[i] = [2]int{p[0], p[1]}
	}
	hull := ConvexHull(points, Full, false)
	res := make([][]int, len(hull))
	for i, p := range hull {
		res[i] = []int{points[p][0], points[p][1]}
	}
	return res
}

// https://judge.yosupo.jp/problem/furthest_pair
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func() {
		var n int32
		fmt.Fscan(in, &n)
		points := make([][2]int, n)
		for i := int32(0); i < n; i++ {
			fmt.Fscan(in, &points[i][0], &points[i][1])
		}
		i, j := FurthestPair(points)
		fmt.Fprintln(out, i, j)
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		solve()
	}
}

// https://atcoder.jp/contests/arc130/tasks/arc130_f
// F - Replace by Average
func arc130_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	points := make([][2]int, n)
	for i := 0; i < n; i++ {
		points[i] = [2]int{i, nums[i]}
	}

	cht := NewConvexHullTrickDeque(false)
	Fenchel(points, Lower, true, func(start, end, a, b int) {
		if start != -INF {
			cht.AddLineMonotone(start, -start*a+b, -1)
		}
		if end != INF {
			cht.AddLineMonotone(end-1, -(end-1)*a+b, -1)
		}
	})
	res := 0
	for i := 0; i < n; i++ {
		v, _ := cht.QueryMonotoneInc(i)
		res += v
	}
	fmt.Fprintln(out, res)
}

type Mode uint8

const (
	Full Mode = iota
	Lower
	Upper
)

const INF int = 4e18

// (凸包/上凸包/下凸包).
func ConvexHull(points [][2]int, mode Mode, isPointsSorted bool) []int32 {
	n := len(points)
	if n == 1 {
		return []int32{0}
	}

	compare := func(i, j int32) int8 {
		x1, y1 := points[i][0], points[i][1]
		x2, y2 := points[j][0], points[j][1]
		if x1 < x2 || (x1 == x2 && y1 < y2) {
			return -1
		}
		if x1 == x2 && y1 == y2 {
			return 0
		}
		return 1
	}

	if n == 2 {
		res := compare(0, 1)
		if res == -1 {
			return []int32{0, 1}
		}
		if res == 1 {
			return []int32{1, 0}
		}
		// !包含共线需要 {0, 1}
		return []int32{0}
	}

	order := make([]int32, n)
	for i := int32(0); i < int32(n); i++ {
		order[i] = i
	}
	if !isPointsSorted {
		sort.Slice(order, func(i, j int) bool { return compare(order[i], order[j]) == -1 })
	}

	check := func(i, j, k int32) bool {
		x1, y1 := points[j][0]-points[i][0], points[j][1]-points[i][1]
		x2, y2 := points[k][0]-points[i][0], points[k][1]-points[i][1]
		return x1*y2 > x2*y1 // !包含共线需要 >=
	}

	calc := func() []int32 {
		var p []int32
		for _, k := range order {
			for len(p) > 1 {
				i, j := p[len(p)-2], p[len(p)-1]
				if check(i, j, k) {
					break
				}
				p = p[:len(p)-1]
			}
			p = append(p, k)
		}
		return p
	}

	var p []int32
	if mode == Full || mode == Lower {
		p = append(p, calc()...)
	}
	if mode == Full || mode == Upper {
		if len(p) > 0 {
			p = p[:len(p)-1]
		}
		reverse(order)
		p = append(p, calc()...)
	}
	if mode == Upper {
		reverse(p)
	}
	for len(p) >= 2 && points[p[0]] == points[p[len(p)-1]] {
		p = p[:len(p)-1]
	}
	return p
}

// 平面最远点对(旋转卡壳).
func FurthestPair(points [][2]int) (int32, int32) {
	best := -1
	resI, resJ := int32(-1), int32(-1)

	update := func(i, j int32) {
		x1, y1 := points[i][0]-points[j][0], points[i][1]-points[j][1]
		d := x1*x1 + y1*y1
		if d > best {
			best = d
			resI, resJ = i, j
		}
	}

	update(0, 1)

	order := ConvexHull(points, Full, false)
	n := int32(len(order))
	if n == 1 {
		return resI, resJ
	}
	if n == 2 {
		return order[0], order[1]
	}
	order = append(order, order...)

	newPoints := reArrage(points, order)
	check := func(i, j int32) bool {
		x1, y1 := newPoints[i+1][0]-newPoints[i][0], newPoints[i+1][1]-newPoints[i][1]
		x2, y2 := newPoints[j+1][0]-newPoints[j][0], newPoints[j+1][1]-newPoints[j][1]
		return x1*y2 > x2*y1
	}

	j := int32(1)
	for i := int32(0); i < n; i++ {
		if i >= j {
			j = i
		}
		for j < 2*n && check(i, j) {
			j++
		}
		update(order[i], order[j])
	}
	return resI, resJ
}

// TODO: not verified
// f: (start, end,a,b)：斜率在 [start,end) 之间时经过点 (a,b).
func Fenchel(points [][2]int, mode Mode, sorted bool, f func(start, end int, a, b int)) {
	if mode == Upper {
		newPoints := make([][2]int, len(points))
		for i, p := range points {
			newPoints[i] = [2]int{p[0], -p[1]}
		}
		res := make([][4]int, 0)
		Fenchel(newPoints, Lower, sorted, func(start, end, a, b int) {
			var l, r int
			if end == INF {
				l = -INF
			} else {
				l = 1 - end
			}
			if start == -INF {
				r = INF
			} else {
				r = 1 - start
			}
			l = max(l, -INF)
			r = min(r, INF)
			res = append(res, [4]int{l, r, a, -b})
		})
		for i := len(res) - 1; i >= 0; i-- {
			f(res[i][0], res[i][1], res[i][2], res[i][3])
		}
		return
	}

	order := ConvexHull(points, Lower, sorted)
	points = reArrage(points, order)

	floor := func(a, b int) int {
		tmp := (a%b != 0) && (a^b) < 0
		if tmp {
			return a/b - 1
		}
		return a / b
	}

	lo := -INF
	for i := 0; i < len(points); i++ {
		hi := INF
		if i+1 < len(points) {
			hi = min(hi, 1+floor(points[i+1][1]-points[i][1], points[i+1][0]-points[i][0]))
		}
		if lo < hi {
			f(lo, hi, points[i][0], points[i][1])
		}
		lo = hi
	}
}

func argSort(nums []int) []int32 {
	order := make([]int32, len(nums))
	for i := int32(0); i < int32(len(nums)); i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func reverse[T any](nums []T) {
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func reArrage[T any](nums []T, order []int32) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
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

func cast[To, From any](v From) To {
	return *(*To)(unsafe.Pointer(&v))
}

type Line struct{ k, b, id int }
type ConvexHullTrickDeque struct {
	isMin bool
	dq    *Deque
}

func NewConvexHullTrickDeque(isMin bool) *ConvexHullTrickDeque {
	return &ConvexHullTrickDeque{
		isMin: isMin,
		dq:    &Deque{},
	}
}

// 追加一条直线.需要保证斜率k是单调递增或者是单调递减的.
func (cht *ConvexHullTrickDeque) AddLineMonotone(k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}

	line := Line{k, b, id}
	if cht.dq.Empty() {
		cht.dq.AppendLeft(line)
		return
	}

	if cht.dq.Front().k <= k {
		if cht.dq.Front().k == k {
			if cht.dq.Front().b <= b {
				return
			}
			cht.dq.PopLeft()
		}
		for cht.dq.Len() >= 2 && cht.check(line, cht.dq.Front(), cht.dq.At(1)) {
			cht.dq.PopLeft()
		}
		cht.dq.AppendLeft(line)
	} else {
		if cht.dq.Back().k == k {
			if cht.dq.Back().b <= b {
				return
			}
			cht.dq.Pop()
		}
		for cht.dq.Len() >= 2 && cht.check(cht.dq.At(cht.dq.Len()-2), cht.dq.Back(), line) {
			cht.dq.Pop()
		}
		cht.dq.Append(line)
	}
}

// O(logn) 查询 k*x + b 的最小(大)值以及对应的直线id.
// 如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) Query(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	left, right := -1, cht.dq.Len()-1
	for left+1 < right {
		mid := (left + right) >> 1
		a, _ := cht.getY(cht.dq.At(mid), x)
		b, _ := cht.getY(cht.dq.At(mid+1), x)
		if a >= b {
			left = mid
		} else {
			right = mid
		}
	}

	res, lineId = cht.getY(cht.dq.At(right), x)
	if !cht.isMin {
		res = -res
	}
	return
}

// O(1) 查询 k*x + b 的最小(大)值以及对应的直线id.
//
//	需要保证x是单调递增的.
//	如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) QueryMonotoneInc(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	for cht.dq.Len() >= 2 {
		a, _ := cht.getY(cht.dq.Front(), x)
		b, _ := cht.getY(cht.dq.At(1), x)
		if a < b {
			break
		}
		cht.dq.PopLeft()
	}

	res, lineId = cht.getY(cht.dq.Front(), x)
	if !cht.isMin {
		res = -res
	}
	return
}

// O(1) 查询 k*x + b 的最小(大)值以及对应的直线id.
//
//	需要保证x是单调递减的.
//	如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) QueryMonotoneDec(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	for cht.dq.Len() >= 2 {
		a, _ := cht.getY(cht.dq.Back(), x)
		b, _ := cht.getY(cht.dq.At(cht.dq.Len()-2), x)
		if a < b {
			break
		}
		cht.dq.Pop()
	}

	res, lineId = cht.getY(cht.dq.Back(), x)
	if !cht.isMin {
		res = -res
	}
	return
}

func (cht *ConvexHullTrickDeque) check(a, b, c Line) bool {
	if b.b == a.b || c.b == b.b {
		return sign(b.k-a.k)*sign(c.b-b.b) >= sign(c.k-b.k)*sign(b.b-a.b)
	}
	return (b.k-a.k)*sign(c.b-b.b)*abs(c.b-b.b) >= (c.k-b.k)*sign(b.b-a.b)*abs(b.b-a.b)
}

func (cht *ConvexHullTrickDeque) getY(line Line, x int) (int, int) {
	return line.k*x + line.b, line.id
}

func sign(x int) int {
	if x == 0 {
		return 0
	} else if x > 0 {
		return 1
	} else {
		return -1
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type E = Line
type Deque struct{ l, r []E }

func (q Deque) Empty() bool     { return len(q.l) == 0 && len(q.r) == 0 }
func (q Deque) Len() int        { return len(q.l) + len(q.r) }
func (q *Deque) AppendLeft(v E) { q.l = append(q.l, v) }
func (q *Deque) Append(v E)     { q.r = append(q.r, v) }
func (q *Deque) PopLeft() (v E) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v E) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() E {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() E {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) E {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
