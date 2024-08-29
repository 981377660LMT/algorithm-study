package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const INF32 int = 1e9 + 10

// No.2242 Cities and Teleporters-爬山
// https://yukicoder.me/problems/no/2242
// 给定N座山，高度为H[i].
// 每座山可以到达高度不超过T[i]的任意一座山.
// 给定q个询问，每次询问两座山a,b，问最少需要多少次操作可以从a到b.
// 如果无法到达，输出-1.
//
// !相当于通过增加T的值来到达目标位置(提升T[i]的等级).
// !每个等级处理出，可以到达的下一个最高等级.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int32
	fmt.Fscan(in, &N)
	H := make([]int, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &H[i])
	}
	T := make([]int, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &T[i])
	}

	// 维护每个高度的最大的T[i].
	seg := NewDynamicSegTreeSparse(0, int(INF32), false)
	root := seg.NewRoot()
	for i := int32(0); i < N; i++ {
		root = seg.Update(root, H[i], T[i])
	}

	newT, origin := Discretize(T)
	m := int32(len(origin))
	D := NewDoubling(m, int(m), func() int { return 0 }, func(e1, e2 int) int { return e1 + e2 })
	for fromLevel := int32(0); fromLevel < m; fromLevel++ {
		x := origin[fromLevel]
		y := seg.Query(root, 0, x+1)
		y = max(x, y)
		toLevel := sort.SearchInts(origin, y)
		D.Add(fromLevel, int32(toLevel), toLevel-int(fromLevel)) // 从fromLevel到toLevel增加的级别.
	}
	D.Build()

	query := func(a, b int32) int32 {
		if a == b {
			return 0
		}
		s, t := newT[a], sort.SearchInts(origin, H[b])
		if s >= t {
			return 1
		}
		check := func(i int32, e int) bool { return e < t-s }
		step, _, _ := D.LastTrue(int32(s), check)
		res := int32(step + 2)
		if res > N {
			res = -1
		}
		return res
	}

	var Q int32
	fmt.Fscan(in, &Q)
	for i := int32(0); i < Q; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)
		a--
		b--
		fmt.Fprintln(out, query(a, b))
	}
}

// PointAddRangeMax
type E = int

func e() E        { return -INF32 }
func op(a, b E) E { return max(a, b) }

type DynamicSegTreeSparse struct {
	L, R       int
	persistent bool
	unit       E
}

type SegNode struct {
	idx       int
	l, r      *SegNode
	data, sum E
}

// 指定 [left,right) 区间建立动态开点线段树.
func NewDynamicSegTreeSparse(left, right int, persistent bool) *DynamicSegTreeSparse {
	return &DynamicSegTreeSparse{
		L:          left,
		R:          right + 5,
		persistent: persistent,
		unit:       e(),
	}
}

func (ds *DynamicSegTreeSparse) NewRoot() *SegNode { return nil }

// 查询区间 [left, right).
// L<=left<=right<=R
func (ds *DynamicSegTreeSparse) Query(root *SegNode, left, right int) E {
	if left == right {
		return ds.unit
	}
	x := ds.unit
	ds._queryRec(root, ds.L, ds.R, left, right, &x)
	return x
}

func (ds *DynamicSegTreeSparse) QueryAll(root *SegNode) E {
	return ds.Query(root, ds.L, ds.R)
}

// L<=index<R
func (ds *DynamicSegTreeSparse) Set(root *SegNode, index int, value E) *SegNode {
	return ds._setRec(root, ds.L, ds.R, index, value)
}

func (ds *DynamicSegTreeSparse) Get(root *SegNode, index int) E {
	return ds._getRec(root, index)
}

// L<=left<R
func (ds *DynamicSegTreeSparse) Update(root *SegNode, index int, value E) *SegNode {
	return ds._updateRec(root, ds.L, ds.R, index, value)
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 check.
// L<=right<=R
func (ds *DynamicSegTreeSparse) MinLeft(root *SegNode, right int, check func(E) bool) int {
	x := ds.unit
	return ds._minLeftRec(root, ds.L, ds.R, right, check, &x)
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 check.
// L<=left<=R
func (ds *DynamicSegTreeSparse) MaxRight(root *SegNode, left int, check func(E) bool) int {
	x := ds.unit
	return ds._maxRightRec(root, ds.L, ds.R, left, check, &x)
}

func (ds *DynamicSegTreeSparse) GetAll(root *SegNode) []struct {
	index int
	value E
} {
	res := make([]struct {
		index int
		value E
	}, 0)
	ds._getAllRec(root, &res)
	return res
}

func (ds *DynamicSegTreeSparse) Copy(node *SegNode) *SegNode {
	if node == nil || !ds.persistent {
		return node
	}
	return &SegNode{idx: node.idx, l: node.l, r: node.r, data: node.data, sum: node.sum}
}

func (ds *DynamicSegTreeSparse) _pushUp(node *SegNode) {
	node.sum = node.data
	if node.l != nil {
		node.sum = op(node.l.sum, node.sum)
	}
	if node.r != nil {
		node.sum = op(node.sum, node.r.sum)
	}
}

func (ds *DynamicSegTreeSparse) _newNode(idx int, x E) *SegNode {
	return &SegNode{idx: idx, data: x, sum: x}
}

func (ds *DynamicSegTreeSparse) _setRec(root *SegNode, l, r, i int, x E) *SegNode {
	if root == nil {
		root = ds._newNode(i, x)
		return root
	}
	root = ds.Copy(root)
	if root.idx == i {
		root.data = x
		ds._pushUp(root)
		return root
	}
	m := (l + r) >> 1
	if i < m {
		if root.idx < i {
			root.idx, i = i, root.idx
			root.data, x = x, root.data
		}
		root.l = ds._setRec(root.l, l, m, i, x)
	} else {
		if i < root.idx {
			root.idx, i = i, root.idx
			root.data, x = x, root.data
		}
		root.r = ds._setRec(root.r, m, r, i, x)
	}
	ds._pushUp(root)
	return root
}

func (ds *DynamicSegTreeSparse) _updateRec(root *SegNode, l, r, i int, x E) *SegNode {
	if root == nil {
		root = ds._newNode(i, x)
		return root
	}
	root = ds.Copy(root)
	if root.idx == i {
		root.data = op(root.data, x)
		ds._pushUp(root)
		return root
	}
	m := (l + r) >> 1
	if i < m {
		if root.idx < i {
			root.idx, i = i, root.idx
			root.data, x = x, root.data
		}
		root.l = ds._updateRec(root.l, l, m, i, x)
	} else {
		if i < root.idx {
			root.idx, i = i, root.idx
			root.data, x = x, root.data
		}
		root.r = ds._updateRec(root.r, m, r, i, x)
	}
	ds._pushUp(root)
	return root
}

func (ds *DynamicSegTreeSparse) _queryRec(root *SegNode, l, r, ql, qr int, x *E) {
	ql = max(ql, l)
	qr = min(qr, r)
	if ql >= qr || root == nil {
		return
	}
	if l == ql && r == qr {
		*x = op(*x, root.sum)
		return
	}
	m := (l + r) >> 1
	ds._queryRec(root.l, l, m, ql, qr, x)
	if ql <= root.idx && root.idx < qr {
		*x = op(*x, root.data)
	}
	ds._queryRec(root.r, m, r, ql, qr, x)
}

func (ds *DynamicSegTreeSparse) _minLeftRec(root *SegNode, l, r, qr int, check func(E) bool, x *E) int {
	if root == nil || qr <= l {
		return ds.L
	}
	if check(op(root.sum, *x)) {
		*x = op(root.sum, *x)
		return ds.L
	}
	m := (l + r) >> 1
	k := ds._minLeftRec(root.r, m, r, qr, check, x)
	if k != ds.L {
		return k
	}
	if root.idx < qr {
		*x = op(root.data, *x)
		if !check(*x) {
			return root.idx + 1
		}
	}
	return ds._minLeftRec(root.l, l, m, qr, check, x)
}

func (ds *DynamicSegTreeSparse) _maxRightRec(root *SegNode, l, r, ql int, check func(E) bool, x *E) int {
	if root == nil || r <= ql {
		return ds.R
	}
	if check(op(*x, root.sum)) {
		*x = op(*x, root.sum)
		return ds.R
	}
	m := (l + r) >> 1
	k := ds._maxRightRec(root.l, l, m, ql, check, x)
	if k != ds.R {
		return k
	}
	if ql <= root.idx {
		*x = op(*x, root.data)
		if !check(*x) {
			return root.idx
		}
	}
	return ds._maxRightRec(root.r, m, r, ql, check, x)
}

func (ds *DynamicSegTreeSparse) _getAllRec(root *SegNode, res *[]struct {
	index int
	value E
}) {
	if root == nil {
		return
	}
	ds._getAllRec(root.l, res)
	*res = append(*res, struct {
		index int
		value E
	}{root.idx, root.data})
	ds._getAllRec(root.r, res)
}

func (ds *DynamicSegTreeSparse) _getRec(root *SegNode, idx int) E {
	if root == nil {
		return ds.unit
	}
	if idx == root.idx {
		return root.data
	}
	if idx < root.idx {
		return ds._getRec(root.l, idx)
	}
	return ds._getRec(root.r, idx)
}

type Doubling[E any] struct {
	prepared bool
	n        int32
	log      int32
	to       []int32

	// 边权
	dp []E
	e  func() E
	op func(e1, e2 E) E
}

func NewDoubling[E any](n int32, maxStep int, e func() E, op func(e1, e2 E) E) *Doubling[E] {
	res := &Doubling[E]{e: e, op: op}
	res.n = n
	res.log = int32(bits.Len(uint(maxStep)))
	size := n * res.log
	res.to = make([]int32, size)
	res.dp = make([]E, size)
	for i := int32(0); i < size; i++ {
		res.to[i] = -1
		res.dp[i] = res.e()
	}
	return res
}

// 初始状态(leaves):从 `from` 状态到 `to` 状态，边权为 `weight`.
//
//	0 <= from, to < n
func (d *Doubling[E]) Add(from, to int32, weight E) {
	if d.prepared {
		panic("Doubling is prepared")
	}
	if to < -1 || to >= d.n {
		panic("to is out of range")
	}

	d.to[from] = to
	d.dp[from] = weight
}

func (d *Doubling[E]) Build() {
	if d.prepared {
		panic("Doubling is prepared")
	}

	d.prepared = true
	n := d.n
	for k := int32(0); k < d.log-1; k++ {
		for v := int32(0); v < n; v++ {
			w := d.to[k*n+v]
			next := (k+1)*n + v
			if w == -1 {
				d.to[next] = -1
				d.dp[next] = d.dp[k*n+v]
				continue
			}
			d.to[next] = d.to[k*n+w]
			d.dp[next] = d.op(d.dp[k*n+v], d.dp[k*n+w])
		}
	}
}

// 从 `from` 状态开始，执行 `step` 次操作，返回最终状态的编号和操作的结果。
//
//	0 <= from < n
//	如果最终状态不存在，返回 -1, e()
func (d *Doubling[E]) Jump(from int32, step int) (to int32, res E) {
	if !d.prepared {
		panic("Doubling is not prepared")
	}
	if step >= 1<<d.log {
		panic("step is over max step")
	}

	res = d.e()
	to = from
	for k := int32(0); k < d.log; k++ {
		if to == -1 {
			break
		}
		if step&(1<<k) != 0 {
			pos := k*d.n + to
			res = d.op(res, d.dp[pos])
			to = d.to[pos]
		}
	}
	return
}

// 求从 `from` 状态开始转移，满足 `check` 为 `true` 的最小的 `step` 以及最终状态的编号和操作的结果。
// 如果不存在，则返回 (-1, -1, e()).
func (d *Doubling[E]) FirstTrue(from int32, check func(next int32, weight E) bool) (step int, to int32, res E) {
	if !d.prepared {
		panic("Doubling is not prepared")
	}

	if e := d.e(); check(from, e) {
		return 0, from, e
	}

	res = d.e()
	for k := d.log - 1; k >= 0; k-- {
		pos := k*d.n + from
		tmp := d.to[pos]
		if tmp == -1 {
			continue
		}
		next := d.op(res, d.dp[pos])
		if !check(tmp, next) {
			step |= 1 << k
			from = tmp
			res = next
		}
	}

	p := d.to[from]
	if p == -1 {
		return -1, -1, d.e()
	} else {
		step++
		to = p
		res = d.op(res, d.dp[from])
	}
	return
}

// 求从 `from` 状态开始转移，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号和操作的结果。
// 如果不存在，则返回 (-1, -1, e()).
func (d *Doubling[E]) LastTrue(from int32, check func(next int32, weight E) bool) (step int, to int32, res E) {
	if !d.prepared {
		panic("Doubling is not prepared")
	}

	if e := d.e(); !check(from, e) {
		return -1, -1, e
	}

	res = d.e()
	for k := d.log - 1; k >= 0; k-- {
		pos := k*d.n + from
		tmp := d.to[pos]
		if tmp == -1 {
			continue
		}
		next := d.op(res, d.dp[pos])
		if check(tmp, next) {
			step |= 1 << k
			from = tmp
			res = next
		}
	}

	to = from
	return
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

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int) (newNums []int, origin []int) {
	newNums = make([]int, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = len(origin) - 1
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}
