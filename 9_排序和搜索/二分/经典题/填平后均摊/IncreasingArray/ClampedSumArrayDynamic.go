// 模拟动态修改的waveletMatrix (multiset区间和).
// waveletMatrixDynamic/waveletMatrixWithSumDynamic/multisetSum/multisetRangeSum
//
// api:
//
// Query:
//  SumWithMin(start, end int32, min int) int
//  SumWithMax(start, end int32, max int) int
//  SumWithMinAndMax(start, end int32, min, max int) int
//  CountRange(start, end int32, min, max int) int
//  SumRange(start, end int32, min, max int) int
//  CountAndSumRange(start, end int32, min, max int) (int, int)
//
// Update:
//  通过修改內部的sl修改.

package main

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
	"time"
)

func main() {
	// test()
	testTime()
}

func demo() {
	sl := NewSortedListWithSum(func(a, b int) bool { return a < b }, 3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5)
	S := NewClampedSumArray(sl)
	fmt.Println(S.Sl)
	fmt.Println(S.SumWithMinAndMax(0, 10, 2, 5))
}

type ClampedSumArrayDynamic struct {
	Sl *SortedListWithSum
}

func NewClampedSumArray(sl *SortedListWithSum) *ClampedSumArrayDynamic {
	return &ClampedSumArrayDynamic{Sl: sl}
}

// [min, ?)
func (a *ClampedSumArrayDynamic) SumWithMin(start, end int, min int) int {
	if start < 0 {
		start = 0
	}
	if n := a.Sl.Len(); end > n {
		end = n
	}
	if start >= end {
		return 0
	}
	lessCount := a.bisectLeft(start, end, min) - start
	largerSum := a.Sl.SumSlice(start+lessCount, end)
	return min*lessCount + largerSum
}

// (?, max]
func (a *ClampedSumArrayDynamic) SumWithMax(start, end int, max int) int {
	if start < 0 {
		start = 0
	}
	if n := a.Sl.Len(); end > n {
		end = n
	}
	if start >= end {
		return 0
	}
	lessCount := a.bisectLeft(start, end, max) - start
	lessSum := a.Sl.SumSlice(start, start+lessCount)
	return lessSum + max*(end-start-lessCount)
}

// [min, max]
func (a *ClampedSumArrayDynamic) SumWithMinAndMax(start, end int, min, max int) int {
	if min > max {
		return 0
	}
	if start < 0 {
		start = 0
	}
	if n := a.Sl.Len(); end > n {
		end = n
	}
	if start >= end {
		return 0
	}
	powLow := a.bisectLeft(start, end, min) - start
	posUp := a.bisectLeft(start, end, max) - start
	return a.Sl.SumSlice(start+powLow, start+posUp) + min*powLow + max*(end-start-posUp)
}

// !WaveletMatrixLike Api

// [start,end) x [y1,y2) 中的数的个数.
func (a *ClampedSumArrayDynamic) CountRange(start, end int, y1, y2 int) int {
	if y1 >= y2 {
		return 0
	}
	if start < 0 {
		start = 0
	}
	if n := a.Sl.Len(); end > n {
		end = n
	}
	if start >= end {
		return 0
	}
	return a.bisectLeft(start, end, y2) - a.bisectLeft(start, end, y1)
}

// [start,end) x [y1,y2) 中的数的和.
func (a *ClampedSumArrayDynamic) SumRange(start, end int, y1, y2 int) int {
	_, sum := a.CountAndSumRange(start, end, y1, y2)
	return sum
}

// [start,end) x [y1,y2) 中的数的个数、和.
func (a *ClampedSumArrayDynamic) CountAndSumRange(start, end int, y1, y2 int) (int, int) {
	if y1 >= y2 {
		return 0, 0
	}
	if start < 0 {
		start = 0
	}
	if n := a.Sl.Len(); end > n {
		end = n
	}
	if start >= end {
		return 0, 0
	}
	left := a.bisectLeft(start, end, y1) - start
	right := a.bisectLeft(start, end, y2) - start
	return right - left, a.Sl.SumSlice(start+left, start+right)
}

func (a *ClampedSumArrayDynamic) bisectLeft(start, end, v int) int {
	return MinLeft(end, func(mid int) bool { return a.Sl.At(mid) >= v }, start)
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含，使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
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
				sl._updateTree(pos, -deleted)
				sl.blocks[pos] = append(sl.blocks[pos][:startIndex], sl.blocks[pos][endIndex:]...)
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

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func test() {
	nums := []int{1, 2, 3, 4, 5}
	sl := NewSortedListWithSum(func(a, b int) bool { return a < b }, nums...)
	A := NewClampedSumArray(sl)

	// for i := 1; i <= 24; i++ {
	// 	fmt.Println(A.Increase(i))
	// }
	// fmt.Println("------")
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(A.SumWithUpClamp(i))
	// }
	// fmt.Println("------")
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(A.SumWithLowClamp(i))
	// }
	// fmt.Println("------")
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(A.DiffSum(i))
	// }
	// fmt.Println("------")
	// for i := 1; i <= 24; i++ {
	// 	fmt.Println(A.Decrease(i))
	// }
	// fmt.Println("------")
	// for i := 0; i < 200; i++ {
	// 	fmt.Println(i, A.IncreaseForArray(i))
	// 	assert(mins(A.IncreaseForArray(i)...) == A.IncreaseForMin(i), "min")
	// }
	// fmt.Println("------")
	// for i := 0; i < 200; i++ {
	// 	fmt.Println(i, A.DecreaseForArray(i))
	// 	assert(maxs(A.DecreaseForArray(i)...) == A.DecreaseForMax(i), "max")
	// }

	upClampRangeBruteForce := func(v int, start, end int) int {
		sum := 0
		for i := start; i < end; i++ {
			sum += min(nums[i], v)
		}
		return sum
	}

	lowClampRangeBruteForce := func(v int, start, end int) int {
		sum := 0
		for i := start; i < end; i++ {
			sum += max(nums[i], v)
		}
		return sum
	}

	upAndLowClampRangeBruteForce := func(low, up int, start, end int) int {
		sum := 0
		for i := start; i < end; i++ {
			sum += max(min(nums[i], up), low)
		}
		return sum
	}

	countRangeBruteForce := func(start, end int, y1, y2 int) int {
		sum := 0
		for i := start; i < end; i++ {
			if y1 <= nums[i] && nums[i] < y2 {
				sum++
			}

		}
		return sum
	}

	countAndSumRangeBruteForce := func(start, end int, y1, y2 int) (int, int) {
		count, sum := 0, 0
		for i := start; i < end; i++ {
			if y1 <= nums[i] && nums[i] < y2 {
				count++
				sum += nums[i]
			}
		}
		return count, sum
	}

	{
		for i := 0; i < len(nums); i++ {
			for j := 0; j < len(nums); j++ {
				for v := -10; v < 10; v++ {
					assert(A.SumWithMax(i, j, v) == upClampRangeBruteForce(v, i, j), "upClampRange")
					assert(A.SumWithMin(i, j, v) == lowClampRangeBruteForce(v, i, j), "lowClampRange")

				}
			}
		}

		for i := 0; i < len(nums); i++ {
			for j := i; j <= len(nums); j++ {
				for y1 := -10; y1 < 10; y1++ {
					for y2 := y1; y2 < 10; y2++ {
						assert(A.SumWithMinAndMax(i, j, y1, y2) == upAndLowClampRangeBruteForce(y1, y2, i, j), "upAndLowClampRange")
					}
				}
			}
		}

		for i := 0; i < len(nums); i++ {
			for j := i; j <= len(nums); j++ {
				for y1 := -10; y1 < 10; y1++ {
					for y2 := y1; y2 < 10; y2++ {
						assert(A.CountRange(i, j, y1, y2) == countRangeBruteForce(i, j, y1, y2), "countRange")
					}
				}
			}
		}
	}

	{
		for i := 0; i < len(nums); i++ {
			for j := i; j <= len(nums); j++ {
				for y1 := -10; y1 < 10; y1++ {
					for y2 := y1; y2 < 10; y2++ {
						count, sum := A.CountAndSumRange(i, j, y1, y2)
						countBrute, sumBrute := countAndSumRangeBruteForce(i, j, y1, y2)
						assert(count == countBrute, "countAndSumRange count")
						assert(sum == sumBrute, "countAndSumRange sum")
					}
				}
			}
		}
	}

	fmt.Println("pass")

}

func testTime() {
	const N = 1e5
	nums := make([]int, N)
	for i := 0; i < N; i++ {
		nums[i] = i
	}
	sl := NewSortedListWithSum(func(a, b int) bool { return a < b }, nums...)
	time1 := time.Now()
	A := NewClampedSumArray(sl)
	for i := 0; i < len(nums); i++ {
		A.SumWithMinAndMax(0, len(nums), 0, N)
	}
	time2 := time.Now()
	fmt.Println(time2.Sub(time1))
}
