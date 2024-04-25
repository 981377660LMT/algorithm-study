// api:
//  1.Insert(index int32, v int8)
//  2.Pop(index int32) int8
//  3.Set(index int32, v int8)
//  4.Get(index int32) int8
//  5.Count0(end int32) int32
//  6.Count1(end int32) int32
//  7.Count(end int32, v int8) int32
//  8.Kth0(k int32) int32
//  9.Kth1(k int32) int32
// 10.Kth(k int32, v int8) int32
// 11.Len() int32
// 12.ToList() []int8
// 13.Debug()

// TODO:SQRTARRAY 需要维护区间和

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	demo()
}

func demo() {
	bv := NewSortedList(10, func(i int32) int8 { return 0 })
	for i := int32(0); i < 10; i++ {
		bv.Insert(i, 1)
	}
	bv.Set(3, 0)
	bv.Set(8, 1)
	bv.Insert(3, 1)
	bv.Pop(0)
	bv.Pop(0)
	bv.Pop(0)
	bv.Pop(0)
	fmt.Println(bv.GetAll())
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int32 = 200

// 使用分块+树状数组维护的动态Bitvector.
type DynamicBitvector struct {
	size              int32
	blocks            [][]int8
	blockOnes         []int16
	preLen            []int32 // bit1
	preOnes           []int32 // bit2
	shouldRebuildTree bool
}

func NewSortedList(n int32, f func(i int32) int8) *DynamicBitvector {
	blocks := make([][]int8, 0, n/_LOAD+1)
	blockOnes := make([]int16, 0, n/_LOAD+1)
	for start := int32(0); start < n; start += _LOAD {
		end := start + _LOAD
		if end > n {
			end = n
		}
		block := make([]int8, end-start)
		ones := int16(0)
		for i := start; i < end; i++ {
			block[i-start] = f(i)
			ones += int16(block[i-start])
		}
		blocks = append(blocks, block)
		blockOnes = append(blockOnes, ones)
	}
	res := &DynamicBitvector{
		size:              n,
		blocks:            blocks,
		blockOnes:         blockOnes,
		shouldRebuildTree: true,
	}
	return res
}

func (sl *DynamicBitvector) Insert(index int32, value int8) {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []int8{value})
		sl.blockOnes = append(sl.blockOnes, int16(value))
		sl.shouldRebuildTree = true
		return
	}

	pos, startIndex := sl._findKth(index)

	if value == 1 {
		sl._updatePreLenAndPreOnes(pos, true)
		sl.blockOnes[pos]++
	} else {
		sl._updatePreLen(pos, true)
	}
	// TODO:TIME
	sl.blocks[pos] = append(sl.blocks[pos][:startIndex], append([]int8{value}, sl.blocks[pos][startIndex:]...)...)

	// n -> load + (n - load)
	// !append 用copy代替
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		totalOnes := sl.blockOnes[pos]
		sl.blocks = append(sl.blocks[:pos+1], append([][]int8{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		ones1 := int16(0)
		for _, v := range sl.blocks[pos] {
			ones1 += int16(v)
		}
		ones2 := totalOnes - ones1
		sl.blockOnes = append(sl.blockOnes, 0)
		// TODO:TIME
		copy(sl.blockOnes[pos+2:], sl.blockOnes[pos+1:])
		sl.blockOnes[pos] = ones1
		sl.blockOnes[pos+1] = ones2
		sl.shouldRebuildTree = true
	}

	return
}

func (sl *DynamicBitvector) Pop(index int32) int8 {
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.size--
	if value == 1 {
		sl._updatePreLenAndPreOnes(pos, false)
		sl.blockOnes[pos]--
	} else {
		sl._updatePreLen(pos, false)
	}
	// TODO:TIME
	sl.blocks[pos] = append(sl.blocks[pos][:index], sl.blocks[pos][index+1:]...)
	if len(sl.blocks[pos]) == 0 {
		// !delete block
		sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
		sl.shouldRebuildTree = true
	}
	return value
}

func (sl *DynamicBitvector) Get(index int32) int8 {
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *DynamicBitvector) Set(index int32, value int8) {
	pos, startIndex := sl._findKth(index)
	curValue := sl.blocks[pos][startIndex]
	if curValue == value {
		return
	}
	sl.blocks[pos][startIndex] = value
	if value == 1 {
		sl._updatePreOnes(pos, true)
		sl.blockOnes[pos]++
	} else {
		sl._updatePreOnes(pos, false)
		sl.blockOnes[pos]--
	}
}

func (sl *DynamicBitvector) Count0(end int32) int32 {

	return 0
}

func (sl *DynamicBitvector) Count1(end int32) int32 {
	return 0
}

func (sl *DynamicBitvector) Count(end int32, value int8) int32 {
	if value == 1 {
		return sl._queryPreOnes(end)
	}
	return end - sl._queryPreOnes(end)
}

func (sl *DynamicBitvector) Kth0(k int32) int32 {
	return 0
}

func (sl *DynamicBitvector) Kth1(k int32) int32 {
	return 0
}

func (sl *DynamicBitvector) Kth(k int32, value int8) int32 {
	if value == 1 {
		return sl.Kth1(k)
	}
	return sl.Kth0(k)
}

func (sl *DynamicBitvector) Len() int32 {
	return sl.size
}

func (sl *DynamicBitvector) GetAll() []int8 {
	res := make([]int8, 0, sl.size)
	for _, block := range sl.blocks {
		res = append(res, block...)
	}
	return res
}

func (sl *DynamicBitvector) _buildPreLenAndPreOnes() {
	sl.preLen = make([]int32, len(sl.blocks))
	sl.preOnes = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.preLen[i] = int32(len(sl.blocks[i]))
		sl.preOnes[i] = int32(sl.blockOnes[i])
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	m := int32(len(tree1))
	for i := int32(0); i < m; i++ {
		j := i | (i + 1)
		if j < m {
			tree1[j] += tree1[i]
			tree2[j] += tree2[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *DynamicBitvector) _updatePreLen(index int32, addOne bool) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.preLen
	m := int32(len(tree))
	if addOne {
		for i := index; i < m; i |= i + 1 {
			tree[i]++
		}
	} else {
		for i := index; i < m; i |= i + 1 {
			tree[i]--
		}
	}
}

func (sl *DynamicBitvector) _updatePreOnes(index int32, addOne bool) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.preOnes
	m := int32(len(tree))
	if addOne {
		for i := index; i < m; i |= i + 1 {
			tree[i]++
		}
	} else {
		for i := index; i < m; i |= i + 1 {
			tree[i]--
		}
	}
}

func (sl *DynamicBitvector) _updatePreLenAndPreOnes(index int32, addOne bool) {
	if sl.shouldRebuildTree {
		return
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	m := int32(len(tree1))
	if addOne {
		for i := index; i < m; i |= i + 1 {
			tree1[i]++
			tree2[i]++
		}
	} else {
		for i := index; i < m; i |= i + 1 {
			tree1[i]--
			tree2[i]--
		}
	}
}

func (sl *DynamicBitvector) _queryPreLen(end int32) int32 {
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree := sl.preLen
	sum := int32(0)
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *DynamicBitvector) _queryPreOnes(end int32) int32 {
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree := sl.preOnes
	sum := int32(0)
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *DynamicBitvector) _queryPreLenAndPreOnes(end int32) (lens, ones int32) {
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	for end > 0 {
		lens += tree1[end-1]
		ones += tree2[end-1]
		end &= end - 1
	}
	return
}

func (sl *DynamicBitvector) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree := sl.preLen
	pos = -1
	m := int32(len(tree))
	bitLen := bits.Len32(uint32(m))
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree[next] {
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
