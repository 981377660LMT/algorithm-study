// AVL树.
// api:
//  NewLazyAVLTreePersistent(n int32, f func(int32) V) *AVLTreePersistent
//  Merge(other *AVLTreePersistent)
//  Insert(k int32, key V)
//  Split(k int32) (*AVLTreePersistent, *AVLTreePersistent)
//  Pop(k int32) (V)
//  Get(k int32) V
//  Set(k int32, key V)
//  Query(start, end int32) E
//  QueryAll() E
//  Reverse(start, end int32)
//  ReverseAll()
//  Clear()
//  ToList() []V
//  Size() int32

// !这里的Monoid可以不满足交换律(commutative)

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	test()
	testTime()
}

type E = int32

func e() E        { return 0 }
func op(x, y E) E { return x + y }

type aNode struct {
	key         E
	data        E
	revData     E
	left, right *aNode
	height      int8
	size        int32
	rev         bool
}

func _newANode(key E) *aNode {
	return &aNode{key: key, data: key, revData: key, height: 1, size: 1}
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

var (
	tmpPath = make([]*aNode, 0, 128)
)

type avlTreeWithSumReversible struct {
	root *aNode
	e    func() E
	op   func(E, E) E
}

func newAVLTreeWithSumReversible(n int32, f func(int32) E, e func() E, op func(E, E) E) *avlTreeWithSumReversible {
	res := &avlTreeWithSumReversible{e: e, op: op}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (t *avlTreeWithSumReversible) Merge(other *avlTreeWithSumReversible) {
	t.root = t._mergeNode(t.root, other.root)
}

func (t *avlTreeWithSumReversible) Insert(k int32, key E) {
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

func (t *avlTreeWithSumReversible) Split(k int32) (*avlTreeWithSumReversible, *avlTreeWithSumReversible) {
	a, b := t._splitNode(t.root, k)
	return _newWithRoot(a, t.e, t.op), _newWithRoot(b, t.e, t.op)
}

func (t *avlTreeWithSumReversible) Pop(k int32) E {
	if k < 0 {
		k += t.Size()
	}
	a, b := t._splitNode(t.root, k+1)
	a, tmp := t._popRight(a)
	t.root = t._mergeNode(a, b)
	return tmp.key
}

func (t *avlTreeWithSumReversible) Set(k int32, key E) {
	if k < 0 {
		k += t.Size()
	}
	node := t.root
	tmpPath = tmpPath[:0]
	path := tmpPath
	for {
		t._propagate(node)
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

func (t *avlTreeWithSumReversible) Clear() { t.root = nil }

func (t *avlTreeWithSumReversible) ToList() []E {
	node := t.root
	stack := make([]*aNode, 0)
	res := make([]E, 0, t.Size())
	for len(stack) > 0 || node != nil {
		if node != nil {
			t._propagate(node)
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

func (t *avlTreeWithSumReversible) Get(k int32) E {
	if k < 0 {
		k += t.Size()
	}
	node := t.root
	for {
		t._propagate(node)
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
func (avl *avlTreeWithSumReversible) Reverse(start, end int32) {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end {
		return
	}
	if start == 0 && end == avl.Size() {
		avl.ReverseAll()
		return
	}
	nm, nr := avl._splitNode(avl.root, end)
	nl, nm := avl._splitNode(nm, start)
	nm.rev = !nm.rev
	nm.data, nm.revData = nm.revData, nm.data
	avl.root = avl._mergeNode(avl._mergeNode(nl, nm), nr)
}

func (avl *avlTreeWithSumReversible) ReverseAll() {
	if avl.root == nil {
		return
	}
	avl.root.rev = !avl.root.rev
}

func (avl *avlTreeWithSumReversible) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end || avl.root == nil {
		return avl.e()
	}
	var dfs func(node *aNode, left, right int32) E
	dfs = func(node *aNode, left, right int32) E {
		if right <= start || end <= left {
			return avl.e()
		}
		avl._propagate(node)
		if start <= left && right < end {
			if node.rev {
				return node.revData
			}
			return node.data
		}
		lsize := int32(0)
		if node.left != nil {
			lsize = node.left.size
		}
		res := avl.e()
		if node.left != nil {
			res = dfs(node.left, left, left+lsize)
		}
		if tmp := left + lsize; start <= tmp && tmp < end {
			res = avl.op(res, node.key)
		}
		if node.right != nil {
			res = avl.op(res, dfs(node.right, left+lsize+1, right))
		}
		return res
	}
	return dfs(avl.root, 0, avl.Size())
}

func (avl *avlTreeWithSumReversible) QueryAll() E {
	if avl.root == nil {
		return avl.e()
	}
	return avl.root.data
}

func (t *avlTreeWithSumReversible) Size() int32 {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func _newWithRoot(root *aNode, e func() E, op func(a, b E) E) *avlTreeWithSumReversible {
	return &avlTreeWithSumReversible{root: root, e: e, op: op}
}

func (t *avlTreeWithSumReversible) _build(n int32, f func(int32) E) {
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

func (t *avlTreeWithSumReversible) _update(node *aNode) {
	node.size = 1
	node.data = node.key
	node.revData = node.key
	node.height = 1
	if left := node.left; left != nil {
		node.size += left.size
		node.data = t.op(left.data, node.data)
		node.revData = t.op(node.revData, left.revData)
		node.height = max8(left.height+1, 1)
	}
	if right := node.right; right != nil {
		node.size += right.size
		node.data = t.op(node.data, right.data)
		node.revData = t.op(right.revData, node.revData)
		node.height = max8(node.height, right.height+1)
	}
}

func (t *avlTreeWithSumReversible) _rotateRight(node *aNode) *aNode {
	u := node.left
	node.left = u.right
	u.right = node
	t._update(node)
	t._update(u)
	return u
}

func (t *avlTreeWithSumReversible) _rotateLeft(node *aNode) *aNode {
	u := node.right
	node.right = u.left
	u.left = node
	t._update(node)
	t._update(u)
	return u
}

func (t *avlTreeWithSumReversible) _balanceLeft(node *aNode) *aNode {
	t._propagate(node.left)
	var u *aNode
	if node.left.left == nil || node.left.left.height+2 == node.left.height {
		u = node.left.right
		t._propagate(u)
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

func (t *avlTreeWithSumReversible) _balanceRight(node *aNode) *aNode {
	t._propagate(node.right)
	var u *aNode
	if node.right.right == nil || node.right.right.height+2 == node.right.height {
		u = node.right.left
		t._propagate(u)
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

func (t *avlTreeWithSumReversible) _mergeWithRoot(l, root, r *aNode) *aNode {
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
		t._propagate(l)
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
		t._propagate(r)
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

func (t *avlTreeWithSumReversible) _mergeNode(l, r *aNode) *aNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	l, root := t._popRight(l)
	return t._mergeWithRoot(l, root, r)
}

func (t *avlTreeWithSumReversible) _popRight(node *aNode) (*aNode, *aNode) {
	t._propagate(node)
	tmpPath = tmpPath[:0]
	path := tmpPath
	mx := node
	for node.right != nil {
		path = append(path, node)
		mx = node.right
		node = node.right
		t._propagate(node)
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

func (t *avlTreeWithSumReversible) _splitNode(node *aNode, k int32) (*aNode, *aNode) {
	if node == nil {
		return nil, nil
	}
	t._propagate(node)
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

func (avl *avlTreeWithSumReversible) _propagate(node *aNode) {
	if node == nil {
		return
	}
	l, r := node.left, node.right
	if node.rev {
		node.left, node.right = r, l
		if l != nil {
			l.rev = !l.rev
			l.data, l.revData = l.revData, l.data
		}
		if r != nil {
			r.rev = !r.rev
			r.data, r.revData = r.revData, r.data
		}
		node.rev = false
	}
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

func test() {
	arr := []int{1, 2, 3, 4, 5}
	tree := newAVLTreeWithSumReversible(
		int32(len(arr)), func(i int32) E { return E(arr[i]) },
		func() E { return 0 }, func(x, y E) E { return x + y },
	)
	tree.ReverseAll()

	for i := 0; i < 500; i++ {
		n := rand.Intn(1000) + 1000
		arr = make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100) + 1
		}
		tree = newAVLTreeWithSumReversible(
			int32(len(arr)), func(i int32) E { return E(arr[i]) },
			func() E { return 0 }, func(x, y E) E { return x + y },
		)

		for j := 0; j < 1000; j++ {
			// set
			{
				index, value := rand.Intn(len(arr)), rand.Intn(100)+1
				arr[index] = value
				tree.Set(int32(index), E(value))
			}

			// get
			{
				index := rand.Intn(len(arr))
				if arr[index] != int(tree.Get(int32(index))) {
					panic("error get")
				}
			}

			// insert
			{
				index, value := rand.Intn(len(arr)+1), rand.Intn(100)+1
				arr = append(arr, 0)
				copy(arr[index+1:], arr[index:])
				arr[index] = value
				tree.Insert(int32(index), E(value))
			}

			// pop
			{
				index := rand.Intn(len(arr))
				tmp1 := arr[index]
				tmp2 := tree.Pop(int32(index))
				arr = append(arr[:index], arr[index+1:]...)
				if tmp1 != int(tmp2) {
					fmt.Println(tmp1, tmp2)
					panic("error pop")
				}
			}

			// reverse
			{
				l, r := rand.Intn(len(arr)), rand.Intn(len(arr))
				if l > r {
					l, r = r, l
				}
				for i, j := l, r-1; i < j; i, j = i+1, j-1 {
					arr[i], arr[j] = arr[j], arr[i]
				}
				tree.Reverse(int32(l), int32(r))
			}

			// reverse all
			{
				for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
					arr[i], arr[j] = arr[j], arr[i]
				}
				tree.ReverseAll()
			}

			// tolist
			{

				tmp1 := make([]int, len(arr))
				copy(tmp1, arr)
				tmp2 := tree.ToList()
				if len(tmp1) != len(tmp2) {
					fmt.Println(tmp1, tmp2)
					panic("error tolist")
				}
				for i := 0; i < len(tmp1); i++ {
					if tmp1[i] != int(tmp2[i]) {
						fmt.Println(tmp1, tmp2)
						panic("error tolist")
					}
				}
			}

			// prod
			{
				l, r := rand.Intn(len(arr)), rand.Intn(len(arr))
				if l > r {
					l, r = r, l
				}
				tmp1 := 0
				for i := l; i < r; i++ {
					tmp1 += arr[i]
				}
				tmp2 := tree.Query(int32(l), int32(r))
				if tmp1 != int(tmp2) {
					fmt.Println(tmp1, tmp2)
					panic("error prod")
				}
			}

			// all prod
			{
				tmp1 := 0
				for i := 0; i < len(arr); i++ {
					tmp1 += arr[i]
				}
				tmp2 := tree.QueryAll()
				if tmp1 != int(tmp2) {
					fmt.Println(tmp1, tmp2)
					panic("error all prod")
				}
			}

			// split + merge
			{
				l, r := rand.Intn(len(arr)), rand.Intn(len(arr))
				if l > r {
					l, r = r, l
				}
				lTree, rTree := tree.Split(int32(l))
				lTree.Merge(rTree)
				tree = lTree
			}
		}

	}

	fmt.Println("Passed")
}

func testTime() {
	n := int(2e5)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(100) + 1
	}

	tree := newAVLTreeWithSumReversible(
		int32(len(arr)), func(i int32) E { return E(arr[i]) },
		func() E { return 0 }, func(x, y E) E { return x + y },
	)

	time1 := time.Now()
	for j := 0; j < int(2e5); j++ {
		tree.Set(int32(j), E(j))
		tree.Get(int32(j))
		tree.Insert(int32(j), E(j))
		tree.Pop(int32(j))
		a, b := tree.Split(int32(j))
		a.Merge(b)
		tree = a
		tree.ReverseAll()
		tree.Reverse(int32(j), int32(n))
		tree.QueryAll()
		tree.Query(int32(j), int32(n))
	}
	fmt.Println(time.Since(time1)) // 313.778916ms
}
