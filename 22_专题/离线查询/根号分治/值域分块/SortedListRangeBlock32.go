// SortedListRangeBlock
// 值域分块，O(1)插入、删除，O(sqrt(n))查询的`SortedList`.
// 一般配合莫队算法使用.

package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	sl := NewSortedListRangeBlock32(10)
	sl.Add(1)
	fmt.Println(sl.Min())
	fmt.Println(sl.Max())
}

// 2163. 删除元素后和的最小差值
// https://leetcode.cn/problems/minimum-difference-in-sums-after-removal-of-elements/
func minimumDifference(nums []int) int64 {
	n := int32(len(nums) / 3)
	pre := NewSortedListRangeBlock32(1e5 + 10)
	for _, v := range nums[:n] {
		pre.Add(int32(v))
	}
	suf := NewSortedListRangeBlock32(1e5 + 10)
	for _, v := range nums[n:] {
		suf.Add(int32(v))
	}
	res := pre.SumSlice(0, n) - suf.SumSlice(suf.Len()-n, suf.Len())
	for i := n; i < 2*n; i++ {
		pre.Add(int32(nums[i]))
		suf.Remove(int32(nums[i]))
		res = min(res, pre.SumSlice(0, n)-suf.SumSlice(suf.Len()-n, suf.Len()))
	}
	return int64(res)
}

const INF int = 1e18

type SortedListRangeBlock32 struct {
	_blockSize  int32   // 每个块的大小.
	_len        int32   // 所有数的个数.
	_counter    []int32 // 每个数出现的次数.
	_blockCount []int32 // 每个块的数的个数.
	_belong     []int32 // 每个数所在的块.
	_blockSum   []int   // 每个块的和.
}

// 值域分块模拟SortedList.
// `O(1)`add/remove，`O(sqrt(n))`查询.
// 一般配合莫队算法使用.
//
//	max:值域的最大值.0 <= max <= 1e6.
//	iterable:初始值.
func NewSortedListRangeBlock32(max int32, nums ...int32) *SortedListRangeBlock32 {
	max += 5
	size := int32(math.Sqrt(float64(max)))
	count := 1 + (max / size)
	sl := &SortedListRangeBlock32{
		_blockSize:  size,
		_counter:    make([]int32, max),
		_blockCount: make([]int32, count),
		_belong:     make([]int32, max),
		_blockSum:   make([]int, count),
	}
	for i := int32(0); i < max; i++ {
		sl._belong[i] = i / size
	}
	if len(nums) > 0 {
		sl.Update(nums...)
	}
	return sl
}

// O(1).
func (sl *SortedListRangeBlock32) Add(value int32) {
	sl._counter[value]++
	pos := sl._belong[value]
	sl._blockCount[pos]++
	sl._blockSum[pos] += int(value)
	sl._len++
}

// O(1).
func (sl *SortedListRangeBlock32) Remove(value int32) {
	sl._counter[value]--
	pos := sl._belong[value]
	sl._blockCount[pos]--
	sl._blockSum[pos] -= int(value)
	sl._len--
}

// O(1).
func (sl *SortedListRangeBlock32) Discard(value int32) bool {
	if !sl.Has(value) {
		return false
	}
	sl.Remove(value)
	return true
}

// O(1).
func (sl *SortedListRangeBlock32) Has(value int32) bool {
	return sl._counter[value] > 0
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) At(index int32) int32 {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	for i := int32(0); i < int32(len(sl._blockCount)); i++ {
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
func (sl *SortedListRangeBlock32) BisectLeft(value int32) int32 {
	pos := sl._belong[value]
	res := int32(0)
	for i := int32(0); i < pos; i++ {
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
func (sl *SortedListRangeBlock32) BisectRight(value int32) int32 {
	return sl.BisectLeft(value + 1)
}

func (sl *SortedListRangeBlock32) Count(value int32) int32 {
	return sl._counter[value]
}

// 返回范围 `[min, max]` 内数的个数.
// O(sqrt(n)).
func (sl *SortedListRangeBlock32) CountRange(min, max int32) int32 {
	if min > max {
		return 0
	}

	minPos := sl._belong[min]
	maxPos := sl._belong[max]
	if minPos == maxPos {
		res := int32(0)
		for i := min; i <= max; i++ {
			res += sl._counter[i]
		}
		return res
	}

	res := int32(0)
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
func (sl *SortedListRangeBlock32) Lower(value int32) (res int32, ok bool) {
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
func (sl *SortedListRangeBlock32) Higher(value int32) (res int32, ok bool) {
	pos := sl._belong[value]
	end := (pos + 1) * sl._blockSize
	for v := value + 1; v < end; v++ {
		if sl._counter[v] > 0 {
			return v, true
		}
	}

	for i := pos + 1; i < int32(len(sl._blockCount)); i++ {
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
func (sl *SortedListRangeBlock32) Floor(value int32) (res int32, ok bool) {
	if sl.Has(value) {
		return value, true
	}
	return sl.Lower(value)
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Ceiling(value int32) (res int32, ok bool) {
	if sl.Has(value) {
		return value, true
	}
	return sl.Higher(value)
}

// 返回区间 `[start, end)` 的和.
// O(sqrt(n)).
func (sl *SortedListRangeBlock32) SumSlice(start, end int32) int {
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
		return int(remain) * int(cur)
	}

	res += int(sufCount) * int(cur)
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
		res += int(real) * int(cur)
		remain -= real
		cur++
	}

	// 以块为单位消耗remain
	pos := sl._belong[cur]
	for pos < int32(len(sl._blockCount)) && remain >= sl._blockCount[pos] {
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
		res += int(real) * int(cur)
		remain -= real
		cur++
	}

	return res
}

// 返回范围 `[min, max]` 的和.
// O(sqrt(n)).
func (sl *SortedListRangeBlock32) SumRange(min, max int32) int {
	minPos := sl._belong[min]
	maxPos := sl._belong[max]
	if minPos == maxPos {
		res := 0
		for i := min; i <= max; i++ {
			res += int(sl._counter[i]) * int(i)
		}
		return res
	}

	res := 0
	minEnd := (minPos + 1) * sl._blockSize
	for v := min; v < minEnd; v++ {
		res += int(sl._counter[v]) * int(v)
	}
	for i := minPos + 1; i < maxPos; i++ {
		res += sl._blockSum[i]
	}
	maxStart := maxPos * sl._blockSize
	for v := maxStart; v <= max; v++ {
		res += int(sl._counter[v]) * int(v)
	}
	return res
}

func (sl *SortedListRangeBlock32) ForEach(f func(value, index int32), reverse bool) {
	if reverse {
		ptr := int32(0)
		for i := int32(len(sl._counter) - 1); i >= 0; i-- {
			count := sl._counter[i]
			for j := int32(0); j < count; j++ {
				f(i, ptr)
				ptr++
			}
		}
	} else {
		ptr := int32(0)
		for i := int32(0); i < int32(len(sl._counter)); i++ {
			count := sl._counter[i]
			for j := int32(0); j < count; j++ {
				f(i, ptr)
				ptr++
			}
		}
	}
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Pop(index int32) int32 {
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

func (sl *SortedListRangeBlock32) Slice(start, end int32) []int32 {
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

	res := make([]int32, end-start)
	count := int32(0)
	sl.Enumerate(start, end, func(value int32) {
		res[count] = value
		count++
	}, false)

	return res
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Erase(start, end int32) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedListRangeBlock32) Enumerate(start, end int32, f func(value int32), erase bool) {
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
		for i := int32(0); i < real; i++ {
			f(cur)
		}
	}
	if erase {
		for i := int32(0); i < real; i++ {
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
			for i := int32(0); i < real; i++ {
				f(cur)
			}
		}
		if erase {
			for i := int32(0); i < real; i++ {
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
				for i := int32(0); i < c; i++ {
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
			for i := int32(0); i < real; i++ {
				f(cur)
			}
		}
		if erase {
			for i := int32(0); i < real; i++ {
				sl.Remove(cur)
			}
		}
		cur++
	}
}

func (sl *SortedListRangeBlock32) Clear() {
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

func (sl *SortedListRangeBlock32) Update(values ...int32) {
	for _, value := range values {
		sl.Add(value)
	}
}

func (sl *SortedListRangeBlock32) Merge(other *SortedListRangeBlock32) {
	other.ForEach(func(value, _ int32) {
		sl.Add(value)
	}, false)
}

func (sl *SortedListRangeBlock32) String() string {
	sb := make([]string, 0, sl._len)
	sl.ForEach(func(value, _ int32) {
		sb = append(sb, fmt.Sprintf("%d", value))
	}, false)
	return fmt.Sprintf("SortedListRangeBlock{%s}", strings.Join(sb, ", "))
}

func (sl *SortedListRangeBlock32) Len() int32 {
	return sl._len
}

func (sl *SortedListRangeBlock32) Min() int32 {
	return sl.At(0)
}

func (sl *SortedListRangeBlock32) Max() int32 {
	if sl._len == 0 {
		panic("empty")
	}

	for i := int32(len(sl._blockCount) - 1); i >= 0; i-- {
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
func (sl *SortedListRangeBlock32) _findKth(kth int32) (value, index int32) {
	for i := int32(0); i < int32(len(sl._blockCount)); i++ {
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
