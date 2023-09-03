// SortedListRangeBlock
// 值域分块，O(1)插入、删除，O(sqrt(n))查询
// 一般配合莫队算法使用.

package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	sl := NewSortedListRangeBlock(10)
	sl.Add(1)
	fmt.Println(sl.Min())
	fmt.Println(sl.Max())
}

func getSubarrayBeauty(nums []int, k int, x int) []int {
	res := []int{}
	sl := NewSortedListRangeBlock(200)
	OFFSET := 100
	n := len(nums)
	for right := 0; right < n; right++ {
		sl.Add(nums[right] + OFFSET)
		if right >= k {
			sl.Discard(nums[right-k] + OFFSET)
		}
		if right >= k-1 {
			xth := sl.At(x-1) - OFFSET
			if xth > 0 {
				xth = 0
			}
			res = append(res, xth)
		}
	}
	return res
}

func minimumDifference(nums []int) int64 {
	n := len(nums) / 3
	pre := NewSortedListRangeBlock(1e5+10, nums[:n]...)
	suf := NewSortedListRangeBlock(1e5+10, nums[n:]...)
	res := pre.SumSlice(0, n) - suf.SumSlice(suf.Len()-n, suf.Len())
	for i := n; i < 2*n; i++ {
		pre.Add(nums[i])
		suf.Remove(nums[i])
		res = min(res, pre.SumSlice(0, n)-suf.SumSlice(suf.Len()-n, suf.Len()))
	}
	return int64(res)
}

const INF int = 1e18

type SortedListRangeBlock struct {
	_blockSize  int   // 每个块的大小.
	_counter    []int // 每个数出现的次数.
	_blockCount []int // 每个块的数的个数.
	_blockSum   []int // 每个块的和.
	_belong     []int // 每个数所在的块.
	_len        int   // 所有数的个数.
}

// 值域分块模拟SortedList.
// `O(1)`add/remove，`O(sqrt(n))`查询.
// 一般配合莫队算法使用.
//
//	max:值域的最大值.0 <= max <= 1e6.
//	iterable:初始值.
func NewSortedListRangeBlock(max int, nums ...int) *SortedListRangeBlock {
	max += 5
	size := int(math.Sqrt(float64(max)))
	count := 1 + (max / size)
	sl := &SortedListRangeBlock{
		_blockSize:  size,
		_counter:    make([]int, max),
		_blockCount: make([]int, count),
		_blockSum:   make([]int, count),
		_belong:     make([]int, max),
	}
	for i := 0; i < max; i++ {
		sl._belong[i] = i / size
	}
	if len(nums) > 0 {
		sl.Update(nums...)
	}
	return sl
}

// O(1).
func (sl *SortedListRangeBlock) Add(value int) {
	sl._counter[value]++
	pos := sl._belong[value]
	sl._blockCount[pos]++
	sl._blockSum[pos] += value
	sl._len++
}

// O(1).
func (sl *SortedListRangeBlock) Remove(value int) {
	sl._counter[value]--
	pos := sl._belong[value]
	sl._blockCount[pos]--
	sl._blockSum[pos] -= value
	sl._len--
}

// O(1).
func (sl *SortedListRangeBlock) Discard(value int) bool {
	if !sl.Has(value) {
		return false
	}
	sl.Remove(value)
	return true
}

// O(1).
func (sl *SortedListRangeBlock) Has(value int) bool {
	return sl._counter[value] > 0
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) At(index int) int {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	for i := 0; i < len(sl._blockCount); i++ {
		count := sl._blockCount[i]
		if index < count {
			num := i * sl._blockSize
			for {
				numCount := sl._counter[num]
				if index < numCount {
					return num
				}
				index -= numCount
				num++
			}
		}
		index -= count
	}
	panic("unreachable")
}

// 严格小于 value 的元素个数.
// 也即第一个大于等于 value 的元素的下标.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) BisectLeft(value int) int {
	pos := sl._belong[value]
	res := 0
	for i := 0; i < pos; i++ {
		res += sl._blockCount[i]
	}
	for v := pos * sl._blockSize; v < value; v++ {
		res += sl._counter[v]
	}
	return res
}

// 小于等于 value 的元素个数.
// 也即第一个大于 value 的元素的下标.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) BisectRight(value int) int {
	return sl.BisectLeft(value + 1)
}

func (sl *SortedListRangeBlock) Count(value int) int {
	return sl._counter[value]
}

// 返回范围 `[min, max]` 内数的个数.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) CountRange(min, max int) int {
	if min > max {
		return 0
	}

	minPos := sl._belong[min]
	maxPos := sl._belong[max]
	if minPos == maxPos {
		res := 0
		for i := min; i <= max; i++ {
			res += sl._counter[i]
		}
		return res
	}

	res := 0
	minEnd := (minPos + 1) * sl._blockSize
	for v := min; v < minEnd; v++ {
		res += sl._counter[v]
	}
	for i := minPos + 1; i < maxPos; i++ {
		res += sl._blockCount[i]
	}
	maxStart := maxPos * sl._blockSize
	for v := maxStart; v <= max; v++ {
		res += sl._counter[v]
	}
	return res
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Lower(value int) (res int, ok bool) {
	pos := sl._belong[value]
	start := pos * sl._blockSize
	for v := value - 1; v >= start; v-- {
		if sl._counter[v] > 0 {
			return v, true
		}
	}

	for i := pos - 1; i >= 0; i-- {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := (i + 1) * sl._blockSize
		for {
			if sl._counter[num] > 0 {
				return num, true
			}
			num--
		}
	}

	return
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Higher(value int) (res int, ok bool) {
	pos := sl._belong[value]
	end := (pos + 1) * sl._blockSize
	for v := value + 1; v < end; v++ {
		if sl._counter[v] > 0 {
			return v, true
		}
	}

	for i := pos + 1; i < len(sl._blockCount); i++ {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := i * sl._blockSize
		for {
			if sl._counter[num] > 0 {
				return num, true
			}
			num++
		}
	}

	return
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Floor(value int) (res int, ok bool) {
	if sl.Has(value) {
		return value, true
	}
	return sl.Lower(value)
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Ceiling(value int) (res int, ok bool) {
	if sl.Has(value) {
		return value, true
	}
	return sl.Higher(value)
}

// 返回区间 `[start, end)` 的和.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) SumSlice(start, end int) int {
	if start < 0 {
		start += sl._len
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl._len
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return 0
	}

	res := 0
	remain := end - start
	cur, index := sl._findKth(start)
	sufCount := sl._counter[cur] - index
	if sufCount >= remain {
		return remain * cur
	}

	res += sufCount * cur
	remain -= sufCount
	cur++

	// 当前块内的和
	blockEnd := (sl._belong[cur] + 1) * sl._blockSize
	for remain > 0 && cur < blockEnd {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		res += real * cur
		remain -= real
		cur++
	}

	// 以块为单位消耗remain
	pos := sl._belong[cur]
	for pos < len(sl._blockCount) && remain >= sl._blockCount[pos] {
		res += sl._blockSum[pos]
		remain -= sl._blockCount[pos]
		pos++
		cur += sl._blockSize
	}

	// 剩余的
	for remain > 0 {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		res += real * cur
		remain -= real
		cur++
	}

	return res
}

// 返回范围 `[min, max]` 的和.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) SumRange(min, max int) int {
	minPos := sl._belong[min]
	maxPos := sl._belong[max]
	if minPos == maxPos {
		res := 0
		for i := min; i <= max; i++ {
			res += sl._counter[i] * i
		}
		return res
	}

	res := 0
	minEnd := (minPos + 1) * sl._blockSize
	for v := min; v < minEnd; v++ {
		res += sl._counter[v] * v
	}
	for i := minPos + 1; i < maxPos; i++ {
		res += sl._blockSum[i]
	}
	maxStart := maxPos * sl._blockSize
	for v := maxStart; v <= max; v++ {
		res += sl._counter[v] * v
	}
	return res
}

func (sl *SortedListRangeBlock) ForEach(f func(value, index int), reverse bool) {
	if reverse {
		ptr := 0
		for i := len(sl._counter) - 1; i >= 0; i-- {
			count := sl._counter[i]
			for j := 0; j < count; j++ {
				f(i, ptr)
				ptr++
			}
		}
	} else {
		ptr := 0
		for i := 0; i < len(sl._counter); i++ {
			count := sl._counter[i]
			for j := 0; j < count; j++ {
				f(i, ptr)
				ptr++
			}
		}
	}
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Pop(index int) int {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	value := sl.At(index)
	sl.Remove(value)
	return value
}

func (sl *SortedListRangeBlock) Slice(start, end int) []int {
	if start < 0 {
		start += sl._len
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl._len
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return nil
	}

	res := make([]int, end-start)
	count := 0
	sl.Enumerate(start, end, func(value int) {
		res[count] = value
		count++
	}, false)

	return res
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Erase(start, end int) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedListRangeBlock) Enumerate(start, end int, f func(value int), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return
	}

	remain := end - start
	cur, index := sl._findKth(start)
	sufCount := sl._counter[cur] - index
	real := sufCount
	if real > remain {
		real = remain
	}
	if f != nil {
		for i := 0; i < real; i++ {
			f(cur)
		}
	}
	if erase {
		for i := 0; i < real; i++ {
			sl.Remove(cur)
		}
	}
	remain -= sufCount
	if remain == 0 {
		return
	}
	cur++

	// 当前块内
	blockEnd := (sl._belong[cur] + 1) * sl._blockSize
	for remain > 0 && cur < blockEnd {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		remain -= real
		if f != nil {
			for i := 0; i < real; i++ {
				f(cur)
			}
		}
		if erase {
			for i := 0; i < real; i++ {
				sl.Remove(cur)
			}
		}
		cur++
	}

	// 以块为单位消耗remain
	pos := sl._belong[cur]
	for count := sl._blockCount[pos]; remain >= count; {
		remain -= count
		if f != nil {
			for v := cur; v < cur+sl._blockSize; v++ {
				c := sl._counter[v]
				for i := 0; i < c; i++ {
					f(v)
				}
			}
		}
		if erase {
			for v := cur; v < cur+sl._blockSize; v++ {
				sl._counter[v] = 0
			}
			sl._len -= count
			sl._blockCount[pos] = 0
			sl._blockSum[pos] = 0
		}
		pos++
		cur += sl._blockSize
	}

	// 剩余的
	for remain > 0 {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		remain -= real
		if f != nil {
			for i := 0; i < real; i++ {
				f(cur)
			}
		}
		if erase {
			for i := 0; i < real; i++ {
				sl.Remove(cur)
			}
		}
		cur++
	}
}

func (sl *SortedListRangeBlock) Clear() {
	for i := range sl._counter {
		sl._counter[i] = 0
	}
	for i := range sl._blockCount {
		sl._blockCount[i] = 0
	}
	for i := range sl._blockSum {
		sl._blockSum[i] = 0
	}
	sl._len = 0
}

func (sl *SortedListRangeBlock) Update(values ...int) {
	for _, value := range values {
		sl.Add(value)
	}
}

func (sl *SortedListRangeBlock) Merge(other *SortedListRangeBlock) {
	other.ForEach(func(value, _ int) {
		sl.Add(value)
	}, false)
}

func (sl *SortedListRangeBlock) String() string {
	sb := make([]string, 0, sl._len)
	sl.ForEach(func(value, _ int) {
		sb = append(sb, fmt.Sprintf("%d", value))
	}, false)
	return fmt.Sprintf("SortedListRangeBlock{%s}", strings.Join(sb, ", "))
}

func (sl *SortedListRangeBlock) Len() int {
	return sl._len
}

func (sl *SortedListRangeBlock) Min() int {
	return sl.At(0)
}

func (sl *SortedListRangeBlock) Max() int {
	if sl._len == 0 {
		panic("empty")
	}

	for i := len(sl._blockCount) - 1; i >= 0; i-- {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := (i+1)*sl._blockSize - 1
		for {
			if sl._counter[num] > 0 {
				return num
			}
			num--
		}
	}

	panic("unreachable")
}

// 返回索引在`kth`处的元素的`value`,以及该元素是`value`中的第几个(`index`).
func (sl *SortedListRangeBlock) _findKth(kth int) (value, index int) {
	for i := 0; i < len(sl._blockCount); i++ {
		count := sl._blockCount[i]
		if kth < count {
			num := i * sl._blockSize
			for {
				numCount := sl._counter[num]
				if kth < numCount {
					return num, kth
				}
				kth -= numCount
				num++
			}
		}
		kth -= count
	}

	panic("unreachable")
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
