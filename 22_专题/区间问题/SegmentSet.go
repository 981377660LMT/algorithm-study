// !注意:比较慢,有时可以用Intervals-珂朵莉树代替

package main

import (
	"fmt"
	"strings"
	"time"
)

const INF int = 1e18

type CountIntervals struct {
	seg *SegmentSet
}

func Constructor() CountIntervals {
	return CountIntervals{seg: NewSegmentSet(1e5)}
}

func (this *CountIntervals) Add(left int, right int) {
	this.seg.Insert(left, right)
}

func (this *CountIntervals) Count() int {
	return this.seg.Count
}

// 管理区间的数据结构.
//  1.所有区间都是闭区间 例如 [1,1] 表示 长为1的区间,起点为1;
//  2.有交集的区间会被合并,例如 [1,2]和[2,3]会被合并为[1,3].
type SegmentSet struct {
	Count int // 区间元素个数
	sl    *_SL
}

func NewSegmentSet(initCap int) *SegmentSet {
	return &SegmentSet{
		sl: _NSL(func(a, b Value) int {
			if a[0] != b[0] {
				return a[0] - b[0]
			}
			return a[1] - b[1]
		}, initCap),
	}
}

// 插入闭区间[left,right].
func (ss *SegmentSet) Insert(left, right int) {
	if left > right {
		return
	}
	it1 := ss.sl.BisectRight(Value{left, INF})
	it2 := ss.sl.BisectRight(Value{right, INF})
	if it1 != 0 && left <= ss.sl.At(it1 - 1)[1] {
		it1--
	}
	if it1 != it2 {
		tmp1 := ss.sl.At(it1)[0]
		if tmp1 < left {
			left = tmp1
		}
		tmp2 := ss.sl.At(it2 - 1)[1]
		if tmp2 > right {
			right = tmp2
		}
		removed := 0
		for i := it1; i < it2; i++ {
			cur := ss.sl.At(i)
			removed += cur[1] - cur[0] + 1
		}
		ss.sl.Erase(it1, it2)
		ss.Count -= removed
	}
	ss.sl.Add(Value{left, right})
	ss.Count += right - left + 1
}

func (ss *SegmentSet) Erase(left, right int) {
	if left > right {
		return
	}
	it1 := ss.sl.BisectLeft(Value{left, -INF})
	it2 := ss.sl.BisectRight(Value{right, INF})
	if it1 != 0 && left <= ss.sl.At(it1 - 1)[1] {
		it1--
	}
	if it1 == it2 {
		return
	}
	nl, nr := ss.sl.At(it1)[0], ss.sl.At(it2 - 1)[1]
	if left < nl {
		nl = left
	}
	if right > nr {
		nr = right
	}
	removed := 0
	for i := it1; i < it2; i++ {
		cur := ss.sl.At(i)
		removed += cur[1] - cur[0] + 1
	}
	ss.sl.Erase(it1, it2)
	ss.Count -= removed
	if nl < left {
		ss.sl.Add(Value{nl, left})
		ss.Count += left - nl + 1
	}
	if right < nr {
		ss.sl.Add(Value{right, nr})
		ss.Count += nr - right + 1
	}
}

// 返回第一个大于等于x的区间起点.
func (ss *SegmentSet) NextStart(x int) (res int, ok bool) {
	it := ss.sl.BisectLeft(Value{x, -INF})
	if it == ss.sl.Len() {
		return
	}
	res = ss.sl.At(it)[0]
	ok = true
	return
}

// 返回最后一个小于等于x的区间起点.
func (ss *SegmentSet) PrevStart(x int) (res int, ok bool) {
	it := ss.sl.BisectRight(Value{x, INF})
	if it == 0 {
		return
	}
	res = ss.sl.At(it - 1)[0]
	ok = true
	return
}

// 返回区间内第一个大于等于x的元素.
func (ss *SegmentSet) Ceiling(x int) (res int, ok bool) {
	it := ss.sl.BisectRight(Value{x, INF})
	if it != 0 && ss.sl.At(it - 1)[1] >= x {
		res = x
		ok = true
		return
	}
	if it != ss.sl.Len() {
		res = ss.sl.At(it)[0]
		ok = true
		return
	}
	return
}

// 返回区间内最后一个小于等于x的元素.
func (ss *SegmentSet) Floor(x int) (res int, ok bool) {
	it := ss.sl.BisectRight(Value{x, INF})
	if it == 0 {
		return
	}
	ok = true
	if ss.sl.At(it - 1)[1] >= x {
		res = x
		return
	}
	res = ss.sl.At(it - 1)[1]
	return
}

// 返回包含x的区间.
func (ss *SegmentSet) GetInterval(x int) (res [2]int, ok bool) {
	it := ss.sl.BisectRight(Value{x, INF})
	if it == 0 || ss.sl.At(it - 1)[1] < x {
		return
	}
	res = ss.sl.At(it - 1)
	ok = true
	return
}

func (ss *SegmentSet) Has(i int) bool {
	it := ss.sl.BisectRight(Value{i, INF})
	return it != 0 && ss.sl.At(it - 1)[1] >= i
}

func (ss *SegmentSet) HasInterval(left, right int) bool {
	if left > right {
		return false
	}
	it1 := ss.sl.BisectRight(Value{left, INF})
	if it1 == 0 {
		return false
	}
	it2 := ss.sl.BisectRight(Value{right, INF})
	if it1 != it2 {
		return false
	}
	return ss.sl.At(it1 - 1)[1] >= right
}

// 返回第index个区间.
func (ss *SegmentSet) At(index int) [2]int {
	return ss.sl.At(index)
}

func (ss *SegmentSet) GetAll() [][2]int {
	return ss.sl.InOrder()
}

func (ss *SegmentSet) String() string {
	res := []string{}
	all := ss.GetAll()
	for _, v := range all {
		res = append(res, fmt.Sprintf("[%d,%d]", v[0], v[1]))
	}
	return strings.Join(res, ",")
}

func (ss *SegmentSet) Len() int {
	return ss.sl.Len()
}

type Value = [2]int // [start,end]

type node struct {
	left, right int
	size        int
	priority    uint64
	value       Value
}

type _SL struct {
	seed       uint64
	root       int
	comparator func(a, b Value) int
	nodes      []node
}

func _NSL(comparator func(a, b Value) int, initCapacity int) *_SL {
	sl := &_SL{
		seed:       uint64(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, 0, max(initCapacity, 16)),
	}
	sl.nodes = append(sl.nodes, node{size: 0, priority: sl.nextRand()}) // dummy node 0
	return sl
}

func (sl *_SL) pushUp(root int) {
	sl.nodes[root].size = sl.nodes[sl.nodes[root].left].size + sl.nodes[sl.nodes[root].right].size + 1
}

func (sl *_SL) Add(value Value) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &y, false)
	z = sl.newNode(value)
	sl.root = sl.merge(sl.merge(x, z), y)
}

func (sl *_SL) At(index int) Value {
	return sl.nodes[sl.kthNode(sl.root, index+1)].value
}

func (sl *_SL) Pop(index int) Value {
	index += 1 // dummy offset
	var x, y, z int
	sl.splitByRank(sl.root, index, &y, &z)
	sl.splitByRank(y, index-1, &x, &y)
	res := sl.nodes[y].value
	sl.root = sl.merge(x, z)
	return res
}

func (sl *_SL) Discard(value Value) {
	var x, y, z int
	sl.splitByValue(sl.root, value, &x, &z, false)
	sl.splitByValue(x, value, &x, &y, true)
	y = sl.merge(sl.nodes[y].left, sl.nodes[y].right)
	sl.root = sl.merge(sl.merge(x, y), z)
}

// Remove [start, stop) from list.
func (sl *_SL) Erase(start, stop int) {
	var x, y, z int
	start++ // dummy offset
	sl.splitByRank(sl.root, stop, &y, &z)
	sl.splitByRank(y, start-1, &x, &y)
	sl.root = sl.merge(x, z)
}

func (sl *_SL) BisectLeft(value Value) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, true)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

func (sl *_SL) BisectRight(value Value) int {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, false)
	res := sl.nodes[x].size
	sl.root = sl.merge(x, y)
	return res
}

// 求小于等于 value 的最大值.
func (sl *_SL) Prev(value Value) (res Value, ok bool) {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, false)
	if x == 0 {
		ok = false
		return
	}
	res = sl.nodes[sl.kthNode(x, sl.nodes[x].size)].value
	sl.root = sl.merge(x, y)
	ok = true
	return
}

// 求大于等于 value 的最小值.
func (sl *_SL) Next(value Value) (res Value, ok bool) {
	var x, y int
	sl.splitByValue(sl.root, value, &x, &y, true)
	if y == 0 {
		ok = false
		return
	}
	res = sl.nodes[sl.kthNode(y, 1)].value
	sl.root = sl.merge(x, y)
	ok = true
	return
}

func (sl *_SL) String() string {
	sb := []string{"SortedList{"}
	values := []string{}
	for i := 0; i < sl.Len(); i++ {
		values = append(values, fmt.Sprintf("%v", sl.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")
}

func (sl *_SL) Len() int {
	return sl.nodes[sl.root].size
}

func (sl *_SL) kthNode(root int, k int) int {
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

func (sl *_SL) splitByValue(root int, value Value, x, y *int, strictLess bool) {
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

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (sl *_SL) splitByRank(root, k int, x, y *int) {
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

func (sl *_SL) merge(x, y int) int {
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
func (sl *_SL) InOrder() []Value {
	res := make([]Value, 0, sl.Len())
	sl.inOrder(sl.root, &res)
	return res
}

func (sl *_SL) inOrder(root int, res *[]Value) {
	if root == 0 {
		return
	}
	sl.inOrder(sl.nodes[root].left, res)
	*res = append(*res, sl.nodes[root].value)
	sl.inOrder(sl.nodes[root].right, res)
}

func (sl *_SL) newNode(value Value) int {
	sl.nodes = append(sl.nodes, node{
		value:    value,
		size:     1,
		priority: sl.nextRand(),
	})
	return len(sl.nodes) - 1
}

// https://nyaannyaan.github.io/library/misc/rng.hpp
func (sl *_SL) nextRand() uint64 {
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
