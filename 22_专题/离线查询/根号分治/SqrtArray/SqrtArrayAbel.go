// 动态分块数组/SqrtArray/SqrtArrayWithSum
//
// api:
//  1.Insert(index int32, v V)
//  2.Pop(index int32) V
//  3.Set(index int32, v V)
//  4.Get(index int32) V
//  5.Sum(start, end int32) V
//    SumAll() V
//  6.Clear()
//  7.Len() int32
//  8.GetAll() []V
//  9.ForEach(f func(i int32, v V) bool)

package main

import (
	"fmt"
	"math"
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
	bv := NewSqrtArrayAbel(10, func(i int32) int { return int(i) }, -1)
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
	bv.Insert(-987, 1)
	bv.Insert(999, 1)
	fmt.Println(bv.GetAll())
}

type E = int

func (*SqrtArrayAbel) e() E        { return 0 }
func (*SqrtArrayAbel) op(a, b E) E { return a + b }
func (*SqrtArrayAbel) inv(a E) E   { return -a }

// 使用分块+树状数组维护的动态数组.
type SqrtArrayAbel struct {
	n                 int32
	blockSize         int32
	threshold         int32
	shouldRebuildTree bool
	totalSum          E
	blocks            [][]E
	blockSum          []E
	tree              []int32 // 每个块块长的前缀和
}

func NewSqrtArrayAbel(n int32, f func(i int32) E, blockSize int32) *SqrtArrayAbel {
	if blockSize == -1 {
		blockSize = int32(math.Sqrt(float64(n))) + 1
	}

	res := &SqrtArrayAbel{n: n, blockSize: blockSize, threshold: blockSize << 1, shouldRebuildTree: true}
	blockCount := (n + blockSize - 1) / blockSize
	blocks, blockSum := make([][]E, blockCount), make([]E, blockCount)
	totalSum := res.e()
	for bid := int32(0); bid < blockCount; bid++ {
		start, end := bid*blockSize, (bid+1)*blockSize
		if end > n {
			end = n
		}
		bucket := make([]E, end-start)
		sum := res.e()
		for i := start; i < end; i++ {
			bucket[i-start] = f(i)
			sum = res.op(sum, bucket[i-start])
		}
		blocks[bid], blockSum[bid] = bucket, sum
		totalSum = res.op(totalSum, sum)
	}
	res.totalSum, res.blocks, res.blockSum = totalSum, blocks, blockSum
	return res
}

func (sl *SqrtArrayAbel) Insert(index int32, value E) {
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []E{value})
		sl.blockSum = append(sl.blockSum, value)
		sl.shouldRebuildTree = true
		sl.n++
		sl.totalSum = value
		return
	}

	if index < 0 {
		index += sl.n
	}
	if index < 0 {
		index = 0
	}
	if index > sl.n {
		index = sl.n
	}

	pos, startIndex := sl._findKth(index)
	sl._updateTree(pos, true)
	sl.blockSum[pos] = sl.op(sl.blockSum[pos], value)
	sl.totalSum = sl.op(sl.totalSum, value)
	sl.blocks[pos] = append(sl.blocks[pos], sl.e())
	copy(sl.blocks[pos][startIndex+1:], sl.blocks[pos][startIndex:])
	sl.blocks[pos][startIndex] = value

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); n > sl.threshold {
		totalSum := sl.blockSum[pos]
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][sl.blockSize:] // !注意max的设置(为了让左右互不影响)
		sl.blocks[pos] = sl.blocks[pos][:sl.blockSize:sl.blockSize]
		sum1 := sl.e()
		for _, v := range sl.blocks[pos] {
			sum1 = sl.op(sum1, v)
		}
		sum2 := sl.op(totalSum, sl.inv(sum1))
		sl.blockSum = append(sl.blockSum, sl.e())
		copy(sl.blockSum[pos+2:], sl.blockSum[pos+1:])
		sl.blockSum[pos] = sum1
		sl.blockSum[pos+1] = sum2
		sl.shouldRebuildTree = true
	}

	sl.n++
	return
}

func (sl *SqrtArrayAbel) Pop(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.n--
	sl._updateTree(pos, false)
	inv := sl.inv(value)
	sl.blockSum[pos] = sl.op(sl.blockSum[pos], inv)
	sl.totalSum = sl.op(sl.totalSum, inv)

	copy(sl.blocks[pos][startIndex:], sl.blocks[pos][startIndex+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]

	if len(sl.blocks[pos]) == 0 {
		// !delete block
		copy(sl.blocks[pos:], sl.blocks[pos+1:])
		sl.blocks = sl.blocks[:len(sl.blocks)-1]
		copy(sl.blockSum[pos:], sl.blockSum[pos+1:])
		sl.blockSum = sl.blockSum[:len(sl.blockSum)-1]
		sl.shouldRebuildTree = true
	}
	return value
}

func (sl *SqrtArrayAbel) Get(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SqrtArrayAbel) Set(index int32, value E) {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	oldValue := sl.blocks[pos][startIndex]
	if oldValue == value {
		return
	}
	sl.blocks[pos][startIndex] = value
	diff := sl.op(value, sl.inv(oldValue))
	sl.blockSum[pos] = sl.op(sl.blockSum[pos], diff)
	sl.totalSum = sl.op(sl.totalSum, diff)
}

func (sl *SqrtArrayAbel) Sum(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > sl.n {
		end = sl.n
	}
	if start >= end {
		return sl.e()
	}

	bid1, startIndex1 := sl._findKth(start)
	bid2, startIndex2 := sl._findKth(end)
	start, end = startIndex1, startIndex2
	res := sl.e()
	if bid1 == bid2 {
		block := sl.blocks[bid1]
		for i := start; i < end; i++ {
			res = sl.op(res, block[i])
		}
	} else {
		if start < int32(len(sl.blocks[bid1])) {
			block1 := sl.blocks[bid1]
			for i := start; i < int32(len(block1)); i++ {
				res = sl.op(res, block1[i])
			}
		}
		for i := bid1 + 1; i < bid2; i++ {
			res = sl.op(res, sl.blockSum[i])
		}
		if m := int32(len(sl.blocks)); bid2 < m && end > 0 {
			block2 := sl.blocks[bid2]
			tmp := sl.e()
			for i := int32(0); i < end; i++ {
				tmp = sl.op(tmp, block2[i])
			}
			res = sl.op(res, tmp)
		}
	}
	return res
}

func (sl *SqrtArrayAbel) SumAll() E {
	return sl.totalSum
}

func (sl *SqrtArrayAbel) Len() int32 {
	return sl.n
}

func (sl *SqrtArrayAbel) Clear() {
	sl.n = 0
	sl.shouldRebuildTree = true
	sl.totalSum = sl.e()
	sl.blocks = sl.blocks[:0]
	sl.blockSum = sl.blockSum[:0]
	sl.tree = sl.tree[:0]
}

func (sl *SqrtArrayAbel) GetAll() []E {
	res := make([]E, 0, sl.n)
	for _, block := range sl.blocks {
		res = append(res, block...)
	}
	return res
}

func (sl *SqrtArrayAbel) ForEach(f func(i int32, v E) (shouldBreak bool)) {
	ptr := int32(0)
	for _, block := range sl.blocks {
		for _, v := range block {
			if f(ptr, v) {
				return
			}
			ptr++
		}
	}
}

func (sl *SqrtArrayAbel) _rebuildTree() {
	sl.tree = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = int32(len(sl.blocks[i]))
	}
	tree := sl.tree
	m := int32(len(tree))
	for i := int32(0); i < m; i++ {
		j := i | (i + 1)
		if j < m {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *SqrtArrayAbel) _updateTree(index int32, addOne bool) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
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

func (sl *SqrtArrayAbel) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.n {
		return last, lastLen
	}
	if k >= sl.n-lastLen {
		return last, k + lastLen - sl.n
	}
	if sl.shouldRebuildTree {
		sl._rebuildTree()
	}
	tree := sl.tree
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
	for i := int32(0); i < 100; i++ {
		n := rand.Int31n(10000) + 1000
		nums := make([]int, n)
		for i := int32(0); i < n; i++ {
			nums[i] = rand.Intn(100)
		}
		seg := NewSqrtArrayAbel(n, func(i int32) E { return E(nums[i]) }, -1)

		for j := 0; j < 1000; j++ {
			// Get
			index := rand.Int31n(n)
			if seg.Get(index) != E(nums[index]) {
				fmt.Println("Get Error")
				panic("Get Error")
			}

			// Set
			index = rand.Int31n(n)
			value := rand.Intn(100)
			nums[index] = value
			seg.Set(index, E(value))
			if seg.Get(index) != E(value) {
				fmt.Println("Set Error")
				panic("Set Error")
			}

			// Query
			start, end := rand.Int31n(n), rand.Int31n(n)
			if start > end {
				start, end = end, start
			}
			sum_ := E(0)
			for i := start; i < end; i++ {
				sum_ += E(nums[i])
			}
			if seg.Sum(start, end) != sum_ {
				fmt.Println("Query Error")
				panic("Query Error")
			}

			// QueryAll
			sum_ = E(0)
			for _, v := range nums {
				sum_ += E(v)
			}
			if seg.SumAll() != sum_ {
				fmt.Println("QueryAll Error")
				panic("QueryAll Error")
			}

			// GetAll
			all := seg.GetAll()
			for i, v := range all {
				if v != E(nums[i]) {
					fmt.Println("GetAll Error")
					panic("GetAll Error")
				}
			}

			// Insert
			index = rand.Int31n(n)
			value = rand.Intn(100)
			nums = append(nums, 0)
			copy(nums[index+1:], nums[index:])
			nums[index] = value
			seg.Insert(index, E(value))

			// Pop
			index = rand.Int31n(n)
			value = E(nums[index])
			nums = append(nums[:index], nums[index+1:]...)
			if seg.Pop(index) != value {
				fmt.Println("Pop Error")
				panic("Pop Error")
			}

			// ForEach
			sum_ = E(0)
			seg.ForEach(func(i int32, v E) bool {
				sum_ += v
				return false
			})
			if sum_ != seg.SumAll() {
				fmt.Println("ForEach Error")
				panic("ForEach Error")
			}
		}
	}
	fmt.Println("Pass")
}

func testTime() {
	// 2e5
	n := int32(2e5)
	nums := make([]int, n)
	for i := 0; i < int(n); i++ {
		nums[i] = rand.Intn(5)
	}

	time1 := time.Now()
	seg := NewSqrtArrayAbel(n, func(i int32) int { return nums[i] }, -1)

	for i := int32(0); i < n; i++ {
		seg.Get(i)
		seg.Set(i, int(E(i)))
		seg.Sum(i, n)
		seg.SumAll()
		seg.Insert(i, int(E(i)))
		if i&1 == 0 {
			seg.Pop(i)
		}
		seg.SumAll()
	}
	fmt.Println("Time1", time.Since(time1)) // Time1 144.308ms
}
