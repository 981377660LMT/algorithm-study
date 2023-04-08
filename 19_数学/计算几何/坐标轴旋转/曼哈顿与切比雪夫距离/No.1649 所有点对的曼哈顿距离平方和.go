// https://yukicoder.me/problems/no/1649
// No.1649 所有点对的曼哈顿距离平方和 mod 998244353
// https://yukicoder.me/problems/no/1649/editorial
// https://maspypy.github.io/library/test/yukicoder/1649_2.test.cpp

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	points := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &points[i][0], &points[i][1])
	}
	fmt.Fprintln(out, ManhattanSquare(points))
}

const MOD int = 998244353
const N = 1e9 + 10

var INV2 = Pow(2, MOD-2, MOD)

// !ManhattanSquare 计算所有点对的曼哈顿距离平方和 (0<=i<j<n)
func ManhattanSquare(points [][2]int) int {

	res := 0
	for t_ := 0; t_ < 4; t_++ {
		for i := range points {
			points[i][0], points[i][1] = -points[i][1], points[i][0]
		}
		seg := NewDynamicSegTreeSparse(-N, N, false)
		root := seg.NewRoot()
		sort.Slice(points, func(i, j int) bool {
			if points[i][0] == points[j][0] {
				return points[i][1] < points[j][1]
			}
			return points[i][0] < points[j][0]
		})
		for _, p := range points {
			x, y := p[0], p[1]
			x2 := (x + y) * (x + y) % MOD
			x1 := (x + y) % MOD
			x0 := 1
			s := seg.Query(root, -N, y)
			res += (x2*s[0] - 2*x1*s[1] + x0*s[2]) % MOD
			res %= MOD
			root = seg.Update(root, y, [3]int{x0, x1, x2})
		}
	}

	res = ((res*INV2)%MOD + MOD) % MOD
	return res
}

// PointAddRangeSum

type E = [3]int

func e() E { return [3]int{0, 0, 0} }
func op(a, b E) E {
	return [3]int{
		((a[0]+b[0])%MOD + MOD) % MOD,
		((a[1]+b[1])%MOD + MOD) % MOD,
		((a[2]+b[2])%MOD + MOD) % MOD,
	}
}

type DynamicSegTreeSparse struct {
	L, R       int
	persistent bool
	unit       E
}
type SegNode struct {
	idx     int
	l, r    *SegNode
	x, prod E
}

// 指定 [left,right) 区间建立动态开点线段树.
func NewDynamicSegTreeSparse(left, right int, persistent bool) *DynamicSegTreeSparse {
	return &DynamicSegTreeSparse{
		L:          left,
		R:          right,
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
func (ds *DynamicSegTreeSparse) _pushUp(node *SegNode) {
	node.prod = node.x
	if node.l != nil {
		node.prod = op(node.l.prod, node.prod)
	}
	if node.r != nil {
		node.prod = op(node.prod, node.r.prod)
	}
}
func (ds *DynamicSegTreeSparse) _newNode(idx int, x E) *SegNode {
	return &SegNode{idx: idx, x: x, prod: x}
}
func (ds *DynamicSegTreeSparse) _copyNode(node *SegNode) *SegNode {
	if node == nil || !ds.persistent {
		return node
	}
	return &SegNode{idx: node.idx, l: node.l, r: node.r, x: node.x, prod: node.prod}
}
func (ds *DynamicSegTreeSparse) _setRec(root *SegNode, l, r, i int, x E) *SegNode {
	if root == nil {
		root = ds._newNode(i, x)
		return root
	}
	root = ds._copyNode(root)
	if root.idx == i {
		root.x = x
		ds._pushUp(root)
		return root
	}
	m := (l + r) >> 1
	if i < m {
		if root.idx < i {
			root.idx, i = i, root.idx
			root.x, x = x, root.x
		}
		root.l = ds._setRec(root.l, l, m, i, x)
	} else {
		if i < root.idx {
			root.idx, i = i, root.idx
			root.x, x = x, root.x
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
	root = ds._copyNode(root)
	if root.idx == i {
		root.x = op(root.x, x)
		ds._pushUp(root)
		return root
	}
	m := (l + r) >> 1
	if i < m {
		if root.idx < i {
			root.idx, i = i, root.idx
			root.x, x = x, root.x
		}
		root.l = ds._updateRec(root.l, l, m, i, x)
	} else {
		if i < root.idx {
			root.idx, i = i, root.idx
			root.x, x = x, root.x
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
		*x = op(*x, root.prod)
		return
	}
	m := (l + r) >> 1
	ds._queryRec(root.l, l, m, ql, qr, x)
	if ql <= root.idx && root.idx < qr {
		*x = op(*x, root.x)
	}
	ds._queryRec(root.r, m, r, ql, qr, x)
}
func (ds *DynamicSegTreeSparse) _minLeftRec(root *SegNode, l, r, qr int, check func(E) bool, x *E) int {
	if root == nil || qr <= l {
		return ds.L
	}
	if check(op(root.prod, *x)) {
		*x = op(root.prod, *x)
		return ds.L
	}
	m := (l + r) >> 1
	k := ds._minLeftRec(root.r, m, r, qr, check, x)
	if k != ds.L {
		return k
	}
	if root.idx < qr {
		*x = op(root.x, *x)
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
	if check(op(*x, root.prod)) {
		*x = op(*x, root.prod)
		return ds.R
	}
	m := (l + r) >> 1
	k := ds._maxRightRec(root.l, l, m, ql, check, x)
	if k != ds.R {
		return k
	}
	if ql <= root.idx {
		*x = op(*x, root.x)
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
	}{root.idx, root.x})
	ds._getAllRec(root.r, res)
}
func (ds *DynamicSegTreeSparse) _getRec(root *SegNode, idx int) E {
	if root == nil {
		return ds.unit
	}
	if idx == root.idx {
		return root.x
	}
	if idx < root.idx {
		return ds._getRec(root.l, idx)
	}
	return ds._getRec(root.r, idx)
}

func abs(x int) int {
	if x < 0 {
		return -x
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

func Pow(base, exp, mod int) int {
	base %= mod
	res := 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}
