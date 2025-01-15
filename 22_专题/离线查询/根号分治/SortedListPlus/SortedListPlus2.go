package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"slices"
	"sort"
	"time"
)

// -------------------- 测试代码 --------------------
func main() {
	for i := 0; i < 100; i++ {
		judge()
	}

}

// CmpFunc 是一个比较函数类型：
// 若 a < b 返回负数，a == b 返回 0，a > b 返回正数。
type CmpFunc[T any] func(a, b T) int

// Ordered 定义支持的有序类型
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// SortedListPlus 试图还原原 C++ sortedlistPlus 的数据结构功能：
type SortedListPlus[T Ordered] struct {
	cmp    CmpFunc[T] // 比较函数
	load1  int        // 对应 C++ 里的 load1
	load2  int        // 对应 C++ 里的 load2
	load1X int        // = 2 * load1
	load2X int        // = 2 * load2
	lg2    uint       // = floor(log2(load2X))
	minVal T          // 用于初始化 blkMax 的最小值

	blocks     [][]T // 分块
	bitCnt     []int // Fenwick 树 / BIT，每个块对应一个计数
	blkMax     []T   // 每个块的最大值，用于块级二分
	segSize    []int // 分段大小信息
	elementCnt int   // 总元素个数
}

// NewSortedListPlus 构造函数，传入比较函数 cmp、以及 load1/load2 两个整形参数。
func NewSortedListPlus[T Ordered](cmp CmpFunc[T], load1, load2 int, minVal T) *SortedListPlus[T] {
	sl := &SortedListPlus[T]{
		cmp:    cmp,
		load1:  load1,
		load2:  load2,
		load1X: load1 * 2,
		load2X: load2 * 2,
		minVal: minVal,
	}
	sl.lg2 = uint(bits.Len(uint(sl.load2X)) - 1) // 等价于 C++ 里的 __lg(load2X)
	sl.clear()
	return sl
}

// NewSortedListPlusFromSlice 从一个切片初始化（类似 C++ 中的带参构造），会先进行排序再分块。
func NewSortedListPlusFromSlice[T Ordered](cmp CmpFunc[T], load1, load2 int, minVal T, arr []T) *SortedListPlus[T] {
	arr = slices.Clone(arr)
	slices.SortFunc(arr, cmp)

	sl := &SortedListPlus[T]{
		cmp:    cmp,
		load1:  load1,
		load2:  load2,
		load1X: load1 * 2,
		load2X: load2 * 2,
		minVal: minVal,
	}
	sl.lg2 = uint(bits.Len(uint(sl.load2X)) - 1)
	sl.elementCnt = len(arr)

	// 分块
	// blocks[0] 只会是空块占位(和原 C++ 代码结构对应)
	sl.blocks = make([][]T, 1, len(arr)/load1+2)

	for start := 0; start < len(arr); start += load1 {
		end := start + load1
		if end > len(arr) {
			end = len(arr)
		}
		chunk := slices.Clone(arr[start:end])
		sl.blocks = append(sl.blocks, chunk)
	}

	// 建立 Fenwick 和 blkMax 等信息
	sl.rebuild()
	return sl
}

// clear 清空所有数据结构
func (sl *SortedListPlus[T]) clear() {
	sl.blocks = sl.blocks[:0]
	sl.blocks = append(sl.blocks, []T{})
	sl.bitCnt = sl.bitCnt[:0]
	sl.bitCnt = append(sl.bitCnt, 0)
	sl.blkMax = sl.blkMax[:0]
	sl.blkMax = append(sl.blkMax, sl.minVal)
	sl.segSize = sl.segSize[:0]
	sl.elementCnt = 0
}

// Size 返回元素总数
func (sl *SortedListPlus[T]) Size() int {
	return sl.elementCnt
}

// chmax 用 cmp 比较更新 a = max(a,b)
func (sl *SortedListPlus[T]) chmax(a *T, b T) {
	if sl.cmp(b, *a) > 0 {
		*a = b
	}
}

// locLeft 返回 (块下标, 块内下标)。
// 寻找第一个块使得 blkMax[块] >= x，然后在该块内二分（用 slices.BinarySearchFunc）。
func (sl *SortedListPlus[T]) locLeft(x T) (int, int) {
	// 若空结构，直接返回 (0, 0) 并让后续自行判断
	if len(sl.segSize) == 0 {
		return 0, 0
	}
	bi := 0
	var maxPow int
	if len(sl.blkMax) <= 1 {
		maxPow = 1
	} else {
		maxPow = 1 << (bits.Len(uint(len(sl.blkMax)-1)) - 1)
	}
	for i := maxPow; i > 0; i >>= 1 {
		next := bi | i
		if next < len(sl.blkMax) && sl.cmp(sl.blkMax[next], x) < 0 {
			bi = next
		}
	}
	if bi+1 < len(sl.blocks) {
		bi++
	}
	// 在第 bi 个块中做 binarySearch
	idx := sl._binarySearchLowerBound(sl.blocks[bi], x)
	return bi, idx
}

// locRight 同理，只是在比较时用 <= x
func (sl *SortedListPlus[T]) locRight(x T) (int, int) {
	if len(sl.segSize) == 0 {
		return 0, 0
	}
	bi := 0
	var maxPow int
	if len(sl.blkMax) <= 1 {
		maxPow = 1
	} else {
		maxPow = 1 << (bits.Len(uint(len(sl.blkMax)-1)) - 1)
	}
	for i := maxPow; i > 0; i >>= 1 {
		next := bi | i
		if next < len(sl.blkMax) && sl.cmp(sl.blkMax[next], x) <= 0 {
			bi = next
		}
	}
	if bi+1 < len(sl.blocks) {
		bi++
	}
	idx := sl._binarySearchUpperBound(sl.blocks[bi], x)
	return bi, idx
}

// Insert 向 SortedListPlus 中插入一个元素 x
func (sl *SortedListPlus[T]) Insert(x T) {
	// 若 segSize 为空，则表示之前是空结构
	if len(sl.segSize) == 0 {
		// blocks[0] 是占位，这里新建块到 blocks[1]
		sl.elementCnt++
		sl.blocks = append(sl.blocks, []T{x})
		sl.bitCnt = append(sl.bitCnt, 1)
		sl.blkMax = append(sl.blkMax, x)
		sl.segSize = append(sl.segSize, 1)
		return
	}
	bi, idx := sl.locRight(x)
	// if idx < len(sl.blocks[bi]) && sl.cmp(sl.blocks[bi][idx], x) == 0 {
	// 	return // 已存在，不插入
	// }
	sl.elementCnt++

	// Fenwick 树计数 +1
	for i := bi; i < len(sl.bitCnt); i += i & -i {
		sl.bitCnt[i]++
	}

	// 在第 bi 个 block 的 idx 位置插入
	sl.blocks[bi] = slices.Insert(sl.blocks[bi], idx, x)

	// 如果插入后块大小 >= load1X，需要分裂
	if len(sl.blocks[bi]) >= sl.load1X {
		sl._splitBlock(bi)
	} else {
		// 否则更新块最大值
		lastVal := sl.blocks[bi][len(sl.blocks[bi])-1]
		sl.blkMax[bi] = lastVal
	}
}

// Erase 从 SortedListPlus 中删除一个元素 x
func (sl *SortedListPlus[T]) Erase(x T) bool {
	if len(sl.segSize) == 0 {
		return false
	}
	bi, idx := sl.locLeft(x)
	if idx >= len(sl.blocks[bi]) || sl.cmp(sl.blocks[bi][idx], x) != 0 {
		return false
	}

	// Fenwick 树计数 -1
	for i := bi; i < len(sl.bitCnt); i += i & -i {
		sl.bitCnt[i]--
	}
	sl.elementCnt--

	// 删除元素
	sl.blocks[bi] = slices.Delete(sl.blocks[bi], idx, idx+1)
	if len(sl.blocks[bi]) == 0 {
		// 该块空了，需要收缩
		sl._shrinkBlock(bi)
	} else {
		// 更新最大值
		lastVal := sl.blocks[bi][len(sl.blocks[bi])-1]
		sl.blkMax[bi] = lastVal
	}
	return true
}

// At 根据 rank（0-based）找到元素
// 类似 operator[](k)，要求 0 <= k < elementCnt
func (sl *SortedListPlus[T]) At(k int) T {
	if k < 0 || k >= sl.elementCnt {
		panic("At: out of range")
	}
	bi := 0
	maxPow := 1 << (bits.Len(uint(len(sl.bitCnt)-1)) - 1)
	if len(sl.bitCnt) <= 1 {
		maxPow = 1
	}
	for i := maxPow; i > 0; i >>= 1 {
		next := bi | i
		if next < len(sl.bitCnt) && sl.bitCnt[next] <= k {
			k -= sl.bitCnt[next]
			bi = next
		}
	}
	bi++
	return sl.blocks[bi][k]
}

// Rank 返回 < x 的元素个数
func (sl *SortedListPlus[T]) Rank(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, idx := sl.locLeft(x)
	// 先计算所有 < bi 的块的总大小
	rk := 0
	i := bi - 1
	for i > 0 {
		rk += sl.bitCnt[i]
		i &= i - 1
	}
	// 再加上该块内的下标 idx
	rk += idx
	return rk
}

// Rank2 返回 <= x 的元素个数
func (sl *SortedListPlus[T]) Rank2(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, idx := sl.locRight(x)
	rk := 0
	i := bi - 1
	for i > 0 {
		rk += sl.bitCnt[i]
		i &= i - 1
	}
	rk += idx
	return rk
}

// Count 返回等于 x 的元素个数
func (sl *SortedListPlus[T]) Count(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	return sl.Rank2(x) - sl.Rank(x)
}

// LessEqual 找到 <= x 的最大值
func (sl *SortedListPlus[T]) LessEqual(x T) (T, bool) {
	bi, idx := sl.locRight(x)
	if !(bi > 1 || (bi == 1 && idx > 0)) {
		return sl.minVal, false
	}
	bi_prev, idx_prev := sl.prev(bi, idx)
	return sl.blocks[bi_prev][idx_prev], true
}

// GreaterEqual 找到 >= x 的最小值
func (sl *SortedListPlus[T]) GreaterEqual(x T) (T, bool) {
	bi, idx := sl.locLeft(x)
	if bi >= len(sl.blocks) || idx >= len(sl.blocks[bi]) {
		return sl.minVal, false
	}
	return sl.blocks[bi][idx], true
}

// prev 返回前一个元素的位置 (bi, idx)
func (sl *SortedListPlus[T]) prev(bi, idx int) (int, int) {
	if bi == 1 && idx == 0 {
		panic("prev out of range")
	}
	if idx > 0 {
		return bi, idx - 1
	} else {
		bi--
		return bi, len(sl.blocks[bi]) - 1
	}
}

// rebuild 重新分块与 Fenwick 树
func (sl *SortedListPlus[T]) rebuild() {
	// 先把原 blocks 移走
	oldBlocks := sl.blocks
	// 统计多少个非空块
	c := 0
	for _, b := range oldBlocks {
		if len(b) > 0 {
			c++
		}
	}

	// 估算要分成多少段
	segn := (c + sl.load2 - 1) / sl.load2

	ec := sl.elementCnt // 记录原本的元素数量

	sl.clear()
	sl.segSize = make([]int, segn)
	sl.elementCnt = ec

	curSeg := 0
	countInSeg := 0
	for _, block := range oldBlocks {
		if len(block) == 0 {
			continue
		}
		sl.blocks = append(sl.blocks, block)
		sl.bitCnt = append(sl.bitCnt, len(block))
		// 更新 blkMax
		sl.blkMax = append(sl.blkMax, block[len(block)-1])
		countInSeg++
		if countInSeg == sl.load2 {
			sl.segSize[curSeg] = countInSeg
			curSeg++
			countInSeg = 0
		}
	}
	if countInSeg > 0 && curSeg < len(sl.segSize) {
		sl.segSize[curSeg] = countInSeg
	}

	for i := 1; i < len(sl.bitCnt); i++ {
		p := i + (i & -i)
		if p < len(sl.bitCnt) {
			sl.bitCnt[p] += sl.bitCnt[i]
		}
	}
}

// _rangeBitModify 初始化 bitCnt 为空的区域
func (sl *SortedListPlus[T]) _rangeBitModify(b1, b2 int) {
	// 初始化 bitCnt 为空的区域
	for i := b1; i <= b2; i++ {
		sl.bitCnt[i] = 0
	}

	// 设置最小值
	mx := sl.minVal
	for i := b1; i <= b2; i++ {
		if len(sl.blocks[i]) > 0 {
			sl.bitCnt[i] += len(sl.blocks[i])
			mx = sl.blocks[i][len(sl.blocks[i])-1] // blkMax 更新为块的最大值
			sl.blkMax[i] = mx
		} else {
			sl.blkMax[i] = mx
		}

		// 更新区间的 BIT
		if (i+(-i&i)) <= b2 && sl.bitCnt[i] > 0 {
			sl.bitCnt[i+(-i&i)] += sl.bitCnt[i]
		}
	}

	// 第二轮更新
	for lowb := (-b2 & b2) / 2; lowb >= sl.load2X; lowb >>= 1 {
		sl.bitCnt[b2] += sl.bitCnt[b2-lowb]
	}
}

// _splitBlock 分裂块：当某个块大小 >= load1X 时，将该块拆分出 load1 大小保持，剩余部分新建块。
func (sl *SortedListPlus[T]) _splitBlock(bi int) {
	block := sl.blocks[bi]
	if len(block) < sl.load1X {
		return
	}
	// 需要拆分出 load1 大小
	newBlock := slices.Clone(block[sl.load1:])
	sl.blocks[bi] = block[:sl.load1]

	// 将新块插入到 blocks[bi+1] 位置
	pos := bi + 1
	sl.blocks = slices.Insert(sl.blocks, pos, newBlock)
	sl.bitCnt = slices.Insert(sl.bitCnt, pos, len(newBlock))
	// 原块 bitCnt[bi] 改为 load1
	sl.bitCnt[bi] = len(sl.blocks[bi])

	// blkMax 同理插入
	sl.blkMax = slices.Insert(sl.blkMax, pos, newBlock[len(newBlock)-1])
	sl.blkMax[bi] = sl.blocks[bi][len(sl.blocks[bi])-1]

	// 更新 segSize
	// 计算该块属于哪一段
	segi := (bi - 1) >> (sl.lg2)
	if segi < len(sl.segSize) {
		sl.segSize[segi]++
		if sl.segSize[segi] == sl.load2X {
			// 到达阈值后执行 _expand
			sl.rebuild()
		} else {
			// 局部修正 Fenwick
			sl._rangeBitModify(segi<<sl.lg2|1, len(sl.blocks)-1)
		}
	}
}

// _shrinkBlock 当某个块被删空时，需要合并/回收
func (sl *SortedListPlus[T]) _shrinkBlock(bi int) {
	// 如果是最后一个块，直接弹掉
	if bi == len(sl.blocks)-1 {
		sl.blocks = sl.blocks[:bi]
		sl.blkMax = sl.blkMax[:bi]
		sl.bitCnt = sl.bitCnt[:bi]
		// 更新 segSize
		segi := (bi - 1) >> sl.lg2
		if segi >= 0 && segi < len(sl.segSize) {
			sl.segSize[segi]--
			if sl.segSize[segi] == 0 {
				sl.segSize = sl.segSize[:len(sl.segSize)-1]
			}
		}
		for len(sl.blocks) > 1 && len(sl.blocks[len(sl.blocks)-1]) == 0 {
			sl.blocks = sl.blocks[:len(sl.blocks)-1]
			sl.blkMax = sl.blkMax[:len(sl.blkMax)-1]
			sl.bitCnt = sl.bitCnt[:len(sl.bitCnt)-1]
		}
		return
	}

	// 否则，如果不是最后一个块：需要把这块删掉
	sl.blocks = slices.Delete(sl.blocks, bi, bi+1)
	sl.blkMax = slices.Delete(sl.blkMax, bi, bi+1)
	sl.bitCnt = slices.Delete(sl.bitCnt, bi, bi+1)

	segi := (bi - 1) >> sl.lg2
	if segi < 0 {
		segi = 0
	}
	if segi >= len(sl.segSize) {
		sl.rebuild()
		return
	}
	sl.segSize[segi]--
	if sl.segSize[segi] == 0 {
		sl.rebuild()
	} else {
		sl._rangeBitModify(segi<<sl.lg2|1, len(sl.blocks)-1)
	}
}

// _binarySearchLowerBound 在有序切片 arr 中找到第一个 >= x 的位置
func (sl *SortedListPlus[T]) _binarySearchLowerBound(arr []T, x T) int {
	idx, ok := slices.BinarySearchFunc(arr, x, func(a, b T) int {
		return sl.cmp(a, b)
	})
	if ok {
		return idx
	}
	return idx
}

// _binarySearchUpperBound 在有序切片 arr 中找到第一个 > x 的位置
func (sl *SortedListPlus[T]) _binarySearchUpperBound(arr []T, x T) int {
	idx, ok := slices.BinarySearchFunc(arr, x, func(a, b T) int {
		res := sl.cmp(a, b)
		if res <= 0 {
			return -1
		}
		return 1
	})
	if ok {
		return idx + 1
	}
	return idx
}

// String 打印调试用
func (sl *SortedListPlus[T]) String() string {
	return fmt.Sprintf(
		"SortedListPlus{elementCnt=%d, blocks=%v}",
		sl.elementCnt, sl.blocks,
	)
}

// -------------------- 朴素对拍结构 --------------------
type naiveSL[T Ordered] struct {
	data []T
}

func NewNaiseSlFrom[T Ordered](sorted []T) *naiveSL[T] {
	sorted = slices.Clone(sorted)
	return &naiveSL[T]{data: sorted}
}

// 若已存在则不再插入；否则插入保持有序
func (n *naiveSL[T]) Insert(x T) {
	pos, _ := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	// if ok {
	// 	return
	// }
	n.data = slices.Insert(n.data, pos, x)
}

func (n *naiveSL[T]) Erase(x T) {
	pos, _ := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	if pos < len(n.data) && n.data[pos] == x {
		n.data = slices.Delete(n.data, pos, pos+1)
	}
}

func (n *naiveSL[T]) Size() int { return len(n.data) }

func (n *naiveSL[T]) BisectLeft(x T) int {
	pos, _ := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	return pos
}

func (n *naiveSL[T]) BisectRight(x T) int {
	pos, _ := slices.BinarySearchFunc(n.data, x, func(a, b T) int {
		if a <= b {
			return -1
		}
		return +1
	})
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

func (n *naiveSL[T]) Floor(x T) (res T, ok bool) {
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

func (n *naiveSL[T]) Ceiling(x T) (res T, ok bool) {
	pos, _ := slices.BinarySearchFunc(n.data, x, cmpOrder[T])
	if pos == len(n.data) {
		return res, false
	}
	return n.data[pos], true
}

func (n *naiveSL[T]) At(k int) T {
	if k < 0 || k >= len(n.data) {
		panic("At: out of range")
	}
	return n.data[k]
}

// cmpOrder 用于比较函数
func cmpOrder[T Ordered](a, b T) int {
	if a < b {
		return -1
	} else if a > b {
		return +1
	}
	return 0
}

// judge: 在此函数中进行随机测试，对拍 SortedListPlus vs naiveSL
func judge() {
	const N = 50000           // 操作次数
	const rangeVal = int(1e9) // 值域
	rand.Seed(time.Now().UnixNano())

	Q := rand.Intn(2e5) + 1
	// 先产生 Q 个随机数，排好序，放入 SortedListPlus & naiveSL
	nums := make([]int, Q)
	for i := 0; i < Q; i++ {
		nums[i] = rand.Intn(rangeVal)
	}
	sort.Ints(nums)

	cmp := func(a, b int) int { return a - b }

	slp := NewSortedListPlusFromSlice(cmp, 200, 64, 1234232323230, nums) // SortedListPlus
	naive := NewNaiseSlFrom(nums)                                        // 朴素做法

	for i := 1; i <= N; i++ {
		op := rand.Intn(9) // 0..8
		x := rand.Intn(rangeVal)
		switch op {

		case 0: // Insert
			// 如需测试插入，可解注释
			slp.Insert(x)
			naive.Insert(x)

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
			rnaive := naive.BisectLeft(x)
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
			rnaive, ok2 := naive.Floor(x)
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
			rnaive, ok2 := naive.Ceiling(x)
			if ok1 != ok2 {
				panic(fmt.Sprintf("Mismatch in GreaterEqual(%d) existence, slpOk=%v naiveOk=%v",
					x, ok1, ok2))
			}
			if ok1 && rslp != rnaive {
				panic(fmt.Sprintf("Mismatch in GreaterEqual(%d), got slp=%d naive=%d",
					x, rslp, rnaive))
			}
		case 7: // bisectRight
			rslp := slp.Rank2(x)
			rnaive := naive.BisectRight(x)
			if rslp != rnaive {
				panic(fmt.Sprintf("Mismatch in bisectRight(%d), got slp=%d naive=%d",
					x, rslp, rnaive))
			}

		case 8: // At
			if slp.Size() > 0 {
				k := rand.Intn(slp.Size())
				rslp := slp.At(k)
				rnaive := naive.At(k)
				if rslp != rnaive {
					panic(fmt.Sprintf("Mismatch in At(%d), got slp=%d naive=%d",
						k, rslp, rnaive))
				}
			}
		}

	}

	fmt.Println("testTime: All tests passed! SortedListPlus and naiveSL match on random data.")
}
