// api:
//   NewLazyAVLTree(n int32, f func(i int32) E) *LazyAVLTree
//   Merge(other *LazyAVLTree)
//   Split(k int32) (*LazyAVLTree, *LazyAVLTree)
//   Insert(k int32, key E)
//   Pop(k int32) E
//   Get(k int32) E
//   Set(k int32, key E)
//   Update(start, end int32, f Id)
//   UpdateAll(f Id)
//   Query(start, end int32) E
//   QueryAll() E
//   Reverse(start, end int32)
//   ReverseAll()
//   Clear()
//   ToList() []E
//   Size() int32

package main

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"
)

func main() {
	// test()
	testTime()
}

type E = int32

func e() E        { return 0 }
func op(a, b E) E { return a + b }

type aNode struct {
	key         E
	data        E
	left, right *aNode
	height      int8
	size        int32
}

func _newANode(key E) *aNode {
	return &aNode{key: key, data: key, height: 1, size: 1}
}

func (n *aNode) Balance() int8 {
	if n.left == nil {
		if n.right == nil {
			return 0
		}
		return -n.right.height
	}
	if n.right == nil {
		return n.left.height
	}
	return n.left.height - n.right.height
}

func (n *aNode) String() string {
	if n.left == nil && n.right == nil {
		return fmt.Sprintf("key=%v, height=%v, size=%v\n", n.key, n.height, n.size)
	}
	return fmt.Sprintf("key=%v, height=%v, size=%v,\n left:%v,\n right:%v\n", n.key, n.height, n.size, n.left, n.right)
}

var (
	tmpPath = make([]*aNode, 0, 128)
)

type AVLTreeWithSum struct {
	root *aNode
}

func NewAVLTreeWithSum(n int32, f func(int32) E) *AVLTreeWithSum {
	res := &AVLTreeWithSum{}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (t *AVLTreeWithSum) Merge(other *AVLTreeWithSum) {
	t.root = t._mergeNode(t.root, other.root)
}

func (t *AVLTreeWithSum) Insert(k int32, key E) {
	n := t.Size()
	if k < 0 {
		k += n
	}
	if k < 0 {
		k = 0
	}
	if k > n {
		k = n
	}
	a, b := t._splitNode(t.root, k)
	t.root = t._mergeWithRoot(a, _newANode(key), b)
}

func (t *AVLTreeWithSum) Split(k int32) (*AVLTreeWithSum, *AVLTreeWithSum) {
	a, b := t._splitNode(t.root, k)
	return _newWithRoot(a), _newWithRoot(b)
}

func (t *AVLTreeWithSum) Pop(k int32) E {
	if k < 0 {
		k += t.Size()
	}
	a, b := t._splitNode(t.root, k+1)
	a, tmp := t._popRight(a)
	t.root = t._mergeNode(a, b)
	return tmp.key
}

func (t *AVLTreeWithSum) Set(k int32, key E) {
	if k < 0 {
		k += t.Size()
	}
	node := t.root
	tmpPath = tmpPath[:0]
	path := tmpPath
	for {
		path = append(path, node)
		t := int32(0)
		if node.left != nil {
			t = node.left.size
		}
		if t == k {
			node.key = key
			break
		}
		if t < k {
			k -= t + 1
			node = node.right
		} else {
			node = node.left
		}
	}
	for i := len(path) - 1; i >= 0; i-- {
		t._update(path[i])
	}
}

func (t *AVLTreeWithSum) Clear() { t.root = nil }

func (t *AVLTreeWithSum) ToList() []E {
	node := t.root
	stack := make([]*aNode, 0)
	res := make([]E, 0, t.Size())
	for len(stack) > 0 || node != nil {
		if node != nil {
			stack = append(stack, node)
			node = node.left
		} else {
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res = append(res, node.key)
			node = node.right
		}
	}
	return res
}

func (t *AVLTreeWithSum) Get(k int32) E {
	if k < 0 {
		k += t.Size()
	}
	node := t.root
	for {
		t := int32(0)
		if node.left != nil {
			t = node.left.size
		}
		if t == k {
			return node.key
		} else if t < k {
			k -= t + 1
			node = node.right
		} else {
			node = node.left
		}
	}
}

func (avl *AVLTreeWithSum) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end || avl.root == nil {
		return e()
	}
	var dfs func(node *aNode, left, right int32) E
	dfs = func(node *aNode, left, right int32) E {
		if right <= start || end <= left {
			return e()
		}
		if start <= left && right < end {
			return node.data
		}
		lsize := int32(0)
		if node.left != nil {
			lsize = node.left.size
		}
		res := e()
		if node.left != nil {
			res = dfs(node.left, left, left+lsize)
		}
		if tmp := left + lsize; start <= tmp && tmp < end {
			res = op(res, node.key)
		}
		if node.right != nil {
			res = op(res, dfs(node.right, left+lsize+1, right))
		}
		return res
	}
	return dfs(avl.root, 0, avl.Size())
}

func (avl *AVLTreeWithSum) QueryAll() E {
	if avl.root == nil {
		return e()
	}
	return avl.root.data
}

func (t *AVLTreeWithSum) Size() int32 {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func _newWithRoot(root *aNode) *AVLTreeWithSum {
	return &AVLTreeWithSum{root: root}
}

func (t *AVLTreeWithSum) _build(n int32, f func(int32) E) {
	var dfs func(l, r int32) *aNode
	dfs = func(l, r int32) *aNode {
		mid := (l + r) >> 1
		node := _newANode(f(mid))
		if l != mid {
			node.left = dfs(l, mid)
		}
		if mid+1 != r {
			node.right = dfs(mid+1, r)
		}
		t._update(node)
		return node
	}
	t.root = dfs(0, n)
}

func (t *AVLTreeWithSum) _update(node *aNode) {
	node.size = 1
	node.data = node.key
	node.height = 1
	if node.left != nil {
		node.size += node.left.size
		node.data = op(node.left.data, node.data)
		node.height = max8(node.left.height+1, 1)
	}
	if node.right != nil {
		node.size += node.right.size
		node.data = op(node.data, node.right.data)
		node.height = max8(node.height, node.right.height+1)
	}
}

func (t *AVLTreeWithSum) _balanceLeft(node *aNode) *aNode {
	var u *aNode
	if node.left.left == nil || node.left.left.height+2 == node.left.height {
		u = node.left.right
		node.left.right = u.left
		u.left = node.left
		node.left = u.right
		u.right = node
		t._update(u.left)
	} else {
		u = node.left
		node.left = u.right
		u.right = node
	}
	t._update(u.right)
	t._update(u)
	return u
}

func (t *AVLTreeWithSum) _balanceRight(node *aNode) *aNode {
	var u *aNode
	if node.right.right == nil || node.right.right.height+2 == node.right.height {
		u = node.right.left
		node.right.left = u.right
		u.right = node.right
		node.right = u.left
		u.left = node
		t._update(u.right)
	} else {
		u = node.right
		node.right = u.left
		u.left = node
	}
	t._update(u.left)
	t._update(u)
	return u
}

func (t *AVLTreeWithSum) _mergeWithRoot(l, root, r *aNode) *aNode {
	diff := int8(0)
	if l == nil {
		if r != nil {
			diff = -r.height
		}
	} else {
		if r == nil {
			diff = l.height
		} else {
			diff = l.height - r.height
		}
	}
	if diff > 1 {
		l.right = t._mergeWithRoot(l.right, root, r)
		t._update(l)
		if l.left == nil {
			if l.right.height == 2 {
				return t._balanceRight(l)
			}
		} else {
			if l.left.height-l.right.height == -2 {
				return t._balanceRight(l)
			}
		}
		return l
	} else if diff < -1 {
		r.left = t._mergeWithRoot(l, root, r.left)
		t._update(r)
		if r.right == nil {
			if r.left.height == 2 {
				return t._balanceLeft(r)
			}
		} else {
			if r.left.height-r.right.height == 2 {
				return t._balanceLeft(r)
			}
		}
		return r
	} else {
		root.left = l
		root.right = r
		t._update(root)
		return root
	}
}

func (t *AVLTreeWithSum) _mergeNode(l, r *aNode) *aNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	l, root := t._popRight(l)
	return t._mergeWithRoot(l, root, r)
}

func (t *AVLTreeWithSum) _popRight(node *aNode) (*aNode, *aNode) {
	tmpPath = tmpPath[:0]
	path := tmpPath
	mx := node
	for node.right != nil {
		path = append(path, node)
		mx = node.right
		node = node.right
	}
	path = append(path, node.left)
	len_ := len(path)
	for i := 0; i < len_-1; i++ {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if node == nil {
			path[len(path)-1].right = nil
			t._update(path[len(path)-1])
			continue
		}
		b := node.Balance()
		if b == 2 {
			path[len(path)-1].right = t._balanceLeft(node)
		} else if b == -2 {
			path[len(path)-1].right = t._balanceRight(node)
		} else {
			path[len(path)-1].right = node
		}
		t._update(path[len(path)-1])
	}
	if path[0] != nil {
		b := path[0].Balance()
		if b == 2 {
			path[0] = t._balanceLeft(path[0])
		} else if b == -2 {
			path[0] = t._balanceRight(path[0])
		}
	}
	mx.left = nil
	t._update(mx)
	return path[0], mx
}

func (t *AVLTreeWithSum) _splitNode(node *aNode, k int32) (*aNode, *aNode) {
	if node == nil {
		return nil, nil
	}
	tmp := k
	if node.left != nil {
		tmp -= node.left.size
	}
	if tmp == 0 {
		left := node.left
		right := t._mergeWithRoot(nil, node, node.right)
		return left, right
	}
	if tmp < 0 {
		left, right := t._splitNode(node.left, k)
		return left, t._mergeWithRoot(right, node, node.right)
	}
	left, right := t._splitNode(node.right, tmp-1)
	return t._mergeWithRoot(node.left, node, left), right
}

func max8(x, y int8) int8 {
	if x > y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func max32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func testTime() {
	n := int(2e5)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(100) + 1
	}

	tree := NewAVLTreeWithSum(int32(len(arr)), func(i int32) E { return E(arr[i]) })

	time1 := time.Now()
	for j := 0; j < int(2e5); j++ {
		tree.Set(int32(j), E(j))
		tree.Get(int32(j))
		tree.Insert(int32(j), E(j))
		tree.Pop(int32(j))
		tree.Query(int32(j), int32(n))
		tree.QueryAll()
		a, b := tree.Split(int32(j))
		a.Merge(b)
		tree = a
	}

	fmt.Println(time.Since(time1)) // 505.323958ms
}

func init() {
	debug.SetGCPercent(-1)
}
