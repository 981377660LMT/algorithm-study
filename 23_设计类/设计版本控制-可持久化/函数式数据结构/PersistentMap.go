// HashArrayMappedTrie(HAMT) 实现的 PersistentMap.
// https://github.com/shanzi/algo-ds/tree/master/hamt
// https://zthinker.com/archives/functional-go-intro-to-hamt
//
// 1.节点内部空间的压缩：bitMap 存储节点信息，onesCount表示节点数量
// 2.压缩 Trie 高度
//
// 查询操作
// 查询操作和在 Trie 上进行查询没有太大的改变，只是因为节点的空间被压缩之后，
// 需要先通过 Bit Map 字段 判断对应的位置是否存在子元素，如果存在，再通过 Pop Count 计算得到子元素所在的实际位置。
// 在获得到子元素之后，鉴于 HAMT 对 Trie 的高度进行了压缩，我们要先判断这个子元素是 Trie 的子节点还是值元素本身。
// 对于 Hash Table，我们将值包装成<Key, Value>这样的一个值节点存放进去。
// 这样取出时就可以先判断 Key 的值是否相同，以此决定是否真的取到了对象。
// 加入操作
// 加入操作类似于查询操作，对于给定的键值对，我们先计算键的 Hash 值，然后从高位开始利用 Hash 值。
// 沿着 Trie 树向下查找，此时存在四种情况：
// 查找到一个位置，这个位置所保存的值对象代表的键和新插入的键相同，则用新插入的值代替原来的值
// 查找到一个子节点，这个子节点再向下目标对象应该所在的位置没有元素，此时只需申请一个新的宽度比当前子节点多1的节点， 将原来子节点上的所有元素连带当前对象一一设置到新的子节点上，维护好元素的顺序和 Bit Map 的值即可
// 查找到一个子节点，这个子节点再向下目标对象所应该在的位置是一个值对象，这时，为这个值对象连带当前对象分配一个宽度为2 的子节点，将原来的值对象和当前对象放置在新的节点中，最后用新的节点在当前子节点中替换掉对应位置
// 查找到一个子节点，这个节点再向下目标对象所应该在的位置是另一个子节点，则在那个子节点上递归进行上述操作
// 删除操作
// 删除操作是加入操作的逆运算，其基本原理也类似。在删除时，现在 Trie 中找到对应的对象，然后依次递归地删除。
// 在从子节点删除一个元素时，我们也一样要分配一个新的子节点，这个子节点的宽度比原来少1。
// 之后将剩余的子元素按照规则放回新的子节点并维护 Bit Map 字段的值。最后用新的子节点替换原来的节点并递归向上进行。
// 一般来讲，在递归时有如下几个情况：
// 在子节点进行删除操作之后，子节点不再包含任何子元素，则删除当前子节点并返回空
// 当前子节点只剩下一个元素，且这个元素为值元素，则应该降低树的高度，直接返回值元素
// 当前子节点只剩下一个元素，且这个元素是子节点，则保留这个子节点，回到上一个节点

package main

import (
	"fmt"
	"math/bits"
	"runtime/debug"
	"time"
	"unsafe"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	demo()
}

func demo() {
	time1 := time.Now()
	mp := NewPersistentMap().WithMutations(func(t ITransientMap) {
		for i := 0; i < 1e6; i++ {
			t.Put(i, i)
			t.Get(i)
		}
	})
	time2 := time.Now()
	fmt.Println(time2.Sub(time1))
	fmt.Println(mp.Size())
}

// 1146. 快照数组
// https://leetcode.cn/problems/snapshot-array/
type SnapshotArray struct {
	mp  IPersistentMap
	git []IPersistentMap
}

func Constructor(length int) SnapshotArray {
	mp := NewPersistentMap().WithMutations(func(t ITransientMap) {
		for i := 0; i < length; i++ {
			t.Put(i, 0)
		}
	})
	return SnapshotArray{mp: mp}
}

func (this *SnapshotArray) Set(index int, val int) {
	this.mp = this.mp.Put(index, val)
}

func (this *SnapshotArray) Snap() int {
	this.git = append(this.git, this.mp)
	return len(this.git) - 1
}

func (this *SnapshotArray) Get(index int, snap_id int) int {
	res, _ := this.git[snap_id].Get(index)
	return res.(int)
}

type Key = int // string

const (
	seed0 uint32 = 0
	seed1 uint32 = 0x00663d81
	seed2 uint32 = 0x7fffffff
)

var globalOwnerId uint32

func nextId() uint32 {
	globalOwnerId++
	return globalOwnerId
	// return atomic.AddUint32(&id, 1)  // concurrent safe
}

// todo: has key to uint32
func keyHash32(s Key, seed uint32) uint32 {
	// h := hashString(([]byte)(s), seed) // hash string to uint32
	h := uint32(s) // hash int to uint32

	// !The most significant two bits will be abandoned during insert and lookup,
	// !Thus we'd better mix them down with the least significant two bits
	return ((h >> 30) & 3) ^ h
}

type IPersistentMap interface {
	Get(key Key) (interface{}, bool)
	Put(key Key, value interface{}) IPersistentMap
	Remove(key Key) (IPersistentMap, interface{})
	Size() int

	AsMutable() ITransientMap
	WithMutations(f func(ITransientMap)) IPersistentMap
}

type PersistentMap TransientMap

var EMPTY_PERSISTENT_MAP = &PersistentMap{}

func NewPersistentMap() IPersistentMap {
	return EMPTY_PERSISTENT_MAP
}

func (self *PersistentMap) Get(key Key) (interface{}, bool) {
	return (*TransientMap)(self).Get(key)
}

func (self *PersistentMap) Put(key Key, value interface{}) IPersistentMap {
	t := self.AsMutable()
	if ok := t.Put(key, value); ok {
		return t.AsImmutable()
	} else {
		return self
	}
}

func (self *PersistentMap) Remove(key Key) (IPersistentMap, interface{}) {
	t := self.AsMutable()
	if value, ok := t.Remove(key); ok {
		return t.AsImmutable(), value
	} else {
		return self, nil
	}
}

func (self *PersistentMap) AsMutable() ITransientMap {
	return &TransientMap{id: nextId(), size: self.size, root: self.root}
}

func (self *PersistentMap) WithMutations(f func(ITransientMap)) IPersistentMap {
	t := self.AsMutable()
	f(t)
	return t.AsImmutable()
}

func (self *PersistentMap) Size() int {
	return int(self.size)
}

type ITransientMap interface {
	Get(key Key) (interface{}, bool)
	Put(key Key, value interface{}) bool
	Remove(key Key) (interface{}, bool)
	Size() int

	AsImmutable() IPersistentMap
}

type TransientMap struct {
	id   uint32
	size int32
	root *node
}

func (self *TransientMap) Get(key Key) (interface{}, bool) {
	if self.root == nil {
		return nil, false
	}

	if value, ok := self.getWithHash(self.root, key); ok {
		return value, true
	}

	return nil, false
}

func (self *TransientMap) Put(key Key, value interface{}) bool {
	e := &entry{key: key, hash: 0, value: value}

	if self.root == nil {
		self.root = newNode(self.id, 0, 1)
	}

	if root, newkey, ok := self.putEntry(self.root, e, 0); ok {
		self.root = root
		if newkey {
			self.size += 1
		}
		return true
	}
	return false
}

func (self *TransientMap) Remove(key Key) (interface{}, bool) {
	if self.root == nil {
		return nil, false
	}

	if root, ent := self.removeEntry(self.root, key, 0, 0); ent != nil {
		self.root = root
		self.size -= 1
		return ent.value, true
	}

	return nil, false
}

func (self *TransientMap) AsImmutable() IPersistentMap {
	m := (*PersistentMap)(self)
	m.id = 0
	return m
}

func (self *TransientMap) Size() int {
	return int(self.size)
}

func (self *TransientMap) getWithHash(root *node, key Key) (interface{}, bool) {
	p := root
	var hash uint32
	for i := 0; i < 18; i++ {
		// At some specific depth, hash need to be recalculate
		switch i {
		case 0:
			hash = keyHash32(key, seed0)
		case 6:
			hash = keyHash32(key, seed1)
		case 12:
			hash = keyHash32(key, seed2)
		}

		d := uint32(i % 6)
		h := uint32((hash >> (d * 5)) & 0x1f)

		c := p.childAt(h)
		if c == nil {
			return nil, false
		}

		if e, ok := c.(*entry); ok {
			if e.key == key {
				// Find object
				return e.value, true
			} else {
				// No match, return nil
				return nil, false
			}
		}

		// c must be a node
		p = c.(*node)
	}

	// Nothing found after drained hash code,
	// return false indicating not found
	return nil, false
}

func (self *TransientMap) putEntry(root *node, e *entry, depth int) (*node, bool, bool) {
	var newkey bool
	// At some specific depth, hash need to be recalculate
	switch depth {
	case 0:
		e.hash = keyHash32(e.key, seed0)
	case 6:
		e.hash = keyHash32(e.key, seed1)
	case 12:
		e.hash = keyHash32(e.key, seed2)
	case 18:
		panic("Inresolvable hash collision!")
	}

	d := uint32(depth % 6)
	h := uint32((e.hash >> (d * 5)) & 0x1f)

	if !root.has(h) {
		// Found a position to put new item in
		return root.putChildAt(self.id, h, e), true, true
	}

	child := root.childAt(h)
	if subnode, ok := child.(*node); ok {
		// Found a sub node, recursively put entry
		if child, newkey, ok = self.putEntry(subnode, e, depth+1); ok {
			return root.putChildAt(self.id, h, child), newkey, ok
		} else {
			return root, false, false
		}
	}

	if olde, ok := child.(*entry); ok {
		// Found an entry
		if olde.key != e.key {
			// !Collision. Create a new node and rehash current entry
			subnode := newNode(self.id, 0, 0)
			subnode, _, _ = self.putEntry(subnode, olde, depth+1)

			if child, newkey, ok = self.putEntry(subnode, e, depth+1); ok {
				return root.putChildAt(self.id, h, child), newkey, ok
			} else {
				return root, false, false
			}
		}

		// Two keys are the same
		if olde.value == e.value {
			// Two values are the same, do nothing and return
			return root, false, false
		} else {
			// Two values are different, overwrite value
			return root.putChildAt(self.id, h, e), false, true
		}
	}

	assertUnreachable()
	return nil, false, false
}

func (self *TransientMap) removeEntry(root *node, key Key, hash uint32, depth int) (*node, *entry) {
	switch depth {
	case 0:
		hash = keyHash32(key, seed0)
	case 6:
		hash = keyHash32(key, seed1)
	case 12:
		hash = keyHash32(key, seed2)
	}

	d := uint32(depth % 6)
	h := uint32((hash >> (d * 5)) & 0x1f)

	if !root.has(h) {
		// The item to remove not found, do nothing
		return root, nil
	}

	child := root.childAt(h)
	if subnode, ok := child.(*node); ok {
		// Found a sub node, recursively remove entry
		if child, ent := self.removeEntry(subnode, key, hash, depth+1); ent != nil {
			switch child.size() {
			case 0:
				// child no longer contains anything, remove it from root
				return root.removeChildAt(self.id, h), ent
			case 1:
				if e, ok := child.children[0].(*entry); ok {
					// child only contains one entry, retrieve it and put it to the root
					// as to reduce height of the trie
					e.hash = hash
					return root.putChildAt(self.id, h, e), ent
				}
				// Although child only contains one item, but it's a node.
				// We choose not to compact the tree in this case as it'll be complex
				// and may increase time cost per remove operation
				fallthrough
			default:
				return root.putChildAt(self.id, h, child), ent
			}
		} else {
			// entry not found, do noting
			return root, nil
		}
	}

	if ent, ok := child.(*entry); ok {
		// Found the entry to remove. Remove it and return the removed entry
		return root.removeChildAt(self.id, h), ent
	}

	assertUnreachable()
	return nil, nil
}

type node struct {
	id       uint32
	mask     uint32
	children []interface{}
}

func newNode(id uint32, mask uint32, size int32) *node {
	cap := size + (size & 1)
	return &node{id, mask, make([](interface{}), size, cap)}
}

func (self *node) has(index uint32) bool {
	if index >= 0 && index < 32 {
		return ((self.mask >> index) & 1) == 1
	}
	return false
}

func (self *node) pos(index uint32) int {
	m := self.mask & ((1 << index) - 1)
	return bits.OnesCount32(m)
}

func (self *node) size() int {
	return len(self.children)
}

func (self *node) shrink() {
	// This method removes unnecessary empty slots in children
	// once the difference between cap and len of child are larger than
	// desired as to save memory
	if cap(self.children)-len(self.children) > 4 {
		size := self.size()
		children := make([](interface{}), size, size+(size&1))
		copy(children, self.children)
		self.children = children
	}
}

func (self *node) childAt(index uint32) interface{} {
	if self.has(index) {
		return self.children[self.pos(index)]
	}
	return nil
}

func (self *node) putChildAt(id uint32, index uint32, child interface{}) *node {
	pos := self.pos(index)
	if id == self.id {
		// Within the same transaction
		if !self.has(index) {
			children := append(self.children, nil)
			for i := len(children) - 1; i > pos; i-- {
				children[i] = children[i-1]
			}
			self.children = children
		}
		self.children[pos] = child
		self.mask |= (1 << index)
		return self
	} else {
		// Not within the same transaction, copy the node
		if self.has(index) {
			cloned := newNode(id, self.mask, int32(len(self.children)))
			copy(cloned.children, self.children)
			cloned.children[pos] = child
			return cloned
		} else {
			cloned := newNode(id, self.mask|(1<<index), int32(len(self.children)+1))
			for i := 0; i < pos; i++ {
				cloned.children[i] = self.children[i]
			}
			cloned.children[pos] = child
			for i := pos; i < len(self.children); i++ {
				cloned.children[i+1] = self.children[i]
			}
			return cloned
		}
	}
}

func (self *node) removeChildAt(id uint32, index uint32) *node {
	if !self.has(index) {
		return self
	}

	pos := self.pos(index)

	if id == self.id {
		// Within the same transaction
		children := self.children
		for i := pos + 1; i < len(children); i++ {
			children[i-1] = children[i]
		}
		// Clear reference to be gc friendly
		children[len(children)-1] = nil
		self.mask ^= (1 << index)
		self.children = children
		self.shrink()
		return self
	} else {
		// Not within the same transaction, copy the node
		cloned := newNode(id, self.mask^(1<<index), int32(len(self.children)-1))
		for i := 0; i < pos; i++ {
			cloned.children[i] = self.children[i]
		}

		for i := pos + 1; i < len(self.children); i++ {
			cloned.children[i-1] = self.children[i]
		}
		return cloned
	}
}

type entry struct {
	key   Key
	hash  uint32
	value interface{}
}

func (self *entry) Key() Key {
	return self.key
}

func (self *entry) Value() interface{} {
	return self.value
}

func assertUnreachable() {
	panic("Should be unreachable!")
}

// Murmur3 32-bit version
func hashString(data []byte, seed uint32) uint32 {
	var h1 = seed

	nblocks := len(data) / 4
	var p uintptr
	if len(data) > 0 {
		p = uintptr(unsafe.Pointer(&data[0]))
	}
	p1 := p + uintptr(4*nblocks)
	for ; p < p1; p += 4 {
		k1 := *(*uint32)(unsafe.Pointer(p))

		k1 *= 0xcc9e2d51
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= 0x1b873593

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19)
		h1 = h1*5 + 0xe6546b64
	}

	tail := data[nblocks*4:]

	var k1 uint32
	switch len(tail) & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= 0xcc9e2d51
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= 0x1b873593
		h1 ^= k1
	}

	h1 ^= uint32(len(data))

	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}
