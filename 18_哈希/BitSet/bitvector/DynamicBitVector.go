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

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"
)

func main() {
	// demo()
	test()
	testTime()
}

func demo() {
	bv := NewDynamicBitvector(10, func(i int32) int8 { return 0 })
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
	fmt.Println(bv.Count1(20))
	fmt.Println(bv.Count0(20))
	fmt.Println(bv.Kth0(0))
	fmt.Println(bv.Kth0(1))
	fmt.Println(bv.Kth0(2))
	fmt.Println(bv.Kth0(3))
	fmt.Println(bv.Kth0(4))
	fmt.Println(bv.Kth0(5))
	fmt.Println(bv.Kth0(6))
	fmt.Println(bv.Kth0(7))
	fmt.Println(bv.Kth0(8))
}

// 1e5 -> 150 , 2e5 -> 250
const _LOAD int32 = 250 // 块尺寸越大，修改越快，查询越慢

// 使用分块+树状数组维护的动态Bitvector.
type DynamicBitvector struct {
	size              int32
	totalOnes         int32
	blocks            [][]int8
	blockOnes         []int32
	preLen            []int32 // 每个块块长的前缀和
	preOnes           []int32 // 每个块01的前缀和
	shouldRebuildTree bool
}

func NewDynamicBitvector(n int32, f func(i int32) int8) *DynamicBitvector {
	blocks := make([][]int8, 0, n/_LOAD+1)
	blockOnes := make([]int32, 0, n/_LOAD+1)
	totalOnes := int32(0)
	for start := int32(0); start < n; start += _LOAD {
		end := start + _LOAD
		if end > n {
			end = n
		}
		block := make([]int8, end-start)
		ones := int32(0)
		for i := start; i < end; i++ {
			block[i-start] = f(i)
			ones += int32(block[i-start])
		}
		blocks = append(blocks, block)
		blockOnes = append(blockOnes, ones)
		totalOnes += ones
	}
	res := &DynamicBitvector{
		size:              n,
		totalOnes:         totalOnes,
		blocks:            blocks,
		blockOnes:         blockOnes,
		shouldRebuildTree: true,
	}
	return res
}

func (sl *DynamicBitvector) Insert(index int32, value int8) {
	if value == 1 {
		sl.totalOnes++
	}
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []int8{value})
		sl.blockOnes = append(sl.blockOnes, int32(value))
		sl.shouldRebuildTree = true
		sl.size++
		return
	}

	pos, startIndex := sl._findKth(index)
	if value == 1 {
		sl._updatePreLenAndPreOnes(pos, true)
		sl.blockOnes[pos]++
	} else {
		sl._updatePreLen(pos, true)
	}
	sl.blocks[pos] = append(sl.blocks[pos], 0)
	copy(sl.blocks[pos][startIndex+1:], sl.blocks[pos][startIndex:])
	sl.blocks[pos][startIndex] = value

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		totalOnes := sl.blockOnes[pos]
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:] // !注意max的设置(为了让左右互不影响)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		ones1 := int32(0)
		for _, v := range sl.blocks[pos] {
			ones1 += int32(v)
		}
		ones2 := totalOnes - ones1
		sl.blockOnes = append(sl.blockOnes, 0)
		copy(sl.blockOnes[pos+2:], sl.blockOnes[pos+1:])
		sl.blockOnes[pos] = ones1
		sl.blockOnes[pos+1] = ones2
		sl.shouldRebuildTree = true
	}

	sl.size++
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
		sl.totalOnes--
	} else {
		sl._updatePreLen(pos, false)
	}

	copy(sl.blocks[pos][startIndex:], sl.blocks[pos][startIndex+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]

	if len(sl.blocks[pos]) == 0 {
		// !delete block
		copy(sl.blocks[pos:], sl.blocks[pos+1:])
		sl.blocks = sl.blocks[:len(sl.blocks)-1]
		copy(sl.blockOnes[pos:], sl.blockOnes[pos+1:])
		sl.blockOnes = sl.blockOnes[:len(sl.blockOnes)-1]
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
		sl.totalOnes++
	} else {
		sl._updatePreOnes(pos, false)
		sl.blockOnes[pos]--
		sl.totalOnes--
	}
}

func (sl *DynamicBitvector) Count0(end int32) int32 {
	if end <= 0 {
		return 0
	}
	if end > sl.size {
		end = sl.size
	}
	return end - sl.Count1(end)
}

func (sl *DynamicBitvector) Count1(end int32) int32 {
	if end <= 0 {
		return 0
	}
	if end > sl.size {
		end = sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	pos, startIndex, preOnes := sl._findKthAndPreOnes(end - 1)
	block := sl.blocks[pos]
	if startIndex < int32(len(block))>>1 {
		// 统计前缀
		onesInBlockPrefix := int16(0)
		for i := int32(0); i <= startIndex; i++ {
			onesInBlockPrefix += int16(block[i])
		}
		return preOnes + int32(onesInBlockPrefix)
	} else {
		onesInBlockSuffix := int16(0)
		for i := int32(len(block) - 1); i > startIndex; i-- {
			onesInBlockSuffix += int16(block[i])
		}
		return preOnes + sl.blockOnes[pos] - int32(onesInBlockSuffix)
	}
}

func (sl *DynamicBitvector) Count(end int32, value int8) int32 {
	if value == 1 {
		return sl.Count1(end)
	}
	return end - sl.Count1(end)
}

func (sl *DynamicBitvector) Kth0(k int32) int32 {
	if k < 0 || sl.size-sl.totalOnes <= k {
		return -1
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	pos, preLen, newK := sl._findKthZero(k)
	block := sl.blocks[pos]
	for i := int32(0); i < int32(len(block)); i++ {
		if block[i] == 0 {
			if newK == 0 {
				return preLen + i
			}
			newK--
		}
	}
	return -1
}

func (sl *DynamicBitvector) Kth1(k int32) int32 {
	if k < 0 || sl.totalOnes <= k {
		return -1
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	pos, preLen, newK := sl._findKthOne(k)
	block := sl.blocks[pos]
	for i := int32(0); i < int32(len(block)); i++ {
		if block[i] == 1 {
			if newK == 0 {
				return preLen + i
			}
			newK--
		}
	}
	return -1
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
	for i := 0; i < len(sl.blocks); i++ {
		sl.preLen[i] = int32(len(sl.blocks[i]))
	}
	sl.preOnes = append(sl.blockOnes[:0:0], sl.blockOnes...)
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

func (sl *DynamicBitvector) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size {
		return last, lastLen
	}
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree := sl.preLen
	pos = -1
	m := int32(len(tree))
	bitLen := int8(bits.Len32(uint32(m)))
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func (sl *DynamicBitvector) _findKthAndPreOnes(k int32) (pos, index, ones int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k, 0
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size, sl.totalOnes - sl.blockOnes[last]
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	pos = -1
	m := int32(len(tree1))
	bitLen := int8(bits.Len32(uint32(m)))
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree1[next] {
			pos = next
			k -= tree1[pos]
			ones += tree2[pos]
		}
	}
	return pos + 1, k, ones
}

func (sl *DynamicBitvector) _findKthOne(k int32) (pos, preLen, newK int32) {
	if k < sl.blockOnes[0] {
		return 0, 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastOnes := sl.blockOnes[last]
	if tmp := sl.totalOnes - lastOnes; k >= tmp {
		return last, sl.size - int32(len(sl.blocks[last])), k - tmp
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	pos = -1
	m := int32(len(tree2))
	bitLen := int8(bits.Len32(uint32(m)))
	newK = k
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && newK >= tree2[next] {
			pos = next
			preLen += tree1[pos]
			newK -= tree2[pos]
		}
	}
	pos++
	return
}

func (sl *DynamicBitvector) _findKthZero(k int32) (pos, preLen, newK int32) {
	if k < int32(len(sl.blocks[0]))-sl.blockOnes[0] {
		return 0, 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	lastZero := lastLen - sl.blockOnes[last]
	if tmp := sl.size - lastLen + lastZero; k >= tmp {
		return last, sl.size - lastLen, k - tmp
	}
	if sl.shouldRebuildTree {
		sl._buildPreLenAndPreOnes()
	}
	tree1, tree2 := sl.preLen, sl.preOnes
	pos = -1
	m := int32(len(tree2))
	bitLen := int8(bits.Len32(uint32(m)))
	newK = k
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && newK >= tree1[next]-tree2[next] {
			pos = next
			preLen += tree1[pos]
			newK -= tree1[pos] - tree2[pos]
		}
	}
	pos++
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

func test() {
	for i := 0; i < 10; i++ {
		n := rand.Intn(1e4) + 50
		nums := make([]int8, n)
		for i := 0; i < n; i++ {
			nums[i] = int8(rand.Intn(2))
		}
		bv := NewDynamicBitvector(int32(n), func(i int32) int8 { return nums[i] })

		count := func(end int32, v int8) int32 {
			res := int32(0)
			for i := int32(0); i < end; i++ {
				if nums[i] == v {
					res++
				}
			}
			return res
		}

		kth := func(k int32, v int8) int32 {
			res := int32(0)
			for i := int32(0); i < int32(len(nums)); i++ {
				if nums[i] == v {
					if res == k {
						return i
					}
					res++
				}
			}
			return -1
		}

		insert := func(index int32, v int8) {
			nums = append(nums, 0)
			copy(nums[index+1:], nums[index:])
			nums[index] = v
		}
		_ = insert

		pop := func(index int32) int8 {
			res := nums[index]
			nums = append(nums[:index], nums[index+1:]...)
			return res
		}
		_ = pop

		for j := 0; j < 1000; j++ {
			// insert
			insertIndex := rand.Intn(n + 1)
			insertValue := int8(rand.Intn(2))
			insert(int32(insertIndex), insertValue)
			bv.Insert(int32(insertIndex), insertValue)

			// fmt.Println(wm.ToList(), nums, "after insert", insertIndex, insertValue)

			// pop
			popIndex := rand.Intn(len(nums))
			if bv.Pop(int32(popIndex)) != pop(int32(popIndex)) {
				panic("error")
			}

			// count
			countIndex := int32(rand.Intn(len(nums) + 1))
			if bv.Count0(countIndex) != count(countIndex, 0) {
				fmt.Println(bv.Count0(countIndex), count(countIndex, 0), bv.GetAll(), nums, countIndex, n, len(nums))
				panic("error1")
			}
			if bv.Count1(countIndex) != count(countIndex, 1) {
				panic("error2")
			}

			// kth
			kthIndex := int32(rand.Intn(len(nums) + 1))
			if bv.Kth0(kthIndex) != kth(kthIndex, 0) {
				fmt.Println(bv.Kth0(kthIndex), kth(kthIndex, 0), bv.GetAll(), nums, kthIndex, n, len(nums))
				panic("error3")
			}
			if bv.Kth1(kthIndex) != kth(kthIndex, 1) {
				panic("error4")
			}

			// fmt.Println(wm.ToList(), nums, "after pop", popIndex)

			// set
			setIndex := rand.Intn(len(nums))
			setValue := int8(rand.Intn(2))
			nums[setIndex] = setValue
			bv.Set(int32(setIndex), setValue)
			// fmt.Println(wm.ToList(), nums, "after set", setIndex, setValue)

			// len
			if bv.Len() != int32(len(nums)) {
				panic("error")
			}

			// get
			for i := 0; i < len(nums); i++ {
				if bv.Get(int32(i)) != nums[i] {
					fmt.Println(bv.GetAll(), nums, i, n, len(nums))
					panic("error get")
				}
			}

			// toList
			list := bv.GetAll()
			for i := 0; i < len(nums); i++ {
				if list[i] != nums[i] {
					fmt.Println(list, nums, i, list[i], nums[i])
					panic("error toList")
				}
			}
		}
	}

	fmt.Println("ok")
}

func testTime() {
	n := int32(2e5)
	startTime := time.Now()
	bv := NewDynamicBitvector(n, func(i int32) int8 {
		if i%2 == 0 {
			return 0
		}
		return 1
	})

	// Count、Kth、Get
	for i := int32(0); i < n; i++ {
		bv.Count(i, 0)
		bv.Count(n-i, 1)
		bv.Kth(i, 0)
		bv.Kth(n-i, 1)
		bv.Get(i)
	}
	time1 := time.Now()

	// Insert、Pop、Set
	for i := int32(0); i < n; i++ {
		bv.Insert(i, 1)
		bv.Insert(i, 0)
	}
	for i := int32(0); i < n; i++ {
		bv.Pop(n - i)
		bv.Set(i, 1)
	}
	bv.GetAll()
	time2 := time.Now()
	fmt.Println("time1", time1.Sub(startTime)) // 143.811ms
	fmt.Println("time2", time2.Sub(time1))     // 62.4144ms
}
