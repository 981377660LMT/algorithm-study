// 3321. 计算子数组的 x-sum II-SortedListWithSun
// https://leetcode.cn/problems/find-x-sum-of-all-k-long-subarrays-ii/description/

package main

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
	"unsafe"
)

func findXSum(nums []int, k int, x int) []int64 {
	sl := NewSortedListWithSum(func(a, b E) bool {
		if a.count != b.count {
			return a.count > b.count
		}
		return a.value > b.value
	})
	counter := make(map[int32]int32, len(nums))

	add := func(v int32) {
		preFreq := counter[v]
		counter[v]++
		if preFreq != 0 {
			sl.Discard(NewE(preFreq, v))
		}
		sl.Add(NewE(preFreq+1, v))
	}
	remove := func(v int32) {
		preFreq := counter[v]
		counter[v]--
		sl.Discard(NewE(preFreq, v))
		if preFreq > 1 {
			sl.Add(NewE(preFreq-1, v))
		}
	}

	n := len(nums)
	res := make([]int64, 0, n-k+1)
	for right := 0; right < n; right++ {
		add(int32(nums[right]))
		if right >= k {
			remove(int32(nums[right-k]))
		}
		if right >= k-1 {
			res = append(res, int64(sl.SumSlice(0, x).sum))
		}
	}
	return res
}

const _LOAD int = 100

type E = struct {
	count, value int32
	sum          int
}

func NewE(count, value int32) E {
	return E{count: count, value: value, sum: int(count) * int(value)}
}

func e() E { return E{} }
func op(a, b E) E {
	a.sum += b.sum
	return a
}
func inv(a E) E {
	a.sum = -a.sum
	return a
}

// 使用分块+树状数组维护的有序序列.
type SortedListWithSum struct {
	less              func(a, b E) bool
	size              int
	blocks            [][]E
	mins              []E
	tree              []int
	shouldRebuildTree bool

	sums   []E
	allSum E
}

func NewSortedListWithSum(less func(a, b E) bool, elements ...E) *SortedListWithSum {
	elements = append(elements[:0:0], elements...)
	res := &SortedListWithSum{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]E{}
	sums := []E{}
	allSum := e()
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
		newBlock := elements[start:end:end] // !各个块互不影响, max参数也需要指定为end
		blocks = append(blocks, newBlock)
		cur := e()
		for _, v := range newBlock {
			cur = op(cur, v)
		}
		sums = append(sums, cur)
		allSum = op(allSum, cur)
	}
	mins := make([]E, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	res.sums = sums
	res.allSum = allSum
	return res
}

// 返回区间`[start, end)`的和.
func (sl *SortedListWithSum) SumSlice(start, end int) E {
	if start < 0 {
		start += sl.size
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl.size
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return e()
	}

	bid1, startIndex1 := sl._findKth(start)
	bid2, startIndex2 := sl._findKth(end)
	start, end = startIndex1, startIndex2
	res := e()
	if bid1 == bid2 {
		block := sl.blocks[bid1]
		for i := start; i < end; i++ {
			res = op(res, block[i])
		}
	} else {
		if start < len(sl.blocks[bid1]) {
			block1 := sl.blocks[bid1]
			for i := start; i < len(block1); i++ {
				res = op(res, block1[i])
			}
		}
		for i := bid1 + 1; i < bid2; i++ {
			res = op(res, sl.sums[i])
		}
		if m := len(sl.blocks); bid2 < m && end > 0 {
			block2 := sl.blocks[bid2]
			tmp := e()
			for i := 0; i < end; i++ {
				tmp = op(tmp, block2[i])
			}
			res = op(res, tmp)
		}
	}
	return res
}

// 返回范围`[min, max]`的和.
func (sl *SortedListWithSum) SumRange(min, max E) E {
	if sl.less(max, min) {
		return e()
	}
	res := e()
	pos, start := sl._locLeft(min)
	for i := pos; i < len(sl.blocks); i++ {
		block := sl.blocks[i]
		if sl.less(max, block[0]) {
			break
		}
		if start == 0 && !sl.less(block[len(block)-1], max) {
			res = op(res, sl.sums[i])
		} else {
			for j := start; j < len(block); j++ {
				cur := block[j]
				if sl.less(max, cur) {
					break
				}
				res = op(res, cur)
			}
		}
		start = 0
	}
	return res
}

func (sl *SortedListWithSum) SumAll() E {
	return sl.allSum
}

func (sl *SortedListWithSum) Add(value E) *SortedListWithSum {
	sl.size++
	sl.allSum = op(sl.allSum, value)
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []E{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		sl.sums = append(sl.sums, value)
		return sl
	}

	pos, index := sl._locRight(value)
	sl._updateTree(pos, 1)
	sl.blocks[pos] = Insert(sl.blocks[pos], index, value)
	sl.mins[pos] = sl.blocks[pos][0]
	sl.sums[pos] = op(sl.sums[pos], value)

	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		oldSum := sl.sums[pos]
		left := append([]E(nil), sl.blocks[pos][:_LOAD]...)
		right := append([]E(nil), sl.blocks[pos][_LOAD:]...)
		sl.blocks = Replace(sl.blocks, pos, pos+1, left, right)
		sl.mins = Replace(sl.mins, pos, pos+1, left[0], right[0])
		sl.shouldRebuildTree = true

		sl._rebuildSum(pos)
		newSum := op(oldSum, inv(sl.sums[pos]))
		sl.sums = Insert(sl.sums, pos+1, newSum)
	}

	return sl
}

func (sl *SortedListWithSum) Has(value E) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < len(sl.blocks[pos]) && sl.blocks[pos][index] == value
}

func (sl *SortedListWithSum) Discard(value E) bool {
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

func (sl *SortedListWithSum) Pop(index int) E {
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

func (sl *SortedListWithSum) At(index int) E {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SortedListWithSum) Erase(start, end int) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedListWithSum) Lower(value E) (res E, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedListWithSum) Higher(value E) (res E, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedListWithSum) Floor(value E) (res E, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedListWithSum) Ceiling(value E) (res E, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedListWithSum) BisectLeft(value E) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedListWithSum) BisectRight(value E) int {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedListWithSum) Count(value E) int {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedListWithSum) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
	sl.sums = sl.sums[:0]
	sl.allSum = e()
}

func (sl *SortedListWithSum) ForEach(f func(value E, index int), reverse bool) {
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

func (sl *SortedListWithSum) Enumerate(start, end int, f func(value E), erase bool) {
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
				sl.blocks = Replace(sl.blocks, pos, pos+1)
				sl.mins = Replace(sl.mins, pos, pos+1)
				sl.allSum = op(sl.allSum, inv(sl.sums[pos]))
				sl.sums = Replace(sl.sums, pos, pos+1)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				for i := startIndex; i < endIndex; i++ {
					inv_ := inv(block[i])
					sl.allSum = op(sl.allSum, inv_)
					sl.sums[pos] = op(sl.sums[pos], inv_)
				}
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

func (sl *SortedListWithSum) Slice(start, end int) []E {
	if start < 0 {
		start += sl.size
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl.size
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return nil
	}
	count := end - start
	res := make([]E, 0, count)
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

func (sl *SortedListWithSum) Range(min, max E) []E {
	if sl.less(max, min) {
		return nil
	}
	res := []E{}
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

func (sl *SortedListWithSum) Min() E {
	if sl.size == 0 {
		panic("Min() called on empty SortedListWithSum")
	}
	return sl.mins[0]
}

func (sl *SortedListWithSum) Max() E {
	if sl.size == 0 {
		panic("Max() called on empty SortedListWithSum")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedListWithSum) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedListWithSum{")
	sl.ForEach(func(value E, index int) {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *SortedListWithSum) Len() int {
	return sl.size
}

func (sl *SortedListWithSum) _delete(pos, index int) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	block := sl.blocks[pos]
	deleted := block[index]
	sl.blocks[pos] = Replace(block, index, index+1)
	sl.allSum = op(sl.allSum, inv(deleted))
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		sl.sums[pos] = op(sl.sums[pos], inv(deleted))
		return
	}

	// !delete block
	sl.blocks = Replace(sl.blocks, pos, pos+1)
	sl.mins = Replace(sl.mins, pos, pos+1)
	sl.shouldRebuildTree = true
	sl.sums = Replace(sl.sums, pos, pos+1)
}

func (sl *SortedListWithSum) _locLeft(value E) (pos, index int) {
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

func (sl *SortedListWithSum) _locRight(value E) (pos, index int) {
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

func (sl *SortedListWithSum) _locBlock(value E) int {
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

func (sl *SortedListWithSum) _buildTree() {
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

func (sl *SortedListWithSum) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SortedListWithSum) _queryTree(end int) int {
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

func (sl *SortedListWithSum) _findKth(k int) (pos, index int) {
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

func (sl *SortedListWithSum) _rebuildSum(pos int) {
	block := sl.blocks[pos]
	cur := e()
	for _, v := range block {
		cur = op(cur, v)
	}
	sl.sums[pos] = cur
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// Insert inserts the values v... into s at index i,
// !returning the modified slice.
// The elements at s[i:] are shifted up to make room.
// In the returned slice r, r[i] == v[0],
// and r[i+len(v)] == value originally at r[i].
// This function is O(len(s) + len(v)).
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

func rotateRight[E any](s []E, r int) {
	rotateLeft(s, len(s)-r)
}

func swap[E any](x, y []E) {
	for i := 0; i < len(x); i++ {
		x[i], y[i] = y[i], x[i]
	}
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

type AddRemoveCounter[T comparable] map[T]int32

func NewAddRemoveCounter[T comparable](cap int32) AddRemoveCounter[T] {
	return make(AddRemoveCounter[T], cap)
}

func (c AddRemoveCounter[T]) Add(x T, onAdd, onRemove func(v T, c /* c > 0 */ int32)) {
	preFreq := c[x]
	c[x]++
	if preFreq != 0 {
		onRemove(x, preFreq)
	}
	onAdd(x, preFreq+1)
}

func (c AddRemoveCounter[T]) Remove(x T, onAdd, onRemove func(v T, c /* c > 0 */ int32)) bool {
	preFreq := c[x]
	if preFreq == 0 {
		return false
	}
	c[x]--
	onRemove(x, preFreq)
	if preFreq > 1 {
		onAdd(x, preFreq-1)
	} else {
		delete(c, x)
	}
	return true
}
