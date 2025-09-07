// 使用缓冲区进行批量/延迟更新
// 写操作频繁的场景下，能通过将多次小的、昂贵的切片插入/删除操作合并为一次大的、高效的 rebuild 操作，来显著提升性能
// 分块和延迟更新，这是一种在数据库和数据结构中常见的优化技术，用于将多次小的写操作（插入、删除）合并为一次大的写操作，以提高整体性能。

package main

import (
	"slices"
)

// 315. 计算右侧小于当前元素的个数
// 给你一个整数数组 nums ，按要求返回一个新数组 counts 。数组 counts 有该性质： counts[i] 的值是  nums[i] 右侧小于 nums[i] 的元素的数量。
func countSmaller(nums []int) []int {
	sl := NewSimpleSortedList(nil, func(a, b int) int { return a - b })
	res := make([]int, len(nums))
	for i := len(nums) - 1; i >= 0; i-- {
		num := nums[i]
		res[i] = sl.Len() - sl.CountGreaterThanOrEqual(num)
		sl.Add(num)
	}
	return res
}

const (
	dataSize   int = 2560
	addBufSize int = 2560
	delBufSize int = 2560
)

// 基于 BatchUpdate 实现的有序列表.
type SimpleSortedList[V any] struct {
	data   []V
	cmp    func(a, b V) int
	addBuf []V
	delBuf []V
}

func NewSimpleSortedList[V any](data []V, cmp func(a, b V) int) *SimpleSortedList[V] {
	sortedData := slices.Clone(data)
	slices.SortFunc(sortedData, cmp)
	return &SimpleSortedList[V]{
		data: sortedData,
		cmp:  cmp,
	}
}

func (s *SimpleSortedList[V]) Add(value V) {
	if s.addBuf == nil && len(s.data) >= dataSize {
		s.addBuf = make([]V, 0, addBufSize)
		s.delBuf = make([]V, 0, delBufSize)
	}

	if s.addBuf != nil { // buffered mode
		if len(s.addBuf) >= addBufSize {
			s.rebuild()
		}
		idx, _ := slices.BinarySearchFunc(s.addBuf, value, s.cmp)
		s.addBuf = slices.Insert(s.addBuf, idx, value)
	} else { // direct mode
		idx, _ := slices.BinarySearchFunc(s.data, value, s.cmp)
		s.data = slices.Insert(s.data, idx, value)
	}
}

// 需要保证 value 存在于列表中.
func (s *SimpleSortedList[V]) Remove(value V) {
	if s.addBuf != nil {
		if len(s.delBuf) >= delBufSize {
			s.rebuild()
		}
		idx, _ := slices.BinarySearchFunc(s.delBuf, value, s.cmp)
		s.delBuf = slices.Insert(s.delBuf, idx, value)
	} else {
		idx, _ := slices.BinarySearchFunc(s.data, value, s.cmp)
		s.data = slices.Delete(s.data, idx, idx+1)
	}
}

// CountGreaterThanOrEqual 计算列表中大于或等于 value 的元素数量.
func (s *SimpleSortedList[V]) CountGreaterThanOrEqual(value V) int {
	idxData, _ := slices.BinarySearchFunc(s.data, value, s.cmp)
	count := len(s.data) - idxData
	if s.addBuf != nil {
		idxAdd, _ := slices.BinarySearchFunc(s.addBuf, value, s.cmp)
		count += len(s.addBuf) - idxAdd
		idxDel, _ := slices.BinarySearchFunc(s.delBuf, value, s.cmp)
		count -= len(s.delBuf) - idxDel
	}
	return count
}

func (s *SimpleSortedList[V]) Len() int {
	if s.addBuf == nil {
		return len(s.data)
	}
	return len(s.data) + len(s.addBuf) - len(s.delBuf)
}

// rebuild 高效地合并主列表和缓冲区，生成新的主列表.
func (s *SimpleSortedList[V]) rebuild() {
	if s.addBuf == nil {
		return
	}

	s.data = slices.Insert(s.data, 0, s.addBuf...)
	n1, n2, n3 := len(s.data), len(s.addBuf), len(s.delBuf)
	ptrData, ptrAdd, ptrDel := n2, 0, 0
	k := 0
	for ptrData < n1 || ptrAdd < n2 {
		if ptrData == n1 || ptrAdd == n2 {
			if ptrAdd == n2 {
				s.data[k] = s.data[ptrData]
				ptrData++
			} else {
				s.data[k] = s.addBuf[ptrAdd]
				ptrAdd++
			}
			k++
		} else {
			if s.cmp(s.data[ptrData], s.addBuf[ptrAdd]) <= 0 {
				s.data[k] = s.data[ptrData]
				ptrData++
			} else {
				s.data[k] = s.addBuf[ptrAdd]
				ptrAdd++
			}
			k++
		}

		if ptrDel < n3 && s.cmp(s.data[k-1], s.delBuf[ptrDel]) == 0 {
			ptrDel++
			k--
		}
	}

	s.data = s.data[:k]
	s.addBuf = s.addBuf[:0]
	s.delBuf = s.delBuf[:0]
}
