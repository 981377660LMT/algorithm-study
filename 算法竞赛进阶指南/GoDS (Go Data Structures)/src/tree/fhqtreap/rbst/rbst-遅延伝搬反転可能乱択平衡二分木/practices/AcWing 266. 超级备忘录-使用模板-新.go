// https://hitonanode.github.io/cplib-cpp/data_structure/lazy_rbst.hpp

// Api:
//  func NewRBST(n int) *RBST

//  func (rb *RBST) Append(e E)
//  func (rb *RBST) Pop(i int) E
//  func (rb *RBST) AppendLeft(e E)
//  func (rb *RBST) PopLeft() E
//  func (rb *RBST) Insert(i int, e E)
//  func (rb *RBST) Erase(start, end int)
//	func (rb *RBST) RotateRight(start, end, k int)
//  func (rb *RBST) Reverse(start, end int)
//  func (rb *RBST) ReverseAll()
//  func (rb *RBST) Get(i int) E
//  func (rb *RBST) Set(i int, e E)
//  func (rb *RBST) Query(start, end int) E
//  func (rb *RBST) QueryAll() E
//  func (rb *RBST) Update(start, end int, lazy Id)
//  func (rb *RBST) MaxRight(left int, f func(E) bool) int
//  func (rb *RBST) MinLeft(right int, f func(E) bool) int
//  func (rb *RBST) Size() int
//  func (rb *RBST) InOrder() []E

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

func main() {
	// https://www.acwing.com/problem/content/268/
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]E, n) // (sum, size, min)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		nums[i] = E{sum: x, size: 1, min: x}
	}

	// !区间更新：加上一个数，区间查询：区间最小值
	T := NewRBST(nums)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "ADD" {
			var left, right, add int
			fmt.Fscan(in, &left, &right, &add)
			left--
			T.Update(left, right, add)
		} else if op == "REVERSE" {
			var left, right int
			fmt.Fscan(in, &left, &right)
			left--
			T.Reverse(left, right)
		} else if op == "REVOLVE" {
			// 区间 轮转k次
			var left, right, k int
			fmt.Fscan(in, &left, &right, &k)
			left--
			T.RotateRight(left, right, k)
		} else if op == "INSERT" {
			// 在pos后插入val
			var pos, val int
			fmt.Fscan(in, &pos, &val)
			pos--
			T.Insert(pos+1, E{sum: val, size: 1, min: val})
		} else if op == "DELETE" {
			var pos int
			fmt.Fscan(in, &pos)
			pos--
			T.Pop(pos)
		} else if op == "MIN" {
			var left, right int
			fmt.Fscan(in, &left, &right)
			left--
			fmt.Fprintln(out, T.Query(left, right).min)
		}
	}
}

type E = struct{ sum, size, min int }
type Id = int

// toggle时翻转左右的行为
func (*RBST) rev(e E) E     { return e }
func (*RBST) id() Id        { return 0 }
func (*RBST) op(e1, e2 E) E { return E{e1.sum + e2.sum, e1.size + e2.size, min(e1.min, e2.min)} }
func (*RBST) mapping(f Id, e E) E {
	return E{e.sum + f*e.size, e.size, e.min + f}
}
func (*RBST) composition(f, g Id) Id { return f + g }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type RNode struct {
	left, right *RNode
	val, sum    E
	lazy        Id
	isReversed  bool
	sz          int
}

func (n *RNode) String() string {
	return fmt.Sprintf("{val: %v, sum: %v, size: %v}", n.val, n.sum, n.sz)
}

type RBST struct {
	x, y, z, w uint32
	data       []*RNode
	dptr       int
	root       *RNode
}

// Lazy randomized binary search tree
func NewRBST(nums []E) *RBST {
	res := &RBST{x: 123456789, y: 362436069, z: 521288629, w: 88675123, data: make([]*RNode, len(nums))}
	if len(nums) > 0 {
		res.root = res.build(0, len(nums), nums)
	}
	return res
}

// 0-indexed.Insert before pos.
func (rb *RBST) Insert(pos int, e E) {
	first, second := rb.split(rb.root, pos)
	rb.root = rb.merge(first, rb.merge(rb.alloc(e), second))
}

func (rb *RBST) Append(e E) {
	rb.Insert(rb.Size(), e)
}

func (rb *RBST) AppendLeft(e E) {
	rb.Insert(0, e)
}

// 0-indexed.Pop at pos.
func (rb *RBST) Pop(pos int) E {
	if pos < 0 {
		pos += rb.Size()
	}
	first, second := rb.split(rb.root, pos)
	res, third := rb.split(second, 1)
	rb.root = rb.merge(first, third)
	return res.val // TODO
}

func (rb *RBST) PopLeft() E {
	return rb.Pop(0)
}

// Remove [start, end) from list.
func (rb *RBST) Erase(start, end int) {
	start++
	var x, y, z *RNode
	y, z = rb.split(rb.root, end)
	x, y = rb.split(y, start-1)
	rb.root = rb.merge(x, z)
}

func (rb *RBST) Reverse(start, end int) {
	if start >= end {
		return
	}
	f1, s1 := rb.split(rb.root, end)
	f2, s2 := rb.split(f1, start)
	rb.toggle(s2)
	rb.root = rb.merge(rb.merge(f2, s2), s1)
}

func (rb *RBST) ReverseAll() { rb.toggle(rb.root) }

// Rotate [start, stop) to the right `k` times.
func (rb *RBST) RotateRight(start, stop, k int) {
	start++
	n := stop - start + 1 - k%(stop-start+1)
	var x, y, z, p *RNode
	x, y = rb.split(rb.root, start-1)
	y, z = rb.split(y, n)
	z, p = rb.split(z, stop-start+1-n)
	rb.root = rb.merge(rb.merge(rb.merge(x, z), y), p)
}

// 0-indexed.Query [start, end)
//  !start must be smaller than end.
func (rb *RBST) Query(start, end int) E {
	f1, s1 := rb.split(rb.root, start)
	f2, s2 := rb.split(s1, end-start)
	rb.push(f2)
	res := f2.sum
	rb.root = rb.merge(f1, rb.merge(f2, s2))
	return res
}

func (rb *RBST) QueryAll() E { return rb.root.sum }

func (rb *RBST) Update(start, end int, lazy Id) {
	if start >= end {
		return
	}
	f1, s1 := rb.split(rb.root, start)
	f2, s2 := rb.split(s1, end-start)
	rb.allApply(f2, lazy)
	rb.root = rb.merge(f1, rb.merge(f2, s2))
}

func (rb *RBST) Get(pos int) E {
	if pos < 0 {
		pos += rb.Size()
	}
	return rb.Query(pos, pos+1)
}

func (rb *RBST) Set(pos int, e E) {
	if pos < 0 {
		pos += rb.Size()
	}
	f1, s1 := rb.split(rb.root, pos)
	f2, s2 := rb.split(s1, 1)
	*f2 = *rb.alloc(e)
	rb.root = rb.merge(f1, rb.merge(f2, s2))
}

func (rb *RBST) Size() int { return rb.size(rb.root) }

// rbst.Query(0, i) が true となるような最大の i を返す．
//  i := rbst.MaxRight(e, func(v E) bool { return v.sum <= k })
//  単調性を仮定．atcoder::lazy_segtree と同じ．
//  e は単位元．
func (rb *RBST) MaxRight(e E, f func(E) bool) int {
	if rb.root == nil {
		return 0
	}
	rb.push(rb.root)
	now := rb.root
	prodNow := e
	res := 0
	for {
		if now.left != nil {
			rb.push(now.left)
			pl := rb.op(prodNow, now.left.sum)
			if f(pl) {
				prodNow = pl
				res += now.left.sz
			} else {
				now = now.left
				continue
			}
		}
		pl := rb.op(prodNow, now.val)
		if !f(pl) {
			return res
		}
		prodNow = pl
		res++
		if now.right == nil {
			return res
		}
		rb.push(now.right)
		now = now.right
	}
}

// rbst.Query(i, rbst.Size()) が true となるような最小の i を返す．
//  i := rbst.MinLeft(e, func(v E) bool { return v.sum >= k })
//  単調性を仮定．atcoder::lazy_segtree と同じ．
//  e は単位元．
func (rb *RBST) MinLeft(e E, f func(E) bool) int {
	if rb.root == nil {
		return 0
	}
	rb.push(rb.root)
	now := rb.root
	prodNow := e
	res := rb.size(rb.root)
	for {
		if now.right != nil {
			rb.push(now.right)
			pr := rb.op(now.right.sum, prodNow)
			if f(pr) {
				prodNow = pr
				res -= now.right.sz
			} else {
				now = now.right
				continue
			}
		}
		pr := rb.op(now.val, prodNow)
		if !f(pr) {
			return res
		}
		prodNow = pr
		res--
		if now.left == nil {
			return res
		}
		rb.push(now.left)
		now = now.left
	}

}

// Return all elements in index order.
func (rb *RBST) InOrder() []E {
	res := make([]E, 0, rb.Size())
	rb.inOrder(rb.root, &res)
	return res
}

func (rb *RBST) inOrder(root *RNode, res *[]E) {
	if root == nil {
		return
	}
	rb.push(root)
	rb.inOrder(root.left, res)
	*res = append(*res, root.val)
	rb.inOrder(root.right, res)
}

func (rb *RBST) String() string {
	nums := rb.InOrder()
	return fmt.Sprintf("rbst%v", nums)
}

// merge l and r, return new root
func (rb *RBST) merge(l, r *RNode) *RNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}

	if rb.rand()%uint32((l.sz+r.sz)) < uint32(l.sz) {
		rb.push(l)
		l.right = rb.merge(l.right, r)
		return rb.update(l)
	} else {
		rb.push(r)
		r.left = rb.merge(l, r.left)
		return rb.update(r)
	}
}

// split root to [0,k) and [k,n)
func (rb *RBST) split(root *RNode, k int) (*RNode, *RNode) {
	if root == nil {
		return nil, nil
	}
	rb.push(root)
	if k <= rb.size(root.left) {
		first, second := rb.split(root.left, k)
		root.left = second
		return first, rb.update(root)
	} else {
		first, second := rb.split(root.right, k-rb.size(root.left)-1)
		root.right = first
		return rb.update(root), second
	}
}

func (rb *RBST) update(t *RNode) *RNode {
	t.sz = 1
	t.sum = t.val
	if t.left != nil {
		t.sz += t.left.sz
		t.sum = rb.op(t.left.sum, t.sum)
	}
	if t.right != nil {
		t.sz += t.right.sz
		t.sum = rb.op(t.sum, t.right.sum)
	}
	return t
}

func (rb *RBST) allApply(t *RNode, f Id) *RNode {
	t.val = rb.mapping(f, t.val)
	t.sum = rb.mapping(f, t.sum)
	t.lazy = rb.composition(f, t.lazy)
	return t
}

func (rb *RBST) toggle(t *RNode) {
	tmp := t.left
	t.left = t.right
	t.right = tmp
	t.sum = rb.rev(t.sum)
	t.isReversed = !t.isReversed
}

func (rb *RBST) push(t *RNode) {
	if t.lazy != rb.id() {
		if t.left != nil {
			rb.allApply(t.left, t.lazy)
		}
		if t.right != nil {
			rb.allApply(t.right, t.lazy)
		}
		t.lazy = rb.id()
	}
	if t.isReversed {
		if t.left != nil {
			rb.toggle(t.left)
		}
		if t.right != nil {
			rb.toggle(t.right)
		}
		t.isReversed = false
	}
}

func (rb *RBST) alloc(v E) *RNode {
	if rb.dptr >= len(rb.data) {
		rb.resize(1 + len(rb.data)*2)
	}
	res := &RNode{val: v, sum: v, sz: 1, lazy: rb.id()}
	rb.data[rb.dptr] = res
	rb.dptr++
	return res
}

func (rb *RBST) resize(n int) {
	newData := make([]*RNode, n)
	copy(newData, rb.data)
	rb.data = newData
}

func (rb *RBST) build(l, r int, nums []E) *RNode {
	if r-l == 1 {
		t := rb.alloc(nums[l])
		return rb.update(t)
	}
	return rb.update(rb.merge(rb.build(l, (l+r)/2, nums), rb.build((l+r)/2, r, nums)))
}

func (rb *RBST) size(node *RNode) int {
	if node == nil {
		return 0
	}
	return node.sz
}

func (rb *RBST) rand() uint32 {
	t := rb.x ^ (rb.x << 11)
	rb.x, rb.y, rb.z = rb.y, rb.z, rb.w
	rb.w = (rb.w ^ (rb.w >> 19)) ^ (t ^ (t >> 8))
	return rb.w
}
