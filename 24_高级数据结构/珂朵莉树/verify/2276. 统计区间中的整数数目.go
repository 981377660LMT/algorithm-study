package main

import (
	"fmt"
	"math/bits"
	"runtime/debug"
	"sort"
	"strings"
)

func init() {
	debug.SetGCPercent(-1)
}

type CountIntervals struct {
	odt *Intervals
}

func Constructor() CountIntervals {
	return CountIntervals{NewIntervals(-INF)}
}

func (this *CountIntervals) Add(left int, right int) {
	this.odt.Set(left, right+1, 1)
}

func (this *CountIntervals) Count() int {
	return this.odt.Count
}

/**
 * Your CountIntervals object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Add(left,right);
 * param_2 := obj.Count();
 */

const INF int = 1e18

type Value = int

type Intervals struct {
	Len       int // 区间数
	Count     int // 区间元素个数之和
	noneValue Value
	mp        *SortedDict
}

func NewIntervals(noneValue Value) *Intervals {
	res := &Intervals{
		noneValue: noneValue,
		mp:        NewSortedDict(func(a, b int) bool { return a < b }),
	}
	res.mp.Set(-INF, noneValue)
	res.mp.Set(INF, noneValue)
	return res
}

// 返回包含 x 的区间的信息.
func (odt *Intervals) Get(x int, erase bool) (start, end int, value Value) {
	iter := odt.mp.UpperBound(x)
	end, _ = iter.Key()
	iter.Prev()
	start, value, _ = iter.Entry()
	if erase && value != odt.noneValue {
		odt.Len--
		odt.Count -= end - start
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
		odt.Len++
		odt.Count += end - start
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
		it1 := odt.mp.UpperBound(start)
		it1.Prev()
		for {
			key1, val1, _ := it1.Entry()
			if key1 >= end {
				break
			}
			it1.Next()
			key2, _ := it1.Key()
			f(max(key1, start), min(key2, end), val1)
		}
		return
	}

	iter1 := odt.mp.UpperBound(start)
	iter1.Prev()
	if key1, val1, _ := iter1.Entry(); key1 < start {
		odt.mp.Set(start, val1)
		if val1 != NONE {
			odt.Len++
		}
	}

	// 分割区间
	iter1 = odt.mp.LowerBound(end)
	if key1, _ := iter1.Key(); key1 > end {
		iter1.Prev()
		val2, _ := iter1.Value()
		odt.mp.Set(end, val2)
		if val2 != NONE {
			odt.Len++
		}
	}

	iter1 = odt.mp.LowerBound(start)
	for {
		key1, val1, _ := iter1.Entry()
		if key1 >= end {
			break
		}
		iter1.Remove()
		key2, _ := iter1.Key()
		f(key1, key2, val1)
		if val1 != NONE {
			odt.Len--
			odt.Count -= key2 - key1
		}
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
	iter1 := odt.mp.LowerBound(p)
	val1, _ := iter1.Value()
	iter1.Prev()
	val2, _ := iter1.Value()
	if val1 == val2 {
		if val1 != odt.noneValue {
			odt.Len--
		}
		iter1.Next()
		iter1.Remove()
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

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 200

type K = int
type V = int
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

func (sd *SortedDict) Slice(start, end int) (res []Entry) {
	res = make([]Entry, 0, end-start)
	sd.Enumerate(start, end, func(key K, value V) {
		res = append(res, Entry{key, value})
	}, false)
	return
}

func (sd *SortedDict) Range(min, max K) (res []Entry) {
	keys := sd.sl.Range(min, max)
	res = make([]Entry, 0, len(keys))
	for _, key := range keys {
		res = append(res, Entry{key, sd.dict[key]})
	}
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

	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]K{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]K{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]K{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意容量的设置.
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
	sl.Enumerate(start, end, func(value K) {}, true)
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
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endIndex := min(len(block), startIndex+count)
		for j := startIndex; j < endIndex; j++ {
			f(block[j])
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == len(block) {
				// !delete block
				sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
				sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				sl._updateTree(pos, endIndex-startIndex)
				block = append(block[:startIndex], block[endIndex:]...)
				sl.mins[pos] = block[0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *SortedList) Slice(start, end int) []K {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return nil
	}
	count := end - start
	res := make([]K, 0, count)
	pos, index := sl._findKth(start)
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endPos := min(len(block), index+count)
		curCount := endPos - index
		res = append(res, block[index:endPos]...)
		count -= curCount
		index = 0
	}
	return res
}

func (sl *SortedList) Range(min, max K) []K {
	if sl.less(max, min) {
		return nil
	}
	res := []K{}
	pos := sl._locBlock(min)
	for i := pos; i < len(sl.blocks); i++ {
		block := sl.blocks[i]
		for j := 0; j < len(block); j++ {
			x := block[j]
			if sl.less(max, x) {
				return res
			}
			if !sl.less(x, min) {
				res = append(res, x)
			}
		}
	}
	return res
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
	sl.blocks[pos] = append(sl.blocks[pos][:index], sl.blocks[pos][index+1:]...)
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
	sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
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
