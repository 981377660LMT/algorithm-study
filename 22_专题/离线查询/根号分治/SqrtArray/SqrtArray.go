// 动态分块数组/SqrtArray
//
// api:
//  1.Insert(index int32, v V)
//  2.Pop(index int32) V
//  3.Set(index int32, v V)
//  4.Get(index int32) V
//  5.Clear()
//  6.Len() int32
//  7.GetAll() []V
//  8.ForEach(f func(i int32, v V) bool)
//
// TODO: 采用两层分块结构优化、删除后合并相邻稀疏的块

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

func main() {
	// abc392_f()
	// demo()
	// test()
	testTime()
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

	arr := NewSqrtArray(0, func(i int32) int { return 0 }, -1)

	for i := int32(0); i < n; i++ {
		pos := int32(NextInt())
		pos--

		arr.Insert(pos, int(i+1))
	}

	arr.ForEach(func(_ int32, v int) bool {
		fmt.Fprint(out, v, " ")
		return false
	})
}

func demo() {
	bv := NewSqrtArray(10, func(i int32) int { return int(i) }, -1)
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

// 使用分块+树状数组维护的动态数组.
type SqrtArray struct {
	n                 int32
	blockSize         int32
	threshold         int32
	blocks            [][]E
	shouldRebuildTree bool
	tree              []int32 // 每个块块长的前缀和
}

func NewSqrtArray(n int32, f func(i int32) E, blockSize int32) *SqrtArray {
	if blockSize < 1 {
		blockSize = 1 << 7
	}

	res := &SqrtArray{n: n, blockSize: blockSize, threshold: blockSize << 1, shouldRebuildTree: true}
	blockCount := (n + blockSize - 1) / blockSize
	blocks := make([][]E, blockCount)
	for bid := int32(0); bid < blockCount; bid++ {
		start, end := bid*blockSize, (bid+1)*blockSize
		if end > n {
			end = n
		}
		bucket := make([]E, end-start)
		for i := start; i < end; i++ {
			bucket[i-start] = f(i)
		}
		blocks[bid] = bucket
	}
	res.blocks = blocks
	return res
}

func (sl *SqrtArray) Insert(index int32, value E) {
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []E{value})
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
	sl._updateTree(pos, 1)
	sl.blocks[pos] = Insert(sl.blocks[pos], int(startIndex), value)

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); n > sl.threshold {
		left, right := Clone(sl.blocks[pos][:sl.blockSize]), sl.blocks[pos][sl.blockSize:]
		sl.blocks = Replace(sl.blocks, int(pos), int(pos+1), left, right)
		sl.shouldRebuildTree = true
	}

	sl.n++
	return
}

func (sl *SqrtArray) Pop(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.n--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = Replace(sl.blocks[pos], int(startIndex), int(startIndex+1))

	if len(sl.blocks[pos]) == 0 {
		// !delete block
		sl.blocks = Replace(sl.blocks, int(pos), int(pos+1))
		sl.shouldRebuildTree = true
	}
	return value
}

func (sl *SqrtArray) Get(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SqrtArray) Set(index int32, value E) {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	sl.blocks[pos][startIndex] = value
}

func (sl *SqrtArray) Len() int32 {
	return sl.n
}

func (sl *SqrtArray) Clear() {
	sl.n = 0
	sl.shouldRebuildTree = false
	sl.blocks = sl.blocks[:0]
	sl.tree = sl.tree[:0]
}

func (sl *SqrtArray) GetAll() []E {
	res := make([]E, 0, sl.n)
	for _, block := range sl.blocks {
		res = append(res, block...)
	}
	return res
}

func (sl *SqrtArray) ForEach(f func(i int32, v E) (shouldBreak bool)) {
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

func (sl *SqrtArray) Erase(start, end int32) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SqrtArray) Enumerate(start, end int32, f func(value E), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl.n {
		end = sl.n
	}
	if start >= end {
		return
	}

	pos, startIndex := sl._findKth(start)
	count := end - start
	m := int32(len(sl.blocks))
	eraseStart := int32(-1)
	eraseCount := int32(0)
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		endIndex := min32(int32(len(block)), startIndex+count)
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == int32(len(block)) {
				// !delete block
				if eraseStart == -1 {
					eraseStart = pos
				}
				eraseCount++
			} else {
				// !delete [index, end)
				sl._updateTree(pos, -deleted)
				sl.blocks[pos] = Replace(sl.blocks[pos], int(startIndex), int(endIndex))
			}
			sl.n -= deleted
		}

		count -= deleted
		startIndex = 0
	}

	if erase && eraseStart != -1 {
		sl.blocks = Replace(sl.blocks, int(eraseStart), int(eraseStart+eraseCount))
		sl.shouldRebuildTree = true
	}
}

func (sl *SqrtArray) _rebuildTree() {
	nb := len(sl.blocks)
	if nb > len(sl.tree) {
		sl.tree = append(sl.tree, make([]int32, nb-len(sl.tree))...)
	} else if nb < len(sl.tree) {
		sl.tree = sl.tree[:nb]
	}
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

func (sl *SqrtArray) _updateTree(index int32, v int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	m := int32(len(tree))
	for i := index; i < m; i |= i + 1 {
		tree[i] += v
	}
}

func (sl *SqrtArray) _findKth(k int32) (pos, index int32) {
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

// Shallow clone.
func Clone[S ~[]E, E any](s S) S {
	return append(s[:0:0], s...)
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
		seg := NewSqrtArray(n, func(i int32) E { return E(nums[i]) }, -1)

		for j := 0; j < 1000; j++ {
			// Get
			index := rand.Int31n(seg.Len())
			if seg.Get(index) != E(nums[index]) {
				fmt.Println("Get Error")
				panic("Get Error")
			}

			// Set
			index = rand.Int31n(seg.Len())
			value := rand.Intn(100)
			nums[index] = value
			seg.Set(index, E(value))
			if seg.Get(index) != E(value) {
				fmt.Println("Set Error")
				panic("Set Error")
			}

			// Query
			start, end := rand.Int31n(seg.Len()), rand.Int31n(seg.Len())
			if start > end {
				start, end = end, start
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
			index = rand.Int31n(seg.Len())
			value = rand.Intn(100)
			nums = append(nums, 0)
			copy(nums[index+1:], nums[index:])
			nums[index] = value
			seg.Insert(index, E(value))

			// Pop
			index = rand.Int31n(seg.Len())
			value = E(nums[index])
			nums = append(nums[:index], nums[index+1:]...)
			if seg.Pop(index) != value {
				fmt.Println("Pop Error")
				panic("Pop Error")
			}

			// Erase
			start, end = rand.Int31n(seg.Len()), rand.Int31n(seg.Len())
			if start > end {
				start, end = end, start
			}
			nums = append(nums[:start], nums[end:]...)
			seg.Erase(start, end)

		}
	}
	fmt.Println("Pass")
}

func testTime() {
	n := int32(1e6)
	rands := make([]int, n)
	for i := int32(0); i < n; i++ {
		rands[i] = rand.Intn(1e9)
	}

	arr := NewSqrtArray(0, func(i int32) E { return E(i) }, 128)

	time1 := time.Now()
	for i := int32(0); i < n; i++ {
		arr.Insert(n-i, E(i))
		arr.Get(i)
		arr.Set(i, E(i))
		arr.Insert(0, E(i))
		arr.Pop(0)
	}

	for i := int32(0); i < n; i++ {
		arr.Get(i)
	}
	for i := int32(0); i < n; i++ {
		arr.Pop(0)
	}

	fmt.Println(time.Since(time1))
}
