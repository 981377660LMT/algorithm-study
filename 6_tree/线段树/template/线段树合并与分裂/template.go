package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// demo()

	P3224()
}

// P3224 [HNOI2012] 永无乡
// https://www.luogu.com.cn/problem/P3224
// 永无乡包含 n座岛，编号从 1到 n，每座岛都有自己的独一无二的重要度，
// 按照重要度可以将这 n座岛排名，名次用1到n来表示.
// 现在有两种操作：
// B x y 在x和y之间建立一座桥，使得x和y之间可以互相到达
// Q x y 询问当前与岛 x连通的所有岛中第 k重要的是哪座岛,如果该岛屿不存在，则输出 −1.
//
// 并查集维护连通性，用权值线段树维护集合的第k小数。
// 当我们合并两个集合的时候，直接合并它们的线段树即可。
func P3224() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)

	uf := NewUnionFindArraySimple(n)
	seg := NewSegmentTreeMerger(0, n-1)
	roots := make([]*Node, n) // 每个岛的线段树
	rankToId := make([]int32, n)
	for i := int32(0); i < n; i++ {
		var rank int32
		fmt.Fscan(in, &rank)
		rank--
		rankToId[rank] = i
		roots[i] = seg.Alloc()
		seg.Update(roots[i], rank, 1)
	}

	// 在x和y之间建立一座桥，使得x和y之间可以互相到达
	addEdge := func(u, v int32) {
		uf.Union(u, v, func(big, small int32) {
			roots[big] = seg.MergeDestructively(roots[big], roots[small])
		})
	}

	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		if u < 1 || v < 1 || u > n || v > n {
			continue
		}
		addEdge(u-1, v-1)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "B" {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			u, v = u-1, v-1
			addEdge(u, v)
		} else {
			var u, k int32
			fmt.Fscan(in, &u, &k)
			u--
			leader := uf.Find(u)
			rank, ok := seg.Kth(roots[leader], k, func(node *Node) int32 {
				return node.value
			})
			if !ok {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, rankToId[rank]+1) // !输出岛的编号
			}
		}
	}
}

type E = int32

func e() E        { return 0 }
func op(a, b E) E { return a + b }

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
		newNode.value = op(a.value, b.value)
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
		a.value = op(a.value, b.value)
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
