// api:
//  1.Insert(index int32, v V) -> log(n)
//  2.Pop(index int32) V -> sqrt(n)
//  3.Set(index int32, v V) -> sqrt(n)
//  4.Get(index int32) V -> O(log(sqrt(n)))
//  5.Query(start, end int32) V -> O(sqrt(n))
//    QueryAll() V
//  6.Clear()
//  7.Len() int32
//  8.GetAll() []V
//  9.ForEach(f func(i int32, v V) bool)
// 10.MaxRight(start int32, predicate func(end int32, sum E) bool) int32
// 11.MinLeft(end int32, predicate func(start int32, sum E) bool) int32

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
	// test()
	testTime()
}

func demo() {
	bv := NewSegmentTreeSqrtDecompositionDynamic(10, func(i int32) int32 { return 0 }, -1)
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

// 不用泛型，泛型会导致性能下降.
type E = int32

func (*SegmentTreeSqrtDecompositionDynamic) e() E        { return 0 }
func (*SegmentTreeSqrtDecompositionDynamic) op(a, b E) E { return max32(a, b) }

// 使用分块+树状数组维护的动态数组.
type SegmentTreeSqrtDecompositionDynamic struct {
	n                 int32
	blockSize         int32
	threshold         int32
	shouldRebuildTree bool
	blocks            [][]E
	blockSum          []E
	tree              []int32 // 每个块块长的前缀和
}

func NewSegmentTreeSqrtDecompositionDynamic(n int32, f func(i int32) E, blockSize int32) *SegmentTreeSqrtDecompositionDynamic {
	if blockSize == -1 {
		blockSize = int32(math.Sqrt(float64(n))) + 1
	}
	if blockSize < 100 {
		blockSize = 100 // 防止 blockSize 过小
	}
	res := &SegmentTreeSqrtDecompositionDynamic{n: n, blockSize: blockSize, threshold: blockSize << 1, shouldRebuildTree: true}
	blockCount := (n + blockSize - 1) / blockSize
	blocks, blockSum := make([][]E, blockCount), make([]E, blockCount)
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
	}
	res.blocks, res.blockSum = blocks, blockSum
	return res
}

func (sl *SegmentTreeSqrtDecompositionDynamic) Insert(index int32, value E) {
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []E{value})
		sl.blockSum = append(sl.blockSum, value)
		sl.shouldRebuildTree = true
		sl.n++
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
	sl.blocks[pos] = append(sl.blocks[pos], sl.e())
	copy(sl.blocks[pos][startIndex+1:], sl.blocks[pos][startIndex:])
	sl.blocks[pos][startIndex] = value

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); n > sl.threshold {
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][sl.blockSize:] // !注意max的设置(为了让左右互不影响)
		sl.blocks[pos] = sl.blocks[pos][:sl.blockSize:sl.blockSize]
		sl.blockSum = append(sl.blockSum, sl.e())
		copy(sl.blockSum[pos+2:], sl.blockSum[pos+1:])
		sl._updateSum(pos)
		sl._updateSum(pos + 1)
		sl.shouldRebuildTree = true
	}

	sl.n++
	return
}

func (sl *SegmentTreeSqrtDecompositionDynamic) Pop(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.n--
	sl._updateTree(pos, false)

	copy(sl.blocks[pos][startIndex:], sl.blocks[pos][startIndex+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if value != sl.e() {
		sl._updateSum(pos)
	}

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

func (sl *SegmentTreeSqrtDecompositionDynamic) Get(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SegmentTreeSqrtDecompositionDynamic) Set(index int32, value E) {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	oldValue := sl.blocks[pos][startIndex]
	if oldValue == value {
		return
	}
	sl.blocks[pos][startIndex] = value
	sl._updateSum(pos)
}

func (sl *SegmentTreeSqrtDecompositionDynamic) Query(start, end int32) E {
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

func (sl *SegmentTreeSqrtDecompositionDynamic) QueryAll() E {
	res := sl.e()
	for _, v := range sl.blockSum {
		res = sl.op(res, v)
	}
	return res
}

func (sl *SegmentTreeSqrtDecompositionDynamic) Len() int32 {
	return sl.n
}

func (sl *SegmentTreeSqrtDecompositionDynamic) Clear() {
	sl.n = 0
	sl.shouldRebuildTree = true
	sl.blocks = sl.blocks[:0]
	sl.blockSum = sl.blockSum[:0]
	sl.tree = sl.tree[:0]
}

func (sl *SegmentTreeSqrtDecompositionDynamic) GetAll() []E {
	res := make([]E, 0, sl.n)
	for _, block := range sl.blocks {
		res = append(res, block...)
	}
	return res
}

func (sl *SegmentTreeSqrtDecompositionDynamic) ForEach(f func(i int32, v E) (shouldBreak bool)) {
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

// 查询最大的 end 使得切片 [start:end] 内的值满足 predicate.
func (st *SegmentTreeSqrtDecompositionDynamic) MaxRight(start int32, predicate func(end int32, sum E) bool) int32 {
	if start >= st.n {
		return st.n
	}

	curSum := st.e()
	res := start
	bid, startPos := st._findKth(start)

	// 散块内
	{
		pos := startPos
		block := st.blocks[bid]
		m := int32(len(block))
		for ; pos < m; pos++ {
			nextRes, nextSum := res+1, st.op(curSum, block[pos])
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				return res
			}
		}
	}
	bid++

	// 整块跳跃
	{
		m := int32(len(st.blocks))
		for ; bid < m; bid++ {
			nextRes := res + int32(len(st.blocks[bid]))
			nextSum := st.op(curSum, st.blockSum[bid])
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				// 答案在这个块内
				block := st.blocks[bid]
				for _, v := range block {
					nextRes, nextSum = res+1, st.op(curSum, v)
					if predicate(nextRes, nextSum) {
						res, curSum = nextRes, nextSum
					} else {
						return res
					}
				}
			}
		}
	}

	return res
}

// 查询最小的 start 使得切片 [start:end] 内的值满足 predicate.
func (st *SegmentTreeSqrtDecompositionDynamic) MinLeft(end int32, predicate func(start int32, sum E) bool) int32 {
	if end <= 0 {
		return 0
	}

	curSum := st.e()
	res := end
	bid, startPos := st._findKth(end - 1)

	// 散块内
	{
		pos := startPos
		block := st.blocks[bid]
		for ; pos >= 0; pos-- {
			nextRes, nextSum := res-1, st.op(block[pos], curSum)
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				return res
			}
		}
	}
	bid--

	// 整块跳跃
	{
		for ; bid >= 0; bid-- {
			nextRes := res - int32(len(st.blocks[bid]))
			nextSum := st.op(st.blockSum[bid], curSum)
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				// 答案在这个块内
				block := st.blocks[bid]
				for i := int32(len(block)) - 1; i >= 0; i-- {
					nextRes, nextSum = res-1, st.op(block[i], curSum)
					if predicate(nextRes, nextSum) {
						res, curSum = nextRes, nextSum
					} else {
						return res
					}
				}
			}
		}
	}
	return res
}

func (sl *SegmentTreeSqrtDecompositionDynamic) _rebuildTree() {
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

func (sl *SegmentTreeSqrtDecompositionDynamic) _updateTree(index int32, addOne bool) {
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

func (sl *SegmentTreeSqrtDecompositionDynamic) _findKth(k int32) (pos, index int32) {
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

func (sl *SegmentTreeSqrtDecompositionDynamic) _updateSum(pos int32) {
	sum := sl.e()
	for _, v := range sl.blocks[pos] {
		sum = sl.op(sum, v)
	}
	sl.blockSum[pos] = sum
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
			nums[i] = rand.Intn(1000)
		}
		seg := NewSegmentTreeSqrtDecompositionDynamic(n, func(i int32) E { return E(nums[i]) }, -1)

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
				sum_ = seg.op(sum_, E(nums[i]))
			}
			if seg.Query(start, end) != sum_ {
				fmt.Println("Query Error")
				panic("Query Error")
			}

			// QueryAll
			sum_ = E(0)
			for _, v := range nums {
				sum_ = seg.op(sum_, E(v))
			}
			if seg.QueryAll() != sum_ {
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
			value = nums[index]
			nums = append(nums[:index], nums[index+1:]...)
			if seg.Pop(index) != E(value) {
				fmt.Println("Pop Error")
				panic("Pop Error")
			}

			// ForEach
			sum_ = E(0)
			seg.ForEach(func(i int32, v E) bool {
				sum_ = seg.op(sum_, v)
				return false
			})
			if sum_ != seg.QueryAll() {
				fmt.Println("ForEach Error")
				panic("ForEach Error")
			}

			// MaxRight
			maxRightBf := func(start int32, predicate func(end int32, sum E) bool) (res int32) {
				res = start
				curSum := seg.e()
				for i := start; i < int32(len(nums)); i++ {
					curSum = seg.op(curSum, E(nums[i]))
					if !predicate(i+1, curSum) {
						return
					}
					res = i + 1
				}
				return
			}

			minLeftBf := func(end int32, predicate func(start int32, sum E) bool) (res int32) {
				res = end
				curSum := seg.e()
				for i := end - 1; i >= 0; i-- {
					curSum = seg.op(curSum, E(nums[i]))
					if !predicate(i, curSum) {
						return
					}
					res = i
				}
				return
			}

			{
				start := rand.Int31n(n)
				upper := rand.Intn(100)
				res1 := seg.MaxRight(start, func(end int32, sum E) bool { return sum <= E(upper) })
				res2 := maxRightBf(start, func(end int32, sum E) bool { return sum <= E(upper) })
				if res1 != res2 {
					fmt.Println("MaxRight Error")
					panic("MaxRight Error")
				}
				res3 := seg.MinLeft(start, func(start int32, sum E) bool { return sum <= E(upper) })
				res4 := minLeftBf(start, func(start int32, sum E) bool { return sum <= E(upper) })
				if res3 != res4 {
					fmt.Println("MinLeft Error")
					panic("MinLeft Error")
				}
			}

		}
	}
	fmt.Println("Pass")
}

func testTime() {
	// 1e5
	n := int32(1e5)
	nums := make([]int, n)
	for i := 0; i < int(n); i++ {
		nums[i] = rand.Intn(5)
	}

	time1 := time.Now()
	seg := NewSegmentTreeSqrtDecompositionDynamic(n, func(i int32) int32 { return E(nums[i]) }, -1)

	for i := int32(0); i < n; i++ {
		seg.Get(i)
		seg.Set(i, i)
		seg.Query(i, n)
		seg.QueryAll()
		seg.Insert(i, i)
		if i&1 == 0 {
			seg.Pop(i)
		}
		seg.MaxRight(i, func(end int32, sum E) bool { return true })
		seg.MinLeft(i, func(start int32, sum E) bool { return true })
	}
	fmt.Println("Time1", time.Since(time1)) // Time1 318.991125ms
}
