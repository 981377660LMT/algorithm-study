/*
Package ctrie provides an implementation of the Ctrie data structure, which is
a concurrent, lock-free hash trie. This data structure was originally presented
in the paper Concurrent Tries with Efficient Non-Blocking Snapshots:

https://axel22.github.io/resources/docs/ctries-snapshot.pdf
*/
package main

import (
	"bytes"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"sync/atomic"
	"unsafe"
)

func main() {
	// 1. 创建空 Ctrie，使用默认的 FNV-1a 哈希
	c := NewCtrie(nil)

	// 2. 插入一些键值
	c.Insert([]byte("foo"), 123)
	c.Insert([]byte("bar"), "hello")

	// 3. 查找
	val, ok := c.Lookup([]byte("foo"))
	fmt.Println("Lookup foo:", val, ok) // 123, true

	// 4. 删除
	removedVal, removedOk := c.Remove([]byte("bar"))
	fmt.Println("Removed bar:", removedVal, removedOk) // "hello", true

	// 5. Snapshot
	snap := c.ReadOnlySnapshot()
	// snap is read-only, c is still modifiable
	// 进行遍历
	for e := range snap.Iterator(nil) {
		fmt.Printf("Key: %s, Value: %v\n", string(e.Key), e.Value)
	}

	// 6. 并发场景下多个 goroutine 可以对 c 执行 Insert、Remove、Lookup
	// 并且snap版本不会变动
}

const (
	// w controls the number of branches at a node (2^w branches).
	w = 5

	// exp2 is 2^w, which is the hashcode space.
	exp2 = 32
)

// HashFactory returns a new Hash32 used to hash keys.
type HashFactory func() hash.Hash32

func defaultHashFactory() hash.Hash32 {
	return fnv.New32a()
}

// Ctrie is a concurrent, lock-free hash trie. By default, keys are hashed
// using FNV-1a unless a HashFactory is provided to New.
type Ctrie struct {
	root        *iNode      // 根节点
	readOnly    bool        // 是否只读模式，若为 `true`，则不允许写操作（插入、删除），试图执行会 panic
	hashFactory HashFactory // 用于对 key 做哈希的函数工厂
}

// generation demarcates Ctrie snapshots. We use a heap-allocated reference
// instead of an integer to avoid integer overflows. Struct must have a field
// on it since two distinct zero-size variables may have the same address in
// memory.
type generation struct{ _ int }

// iNode is an indirection node. I-nodes remain present in the Ctrie even as
// nodes above and below change. Thread-safety is achieved in part by
// performing CAS operations on the I-node instead of the internal node array.
type iNode struct {
	main *mainNode
	gen  *generation

	// rdcss is set during an RDCSS operation. The I-node is actually a wrapper
	// around the descriptor in this case so that a single type is used during
	// CAS operations on the root.
	rdcss *rdcssDescriptor
}

// copyToGen returns a copy of this I-node copied to the given generation.
func (i *iNode) copyToGen(gen *generation, ctrie *Ctrie) *iNode {
	nin := &iNode{gen: gen}
	main := gcasRead(i, ctrie)
	atomic.StorePointer(
		(*unsafe.Pointer)(unsafe.Pointer(&nin.main)), unsafe.Pointer(main))
	return nin
}

// mainNode is either a cNode, tNode, lNode, or failed node which makes up an
// I-node.
type mainNode struct {
	cNode  *cNode    // 内部节点（包含若干分支）
	tNode  *tNode    // tomb node，删除后留的标记
	lNode  *lNode    // list node，表示出现哈希冲突时将这些键值放在不可变链表里
	failed *mainNode // 在 GCAS 操作失败或回退时使用

	prev *mainNode
}

// cNode is an internal main node containing a bitmap and the array with
// references to branch nodes. A branch node is either another I-node or a
// singleton S-node.
type cNode struct {
	bmp   uint32   // bitmap，标记哪几个下标有效
	array []branch // branch 是 *iNode 或 *sNode
	gen   *generation
}

// newMainNode is a recursive constructor which creates a new mainNode. This
// mainNode will consist of cNodes as long as the hashcode chunks of the two
// keys are equal at the given level. If the level exceeds 2^w, an lNode is
// created.
func newMainNode(x *sNode, xhc uint32, y *sNode, yhc uint32, lev uint, gen *generation) *mainNode {
	if lev < exp2 {
		xidx := (xhc >> lev) & 0x1f
		yidx := (yhc >> lev) & 0x1f
		bmp := uint32((1 << xidx) | (1 << yidx))

		if xidx == yidx {
			// Recurse when indexes are equal.
			main := newMainNode(x, xhc, y, yhc, lev+w, gen)
			iNode := &iNode{main: main, gen: gen}
			return &mainNode{cNode: &cNode{bmp, []branch{iNode}, gen}}
		}
		if xidx < yidx {
			return &mainNode{cNode: &cNode{bmp, []branch{x, y}, gen}}
		}
		return &mainNode{cNode: &cNode{bmp, []branch{y, x}, gen}}
	}
	l := Empty.Add(x).Add(y)
	return &mainNode{lNode: &lNode{l}}
}

// inserted returns a copy of this cNode with the new entry at the given
// position.
func (c *cNode) inserted(pos, flag uint32, br branch, gen *generation) *cNode {
	length := uint32(len(c.array))
	bmp := c.bmp
	array := make([]branch, length+1)
	copy(array, c.array)
	array[pos] = br
	for i, x := pos, uint32(0); x < length-pos; i++ {
		array[i+1] = c.array[i]
		x++
	}
	ncn := &cNode{bmp: bmp | flag, array: array, gen: gen}
	return ncn
}

// updated returns a copy of this cNode with the entry at the given index
// updated.
func (c *cNode) updated(pos uint32, br branch, gen *generation) *cNode {
	array := make([]branch, len(c.array))
	copy(array, c.array)
	array[pos] = br
	ncn := &cNode{bmp: c.bmp, array: array, gen: gen}
	return ncn
}

// removed returns a copy of this cNode with the entry at the given index
// removed.
func (c *cNode) removed(pos, flag uint32, gen *generation) *cNode {
	length := uint32(len(c.array))
	bmp := c.bmp
	array := make([]branch, length-1)
	for i := uint32(0); i < pos; i++ {
		array[i] = c.array[i]
	}
	for i, x := pos, uint32(0); x < length-pos-1; i++ {
		array[i] = c.array[i+1]
		x++
	}
	ncn := &cNode{bmp: bmp ^ flag, array: array, gen: gen}
	return ncn
}

// renewed returns a copy of this cNode with the I-nodes below it copied to the
// given generation.
func (c *cNode) renewed(gen *generation, ctrie *Ctrie) *cNode {
	array := make([]branch, len(c.array))
	for i, br := range c.array {
		switch t := br.(type) {
		case *iNode:
			array[i] = t.copyToGen(gen, ctrie)
		default:
			array[i] = br
		}
	}
	return &cNode{bmp: c.bmp, array: array, gen: gen}
}

// tNode is tomb node which is a special node used to ensure proper ordering
// during removals.
type tNode struct {
	*sNode
}

// untombed returns the S-node contained by the T-node.
func (t *tNode) untombed() *sNode {
	return &sNode{&Entry{Key: t.Key, hash: t.hash, Value: t.Value}}
}

// 当哈希前缀超过一定层次（或有碰撞），
// 可能把一组键值对存在一个不可变持久化链表 `PersistentList` 中
type lNode struct {
	PersistentList
}

// entry returns the first S-node contained in the L-node.
func (l *lNode) entry() *sNode {
	head, _ := l.Head()
	return head.(*sNode)
}

// lookup returns the value at the given entry in the L-node or returns false
// if it's not contained.
func (l *lNode) lookup(e *Entry) (interface{}, bool) {
	found, ok := l.Find(func(sn interface{}) bool {
		return bytes.Equal(e.Key, sn.(*sNode).Key)
	})
	if !ok {
		return nil, false
	}
	return found.(*sNode).Value, true
}

// inserted creates a new L-node with the added entry.
func (l *lNode) inserted(entry *Entry) *lNode {
	return &lNode{l.removed(entry).Add(&sNode{entry})}
}

// removed creates a new L-node with the entry removed.
func (l *lNode) removed(e *Entry) *lNode {
	idx := l.FindIndex(func(sn interface{}) bool {
		return bytes.Equal(e.Key, sn.(*sNode).Key)
	})
	if idx < 0 {
		return l
	}
	nl, _ := l.Remove(uint(idx))
	return &lNode{nl}
}

// length returns the L-node list length.
func (l *lNode) length() uint {
	return l.Length()
}

// branch is either an iNode or sNode.
type branch interface{}

// Entry contains a Ctrie key-value pair.
type Entry struct {
	Key   []byte
	Value interface{}
	hash  uint32
}

// sNode is a singleton node which contains a single key and value.
type sNode struct {
	*Entry
}

// NewCtrie creates an empty Ctrie which uses the provided HashFactory for key
// hashing. If nil is passed in, it will default to FNV-1a hashing.
func NewCtrie(hashFactory HashFactory) *Ctrie {
	if hashFactory == nil {
		hashFactory = defaultHashFactory
	}
	root := &iNode{main: &mainNode{cNode: &cNode{}}}
	return newCtrie(root, hashFactory, false)
}

func newCtrie(root *iNode, hashFactory HashFactory, readOnly bool) *Ctrie {
	return &Ctrie{
		root:        root,
		hashFactory: hashFactory,
		readOnly:    readOnly,
	}
}

// Insert adds the key-value pair to the Ctrie, replacing the existing value if
// the key already exists.
func (c *Ctrie) Insert(key []byte, value interface{}) {
	c.assertReadWrite()
	c.insert(&Entry{
		Key:   key,
		Value: value,
		hash:  c.hash(key),
	})
}

// Lookup returns the value for the associated key or returns false if the key
// doesn't exist.
func (c *Ctrie) Lookup(key []byte) (interface{}, bool) {
	return c.lookup(&Entry{Key: key, hash: c.hash(key)})
}

// Remove deletes the value for the associated key, returning true if it was
// removed or false if the entry doesn't exist.
func (c *Ctrie) Remove(key []byte) (interface{}, bool) {
	c.assertReadWrite()
	return c.remove(&Entry{Key: key, hash: c.hash(key)})
}

// Snapshot returns a stable, point-in-time snapshot of the Ctrie. If the Ctrie
// is read-only, the returned Ctrie will also be read-only.
func (c *Ctrie) Snapshot() *Ctrie {
	return c.snapshot(c.readOnly)
}

// ReadOnlySnapshot returns a stable, point-in-time snapshot of the Ctrie which
// is read-only. Write operations on a read-only snapshot will panic.
func (c *Ctrie) ReadOnlySnapshot() *Ctrie {
	return c.snapshot(true)
}

// **快照**允许在无锁情况下进行安全遍历，而不会受到并发修改的影响
func (c *Ctrie) snapshot(readOnly bool) *Ctrie {
	if readOnly && c.readOnly {
		return c
	}
	for {
		root := c.readRoot()
		main := gcasRead(root, c)
		if c.rdcssRoot(root, main, root.copyToGen(&generation{}, c)) {
			if readOnly {
				// For a read-only snapshot, we can share the old generation
				// root.
				return newCtrie(root, c.hashFactory, readOnly)
			}
			// For a read-write snapshot, we need to take a copy of the root
			// in the new generation.
			return newCtrie(c.readRoot().copyToGen(&generation{}, c), c.hashFactory, readOnly)
		}
	}
}

// Clear removes all keys from the Ctrie.
func (c *Ctrie) Clear() {
	for {
		root := c.readRoot()
		gen := &generation{}
		newRoot := &iNode{
			main: &mainNode{cNode: &cNode{array: make([]branch, 0), gen: gen}},
			gen:  gen,
		}
		if c.rdcssRoot(root, gcasRead(root, c), newRoot) {
			return
		}
	}
}

// Iterator returns a channel which yields the Entries of the Ctrie. If a
// cancel channel is provided, closing it will terminate and close the iterator
// channel. Note that if a cancel channel is not used and not every entry is
// read from the iterator, a goroutine will leak.
func (c *Ctrie) Iterator(cancel <-chan struct{}) <-chan *Entry {
	ch := make(chan *Entry)
	snapshot := c.ReadOnlySnapshot()
	go func() {
		snapshot.traverse(snapshot.readRoot(), ch, cancel)
		close(ch)
	}()
	return ch
}

// Size returns the number of keys in the Ctrie.
func (c *Ctrie) Size() uint {
	// TODO: The size operation can be optimized further by caching the size
	// information in main nodes of a read-only Ctrie – this reduces the
	// amortized complexity of the size operation to O(1) because the size
	// computation is amortized across the update operations that occurred
	// since the last snapshot.
	size := uint(0)
	for _ = range c.Iterator(nil) {
		size++
	}
	return size
}

var errCanceled = errors.New("canceled")

func (c *Ctrie) traverse(i *iNode, ch chan<- *Entry, cancel <-chan struct{}) error {
	main := gcasRead(i, c)
	switch {
	case main.cNode != nil:
		for _, br := range main.cNode.array {
			switch b := br.(type) {
			case *iNode:
				if err := c.traverse(b, ch, cancel); err != nil {
					return err
				}
			case *sNode:
				select {
				case ch <- b.Entry:
				case <-cancel:
					return errCanceled
				}
			}
		}
	case main.lNode != nil:
		for _, e := range main.lNode.Map(func(sn interface{}) interface{} {
			return sn.(*sNode).Entry
		}) {
			select {
			case ch <- e.(*Entry):
			case <-cancel:
				return errCanceled
			}
		}
	case main.tNode != nil:
		select {
		case ch <- main.tNode.Entry:
		case <-cancel:
			return errCanceled
		}
	}
	return nil
}

func (c *Ctrie) assertReadWrite() {
	if c.readOnly {
		panic("Cannot modify read-only snapshot")
	}
}

func (c *Ctrie) insert(entry *Entry) {
	root := c.readRoot()
	if !c.iinsert(root, entry, 0, nil, root.gen) {
		c.insert(entry)
	}
}

func (c *Ctrie) lookup(entry *Entry) (interface{}, bool) {
	root := c.readRoot()
	result, exists, ok := c.ilookup(root, entry, 0, nil, root.gen)
	for !ok {
		return c.lookup(entry)
	}
	return result, exists
}

func (c *Ctrie) remove(entry *Entry) (interface{}, bool) {
	root := c.readRoot()
	result, exists, ok := c.iremove(root, entry, 0, nil, root.gen)
	for !ok {
		return c.remove(entry)
	}
	return result, exists
}

func (c *Ctrie) hash(k []byte) uint32 {
	hasher := c.hashFactory()
	hasher.Write(k)
	return hasher.Sum32()
}

// iinsert attempts to insert the entry into the Ctrie. If false is returned,
// the operation should be retried.
func (c *Ctrie) iinsert(i *iNode, entry *Entry, lev uint, parent *iNode, startGen *generation) bool {
	// Linearization point.
	main := gcasRead(i, c)
	switch {
	case main.cNode != nil:
		cn := main.cNode
		flag, pos := flagPos(entry.hash, lev, cn.bmp)
		if cn.bmp&flag == 0 {
			// If the relevant bit is not in the bitmap, then a copy of the
			// cNode with the new entry is created. The linearization point is
			// a successful CAS.
			rn := cn
			if cn.gen != i.gen {
				rn = cn.renewed(i.gen, c)
			}
			ncn := &mainNode{cNode: rn.inserted(pos, flag, &sNode{entry}, i.gen)}
			return gcas(i, main, ncn, c)
		}
		// If the relevant bit is present in the bitmap, then its corresponding
		// branch is read from the array.
		branch := cn.array[pos]
		switch branch.(type) {
		case *iNode:
			// If the branch is an I-node, then iinsert is called recursively.
			in := branch.(*iNode)
			if startGen == in.gen {
				return c.iinsert(in, entry, lev+w, i, startGen)
			}
			if gcas(i, main, &mainNode{cNode: cn.renewed(startGen, c)}, c) {
				return c.iinsert(i, entry, lev, parent, startGen)
			}
			return false
		case *sNode:
			sn := branch.(*sNode)
			if !bytes.Equal(sn.Key, entry.Key) {
				// If the branch is an S-node and its key is not equal to the
				// key being inserted, then the Ctrie has to be extended with
				// an additional level. The C-node is replaced with its updated
				// version, created using the updated function that adds a new
				// I-node at the respective position. The new Inode has its
				// main node pointing to a C-node with both keys. The
				// linearization point is a successful CAS.
				rn := cn
				if cn.gen != i.gen {
					rn = cn.renewed(i.gen, c)
				}
				nsn := &sNode{entry}
				nin := &iNode{main: newMainNode(sn, sn.hash, nsn, nsn.hash, lev+w, i.gen), gen: i.gen}
				ncn := &mainNode{cNode: rn.updated(pos, nin, i.gen)}
				return gcas(i, main, ncn, c)
			}
			// If the key in the S-node is equal to the key being inserted,
			// then the C-node is replaced with its updated version with a new
			// S-node. The linearization point is a successful CAS.
			ncn := &mainNode{cNode: cn.updated(pos, &sNode{entry}, i.gen)}
			return gcas(i, main, ncn, c)
		default:
			panic("Ctrie is in an invalid state")
		}
	case main.tNode != nil:
		clean(parent, lev-w, c)
		return false
	case main.lNode != nil:
		nln := &mainNode{lNode: main.lNode.inserted(entry)}
		return gcas(i, main, nln, c)
	default:
		panic("Ctrie is in an invalid state")
	}
}

// ilookup attempts to fetch the entry from the Ctrie. The first two return
// values are the entry value and whether or not the entry was contained in the
// Ctrie. The last bool indicates if the operation succeeded. False means it
// should be retried.
func (c *Ctrie) ilookup(i *iNode, entry *Entry, lev uint, parent *iNode, startGen *generation) (interface{}, bool, bool) {
	// Linearization point.
	main := gcasRead(i, c)
	switch {
	case main.cNode != nil:
		cn := main.cNode
		flag, pos := flagPos(entry.hash, lev, cn.bmp)
		if cn.bmp&flag == 0 {
			// If the bitmap does not contain the relevant bit, a key with the
			// required hashcode prefix is not present in the trie.
			return nil, false, true
		}
		// Otherwise, the relevant branch at index pos is read from the array.
		branch := cn.array[pos]
		switch branch.(type) {
		case *iNode:
			// If the branch is an I-node, the ilookup procedure is called
			// recursively at the next level.
			in := branch.(*iNode)
			if c.readOnly || startGen == in.gen {
				return c.ilookup(in, entry, lev+w, i, startGen)
			}
			if gcas(i, main, &mainNode{cNode: cn.renewed(startGen, c)}, c) {
				return c.ilookup(i, entry, lev, parent, startGen)
			}
			return nil, false, false
		case *sNode:
			// If the branch is an S-node, then the key within the S-node is
			// compared with the key being searched – these two keys have the
			// same hashcode prefixes, but they need not be equal. If they are
			// equal, the corresponding value from the S-node is
			// returned and a NOTFOUND value otherwise.
			sn := branch.(*sNode)
			if bytes.Equal(sn.Key, entry.Key) {
				return sn.Value, true, true
			}
			return nil, false, true
		default:
			panic("Ctrie is in an invalid state")
		}
	case main.tNode != nil:
		return cleanReadOnly(main.tNode, lev, parent, c, entry)
	case main.lNode != nil:
		// Hash collisions are handled using L-nodes, which are essentially
		// persistent linked lists.
		val, ok := main.lNode.lookup(entry)
		return val, ok, true
	default:
		panic("Ctrie is in an invalid state")
	}
}

// iremove attempts to remove the entry from the Ctrie. The first two return
// values are the entry value and whether or not the entry was contained in the
// Ctrie. The last bool indicates if the operation succeeded. False means it
// should be retried.
func (c *Ctrie) iremove(i *iNode, entry *Entry, lev uint, parent *iNode, startGen *generation) (interface{}, bool, bool) {
	// Linearization point.
	main := gcasRead(i, c)
	switch {
	case main.cNode != nil:
		cn := main.cNode
		flag, pos := flagPos(entry.hash, lev, cn.bmp)
		if cn.bmp&flag == 0 {
			// If the bitmap does not contain the relevant bit, a key with the
			// required hashcode prefix is not present in the trie.
			return nil, false, true
		}
		// Otherwise, the relevant branch at index pos is read from the array.
		branch := cn.array[pos]
		switch branch.(type) {
		case *iNode:
			// If the branch is an I-node, the iremove procedure is called
			// recursively at the next level.
			in := branch.(*iNode)
			if startGen == in.gen {
				return c.iremove(in, entry, lev+w, i, startGen)
			}
			if gcas(i, main, &mainNode{cNode: cn.renewed(startGen, c)}, c) {
				return c.iremove(i, entry, lev, parent, startGen)
			}
			return nil, false, false
		case *sNode:
			// If the branch is an S-node, its key is compared against the key
			// being removed.
			sn := branch.(*sNode)
			if !bytes.Equal(sn.Key, entry.Key) {
				// If the keys are not equal, the NOTFOUND value is returned.
				return nil, false, true
			}
			//  If the keys are equal, a copy of the current node without the
			//  S-node is created. The contraction of the copy is then created
			//  using the toContracted procedure. A successful CAS will
			//  substitute the old C-node with the copied C-node, thus removing
			//  the S-node with the given key from the trie – this is the
			//  linearization point
			ncn := cn.removed(pos, flag, i.gen)
			cntr := toContracted(ncn, lev)
			if gcas(i, main, cntr, c) {
				if parent != nil {
					main = gcasRead(i, c)
					if main.tNode != nil {
						cleanParent(parent, i, entry.hash, lev-w, c, startGen)
					}
				}
				return sn.Value, true, true
			}
			return nil, false, false
		default:
			panic("Ctrie is in an invalid state")
		}
	case main.tNode != nil:
		clean(parent, lev-w, c)
		return nil, false, false
	case main.lNode != nil:
		nln := &mainNode{lNode: main.lNode.removed(entry)}
		if nln.lNode.length() == 1 {
			nln = entomb(nln.lNode.entry())
		}
		if gcas(i, main, nln, c) {
			val, ok := main.lNode.lookup(entry)
			return val, ok, true
		}
		return nil, false, true
	default:
		panic("Ctrie is in an invalid state")
	}
}

// toContracted ensures that every I-node except the root points to a C-node
// with at least one branch. If a given C-Node has only a single S-node below
// it and is not at the root level, a T-node which wraps the S-node is
// returned.
func toContracted(cn *cNode, lev uint) *mainNode {
	if lev > 0 && len(cn.array) == 1 {
		branch := cn.array[0]
		switch branch.(type) {
		case *sNode:
			return entomb(branch.(*sNode))
		default:
			return &mainNode{cNode: cn}
		}
	}
	return &mainNode{cNode: cn}
}

// toCompressed compacts the C-node as a performance optimization.
func toCompressed(cn *cNode, lev uint) *mainNode {
	tmpArray := make([]branch, len(cn.array))
	for i, sub := range cn.array {
		switch sub.(type) {
		case *iNode:
			inode := sub.(*iNode)
			mainPtr := (*unsafe.Pointer)(unsafe.Pointer(&inode.main))
			main := (*mainNode)(atomic.LoadPointer(mainPtr))
			tmpArray[i] = resurrect(inode, main)
		case *sNode:
			tmpArray[i] = sub
		default:
			panic("Ctrie is in an invalid state")
		}
	}

	return toContracted(&cNode{bmp: cn.bmp, array: tmpArray}, lev)
}

func entomb(m *sNode) *mainNode {
	return &mainNode{tNode: &tNode{m}}
}

func resurrect(iNode *iNode, main *mainNode) branch {
	if main.tNode != nil {
		return main.tNode.untombed()
	}
	return iNode
}

func clean(i *iNode, lev uint, ctrie *Ctrie) bool {
	main := gcasRead(i, ctrie)
	if main.cNode != nil {
		return gcas(i, main, toCompressed(main.cNode, lev), ctrie)
	}
	return true
}

func cleanReadOnly(tn *tNode, lev uint, p *iNode, ctrie *Ctrie, entry *Entry) (val interface{}, exists bool, ok bool) {
	if !ctrie.readOnly {
		clean(p, lev-5, ctrie)
		return nil, false, false
	}
	if tn.hash == entry.hash && bytes.Equal(tn.Key, entry.Key) {
		return tn.Value, true, true
	}
	return nil, false, true
}

func cleanParent(p, i *iNode, hc uint32, lev uint, ctrie *Ctrie, startGen *generation) {
	var (
		mainPtr  = (*unsafe.Pointer)(unsafe.Pointer(&i.main))
		main     = (*mainNode)(atomic.LoadPointer(mainPtr))
		pMainPtr = (*unsafe.Pointer)(unsafe.Pointer(&p.main))
		pMain    = (*mainNode)(atomic.LoadPointer(pMainPtr))
	)
	if pMain.cNode != nil {
		flag, pos := flagPos(hc, lev, pMain.cNode.bmp)
		if pMain.cNode.bmp&flag != 0 {
			sub := pMain.cNode.array[pos]
			if sub == i && main.tNode != nil {
				ncn := pMain.cNode.updated(pos, resurrect(i, main), i.gen)
				if !gcas(p, pMain, toContracted(ncn, lev), ctrie) && ctrie.readRoot().gen == startGen {
					cleanParent(p, i, hc, lev, ctrie, startGen)
				}
			}
		}
	}
}

func flagPos(hashcode uint32, lev uint, bmp uint32) (uint32, uint32) {
	idx := (hashcode >> lev) & 0x1f
	flag := uint32(1) << uint32(idx)
	mask := uint32(flag - 1)
	pos := bitCount(bmp & mask)
	return flag, pos
}

func bitCount(x uint32) uint32 {
	x -= (x >> 1) & 0x55555555
	x = ((x >> 2) & 0x33333333) + (x & 0x33333333)
	x = ((x >> 4) + x) & 0x0f0f0f0f
	x *= 0x01010101
	return x >> 24
}

// gcas is a generation-compare-and-swap which has semantics similar to RDCSS,
// but it does not create the intermediate object except in the case of
// failures that occur due to the snapshot being taken. This ensures that the
// write occurs only if the Ctrie root generation has remained the same in
// addition to the I-node having the expected value.
func gcas(in *iNode, old, n *mainNode, ct *Ctrie) bool {
	prevPtr := (*unsafe.Pointer)(unsafe.Pointer(&n.prev))
	atomic.StorePointer(prevPtr, unsafe.Pointer(old))
	if atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&in.main)),
		unsafe.Pointer(old), unsafe.Pointer(n)) {
		gcasComplete(in, n, ct)
		return atomic.LoadPointer(prevPtr) == nil
	}
	return false
}

// gcasRead performs a GCAS-linearizable read of the I-node's main node.
func gcasRead(in *iNode, ctrie *Ctrie) *mainNode {
	m := (*mainNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&in.main))))
	prev := (*mainNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&m.prev))))
	if prev == nil {
		return m
	}
	return gcasComplete(in, m, ctrie)
}

// gcasComplete commits the GCAS operation.
func gcasComplete(i *iNode, m *mainNode, ctrie *Ctrie) *mainNode {
	for {
		if m == nil {
			return nil
		}
		prev := (*mainNode)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&m.prev))))
		root := ctrie.rdcssReadRoot(true)
		if prev == nil {
			return m
		}

		if prev.failed != nil {
			// Signals GCAS failure. Swap old value back into I-node.
			fn := prev.failed
			if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&i.main)),
				unsafe.Pointer(m), unsafe.Pointer(fn)) {
				return fn
			}
			m = (*mainNode)(atomic.LoadPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&i.main))))
			continue
		}

		if root.gen == i.gen && !ctrie.readOnly {
			// Commit GCAS.
			if atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&m.prev)), unsafe.Pointer(prev), nil) {
				return m
			}
			continue
		}

		// Generations did not match. Store failed node on prev to signal
		// I-node's main node must be set back to the previous value.
		atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&m.prev)),
			unsafe.Pointer(prev),
			unsafe.Pointer(&mainNode{failed: prev}))
		m = (*mainNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&i.main))))
		return gcasComplete(i, m, ctrie)
	}
}

// rdcssDescriptor is an intermediate struct which communicates the intent to
// replace the value in an I-node and check that the root's generation has not
// changed before committing to the new value.
type rdcssDescriptor struct {
	old       *iNode
	expected  *mainNode
	nv        *iNode
	committed int32
}

// readRoot performs a linearizable read of the Ctrie root. This operation is
// prioritized so that if another thread performs a GCAS on the root, a
// deadlock does not occur.
func (c *Ctrie) readRoot() *iNode {
	return c.rdcssReadRoot(false)
}

// rdcssReadRoot performs a RDCSS-linearizable read of the Ctrie root with the
// given priority.
func (c *Ctrie) rdcssReadRoot(abort bool) *iNode {
	r := (*iNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&c.root))))
	if r.rdcss != nil {
		return c.rdcssComplete(abort)
	}
	return r
}

// rdcssRoot performs a RDCSS on the Ctrie root. This is used to create a
// snapshot of the Ctrie by copying the root I-node and setting it to a new
// generation.
func (c *Ctrie) rdcssRoot(old *iNode, expected *mainNode, nv *iNode) bool {
	desc := &iNode{
		rdcss: &rdcssDescriptor{
			old:      old,
			expected: expected,
			nv:       nv,
		},
	}
	if c.casRoot(old, desc) {
		c.rdcssComplete(false)
		return atomic.LoadInt32(&desc.rdcss.committed) == 1
	}
	return false
}

// rdcssComplete commits the RDCSS operation.
func (c *Ctrie) rdcssComplete(abort bool) *iNode {
	for {
		r := (*iNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&c.root))))
		if r.rdcss == nil {
			return r
		}

		var (
			desc = r.rdcss
			ov   = desc.old
			exp  = desc.expected
			nv   = desc.nv
		)

		if abort {
			if c.casRoot(r, ov) {
				return ov
			}
			continue
		}

		oldeMain := gcasRead(ov, c)
		if oldeMain == exp {
			// Commit the RDCSS.
			if c.casRoot(r, nv) {
				atomic.StoreInt32(&desc.committed, 1)
				return nv
			}
			continue
		}
		if c.casRoot(r, ov) {
			return ov
		}
		continue
	}
}

// casRoot performs a CAS on the Ctrie root.
func (c *Ctrie) casRoot(ov, nv *iNode) bool {
	c.assertReadWrite()
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&c.root)), unsafe.Pointer(ov), unsafe.Pointer(nv))
}

// #region list

var (
	// Empty is an empty PersistentList.
	Empty PersistentList = &emptyList{}

	// ErrEmptyList is returned when an invalid operation is performed on an
	// empty list.
	ErrEmptyList = errors.New("Empty list")
)

// PersistentList is an immutable, persistent linked list.
type PersistentList interface {
	// Head returns the head of the list. The bool will be false if the list is
	// empty.
	Head() (interface{}, bool)

	// Tail returns the tail of the list. The bool will be false if the list is
	// empty.
	Tail() (PersistentList, bool)

	// IsEmpty indicates if the list is empty.
	IsEmpty() bool

	// Length returns the number of items in the list.
	Length() uint

	// Add will add the item to the list, returning the new list.
	Add(head interface{}) PersistentList

	// Insert will insert the item at the given position, returning the new
	// list or an error if the position is invalid.
	Insert(val interface{}, pos uint) (PersistentList, error)

	// Get returns the item at the given position or an error if the position
	// is invalid.
	Get(pos uint) (interface{}, bool)

	// Remove will remove the item at the given position, returning the new
	// list or an error if the position is invalid.
	Remove(pos uint) (PersistentList, error)

	// Find applies the predicate function to the list and returns the first
	// item which matches.
	Find(func(interface{}) bool) (interface{}, bool)

	// FindIndex applies the predicate function to the list and returns the
	// index of the first item which matches or -1 if there is no match.
	FindIndex(func(interface{}) bool) int

	// Map applies the function to each entry in the list and returns the
	// resulting slice.
	Map(func(interface{}) interface{}) []interface{}
}

type emptyList struct{}

// Head returns the head of the list. The bool will be false if the list is
// empty.
func (e *emptyList) Head() (interface{}, bool) {
	return nil, false
}

// Tail returns the tail of the list. The bool will be false if the list is
// empty.
func (e *emptyList) Tail() (PersistentList, bool) {
	return nil, false
}

// IsEmpty indicates if the list is empty.
func (e *emptyList) IsEmpty() bool {
	return true
}

// Length returns the number of items in the list.
func (e *emptyList) Length() uint {
	return 0
}

// Add will add the item to the list, returning the new list.
func (e *emptyList) Add(head interface{}) PersistentList {
	return &list{head, e}
}

// Insert will insert the item at the given position, returning the new list or
// an error if the position is invalid.
func (e *emptyList) Insert(val interface{}, pos uint) (PersistentList, error) {
	if pos == 0 {
		return e.Add(val), nil
	}
	return nil, ErrEmptyList
}

// Get returns the item at the given position or an error if the position is
// invalid.
func (e *emptyList) Get(pos uint) (interface{}, bool) {
	return nil, false
}

// Remove will remove the item at the given position, returning the new list or
// an error if the position is invalid.
func (e *emptyList) Remove(pos uint) (PersistentList, error) {
	return nil, ErrEmptyList
}

// Find applies the predicate function to the list and returns the first item
// which matches.
func (e *emptyList) Find(func(interface{}) bool) (interface{}, bool) {
	return nil, false
}

// FindIndex applies the predicate function to the list and returns the index
// of the first item which matches or -1 if there is no match.
func (e *emptyList) FindIndex(func(interface{}) bool) int {
	return -1
}

// Map applies the function to each entry in the list and returns the resulting
// slice.
func (e *emptyList) Map(func(interface{}) interface{}) []interface{} {
	return nil
}

type list struct {
	head interface{}
	tail PersistentList
}

// Head returns the head of the list. The bool will be false if the list is
// empty.
func (l *list) Head() (interface{}, bool) {
	return l.head, true
}

// Tail returns the tail of the list. The bool will be false if the list is
// empty.
func (l *list) Tail() (PersistentList, bool) {
	return l.tail, true
}

// IsEmpty indicates if the list is empty.
func (l *list) IsEmpty() bool {
	return false
}

// Length returns the number of items in the list.
func (l *list) Length() uint {
	curr := l
	length := uint(0)
	for {
		length += 1
		tail, _ := curr.Tail()
		if tail.IsEmpty() {
			return length
		}
		curr = tail.(*list)
	}
}

// Add will add the item to the list, returning the new list.
func (l *list) Add(head interface{}) PersistentList {
	return &list{head, l}
}

// Insert will insert the item at the given position, returning the new list or
// an error if the position is invalid.
func (l *list) Insert(val interface{}, pos uint) (PersistentList, error) {
	if pos == 0 {
		return l.Add(val), nil
	}
	nl, err := l.tail.Insert(val, pos-1)
	if err != nil {
		return nil, err
	}
	return nl.Add(l.head), nil
}

// Get returns the item at the given position or an error if the position is
// invalid.
func (l *list) Get(pos uint) (interface{}, bool) {
	if pos == 0 {
		return l.head, true
	}
	return l.tail.Get(pos - 1)
}

// Remove will remove the item at the given position, returning the new list or
// an error if the position is invalid.
func (l *list) Remove(pos uint) (PersistentList, error) {
	if pos == 0 {
		nl, _ := l.Tail()
		return nl, nil
	}

	nl, err := l.tail.Remove(pos - 1)
	if err != nil {
		return nil, err
	}
	return &list{l.head, nl}, nil
}

// Find applies the predicate function to the list and returns the first item
// which matches.
func (l *list) Find(pred func(interface{}) bool) (interface{}, bool) {
	if pred(l.head) {
		return l.head, true
	}
	return l.tail.Find(pred)
}

// FindIndex applies the predicate function to the list and returns the index
// of the first item which matches or -1 if there is no match.
func (l *list) FindIndex(pred func(interface{}) bool) int {
	curr := l
	idx := 0
	for {
		if pred(curr.head) {
			return idx
		}
		tail, _ := curr.Tail()
		if tail.IsEmpty() {
			return -1
		}
		curr = tail.(*list)
		idx += 1
	}
}

// Map applies the function to each entry in the list and returns the resulting
// slice.
func (l *list) Map(f func(interface{}) interface{}) []interface{} {
	return append(l.tail.Map(f), f(l.head))
}

// #endregion
