// Copyright 2020 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
//
// !类似python里的SortedSet. 有序集合，不允许重复元素.
//
// - Basic
// Set(item)               // insert or replace an item
// Get(item)               // get an existing item
// Delete(item)            // delete an item
// Len()                   // return the number of items in the btree
//
// - Iteration
// Scan(iter)              // scan items in ascending order
// Reverse(iter)           // scan items in descending order
// Ascend(key, iter)       // scan items in ascending order that are >= to key
// Descend(key, iter)      // scan items in descending order that are <= to key.
// Iter()                  // returns a read-only iterator for for-loops.
//
// - Array-like operations
// GetAt(index)            // returns the item at index
// DeleteAt(index)         // deletes the item at index
// PopMin()                // removes and returns the smallest item
// PopMax()                // removes and returns the largest item
//
// - Bulk-loading
// Load(item)              // load presorted items into tree
//
// - Path hinting
// SetHint(item, *hint)    // insert or replace an existing item
// GetHint(item, *hint)    // get an existing item
// DeleteHint(item, *hint) // delete an item
// AscendHint(key, iter, *hint)
// DescendHint(key, iter, *hint)
// SeekHint(key, iter, *hint)
//
// - Copy-on-write
// Copy()                  // copy the btree
//
// -Iter
// Iter/IterMut
// Seek/SeekHint
// First
// Last
// Next
// Prev
// Item
// Release

package main

import "fmt"

func main() {
	pathHint()
}

func pathHint() {
	// 创建一个 B-树，键类型为 int
	less := func(a, b int) bool { return a < b }
	tree := NewBTreeG(less)

	// 初始化 PathHint
	var hint PathHint

	// 插入元素，并使用 PathHint
	keys := []int{10, 20, 30, 25, 15, 5, 35}
	for _, key := range keys {
		tree.SetHint(key, &hint)
	}

	// 搜索元素，并使用 PathHint
	searchKeys := []int{15, 25, 35, 40}
	for _, key := range searchKeys {
		value, found := tree.GetHint(key, &hint)
		if found {
			fmt.Printf("Found key %d with value %d\n", key, value)
		} else {
			fmt.Printf("Key %d not found\n", key)
		}
	}

	// 使用迭代器进行遍历，并利用 PathHint
	iter := tree.Iter()
	defer iter.Release()
	if iter.SeekHint(15, &hint) {
		fmt.Println("Iterator found:", iter.Item())
		for iter.Next() {
			fmt.Println("Iterator next:", iter.Item())
		}
	}
}

var gisoid uint64

func newIsoID() uint64 {
	gisoid++
	return gisoid
}

type copier[T any] interface {
	Copy() T
}

type isoCopier[T any] interface {
	IsoCopy() T // iso: isolate
}

func degreeToMinMax(deg int) (min, max int) {
	if deg <= 0 {
		deg = 32
	} else if deg == 1 {
		deg = 2 // must have at least 2
	}
	max = deg*2 - 1 // max items per node. max children is +1
	min = max / 2
	return min, max
}

type BTreeG[T any] struct {
	isoid uint64 // 树的“隔离ID”，用于COW判断

	root  *node[T]
	count int

	copyItems    bool // T是否具有Copy方法
	isoCopyItems bool // T是否具有IsoCopy方法

	less  func(a, b T) bool
	empty T

	max int // 每个节点最大元素数
	min int
}

type node[T any] struct {
	isoid    uint64
	count    int
	items    []T
	children *[]*node[T]
}

// PathHint 是一个用于 *Hint() 函数的工具类型。Hints 为聚集键提供更快的操作。
// 最多记录 8 层深度（对绝大部分实际 B-Tree 已足够）。
// 在搜索时可优先使用 hint 里的 index，而不是做二分。
// 参考 hintsearch() 函数。
type PathHint struct {
	used [8]bool  // 每一层是否已经使用了路径提示
	path [8]uint8 // 存储每一层的索引位置，用于优化搜索路径
}

func NewBTreeG[T any](less func(a, b T) bool) *BTreeG[T] {
	return NewBTreeGWithDegreee(less, 32)
}

// Degree 用于定义每个内部节点在分支之前可以包含多少个元素和子节点。
// 例如，Degree 为2将创建一个2-3-4树，每个节点可以包含1-3个元素和2-4个子节点。参见 https://en.wikipedia.org/wiki/2–3–4_tree。
// 默认值为32.
func NewBTreeGWithDegreee[T any](less func(a, b T) bool, degree int) *BTreeG[T] {
	tr := new(BTreeG[T])
	tr.isoid = newIsoID()
	tr.less = less
	tr.init(degree)
	return tr
}

func (tr *BTreeG[T]) init(degree int) {
	if tr.min != 0 {
		return
	}
	tr.min, tr.max = degreeToMinMax(degree)
	_, tr.copyItems = ((interface{})(tr.empty)).(copier[T])
	if !tr.copyItems {
		_, tr.isoCopyItems = ((interface{})(tr.empty)).(isoCopier[T])
	}
}

func (tr *BTreeG[T]) Less(a, b T) bool {
	return tr.less(a, b)
}

func (tr *BTreeG[T]) newNode(leaf bool) *node[T] {
	n := &node[T]{isoid: tr.isoid}
	if !leaf {
		n.children = new([]*node[T])
	}
	return n
}

func (n *node[T]) leaf() bool {
	return n.children == nil
}

func (tr *BTreeG[T]) bsearch(n *node[T], key T) (index int, found bool) {
	low, high := 0, len(n.items)
	for low < high {
		h := (low + high) / 2
		if !tr.less(key, n.items[h]) {
			low = h + 1
		} else {
			high = h
		}
	}
	if low > 0 && !tr.less(n.items[low-1], key) {
		return low - 1, true
	}
	return low, false
}

func (tr *BTreeG[T]) find(n *node[T], key T, hint *PathHint, depth int) (index int, found bool) {
	if hint == nil {
		return tr.bsearch(n, key)
	}
	return tr.hintsearch(n, key, hint, depth)
}

func (tr *BTreeG[T]) hintsearch(n *node[T], key T, hint *PathHint, depth int) (index int, found bool) {
	// 最佳情况找到精确匹配，更新hint并返回。
	// 最坏情况，更新low和high边界进行二分查找。
	low := 0
	high := len(n.items) - 1
	if depth < 8 && hint.used[depth] {
		index = int(hint.path[depth])
		if index >= len(n.items) {
			// tail item
			if tr.Less(n.items[len(n.items)-1], key) {
				index = len(n.items)
				goto path_match
			}
			index = len(n.items) - 1
		}
		if tr.Less(key, n.items[index]) {
			if index == 0 || tr.Less(n.items[index-1], key) {
				goto path_match
			}
			high = index - 1
		} else if tr.Less(n.items[index], key) {
			low = index + 1
		} else {
			found = true
			goto path_match
		}
	}

	for low <= high {
		mid := low + ((high+1)-low)/2
		if !tr.Less(key, n.items[mid]) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	if low > 0 && !tr.Less(n.items[low-1], key) {
		index = low - 1
		found = true
	} else {
		index = low
		found = false
	}

path_match:
	if depth < 8 {
		hint.used[depth] = true
		var pathIndex uint8
		if n.leaf() && found {
			pathIndex = uint8(index + 1)
		} else {
			pathIndex = uint8(index)
		}
		if pathIndex != hint.path[depth] {
			hint.path[depth] = pathIndex
			for i := depth + 1; i < 8; i++ {
				hint.used[i] = false
			}
		}
	}
	return index, found
}

func (tr *BTreeG[T]) SetHint(item T, hint *PathHint) (prev T, replaced bool) {
	if tr.root == nil {
		tr.init(0)
		tr.root = tr.newNode(true)
		tr.root.items = append([]T{}, item)
		tr.root.count = 1
		tr.count = 1
		return tr.empty, false
	}
	// 如果节点分裂，调整树的结构，创建新的根节点
	prev, replaced, split := tr.nodeSet(&tr.root, item, hint, 0)
	if split {
		left := tr.isoLoad(&tr.root, true)
		right, median := tr.nodeSplit(left)
		tr.root = tr.newNode(false)
		*tr.root.children = make([]*node[T], 0, tr.max+1)
		*tr.root.children = append([]*node[T]{}, left, right)
		tr.root.items = append([]T{}, median)
		tr.root.updateCount()
		return tr.SetHint(item, hint)
	}
	if replaced {
		return prev, true
	}
	tr.count++
	return tr.empty, false
}

// Set or replace a value for a key
func (tr *BTreeG[T]) Set(item T) (T, bool) {
	return tr.SetHint(item, nil)
}

func (tr *BTreeG[T]) nodeSplit(n *node[T]) (right *node[T], median T) {
	i := tr.max / 2
	median = n.items[i]

	// right node
	right = tr.newNode(n.leaf())
	right.items = n.items[i+1:]
	if !n.leaf() {
		*right.children = (*n.children)[i+1:]
	}
	right.updateCount()

	// left node
	n.items[i] = tr.empty
	n.items = n.items[:i:i]
	if !n.leaf() {
		*n.children = (*n.children)[: i+1 : i+1]
	}
	n.updateCount()

	return right, median
}

func (n *node[T]) updateCount() {
	n.count = len(n.items)
	if !n.leaf() {
		for i := 0; i < len(*n.children); i++ {
			n.count += (*n.children)[i].count
		}
	}
}

// Copy the node for safe isolation.
func (tr *BTreeG[T]) copy(n *node[T]) *node[T] {
	n2 := new(node[T])
	n2.isoid = tr.isoid
	n2.count = n.count
	n2.items = make([]T, len(n.items), cap(n.items))
	copy(n2.items, n.items)
	if tr.copyItems {
		for i := 0; i < len(n2.items); i++ {
			n2.items[i] = ((interface{})(n2.items[i])).(copier[T]).Copy()
		}
	} else if tr.isoCopyItems {
		for i := 0; i < len(n2.items); i++ {
			n2.items[i] = ((interface{})(n2.items[i])).(isoCopier[T]).IsoCopy()
		}
	}
	if !n.leaf() {
		n2.children = new([]*node[T])
		*n2.children = make([]*node[T], len(*n.children), tr.max+1)
		copy(*n2.children, *n.children)
	}
	return n2
}

// isoLoad loads the provided node and, if needed, performs a copy-on-write.
func (tr *BTreeG[T]) isoLoad(cn **node[T], mut bool) *node[T] {
	if mut && (*cn).isoid != tr.isoid {
		*cn = tr.copy(*cn)
	}
	return *cn
}

// 递归地在节点中插入元素，处理节点的分裂
func (tr *BTreeG[T]) nodeSet(cn **node[T], item T,
	hint *PathHint, depth int,
) (prev T, replaced bool, split bool) {
	if (*cn).isoid != tr.isoid {
		*cn = tr.copy(*cn)
	}
	n := *cn
	var i int
	var found bool
	if hint == nil {
		i, found = tr.bsearch(n, item)
	} else {
		i, found = tr.hintsearch(n, item, hint, depth)
	}

	// 处理重复元素
	if found {
		prev = n.items[i]
		n.items[i] = item
		return prev, true, false
	}

	if n.leaf() {
		if len(n.items) == tr.max {
			return tr.empty, false, true
		}
		n.items = append(n.items, tr.empty)
		copy(n.items[i+1:], n.items[i:])
		n.items[i] = item
		n.count++
		return tr.empty, false, false
	}

	prev, replaced, split = tr.nodeSet(&(*n.children)[i], item, hint, depth+1)
	if split {
		if len(n.items) == tr.max {
			return tr.empty, false, true
		}
		right, median := tr.nodeSplit((*n.children)[i])
		*n.children = append(*n.children, nil)
		copy((*n.children)[i+1:], (*n.children)[i:])
		(*n.children)[i+1] = right
		n.items = append(n.items, tr.empty)
		copy(n.items[i+1:], n.items[i:])
		n.items[i] = median
		return tr.nodeSet(&n, item, hint, depth)
	}
	if !replaced {
		n.count++
	}
	return prev, replaced, false
}

func (tr *BTreeG[T]) Scan(iter func(item T) bool) {
	tr.scan(iter, false)
}
func (tr *BTreeG[T]) ScanMut(iter func(item T) bool) {
	tr.scan(iter, true)
}

func (tr *BTreeG[T]) scan(iter func(item T) bool, mut bool) {
	if tr.root == nil {
		return
	}
	tr.nodeScan(&tr.root, iter, mut)
}

func (tr *BTreeG[T]) nodeScan(cn **node[T], iter func(item T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	if n.leaf() {
		for i := 0; i < len(n.items); i++ {
			if !iter(n.items[i]) {
				return false
			}
		}
		return true
	}
	for i := 0; i < len(n.items); i++ {
		if !tr.nodeScan(&(*n.children)[i], iter, mut) {
			return false
		}
		if !iter(n.items[i]) {
			return false
		}
	}
	return tr.nodeScan(&(*n.children)[len(*n.children)-1], iter, mut)
}

// Get a value for key
func (tr *BTreeG[T]) Get(key T) (T, bool) {
	return tr.getHint(key, nil, false)
}

func (tr *BTreeG[T]) GetMut(key T) (T, bool) {
	return tr.getHint(key, nil, true)
}

// GetHint gets a value for key using a path hint
func (tr *BTreeG[T]) GetHint(key T, hint *PathHint) (value T, ok bool) {
	return tr.getHint(key, hint, false)
}
func (tr *BTreeG[T]) GetHintMut(key T, hint *PathHint) (value T, ok bool) {
	return tr.getHint(key, hint, true)
}

// GetHint gets a value for key using a path hint
func (tr *BTreeG[T]) getHint(key T, hint *PathHint, mut bool) (T, bool) {
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, mut)
	depth := 0
	for {
		i, found := tr.find(n, key, hint, depth)
		if found {
			return n.items[i], true
		}
		if n.children == nil {
			return tr.empty, false
		}
		n = tr.isoLoad(&(*n.children)[i], mut)
		depth++
	}
}

func (tr *BTreeG[T]) Len() int {
	return tr.count
}

func (tr *BTreeG[T]) Delete(key T) (T, bool) {
	return tr.DeleteHint(key, nil)
}

// 使用路径提示删除 key 的值并返回被删除的值。
// 如果未找到 key，则返回 false。
func (tr *BTreeG[T]) DeleteHint(key T, hint *PathHint) (T, bool) {
	if tr.root == nil {
		return tr.empty, false
	}
	prev, deleted := tr.delete(&tr.root, false, key, hint, 0)
	if !deleted {
		return tr.empty, false
	}
	if len(tr.root.items) == 0 && !tr.root.leaf() {
		tr.root = (*tr.root.children)[0]
	}
	tr.count--
	if tr.count == 0 {
		tr.root = nil
	}
	return prev, true
}

func (tr *BTreeG[T]) delete(cn **node[T], max bool, key T,
	hint *PathHint, depth int,
) (T, bool) {
	n := tr.isoLoad(cn, true)
	var i int
	var found bool
	if max {
		i, found = len(n.items)-1, true
	} else {
		i, found = tr.find(n, key, hint, depth)
	}
	if n.leaf() {
		if found {
			// found the items at the leaf, remove it and return.
			prev := n.items[i]
			copy(n.items[i:], n.items[i+1:])
			n.items[len(n.items)-1] = tr.empty
			n.items = n.items[:len(n.items)-1]
			n.count--
			return prev, true
		}
		return tr.empty, false
	}

	var prev T
	var deleted bool
	if found {
		if max {
			i++
			prev, deleted = tr.delete(&(*n.children)[i], true, tr.empty, nil, 0)
		} else {
			prev = n.items[i]
			maxItem, _ := tr.delete(&(*n.children)[i], true, tr.empty, nil, 0)
			deleted = true
			n.items[i] = maxItem
		}
	} else {
		prev, deleted = tr.delete(&(*n.children)[i], max, key, hint, depth+1)
	}
	if !deleted {
		return tr.empty, false
	}
	n.count--
	// 如果子节点删除后元素数少于 min，调用 nodeRebalance 进行重平衡
	if len((*n.children)[i].items) < tr.min {
		tr.nodeRebalance(n, i)
	}
	return prev, true
}

// nodeRebalance rebalances the child nodes following a delete operation.
// Provide the index of the child node with the number of items that fell
// below minItems.
func (tr *BTreeG[T]) nodeRebalance(n *node[T], i int) {
	if i == len(n.items) {
		i--
	}

	// ensure copy-on-write
	left := tr.isoLoad(&(*n.children)[i], true)
	right := tr.isoLoad(&(*n.children)[i+1], true)

	if len(left.items)+len(right.items) < tr.max {
		// Merges the left and right children nodes together as a single node
		// that includes (left,item,right), and places the contents into the
		// existing left node. Delete the right node altogether and move the
		// following items and child nodes to the left by one slot.

		// merge (left,item,right)
		left.items = append(left.items, n.items[i])
		left.items = append(left.items, right.items...)
		if !left.leaf() {
			*left.children = append(*left.children, *right.children...)
		}
		left.count += right.count + 1

		// move the items over one slot
		copy(n.items[i:], n.items[i+1:])
		n.items[len(n.items)-1] = tr.empty
		n.items = n.items[:len(n.items)-1]

		// move the children over one slot
		copy((*n.children)[i+1:], (*n.children)[i+2:])
		(*n.children)[len(*n.children)-1] = nil
		(*n.children) = (*n.children)[:len(*n.children)-1]
	} else if len(left.items) > len(right.items) {
		// move left -> right over one slot

		// Move the item of the parent node at index into the right-node first
		// slot, and move the left-node last item into the previously moved
		// parent item slot.
		right.items = append(right.items, tr.empty)
		copy(right.items[1:], right.items)
		right.items[0] = n.items[i]
		right.count++
		n.items[i] = left.items[len(left.items)-1]
		left.items[len(left.items)-1] = tr.empty
		left.items = left.items[:len(left.items)-1]
		left.count--

		if !left.leaf() {
			// 将左子节点的最后一个子节点移动到右子节点的第一个位置
			*right.children = append(*right.children, nil)
			copy((*right.children)[1:], *right.children)
			(*right.children)[0] = (*left.children)[len(*left.children)-1]
			(*left.children)[len(*left.children)-1] = nil
			(*left.children) = (*left.children)[:len(*left.children)-1]
			left.count -= (*right.children)[0].count
			right.count += (*right.children)[0].count
		}
	} else {
		// 将右子节点的第一个项目移动到左子节点的最后一个位置

		// Same as above but the other direction
		left.items = append(left.items, n.items[i])
		left.count++
		n.items[i] = right.items[0]
		copy(right.items, right.items[1:])
		right.items[len(right.items)-1] = tr.empty
		right.items = right.items[:len(right.items)-1]
		right.count--

		if !left.leaf() {
			*left.children = append(*left.children, (*right.children)[0])
			copy(*right.children, (*right.children)[1:])
			(*right.children)[len(*right.children)-1] = nil
			*right.children = (*right.children)[:len(*right.children)-1]
			left.count += (*left.children)[len(*left.children)-1].count
			right.count -= (*left.children)[len(*left.children)-1].count
		}
	}
}

// Ascend 在范围 [pivot, last] 内以升序遍历树
// 传入 nil 作为 pivot 则扫描所有项目，按升序排列
// 返回 false 来停止迭代
func (tr *BTreeG[T]) Ascend(pivot T, iter func(item T) bool) {
	tr.ascend(pivot, iter, false, nil)
}
func (tr *BTreeG[T]) AscendMut(pivot T, iter func(item T) bool) {
	tr.ascend(pivot, iter, true, nil)
}
func (tr *BTreeG[T]) ascend(pivot T, iter func(item T) bool, mut bool,
	hint *PathHint,
) {
	if tr.root == nil {
		return
	}
	tr.nodeAscend(&tr.root, pivot, hint, 0, iter, mut)
}
func (tr *BTreeG[T]) AscendHint(pivot T, iter func(item T) bool, hint *PathHint,
) {
	tr.ascend(pivot, iter, false, hint)
}
func (tr *BTreeG[T]) AscendHintMut(pivot T, iter func(item T) bool,
	hint *PathHint,
) {
	tr.ascend(pivot, iter, true, hint)
}

func (tr *BTreeG[T]) nodeAscend(cn **node[T], pivot T, hint *PathHint,
	depth int, iter func(item T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	i, found := tr.find(n, pivot, hint, depth)
	if !found {
		if !n.leaf() {
			if !tr.nodeAscend(&(*n.children)[i], pivot, hint, depth+1, iter,
				mut) {
				return false
			}
		}
	}
	// 我们处于以下情况之一：
	// - 找到了节点，应该从 `i` 开始迭代。
	// - 没找到节点，需要处理
	for ; i < len(n.items); i++ {
		if !iter(n.items[i]) {
			return false
		}
		if !n.leaf() {
			if !tr.nodeScan(&(*n.children)[i+1], iter, mut) {
				return false
			}
		}
	}
	return true
}

func (tr *BTreeG[T]) Reverse(iter func(item T) bool) {
	tr.reverse(iter, false)
}
func (tr *BTreeG[T]) ReverseMut(iter func(item T) bool) {
	tr.reverse(iter, true)
}
func (tr *BTreeG[T]) reverse(iter func(item T) bool, mut bool) {
	if tr.root == nil {
		return
	}
	tr.nodeReverse(&tr.root, iter, mut)
}

func (tr *BTreeG[T]) nodeReverse(cn **node[T], iter func(item T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	if n.leaf() {
		for i := len(n.items) - 1; i >= 0; i-- {
			if !iter(n.items[i]) {
				return false
			}
		}
		return true
	}
	if !tr.nodeReverse(&(*n.children)[len(*n.children)-1], iter, mut) {
		return false
	}
	for i := len(n.items) - 1; i >= 0; i-- {
		if !iter(n.items[i]) {
			return false
		}
		if !tr.nodeReverse(&(*n.children)[i], iter, mut) {
			return false
		}
	}
	return true
}

// 范围 [pivot, first] 内以降序遍历树
// 传入 nil 作为 pivot 则扫描所有项目，按降序排列
// 返回 false 来停止迭代
func (tr *BTreeG[T]) Descend(pivot T, iter func(item T) bool) {
	tr.descend(pivot, iter, false, nil)
}
func (tr *BTreeG[T]) DescendMut(pivot T, iter func(item T) bool) {
	tr.descend(pivot, iter, true, nil)
}
func (tr *BTreeG[T]) descend(pivot T, iter func(item T) bool, mut bool, hint *PathHint) {
	if tr.root == nil {
		return
	}
	tr.nodeDescend(&tr.root, pivot, hint, 0, iter, mut)
}

func (tr *BTreeG[T]) DescendHint(pivot T, iter func(item T) bool,
	hint *PathHint,
) {
	tr.descend(pivot, iter, false, hint)
}
func (tr *BTreeG[T]) DescendHintMut(pivot T, iter func(item T) bool,
	hint *PathHint,
) {
	tr.descend(pivot, iter, true, hint)
}

func (tr *BTreeG[T]) nodeDescend(cn **node[T], pivot T, hint *PathHint,
	depth int, iter func(item T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	i, found := tr.find(n, pivot, hint, depth)
	if !found {
		if !n.leaf() {
			if !tr.nodeDescend(&(*n.children)[i], pivot, hint, depth+1, iter,
				mut) {
				return false
			}
		}
		i--
	}
	for ; i >= 0; i-- {
		if !iter(n.items[i]) {
			return false
		}
		if !n.leaf() {
			if !tr.nodeReverse(&(*n.children)[i], iter, mut) {
				return false
			}
		}
	}
	return true
}

// 在已经预排序的情况下，将元素乐观插入 B-树的右侧叶子节点.
// 返回(旧值, 是否替换).
func (tr *BTreeG[T]) Load(item T) (T, bool) {
	if tr.root == nil {
		return tr.SetHint(item, nil)
	}
	n := tr.isoLoad(&tr.root, true)
	for {
		n.count++ // 乐观更新计数
		if n.leaf() {
			if len(n.items) < tr.max {
				if tr.Less(n.items[len(n.items)-1], item) {
					n.items = append(n.items, item)
					tr.count++
					return tr.empty, false
				}
			}
			break
		}
		n = tr.isoLoad(&(*n.children)[len(*n.children)-1], true)
	}
	// 回滚计数
	n = tr.root
	for {
		n.count--
		if n.leaf() {
			break
		}
		n = (*n.children)[len(*n.children)-1]
	}
	return tr.SetHint(item, nil)
}

func (tr *BTreeG[T]) Min() (T, bool) {
	return tr.minMut(false)
}

func (tr *BTreeG[T]) MinMut() (T, bool) {
	return tr.minMut(true)
}

func (tr *BTreeG[T]) minMut(mut bool) (T, bool) {
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, mut)
	for {
		if n.leaf() {
			return n.items[0], true
		}
		n = tr.isoLoad(&(*n.children)[0], mut)
	}
}

func (tr *BTreeG[T]) Max() (T, bool) {
	return tr.maxMut(false)
}

func (tr *BTreeG[T]) MaxMut() (T, bool) {
	return tr.maxMut(true)
}

func (tr *BTreeG[T]) maxMut(mut bool) (T, bool) {
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, mut)
	for {
		if n.leaf() {
			return n.items[len(n.items)-1], true
		}
		n = tr.isoLoad(&(*n.children)[len(*n.children)-1], mut)
	}
}

// 移除并返回树中的最小项目。
// 如果树没有项目，则返回 false。
func (tr *BTreeG[T]) PopMin() (T, bool) {
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, true)
	var item T
	for {
		n.count-- // optimistically update counts
		if n.leaf() {
			item = n.items[0]
			if len(n.items) == tr.min {
				break
			}
			copy(n.items[:], n.items[1:])
			n.items[len(n.items)-1] = tr.empty
			n.items = n.items[:len(n.items)-1]
			tr.count--
			if tr.count == 0 {
				tr.root = nil
			}
			return item, true
		}
		n = tr.isoLoad(&(*n.children)[0], true)
	}
	// revert the counts
	n = tr.root
	for {
		n.count++
		if n.leaf() {
			break
		}
		n = (*n.children)[0]
	}
	return tr.DeleteHint(item, nil)
}

// 移除并返回树中的最大项目。
// 如果树没有项目，则返回 false。
func (tr *BTreeG[T]) PopMax() (T, bool) {
	if tr.root == nil {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, true)
	var item T
	for {
		n.count-- // optimistically update counts
		if n.leaf() {
			item = n.items[len(n.items)-1]
			if len(n.items) == tr.min {
				break
			}
			n.items[len(n.items)-1] = tr.empty
			n.items = n.items[:len(n.items)-1]
			tr.count--
			if tr.count == 0 {
				tr.root = nil
			}
			return item, true
		}
		n = tr.isoLoad(&(*n.children)[len(*n.children)-1], true)
	}
	// revert the counts
	n = tr.root
	for {
		n.count++
		if n.leaf() {
			break
		}
		n = (*n.children)[len(*n.children)-1]
	}
	return tr.DeleteHint(item, nil)
}

// 返回指定索引的值。
// 如果树为空或索引超出范围，返回 false。
func (tr *BTreeG[T]) GetAt(index int) (T, bool) {
	return tr.getAt(index, false)
}
func (tr *BTreeG[T]) GetAtMut(index int) (T, bool) {
	return tr.getAt(index, true)
}
func (tr *BTreeG[T]) getAt(index int, mut bool) (T, bool) {
	if tr.root == nil || index < 0 || index >= tr.count {
		return tr.empty, false
	}
	n := tr.isoLoad(&tr.root, mut)
	for {
		if n.leaf() {
			return n.items[index], true
		}
		i := 0
		for ; i < len(n.items); i++ {
			if index < (*n.children)[i].count {
				break
			} else if index == (*n.children)[i].count {
				return n.items[i], true
			}
			index -= (*n.children)[i].count + 1
		}
		n = tr.isoLoad(&(*n.children)[i], mut)
	}
}

// 删除指定索引的项目。
// 如果树为空或索引超出范围，返回 false。
func (tr *BTreeG[T]) DeleteAt(index int) (T, bool) {
	if tr.root == nil || index < 0 || index >= tr.count {
		return tr.empty, false
	}
	var pathbuf [8]uint8 // track the path
	path := pathbuf[:0]
	var item T
	n := tr.isoLoad(&tr.root, true)
outer:
	for {
		n.count-- // optimistically update counts
		if n.leaf() {
			// the index is the item position
			item = n.items[index]
			if len(n.items) == tr.min {
				path = append(path, uint8(index))
				break outer
			}
			copy(n.items[index:], n.items[index+1:])
			n.items[len(n.items)-1] = tr.empty
			n.items = n.items[:len(n.items)-1]
			tr.count--
			if tr.count == 0 {
				tr.root = nil
			}
			return item, true
		}
		i := 0
		for ; i < len(n.items); i++ {
			if index < (*n.children)[i].count {
				break
			} else if index == (*n.children)[i].count {
				item = n.items[i]
				path = append(path, uint8(i))
				break outer
			}
			index -= (*n.children)[i].count + 1
		}
		path = append(path, uint8(i))
		n = tr.isoLoad(&(*n.children)[i], true)
	}
	// revert the counts
	var hint PathHint
	n = tr.root
	for i := 0; i < len(path); i++ {
		if i < len(hint.path) {
			hint.path[i] = uint8(path[i])
			hint.used[i] = true
		}
		n.count++
		if !n.leaf() {
			n = (*n.children)[uint8(path[i])]
		}
	}
	return tr.DeleteHint(item, &hint)
}

// 返回树的高度。
// 如果树没有项目，返回0。
func (tr *BTreeG[T]) Height() int {
	var height int
	if tr.root != nil {
		n := tr.root
		for {
			height++
			if n.leaf() {
				break
			}
			n = (*n.children)[0]
		}
	}
	return height
}

// 按顺序迭代树中的所有项目。
// items 参数将包含一个或多个项目。
func (tr *BTreeG[T]) Walk(iter func(item []T) bool) {
	tr.walk(iter, false)
}
func (tr *BTreeG[T]) WalkMut(iter func(item []T) bool) {
	tr.walk(iter, true)
}
func (tr *BTreeG[T]) walk(iter func(item []T) bool, mut bool) {
	if tr.root == nil {
		return
	}
	tr.nodeWalk(&tr.root, iter, mut)
}

func (tr *BTreeG[T]) nodeWalk(cn **node[T], iter func(item []T) bool, mut bool,
) bool {
	n := tr.isoLoad(cn, mut)
	if n.leaf() {
		if !iter(n.items) {
			return false
		}
	} else {
		for i := 0; i < len(n.items); i++ {
			if !tr.nodeWalk(&(*n.children)[i], iter, mut) {
				return false
			}
			if !iter(n.items[i : i+1]) {
				return false
			}
		}
		if !tr.nodeWalk(&(*n.children)[len(n.items)], iter, mut) {
			return false
		}
	}
	return true
}

// Copy 复制树。这是一个 copy-on-write 操作，非常快速，因为它只执行影子复制。
func (tr *BTreeG[T]) Copy() *BTreeG[T] {
	return tr.IsoCopy()
}

func (tr *BTreeG[T]) IsoCopy() *BTreeG[T] {
	tr.isoid = newIsoID()
	tr2 := new(BTreeG[T])
	*tr2 = *tr
	tr2.isoid = newIsoID()
	return tr2
}

// B树迭代器.
type IterG[T any] struct {
	tr  *BTreeG[T]
	mut bool // 迭代器是否需要进行修改（可变迭代器）

	seeked  bool // 是否已经定位到一个项目
	atstart bool // 是否位于树的起始位置之前
	atend   bool // 是否位于末尾位置之后

	stack0 [4]iterStackItemG[T]
	stack  []iterStackItemG[T] // 保存从根节点到当前节点的路径，支持 Next 和 Prev 操作

	item T // 当前迭代器指向的项目
}

type iterStackItemG[T any] struct {
	n *node[T] // 当前节点
	i int      // 当前节点中的索引
}

// 返回一个只读的迭代器。
// 必须在完成迭代器后调用 Release 方法。
func (tr *BTreeG[T]) Iter() IterG[T] {
	return tr.iter(false)
}

// 返回一个可变的迭代器，允许在遍历过程中修改树（例如，边遍历边删除符合条件的元素）
func (tr *BTreeG[T]) IterMut() IterG[T] {
	return tr.iter(true)
}

func (tr *BTreeG[T]) iter(mut bool) IterG[T] {
	var iter IterG[T]
	iter.tr = tr
	iter.mut = mut
	iter.stack = iter.stack0[:0]
	return iter
}

// 定位到大于或等于 key 的项目。
// 如果没有找到项目，返回 false。
func (iter *IterG[T]) Seek(key T) bool {
	return iter.seek(key, nil)
}

func (iter *IterG[T]) SeekHint(key T, hint *PathHint) bool {
	return iter.seek(key, hint)
}

func (iter *IterG[T]) seek(key T, hint *PathHint) bool {
	if iter.tr == nil {
		return false
	}
	iter.seeked = true
	iter.stack = iter.stack[:0]
	if iter.tr.root == nil {
		return false
	}
	n := iter.tr.isoLoad(&iter.tr.root, iter.mut)
	var depth int
	for {
		i, found := iter.tr.find(n, key, hint, depth)
		iter.stack = append(iter.stack, iterStackItemG[T]{n, i})
		if found {
			iter.item = n.items[i]
			return true
		}
		if n.leaf() {
			iter.stack[len(iter.stack)-1].i--
			return iter.Next()
		}
		n = iter.tr.isoLoad(&(*n.children)[i], iter.mut)
		depth++
	}
}

// 将迭代器移动到树中的第一个项目。
// 如果树为空，返回 false。
func (iter *IterG[T]) First() bool {
	if iter.tr == nil {
		return false
	}
	iter.atend = false
	iter.atstart = false // 一旦迭代器移动到第一个元素，atstart 被设置为 false
	iter.seeked = true
	iter.stack = iter.stack[:0]
	if iter.tr.root == nil {
		return false
	}
	n := iter.tr.isoLoad(&iter.tr.root, iter.mut)
	for {
		iter.stack = append(iter.stack, iterStackItemG[T]{n, 0})
		if n.leaf() {
			break
		}
		n = iter.tr.isoLoad(&(*n.children)[0], iter.mut)
	}
	s := &iter.stack[len(iter.stack)-1]
	iter.item = s.n.items[s.i]
	return true
}

// Last 将迭代器移动到树中的最后一个项目。
// 如果树为空，返回 false。
func (iter *IterG[T]) Last() bool {
	if iter.tr == nil {
		return false
	}
	iter.seeked = true
	iter.stack = iter.stack[:0]
	if iter.tr.root == nil {
		return false
	}
	n := iter.tr.isoLoad(&iter.tr.root, iter.mut)
	for {
		iter.stack = append(iter.stack, iterStackItemG[T]{n, len(n.items)})
		if n.leaf() {
			iter.stack[len(iter.stack)-1].i--
			break
		}
		n = iter.tr.isoLoad(&(*n.children)[len(n.items)], iter.mut)
	}
	s := &iter.stack[len(iter.stack)-1]
	iter.item = s.n.items[s.i]
	return true
}

func (iter *IterG[T]) Release() {
	if iter.tr == nil {
		return
	}
	iter.stack = nil
	iter.tr = nil
}

// 将迭代器移动到下一个项目。
// 如果树为空或迭代器已经到达树的末尾，返回 false。
func (iter *IterG[T]) Next() bool {
	if iter.tr == nil {
		return false
	}
	if !iter.seeked {
		return iter.First()
	}
	if len(iter.stack) == 0 {
		if iter.atstart {
			return iter.First() && iter.Next()
		}
		return false
	}
	s := &iter.stack[len(iter.stack)-1]
	s.i++
	if s.n.leaf() {
		if s.i == len(s.n.items) {
			for {
				iter.stack = iter.stack[:len(iter.stack)-1]
				if len(iter.stack) == 0 {
					iter.atend = true // 一旦迭代器移动到最后一个元素，atend 被设置为 true
					return false
				}
				s = &iter.stack[len(iter.stack)-1]
				if s.i < len(s.n.items) {
					break
				}
			}
		}
	} else {
		n := iter.tr.isoLoad(&(*s.n.children)[s.i], iter.mut)
		for {
			iter.stack = append(iter.stack, iterStackItemG[T]{n, 0})
			if n.leaf() {
				break
			}
			n = iter.tr.isoLoad(&(*n.children)[0], iter.mut)
		}
	}
	s = &iter.stack[len(iter.stack)-1]
	iter.item = s.n.items[s.i]
	return true
}

// 将迭代器移动到上一个项目。
// 如果树为空或迭代器已经到达树的起始位置，返回 false。
func (iter *IterG[T]) Prev() bool {
	if iter.tr == nil {
		return false
	}
	if !iter.seeked {
		return false
	}
	if len(iter.stack) == 0 {
		if iter.atend {
			return iter.Last() && iter.Prev()
		}
		return false
	}
	s := &iter.stack[len(iter.stack)-1]
	if s.n.leaf() {
		s.i--
		if s.i == -1 {
			// 需要回溯到上层节点
			for {
				iter.stack = iter.stack[:len(iter.stack)-1]
				if len(iter.stack) == 0 {
					iter.atstart = true // 一旦迭代器移动到第一个元素，atstart 被设置为 true
					return false
				}
				s = &iter.stack[len(iter.stack)-1]
				s.i--
				if s.i > -1 {
					break
				}
			}
		}
	} else {
		n := iter.tr.isoLoad(&(*s.n.children)[s.i], iter.mut)
		for {
			iter.stack = append(iter.stack, iterStackItemG[T]{n, len(n.items)})
			if n.leaf() {
				iter.stack[len(iter.stack)-1].i--
				break
			}
			n = iter.tr.isoLoad(&(*n.children)[len(n.items)], iter.mut)
		}
	}
	s = &iter.stack[len(iter.stack)-1]
	iter.item = s.n.items[s.i]
	return true
}

func (iter *IterG[T]) Item() T {
	return iter.item
}

// 返回所有按顺序排列的项目.
func (tr *BTreeG[T]) Items() []T {
	return tr.items(false)
}

func (tr *BTreeG[T]) ItemsMut() []T {
	return tr.items(true)
}

func (tr *BTreeG[T]) items(mut bool) []T {
	items := make([]T, 0, tr.Len())
	if tr.root != nil {
		items = tr.nodeItems(&tr.root, items, mut)
	}
	return items
}

func (tr *BTreeG[T]) nodeItems(cn **node[T], items []T, mut bool) []T {
	n := tr.isoLoad(cn, mut)
	if n.leaf() {
		return append(items, n.items...)
	}
	for i := 0; i < len(n.items); i++ {
		items = tr.nodeItems(&(*n.children)[i], items, mut)
		items = append(items, n.items[i])
	}
	return tr.nodeItems(&(*n.children)[len(*n.children)-1], items, mut)
}

func (tr *BTreeG[T]) Clear() {
	tr.root = nil
	tr.count = 0
}
