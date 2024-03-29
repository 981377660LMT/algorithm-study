package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"sort"
	"strings"
)

// 给你一个下标从 0 开始的二进制数组 nums，其长度为 n ；另给你一个 正整数 k 以及一个 非负整数 maxChanges 。

// 灵茶山艾府在玩一个游戏，游戏的目标是让灵茶山艾府使用 最少 数量的 行动 次数从 nums 中拾起 k 个 1 。游戏开始时，灵茶山艾府可以选择数组 [0, n - 1] 范围内的任何索引index 站立。如果 nums[index] == 1 ，灵茶山艾府就会拾起一个 1 ，并且 nums[index] 变成0（这 不算 作一次行动）。之后，灵茶山艾府可以执行 任意数量 的 行动（包括零次），在每次行动中灵茶山艾府必须 恰好 执行以下动作之一：

// 选择任意一个下标 j != index 且满足 nums[j] == 0 ，然后将 nums[j] 设置为 1 。这个动作最多可以执行 maxChanges 次。
// 选择任意两个相邻的下标 x 和 y（|x - y| == 1）且满足 nums[x] == 1, nums[y] == 0 ，然后交换它们的值（将 nums[y] = 1 和 nums[x] = 0）。如果 y == index，在这次行动后灵茶山艾府拾起一个 1 ，并且 nums[y] 变成 0 。
// 返回灵茶山艾府拾起 恰好 k 个 1 所需的 最少 行动次数。
func main() {
	// nums = [0,0,0,0], k = 2, maxChanges = 3
	fmt.Println(minimumMoves([]int{0, 1, 0, 1, 0, 1, 0, 1}, 5, 1)) // 2
}

const INF int = 1e18

func minimumMoves(nums []int, k int, maxChanges int) int64 {
	fastPath := func(startIndex int) int {
		ones := 0
		step := 0
		if nums[startIndex] == 1 {
			ones++
		}
		if ones >= k {
			return 0
		}
		if startIndex-1 >= 0 && nums[startIndex-1] == 1 {
			ones++
			step++
			if ones >= k {
				return step
			}
		}
		if startIndex+1 < len(nums) && nums[startIndex+1] == 1 {
			ones++
			step++
			if ones >= k {
				return step
			}
		}
		diff := k - ones
		if diff > maxChanges {
			return INF
		}
		return step + diff*2
	}

	n := len(nums)
	res := INF
	for i := 0; i < n; i++ {
		tmp := fastPath(i)
		if tmp != INF {
			res = min(res, tmp)
		}
	}
	if res != INF {
		return int64(res)
	}

	// 枚举起点，把maxChanges用完还不能成功
	// 这时候一定存在1
	baseStep := maxChanges * 2
	need := k - maxChanges // need个1到startIndex(包括startIndex自己的1)需要的最小总移动次数
	minDist := make([]int, n)
	for i := 0; i < n; i++ {
		minDist[i] = INF
	}

	sl := NewSortedListWithSum(func(a, b int) bool { return a < b })
	mf := NewMedianFinderSortedList(sl)
	// 连续k个1对答案的贡献
	onesIndex := []int{}
	for i, v := range nums {
		if v == 1 {
			onesIndex = append(onesIndex, i)
		}
	}

	right := need - 1
	for i := 0; i < need; i++ {
		sl.Add(onesIndex[i])
	}
	for i := 0; i < n; i++ {
		cost1 := mf.DistSum(i)
		minDist[i] = cost1
		// 右边加一个，左边退一个是否更优
		if right+1 < len(onesIndex) {
			added := onesIndex[right+1]
			sl.Add(added)
			right++
			popped := sl.Pop(0)
			cost2 := mf.DistSum(i)
			if cost2 < cost1 {
				minDist[i] = cost2
			} else {
				// rollback
				sl.Pop(-1)
				sl.Add(popped)
				right--
			}
		}
	}

	for i := 0; i < n; i++ {
		res = min(res, minDist[i])
	}
	return int64(res + baseStep)
}

// SortedList 动态维护中位数信息.
// `Proxy 的内部持有一个对 SortedList 的引用`.
type MedianFinderSortedList struct {
	List *SortedListWithSum
}

func NewMedianFinderSortedList(sortedList *SortedListWithSum) *MedianFinderSortedList {
	return &MedianFinderSortedList{List: sortedList}
}

// 返回向下取整的中位数.
func (mfs *MedianFinderSortedList) Median() (res int) {
	if mfs.List.Len() == 0 {
		panic("Median() called on empty MedianFinderSortedList")
	}
	if mfs.List.Len()&1 == 1 {
		return mfs.List.At(mfs.List.Len() >> 1)
	} else {
		mid := mfs.List.Len() >> 1
		return (mfs.List.At(mid-1) + mfs.List.At(mid)) >> 1
	}
}

// 返回所有数到`to`的距离和.
func (mfs *MedianFinderSortedList) DistSum(to int) int {
	sl := mfs.List
	pos := sl.BisectRight(to)
	allSum := sl.SumAll()
	var sum1, sum2 int
	if pos < sl.Len()>>1 {
		sum1 = sl.SumSlice(0, pos)
		sum2 = allSum - sum1
	} else {
		sum2 = sl.SumSlice(pos, sl.Len())
		sum1 = allSum - sum2
	}
	leftSum := to*pos - sum1
	rightSum := sum2 - to*(sl.Len()-pos)
	return leftSum + rightSum
}

// 返回切片`[start,end)`中所有数到`to`的距离和.
func (mfs *MedianFinderSortedList) DistSumRange(to int, start, end int) int {
	sl := mfs.List
	if start < 0 {
		start = 0
	}
	if end > sl.Len() {
		end = sl.Len()
	}
	if start >= end {
		return 0
	}
	pos := sl.BisectLeft(to)
	if pos <= start {
		return sl.SumSlice(start, end) - to*(end-start)
	}
	if pos >= end {
		return to*(end-start) - sl.SumSlice(start, end)
	}
	leftSum := to*(pos-start) - sl.SumSlice(start, pos)
	rightSum := sl.SumSlice(pos, end) - to*(end-pos)
	return leftSum + rightSum
}

// 返回所有数到中位数的距离和.
func (mfs *MedianFinderSortedList) DistSumToMedian() int {
	if mfs.List.Len() == 0 {
		return 0
	}
	return mfs.DistSum(mfs.Median())
}

// 返回切片`[start,end)`中所有数到中位数的距离和.
func (mfs *MedianFinderSortedList) DistSumToMedianRange(start, end int) int {
	if mfs.List.Len() == 0 {
		return 0
	}
	return mfs.DistSumRange(mfs.Median(), start, end)
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 500

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

// RangeMinPointGet

type Id = int

const COMMUTATIVE = true

func (*SegmentTreeDual) id() Id                 { return INF }
func (*SegmentTreeDual) composition(f, g Id) Id { return min(f, g) }

type SegmentTreeDual struct {
	n            int
	size, height int
	lazy         []Id
	unit         Id
}

func NewSegmentTreeDual(n int) *SegmentTreeDual {
	res := &SegmentTreeDual{}
	size := 1
	height := 0
	for size < n {
		size <<= 1
		height++
	}
	lazy := make([]Id, 2*size)
	unit := res.id()
	for i := 0; i < 2*size; i++ {
		lazy[i] = unit
	}
	res.n = n
	res.size = size
	res.height = height
	res.lazy = lazy
	res.unit = unit
	return res
}
func (seg *SegmentTreeDual) Get(index int) Id {
	index += seg.size
	for i := seg.height; i > 0; i-- {
		seg.propagate(index >> i)
	}
	return seg.lazy[index]
}
func (seg *SegmentTreeDual) GetAll() []Id {
	for i := 0; i < seg.size; i++ {
		seg.propagate(i)
	}
	res := make([]Id, seg.n)
	copy(res, seg.lazy[seg.size:seg.size+seg.n])
	return res
}
func (seg *SegmentTreeDual) Update(left, right int, value Id) {
	if left < 0 {
		left = 0
	}
	if right > seg.n {
		right = seg.n
	}
	if left >= right {
		return
	}
	left += seg.size
	right += seg.size
	if !COMMUTATIVE {
		for i := seg.height; i > 0; i-- {
			if (left>>i)<<i != left {
				seg.propagate(left >> i)
			}
			if (right>>i)<<i != right {
				seg.propagate((right - 1) >> i)
			}
		}
	}
	for left < right {
		if left&1 > 0 {
			seg.lazy[left] = seg.composition(value, seg.lazy[left])
			left++
		}
		if right&1 > 0 {
			right--
			seg.lazy[right] = seg.composition(value, seg.lazy[right])
		}
		left >>= 1
		right >>= 1
	}
}
func (seg *SegmentTreeDual) propagate(k int) {
	if seg.lazy[k] != seg.unit {
		seg.lazy[k<<1] = seg.composition(seg.lazy[k], seg.lazy[k<<1])
		seg.lazy[k<<1|1] = seg.composition(seg.lazy[k], seg.lazy[k<<1|1])
		seg.lazy[k] = seg.unit
	}
}
func (st *SegmentTreeDual) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := 0; i < st.n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprint(st.Get(i)))
	}
	buf.WriteByte(']')
	return buf.String()
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
