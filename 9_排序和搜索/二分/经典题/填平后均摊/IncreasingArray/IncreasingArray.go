// IncreasingArrayUtils/SortedArrayUtils
//
// api:
//	NewIncreasingArray(increasingArray []int) *IncreasingArray
//	Increase(k int) (value, pos int)
//	IncreaseForMin(k int) int
//	IncreaseForArray(k int) []int
//	Decrease(k int) (value, pos int)
//	DecreaseForMax(k int) int
//	DecreaseForArray(k int) []int
//
//	SumWithUpClamp(v int) int
//  SumWithUpClampRange(v int, start, end int) int
//	SumWithLowClamp(v int) int
//  SumWithLowClampRange(v int, start, end int) int
//  SumWithUpAndLowClamp(low, up int) int
//  SumWithUpAndLowClampRange(low, up int, start, end int) int
//
//	DiffSum(v int) int
//  DiffSumRange(v int, start, end int) int
//
//	CountRange(start, end int, y1, y2 int) int
//	SumRange(start, end int, y1, y2 int) int
//	CountAndSumRange(start, end int, y1, y2 int) (int, int)

package main

import (
	"fmt"
	"sort"
)

func main() {
	test()
}

// 2233. K 次增加后的最大乘积
// https://leetcode.cn/problems/maximum-product-after-k-increments/description/
func maximumProduct(nums []int, k int) int {
	mod := int(1e9 + 7)
	sort.Ints(nums)
	A := NewIncreasingArray(nums)
	arr := A.IncreaseForArray(k)
	res := 1
	for _, v := range arr {
		res = res * v % mod
	}
	return res
}

// 2333. 最小差值平方和
// https://leetcode.cn/problems/minimum-sum-of-squared-difference/description/
func minSumSquareDiff(nums1 []int, nums2 []int, k1 int, k2 int) int64 {
	diff := make([]int, len(nums1))
	for i := 0; i < len(nums1); i++ {
		diff[i] = abs(nums1[i] - nums2[i])
	}
	sort.Ints(diff)
	A := NewIncreasingArray(diff)
	arr := A.DecreaseForArray(k1 + k2)
	for i := range arr {
		if arr[i] < 0 {
			arr[i] = 0
		}
	}
	res := 0
	for i := 0; i < len(arr); i++ {
		res += arr[i] * arr[i]
	}
	return int64(res)
}

// 3107. 使数组中位数等于 K 的最少操作数
// https://leetcode.cn/problems/minimum-operations-to-make-median-of-array-equal-to-k/
// 一次操作中，你可以选择任一元素 加 1 或者减 1 。
// 请你返回将 nums 中位数 变为 k 所需要的 最少 操作次数。
// 把中位数左边的数(含自己)都变成 ≤k 的，右边的数(含自己)都变成 ≥k 的
func minOperationsToMakeMedianK(nums []int, k int) int64 {
	const INF int = 1e18
	sort.Ints(nums)
	m := len(nums) >> 1
	A := NewIncreasingArray(nums)

	res := 0
	largerCount, largerSum := A.CountAndSumRange(0, m+1, k+1, INF)
	res += (largerSum - k*largerCount)
	lessCount, lessSum := A.CountAndSumRange(m, len(nums), -INF, k)
	res += (k*lessCount - lessSum)
	return int64(res)
}

type IncreasingArray struct {
	Arr    []int
	Presum []int
}

func NewIncreasingArray(increasingArray []int) *IncreasingArray {
	if !sort.IntsAreSorted(increasingArray) {
		panic("input array should be increasing")
	}
	presum := make([]int, len(increasingArray)+1)
	for i := 0; i < len(increasingArray); i++ {
		presum[i+1] = presum[i] + increasingArray[i]
	}
	return &IncreasingArray{Arr: increasingArray, Presum: presum}
}

// 每次选取最矮的矩形中编号最小的，并把它的高度+1
// 返回第K次选取的矩形的高度和编号.
// k>=1.
func (a *IncreasingArray) Increase(k int) (value, pos int) {
	if k <= 0 {
		panic("k should be >=1")
	}
	// !二分无法与哪个数齐平
	right := MaxRight(0, func(r int) bool { return a.Arr[r-1]*r-a.Presum[r] < k }, len(a.Arr))
	filled := a.Arr[right-1]*right - a.Presum[right]
	remain := k - filled
	div, mod := remain/right, remain%right
	value = a.Arr[right-1] + div
	if mod > 0 {
		value++
	}
	pos = mod - 1
	if pos < 0 {
		pos += right
	}
	return
}

func (a *IncreasingArray) IncreaseForMin(k int) int {
	if k <= 0 {
		return a.Arr[0]
	}
	right := MaxRight(0, func(r int) bool { return a.Arr[r-1]*r-a.Presum[r] < k }, len(a.Arr))
	filled := a.Arr[right-1]*right - a.Presum[right]
	remain := k - filled
	return a.Arr[right-1] + remain/right
}

// 每次选取最矮的矩形中编号最小的，并把它的高度+1
// 返回操作后的数组.
func (a *IncreasingArray) IncreaseForArray(k int) []int {
	if k <= 0 {
		return a.Arr
	}
	right := MaxRight(0, func(r int) bool { return a.Arr[r-1]*r-a.Presum[r] < k }, len(a.Arr))
	filled := a.Arr[right-1]*right - a.Presum[right]
	remain := k - filled
	div, mod := remain/right, remain%right
	base := a.Arr[right-1] + div
	res := append(a.Arr[:0:0], a.Arr...)
	for i := 0; i < right; i++ {
		res[i] = base
		if i < mod {
			res[i]++
		}
	}
	return res
}

// 每次选取最高的矩形中编号最小的，并把它的高度-1
// 返回第K次选取的矩形的高度和编号.
// k>=1.
func (a *IncreasingArray) Decrease(k int) (value, pos int) {
	if k <= 0 {
		panic("k should be >=1")
	}
	left := MinLeft(
		len(a.Arr),
		func(l int) bool {
			return (a.Presum[len(a.Arr)]-a.Presum[l])-a.Arr[l]*(len(a.Arr)-l) < k
		},
		0,
	)
	filled := a.Presum[len(a.Arr)] - a.Presum[left] - a.Arr[left]*(len(a.Arr)-left)
	remain := k - filled
	div, mod := remain/(len(a.Arr)-left), remain%(len(a.Arr)-left)
	value = a.Arr[left] - div
	if mod > 0 {
		value--
	}
	pos = mod - 1
	if pos < 0 {
		pos += len(a.Arr) - left
	}
	pos += left
	return
}

func (a *IncreasingArray) DecreaseForMax(k int) int {
	if k <= 0 {
		return a.Arr[len(a.Arr)-1]
	}
	left := MinLeft(
		len(a.Arr),
		func(l int) bool {
			return (a.Presum[len(a.Arr)]-a.Presum[l])-a.Arr[l]*(len(a.Arr)-l) < k
		},
		0,
	)
	filled := a.Presum[len(a.Arr)] - a.Presum[left] - a.Arr[left]*(len(a.Arr)-left)
	remain := k - filled
	return a.Arr[left] - remain/(len(a.Arr)-left)
}

// 每次选取最高的矩形中编号最小的，并把它的高度-1
// 返回操作后的数组.
func (a *IncreasingArray) DecreaseForArray(k int) []int {
	if k <= 0 {
		return a.Arr
	}
	left := MinLeft(
		len(a.Arr),
		func(l int) bool {
			return (a.Presum[len(a.Arr)]-a.Presum[l])-a.Arr[l]*(len(a.Arr)-l) < k
		},
		0,
	)
	filled := a.Presum[len(a.Arr)] - a.Presum[left] - a.Arr[left]*(len(a.Arr)-left)
	remain := k - filled
	div, mod := remain/(len(a.Arr)-left), remain%(len(a.Arr)-left)
	base := a.Arr[left] - div
	res := append(a.Arr[:0:0], a.Arr...)
	for i := left; i < len(a.Arr); i++ {
		res[i] = base
		if i-left < mod {
			res[i]--
		}
	}
	return res
}

// 求所有数与v取min的和.
func (a *IncreasingArray) SumWithUpClamp(v int) int {
	pos := sort.SearchInts(a.Arr, v)
	return a.Presum[pos] + (len(a.Arr)-pos)*v
}

func (a *IncreasingArray) SumWithUpClampRange(v int, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0
	}
	lessCount := sort.SearchInts(a.Arr[start:end], v)
	lessSum := a.Presum[start+lessCount] - a.Presum[start]
	return lessSum + v*(end-start-lessCount)
}

// 求所有数与v取max的和.
func (a *IncreasingArray) SumWithLowClamp(v int) int {
	pos := sort.SearchInts(a.Arr, v)
	return pos*v + a.Presum[len(a.Arr)] - a.Presum[pos]
}

func (a *IncreasingArray) SumWithLowClampRange(v int, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0
	}
	lessCount := sort.SearchInts(a.Arr[start:end], v)
	largerSum := a.Presum[end] - a.Presum[start+lessCount]
	return v*lessCount + largerSum
}

func (a *IncreasingArray) SumWithUpAndLowClamp(low, up int) int {
	if low > up {
		return 0
	}
	posLow := sort.SearchInts(a.Arr, low)
	posUp := sort.SearchInts(a.Arr, up)
	return a.Presum[posUp] - a.Presum[posLow] + low*posLow + up*(len(a.Arr)-posUp)
}

func (a *IncreasingArray) SumWithUpAndLowClampRange(low, up int, start, end int) int {
	if low > up {
		return 0
	}
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0
	}
	posLow := sort.SearchInts(a.Arr[start:end], low)
	posUp := sort.SearchInts(a.Arr[start:end], up)
	return a.Presum[start+posUp] - a.Presum[start+posLow] + low*posLow + up*(end-start-posUp)
}

// 求所有数与v的绝对值差的和.
func (a *IncreasingArray) DiffSum(v int) int {
	pos := sort.SearchInts(a.Arr, v)
	n := len(a.Arr)
	leftSum := v*pos - a.Presum[pos]
	rightSum := a.Presum[n] - a.Presum[pos] - v*(n-pos)
	return leftSum + rightSum
}

func (a *IncreasingArray) DiffSumRange(v int, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0
	}
	presum := a.Presum
	pos := sort.SearchInts(a.Arr, v)
	if pos <= start {
		return (presum[end] - presum[start]) - v*(end-start)
	}
	if pos >= end {
		return v*(end-start) - (presum[end] - presum[start])
	}
	leftSum := v*(pos-start) - (presum[pos] - presum[start])
	rightSum := presum[end] - presum[pos] - v*(end-pos)
	return leftSum + rightSum
}

// !WaveletMatrixLike Api

// [start,end) x [y1,y2) 中的数的个数.
func (a *IncreasingArray) CountRange(start, end int, y1, y2 int) int {
	count, _ := a.CountAndSumRange(start, end, y1, y2)
	return count
}

// [start,end) x [y1,y2) 中的数的和.
func (a *IncreasingArray) SumRange(start, end int, y1, y2 int) int {
	_, sum := a.CountAndSumRange(start, end, y1, y2)
	return sum
}

// [start,end) x [y1,y2) 中的数的个数、和.
func (a *IncreasingArray) CountAndSumRange(start, end int, y1, y2 int) (int, int) {
	if y1 >= y2 {
		return 0, 0
	}
	if start < 0 {
		start = 0
	}
	if end > len(a.Arr) {
		end = len(a.Arr)
	}
	if start >= end {
		return 0, 0
	}
	nums := a.Arr[start:end]
	left := sort.SearchInts(nums, y1)
	right := sort.SearchInts(nums, y2)
	return right - left, a.Presum[start+right] - a.Presum[start+left]
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含，使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func test() {
	nums := []int{1, 2, 3, 4, 5}
	A := NewIncreasingArray(nums)

	// for i := 1; i <= 24; i++ {
	// 	fmt.Println(A.Increase(i))
	// }
	// fmt.Println("------")
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(A.SumWithUpClamp(i))
	// }
	// fmt.Println("------")
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(A.SumWithLowClamp(i))
	// }
	// fmt.Println("------")
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(A.DiffSum(i))
	// }
	// fmt.Println("------")
	// for i := 1; i <= 24; i++ {
	// 	fmt.Println(A.Decrease(i))
	// }
	// fmt.Println("------")
	// for i := 0; i < 200; i++ {
	// 	fmt.Println(i, A.IncreaseForArray(i))
	// 	assert(mins(A.IncreaseForArray(i)...) == A.IncreaseForMin(i), "min")
	// }
	// fmt.Println("------")
	// for i := 0; i < 200; i++ {
	// 	fmt.Println(i, A.DecreaseForArray(i))
	// 	assert(maxs(A.DecreaseForArray(i)...) == A.DecreaseForMax(i), "max")
	// }

	upClampRangeBruteForce := func(v int, start, end int) int {
		sum := 0
		for i := start; i < end; i++ {
			sum += min(nums[i], v)
		}
		return sum
	}

	lowClampRangeBruteForce := func(v int, start, end int) int {
		sum := 0
		for i := start; i < end; i++ {
			sum += max(nums[i], v)
		}
		return sum
	}

	upAndLowClampBruteForce := func(low, up int) int {
		sum := 0
		for i := 0; i < len(nums); i++ {
			sum += max(min(nums[i], up), low)
		}
		return sum
	}

	upAndLowClampRangeBruteForce := func(low, up int, start, end int) int {
		sum := 0
		for i := start; i < end; i++ {
			sum += max(min(nums[i], up), low)
		}
		return sum
	}

	countRangeBruteForce := func(start, end int, y1, y2 int) int {
		sum := 0
		for i := start; i < end; i++ {
			if y1 <= nums[i] && nums[i] < y2 {
				sum++
			}

		}
		return sum
	}

	countAndSumRangeBruteForce := func(start, end int, y1, y2 int) (int, int) {
		count, sum := 0, 0
		for i := start; i < end; i++ {
			if y1 <= nums[i] && nums[i] < y2 {
				count++
				sum += nums[i]
			}
		}
		return count, sum
	}

	{
		for i := 0; i < len(nums); i++ {
			for j := 0; j < len(nums); j++ {
				for v := -10; v < 10; v++ {
					assert(A.SumWithUpClampRange(v, i, j) == upClampRangeBruteForce(v, i, j), "upClampRange")
					assert(A.SumWithLowClampRange(v, i, j) == lowClampRangeBruteForce(v, i, j), "lowClampRange")

				}
			}
		}

		for min := -10; min < 10; min++ {
			for max := min; max < 10; max++ {
				assert(A.SumWithUpAndLowClamp(min, max) == upAndLowClampBruteForce(min, max), "upAndLowClamp")
			}
		}

		for i := 0; i < len(nums); i++ {
			for j := i; j <= len(nums); j++ {
				for y1 := -10; y1 < 10; y1++ {
					for y2 := y1; y2 < 10; y2++ {
						assert(A.SumWithUpAndLowClampRange(y1, y2, i, j) == upAndLowClampRangeBruteForce(y1, y2, i, j), "upAndLowClampRange")
					}
				}
			}
		}

		for i := 0; i < len(nums); i++ {
			for j := i; j <= len(nums); j++ {
				for y1 := -10; y1 < 10; y1++ {
					for y2 := y1; y2 < 10; y2++ {
						assert(A.CountRange(i, j, y1, y2) == countRangeBruteForce(i, j, y1, y2), "countRange")
					}
				}
			}
		}
	}

	{
		for i := 0; i < len(nums); i++ {
			for j := i; j <= len(nums); j++ {
				for y1 := -10; y1 < 10; y1++ {
					for y2 := y1; y2 < 10; y2++ {
						count, sum := A.CountAndSumRange(i, j, y1, y2)
						countBrute, sumBrute := countAndSumRangeBruteForce(i, j, y1, y2)
						assert(count == countBrute, "countAndSumRange count")
						assert(sum == sumBrute, "countAndSumRange sum")
					}
				}
			}
		}
	}

	fmt.Println("pass")

}
