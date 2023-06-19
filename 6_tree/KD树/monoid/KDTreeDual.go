package main

import (
	"fmt"
	"sort"
	"time"
)

func rangeAddQueries(n int, queries [][]int) [][]int {
	xs := make([]int, n*n)
	ys := make([]int, n*n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			xs[i*n+j] = i
			ys[i*n+j] = j
		}
	}

	tree := NewKDTreeDual(xs, ys)
	for _, q := range queries {
		x1, y1, x2, y2 := q[0], q[1], q[2], q[3]
		tree.Update(x1, x2+1, y1, y2+1, 1)
	}
	res := make([][]int, n)

	for i := 0; i < n; i++ {
		res[i] = make([]int, n)
		for j := 0; j < n; j++ {
			res[i][j] = tree.Get(i*n + j)
		}
	}
	return res
}

func main() {
	n := int(2e5)
	xs := make([]int, n)
	ys := make([]int, n)
	for i := 0; i < n; i++ {
		xs[i] = i
		ys[i] = i
	}
	kd := NewKDTreeDual(xs, ys)

	time1 := time.Now()
	for i := 0; i < n; i++ {
		kd.Update(0, i, 0, i, 1)
		kd.Get(i)
	}
	fmt.Println(time.Since(time1))
}

type P = int // 点的坐标类型
const INF P = 1e18

type Id = int

func id() Id { return 0 }
func composition(f, g Id) Id {
	return f + g
}

type KDTreeDual struct {
	n, log      int
	closedRange [][4]P
	lazy        []Id
	size        []int
	where       []int
}

// 数据集中所有点的横坐标、纵坐标.
func NewKDTreeDual(xs, ys []P) *KDTreeDual {
	n := len(xs)
	log := 0
	for (1 << log) < n {
		log++
	}
	lazy := make([]Id, 1<<(log+1))
	for i := range lazy {
		lazy[i] = id()
	}
	closedRange := make([][4]P, 1<<(log+1))
	size := make([]int, 1<<(log+1))
	where := make([]int, n)
	res := &KDTreeDual{
		n:           n,
		log:         log,
		closedRange: closedRange,
		lazy:        lazy,
		size:        size,
		where:       where,
	}
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	if n > 0 {
		res._build(1, xs, ys, order, true)
	}
	return res
}

// 查询第i个点的值.(顺序与构造函数中的顺序一致)
func (kd *KDTreeDual) Get(i int) Id {
	i = kd.where[i]
	for k := kd.log; k >= 1; k-- {
		kd._pushDown(i >> k)
	}
	return kd.lazy[i]
}

// [xl, xr) x [yl, yr)。
func (kd *KDTreeDual) Update(xl, xr, yl, yr P, lazy Id) {
	if xr <= xl || yr <= yl {
		return
	}
	kd._updateRec(1, xl, xr, yl, yr, lazy)
}

func (kd *KDTreeDual) _pushDown(index int) {
	lazy := kd.lazy[index]
	if lazy == id() {
		return
	}
	kd._propagate(index<<1, lazy)
	kd._propagate(index<<1|1, lazy)
	kd.lazy[index] = id()
}

func (kd *KDTreeDual) _propagate(index int, lazy Id) {
	kd.lazy[index] = composition(lazy, kd.lazy[index])
}

func (kd *KDTreeDual) _build(idx int, xs, ys []P, rawIndex []int, divx bool) {
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
		kd.where[rawIndex[0]] = idx
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

	xs, ys, rawIndex = reArrage(xs, order), reArrage(ys, order), reArrage2(rawIndex, order)
	kd._build(2*idx, xs[:m], ys[:m], rawIndex[:m], !divx)
	kd._build(2*idx+1, xs[m:], ys[m:], rawIndex[m:], !divx)
}

func (kd *KDTreeDual) _isLeaf(idx int) bool {
	range4 := kd.closedRange[idx]
	xmin, xmax, ymin, ymax := range4[0], range4[1], range4[2], range4[3]
	return xmin == xmax && ymin == ymax
}

func (kd *KDTreeDual) _isin(x, y P, idx int) bool {
	range4 := kd.closedRange[idx]
	xmin, xmax, ymin, ymax := range4[0], range4[1], range4[2], range4[3]
	return xmin <= x && x <= xmax && ymin <= y && y <= ymax
}

func (kd *KDTreeDual) _updateRec(index int, x1, x2, y1, y2 P, lazy Id) {
	range4 := kd.closedRange[index]
	xmin, xmax, ymin, ymax := range4[0], range4[1], range4[2], range4[3]
	if x2 <= xmin || xmax < x1 || y2 <= ymin || ymax < y1 {
		return
	}
	if x1 <= xmin && xmax < x2 && y1 <= ymin && ymax < y2 {
		kd._propagate(index, lazy)
		return
	}
	kd._pushDown(index)
	kd._updateRec(index<<1, x1, x2, y1, y2, lazy)
	kd._updateRec(index<<1|1, x1, x2, y1, y2, lazy)
}

func reArrage(nums []P, order []int) []P {
	res := make([]P, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}

func reArrage2(nums []Id, order []int) []Id {
	res := make([]Id, len(order))
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
