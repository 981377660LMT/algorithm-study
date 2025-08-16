package main

import (
	"fmt"
	"sort"
)

func main() {
	testRangeFreq()
}

type RangeFreq[T comparable] struct {
	valueToIndexes map[T][]int
}

func NewRangeFreq[T comparable](nums []T) *RangeFreq[T] {
	rf := &RangeFreq[T]{
		valueToIndexes: make(map[T][]int),
	}
	for i, v := range nums {
		rf.valueToIndexes[v] = append(rf.valueToIndexes[v], i)
	}
	return rf
}

// Query 查询区间 [start, end) 中值为 value 的元素个数.
func (rf *RangeFreq[T]) Query(start, end int, value T) int {
	if start >= end {
		return 0
	}
	pos, exists := rf.valueToIndexes[value]
	if !exists {
		return 0
	}
	return sort.SearchInts(pos, end) - sort.SearchInts(pos, start)
}

// FindFirst 查找区间 [start, end) 中值为 value 的第一个位置.
// 返回第一个位置的索引，如果不存在返回 -1.
func (rf *RangeFreq[T]) FindFirst(start, end int, value T) int {
	if start >= end {
		return -1
	}
	pos, exists := rf.valueToIndexes[value]
	if !exists {
		return -1
	}
	idx := sort.SearchInts(pos, start)
	if idx < len(pos) && pos[idx] < end {
		return pos[idx]
	}
	return -1
}

// FindLast 查找区间 [start, end) 中值为 value 的最后一个位置.
// 返回最后一个位置的索引，如果不存在返回 -1.
func (rf *RangeFreq[T]) FindLast(start, end int, value T) int {
	if start >= end {
		return -1
	}
	pos, exists := rf.valueToIndexes[value]
	if !exists {
		return -1
	}
	idx := sort.SearchInts(pos, end)
	if idx > 0 && pos[idx-1] >= start {
		return pos[idx-1]
	}
	return -1
}

func testRangeFreq() {
	// 测试基本功能
	nums := []int{1, 2, 3, 2, 4, 2, 5}
	rf := NewRangeFreq(nums)

	// 测试Query方法
	assert(rf.Query(0, 7, 2) == 3, "元素2在整个数组中出现3次")
	assert(rf.Query(1, 5, 2) == 2, "元素2在区间[1,5)中出现2次")
	assert(rf.Query(0, 3, 2) == 1, "元素2在区间[0,3)中出现1次")
	assert(rf.Query(0, 7, 6) == 0, "元素6不存在")

	// 测试FindFirst方法
	assert(rf.FindFirst(0, 7, 2) == 1, "元素2的第一个位置是1")
	assert(rf.FindFirst(2, 7, 2) == 3, "在区间[2,7)中元素2的第一个位置是3")
	assert(rf.FindFirst(4, 7, 2) == 5, "在区间[4,7)中元素2的第一个位置是5")
	assert(rf.FindFirst(0, 7, 6) == -1, "元素6不存在")
	assert(rf.FindFirst(6, 7, 2) == -1, "区间[6,7)中没有元素2")

	// 测试FindLast方法
	assert(rf.FindLast(0, 7, 2) == 5, "元素2的最后一个位置是5")
	assert(rf.FindLast(0, 4, 2) == 3, "在区间[0,4)中元素2的最后一个位置是3")
	assert(rf.FindLast(0, 2, 2) == 1, "在区间[0,2)中元素2的最后一个位置是1")
	assert(rf.FindLast(0, 7, 6) == -1, "元素6不存在")
	assert(rf.FindLast(0, 1, 2) == -1, "区间[0,1)中没有元素2")

	// 测试边界情况
	assert(rf.Query(3, 3, 2) == 0, "空区间")
	assert(rf.FindFirst(3, 3, 2) == -1, "空区间")
	assert(rf.FindLast(3, 3, 2) == -1, "空区间")

	// 测试字符串类型
	strNums := []string{"a", "b", "c", "b", "d", "b"}
	strRf := NewRangeFreq(strNums)
	assert(strRf.Query(0, 6, "b") == 3, "字符串测试")
	assert(strRf.FindFirst(0, 6, "b") == 1, "字符串测试")
	assert(strRf.FindLast(0, 6, "b") == 5, "字符串测试")

	fmt.Println("所有测试通过!")
}

func assert(condition bool, message string) {
	if !condition {
		panic(fmt.Sprintf("断言失败: %s", message))
	}
}
