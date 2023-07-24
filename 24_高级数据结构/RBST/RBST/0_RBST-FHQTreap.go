// https://hitonanode.github.io/cplib-cpp/data_structure/lazy_rbst.hpp (有一些问题)
// https://nyaannyaan.github.io/library/rbst/rbst-base.hpp  (没问题)

// Api:
//  func NewRBST(n int) *RBST

//  func (rb *RBST) Append(e E)
//  func (rb *RBST) Pop(i int) E
//  func (rb *RBST) AppendLeft(e E)
//  func (rb *RBST) PopLeft() E
//  func (rb *RBST) Insert(i int, e E)
//  func (rb *RBST) Erase(start, end int)
//	func (rb *RBST) RotateRight(start, end, k int)
//  func (rb *RBST) RotateLeft(start, end, k int)
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
	"fmt"
	"time"
)

// RangeChmaxRangeMax
func lengthOfLIS(nums []int, k int) int {
	initNums := make([]E, 1e5+1)
	rbst := NewRBST(initNums)
	for _, num := range nums {
		preMax := rbst.Query(max(0, num-k), num)
		rbst.Update(num, num+1, preMax+1)
	}
	return rbst.QueryAll()
}

func main() {
	nums := []int{1, 3, 5, 4, 7}
	R := NewRBST(nums)
	fmt.Println(R.QueryAll())
	R.Set(0, 100)
	fmt.Println(R)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const INF int = 1e18

type E = int
type Id = int

// toggle时翻转左右的行为
// toggle时翻转左右的行为
func (*RBST) rev(e E) E     { return e }
func (*RBST) id() Id        { return -INF }
func (*RBST) op(e1, e2 E) E { return max(e1, e2) }

func (*RBST) mapping(f Id, e E) E {
	return max(f, e)
}
func (*RBST) composition(f, g Id) Id {
	return max(f, g)
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
	seed uint64
	root *RNode
}

// Lazy randomized binary search tree
func NewRBST(nums []E) *RBST {
	res := &RBST{seed: uint64(time.Now().UnixNano()/2 + 1)}
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
	p21, p22 := rb.split(rb.root, end)
	p11, p12 := rb.split(p21, start)
	rb.toggle(p12)
	rb.root = rb.merge(rb.merge(p11, p12), p22)
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

// Rotate [start, stop) to the left `k` times.
func (rb *RBST) RotateLeft(start, stop, k int) {
	start++
	k %= (stop - start + 1)
	var x, y, z, p *RNode
	x, y = rb.split(rb.root, start-1)
	y, z = rb.split(y, k)
	z, p = rb.split(z, stop-start+1-k)
	rb.root = rb.merge(rb.merge(rb.merge(x, z), y), p)
}

// 0-indexed.Query [start, end)
//  !start must be smaller than end.
func (rb *RBST) Query(start, end int) E {
	f1, s1 := rb.split(rb.root, start)
	f2, s2 := rb.split(s1, end-start)
	// rb.push(f2)  // TODO
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
	return rb.Query(pos, pos+1)
}

func (rb *RBST) Set(pos int, e E) {
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
	if l == nil || r == nil {
		if l == nil {
			return r
		}
		return l
	}

	if int(rb.nextRand()*(uint64(l.sz)+uint64(r.sz))>>32) < l.sz {
		rb.push(l)
		l.right = rb.merge(l.right, r)
		return rb.update(l)
	}
	rb.push(r)
	r.left = rb.merge(l, r.left)
	return rb.update(r)
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
	}
	first, second := rb.split(root.right, k-rb.size(root.left)-1)
	root.right = first
	return rb.update(root), second
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
	res := &RNode{val: v, sum: v, sz: 1, lazy: rb.id()}
	return res
}

func (rb *RBST) build(l, r int, nums []E) *RNode {
	if r-l == 1 {
		return rb.alloc(nums[l])
	}
	mid := (l + r) >> 1
	root := rb.alloc(nums[mid])
	if l < mid {
		root.left = rb.build(l, mid, nums)
	}
	if mid+1 < r {
		root.right = rb.build(mid+1, r, nums)
	}
	return rb.update(root)
}

func (rb *RBST) size(node *RNode) int {
	if node == nil {
		return 0
	}
	return node.sz
}

// XORShift
func (rb *RBST) nextRand() uint64 {
	rb.seed ^= rb.seed << 7
	rb.seed ^= rb.seed >> 9
	return rb.seed & 0xFFFFFFFF
}
