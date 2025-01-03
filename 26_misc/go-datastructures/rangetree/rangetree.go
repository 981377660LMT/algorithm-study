/*
Package rangetree is designed to store n-dimensional data in an easy-to-query
way.  Given this package's primary use as representing cartesian data, this
information is represented by int64s at n-dimensions.  This implementation
is not actually a tree but a sparse n-dimensional list.  This package also
includes two implementations of this sparse list, one mutable (and not threadsafe)
and another that is immutable copy-on-write which is threadsafe.  The mutable
version is obviously faster but will likely have write contention for any
consumer that needs a threadsafe rangetree.

TODO: unify both implementations with the same interface.
*/
package rangetree

import (
	"fmt"
	"sort"
	"sync"
)

// #region interface
// Entry defines items that can be added to the rangetree.
type Entry interface {
	// ValueAtDimension returns the value of this entry
	// at the specified dimension.
	ValueAtDimension(dimension uint64) int64
}

// Interval describes the methods required to query the rangetree.  Note that
// all ranges are inclusive.
type Interval interface {
	// LowAtDimension returns an integer representing the lower bound
	// at the requested dimension.
	LowAtDimension(dimension uint64) int64
	// HighAtDimension returns an integer representing the higher bound
	// at the request dimension.
	HighAtDimension(dimension uint64) int64
}

// RangeTree describes the methods available to the rangetree.
type RangeTree interface {
	// Add will add the provided entries to the tree.  Any entries that
	// were overwritten will be returned in the order in which they
	// were overwritten.  If an entry's addition does not overwrite, a nil
	// is returned for that entry's index in the provided cells.
	Add(entries ...Entry) Entries
	// Len returns the number of entries in the tree.
	Len() uint64
	// Delete will remove the provided entries from the tree.
	// Any entries that were deleted will be returned in the order in
	// which they were deleted.  If an entry does not exist to be deleted,
	// a nil is returned for that entry's index in the provided cells.
	Delete(entries ...Entry) Entries
	// Query will return a list of entries that fall within
	// the provided interval.  The values at dimensions are inclusive.
	Query(interval Interval) Entries
	// Apply will call the provided function with each entry that exists
	// within the provided range, in order.  Return false at any time to
	// cancel iteration.  Altering the entry in such a way that its location
	// changes will result in undefined behavior.
	Apply(interval Interval, fn func(Entry) bool)
	// Get returns any entries that exist at the addresses provided by the
	// given entries.  Entries are returned in the order in which they are
	// received.  If an entry cannot be found, a nil is returned in its
	// place.
	Get(entries ...Entry) Entries
	// InsertAtDimension will increment items at and above the given index
	// by the number provided.  Provide a negative number to to decrement.
	// Returned are two lists.  The first list is a list of entries that
	// were moved.  The second is a list entries that were deleted.  These
	// lists are exclusive.
	InsertAtDimension(dimension uint64, index, number int64) (Entries, Entries)
}

// #endregion

// #region entry

var entriesPool = sync.Pool{
	New: func() interface{} {
		return make(Entries, 0, 10)
	},
}

// Entries is a typed list of Entry that can be reused if Dispose
// is called.
type Entries []Entry

// Dispose will free the resources consumed by this list and
// allow the list to be reused.
func (entries *Entries) Dispose() {
	for i := 0; i < len(*entries); i++ {
		(*entries)[i] = nil
	}

	*entries = (*entries)[:0]
	entriesPool.Put(*entries)
}

// NewEntries will return a reused list of entries.
func NewEntries() Entries {
	return entriesPool.Get().(Entries)
}

// #endregion

// #region error

// NoEntriesError is returned from an operation that requires
// existing entries when none are found.
type NoEntriesError struct{}

func (nee NoEntriesError) Error() string {
	return `No entries in this tree.`
}

// OutOfDimensionError is returned when a requested operation
// doesn't meet dimensional requirements.
type OutOfDimensionError struct {
	provided, max uint64
}

func (oode OutOfDimensionError) Error() string {
	return fmt.Sprintf(`Provided dimension: %d is 
		greater than max dimension: %d`,
		oode.provided, oode.max,
	)
}

// #endregion

// #region immutable

type immutableRangeTree struct {
	number     uint64
	top        orderedNodes
	dimensions uint64
}

func newCache(dimensions uint64) []Int64Slice {
	cache := make([]Int64Slice, 0, dimensions-1)
	for i := uint64(0); i < dimensions; i++ {
		cache = append(cache, Int64Slice{})
	}
	return cache
}

func (irt *immutableRangeTree) needNextDimension() bool {
	return irt.dimensions > 1
}

func (irt *immutableRangeTree) add(nodes *orderedNodes, cache []Int64Slice, entry Entry, added *uint64) {
	var node *node
	list := nodes

	for i := uint64(1); i <= irt.dimensions; i++ {
		if isLastDimension(irt.dimensions, i) {
			if i != 1 && !cache[i-1].Exists(node.value) {
				nodes := make(orderedNodes, len(*list))
				copy(nodes, *list)
				list = &nodes
				cache[i-1].Insert(node.value)
			}

			newNode := newNode(entry.ValueAtDimension(i), entry, false)
			overwritten := list.add(newNode)
			if overwritten == nil {
				*added++
			}
			if node != nil {
				node.orderedNodes = *list
			}
			break
		}

		if i != 1 && !cache[i-1].Exists(node.value) {
			nodes := make(orderedNodes, len(*list))
			copy(nodes, *list)
			list = &nodes
			cache[i-1].Insert(node.value)
			node.orderedNodes = *list
		}

		node, _ = list.getOrAdd(entry, i, irt.dimensions)
		list = &node.orderedNodes
	}
}

// Add will add the provided entries into the tree and return
// a new tree with those entries added.
func (irt *immutableRangeTree) Add(entries ...Entry) *immutableRangeTree {
	if len(entries) == 0 {
		return irt
	}

	cache := newCache(irt.dimensions)
	top := make(orderedNodes, len(irt.top))
	copy(top, irt.top)
	added := uint64(0)
	for _, entry := range entries {
		irt.add(&top, cache, entry, &added)
	}

	tree := NewImmutableRangeTree(irt.dimensions)
	tree.top = top
	tree.number = irt.number + added
	return tree
}

// InsertAtDimension will increment items at and above the given index
// by the number provided.  Provide a negative number to to decrement.
// Returned are two lists and the modified tree.  The first list is a
// list of entries that were moved.  The second is a list entries that
// were deleted.  These lists are exclusive.
func (irt *immutableRangeTree) InsertAtDimension(dimension uint64,
	index, number int64) (*immutableRangeTree, Entries, Entries) {

	if dimension > irt.dimensions || number == 0 {
		return irt, nil, nil
	}

	modified, deleted := make(Entries, 0, 100), make(Entries, 0, 100)

	tree := NewImmutableRangeTree(irt.dimensions)
	tree.top = irt.top.immutableInsert(
		dimension, 1, irt.dimensions,
		index, number,
		&modified, &deleted,
	)
	tree.number = irt.number - uint64(len(deleted))

	return tree, modified, deleted
}

type immutableNodeBundle struct {
	list         *orderedNodes
	index        int
	previousNode *node
	newNode      *node
}

func (irt *immutableRangeTree) Delete(entries ...Entry) *immutableRangeTree {
	cache := newCache(irt.dimensions)
	top := make(orderedNodes, len(irt.top))
	copy(top, irt.top)
	deleted := uint64(0)
	for _, entry := range entries {
		irt.delete(&top, cache, entry, &deleted)
	}

	tree := NewImmutableRangeTree(irt.dimensions)
	tree.top = top
	tree.number = irt.number - deleted
	return tree
}

func (irt *immutableRangeTree) delete(top *orderedNodes,
	cache []Int64Slice, entry Entry, deleted *uint64) {

	path := make([]*immutableNodeBundle, 0, 5)
	var index int
	var n *node
	var local *node
	list := top

	for i := uint64(1); i <= irt.dimensions; i++ {
		value := entry.ValueAtDimension(i)
		local, index = list.get(value)
		if local == nil { // there's nothing to delete
			return
		}

		nb := &immutableNodeBundle{
			list:         list,
			index:        index,
			previousNode: n,
		}
		path = append(path, nb)
		n = local
		list = &n.orderedNodes
	}

	*deleted++

	for i := len(path) - 1; i >= 0; i-- {
		nb := path[i]
		if nb.previousNode != nil {
			nodes := make(orderedNodes, len(*nb.list))
			copy(nodes, *nb.list)
			nb.list = &nodes
			if len(*nb.list) == 1 {
				continue
			}
			nn := newNode(
				nb.previousNode.value,
				nb.previousNode.entry,
				!isLastDimension(irt.dimensions, uint64(i)+1),
			)
			nn.orderedNodes = nodes
			path[i-1].newNode = nn
		}
	}

	for _, nb := range path {
		if nb.newNode == nil {
			nb.list.deleteAt(nb.index)
		} else {
			(*nb.list)[nb.index] = nb.newNode
		}
	}
}

func (irt *immutableRangeTree) apply(list orderedNodes, interval Interval,
	dimension uint64, fn func(*node) bool) bool {

	low, high := interval.LowAtDimension(dimension), interval.HighAtDimension(dimension)

	if isLastDimension(irt.dimensions, dimension) {
		if !list.apply(low, high, fn) {
			return false
		}
	} else {
		if !list.apply(low, high, func(n *node) bool {
			if !irt.apply(n.orderedNodes, interval, dimension+1, fn) {
				return false
			}
			return true
		}) {
			return false
		}
		return true
	}

	return true
}

// Query will return an ordered list of results in the given
// interval.
func (irt *immutableRangeTree) Query(interval Interval) Entries {
	entries := NewEntries()

	irt.apply(irt.top, interval, 1, func(n *node) bool {
		entries = append(entries, n.entry)
		return true
	})

	return entries
}

func (irt *immutableRangeTree) get(entry Entry) Entry {
	on := irt.top
	for i := uint64(1); i <= irt.dimensions; i++ {
		n, _ := on.get(entry.ValueAtDimension(i))
		if n == nil {
			return nil
		}
		if i == irt.dimensions {
			return n.entry
		}
		on = n.orderedNodes
	}

	return nil
}

// Get returns any entries that exist at the addresses provided by the
// given entries.  Entries are returned in the order in which they are
// received.  If an entry cannot be found, a nil is returned in its
// place.
func (irt *immutableRangeTree) Get(entries ...Entry) Entries {
	result := make(Entries, 0, len(entries))
	for _, entry := range entries {
		result = append(result, irt.get(entry))
	}

	return result
}

// Len returns the number of items in this tree.
func (irt *immutableRangeTree) Len() uint64 {
	return irt.number
}

func NewImmutableRangeTree(dimensions uint64) *immutableRangeTree {
	return &immutableRangeTree{
		dimensions: dimensions,
	}
}

// #endregion

// #region node

type nodes []*node

type node struct {
	value        int64
	entry        Entry
	orderedNodes orderedNodes
}

func newNode(value int64, entry Entry, needNextDimension bool) *node {
	n := &node{}
	n.value = value
	if needNextDimension {
		n.orderedNodes = make(orderedNodes, 0, 10)
	} else {
		n.entry = entry
	}

	return n
}

// #endregion

// #region ordered

// orderedNodes represents an ordered list of points living
// at the last dimension.  No duplicates can be inserted here.
type orderedNodes nodes

func (nodes orderedNodes) search(value int64) int {
	return sort.Search(
		len(nodes),
		func(i int) bool { return nodes[i].value >= value },
	)
}

// addAt will add the provided node at the provided index.  Returns
// a node if one was overwritten.
func (nodes *orderedNodes) addAt(i int, node *node) *node {
	if i == len(*nodes) {
		*nodes = append(*nodes, node)
		return nil
	}

	if (*nodes)[i].value == node.value {
		overwritten := (*nodes)[i]
		// this is a duplicate, there can't be a duplicate
		// point in the last dimension
		(*nodes)[i] = node
		return overwritten
	}

	*nodes = append(*nodes, nil)
	copy((*nodes)[i+1:], (*nodes)[i:])
	(*nodes)[i] = node
	return nil
}

func (nodes *orderedNodes) add(node *node) *node {
	i := nodes.search(node.value)
	return nodes.addAt(i, node)
}

func (nodes *orderedNodes) deleteAt(i int) *node {
	if i >= len(*nodes) { // no matching found
		return nil
	}

	deleted := (*nodes)[i]
	copy((*nodes)[i:], (*nodes)[i+1:])
	(*nodes)[len(*nodes)-1] = nil
	*nodes = (*nodes)[:len(*nodes)-1]
	return deleted
}

func (nodes *orderedNodes) delete(value int64) *node {
	i := nodes.search(value)

	if (*nodes)[i].value != value || i == len(*nodes) {
		return nil
	}

	return nodes.deleteAt(i)
}

func (nodes orderedNodes) apply(low, high int64, fn func(*node) bool) bool {
	index := nodes.search(low)
	if index == len(nodes) {
		return true
	}

	for ; index < len(nodes); index++ {
		if nodes[index].value > high {
			break
		}

		if !fn(nodes[index]) {
			return false
		}
	}

	return true
}

func (nodes orderedNodes) get(value int64) (*node, int) {
	i := nodes.search(value)
	if i == len(nodes) {
		return nil, i
	}

	if nodes[i].value == value {
		return nodes[i], i
	}

	return nil, i
}

func (nodes *orderedNodes) getOrAdd(entry Entry,
	dimension, lastDimension uint64) (*node, bool) {

	isLastDimension := isLastDimension(lastDimension, dimension)
	value := entry.ValueAtDimension(dimension)

	i := nodes.search(value)
	if i == len(*nodes) {
		node := newNode(value, entry, !isLastDimension)
		*nodes = append(*nodes, node)
		return node, true
	}

	if (*nodes)[i].value == value {
		return (*nodes)[i], false
	}

	node := newNode(value, entry, !isLastDimension)
	*nodes = append(*nodes, nil)
	copy((*nodes)[i+1:], (*nodes)[i:])
	(*nodes)[i] = node
	return node, true
}

func (nodes orderedNodes) flatten(entries *Entries) {
	for _, node := range nodes {
		if node.orderedNodes != nil {
			node.orderedNodes.flatten(entries)
		} else {
			*entries = append(*entries, node.entry)
		}
	}
}

func (nodes *orderedNodes) insert(insertDimension, dimension, maxDimension uint64,
	index, number int64, modified, deleted *Entries) {

	lastDimension := isLastDimension(maxDimension, dimension)

	if insertDimension == dimension {
		i := nodes.search(index)
		var toDelete []int

		for j := i; j < len(*nodes); j++ {
			(*nodes)[j].value += number
			if (*nodes)[j].value < index {
				toDelete = append(toDelete, j)
				if lastDimension {
					*deleted = append(*deleted, (*nodes)[j].entry)
				} else {
					(*nodes)[j].orderedNodes.flatten(deleted)
				}
				continue
			}
			if lastDimension {
				*modified = append(*modified, (*nodes)[j].entry)
			} else {
				(*nodes)[j].orderedNodes.flatten(modified)
			}
		}

		for i, index := range toDelete {
			nodes.deleteAt(index - i)
		}

		return
	}

	for _, node := range *nodes {
		node.orderedNodes.insert(
			insertDimension, dimension+1, maxDimension,
			index, number, modified, deleted,
		)
	}
}

func (nodes orderedNodes) immutableInsert(insertDimension, dimension, maxDimension uint64,
	index, number int64, modified, deleted *Entries) orderedNodes {

	lastDimension := isLastDimension(maxDimension, dimension)

	cp := make(orderedNodes, len(nodes))
	copy(cp, nodes)

	if insertDimension == dimension {
		i := cp.search(index)
		var toDelete []int

		for j := i; j < len(cp); j++ {
			nn := newNode(cp[j].value+number, cp[j].entry, !lastDimension)
			nn.orderedNodes = cp[j].orderedNodes
			cp[j] = nn
			if cp[j].value < index {
				toDelete = append(toDelete, j)
				if lastDimension {
					*deleted = append(*deleted, cp[j].entry)
				} else {
					cp[j].orderedNodes.flatten(deleted)
				}
				continue
			}
			if lastDimension {
				*modified = append(*modified, cp[j].entry)
			} else {
				cp[j].orderedNodes.flatten(modified)
			}
		}

		for _, index := range toDelete {
			cp.deleteAt(index)
		}

		return cp
	}

	for i := 0; i < len(cp); i++ {
		oldNode := nodes[i]
		nn := newNode(oldNode.value, oldNode.entry, !lastDimension)
		nn.orderedNodes = oldNode.orderedNodes.immutableInsert(
			insertDimension, dimension+1,
			maxDimension,
			index, number,
			modified, deleted,
		)
		cp[i] = nn
	}

	return cp
}

// #endregion

// #region orderedtree

func isLastDimension(value, test uint64) bool {
	return test >= value
}

type nodeBundle struct {
	list  *orderedNodes
	index int
}

type orderedTree struct {
	top        orderedNodes
	number     uint64
	dimensions uint64
	path       []*nodeBundle
}

func (ot *orderedTree) resetPath() {
	ot.path = ot.path[:0]
}

func (ot *orderedTree) needNextDimension() bool {
	return ot.dimensions > 1
}

// add will add the provided entry to the rangetree and return an
// entry if one was overwritten.
func (ot *orderedTree) add(entry Entry) *node {
	var node *node
	list := &ot.top

	for i := uint64(1); i <= ot.dimensions; i++ {
		if isLastDimension(ot.dimensions, i) {
			overwritten := list.add(
				newNode(entry.ValueAtDimension(i), entry, false),
			)
			if overwritten == nil {
				ot.number++
			}
			return overwritten
		}
		node, _ = list.getOrAdd(entry, i, ot.dimensions)
		list = &node.orderedNodes
	}

	return nil
}

// Add will add the provided entries to the tree.  This method
// returns a list of entries that were overwritten in the order
// in which entries were received.  If an entry doesn't overwrite
// anything, a nil will be returned for that entry in the returned
// slice.
func (ot *orderedTree) Add(entries ...Entry) Entries {
	if len(entries) == 0 {
		return nil
	}

	overwrittens := make(Entries, len(entries))
	for i, entry := range entries {
		if entry == nil {
			continue
		}

		overwritten := ot.add(entry)
		if overwritten != nil {
			overwrittens[i] = overwritten.entry
		}
	}

	return overwrittens
}

func (ot *orderedTree) delete(entry Entry) *node {
	ot.resetPath()
	var index int
	var node *node
	list := &ot.top

	for i := uint64(1); i <= ot.dimensions; i++ {
		value := entry.ValueAtDimension(i)
		node, index = list.get(value)
		if node == nil { // there's nothing to delete
			return nil
		}

		nb := &nodeBundle{list: list, index: index}
		ot.path = append(ot.path, nb)

		list = &node.orderedNodes
	}

	ot.number--

	for i := len(ot.path) - 1; i >= 0; i-- {
		nb := ot.path[i]
		nb.list.deleteAt(nb.index)
		if len(*nb.list) > 0 {
			break
		}
	}

	return node
}

func (ot *orderedTree) get(entry Entry) Entry {
	on := ot.top
	for i := uint64(1); i <= ot.dimensions; i++ {
		n, _ := on.get(entry.ValueAtDimension(i))
		if n == nil {
			return nil
		}
		if i == ot.dimensions {
			return n.entry
		}
		on = n.orderedNodes
	}

	return nil
}

// Get returns any entries that exist at the addresses provided by the
// given entries.  Entries are returned in the order in which they are
// received.  If an entry cannot be found, a nil is returned in its
// place.
func (ot *orderedTree) Get(entries ...Entry) Entries {
	result := make(Entries, 0, len(entries))
	for _, entry := range entries {
		result = append(result, ot.get(entry))
	}

	return result
}

// Delete will remove the provided entries from the tree.
// Any entries that were deleted will be returned in the order in
// which they were deleted.  If an entry does not exist to be deleted,
// a nil is returned for that entry's index in the provided cells.
func (ot *orderedTree) Delete(entries ...Entry) Entries {
	if len(entries) == 0 {
		return nil
	}

	deletedEntries := make(Entries, len(entries))
	for i, entry := range entries {
		if entry == nil {
			continue
		}

		deleted := ot.delete(entry)
		if deleted != nil {
			deletedEntries[i] = deleted.entry
		}
	}

	return deletedEntries
}

// Len returns the number of items in the tree.
func (ot *orderedTree) Len() uint64 {
	return ot.number
}

func (ot *orderedTree) apply(list orderedNodes, interval Interval,
	dimension uint64, fn func(*node) bool) bool {

	low, high := interval.LowAtDimension(dimension), interval.HighAtDimension(dimension)

	if isLastDimension(ot.dimensions, dimension) {
		if !list.apply(low, high, fn) {
			return false
		}
	} else {
		if !list.apply(low, high, func(n *node) bool {
			if !ot.apply(n.orderedNodes, interval, dimension+1, fn) {
				return false
			}
			return true
		}) {
			return false
		}
		return true
	}

	return true
}

// Apply will call (in order) the provided function to every
// entry that falls within the provided interval.  Any alteration
// the the entry that would result in different answers to the
// interface methods results in undefined behavior.
func (ot *orderedTree) Apply(interval Interval, fn func(Entry) bool) {
	ot.apply(ot.top, interval, 1, func(n *node) bool {
		return fn(n.entry)
	})
}

// Query will return an ordered list of results in the given
// interval.
func (ot *orderedTree) Query(interval Interval) Entries {
	entries := NewEntries()

	ot.apply(ot.top, interval, 1, func(n *node) bool {
		entries = append(entries, n.entry)
		return true
	})

	return entries
}

// InsertAtDimension will increment items at and above the given index
// by the number provided.  Provide a negative number to to decrement.
// Returned are two lists.  The first list is a list of entries that
// were moved.  The second is a list entries that were deleted.  These
// lists are exclusive.
func (ot *orderedTree) InsertAtDimension(dimension uint64,
	index, number int64) (Entries, Entries) {

	// TODO: perhaps return an error here?
	if dimension > ot.dimensions || number == 0 {
		return nil, nil
	}

	modified := make(Entries, 0, 100)
	deleted := make(Entries, 0, 100)

	ot.top.insert(dimension, 1, ot.dimensions,
		index, number, &modified, &deleted,
	)

	ot.number -= uint64(len(deleted))

	return modified, deleted
}

func NewOrderedTree(dimensions uint64) *orderedTree {
	return &orderedTree{
		dimensions: dimensions,
		path:       make([]*nodeBundle, 0, dimensions),
	}
}

// #endregion

// #region slice

// Int64Slice is a slice that fulfills the sort.Interface interface.
type Int64Slice []int64

// Len returns the len of this slice.  Required by sort.Interface.
func (s Int64Slice) Len() int {
	return len(s)
}

// Less returns a bool indicating if the value at position i
// is less than at position j.  Required by sort.Interface.
func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

// Search will search this slice and return an index that corresponds
// to the lowest position of that value.  You'll need to check
// separately if the value at that position is equal to x.  The
// behavior of this method is undefinited if the slice is not sorted.
func (s Int64Slice) Search(x int64) int {
	return sort.Search(len(s), func(i int) bool {
		return s[i] >= x
	})
}

// Sort will in-place sort this list of int64s.
func (s Int64Slice) Sort() {
	sort.Sort(s)
}

// Swap will swap the elements at positions i and j.  This is required
// by sort.Interface.
func (s Int64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Exists returns a bool indicating if the provided value exists
// in this list.  This has undefined behavior if the list is not
// sorted.
func (s Int64Slice) Exists(x int64) bool {
	i := s.Search(x)
	if i == len(s) {
		return false
	}

	return s[i] == x
}

// Insert will insert x into the sorted position in this list
// and return a list with the value added.  If this slice has not
// been sorted Insert's behavior is undefined.
func (s Int64Slice) Insert(x int64) Int64Slice {
	i := s.Search(x)
	if i == len(s) {
		return append(s, x)
	}

	if s[i] == x {
		return s
	}

	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = x
	return s
}

// #endregion
