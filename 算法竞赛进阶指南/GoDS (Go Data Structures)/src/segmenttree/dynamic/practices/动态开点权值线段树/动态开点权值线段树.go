// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/segment_tree.go

// 动态开点权值线段树
// Update
// Query/QueryAll
// Set/Get
// Merge
// Split

package main

import "fmt"

func main() {
	segmentTree1 := CreateSegmentTree(1, 1e9)
	fmt.Println(segmentTree1.QueryAll())
	segmentTree1.Update(1, 5, 1)
	fmt.Println(segmentTree1.QueryAll())

	segmentTree2 := CreateSegmentTree(1, 1e9)
	segmentTree1.Update(1, 5, 1)
	segmentTree1.Merge(segmentTree2)
	fmt.Println(segmentTree1.Query(1, 3))

	root1, root2 := segmentTree1.Split(1, 3) // !将 [1,3] 从 segmentTree1 分离成 root1 和 root2
	fmt.Println(root1.QueryAll(), root2.QueryAll())

	root1.Set(1, 1)
	fmt.Println(root1.QueryAll(), root2.QueryAll())
}

// 指定区间上下界建立权值线段树.
func CreateSegmentTree(lower, upper int) *lazyNode {
	root := &lazyNode{left: lower, right: upper}
	return root
}

type lazyNode struct {
	left, right           int
	sum                   int
	lazy                  int
	leftChild, rightChild *lazyNode
}

func (lazyNode) op(a, b int) int {
	return a + b
}

func (o *lazyNode) propagate(add int) {
	o.lazy += add                         // % mod
	o.sum += (o.right - o.left + 1) * add // % mod
}

func (o *lazyNode) pushDown() {
	m := (o.left + o.right) >> 1
	if o.leftChild == nil {
		o.leftChild = &lazyNode{left: o.left, right: m}
	}
	if o.rightChild == nil {
		o.rightChild = &lazyNode{left: m + 1, right: o.right}
	}
	if add := o.lazy; add != 0 {
		o.leftChild.propagate(add)
		o.rightChild.propagate(add)
		o.lazy = 0
	}
}

func (o *lazyNode) pushUp() {
	o.sum = o.op(o.leftChild.QueryAll(), o.rightChild.QueryAll())
}

// [left, right]
func (o *lazyNode) Update(left, right int, add int) {
	if left <= o.left && o.right <= right {
		o.propagate(add)
		return
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if left <= m {
		o.leftChild.Update(left, right, add)
	}
	if m < right {
		o.rightChild.Update(left, right, add)
	}
	o.pushUp()
}

// [left, right]
func (o *lazyNode) Query(left, right int) int {
	if o == nil || left > o.right || right < o.left {
		return 0
	}
	if left <= o.left && o.right <= right {
		return o.sum
	}
	o.pushDown()
	return o.op(o.leftChild.Query(left, right), o.rightChild.Query(left, right))
}

func (o *lazyNode) Set(pos int, val int) {
	if o.left == o.right {
		o.sum = val
		return
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if pos <= m {
		o.leftChild.Set(pos, val)
	} else {
		o.rightChild.Set(pos, val)
	}
	o.pushUp()
}

func (o *lazyNode) Get(pos int) int {
	if o.left == o.right {
		return o.sum
	}
	o.pushDown()
	m := (o.left + o.right) >> 1
	if pos <= m {
		return o.leftChild.Get(pos)
	}
	return o.rightChild.Get(pos)
}

func (o *lazyNode) QueryAll() int {
	if o != nil {
		return o.sum
	}
	return 0
}

// 线段树合并
func (o *lazyNode) Merge(b *lazyNode) *lazyNode {
	if o == nil {
		return b
	}
	if b == nil {
		return o
	}
	if o.left == o.right {
		o.sum += b.sum
		return o
	}
	o.leftChild = o.leftChild.Merge(b.leftChild)
	o.rightChild = o.rightChild.Merge(b.rightChild)
	o.pushUp()
	return o
}

// 线段树分裂
//  将区间 [l,r] 从原树分离到 other 上, this 为原树的剩余部分
func (o *lazyNode) Split(left, right int) (this, other *lazyNode) {
	this, other = o.split(nil, left, right)
	return
}

func (o *lazyNode) split(b *lazyNode, l, r int) (*lazyNode, *lazyNode) {
	if o == nil || l > o.right || r < o.left {
		return o, nil
	}
	if l <= o.left && o.right <= r {
		return nil, o
	}
	if b == nil {
		b = &lazyNode{left: o.left, right: o.right}
	}
	o.leftChild, b.leftChild = o.leftChild.split(b.leftChild, l, r)
	o.rightChild, b.rightChild = o.rightChild.split(b.rightChild, l, r)
	o.pushUp()
	b.pushUp()
	return o, b
}

// 权值线段树求第 k 小
// 调用前需保证 1 <= k <= root.QueryAll()
func (o *lazyNode) kth(k int) int {
	if o.left == o.right {
		return o.left
	}
	if lc := o.leftChild.QueryAll(); k <= lc {
		return o.leftChild.kth(k)
	} else {
		return o.rightChild.kth(k - lc)
	}
}
