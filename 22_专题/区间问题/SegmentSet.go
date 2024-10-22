package main

import (
	"fmt"
	"math/bits"
	"unsafe"

	"sort"
	"strings"
)

const INF int = 1e18

// https://leetcode.cn/problems/range-module/submissions/
// Range模块/Range 模块
type RangeModule struct {
	ss *SegmentSet
}

func Constructor() RangeModule {
	return RangeModule{ss: NewSegmentSet()}
}

func (this *RangeModule) AddRange(left int, right int) {
	this.ss.Insert(left, right)
}

func (this *RangeModule) QueryRange(left int, right int) bool {
	return this.ss.HasInterval(left, right)
}

func (this *RangeModule) RemoveRange(left int, right int) {
	this.ss.Erase(left, right)
}

// 管理区间的数据结构.
//
//	1.所有区间都是闭区间 例如 [1,1] 表示 长为1的区间,起点为1;
//	2.有交集的区间会被合并,例如 [1,2]和[2,3]会被合并为[1,3].
type SegmentSet struct {
	sl        *sortedList
	count     int // 区间元素个数
	flyWeight S
}

func NewSegmentSet() *SegmentSet {
	return &SegmentSet{sl: NewSortedList(func(a, b S) bool { return a.left < b.left }), flyWeight: Interval{}}
}

// 插入闭区间[left,right].
func (ss *SegmentSet) Insert(left, right int) {
	if left > right {
		return
	}
	it1 := ss.sl.BisectRight(ss.getFlyWeight(left, INF))
	it2 := ss.sl.BisectRight(ss.getFlyWeight(right, INF))
	if it1 != 0 && left <= ss.sl.At(it1-1).right {
		it1--
	}
	if it1 != it2 {
		tmp1 := ss.sl.At(it1).left
		if tmp1 < left {
			left = tmp1
		}
		tmp2 := ss.sl.At(it2 - 1).right
		if tmp2 > right {
			right = tmp2
		}
		removed := 0
		ss.sl.Enumerate(it1, it2, func(value S) {
			removed += value.right - value.left + 1
		}, true)
		ss.count -= removed
	}
	ss.sl.Add(Interval{left, right})
	ss.count += right - left + 1
}

func (ss *SegmentSet) Erase(left, right int) {
	if left > right {
		return
	}
	it1 := ss.sl.BisectRight(ss.getFlyWeight(left, INF))
	it2 := ss.sl.BisectRight(ss.getFlyWeight(right, INF))
	if it1 != 0 && left <= ss.sl.At(it1-1).right {
		it1--
	}
	if it1 == it2 {
		return
	}
	nl, nr := ss.sl.At(it1).left, ss.sl.At(it2-1).right
	if left < nl {
		nl = left
	}
	if right > nr {
		nr = right
	}

	removed := 0
	ss.sl.Enumerate(it1, it2, func(value S) {
		removed += value.right - value.left + 1
	}, true)
	ss.count -= removed
	if nl < left {
		ss.sl.Add(Interval{nl, left})
		ss.count += left - nl + 1
	}
	if right < nr {
		ss.sl.Add(Interval{right, nr})
		ss.count += nr - right + 1
	}
}

// 返回第一个大于等于x的区间起点.
func (ss *SegmentSet) NextStart(x int) (res int, ok bool) {
	it := ss.sl.BisectLeft(ss.getFlyWeight(x, -INF))
	if it == ss.sl.Len() {
		return
	}
	res = ss.sl.At(it).left
	ok = true
	return
}

// 返回最后一个小于等于x的区间起点.
func (ss *SegmentSet) PrevStart(x int) (res int, ok bool) {
	it := ss.sl.BisectRight(ss.getFlyWeight(x, INF)) - 1
	if it == 0 {
		return
	}
	res = ss.sl.At(it - 1).left
	ok = true
	return
}

// 返回区间内第一个大于等于x的元素.
func (ss *SegmentSet) Ceiling(x int) (res int, ok bool) {
	it := ss.sl.BisectRight(ss.getFlyWeight(x, INF))
	if it != 0 && ss.sl.At(it-1).right >= x {
		res = x
		ok = true
		return
	}
	if it != ss.sl.Len() {
		res = ss.sl.At(it).left
		ok = true
		return
	}
	return
}

// 返回区间内最后一个小于等于x的元素.
func (ss *SegmentSet) Floor(x int) (res int, ok bool) {
	it := ss.sl.BisectRight(ss.getFlyWeight(x, INF))
	if it == 0 {
		return
	}
	ok = true
	if ss.sl.At(it-1).right >= x {
		res = x
		return
	}
	res = ss.sl.At(it - 1).right
	return
}

// 返回包含x的区间.
func (ss *SegmentSet) GetInterval(x int) (res S, ok bool) {
	it := ss.sl.BisectRight(ss.getFlyWeight(x, INF))
	if it == 0 || ss.sl.At(it-1).right < x {
		return
	}
	res = ss.sl.At(it - 1)
	ok = true
	return
}

func (ss *SegmentSet) Has(i int) bool {
	it := ss.sl.BisectRight(ss.getFlyWeight(i, INF))
	return it != 0 && ss.sl.At(it-1).right >= i
}

func (ss *SegmentSet) HasInterval(left, right int) bool {
	if left > right {
		return false
	}
	it1 := ss.sl.BisectRight(ss.getFlyWeight(left, INF))
	if it1 == 0 {
		return false
	}
	it2 := ss.sl.BisectRight(ss.getFlyWeight(right, INF))
	if it1 != it2 {
		return false
	}
	return ss.sl.At(it1-1).right >= right
}

// 返回第index个区间.
func (ss *SegmentSet) At(index int) S {
	return ss.sl.At(index)
}

func (ss *SegmentSet) GetAll() []S {
	res := make([]S, 0, ss.sl.Len())
	ss.sl.ForEach(func(value S, _ int) bool {
		res = append(res, value)
		return false
	}, false)
	return res
}

// 遍历 [L,R] 内的所有区间范围.
func (ss *SegmentSet) EnumerateRange(L, R int, f func(left, right int)) {
	if L > R {
		return
	}
	it := ss.sl.BisectRight(ss.getFlyWeight(L, INF)) - 1
	if it < 0 {
		it++
	}
	ss.sl.EnumerateSlice(it, ss.sl.Len(), func(interval S) bool {
		if interval.left > R {
			return true
		}
		f(max(interval.left, L), min(interval.right, R))
		return false
	})
}

func (ss *SegmentSet) String() string {
	res := []string{}
	all := ss.GetAll()
	for _, v := range all {
		res = append(res, fmt.Sprintf("[%d,%d]", v.left, v.right))
	}
	return strings.Join(res, ",")
}

func (ss *SegmentSet) Len() int {
	return ss.sl.Len()
}

func (ss *SegmentSet) Count() int {
	return ss.count
}

func (ss *SegmentSet) Clear() {
	ss.sl.Clear()
	ss.count = 0
}

func (ss *SegmentSet) getFlyWeight(left, right int) S {
	ss.flyWeight.left = left
	ss.flyWeight.right = right
	return ss.flyWeight
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 100

type Interval = struct{ left, right int }
type S = Interval

// 使用分块+树状数组维护的有序序列.
type sortedList struct {
	less              func(a, b S) bool
	size              int
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func NewSortedList(less func(a, b S) bool, elements ...S) *sortedList {
	elements = append(elements[:0:0], elements...)
	res := &sortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]S{}
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]S, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	return res
}

func (sl *sortedList) Add(value S) *sortedList {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return sl
	}

	pos, index := sl._locRight(value)

	sl._updateTree(pos, 1)
	sl.blocks[pos] = Insert(sl.blocks[pos], index, value)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		left := append([]S(nil), sl.blocks[pos][:_LOAD]...)
		right := append([]S(nil), sl.blocks[pos][_LOAD:]...)
		sl.blocks = Replace(sl.blocks, pos, pos+1, left, right)
		sl.mins = Insert(sl.mins, pos+1, right[0])
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *sortedList) Has(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < len(sl.blocks[pos]) && sl.blocks[pos][index] == value
}

func (sl *sortedList) Discard(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locRight(value)
	if index > 0 && sl.blocks[pos][index-1] == value {
		sl._delete(pos, index-1)
		return true
	}
	return false
}

func (sl *sortedList) Pop(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}

func (sl *sortedList) At(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *sortedList) Erase(start, end int) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *sortedList) Lower(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *sortedList) Higher(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *sortedList) Floor(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *sortedList) Ceiling(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *sortedList) BisectLeft(value S) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *sortedList) BisectRight(value S) int {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *sortedList) Count(value S) int {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *sortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}
func (sl *sortedList) ForEach(f func(value S, index int) bool, reverse bool) {
	if !reverse {
		count := 0
		for i := 0; i < len(sl.blocks); i++ {
			block := sl.blocks[i]
			for j := 0; j < len(block); j++ {
				if f(block[j], count) {
					return
				}
				count++
			}
		}
		return
	}
	count := 0
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		for j := len(block) - 1; j >= 0; j-- {
			if f(block[j], count) {
				return
			}
			count++
		}
	}
}
func (sl *sortedList) Enumerate(start, end int, f func(value S), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return
	}

	pos, startIndex := sl._findKth(start)
	count := end - start
	m := len(sl.blocks)
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		endIndex := min(len(block), startIndex+count)
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == len(block) {
				// !delete block
				sl.blocks = Replace(sl.blocks, pos, pos+1)
				sl.mins = Replace(sl.mins, pos, pos+1)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				sl._updateTree(pos, -deleted)
				sl.blocks[pos] = Replace(sl.blocks[pos], startIndex, endIndex)
				sl.mins[pos] = sl.blocks[pos][0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *sortedList) Min() S {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *sortedList) Max() S {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *sortedList) EnumerateSlice(start, end int, f func(value S) bool) {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return
	}
	count := end - start
	pos, index := sl._findKth(start)
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endPos := min(len(block), index+count)
		for j := index; j < endPos; j++ {
			if f(block[j]) {
				return
			}
		}
		count -= endPos - index
		index = 0
	}
}

func (sl *sortedList) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedList{")
	sl.ForEach(func(value S, index int) bool {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
		return false
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *sortedList) Len() int {
	return sl.size
}

func (sl *sortedList) _delete(pos, index int) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = Replace(sl.blocks[pos], index, index+1)
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	sl.blocks = Replace(sl.blocks, pos, pos+1)
	sl.mins = Replace(sl.mins, pos, pos+1)
	sl.shouldRebuildTree = true
}

func (sl *sortedList) _locLeft(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := -1
	right := len(sl.blocks) - 1
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	pos = right

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = len(cur)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(cur[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *sortedList) _locRight(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := 0
	right := len(sl.blocks)
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, sl.mins[mid]) {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = len(cur)
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, cur[mid]) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *sortedList) _locBlock(value S) int {
	left, right := -1, len(sl.blocks)-1
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	return right
}

func (sl *sortedList) _buildTree() {
	sl.tree = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = int32(len(sl.blocks[i]))
	}
	tree := sl.tree
	for i := 0; i < len(tree); i++ {
		j := i | (i + 1)
		if j < len(tree) {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *sortedList) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	d32 := int32(delta)
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += d32
	}
}

func (sl *sortedList) _queryTree(end int) int {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := int32(0)
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return int(sum)
}

func (sl *sortedList) _findKth(k int) (pos, index int) {
	if k < len(sl.blocks[0]) {
		return 0, k
	}
	last := len(sl.blocks) - 1
	lastLen := len(sl.blocks[last])
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	pos = -1
	k32 := int32(k)
	bitLength := bits.Len32(uint32(len(tree)))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < len(tree) && k32 >= tree[next] {
			pos = next
			k32 -= tree[pos]
		}
	}
	return pos + 1, int(k32)
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

// Replace replaces the elements s[i:j] by the given v, and returns the modified slice.
// !Like JavaScirpt's Array.prototype.splice.
func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	if i < 0 {
		i = 0
	}
	if j > len(s) {
		j = len(s)
	}
	if i == j {
		return Insert(s, i, v...)
	}
	if j == len(s) {
		return append(s[:i], v...)
	}
	tot := len(s[:i]) + len(v) + len(s[j:])
	if tot > cap(s) {
		s2 := append(s[:i], make(S, tot-i)...)
		copy(s2[i:], v)
		copy(s2[i+len(v):], s[j:])
		return s2
	}
	r := s[:tot]
	if i+len(v) <= j {
		copy(r[i:], v)
		copy(r[i+len(v):], s[j:])
		clear(s[tot:])
		return r
	}
	if !overlaps(r[i+len(v):], v) {
		copy(r[i+len(v):], s[j:])
		copy(r[i:], v)
		return r
	}
	y := len(v) - (j - i)
	if !overlaps(r[i:j], v) {
		copy(r[i:j], v[y:])
		copy(r[len(s):], v[:y])
		rotateRight(r[i:], y)
		return r
	}
	if !overlaps(r[len(s):], v) {
		copy(r[len(s):], v[:y])
		copy(r[i:j], v[y:])
		rotateRight(r[i:], y)
		return r
	}
	k := startIdx(v, s[j:])
	copy(r[i:], v)
	copy(r[i+len(v):], r[i+k:])
	return r
}

func rotateLeft[E any](s []E, r int) {
	for r != 0 && r != len(s) {
		if r*2 <= len(s) {
			swap(s[:r], s[len(s)-r:])
			s = s[:len(s)-r]
		} else {
			swap(s[:len(s)-r], s[r:])
			s, r = s[len(s)-r:], r*2-len(s)
		}
	}
}

func rotateRight[E any](s []E, r int) {
	rotateLeft(s, len(s)-r)
}

func swap[E any](x, y []E) {
	for i := 0; i < len(x); i++ {
		x[i], y[i] = y[i], x[i]
	}
}

func overlaps[E any](a, b []E) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	elemSize := unsafe.Sizeof(a[0])
	if elemSize == 0 {
		return false
	}
	return uintptr(unsafe.Pointer(&a[0])) <= uintptr(unsafe.Pointer(&b[len(b)-1]))+(elemSize-1) &&
		uintptr(unsafe.Pointer(&b[0])) <= uintptr(unsafe.Pointer(&a[len(a)-1]))+(elemSize-1)
}

func startIdx[E any](haystack, needle []E) int {
	p := &needle[0]
	for i := range haystack {
		if p == &haystack[i] {
			return i
		}
	}
	panic("needle not found")
}

func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	if i < 0 {
		i = 0
	}
	if i > len(s) {
		i = len(s)
	}

	m := len(v)
	if m == 0 {
		return s
	}
	n := len(s)
	if i == n {
		return append(s, v...)
	}
	if n+m > cap(s) {
		s2 := append(s[:i], make(S, n+m-i)...)
		copy(s2[i:], v)
		copy(s2[i+m:], s[i:])
		return s2
	}
	s = s[:n+m]
	if !overlaps(v, s[i+m:]) {
		copy(s[i+m:], s[i:])
		copy(s[i:], v)
		return s
	}
	copy(s[n:], v)
	rotateRight(s[i:], m)
	return s
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
