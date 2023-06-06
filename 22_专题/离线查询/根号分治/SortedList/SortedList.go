// !速度比Treap实现的名次树慢, 谨慎使用

package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// https://leetcode.cn/problems/find-score-of-an-array-after-marking-all-elements/submissions/
func findScore(nums []int) int64 {
	pairs := make([][2]int, len(nums))
	for i, v := range nums {
		pairs[i] = [2]int{v, i}
	}
	sl := _NSL(func(a, b S) bool {
		if a[0] != b[0] {
			return a[0] < b[0]
		}
		return a[1] < b[1]
	}, pairs)
	res := 0
	for sl.Len() > 0 {
		a := sl.Pop(0)
		v, i := a[0], a[1]
		res += v
		if i-1 >= 0 {
			sl.Discard([2]int{nums[i-1], i - 1})
		}
		if i+1 < len(nums) {
			sl.Discard([2]int{nums[i+1], i + 1})
		}
	}
	return int64(res)
}

func main() {
	sl := _NSL(func(a, b S) bool { return a[0] < b[0] }, nil)
	sl.Add([2]int{2, 3})
	sl.Add([2]int{1, 2})
	fmt.Println(sl.Pop(0))
	sl.Erase(0, 2)
	fmt.Println(sl, sl.Len())
}

type S = [2]int

// SortedListSQRT
type _SL struct {
	less   func(a, b S) bool
	size   int
	blocks [][]S
}

func _NSL(less func(a, b S) bool, items []S) *_SL {
	res := &_SL{less: less}
	if len(items) > 0 {
		items = append(items[:0:0], items...)
		sort.Slice(items, func(i, j int) bool {
			return less(items[i], items[j])
		})
		res.blocks = res._initBlocks(items)
	}
	res.size = len(items)
	return res
}

const (
	_BLOCK_RATIO   = 25
	_REBUILD_RATIO = 70
)

func (sl *_SL) Add(value S) {
	if sl.size == 0 {
		sl.blocks = append(sl.blocks[:0], []S{value})
		sl.size = 1
		return
	}

	hitIndex := sl._findBlockIndex(value)
	if hitIndex == -1 {
		sl.blocks[len(sl.blocks)-1] = append(sl.blocks[len(sl.blocks)-1], value)
		sl.size++
		if len(sl.blocks[len(sl.blocks)-1]) > _REBUILD_RATIO*len(sl.blocks) {
			sl._rebuild()
		}
		return
	}

	hitted := sl.blocks[hitIndex]
	pos := sl._bisectRight(hitted, value)
	sl.blocks[hitIndex] = append(hitted[:pos], append([]S{value}, hitted[pos:]...)...)
	sl.size++
	if len(hitted) > _REBUILD_RATIO*len(sl.blocks) {
		sl._rebuild()
	}
}

func (sl *_SL) Has(value S) bool {
	if sl.size == 0 {
		return false
	}
	hitIndex := sl._findBlockIndex(value)
	if hitIndex == -1 {
		return false
	}
	hitted := sl.blocks[hitIndex]
	pos := sl._bisectLeft(hitted, value)
	return pos < len(hitted) && hitted[pos] == value
}

func (sl *_SL) Discard(value S) bool {
	if sl.size == 0 {
		return false
	}
	hitIndex := sl._findBlockIndex(value)
	if hitIndex == -1 {
		return false
	}
	hitted := sl.blocks[hitIndex]
	pos := sl._bisectLeft(hitted, value)
	if pos == len(hitted) || hitted[pos] != value {
		return false
	}
	sl.blocks[hitIndex] = append(hitted[:pos], hitted[pos+1:]...)
	sl.size--
	if len(sl.blocks[hitIndex]) == 0 {
		// !Splice When Empty, Do Not Rebuild
		sl.blocks = append(sl.blocks[:hitIndex], sl.blocks[hitIndex+1:]...)
	}
	return true
}

func (sl *_SL) Pop(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	for i, block := range sl.blocks {
		if index < len(block) {
			res := block[index]
			sl.blocks[i] = append(block[:index], block[index+1:]...)
			sl.size--
			if len(sl.blocks[i]) == 0 {
				// !Splice When Empty, Do Not Rebuild
				sl.blocks = append(sl.blocks[:i], sl.blocks[i+1:]...)
			}
			return res
		}
		index -= len(block)
	}
	panic("impossible")
}

// 删除区间 [start, end) 内的元素.
func (sl *_SL) Erase(start, end int) {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return
	}

	bid, startPos := sl._moveTo(start)
	deleteCount := end - start
	for ; bid < len(sl.blocks) && deleteCount > 0; bid++ {
		block := sl.blocks[bid]
		endPos := min(len(block), startPos+deleteCount)
		curDeleteCount := endPos - startPos
		if curDeleteCount == len(block) {
			sl.blocks = append(sl.blocks[:bid], sl.blocks[bid+1:]...)
			bid--
		} else {
			sl.blocks[bid] = append(block[:startPos], block[endPos:]...)
		}
		deleteCount -= curDeleteCount
		sl.size -= curDeleteCount
		startPos = 0
	}
}

func (sl *_SL) At(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	for _, block := range sl.blocks {
		if index < len(block) {
			return block[index]
		}
		index -= len(block)
	}
	panic("impossible")
}

// Count the number of elements < value or
// returns the index of the first element >= value.
func (sl *_SL) BisectLeft(value S) int {
	res := 0
	for _, block := range sl.blocks {
		if !sl.less(block[len(block)-1], value) {
			return res + sl._bisectLeft(block, value)
		}
		res += len(block)
	}
	return res
}

// Count the number of elements <= value or
// returns the index of the first element > value.
func (sl *_SL) BisectRight(value S) int {
	res := 0
	for _, block := range sl.blocks {
		if sl.less(value, block[len(block)-1]) {
			return res + sl._bisectRight(block, value)
		}
		res += len(block)
	}
	return res
}

func (sl *_SL) Clear() {
	sl.blocks = sl.blocks[:0]
	sl.size = 0
}

func (sl *_SL) Lower(value S) (res S, ok bool) {
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		if sl.less(block[0], value) {
			pos := sl._bisectLeft(block, value)
			return block[pos-1], true
		}
	}
	return
}

func (sl *_SL) Higher(value S) (res S, ok bool) {
	for _, block := range sl.blocks {
		if sl.less(value, block[len(block)-1]) {
			pos := sl._bisectRight(block, value)
			return block[pos], true
		}
	}
	return
}

func (sl *_SL) Floor(value S) (res S, ok bool) {
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		if !sl.less(value, block[0]) {
			pos := sl._bisectRight(block, value)
			return block[pos-1], true
		}
	}
	return
}

func (sl *_SL) Ceiling(value S) (res S, ok bool) {
	for _, block := range sl.blocks {
		if !sl.less(block[len(block)-1], value) {
			pos := sl._bisectLeft(block, value)
			return block[pos], true
		}
	}
	return
}

func (sl *_SL) ForEach(f func(value S, index int)) {
	pos := 0
	for _, block := range sl.blocks {
		for _, value := range block {
			f(value, pos)
			pos++
		}
	}
}

func (sl *_SL) Len() int {
	return sl.size
}

func (sl *_SL) String() string {
	res := make([]string, 0)
	sl.ForEach(func(value S, _ int) {
		res = append(res, fmt.Sprintf("%v", value))
	})
	return fmt.Sprintf("SortedList{%v}", strings.Join(res, ", "))
}

func (sl *_SL) Min() (res S, ok bool) {
	if len(sl.blocks) == 0 {
		return
	}
	return sl.blocks[0][0], true
}

func (sl *_SL) Max() (res S, ok bool) {
	if len(sl.blocks) == 0 {
		return
	}
	last := sl.blocks[len(sl.blocks)-1]
	return last[len(last)-1], true
}

func (sl *_SL) _rebuild() {
	if sl.size == 0 {
		return
	}
	bc := int(math.Ceil(math.Sqrt(float64(sl.size) / _BLOCK_RATIO)))
	bs := (sl.size + bc - 1) / bc
	newB := make([][]S, bc)
	ptr := 0
	for i := 0; i < len(sl.blocks); i++ {
		b := sl.blocks[i]
		for j := 0; j < len(b); j++ {
			tmp := ptr / bs
			newB[tmp] = append(newB[tmp], b[j])
			ptr++
		}
	}
	sl.blocks = newB
}

func (sl *_SL) _initBlocks(sorted []S) [][]S {
	bc := int(math.Ceil(math.Sqrt(float64(len(sorted)) / _BLOCK_RATIO)))
	bs := (len(sorted) + bc - 1) / bc
	res := make([][]S, bc)
	for i := 0; i < bc; i++ {
		res[i] = append(res[i], sorted[i*bs:min((i+1)*bs, len(sorted))]...)
	}
	return res
}

func (sl *_SL) _bisectLeft(nums []S, value S) int {
	return sort.Search(len(nums), func(i int) bool {
		return !sl.less(nums[i], value)
	})
}

func (sl *_SL) _bisectRight(nums []S, value S) int {
	return sort.Search(len(nums), func(i int) bool {
		return sl.less(value, nums[i])
	})
}

// 如果没有找到,返回-1
func (sl *_SL) _findBlockIndex(x S) int {
	for i, block := range sl.blocks {
		if !sl.less(block[len(block)-1], x) {
			return i
		}
	}
	return -1
}

func (sl *_SL) _moveTo(index int) (blockId, startPos int) {
	for i, block := range sl.blocks {
		if index < len(block) {
			return i, index
		}
		index -= len(block)
	}
	return len(sl.blocks), 0
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
