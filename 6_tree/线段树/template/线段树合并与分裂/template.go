package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

// 允许的空间很大时，禁用gc加速
func init() {
	debug.SetGCPercent(-1)
}

func main() {
	CF600E()
}

// Lomsat gelral
// https://www.luogu.com.cn/problem/CF600E
// 给你一棵有n个点的树(n≤1e5)，树上每个节点都有一种颜色ci(ci≤n)，
// 求每个点子树出现最多的颜色的编号的和.
//
// 每个结点开一个线段树维护(出现最多次颜色的个数，最多次颜色的和).
func CF600E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	colors := make([]int32, n)
	for i := range colors {
		fmt.Fscan(in, &colors[i])
	}

	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	res := make([]int, n)
	roots := make([]*Node, n)
	seg := NewSegmentTreeMerger(0, n)
	for i := int32(0); i < n; i++ {
		roots[i] = seg.Alloc()
		seg.Set(roots[i], colors[i], E{score: int(colors[i]), maxCount: 1})
	}

	var dfs func(cur, pre int32)
	dfs = func(cur, pre int32) {
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			dfs(next, cur)
			roots[cur] = seg.MergeDestructively(roots[cur], roots[next])
		}
		res[cur] = seg.QueryAll(roots[cur]).score
	}
	dfs(0, -1)

	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

type E = struct {
	score    int
	maxCount int32
}

func e() E { return E{} }
func op(a, b E) E {
	if a.maxCount > b.maxCount {
		return a
	}
	if a.maxCount < b.maxCount {
		return b
	}
	a.score += b.score
	return a
}
func merge(a, b E) E { // 合并两个不同的树的结点的函数
	a.maxCount += b.maxCount
	return a
}

type Node struct {
	value                 E
	leftChild, rightChild *Node
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.value)
}

type SegmentTreeMerger struct {
	left, right int32
}

// 指定闭区间[left,right]建立Merger.
func NewSegmentTreeMerger(left, right int32) *SegmentTreeMerger {
	return &SegmentTreeMerger{left: left, right: right}
}

// NewRoot().
func (sm *SegmentTreeMerger) Alloc() *Node {
	return &Node{value: e()}
}

// 权值线段树求第 k 小.
// 调用前需保证 1 <= k <= node.value.
func (sm *SegmentTreeMerger) Kth(node *Node, k int32, getCount func(node *Node) int32) (res int32, ok bool) {
	if k < 1 || k > getCount(node) {
		return
	}
	return sm._kth(k, node, sm.left, sm.right, getCount), true
}

func (sm *SegmentTreeMerger) Get(node *Node, index int32) E {
	return sm._get(node, index, sm.left, sm.right)
}

func (sm *SegmentTreeMerger) Set(node *Node, index int32, value E) {
	sm._set(node, index, value, sm.left, sm.right)
}

func (sm *SegmentTreeMerger) Query(node *Node, left, right int32) E {
	return sm._query(node, left, right, sm.left, sm.right)
}

func (sm *SegmentTreeMerger) QueryAll(node *Node) E {
	return sm._eval(node)
}

func (sm *SegmentTreeMerger) Update(node *Node, index int32, value E) {
	sm._update(node, index, value, sm.left, sm.right)
}

// 用一个新的节点存合并的结果，会生成重合节点数量的新节点.
func (sm *SegmentTreeMerger) Merge(a, b *Node) *Node {
	return sm._merge(a, b, sm.left, sm.right)
}

// 把第二棵树直接合并到第一棵树上，比较省空间，缺点是会丢失合并前树的信息.
func (sm *SegmentTreeMerger) MergeDestructively(a, b *Node) *Node {
	return sm._mergeDestructively(a, b, sm.left, sm.right)
}

// 线段树分裂，将区间 [left,right] 从原树分离到 other 上, this 为原树的剩余部分.
func (sm *SegmentTreeMerger) Split(node *Node, left, right int32) (this, other *Node) {
	this, other = sm._split(node, nil, left, right, sm.left, sm.right)
	return
}

func (sm *SegmentTreeMerger) _kth(k int32, node *Node, left, right int32, getCount func(*Node) int32) int32 {
	if left == right {
		return left
	}
	mid := (left + right) >> 1
	leftCount := int32(0)
	if node.leftChild != nil {
		leftCount = getCount(node.leftChild)
	}
	if leftCount >= k {
		return sm._kth(k, node.leftChild, left, mid, getCount)
	} else {
		return sm._kth(k-leftCount, node.rightChild, mid+1, right, getCount)
	}
}

func (sm *SegmentTreeMerger) _get(node *Node, index int32, left, right int32) E {
	if node == nil {
		return e()
	}
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
func (sm *SegmentTreeMerger) _query(node *Node, L, R int32, left, right int32) E {
	if node == nil {
		return e()
	}
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
	return op(sm._query(node.leftChild, L, R, left, mid), sm._query(node.rightChild, L, R, mid+1, right))
}

func (sm *SegmentTreeMerger) _set(node *Node, index int32, value E, left, right int32) {
	if left == right {
		node.value = value
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._set(node.leftChild, index, value, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._set(node.rightChild, index, value, mid+1, right)
	}
	node.value = op(sm._eval(node.leftChild), sm._eval(node.rightChild))
}

func (sm *SegmentTreeMerger) _update(node *Node, index int32, value E, left, right int32) {
	if left == right {
		node.value = op(node.value, value)
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._update(node.leftChild, index, value, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._update(node.rightChild, index, value, mid+1, right)
	}
	node.value = op(sm._eval(node.leftChild), sm._eval(node.rightChild))
}

func (sm *SegmentTreeMerger) _merge(a, b *Node, left, right int32) *Node {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	newNode := sm.Alloc()
	if left == right {
		newNode.value = merge(a.value, b.value)
		return newNode
	}
	mid := (left + right) >> 1
	newNode.leftChild = sm._merge(a.leftChild, b.leftChild, left, mid)
	newNode.rightChild = sm._merge(a.rightChild, b.rightChild, mid+1, right)
	newNode.value = op(sm._eval(newNode.leftChild), sm._eval(newNode.rightChild))
	return newNode
}

func (sm *SegmentTreeMerger) _mergeDestructively(a, b *Node, left, right int32) *Node {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	if left == right {
		a.value = merge(a.value, b.value)
		return a
	}
	mid := (left + right) >> 1
	a.leftChild = sm._mergeDestructively(a.leftChild, b.leftChild, left, mid)
	a.rightChild = sm._mergeDestructively(a.rightChild, b.rightChild, mid+1, right)
	a.value = op(sm._eval(a.leftChild), sm._eval(a.rightChild))
	return a
}

func (sm *SegmentTreeMerger) _split(a, b *Node, L, R int32, left, right int32) (*Node, *Node) {
	if a == nil || L > right || R < left {
		return a, nil
	}
	if L <= left && right <= R {
		return nil, a
	}
	if b == nil {
		b = sm.Alloc()
	}
	mid := (left + right) >> 1
	a.leftChild, b.leftChild = sm._split(a.leftChild, b.leftChild, L, R, left, mid)
	a.rightChild, b.rightChild = sm._split(a.rightChild, b.rightChild, L, R, mid+1, right)
	a.value = op(sm._eval(a.leftChild), sm._eval(a.rightChild))
	b.value = op(sm._eval(b.leftChild), sm._eval(b.rightChild))
	return a, b
}

func (sm *SegmentTreeMerger) _eval(node *Node) E {
	if node == nil {
		return e()
	}
	return node.value
}
