// 可持久化AVL树.
// api:
//  NewLazyAVLTreePersistent(n int32, f func(int32) V) *AVLTreePersistent
//  Merge(other *AVLTreePersistent) *AVLTreePersistent
//  Insert(k int32, key V) *AVLTreePersistent
//  Split(k int32) (*AVLTreePersistent, *AVLTreePersistent)
//  Pop(k int32) (*AVLTreePersistent, V)
//  Get(k int32) V
//  Set(k int32, key V) *AVLTreePersistent
//  ToList() []V
//  Size() int32
//	CopyWithin(tree *AVLTreePersistent, target, start, end int32) *AVLTreePersistent

package main

import (
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	test()
	testTime()
	testCopyWithin()
	// arc030_4()
}

func CopyWithin(tree *AVLTreePersistent, target, start, end int32) *AVLTreePersistent {
	len_ := end - start
	p1Left, p1Right := tree.Split(start)
	p2Left, p2Right := p1Right.Split(len_)
	newTree := p1Left.Merge(p2Left).Merge(p2Right)
	p3Left, p3Right := newTree.Split(target)
	_, p4Right := p3Right.Split(len_)
	newTree = p3Left.Merge(p2Left).Merge(p4Right)
	return newTree
}

type V = int32
type aNode struct {
	key         V
	left, right *aNode
	height      int8
	size        int32
}

func _newANode(key V) *aNode {
	return &aNode{key: key, height: 1, size: 1}
}

func (n *aNode) Copy() *aNode {
	return &aNode{key: n.key, left: n.left, right: n.right, height: n.height, size: n.size}
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

type AVLTreePersistent struct {
	root *aNode
}

func NewLazyAVLTreePersistent(n int32, f func(int32) V) *AVLTreePersistent {
	res := &AVLTreePersistent{}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (t *AVLTreePersistent) Merge(other *AVLTreePersistent) *AVLTreePersistent {
	root := t._mergeNode(t.root, other.root)
	return _newWithRoot(root)
}

func (t *AVLTreePersistent) Insert(k int32, key V) *AVLTreePersistent {
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
	root := t._mergeWithRoot(a, _newANode(key), b)
	return _newWithRoot(root)
}

func (t *AVLTreePersistent) Split(k int32) (*AVLTreePersistent, *AVLTreePersistent) {
	a, b := t._splitNode(t.root, k)
	return _newWithRoot(a), _newWithRoot(b)
}

func (t *AVLTreePersistent) Pop(k int32) (*AVLTreePersistent, V) {
	if k < 0 {
		k += t.Size()
	}
	a, b := t._splitNode(t.root, k+1)
	a, tmp := t._popRight(a)
	root := t._mergeNode(a, b)
	return _newWithRoot(root), tmp.key
}

func (t *AVLTreePersistent) Set(k int32, key V) *AVLTreePersistent {
	if k < 0 {
		k += t.Size()
	}
	newRoot := t.root.Copy()
	node := newRoot
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
		t._update(path[i])
	}
	return _newWithRoot(newRoot)
}

func (t *AVLTreePersistent) ToList() []V {
	node := t.root
	stack := make([]*aNode, 0)
	res := make([]V, 0, t.Size())
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

func (t *AVLTreePersistent) Get(k int32) V {
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

func (t *AVLTreePersistent) Size() int32 {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func _newWithRoot(root *aNode) *AVLTreePersistent {
	return &AVLTreePersistent{root: root}
}

func (t *AVLTreePersistent) _build(n int32, f func(int32) V) {
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

func (t *AVLTreePersistent) _update(node *aNode) {
	node.size = 1
	node.height = 1
	if node.left != nil {
		node.size += node.left.size
		node.height = max8(node.left.height+1, 1)
	}
	if node.right != nil {
		node.size += node.right.size
		node.height = max8(node.height, node.right.height+1)
	}
}

func (t *AVLTreePersistent) _rotateRight(node *aNode) *aNode {
	u := node.left.Copy()
	node.left = u.right
	u.right = node
	t._update(node)
	t._update(u)
	return u
}

func (t *AVLTreePersistent) _rotateLeft(node *aNode) *aNode {
	u := node.right.Copy()
	node.right = u.left
	u.left = node
	t._update(node)
	t._update(u)
	return u
}

func (t *AVLTreePersistent) _balanceLeft(node *aNode) *aNode {
	node.right = node.right.Copy()
	u := node.right
	if u.Balance() == 1 {
		node.right = t._rotateRight(u)
	}
	u = t._rotateLeft(node)
	return u
}

func (t *AVLTreePersistent) _balanceRight(node *aNode) *aNode {
	node.left = node.left.Copy()
	u := node.left
	if u.Balance() == -1 {
		node.left = t._rotateLeft(u)
	}
	u = t._rotateRight(node)
	return u
}

func (t *AVLTreePersistent) _mergeWithRoot(l, root, r *aNode) *aNode {
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
		l = l.Copy()
		l.right = t._mergeWithRoot(l.right, root, r)
		t._update(l)
		if l.Balance() == -2 {
			return t._balanceLeft(l)
		}
		return l
	}
	if diff < -1 {
		r = r.Copy()
		r.left = t._mergeWithRoot(l, root, r.left)
		t._update(r)
		if r.Balance() == 2 {
			return t._balanceRight(r)
		}
		return r
	}
	root = root.Copy()
	root.left = l
	root.right = r
	t._update(root)
	return root
}

func (t *AVLTreePersistent) _mergeNode(l, r *aNode) *aNode {
	if l == nil && r == nil {
		return nil
	}
	if l == nil {
		return r.Copy()
	}
	if r == nil {
		return l.Copy()
	}
	l = l.Copy()
	r = r.Copy()
	l, root := t._popRight(l)
	return t._mergeWithRoot(l, root, r)
}

func (t *AVLTreePersistent) _popRight(node *aNode) (*aNode, *aNode) {
	tmpPath = tmpPath[:0]
	path := tmpPath
	node = node.Copy()
	mx := node
	for node.right != nil {
		path = append(path, node)
		node = node.right.Copy()
		mx = node
	}
	if node.left != nil {
		path = append(path, node.left.Copy())
	} else {
		path = append(path, nil)
	}
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
			path[len(path)-1].right = t._balanceRight(node)
		} else if b == -2 {
			path[len(path)-1].right = t._balanceLeft(node)
		} else {
			path[len(path)-1].right = node
		}
		t._update(path[len(path)-1])
	}
	if path[0] != nil {
		b := path[0].Balance()
		if b == 2 {
			path[0] = t._balanceRight(path[0])
		} else if b == -2 {
			path[0] = t._balanceLeft(path[0])
		}
	}
	mx.left = nil
	t._update(mx)
	return path[0], mx
}

func (t *AVLTreePersistent) _splitNode(node *aNode, k int32) (*aNode, *aNode) {
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
func test() {
	arr := []int{1, 2, 3, 4, 5}
	tree := NewLazyAVLTreePersistent(int32(len(arr)), func(i int32) V { return V(arr[i]) })
	a, b := tree.Split(1)
	fmt.Println(a.ToList(), b.ToList())
	a = a.Set(0, 100)
	fmt.Println(a.ToList(), b.ToList())

	for i := 0; i < 1000; i++ {
		n := rand.Intn(1000) + 500
		arr = make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100) + 1
		}
		tree = NewLazyAVLTreePersistent(int32(len(arr)), func(i int32) V { return V(arr[i]) })

		for j := 0; j < 500; j++ {
			// set
			{
				index, value := rand.Intn(len(arr)), rand.Intn(100)+1
				arr[index] = value
				tree = tree.Set(int32(index), V(value))
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
				tree = tree.Insert(int32(index), V(value))
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

	tree := NewLazyAVLTreePersistent(int32(len(arr)), func(i int32) V { return V(arr[i]) })

	time1 := time.Now()
	for j := 0; j < int(2e5); j++ {
		tree = tree.Set(int32(j), V(j))
		tree.Get(int32(j))
		tree = tree.Insert(int32(j), V(j))
		tree, _ = tree.Pop(int32(j))
		a, b := tree.Split(int32(j))
		a = a.Merge(b)
		tree = a
	}
	fmt.Println(time.Since(time1)) // 819.008833ms
}

func testCopyWithin() {
	for i := 0; i < 100; i++ {
		n := rand.Intn(1000) + 500
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100) + 1
		}
		tree := NewLazyAVLTreePersistent(int32(len(arr)), func(i int32) V { return V(arr[i]) })

		copyWithin := func(target, start, end int32) {
			copy(arr[target:], arr[start:end])
		}

		for j := 0; j < 500; j++ {
			target, start, end := int32(rand.Intn(n)), int32(rand.Intn(n)), int32(rand.Intn(n))
			if start > end {
				start, end = end, start
			}
			if target+end-start > int32(len(arr)) {
				continue
			}
			copyWithin(target, start, end)
			tree = CopyWithin(tree, target, start, end)
			if len(arr) != int(tree.Size()) {
				fmt.Println(arr, tree.ToList())
				panic("error size")
			}
			for k := 0; k < len(arr); k++ {
				if arr[k] != int(tree.Get(int32(k))) {
					fmt.Println(arr, tree.ToList())
					panic("error get")
				}
			}
		}
	}

	fmt.Println("Passed")
}
