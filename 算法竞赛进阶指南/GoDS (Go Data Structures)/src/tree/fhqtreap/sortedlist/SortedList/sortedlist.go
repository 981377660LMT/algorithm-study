// 基于fhq-treap实现
package main

import (
	"fmt"
	"strings"
	"time"
)

// TODO O(n) 笛卡尔建树
func NewSortedList(comparator func(a, b interface{}) int, initCapacity int) *SortedList {
	return &SortedList{
		seed:       uint(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, max(initCapacity, 16)),
	}
}

type node struct {
	left, right int
	priority    uint
	size        int
	value       interface{}
}

type SortedList struct {
	seed       uint
	nodeId     int
	root       int
	comparator func(a, b interface{}) int
	nodes      []node
}

func (sl *SortedList) Add(value interface{}) {
	sl.resureCapacity(sl.nodeId + 2)
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &y, false)
	z = sl.newNode(value)
	sl.root = sl.merge(sl.merge(x, z), y)
}

func (sl *SortedList) At(index int) interface{} {
	n := sl.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("%d index out of range: [%d,%d]", index, 0, n-1))
	}
	return sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (sl *SortedList) Pop(index int) interface{} {
	popped := sl.At(index)
	sl.Discard(popped)
	return popped
}

func (sl *SortedList) Discard(value interface{}) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &z, false)
	sl.splitByValue(x, value, &x, &y, true)
	y = sl.merge(sl.nodes[y].left, sl.nodes[y].right)
	sl.root = sl.merge(sl.merge(x, y), z)
}

func (sl *SortedList) BisectLeft(value interface{}) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, true)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) BisectRight(value interface{}) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, false)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *SortedList) String() string {
	sb := []string{"SortedList{"}
	values := []string{}
	for i := 0; i < sl.Len(); i++ {
		values = append(values, fmt.Sprintf("%d", sl.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")
}

func (sl *SortedList) Len() int {
	return sl.nodes[sl.root].size
}

func (sl *SortedList) kthNode(root int, k int) int {
	cur := root
	for cur != 0 {
		if sl.nodes[sl.nodes[cur].left].size+1 == k {
			break
		} else if sl.nodes[sl.nodes[cur].left].size >= k {
			cur = sl.nodes[cur].left
		} else {
			k -= sl.nodes[sl.nodes[cur].left].size + 1
			cur = sl.nodes[cur].right
		}
	}
	return cur
}

func (sl *SortedList) splitByValue(root int, value interface{}, x, y *int, strictLess bool) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	if strictLess {
		if sl.comparator(sl.nodes[root].value, value) < 0 {
			*x = root
			sl.splitByValue(sl.nodes[root].right, value, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, value, x, &sl.nodes[root].left, strictLess)
		}
	} else {
		if sl.comparator(sl.nodes[root].value, value) <= 0 {
			*x = root
			sl.splitByValue(sl.nodes[root].right, value, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, value, x, &sl.nodes[root].left, strictLess)
		}
	}

	sl.pushUp(root)
}

func (sl *SortedList) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}
	if sl.nodes[x].priority < sl.nodes[y].priority {
		sl.nodes[x].right = sl.merge(sl.nodes[x].right, y)
		sl.pushUp(x)
		return x
	}
	sl.nodes[y].left = sl.merge(x, sl.nodes[y].left)
	sl.pushUp(y)
	return y
}

func (sl *SortedList) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (sl *SortedList) newNode(value interface{}) int {
	sl.nodeId++
	index := sl.nodeId
	sl.nodes[index].value = value
	sl.nodes[index].priority = sl.fastRand()
	sl.nodes[index].size = 1
	return index
}

func (sl *SortedList) fastRand() uint {
	sl.seed ^= sl.seed << 13
	sl.seed ^= sl.seed >> 17
	sl.seed ^= sl.seed << 5
	return sl.seed
}

func (sl *SortedList) resureCapacity(needCap int) {
	if needCap > len(sl.nodes) {
		sl.resize(needCap)
	}
}

func (sl *SortedList) resize(needCap int) {
	newCap := len(sl.nodes) * 2
	if newCap < needCap {
		newCap = needCap
	}
	newNodes := make([]node, newCap)
	copy(newNodes, sl.nodes)
	sl.nodes = newNodes
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
