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
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

func main() {
	abc392_f()
	// demo()
	// test()
	// testTime()
}

// https://atcoder.jp/contests/abc392/tasks/abc392_f
func abc392_f() {
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

	n := int32(NextInt())
	sqrt := int32(math.Sqrt(float64(n))) + 1
	arr := NewSqrtArrayMonoid(0, func(i int32) int32 { return 0 }, sqrt)

	for i := int32(0); i < n; i++ {
		pos := int32(NextInt())
		pos--
		arr.Insert(pos, int32(i+1))
	}

	arr.ForEach(func(_ int32, v int32) bool {
		fmt.Fprint(out, v, " ")
		return false
	})
}

func demo() {
	bv := NewSqrtArrayMonoid(10, func(i int32) int32 { return 0 }, -1)
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

type E = int32

func (*SqrtArrayMonoid) e() E        { return 0 }
func (*SqrtArrayMonoid) op(a, b E) E { return max32(a, b) }

// 使用分块+树状数组维护的动态数组.
type SqrtArrayMonoid struct {
	n                 int32
	blockSize         int32
	threshold         int32
	shouldRebuildTree bool
	blocks            [][]E
	blockSum          []E
	tree              []int32 // 每个块块长的前缀和
}

func NewSqrtArrayMonoid(n int32, f func(i int32) E, blockSize int32) *SqrtArrayMonoid {
	if blockSize < 1 {
		blockSize = int32(math.Sqrt(float64(n))) + 1
	}

	res := &SqrtArrayMonoid{n: n, blockSize: blockSize, threshold: blockSize << 1, shouldRebuildTree: true}
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

func (sl *SqrtArrayMonoid) Insert(index int32, value E) {
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
	sl.blocks[pos] = Insert(sl.blocks[pos], int(startIndex), value)

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); n > sl.threshold {
		left, right := Clone(sl.blocks[pos][:sl.blockSize]), sl.blocks[pos][sl.blockSize:]
		sl.blocks = Replace(sl.blocks, int(pos), int(pos+1), left, right)

		sl.blockSum = Replace(sl.blockSum, int(pos), int(pos+1), sl.e(), sl.e())
		sl._updateSum(pos)
		sl._updateSum(pos + 1)
		sl.shouldRebuildTree = true
	}

	sl.n++
	return
}

func (sl *SqrtArrayMonoid) Pop(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.n--
	sl._updateTree(pos, false)

	sl.blocks[pos] = Replace(sl.blocks[pos], int(startIndex), int(startIndex+1))
	if value != sl.e() {
		sl._updateSum(pos)
	}

	if len(sl.blocks[pos]) == 0 {
		// !delete block
		sl.blocks = Replace(sl.blocks, int(pos), int(pos+1))
		sl.blockSum = Replace(sl.blockSum, int(pos), int(pos+1))
		sl.shouldRebuildTree = true
	}
	return value
}

func (sl *SqrtArrayMonoid) Get(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SqrtArrayMonoid) Set(index int32, value E) {
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

func (sl *SqrtArrayMonoid) Sum(start, end int32) E {
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

func (sl *SqrtArrayMonoid) SumAll() E {
	res := sl.e()
	for _, v := range sl.blockSum {
		res = sl.op(res, v)
	}
	return res
}

func (sl *SqrtArrayMonoid) Len() int32 {
	return sl.n
}

func (sl *SqrtArrayMonoid) Clear() {
	sl.n = 0
	sl.shouldRebuildTree = true
	sl.blocks = sl.blocks[:0]
	sl.blockSum = sl.blockSum[:0]
	sl.tree = sl.tree[:0]
}

func (sl *SqrtArrayMonoid) GetAll() []E {
	res := make([]E, 0, sl.n)
	for _, block := range sl.blocks {
		res = append(res, block...)
	}
	return res
}

func (sl *SqrtArrayMonoid) ForEach(f func(i int32, v E) (shouldBreak bool)) {
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

func (sl *SqrtArrayMonoid) _rebuildTree() {
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

func (sl *SqrtArrayMonoid) _updateTree(index int32, addOne bool) {
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

func (sl *SqrtArrayMonoid) _findKth(k int32) (pos, index int32) {
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

func (sl *SqrtArrayMonoid) _updateSum(pos int32) {
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

func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	if i < 0 {
		i = 0
	}
	if j > len(s) {
		j = len(s)
	}
	if i == j {
		return Insert(s, i, v...)
	}
	if j == len(s) {
		return append(s[:i], v...)
	}
	tot := len(s[:i]) + len(v) + len(s[j:])
	if tot > cap(s) {
		s2 := append(s[:i], make(S, tot-i)...)
		copy(s2[i:], v)
		copy(s2[i+len(v):], s[j:])
		return s2
	}
	r := s[:tot]
	if i+len(v) <= j {
		copy(r[i:], v)
		copy(r[i+len(v):], s[j:])
		// clear(s[tot:])
		return r
	}
	if !overlaps(r[i+len(v):], v) {
		copy(r[i+len(v):], s[j:])
		copy(r[i:], v)
		return r
	}
	y := len(v) - (j - i)
	if !overlaps(r[i:j], v) {
		copy(r[i:j], v[y:])
		copy(r[len(s):], v[:y])
		rotateRight(r[i:], y)
		return r
	}
	if !overlaps(r[len(s):], v) {
		copy(r[len(s):], v[:y])
		copy(r[i:j], v[y:])
		rotateRight(r[i:], y)
		return r
	}
	k := startIdx(v, s[j:])
	copy(r[i:], v)
	copy(r[i+len(v):], r[i+k:])
	return r
}

func startIdx[E any](haystack, needle []E) int {
	p := &needle[0]
	for i := range haystack {
		if p == &haystack[i] {
			return i
		}
	}
	panic("needle not found")
}

func test() {
	for i := int32(0); i < 100; i++ {
		n := rand.Int31n(10000) + 1000
		nums := make([]int, n)
		for i := int32(0); i < n; i++ {
			nums[i] = rand.Intn(100)
		}
		seg := NewSqrtArrayMonoid(n, func(i int32) E { return E(nums[i]) }, -1)

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
			if seg.Sum(start, end) != sum_ {
				fmt.Println("Query Error")
				panic("Query Error")
			}

			// QueryAll
			sum_ = E(0)
			for _, v := range nums {
				sum_ = seg.op(sum_, E(v))
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
	seg := NewSqrtArrayMonoid(n, func(i int32) int32 { return E(nums[i]) }, -1)

	for i := int32(0); i < n; i++ {
		seg.Get(i)
		seg.Set(i, i)
		seg.Sum(i, n)
		seg.SumAll()
		seg.Insert(i, i)
		if i&1 == 0 {
			seg.Pop(i)
		}
		seg.SumAll()
	}
	fmt.Println("Time1", time.Since(time1)) // Time1 336.550792ms
}
