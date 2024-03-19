// API:
//  func (sl *SortedListWithSum) Add(value S) *SortedListWithSum                               {}
//  func (sl *SortedListWithSum) Has(value S) bool                                      {}
//  func (sl *SortedListWithSum) Discard(value S) bool                                  {}
//  func (sl *SortedListWithSum) Pop(index int) S                                       {}
//  func (sl *SortedListWithSum) At(index int) S                                        {}
//  func (sl *SortedListWithSum) Erase(start, end int)                                  {}
//
//  func (sl *SortedListWithSum) Lower(value S) S                                       {}
//  func (sl *SortedListWithSum) Higher(value S) S                                      {}
//  func (sl *SortedListWithSum) Floor(value S) S                                       {}
//  func (sl *SortedListWithSum) Ceiling(value S) S                                     {}
//
//  func (sl *SortedListWithSum) BisectLeft(value S) int                                {}
//  func (sl *SortedListWithSum) BisectRight(value S) int                               {}
//  func (sl *SortedListWithSum) Count(value S) int                                     {}
//
//  func (sl *SortedListWithSum) Clear()                                                {}
//
//  func (sl *SortedListWithSum) ForEach(f func(value S, index int) bool, reverse bool) {}
//  func (sl *SortedListWithSum) Enumerate(start, end int, f func(value S), erase bool) {}
//
//  func (sl *SortedListWithSum) Slice(start, end int) []S                              {}
//  func (sl *SortedListWithSum) Range(min, max S) []S                                  {}
//
//  func (sl *SortedListWithSum) IteratorAt(index int) *Iterator                        {}
//  func (sl *SortedListWithSum) LowerBound(value S) *Iterator                          {}
//  func (sl *SortedListWithSum) UpperBound(value S) *Iterator                          {}
//
//  func (sl *SortedListWithSum) Min() S                                                {}
//  func (sl *SortedListWithSum) Max() S                                                {}
//  func (sl *SortedListWithSum) String() string                                        {}
//  func (sl *SortedListWithSum) Len() int                                              {}

//  !func (sl *SortedListWithSum) SumSlice(start, end int) S 													  {}
//  !func (sl *SortedListWithSum) SumRange(min, max S) S
//  !func (sl *SortedListWithSum) SumAll() S

// test:
// https://leetcode.cn/problems/smallest-missing-genetic-value-in-each-subtree/submissions/
// https://leetcode.cn/problems/sliding-subarray-beauty/
// https://leetcode.cn/problems/count-the-number-of-fair-pairs/
// https://leetcode.cn/problems/minimum-difference-in-sums-after-removal-of-elements/
// https://atcoder.jp/contests/abc281/tasks/abc281_e

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func main() {
	// abc241_d()
	abc281_e()
}

const INF int = 1e18

// https://atcoder.jp/contests/abc281/tasks/abc281_e
func abc281_e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	sl := NewSortedListWithSum(func(a, b int) bool { return a < b })

	res := make([]int, 0, n-m+1)
	for i := 0; i < n; i++ {
		sl.Add(nums[i])
		if i >= m {
			sl.Discard(nums[i-m])
		}
		if i >= m-1 {
			res = append(res, sl.SumSlice(0, k))
		}
	}

	for _, x := range res {
		fmt.Fprint(out, x, " ")
	}
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 200

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(a E) E   { return -a }

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

	res := e()
	pos, index := sl._findKth(start)
	count := end - start
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endIndex := min(len(block), index+count)
		curCount := endIndex - index
		if curCount == len(block) {
			res = op(res, sl.sums[pos])
		} else {
			for j := index; j < endIndex; j++ {
				res = op(res, block[j])
			}
		}
		count -= curCount
		index = 0
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
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]E{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]
	sl.sums[pos] = op(sl.sums[pos], value)

	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		oldSum := sl.sums[pos]
		left, right := make([]E, _LOAD), make([]E, n-_LOAD)
		copy(left, sl.blocks[pos][:_LOAD])
		copy(right, sl.blocks[pos][_LOAD:])
		sl.blocks = append(sl.blocks[:pos], append([][]E{left, right}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos], append([]E{left[0], right[0]}, sl.mins[pos+1:]...)...)
		sl.shouldRebuildTree = true

		sl._rebuildSum(pos)
		newSum := op(oldSum, inv(sl.sums[pos]))
		sl.sums = append(sl.sums[:pos+1], append([]E{newSum}, sl.sums[pos+1:]...)...)
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
				sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
				sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
				sl.allSum = op(sl.allSum, inv(sl.sums[pos]))
				sl.sums = append(sl.sums[:pos], sl.sums[pos+1:]...)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				for i := startIndex; i < endIndex; i++ {
					inv_ := inv(block[i])
					sl.allSum = op(sl.allSum, inv_)
					sl.sums[pos] = op(sl.sums[pos], inv_)
				}
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

func (sl *SortedListWithSum) IteratorAt(index int) *Iterator {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("Index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl._iteratorAt(pos, startIndex)
}

func (sl *SortedListWithSum) LowerBound(value E) *Iterator {
	pos, index := sl._locLeft(value)
	return sl._iteratorAt(pos, index)
}

func (sl *SortedListWithSum) UpperBound(value E) *Iterator {
	pos, index := sl._locRight(value)
	return sl._iteratorAt(pos, index)
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
	sl.blocks[pos] = append(block[:index], block[index+1:]...)
	sl.allSum = op(sl.allSum, inv(deleted))
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		sl.sums[pos] = op(sl.sums[pos], inv(deleted))
		return
	}

	// !delete block
	sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
	sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
	sl.shouldRebuildTree = true
	sl.sums = append(sl.sums[:pos], sl.sums[pos+1:]...)
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

func (sl *SortedListWithSum) _iteratorAt(pos, index int) *Iterator {
	return &Iterator{sl: sl, pos: pos, index: index}
}

type Iterator struct {
	sl    *SortedListWithSum
	pos   int
	index int
}

func (it *Iterator) HasNext() bool {
	return it.pos < len(it.sl.blocks)-1 || it.index < len(it.sl.blocks[it.pos])-1
}

func (it *Iterator) Next() (res E, ok bool) {
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

func (it *Iterator) HasPrev() bool {
	return it.pos > 0 || it.index > 0
}

func (it *Iterator) Prev() (res E, ok bool) {
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

func (it *Iterator) Remove() {
	it.sl._delete(it.pos, it.index)
}

func (it *Iterator) Value() (res E, ok bool) {
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
