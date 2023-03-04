// !动态开点线段树(支持线段树的合并与分裂)
// https://github.dev/EndlessCheng/codeforces-go/blob/551e365da1be6ff2875955a8ededc6479e336528/copypasta/segment_tree.go#L459
// https://www.luogu.com.cn/blog/styx-ferryman/xian-duan-shu-ge-bing-zong-ru-men-dao-fang-qi
// https://www.luogu.com.cn/problem/P4556
// https://www.luogu.com.cn/problem/P5494

// 动态开点线段树
// Update
// Query/QueryAll
// Set/Get
// Merge
// Split

package main

import (
	"fmt"
)

func main() {
	segmentTree1 := CreateSegmentTree(1, 1e9)
	fmt.Println(segmentTree1.QueryAll().sum)
	segmentTree1.Update(1, 5, 1)
	fmt.Println(segmentTree1.QueryAll().sum)

	segmentTree2 := CreateSegmentTree(1, 1e9)
	segmentTree1.Update(1, 5, 1)
	segmentTree1.Merge(segmentTree2)
	fmt.Println(segmentTree1.Query(1, 3).sum)

	root1, root2 := segmentTree1.Split(1, 3) // !将 [1,3] 从 segmentTree1 分离成 root1 和 root2
	fmt.Println(root1.QueryAll(), root2.QueryAll())
	root1.Set(1, E{1, 1})
	fmt.Println(root1.QueryAll(), root2.QueryAll())
}

// !线段树维护的数据类型 示例: 区间和
type E = struct{ size, sum int }
type Id = int

func e(left, right int) E             { return E{size: right - left + 1} }
func id() Id                          { return 0 }
func op(e1, e2 E) E                   { return E{e1.size + e2.size, e1.sum + e2.sum} }
func mapping(f Id, g E) E             { return E{g.size, g.sum + f*g.size} }
func composition(parent, child Id) Id { return parent + child }

//
//
//
// 指定区间上下界建立线段树
func CreateSegmentTree(lower, upper int) *Node {
	root := newNode(lower, upper)
	return root
}

type Node struct {
	left, right           int
	leftChild, rightChild *Node

	data E
	lazy Id
}

// lower<=left<=right<=upper
func (o *Node) Update(left, right int, lazy Id) {
	if left <= o.left && o.right <= right {
		o.propagate(lazy)
		return
	}

	o.pushDown()
	mid := (o.left + o.right) >> 1
	if left <= mid {
		o.leftChild.Update(left, right, lazy)
	}
	if right > mid {
		o.rightChild.Update(left, right, lazy)
	}
	o.pushUp()
}

// lower<=left<=right<=upper
func (o *Node) Query(left, right int) E {
	if left <= o.left && o.right <= right {
		return o.data
	}
	o.pushDown()
	mid := (o.left + o.right) >> 1
	res := e(left, right)
	if left <= mid {
		res = op(res, o.leftChild.Query(left, right))
	}
	if right > mid {
		res = op(res, o.rightChild.Query(left, right))
	}
	return res
}

func (o *Node) QueryAll() E {
	if o == nil {
		return e(0, -1)
	}
	return o.data
}

// lower<=pos<=upper
func (o *Node) Set(pos int, val E) {
	if o.left == o.right {
		o.data = val
		return
	}
	o.pushDown()
	mid := (o.left + o.right) >> 1
	if pos <= mid {
		o.leftChild.Set(pos, val)
	} else {
		o.rightChild.Set(pos, val)
	}
	o.pushUp()
}

// lower<=pos<=upper
func (o *Node) Get(pos int) E {
	if o.left == o.right {
		return o.data
	}
	o.pushDown()
	mid := (o.left + o.right) >> 1
	if pos <= mid {
		return o.leftChild.Get(pos)
	}
	return o.rightChild.Get(pos)
}

// Build from array. [1,len(nums))]
func (o *Node) Build(nums []E) {
	o.build(nums, 1, len(nums))
}

func (o *Node) build(nums []E, left, right int) {
	o.left, o.right = left, right
	if left == right {
		o.data = nums[left-1]
		return
	}
	m := (left + right) >> 1
	o.leftChild = newNode(left, m)
	o.leftChild.build(nums, left, m)
	o.rightChild = newNode(m+1, right)
	o.rightChild.build(nums, m+1, right)
	o.pushUp()
}

// 线段树合并
//  root = root.Merge(other)
func (this *Node) Merge(other *Node) *Node {
	if this == nil {
		return other
	}
	if other == nil {
		return this
	}
	if this.left == this.right {
		this.data = op(this.data, other.data)
		return this
	}
	this.leftChild = this.leftChild.Merge(other.leftChild)
	this.rightChild = this.rightChild.Merge(other.rightChild)
	this.pushUp()
	return this
}

// 线段树分裂
// 将区间 [left,right] 从原树分离到 other 上, this 为原树的剩余部分
//  root1, root2 := segmentTree1.Split(1, 3)
func (o *Node) Split(left, right int) (this, other *Node) {
	this, other = o.split(nil, left, right)
	return
}

func (this *Node) split(other *Node, left, right int) (*Node, *Node) {
	if this == nil || left > this.right || right < this.left {
		return this, nil
	}
	if left <= this.left && this.right <= right {
		return nil, this
	}
	if other == nil {
		other = newNode(this.left, this.right)
	}
	this.leftChild, other.leftChild = this.leftChild.split(other.leftChild, left, right)
	this.rightChild, other.rightChild = this.rightChild.split(other.rightChild, left, right)
	this.pushUp()
	other.pushUp()
	return this, other
}

// 权值线段树求第 k 小
// 调用前需保证 1 <= k <= root.QueryAll()
func (o *Node) kth(k int) int {
	if o.left == o.right {
		return o.left
	}
	if lc := o.leftChild.QueryAll().sum; k <= lc {
		return o.leftChild.kth(k)
	} else {
		return o.rightChild.kth(k - lc)
	}
}

func newNode(left, right int) *Node {
	return &Node{left: left, right: right, lazy: id(), data: e(left, right)}
}

// op
func (o *Node) pushUp() {
	o.data = op(o.leftChild.QueryAll(), o.rightChild.QueryAll())
}

func (o *Node) pushDown() {
	mid := (o.left + o.right) >> 1
	if o.leftChild == nil {
		o.leftChild = newNode(o.left, mid)
	}
	if o.rightChild == nil {
		o.rightChild = newNode(mid+1, o.right)
	}
	if o.lazy != id() {
		o.leftChild.propagate(o.lazy)
		o.rightChild.propagate(o.lazy)
		o.lazy = id()
	}
}

// mapping + composition
func (o *Node) propagate(lazy Id) {
	o.data = mapping(lazy, o.data)
	o.lazy = composition(lazy, o.lazy)
}
