// MultisetSum
// api:
//   NewSegmentTreeDynamicMultiSet() *SegmentTreeDynamicMultiSet
//   NewRoot() *SegNode
//   Add(node *SegNode, value, count int) *SegNode
//   Query(node *SegNode, start, end int) E
//   QueryAll(node *SegNode) E
//   PrefixKth(node *SegNode, k int) (int, int)
//   SuffixKth(node *SegNode, k int) (int, int)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()
	// abc241d()
	abc281e()
}

// https://atcoder.jp/contests/abc241/tasks/abc241_d
func abc241d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	S := NewSegmentTreeDynamicMultiSet()
	root := S.NewRoot()
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var x int
			fmt.Fscan(in, &x)
			root = S.Add(root, x, 1)
		}
		if t == 2 {
			var x, k int
			fmt.Fscan(in, &x, &k)
			e := S.Query(root, -INF, x+1)
			if e.count < k {
				fmt.Fprintln(out, -1)
			} else {
				res, _ := S.PrefixKth(root, e.count-k)
				fmt.Fprintln(out, res)
			}
		}
		if t == 3 {
			var x, k int
			fmt.Fscan(in, &x, &k)
			e := S.Query(root, x, INF+1)
			if e.count < k {
				fmt.Fprintln(out, -1)
			} else {
				a, _ := S.SuffixKth(root, e.count-k)
				fmt.Fprintln(out, a)
			}
		}
	}
}

// https://atcoder.jp/contests/abc281/tasks/abc281_e
func abc281e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	S := NewSegmentTreeDynamicMultiSet()
	root := S.NewRoot()
	res := make([]int, 0, n-m+1)
	for i := 0; i < m-1; i++ {
		root = S.Add(root, nums[i], 1)
	}
	for i := m - 1; i < n; i++ {
		root = S.Add(root, nums[i], 1)
		_, sum := S.PrefixKth(root, k)
		res = append(res, sum)
		root = S.Add(root, nums[i-m+1], -1)
	}
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

func demo() {
	S := NewSegmentTreeDynamicMultiSet()
	root := S.NewRoot()
	root = S.Add(root, 1, 1)
	root = S.Add(root, 1000000000, 100)
	fmt.Println(S.Query(root, 0, 10000000001))
	fmt.Println(S.PrefixKth(root, 1))
	fmt.Println(S.SuffixKth(root, 100))
}

const INF int = 1e18

type SegmentTreeDynamicMultiSet struct {
	seg *DynamicSegTreeSparse
}

func NewSegmentTreeDynamicMultiSet() *SegmentTreeDynamicMultiSet {
	return &SegmentTreeDynamicMultiSet{
		seg: NewDynamicSegTreeSparse(-INF, INF+1, false),
	}
}

func (s *SegmentTreeDynamicMultiSet) NewRoot() *SegNode {
	return s.seg.NewRoot()
}

func (s *SegmentTreeDynamicMultiSet) Add(node *SegNode, value, count int) *SegNode {
	return s.seg.Update(node, value, E{count: count, sum: value * count})
}

func (s *SegmentTreeDynamicMultiSet) Query(node *SegNode, start, end int) E {
	return s.seg.Query(node, start, end)
}

func (s *SegmentTreeDynamicMultiSet) QueryAll(node *SegNode) E {
	return s.seg.QueryAll(node)
}

// 返回前缀第 k 小的值，前k个数的和.
// 如果 k 大于当前节点的总数，返回 (INF, sum).
func (s *SegmentTreeDynamicMultiSet) PrefixKth(node *SegNode, k int) (int, int) {
	e := s.seg.QueryAll(node)
	if k >= e.count {
		return INF, e.sum
	}
	return s.prefixKth(node, k)
}

// 返回后缀第 k 小的值，后k个数的和.
// 如果 k 大于当前节点的总数，返回 (-INF, sum).
func (s *SegmentTreeDynamicMultiSet) SuffixKth(node *SegNode, k int) (int, int) {
	e := s.seg.QueryAll(node)
	count, sum := e.count, e.sum
	if k >= count {
		return -INF, sum
	}
	a, b := s.prefixKth(node, count-1-k)
	return a, sum - b - a
}

func (s *SegmentTreeDynamicMultiSet) prefixKth(node *SegNode, k int) (int, int) {
	key := s.seg.MaxRight(node, -INF, func(e E) bool { return e.count <= k })
	e := s.seg.Query(node, -INF, key)
	count, sum := e.count, e.sum
	return key, sum + key*(k-count)
}

// PointAddRangeSum
type E = struct {
	count int
	sum   int
}

func e() E { return E{} }
func op(a, b E) E {
	a.count += b.count
	a.sum += b.sum
	return a
}

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
