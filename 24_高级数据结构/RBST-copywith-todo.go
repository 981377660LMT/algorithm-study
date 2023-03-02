// 有问题
// https://atcoder.jp/contests/arc030/submissions/24709914

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]E, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i].sum)
		nums[i].size = 1
	}

	rbst := NewRBST(n)
	tree := rbst.NewTree()
	rbst.Build(tree, nums)
	fmt.Println(rbst.Query(tree, 0, 1), 999)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var a, b, v int
			fmt.Fscan(in, &a, &b, &v)
			rbst.Update(tree, a-1, b, v)
		} else if op == 2 {
			var a, b, c, d int
			fmt.Fscan(in, &a, &b, &c, &d)
			rbst.CopyWith(tree, a-1, c-1, b-a+1)
		} else if op == 3 {
			var a, b int
			fmt.Fscan(in, &a, &b)
			fmt.Fprintln(out, rbst.Query(tree, a-1, b).sum)
		}
	}
}

// func main() {
// 	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
// 	leaves := make([]E, len(nums))
// 	for i := 0; i < len(nums); i++ {
// 		leaves[i] = E{nums[i], 1}
// 	}
// 	rb := NewRBST(10)
// 	root := rb.NewTree()
// 	rb.Build(root, leaves)
// 	fmt.Println(rb)
// 	// rotate right
// 	rb.RotateRight(root, 1, 4, 2)
// 	fmt.Println(rb)
// 	rb.Erase(root, 1, 3)
// 	fmt.Println(rb)
// 	fmt.Println(rb.QueryAll(root))
// 	fmt.Println(root)
// }

type E = struct{ sum, size int }
type Id = int

// toggle时翻转左右的行为
func (*RBST) rev(e E) E              { return e }
func (*RBST) id() Id                 { return 0 }
func (*RBST) op(e1, e2 E) E          { return E{e1.sum + e2.sum, e1.size + e2.size} }
func (*RBST) mapping(f Id, e E) E    { return E{e.sum + e.size*f, e.size} }
func (*RBST) composition(f, g Id) Id { return f + g }

type RNode struct {
	left, right *RNode
	val, sum    E
	lazy        Id
	isReversed  bool
	sz          int
}

type RBST struct {
	x, y, z, w uint32
	data       []*RNode
}

// Lazy randomized binary search tree
func NewRBST(initCapacity int) *RBST {
	res := &RBST{x: 123456789, y: 362436069, z: 521288629, w: 88675123, data: make([]*RNode, 0, initCapacity)}
	return res
}

func (rb *RBST) Build(root *RNode, nums []E) {
	*root = *rb.build(0, len(nums), nums)
}

func (rb *RBST) NewTree() *RNode { return &RNode{} }

// copy from [start, start+len) to [target, )
func (rb *RBST) CopyWith(root *RNode, target, start, len int) {
	f1, s1 := rb.split(root, start)
	f2, s2 := rb.split(s1, len)
	*root = *rb.merge(f1, rb.merge(f2, s2))
	f3, s3 := rb.split(root, target)
	_, s4 := rb.split(s3, len)
	*root = *rb.merge(f3, rb.merge(f2, s4))
}

// 0-indexed.Insert before pos.
func (rb *RBST) Insert(root *RNode, pos int, e E) {
	first, second := rb.split(root, pos)
	*root = *rb.merge(first, rb.merge(rb.alloc(e), second))
}

func (rb *RBST) Append(root *RNode, e E) {
	rb.Insert(root, rb.Size(root), e)
}

func (rb *RBST) AppendLeft(root *RNode, e E) {
	rb.Insert(root, 0, e)
}

// 0-indexed.Pop at pos.
func (rb *RBST) Pop(root *RNode, pos int) E {
	if pos < 0 {
		pos += rb.Size(root)
	}
	first, second := rb.split(root, pos)
	res, third := rb.split(second, 1)
	*root = *rb.merge(first, third)
	return res.val // TODO
}

func (rb *RBST) PopLeft(root *RNode) E {
	return rb.Pop(root, 0)
}

// Remove [start, end) from list.
func (rb *RBST) Erase(root *RNode, start, end int) {
	start++
	var x, y, z *RNode
	y, z = rb.split(root, end)
	x, y = rb.split(y, start-1)
	*root = *rb.merge(x, z)
}

func (rb *RBST) Reverse(root *RNode, start, end int) {
	if start >= end {
		return
	}
	f1, s1 := rb.split(root, end)
	f2, s2 := rb.split(f1, start)
	rb.toggle(s2)
	*root = *rb.merge(rb.merge(f2, s2), s1)
}

func (rb *RBST) ReverseAll(root *RNode) { rb.toggle(root) }

// Rotate [start, stop) to the right `k` times.
func (rb *RBST) RotateRight(root *RNode, start, stop, k int) {
	start++
	n := stop - start + 1 - k%(stop-start+1)
	var x, y, z, p *RNode
	x, y = rb.split(root, start-1)
	y, z = rb.split(y, n)
	z, p = rb.split(z, stop-start+1-n)
	*root = *rb.merge(rb.merge(rb.merge(x, z), y), p)
}

// 0-indexed.Query [start, end)
//  !start must be smaller than end.
func (rb *RBST) Query(root *RNode, start, end int) E {
	f1, s1 := rb.split(root, start)
	f2, s2 := rb.split(s1, end-start)
	rb.push(f2)
	res := f2.sum
	*root = *rb.merge(f1, rb.merge(f2, s2))
	return res
}

func (rb *RBST) QueryAll(root *RNode) E { return root.sum }

func (rb *RBST) Update(root *RNode, start, end int, lazy Id) {
	if start >= end {
		return
	}
	f1, s1 := rb.split(root, start)
	f2, s2 := rb.split(s1, end-start)
	rb.allApply(f2, lazy)
	*root = *rb.merge(f1, rb.merge(f2, s2))
}

func (rb *RBST) Get(root *RNode, pos int) E {
	if pos < 0 {
		pos += rb.Size(root)
	}
	return rb.Query(root, pos, pos+1)
}

func (rb *RBST) Set(root *RNode, pos int, e E) {
	if pos < 0 {
		pos += rb.Size(root)
	}
	f1, s1 := rb.split(root, pos)
	f2, s2 := rb.split(s1, 1)
	rb.duplicateNode(f2)
	*f2 = *rb.newNode(e)
	*root = *rb.merge(f1, rb.merge(f2, s2))
}

func (rb *RBST) Size(root *RNode) int { return rb.size(root) }

// merge l and r, return new root
func (rb *RBST) merge(l, r *RNode) *RNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}

	if rb.nextRand()%uint32((l.sz+r.sz)) < uint32(l.sz) {
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
	rb.duplicateNode(t)
	if t.lazy != rb.id() {
		if t.left != nil {
			rb.duplicateNode(t.left)
			rb.allApply(t.left, t.lazy)
		}
		if t.right != nil {
			rb.duplicateNode(t.right)
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
	rb.data = append(rb.data, res)
	return res
}

func (rb *RBST) newNode(v E) *RNode {
	return &RNode{val: v, sum: v, sz: 1, lazy: rb.id()}
}

func (rb *RBST) build(l, r int, nums []E) *RNode {
	if r-l == 1 {
		t := rb.alloc(nums[l])
		return rb.update(t)
	}
	return rb.update(rb.merge(rb.build(l, (l+r)/2, nums), rb.build((l+r)/2, r, nums)))
}

func (rb *RBST) duplicateNode(t *RNode) {
	if t == nil {
		return
	}
	newNode := &RNode{val: t.val, sum: t.sum, sz: t.sz, lazy: t.lazy, isReversed: t.isReversed}
	rb.data = append(rb.data, newNode)
	*t = *newNode
}

func (rb *RBST) size(node *RNode) int {
	if node == nil {
		return 0
	}
	return node.sz
}

func (rb *RBST) nextRand() uint32 {
	t := rb.x ^ (rb.x << 11)
	rb.x, rb.y, rb.z = rb.y, rb.z, rb.w
	rb.w = (rb.w ^ (rb.w >> 19)) ^ (t ^ (t >> 8))
	return rb.w
}
