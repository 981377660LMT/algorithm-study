// https://github.com/huandu/skiplist
// !Not thread safe.
//
// Package skiplist 是一个跳表的实现，用于按递增顺序存储元素。
// 它允许在大约 O(n log(n)) 的时间复杂度内进行查找、插入和删除操作。
// 此外，还有用于获取下一个和上一个元素的方法，以及在无需重新插入的情况下更改实际值的方法
// （前提是键保持不变！）
// 跳表是平衡树的快速替代方案。

package main

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"time"
)

type Element int

func (e Element) ExtractKey() float64 {
	return float64(e)
}
func (e Element) String() string {
	return fmt.Sprintf("%03d", e)
}

func main() {
	list := New()

	// Insert some elements
	for i := 0; i < 10; i++ {
		list.Insert(Element(i))
	}

	// Find an element
	if e, ok := list.Find(Element(53)); ok {
		fmt.Println(e.GetValue())
	}

	// Delete all elements
	for i := 0; i < 10; i++ {
		list.Delete(Element(i))
	}
}

const (
	// 跳表的最大层数为 25。这一高度使得跳表对于最多 3400 万（34m）条目仍然高效。如果需要存储更多元素，可以适当调整此常量
	maxLevel = 25
	// 定义一个极小值（ε），用于在比较键值时判断两个键是否相等。由于浮点数计算可能存在精度问题，使用 eps 来确定键值的近似相等
	eps = 0.00001
)

// 跳表中的元素需要实现的接口.
type ListElement interface {
	ExtractKey() float64
	String() string
}

type SkipListElement struct {
	next  [maxLevel]*SkipListElement // 指向各层下一个节点的指针
	level int
	key   float64
	value ListElement
	prev  *SkipListElement // 指向前一个节点的指针，仅在底层（第 0 层）中使用，便于双向遍历
}

type SkipList struct {
	startLevels  [maxLevel]*SkipListElement
	endLevels    [maxLevel]*SkipListElement
	maxNewLevel  int
	maxLevel     int
	elementCount int
	eps          float64
}

func NewSeedEps(seed int64, eps float64) SkipList {

	// Initialize random number generator.
	rand.Seed(seed)
	//fmt.Printf("SkipList seed: %v\n", seed)

	list := SkipList{
		startLevels:  [maxLevel]*SkipListElement{},
		endLevels:    [maxLevel]*SkipListElement{},
		maxNewLevel:  maxLevel,
		maxLevel:     0,
		elementCount: 0,
		eps:          eps,
	}

	return list
}

func NewEps(eps float64) SkipList {
	return NewSeedEps(time.Now().UTC().UnixNano(), eps)
}

func NewSeed(seed int64) SkipList {
	return NewSeedEps(seed, eps)
}

func New() SkipList {
	return NewSeedEps(time.Now().UTC().UnixNano(), eps)
}

// IsEmpty checks, if the skiplist is empty.
func (t *SkipList) IsEmpty() bool {
	return t.startLevels[0] == nil
}

// 生成一个随机层数，用于新插入节点的层高.
func (t *SkipList) generateLevel(maxLevel int) int {
	level := maxLevel - 1
	// First we apply some mask which makes sure that we don't get a level
	// above our desired level. Then we find the first set bit.
	var x uint64 = rand.Uint64() & ((1 << uint(maxLevel-1)) - 1)
	zeroes := bits.TrailingZeros64(x)
	if zeroes <= maxLevel {
		level = zeroes
	}

	return level
}

func (t *SkipList) findEntryIndex(key float64, level int) int {
	// Find good entry point so we don't accidentally skip half the list.
	for i := t.maxLevel; i >= 0; i-- {
		if t.startLevels[i] != nil && t.startLevels[i].key <= key || i <= level {
			return i
		}
	}
	return 0
}

func (t *SkipList) findExtended(key float64, findGreaterOrEqual bool) (foundElem *SkipListElement, ok bool) {
	foundElem = nil
	ok = false
	if t.IsEmpty() {
		return
	}

	index := t.findEntryIndex(key, 0)
	var currentNode *SkipListElement

	currentNode = t.startLevels[index]
	nextNode := currentNode

	// In case, that our first element is already greater-or-equal!
	if findGreaterOrEqual && currentNode.key > key {
		foundElem = currentNode
		ok = true
		return
	}

	for {
		if math.Abs(currentNode.key-key) <= t.eps {
			foundElem = currentNode
			ok = true
			return
		}

		nextNode = currentNode.next[index]

		// Which direction are we continuing next time?
		if nextNode != nil && nextNode.key <= key {
			// Go right
			currentNode = nextNode
		} else {
			if index > 0 {
				// Early exit
				if currentNode.next[0] != nil && math.Abs(currentNode.next[0].key-key) <= t.eps {
					foundElem = currentNode.next[0]
					ok = true
					return
				}
				// Go down
				index--
			} else {
				// Element is not found and we reached the bottom.
				if findGreaterOrEqual {
					foundElem = nextNode
					ok = nextNode != nil
				}

				return
			}
		}
	}
}

// Find tries to find an element in the skiplist based on the key from the given ListElement.
// elem can be used, if ok is true.
// Find runs in approx. O(log(n))
func (t *SkipList) Find(e ListElement) (elem *SkipListElement, ok bool) {
	if t == nil || e == nil {
		return
	}

	elem, ok = t.findExtended(e.ExtractKey(), false)
	return
}

// FindGreaterOrEqual finds the first element, that is greater or equal to the given ListElement e.
// The comparison is done on the keys (So on ExtractKey()).
// FindGreaterOrEqual runs in approx. O(log(n))
func (t *SkipList) FindGreaterOrEqual(e ListElement) (elem *SkipListElement, ok bool) {
	if t == nil || e == nil {
		return
	}

	elem, ok = t.findExtended(e.ExtractKey(), true)
	return
}

// Delete removes an element equal to e from the skiplist, if there is one.
// If there are multiple entries with the same value, Delete will remove one of them
// (Which one will change based on the actual skiplist layout)
// Delete runs in approx. O(log(n))
func (t *SkipList) Delete(e ListElement) {
	if t == nil || t.IsEmpty() || e == nil {
		return
	}

	key := e.ExtractKey()

	index := t.findEntryIndex(key, 0)

	var currentNode *SkipListElement
	nextNode := currentNode

	for {
		if currentNode == nil {
			nextNode = t.startLevels[index]
		} else {
			nextNode = currentNode.next[index]
		}

		// Found and remove!
		if nextNode != nil && math.Abs(nextNode.key-key) <= t.eps {

			if currentNode != nil {
				currentNode.next[index] = nextNode.next[index]
			}

			if index == 0 {
				if nextNode.next[index] != nil {
					nextNode.next[index].prev = currentNode
				}
				t.elementCount--
			}

			// Link from start needs readjustments.
			if t.startLevels[index] == nextNode {
				t.startLevels[index] = nextNode.next[index]
				// This was our currently highest node!
				if t.startLevels[index] == nil {
					t.maxLevel = index - 1
				}
			}

			// Link from end needs readjustments.
			if nextNode.next[index] == nil {
				t.endLevels[index] = currentNode
			}
			nextNode.next[index] = nil
		}

		if nextNode != nil && nextNode.key < key {
			// Go right
			currentNode = nextNode
		} else {
			// Go down
			index--
			if index < 0 {
				break
			}
		}
	}

}

// Insert inserts the given ListElement into the skiplist.
// Insert runs in approx. O(log(n))
func (t *SkipList) Insert(e ListElement) {
	if t == nil || e == nil {
		return
	}

	level := t.generateLevel(t.maxNewLevel)

	// Only grow the height of the skiplist by one at a time!
	if level > t.maxLevel {
		level = t.maxLevel + 1
		t.maxLevel = level
	}

	elem := &SkipListElement{
		next:  [maxLevel]*SkipListElement{},
		level: level,
		key:   e.ExtractKey(),
		value: e,
	}

	t.elementCount++

	newFirst := true
	newLast := true
	if !t.IsEmpty() {
		newFirst = elem.key < t.startLevels[0].key
		newLast = elem.key > t.endLevels[0].key
	}

	normallyInserted := false
	if !newFirst && !newLast {

		normallyInserted = true

		index := t.findEntryIndex(elem.key, level)

		var currentNode *SkipListElement
		nextNode := t.startLevels[index]

		for {
			if currentNode == nil {
				nextNode = t.startLevels[index]
			} else {
				nextNode = currentNode.next[index]
			}

			// Connect node to next
			if index <= level && (nextNode == nil || nextNode.key > elem.key) {
				elem.next[index] = nextNode
				if currentNode != nil {
					currentNode.next[index] = elem
				}
				if index == 0 {
					elem.prev = currentNode
					if nextNode != nil {
						nextNode.prev = elem
					}
				}
			}

			if nextNode != nil && nextNode.key <= elem.key {
				// Go right
				currentNode = nextNode
			} else {
				// Go down
				index--
				if index < 0 {
					break
				}
			}
		}
	}

	// Where we have a left-most position that needs to be referenced!
	for i := level; i >= 0; i-- {
		didSomething := false

		if newFirst || normallyInserted {
			if t.startLevels[i] == nil || t.startLevels[i].key > elem.key {
				if i == 0 && t.startLevels[i] != nil {
					t.startLevels[i].prev = elem
				}
				elem.next[i] = t.startLevels[i]
				t.startLevels[i] = elem
			}

			// link the endLevels to this element!
			if elem.next[i] == nil {
				t.endLevels[i] = elem
			}

			didSomething = true
		}

		if newLast {
			// Places the element after the very last element on this level!
			// This is very important, so we are not linking the very first element (newFirst AND newLast) to itself!
			if !newFirst {
				if t.endLevels[i] != nil {
					t.endLevels[i].next[i] = elem
				}
				if i == 0 {
					elem.prev = t.endLevels[i]
				}
				t.endLevels[i] = elem
			}

			// Link the startLevels to this element!
			if t.startLevels[i] == nil || t.startLevels[i].key > elem.key {
				t.startLevels[i] = elem
			}

			didSomething = true
		}

		if !didSomething {
			break
		}
	}
}

// GetValue extracts the ListElement value from a skiplist node.
func (e *SkipListElement) GetValue() ListElement {
	return e.value
}

// GetSmallestNode returns the very first/smallest node in the skiplist.
// GetSmallestNode runs in O(1)
func (t *SkipList) GetSmallestNode() *SkipListElement {
	return t.startLevels[0]
}

// GetLargestNode returns the very last/largest node in the skiplist.
// GetLargestNode runs in O(1)
func (t *SkipList) GetLargestNode() *SkipListElement {
	return t.endLevels[0]
}

// Next returns the next element based on the given node.
// Next will loop around to the first node, if you call it on the last!
func (t *SkipList) Next(e *SkipListElement) *SkipListElement {
	if e.next[0] == nil {
		return t.startLevels[0]
	}
	return e.next[0]
}

// Prev returns the previous element based on the given node.
// Prev will loop around to the last node, if you call it on the first!
func (t *SkipList) Prev(e *SkipListElement) *SkipListElement {
	if e.prev == nil {
		return t.endLevels[0]
	}
	return e.prev
}

// GetNodeCount returns the number of nodes currently in the skiplist.
func (t *SkipList) GetNodeCount() int {
	return t.elementCount
}

// ChangeValue can be used to change the actual value of a node in the skiplist
// without the need of Deleting and reinserting the node again.
// Be advised, that ChangeValue only works, if the actual key from ExtractKey() will stay the same!
// ok is an indicator, wether the value is actually changed.
func (t *SkipList) ChangeValue(e *SkipListElement, newValue ListElement) (ok bool) {
	// The key needs to stay correct, so this is very important!
	if math.Abs(newValue.ExtractKey()-e.key) <= t.eps {
		e.value = newValue
		ok = true
	} else {
		ok = false
	}
	return
}

// String returns a string format of the skiplist. Useful to get a graphical overview and/or debugging.
func (t *SkipList) String() string {
	s := ""

	s += " --> "
	for i, l := range t.startLevels {
		if l == nil {
			break
		}
		if i > 0 {
			s += " -> "
		}
		next := "---"
		if l != nil {
			next = l.value.String()
		}
		s += fmt.Sprintf("[%v]", next)

		if i == 0 {
			s += "    "
		}
	}
	s += "\n"

	node := t.startLevels[0]
	for node != nil {
		s += fmt.Sprintf("%v: ", node.value)
		for i := 0; i <= node.level; i++ {

			l := node.next[i]

			next := "---"
			if l != nil {
				next = l.value.String()
			}

			if i == 0 {
				prev := "---"
				if node.prev != nil {
					prev = node.prev.value.String()
				}
				s += fmt.Sprintf("[%v|%v]", prev, next)
			} else {
				s += fmt.Sprintf("[%v]", next)
			}
			if i < node.level {
				s += " -> "
			}

		}
		s += "\n"
		node = node.next[0]
	}

	s += " --> "
	for i, l := range t.endLevels {
		if l == nil {
			break
		}
		if i > 0 {
			s += " -> "
		}
		next := "---"
		if l != nil {
			next = l.value.String()
		}
		s += fmt.Sprintf("[%v]", next)
		if i == 0 {
			s += "    "
		}
	}
	s += "\n"
	return s
}
