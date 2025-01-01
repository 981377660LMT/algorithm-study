package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

func main() {
	// Create a new generic tree with int as Key and string as Value.
	tree := NewGTree[int, string]()

	// Insert some values into the tree.
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(3, "three")

	// Search for a value.
	if value, found := tree.Search(2); found {
		fmt.Printf("Found: %s\n", value) // Output: Found: two
	} else {
		fmt.Println("Not found")
	}

	// Delete a value.
	if value, deleted := tree.Delete(3); deleted {
		fmt.Printf("Deleted: %s\n", value) // Output: Deleted: three
	} else {
		fmt.Println("Not found for deletion")
	}

	// Check the size of the tree.
	fmt.Printf("Tree Size: %d\n", tree.Size()) // Output: Tree Size: 2

	// Traverse the tree using ForEach.
	tree.ForEach(func(node Node) bool {
		fmt.Printf("Node Key: %s, Node Value: %s\n", string(node.Key()), node.Value().(string))
		return true // Continue iteration
	}, TraverseLeaf)
}

// GTree is a generic tree that supports any type for keys and values.
type GTree[K comparable, V any] struct {
	tree Tree
}

// NewGTree creates a new generic adaptive radix tree.
func NewGTree[K comparable, V any]() *GTree[K, V] {
	return &GTree[K, V]{
		tree: NewTree(),
	}
}

// Insert a new key-value pair into the tree.
func (gt *GTree[K, V]) Insert(key K, value V) (oldValue V, updated bool, err error) {
	// Convert key to []byte
	keyBytes, err := convertKeyToBytes(key)
	if err != nil {
		return
	}

	oldVal, updated := gt.tree.Insert(Key(keyBytes), value)

	if oldVal != nil {
		oldValue = oldVal.(V)
	}
	return
}

// Delete removes a key from the tree.
func (gt *GTree[K, V]) Delete(key K) (value V, deleted bool) {
	keyBytes, err := convertKeyToBytes(key)
	if err != nil {
		return
	}

	val, deleted := gt.tree.Delete(Key(keyBytes))
	if val != nil {
		value = val.(V)
	}
	return
}

// Search for a key in the tree.
func (gt *GTree[K, V]) Search(key K) (value V, found bool) {
	keyBytes, err := convertKeyToBytes(key)
	if err != nil {
		return
	}

	val, found := gt.tree.Search(Key(keyBytes))
	if val != nil {
		value = val.(V)
	}
	return
}

// Size returns the number of elements in the tree.
func (gt *GTree[K, V]) Size() int {
	return gt.tree.Size()
}

// ForEach performs the given callback on each node.
func (gt *GTree[K, V]) ForEach(callback func(node Node) bool, options ...int) {
	gt.tree.ForEach(callback, options...)
}

// Helper function to convert an integer key to a byte slice.
func convertKeyToBytes[K comparable](key K) ([]byte, error) {
	switch v := any(key).(type) {
	case int:
		return []byte(fmt.Sprintf("%d", v)), nil // Simple conversion
	default:
		return nil, errors.New("unsupported key type")
	}
}

// #region api

// Key represents the type used for keys in the Adaptive Radix Tree.
// It can consist of any byte sequence, including Unicode characters and null bytes.
type Key []byte

// Value is an interface representing the value type stored in the tree.
// Any type of data can be stored as a Value.
type Value any

// Node types.
const (
	Leaf    Kind = 0
	Node4   Kind = 1
	Node16  Kind = 2
	Node48  Kind = 3
	Node256 Kind = 4
)

// Traverse Options.
const (
	// Iterate only over leaf nodes.
	TraverseLeaf = 1

	// Iterate only over non-leaf nodes.
	TraverseNode = 2

	// Iterate over all nodes in the tree.
	TraverseAll = TraverseLeaf | TraverseNode

	// Iterate in reverse order.
	TraverseReverse = 4
)

// These errors can be returned when iteration over the tree.
var (
	ErrConcurrentModification = errors.New("concurrent modification has been detected")
	ErrNoMoreNodes            = errors.New("there are no more nodes in the tree")
)

// Kind is a node type.
type Kind int

// String returns string representation of the Kind value.
func (k Kind) String() string {
	return []string{"Leaf", "Node4", "Node16", "Node48", "Node256"}[k]
}

// Callback defines the function type used during tree traversal.
// It is invoked for each node visited in the traversal.
// If the callback function returns false, the iteration is terminated early.
type Callback func(node Node) (cont bool)

// Node represents a node within the Adaptive Radix Tree.
type Node interface {
	// Kind returns the type of the node, distinguishing between leaf and internal nodes.
	Kind() Kind

	// Key returns the key associated with a leaf node.
	// This method should only be called on leaf nodes.
	// Calling this on a non-leaf node will return nil.
	Key() Key

	// Value returns the value stored in a leaf node.
	// This method should only be called on leaf nodes.
	// Calling this on a non-leaf node will return nil.
	Value() Value
}

// Iterator provides a mechanism to traverse nodes in key order within the tree.
type Iterator interface {
	// HasNext returns true if there are more nodes to visit during the iteration.
	// Use this method to check for remaining nodes before calling Next.
	HasNext() bool

	// Next returns the next node in the iteration and advances the iterator's position.
	// If the iteration has no more nodes, it returns ErrNoMoreNodes error.
	// Ensure you call HasNext before invoking Next to avoid errors.
	// If the tree has been structurally modified since the iterator was created,
	// it returns an ErrConcurrentModification error.
	Next() (Node, error)
}

// Tree is an Adaptive Radix Tree interface.
type Tree interface {
	// Insert adds a new key-value pair into the tree.
	// If the key already exists in the tree, it updates its value and returns the old value along with true.
	// If the key is new, it returns nil and false.
	Insert(key Key, value Value) (oldValue Value, updated bool)

	// Delete removes the specified key and its associated value from the tree.
	// If the key is found and deleted, it returns the removed value and true.
	// If the key does not exist, it returns nil and false.
	Delete(key Key) (value Value, deleted bool)

	// Search retrieves the value associated with the specified key in the tree.
	// If the key exists, it returns the value and true.
	// If the key does not exist, it returns nil and false.
	Search(key Key) (value Value, found bool)

	// ForEach iterates over all the nodes in the tree, invoking a provided callback function for each node.
	// By default, it processes leaf nodes in ascending order.
	// The iteration can be customized using options:
	// - Pass TraverseReverse to iterate over nodes in descending order.
	// The iteration stops if the callback function returns false, allowing for early termination.
	ForEach(callback Callback, options ...int)

	// ForEachPrefix iterates over all leaf nodes whose keys start with the specified keyPrefix,
	// invoking a provided callback function for each matching node.
	// By default, the iteration processes nodes in ascending order.
	// Use the TraverseReverse option to iterate over nodes in descending order.
	// Iteration stops if the callback function returns false, allowing for early termination.
	ForEachPrefix(keyPrefix Key, callback Callback, options ...int)

	// Iterator returns an iterator for traversing leaf nodes in the tree.
	// By default, the iteration occurs in ascending order.
	// To traverse nodes in reverse (descending) order, pass the TraverseReverse option.
	Iterator(options ...int) Iterator

	// Minimum retrieves the leaf node with the smallest key in the tree.
	// If such a leaf is found, it returns its value and true.
	// If the tree is empty, it returns nil and false.
	Minimum() (Value, bool)

	// Maximum retrieves the leaf node with the largest key in the tree.
	// If such a leaf is found, it returns its value and true.
	// If the tree is empty, it returns nil and false.
	Maximum() (Value, bool)

	// Size returns the number of key-value pairs stored in the tree.
	Size() int
}

// NewTree creates a new adaptive radix tree.
func NewTree() Tree {
	return newTree()
}

// #endregion

// #region consts
// node constraints.
const (
	node4Min = 2 // minimum number of children for node4.
	node4Max = 4 // maximum number of children for node4.

	node16Min = node4Max + 1 // minimum number of children for node16.
	node16Max = 16           // maximum number of children for node16.

	node48Min = node16Max + 1 // minimum number of children for node48.
	node48Max = 48            // maximum number of children for node48.

	node256Min = node48Max + 1 // minimum number of children for node256.
	node256Max = 256           // maximum number of children for node256.
)

const (
	// maxPrefixLen is maximum prefix length for internal nodes.
	maxPrefixLen = 10
)

// #endregion

// #region factory

// nodeFactory is an interface for creating various types of ART nodes,
// including nodes with different capacities and leaf nodes.
type nodeFactory interface {
	newNode4() *nodeRef
	newNode16() *nodeRef
	newNode48() *nodeRef
	newNode256() *nodeRef

	newLeaf(key Key, value interface{}) *nodeRef
}

// make sure that objFactory implements all methods of nodeFactory interface.
var _ nodeFactory = &objFactory{}

//nolint:gochecknoglobals
var (
	factory = newObjFactory()
)

// newTree creates a new tree.
func newTree() *tree {
	return &tree{
		version: 0,
		root:    nil,
		size:    0,
	}
}

// objFactory implements nodeFactory interface.
type objFactory struct{}

// newObjFactory creates a new objFactory.
func newObjFactory() nodeFactory {
	return &objFactory{}
}

// Simple obj factory implementation.
func (f *objFactory) newNode4() *nodeRef {
	return &nodeRef{
		kind: Node4,
		ref:  unsafe.Pointer(new(node4)), //#nosec:G103
	}
}

// newNode16 creates a new node16 as a nodeRef.
func (f *objFactory) newNode16() *nodeRef {
	return &nodeRef{
		kind: Node16,
		ref:  unsafe.Pointer(new(node16)), //#nosec:G103
	}
}

// newNode48 creates a new node48 as a nodeRef.
func (f *objFactory) newNode48() *nodeRef {
	return &nodeRef{
		kind: Node48,
		ref:  unsafe.Pointer(new(node48)), //#nosec:G103
	}
}

// newNode256 creates a new node256 as a nodeRef.
func (f *objFactory) newNode256() *nodeRef {
	return &nodeRef{
		kind: Node256,
		ref:  unsafe.Pointer(new(node256)), //#nosec:G103
	}
}

// newLeaf creates a new leaf node as a nodeRef.
// It clones the key to avoid any source key mutation.
func (f *objFactory) newLeaf(key Key, value interface{}) *nodeRef {
	keyClone := make(Key, len(key))
	copy(keyClone, key)

	return &nodeRef{
		kind: Leaf,
		ref: unsafe.Pointer(&leaf{ //#nosec:G103
			key:   keyClone,
			value: value,
		}),
	}
}

// #endregion

// #region nodeRef

// indexNotFound is a special index value
// that indicates that the index is not found.
const indexNotFound = -1

// nodeNotFound is a special node pointer
// that indicates that the node is not found
// for different internal tree operations.
var nodeNotFound *nodeRef //nolint:gochecknoglobals

// nodeRef stores all available tree nodes leaf and nodeX types
// as a ref to *unsafe* pointer.
// The kind field is used to determine the type of the node.
type nodeRef struct {
	ref  unsafe.Pointer
	kind Kind
}

type nodeLeafer interface {
	minimum() *leaf
	maximum() *leaf
}

type nodeSizeManager interface {
	hasCapacityForChild() bool
	grow() *nodeRef

	isReadyToShrink() bool
	shrink() *nodeRef
}

type nodeOperations interface {
	addChild(kc keyChar, child *nodeRef)
	deleteChild(kc keyChar) int
}

type nodeChildren interface {
	childAt(idx int) **nodeRef
	allChildren() []*nodeRef
}

type nodeKeyIndexer interface {
	index(kc keyChar) int
}

// noder is an interface that defines methods that
// must be implemented by nodeRef and all node types.
// extra interfaces are used to group methods by their purpose
// and help with code readability.
type noder interface {
	nodeLeafer
	nodeOperations
	nodeChildren
	nodeKeyIndexer
	nodeSizeManager
}

// toNode converts the nodeRef to specific node type.
// the idea is to avoid type assertion in the code in multiple places.
func toNode(nr *nodeRef) noder {
	if nr == nil {
		return noopNoder
	}

	switch nr.kind { //nolint:exhaustive
	case Node4:
		return nr.node4()
	case Node16:
		return nr.node16()
	case Node48:
		return nr.node48()
	case Node256:
		return nr.node256()
	default:
		return noopNoder
	}
}

// noop is a no-op noder implementation.
type noop struct{}

func (*noop) minimum() *leaf             { return nil }
func (*noop) maximum() *leaf             { return nil }
func (*noop) index(keyChar) int          { return indexNotFound }
func (*noop) childAt(int) **nodeRef      { return &nodeNotFound }
func (*noop) allChildren() []*nodeRef    { return nil }
func (*noop) hasCapacityForChild() bool  { return true }
func (*noop) grow() *nodeRef             { return nil }
func (*noop) isReadyToShrink() bool      { return false }
func (*noop) shrink() *nodeRef           { return nil }
func (*noop) addChild(keyChar, *nodeRef) {}
func (*noop) deleteChild(keyChar) int    { return 0 }

// noopNoder is the default Noder implementation.
var noopNoder noder = &noop{} //nolint:gochecknoglobals

// assert that all node types implement noder interface.
var _ noder = (*node4)(nil)
var _ noder = (*node16)(nil)
var _ noder = (*node48)(nil)
var _ noder = (*node256)(nil)

// assert that nodeRef implements public Node interface.
var _ Node = (*nodeRef)(nil)

// Kind returns the node kind.
func (nr *nodeRef) Kind() Kind {
	return nr.kind
}

// Key returns the node key for leaf nodes.
// for nodeX types, it returns nil.
func (nr *nodeRef) Key() Key {
	if nr.isLeaf() {
		return nr.leaf().key
	}

	return nil
}

// Value returns the node value for leaf nodes.
// for nodeX types, it returns nil.
func (nr *nodeRef) Value() Value {
	if nr.isLeaf() {
		return nr.leaf().value
	}

	return nil
}

// isLeaf returns true if the node is a leaf node.
func (nr *nodeRef) isLeaf() bool {
	return nr.kind == Leaf
}

// setPrefix sets the node prefix with the new prefix and prefix length.
func (nr *nodeRef) setPrefix(newPrefix []byte, prefixLen int) {
	n := nr.node()

	n.prefixLen = uint16(prefixLen) //#nosec:G115
	for i := 0; i < minInt(prefixLen, maxPrefixLen); i++ {
		n.prefix[i] = newPrefix[i]
	}
}

// minimum returns itself if the node is a leaf node.
// otherwise it returns the minimum leaf node under the current node.
func (nr *nodeRef) minimum() *leaf {
	if nr.kind == Leaf {
		return nr.leaf()
	}

	return toNode(nr).minimum()
}

// maximum returns itself if the node is a leaf node.
// otherwise it returns the maximum leaf node under the current node.
func (nr *nodeRef) maximum() *leaf {
	if nr.kind == Leaf {
		return nr.leaf()
	}

	return toNode(nr).maximum()
}

// findChildByKey returns the child node reference for the given key.
func (nr *nodeRef) findChildByKey(key Key, keyOffset int) **nodeRef {
	n := toNode(nr)
	idx := n.index(key.charAt(keyOffset))

	return n.childAt(idx)
}

// nodeX/leaf casts the nodeRef to the specific nodeX/leaf type.
func (nr *nodeRef) node() *node       { return (*node)(nr.ref) }    // node casts nodeRef to node.
func (nr *nodeRef) node4() *node4     { return (*node4)(nr.ref) }   // node4 casts nodeRef to node4.
func (nr *nodeRef) node16() *node16   { return (*node16)(nr.ref) }  // node16 casts nodeRef to node16.
func (nr *nodeRef) node48() *node48   { return (*node48)(nr.ref) }  // node48 casts nodeRef to node48.
func (nr *nodeRef) node256() *node256 { return (*node256)(nr.ref) } // node256 casts nodeRef to node256.
func (nr *nodeRef) leaf() *leaf       { return (*leaf)(nr.ref) }    // leaf casts nodeRef to leaf.

// addChild adds a new child node to the current node.
// If the node is full, it grows to the next node type.
func (nr *nodeRef) addChild(kc keyChar, child *nodeRef) {
	n := toNode(nr)

	if n.hasCapacityForChild() {
		n.addChild(kc, child)
	} else {
		bigNode := n.grow()         // grow to the next node type
		bigNode.addChild(kc, child) // recursively add the child to the new node
		replaceNode(nr, bigNode)    // replace the current node with the new node
	}
}

// deleteChild deletes the child node from the current node.
// If the node can shrink after, it shrinks to the previous node type.
func (nr *nodeRef) deleteChild(kc keyChar) bool {
	shrank := false
	n := toNode(nr)
	n.deleteChild(kc)

	if n.isReadyToShrink() {
		shrank = true
		smallNode := n.shrink()    // shrink to the previous node type
		replaceNode(nr, smallNode) // replace the current node with the shrank node
	}

	return shrank
}

// match finds the first mismatched index between
// the node's prefix and the specified key prefix.
// This approach efficiently identifies the mismatch by
// leveraging the node's existing prefix data.
func (nr *nodeRef) match(key Key, keyOffset int) int /* 1st mismatch index*/ {
	// calc the remaining key length from offset
	keyRemaining := len(key) - keyOffset
	if keyRemaining < 0 {
		return 0
	}

	n := nr.node()

	// the maximum length we can check against the node's prefix
	maxPrefixLen := minInt(int(n.prefixLen), maxPrefixLen)
	limit := minInt(maxPrefixLen, keyRemaining)

	// compare the key against the node's prefix
	for i := 0; i < limit; i++ {
		if n.prefix[i] != key[keyOffset+i] {
			return i
		}
	}

	return limit
}

// matchDeep returns the first index where the key mismatches,
// starting with the node's prefix(see match) and continuing with the minimum leaf's key.
// It returns the mismatch index or matches up to the key's end.
func (nr *nodeRef) matchDeep(key Key, keyOffset int) int /* mismatch index*/ {
	mismatchIdx := nr.match(key, keyOffset)
	if mismatchIdx < maxPrefixLen {
		return mismatchIdx
	}

	leafKey := nr.minimum().key
	limit := minInt(len(leafKey), len(key)) - keyOffset

	for ; mismatchIdx < limit; mismatchIdx++ {
		if leafKey[keyOffset+mismatchIdx] != key[keyOffset+mismatchIdx] {
			break
		}
	}

	return mismatchIdx
}

// #endregion

// #region nodeLeaf
// Leaf node stores the key-value pair.
type leaf struct {
	key   Key
	value interface{}
}

// match returns true if the leaf node's key matches the given key.
func (l *leaf) match(key Key) bool {
	return len(l.key) == len(key) && bytes.Equal(l.key, key)
}

// prefixMatch returns true if the leaf node's key has the given key as a prefix.
func (l *leaf) prefixMatch(key Key) bool {
	if key == nil || len(l.key) < len(key) {
		return false
	}

	return bytes.Equal(l.key[:len(key)], key)
}

// #endregion

// #region node
// prefix used in the node to store the key prefix.
// it is used to improve leaf key comparison performance.
type prefix [maxPrefixLen]byte

// node is the base struct for all node types.
// it contains the common fields for all nodeX types.
type node struct {
	prefix      prefix // prefix of the node
	prefixLen   uint16 // length of the prefix
	childrenLen uint16 // number of children in the node4, node16, node48, node256
}

// replaceRef is used to replace node in-place by updating the reference.
func replaceRef(oldNode **nodeRef, newNode *nodeRef) {
	*oldNode = newNode
}

// replaceNode is used to replace node in-place by updating the node.
func replaceNode(oldNode *nodeRef, newNode *nodeRef) {
	*oldNode = *newNode
}

// #endregion

// #region node4

// node4 represents a node with 4 children.
type node4 struct {
	node
	children [node4Max + 1]*nodeRef // pointers to the child nodes, +1 is for the zero byte child
	keys     [node4Max]byte         // keys for the children
	present  [node4Max]byte         // present bits for the keys
}

// minimum returns the minimum leaf node.
func (n *node4) minimum() *leaf {
	return nodeMinimum(n.children[:])
}

// maximum returns the maximum leaf node.
func (n *node4) maximum() *leaf {
	return nodeMaximum(n.children[:n.childrenLen])
}

// index returns the index of the given character.
func (n *node4) index(kc keyChar) int {
	if kc.invalid {
		return node4Max
	}

	return findIndex(n.keys[:n.childrenLen], kc.ch)
}

// childAt returns the child at the given index.
func (n *node4) childAt(idx int) **nodeRef {
	if idx < 0 || idx >= len(n.children) {
		return &nodeNotFound
	}

	return &n.children[idx]
}

func (n *node4) allChildren() []*nodeRef {
	return n.children[:]
}

// hasCapacityForChild returns true if the node has room for more children.
func (n *node4) hasCapacityForChild() bool {
	return n.childrenLen < node4Max
}

// grow converts the node4 into the node16.
func (n *node4) grow() *nodeRef {
	an16 := factory.newNode16()
	n16 := an16.node16()

	copyNode(&n16.node, &n.node)
	n16.children[node16Max] = n.children[node4Max] // copy zero byte child

	for i := 0; i < int(n.childrenLen); i++ {
		// skip if the key is not present
		if n.present[i] == 0 {
			continue
		}

		// copy elements from n4 to n16 to the last position
		n16.insertChildAt(i, n.keys[i], n.children[i])
	}

	return an16
}

// isReadyToShrink returns true if the node is under-utilized and ready to shrink.
func (n *node4) isReadyToShrink() bool {
	// we have to return the number of children for the current node(node4) as
	// `node.numChildren` plus one if zero node is not nil.
	// For all higher nodes(16/48/256) we simply copy zero node to a smaller node
	// see deleteChild() and shrink() methods for implementation details
	numChildren := n.childrenLen
	if n.children[node4Max] != nil {
		numChildren++
	}

	return numChildren < node4Min
}

// shrink converts the node4 into the leaf node or a node with fewer children.
func (n *node4) shrink() *nodeRef {
	// Select the non-nil child node
	var nonNilChild *nodeRef
	if n.children[0] != nil {
		nonNilChild = n.children[0]
	} else {
		nonNilChild = n.children[node4Max]
	}

	// if the only child is a leaf node, return it
	if nonNilChild.isLeaf() {
		return nonNilChild
	}

	// update the prefix of the child node
	n.adjustPrefix(nonNilChild.node())

	return nonNilChild
}

// adjustPrefix handles prefix adjustments for a non-leaf child.
func (n *node4) adjustPrefix(childNode *node) {
	nodePrefLen := int(n.prefixLen)

	// at this point, the node has only one child
	// copy the key part of the current node as prefix
	if nodePrefLen < maxPrefixLen {
		n.prefix[nodePrefLen] = n.keys[0]
		nodePrefLen++
	}

	// copy the part of child prefix that fits into the current node
	if nodePrefLen < maxPrefixLen {
		childPrefLen := minInt(int(childNode.prefixLen), maxPrefixLen-nodePrefLen)
		copy(n.prefix[nodePrefLen:], childNode.prefix[:childPrefLen])
		nodePrefLen += childPrefLen
	}

	// copy the part of the current node prefix that fits into the child node
	prefixLen := minInt(nodePrefLen, maxPrefixLen)
	copy(childNode.prefix[:], n.prefix[:prefixLen])
	childNode.prefixLen += n.prefixLen + 1
}

// addChild adds a new child to the node.
func (n *node4) addChild(kc keyChar, child *nodeRef) {
	pos := n.findInsertPos(kc)
	n.makeRoom(pos)
	n.insertChildAt(pos, kc.ch, child)
}

// find the insert position for the new child.
func (n *node4) findInsertPos(kc keyChar) int {
	if kc.invalid {
		return node4Max
	}

	numChildren := int(n.childrenLen)
	for i := 0; i < numChildren; i++ {
		if n.keys[i] > kc.ch {
			return i
		}
	}

	return numChildren
}

// makeRoom creates space for the new child by shifting the elements to the right.
func (n *node4) makeRoom(pos int) {
	if pos < 0 || pos >= int(n.childrenLen) {
		return
	}

	for i := int(n.childrenLen); i > pos; i-- {
		n.keys[i] = n.keys[i-1]
		n.present[i] = n.present[i-1]
		n.children[i] = n.children[i-1]
	}
}

// insertChildAt inserts the child at the given position.
func (n *node4) insertChildAt(pos int, ch byte, child *nodeRef) {
	if pos == node4Max {
		n.children[pos] = child
	} else {
		n.keys[pos] = ch
		n.present[pos] = 1
		n.children[pos] = child
		n.childrenLen++
	}
}

// deleteChild deletes the child from the node.
func (n *node4) deleteChild(kc keyChar) int {
	if kc.invalid {
		// clear the zero byte child reference
		n.children[node4Max] = nil
	} else if idx := n.index(kc); idx >= 0 {
		n.deleteChildAt(idx)
		n.clearLastElement()
	}

	// we have to return the number of children for the current node(node4) as
	// `n.numChildren` plus one if null node is not nil.
	// `Shrink` method can be invoked after this method,
	// `Shrink` can convert this node into a leaf node type.
	// For all higher nodes(16/48/256) we simply copy null node to a smaller node
	// see deleteChild() and shrink() methods for implementation details
	numChildren := int(n.childrenLen)
	if n.children[node4Max] != nil {
		numChildren++
	}

	return numChildren
}

// deleteChildAt deletes the child at the given index
// by shifting the elements to the left to overwrite deleted child.
func (n *node4) deleteChildAt(idx int) {
	for i := idx; i < int(n.childrenLen) && i+1 < node4Max; i++ {
		n.keys[i] = n.keys[i+1]
		n.present[i] = n.present[i+1]
		n.children[i] = n.children[i+1]
	}

	n.childrenLen--
}

// clearLastElement clears the last element in the node.
func (n *node4) clearLastElement() {
	lastIdx := int(n.childrenLen)
	n.keys[lastIdx] = 0
	n.present[lastIdx] = 0
	n.children[lastIdx] = nil
}

// #endregion

// #region node16

// present16 is a bitfield to store the presence of keys in the node16.
// node16 needs 16 bits to store the presence of keys.
type present16 uint16

func (p present16) hasChild(idx int) bool {
	return p&(1<<idx) != 0
}

func (p *present16) setAt(idx int) {
	*p |= 1 << idx
}

func (p *present16) clearAt(idx int) {
	*p &= ^(1 << idx)
}

func (p *present16) shiftRight(idx int) {
	p.clearAt(idx)
	*p |= ((*p & (1 << (idx - 1))) << 1)
}

func (p *present16) shiftLeft(idx int) {
	p.clearAt(idx)
	*p |= ((*p & (1 << (idx + 1))) >> 1)
}

// node16 represents a node with 16 children.
type node16 struct {
	node
	children [node16Max + 1]*nodeRef // +1 is for the zero byte child
	keys     [node16Max]byte
	present  present16
}

// minimum returns the minimum leaf node.
func (n *node16) minimum() *leaf {
	return nodeMinimum(n.children[:])
}

// maximum returns the maximum leaf node.
func (n *node16) maximum() *leaf {
	return nodeMaximum(n.children[:n.childrenLen])
}

// index returns the child index for the given key.
func (n *node16) index(kc keyChar) int {
	if kc.invalid {
		return node16Max
	}

	return findIndex(n.keys[:n.childrenLen], kc.ch)
}

// childAt returns the child at the given index.
func (n *node16) childAt(idx int) **nodeRef {
	if idx < 0 || idx >= len(n.children) {
		return &nodeNotFound
	}

	return &n.children[idx]
}

func (n *node16) allChildren() []*nodeRef {
	return n.children[:]
}

// hasCapacityForChild returns true if the node has room for more children.
func (n *node16) hasCapacityForChild() bool {
	return n.childrenLen < node16Max
}

// grow converts the node to a node48.
func (n *node16) grow() *nodeRef {
	an48 := factory.newNode48()
	n48 := an48.node48()

	copyNode(&n48.node, &n.node)
	n48.children[node48Max] = n.children[node16Max] // copy zero byte child

	for numChildren, i := 0, 0; i < int(n.childrenLen); i++ {
		if !n.hasChild(i) {
			continue // skip if the key is not present
		}

		n48.insertChildAt(numChildren, n.keys[i], n.children[i])

		numChildren++
	}

	return an48
}

// caShrinkNode returns true if the node can be shriken.
func (n *node16) isReadyToShrink() bool {
	return n.childrenLen < node16Min
}

// shrink converts the node16 into the node4.
func (n *node16) shrink() *nodeRef {
	an4 := factory.newNode4()
	n4 := an4.node4()

	copyNode(&n4.node, &n.node)
	n4.children[node4Max] = n.children[node16Max]

	for i := 0; i < node4Max; i++ {
		n4.keys[i] = n.keys[i]

		if n.hasChild(i) {
			n4.present[i] = 1
		}

		n4.children[i] = n.children[i]
		if n4.children[i] != nil {
			n4.childrenLen++
		}
	}

	return an4
}

func (n *node16) hasChild(idx int) bool {
	return n.present.hasChild(idx)
}

// addChild adds a new child to the node.
func (n *node16) addChild(kc keyChar, child *nodeRef) {
	pos := n.findInsertPos(kc)
	n.makeRoom(pos)
	n.insertChildAt(pos, kc.ch, child)
}

// find the insert position for the new child.
func (n *node16) findInsertPos(kc keyChar) int {
	if kc.invalid {
		return node16Max
	}

	for i := 0; i < int(n.childrenLen); i++ {
		if n.keys[i] > kc.ch {
			return i
		}
	}

	return int(n.childrenLen)
}

// makeRoom makes room for a new child at the given position.
func (n *node16) makeRoom(pos int) {
	if pos < 0 || pos >= int(n.childrenLen) {
		return
	}

	// Shift keys and children to the right starting from the position
	copy(n.keys[pos+1:], n.keys[pos:int(n.childrenLen)])
	copy(n.children[pos+1:], n.children[pos:int(n.childrenLen)])

	for i := int(n.childrenLen); i > pos; i-- {
		n.present.shiftRight(i)
	}
}

// insertChildAt inserts a new child at the given position.
func (n *node16) insertChildAt(pos int, ch byte, child *nodeRef) {
	if pos < 0 || pos > node16Max {
		return
	}

	if pos == node16Max {
		n.children[pos] = child
	} else {
		n.keys[pos] = ch
		n.present.setAt(pos)
		n.children[pos] = child
		n.childrenLen++
	}
}

// deleChild removes a child from the node.
func (n *node16) deleteChild(kc keyChar) int {
	if kc.invalid {
		// clear the zero byte child reference
		n.children[node16Max] = nil
	} else if idx := n.index(kc); idx >= 0 {
		n.deleteChildAt(idx)
		n.clearLastElement()
	}

	return int(n.childrenLen)
}

// deleteChildAt removes a child at the given position.
func (n *node16) deleteChildAt(idx int) {
	childrenLen := int(n.childrenLen)
	if idx >= childrenLen {
		return
	}

	// Shift keys and children to the left, overwriting the deleted index
	copy(n.keys[idx:], n.keys[idx+1:childrenLen])
	copy(n.children[idx:], n.children[idx+1:childrenLen])

	// shift elements to the left to fill the gap
	for i := idx; i < childrenLen && i+1 < node16Max; i++ {
		n.present.shiftLeft(i)
	}

	n.childrenLen--
}

// clearLastElement clears the last element in the node.
func (n *node16) clearLastElement() {
	lastIdx := int(n.childrenLen)
	n.keys[lastIdx] = 0
	n.present.clearAt(lastIdx)
	n.children[lastIdx] = nil
}

// #endregion
// #region node48

// Node with 48 children.
const (
	n48bitShift = 6  // 2^n48bitShift == n48maskLen
	n48maskLen  = 64 // it should be sizeof(node48.present[0])
)

// present48 is a bitfield to store the presence of keys in the node48.
// It is a bitfield of 256 bits, so it is stored in 4 uint64.
type present48 [4]uint64

func (p present48) hasChild(ch int) bool {
	return p[ch>>n48bitShift]&(1<<(ch%n48maskLen)) != 0
}

func (p *present48) setAt(ch int) {
	(*p)[ch>>n48bitShift] |= (1 << (ch % n48maskLen))
}

func (p *present48) clearAt(ch int) {
	(*p)[ch>>n48bitShift] &= ^(1 << (ch % n48maskLen))
}

type node48 struct {
	node
	children [node48Max + 1]*nodeRef // +1 is for the zero byte child
	keys     [node256Max]byte
	present  present48 // need 256 bits for keys
}

// minimum returns the minimum leaf node.
func (n *node48) minimum() *leaf {
	if n.children[node48Max] != nil {
		return n.children[node48Max].minimum()
	}

	idx := 0
	for !n.hasChild(idx) {
		idx++
	}

	if n.children[n.keys[idx]] != nil {
		return n.children[n.keys[idx]].minimum()
	}

	return nil
}

// maximum returns the maximum leaf node.
func (n *node48) maximum() *leaf {
	idx := node256Max - 1
	for !n.hasChild(idx) {
		idx--
	}

	return n.children[n.keys[idx]].maximum()
}

// index returns the index of the child with the given key.
func (n *node48) index(kc keyChar) int {
	if kc.invalid {
		return node48Max
	}

	if n.hasChild(int(kc.ch)) {
		idx := int(n.keys[kc.ch])
		if idx < node48Max && n.children[idx] != nil {
			return idx
		}
	}

	return indexNotFound
}

// childAt returns the child at the given index.
func (n *node48) childAt(idx int) **nodeRef {
	if idx < 0 || idx >= len(n.children) {
		return &nodeNotFound
	}

	return &n.children[idx]
}

func (n *node48) allChildren() []*nodeRef {
	return n.children[:]
}

// hasCapacityForChild returns true if the node has room for more children.
func (n *node48) hasCapacityForChild() bool {
	return n.childrenLen < node48Max
}

// grow converts the node to a node256.
func (n *node48) grow() *nodeRef {
	an256 := factory.newNode256()
	n256 := an256.node256()

	copyNode(&n256.node, &n.node)
	n256.children[node256Max] = n.children[node48Max] // copy zero byte child

	for i := 0; i < node256Max; i++ {
		if n.hasChild(i) {
			n256.addChild(keyChar{ch: byte(i)}, n.children[n.keys[i]])
		}
	}

	return an256
}

// isReadyToShrink returns true if the node can be shrunk to a smaller node type.
func (n *node48) isReadyToShrink() bool {
	return n.childrenLen < node48Min
}

// shrink converts the node to a node16.
func (n *node48) shrink() *nodeRef {
	an16 := factory.newNode16()
	n16 := an16.node16()

	copyNode(&n16.node, &n.node)
	n16.children[node16Max] = n.children[node48Max]
	numChildren := 0

	for i, idx := range n.keys {
		if !n.hasChild(i) {
			continue // skip if the key is not present
		}

		child := n.children[idx]
		if child == nil {
			continue // skip if the child is nil
		}

		// copy elements from n48 to n16 to the last position
		n16.insertChildAt(numChildren, byte(i), child)

		numChildren++
	}

	return an16
}

func (n *node48) hasChild(idx int) bool {
	return n.present.hasChild(idx)
}

// addChild adds a new child to the node.
func (n *node48) addChild(kc keyChar, child *nodeRef) {
	pos := n.findInsertPos(kc)
	n.insertChildAt(pos, kc.ch, child)
}

// find the insert position for the new child.
func (n *node48) findInsertPos(kc keyChar) int {
	if kc.invalid {
		return node48Max
	}

	var i int
	for i < node48Max && n.children[i] != nil {
		i++
	}

	return i
}

// insertChildAt inserts a child at the given position.
func (n *node48) insertChildAt(pos int, ch byte, child *nodeRef) {
	if pos == node48Max {
		// insert the child at the zero byte child reference
		n.children[node48Max] = child
	} else {
		// insert the child at the given index
		n.keys[ch] = byte(pos)
		n.present.setAt(int(ch))
		n.children[pos] = child
		n.childrenLen++
	}
}

// deleteChild removes the child with the given key.
func (n *node48) deleteChild(kc keyChar) int {
	if kc.invalid {
		// clear the zero byte child reference
		n.children[node48Max] = nil
	} else if idx := n.index(kc); idx >= 0 && n.children[idx] != nil {
		// clear the child at the given index
		n.keys[kc.ch] = 0
		n.present.clearAt(int(kc.ch))
		n.children[idx] = nil
		n.childrenLen--
	}

	return int(n.childrenLen)
}

// #endregion

// #region node256

// Node with 256 children.
type node256 struct {
	node
	children [node256Max + 1]*nodeRef // +1 is for the zero byte child
}

// minimum returns the minimum leaf node.
func (n *node256) minimum() *leaf {
	return nodeMinimum(n.children[:])
}

// maximum returns the maximum leaf node.
func (n *node256) maximum() *leaf {
	return nodeMaximum(n.children[:node256Max])
}

// index returns the index of the child with the given key.
func (n *node256) index(kc keyChar) int {
	if kc.invalid { // handle zero byte in the key
		return node256Max
	}

	return int(kc.ch)
}

// childAt returns the child at the given index.
func (n *node256) childAt(idx int) **nodeRef {
	if idx < 0 || idx >= len(n.children) {
		return &nodeNotFound
	}

	return &n.children[idx]
}

func (n *node256) allChildren() []*nodeRef {
	return n.children[:]
}

// addChild adds a new child to the node.
func (n *node256) addChild(kc keyChar, child *nodeRef) {
	if kc.invalid {
		// handle zero byte in the key
		n.children[node256Max] = child
	} else {
		// insert new child
		n.children[kc.ch] = child
		n.childrenLen++
	}
}

// hasCapacityForChild for node256 always returns true.
func (n *node256) hasCapacityForChild() bool {
	return true
}

// grow for node256 always returns nil,
// because node256 has the maximum capacity.
func (n *node256) grow() *nodeRef {
	return nil
}

// isReadyToShrink returns true if the node can be shrunk.
func (n *node256) isReadyToShrink() bool {
	return n.childrenLen < node256Min
}

// shrink shrinks the node to a smaller type.
func (n *node256) shrink() *nodeRef {
	an48 := factory.newNode48()
	n48 := an48.node48()

	copyNode(&n48.node, &n.node)
	n48.children[node48Min] = n.children[node256Max] // copy zero byte child

	for numChildren, i := 0, 0; i < node256Max; i++ {
		if n.children[i] == nil {
			continue // skip if the child is nil
		}
		// copy elements from n256 to n48 to the last position
		n48.insertChildAt(numChildren, byte(i), n.children[i])

		numChildren++
	}

	return an48
}

// deleteChild removes the child with the given key.
func (n *node256) deleteChild(kc keyChar) int {
	if kc.invalid {
		// clear the zero byte child reference
		n.children[node256Max] = nil
	} else if idx := n.index(kc); n.children[idx] != nil {
		// clear the child at the given index
		n.children[idx] = nil
		n.childrenLen--
	}

	return int(n.childrenLen)
}

// #endregion

// #region tree

// treeOpResult represents the result of the tree operation.
type treeOpResult int

const (
	// treeOpNoChange indicates that the key was not found.
	treeOpNoChange treeOpResult = iota

	// treeOpInserted indicates that the key/value was inserted.
	treeOpInserted

	// treeOpUpdated indicates that the existing key was updated with a new value.
	treeOpUpdated

	// treeOpDeleted indicates that the key was deleted.
	treeOpDeleted
)

// keyChar stores the key character and an flag
// to indicate if the key char is invalid.
type keyChar struct {
	ch      byte
	invalid bool
}

// singleton keyChar instance to indicate
// that the key char is invalid.
//
//nolint:gochecknoglobals
var keyCharInvalid = keyChar{ch: 0, invalid: true}

// charAt returns the character at the given index.
// If the index is out of bounds, it returns 0 and false.
func (k Key) charAt(idx int) keyChar {
	if k.isValid(idx) {
		return keyChar{ch: k[idx]}
	}

	return keyCharInvalid
}

// isValid checks if the given index is within the bounds of the key.
func (k Key) isValid(idx int) bool {
	return idx >= 0 && idx < len(k)
}

// tree is the main data structure of the ART tree.
type tree struct {
	version int      // version is used to detect concurrent modifications
	size    int      // size is the number of elements in the tree
	root    *nodeRef // root is the root node of the tree
}

// make sure that tree implements all methods from the Tree interface.
var _ Tree = (*tree)(nil)

// Insert inserts the given key and value into the tree.
// If the key already exists, it updates the value and
// returns the old value with second return value set to true.
func (tr *tree) Insert(key Key, value Value) (Value, bool) {
	oldVal, status := tr.insertRecursively(&tr.root, key, value, 0)
	if status == treeOpInserted {
		tr.version++
		tr.size++
	}

	return oldVal, status == treeOpUpdated
}

// Delete deletes the given key from the tree.
func (tr *tree) Delete(key Key) (Value, bool) {
	val, status := tr.deleteRecursively(&tr.root, key, 0)
	if status == treeOpDeleted {
		tr.version++
		tr.size--

		return val, true
	}

	return nil, false
}

// Search searches for the given key in the tree.
func (tr *tree) Search(key Key) (Value, bool) {
	keyOffset := 0

	current := tr.root
	for current != nil {
		if current.isLeaf() {
			leaf := current.leaf()
			if leaf.match(key) {
				return leaf.value, true
			}

			return nil, false
		}

		curNode := current.node()
		if curNode.prefixLen > 0 {
			prefixLen := current.match(key, keyOffset)
			if prefixLen != minInt(int(curNode.prefixLen), maxPrefixLen) {
				return nil, false
			}

			keyOffset += int(curNode.prefixLen)
		}

		next := current.findChildByKey(key, keyOffset)
		if *next != nil {
			current = *next
		} else {
			current = nil
		}

		keyOffset++
	}

	return nil, false
}

// Minimum returns the minimum key in the tree.
func (tr *tree) Minimum() (Value, bool) {
	if tr == nil || tr.root == nil {
		return nil, false
	}

	return tr.root.minimum().value, true
}

// Maximum returns the maximum key in the tree.
func (tr *tree) Maximum() (Value, bool) {
	if tr == nil || tr.root == nil {
		return nil, false
	}

	return tr.root.maximum().value, true
}

// Size returns the number of elements in the tree.
func (tr *tree) Size() int {
	if tr == nil || tr.root == nil {
		return 0
	}

	return tr.size
}

// ForEach iterates over all keys in the tree and calls the callback function.
func (tr *tree) ForEach(callback Callback, opts ...int) {
	options := traverseOptions(opts...)
	tr.forEachRecursively(tr.root, traverseFilter(options, callback), options.hasReverse())
}

// ForEachPrefix iterates over all keys with the given prefix.
func (tr *tree) ForEachPrefix(key Key, callback Callback, opts ...int) {
	options := mergeOptions(opts...)
	tr.forEachPrefix(key, callback, options)
}

// Iterator returns a new tree iterator.
func (tr *tree) Iterator(opts ...int) Iterator {
	return newTreeIterator(tr, traverseOptions(opts...))
}

// String returns tree in the human readable format, see DumpNode for examples.
func (tr *tree) String() string {
	return DumpNode(tr.root)
}

// #endregion

// #region treeTraverse

// traverseAction is an action to be taken during tree traversal.
type traverseAction int

const (
	traverseStop     traverseAction = iota // traverseStop stops the tree traversal.
	traverseContinue                       // traverseContinue continues the tree traversal.
)

// traverseFunc defines the function for tree traversal.
// It returns the index of the next child node to traverse.
// The second return value indicates whether there are more child nodes to traverse.
type traverseFunc func() (int, bool)

// noopTraverseFunc is a no-op function for tree traversal.
func noopTraverseFunc() (int, bool) {
	return 0, false
}

// traverseOpts defines the options for tree traversal.
type traverseOpts int

func (opts traverseOpts) hasLeaf() bool {
	return opts&TraverseLeaf == TraverseLeaf
}

func (opts traverseOpts) hasNode() bool {
	return opts&TraverseNode == TraverseNode
}

func (opts traverseOpts) hasAll() bool {
	return opts&TraverseAll == TraverseAll
}

func (opts traverseOpts) hasReverse() bool {
	return opts&TraverseReverse == TraverseReverse
}

// traverseContext is a context for traversing nodes with 4, 16, or 256 children.
type traverseContext struct {
	numChildren   int
	zeroChildDone bool
	curChildIdx   int
}

// ascTraversal traverses the children in ascending order.
func (ctx *traverseContext) ascTraversal() (int, bool) {
	if !ctx.zeroChildDone {
		ctx.zeroChildDone = true

		return ctx.numChildren, true
	}

	idx := ctx.curChildIdx
	ctx.curChildIdx++

	return idx, idx < ctx.numChildren
}

// descTraversal traverses the children in descending order.
func (ctx *traverseContext) descTraversal() (int, bool) {
	if ctx.curChildIdx >= 0 {
		idx := ctx.curChildIdx
		ctx.curChildIdx--

		return idx, true
	}

	if !ctx.zeroChildDone {
		ctx.zeroChildDone = true

		return ctx.numChildren, true
	}

	return 0, false
}

// newTraverseGenericFunc creates a new traverseFunc for nodes with 4, 16, or 256 children.
// The reverse parameter indicates whether to traverse the children in reverse order.
func newTraverseGenericFunc(numChildren int, reverse bool) traverseFunc {
	ctx := &traverseContext{
		numChildren:   numChildren,
		zeroChildDone: false,
		curChildIdx:   ternary(reverse, numChildren-1, 0),
	}

	return ternary(reverse, ctx.descTraversal, ctx.ascTraversal)
}

// traverse48Context is a context for traversing nodes with 48 children.
type traverse48Context struct {
	curKeyIdx     int
	curKeyCh      byte
	zeroChildDone bool
	n48           *node48
}

// ascTraversal traverses the children in ascending order.
func (ctx *traverse48Context) ascTraversal() (int, bool) {
	if !ctx.zeroChildDone {
		ctx.zeroChildDone = true

		return node48Max, true
	}

	for ; ctx.curKeyIdx < node256Max; ctx.curKeyIdx++ {
		if ctx.n48.hasChild(ctx.curKeyIdx) {
			ctx.curKeyCh = ctx.n48.keys[ctx.curKeyIdx]
			ctx.curKeyIdx++

			return int(ctx.curKeyCh), true
		}
	}

	return 0, false
}

// descTraversal traverses the children in descending order.
func (ctx *traverse48Context) descTraversal() (int, bool) {
	for ; ctx.curKeyIdx > 0; ctx.curKeyIdx-- {
		if ctx.n48.hasChild(ctx.curKeyIdx) {
			ctx.curKeyCh = ctx.n48.keys[ctx.curKeyIdx]
			ctx.curKeyIdx--

			return int(ctx.curKeyCh), true
		}
	}

	if !ctx.zeroChildDone {
		ctx.zeroChildDone = true

		return node48Max, true
	}

	return 0, false
}

// newTraverse48Func creates a new traverseFunc for nodes with 48 children.
// The reverse parameter indicates whether to traverse the children in reverse order.
func newTraverse48Func(n48 *node48, reverse bool) traverseFunc {
	ctx := &traverse48Context{
		curKeyIdx: ternary(reverse, node256Max-1, 0),
		n48:       n48,
	}

	return ternary(reverse, ctx.descTraversal, ctx.ascTraversal)
}

func newTraverseFunc(n *nodeRef, reverse bool) traverseFunc {
	if n == nil {
		return noopTraverseFunc
	}

	switch n.kind { //nolint:exhaustive
	case Node4:
		return newTraverseGenericFunc(node4Max, reverse)
	case Node16:
		return newTraverseGenericFunc(node16Max, reverse)
	case Node48:
		return newTraverse48Func(n.node48(), reverse)
	case Node256:
		return newTraverseGenericFunc(node256Max, reverse)
	default:
		return noopTraverseFunc
	}
}

func mergeOptions(options ...int) int {
	opts := 0
	for _, opt := range options {
		opts |= opt
	}

	return opts
}

func traverseOptions(options ...int) traverseOpts {
	opts := mergeOptions(options...)

	typeOpts := opts & TraverseAll
	if typeOpts == 0 {
		typeOpts = TraverseLeaf // By default filter only leafs
	}

	orderOpts := opts & TraverseReverse

	return traverseOpts(typeOpts | orderOpts)
}

func traverseFilter(opts traverseOpts, callback Callback) Callback {
	if opts.hasAll() {
		return callback
	}

	return func(node Node) bool {
		if opts.hasLeaf() && node.Kind() == Leaf {
			return callback(node)
		}

		if opts.hasNode() && node.Kind() != Leaf {
			return callback(node)
		}

		return true
	}
}

func (tr *tree) forEachRecursively(current *nodeRef, callback Callback, reverse bool) traverseAction {
	if current == nil {
		return traverseContinue
	}

	if !callback(current) {
		return traverseStop
	}

	nextFn := newTraverseFunc(current, reverse)
	children := toNode(current).allChildren()

	return tr.traverseChildren(nextFn, children, callback, reverse)
}

func (tr *tree) traverseChildren(nextFn traverseFunc, children []*nodeRef, cb Callback, reverse bool) traverseAction {
	for {
		idx, hasMore := nextFn()
		if !hasMore {
			break
		}

		if child := children[idx]; child != nil {
			if tr.forEachRecursively(child, cb, reverse) == traverseStop {
				return traverseStop
			}
		}
	}

	return traverseContinue
}

func (tr *tree) forEachPrefix(key Key, callback Callback, opts int) traverseAction {
	opts &= (TraverseLeaf | TraverseReverse) // keep only leaf and reverse options

	tr.ForEach(func(n Node) bool {
		current, ok := n.(*nodeRef)
		if !ok {
			return false
		}

		if leaf := current.leaf(); leaf.prefixMatch(key) {
			return callback(current)
		}

		return true
	}, opts)

	return traverseContinue
}

// #endregion
// #region treeInsert

// insertRecursively inserts a new key-value pair into the tree.
// nrp means Node Reference Pointer.
func (tr *tree) insertRecursively(nrp **nodeRef, key Key, value Value, keyOffset int) (Value, treeOpResult) {
	nr := *nrp
	if nr == nil {
		return tr.insertNewLeaf(nrp, key, value)
	}

	if nr.isLeaf() {
		return tr.handleLeafInsertion(nrp, key, value, keyOffset)
	}

	return tr.handleNodeInsertion(nrp, key, value, keyOffset)
}

func (tr *tree) insertNewLeaf(nrp **nodeRef, key Key, value Value) (Value, treeOpResult) {
	replaceRef(nrp, factory.newLeaf(key, value))

	return nil, treeOpInserted
}

func (tr *tree) handleLeafInsertion(nrp **nodeRef, key Key, value Value, keyOffset int) (Value, treeOpResult) {
	nr := *nrp

	if leaf := nr.leaf(); leaf.match(key) {
		oldValue := leaf.value
		leaf.value = value

		return oldValue, treeOpUpdated
	}

	// Insert a new leaf by splitting
	// the old leaf to a node4 and adding the new leaf
	return tr.splitLeaf(nrp, key, value, keyOffset)
}

func (tr *tree) splitLeaf(nrpCurLeaf **nodeRef, key Key, value Value, keyOffset int) (Value, treeOpResult) {
	nrCurLeaf := *nrpCurLeaf
	curLeaf := nrCurLeaf.leaf()

	keysLCP := findLongestCommonPrefix(curLeaf.key, key, keyOffset)

	// Create a new node4 with the longest common prefix
	// between the old leaf and the new leaf key.
	nr4 := factory.newNode4()
	nr4.setPrefix(key[keyOffset:], keysLCP)
	keyOffset += keysLCP

	// branch by the first differing character
	// add the old leaf and the new leaf as children
	// to a newly created node4.
	nr4.addChild(curLeaf.key.charAt(keyOffset), nrCurLeaf)           // old leaf
	nr4.addChild(key.charAt(keyOffset), factory.newLeaf(key, value)) // new leaf

	// replace the old leaf with the new node4
	replaceRef(nrpCurLeaf, nr4)

	return nil, treeOpInserted
}

func (tr *tree) handleNodeInsertion(nrp **nodeRef, key Key, value Value, keyOffset int) (Value, treeOpResult) {
	nr := *nrp

	n := nr.node()
	if n.prefixLen > 0 {
		prefixMismatchIdx := nr.matchDeep(key, keyOffset)
		if prefixMismatchIdx < int(n.prefixLen) {
			return tr.splitNode(nrp, key, value, keyOffset, prefixMismatchIdx)
		}

		keyOffset += int(n.prefixLen)
	}

	return tr.continueInsertion(nrp, key, value, keyOffset)
}

func (tr *tree) splitNode(nrp **nodeRef, key Key, value Value, keyOffset int, mismatchIdx int) (Value, treeOpResult) {
	nr := *nrp
	n := nr.node()

	nr4 := factory.newNode4()
	nr4.setPrefix(n.prefix[:], mismatchIdx)

	tr.reassignPrefix(nr4, nr, key, value, keyOffset, mismatchIdx)

	replaceRef(nrp, nr4)

	return nil, treeOpInserted
}

func (tr *tree) reassignPrefix(newNRP *nodeRef, curNRP *nodeRef, key Key, value Value, keyOffset int, mismatchIdx int) {
	curNode := curNRP.node()
	curNode.prefixLen -= uint16(mismatchIdx + 1) //#nosec:G115

	idx := keyOffset + mismatchIdx

	// Adjust prefix and add children
	leaf := curNRP.minimum()
	newNRP.addChild(leaf.key.charAt(idx), curNRP)

	for i := 0; i < minInt(int(curNode.prefixLen), maxPrefixLen); i++ {
		curNode.prefix[i] = leaf.key[keyOffset+mismatchIdx+i+1]
	}

	// Insert the new leaf
	newNRP.addChild(key.charAt(idx), factory.newLeaf(key, value))
}

func (tr *tree) continueInsertion(nrp **nodeRef, key Key, value Value, keyOffset int) (Value, treeOpResult) {
	nr := *nrp

	nextNRP := nr.findChildByKey(key, keyOffset)
	if *nextNRP != nil {
		// Found a partial match, continue inserting
		return tr.insertRecursively(nextNRP, key, value, keyOffset+1)
	}

	// No child found, create a new leaf node
	nr.addChild(key.charAt(keyOffset), factory.newLeaf(key, value))

	return nil, treeOpInserted
}

// #endregion

// #region treeDelete

// deleteRecursively removes a node associated with the key from the tree.
func (tr *tree) deleteRecursively(nrp **nodeRef, key Key, keyOffset int) (Value, treeOpResult) {
	if tr == nil || *nrp == nil || len(key) == 0 {
		return nil, treeOpNoChange
	}

	nr := *nrp
	if nr.isLeaf() {
		return tr.handleLeafDeletion(nrp, key)
	}

	return tr.handleInternalNodeDeletion(nr, key, keyOffset)
}

// handleLeafDeletion removes a leaf node associated with the key from the tree.
func (tr *tree) handleLeafDeletion(nrp **nodeRef, key Key) (Value, treeOpResult) {
	if leaf := (*nrp).leaf(); leaf.match(key) {
		replaceRef(nrp, nil)

		return leaf.value, treeOpDeleted
	}

	return nil, treeOpNoChange
}

// handleInternalNodeDeletion removes a node associated with the key from the node.
func (tr *tree) handleInternalNodeDeletion(nr *nodeRef, key Key, keyOffset int) (Value, treeOpResult) {
	n := nr.node()

	if n.prefixLen > 0 {
		if mismatchIdx := nr.match(key, keyOffset); mismatchIdx != minInt(int(n.prefixLen), maxPrefixLen) {
			return nil, treeOpNoChange
		}

		keyOffset += int(n.prefixLen)
	}

	next := nr.findChildByKey(key, keyOffset)
	if *next == nil {
		return nil, treeOpNoChange
	}

	if (*next).isLeaf() {
		return tr.handleDeletionInChild(nr, *next, key, keyOffset)
	}

	return tr.deleteRecursively(next, key, keyOffset+1)
}

// handleDeletionInChild removes a leaf node from the child node.
func (tr *tree) handleDeletionInChild(curNR, nextNR *nodeRef, key Key, keyOffset int) (Value, treeOpResult) {
	leaf := (*nextNR).leaf()
	if !leaf.match(key) {
		return nil, treeOpNoChange
	}

	curNR.deleteChild(key.charAt(keyOffset))

	return leaf.value, treeOpDeleted
}

// #endregion

// #region treeIterator

// state represents the iteration state during tree traversal.
type state struct {
	items []*iteratorContext
}

// push adds a new iterator context to the state.
func (s *state) push(ctx *iteratorContext) {
	s.items = append(s.items, ctx)
}

// current returns the current iterator context and a flag indicating if there is any.
func (s *state) current() (*iteratorContext, bool) {
	if len(s.items) == 0 {
		return nil, false
	}

	return s.items[len(s.items)-1], true
}

// discard removes the last iterator context from the state.
func (s *state) discard() {
	if len(s.items) == 0 {
		return
	}

	s.items = s.items[:len(s.items)-1]
}

// iteratorContext represents the context of the tree iterator for one node.
type iteratorContext struct {
	nextChildFn traverseFunc
	children    []*nodeRef
}

// newIteratorContext creates a new iterator context for the given node.
func newIteratorContext(nr *nodeRef, reverse bool) *iteratorContext {
	return &iteratorContext{
		nextChildFn: newTraverseFunc(nr, reverse),
		children:    toNode(nr).allChildren(),
	}
}

// next returns the next node reference and a flag indicating if there are more nodes.
func (ic *iteratorContext) next() (*nodeRef, bool) {
	for {
		idx, ok := ic.nextChildFn()
		if !ok {
			break
		}

		if child := ic.children[idx]; child != nil {
			return child, true
		}
	}

	return nil, false
}

// iterator is a struct for tree traversal iteration.
type iterator struct {
	version  int      // tree version at the time of iterator creation
	tree     *tree    // tree to iterate
	state    *state   // iteration state
	nextNode *nodeRef // next node to iterate
	reverse  bool     // indicates if the iteration is in reverse order
}

// assert that iterator implements the Iterator interface.
var _ Iterator = (*iterator)(nil)

// newTreeIterator creates a new tree iterator.
func newTreeIterator(tr *tree, opts traverseOpts) Iterator {
	state := &state{}
	state.push(newIteratorContext(tr.root, opts.hasReverse()))

	it := &iterator{
		version:  tr.version,
		tree:     tr,
		nextNode: tr.root,
		state:    state,
		reverse:  opts.hasReverse(),
	}

	if opts&TraverseAll == TraverseAll {
		return it
	}

	bit := &bufferedIterator{
		opts: opts,
		it:   it,
	}

	// peek the first node or leaf
	bit.peek()

	return bit
}

// hasConcurrentModification checks if the tree has been modified concurrently.
func (it *iterator) hasConcurrentModification() bool {
	return it.version != it.tree.version
}

// HasNext returns true if there are more nodes to iterate.
func (it *iterator) HasNext() bool {
	return it.nextNode != nil
}

// Next returns the next node and an error if any.
// It returns ErrNoMoreNodes if there are no more nodes to iterate.
// It returns ErrConcurrentModification if the tree has been modified concurrently.
func (it *iterator) Next() (Node, error) {
	if !it.HasNext() {
		return nil, ErrNoMoreNodes
	}

	if it.hasConcurrentModification() {
		return nil, ErrConcurrentModification
	}

	current := it.nextNode
	it.next()

	return current, nil
}

// next moves the iterator to the next node.
func (it *iterator) next() {
	for {
		ctx, ok := it.state.current()
		if !ok {
			it.nextNode = nil // no more nodes to iterate

			return
		}

		nextNode, hasMore := ctx.next()
		if hasMore {
			it.nextNode = nextNode
			it.state.push(newIteratorContext(nextNode, it.reverse))

			return
		}

		it.state.discard() // discard the current context as exhausted
	}
}

// BufferedIterator implements HasNext and Next methods for buffered iteration.
// It allows to iterate over leaf or non-leaf nodes only.
type bufferedIterator struct {
	opts     traverseOpts
	it       Iterator
	nextNode Node
	nextErr  error
}

// HasNext returns true if there are more nodes to iterate.
func (bit *bufferedIterator) HasNext() bool {
	return bit.nextNode != nil
}

// Next returns the next node or leaf node and an error if any.
// ErrNoMoreNodes is returned if there are no more nodes to iterate.
// ErrConcurrentModification is returned if the tree has been modified concurrently.
func (bit *bufferedIterator) Next() (Node, error) {
	current := bit.nextNode

	if !bit.HasNext() {
		return nil, bit.nextErr
	}

	bit.peek()

	// ErrConcurrentModification should be returned immediately.
	// ErrNoMoreNodes will be return on the next call.
	if errors.Is(bit.nextErr, ErrConcurrentModification) {
		return nil, bit.nextErr
	}

	return current, nil
}

// hasLeafIterator checks if the iterator is for leaf nodes.
func (bit *bufferedIterator) hasLeafIterator() bool {
	return bit.opts&TraverseLeaf == TraverseLeaf
}

// hasNodeIterator checks if the iterator is for non-leaf nodes.
func (bit *bufferedIterator) hasNodeIterator() bool {
	return bit.opts&TraverseNode == TraverseNode
}

// peek looks for the next node or leaf node to iterate.
func (bit *bufferedIterator) peek() {
	for {
		bit.nextNode, bit.nextErr = bit.it.Next()
		if bit.nextErr != nil {
			return
		}

		if bit.matchesFilter() {
			return
		}
	}
}

// matchesFilter checks if the next node matches the iterator filter.
func (bit *bufferedIterator) matchesFilter() bool {
	// check if the iterator is looking for leaf nodes
	if bit.hasLeafIterator() && bit.nextNode.Kind() == Leaf {
		return true
	}

	// check if the iterator is looking for non-leaf nodes
	if bit.hasNodeIterator() && bit.nextNode.Kind() != Leaf {
		return true
	}

	return false
}

// #endregion

// #region treeDump

const (
	printValuesAsChar = 1 << iota
	printValuesAsDecimal
	printValuesAsHex

	printValueDefault = printValuesAsChar
)

// refFormatter is a function that formats an artNodeRef.
type refFormatter func(*dumpNodeRef) string

// RefFullFormatter returns the full address of the node, including the ID and the pointer.
func RefFullFormatter(a *dumpNodeRef) string {
	if a.ptr == nil {
		return "-"
	}

	return fmt.Sprintf("#%d/%p", a.id, a.ptr)
}

// RefShortFormatter returns only the ID of the node.
func RefShortFormatter(a *dumpNodeRef) string {
	if a.ptr == nil {
		return "-"
	}

	return fmt.Sprintf("#%d", a.id)
}

// RefAddrFormatter returns only the pointer address of the node (legacy).
func RefAddrFormatter(a *dumpNodeRef) string {
	if a.ptr == nil {
		return "-"
	}

	return fmt.Sprintf("%p", a.ptr)
}

// dumpNodeRef represents the address of a nodeRef in the tree,
// composed of a unique, sequential ID and a pointer to the node.
// The ID remains consistent for trees built with the same keys
// while the pointer may change with each build.
// This structure helps identify and compare nodes across different tree instances.
// It is also helpful for debugging and testing.
//
// For example: if you inserted the same keys in two different trees (or rerun the same test),
// you can compare the nodes of the two trees by their IDs.
// The IDs will be the same for the same keys, but the pointers will be different.
type dumpNodeRef struct {
	id  int          // unique ID
	ptr *nodeRef     // pointer to the node
	fmt refFormatter // function to format the address
}

// String returns the string representation of the address.
func (a dumpNodeRef) String() string {
	if a.fmt == nil {
		return RefFullFormatter(&a)
	}

	return a.fmt(&a)
}

// NodeRegistry maintains a mapping between nodeRef pointers and their unique IDs.
type nodeRegistry struct {
	ptrToID   map[*nodeRef]int // Maps a node pointer to its unique ID
	addresses []dumpNodeRef    // List of node references
	formatter refFormatter     // Function to format node references
}

// register adds a nodeRef to the registry and returns its reference.
func (nr *nodeRegistry) register(node *nodeRef) dumpNodeRef {
	// Check if the node is already registered.
	if id, exists := nr.ptrToID[node]; exists {
		return nr.addresses[id]
	}

	// Create a new reference for the node.
	id := len(nr.addresses)
	ref := dumpNodeRef{
		id:  id,
		ptr: node,
		fmt: nr.formatter,
	}

	// Register the node and its reference.
	nr.ptrToID[node] = id
	nr.addresses = append(nr.addresses, ref)

	return ref
}

// depthStorage stores information about the depth of the tree.
type depthStorage struct {
	childNum      int
	childrenTotal int
}

// treeStringer is a helper struct for generating a human-readable representation of the tree.
type treeStringer struct {
	storage      []depthStorage // Storage for depth information
	buf          *bytes.Buffer  // Buffer for building the string representation
	nodeRegistry *nodeRegistry  // Registry for node references
}

// String returns the string representation of the tree.
func (ts *treeStringer) String() string {
	s := ts.buf.String()
	// trim trailing whitespace and newlines.
	s = strings.TrimRight(s, "\n")
	s = strings.TrimRight(s, " ")

	return s
}

// regNode registers a nodeRef and returns its reference.
func (ts *treeStringer) regNode(node *nodeRef) dumpNodeRef {
	addr := ts.nodeRegistry.register(node)

	return addr
}

// regNodes registers a slice of artNodes and returns their references.
func (ts *treeStringer) regNodes(nodes []*nodeRef) []dumpNodeRef {
	if nodes == nil {
		return nil
	}

	addrs := make([]dumpNodeRef, 0, len(nodes))
	for _, n := range nodes {
		addrs = append(addrs, ts.nodeRegistry.register(n))
	}

	return addrs
}

// generatePads generates padding strings for the tree representation.
func (ts *treeStringer) generatePads(depth int, childNum int, childrenTotal int) (pad0, pad string) {
	ts.storage[depth] = depthStorage{childNum, childrenTotal}

	for d := 0; d <= depth; d++ {
		if d < depth {
			if ts.storage[d].childNum+1 < ts.storage[d].childrenTotal {
				pad0 += "│   "
			} else {
				pad0 += "    "
			}
		} else {
			if childrenTotal == 0 {
				pad0 += "─"
			} else if ts.storage[d].childNum+1 < ts.storage[d].childrenTotal {
				pad0 += "├"
			} else {
				pad0 += "└"
			}
			pad0 += "──"
		}

	}
	pad0 += " "

	for d := 0; d <= depth; d++ {
		if childNum+1 < childrenTotal && childrenTotal > 0 {
			if ts.storage[d].childNum+1 < ts.storage[d].childrenTotal {
				pad += "│   "
			} else {
				pad += "    "
			}
		} else if d < depth && ts.storage[d].childNum+1 < ts.storage[d].childrenTotal {
			pad += "│   "
		} else {
			pad += "    "
		}

	}

	return
}

// append adds a string representation of a value to the buffer.
// opts is a list of options for formatting the value.
// If no options are provided, the default is to print the value as a character.
// The available options are:
// - printValuesAsChar: print values as characters
// - printValuesAsDecimal: print values as decimal numbers
// - printValuesAsHex: print values as hexadecimal numbers
func (ts *treeStringer) append(v interface{}, opts ...int) *treeStringer {
	options := 0
	for _, opt := range opts {
		options |= opt
	}

	if options == 0 {
		options = printValueDefault
	}

	switch v := v.(type) {

	case string:
		ts.buf.WriteString(v)

	case []byte:
		ts.append("[")
		for i, b := range v {
			if (options & printValuesAsChar) != 0 {
				if b > 0 {
					ts.append(fmt.Sprintf("%c", b))
				} else {
					ts.append("·")
				}

			} else if (options & printValuesAsDecimal) != 0 {
				ts.append(fmt.Sprintf("%d", b))
			}
			if (options&printValuesAsDecimal) != 0 && i+1 < len(v) {
				ts.append(" ")
			}
		}
		ts.append("]")

	case Key:
		ts.append([]byte(v))

	default:
		ts.append("[")
		ts.append(fmt.Sprintf("%#v", v))
		ts.append("]")
	}

	return ts
}

// appendKey adds a string representation of a nodeRef's key to the buffer.
// see append for the list of available options.
func (ts *treeStringer) appendKey(keys []byte, present []byte, opts ...int) *treeStringer {
	options := 0
	for _, opt := range opts {
		options |= opt
	}

	if options == 0 {
		options = printValueDefault
	}

	ts.append("[")
	for i, b := range keys {
		if (options & printValuesAsChar) != 0 {
			if present[i] != 0 {
				ts.append(fmt.Sprintf("%c", b))
			} else {
				ts.append("·")
			}

		} else if (options & printValuesAsDecimal) != 0 {
			if present[i] != 0 {
				ts.append(fmt.Sprintf("%2d", b))
			} else {
				ts.append("·")
			}
		} else if (options & printValuesAsHex) != 0 {
			if present[i] != 0 {
				ts.append(fmt.Sprintf("%2x", b))
			} else {
				ts.append("·")
			}
		}
		if (options&(printValuesAsDecimal|printValuesAsHex)) != 0 && i+1 < len(keys) {
			ts.append(" ")
		}
	}
	ts.append("]")

	return ts
}

// children generates a string representation of the children of a nodeRef.
func (ts *treeStringer) children(children []*nodeRef, _ /*numChildred*/ uint16, keyOffset int, zeroChild *nodeRef) {
	for i, child := range children {
		ts.baseNode(child, keyOffset, i, len(children)+1)
	}

	ts.baseNode(zeroChild, keyOffset, len(children)+1, len(children)+1)
}

// node generates a string representation of a nodeRef.
func (ts *treeStringer) node(pad string, prefixLen uint16, prefix []byte, keys []byte, present []byte, children []*nodeRef, numChildren uint16, keyOffset int, zeroChild *nodeRef) {
	if prefix != nil {
		ts.append(pad).
			append(fmt.Sprintf("prefix(%x): ", prefixLen)).
			append(prefix).
			append(" ").
			append(fmt.Sprintf("%v", prefix)).
			append("\n")
	}

	if keys != nil {
		ts.append(pad).
			append("keys: ").
			appendKey(keys, present, printValuesAsChar).
			append(" ").
			appendKey(keys, present, printValuesAsDecimal).
			append("\n")
	}

	ts.append(pad).
		append(fmt.Sprintf("children(%v): %+v <%v>\n",
			numChildren,
			ts.regNodes(children),
			ts.regNode(zeroChild)))

	ts.children(children, numChildren, keyOffset+1, zeroChild)
}

func (ts *treeStringer) baseNode(an *nodeRef, depth int, childNum int, childrenTotal int) {
	padHeader, pad := ts.generatePads(depth, childNum, childrenTotal)
	if an == nil {
		ts.append(padHeader).
			append("nil").
			append("\n")
		return
	}

	ts.append(padHeader).
		append(fmt.Sprintf("%v (%v)\n",
			an.kind,
			ts.regNode(an)))

	switch an.kind {
	case Node4:
		nn := an.node4()

		ts.node(pad,
			nn.prefixLen,
			nn.prefix[:],
			nn.keys[:],
			nn.present[:],
			nn.children[:node4Max],
			nn.childrenLen,
			depth,
			nn.children[node4Max])

	case Node16:
		nn := an.node16()

		var present []byte
		for i := 0; i < len(nn.keys); i++ {
			var b byte
			if nn.hasChild(i) {
				b = 1
			}
			present = append(present, b)
		}

		ts.node(pad,
			nn.prefixLen,
			nn.prefix[:],
			nn.keys[:],
			present,
			nn.children[:node16Max],
			nn.childrenLen,
			depth,
			nn.children[node16Max])

	case Node48:
		nn := an.node48()

		var present []byte
		for i := 0; i < len(nn.keys); i++ {
			var b byte
			if nn.hasChild(i) {
				b = 1
			}
			present = append(present, b)
		}

		ts.node(pad,
			nn.prefixLen,
			nn.prefix[:],
			nn.keys[:],
			present,
			nn.children[:node48Max],
			nn.childrenLen,
			depth,
			nn.children[node48Max])

	case Node256:
		nn := an.node256()

		ts.node(pad,
			nn.prefixLen,
			nn.prefix[:],
			nil,
			nil,
			nn.children[:node256Max],
			nn.childrenLen,
			depth,
			nn.children[node256Max])

	case Leaf:
		n := an.leaf()

		ts.append(pad).
			append(fmt.Sprintf("key(%d): ", len(n.key))).
			append(n.key).
			append(" ").
			append(fmt.Sprintf("%v", n.key)).
			append("\n")

		if s, ok := n.value.(string); ok {
			ts.append(pad).
				append(fmt.Sprintf("val: %v\n",
					s))
		} else if b, ok := n.value.([]byte); ok {
			ts.append(pad).
				append(fmt.Sprintf("val: %v\n",
					string(b)))
		} else {
			ts.append(pad).
				append(fmt.Sprintf("val: %v\n",
					n.value))
		}

	}

	ts.append(pad).
		append("\n")
}

func (ts *treeStringer) startFromNode(an *nodeRef) {
	ts.baseNode(an, 0, 0, 0)
}

/*
DumpNode returns Tree in the human readable format:

--8<-- // main.go

	package main

	import (
		"fmt"
		art "github.com/plar/go-adaptive-radix-tree"
	)

	func main() {
		tree := art.New()
		terms := []string{"A", "a", "aa"}
		for _, term := range terms {
			tree.Insert(art.Key(term), term)
		}
		fmt.Println(tree)
	}

--8<--

	$ go run main.go

	─── Node4 (0xc00011c2d0)
		prefix(0): [··········] [0 0 0 0 0 0 0 0 0 0]
		keys: [Aa··] [65 97 · ·]
		children(2): [0xc00011c2a0 0xc00011c300 - -] <->
		├── Leaf (0xc00011c2a0)
		│   key(1): [A] [65]
		│   val: A
		│
		├── Node4 (0xc00011c300)
		│   prefix(0): [··········] [0 0 0 0 0 0 0 0 0 0]
		│   keys: [a···] [97 · · ·]
		│   children(1): [0xc00011c2f0 - - -] <0xc00011c2c0>
		│   ├── Leaf (0xc00011c2f0)
		│   │   key(2): [aa] [97 97]
		│   │   val: aa
		│   │
		│   ├── nil
		│   ├── nil
		│   ├── nil
		│   └── Leaf (0xc00011c2c0)
		│       key(1): [a] [97]
		│       val: a
		│
		│
		├── nil
		├── nil
		└── nil
*/
func DumpNode(root *nodeRef) string {
	opts := createTreeStringerOptions(WithRefFormatter(RefAddrFormatter))
	trs := newTreeStringer(opts)
	trs.startFromNode(root)
	return trs.String()
}

// treeStringerOptions contains options for DumpTree function.
type treeStringerOptions struct {
	storageSize int
	formatter   refFormatter
}

// treeStringerOption is a function that sets an option for DumpTree.
type treeStringerOption func(opts *treeStringerOptions)

// WithStorageSize sets the size of the storage for depth information.
func WithStorageSize(size int) treeStringerOption {
	return func(opts *treeStringerOptions) {
		opts.storageSize = size
	}
}

// WithRefFormatter sets the formatter for node references.
func WithRefFormatter(formatter refFormatter) treeStringerOption {
	return func(opts *treeStringerOptions) {
		opts.formatter = formatter
	}
}

// TreeStringer returns the string representation of the tree.
// The tree must be of type *art.tree.
func TreeStringer(t Tree, opts ...treeStringerOption) string {
	tr, ok := t.(*tree)
	if !ok {
		return "expected *art.tree"
	}

	trs := newTreeStringer(createTreeStringerOptions(opts...))
	trs.startFromNode(tr.root)
	return trs.String()
}

func createTreeStringerOptions(opts ...treeStringerOption) treeStringerOptions {
	defOpts := treeStringerOptions{
		storageSize: 4096,
		formatter:   RefShortFormatter,
	}

	for _, opt := range opts {
		opt(&defOpts)
	}

	return defOpts
}

func newTreeStringer(opts treeStringerOptions) *treeStringer {
	return &treeStringer{
		storage: make([]depthStorage, opts.storageSize),
		buf:     bytes.NewBufferString(""),
		nodeRegistry: &nodeRegistry{
			ptrToID:   make(map[*nodeRef]int),
			formatter: opts.formatter,
		},
	}
}

func defaultTreeStringer() *treeStringer {
	return newTreeStringer(createTreeStringerOptions())
}

// #endregion

// #region utils

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// copy the node from src to dst.
func copyNode(dst *node, src *node) {
	if dst == nil || src == nil {
		return
	}

	dst.prefixLen = src.prefixLen
	dst.prefix = src.prefix
}

// find the child node index by key.
func findIndex(keys []byte, ch byte) int {
	for i, key := range keys {
		if key == ch {
			return i
		}
	}

	return indexNotFound
}

// findLongestCommonPrefix returns the longest common prefix of key1 and key2.
func findLongestCommonPrefix(key1 Key, key2 Key, keyOffset int) int {
	limit := minInt(len(key1), len(key2))

	idx := keyOffset
	for ; idx < limit; idx++ {
		if key1[idx] != key2[idx] {
			break
		}
	}

	return idx - keyOffset
}

// nodeMinimum returns the minimum leaf node.
func nodeMinimum(children []*nodeRef) *leaf {
	numChildren := len(children)
	if numChildren == 0 {
		return nil
	}

	// zero byte key
	if children[numChildren-1] != nil {
		return children[numChildren-1].minimum()
	}

	for i := 0; i < numChildren-1; i++ {
		if children[i] != nil {
			return children[i].minimum()
		}
	}

	return nil
}

// nodeMaximum returns the maximum leaf node.
func nodeMaximum(children []*nodeRef) *leaf {
	for i := len(children) - 1; i >= 0; i-- {
		if children[i] != nil {
			return children[i].maximum()
		}
	}

	return nil
}

// ternary is a generic ternary operator.
func ternary[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}

	return ifFalse
}

// #endregion
