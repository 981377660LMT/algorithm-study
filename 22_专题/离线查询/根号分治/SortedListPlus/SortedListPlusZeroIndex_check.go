package main

import (
	"cmp"
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"time"
)

// -------------------- 测试入口 --------------------
func main() {
	fmt.Println("Begin testTime for SortedListPlus vs naiveSL ...")
	for i := 0; i < 100; i++ {
		judge() // 进行多次随机对拍测试
	}
}

// judge: 在此函数中进行随机测试，对拍 SortedListPlus vs naiveSL
func judge() {
	const N = 50000           // 操作次数
	const rangeVal = int(1e9) // 值域
	rand.Seed(time.Now().UnixNano())

	Q := rand.Intn(1000) + 1
	// 先产生 N 个随机数，排好序，放入 SortedListPlus & naiveSL
	nums := make([]int, N)
	for i := 0; i < N; i++ {
		nums[i] = rand.Intn(rangeVal)
	}
	sort.Ints(nums)

	slp := NewSortedListPlusFrom(nums)
	naive := NewNaiseSlFrom(nums) // 朴素做法

	for i := 1; i <= N; i++ {
		op := rand.Intn(7) // 0..6 => 7 种操作
		x := rand.Intn(rangeVal)
		switch op {

		case 0: // Insert
			// 如需测试插入，可解注释
			// slp.Insert(x)
			// naive.Insert(x)

		case 1: // Erase
			// 如需测试删除，可解注释
			slp.Erase(x)
			naive.Erase(x)

		case 2: // Size
			if slp.Size() != naive.Size() {
				panic(fmt.Sprintf("Mismatch in Size, got slp=%d naive=%d",
					slp.Size(), naive.Size()))
			}

		case 3: // Rank
			rslp := slp.Rank(x)
			rnaive := naive.Rank(x)
			if rslp != rnaive {
				panic(fmt.Sprintf("Mismatch in Rank(%d), got slp=%d naive=%d",
					x, rslp, rnaive))
			}

		case 4: // Count
			cslp := slp.Count(x)
			cnaive := naive.Count(x)
			if cslp != cnaive {
				panic(fmt.Sprintf("Mismatch in Count(%d), got slp=%d naive=%d",
					x, cslp, cnaive))
			}

		case 5: // LessEqual
			rslp, ok1 := slp.LessEqual(x)
			rnaive, ok2 := naive.LessEqual(x)
			if ok1 != ok2 {
				panic(fmt.Sprintf("Mismatch in LessEqual(%d) existence, slpOk=%v naiveOk=%v",
					x, ok1, ok2))
			}
			if ok1 && rslp != rnaive {
				panic(fmt.Sprintf("Mismatch in LessEqual(%d), got slp=%d naive=%d",
					x, rslp, rnaive))
			}

		case 6: // GreaterEqual
			rslp, ok1 := slp.GreaterEqual(x)
			rnaive, ok2 := naive.GreaterEqual(x)
			if ok1 != ok2 {
				panic(fmt.Sprintf("Mismatch in GreaterEqual(%d) existence, slpOk=%v naiveOk=%v",
					x, ok1, ok2))
			}
			if ok1 && rslp != rnaive {
				panic(fmt.Sprintf("Mismatch in GreaterEqual(%d), got slp=%d naive=%d",
					x, rslp, rnaive))
			}
		}
	}

	fmt.Println("testTime: All tests passed! SortedListPlus and naiveSL match on random data.")
}

// -------------------- 朴素对拍结构 --------------------
type naiveSL[T cmp.Ordered] struct {
	data []T
}

func NewNaiseSlFrom[T cmp.Ordered](sorted []T) *naiveSL[T] {
	sorted = slices.Clone(sorted)
	return &naiveSL[T]{data: sorted}
}

// 若已存在则不再插入；否则插入保持有序
func (n *naiveSL[T]) Insert(x T) {
	pos, ok := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	if ok {
		return
	}
	n.data = slices.Insert(n.data, pos, x)
}

func (n *naiveSL[T]) Erase(x T) {
	pos, _ := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	if pos < len(n.data) && n.data[pos] == x {
		n.data = slices.Delete(n.data, pos, pos+1)
	}
}

func (n *naiveSL[T]) Size() int { return len(n.data) }

func (n *naiveSL[T]) Rank(x T) int {
	pos, _ := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	return pos
}

// Count(x) = upper_bound(x) - lower_bound(x)
func (n *naiveSL[T]) Count(x T) int {
	lpos, _ := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	rpos, _ := slices.BinarySearchFunc(n.data, x, func(a, b T) int {
		// upper_bound => 找第一个 > x
		if a <= b {
			return -1
		}
		return +1
	})
	return rpos - lpos
}

func (n *naiveSL[T]) LessEqual(x T) (res T, ok bool) {
	// upper_bound(x) 前驱
	pos, _ := slices.BinarySearchFunc(n.data, x, func(a, b T) int {
		if a <= b {
			return -1
		}
		return +1
	})
	if pos == 0 {
		return res, false
	}
	return n.data[pos-1], true
}

func (n *naiveSL[T]) GreaterEqual(x T) (res T, ok bool) {
	pos, _ := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	if pos == len(n.data) {
		return res, false
	}
	return n.data[pos], true
}

func cmpOrder[T cmp.Ordered](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return +1
	default:
		return 0
	}
}

// -------------------- SortedListPlus 实现 --------------------
type SortedListPlus[T cmp.Ordered] struct {
	load1  int
	load2  int
	load1X int // = load1 * 2
	load2X int // = load2 * 2
	lg2    int // = floor(log2(load2X))
	minVal T

	blocks [][]T // blocks[i] 表示第 i 个小块
	bitCnt []int // Fenwick 0-based: bitCnt[i] 存小块 i 的大小
	blkMax []T   // blkMax[i] 存 blocks[i] 的最大值

	segMax  []T   // segMax[j] 表示第 j 小段的最大值
	segSize []int // segSize[j] 表示第 j 小段非空块个数

	elemCnt int // 整个结构中元素总数
}

// NewSortedListPlusFrom: 传入一个可能乱序的切片 `arr`，先排序，再构建
func NewSortedListPlusFrom[T cmp.Ordered](arr []T) *SortedListPlus[T] {
	slices.Sort(arr)

	sl := &SortedListPlus[T]{
		load1:  200,
		load2:  64,
		minVal: getMinValue[T](),
	}
	sl.load1X = sl.load1 * 2    // 400
	sl.load2X = sl.load2 * 2    // 128
	sl.lg2 = intLog2(sl.load2X) // 7
	sl.elemCnt = len(arr)

	// 按 load1 分块
	blockCount := (len(arr) + sl.load1 - 1) / sl.load1
	sl.blocks = make([][]T, 0, blockCount)
	start := 0
	for start < len(arr) {
		end := start + sl.load1
		if end > len(arr) {
			end = len(arr)
		}
		sl.blocks = append(sl.blocks, Clone(arr[start:end]))
		start = end
	}

	sl._expand()
	return sl
}

func (sl *SortedListPlus[T]) Clear() {
	sl.blocks = sl.blocks[:0]
	sl.bitCnt = sl.bitCnt[:0]
	sl.blkMax = sl.blkMax[:0]
	sl.segMax = sl.segMax[:0]
	sl.segSize = sl.segSize[:0]
	sl.elemCnt = 0
}

func (sl *SortedListPlus[T]) Size() int {
	return sl.elemCnt
}

// Insert 插入元素 x（若已存在则跳过）
func (sl *SortedListPlus[T]) Insert(x T) {
	// 若空
	if len(sl.segSize) == 0 {
		sl.blocks = append(sl.blocks, []T{x})
		sl.bitCnt = append(sl.bitCnt, 1)
		sl.blkMax = append(sl.blkMax, x)
		sl.segMax = append(sl.segMax, x)
		sl.segSize = append(sl.segSize, 1)
		sl.elemCnt = 1
		return
	}

	bi, pos := sl.lowerBound(x)

	// 若已存在则直接返回
	if pos < len(sl.blocks[bi]) && sl.blocks[bi][pos] == x {
		return
	}
	sl.elemCnt++

	sl.fenwickUpdate(bi, 1)
	sl.blocks[bi] = slices.Insert(sl.blocks[bi], pos, x)

	if x > sl.blkMax[bi] {
		sl.blkMax[bi] = x
	}
	segi := bi >> sl.lg2
	if x > sl.segMax[segi] {
		sl.segMax[segi] = x
	}

	if len(sl.blocks[bi]) >= sl.load1X {
		sl.splitBlock(segi, bi)
	}
}

// Erase 删除值为 x（若不存在则忽略）
func (sl *SortedListPlus[T]) Erase(x T) {
	if len(sl.segMax) == 0 {
		return
	}
	if sl.segMax[len(sl.segMax)-1] < x {
		return
	}
	bi, pos := sl.lowerBound(x)
	if pos == len(sl.blocks[bi]) || sl.blocks[bi][pos] != x {
		return
	}
	sl.fenwickUpdate(bi, -1)
	sl.elemCnt--

	sl.blocks[bi] = slices.Delete(sl.blocks[bi], pos, pos+1)
	if len(sl.blocks[bi]) == 0 {
		sl.eraseBlock(bi)
	} else {
		sl.blkMax[bi] = sl.blocks[bi][len(sl.blocks[bi])-1]
		segi := bi >> sl.lg2
		bj := (segi << sl.lg2) + sl.segSize[segi] - 1
		if bi == bj {
			sl.segMax[segi] = sl.blkMax[bi]
		}
	}
}

// At(k) => 第 k(0-based) 小元素
func (sl *SortedListPlus[T]) At(k int) T {
	if k < 0 || k >= sl.elemCnt {
		panic("At: out of range")
	}
	blkIndex := sl.findFenwick(k)
	offset := k - sl.prefixSum(blkIndex-1)
	return sl.blocks[blkIndex][offset]
}

// Rank(x) => lower_bound(x) 在全局的下标（有多少元素 < x）
func (sl *SortedListPlus[T]) Rank(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, pos := sl.lowerBound(x)
	return sl.prefixSum(bi-1) + pos
}

// Count(x) => upper_bound(x) - lower_bound(x)
func (sl *SortedListPlus[T]) Count(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bu, pu := sl.upperBound(x)
	bl, pl := sl.lowerBound(x)
	rkU := sl.prefixSum(bu-1) + pu
	rkL := sl.prefixSum(bl-1) + pl
	return rkU - rkL
}

// LessEqual => 最大 <= x
func (sl *SortedListPlus[T]) LessEqual(x T) (res T, ok bool) {
	if len(sl.segSize) == 0 {
		return
	}
	bi, pos := sl.upperBound(x)
	if bi == 0 && pos == 0 {
		return res, false
	}
	pb, pe := sl.prev(bi, pos)
	return sl.blocks[pb][pe], true
}

// GreaterEqual => 最小 >= x
func (sl *SortedListPlus[T]) GreaterEqual(x T) (res T, ok bool) {
	if len(sl.segSize) == 0 {
		return
	}
	bi, pos := sl.lowerBound(x)
	if bi >= len(sl.blocks) || pos == len(sl.blocks[bi]) {
		return res, false
	}
	return sl.blocks[bi][pos], true
}

// ------------------ 以下是内部逻辑 ------------------

func (sl *SortedListPlus[T]) lowerBound(x T) (bi, pos int) {
	// 1. 小段二分
	l, r := 0, len(sl.segSize)-1
	for l < r {
		mid := (l + r) >> 1
		if sl.segMax[mid] >= x {
			r = mid
		} else {
			l = mid + 1
		}
	}
	segi := l

	// 2. 在该小段内做块二分 (每段最多 load2X=128 块)
	startBlk := segi << sl.lg2
	endBlk := startBlk + sl.segSize[segi] - 1
	lb, rb := startBlk, endBlk
	for lb < rb {
		mid := (lb + rb) >> 1
		if sl.blkMax[mid] >= x {
			rb = mid
		} else {
			lb = mid + 1
		}
	}
	bi = lb

	// 3. 块内二分 => 找第一个 >= x
	pos, _ = slices.BinarySearchFunc(sl.blocks[bi], x, cmpOrder[T])
	return
}

func (sl *SortedListPlus[T]) upperBound(x T) (bi, pos int) {
	// 1. 小段二分
	l, r := 0, len(sl.segSize)-1
	for l < r {
		mid := (l + r) >> 1
		if sl.segMax[mid] > x {
			r = mid
		} else {
			l = mid + 1
		}
	}
	segi := l

	// 2. 在该小段内做块二分
	startBlk := segi << sl.lg2
	endBlk := startBlk + sl.segSize[segi] - 1
	lb, rb := startBlk, endBlk
	for lb < rb {
		mid := (lb + rb) >> 1
		if sl.blkMax[mid] > x {
			rb = mid
		} else {
			lb = mid + 1
		}
	}
	bi = lb

	// 3. 块内二分 => 找第一个 > x
	pos, _ = slices.BinarySearchFunc(sl.blocks[bi], x, func(a, b T) int {
		// upper_bound => a<=b 要继续右移
		if a <= b {
			return -1
		}
		return +1
	})
	return
}

// prev(bi,pos) => 前驱
func (sl *SortedListPlus[T]) prev(bi, pos int) (int, int) {
	if pos > 0 {
		return bi, pos - 1
	}
	if bi == 0 {
		return 0, 0
	}
	return bi - 1, len(sl.blocks[bi-1]) - 1
}

// splitBlock => 分裂块
func (sl *SortedListPlus[T]) splitBlock(segi, bi int) {
	oldBlk := sl.blocks[bi]
	newBlk := slices.Clone(oldBlk[sl.load1:])
	sl.blocks[bi] = oldBlk[:sl.load1]

	biNew := bi + 1
	sl.blocks = slices.Insert(sl.blocks, biNew, []T{})
	sl.blocks[biNew] = newBlk

	sl.blkMax = slices.Insert(sl.blkMax, biNew, newBlk[len(newBlk)-1])
	sl.bitCnt = slices.Insert(sl.bitCnt, biNew, len(newBlk))

	sl.fenwicksumRebuild()

	sl.segSize[segi]++
	if sl.segSize[segi] == sl.load2X {
		sl._expand()
	}
}

// eraseBlock => 某块为空，删除之
func (sl *SortedListPlus[T]) eraseBlock(bi int) {
	sl.blocks = slices.Delete(sl.blocks, bi, bi+1)
	sl.blkMax = slices.Delete(sl.blkMax, bi, bi+1)
	sl.bitCnt = slices.Delete(sl.bitCnt, bi, bi+1)

	segi := bi >> sl.lg2
	sl.segSize[segi]--
	if sl.segSize[segi] == 0 {
		if segi == len(sl.segSize)-1 {
			sl.segMax = sl.segMax[:segi]
			sl.segSize = sl.segSize[:segi]
		} else {
			sl._expand()
		}
	} else {
		sl.fenwicksumRebuild()
		bj := (segi << sl.lg2) + sl.segSize[segi] - 1
		if bi == bj+1 {
			sl.segMax[segi] = sl.blkMax[bj]
		}
	}
}

// _expand 全局重构 => 全部非空块连续存储
func (sl *SortedListPlus[T]) _expand() {
	oldBlocks := sl.blocks
	if len(oldBlocks) == 0 {
		sl.Clear()
		return
	}
	// 统计非空块数
	c := 0
	for _, b := range oldBlocks {
		if len(b) > 0 {
			c++
		}
	}
	if c == 0 {
		sl.Clear()
		return
	}
	// ***修正关键处***：每段可容纳 load2X 个块
	segn := (c + sl.load2X - 1) / sl.load2X

	ec := sl.elemCnt
	sl.Clear()
	sl.elemCnt = ec

	sl.segMax = make([]T, segn)
	sl.segSize = make([]int, segn)

	j := 0
	for _, block := range oldBlocks {
		if len(block) == 0 {
			continue
		}
		segi := j >> sl.lg2 // = j/load2X
		sl.blocks = append(sl.blocks, block)
		sl.bitCnt = append(sl.bitCnt, len(block))
		sl.blkMax = append(sl.blkMax, block[len(block)-1])

		sl.segSize[segi]++
		if block[len(block)-1] > sl.segMax[segi] {
			sl.segMax[segi] = block[len(block)-1]
		}
		j++
	}
	sl.fenwicksumRebuild()
}

// fenwicksumRebuild => 重新构造 Fenwick
func (sl *SortedListPlus[T]) fenwicksumRebuild() {
	// sl.bitCnt = slices.Clone(sl.bitCnt) // 保守起见，可清空再写
	for i := range sl.bitCnt {
		sl.bitCnt[i] = len(sl.blocks[i])
	}
	for i := 0; i < len(sl.bitCnt); i++ {
		j := i | (i + 1)
		if j < len(sl.bitCnt) {
			sl.bitCnt[j] += sl.bitCnt[i]
		}
	}
}

// fenwickUpdate => 对第 idx 块计数增量更新
func (sl *SortedListPlus[T]) fenwickUpdate(idx, delta int) {
	for idx < len(sl.bitCnt) {
		sl.bitCnt[idx] += delta
		idx |= (idx + 1)
	}
}

// prefixSum(idx) => bitCnt[0..idx] 之和
func (sl *SortedListPlus[T]) prefixSum(idx int) int {
	if idx < 0 {
		return 0
	}
	res := 0
	for idx >= 0 {
		res += sl.bitCnt[idx]
		idx = (idx & (idx + 1)) - 1
	}
	return res
}

// findFenwick(k) => 找到块下标，使 prefixSum(块下标-1) <= k < prefixSum(块下标)
func (sl *SortedListPlus[T]) findFenwick(k int) int {
	idx := 0
	bitMask := 1
	for bitMask < len(sl.bitCnt) {
		bitMask <<= 1
	}
	half := bitMask >> 1
	for half > 0 {
		tmp := idx + half - 1
		if tmp < len(sl.bitCnt) && sl.bitCnt[tmp] <= k {
			k -= sl.bitCnt[tmp]
			idx += half
		}
		half >>= 1
	}
	return idx
}

// 工具函数
func intLog2(x int) int {
	e := 0
	for (1 << e) <= x {
		e++
	}
	return e - 1
}

func getMinValue[T cmp.Ordered]() T {
	var zero T
	return zero
}

func Clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
}
