// https://qoj.ac/problem/382
// 二进制数加法/减法, 二进制数第k位, 二进制数的符号.

package main

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
	"unsafe"
)

func main() {
	b := NewBinaryNumber()
	b.Add(0, 100)
	b.Add(0, 2333)
	b.Add(0, -233)
	fmt.Println(b.Kth(5))
	fmt.Println(b.Kth(7))
	fmt.Println(b.Kth(15))
	b.Add(15, 5)
	fmt.Println(b.Kth(15))
	b.Add(12, -1)
	fmt.Println(b.Kth(15))
}

const INF int = 1e18

type BinaryNumber struct {
	data *SortedDict
}

func NewBinaryNumber() *BinaryNumber {
	return &BinaryNumber{data: NewSortedDict(func(a, b K) bool { return a < b })}
}

func (bn *BinaryNumber) Sgn() int32 {
	_, x := bn.prev(INF)
	return x
}

// 二进制第k位.
func (bn *BinaryNumber) Kth(k int) int32 {
	x := int32(0)
	if v, ok := bn.data.Get(k); ok {
		x = v
	}
	_, y := bn.prev(k - 1)
	if x == 0 {
		if y >= 0 {
			return 0
		}
		return 1
	}
	if y >= 0 {
		return 1
	}
	return 0
}

// 加上 2^k * x
func (bn *BinaryNumber) Add(k int, x int) {
	for x != 0 {
		if v, ok := bn.data.Get(k); ok {
			x += int(v)
		}
		if x&1 == 0 {
			bn.data.Delete(k)
		} else {
			bn.data.Set(k, int32(x&1))
		}
		k++
		x >>= 1
	}
}

func (bn *BinaryNumber) prev(k int) (key K, value V) {
	it := bn.data.UpperBound(k)
	if !it.HasPrev() {
		return -1, 0
	}
	it.Prev()
	key, value, _ = it.Entry()
	return
}

const _LOAD int = 100

type K = int
type V = int32
type Entry struct {
	key   K
	value V
}

type Iterator struct {
	sd    *SortedDict
	sIter *_SIterator
}

func (it *Iterator) Next() (key K, value V, ok bool) {
	next, ok := it.sIter.Next()
	if !ok {
		return
	}
	key = next
	value, ok = it.sd.dict[key]
	return
}

func (it *Iterator) Prev() (key K, value V, ok bool) {
	prev, ok := it.sIter.Prev()
	if !ok {
		return
	}
	key = prev
	value, ok = it.sd.dict[key]
	return
}

func (it *Iterator) HasNext() bool {
	return it.sIter.HasNext()
}

func (it *Iterator) HasPrev() bool {
	return it.sIter.HasPrev()
}

func (it *Iterator) Key() (key K, ok bool) {
	return it.sIter.Value()
}

func (it *Iterator) Value() (value V, ok bool) {
	key, ok := it.sIter.Value()
	if !ok {
		return
	}
	value, ok = it.sd.dict[key]
	return
}

func (it *Iterator) Entry() (key K, value V, ok bool) {
	key, ok = it.sIter.Value()
	if !ok {
		return
	}
	value, ok = it.sd.dict[key]
	return
}

func (it *Iterator) Remove() {
	key, ok := it.sIter.Value()
	if !ok {
		return
	}
	it.sIter.Remove()
	it.sd.Delete(key)
}

type SortedDict struct {
	sl   *SortedList
	dict map[K]V
}

func NewSortedDict(less func(a, b K) bool) *SortedDict {
	return &SortedDict{sl: NewSortedList(less), dict: map[K]V{}}
}

func (sd *SortedDict) Set(key K, value V) *SortedDict {
	if _, ok := sd.dict[key]; !ok {
		sd.sl.Add(key)
	}
	sd.dict[key] = value
	return sd
}

func (sd *SortedDict) SetDefault(key K, defaultValue V) V {
	if v, ok := sd.dict[key]; ok {
		return v
	}
	sd.sl.Add(key)
	sd.dict[key] = defaultValue
	return defaultValue
}

func (sd *SortedDict) Has(key K) bool {
	_, ok := sd.dict[key]
	return ok
}

func (sd *SortedDict) Get(key K) (value V, ok bool) {
	value, ok = sd.dict[key]
	return
}

func (sd *SortedDict) Delete(key K) bool {
	if _, ok := sd.dict[key]; !ok {
		return false
	}
	sd.sl.Discard(key)
	delete(sd.dict, key)
	return true
}

func (sd *SortedDict) Pop(key K, defaultValue V) V {
	if v, ok := sd.dict[key]; ok {
		sd.sl.Discard(key)
		delete(sd.dict, key)
		return v
	}
	return defaultValue
}

func (sd *SortedDict) PopItem(index int) (key K, value V, ok bool) {
	if len(sd.dict) == 0 {
		return
	}
	key = sd.sl.Pop(index)
	value, ok = sd.dict[key]
	delete(sd.dict, key)
	return
}

func (sd *SortedDict) PeekItem(index int) (key K, value V, ok bool) {
	if len(sd.dict) == 0 {
		return
	}
	key = sd.sl.At(index)
	value, ok = sd.dict[key]
	return
}

func (sd *SortedDict) PeekMinItem() (key K, value V, ok bool) {
	if len(sd.dict) == 0 {
		return
	}
	key = sd.sl.Min()
	value, ok = sd.dict[key]
	return
}

func (sd *SortedDict) PeekMaxItem() (key K, value V, ok bool) {
	if len(sd.dict) == 0 {
		return
	}
	key = sd.sl.Max()
	value, ok = sd.dict[key]
	return
}

func (sd *SortedDict) ForEach(callbackfn func(value V, key K)) {
	sd.sl.ForEach(func(key K, _ int) {
		callbackfn(sd.dict[key], key)
	}, false)
}

func (sd *SortedDict) Enumerate(start, end int, f func(key K, value V), erase bool) {
	sd.sl.Enumerate(start, end, func(key K) {
		f(key, sd.dict[key])
		if erase {
			delete(sd.dict, key)
		}
	}, erase)
}

func (sd *SortedDict) BisectLeft(key K) (index int) {
	return sd.sl.BisectLeft(key)
}

func (sd *SortedDict) BisectRight(key K) (index int) {
	return sd.sl.BisectRight(key)
}

func (sd *SortedDict) Floor(key K) (floorKey K, floorValue V, ok bool) {
	floorKey, ok = sd.sl.Floor(key)
	if !ok {
		return
	}
	floorValue, ok = sd.dict[floorKey]
	return
}

func (sd *SortedDict) Ceiling(key K) (ceilKey K, ceilValue V, ok bool) {
	ceilKey, ok = sd.sl.Ceiling(key)
	if !ok {
		return
	}
	ceilValue, ok = sd.dict[ceilKey]
	return
}

func (sd *SortedDict) Lower(key K) (lowerKey K, lowerValue V, ok bool) {
	lowerKey, ok = sd.sl.Lower(key)
	if !ok {
		return
	}
	lowerValue, ok = sd.dict[lowerKey]
	return
}

func (sd *SortedDict) Higher(key K) (higherKey K, higherValue V, ok bool) {
	higherKey, ok = sd.sl.Higher(key)
	if !ok {
		return
	}
	higherValue, ok = sd.dict[higherKey]
	return
}

func (sd *SortedDict) Size() int {
	return len(sd.dict)
}

func (sd *SortedDict) Keys() (res []K) {
	res = make([]K, 0, len(sd.dict))
	sd.sl.ForEach(func(key K, _ int) {
		res = append(res, key)
	}, false)
	return
}

func (sd *SortedDict) Values() (res []V) {
	res = make([]V, 0, len(sd.dict))
	sd.sl.ForEach(func(_ K, value int) {
		res = append(res, sd.dict[value])
	}, false)
	return
}

func (sd *SortedDict) Entries() (res []Entry) {
	res = make([]Entry, 0, len(sd.dict))
	sd.sl.ForEach(func(key K, value int) {
		res = append(res, Entry{key, sd.dict[value]})
	}, false)
	return
}

func (sd *SortedDict) Clear() {
	sd.sl.Clear()
	sd.dict = map[K]V{}
}

func (sd *SortedDict) String() string {
	sb := []string{}
	sb = append(sb, fmt.Sprintf("SortedDict(%d) {", sd.Size()))
	sd.ForEach(func(value V, key K) {
		sb = append(sb, fmt.Sprintf("  %v => %v,", key, value))
	})
	sb = append(sb, "}")
	return strings.Join(sb, "\n")
}

func (sd *SortedDict) IteratorAt(index int) *Iterator {
	return &Iterator{sd: sd, sIter: sd.sl.IteratorAt(index)}
}

func (sd *SortedDict) LowerBound(key K) *Iterator {
	return &Iterator{sd: sd, sIter: sd.sl.LowerBound(key)}
}

func (sd *SortedDict) UpperBound(key K) *Iterator {
	return &Iterator{sd: sd, sIter: sd.sl.UpperBound(key)}
}

// 使用分块+树状数组维护的有序序列.
type SortedList struct {
	less              func(a, b K) bool
	size              int
	blocks            [][]K
	mins              []K
	tree              []int
	shouldRebuildTree bool
}

func NewSortedList(less func(a, b K) bool, elements ...K) *SortedList {
	elements = append(elements[:0:0], elements...)
	res := &SortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]K{}
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]K, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	return res
}

func (sl *SortedList) Add(value K) *SortedList {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []K{value})
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
		left := append([]K(nil), sl.blocks[pos][:_LOAD]...)
		right := append([]K(nil), sl.blocks[pos][_LOAD:]...)
		sl.blocks = Replace(sl.blocks, pos, pos+1, left, right)
		sl.mins = Insert(sl.mins, pos+1, right[0])
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SortedList) Has(value K) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < len(sl.blocks[pos]) && sl.blocks[pos][index] == value
}

func (sl *SortedList) Discard(value K) bool {
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

func (sl *SortedList) Pop(index int) K {
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

func (sl *SortedList) At(index int) K {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SortedList) Erase(start, end int) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedList) Lower(value K) (res K, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList) Higher(value K) (res K, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedList) Floor(value K) (res K, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList) Ceiling(value K) (res K, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedList) BisectLeft(value K) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedList) BisectRight(value K) int {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedList) Count(value K) int {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SortedList) ForEach(f func(value K, index int), reverse bool) {
	if !reverse {
		count := 0
		for i := 0; i < len(sl.blocks); i++ {
			block := sl.blocks[i]
			for j := 0; j < len(block); j++ {
				f(block[j], count)
				count++
			}
		}
		return
	}

	count := 0
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		for j := len(block) - 1; j >= 0; j-- {
			f(block[j], count)
			count++
		}
	}
}

func (sl *SortedList) Enumerate(start, end int, f func(value K), erase bool) {
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

func (sl *SortedList) IteratorAt(index int) *_SIterator {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("Index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl._iteratorAt(pos, startIndex)
}

func (sl *SortedList) LowerBound(value K) *_SIterator {
	pos, index := sl._locLeft(value)
	return sl._iteratorAt(pos, index)
}

func (sl *SortedList) UpperBound(value K) *_SIterator {
	pos, index := sl._locRight(value)
	return sl._iteratorAt(pos, index)
}

func (sl *SortedList) Min() K {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *SortedList) Max() K {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedList) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedList{")
	sl.ForEach(func(value K, index int) {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *SortedList) Len() int {
	return sl.size
}

func (sl *SortedList) _delete(pos, index int) {
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

func (sl *SortedList) _locLeft(value K) (pos, index int) {
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

func (sl *SortedList) _locRight(value K) (pos, index int) {
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

func (sl *SortedList) _locBlock(value K) int {
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

func (sl *SortedList) _buildTree() {
	sl.tree = make([]int, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = len(sl.blocks[i])
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

func (sl *SortedList) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SortedList) _queryTree(end int) int {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := 0
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *SortedList) _findKth(k int) (pos, index int) {
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
	bitLength := bits.Len32(uint32(len(tree)))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < len(tree) && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func (sl *SortedList) _iteratorAt(pos, index int) *_SIterator {
	return &_SIterator{sl: sl, pos: pos, index: index}
}

type _SIterator struct {
	sl    *SortedList
	pos   int
	index int
}

func (it *_SIterator) HasNext() bool {
	return it.pos < len(it.sl.blocks)-1 || it.index < len(it.sl.blocks[it.pos])-1
}

func (it *_SIterator) Next() (res K, ok bool) {
	if !it.HasNext() {
		return
	}
	it.index++
	if it.index == len(it.sl.blocks[it.pos]) {
		it.pos++
		it.index = 0
	}
	res = it.sl.blocks[it.pos][it.index]
	ok = true
	return
}

func (it *_SIterator) HasPrev() bool {
	return it.pos > 0 || it.index > 0
}

func (it *_SIterator) Prev() (res K, ok bool) {
	if !it.HasPrev() {
		return
	}
	it.index--
	if it.index == -1 {
		it.pos--
		it.index = len(it.sl.blocks[it.pos]) - 1
	}
	res = it.sl.blocks[it.pos][it.index]
	ok = true
	return
}

func (it *_SIterator) Remove() {
	it.sl._delete(it.pos, it.index)
}

func (it *_SIterator) Value() (res K, ok bool) {
	if it.pos < 0 || it.pos >= it.sl.Len() {
		return
	}
	block := it.sl.blocks[it.pos]
	if it.index < 0 || it.index >= len(block) {
		return
	}
	res = block[it.index]
	ok = true
	return
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
