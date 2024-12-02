// 区间查询区间修改的 KD 树.
// Api:
//   NewKDTreeLazy(xs, ys []P, vs []E) *KDTreeLazy
//   (kd *KDTreeLazy) Set(i int32, v E)
//   (kd *KDTreeLazy) Update(i int32, v E)
//   (kd *KDTreeLazy) Query(xl, xr, yl, yr P) E
//   (kd *KDTreeLazy) QueryAll() E
//   (kd *KDTreeLazy) UpdateRange(xl, xr, yl, yr P, lazy Id)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/dynamic_point_set_rectangle_affine_rectangle_sum
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	xs, ys := make([]int, n), make([]int, n)
	vs := make([]E, n)
	for i := int32(0); i < n; i++ {
		var x, y, w int
		fmt.Fscan(in, &x, &y, &w)
		xs[i], ys[i], vs[i] = x, y, E{1, w}
	}

	type query = [7]int
	queries := make([]query, 0, q)
	for i := int32(0); i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 0 {
			var x, y, w int
			fmt.Fscan(in, &x, &y, &w)
			k := len(xs)
			xs, ys, vs = append(xs, x), append(ys, y), append(vs, E{})
			queries = append(queries, query{0, k, w, 0, 0, 0, 0})
		}
		if t == 1 {
			var k, w int
			fmt.Fscan(in, &k, &w)
			queries = append(queries, query{1, k, w, 0, 0, 0, 0})
		}
		if t == 2 {
			var a, b, c, d int
			fmt.Fscan(in, &a, &b, &c, &d)
			queries = append(queries, query{2, a, c, b, d, 0, 0})
		}
		if t == 3 {
			var a, b, c, d, p, q int
			fmt.Fscan(in, &a, &b, &c, &d, &p, &q)
			queries = append(queries, query{3, a, c, b, d, p, q})
		}
	}

	if len(xs) != len(ys) || len(xs) != len(vs) {
		panic("invalid input")
	}

	kdt := NewKDTreeLazy(xs, ys, vs)
	for _, q := range queries {
		t := q[0]
		if t == 0 {
			kdt.Set(int32(q[1]), E{1, q[2]})
		}
		if t == 1 {
			kdt.Set(int32(q[1]), E{1, q[2]})
		}
		if t == 2 {
			res := kdt.Query(q[1], q[2], q[3], q[4])
			fmt.Fprintln(out, res.sum2)
		}
		if t == 3 {
			kdt.UpdateRange(q[1], q[2], q[3], q[4], Id{mul: q[5], add: q[6]})
		}
	}
}

const MOD int = 998244353

type P = int // 点的坐标类型
const INF P = 1e18

type E = struct {
	sum1 int
	sum2 int
}
type Id = struct {
	mul int
	add int
}

func e() E   { return E{} }
func id() Id { return Id{mul: 1} }
func op(e1, e2 E) E {
	e1.sum1 = (e1.sum1 + e2.sum1) % MOD
	e1.sum2 = (e1.sum2 + e2.sum2) % MOD
	return e1
}
func mapping(f Id, x E, _ int32) E {
	res := E{}
	res.sum1 = x.sum1
	res.sum2 = (f.mul*x.sum2 + f.add*x.sum1) % MOD
	return res
}

func composition(f, g Id) Id {
	return Id{
		mul: f.mul * g.mul % MOD,
		add: (f.mul*g.add + f.add) % MOD,
	}
}

type KDTreeLazy struct {
	closedRange [][4]P
	data        []E
	lazy        []Id
	size        []int32
	pos         []int32 // raw data -> index
	n, log      int32
}

// 数据集中所有点的横坐标、纵坐标和对应的值.
func NewKDTreeLazy(xs, ys []P, vs []E) *KDTreeLazy {
	n := int32(len(xs))
	if n == 0 {
		panic("empty data")
	}
	log := int32(0)
	for (1 << log) < n {
		log++
	}
	data := make([]E, 1<<(log+1))
	lazy := make([]Id, 1<<log)
	for i := range lazy {
		lazy[i] = id()
	}
	closedRange := make([][4]P, 1<<(log+1))
	for i := range closedRange {
		closedRange[i] = [4]P{INF, -INF, INF, -INF}
	}
	size := make([]int32, 1<<(log+1))
	pos := make([]int32, n)
	ids := make([]int32, n)
	for i := int32(0); i < n; i++ {
		ids[i] = i
	}
	res := &KDTreeLazy{
		closedRange: closedRange,
		data:        data,
		lazy:        lazy,
		size:        size,
		pos:         pos,
		n:           n,
		log:         log,
	}
	res._build(1, xs, ys, vs, ids, true)
	return res
}

func (kd *KDTreeLazy) Set(i int32, v E) {
	i = kd.pos[i]
	for k := kd.log; k >= 1; k-- {
		kd._pushDown(i >> k)
	}
	kd.data[i] = v
	for i > 1 {
		i >>= 1
		kd.data[i] = op(kd.data[i<<1], kd.data[i<<1|1])
	}
}

func (kd *KDTreeLazy) Update(i int32, v E) {
	i = kd.pos[i]
	for k := kd.log; k >= 1; k-- {
		kd._pushDown(i >> k)
	}
	kd.data[i] = op(kd.data[i], v)
	for i > 1 {
		i >>= 1
		kd.data[i] = op(kd.data[i<<1], kd.data[i<<1|1])
	}
}

// [xl, xr) x [yl, yr)
func (kd *KDTreeLazy) Query(xl, xr, yl, yr P) E {
	if xr <= xl || yr <= yl {
		return e()
	}
	return kd._queryRec(1, xl, xr, yl, yr)
}

func (kd *KDTreeLazy) QueryAll() E {
	return kd.data[1]
}

// [xl, xr) x [yl, yr)
func (kd *KDTreeLazy) UpdateRange(xl, xr, yl, yr P, lazy Id) {
	if xr <= xl || yr <= yl {
		return
	}
	kd._updateRangeRec(1, xl, xr, yl, yr, lazy)
}

func (kd *KDTreeLazy) _pushDown(index int32) {
	lazy := kd.lazy[index]
	if lazy == id() {
		return
	}
	kd._propagate(index<<1, lazy)
	kd._propagate(index<<1|1, lazy)
	kd.lazy[index] = id()
}

func (kd *KDTreeLazy) _propagate(index int32, lazy Id) {
	kd.data[index] = mapping(lazy, kd.data[index], kd.size[index])
	if index < (1 << kd.log) {
		kd.lazy[index] = composition(lazy, kd.lazy[index])
	}
}

func (kd *KDTreeLazy) _build(idx int32, xs, ys []P, vs []E, ids []int32, divx bool) {
	n := int32(len(xs))
	kd.size[idx] = n
	range4 := &kd.closedRange[idx]
	xmin, xmax, ymin, ymax := &range4[0], &range4[1], &range4[2], &range4[3]
	*xmin, *ymin = INF, INF
	*xmax, *ymax = -INF, -INF

	for i := int32(0); i < n; i++ {
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

	if n == 1 {
		kd.data[idx] = vs[0]
		kd.pos[ids[0]] = idx
		return
	}

	m := n / 2
	order := make([]int32, n)
	for i := int32(0); i < n; i++ {
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

	xs, ys, vs, ids = reArrage(xs, order), reArrage(ys, order), reArrage(vs, order), reArrage(ids, order)
	kd._build(2*idx, xs[:m], ys[:m], vs[:m], ids[:m], !divx)
	kd._build(2*idx+1, xs[m:], ys[m:], vs[m:], ids[m:], !divx)
	kd.data[idx] = op(kd.data[idx<<1], kd.data[idx<<1|1])
}

func (kd *KDTreeLazy) _queryRec(index int32, x1, x2, y1, y2 P) E {
	if index >= int32(len(kd.closedRange)) {
		return e()
	}
	xmin, xmax, ymin, ymax := kd.closedRange[index][0], kd.closedRange[index][1], kd.closedRange[index][2], kd.closedRange[index][3]
	if xmin > xmax {
		return e()
	}
	if x2 <= xmin || xmax < x1 {
		return e()
	}
	if y2 <= ymin || ymax < y1 {
		return e()
	}
	if x1 <= xmin && xmax < x2 && y1 <= ymin && ymax < y2 {
		return kd.data[index]
	}
	kd._pushDown(index)
	return op(kd._queryRec(index<<1, x1, x2, y1, y2), kd._queryRec(index<<1|1, x1, x2, y1, y2))
}

func (kd *KDTreeLazy) _updateRangeRec(index int32, x1, x2, y1, y2 P, lazy Id) {
	if index >= int32(len(kd.closedRange)) {
		return
	}
	xmin, xmax, ymin, ymax := kd.closedRange[index][0], kd.closedRange[index][1], kd.closedRange[index][2], kd.closedRange[index][3]
	if xmin > xmax {
		return
	}
	if x2 <= xmin || xmax < x1 {
		return
	}
	if y2 <= ymin || ymax < y1 {
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

func reArrage[T any](arr []T, order []int32) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = arr[order[i]]
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
