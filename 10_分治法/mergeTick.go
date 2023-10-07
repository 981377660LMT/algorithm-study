package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {

	bigArr := make([]int, 1e3)
	for i := 0; i < 1e3; i++ {
		bigArr[i] = rand.Intn(1e3)
	}
	mt := NewMergeTrick(bigArr, true)
	time1 := time.Now()
	for i := 0; i < 1e5; i++ {
		mt.Add(0, 500, 1)
		mt.GetSortedNums()
	}
	fmt.Println(time.Since(time1))

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	mt = NewMergeTrick(arr, true)
	mt.Add(0, 5, 1)
	fmt.Println(mt.GetSortedNums())

}

type SortedItem = struct{ value, index int }
type MergeTrick struct {
	_nums        []int
	_sortedItems []*SortedItem
	_sortedNums  []int
	_dirty       bool
}

// O(n)区间加, O(n)整体排序.
//
//	shouldCopy: 是否复制nums.
func NewMergeTrick(nums []int, shouldCopy bool) *MergeTrick {
	if shouldCopy {
		nums = append(nums[:0:0], nums...)
	}
	sortedItems := make([]*SortedItem, len(nums))
	for i := range nums {
		sortedItems[i] = &SortedItem{value: nums[i], index: i}
	}
	sort.Slice(sortedItems, func(i, j int) bool { return sortedItems[i].value < sortedItems[j].value })
	sortedNums := make([]int, len(nums))
	for i := range sortedItems {
		sortedNums[i] = sortedItems[i].value
	}
	return &MergeTrick{
		_nums:        nums,
		_sortedItems: sortedItems,
		_sortedNums:  sortedNums,
	}
}

func (mt *MergeTrick) Add(start, end, delta int) {
	mt._dirty = true
	n := len(mt._nums)
	modified := make([]*SortedItem, end-start)
	unmodified := make([]*SortedItem, n-(end-start))
	for i, ptr1, ptr2 := 0, 0, 0; i < n; i++ {
		item := mt._sortedItems[i]
		if index := item.index; index >= start && index < end {
			item.value += delta
			modified[ptr1] = item
			ptr1++
			mt._nums[index] += delta
		} else {
			unmodified[ptr2] = item
			ptr2++
		}
	}

	// 归并
	i1, i2, k := 0, 0, 0
	for i1 < len(modified) && i2 < len(unmodified) {
		if modified[i1].value < unmodified[i2].value {
			mt._sortedItems[k] = modified[i1]
			i1++
		} else {
			mt._sortedItems[k] = unmodified[i2]
			i2++
		}
		k++
	}

	for i1 < len(modified) {
		mt._sortedItems[k] = modified[i1]
		i1++
		k++
	}

	for i2 < len(unmodified) {
		mt._sortedItems[k] = unmodified[i2]
		i2++
		k++
	}
}

// 返回原始数组.
func (mt *MergeTrick) GetNums() []int {
	return mt._nums
}

// 返回排序后的数组.
func (mt *MergeTrick) GetSortedNums() []int {
	if !mt._dirty {
		return mt._sortedNums
	}
	mt._dirty = false
	res := make([]int, len(mt._nums))
	for i := range res {
		res[i] = mt._sortedItems[i].value
	}
	mt._sortedNums = res
	return res
}
