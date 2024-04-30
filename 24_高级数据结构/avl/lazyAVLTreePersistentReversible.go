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
//   Reverse(start, end int32)
//   ReverseAll()
//   Query(start, end int32) E
//   QueryAll() E
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
	test()
	testTime()
}

func init() {
	debug.SetGCPercent(-1)
}

type E = int32
type Id = int32

func e() E                   { return 0 }
func id() Id                 { return 0 }
func op(x, y E) E            { return max32(x, y) }
func mapping(f Id, x E) E    { return x + f }
func composition(f, g Id) Id { return f + g }

type avlNode struct {
	key         E
	data        E
	left, right *avlNode
	lazy        Id
	rev         bool
	height      int8
	size        int32
}

func newAVLNode(key E, id Id) *avlNode {
	return &avlNode{key: key, data: key, lazy: id, height: 1, size: 1}
}

func (n *avlNode) Balance() int8 {
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

func (n *avlNode) Copy() *avlNode {
	if n == nil {
		return nil
	}
	return &avlNode{
		key: n.key, data: n.data, left: n.left, right: n.right,
		lazy: n.lazy, rev: n.rev, height: n.height, size: n.size,
	}
}

var (
	tmpPath  = make([]*avlNode, 0, 128)
	tmpStack = make([]stackItem, 0, 128)
)

type stackItem struct {
	node        *avlNode
	left, right int32
}
type LazyAVLTreePersistentReverible struct {
	root *avlNode
}

func NewLazyAVLTreePersistentReversible(n int32, f func(i int32) E) *LazyAVLTreePersistentReverible {
	res := &LazyAVLTreePersistentReverible{}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func _newWithRoot(root *avlNode) *LazyAVLTreePersistentReverible {
	return &LazyAVLTreePersistentReverible{root: root}
}

func (avl *LazyAVLTreePersistentReverible) Merge(other *LazyAVLTreePersistentReverible) *LazyAVLTreePersistentReverible {
	root := avl._mergeNode(avl.root, other.root)
	return _newWithRoot(root)
}

func (avl *LazyAVLTreePersistentReverible) Split(k int32) (*LazyAVLTreePersistentReverible, *LazyAVLTreePersistentReverible) {
	l, r := avl._splitNode(avl.root, k)
	return _newWithRoot(l), _newWithRoot(r)
}

func (avl *LazyAVLTreePersistentReverible) Insert(k int32, key E) *LazyAVLTreePersistentReverible {
	if k < 0 {
		k += avl.Size()
	}
	if k < 0 {
		k = 0
	}
	if n := avl.Size(); k > n {
		k = n
	}
	s, t := avl._splitNode(avl.root, k)
	root := avl._mergeWithRoot(s, newAVLNode(key, id()), t)
	return _newWithRoot(root)
}

func (avl *LazyAVLTreePersistentReverible) Pop(k int32) (*LazyAVLTreePersistentReverible, E) {
	if k < 0 {
		k += avl.Size()
	}
	s, t := avl._splitNode(avl.root, k+1)
	s, tmp := avl._popMax(s)
	root := avl._mergeNode(s, t)
	return _newWithRoot(root), tmp.key
}

func (avl *LazyAVLTreePersistentReverible) Get(k int32) E {
	return avl._kthElm(k)
}

func (avl *LazyAVLTreePersistentReverible) Set(k int32, key E) *LazyAVLTreePersistentReverible {
	if k < 0 {
		k += avl.Size()
	}
	newRoot := avl.root.Copy()
	node := newRoot
	tmpPath = tmpPath[:0]
	path := tmpPath
	for {
		avl._propagate(node)
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
			newNode := node.right.Copy()
			node.right = newNode
			node = newNode
		} else {
			newNode := node.left.Copy()
			node.left = newNode
			node = newNode
		}
	}
	for i := len(path) - 1; i >= 0; i-- {
		avl._update(path[i])
	}
	return _newWithRoot(newRoot)
}

func (avl *LazyAVLTreePersistentReverible) Update(start, end int32, f Id) *LazyAVLTreePersistentReverible {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end || avl.root == nil {
		return _newWithRoot(avl.root.Copy())
	}
	root := avl.root.Copy()
	tmpStack = tmpStack[:0]
	stack := tmpStack
	stack = append(stack, stackItem{root, -1, -1}, stackItem{root, 0, avl.Size()})
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		node, left, right := top.node, top.left, top.right
		if left == -1 {
			avl._update(node)
			continue
		}
		if right <= start || end <= left {
			continue
		}
		avl._propagate(node)
		if start <= left && right < end {
			node.key = mapping(f, node.key)
			node.data = mapping(f, node.data)
			if node.lazy == id() {
				node.lazy = f
			} else {
				node.lazy = composition(f, node.lazy)
			}
		} else {
			lsize := int32(0)
			if node.left != nil {
				lsize = node.left.size
			}
			stack = append(stack, stackItem{node, -1, -1})
			if node.left != nil {
				leftCopy := node.left.Copy()
				node.left = leftCopy
				stack = append(stack, stackItem{leftCopy, left, left + lsize})
			}
			if tmp := left + lsize; start <= tmp && tmp < end {
				node.key = mapping(f, node.key)
			}
			if node.right != nil {
				rightCopy := node.right.Copy()
				node.right = rightCopy
				stack = append(stack, stackItem{rightCopy, left + lsize + 1, right})
			}
		}
	}
	return _newWithRoot(root)
}

func (avl *LazyAVLTreePersistentReverible) UpdateAll(f Id) *LazyAVLTreePersistentReverible {
	if avl.root == nil {
		return _newWithRoot(nil)
	}
	root := avl.root.Copy()
	root.key = mapping(f, root.key)
	root.data = mapping(f, root.data)
	if root.lazy == id() {
		root.lazy = f
	} else {
		root.lazy = composition(f, root.lazy)
	}
	return _newWithRoot(root)
}

func (avl *LazyAVLTreePersistentReverible) Reverse(start, end int32) *LazyAVLTreePersistentReverible {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end {
		return _newWithRoot(avl.root.Copy())
	}
	s, t := avl._splitNode(avl.root, end)
	r, s := avl._splitNode(s, start)
	s.rev = !s.rev
	root := avl._mergeNode(avl._mergeNode(r, s), t)
	return _newWithRoot(root)
}

func (avl *LazyAVLTreePersistentReverible) ReverseAll() *LazyAVLTreePersistentReverible {
	if avl.root == nil {
		return _newWithRoot(nil)
	}
	root := avl.root.Copy()
	root.rev = !root.rev
	return _newWithRoot(root)
}

func (avl *LazyAVLTreePersistentReverible) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end || avl.root == nil {
		return e()
	}
	var dfs func(node *avlNode, left, right int32) E
	dfs = func(node *avlNode, left, right int32) E {
		if right <= start || end <= left {
			return e()
		}
		avl._propagate(node)
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

func (avl *LazyAVLTreePersistentReverible) QueryAll() E {
	if avl.root == nil {
		return e()
	}
	return avl.root.data
}

func (avl *LazyAVLTreePersistentReverible) ToList() []E {
	node := avl.root
	stack := []*avlNode{}
	res := make([]E, 0, avl.Size())
	for len(stack) > 0 || node != nil {
		if node != nil {
			avl._propagate(node)
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

func (avl *LazyAVLTreePersistentReverible) Size() int32 {
	if avl.root == nil {
		return 0
	}
	return avl.root.size
}

func (avl *LazyAVLTreePersistentReverible) _build(n int32, f func(i int32) E) {
	var dfs func(l, r int32) *avlNode
	dfs = func(l, r int32) *avlNode {
		mid := (l + r) >> 1
		node := newAVLNode(f(mid), id())
		if l != mid {
			node.left = dfs(l, mid)
		}
		if mid+1 != r {
			node.right = dfs(mid+1, r)
		}
		avl._update(node)
		return node
	}
	avl.root = dfs(0, n)
}

func (avl *LazyAVLTreePersistentReverible) _propagate(node *avlNode) {
	hasCopy := false
	if node.rev {
		node.rev = false
		l := node.left.Copy()
		r := node.right.Copy()
		hasCopy = true
		node.left, node.right = r, l
		if l != nil {
			l.rev = !l.rev
		}
		if r != nil {
			r.rev = !r.rev
		}
	}

	if !hasCopy {
		if node.lazy != id() {
			lazy := node.lazy
			node.lazy = id()
			if node.left != nil {
				l := node.left.Copy()
				l.data = mapping(lazy, l.data)
				l.key = mapping(lazy, l.key)
				l.lazy = composition(lazy, l.lazy)
				node.left = l
			}
			if node.right != nil {
				r := node.right.Copy()
				r.data = mapping(lazy, r.data)
				r.key = mapping(lazy, r.key)
				r.lazy = composition(lazy, r.lazy)
				node.right = r
			}
		}
	} else {
		if node.lazy != id() {
			lazy := node.lazy
			node.lazy = id()
			if node.left != nil {
				l := node.left
				l.data = mapping(lazy, l.data)
				l.key = mapping(lazy, l.key)
				l.lazy = composition(lazy, l.lazy)
			}
			if node.right != nil {
				r := node.right
				r.data = mapping(lazy, r.data)
				r.key = mapping(lazy, r.key)
				r.lazy = composition(lazy, r.lazy)
			}
		}
	}
}

func (avl *LazyAVLTreePersistentReverible) _update(node *avlNode) {
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

func (t *LazyAVLTreePersistentReverible) _rotateRight(node *avlNode) *avlNode {
	u := node.left.Copy()
	node.left = u.right
	u.right = node
	t._update(node)
	t._update(u)
	return u
}

func (t *LazyAVLTreePersistentReverible) _rotateLeft(node *avlNode) *avlNode {
	u := node.right.Copy()
	node.right = u.left
	u.left = node
	t._update(node)
	t._update(u)
	return u
}

func (t *LazyAVLTreePersistentReverible) _balanceLeft(node *avlNode) *avlNode {
	t._propagate(node.right)
	node.right = node.right.Copy()
	u := node.right
	if u.Balance() == 1 {
		t._propagate(u.left)
		node.right = t._rotateRight(u)
	}
	u = t._rotateLeft(node)
	return u
}

func (t *LazyAVLTreePersistentReverible) _balanceRight(node *avlNode) *avlNode {
	t._propagate(node.left)
	node.left = node.left.Copy()
	u := node.left
	if u.Balance() == -1 {
		t._propagate(u.right)
		node.left = t._rotateLeft(u)
	}
	u = t._rotateRight(node)
	return u
}

func (avl *LazyAVLTreePersistentReverible) _kthElm(k int32) E {
	if k < 0 {
		k += avl.Size()
	}
	node := avl.root
	for {
		avl._propagate(node)
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

func (avl *LazyAVLTreePersistentReverible) _mergeWithRoot(l, root, r *avlNode) *avlNode {
	diff := int8(0)
	if l != nil {
		diff += l.height
	}
	if r != nil {
		diff -= r.height
	}
	if diff > 1 {
		avl._propagate(l)
		l = l.Copy()
		l.right = avl._mergeWithRoot(l.right, root, r)
		avl._update(l)
		if l.Balance() == -2 {
			return avl._balanceLeft(l)
		}
		return l
	}
	if diff < -1 {
		avl._propagate(r)
		r = r.Copy()
		r.left = avl._mergeWithRoot(l, root, r.left)
		avl._update(r)
		if r.Balance() == 2 {
			return avl._balanceRight(r)
		}
		return r
	}
	root = root.Copy()
	root.left = l
	root.right = r
	avl._update(root)
	return root

}

func (avl *LazyAVLTreePersistentReverible) _mergeNode(l, r *avlNode) *avlNode {
	if l == nil && r == nil {
		return nil
	}
	if l == nil {
		return r.Copy()
	}
	if r == nil {
		return l.Copy()
	}
	l, r = l.Copy(), r.Copy()
	l, tmp := avl._popMax(l)
	return avl._mergeWithRoot(l, tmp, r)
}

func (avl *LazyAVLTreePersistentReverible) _popMax(node *avlNode) (*avlNode, *avlNode) {
	avl._propagate(node)
	node = node.Copy()
	tmpPath = tmpPath[:0]
	path := tmpPath
	mx := node
	for node.right != nil {
		path = append(path, node)
		avl._propagate(node.right)
		node = node.right.Copy()
		mx = node
	}
	path = append(path, node.left.Copy())
	len_ := len(path)
	for i := 0; i < len_-1; i++ {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if node == nil {
			path[len(path)-1].right = nil
			avl._update(path[len(path)-1])
			continue
		}
		b := node.Balance()
		if b == 2 {
			path[len(path)-1].right = avl._balanceRight(node)
		} else if b == -2 {
			path[len(path)-1].right = avl._balanceLeft(node)
		} else {
			path[len(path)-1].right = node
		}
		avl._update(path[len(path)-1])
	}

	if path[0] != nil {
		b := path[0].Balance()
		if b == 2 {
			path[0] = avl._balanceRight(path[0])
		} else if b == -2 {
			path[0] = avl._balanceLeft(path[0])
		}
	}
	mx.left = nil
	avl._update(mx)
	return path[0], mx
}

func (avl *LazyAVLTreePersistentReverible) _splitNode(node *avlNode, k int32) (*avlNode, *avlNode) {
	if node == nil {
		return nil, nil
	}
	avl._propagate(node)
	tmp := k
	if node.left != nil {
		tmp -= node.left.size
	}
	if tmp == 0 {
		res1 := node.left
		res2 := avl._mergeWithRoot(nil, node, node.right)
		return res1, res2
	}
	if tmp < 0 {
		s, t := avl._splitNode(node.left, k)
		return s, avl._mergeWithRoot(t, node, node.right)
	}
	s, t := avl._splitNode(node.right, tmp-1)
	return avl._mergeWithRoot(node.left, node, s), t
}

func max8(x, y int8) int8 {
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

func min32(x, y int32) int32 {
	if x < y {
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

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func test() {
	arr := []int{1, 2, 3, 4, 5}
	tree := NewLazyAVLTreePersistentReversible(int32(len(arr)), func(i int32) E { return E(arr[i]) })
	fmt.Println(tree.ToList())
	tree = tree.Set(0, 10)
	fmt.Println(tree.ToList())

	for i := 0; i < 200; i++ {
		n := rand.Intn(1000) + 500
		arr = make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100) + 1
		}
		tree = NewLazyAVLTreePersistentReversible(int32(len(arr)), func(i int32) E { return E(arr[i]) })

		for j := 0; j < 500; j++ {
			// set
			{
				index, value := rand.Intn(len(arr)), rand.Intn(100)+1
				arr[index] = value
				tree = tree.Set(int32(index), E(value))
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
				tree = tree.Insert(int32(index), E(value))
			}

			// pop
			{
				index := rand.Intn(len(arr))
				tmp1 := arr[index]
				newTree, tmp2 := tree.Pop(int32(index))
				tree = newTree
				arr = append(arr[:index], arr[index+1:]...)
				if tmp1 != int(tmp2) {
					fmt.Println(tmp1, tmp2)
					panic("error pop")
				}
			}

			// apply
			{
				l, r, f := rand.Intn(len(arr)), rand.Intn(len(arr)), rand.Intn(100)+1
				for k := l; k < r; k++ {
					arr[k] += f
				}
				tree = tree.Update(int32(l), int32(r), Id(f))
			}

			// apply all
			{
				f := rand.Intn(100) + 1
				for k := 0; k < len(arr); k++ {
					arr[k] += f
				}
				tree = tree.UpdateAll(Id(f))
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
				tree = tree.Reverse(int32(l), int32(r))
			}

			// reverse all
			{
				for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
					arr[i], arr[j] = arr[j], arr[i]
				}
				tree = tree.ReverseAll()
			}

			// prod
			{
				l, r := rand.Intn(len(arr)), rand.Intn(len(arr))
				if l > r {
					l, r = r, l
				}
				tmp1 := 0
				for i := l; i < r; i++ {
					tmp1 = max(tmp1, arr[i])
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
					tmp1 = max(tmp1, arr[i])
				}
				tmp2 := tree.QueryAll()
				if tmp1 != int(tmp2) {
					fmt.Println(tmp1, tmp2)
					panic("error all prod")
				}
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

			// split + merge
			{

				l, r := rand.Intn(len(arr)), rand.Intn(len(arr))
				if l > r {
					l, r = r, l
				}
				lTree, rTree := tree.Split(int32(l))
				lTree = lTree.Merge(rTree)
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

	tree := NewLazyAVLTreePersistentReversible(int32(len(arr)), func(i int32) E { return E(arr[i]) })

	time1 := time.Now()
	for j := 0; j < int(2e5); j++ {
		tree = tree.Set(int32(j), E(j))
		tree.Get(int32(j))
		tree = tree.Insert(int32(j), E(j))
		newTree, _ := tree.Pop(int32(j))
		tree = newTree
		tree = tree.Update(int32(j), int32(n), Id(j))
		tree = tree.Reverse(int32(j), int32(n))
		tree = tree.ReverseAll()
		tree.Query(int32(j), int32(n))
		tree.QueryAll()
		a, b := tree.Split(int32(j))
		a = a.Merge(b)
		tree = a
	}
	fmt.Println(time.Since(time1)) // 2.903513375s
}
