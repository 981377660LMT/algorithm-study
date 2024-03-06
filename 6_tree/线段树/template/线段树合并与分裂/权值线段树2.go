// https://zhuanlan.zhihu.com/p/575513452
// https://www.luogu.com/article/bsm4zrgr

package main

import (
	"fmt"
)

func main() {

}

// P4556 [Vani有约会] 雨天的尾巴 /【模板】线段树合并
// https://www.luogu.com.cn/problem/P4556
// 村落里的一共有n座房屋，并形成一个树状结构。
// 然后救济粮分m次发放，每次选择两个房屋x和y ，然后对于x到y的路径上(含x和y)每座房子里发放一袋z类型的救济粮。
// 深绘里想知道，当所有的救济粮发放完毕后，每座房子里存放的最多的是哪种救济粮。
// 如果有多种救济粮都是存放最多的，输出种类编号最小的一种。
// 如果某座房屋没有救济粮，则输出 0。
//
// 每个节点开一棵权值线段树，树上差分
func P4556() {}

// Lomsat gelral
// https://www.luogu.com.cn/problem/CF600E
// 给你一棵有n个点的树(n≤1e5)，树上每个节点都有一种颜色ci(ci≤n)，
// 求每个点子树出现最多的颜色/编号的和.
func CF600E() {}

// Dominant Indices
// https://www.luogu.com.cn/problem/CF1009F
func CF1009F() {}

// TODO 改成mutable操作，而不是返回新数据
type E = int32

type Node struct {
	value                 E
	leftChild, rightChild *Node
}

func (n *Node) String() string {
	return fmt.Sprintf("%d", n.value)
}

type SegmentTreeOnRange struct {
	min, max int32
}

// 指定闭区间[min,max]建立权值线段树.
func NewSegmentTreeOnRange(min, max int32) *SegmentTreeOnRange {
	return &SegmentTreeOnRange{min: min, max: max}
}

func (sm *SegmentTreeOnRange) Build(n int32, f func(i int32) E) []*Node {
	res := make([]*Node, n)
	for i := int32(0); i < n; i++ {
		res[i] = sm.Alloc(i, f(i))
	}
	return res
}

func (sm *SegmentTreeOnRange) Alloc(index int32, value E) *Node {
	return sm._alloc(index, value, sm.min, sm.max)
}

// 权值线段树求第 k 小.
// 调用前需保证 1 <= k <= node.value.
func (sm *SegmentTreeOnRange) Kth(node *Node, k int32) (value int32, ok bool) {
	if k < 1 || k > sm._eval(node) {
		return 0, false
	}
	return sm._kth(k, node, sm.min, sm.max), true
}

func (sm *SegmentTreeOnRange) Get(node *Node, index int32) E {
	return sm._get(node, index, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) Set(node *Node, index int32, value E) {
	sm._set(node, index, value, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) Query(node *Node, left, right int32) E {
	return sm._query(node, left, right, sm.min, sm.max)
}

func (sm *SegmentTreeOnRange) QueryAll(node *Node) E {
	return sm._eval(node)
}

func (sm *SegmentTreeOnRange) Update(node *Node, index int32, value E) {
	sm._update(node, index, value, sm.min, sm.max)
}

// 用一个新的节点存合并的结果，会生成重合节点数量的新节点.
func (sm *SegmentTreeOnRange) Merge(a, b *Node) *Node {
	return sm._merge(a, b, sm.min, sm.max)
}

// 把第二棵树直接合并到第一棵树上，比较省空间，缺点是会丢失合并前树的信息.
func (sm *SegmentTreeOnRange) MergeDestructively(a, b *Node) *Node {
	return sm._mergeDestructively(a, b, sm.min, sm.max)
}

// 线段树分裂，将区间 [left,right] 从原树分离到 other 上, this 为原树的剩余部分.
func (sm *SegmentTreeOnRange) Split(node *Node, left, right int32) (this, other *Node) {
	this, other = sm._split(node, nil, left, right, sm.min, sm.max)
	return
}

func (sm *SegmentTreeOnRange) _alloc(index int32, value E, left, right int32) *Node {
	if left == right {
		return &Node{value: value}
	}
	mid := (left + right) >> 1
	node := &Node{}
	if index <= mid {
		node.leftChild = sm._alloc(index, value, left, mid)
	} else {
		node.rightChild = sm._alloc(index, value, mid+1, right)
	}
	node.value = sm._eval(node.leftChild) + sm._eval(node.rightChild)
	return node
}

func (sm *SegmentTreeOnRange) _kth(k int32, node *Node, left, right int32) int32 {
	if left == right {
		return left
	}
	mid := (left + right) >> 1
	if leftCount := sm._eval(node.leftChild); leftCount >= k {
		return sm._kth(k, node.leftChild, left, mid)
	} else {
		return sm._kth(k-leftCount, node.rightChild, mid+1, right)
	}
}

func (sm *SegmentTreeOnRange) _get(node *Node, index int32, left, right int32) E {
	if left == right {
		return node.value
	}
	mid := (left + right) >> 1
	if index <= mid {
		return sm._get(node.leftChild, index, left, mid)
	} else {
		return sm._get(node.rightChild, index, mid+1, right)
	}
}
func (sm *SegmentTreeOnRange) _query(node *Node, L, R int32, left, right int32) E {
	if L <= left && right <= R {
		return node.value
	}
	mid := (left + right) >> 1
	if R <= mid {
		return sm._query(node.leftChild, L, R, left, mid)
	}
	if L > mid {
		return sm._query(node.rightChild, L, R, mid+1, right)
	}
	return sm._query(node.leftChild, L, R, left, mid) + sm._query(node.rightChild, L, R, mid+1, right)
}

func (sm *SegmentTreeOnRange) _set(node *Node, index int32, value E, left, right int32) {
	if left == right {
		node.value = value
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		sm._set(node.leftChild, index, value, left, mid)
	} else {
		sm._set(node.rightChild, index, value, mid+1, right)
	}
	node.value = sm._eval(node.leftChild) + sm._eval(node.rightChild)
}

func (sm *SegmentTreeOnRange) _update(node *Node, index int32, value E, left, right int32) {
	if left == right {
		node.value += value
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		sm._update(node.leftChild, index, value, left, mid)
	} else {
		sm._update(node.rightChild, index, value, mid+1, right)
	}
	node.value = sm._eval(node.leftChild) + sm._eval(node.rightChild)
}

func (sm *SegmentTreeOnRange) _merge(a, b *Node, left, right int32) *Node {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	newNode := &Node{}
	if left == right {
		newNode.value = a.value + b.value
		return newNode
	}
	mid := (left + right) >> 1
	newNode.leftChild = sm._merge(a.leftChild, b.leftChild, left, mid)
	newNode.rightChild = sm._merge(a.rightChild, b.rightChild, mid+1, right)
	newNode.value = sm._eval(newNode.leftChild) + sm._eval(newNode.rightChild)
	return newNode
}

func (sm *SegmentTreeOnRange) _mergeDestructively(a, b *Node, left, right int32) *Node {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	if left == right {
		a.value += b.value
		return a
	}
	mid := (left + right) >> 1
	a.leftChild = sm._mergeDestructively(a.leftChild, b.leftChild, left, mid)
	a.rightChild = sm._mergeDestructively(a.rightChild, b.rightChild, mid+1, right)
	a.value = sm._eval(a.leftChild) + sm._eval(a.rightChild)
	return a
}

func (sm *SegmentTreeOnRange) _split(a, b *Node, L, R int32, left, right int32) (*Node, *Node) {
	if a == nil || L > right || R < left {
		return a, nil
	}
	if L <= left && right <= R {
		return nil, a
	}
	if b == nil {
		b = &Node{}
	}
	mid := (left + right) >> 1
	a.leftChild, b.leftChild = sm._split(a.leftChild, b.leftChild, L, R, left, mid)
	a.rightChild, b.rightChild = sm._split(a.rightChild, b.rightChild, L, R, mid+1, right)
	a.value = sm._eval(a.leftChild) + sm._eval(a.rightChild)
	b.value = sm._eval(b.leftChild) + sm._eval(b.rightChild)
	return a, b
}
func (sm *SegmentTreeOnRange) _eval(node *Node) E {
	if node == nil {
		return 0
	}
	return node.value
}

type UnionFindArraySimple struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple(n int32) *UnionFindArraySimple {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple) Union(key1, key2 int32, cb func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	if cb != nil {
		cb(root1, root2)
	}
	return true
}

func (u *UnionFindArraySimple) Find(key int32) int32 {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = u.Find(u.data[key])
	return u.data[key]
}

func (u *UnionFindArraySimple) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}
