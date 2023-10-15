package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {
	ps := NewPointSetRangeFreq([]int{12, 33, 4, 56, 22, 2, 34, 33, 22, 12, 34, 56}, 3)

	countNavie := func(nums []int, target int) int {
		res := 0
		for _, v := range nums {
			if v == target {
				res++
			}
		}
		return res
	}

	countLowerNavie := func(nums []int, target int) int {
		res := 0
		for _, v := range nums {
			if v < target {
				res++
			}
		}
		return res
	}

	countFloorNavie := func(nums []int, target int) int {
		res := 0
		for _, v := range nums {
			if v <= target {
				res++
			}
		}
		return res
	}

	for i := 0; i < len(ps.nums); i++ {
		for j := i; j < len(ps.nums); j++ {
			for k := 0; k <= 1000; k++ {
				if ps.Count(i, j+1, k) != countNavie(ps.nums[i:j+1], k) {
					fmt.Println(i, j+1, k)
					fmt.Println(ps.Count(i, j+1, k), countNavie(ps.nums[i:j+1], k))
					panic("")
				}
				if ps.CountLower(i, j+1, k) != countLowerNavie(ps.nums[i:j+1], k) {
					fmt.Println(i, j+1, k)
					fmt.Println(ps.CountLower(i, j+1, k), countLowerNavie(ps.nums[i:j+1], k))
					panic("")
				}
				if ps.CountFloor(i, j+1, k) != countFloorNavie(ps.nums[i:j+1], k) {
					fmt.Println(i, j+1, k)
					fmt.Println(ps.CountFloor(i, j+1, k), countFloorNavie(ps.nums[i:j+1], k))
					panic("")
				}

				randPos := i + rand.Intn(j-i+1)
				randValue := rand.Intn(1000)
				ps.Set(randPos, randValue)
			}
		}
	}

	fmt.Println("OK")

}

// https://leetcode.cn/problems/range-frequency-queries/description/
type RangeFreqQuery struct {
	ps *PointSetRangeFreq
}

func Constructor(arr []int) RangeFreqQuery {
	return RangeFreqQuery{}
}

func (this *RangeFreqQuery) Query(left int, right int, value int) int {
	return this.ps.Count(left, right+1, value)
}

// 单点修改，区间频率查询.
//
//	单次修改复杂度 O(sqrt(n)), 单次查询复杂度 O(sqrt(n) * log(sqrt(n))).
type PointSetRangeFreq struct {
	nums        []int
	belong      []int
	blockStart  []int
	blockEnd    []int
	blockCount  int
	blockSorted [][]int
}

// ps := NewPointSetRangeFreq(arr, 2*int(math.Sqrt(float64(len(arr)))+1))
func NewPointSetRangeFreq(nums []int, blockSize int) *PointSetRangeFreq {
	nums = append(nums[:0:0], nums...)
	res := &PointSetRangeFreq{}
	block := UseBlock(len(nums), blockSize)
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount

	blockSorted := make([][]int, blockCount)
	for bid := 0; bid < blockCount; bid++ {
		curSorted := make([]int, blockEnd[bid]-blockStart[bid])
		copy(curSorted, nums[blockStart[bid]:blockEnd[bid]])
		sort.Ints(curSorted)
		blockSorted[bid] = curSorted
	}

	res.nums = nums
	res.belong = belong
	res.blockStart = blockStart
	res.blockEnd = blockEnd
	res.blockCount = blockCount
	res.blockSorted = blockSorted
	return res
}

func (ps *PointSetRangeFreq) Set(pos, newValue int) {
	if ps.nums[pos] == newValue {
		return
	}
	pre := ps.nums[pos]
	ps.nums[pos] = newValue

	bid := ps.belong[pos]
	removeIndex := ps._bisectRight(ps.blockSorted[bid], pre, 0, len(ps.blockSorted[bid])-1) - 1
	ps.blockSorted[bid] = append(ps.blockSorted[bid][:removeIndex], ps.blockSorted[bid][removeIndex+1:]...)
	insertIndex := ps._bisectRight(ps.blockSorted[bid], newValue, 0, len(ps.blockSorted[bid])-1)
	ps.blockSorted[bid] = append(ps.blockSorted[bid], 0)
	copy(ps.blockSorted[bid][insertIndex+1:], ps.blockSorted[bid][insertIndex:])
	ps.blockSorted[bid][insertIndex] = newValue
}

// 统计 [start, end) 中等于 target 的元素个数.
func (ps *PointSetRangeFreq) Count(start, end int, target int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ps.nums) {
		end = len(ps.nums)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ps.belong[start], ps.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if ps.nums[i] == target {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ps.blockEnd[bid1]; i++ {
		if ps.nums[i] == target {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res += ps._count(ps.blockSorted[bid], target, 0, len(ps.blockSorted[bid])-1)
	}
	for i := ps.blockStart[bid2]; i < end; i++ {
		if ps.nums[i] == target {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中严格小于 target 的元素个数.
func (ps *PointSetRangeFreq) CountLower(start, end int, target int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ps.nums) {
		end = len(ps.nums)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ps.belong[start], ps.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if ps.nums[i] < target {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ps.blockEnd[bid1]; i++ {
		if ps.nums[i] < target {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res += ps._bisectLeft(ps.blockSorted[bid], target, 0, len(ps.blockSorted[bid])-1)
	}
	for i := ps.blockStart[bid2]; i < end; i++ {
		if ps.nums[i] < target {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中小于等于 target 的元素个数.
func (ps *PointSetRangeFreq) CountFloor(start, end int, target int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ps.nums) {
		end = len(ps.nums)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ps.belong[start], ps.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if ps.nums[i] <= target {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ps.blockEnd[bid1]; i++ {
		if ps.nums[i] <= target {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res += ps._bisectRight(ps.blockSorted[bid], target, 0, len(ps.blockSorted[bid])-1)
	}
	for i := ps.blockStart[bid2]; i < end; i++ {
		if ps.nums[i] <= target {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中大于等于 target 的元素个数.
func (ps *PointSetRangeFreq) CountCeiling(start, end int, target int) int {
	return (end - start) - ps.CountLower(start, end, target)
}

// 统计 [start, end) 中严格大于 target 的元素个数.
func (ps *PointSetRangeFreq) CountHigher(start, end int, target int) int {
	return (end - start) - ps.CountFloor(start, end, target)
}

func (ps *PointSetRangeFreq) _bisectLeft(nums []int, target int, left, right int) int {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] >= target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return left
}

func (ps *PointSetRangeFreq) _bisectRight(nums []int, target int, left, right int) int {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return left
}

func (ps *PointSetRangeFreq) _count(nums []int, target int, left, right int) int {
	return ps._bisectRight(nums, target, left, right) - ps._bisectLeft(nums, target, left, right)
}

// blockSize := int(math.Sqrt(float64(len(nums))) + 1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}
