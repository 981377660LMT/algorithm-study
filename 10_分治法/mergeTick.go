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
	mt := NewMergeTrick(bigArr)
	time1 := time.Now()
	for i := 0; i < 1e5; i++ {
		mt.Add(0, 500, 1)
		mt.GetSortedNums()
	}
	fmt.Println(time.Since(time1))

	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	mt = NewMergeTrick(arr)
	mt.Add(0, 5, 1)
	fmt.Println(mt.GetSortedNums())

}

type MergeTrick struct {
	_nums        []int
	_originIndex []int
	_sortedNums  []int
	_dirty       bool
}

func NewMergeTrick(nums []int) *MergeTrick {
	nums = append(nums[:0:0], nums...)
	order := make([]int, len(nums))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		return nums[order[i]] < nums[order[j]]
	})
	originIndex := make([]int, len(nums))
	for i := range originIndex {
		originIndex[order[i]] = i
	}
	sortedNums := append(nums[:0:0], nums...)
	sort.Ints(sortedNums)
	return &MergeTrick{_nums: nums, _originIndex: originIndex, _sortedNums: sortedNums}
}

func (mt *MergeTrick) Add(start, end, delta int) {
	mt._dirty = true
	n := len(mt._nums)
	modified := make([]int, end-start)
	unmodified := make([]int, n-(end-start))
	for i, ptr1, ptr2 := 0, 0, 0; i < n; i++ {
		index := mt._originIndex[i]
		if index >= start && index < end {
			modified[ptr1] = index
			mt._nums[index] += delta
			ptr1++
		} else {
			unmodified[ptr2] = index
			ptr2++
		}
	}

	// 归并
	i1, i2, k := 0, 0, 0
	for i1 < len(modified) && i2 < len(unmodified) {
		if mt._nums[modified[i1]] < mt._nums[unmodified[i2]] {
			mt._originIndex[k] = modified[i1]
			i1++
		} else {
			mt._originIndex[k] = unmodified[i2]
			i2++
		}
		k++
	}

	for i1 < len(modified) {
		mt._originIndex[k] = modified[i1]
		i1++
		k++
	}

	for i2 < len(unmodified) {
		mt._originIndex[k] = unmodified[i2]
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
		res[i] = mt._nums[mt._originIndex[i]]
	}
	mt._sortedNums = res
	return res
}
