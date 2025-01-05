/*
Package hilbert implements a Hilbert R-tree based on PALM principles
to improve multithreaded performance.  This package is not quite complete
and some optimization and delete codes remain to be completed.

This serves as a potential replacement for the interval tree and rangetree.

Benchmarks:
BenchmarkBulkAddPoints-8	     500	   2589270 ns/op
BenchmarkBulkUpdatePoints-8	    2000	   1212641 ns/op
BenchmarkPointInsertion-8	  200000	      9135 ns/op
BenchmarkQueryPoints-8	  	  500000	      3122 ns/op
*/
package main

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// Rectangle implementation
type MyRectangle struct {
	xlow, ylow, xhigh, yhigh int32
}

func (r *MyRectangle) LowerLeft() (int32, int32) {
	return r.xlow, r.ylow
}

func (r *MyRectangle) UpperRight() (int32, int32) {
	return r.xhigh, r.yhigh
}

// Utility function to create a new MyRectangle
func NewMyRectangle(xlow, ylow, xhigh, yhigh int32) *MyRectangle {
	return &MyRectangle{
		xlow:  xlow,
		ylow:  ylow,
		xhigh: xhigh,
		yhigh: yhigh,
	}
}

func main() {
	// Initialize the R-Tree
	bufferSize := uint64(1024) // Adjust based on expected workload
	arity := uint64(16)        // Branching factor
	tree := NewTree(bufferSize, arity)
	defer tree.Dispose() // Ensure resources are cleaned up

	// Create some rectangles
	rects := []*MyRectangle{
		NewMyRectangle(0, 0, 10, 10),
		NewMyRectangle(5, 5, 15, 15),
		NewMyRectangle(10, 10, 20, 20),
		NewMyRectangle(15, 15, 25, 25),
	}

	// Insert rectangles into the tree
	for _, rect := range rects {
		tree.Insert(rect)
		fmt.Printf("Inserted rectangle: %+v\n", rect)
	}

	// Perform a search for rectangles intersecting with a given rectangle
	searchRect := NewMyRectangle(7, 7, 12, 12)
	results := tree.Search(searchRect)
	fmt.Printf("\nSearch results for rectangle %+v:\n", searchRect)
	for _, r := range results {
		fmt.Printf(" - %+v\n", r)
	}

	// Delete a rectangle from the tree
	toDelete := rects[1] // Rectangle with coordinates (5,5,15,15)
	tree.Delete(toDelete)
	fmt.Printf("\nDeleted rectangle: %+v\n", toDelete)

	// Perform the same search again to see updated results
	results = tree.Search(searchRect)
	fmt.Printf("\nSearch results after deletion for rectangle %+v:\n", searchRect)
	for _, r := range results {
		fmt.Printf(" - %+v\n", r)
	}

	// Display the total number of rectangles in the tree
	fmt.Printf("\nTotal rectangles in the tree: %d\n", tree.Len())
}

// #region interface
// Rectangles is a typed list of Rectangle.
type Rectangles []Rectangle

// Rectangle describes a two-dimensional bound.
type Rectangle interface {
	// LowerLeft describes the lower left coordinate of this rectangle.
	LowerLeft() (int32, int32)
	// UpperRight describes the upper right coordinate of this rectangle.
	UpperRight() (int32, int32)
}

// RTree defines an object that can be returned from any subpackage
// of this package.
type RTree interface {
	// Search will perform an intersection search of the given
	// rectangle and return any rectangles that intersect.
	Search(Rectangle) Rectangles
	// Len returns in the number of items in the RTree.
	Len() uint64
	// Dispose will clean up any objects used by the RTree.
	Dispose()
	// Delete will remove the provided rectangles from the RTree.
	Delete(...Rectangle)
	// Insert will add the provided rectangles to the RTree.
	Insert(...Rectangle)
}

// #endregion

// #region tree
type operation int

const (
	get operation = iota
	add
	remove
)

const multiThreadAt = 1000 // number of keys before we multithread lookups

type keyBundle struct {
	key         hilbert
	left, right Rectangle
}

type tree struct {
	root            *node
	_               [8]uint64
	number          uint64
	_               [8]uint64
	ary, bufferSize uint64
	actions         *RingBuffer
	cache           []interface{}
	_               [8]uint64
	disposed        uint64
	_               [8]uint64
	running         uint64
}

func (tree *tree) checkAndRun(action action) {
	if tree.actions.Len() > 0 {
		if action != nil {
			tree.actions.Put(action)
		}
		if atomic.CompareAndSwapUint64(&tree.running, 0, 1) {
			var a interface{}
			var err error
			for tree.actions.Len() > 0 {
				a, err = tree.actions.Get()
				if err != nil {
					return
				}
				tree.cache = append(tree.cache, a)
				if uint64(len(tree.cache)) >= tree.bufferSize {
					break
				}
			}

			go tree.operationRunner(tree.cache, true)
		}
	} else if action != nil {
		if atomic.CompareAndSwapUint64(&tree.running, 0, 1) {
			switch action.operation() {
			case get:
				ga := action.(*getAction)
				result := tree.search(ga.lookup)
				ga.result = result
				action.complete()
				tree.reset()
			case add, remove:
				if len(action.keys()) > multiThreadAt {
					tree.operationRunner(interfaces{action}, true)
				} else {
					tree.operationRunner(interfaces{action}, false)
				}
			}
		} else {
			tree.actions.Put(action)
			tree.checkAndRun(nil)
		}
	}
}

func (tree *tree) init(bufferSize, ary uint64) {
	tree.bufferSize = bufferSize
	tree.ary = ary
	tree.cache = make([]interface{}, 0, bufferSize)
	tree.root = newNode(true, newKeys(ary), newNodes(ary))
	tree.root.mbr = &rectangle{}
	tree.actions = NewRingBuffer(tree.bufferSize)
}

func (tree *tree) operationRunner(xns interfaces, threaded bool) {
	writeOperations, deleteOperations, toComplete := tree.fetchKeys(xns, threaded)
	tree.recursiveMutate(writeOperations, deleteOperations, false, threaded)
	for _, a := range toComplete {
		a.complete()
	}

	tree.reset()
}

func (tree *tree) fetchKeys(xns interfaces, inParallel bool) (map[*node][]*keyBundle, map[*node][]*keyBundle, actions) {
	if inParallel {
		tree.fetchKeysInParallel(xns)
	} else {
		tree.fetchKeysInSerial(xns)
	}

	writeOperations := make(map[*node][]*keyBundle)
	deleteOperations := make(map[*node][]*keyBundle)
	toComplete := make(actions, 0, len(xns)/2)
	for _, ifc := range xns {
		action := ifc.(action)
		switch action.operation() {
		case add:
			for i, n := range action.nodes() {
				writeOperations[n] = append(writeOperations[n], &keyBundle{key: action.rects()[i].hilbert, left: action.rects()[i].rect})
			}
			toComplete = append(toComplete, action)
		case remove:
			for i, n := range action.nodes() {
				deleteOperations[n] = append(deleteOperations[n], &keyBundle{key: action.rects()[i].hilbert, left: action.rects()[i].rect})
			}
			toComplete = append(toComplete, action)
		case get:
			action.complete()
		}
	}

	return writeOperations, deleteOperations, toComplete
}

func (tree *tree) fetchKeysInSerial(xns interfaces) {
	for _, ifc := range xns {
		action := ifc.(action)
		switch action.operation() {
		case add, remove:
			for i, key := range action.rects() {
				n := getParent(tree.root, key.hilbert, key.rect)
				action.addNode(int64(i), n)
			}
		case get:
			ga := action.(*getAction)
			rects := tree.search(ga.lookup)
			ga.result = rects
		}
	}
}

func (tree *tree) reset() {
	for i := range tree.cache {
		tree.cache[i] = nil
	}

	tree.cache = tree.cache[:0]
	atomic.StoreUint64(&tree.running, 0)
	tree.checkAndRun(nil)
}

func (tree *tree) fetchKeysInParallel(xns []interface{}) {
	var forCache struct {
		i      int64
		buffer [8]uint64 // different cache lines
		js     []int64
	}

	for j := 0; j < len(xns); j++ {
		forCache.js = append(forCache.js, -1)
	}
	numCPU := runtime.NumCPU()
	if numCPU > 1 {
		numCPU--
	}
	var wg sync.WaitGroup
	wg.Add(numCPU)

	for k := 0; k < numCPU; k++ {
		go func() {
			for {
				index := atomic.LoadInt64(&forCache.i)
				if index >= int64(len(xns)) {
					break
				}
				action := xns[index].(action)

				j := atomic.AddInt64(&forCache.js[index], 1)
				if j > int64(len(action.rects())) { // someone else is updating i
					continue
				} else if j == int64(len(action.rects())) {
					atomic.StoreInt64(&forCache.i, index+1)
					continue
				}

				switch action.operation() {
				case add, remove:
					hb := action.rects()[j]
					n := getParent(tree.root, hb.hilbert, hb.rect)
					action.addNode(j, n)
				case get:
					ga := action.(*getAction)
					result := tree.search(ga.lookup)
					ga.result = result
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (tree *tree) splitNode(n, parent *node, nodes *[]*node, keys *hilberts) {
	if !n.needsSplit(tree.ary) {
		return
	}

	length := n.keys.len()
	splitAt := tree.ary - 1

	for i := splitAt; i < length; i += splitAt {
		offset := length - i
		k, left, right := n.split(offset, tree.ary)
		left.right = right
		*keys = append(*keys, k)
		*nodes = append(*nodes, left, right)
		left.parent = parent
		right.parent = parent
	}
}

func (tree *tree) applyNode(n *node, adds, deletes []*keyBundle) {
	for _, kb := range deletes {
		if n.keys.len() == 0 {
			break
		}

		deleted := n.delete(kb)
		if deleted != nil {
			atomic.AddUint64(&tree.number, ^uint64(0))
		}
	}

	for _, kb := range adds {
		old := n.insert(kb)
		if n.isLeaf && old == nil {
			atomic.AddUint64(&tree.number, 1)
		}
	}
}

func (tree *tree) recursiveMutate(adds, deletes map[*node][]*keyBundle, setRoot, inParallel bool) {
	if len(adds) == 0 && len(deletes) == 0 {
		return
	}

	if setRoot && len(adds) > 1 {
		panic(`SHOULD ONLY HAVE ONE ROOT`)
	}

	ifs := make(interfaces, 0, len(adds))
	for n := range adds {
		if n.parent == nil {
			setRoot = true
		}
		ifs = append(ifs, n)
	}

	for n := range deletes {
		if n.parent == nil {
			setRoot = true
		}

		if _, ok := adds[n]; !ok {
			ifs = append(ifs, n)
		}
	}

	var dummyRoot *node
	if setRoot {
		dummyRoot = &node{
			keys:  newKeys(tree.ary),
			nodes: newNodes(tree.ary),
			mbr:   &rectangle{},
		}
	}

	var write sync.Mutex
	nextLayerWrite := make(map[*node][]*keyBundle)
	nextLayerDelete := make(map[*node][]*keyBundle)

	var mutate func(interfaces, func(interface{}))
	if inParallel {
		mutate = executeInterfacesInParallel
	} else {
		mutate = executeInterfacesInSerial
	}

	mutate(ifs, func(ifc interface{}) {
		n := ifc.(*node)
		adds := adds[n]
		deletes := deletes[n]

		if len(adds) == 0 && len(deletes) == 0 {
			return
		}

		if setRoot {
			tree.root = n
		}

		parent := n.parent
		if parent == nil {
			parent = dummyRoot
			setRoot = true
		}

		tree.applyNode(n, adds, deletes)

		if n.needsSplit(tree.ary) {
			keys := make(hilberts, 0, n.keys.len())
			nodes := make([]*node, 0, n.nodes.len())
			tree.splitNode(n, parent, &nodes, &keys)
			write.Lock()
			for i, k := range keys {
				nextLayerWrite[parent] = append(nextLayerWrite[parent], &keyBundle{key: k, left: nodes[i*2], right: nodes[i*2+1]})
			}
			write.Unlock()
		}
	})

	tree.recursiveMutate(nextLayerWrite, nextLayerDelete, setRoot, inParallel)
}

// Insert will add the provided keys to the tree.
func (tree *tree) Insert(rects ...Rectangle) {
	ia := newInsertAction(rects)
	tree.checkAndRun(ia)
	ia.completer.Wait()
}

// Delete will remove the provided keys from the tree.  If no
// matching key is found, this is a no-op.
func (tree *tree) Delete(rects ...Rectangle) {
	ra := newRemoveAction(rects)
	tree.checkAndRun(ra)
	ra.completer.Wait()
}

func (tree *tree) search(r *rectangle) Rectangles {
	if tree.root == nil {
		return Rectangles{}
	}

	result := make(Rectangles, 0, 10)
	whs := tree.root.searchRects(r)
	for len(whs) > 0 {
		wh := whs[0]
		if n, ok := wh.(*node); ok {
			whs = append(whs, n.searchRects(r)...)
		} else {
			result = append(result, wh)
		}
		whs = whs[1:]
	}

	return result
}

// Search will return a list of rectangles that intersect the provided
// rectangle.
func (tree *tree) Search(rect Rectangle) Rectangles {
	ga := newGetAction(rect)
	tree.checkAndRun(ga)
	ga.completer.Wait()
	return ga.result
}

// Len returns the number of items in the tree.
func (tree *tree) Len() uint64 {
	return atomic.LoadUint64(&tree.number)
}

// Dispose will clean up any resources used by this tree.  This
// must be called to prevent a memory leak.
func (tree *tree) Dispose() {
	tree.actions.Dispose()
	atomic.StoreUint64(&tree.disposed, 1)
}

func NewTree(bufferSize, ary uint64) *tree {
	tree := &tree{}
	tree.init(bufferSize, ary)
	return tree
}

// #endregion

type actions []action

type action interface {
	operation() operation
	keys() hilberts
	rects() []*hilbertBundle
	complete()
	addNode(int64, *node)
	nodes() []*node
}

type getAction struct {
	result    Rectangles
	completer *sync.WaitGroup
	lookup    *rectangle
}

func (ga *getAction) complete() {
	ga.completer.Done()
}

func (ga *getAction) operation() operation {
	return get
}

func (ga *getAction) keys() hilberts {
	return nil
}

func (ga *getAction) addNode(i int64, n *node) {
	return // not necessary for gets
}

func (ga *getAction) nodes() []*node {
	return nil
}

func (ga *getAction) rects() []*hilbertBundle {
	return []*hilbertBundle{&hilbertBundle{}}
}

func newGetAction(rect Rectangle) *getAction {
	r := newRectangeFromRect(rect)
	ga := &getAction{
		completer: new(sync.WaitGroup),
		lookup:    r,
	}
	ga.completer.Add(1)
	return ga
}

type insertAction struct {
	rs        []*hilbertBundle
	completer *sync.WaitGroup
	ns        []*node
}

func (ia *insertAction) complete() {
	ia.completer.Done()
}

func (ia *insertAction) operation() operation {
	return add
}

func (ia *insertAction) keys() hilberts {
	return nil
}

func (ia *insertAction) addNode(i int64, n *node) {
	ia.ns[i] = n
}

func (ia *insertAction) nodes() []*node {
	return ia.ns
}

func (ia *insertAction) rects() []*hilbertBundle {
	return ia.rs
}

func newInsertAction(rects Rectangles) *insertAction {
	ia := &insertAction{
		rs:        bundlesFromRects(rects...),
		completer: new(sync.WaitGroup),
		ns:        make([]*node, len(rects)),
	}
	ia.completer.Add(1)
	return ia
}

type removeAction struct {
	*insertAction
}

func (ra *removeAction) operation() operation {
	return remove
}

func newRemoveAction(rects Rectangles) *removeAction {
	return &removeAction{
		newInsertAction(rects),
	}
}

func minUint64(choices ...uint64) uint64 {
	min := choices[0]
	for i := 1; i < len(choices); i++ {
		if choices[i] < min {
			min = choices[i]
		}
	}

	return min
}

type interfaces []interface{}

func executeInterfacesInParallel(ifs interfaces, fn func(interface{})) {
	if len(ifs) == 0 {
		return
	}

	done := int64(-1)
	numCPU := uint64(runtime.NumCPU())
	if numCPU > 1 {
		numCPU--
	}

	numCPU = minUint64(numCPU, uint64(len(ifs)))

	var wg sync.WaitGroup
	wg.Add(int(numCPU))

	for i := uint64(0); i < numCPU; i++ {
		go func() {
			defer wg.Done()

			for {
				i := atomic.AddInt64(&done, 1)
				if i >= int64(len(ifs)) {
					return
				}

				fn(ifs[i])
			}
		}()
	}

	wg.Wait()
}

func executeInterfacesInSerial(ifs interfaces, fn func(interface{})) {
	if len(ifs) == 0 {
		return
	}

	for _, ifc := range ifs {
		fn(ifc)
	}
}

// #region hilbert
func getCenter(rect Rectangle) (int32, int32) {
	xlow, ylow := rect.LowerLeft()
	xhigh, yhigh := rect.UpperRight()

	return (xhigh + xlow) / 2, (yhigh + ylow) / 2
}

type hilbertBundle struct {
	hilbert hilbert
	rect    Rectangle
}

func bundlesFromRects(rects ...Rectangle) []*hilbertBundle {
	chunks := chunkRectangles(rects, int64(runtime.NumCPU()))
	bundleChunks := make([][]*hilbertBundle, len(chunks))
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for i := 0; i < runtime.NumCPU(); i++ {
		if len(chunks[i]) == 0 {
			bundleChunks[i] = []*hilbertBundle{}
			wg.Done()
			continue
		}
		go func(i int) {
			bundles := make([]*hilbertBundle, 0, len(chunks[i]))
			for _, r := range chunks[i] {
				h := Encode(getCenter(r))
				bundles = append(bundles, &hilbertBundle{hilbert(h), r})
			}
			bundleChunks[i] = bundles
			wg.Done()
		}(i)
	}

	wg.Wait()

	bundles := make([]*hilbertBundle, 0, len(rects))
	for _, bc := range bundleChunks {
		bundles = append(bundles, bc...)
	}

	return bundles
}

// chunkRectangles takes a slice of rtree.Rectangle values and chunks it into `numParts` subslices.
func chunkRectangles(slice Rectangles, numParts int64) []Rectangles {
	parts := make([]Rectangles, numParts)
	for i := int64(0); i < numParts; i++ {
		parts[i] = slice[i*int64(len(slice))/numParts : (i+1)*int64(len(slice))/numParts]
	}
	return parts
}

// n defines the maximum power of 2 that can define a bound,
// this is the value for 2-d space if you want to support
// all hilbert ids with a single integer variable
const n = 1 << 31

// Encode will encode the provided x and y coordinates into a Hilbert
// distance.
func Encode(x, y int32) int64 {
	var rx, ry int32
	var d int64
	for s := int32(n / 2); s > 0; s /= 2 {
		rx = boolToInt(x&s > 0)
		ry = boolToInt(y&s > 0)
		d += int64(int64(s) * int64(s) * int64(((3 * rx) ^ ry)))
		rotate(s, rx, ry, &x, &y)
	}

	return d
}

// Decode will decode the provided Hilbert distance into a corresponding
// x and y value, respectively.
func Decode(h int64) (int32, int32) {
	var ry, rx int64
	var x, y int32
	t := h

	for s := int64(1); s < int64(n); s *= 2 {
		rx = 1 & (t / 2)
		ry = 1 & (t ^ rx)
		rotate(int32(s), int32(rx), int32(ry), &x, &y)
		x += int32(s * rx)
		y += int32(s * ry)
		t /= 4
	}

	return x, y
}

func boolToInt(value bool) int32 {
	if value {
		return int32(1)
	}

	return int32(0)
}

func rotate(n, rx, ry int32, x, y *int32) {
	if ry == 0 {
		if rx == 1 {
			*x = n - 1 - *x
			*y = n - 1 - *y
		}

		t := *x
		*x = *y
		*y = t
	}
}

// #endregion

// #region node
type hilbert int64

type hilberts []hilbert

func getParent(parent *node, key hilbert, r1 Rectangle) *node {
	var n *node
	for parent != nil && !parent.isLeaf {
		n = parent.searchNode(key)
		parent = n
	}

	if parent != nil && r1 != nil { // must be leaf and we need exact match
		// we are safe to travel to the right
		i := parent.search(key)
		for parent.keys.byPosition(i) == key {
			if equal(parent.nodes.list[i], r1) {
				break
			}

			i++
			if i == parent.keys.len() {
				if parent.right == nil { // we are far to the right
					break
				}

				if parent.right.keys.byPosition(0) != key {
					break
				}

				parent = parent.right
				i = 0
			}
		}
	}

	return parent
}

type nodes struct {
	list Rectangles
}

func (ns *nodes) push(n Rectangle) {
	ns.list = append(ns.list, n)
}

func (ns *nodes) splitAt(i, capacity uint64) (*nodes, *nodes) {
	i++
	right := make(Rectangles, uint64(len(ns.list))-i, capacity)
	copy(right, ns.list[i:])
	for j := i; j < uint64(len(ns.list)); j++ {
		ns.list[j] = nil
	}
	ns.list = ns.list[:i]
	return ns, &nodes{list: right}
}

func (ns *nodes) byPosition(pos uint64) *node {
	if pos >= uint64(len(ns.list)) {
		return nil
	}

	return ns.list[pos].(*node)
}

func (ns *nodes) insertAt(i uint64, n Rectangle) {
	ns.list = append(ns.list, nil)
	copy(ns.list[i+1:], ns.list[i:])
	ns.list[i] = n
}

func (ns *nodes) replaceAt(i uint64, n Rectangle) {
	ns.list[i] = n
}

func (ns *nodes) len() uint64 {
	return uint64(len(ns.list))
}

func (ns *nodes) deleteAt(i uint64) {
	copy(ns.list[i:], ns.list[i+1:])
	ns.list = ns.list[:len(ns.list)-1]
}

func newNodes(size uint64) *nodes {
	return &nodes{
		list: make(Rectangles, 0, size),
	}
}

type keys struct {
	list hilberts
}

func (ks *keys) splitAt(i, capacity uint64) (*keys, *keys) {
	i++
	right := make(hilberts, uint64(len(ks.list))-i, capacity)
	copy(right, ks.list[i:])
	ks.list = ks.list[:i]
	return ks, &keys{list: right}
}

func (ks *keys) len() uint64 {
	return uint64(len(ks.list))
}

func (ks *keys) byPosition(i uint64) hilbert {
	if i >= uint64(len(ks.list)) {
		return -1
	}
	return ks.list[i]
}

func (ks *keys) deleteAt(i uint64) {
	copy(ks.list[i:], ks.list[i+1:])
	ks.list = ks.list[:len(ks.list)-1]
}

func (ks *keys) delete(k hilbert) hilbert {
	i := ks.search(k)
	if i >= uint64(len(ks.list)) {
		return -1
	}

	if ks.list[i] != k {
		return -1
	}
	old := ks.list[i]
	ks.deleteAt(i)
	return old
}

func (ks *keys) search(key hilbert) uint64 {
	i := sort.Search(len(ks.list), func(i int) bool {
		return ks.list[i] >= key
	})

	return uint64(i)
}

func (ks *keys) insert(key hilbert) (hilbert, uint64) {
	i := ks.search(key)
	if i == uint64(len(ks.list)) {
		ks.list = append(ks.list, key)
		return -1, i
	}

	var old hilbert
	if ks.list[i] == key {
		old = ks.list[i]
		ks.list[i] = key
	} else {
		ks.insertAt(i, key)
	}

	return old, i
}

func (ks *keys) last() hilbert {
	return ks.list[len(ks.list)-1]
}

func (ks *keys) insertAt(i uint64, k hilbert) {
	ks.list = append(ks.list, -1)
	copy(ks.list[i+1:], ks.list[i:])
	ks.list[i] = k
}

func (ks *keys) withPosition(k hilbert) (hilbert, uint64) {
	i := ks.search(k)
	if i == uint64(len(ks.list)) {
		return -1, i
	}
	if ks.list[i] == k {
		return ks.list[i], i
	}

	return -1, i
}

func newKeys(size uint64) *keys {
	return &keys{
		list: make(hilberts, 0, size),
	}
}

type node struct {
	keys          *keys
	nodes         *nodes
	isLeaf        bool
	parent, right *node
	mbr           *rectangle
	maxHilbert    hilbert
}

func (n *node) insert(kb *keyBundle) Rectangle {
	i := n.keys.search(kb.key)
	if n.isLeaf { // we can have multiple keys with the same hilbert number
		for i < n.keys.len() && n.keys.list[i] == kb.key {
			if equal(n.nodes.list[i], kb.left) {
				old := n.nodes.list[i]
				n.nodes.list[i] = kb.left
				return old
			}
			i++
		}
	}

	if i == n.keys.len() {
		n.maxHilbert = kb.key
	}

	n.keys.insertAt(i, kb.key)
	if n.isLeaf {
		n.nodes.insertAt(i, kb.left)
	} else {
		if n.nodes.len() == 0 {
			n.nodes.push(kb.left)
			n.nodes.push(kb.right)
		} else {
			n.nodes.replaceAt(i, kb.left)
			n.nodes.insertAt(i+1, kb.right)
		}
		n.mbr.adjust(kb.left)
		n.mbr.adjust(kb.right)
		if kb.right.(*node).maxHilbert > n.maxHilbert {
			n.maxHilbert = kb.right.(*node).maxHilbert
		}
	}

	return nil
}

func (n *node) delete(kb *keyBundle) Rectangle {
	i := n.keys.search(kb.key)
	if n.keys.byPosition(i) != kb.key { // hilbert value not found
		return nil
	}

	if !equal(n.nodes.list[i], kb.left) {
		return nil
	}

	old := n.nodes.list[i]
	n.keys.deleteAt(i)
	n.nodes.deleteAt(i)
	return old
}

func (n *node) LowerLeft() (int32, int32) {
	return n.mbr.xlow, n.mbr.ylow
}

func (n *node) UpperRight() (int32, int32) {
	return n.mbr.xhigh, n.mbr.yhigh
}

func (n *node) needsSplit(ary uint64) bool {
	return n.keys.len() >= ary
}

func (n *node) splitLeaf(i, capacity uint64) (hilbert, *node, *node) {
	key := n.keys.byPosition(i)
	_, rightKeys := n.keys.splitAt(i, capacity)
	_, rightNodes := n.nodes.splitAt(i, capacity)
	nn := &node{
		keys:   rightKeys,
		nodes:  rightNodes,
		isLeaf: true,
		right:  n.right,
		parent: n.parent,
	}
	n.right = nn
	nn.mbr = newRectangleFromRects(rightNodes.list)
	n.mbr = newRectangleFromRects(n.nodes.list)
	nn.maxHilbert = rightKeys.last()
	n.maxHilbert = n.keys.last()
	return key, n, nn
}

func (n *node) splitInternal(i, capacity uint64) (hilbert, *node, *node) {
	key := n.keys.byPosition(i)
	n.keys.delete(key)

	_, rightKeys := n.keys.splitAt(i-1, capacity)
	_, rightNodes := n.nodes.splitAt(i, capacity)

	nn := newNode(false, rightKeys, rightNodes)
	for _, n := range rightNodes.list {
		n.(*node).parent = nn
	}
	nn.mbr = newRectangleFromRects(rightNodes.list)
	n.mbr = newRectangleFromRects(n.nodes.list)
	nn.maxHilbert = nn.keys.last()
	n.maxHilbert = n.keys.last()

	return key, n, nn
}

func (n *node) split(i, capacity uint64) (hilbert, *node, *node) {
	if n.isLeaf {
		return n.splitLeaf(i, capacity)
	}

	return n.splitInternal(i, capacity)
}

func (n *node) search(key hilbert) uint64 {
	return n.keys.search(key)
}

func (n *node) searchNode(key hilbert) *node {
	i := n.search(key)

	return n.nodes.byPosition(uint64(i))
}

func (n *node) searchRects(r *rectangle) Rectangles {
	rects := make(Rectangles, 0, n.nodes.len())
	for _, child := range n.nodes.list {
		if intersect(r, child) {
			rects = append(rects, child)
		}
	}

	return rects
}

func (n *node) key() hilbert {
	return n.keys.last()
}

func newNode(isLeaf bool, keys *keys, ns *nodes) *node {
	return &node{
		isLeaf: isLeaf,
		keys:   keys,
		nodes:  ns,
	}
}

// #endregion

// #region rectangle
type rectangle struct {
	xlow, xhigh, ylow, yhigh int32
}

func (r *rectangle) adjust(rect Rectangle) {
	x, y := rect.LowerLeft()
	if x < r.xlow {
		r.xlow = x
	}
	if y < r.ylow {
		r.ylow = y
	}

	x, y = rect.UpperRight()
	if x > r.xhigh {
		r.xhigh = x
	}

	if y > r.yhigh {
		r.yhigh = y
	}
}

func equal(r1, r2 Rectangle) bool {
	xlow1, ylow1 := r1.LowerLeft()
	xhigh2, yhigh2 := r2.UpperRight()

	xhigh1, yhigh1 := r1.UpperRight()
	xlow2, ylow2 := r2.LowerLeft()

	return xlow1 == xlow2 && xhigh1 == xhigh2 && ylow1 == ylow2 && yhigh1 == yhigh2
}

func intersect(rect1 *rectangle, rect2 Rectangle) bool {
	xhigh2, yhigh2 := rect2.UpperRight()
	xlow2, ylow2 := rect2.LowerLeft()

	return xhigh2 >= rect1.xlow && xlow2 <= rect1.xhigh && yhigh2 >= rect1.ylow && ylow2 <= rect1.yhigh
}

func newRectangeFromRect(rect Rectangle) *rectangle {
	r := &rectangle{}
	x, y := rect.LowerLeft()
	r.xlow = x
	r.ylow = y

	x, y = rect.UpperRight()
	r.xhigh = x
	r.yhigh = y

	return r
}

func newRectangleFromRects(rects Rectangles) *rectangle {
	if len(rects) == 0 {
		panic(`Cannot construct rectangle with no dimensions.`)
	}

	xlow, ylow := rects[0].LowerLeft()
	xhigh, yhigh := rects[0].UpperRight()
	r := &rectangle{
		xlow:  xlow,
		xhigh: xhigh,
		ylow:  ylow,
		yhigh: yhigh,
	}

	for i := 1; i < len(rects); i++ {
		r.adjust(rects[i])
	}

	return r
}

// #endregion

// #region ringbuffer

// roundUp takes a uint64 greater than 0 and rounds it up to the next
// power of 2.
func roundUp(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

type rnode struct {
	position uint64
	data     interface{}
}

type rnodes []rnode

// RingBuffer is a MPMC buffer that achieves threadsafety with CAS operations
// only.  A put on full or get on empty call will block until an item
// is put or retrieved.  Calling Dispose on the RingBuffer will unblock
// any blocked threads with an error.  This buffer is similar to the buffer
// described here: http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
// with some minor additions.
type RingBuffer struct {
	_padding0      [8]uint64
	queue          uint64
	_padding1      [8]uint64
	dequeue        uint64
	_padding2      [8]uint64
	mask, disposed uint64
	_padding3      [8]uint64
	nodes          rnodes
}

func (rb *RingBuffer) init(size uint64) {
	size = roundUp(size)
	rb.nodes = make(rnodes, size)
	for i := uint64(0); i < size; i++ {
		rb.nodes[i] = rnode{position: i}
	}
	rb.mask = size - 1 // so we don't have to do this with every put/get operation
}

// Put adds the provided item to the queue.  If the queue is full, this
// call will block until an item is added to the queue or Dispose is called
// on the queue.  An error will be returned if the queue is disposed.
func (rb *RingBuffer) Put(item interface{}) error {
	_, err := rb.put(item, false)
	return err
}

// Offer adds the provided item to the queue if there is space.  If the queue
// is full, this call will return false.  An error will be returned if the
// queue is disposed.
func (rb *RingBuffer) Offer(item interface{}) (bool, error) {
	return rb.put(item, true)
}

func (rb *RingBuffer) put(item interface{}, offer bool) (bool, error) {
	var n *rnode
	pos := atomic.LoadUint64(&rb.queue)
L:
	for {
		if atomic.LoadUint64(&rb.disposed) == 1 {
			return false, ErrDisposed
		}

		n = &rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := seq - pos; {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.queue, pos, pos+1) {
				break L
			}
		case dif < 0:
			panic(`Ring buffer in a compromised state during a put operation.`)
		default:
			pos = atomic.LoadUint64(&rb.queue)
		}

		if offer {
			return false, nil
		}

		runtime.Gosched() // free up the cpu before the next iteration
	}

	n.data = item
	atomic.StoreUint64(&n.position, pos+1)
	return true, nil
}

// Get will return the next item in the queue.  This call will block
// if the queue is empty.  This call will unblock when an item is added
// to the queue or Dispose is called on the queue.  An error will be returned
// if the queue is disposed.
func (rb *RingBuffer) Get() (interface{}, error) {
	return rb.Poll(0)
}

// Poll will return the next item in the queue.  This call will block
// if the queue is empty.  This call will unblock when an item is added
// to the queue, Dispose is called on the queue, or the timeout is reached. An
// error will be returned if the queue is disposed or a timeout occurs. A
// non-positive timeout will block indefinitely.
func (rb *RingBuffer) Poll(timeout time.Duration) (interface{}, error) {
	var (
		n     *rnode
		pos   = atomic.LoadUint64(&rb.dequeue)
		start time.Time
	)
	if timeout > 0 {
		start = time.Now()
	}
L:
	for {
		if atomic.LoadUint64(&rb.disposed) == 1 {
			return nil, ErrDisposed
		}

		n = &rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := seq - (pos + 1); {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.dequeue, pos, pos+1) {
				break L
			}
		case dif < 0:
			panic(`Ring buffer in compromised state during a get operation.`)
		default:
			pos = atomic.LoadUint64(&rb.dequeue)
		}

		if timeout > 0 && time.Since(start) >= timeout {
			return nil, ErrTimeout
		}

		runtime.Gosched() // free up the cpu before the next iteration
	}
	data := n.data
	n.data = nil
	atomic.StoreUint64(&n.position, pos+rb.mask+1)
	return data, nil
}

// Len returns the number of items in the queue.
func (rb *RingBuffer) Len() uint64 {
	return atomic.LoadUint64(&rb.queue) - atomic.LoadUint64(&rb.dequeue)
}

// Cap returns the capacity of this ring buffer.
func (rb *RingBuffer) Cap() uint64 {
	return uint64(len(rb.nodes))
}

// Dispose will dispose of this queue and free any blocked threads
// in the Put and/or Get methods.  Calling those methods on a disposed
// queue will return an error.
func (rb *RingBuffer) Dispose() {
	atomic.CompareAndSwapUint64(&rb.disposed, 0, 1)
}

// IsDisposed will return a bool indicating if this queue has been
// disposed.
func (rb *RingBuffer) IsDisposed() bool {
	return atomic.LoadUint64(&rb.disposed) == 1
}

// NewRingBuffer will allocate, initialize, and return a ring buffer
// with the specified size.
func NewRingBuffer(size uint64) *RingBuffer {
	rb := &RingBuffer{}
	rb.init(size)
	return rb
}

var (
	// ErrDisposed is returned when an operation is performed on a disposed
	// queue.
	ErrDisposed = errors.New(`queue: disposed`)

	// ErrTimeout is returned when an applicable queue operation times out.
	ErrTimeout = errors.New(`queue: poll timed out`)

	// ErrEmptyQueue is returned when an non-applicable queue operation was called
	// due to the queue's empty item state
	ErrEmptyQueue = errors.New(`queue: empty queue`)
)

// #endregion
