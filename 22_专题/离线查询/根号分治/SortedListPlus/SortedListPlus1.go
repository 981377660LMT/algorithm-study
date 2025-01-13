package main

import (
	"bufio"

	"fmt"
	"os"
	"sort"
	"unsafe"
)

func main() {
	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	NextByte := func() byte {
		b := rc()
		for ; '0' > b; b = rc() {
		}
		return b
	}
	_ = NextByte

	// 读一个整数，支持负数
	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt

	n, q := NextInt(), NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = NextInt()
	}
	sort.Ints(nums)
	sl := NewSortedListPlusFrom(nums)
	for i := 0; i < q; i++ {
		op, x := NextInt(), NextInt()
		switch op {
		case 0:
			sl.Insert(x)
		case 1:
			sl.Erase(x)
		case 2:
			if x > sl.Size() {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, sl.At(x-1))
			}
		case 3:
			fmt.Fprintln(out, sl.BisectRight(x))
		case 4:
			if res, ok := sl.Floor(x); ok {
				fmt.Fprintln(out, res)
			} else {
				fmt.Fprintln(out, -1)
			}
		case 5:
			if res, ok := sl.Ceiling(x); ok {
				fmt.Fprintln(out, res)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

type SortedListPlus[T Ordered] struct {
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

func NewSortedListPlusFrom[T Ordered](sorted []T) *SortedListPlus[T] {

	// TODO: sort
	sl := &SortedListPlus[T]{
		load1:  200,
		load2:  64,
		minVal: getMinValue[T](),
	}
	sl.load1X = sl.load1 * 2
	sl.load2X = sl.load2 * 2
	sl.lg2 = intLog2(sl.load2X)

	sl.elemCnt = len(sorted)
	blockCount := (len(sorted) + int(sl.load1) - 1) / int(sl.load1)
	sl.blocks = make([][]T, 0, blockCount+1)
	sl.blocks = append(sl.blocks, nil)

	for i := 0; i < len(sorted); i += int(sl.load1) {
		end := i + int(sl.load1)
		if end > len(sorted) {
			end = len(sorted)
		}
		sl.blocks = append(sl.blocks, Clone(sorted[i:end]))
	}
	sl._expand()
	return sl
}

// NewEmptySortedListPlus 构造一个空的 SortedListPlus
func NewEmptySortedListPlus[T Ordered]() *SortedListPlus[T] {
	sl := &SortedListPlus[T]{
		load1:  200,
		load2:  64,
		minVal: getMinValue[T](),
	}
	sl.load1X = sl.load1 * 2
	sl.load2X = sl.load2 * 2
	sl.lg2 = intLog2(sl.load2X)
	sl.Clear()
	return sl
}

// Clear 清空数据结构
func (sl *SortedListPlus[T]) Clear() {
	sl.blocks = sl.blocks[:0]
	sl.bitCnt = sl.bitCnt[:0]
	sl.blkMax = sl.blkMax[:0]
	sl.segMax = sl.segMax[:0]
	sl.segSize = sl.segSize[:0]
	sl.elemCnt = 0
}

// Size 返回当前结构里的元素个数
func (sl *SortedListPlus[T]) Size() int {
	return sl.elemCnt
}

// Insert 插入元素 x
func (sl *SortedListPlus[T]) Insert(x T) {
	// 如果 segSize 为空 => 整个结构空
	if len(sl.segSize) == 0 {
		sl.blocks = append(sl.blocks, []T{x})
		sl.bitCnt = append(sl.bitCnt, 1)
		sl.blkMax = append(sl.blkMax, x)
		sl.segMax = append(sl.segMax, x)
		sl.segSize = append(sl.segSize, 1)
		sl.elemCnt = 1
		return
	}

	bi, pos := sl.lowerBound(x) // 找到要插入的块及在块内的位置
	// do nothing if x already exists
	if pos < len(sl.blocks[bi]) && sl.blocks[bi][pos] == x {
		return
	}
	sl.elemCnt++

	// Fenwick 树更新
	sl.fenwickUpdate(bi, 1)

	// 在块内插入
	sl.blocks[bi] = Insert(sl.blocks[bi], pos, x)

	// 更新块最大值
	if x > sl.blkMax[bi] {
		sl.blkMax[bi] = x
	}

	// 更新 segMax
	segi := bi >> sl.lg2
	if x > sl.segMax[segi] {
		sl.segMax[segi] = x
	}

	// 如果块大小已达阈值 => 分裂
	if len(sl.blocks[bi]) >= sl.load1X {
		sl.splitBlock(segi, bi)
	}
}

// Erase 删除值为 x 的元素（若不存在则忽略）
func (sl *SortedListPlus[T]) Erase(x T) {
	if len(sl.segMax) == 0 {
		return
	}
	// 如果全局最大值 < x，直接无
	if sl.segMax[len(sl.segMax)-1] < x {
		return
	}
	bi, pos := sl.lowerBound(x)
	// 如果插入点已经到块尾 || 找到的值 != x => 不存在
	if pos == len(sl.blocks[bi]) || sl.blocks[bi][pos] != x {
		return
	}

	// Fenwick 树更新
	sl.fenwickUpdate(bi, -1)
	sl.elemCnt--

	// 块内删除
	sl.blocks[bi] = Delete(sl.blocks[bi], pos, pos+1)

	// 若块变空，删除这个块
	if len(sl.blocks[bi]) == 0 {
		sl.eraseBlock(bi)
	} else {
		// 否则更新块最大值
		sl.blkMax[bi] = sl.blocks[bi][len(sl.blocks[bi])-1]
		// 如果这个块正好是这个 segment 的最后一个非空块，需要更新 segMax
		segi := bi >> sl.lg2
		bj := (segi << sl.lg2) + sl.segSize[segi] - 1 // 最后一个非空块的下标
		if bi == bj {
			sl.segMax[segi] = sl.blkMax[bi]
		}
	}
}

// At 返回第 k (0-based) 小的元素
func (sl *SortedListPlus[T]) At(k int) T {
	if k < 0 || k >= sl.elemCnt {
		panic("At: out of range")
	}
	// 通过 Fenwick 找到所在块
	blkIndex := sl.findFenwick(k)
	// 块内下标 = k - 前面块的元素数量和
	offset := k - sl.prefixSum(blkIndex-1)
	return sl.blocks[blkIndex][offset]
}

// BisectLeft(x) = lower_bound(x) 在整个序列中的下标（即有多少元素 < x）
func (sl *SortedListPlus[T]) BisectLeft(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, pos := sl.lowerBound(x)
	return sl.prefixSum(bi-1) + pos
}

// BisectRight(x) = upper_bound(x) 在整个序列中的下标（即有多少元素 <= x）
func (sl *SortedListPlus[T]) BisectRight(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, pos := sl.upperBound(x)
	return sl.prefixSum(bi-1) + pos
}

// Count(x) = upper_bound(x) - lower_bound(x)
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

// Floor(x) 返回最大 <= x 的元素，不存在则 (零值, false)
func (sl *SortedListPlus[T]) Floor(x T) (res T, ok bool) {
	if len(sl.segSize) == 0 {
		return
	}
	bi, pos := sl.upperBound(x)
	if bi == 0 && pos == 0 {
		return res, false
	}
	// 前驱
	pb, pe := sl.prev(bi, pos)
	return sl.blocks[pb][pe], true
}

// Ceiling(x) 返回最小 >= x 的元素，不存在则 (零值, false)
func (sl *SortedListPlus[T]) Ceiling(x T) (res T, ok bool) {
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

// lowerBound(x) -> (块编号, 块内下标)，使 blocks[bi][pos] 为第一个 >= x
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

	// 2. 小段内非空块之间二分 (连续下标)
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

	// 3. 块内二分
	pos, _ = BinarySearchFunc(sl.blocks[bi], x, func(a, b T) int {
		if a < b {
			return -1
		} else if a > b {
			return +1
		}
		return 0
	})
	return
}

// upperBound(x) -> (块编号, 块内下标)，使 blocks[bi][pos] 为第一个 > x
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

	// 2. 小段内非空块之间二分
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
	pos, _ = BinarySearchFunc(sl.blocks[bi], x, func(a, b T) int {
		if a <= b {
			return -1
		}
		return +1
	})
	return
}

// prev(bi, pos) => 同一块前移 或 跑到前一块
func (sl *SortedListPlus[T]) prev(bi, pos int) (int, int) {
	if pos > 0 {
		return bi, pos - 1
	}
	if bi == 0 {
		return 0, 0
	}
	return bi - 1, len(sl.blocks[bi-1]) - 1
}

// splitBlock 当块大小 >= load1X 时分裂
func (sl *SortedListPlus[T]) splitBlock(segi, bi int) {
	oldBlk := sl.blocks[bi]
	newBlk := Clone(oldBlk[sl.load1:])
	sl.blocks[bi] = oldBlk[:sl.load1]

	biNew := bi + 1
	sl.blocks = Insert(sl.blocks, biNew, []T{})
	sl.blocks[biNew] = newBlk

	sl.blkMax = Insert(sl.blkMax, biNew, newBlk[len(newBlk)-1])
	sl.bitCnt = Insert(sl.bitCnt, biNew, len(newBlk))

	sl.fenwicksumRebuild()

	sl.segSize[segi]++
	if sl.segSize[segi] == sl.load2X {
		sl._expand()
	}
}

// eraseBlock 当某块为空
func (sl *SortedListPlus[T]) eraseBlock(bi int) {
	sl.blocks = Delete(sl.blocks, bi, bi+1)
	sl.blkMax = Delete(sl.blkMax, bi, bi+1)
	sl.bitCnt = Delete(sl.bitCnt, bi, bi+1)

	segi := bi >> sl.lg2
	sl.segSize[segi]--
	if sl.segSize[segi] == 0 {
		// 若正好是最后一小段
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

// _expand 全局重构 => 彻底保证所有非空块紧密存储，不插空块
func (sl *SortedListPlus[T]) _expand() {
	oldBlocks := sl.blocks
	if len(oldBlocks) == 0 {
		sl.Clear()
		return
	}
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
	segn := (c + sl.load2 - 1) / sl.load2

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
		segi := j >> sl.lg2
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

// fenwicksumRebuild => 重新构建 Fenwick
func (sl *SortedListPlus[T]) fenwicksumRebuild() {
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

// fenwickUpdate => 增量更新
func (sl *SortedListPlus[T]) fenwickUpdate(idx, delta int) {
	for idx < len(sl.bitCnt) {
		sl.bitCnt[idx] += delta
		idx |= (idx + 1)
	}
}

// prefixSum(idx) => bitCnt[0..idx]之和
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

// findFenwick(k) => 找前缀和 >= k+1 的块下标
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
func getMinValue[T Ordered]() T {
	var zero T
	return zero
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (sl *SortedListPlus[T]) DebugPrint() {
	fmt.Printf("Blocks:\n")
	for i := 1; i < len(sl.blocks); i++ {
		fmt.Printf("  i=%d, blk=%v\n", i, sl.blocks[i])
	}
	fmt.Printf("bitCnt=%v\n", sl.bitCnt)
	fmt.Printf("blkMax=%v\n", sl.blkMax)
	fmt.Printf("segMax=%v\n", sl.segMax)
	fmt.Printf("segSize=%v\n", sl.segSize)
	fmt.Printf("elemCnt=%d\n", sl.elemCnt)
	fmt.Println("--------------------------------")
}

// Insert inserts the values v... into s at index i,
// !returning the modified slice.
// The elements at s[i:] are shifted up to make room.
// In the returned slice r, r[i] == v[0],
// and r[i+len(v)] == value originally at r[i].
// This function is O(len(s) + len(v)).
func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	if i < 0 {
		i = 0
	}
	if i > len(s) {
		i = len(s)
	}

	m := len(v)
	if m == 0 {
		return s
	}
	n := len(s)
	if i == n {
		return append(s, v...)
	}
	if n+m > cap(s) {
		s2 := append(s[:i], make(S, n+m-i)...)
		copy(s2[i:], v)
		copy(s2[i+m:], s[i:])
		return s2
	}
	s = s[:n+m]
	if !overlaps(v, s[i+m:]) {
		copy(s[i+m:], s[i:])
		copy(s[i:], v)
		return s
	}
	copy(s[n:], v)
	rotateRight(s[i:], m)
	return s
}

func overlaps[E any](a, b []E) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	elemSize := unsafe.Sizeof(a[0])
	if elemSize == 0 {
		return false
	}
	return uintptr(unsafe.Pointer(&a[0])) <= uintptr(unsafe.Pointer(&b[len(b)-1]))+(elemSize-1) &&
		uintptr(unsafe.Pointer(&b[0])) <= uintptr(unsafe.Pointer(&a[len(a)-1]))+(elemSize-1)
}

func rotateLeft[E any](s []E, r int) {
	for r != 0 && r != len(s) {
		if r*2 <= len(s) {
			swap(s[:r], s[len(s)-r:])
			s = s[:len(s)-r]
		} else {
			swap(s[:len(s)-r], s[r:])
			s, r = s[len(s)-r:], r*2-len(s)
		}
	}
}

func rotateRight[E any](s []E, r int) {
	rotateLeft(s, len(s)-r)
}

func swap[E any](x, y []E) {
	for i := 0; i < len(x); i++ {
		x[i], y[i] = y[i], x[i]
	}
}

// Shallow clone.
func Clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
}

func Delete[S ~[]E, E any](s S, i, j int) S {
	if i < 0 {
		i = 0
	}
	if j > len(s) {
		j = len(s)
	}
	if i >= j {
		return s
	}

	s = append(s[:i], s[j:]...)
	return s
}

func DeleteFunc[S ~[]E, E any](s S, del func(E) bool) S {
	i := func() int {
		for i, v := range s {
			if del(v) {
				return i
			}
		}
		return -1
	}()
	if i == -1 {
		return s
	}
	for j := i + 1; j < len(s); j++ {
		if v := s[j]; !del(v) {
			s[i] = v
			i++
		}
	}
	return s[:i]
}

// BinarySearchFunc works like [BinarySearch], but uses a custom comparison
// function. The slice must be sorted in increasing order, where "increasing"
// is defined by  cmp should return 0 if the slice element matches
// the target, a negative number if the slice element precedes the target,
// or a positive number if the slice element follows the target.
// cmp must implement the same ordering as the slice, such that if
// cmp(a, t) < 0 and cmp(b, t) >= 0, then a must precede b in the slice.
func BinarySearchFunc[S ~[]E, E, T any](x S, target T, cmp func(E, T) int) (int, bool) {
	n := len(x)
	// Define cmp(x[-1], target) < 0 and cmp(x[n], target) >= 0 .
	// Invariant: cmp(x[i - 1], target) < 0, cmp(x[j], target) >= 0.
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if cmp(x[h], target) < 0 {
			i = h + 1 // preserves cmp(x[i - 1], target) < 0
		} else {
			j = h // preserves cmp(x[j], target) >= 0
		}
	}
	// i == j, cmp(x[i-1], target) < 0, and cmp(x[j], target) (= cmp(x[i], target)) >= 0  =>  answer is i.
	return i, i < n && cmp(x[i], target) == 0
}

// #region

// #endregion
