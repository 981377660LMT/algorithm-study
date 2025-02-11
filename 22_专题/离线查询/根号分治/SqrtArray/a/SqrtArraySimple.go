package main

import (
	"fmt"
	"math"
	"unsafe"
)

func main() {

}

type SqrtArray[T any] struct {
	n      int
	blocks [][]T
}

// NewSqrtArray 创建一个新的SqrtArray
func NewSqrtArray[T any](n int, f func(int) T, blockSize int) *SqrtArray[T] {
	if n <= 0 {
		return &SqrtArray[T]{}
	}

	blockCount := int(math.Sqrt(float64(n)))
	if blockSize < 1 {
		blockSize = (n + blockCount - 1) / blockCount
	}

	blocks := make([][]T, blockCount)
	for i := 0; i < blockCount; i++ {
		start, end := i*blockSize, min((i+1)*blockSize, n)
		blocks[i] = make([]T, 0, end-start)
		for j := start; j < end; j++ {
			blocks[i] = append(blocks[i], f(j))
		}
	}

	return &SqrtArray[T]{
		n:      n,
		blocks: blocks,
	}
}

func (s *SqrtArray[T]) Set(i int, v T) {
	if i < 0 || i >= s.n {
		panic("index out of range")
	}
	bid, pos := s.findKth(i)
	s.blocks[bid][pos] = v
}

func (s *SqrtArray[T]) Get(i int) T {
	if i < 0 {
		i += s.n
	}
	if i < 0 || i >= s.n {
		panic("index out of range")
	}
	bid, pos := s.findKth(i)
	return s.blocks[bid][pos]
}

// At 获取指定索引的值，支持负数索引
func (s *SqrtArray[T]) At(i int) (T, bool) {
	if i < 0 {
		i += s.n
	}
	if i < 0 || i >= s.n {
		var zero T
		return zero, false
	}
	return s.Get(i), true
}

// Push 在数组末尾添加元素
func (s *SqrtArray[T]) Push(v T) {
	s.Insert(s.n, v)
}

// Pop 删除并返回指定索引的元素
func (s *SqrtArray[T]) Pop(i int) (T, bool) {
	if i < 0 {
		i += s.n
	}
	if i < 0 || i >= s.n {
		var zero T
		return zero, false
	}

	var bi int
	var res T
	if i == s.n-1 {
		bi = len(s.blocks) - 1
		res = s.blocks[bi][len(s.blocks[bi])-1]
		s.blocks[bi] = s.blocks[bi][:len(s.blocks[bi])-1]
	} else {
		for ; i >= len(s.blocks[bi]); i -= len(s.blocks[bi]) {
			bi++
		}
		res = s.blocks[bi][i]
		s.blocks[bi] = append(s.blocks[bi][:i], s.blocks[bi][i+1:]...)
	}
	s.n--
	if len(s.blocks[bi]) == 0 {
		s.blocks = append(s.blocks[:bi], s.blocks[bi+1:]...)
	}
	return res, true
}

// Shift 删除并返回数组第一个元素
func (s *SqrtArray[T]) Shift() (T, bool) {
	return s.Pop(0)
}

// Unshift 在数组开头添加元素
func (s *SqrtArray[T]) Unshift(v T) {
	s.Insert(0, v)
}

// Erase 删除区间 [start, end) 内的元素
func (s *SqrtArray[T]) Erase(start, end int) {
	if start < 0 {
		start = 0
	}
	if end > s.n {
		end = s.n
	}
	if start >= end {
		return
	}

	bid, startPos := s.findKth(start)
	deleteCount := end - start
	for bid < len(s.blocks) && deleteCount > 0 {
		block := s.blocks[bid]
		endPos := min(len(block), startPos+deleteCount)
		curDeleteCount := endPos - startPos
		if curDeleteCount == len(block) {
			s.blocks = append(s.blocks[:bid], s.blocks[bid+1:]...)
			bid--
		} else {
			s.blocks[bid] = append(block[:startPos], block[endPos:]...)
		}
		deleteCount -= curDeleteCount
		s.n -= curDeleteCount
		startPos = 0
		bid++
	}
}

// Insert 在指定位置插入元素
func (s *SqrtArray[T]) Insert(i int, v T) {
	if s.n == 0 {
		s.blocks = append(s.blocks, []T{v})
		s.n++
		return
	}

	bi := 0
	if i >= s.n {
		bi = len(s.blocks) - 1
		s.blocks[bi] = append(s.blocks[bi], v)
	} else {
		for ; bi < len(s.blocks) && i >= len(s.blocks[bi]); i -= len(s.blocks[bi]) {
			bi++
		}
		s.blocks[bi] = append(s.blocks[bi][:i], append([]T{v}, s.blocks[bi][i:]...)...)
	}
	s.n++

	sqrtn2 := int(math.Sqrt(float64(s.n))) * 3
	// 定期重构
	if len(s.blocks[bi]) > 2*sqrtn2 {
		y := s.blocks[bi][sqrtn2:]
		s.blocks[bi] = s.blocks[bi][:sqrtn2]
		s.blocks = append(s.blocks[:bi+1], append([][]T{y}, s.blocks[bi+1:]...)...)
	}
}

// Enumerate 遍历区间 [start, end) 内的元素，并选择是否在遍历后删除
func (s *SqrtArray[T]) Enumerate(start, end int, f func(T), erase bool) {
	bid, startPos := s.findKth(start)
	count := end - start

	for bid < len(s.blocks) && count > 0 {
		block := s.blocks[bid]
		endPos := min(len(block), startPos+count)
		for j := startPos; j < endPos; j++ {
			f(block[j])
		}

		curDeleteCount := endPos - startPos
		if erase {
			if curDeleteCount == len(block) {
				s.blocks = append(s.blocks[:bid], s.blocks[bid+1:]...)
				bid--
			} else {
				s.blocks[bid] = append(block[:startPos], block[endPos:]...)
			}
			s.n -= curDeleteCount
		}

		count -= curDeleteCount
		startPos = 0
		bid++
	}
}

// Slice 返回指定区间的切片
func (s *SqrtArray[T]) Slice(start, end int) []T {
	if start < 0 {
		start = 0
	}
	if end > s.n {
		end = s.n
	}
	if start >= end {
		return []T{}
	}
	count := end - start
	res := make([]T, count)
	bid, startPos := s.findKth(start)
	ptr := 0
	for bid < len(s.blocks) && count > 0 {
		block := s.blocks[bid]
		endPos := min(len(block), startPos+count)
		curCount := endPos - startPos
		copy(res[ptr:], block[startPos:endPos])
		ptr += curCount
		count -= curCount
		startPos = 0
		bid++
	}
	return res
}

// Fill 用指定值填充整个数组
func (s *SqrtArray[T]) Fill(v T) {
	for i := range s.blocks {
		for j := range s.blocks[i] {
			s.blocks[i][j] = v
		}
	}
}

// ISlice 返回一个迭代器，用于遍历区间 [start, end) 内的元素
func (s *SqrtArray[T]) ISlice(start, end int, reverse bool) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		if start < 0 {
			start = 0
		}
		if end > s.n {
			end = s.n
		}
		if start >= end {
			return
		}
		count := end - start

		if reverse {
			bid, endPos := s.findKth(end - 1)
			for bid >= 0 && count > 0 {
				block := s.blocks[bid]
				startPos := max(0, endPos-count)
				curCount := endPos - startPos
				for j := endPos - 1; j >= startPos; j-- {
					ch <- block[j]
				}
				count -= curCount
				bid--
				if bid >= 0 {
					endPos = len(s.blocks[bid])
				}
			}
		} else {
			bid, startPos := s.findKth(start)
			for bid < len(s.blocks) && count > 0 {
				block := s.blocks[bid]
				endPos := min(len(block), startPos+count)
				curCount := endPos - startPos
				for j := startPos; j < endPos; j++ {
					ch <- block[j]
				}
				count -= curCount
				startPos = 0
				bid++
			}
		}
	}()
	return ch
}

// ForEach 遍历整个数组
func (s *SqrtArray[T]) ForEach(callback func(T, int)) {
	ptr := 0
	for bi := 0; bi < len(s.blocks); bi++ {
		for j := 0; j < len(s.blocks[bi]); j++ {
			callback(s.blocks[bi][j], ptr)
			ptr++
		}
	}
}

// Entries 返回一个迭代器，用于遍历数组的索引和值
func (s *SqrtArray[T]) Entries() <-chan [2]interface{} {
	ch := make(chan [2]interface{})
	go func() {
		defer close(ch)
		ptr := 0
		for i := 0; i < len(s.blocks); i++ {
			block := s.blocks[i]
			for j := 0; j < len(block); j++ {
				ch <- [2]interface{}{ptr, block[j]}
				ptr++
			}
		}
	}()
	return ch
}

// Iterator 返回一个迭代器，用于遍历数组的值
func (s *SqrtArray[T]) Iterator() <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for i := 0; i < len(s.blocks); i++ {
			block := s.blocks[i]
			for j := 0; j < len(block); j++ {
				ch <- block[j]
			}
		}
	}()
	return ch
}

func (s *SqrtArray[T]) String() string {
	return fmt.Sprintf("SqrtArray{%v}", s.blocks)
}

func (s *SqrtArray[T]) findKth(index int) (int, int) {
	var bi, pos int
	if index < s.n>>1 {
		bi, pos = s.findFromStart(index)
	} else {
		bi, pos = s.findFromEnd(s.n - index - 1)
	}
	return bi, pos
}

func (s *SqrtArray[T]) findFromStart(step int) (bi, pos int) {
	for j := 0; j < len(s.blocks); j++ {
		if step < len(s.blocks[j]) {
			return j, step
		}
		step -= len(s.blocks[j])
	}
	panic("index out of range")
}

func (s *SqrtArray[T]) findFromEnd(step int) (bi, pos int) {
	for j := len(s.blocks) - 1; j >= 0; j-- {
		if step < len(s.blocks[j]) {
			return j, len(s.blocks[j]) - 1 - step
		}
		step -= len(s.blocks[j])
	}
	panic("index out of range")
}

func (s *SqrtArray[T]) Len() int {
	return s.n
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

// Replace replaces the elements s[i:j] by the given v, and returns the modified slice.
// !Like JavaScirpt's Array.prototype.splice.
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

func startIdx[E any](haystack, needle []E) int {
	p := &needle[0]
	for i := range haystack {
		if p == &haystack[i] {
			return i
		}
	}
	panic("needle not found")
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
