package main

type RangeModule struct {
	segmentTree *Node
}

func Constructor() RangeModule {
	return RangeModule{CreateSegmentTree(0, 1e9)}
}

func (this *RangeModule) AddRange(left int, right int) {

	this.segmentTree.Update(left, right-1, 1)
}

func (this *RangeModule) QueryRange(left int, right int) bool {
	return this.segmentTree.Query(left, right-1).sum == right-left
}

func (this *RangeModule) RemoveRange(left int, right int) {
	this.segmentTree.Update(left, right-1, 0)
}

/**
 * Your RangeModule object will be instantiated and called as such:
 * obj := Constructor();
 * obj.AddRange(left,right);
 * param_2 := obj.QueryRange(left,right);
 * obj.RemoveRange(left,right);
 */

// !线段树维护的数据类型 区间和
type Data = struct{ size, sum int }
type Lazy = int

func e(left, right int) Data { return Data{size: right - left + 1} }
func id() Lazy               { return -1 } // !-1表示monoid 0表示需要染成0 1表示需要染成1
func op(leftData, rightData Data) Data {
	return Data{leftData.size + rightData.size, leftData.sum + rightData.sum}
}
func mapping(parentLazy Lazy, childData Data) Data {
	if parentLazy == -1 {
		return childData
	}
	if parentLazy == 0 {
		return Data{childData.size, 0}
	}
	return Data{childData.size, childData.size}
}
func composition(parentLazy, childLazy Lazy) Lazy {
	if parentLazy == -1 {
		return childLazy
	}
	return parentLazy
}

// 指定区间上下界建立线段树
func CreateSegmentTree(lower, upper int) *Node {
	root := newNode(lower, upper)
	return root
}

type Node struct {
	left, right           int
	leftChild, rightChild *Node

	data Data
	lazy Lazy
}

func (o *Node) Update(left, right int, lazy Lazy) {
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

func (o *Node) Query(left, right int) Data {
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

func (o *Node) QueryAll() Data {
	return o.data
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
func (o *Node) propagate(lazy Lazy) {
	o.data = mapping(lazy, o.data)
	o.lazy = composition(lazy, o.lazy)
}
