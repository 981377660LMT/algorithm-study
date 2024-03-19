// https://atcoder.jp/contests/abc324/editorial/7399

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

// 给定一个数组，数组元素为1-n的排列
// 有两种操作：
// 1.把 A[version]中下标大于等于 x 的元素分裂成一个新的数组 Ai(A[version]中保留x个)。
// 2.把 A[version]中值大于 x 的元素分裂成一个新的数组 Ai。
// 这两种操作都不会改变元素相对顺序。
// 输出每次分裂出的数组大小。
//
// 两个 SortedList 维护.
// SortedList 启发式分裂：每次分裂出较小的那一半
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)

	history := make([]*SortedListWithEntry, q+1)
	entries := make([]Entry, n)
	for i := range entries {
		entries[i] = Entry{index: int32(i), value: int32(nums[i])}
	}
	history[0] = NewSortedListWithEntry(
		func(a, b Entry) bool { return a.index < b.index },
		func(a, b Entry) bool { return a.value < b.value },
		entries...,
	)
	for i := 1; i < len(history); i++ {
		history[i] = NewSortedListWithEntry(
			func(a, b Entry) bool { return a.index < b.index },
			func(a, b Entry) bool { return a.value < b.value },
		)
	}

	for cur := 1; cur < q+1; cur++ {
		var kind, pre, x int
		fmt.Fscan(in, &kind, &pre, &x)

		if kind == 1 { // 将 A[pre] 中下标大于等于 x 的元素分裂成一个新的数组 Ai
			len1 := x
			len2 := history[pre].Len() - x
			if len1 < len2 { // 前面少，拆到前面
				history[cur], history[pre] = history[pre], history[cur]
				for j := 0; j < len1; j++ {
					history[pre].Add(history[cur].PopFront())
				}
			} else { // 后面少，拆到后面
				for j := 0; j < len2; j++ {
					history[cur].Add(history[pre].PopBack())
				}
			}
		} else { // 将 A[pre] 中值大于 x 的元素分裂成一个新的数组 Ai
			len1 := history[pre].SlValue.BisectRight(Entry{value: int32(x)})
			len2 := history[pre].Len() - len1
			if len1 < len2 { // 前面少，拆到前面
				history[cur], history[pre] = history[pre], history[cur]
				for j := 0; j < len1; j++ {
					history[pre].Add(history[cur].PopMin())
				}
			} else { // 后面少，拆到后面
				for j := 0; j < len2; j++ {
					history[cur].Add(history[pre].PopMax())
				}
			}
		}

		fmt.Fprintln(out, history[cur].Len())
	}

}

const _LOAD int = 50 // !块大小较小，适合频繁分裂的情形

type Value = int32 // value 必须是基本类型
type Entry = struct {
	index int32
	value Value
}

type SortedListWithEntry struct {
	SlIndex *SortedList // 按照index排序
	SlValue *SortedList // 按照value排序
}

func NewSortedListWithEntry(
	indexLess func(a, b Entry) bool,
	valueLess func(a, b Entry) bool,
	entries ...Entry,
) *SortedListWithEntry {
	sl1 := NewSortedList(indexLess, entries...)
	sl2 := NewSortedList(valueLess, entries...)
	return &SortedListWithEntry{sl1, sl2}
}

func (sl *SortedListWithEntry) Add(value Entry) *SortedListWithEntry {
	sl.SlIndex.Add(value)
	sl.SlValue.Add(value)
	return sl
}

func (sl *SortedListWithEntry) GetKthEntryByIndex(kth int) Entry {
	return sl.SlIndex.At(kth)
}

func (sl *SortedListWithEntry) GetKthEntryByValue(kth int) Entry {
	return sl.SlValue.At(kth)
}

func (sl *SortedListWithEntry) Front() Entry {
	return sl.SlIndex.Min()
}

func (sl *SortedListWithEntry) Back() Entry {
	return sl.SlIndex.Max()
}

func (sl *SortedListWithEntry) Min() Entry {
	return sl.SlValue.Min()
}

func (sl *SortedListWithEntry) Max() Entry {
	return sl.SlValue.Max()
}

func (sl *SortedListWithEntry) PopFront() Entry {
	popped := sl.SlIndex.Pop(0)
	sl.SlValue.Discard(popped)
	return popped
}

func (sl *SortedListWithEntry) PopBack() Entry {
	popped := sl.SlIndex.Pop(sl.Len() - 1)
	sl.SlValue.Discard(popped)
	return popped
}

func (sl *SortedListWithEntry) PopMin() Entry {
	popped := sl.SlValue.Pop(0)
	sl.SlIndex.Discard(popped)
	return popped
}

func (sl *SortedListWithEntry) PopMax() Entry {
	popped := sl.SlValue.Pop(sl.Len() - 1)
	sl.SlIndex.Discard(popped)
	return popped
}

func (sl *SortedListWithEntry) Discard(value Entry) bool {
	return sl.SlIndex.Discard(value) && sl.SlValue.Discard(value)
}

func (sl *SortedListWithEntry) Len() int {
	return sl.SlIndex.Len()
}

// 使用分块+树状数组维护的有序序列.
type SortedList struct {
	less              func(a, b Entry) bool
	size              int
	blocks            [][]Entry
	mins              []Entry
	tree              []int
	shouldRebuildTree bool
}

func NewSortedList(less func(a, b Entry) bool, elements ...Entry) *SortedList {
	elements = append(elements[:0:0], elements...)
	res := &SortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]Entry{}
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]Entry, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	return res
}

func (sl *SortedList) Add(value Entry) *SortedList {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []Entry{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return sl
	}

	pos, index := sl._locRight(value)

	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]Entry{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]Entry{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]Entry{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SortedList) Has(value Entry) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < len(sl.blocks[pos]) && sl.blocks[pos][index] == value
}

func (sl *SortedList) Discard(value Entry) bool {
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

func (sl *SortedList) Pop(index int) Entry {
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

func (sl *SortedList) At(index int) Entry {
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

func (sl *SortedList) Lower(value Entry) (res Entry, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList) Higher(value Entry) (res Entry, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedList) Floor(value Entry) (res Entry, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList) Ceiling(value Entry) (res Entry, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedList) BisectLeft(value Entry) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedList) BisectRight(value Entry) int {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedList) Count(value Entry) int {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SortedList) ForEach(f func(value Entry, index int), reverse bool) {
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

func (sl *SortedList) Enumerate(start, end int, f func(value Entry), erase bool) {
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
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
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
				sl._updateTree(pos, -deleted)
				block = append(block[:startIndex], block[endIndex:]...)
				sl.mins[pos] = block[0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *SortedList) Slice(start, end int) []Entry {
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
	res := make([]Entry, 0, count)
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

func (sl *SortedList) Range(min, max Entry) []Entry {
	if sl.less(max, min) {
		return nil
	}
	res := []Entry{}
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

func (sl *SortedList) Min() Entry {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *SortedList) Max() Entry {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedList) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedList{")
	sl.ForEach(func(value Entry, index int) {
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

func (sl *SortedList) _locLeft(value Entry) (pos, index int) {
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

func (sl *SortedList) _locRight(value Entry) (pos, index int) {
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

func (sl *SortedList) _locBlock(value Entry) int {
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
