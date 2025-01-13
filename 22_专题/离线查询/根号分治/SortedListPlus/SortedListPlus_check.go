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
	sl := NewSortedListPlus(nums)
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
			if res, ok := sl.LessEqual(x); ok {
				fmt.Fprintln(out, res)
			} else {
				fmt.Fprintln(out, -1)
			}
		case 5:
			if res, ok := sl.GreaterEqual(x); ok {
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
	load1   uint
	load2   uint
	load1X  uint
	load2X  uint
	lg2     uint
	minVal  T
	blocks  [][]T
	bitCnt  []int
	blkMax  []T
	segMax  []T
	segSize []int
	elemCnt int
}

func NewSortedListPlus[T Ordered](sorted []T) *SortedListPlus[T] {
	// TODO: sort

	sl := &SortedListPlus[T]{
		load1:  200,
		load2:  64,
		minVal: getMinValue[T](),
	}
	sl.load1X = sl.load1 * 2
	sl.load2X = sl.load2 * 2
	sl.lg2 = log2(sl.load2X)

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

func NewEmptySortedListPlus[T Ordered]() *SortedListPlus[T] {
	sl := &SortedListPlus[T]{
		load1:  200,
		load2:  64,
		minVal: getMinValue[T](),
	}
	sl.load1X = sl.load1 * 2
	sl.load2X = sl.load2 * 2
	sl.lg2 = log2(sl.load2X)
	sl.Clear()
	return sl
}

func (sl *SortedListPlus[T]) Clear() {
	sl.blocks = make([][]T, 1)
	sl.bitCnt = make([]int, 1)
	sl.blkMax = make([]T, 1)
	sl.segMax = sl.segMax[:0]
	sl.segSize = sl.segSize[:0]
	sl.elemCnt = 0
}

func (sl *SortedListPlus[T]) Insert(x T) {
	if len(sl.segSize) == 0 {
		sl.elemCnt++
		sl.blocks = append(sl.blocks, []T{x})
		sl.bitCnt = append(sl.bitCnt, 1)
		sl.blkMax = append(sl.blkMax, x)
		sl.segMax = append(sl.segMax, x)
		sl.segSize = append(sl.segSize, 1)
		return
	}

	bi, ei := sl.lowerBound(x)
	// If x is already in the list, do nothing.
	if ei < len(sl.blocks[bi]) && sl.blocks[bi][ei] == x {
		return
	}
	sl.elemCnt++

	for i := bi; i < len(sl.bitCnt); i += i & -i {
		sl.bitCnt[i]++
	}

	sl.blocks[bi] = Insert(sl.blocks[bi], ei, x)
	if len(sl.blocks[bi]) > 0 {
		if x > sl.blkMax[bi] {
			sl.blkMax[bi] = x
		}
	}
	segi := (bi - 1) >> sl.lg2
	if x > sl.segMax[segi] {
		sl.segMax[segi] = x
	}
	if len(sl.blocks[bi]) >= int(sl.load1X) {
		sl.splitBlock(segi, bi)
	}
}

func (sl *SortedListPlus[T]) Erase(x T) {
	if len(sl.segMax) == 0 {
		return
	}
	if sl.segMax[len(sl.segMax)-1] < x {
		return
	}
	bi, ei := sl.lowerBound(x)
	if bi == 0 || ei == len(sl.blocks[bi]) {
		return
	}
	if sl.blocks[bi][ei] != x {
		return
	}
	sl.elemCnt--
	for i := bi; i < len(sl.bitCnt); i += i & -i {
		sl.bitCnt[i]--
	}
	sl.blocks[bi] = Delete(sl.blocks[bi], ei, ei+1)

	if len(sl.blocks[bi]) == 0 {
		sl.eraseBlock(bi)
	} else {
		sl.blkMax[bi] = sl.blocks[bi][len(sl.blocks[bi])-1]
		segi := (bi - 1) >> sl.lg2
		bj := (segi << sl.lg2) | sl.segSize[segi]
		if bi == bj {
			sl.segMax[segi] = sl.blkMax[bi]
		}
	}
}

func (sl *SortedListPlus[T]) Size() int {
	return sl.elemCnt
}

func (sl *SortedListPlus[T]) At(k int) T {
	if k < 0 || k >= sl.elemCnt {
		panic("At: index out of range")
	}
	bi := 0
	i := 1
	for i < len(sl.bitCnt) {
		i <<= 1
	}
	i >>= 1
	for ; i > 0; i >>= 1 {
		if (bi|i) < len(sl.bitCnt) && sl.bitCnt[bi|i] <= k {
			k -= sl.bitCnt[bi|i]
			bi |= i
		}
	}
	blockIdx := bi + 1
	return sl.blocks[blockIdx][k]
}

func (sl *SortedListPlus[T]) BisectLeft(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, ei := sl.lowerBound(x)
	rk := ei
	for i := bi - 1; i > 0; i -= i & -i {
		rk += sl.bitCnt[i]
	}
	return rk
}

func (sl *SortedListPlus[T]) BisectRight(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, ei := sl.upperBound(x)
	rk := ei
	for i := bi - 1; i > 0; i -= i & -i {
		rk += sl.bitCnt[i]
	}
	return rk
}

func (sl *SortedListPlus[T]) Count(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	biU, eiU := sl.upperBound(x)
	biL, eiL := sl.lowerBound(x)
	rkU := eiU
	for i := biU - 1; i > 0; i -= i & -i {
		rkU += sl.bitCnt[i]
	}
	rkL := eiL
	for i := biL - 1; i > 0; i -= i & -i {
		rkL += sl.bitCnt[i]
	}
	return rkU - rkL
}

func (sl *SortedListPlus[T]) LessEqual(x T) (res T, ok bool) {
	bi, ei := sl.upperBound(x)
	if bi == 0 && ei == 0 {
		return res, false
	}
	prevBi, prevEi := sl.prev(bi, ei)
	return sl.blocks[prevBi][prevEi], true
}

func (sl *SortedListPlus[T]) GreaterEqual(x T) (res T, ok bool) {
	bi, ei := sl.lowerBound(x)
	if bi == 0 || ei == len(sl.blocks[bi]) {
		return res, false
	}
	return sl.blocks[bi][ei], true
}

func (sl *SortedListPlus[T]) lowerBound(x T) (int, int) {
	if len(sl.segSize) == 0 {
		return 0, 0
	}
	l, r := -1, len(sl.segSize)-1
	for r-l > 1 {
		mid := (l + r) >> 1
		if sl.segMax[mid] >= x {
			r = mid
		} else {
			l = mid
		}
	}
	segi := r

	l = segi<<sl.lg2 | 1
	r = segi<<sl.lg2 | sl.segSize[segi]
	for r-l > 1 {
		mid := (l + r) >> 1
		if sl.blkMax[mid] >= x {
			r = mid
		} else {
			l = mid
		}
	}
	bi := r
	ei, _ := BinarySearchFunc(sl.blocks[bi], x, func(a, b T) int {
		if a < b {
			return -1
		} else if a > b {
			return +1
		}
		return 0
	})
	return bi, ei
}

func (sl *SortedListPlus[T]) upperBound(x T) (int, int) {
	if len(sl.segSize) == 0 {
		return 0, 0
	}
	l, r := -1, len(sl.segSize)-1
	for r-l > 1 {
		mid := (l + r) >> 1
		if sl.segMax[mid] > x {
			r = mid
		} else {
			l = mid
		}
	}
	segi := r

	l = segi<<sl.lg2 | 1
	r = segi<<sl.lg2 | sl.segSize[segi]
	for r-l > 1 {
		mid := (l + r) >> 1
		if sl.blkMax[mid] > x {
			r = mid
		} else {
			l = mid
		}
	}
	bi := r
	ei, _ := BinarySearchFunc(sl.blocks[bi], x, func(a, b T) int {
		if a <= b {
			return -1
		} else {
			return +1
		}
	})
	return bi, ei
}

func (sl *SortedListPlus[T]) prev(bi, ei int) (int, int) {
	if ei > 0 {
		return bi, ei - 1
	}
	if (bi & (int(sl.load2X) - 1)) != 1 {
		return bi - 1, len(sl.blocks[bi-1]) - 1
	}
	segi := (bi - 1) >> sl.lg2
	if segi == 0 {
		return 1, 0
	}
	newBi := ((segi - 1) << sl.lg2) | sl.segSize[segi-1]
	return newBi, len(sl.blocks[newBi]) - 1
}

func (sl *SortedListPlus[T]) splitBlock(segi, bi int) {
	oldBlock := sl.blocks[bi]
	newBlock := Clone(oldBlock[sl.load1:])
	sl.blocks[bi] = oldBlock[:sl.load1]

	biPlus1 := bi + 1
	emptySlice := make([][]T, 1)
	sl.blocks = Insert(sl.blocks, biPlus1, emptySlice...)
	sl.blocks[biPlus1] = newBlock

	sl.blkMax = Insert(sl.blkMax, biPlus1, newBlock[len(newBlock)-1])
	sl.bitCnt = Insert(sl.bitCnt, biPlus1, len(newBlock))

	sl.segSize[segi]++
	sl._rangeBitModify(segi<<sl.lg2|1, minInt(((segi+1)<<sl.lg2), len(sl.bitCnt)-1))

	if sl.segSize[segi] == int(sl.load2X) {
		sl._expand()
	}
}

func (sl *SortedListPlus[T]) eraseBlock(bi int) {
	sl.blocks = Delete(sl.blocks, bi, bi+1)
	sl.blkMax = Delete(sl.blkMax, bi, bi+1)
	sl.bitCnt = Delete(sl.bitCnt, bi, bi+1)

	segi := (bi - 1) >> sl.lg2
	sl.segSize[segi]--
	if sl.segSize[segi] == 0 {
		if segi == len(sl.segSize)-1 {
			sl.segMax = sl.segMax[:segi]
			sl.segSize = sl.segSize[:segi]
			return
		}
		sl._expand()
	} else {
		bn := (segi + 1) << sl.lg2
		if bn > len(sl.bitCnt)-1 {
			bn = len(sl.bitCnt) - 1
		}
		sl._rangeBitModify(segi<<sl.lg2|1, bn)

		bj := (segi << sl.lg2) | sl.segSize[segi]
		if bi == bj {
			sl.segMax[segi] = sl.blkMax[bj]
		}
	}
}

func (sl *SortedListPlus[T]) _rangeBitModify(b1, b2 int) {
	if b1 > b2 || b2 >= len(sl.bitCnt) {
		return
	}
	for i := b1; i <= b2; i++ {
		sl.bitCnt[i] = 0
		if len(sl.blocks[i]) > 0 {
			sl.bitCnt[i] = len(sl.blocks[i])
			sl.blkMax[i] = sl.blocks[i][len(sl.blocks[i])-1]
		}
	}
	for i := b1; i <= b2; i++ {
		j := i + (i & -i)
		if j <= b2 {
			sl.bitCnt[j] += sl.bitCnt[i]
		}
	}
}

func (sl *SortedListPlus[T]) _expand() {
	oldBlocks := sl.blocks
	c := 0
	for i := 1; i < len(oldBlocks); i++ {
		if len(oldBlocks[i]) > 0 {
			c++
		}
	}
	if c == 0 {
		sl.Clear()
		return
	}
	segn := (c + int(sl.load2) - 1) / int(sl.load2)

	ec := sl.elemCnt
	sl.Clear()
	sl.elemCnt = ec

	sl.segMax = make([]T, segn)
	sl.segSize = make([]int, segn)

	for i, j := 1, 0; i < len(oldBlocks); i++ {
		if len(oldBlocks[i]) > 0 {
			segi := j >> sl.lg2
			sl.segSize[segi]++
			blk := oldBlocks[i]
			sl.blocks = append(sl.blocks, blk)
			sl.bitCnt = append(sl.bitCnt, len(blk))
			sl.blkMax = append(sl.blkMax, blk[len(blk)-1])
			if blk[len(blk)-1] > sl.segMax[segi] {
				sl.segMax[segi] = blk[len(blk)-1]
			}
			j++
			if j&int(sl.load2-1) == 0 && j < c {
				for k := 0; k < int(sl.load2); k++ {
					sl.blocks = append(sl.blocks, nil)
					sl.bitCnt = append(sl.bitCnt, 0)
					sl.blkMax = append(sl.blkMax, sl.minVal)
				}
			}
		}
	}

	for i := 1; i < len(sl.bitCnt); i++ {
		j := i + (i & -i)
		if j < len(sl.bitCnt) {
			sl.bitCnt[j] += sl.bitCnt[i]
		}
	}
}

func chmax[T Ordered](a *T, b T) {
	if b > *a {
		*a = b
	}
}

func log2(x uint) uint {
	var r uint
	for (1 << r) <= x {
		r++
	}
	if r > 0 {
		r--
	}
	return r
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
// is defined by cmp. cmp should return 0 if the slice element matches
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
