// Package dtrie provides an implementation of the dtrie data structure, which
// is a persistent hash trie that dynamically expands or shrinks to provide
// efficient memory allocation. This data structure is based on the papers
// Ideal Hash Trees by Phil Bagwell and Optimizing Hash-Array Mapped Tries for
// Fast and Lean Immutable JVM Collections by Michael J. Steindorfer and
// Jurgen J. Vinju
package dtrie

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// Dtrie is a persistent hash trie that dynamically expands or shrinks
// to provide efficient memory allocation.
type Dtrie struct {
	root   *node
	hasher func(v interface{}) uint32
}

type entry struct {
	hash  uint32
	key   interface{}
	value interface{}
}

func (e *entry) KeyHash() uint32 {
	return e.hash
}

func (e *entry) Key() interface{} {
	return e.key
}

func (e *entry) Value() interface{} {
	return e.value
}

// New creates an empty DTrie with the given hashing function.
// If nil is passed in, the default hashing function will be used.
func New(hasher func(v interface{}) uint32) *Dtrie {
	if hasher == nil {
		hasher = defaultHasher
	}
	return &Dtrie{
		root:   emptyNode(0, 32),
		hasher: hasher,
	}
}

// Size returns the number of entries in the Dtrie.
func (d *Dtrie) Size() (size int) {
	for _ = range iterate(d.root, nil) {
		size++
	}
	return size
}

// Get returns the value for the associated key or returns nil if the
// key does not exist.
func (d *Dtrie) Get(key interface{}) interface{} {
	node := get(d.root, d.hasher(key), key)
	if node != nil {
		return node.Value()
	}
	return nil
}

// Insert adds a key value pair to the Dtrie, replacing the existing value if
// the key already exists and returns the resulting Dtrie.
func (d *Dtrie) Insert(key, value interface{}) *Dtrie {
	root := insert(d.root, &entry{d.hasher(key), key, value})
	return &Dtrie{root, d.hasher}
}

// Remove deletes the value for the associated key if it exists and returns
// the resulting Dtrie.
func (d *Dtrie) Remove(key interface{}) *Dtrie {
	root := remove(d.root, d.hasher(key), key)
	return &Dtrie{root, d.hasher}
}

// Iterator returns a read-only channel of Entries from the Dtrie. If a stop
// channel is provided, closing it will terminate and close the iterator
// channel. Note that if a cancel channel is not used and not every entry is
// read from the iterator, a goroutine will leak.
func (d *Dtrie) Iterator(stop <-chan struct{}) <-chan Entry {
	return iterate(d.root, stop)
}

func mask(hash uint32, level uint8) uint32 {
	return (hash >> (5 * level)) & 0x01f
}

func defaultHasher(value interface{}) uint32 {
	switch v := value.(type) {
	case uint8:
		return uint32(v)
	case uint16:
		return uint32(v)
	case uint32:
		return v
	case uint64:
		return uint32(v)
	case int8:
		return uint32(v)
	case int16:
		return uint32(v)
	case int32:
		return uint32(v)
	case int64:
		return uint32(v)
	case uint:
		return uint32(v)
	case int:
		return uint32(v)
	case uintptr:
		return uint32(v)
	case float32:
		return uint32(v)
	case float64:
		return uint32(v)
	}
	hasher := fnv.New32a()
	hasher.Write([]byte(fmt.Sprintf("%#v", value)))
	return hasher.Sum32()
}

type node struct {
	entries []Entry
	nodeMap Bitmap32
	dataMap Bitmap32
	level   uint8 // level starts at 0
}

func (n *node) KeyHash() uint32    { return 0 }
func (n *node) Key() interface{}   { return nil }
func (n *node) Value() interface{} { return nil }

func (n *node) String() string {
	return fmt.Sprint(n.entries)
}

type collisionNode struct {
	entries []Entry
}

func (n *collisionNode) KeyHash() uint32    { return 0 }
func (n *collisionNode) Key() interface{}   { return nil }
func (n *collisionNode) Value() interface{} { return nil }

func (n *collisionNode) String() string {
	return fmt.Sprintf("<COLLISIONS %v>%v", len(n.entries), n.entries)
}

// Entry defines anything held within the data structure
type Entry interface {
	KeyHash() uint32
	Key() interface{}
	Value() interface{}
}

func emptyNode(level uint8, capacity int) *node {
	return &node{entries: make([]Entry, capacity), level: level}
}

func insert(n *node, entry Entry) *node {
	index := uint(mask(entry.KeyHash(), n.level))
	newNode := n
	if newNode.level == 6 { // handle hash collisions on 6th level
		if newNode.entries[index] == nil {
			newNode.entries[index] = entry
			newNode.dataMap = newNode.dataMap.SetBit(index)
			return newNode
		}
		if newNode.dataMap.GetBit(index) {
			if newNode.entries[index].Key() == entry.Key() {
				newNode.entries[index] = entry
				return newNode
			}
			cNode := &collisionNode{entries: make([]Entry, 2)}
			cNode.entries[0] = newNode.entries[index]
			cNode.entries[1] = entry
			newNode.entries[index] = cNode
			newNode.dataMap = newNode.dataMap.ClearBit(index)
			return newNode
		}
		cNode := newNode.entries[index].(*collisionNode)
		cNode.entries = append(cNode.entries, entry)
		return newNode
	}
	if !newNode.dataMap.GetBit(index) && !newNode.nodeMap.GetBit(index) { // insert directly
		newNode.entries[index] = entry
		newNode.dataMap = newNode.dataMap.SetBit(index)
		return newNode
	}
	if newNode.nodeMap.GetBit(index) { // insert into sub-node
		newNode.entries[index] = insert(newNode.entries[index].(*node), entry)
		return newNode
	}
	if newNode.entries[index].Key() == entry.Key() {
		newNode.entries[index] = entry
		return newNode
	}
	// create new node with the new and existing entries
	var subNode *node
	if newNode.level == 5 { // only 2 bits left at level 6 (4 possible indices)
		subNode = emptyNode(newNode.level+1, 4)
	} else {
		subNode = emptyNode(newNode.level+1, 32)
	}
	subNode = insert(subNode, newNode.entries[index])
	subNode = insert(subNode, entry)
	newNode.dataMap = newNode.dataMap.ClearBit(index)
	newNode.nodeMap = newNode.nodeMap.SetBit(index)
	newNode.entries[index] = subNode
	return newNode
}

// returns nil if not found
func get(n *node, keyHash uint32, key interface{}) Entry {
	index := uint(mask(keyHash, n.level))
	if n.dataMap.GetBit(index) {
		return n.entries[index]
	}
	if n.nodeMap.GetBit(index) {
		return get(n.entries[index].(*node), keyHash, key)
	}
	if n.level == 6 { // get from collisionNode
		if n.entries[index] == nil {
			return nil
		}
		cNode := n.entries[index].(*collisionNode)
		for _, e := range cNode.entries {
			if e.Key() == key {
				return e
			}
		}
	}
	return nil
}

func remove(n *node, keyHash uint32, key interface{}) *node {
	index := uint(mask(keyHash, n.level))
	newNode := n
	if n.dataMap.GetBit(index) {
		newNode.entries[index] = nil
		newNode.dataMap = newNode.dataMap.ClearBit(index)
		return newNode
	}
	if n.nodeMap.GetBit(index) {
		subNode := newNode.entries[index].(*node)
		subNode = remove(subNode, keyHash, key)
		// compress if only 1 entry exists in sub-node
		if subNode.nodeMap.PopCount() == 0 && subNode.dataMap.PopCount() == 1 {
			var e Entry
			for i := uint(0); i < 32; i++ {
				if subNode.dataMap.GetBit(i) {
					e = subNode.entries[i]
					break
				}
			}
			newNode.entries[index] = e
			newNode.nodeMap = newNode.nodeMap.ClearBit(index)
			newNode.dataMap = newNode.dataMap.SetBit(index)
		}
		newNode.entries[index] = subNode
		return newNode
	}
	if n.level == 6 { // delete from collisionNode
		cNode := newNode.entries[index].(*collisionNode)
		for i, e := range cNode.entries {
			if e.Key() == key {
				cNode.entries = append(cNode.entries[:i], cNode.entries[i+1:]...)
				break
			}
		}
		// compress if only 1 entry exists in collisionNode
		if len(cNode.entries) == 1 {
			newNode.entries[index] = cNode.entries[0]
			newNode.dataMap = newNode.dataMap.SetBit(index)
		}
		return newNode
	}
	return n
}

func iterate(n *node, stop <-chan struct{}) <-chan Entry {
	out := make(chan Entry)
	go func() {
		defer close(out)
		pushEntries(n, stop, out)
	}()
	return out
}

func pushEntries(n *node, stop <-chan struct{}, out chan Entry) {
	var wg sync.WaitGroup
	for i, e := range n.entries {
		select {
		case <-stop:
			return
		default:
			index := uint(i)
			switch {
			case n.dataMap.GetBit(index):
				out <- e
			case n.nodeMap.GetBit(index):
				wg.Add(1)
				go func() {
					defer wg.Done()
					pushEntries(e.(*node), stop, out)
				}()
				wg.Wait()
			case n.level == 6 && e != nil:
				for _, ce := range n.entries[index].(*collisionNode).entries {
					select {
					case <-stop:
						return
					default:
						out <- ce
					}
				}
			}
		}
	}
}

// #region bitmap32

// #region bitmap

// Bitmap32 tracks 32 bool values within a uint32
type Bitmap32 uint32

// SetBit returns a Bitmap32 with the bit at the given position set to 1
func (b Bitmap32) SetBit(pos uint) Bitmap32 {
	return b | (1 << pos)
}

// ClearBit returns a Bitmap32 with the bit at the given position set to 0
func (b Bitmap32) ClearBit(pos uint) Bitmap32 {
	return b & ^(1 << pos)
}

// GetBit returns true if the bit at the given position in the Bitmap32 is 1
func (b Bitmap32) GetBit(pos uint) bool {
	return (b & (1 << pos)) != 0
}

// PopCount returns the amount of bits set to 1 in the Bitmap32
func (b Bitmap32) PopCount() int {
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	b -= (b >> 1) & 0x55555555
	b = (b>>2)&0x33333333 + b&0x33333333
	b += b >> 4
	b &= 0x0f0f0f0f
	b *= 0x01010101
	return int(byte(b >> 24))
}

// Bitmap64 tracks 64 bool values within a uint64
type Bitmap64 uint64

// SetBit returns a Bitmap64 with the bit at the given position set to 1
func (b Bitmap64) SetBit(pos uint) Bitmap64 {
	return b | (1 << pos)
}

// ClearBit returns a Bitmap64 with the bit at the given position set to 0
func (b Bitmap64) ClearBit(pos uint) Bitmap64 {
	return b & ^(1 << pos)
}

// GetBit returns true if the bit at the given position in the Bitmap64 is 1
func (b Bitmap64) GetBit(pos uint) bool {
	return (b & (1 << pos)) != 0
}

// PopCount returns the amount of bits set to 1 in the Bitmap64
func (b Bitmap64) PopCount() int {
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	b -= (b >> 1) & 0x5555555555555555
	b = (b>>2)&0x3333333333333333 + b&0x3333333333333333
	b += b >> 4
	b &= 0x0f0f0f0f0f0f0f0f
	b *= 0x0101010101010101
	return int(byte(b >> 56))
}

// #endregion

// #endregion
