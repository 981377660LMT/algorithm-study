// https://www.luogu.com.cn/problem/P1975
// P1975 [国家集训队] 排队
// 区间动态逆序对
// n<=2e4 q<=2e3 nums[i]<=1e9

// !update/set = remove + add

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// 幼儿园阿姨每次会选出两个小朋友，交换他们的位置，请你帮忙计算出每次交换后，序列的逆序对数。
func InversePairs(nums []int, operations [][2]int) []int {
	res := make([]int, len(operations)+1)
	res[0] = countInvMergeSort(nums)
	ps := _NewPointSetRangeFreq(nums, 2*int(math.Sqrt(float64(len(nums)))+1))
	for i, operation := range operations {
		pos1, pos2 := operation[0], operation[1]
		pre1, pre2 := nums[pos1], nums[pos2]
		nums[pos1], nums[pos2] = nums[pos2], nums[pos1]
		if pre1 == pre2 {
			res[i+1] = res[i]
			continue
		}

		diff := 0
		if pos1 > pos2 {
			pos1, pos2 = pos2, pos1
			pre1, pre2 = pre2, pre1
		}

		// remove pre1, add pre2
		diff -= ps.CountLower(pos1+1, pos2+1, pre1)
		ps.Set(pos1, pre2)
		diff += ps.CountLower(pos1+1, pos2+1, pre2)

		// remove pre2, add pre1
		diff -= ps.CountHigher(pos1, pos2, pre2)
		ps.Set(pos2, pre1)
		diff += ps.CountHigher(pos1, pos2, pre1)

		res[i+1] = res[i] + diff
	}

	return res
}

func InversePairsNaive(nums []int, operations [][2]int) []int {
	res := make([]int, len(operations)+1)
	res[0] = countInvMergeSort(nums)
	for i, operation := range operations {
		pos1, pos2 := operation[0], operation[1]
		nums[pos1], nums[pos2] = nums[pos2], nums[pos1]
		res[i+1] = countInvMergeSort(nums)
	}
	return res
}

func main() {
	// debug()
	P1975()
}

func debug() {
	// rand.Seed(112)
	// nums := make([]int, 3)
	// for i := range nums {
	// 	nums[i] = rand.Intn(3)
	// }
	// operations := make([][2]int, 3)
	// for i := range operations {
	// 	operations[i] = [2]int{rand.Intn(3), rand.Intn(3)}
	// }
	// fmt.Println(nums, operations)
	// res1 := InversePairs(append(nums[:0:0], nums...), operations)
	// res2 := InversePairsNaive(append(nums[:0:0], nums...), operations)
	// for i := range res1 {
	// 	if res1[i] != res2[i] {
	// 		fmt.Println(i, res1[i], res2[i], res1, res2)
	// 		panic("")
	// 	}
	// }
	// [0 1 1] [[0 1] [1 0] [0 2]]  ->  [0 1 0 2]
	fmt.Println(InversePairs([]int{0, 1, 1}, [][2]int{{0, 1}, {1, 0}, {0, 2}}))
}

func P1975() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)
	operations := make([][2]int, q)
	for i := range operations {
		var pos1, pos2 int
		fmt.Fscan(in, &pos1, &pos2)
		pos1--
		pos2--
		operations[i] = [2]int{pos1, pos2}
	}

	res := InversePairs(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 归并排序求逆序对.
func countInvMergeSort(arr []int) int {
	if len(arr) < 2 {
		return 0
	}
	if len(arr) == 2 {
		if arr[0] > arr[1] {
			return 1
		}
		return 0
	}
	res := 0
	midCount := 0
	upper := make([]int, 0)
	lower := make([]int, 0)
	mid := arr[len(arr)/2]
	for i := 0; i < len(arr); i++ {
		num := arr[i]
		if num < mid {
			lower = append(lower, num)
			res += len(upper)
			res += midCount
		} else if num > mid {
			upper = append(upper, num)
		} else {
			midCount++
			res += len(upper)
		}
	}
	res += countInvMergeSort(lower)
	res += countInvMergeSort(upper)
	return res
}

// 单点修改，区间频率查询.
//
//	单次修改复杂度 O(sqrt(n)), 单次查询复杂度 O(sqrt(n) * log(sqrt(n))).
type _PointSetRangeFreq struct {
	nums        []int
	belong      []int
	blockStart  []int
	blockEnd    []int
	blockCount  int
	blockSorted [][]int
}

// ps := _NewPointSetRangeFreq(arr, 2*int(math.Sqrt(float64(len(arr)))+1))
func _NewPointSetRangeFreq(nums []int, blockSize int) *_PointSetRangeFreq {
	nums = append(nums[:0:0], nums...)
	res := &_PointSetRangeFreq{}
	block := UseBlock(nums, blockSize)
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

func (ps *_PointSetRangeFreq) Set(pos, newValue int) {
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
func (ps *_PointSetRangeFreq) Count(start, end int, target int) int {
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
func (ps *_PointSetRangeFreq) CountLower(start, end int, target int) int {
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
func (ps *_PointSetRangeFreq) CountFloor(start, end int, target int) int {
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
func (ps *_PointSetRangeFreq) CountCeiling(start, end int, target int) int {
	return (end - start) - ps.CountLower(start, end, target)
}

// 统计 [start, end) 中严格大于 target 的元素个数.
func (ps *_PointSetRangeFreq) CountHigher(start, end int, target int) int {
	return (end - start) - ps.CountFloor(start, end, target)
}

func (ps *_PointSetRangeFreq) _bisectLeft(nums []int, target int, left, right int) int {
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

func (ps *_PointSetRangeFreq) _bisectRight(nums []int, target int, left, right int) int {
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

func (ps *_PointSetRangeFreq) _count(nums []int, target int, left, right int) int {
	return ps._bisectRight(nums, target, left, right) - ps._bisectLeft(nums, target, left, right)
}

// blockSize := int(math.Sqrt(float64(len(nums))) + 1)
func UseBlock(nums []int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	n := len(nums)
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
