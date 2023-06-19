package main

import (
	"fmt"
	"sort"
	"time"
)

func rangeAddQueries(n int, queries [][]int) [][]int {
	xs := make([]int, n*n)
	ys := make([]int, n*n)
	vs := make([]int, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			xs[i*n+j] = i
			ys[i*n+j] = j
		}
	}

	tree := NewKDTreeLazy(xs, ys, vs)
	for _, q := range queries {
		x1, y1, x2, y2 := q[0], q[1], q[2], q[3]
		tree.UpdateRange(x1, x2+1, y1, y2+1, 1)
	}
	res := make([][]int, n)
	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		for j := 0; j < n; j++ {
			res[i][j] = tree.Query(i, i+1, j, j+1)
		}
	}
	return res
}

func main() {
	n := int(2e5)
	xs := make([]int, n)
	ys := make([]int, n)
	vs := make([]int, n)
	for i := 0; i < n; i++ {
		xs[i] = i
		ys[i] = i
		vs[i] = i
	}
	kd := NewKDTreeLazy(xs, ys, vs)

	time1 := time.Now()
	for i := 0; i < n; i++ {
		kd.UpdateRange(0, i, 0, i, 1)
		kd.Query(i, i+1, i, i+1)
		kd.QueryAll()
		kd.Update(i, i, 0)
	}
	fmt.Println(time.Since(time1))
}

type P = int // 点的坐标类型
const INF P = 1e18

type E = int
type Id = int

func e() E   { return 0 }
func id() Id { return 0 }
func op(e1, e2 E) E {
	return e1 + e2
}
func mapping(f Id, x E, size int) E {
	return f*size + x
}
func composition(f, g Id) Id {
	return f + g
}

type KDTreeLazy struct {
	n           int
	closedRange [][4]P
	data        []E
	lazy        []Id
	size        []int
}

// 数据集中所有点的横坐标、纵坐标和对应的值.
func NewKDTreeLazy(xs, ys []P, vs []E) *KDTreeLazy {
	n := len(xs)
	log := 0
	for (1 << log) < n {
		log++
	}
	data := make([]E, 1<<(log+1))
	lazy := make([]Id, 1<<(log+1))
	for i := range lazy {
		lazy[i] = id()
	}
	closedRange := make([][4]P, 1<<(log+1))
	size := make([]int, 1<<(log+1))
	res := &KDTreeLazy{
		n:           n,
		closedRange: closedRange,
		data:        data,
		lazy:        lazy,
		size:        size,
	}
	if n > 0 {
		res._build(1, xs, ys, vs, true)
	}
	return res
}

func (kd *KDTreeLazy) Update(x, y P, v E) {
	kd._updateRec(1, x, y, v)
}

// [xl, xr) x [yl, yr)。
func (kd *KDTreeLazy) UpdateRange(xl, xr, yl, yr P, lazy Id) {
	if xr <= xl || yr <= yl {
		return
	}
	kd._updateRangeRec(1, xl, xr, yl, yr, lazy)
}

// [xl, xr) x [yl, yr)。
func (kd *KDTreeLazy) Query(xl, xr, yl, yr P) E {
	if xr <= xl || yr <= yl {
		return e()
	}
	return kd._queryRec(1, xl, xr, yl, yr)
}

func (kd *KDTreeLazy) QueryAll() E {
	return kd.data[1]
}

func (kd *KDTreeLazy) _pushDown(index int) {
	lazy := kd.lazy[index]
	if lazy == id() {
		return
	}
	kd._propagate(index<<1, lazy)
	kd._propagate(index<<1|1, lazy)
	kd.lazy[index] = id()
}

func (kd *KDTreeLazy) _propagate(index int, lazy Id) {
	kd.data[index] = mapping(lazy, kd.data[index], kd.size[index])
	if !kd._isLeaf(index) {
		kd.lazy[index] = composition(lazy, kd.lazy[index])
	}
}

func (kd *KDTreeLazy) _build(idx int, xs, ys []P, vs []E, divx bool) {
	n := len(xs)
	kd.size[idx] = n
	range4 := &kd.closedRange[idx]
	xmin, xmax, ymin, ymax := &range4[0], &range4[1], &range4[2], &range4[3]
	*xmin, *ymin = INF, INF
	*xmax, *ymax = -INF, -INF

	for i := 0; i < n; i++ {
		x, y := xs[i], ys[i]
		if x < *xmin {
			*xmin = x
		}
		if x > *xmax {
			*xmax = x
		}
		if y < *ymin {
			*ymin = y
		}
		if y > *ymax {
			*ymax = y
		}
	}

	if *xmin == *xmax && *ymin == *ymax {
		x := e()
		for _, v := range vs {
			x = op(x, v)
		}
		kd.data[idx] = x
		return
	}

	m := n / 2
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}

	if divx {
		sort.Slice(order, func(i, j int) bool {
			return xs[order[i]] < xs[order[j]]
		})
	} else {
		sort.Slice(order, func(i, j int) bool {
			return ys[order[i]] < ys[order[j]]
		})
	}

	xs, ys, vs = reArrage(xs, order), reArrage(ys, order), reArrage2(vs, order)
	kd._build(2*idx, xs[:m], ys[:m], vs[:m], !divx)
	kd._build(2*idx+1, xs[m:], ys[m:], vs[m:], !divx)
	kd.data[idx] = op(kd.data[idx<<1], kd.data[idx<<1|1])
}
func (kd *KDTreeLazy) _updateRec(index int, x, y P, v E) bool {
	if !kd._isin(x, y, index) {
		return false
	}
	if kd._isLeaf(index) {
		kd.data[index] = op(kd.data[index], v)
		kd.size[index] += 1
		return true
	}
	kd._pushDown(index)
	done := kd._updateRec(index<<1, x, y, v)
	if !done && kd._updateRec(index<<1|1, x, y, v) {
		done = true
	}
	if done {
		kd.data[index] = op(kd.data[index<<1], kd.data[index<<1|1])
		kd.size[index] = kd.size[index<<1] + kd.size[index<<1|1]
	}
	return done
}
func (kd *KDTreeLazy) _updateRangeRec(index int, x1, x2, y1, y2 P, lazy Id) {

	xmin, xmax, ymin, ymax := kd.closedRange[index][0], kd.closedRange[index][1], kd.closedRange[index][2], kd.closedRange[index][3]
	if x2 <= xmin || xmax < x1 || y2 <= ymin || ymax < y1 {
		return
	}
	if x1 <= xmin && xmax < x2 && y1 <= ymin && ymax < y2 {
		kd._propagate(index, lazy)
		return
	}
	kd._pushDown(index)
	kd._updateRangeRec(index<<1, x1, x2, y1, y2, lazy)
	kd._updateRangeRec(index<<1|1, x1, x2, y1, y2, lazy)
	kd.data[index] = op(kd.data[index<<1], kd.data[index<<1|1])
}
func (kd *KDTreeLazy) _queryRec(index int, x1, x2, y1, y2 P) E {

	xmin, xmax, ymin, ymax := kd.closedRange[index][0], kd.closedRange[index][1], kd.closedRange[index][2], kd.closedRange[index][3]
	if x2 <= xmin || xmax < x1 || y2 <= ymin || ymax < y1 {
		return e()
	}
	if x1 <= xmin && xmax < x2 && y1 <= ymin && ymax < y2 {
		return kd.data[index]
	}
	kd._pushDown(index)
	return op(
		kd._queryRec(index<<1, x1, x2, y1, y2),
		kd._queryRec(index<<1|1, x1, x2, y1, y2),
	)
}

func (kd *KDTreeLazy) _isLeaf(idx int) bool {

	xmin, xmax, ymin, ymax := kd.closedRange[idx][0], kd.closedRange[idx][1], kd.closedRange[idx][2], kd.closedRange[idx][3]
	return xmin == xmax && ymin == ymax
}

func (kd *KDTreeLazy) _isin(x, y P, idx int) bool {
	xmin, xmax, ymin, ymax := kd.closedRange[idx][0], kd.closedRange[idx][1], kd.closedRange[idx][2], kd.closedRange[idx][3]
	return xmin <= x && x <= xmax && ymin <= y && y <= ymax
}

func reArrage(nums []P, order []int) []P {
	res := make([]P, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}

func reArrage2(nums []E, order []int) []E {
	res := make([]E, len(order))
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
