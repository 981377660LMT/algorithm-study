// !Treap 的速度不如 RBST

// 暗号文 S が与えられます。この暗号文は、以下の操作で解読することが出来ます。

// T を空文字列とする。
// i=1,2,…,∣S∣ について、順番に以下を行う。 (∣S∣ は S の長さを表す)
// S の i 文字目が R のとき、T を反転させる。
// S の i 文字目が R でないとき、その文字を T の末尾に加える。
// この操作の後、T の中に同じ文字が 2 つ連続で並んでいたら、その 2 文字を取り除く。この操作を出来る限り続ける。 (最終的に得られる文字列は取り除く順番によらないことが証明できる)
// この操作で得られる文字列 T を出力してください。

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	const INF int = 1e18

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	treap := NewFHQTreap(Operations{
		elementMonoid: func() Element {
			return INF
		},
		dataMonoid: func(element Element) Data { // 每个位置处最小值
			return element
		},
		lazyMonoid: func() Lazy {
			return 0
		},
		op: func(a, b Data, e Element) Data {
			return min(min(a, b), e)
		},
		mappingElement: func(lazy Lazy, element Element) Element {
			return element + lazy
		},
		mappingData: func(lazy Lazy, data Data) Data {
			return data + lazy
		},
		composition: func(lazy1 Lazy, lazy2 Lazy) Lazy {
			return lazy1 + lazy2
		},
	}, []int{}, len(s))

	for _, c := range s {
		if c == 'R' {
			treap.Reverse(0, treap.Size())
		} else {
			treap.Append(int(c))
		}
	}

	stack := []int{}
	for i := 0; i < treap.Size(); i++ {
		num := treap.At(i)
		if len(stack) > 0 && stack[len(stack)-1] == num {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, num)
		}
	}

	res := make([]byte, len(stack))
	for i, num := range stack {
		res[i] = byte(num)
	}

	fmt.Fprintln(out, string(res))
}

// TODO
// Add a new node and return its nodeId.
func (t *FHQTreap) newNode(ele Element) int {
	node := Node{
		size:     1,
		element:  ele,
		data:     t.dataMonoid(ele),
		lazy:     t.lazyMonoid(),
		priority: t.nextRand(),
	}
	t.nodes = append(t.nodes, node)
	return len(t.nodes) - 1
}

// !Type and functions to be implemented.
// type Element = interface{}
// type Data = interface{}
// type Lazy = interface{}
type Element = int
type Data = int
type Lazy = int
type Operations struct {
	elementMonoid  func() Element
	dataMonoid     func(element Element) Data
	lazyMonoid     func() Lazy
	op             func(leftData, rightData Data, element Element) Data
	mappingData    func(lazy Lazy, data Data) Data
	mappingElement func(lazy Lazy, element Element) Element
	composition    func(parentLazy Lazy, childLazy Lazy) Lazy
}

// !Template
//
func NewFHQTreap(operations Operations, nums []Element, capacity int) *FHQTreap {
	treap := &FHQTreap{
		seed:       uint64(time.Now().UnixNano()/2 + 1),
		nodes:      make([]Node, 0, max(capacity, 16)),
		Operations: operations,
	}

	// 0 dummy
	dummy := &Node{
		size:     0,
		data:     operations.dataMonoid(operations.elementMonoid()),
		lazy:     operations.lazyMonoid(),
		element:  operations.elementMonoid(),
		priority: treap.nextRand(),
	}
	treap.nodes = append(treap.nodes, *dummy)
	if len(nums) > 0 {
		treap.root = treap.build(nums)
	}
	return treap
}

type FHQTreap struct {
	seed  uint64
	root  int
	nodes []Node

	// Segment-tree like operations
	Operations
}

type Node struct {
	// !Raw value
	element Element

	// !Data and lazy tag maintained by segment tree
	data Data
	lazy Lazy

	// FHQTreap inner attributes
	left, right int
	size        int
	priority    uint64
	isReversed  bool
}

// !op
func (t *FHQTreap) pushUp(root int) {
	if root == 0 {
		return
	}

	rootRef := &t.nodes[root]
	// If left or right is 0(dummy), it will update with monoid.
	rootRef.data = t.op(t.nodes[rootRef.left].data, t.nodes[rootRef.right].data, rootRef.element)
	rootRef.size = t.nodes[rootRef.left].size + t.nodes[rootRef.right].size + 1
}

// !Reverse first and then push down the lazy tag.
func (t *FHQTreap) pushDown(root int) {
	if t.nodes[root].isReversed {
		if t.nodes[root].left != 0 {
			t.toggle(t.nodes[root].left)
		}
		if t.nodes[root].right != 0 {
			t.toggle(t.nodes[root].right)
		}
		t.nodes[root].isReversed = false
	}

	if t.nodes[root].lazy != t.lazyMonoid() {
		if t.nodes[root].left != 0 {
			t.propagate(t.nodes[root].left, t.nodes[root].lazy)
		}
		if t.nodes[root].right != 0 {
			t.propagate(t.nodes[root].right, t.nodes[root].lazy)
		}
		t.nodes[root].lazy = t.lazyMonoid()
	}
}

// !mapping + composition
func (t *FHQTreap) propagate(root int, lazy Lazy) {
	t.nodes[root].element = t.mappingElement(lazy, t.nodes[root].element)
	t.nodes[root].data = t.mappingData(lazy, t.nodes[root].data)
	t.nodes[root].lazy = t.composition(lazy, t.nodes[root].lazy)
}

// Return the element at the k-th position (0-indexed).
func (t *FHQTreap) At(index int) Element {
	n := t.Size()
	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		panic(fmt.Sprintf("index %d out of range [0,%d]", index, n-1))
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index, &y, &z)
	t.splitByRank(y, index-1, &x, &y)
	res := &t.nodes[y].element
	t.root = t.merge(t.merge(x, y), z)
	return *res
}

// Set the k-th position (0-indexed) to the given value.
func (t *FHQTreap) Set(index int, element Element) {
	n := t.Size()
	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		panic(fmt.Sprintf("index %d out of range [0,%d]", index, n-1))
	}

	index++ // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index, &y, &z)
	t.splitByRank(y, index-1, &x, &y)
	t.nodes[y].element = element
	t.root = t.merge(t.merge(x, y), z)
}

// Reverse [start, stop) in place.
func (t *FHQTreap) Reverse(start, stop int) {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start-1, &x, &y)
	t.toggle(y)
	t.root = t.merge(x, t.merge(y, z))
}

// Rotate [start, stop) to the right step times.
func (t *FHQTreap) RotateRight(start, stop, step int) {
	start++
	n := stop - start + 1 - step%(stop-start+1)
	var x, y, z, p int
	t.splitByRank(t.root, start-1, &x, &y)
	t.splitByRank(y, n, &y, &z)
	t.splitByRank(z, stop-start+1-n, &z, &p)
	t.root = t.merge(t.merge(t.merge(x, z), y), p)
}

// Append element to the end of the list.
func (t *FHQTreap) Append(element Element) {
	t.Insert(t.Size(), element)
}

// Insert element before index.
func (t *FHQTreap) Insert(index int, element Element) {
	n := t.Size()
	if index < 0 {
		index += n
	}

	var x, y int
	t.splitByRank(t.root, index, &x, &y)
	t.root = t.merge(t.merge(x, t.newNode(element)), y)
}

// Remove and return element at index.
func (t *FHQTreap) Pop(index int) (element Element) {
	n := t.Size()
	if index < 0 {
		index += n
	}

	index++
	var x, y, z int
	t.splitByRank(t.root, index, &y, &z)
	t.splitByRank(y, index-1, &x, &y)
	element = t.nodes[y].element
	t.root = t.merge(x, z)
	return
}

// Remove [start, stop) from list.
func (t *FHQTreap) Erase(start, stop int) {
	var x, y, z int
	start++
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	t.root = t.merge(x, z)
}

// Update [start, stop) with value .
//  0 <= start <= stop <= n
func (t *FHQTreap) Update(start, stop int, lazy Lazy) {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start-1, &x, &y)
	t.propagate(y, lazy)
	t.root = t.merge(x, t.merge(y, z))
}

// Query data in [start, stop).
//  0 <= start <= stop <= n
func (t *FHQTreap) Query(start, stop int) Data {
	start++
	var x, y, z int
	t.splitByRank(t.root, stop, &x, &z)
	t.splitByRank(x, start-1, &x, &y)
	res := t.nodes[y].data
	t.root = t.merge(x, t.merge(y, z))
	return res
}

// Query data in [0, n).
func (t *FHQTreap) QueryAll() Data {
	return t.nodes[t.root].data
}

// Return the number of items in the list.
func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

// Return all elements in index order.
func (t *FHQTreap) InOrder() []Element {
	res := make([]Element, 0, t.Size())
	t.inOrder(t.root, &res)
	return res
}

func (t *FHQTreap) inOrder(root int, res *[]Element) {
	if root == 0 {
		return
	}
	t.pushDown(root) // !pushDown lazy tag
	t.inOrder(t.nodes[root].left, res)
	*res = append(*res, t.nodes[root].element)
	t.inOrder(t.nodes[root].right, res)
}

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (t *FHQTreap) splitByRank(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	t.pushDown(root)
	leftSize := t.nodes[t.nodes[root].left].size
	if leftSize+1 <= k {
		*x = root
		t.splitByRank(t.nodes[root].right, k-leftSize-1, &t.nodes[root].right, y)
		t.pushUp(*x)
	} else {
		*y = root
		t.splitByRank(t.nodes[root].left, k, x, &t.nodes[root].left)
		t.pushUp(*y)
	}
}

// Make sure that the height of the resulting tree is at most O(log n).
// A random priority is introduced to determine who is the root after merge operation.
// If left subtree is smaller, merge right subtree with the right child of the left subtree.
// Otherwise, merge left subtree with the left child of the right subtree.
func (t *FHQTreap) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}

	// !小根堆
	if t.nodes[x].priority < t.nodes[y].priority {
		t.pushDown(x)
		t.nodes[x].right = t.merge(t.nodes[x].right, y)
		t.pushUp(x)
		return x
	} else {
		t.pushDown(y)
		t.nodes[y].left = t.merge(x, t.nodes[y].left)
		t.pushUp(y)
		return y
	}
}

// 返回build后的根节点版本编号
func (t *FHQTreap) build(nums []int) int {
	n := len(nums)
	nodes := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nodes = append(nodes, t.newNode(nums[i]))
	}

	stack := []int{}
	pre := make([]int, n)
	for i := 0; i < n; i++ {
		pre[i] = -1
	}

	for i := 0; i < n; i++ {
		last := -1
		for len(stack) > 0 && t.nodes[stack[len(stack)-1]].priority > t.nodes[nodes[i]].priority {
			last = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}

		if len(stack) > 0 {
			pre[i] = stack[len(stack)-1]
		}
		if last != -1 {
			pre[last] = i
		}

		stack = append(stack, i)
	}

	root := -1
	for i := 0; i < n; i++ {
		if pre[i] != -1 {
			if i < pre[i] {
				t.nodes[nodes[pre[i]]].left = nodes[i]
			} else {
				t.nodes[nodes[pre[i]]].right = nodes[i]
			}
		} else {
			root = i
		}
	}

	t.root = nodes[root]
	t.buildHelper(nodes[root])
	return nodes[root]
}

func (t *FHQTreap) buildHelper(root int) {
	nodeRef := t.nodes[root]
	if nodeRef.left != 0 {
		t.buildHelper(nodeRef.left)
	}
	if nodeRef.right != 0 {
		t.buildHelper(nodeRef.right)
	}
	t.pushUp(root)
}

func (t *FHQTreap) String() string {
	sb := []string{"rbstarray{"}
	values := []string{}
	for i := 0; i < t.Size(); i++ {
		values = append(values, fmt.Sprintf("%d", t.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")
}

func (t *FHQTreap) toggle(root int) {
	t.nodes[root].left, t.nodes[root].right = t.nodes[root].right, t.nodes[root].left
	t.nodes[root].isReversed = !t.nodes[root].isReversed
}

func (t *FHQTreap) nextRand() uint64 {
	t.seed ^= t.seed << 7
	t.seed ^= t.seed >> 9
	return t.seed & 0xFFFFFFFF
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
