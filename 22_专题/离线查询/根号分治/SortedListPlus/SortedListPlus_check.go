package main

import (
	"bufio"
	"math/rand"
	"time"

	"fmt"
	"os"
	"sort"
	"unsafe"
)

func main() {
	// for i := 0; i < 100; i++ {
	// 	judge()
	// }
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
	sort.Ints(nums)
	sl := NewNaiseSlFrom(nums)
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

func cmpOrder[T Ordered](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return +1
	default:
		return 0
	}
}

type SortedListPlus[T Ordered] struct {
	load1  int
	load2  int
	load1X int
	load2X int
	lg2    int
	minVal T

	blocks [][]T
	bitCnt []int
	blkMax []T

	segMax  []T
	segSize []int

	elemCnt int
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

	blockCount := (len(sorted) + sl.load1 - 1) / sl.load1
	sl.blocks = make([][]T, 0, blockCount)
	start := 0
	for start < len(sorted) {
		end := start + sl.load1
		if end > len(sorted) {
			end = len(sorted)
		}
		sl.blocks = append(sl.blocks, Clone(sorted[start:end]))
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

func (sl *SortedListPlus[T]) Insert(x T) {
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

	if pos < len(sl.blocks[bi]) && sl.blocks[bi][pos] == x {
		return
	}
	sl.elemCnt++

	sl.fenwickUpdate(bi, 1)
	sl.blocks[bi] = Insert(sl.blocks[bi], pos, x)

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

	sl.blocks[bi] = Delete(sl.blocks[bi], pos, pos+1)
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

func (sl *SortedListPlus[T]) At(k int) T {
	if k < 0 || k >= sl.elemCnt {
		panic("At: out of range")
	}
	blkIndex := sl.findFenwick(k)
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

func (sl *SortedListPlus[T]) Floor(x T) (res T, ok bool) {
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

func (sl *SortedListPlus[T]) lowerBound(x T) (bi, pos int) {
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

	pos, _ = BinarySearchFunc(sl.blocks[bi], x, cmpOrder[T])
	return
}

func (sl *SortedListPlus[T]) upperBound(x T) (bi, pos int) {
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

	pos, _ = BinarySearchFunc(sl.blocks[bi], x, func(a, b T) int {
		if a <= b {
			return -1
		}
		return +1
	})
	return
}

func (sl *SortedListPlus[T]) prev(bi, pos int) (int, int) {
	if pos > 0 {
		return bi, pos - 1
	}
	if bi == 0 {
		return 0, 0
	}
	return bi - 1, len(sl.blocks[bi-1]) - 1
}

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

func (sl *SortedListPlus[T]) eraseBlock(bi int) {
	sl.blocks = Delete(sl.blocks, bi, bi+1)
	sl.blkMax = Delete(sl.blkMax, bi, bi+1)
	sl.bitCnt = Delete(sl.bitCnt, bi, bi+1)

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

func (sl *SortedListPlus[T]) fenwickUpdate(idx, delta int) {
	for idx < len(sl.bitCnt) {
		sl.bitCnt[idx] += delta
		idx |= (idx + 1)
	}
}

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

func Clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
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

// judge: 在此函数中进行随机测试，对拍 SortedListPlus vs naiveSL
func judge() {
	const N = 50000           // 操作次数
	const rangeVal = int(1e9) // 值域
	rand.Seed(time.Now().UnixNano())

	// 先产生 N 个随机数，排好序，放入 SortedListPlus & naiveSL
	nums := make([]int, N)
	for i := 0; i < N; i++ {
		nums[i] = rand.Intn(rangeVal)
	}
	sort.Ints(nums)

	slp := NewSortedListPlusFrom(nums)
	naive := NewNaiseSlFrom(nums) // 朴素做法

	for i := 1; i <= N; i++ {
		op := rand.Intn(9) // 0..7
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
			rslp := slp.BisectLeft(x)
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
			rslp, ok1 := slp.Floor(x)
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
			rslp, ok1 := slp.Ceiling(x)
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
			rslp := slp.BisectRight(x)
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

// -------------------- 朴素对拍结构 --------------------
type naiveSL[T Ordered] struct {
	data []T
}

func NewNaiseSlFrom[T Ordered](sorted []T) *naiveSL[T] {
	sorted = Clone(sorted)
	return &naiveSL[T]{data: sorted}
}

// 若已存在则不再插入；否则插入保持有序
func (n *naiveSL[T]) Insert(x T) {
	pos, ok := BinarySearchFunc(n.data, x, cmpOrder[T])
	if ok {
		return
	}
	n.data = Insert(n.data, pos, x)
}

func (n *naiveSL[T]) Erase(x T) {
	pos, _ := BinarySearchFunc(n.data, x, cmpOrder[T])
	if pos < len(n.data) && n.data[pos] == x {
		n.data = Delete(n.data, pos, pos+1)
	}
}

func (n *naiveSL[T]) Size() int { return len(n.data) }

func (n *naiveSL[T]) BisectLeft(x T) int {
	pos, _ := BinarySearchFunc(n.data, x, cmpOrder[T])
	return pos
}

func (n *naiveSL[T]) BisectRight(x T) int {
	pos, _ := BinarySearchFunc(n.data, x, func(a, b T) int {
		if a <= b {
			return -1
		}
		return +1
	})
	return pos
}

// Count(x) = upper_bound(x) - lower_bound(x)
func (n *naiveSL[T]) Count(x T) int {
	lpos, _ := BinarySearchFunc(n.data, x, cmpOrder[T])
	rpos, _ := BinarySearchFunc(n.data, x, func(a, b T) int {
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
	pos, _ := BinarySearchFunc(n.data, x, func(a, b T) int {
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
	pos, _ := BinarySearchFunc(n.data, x, cmpOrder[T])
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
