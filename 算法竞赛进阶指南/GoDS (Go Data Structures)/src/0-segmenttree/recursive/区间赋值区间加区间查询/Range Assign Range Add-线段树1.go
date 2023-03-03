// 区间赋值 区间加 区间和查询

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	tree := NewLazySegmentTree(nums)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var x int
			fmt.Fscan(in, &x)
			tree.Assign(1, n, x) // 区间赋值
		} else if op == 2 {
			var i, x int
			fmt.Fscan(in, &i, &x)
			tree.Add(i, i, x) // 单点加
		} else {
			var i int
			fmt.Fscan(in, &i)
			fmt.Fprintln(out, tree.Query(i, i)) // 单点查询
		}
	}
}

type LazySegmentTree struct {
	n     int
	sum   []int // 幺元为0
	add   []int // 幺元为0
	cover []int // 幺元为-INF
}

func NewLazySegmentTree(leaves []int) *LazySegmentTree {
	cap := 1 << (bits.Len(uint(len(leaves)-1)) + 1)
	// !初始化data和lazy数组 然后建树
	cover := make([]int, cap)
	for i := 0; i < cap; i++ {
		cover[i] = -INF
	}
	tree := &LazySegmentTree{
		n:     len(leaves),
		sum:   make([]int, cap),
		add:   make([]int, cap),
		cover: cover,
	}
	tree._build(1, 1, tree.n, leaves)
	return tree
}

func (t *LazySegmentTree) _pushUp(root, left, right int) {
	t.sum[root] = t.sum[root<<1] + t.sum[root<<1|1] // !op
}

// 在执行pushdown函数时，要先下放cover标记，再下放add标记
func (t *LazySegmentTree) _pushDown(root, left, right int) {
	mid := (left + right) >> 1

	if t.cover[root] != -INF { // !monoid
		t._propagateCover(root<<1, left, mid, t.cover[root])
		t._propagateCover(root<<1|1, mid+1, right, t.cover[root])
	}

	if t.add[root] != 0 { // !monoid
		t._propagateAdd(root<<1, left, mid, t.add[root])
		t._propagateAdd(root<<1|1, mid+1, right, t.add[root])
	}

	t.cover[root] = -INF // !monoid
	t.add[root] = 0      // !monoid
}

// !mapping + composition
func (t *LazySegmentTree) _propagateCover(root, left, right, cover int) {
	t.sum[root] = cover * (right - left + 1)
	t.cover[root] = cover
	t.add[root] = 0
}

// !mapping + composition
func (t *LazySegmentTree) _propagateAdd(root, left, right, add int) {
	t.sum[root] += (right - left + 1) * add
	t.add[root] += add
}

func (t *LazySegmentTree) _build(root, left, right int, leaves []int) {
	if left == right {
		t.sum[root] = leaves[left-1]
		return
	}

	mid := (left + right) >> 1
	t._build(root<<1, left, mid, leaves)
	t._build(root<<1|1, mid+1, right, leaves)
	t._pushUp(root, left, right)
}

func (t *LazySegmentTree) _query(root, L, R, l, r int) int {
	if L <= l && r <= R {
		return t.sum[root]
	}

	t._pushDown(root, l, r)
	mid := (l + r) >> 1
	res := 0 // !monoid
	if L <= mid {
		res += t._query(root<<1, L, R, l, mid) // !op
	}
	if R > mid {
		res += t._query(root<<1|1, L, R, mid+1, r) // !op
	}
	return res
}

func (t *LazySegmentTree) _add(root, L, R, l, r, val int) {
	if L <= l && r <= R {
		t._propagateAdd(root, l, r, val)
		return
	}

	t._pushDown(root, l, r)
	mid := (l + r) >> 1
	if L <= mid {
		t._add(root<<1, L, R, l, mid, val)
	}
	if R > mid {
		t._add(root<<1|1, L, R, mid+1, r, val)
	}
	t._pushUp(root, l, r)
}

func (t *LazySegmentTree) _assign(root, L, R, l, r, val int) {
	if L <= l && r <= R {
		t._propagateCover(root, l, r, val)
		return
	}

	t._pushDown(root, l, r)
	mid := (l + r) >> 1
	if L <= mid {
		t._assign(root<<1, L, R, l, mid, val)
	}
	if R > mid {
		t._assign(root<<1|1, L, R, mid+1, r, val)
	}
	t._pushUp(root, l, r)
}

// public api

// !1 <= left <= right <= n
func (t *LazySegmentTree) Query(left, right int) int { return t._query(1, left, right, 1, t.n) }

// !1 <= left <= right <= n
func (t *LazySegmentTree) Add(left, right, delta int) { t._add(1, left, right, 1, t.n, delta) }

// !1 <= left <= right <= n
func (t *LazySegmentTree) Assign(left, right, value int) { t._assign(1, left, right, 1, t.n, value) }

func (t *LazySegmentTree) QueryAll() int { return t.sum[1] }
