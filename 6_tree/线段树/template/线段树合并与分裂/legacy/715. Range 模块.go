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
type Size = int32

// !线段树维护的数据类型 区间和
type E = int
type Id = int

func e1() E                { return 0 }
func e2(start, end Size) E { return int(end - start + 1) } // 区间[start,end)的初始值.
func id() Id               { return -1 }                   // !-1表示monoid 0表示需要染成0 1表示需要染成1
func op(a, b E) E {
	return a + b
}
func mapping(f Id, g E, size Size) E {
	if f == -1 {
		return g
	}
	return f * int(size)
}
func composition(f, g Id) Id {
	if f == -1 {
		return g
	}
	return f
}

// 指定区间上下界建立线段树
func CreateSegmentTree(lower, upper Size) *Node {
	root := newNode(lower, upper)
	return root
}

type Node struct {
	left, right           Size
	leftChild, rightChild *Node

	data E
	lazy Id
}

func (o *Node) Update(left, right Size, lazy Id) {
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
	res := e2(left, right)
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

func newNode(left, right Size) *Node {
	return &Node{left: left, right: right, lazy: id(), data: e2(left, right)}
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
