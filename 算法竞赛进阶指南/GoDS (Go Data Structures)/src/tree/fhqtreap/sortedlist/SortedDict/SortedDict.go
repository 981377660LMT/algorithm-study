package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	mp := NewSortedDict(func(a, b Key) int {
		if a[0] != b[0] {
			return a[0] - b[0]
		}
		return a[1] - b[1]
	}, 16)
	mp.Set(Key{0, 1}, Value{"foo1", "bar1"})
	fmt.Println(mp)
	mp.Set(Key{-1, 2}, Value{"foo2", "bar2"})
	fmt.Println(mp.PeekItem(0))
	fmt.Println(mp.PopItem(1))
	fmt.Println(mp)
}

type Key = [2]int
type Value = struct{ foo, bar string } // 不能是指针类型

type SortedDict struct {
	sl *sortedList
	mp map[Key]Value
}

// key 为 int, value 为 Item 的有序映射.
func NewSortedDict(comparator func(a, b Key) int, initCapacity int) *SortedDict {
	return &SortedDict{
		sl: newSortedList(comparator, initCapacity),
		mp: make(map[Key]Value, initCapacity),
	}
}

func (sm *SortedDict) Set(key Key, value Value) {
	if _, has := sm.mp[key]; has {
		sm.mp[key] = value
		return
	}
	sm.mp[key] = value
	sm.sl.Add(key)
}

func (sm *SortedDict) Has(key Key) bool {
	_, has := sm.mp[key]
	return has
}

func (sm *SortedDict) Get(key Key) Value {
	return sm.mp[key]
}

func (sm *SortedDict) GetOrDefault(key Key, defaultValue Value) Value {
	if v, has := sm.mp[key]; has {
		return v
	}
	return defaultValue
}

func (sm *SortedDict) Discard(key Key) bool {
	if _, has := sm.mp[key]; has {
		delete(sm.mp, key)
		sm.sl.Discard(key)
		return true
	}
	return false
}

// `Erase` in cpp
func (sm *SortedDict) PopItem(iter int) (key Key, value Value) {
	key = sm.sl.At(iter)
	sm.sl.Pop(iter)
	value = sm.mp[key]
	delete(sm.mp, key)
	return
}

func (sm *SortedDict) PeekItem(iter int) (key Key, value Value) {
	key = sm.sl.At(iter)
	value = sm.mp[key]
	return
}

func (sm *SortedDict) LowerBound(key_ Key) (iter int, key Key, value Value) {
	iter = sm.sl.BisectLeft(key_)
	key = sm.sl.At(iter)
	value = sm.mp[key]
	return
}

func (sm *SortedDict) UpperBound(key_ Key) (iter int, key Key, value Value) {
	iter = sm.sl.BisectRight(key_)
	key = sm.sl.At(iter)
	value = sm.mp[key]
	return
}

func (sm *SortedDict) Size() int {
	return len(sm.mp)
}

func (sm *SortedDict) String() string {
	var b []string
	for i := 0; i < sm.sl.Len(); i++ {
		k, v := sm.PeekItem(i)
		b = append(b, fmt.Sprintf("%v:%v", k, v))
	}
	content := strings.Join(b, ", ")
	return fmt.Sprint("SortedDict{", content, "}")
}

type node struct {
	left, right int
	size        int
	priority    uint64
	key         Key
}

type sortedList struct {
	seed       uint64
	root       int
	comparator func(a, b Key) int
	nodes      []node
}

func newSortedList(comparator func(a, b Key) int, initCapacity int) *sortedList {
	sl := &sortedList{
		seed:       uint64(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, 0, max(initCapacity, 16)),
	}
	sl.nodes = append(sl.nodes, node{size: 0, priority: sl.nextRand()}) // dummy node 0
	return sl
}

func (sl *sortedList) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (sl *sortedList) Add(key Key) {
	var x, y, z int
	sl.splitByValue(sl.root, key, &x, &y, false)
	z = sl.newNode(key)
	sl.root = sl.merge(sl.merge(x, z), y)
}

func (sl *sortedList) At(index int) Key {
	n := sl.Len()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("%d index out of range: [%d,%d]", index, 0, n-1))
	}
	return sl.nodes[sl.kthNode(sl.root, index+1)].key
}

func (sl *sortedList) Pop(index int) Key {
	n := sl.Len()
	if index < 0 {
		index += n
	}

	index += 1 // dummy offset
	var x, y, z int
	sl.splitByRank(sl.root, index, &y, &z)
	sl.splitByRank(y, index-1, &x, &y)
	res := sl.nodes[y].key
	sl.root = sl.merge(x, z)
	return res
}

func (sl *sortedList) Discard(key Key) {
	var x, y, z int
	sl.splitByValue(sl.root, key, &x, &z, false)
	sl.splitByValue(x, key, &x, &y, true)
	y = sl.merge(sl.nodes[y].left, sl.nodes[y].right)
	sl.root = sl.merge(sl.merge(x, y), z)
}

// Remove [start, stop) from list.
func (sl *sortedList) Erase(start, stop int) {
	var x, y, z int
	start++ // dummy offset
	sl.splitByRank(sl.root, stop, &y, &z)
	sl.splitByRank(y, start-1, &x, &y)
	sl.root = sl.merge(x, z)
}

func (sl *sortedList) BisectLeft(key Key) int {
	var x, y int
	sl.splitByValue(sl.root, key, &x, &y, true)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *sortedList) BisectRight(key Key) int {
	var x, y int
	sl.splitByValue(sl.root, key, &x, &y, false)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *sortedList) String() string {
	sb := []string{"SortedList{"}
	values := []string{}
	for i := 0; i < sl.Len(); i++ {
		values = append(values, fmt.Sprintf("%v", sl.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")
}

func (sl *sortedList) Len() int {
	return sl.nodes[sl.root].size
}

func (sl *sortedList) kthNode(root int, k int) int {
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

func (sl *sortedList) splitByValue(root int, key Key, x, y *int, strictLess bool) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}
	if strictLess {
		if sl.comparator(sl.nodes[root].key, key) < 0 {
			*x = root
			sl.splitByValue(sl.nodes[root].right, key, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, key, x, &sl.nodes[root].left, strictLess)
		}
	} else {
		if sl.comparator(sl.nodes[root].key, key) <= 0 {
			*x = root
			sl.splitByValue(sl.nodes[root].right, key, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, key, x, &sl.nodes[root].left, strictLess)
		}
	}
	sl.pushUp(root)
}

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (sl *sortedList) splitByRank(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}
	if k <= sl.nodes[sl.nodes[root].left].size {
		*y = root
		sl.splitByRank(sl.nodes[root].left, k, x, &sl.nodes[root].left)
		sl.pushUp(*y)
	} else {
		*x = root
		sl.splitByRank(sl.nodes[root].right, k-sl.nodes[sl.nodes[root].left].size-1, &sl.nodes[root].right, y)
		sl.pushUp(*x)
	}
}

func (sl *sortedList) merge(x, y int) int {
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

// Return all elements in index order.
func (sl *sortedList) InOrder() []Key {
	res := make([]Key, 0, sl.Len())
	sl.inOrder(sl.root, &res)
	return res
}

func (sl *sortedList) inOrder(root int, res *[]Key) {
	if root == 0 {
		return
	}
	sl.inOrder(sl.nodes[root].left, res)
	*res = append(*res, sl.nodes[root].key)
	sl.inOrder(sl.nodes[root].right, res)
}

func (sl *sortedList) newNode(key Key) int {
	node := &node{
		key:      key,
		size:     1,
		priority: sl.nextRand(),
	}
	sl.nodes = append(sl.nodes, *node)
	return len(sl.nodes) - 1
}

// https://nyaannyaan.github.io/library/misc/rng.hpp
func (sl *sortedList) nextRand() uint64 {
	sl.seed ^= sl.seed << 7
	sl.seed ^= sl.seed >> 9
	return sl.seed
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
