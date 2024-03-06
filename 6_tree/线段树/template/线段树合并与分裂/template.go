//	struct {
//		int ls, rs, v;
//	} T[M];
//
// int cnt, nn, R[N], P[N];
// namespace segt {
//
//	int query(int x, int p, int l = 1, int r = nn) {
//		if (l == r) return T[p].v;
//		int mid = l + r >> 1;
//		if (x <= mid)
//				return query(x, T[p].ls, l, mid);
//		else
//				return query(x, T[p].rs, mid + 1, r);
//	}
//
//	int kth(int k, int p, int l = 1, int r = nn) {
//		if (l == r) return T[p].v >= k ? l : 0;
//		int mid = l + r >> 1;
//		if (T[T[p].ls].v >= k)
//				return kth(k, T[p].ls, l, mid);
//		else
//				return kth(k - T[T[p].ls].v, T[p].rs, mid + 1, r);
//	}
//
//	void pushup(int p) {
//		T[p].v = T[T[p].ls].v + T[T[p].rs].v;
//	}
//
//	void add(int x, int d, int &p, int l = 1, int r = nn) {
//		if (!p) p = ++cnt;
//		if (l == r) return (void)(T[p].v += d);
//		int mid = l + r >> 1;
//		if (x <= mid)
//				add(x, d, T[p].ls, l, mid);
//		else
//				add(x, d, T[p].rs, mid + 1, r);
//		pushup(p);
//	}
//
//	int merge(int p, int q, int l = 1, int r = nn) {
//		if (!p || !q) return p + q;
//		if (l == r) return T[p].v += T[q].v, p;
//		int mid = l + r >> 1;
//		T[p].ls = merge(T[p].ls, T[q].ls, l, mid);
//		T[p].rs = merge(T[p].rs, T[q].rs, mid + 1, r);
//		pushup(p);
//		return p;
//	}

package main

import "fmt"

func main() {
	sm := NewSegmentTreeMerger(0, 10)
	nodes := sm.Build(10, func(i int32) int32 { return i })
	fmt.Println(nodes)
	nodes[1] = sm.Merge(nodes[2], nodes[1])
	fmt.Println(sm.QueryAll(nodes[1]))
	fmt.Println(nodes)
}

type E = int32

func e() E        { return 0 }
func op(a, b E) E { return a + b }

type Node struct {
	value                 E
	leftChild, rightChild *Node
}

func (n *Node) String() string {
	return fmt.Sprintf("%d", n.value)
}

type SegmentTreeMerger struct {
	left, right int32
}

// 指定闭区间[left,right]建立Merger.
func NewSegmentTreeMerger(left, right int32) *SegmentTreeMerger {
	return &SegmentTreeMerger{left: left, right: right}
}

func (sm *SegmentTreeMerger) Build(n int32, f func(i int32) E) []*Node {
	res := make([]*Node, n)
	for i := int32(0); i < n; i++ {
		res[i] = sm.Alloc(i, f(i))
	}
	return res
}

func (sm *SegmentTreeMerger) Alloc(index int32, value E) *Node {
	return sm._alloc(index, value, sm.left, sm.right)
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

func (sm *SegmentTreeMerger) _alloc(index int32, value E, left, right int32) *Node {
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
	node.value = op(sm._eval(node.leftChild), sm._eval(node.rightChild))
	return node
}

func (sm *SegmentTreeMerger) _get(node *Node, index int32, left, right int32) E {
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
		sm._set(node.leftChild, index, value, left, mid)
	} else {
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
		sm._update(node.leftChild, index, value, left, mid)
	} else {
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
	newNode := &Node{}
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
		b = &Node{leftChild: a.leftChild, rightChild: a.rightChild}
	}
	mid := (left + right) >> 1
	a.leftChild, b.leftChild = sm._split(a.leftChild, b.leftChild, L, R, left, mid)
	a.rightChild, b.rightChild = sm._split(a.rightChild, b.rightChild, L, R, mid+1, right)
	a.value = sm._eval(a.leftChild) + sm._eval(a.rightChild)
	b.value = sm._eval(b.leftChild) + sm._eval(b.rightChild)
	return a, b
}

func (sm *SegmentTreeMerger) _eval(node *Node) E {
	if node == nil {
		return e()
	}
	return node.value
}
