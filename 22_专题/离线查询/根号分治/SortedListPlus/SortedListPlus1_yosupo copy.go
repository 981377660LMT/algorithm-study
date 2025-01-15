package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"unsafe"
)

func main() {
	yosupo()
}

func yosupo() {
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

	sl := NewSortedListPlusFromSlice(func(a, b int) int { return a - b }, 200, 64, 0, nums)
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

type CmpFunc[T any] func(a, b T) int

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

type SortedListPlus[T Ordered] struct {
	cmp    CmpFunc[T]
	load1  int
	load2  int
	load1X int
	load2X int
	lg2    uint
	minVal T

	blocks     [][]T
	bitCnt     []int
	blkMax     []T
	segMax     []T
	segSize    []int
	elementCnt int
}

func NewSortedListPlus[T Ordered](cmp CmpFunc[T], load1, load2 int, minVal T) *SortedListPlus[T] {
	sl := &SortedListPlus[T]{
		cmp:    cmp,
		load1:  load1,
		load2:  load2,
		load1X: load1 * 2,
		load2X: load2 * 2,
		minVal: minVal,
	}
	sl.lg2 = uint(bits.Len(uint(sl.load2X)) - 1)
	sl.clear()
	return sl
}

func NewSortedListPlusFromSlice[T Ordered](cmp CmpFunc[T], load1, load2 int, minVal T, arr []T) *SortedListPlus[T] {
	arr = Clone(arr)
	// SortFunc(arr, cmp)
	less := func(i, j int) bool { return cmp(arr[i], arr[j]) < 0 }
	sort.Slice(arr, less)

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

	sl.blocks = make([][]T, 1, len(arr)/load1+2)

	for start := 0; start < len(arr); start += load1 {
		end := start + load1
		if end > len(arr) {
			end = len(arr)
		}
		chunk := Clone(arr[start:end])
		sl.blocks = append(sl.blocks, chunk)
	}

	sl._expand()
	return sl
}

func (sl *SortedListPlus[T]) clear() {
	sl.blocks = make([][]T, 1)
	sl.blocks[0] = nil
	sl.bitCnt = make([]int, 1)
	sl.blkMax = make([]T, 1)
	sl.segMax = sl.segMax[:0]
	sl.segSize = sl.segSize[:0]
	sl.elementCnt = 0
}

func (sl *SortedListPlus[T]) Size() int {
	return sl.elementCnt
}

func (sl *SortedListPlus[T]) chmax(a *T, b T) {
	if sl.cmp(b, *a) > 0 {
		*a = b
	}
}

func (sl *SortedListPlus[T]) lower_bound(x T) (int, int) {
	if len(sl.segSize) == 0 {
		return 0, 0
	}
	// 1. 在小段之间进行二分查找，找到第一个 segMax >= x 的小段
	l, r := -1, len(sl.segMax)-1
	var mid int
	for r-l > 1 {
		mid = (l + r) >> 1
		if sl.segMax[mid] >= x {
			r = mid
		} else {
			l = mid
		}
	}
	segi := r

	// 2. 在小段内的非空数据块之间进行二分查找，找到第一个 blkMax >= x 的数据块
	blockStart := (segi << sl.lg2) | 0 // segi * load2X
	blockEnd := (segi << sl.lg2) | sl.segSize[segi]
	l, r = blockStart-1, blockEnd
	for r-l > 1 {
		mid = (l + r) >> 1
		if sl.blkMax[mid] >= x {
			r = mid
		} else {
			l = mid
		}
	}
	bi := r

	idx := sl._binarySearchLowerBound(sl.blocks[bi], x)
	return bi, idx
}

func (sl *SortedListPlus[T]) upper_bound(x T) (int, int) {
	if len(sl.segSize) == 0 {
		return 0, 0
	}

	// 1. 在小段之间进行二分查找，找到第一个 segMax > x 的小段
	l, r := -1, len(sl.segMax)-1
	var mid int
	for r-l > 1 {
		mid = (l + r) >> 1
		if sl.segMax[mid] > x {
			r = mid
		} else {
			l = mid
		}
	}
	segi := r

	// 2. 在小段内的非空数据块之间进行二分查找，找到第一个 blkMax > x 的数据块
	blockStart := (segi << sl.lg2) | 0 // segi * load2X
	blockEnd := (segi << sl.lg2) | sl.segSize[segi]
	l, r = blockStart-1, blockEnd
	for r-l > 1 {
		mid = (l + r) >> 1
		if sl.blkMax[mid] > x {
			r = mid
		} else {
			l = mid
		}
	}
	bi := r

	idx := sl._binarySearchUpperBound(sl.blocks[bi], x)
	return bi, idx
}

func (sl *SortedListPlus[T]) Insert(x T) {
	if len(sl.segSize) == 0 {
		sl.elementCnt++
		sl.blocks = append(sl.blocks, []T{x})
		sl.bitCnt = append(sl.bitCnt, 1)
		sl.blkMax = append(sl.blkMax, x)
		sl.segSize = append(sl.segSize, 1)
		sl.segMax = append(sl.segMax, x)
		return
	}
	bi, idx := sl.lower_bound(x)
	if idx < len(sl.blocks[bi]) && sl.cmp(sl.blocks[bi][idx], x) == 0 {
		return
	}
	sl.elementCnt++

	for i := bi; i < len(sl.bitCnt); i += i & -i {
		sl.bitCnt[i]++
	}

	sl.blocks[bi] = Insert(sl.blocks[bi], idx, x)

	if len(sl.blocks[bi]) >= sl.load1X {
		sl._splitBlock(bi)
	} else {
		lastVal := sl.blocks[bi][len(sl.blocks[bi])-1]
		sl.blkMax[bi] = lastVal
	}
}

func (sl *SortedListPlus[T]) Erase(x T) {
	if len(sl.segSize) == 0 {
		return
	}
	bi, idx := sl.lower_bound(x)
	if bi >= len(sl.blocks) ||
		idx >= len(sl.blocks[bi]) ||
		sl.cmp(sl.blocks[bi][idx], x) != 0 {
		return
	}
	for i := bi; i < len(sl.bitCnt); i += i & -i {
		sl.bitCnt[i]--
	}
	sl.elementCnt--

	sl.blocks[bi] = Delete(sl.blocks[bi], idx, idx+1)
	if len(sl.blocks[bi]) == 0 {
		sl._shrinkBlock(bi)
	} else {
		lastVal := sl.blocks[bi][len(sl.blocks[bi])-1]
		sl.blkMax[bi] = lastVal
	}
}

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
	if bi >= len(sl.blocks) {
		panic("At: block index out of range")
	}
	if k >= len(sl.blocks[bi]) {
		panic("At: index within block out of range")
	}
	return sl.blocks[bi][k]
}

func (sl *SortedListPlus[T]) BisectLeft(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, idx := sl.lower_bound(x)
	rk := 0
	i := bi - 1
	for i > 0 {
		rk += sl.bitCnt[i]
		i -= i & -i
	}
	rk += idx
	return rk
}

func (sl *SortedListPlus[T]) BisectRight(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	bi, idx := sl.upper_bound(x)
	rk := 0
	i := bi - 1
	for i > 0 {
		rk += sl.bitCnt[i]
		i -= i & -i
	}
	rk += idx
	return rk
}

func (sl *SortedListPlus[T]) Count(x T) int {
	if len(sl.segSize) == 0 {
		return 0
	}
	return sl.BisectRight(x) - sl.BisectLeft(x)
}

func (sl *SortedListPlus[T]) Floor(x T) (T, bool) {
	bi, idx := sl.upper_bound(x)
	if !(bi > 1 || (bi == 1 && idx > 0)) {
		return sl.minVal, false
	}
	bi_prev, idx_prev := sl.prev(bi, idx)
	return sl.blocks[bi_prev][idx_prev], true
}

func (sl *SortedListPlus[T]) Ceiling(x T) (T, bool) {
	bi, idx := sl.lower_bound(x)
	if bi >= len(sl.blocks) || idx >= len(sl.blocks[bi]) {
		return sl.minVal, false
	}
	return sl.blocks[bi][idx], true
}

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

// TODO
func (sl *SortedListPlus[T]) _expand() {
	oldBlocks := sl.blocks
	c := 0
	for _, b := range oldBlocks {
		if len(b) > 0 {
			c++
		}
	}

	segn := (c + sl.load2 - 1) / sl.load2

	ec := sl.elementCnt
	sl.clear()
	sl.segSize = make([]int, segn)
	sl.blocks = sl.blocks[:1]
	sl.bitCnt = sl.bitCnt[:1]
	sl.blkMax = sl.blkMax[:1]
	sl.elementCnt = ec

	curSeg := 0
	countInSeg := 0

	for _, block := range oldBlocks {
		if len(block) == 0 {
			continue
		}
		sl.blocks = append(sl.blocks, block)
		sl.bitCnt = append(sl.bitCnt, len(block))
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

func (sl *SortedListPlus[T]) _rangeBitModify(b1, b2 int) {
	for i := b1; i <= b2; i++ {
		sl.bitCnt[i] = 0
	}

	// 2. 遍历 [b1, b2] 范围内的每个块
	for i := b1; i <= b2; i++ {
		if len(sl.blocks[i]) > 0 {
			sl.bitCnt[i] += len(sl.blocks[i])                // 更新 bitCnt[i] 为块的大小
			sl.blkMax[i] = sl.blocks[i][len(sl.blocks[i])-1] // 更新 blkMax[i] 为块的最后一个元素
		}
		// 计算 i + (i & -i)，这是最低位的 1 的位置
		next := i + (i & (-i))
		if next <= b2 && sl.bitCnt[i] > 0 {
			sl.bitCnt[next] += sl.bitCnt[i] // 更新 bitCnt[next] += bitCnt[i]
		}
	}

	for lowb := (-b2 & b2) / 2; lowb >= sl.load2X; lowb >>= 1 {
		sl.bitCnt[b2] += sl.bitCnt[b2-lowb]
	}
}

func (sl *SortedListPlus[T]) _splitBlock(bi int) {
	block := sl.blocks[bi]
	if len(block) < sl.load1X {
		return
	}
	newBlock := Clone(block[sl.load1:])
	sl.blocks[bi] = block[:sl.load1]

	pos := bi + 1
	sl.blocks = Insert(sl.blocks, pos, newBlock)
	sl.bitCnt = Insert(sl.bitCnt, pos, len(newBlock))
	sl.bitCnt[bi] = len(sl.blocks[bi])

	sl.blkMax = Insert(sl.blkMax, pos, newBlock[len(newBlock)-1])
	sl.blkMax[bi] = sl.blocks[bi][len(sl.blocks[bi])-1]

	segi := (bi - 1) >> (sl.lg2)
	if segi < len(sl.segSize) {
		sl.segSize[segi]++
		if sl.segSize[segi] == sl.load2X {
			sl._expand()
		} else {
			sl._rangeBitModify(segi<<sl.lg2|1, len(sl.blocks)-1)
		}
	}
}

func (sl *SortedListPlus[T]) _shrinkBlock(bi int) {
	if bi == len(sl.blocks)-1 {
		sl.blocks = sl.blocks[:bi]
		sl.blkMax = sl.blkMax[:bi]
		sl.bitCnt = sl.bitCnt[:bi]
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

	sl.blocks = Delete(sl.blocks, bi, bi+1)
	sl.blkMax = Delete(sl.blkMax, bi, bi+1)
	sl.bitCnt = Delete(sl.bitCnt, bi, bi+1)

	segi := (bi - 1) >> sl.lg2
	if segi < 0 {
		segi = 0
	}
	if segi >= len(sl.segSize) {
		sl._expand()
		return
	}
	sl.segSize[segi]--
	if sl.segSize[segi] == 0 {
		sl._expand()
	} else {
		sl._rangeBitModify(segi<<sl.lg2|1, len(sl.blocks)-1)
	}
}

func (sl *SortedListPlus[T]) _binarySearchLowerBound(arr []T, x T) int {
	idx, ok := BinarySearchFunc(arr, x, func(a, b T) int {
		return sl.cmp(a, b)
	})
	if ok {
		return idx
	}
	return idx
}

func (sl *SortedListPlus[T]) _binarySearchUpperBound(arr []T, x T) int {
	idx, ok := BinarySearchFunc(arr, x, func(a, b T) int {
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

func (sl *SortedListPlus[T]) String() string {
	return fmt.Sprintf(
		"SortedListPlus{elementCnt=%d, blocks=%v}",
		sl.elementCnt, sl.blocks,
	)
}

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

func BinarySearchFunc[S ~[]E, E, T any](x S, target T, cmp func(E, T) int) (int, bool) {
	n := len(x)
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		if cmp(x[h], target) < 0 {
			i = h + 1
		} else {
			j = h
		}
	}
	return i, i < n && cmp(x[i], target) == 0
}

// Shallow clone.
func Clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
}
