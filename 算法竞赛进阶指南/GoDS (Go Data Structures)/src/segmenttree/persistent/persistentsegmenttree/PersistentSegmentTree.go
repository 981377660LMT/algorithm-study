// 可持久化线段树（又称函数式线段树、主席树） Persistent Segment Tree
// !支持区间更新单点查询/查询区间第k大/区间众数位置/区间元素种类数

package main

import "fmt"

func main() {
	root0 := Build(1, 5, []int{1, 2, 3, 4, 5})
	fmt.Println(root0.QueryRangeDangerously(1, 5)) // 15
	root1 := root0.Update(1, 3)
	fmt.Println(root1.QueryRangeDangerously(1, 5)) // 18
	fmt.Println(root0, root1)                      // PersistentArray [1 2 3 4 5] PersistentArray [4 2 3 4 5]

	// !区间更新之后只能单点查询了
	roo2 := root1.UpdateRangeDangerously(1, 5, 3)
	fmt.Println(roo2.Query(3)) // 6
}

// !线段树维护的数据类型
type Data = int
type Lazy = int

func e() Data                                      { return 0 }
func id() Lazy                                     { return 0 }
func op(leftData, rightData Data) Data             { return leftData + rightData }
func mapping(parentLazy Lazy, childData Data) Data { return childData + parentLazy }
func composition(parentLazy, childLazy Lazy) Lazy  { return parentLazy + childLazy }

type Node struct {
	left, right           int
	size                  int
	leftChild, rightChild *Node

	data Data
	lazy Lazy
}

func newNode(left, right int, data Data, size int) *Node {
	return &Node{data: data, lazy: id(), size: size, left: left, right: right}
}

// op
func (o *Node) pushUp() {
	o.size = o.leftChild.size + o.rightChild.size
	o.data = op(o.leftChild.data, o.rightChild.data)
}

func (o *Node) pushDown() {
	if o.lazy != id() {
		if o.leftChild != nil {
			o.leftChild.propagate(o.lazy)
		}
		if o.rightChild != nil {
			o.rightChild.propagate(o.lazy)
		}
		o.lazy = id()
	}
}

// mapping + composition
func (o *Node) propagate(lazy Lazy) {
	o.data = mapping(lazy, o.data)
	o.lazy = composition(lazy, o.lazy)
}

// usage:
//  roots := make([]*Node, 1, maxVersion+1)
//  roots[0] = Build(1, len(nums), nums)
func Build(left, right int, nums []Data) *Node {
	if left == right {
		return newNode(left, right, nums[left-1], 1)
	}

	node := newNode(left, right, nums[left-1], 1)
	mid := (left + right) >> 1
	node.leftChild = Build(left, mid, nums)
	node.rightChild = Build(mid+1, right, nums)
	node.pushUp()
	return node
}

// usage:
//  newRoot := roots[version].Update(index, value)
//  roots = append(roots, newRoot)
//  1 <= index <= len(nums)
// !注意：为了拷贝一份 Node，这里的接收器不是指针
func (o Node) Update(index int, lazy Lazy) *Node {
	if o.left == o.right {
		o.propagate(lazy)
		return &o
	}

	o.pushDown()
	mid := (o.left + o.right) >> 1
	if index <= mid {
		o.leftChild = o.leftChild.Update(index, lazy)
	} else {
		o.rightChild = o.rightChild.Update(index, lazy)
	}
	o.pushUp()
	return &o
}

func (o *Node) Query(index int) Data {
	if o.left == o.right {
		return o.data
	}

	o.pushDown()
	mid := (o.left + o.right) >> 1
	if index <= mid {
		return o.leftChild.Query(index)
	}
	return o.rightChild.Query(index)
}

// !如果进行了区间更新，那么就不能使用区间查询了
//   1 <= left <= right <= len(nums)
func (o *Node) QueryRangeDangerously(left, right int) int {
	if left <= o.left && o.right <= right {
		return o.data
	}

	mid := (o.left + o.right) >> 1
	if right <= mid {
		return o.leftChild.QueryRangeDangerously(left, right)
	}

	if left > mid {
		return o.rightChild.QueryRangeDangerously(left, right)
	}

	return op(o.leftChild.QueryRangeDangerously(left, right), o.rightChild.QueryRangeDangerously(left, right))
}

// !区间更新（只能配合单点查询）
//   每次更新需要拷贝一份 root 节点
//   1 <= left <= right <= len(nums)
//   https://atcoder.jp/contests/abc253/tasks/abc253_f
func (o *Node) UpdateRangeDangerously(left, right int, lazy Lazy) *Node {
	if left <= o.left && o.right <= right {
		o.propagate(lazy)
		return o
	}

	leftCopy := *o.leftChild
	o.leftChild = &leftCopy
	rightCopy := *o.rightChild
	o.rightChild = &rightCopy
	o.pushDown()
	mid := (o.left + o.right) >> 1
	if left <= mid {
		o.leftChild.UpdateRangeDangerously(left, right, lazy)
	}
	if right > mid {
		o.rightChild.UpdateRangeDangerously(left, right, lazy)
	}

	// !因为使用了区间修改后只支持单点查询，所以不需要 pushUp
	return o
}

func (o *Node) String() string {
	res := o.dfs()
	return fmt.Sprintf("PersistentArray %v", res)
}

func (o *Node) dfs() []Data {
	res := make([]Data, 0, o.size)
	if o.left == o.right {
		res = append(res, o.data)
		return res
	}
	res = append(res, o.leftChild.dfs()...)
	res = append(res, o.rightChild.dfs()...)
	return res
}

// !一些api
// !1. 区间第k小
// 主席树相当于对`离散化后`数组的每个前缀建立一颗线段树
// !离散化时，求 kth 需要`将相同元素也视作不同的`

// 查询区间 [left,right] 中第 k 小在整个数组上的名次（从 1 开始）.
// !注意返回的是（排序去重后的数组的）`下标`，不是元素值.
// usage:
//     roots := make([]*Node, 1, maxVersion+1)
//     roots[0] = Build(1, len(A))
//     !roots[i+1] = roots[i].update(kth[i], 1)  // kth[i] 为 A[i] 离散化后的值（从 1 开始） , 1 表示增加1个计数
//     roots[right].kth(roots[left-1], k)  // 类似前缀和 [left,right] 1<=left<=right<=n
func (o *Node) kth(oldNode *Node, k int) int {
	if o.left == o.right {
		return o.left
	}

	leftCount := o.leftChild.size - oldNode.leftChild.size // 值域左半边的数字个数
	if k <= leftCount {
		return o.leftChild.kth(oldNode.leftChild, k)
	}
	return o.rightChild.kth(oldNode.rightChild, k-leftCount)
}

// !2. 区间[left,right] 中 在 [lower,upper] 范围内元素的个数
// !lower和higher是离散化后的值(从1开始)
// usage:
//     roots := make([]*Node, 1, maxVersion+1)
//     roots[0] = Build(1, len(A))
//     roots[i+1] = roots[i].update(kth[i], 1)  // kth[i] 为 A[i] 离散化后的值（从 1 开始）
//     roots[right].count(roots[left-1], lower, upper)  // 类似前缀和 [left,right] 1<=left<=right<=n
// https://github.dev/EndlessCheng/codeforces-go/blob/551e365da1be6ff2875955a8ededc6479e336528/copypasta/segment_tree.go#L753
func (o *Node) countRangeTypes(oldNode *Node, lower, upper int) int {
	if upper < o.left || lower > o.right {
		return 0
	}

	if lower <= o.left && o.right <= upper {
		return o.data - oldNode.data // o.sum - oldNode.sum  // !区间和
	}

	mid := (o.left + o.right) >> 1
	if upper <= mid {
		return o.leftChild.countRangeTypes(oldNode.leftChild, lower, upper)
	}

	if lower > mid {
		return o.rightChild.countRangeTypes(oldNode.rightChild, lower, upper)
	}

	return o.leftChild.countRangeTypes(oldNode.leftChild, lower, upper) + o.rightChild.countRangeTypes(oldNode.rightChild, lower, upper)
}

// !区间绝对众数以及出现次数 (注意返回的众数是离散化后的值）
//  threshold: 众数出现次数的阈值(>=threshold)
func (o *Node) findMajority(old *Node, threshold int) (majority, count int) {
	if o.left == o.right {
		return o.left, o.data - old.data // o.sum - old.sum
	}

	if o.leftChild.data-old.leftChild.data >= threshold {
		return o.leftChild.findMajority(old.leftChild, threshold)
	}
	if o.rightChild.data-old.rightChild.data >= threshold {
		return o.rightChild.findMajority(old.rightChild, threshold)
	}

	return -1, 0
}
