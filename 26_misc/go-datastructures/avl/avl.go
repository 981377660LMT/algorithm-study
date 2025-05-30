/*
Package avl includes an immutable AVL tree.

AVL trees can be used as the foundation for many functional data types.
Combined with a B+ tree, you can make an immutable index which serves as the
backbone for many different kinds of key/value stores.

Time complexities:
Space: O(n)
Insert: O(log n)
Delete: O(log n)
Get: O(log n)

The immutable version of the AVL tree is obviously going to be slower than
the mutable version but should offer higher read availability.
*/

// 可持久化AVL树
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	test()
}

func test() {
	n := int(1e5)
	slice := make([]int, n)
	for i := 0; i < n; i++ {
		slice[i] = n - i
	}
	rand.Shuffle(n, func(i, j int) { slice[i], slice[j] = slice[j], slice[i] })

	start := time.Now()
	tree := NewImmutableAVL()
	for _, v := range slice {
		tree, _ = tree.Insert(intEntry(v))
		tree.Get(intEntry(v))
	}
	fmt.Println("insert time:", time.Since(start).Milliseconds())
}

type intEntry int

func (i intEntry) Compare(other Entry) int {
	o := other.(intEntry)
	switch {
	case i < o:
		return -1
	case i > o:
		return 1
	default:
		return 0
	}
}

func demo() {
	// 1. 创建不可变AVL树
	tree := NewImmutableAVL()

	// 2. 插入
	// Insert返回 (新树, 被覆盖的旧值列表)
	nums := []Entry{intEntry(10), intEntry(5), intEntry(15), intEntry(10)}
	newTree, overwritten := tree.Insert(nums...)
	// 因为 10 重复插入一次，所以可能返回被覆盖的旧值
	fmt.Printf("overwritten: %v\n", overwritten)
	// oldTree 保持不变，newTree 是新的版本

	// 3. 查找
	// Get返回[]Entry，如果找不到则返回nil作为对应下标
	found := newTree.Get(intEntry(5), intEntry(10), intEntry(999))
	fmt.Printf("found: %v\n", found)
	// found[2] 可能就是 nil，因为 999 不存在

	// 4. 删除
	// Delete返回 (新树, 被删除的旧值列表)，无法删除就返回nil
	newTree2, deleted := newTree.Delete(intEntry(10))
	fmt.Printf("deleted: %v\n", deleted)

	// 5. 再查一次看看
	found2 := newTree2.Get(intEntry(10))
	fmt.Printf("found2 after delete: %v\n", found2) // 应该是 [nil]

	// 注意：tree/newTree/newTree2 是三个版本的树，并且可以并行安全地使用它们
}

// Immutable represents an immutable AVL tree.  This is acheived
// by branch copying.
type Immutable struct {
	root   *node
	number uint64
	dummy  node // helper for inserts.
}

// copy returns a copy of this immutable tree with a copy
// of the root and a new dummy helper for the insertion operation.
func (immutable *Immutable) copy() *Immutable {
	var root *node
	if immutable.root != nil {
		root = immutable.root.copy()
	}
	cp := &Immutable{
		root:   root,
		number: immutable.number,
		dummy:  *newNode(nil),
	}
	return cp
}

func (immutable *Immutable) resetDummy() {
	immutable.dummy.children[0], immutable.dummy.children[1] = nil, nil
	immutable.dummy.balance = 0
}

func (immutable *Immutable) init() {
	immutable.dummy = node{
		children: [2]*node{},
	}
}

func (immutable *Immutable) get(entry Entry) Entry {
	n := immutable.root
	var result int
	for n != nil {
		switch result = n.entry.Compare(entry); {
		case result == 0:
			return n.entry
		case result > 0:
			n = n.children[0]
		case result < 0:
			n = n.children[1]
		}
	}

	return nil
}

// Get will get the provided Entries from the tree.  If no matching
// Entry is found, a nil is returned in its place.
func (immutable *Immutable) Get(entries ...Entry) Entries {
	result := make(Entries, 0, len(entries))
	for _, e := range entries {
		result = append(result, immutable.get(e))
	}

	return result
}

// Len returns the number of items in this immutable.
func (immutable *Immutable) Len() uint64 {
	return immutable.number
}

func (immutable *Immutable) insert(entry Entry) Entry {
	// TODO: check cache to see if a node has already been copied.
	if immutable.root == nil {
		immutable.root = newNode(entry)
		immutable.number++
		return nil
	}

	immutable.resetDummy()
	var (
		dummy           = immutable.dummy
		p, s, q         *node
		dir, normalized int
		helper          = &dummy
	)

	// set this AFTER clearing dummy
	helper.children[1] = immutable.root
	// we'll go ahead and copy on the way down as we'll need to branch
	// copy no matter what.
	for s, p = helper.children[1], helper.children[1]; ; {
		dir = p.entry.Compare(entry)

		normalized = normalizeComparison(dir)
		if dir > 0 { // go left
			if p.children[0] != nil {
				q = p.children[0].copy()
				p.children[0] = q
			} else {
				q = nil
			}
		} else if dir < 0 { // go right
			if p.children[1] != nil {
				q = p.children[1].copy()
				p.children[1] = q
			} else {
				q = nil
			}
		} else { // equality
			oldEntry := p.entry
			p.entry = entry
			return oldEntry
		}
		if q == nil {
			break
		}

		if q.balance != 0 {
			helper = p
			s = q
		}
		p = q
	}

	immutable.number++
	q = newNode(entry)
	p.children[normalized] = q

	immutable.root = dummy.children[1]
	for p = s; p != q; p = p.children[normalized] {
		normalized = normalizeComparison(p.entry.Compare(entry))
		if normalized == 0 {
			p.balance--
		} else {
			p.balance++
		}
	}

	q = s

	if math.Abs(float64(s.balance)) > 1 {
		normalized = normalizeComparison(s.entry.Compare(entry))
		s = insertBalance(s, normalized)
	}

	if q == dummy.children[1] {
		immutable.root = s
	} else {
		helper.children[intFromBool(helper.children[1] == q)] = s
	}
	return nil
}

// Insert will add the provided entries into the tree and return the new
// state.  Also returned is a list of Entries that were overwritten.  If
// nothing was overwritten for an Entry, a nil is returned in its place.
func (immutable *Immutable) Insert(entries ...Entry) (*Immutable, Entries) {
	if len(entries) == 0 {
		return immutable, Entries{}
	}

	overwritten := make(Entries, 0, len(entries))
	cp := immutable.copy()
	for _, e := range entries {
		overwritten = append(overwritten, cp.insert(e))
	}

	return cp, overwritten
}

func (immutable *Immutable) delete(entry Entry) Entry {
	// TODO: reuse cache and dirs, check cache to see if nodes
	// really need to be copied.
	if immutable.root == nil { // easy case, nothing to remove
		return nil
	}

	var (
		// we are going to make a list here representing our stack.
		// This means we don't have to copy if a value wasn't found.
		cache                      = make(nodes, 64)
		it, p, q                   *node
		top, done, dir, normalized int
		dirs                       = make([]int, 64)
		oldEntry                   Entry
	)

	it = immutable.root

	for {
		if it == nil {
			return nil
		}

		dir = it.entry.Compare(entry)
		if dir == 0 {
			break
		}
		normalized = normalizeComparison(dir)
		dirs[top] = normalized
		cache[top] = it
		top++
		it = it.children[normalized]
	}
	immutable.number--
	oldEntry = it.entry // we need to return this

	// we need to make a branch copy now
	for i := 0; i < top; i++ { // first item will be root
		p = cache[i]
		if p.children[dirs[i]] != nil {
			q = p.children[dirs[i]].copy()
			p.children[dirs[i]] = q
			if i != top-1 {
				cache[i+1] = q
			}
		}
	}
	it = it.copy() // the node we found needs to be copied

	oldTop := top
	if it.children[0] == nil || it.children[1] == nil { // need to set children on parent, splicing out
		dir = intFromBool(it.children[0] == nil)
		if top != 0 {
			cache[top-1].children[dirs[top-1]] = it.children[dir]
		} else {
			immutable.root = it.children[dir]
		}
	} else { // climb up and set heirs
		heir := it.children[1]
		dirs[top] = 1
		cache[top] = it
		top++

		for heir.children[0] != nil {
			dirs[top] = 0
			cache[top] = heir
			top++
			heir = heir.children[0]
		}

		it.entry = heir.entry
		if oldTop != 0 {
			cache[oldTop-1].children[dirs[oldTop-1]] = it
		} else {
			immutable.root = it
		}
		cache[top-1].children[intFromBool(cache[top-1] == it)] = heir.children[1]
	}

	for top-1 >= 0 && done == 0 {
		top--
		// set bounded balance
		if dirs[top] != 0 {
			cache[top].balance--
		} else {
			cache[top].balance++
		}

		if math.Abs(float64(cache[top].balance)) == 1 {
			break
		} else if math.Abs(float64(cache[top].balance)) > 1 {
			// any rotations done here
			cache[top] = removeBalance(cache[top], dirs[top], &done)

			if top != 0 {
				cache[top-1].children[dirs[top-1]] = cache[top]
			} else {
				immutable.root = cache[0]
			}
		}

	}

	return oldEntry
}

// Delete will remove the provided entries from this AVL tree and
// return a new tree and any entries removed.  If an entry could not
// be found, nil is returned in its place.
func (immutable *Immutable) Delete(entries ...Entry) (*Immutable, Entries) {
	if len(entries) == 0 {
		return immutable, Entries{}
	}

	deleted := make(Entries, 0, len(entries))
	cp := immutable.copy()
	for _, e := range entries {
		deleted = append(deleted, cp.delete(e))
	}

	return cp, deleted
}

func insertBalance(root *node, dir int) *node {
	n := root.children[dir]
	var bal int8
	if dir == 0 {
		bal = -1
	} else {
		bal = 1
	}

	if n.balance == bal {
		root.balance, n.balance = 0, 0
		root = rotate(root, takeOpposite(dir))
	} else {
		adjustBalance(root, dir, int(bal))
		root = doubleRotate(root, takeOpposite(dir))
	}

	return root
}

func removeBalance(root *node, dir int, done *int) *node {
	n := root.children[takeOpposite(dir)].copy()
	root.children[takeOpposite(dir)] = n
	var bal int8
	if dir == 0 {
		bal = -1
	} else {
		bal = 1
	}

	if n.balance == -bal {
		root.balance, n.balance = 0, 0
		root = rotate(root, dir)
	} else if n.balance == bal {
		adjustBalance(root, takeOpposite(dir), int(-bal))
		root = doubleRotate(root, dir)
	} else {
		root.balance = -bal
		n.balance = bal
		root = rotate(root, dir)
		*done = 1
	}

	return root
}

func intFromBool(value bool) int {
	if value {
		return 1
	}

	return 0
}

func takeOpposite(value int) int {
	return 1 - value
}

func adjustBalance(root *node, dir, bal int) {
	n := root.children[dir]
	nn := n.children[takeOpposite(dir)]

	if nn.balance == 0 {
		root.balance, n.balance = 0, 0
	} else if int(nn.balance) == bal {
		root.balance = int8(-bal)
		n.balance = 0
	} else {
		root.balance = 0
		n.balance = int8(bal)
	}
	nn.balance = 0
}

func rotate(parent *node, dir int) *node {
	otherDir := takeOpposite(dir)

	child := parent.children[otherDir]
	parent.children[otherDir] = child.children[dir]
	child.children[dir] = parent

	return child
}

func doubleRotate(parent *node, dir int) *node {
	otherDir := takeOpposite(dir)

	parent.children[otherDir] = rotate(parent.children[otherDir], otherDir)
	return rotate(parent, dir)
}

// normalizeComparison converts the value returned from Entry.Compare
// into a direction, ie, left or right, 0 or 1.
func normalizeComparison(i int) int {
	if i < 0 {
		return 1
	}

	if i > 0 {
		return 0
	}

	return -1
}

// NewImmutableAVL allocates, initializes, and returns a new immutable
// AVL tree.
func NewImmutableAVL() *Immutable {
	immutable := &Immutable{}
	immutable.init()
	return immutable
}

// Entries is a list of type Entry.
type Entries []Entry

// Entry represents all items that can be placed into the AVL tree.
// They must implement a Compare method that can be used to determine
// the Entry's correct place in the tree.  Any object can implement
// Compare.
type Entry interface {
	// Compare should return a value indicating the relationship
	// of this Entry to the provided Entry.  A -1 means this entry
	// is less than, 0 means equality, and 1 means greater than.
	Compare(Entry) int
}
type nodes []*node

func (ns nodes) reset() {
	for i := range ns {
		ns[i] = nil
	}
}

type node struct {
	balance  int8 // bounded, |balance| should be <= 1
	children [2]*node
	entry    Entry
}

// copy returns a copy of this node with pointers to the original
// children.
func (n *node) copy() *node {
	return &node{
		balance:  n.balance,
		children: [2]*node{n.children[0], n.children[1]},
		entry:    n.entry,
	}
}

// newNode returns a new node for the provided entry.  A nil
// entry is used to represent the dummy node.
func newNode(entry Entry) *node {
	return &node{
		entry:    entry,
		children: [2]*node{},
	}
}
