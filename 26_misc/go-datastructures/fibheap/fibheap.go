package main

import (
	"fmt"
	"math"
)

/******************************************
 ************** INTERFACE *****************
 ******************************************/

// FloatingFibonacciHeap is an implementation of a fibonacci heap
// with only floating-point priorities and no user data attached.
type FloatingFibonacciHeap struct {
	min  *Entry // The minimal element
	size uint   // Size of the heap
}

// Entry is the entry type that will be used
// for each node of the Fibonacci heap
type Entry struct {
	degree                    int
	marked                    bool
	next, prev, child, parent *Entry
	// Priority is the numerical priority of the node
	Priority float64
}

// EmptyHeapError fires when the heap is empty and an operation could
// not be completed for that reason. Its string holds additional data.
type EmptyHeapError string

func (e EmptyHeapError) Error() string {
	return string(e)
}

// NilError fires when a heap or entry is nil and an operation could
// not be completed for that reason. Its string holds additional data.
type NilError string

func (e NilError) Error() string {
	return string(e)
}

// NewFloatFibHeap creates a new, empty, Fibonacci heap object.
func NewFloatFibHeap() FloatingFibonacciHeap { return FloatingFibonacciHeap{nil, 0} }

// Enqueue adds and element to the heap
func (heap *FloatingFibonacciHeap) Enqueue(priority float64) *Entry {
	singleton := newEntry(priority)

	// Merge singleton list with heap
	heap.min = mergeLists(heap.min, singleton)
	heap.size++
	return singleton
}

// Min returns the minimum element in the heap
func (heap *FloatingFibonacciHeap) Min() (*Entry, error) {
	if heap.IsEmpty() {
		return nil, EmptyHeapError("Trying to get minimum element of empty heap")
	}
	return heap.min, nil
}

// IsEmpty answers: is the heap empty?
func (heap *FloatingFibonacciHeap) IsEmpty() bool {
	return heap.size == 0
}

// Size gives the number of elements in the heap
func (heap *FloatingFibonacciHeap) Size() uint {
	return heap.size
}

// DequeueMin removes and returns the
// minimal element in the heap
func (heap *FloatingFibonacciHeap) DequeueMin() (*Entry, error) {
	if heap.IsEmpty() {
		return nil, EmptyHeapError("Cannot dequeue minimum of empty heap")
	}

	heap.size--

	// Copy pointer. Will need it later.
	min := heap.min

	if min.next == min { // This is the only root node
		heap.min = nil
	} else { // There are more root nodes
		heap.min.prev.next = heap.min.next
		heap.min.next.prev = heap.min.prev
		heap.min = heap.min.next // Arbitrary element of the root list
	}

	if min.child != nil {
		// Keep track of the first visited node
		curr := min.child
		for ok := true; ok; ok = (curr != min.child) {
			curr.parent = nil
			curr = curr.next
		}
	}

	heap.min = mergeLists(heap.min, min.child)

	if heap.min == nil {
		// If there are no entries left, we're done.
		return min, nil
	}

	treeSlice := make([]*Entry, 0, heap.size)
	toVisit := make([]*Entry, 0, heap.size)

	for curr := heap.min; len(toVisit) == 0 || toVisit[0] != curr; curr = curr.next {
		toVisit = append(toVisit, curr)
	}

	for _, curr := range toVisit {
		for {
			for curr.degree >= len(treeSlice) {
				treeSlice = append(treeSlice, nil)
			}

			if treeSlice[curr.degree] == nil {
				treeSlice[curr.degree] = curr
				break
			}

			other := treeSlice[curr.degree]
			treeSlice[curr.degree] = nil

			// Determine which of two trees has the smaller root
			var minT, maxT *Entry
			if other.Priority < curr.Priority {
				minT = other
				maxT = curr
			} else {
				minT = curr
				maxT = other
			}

			// Break max out of the root list,
			// then merge it into min's child list
			maxT.next.prev = maxT.prev
			maxT.prev.next = maxT.next

			// Make it a singleton so that we can merge it
			maxT.prev = maxT
			maxT.next = maxT
			minT.child = mergeLists(minT.child, maxT)

			// Reparent max appropriately
			maxT.parent = minT

			// Clear max's mark, since it can now lose another child
			maxT.marked = false

			// Increase min's degree. It has another child.
			minT.degree++

			// Continue merging this tree
			curr = minT
		}

		/* Update the global min based on this node.  Note that we compare
		 * for <= instead of < here.  That's because if we just did a
		 * reparent operation that merged two different trees of equal
		 * priority, we need to make sure that the min pointer points to
		 * the root-level one.
		 */
		if curr.Priority <= heap.min.Priority {
			heap.min = curr
		}
	}

	return min, nil
}

// DecreaseKey decreases the key of the given element, sets it to the new
// given priority and returns the node if successfully set
func (heap *FloatingFibonacciHeap) DecreaseKey(node *Entry, newPriority float64) (*Entry, error) {

	if heap.IsEmpty() {
		return nil, EmptyHeapError("Cannot decrease key in an empty heap")
	}

	if node == nil {
		return nil, NilError("Cannot decrease key: given node is nil")
	}

	if newPriority >= node.Priority {
		return nil, fmt.Errorf("The given new priority: %v, is larger than or equal to the old: %v",
			newPriority, node.Priority)
	}

	decreaseKeyUnchecked(heap, node, newPriority)
	return node, nil
}

// Delete deletes the given element in the heap
func (heap *FloatingFibonacciHeap) Delete(node *Entry) error {

	if heap.IsEmpty() {
		return EmptyHeapError("Cannot delete element from an empty heap")
	}

	if node == nil {
		return NilError("Cannot delete node: given node is nil")
	}

	decreaseKeyUnchecked(heap, node, -math.MaxFloat64)
	heap.DequeueMin()
	return nil
}

// Merge returns a new Fibonacci heap that contains
// all of the elements of the two heaps.  Each of the input heaps is
// destructively modified by having all its elements removed.  You can
// continue to use those heaps, but be aware that they will be empty
// after this call completes.
func (heap *FloatingFibonacciHeap) Merge(other *FloatingFibonacciHeap) (FloatingFibonacciHeap, error) {

	if heap == nil || other == nil {
		return FloatingFibonacciHeap{}, NilError("One of the heaps to merge is nil. Cannot merge")
	}

	resultSize := heap.size + other.size

	resultMin := mergeLists(heap.min, other.min)

	heap.min = nil
	other.min = nil
	heap.size = 0
	other.size = 0

	return FloatingFibonacciHeap{resultMin, resultSize}, nil
}

/******************************************
 ************** END INTERFACE *************
 ******************************************/

// ****************
// HELPER FUNCTIONS
// ****************

func newEntry(priority float64) *Entry {
	result := new(Entry)
	result.degree = 0
	result.marked = false
	result.child = nil
	result.parent = nil
	result.next = result
	result.prev = result
	result.Priority = priority
	return result
}

func mergeLists(one, two *Entry) *Entry {
	if one == nil && two == nil {
		return nil
	} else if one != nil && two == nil {
		return one
	} else if one == nil && two != nil {
		return two
	}
	// Both trees non-null; actually do the merge.
	oneNext := one.next
	one.next = two.next
	one.next.prev = one
	two.next = oneNext
	two.next.prev = two

	if one.Priority < two.Priority {
		return one
	}
	return two

}

func decreaseKeyUnchecked(heap *FloatingFibonacciHeap, node *Entry, priority float64) {
	node.Priority = priority

	if node.parent != nil && node.Priority <= node.parent.Priority {
		cutNode(heap, node)
	}

	if node.Priority <= heap.min.Priority {
		heap.min = node
	}
}

func cutNode(heap *FloatingFibonacciHeap, node *Entry) {
	node.marked = false

	if node.parent == nil {
		return
	}

	// Rewire siblings if it has any
	if node.next != node {
		node.next.prev = node.prev
		node.prev.next = node.next
	}

	// Rewrite pointer if this is the representative child node
	if node.parent.child == node {
		if node.next != node {
			node.parent.child = node.next
		} else {
			node.parent.child = nil
		}
	}

	node.parent.degree--

	node.prev = node
	node.next = node
	heap.min = mergeLists(heap.min, node)

	// cut parent recursively if marked
	if node.parent.marked {
		cutNode(heap, node.parent)
	} else {
		node.parent.marked = true
	}

	node.parent = nil
}
