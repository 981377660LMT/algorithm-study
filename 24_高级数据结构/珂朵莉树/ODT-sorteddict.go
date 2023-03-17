// 珂朵莉树(ODT)/Intervals
// 珂朵莉树(ODT)/Intervals

package main

import (
	"fmt"
	"strings"
	"time"
)

const INF int = 1e18

func demo() {
	sd := nsd(10)
	sd.Set(1, 1)
	sd.Set(2, 2)
	fmt.Println(sd)
	sd.Discard(1)
	fmt.Println(sd)
	sd.Set(3, 3)
	fmt.Println(sd)
	sd.Set(4, 4)
	fmt.Println(sd)
	fmt.Println(sd.PeekItem(0))
	sd.Erase(1)
}

func main() {
	demo()
}

type Value = int // !不能是指针

type Intervals struct {
	Count     int // 区间数
	Len       int // 区间长度和
	noneValue Value
	mp        *sd
}

func NewIntervals(noneValue Value) *Intervals {
	res := &Intervals{
		noneValue: noneValue,
		mp:        nsd(16),
	}
	res.mp.Set(-INF, noneValue)
	res.mp.Set(INF, noneValue)
	return res
}

// 返回包含 x 的区间的信息.
func (odt *Intervals) Get(x int, erase bool) (start, end int, value Value) {
	iter2, key2, _ := odt.mp.UpperBound(x)
	iter1 := iter2 - 1
	key1, val1 := odt.mp.PeekItem(iter1)
	start, end, value = key1, key2, val1
	if erase && value != odt.noneValue {
		odt.Count--
		odt.Len -= end - start
		odt.mp.Set(start, odt.noneValue)
		odt.mergeAt(start)
		odt.mergeAt(end)
	}
	return
}

func (odt *Intervals) Set(start, end int, value Value) {
	odt.EnumerateRange(start, end, func(l, r int, x Value) {}, true)
	odt.mp.Set(start, value)
	if value != odt.noneValue {
		odt.Count++
		odt.Len += end - start
	}
	odt.mergeAt(start)
	odt.mergeAt(end)
}

func (odt *Intervals) EnumerateAll(f func(start, end int, value Value)) {
	odt.EnumerateRange(-INF, INF, f, false)
}

// 遍历范围 [L, R) 内的所有数据.
func (odt *Intervals) EnumerateRange(start, end int, f func(start, end int, value Value), erase bool) {
	if !(-INF <= start && start <= end && end <= INF) {
		panic(fmt.Sprintf("invalid range [%d, %d)", start, end))
	}

	NONE := odt.noneValue

	if !erase {
		it1, _, _ := odt.mp.UpperBound(start)
		it1--
		for {
			key1, val1 := odt.mp.PeekItem(it1)
			if key1 >= end {
				break
			}
			it2 := it1 + 1
			key2, _ := odt.mp.PeekItem(it2)
			f(max(key1, start), min(key2, end), val1)
			it1 = it2
		}
		return
	}

	iter1, _, _ := odt.mp.UpperBound(start)
	iter1--
	if key1, val1 := odt.mp.PeekItem(iter1); key1 < start {
		odt.mp.Set(start, val1)
		if val1 != NONE {
			odt.Count++
		}
	}

	// 分割区间
	iter1, _, _ = odt.mp.LowerBound(end)
	if key1, _ := odt.mp.PeekItem(iter1); key1 > end {
		_, val2 := odt.mp.PeekItem(iter1 - 1)
		odt.mp.Set(end, val2)
		if val2 != NONE {
			odt.Count++
		}
	}

	iter1, _, _ = odt.mp.LowerBound(start)
	for {
		key1, val1 := odt.mp.PeekItem(iter1)
		if key1 >= end {
			break
		}
		iter2 := iter1 + 1
		key2, _ := odt.mp.PeekItem(iter2)
		f(key1, key2, val1)
		if val1 != NONE {
			odt.Count--
			odt.Len -= key2 - key1
		}
		odt.mp.Erase(iter1)
	}

	odt.mp.Set(start, NONE)
}

func (odt *Intervals) String() string {
	sb := []string{}
	odt.EnumerateAll(func(start, end int, value Value) {
		var v interface{} = value
		if value == odt.noneValue {
			v = "nil"
		}
		sb = append(sb, fmt.Sprintf("[%d,%d):%v", start, end, v))
	})
	return fmt.Sprintf("ODT{%v}", strings.Join(sb, ", "))
}

func (odt *Intervals) mergeAt(p int) {
	if p == -INF || p == INF {
		return
	}
	iter1, _, val1 := odt.mp.LowerBound(p)
	iter2 := iter1 - 1
	_, val2 := odt.mp.PeekItem(iter2)
	if val1 == val2 {
		if val1 != odt.noneValue {
			odt.Count--
		}
		odt.mp.Erase(iter1)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type sd struct {
	sl *stl
	mp map[int]Value
}

// key 为 int, value 为 Item 的有序映射.
func nsd(initCapacity int) *sd {
	return &sd{
		sl: nstl(initCapacity),
		mp: make(map[int]Value, initCapacity),
	}
}

func (sm *sd) Set(key int, value Value) {
	if _, has := sm.mp[key]; has {
		sm.mp[key] = value
		return
	}
	sm.mp[key] = value
	sm.sl.Add(key)
}

func (sm *sd) Get(key int) Value {
	return sm.mp[key]
}

func (sm *sd) Discard(key int) {
	if _, has := sm.mp[key]; has {
		delete(sm.mp, key)
		sm.sl.Discard(key)
	}
}

func (sm *sd) Erase(iter int) {
	key := sm.sl.At(iter)
	sm.sl.Pop(iter)
	delete(sm.mp, key)
}

func (sm *sd) PeekItem(iter int) (key int, value Value) {
	key = sm.sl.At(iter)
	value = sm.mp[key]
	return
}

func (sm *sd) LowerBound(key_ int) (iter int, key int, value Value) {
	iter = sm.sl.BisectLeft(key_)
	key = sm.sl.At(iter)
	value = sm.mp[key]
	return
}

func (sm *sd) UpperBound(key_ int) (iter int, key int, value Value) {
	iter = sm.sl.BisectRight(key_)
	key = sm.sl.At(iter)
	value = sm.mp[key]
	return
}

func (sm *sd) String() string {
	var b []string
	for i := 0; i < sm.sl.Len(); i++ {
		k, v := sm.PeekItem(i)
		b = append(b, fmt.Sprintf("%v:%v", k, v))
	}
	content := strings.Join(b, ", ")
	return fmt.Sprint("sortedDict{", content, "}")
}

type node struct {
	left, right int
	size        int
	priority    uint64
	value       int
}

type stl struct {
	seed  uint64
	root  int
	nodes []node
}

// value 为 int 的有序列表.
func nstl(initCapacity int) *stl {
	sl := &stl{
		seed:  uint64(time.Now().UnixNano()/2 + 1),
		nodes: make([]node, 0, max(initCapacity, 16)),
	}
	sl.nodes = append(sl.nodes, node{size: 0, priority: sl.nextRand()}) // dummy node 0
	return sl
}

func (sl *stl) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (sl *stl) Add(value int) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &y, false)
	z = sl.newNode(value)
	sl.root = sl.merge(sl.merge(x, z), y)
}

func (sl *stl) At(index int) int {
	return sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (sl *stl) Pop(index int) int {
	index += 1 // dummy offset
	var x, y, z int
	sl.splitByRank(sl.root, index, &y, &z)
	sl.splitByRank(y, index-1, &x, &y)
	res := sl.nodes[y].value
	sl.root = sl.merge(x, z)
	return res
}

func (sl *stl) Discard(value int) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &z, false)
	sl.splitByValue(x, value, &x, &y, true)
	y = sl.merge(sl.nodes[y].left, sl.nodes[y].right)
	sl.root = sl.merge(sl.merge(x, y), z)
}

func (sl *stl) BisectLeft(value int) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, true)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *stl) BisectRight(value int) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, false)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *stl) String() string {
	sb := []string{"SortedList{"}
	values := []string{}
	for i := 0; i < sl.Len(); i++ {
		values = append(values, fmt.Sprintf("%v", sl.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")
}

func (sl *stl) Len() int {
	return sl.nodes[sl.root].size
}

func (sl *stl) kthNode(root int, k int) int {
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

func (sl *stl) splitByValue(root int, value int, x, y *int, strictLess bool) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}
	if strictLess {
		if sl.nodes[root].value < value {
			*x = root
			sl.splitByValue(sl.nodes[root].right, value, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, value, x, &sl.nodes[root].left, strictLess)
		}
	} else {
		if sl.nodes[root].value <= value {
			*x = root
			sl.splitByValue(sl.nodes[root].right, value, &sl.nodes[root].right, y, strictLess)
		} else {
			*y = root
			sl.splitByValue(sl.nodes[root].left, value, x, &sl.nodes[root].left, strictLess)
		}
	}
	sl.pushUp(root)
}

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (sl *stl) splitByRank(root, k int, x, y *int) {
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

func (sl *stl) merge(x, y int) int {
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

func (sl *stl) newNode(value int) int {
	node := &node{
		value:    value,
		size:     1,
		priority: sl.nextRand(),
	}
	sl.nodes = append(sl.nodes, *node)
	return len(sl.nodes) - 1
}

// https://nyaannyaan.github.io/library/misc/rng.hpp
func (sl *stl) nextRand() uint64 {
	sl.seed ^= sl.seed << 7
	sl.seed ^= sl.seed >> 9
	return sl.seed
}
