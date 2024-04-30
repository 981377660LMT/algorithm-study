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
	test()
	testTime()
}

func init() {
	debug.SetGCPercent(-1)
}

type E = struct {
	sum  int
	size int32
}
type Id = int

func e() E   { return E{} }
func id() Id { return 0 }
func op(x, y E) E {
	x.sum += y.sum
	x.size += y.size
	return x
}
func mapping(f Id, x E) E {
	x.sum += int(x.size) * f
	return x
}
func composition(f, g Id) Id { return f + g }

type avlNode struct {
	key         E
	data        E
	left, right *avlNode
	lazy        Id
	height      int8
	size        int32
}

func _newAVLNode(key E, id Id) *avlNode {
	return &avlNode{key: key, data: key, lazy: id, height: 1, size: 1}
}

var (
	tmpPath  = make([]*avlNode, 0, 128)
	tmpStack = make([]stackItem, 0, 128)
)

type stackItem struct {
	node        *avlNode
	left, right int32
}
type LazyAVLTree struct {
	root *avlNode
}

func NewLazyAVLTree(n int32, f func(i int32) E) *LazyAVLTree {
	res := &LazyAVLTree{}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func _newWithRoot(root *avlNode) *LazyAVLTree {
	return &LazyAVLTree{root: root}
}

func (avl *LazyAVLTree) Merge(other *LazyAVLTree) {
	avl.root = avl._mergeNode(avl.root, other.root)
}

func (avl *LazyAVLTree) Split(k int32) (*LazyAVLTree, *LazyAVLTree) {
	l, r := avl._splitNode(avl.root, k)
	return _newWithRoot(l), _newWithRoot(r)
}

func (avl *LazyAVLTree) Insert(k int32, key E) {
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
	avl.root = avl._mergeWithRoot(s, _newAVLNode(key, id()), t)
}

func (avl *LazyAVLTree) Pop(k int32) E {
	if k < 0 {
		k += avl.Size()
	}
	s, t := avl._splitNode(avl.root, k+1)
	s, tmp := avl._popMax(s)
	avl.root = avl._mergeNode(s, t)
	return tmp.key
}

func (avl *LazyAVLTree) Get(k int32) E {
	return avl._kthElm(k)
}

func (avl *LazyAVLTree) Set(k int32, key E) {
	if k < 0 {
		k += avl.Size()
	}
	node := avl.root
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
			node = node.right
		} else {
			node = node.left
		}
	}
	for i := len(path) - 1; i >= 0; i-- {
		avl._update(path[i])
	}
}

func (avl *LazyAVLTree) Update(start, end int32, f Id) {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end || avl.root == nil {
		return
	}

	tmpStack = tmpStack[:0]
	stack := tmpStack
	root := avl.root
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
				stack = append(stack, stackItem{node.left, left, left + lsize})
			}
			if tmp := left + lsize; start <= tmp && tmp < end {
				node.key = mapping(f, node.key)
			}
			if node.right != nil {
				stack = append(stack, stackItem{node.right, left + lsize + 1, right})
			}
		}
	}
}

func (avl *LazyAVLTree) UpdateAll(f Id) {
	if avl.root == nil {
		return
	}
	root := avl.root
	root.key = mapping(f, root.key)
	root.data = mapping(f, root.data)
	if root.lazy == id() {
		root.lazy = f
	} else {
		root.lazy = composition(f, root.lazy)
	}
}

func (avl *LazyAVLTree) Query(start, end int32) E {
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

func (avl *LazyAVLTree) QueryAll() E {
	if avl.root == nil {
		return e()
	}
	return avl.root.data
}

func (avl *LazyAVLTree) Clear() { avl.root = nil }

func (avl *LazyAVLTree) ToList() []E {
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

func (avl *LazyAVLTree) Size() int32 {
	if avl.root == nil {
		return 0
	}
	return avl.root.size
}

func (avl *LazyAVLTree) _build(n int32, f func(i int32) E) {
	var dfs func(l, r int32) *avlNode
	dfs = func(l, r int32) *avlNode {
		mid := (l + r) >> 1
		node := _newAVLNode(f(mid), id())
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

func (avl *LazyAVLTree) _propagate(node *avlNode) {
	l, r := node.left, node.right
	if node.lazy != id() {
		lazy := node.lazy
		if l != nil {
			l.data = mapping(lazy, l.data)
			l.key = mapping(lazy, l.key)
			if l.lazy == id() {
				l.lazy = lazy
			} else {
				l.lazy = composition(lazy, l.lazy)
			}
		}
		if r != nil {
			r.data = mapping(lazy, r.data)
			r.key = mapping(lazy, r.key)
			if r.lazy == id() {
				r.lazy = lazy
			} else {
				r.lazy = composition(lazy, r.lazy)
			}
		}
		node.lazy = id()
	}
}

func (avl *LazyAVLTree) _update(node *avlNode) {
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

func (avl *LazyAVLTree) _getBalance(node *avlNode) int8 {
	if node.left == nil {
		if node.right == nil {
			return 0
		}
		return -node.right.height
	}
	if node.right == nil {
		return node.left.height
	}
	return node.left.height - node.right.height
}

func (avl *LazyAVLTree) _balanceLeft(node *avlNode) *avlNode {
	avl._propagate(node.left)
	var u *avlNode
	if node.left.left == nil || node.left.left.height+2 == node.left.height {
		u = node.left.right
		avl._propagate(u)
		node.left.right = u.left
		u.left = node.left
		node.left = u.right
		u.right = node
		avl._update(u.left)
	} else {
		u = node.left
		node.left = u.right
		u.right = node
	}
	avl._update(u.right)
	avl._update(u)
	return u
}

func (avl *LazyAVLTree) _balanceRight(node *avlNode) *avlNode {
	avl._propagate(node.right)
	var u *avlNode
	if node.right.right == nil || node.right.right.height+2 == node.right.height {
		u = node.right.left
		avl._propagate(u)
		node.right.left = u.right
		u.right = node.right
		node.right = u.left
		u.left = node
		avl._update(u.right)
	} else {
		u = node.right
		node.right = u.left
		u.left = node
	}
	avl._update(u.left)
	avl._update(u)
	return u
}

func (avl *LazyAVLTree) _kthElm(k int32) E {
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

func (avl *LazyAVLTree) _mergeWithRoot(l, root, r *avlNode) *avlNode {
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
		avl._propagate(l)
		l.right = avl._mergeWithRoot(l.right, root, r)
		avl._update(l)
		if l.left == nil {
			if l.right.height == 2 {
				return avl._balanceRight(l)
			}
		} else {
			if l.left.height-l.right.height == -2 {
				return avl._balanceRight(l)
			}
		}
		return l
	} else if diff < -1 {
		avl._propagate(r)
		r.left = avl._mergeWithRoot(l, root, r.left)
		avl._update(r)
		if r.right == nil {
			if r.left.height == 2 {
				return avl._balanceLeft(r)
			}
		} else {
			if r.left.height-r.right.height == 2 {
				return avl._balanceLeft(r)
			}
		}
		return r
	} else {
		root.left = l
		root.right = r
		avl._update(root)
		return root
	}
}

func (avl *LazyAVLTree) _mergeNode(l, r *avlNode) *avlNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	l, tmp := avl._popMax(l)
	return avl._mergeWithRoot(l, tmp, r)
}

func (avl *LazyAVLTree) _popMax(node *avlNode) (*avlNode, *avlNode) {
	avl._propagate(node)
	tmpPath = tmpPath[:0]
	path := tmpPath
	mx := node
	for node.right != nil {
		path = append(path, node)
		mx = node.right
		node = node.right
		avl._propagate(node)
	}
	path = append(path, node.left)
	len_ := len(path)
	for i := 0; i < len_-1; i++ {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if node == nil {
			path[len(path)-1].right = nil
			avl._update(path[len(path)-1])
			continue
		}
		b := avl._getBalance(node)
		if b == 2 {
			path[len(path)-1].right = avl._balanceLeft(node)
		} else if b == -2 {
			path[len(path)-1].right = avl._balanceRight(node)
		} else {
			path[len(path)-1].right = node
		}
		avl._update(path[len(path)-1])
	}
	if path[0] != nil {
		b := avl._getBalance(path[0])
		if b == 2 {
			path[0] = avl._balanceLeft(path[0])
		} else if b == -2 {
			path[0] = avl._balanceRight(path[0])
		}
	}
	mx.left = nil
	avl._update(mx)
	return path[0], mx
}

func (avl *LazyAVLTree) _splitNode(node *avlNode, k int32) (*avlNode, *avlNode) {
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
	tree := NewLazyAVLTree(int32(len(arr)), func(i int32) E { return E{arr[i], 1} })
	fmt.Println(tree.ToList())
	fmt.Println(tree.Pop(1))
	fmt.Println(tree.ToList())

	for i := 0; i < 1000; i++ {
		n := rand.Intn(1000) + 500
		arr = make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100) + 1
		}
		tree = NewLazyAVLTree(int32(len(arr)), func(i int32) E { return E{arr[i], 1} })

		for j := 0; j < 500; j++ {
			// set
			{
				index, value := rand.Intn(len(arr)), rand.Intn(100)+1
				arr[index] = value
				tree.Set(int32(index), E{value, 1})
			}

			// get
			{
				index := rand.Intn(len(arr))
				if arr[index] != int(tree.Get(int32(index)).sum) {
					panic("error get")
				}
			}

			// insert
			{
				index, value := rand.Intn(len(arr)+1), rand.Intn(100)+1
				arr = append(arr, 0)
				copy(arr[index+1:], arr[index:])
				arr[index] = value
				tree.Insert(int32(index), E{value, 1})

			}

			// pop
			{
				index := rand.Intn(len(arr))
				tmp1 := arr[index]
				tmp2 := tree.Pop(int32(index)).sum
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
				tree.Update(int32(l), int32(r), Id(f))
			}

			// apply all
			{
				f := rand.Intn(100) + 1
				for k := 0; k < len(arr); k++ {
					arr[k] += f
				}
				tree.UpdateAll(Id(f))
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
				tmp2 := tree.Query(int32(l), int32(r)).sum
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
				tmp2 := tree.QueryAll().sum
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
					if tmp1[i] != int(tmp2[i].sum) {
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

	tree := NewLazyAVLTree(int32(len(arr)), func(i int32) E { return E{arr[i], 1} })

	time1 := time.Now()
	for j := 0; j < int(2e5); j++ {
		tree.Set(int32(j), E{arr[j], 1})
		tree.Get(int32(j))
		tree.Insert(int32(j), E{arr[j], 1})
		tree.Pop(int32(j))
		tree.Update(int32(j), int32(n), Id(j))
		tree.UpdateAll(Id(j))
		tree.Query(int32(j), int32(n))
		tree.QueryAll()
		a, b := tree.Split(int32(j))
		a.Merge(b)
		tree = a
	}
	fmt.Println(time.Since(time1)) // 297.617458ms
}
