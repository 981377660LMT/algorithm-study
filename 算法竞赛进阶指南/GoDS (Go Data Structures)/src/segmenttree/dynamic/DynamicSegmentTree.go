// !动态开点线段树(支持线段树的合并与分裂)
// https://github.dev/EndlessCheng/codeforces-go/blob/551e365da1be6ff2875955a8ededc6479e336528/copypasta/segment_tree.go#L459

package dynamicsegmenttree

import (
	"fmt"
	"runtime/debug"
)

// atcoder等使用单组样例测试的oj上,禁用gc会快很多
func init() {
	debug.SetGCPercent(-1)
}

func main() {
	segmentTree1 := CreateSegmentTree(1, 1e9)
	fmt.Println(segmentTree1.QueryAll().sum)
	segmentTree1.Update(1, 5, 1)
	fmt.Println(segmentTree1.QueryAll().sum)

	segmentTree2 := CreateSegmentTree(1, 1e9)
	segmentTree1.Update(1, 5, 1)
	segmentTree1.Merge(segmentTree2)
	fmt.Println(segmentTree1.Query(1, 3).sum)

	root1, root2 := segmentTree1.Split(segmentTree2, 1, 3) // !将 [1,3] 从 segmentTree1 分离成 root1 和 root2
	fmt.Println(root1.QueryAll(), root2.QueryAll())
}

// !线段树维护的数据类型 示例: 区间和
type E = struct{ size, sum int }
type Id = int

func e(left, right int) E             { return E{size: right - left + 1} }
func id() Id                          { return 0 }
func op(left, right E) E              { return E{left.size + right.size, left.sum + right.sum} }
func mapping(parent Id, child E) E    { return E{child.size, child.sum + parent*child.size} }
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
	return o.data
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
// 将区间 [left,right] 从 this 分离到 other 中
//  root1, root2 := segmentTree1.Split(segmentTree2, 1, 3) // 将 segmentTree1 的 [1,3] 区间 分离成 root1 和 root2
func (this *Node) Split(other *Node, left, right int) (*Node, *Node) {
	if this == nil || left > this.right || right < this.left {
		return this, nil
	}
	if left <= this.left && this.right <= right {
		return nil, this
	}
	if other == nil {
		other = newNode(this.left, this.right)
	}

	this.leftChild, other.leftChild = this.leftChild.Split(other.leftChild, left, right)
	this.rightChild, other.rightChild = this.rightChild.Split(other.rightChild, left, right)
	this.checkNilAndPushUp()
	other.checkNilAndPushUp()
	return this, other
}

func (o *Node) checkNilAndPushUp() {
	var leftData, rightData E
	if o.leftChild != nil {
		leftData = o.leftChild.data
	} else {
		leftData = e(o.left, o.left)
	}
	if o.rightChild != nil {
		rightData = o.rightChild.data
	} else {
		rightData = e(o.right, o.right)
	}
	o.data = op(leftData, rightData)
}

func newNode(left, right int) *Node {
	return &Node{left: left, right: right, lazy: id(), data: e(left, right)}
}

// op
func (o *Node) pushUp() {
	o.data = op(o.leftChild.data, o.rightChild.data)
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
