// api:
//	NewIncreasingArray(increasingArray []int) *IncreasingArray
//	Increase(k int) (value, pos int)
//	IncreaseForMin(k int) int
//	IncreaseForArray(k int) []int
//	Decrease(k int) (value, pos int)
//	DecreaseForMax(k int) int
//	DecreaseForArray(k int) []int
//	SumWithUpClamp(v int) int
//	SumWithLowClamp(v int) int
//	DiffSum(v int) int

package main

import (
	"fmt"
	"sort"
)

func main() {
	demo()
}

func demo() {
	A := NewIncreasingArray([]int{1, 2, 3, 4, 5})
	for i := 1; i <= 24; i++ {
		fmt.Println(A.Increase(i))
	}
	fmt.Println("------")
	for i := 0; i < 10; i++ {
		fmt.Println(A.SumWithUpClamp(i))
	}
	fmt.Println("------")
	for i := 0; i < 10; i++ {
		fmt.Println(A.SumWithLowClamp(i))
	}
	fmt.Println("------")
	for i := 0; i < 10; i++ {
		fmt.Println(A.DiffSum(i))
	}
	fmt.Println("------")
	for i := 1; i <= 24; i++ {
		fmt.Println(A.Decrease(i))
	}
	fmt.Println("------")
	for i := 0; i < 200; i++ {
		fmt.Println(i, A.IncreaseForArray(i))
		assert(mins(A.IncreaseForArray(i)...) == A.IncreaseForMin(i), "min")
	}
	fmt.Println("------")
	for i := 0; i < 200; i++ {
		fmt.Println(i, A.DecreaseForArray(i))
		assert(maxs(A.DecreaseForArray(i)...) == A.DecreaseForMax(i), "max")
	}
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

type IncreasingArray struct {
	arr    []int
	presum []int
}

func NewIncreasingArray(increasingArray []int) *IncreasingArray {
	if !sort.IntsAreSorted(increasingArray) {
		panic("input array should be increasing")
	}
	presum := make([]int, len(increasingArray)+1)
	for i := 0; i < len(increasingArray); i++ {
		presum[i+1] = presum[i] + increasingArray[i]
	}
	return &IncreasingArray{arr: increasingArray, presum: presum}
}

// 每次选取最矮的矩形中编号最小的，并把它的高度+1
// 返回第K次选取的矩形的高度和编号.
// k>=1.
func (a *IncreasingArray) Increase(k int) (value, pos int) {
	if k <= 0 {
		panic("k should be >=1")
	}
	// !二分无法与哪个数齐平
	right := MaxRight(0, func(r int) bool { return a.arr[r-1]*r-a.presum[r] < k }, len(a.arr))
	filled := a.arr[right-1]*right - a.presum[right]
	remain := k - filled
	div, mod := remain/right, remain%right
	value = a.arr[right-1] + div
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
		return a.arr[0]
	}
	right := MaxRight(0, func(r int) bool { return a.arr[r-1]*r-a.presum[r] < k }, len(a.arr))
	filled := a.arr[right-1]*right - a.presum[right]
	remain := k - filled
	return a.arr[right-1] + remain/right
}

// 每次选取最矮的矩形中编号最小的，并把它的高度+1
// 返回操作后的数组.
func (a *IncreasingArray) IncreaseForArray(k int) []int {
	if k <= 0 {
		return a.arr
	}
	right := MaxRight(0, func(r int) bool { return a.arr[r-1]*r-a.presum[r] < k }, len(a.arr))
	filled := a.arr[right-1]*right - a.presum[right]
	remain := k - filled
	div, mod := remain/right, remain%right
	base := a.arr[right-1] + div
	res := append(a.arr[:0:0], a.arr...)
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
		len(a.arr),
		func(l int) bool {
			return (a.presum[len(a.arr)]-a.presum[l])-a.arr[l]*(len(a.arr)-l) < k
		},
		0,
	)
	filled := a.presum[len(a.arr)] - a.presum[left] - a.arr[left]*(len(a.arr)-left)
	remain := k - filled
	div, mod := remain/(len(a.arr)-left), remain%(len(a.arr)-left)
	value = a.arr[left] - div
	if mod > 0 {
		value--
	}
	pos = mod - 1
	if pos < 0 {
		pos += len(a.arr) - left
	}
	pos += left
	return
}

func (a *IncreasingArray) DecreaseForMax(k int) int {
	if k <= 0 {
		return a.arr[len(a.arr)-1]
	}
	left := MinLeft(
		len(a.arr),
		func(l int) bool {
			return (a.presum[len(a.arr)]-a.presum[l])-a.arr[l]*(len(a.arr)-l) < k
		},
		0,
	)
	filled := a.presum[len(a.arr)] - a.presum[left] - a.arr[left]*(len(a.arr)-left)
	remain := k - filled
	return a.arr[left] - remain/(len(a.arr)-left)
}

// 每次选取最高的矩形中编号最小的，并把它的高度-1
// 返回操作后的数组.
func (a *IncreasingArray) DecreaseForArray(k int) []int {
	if k <= 0 {
		return a.arr
	}
	left := MinLeft(
		len(a.arr),
		func(l int) bool {
			return (a.presum[len(a.arr)]-a.presum[l])-a.arr[l]*(len(a.arr)-l) < k
		},
		0,
	)
	filled := a.presum[len(a.arr)] - a.presum[left] - a.arr[left]*(len(a.arr)-left)
	remain := k - filled
	div, mod := remain/(len(a.arr)-left), remain%(len(a.arr)-left)
	base := a.arr[left] - div
	res := append(a.arr[:0:0], a.arr...)
	for i := left; i < len(a.arr); i++ {
		res[i] = base
		if i-left < mod {
			res[i]--
		}
	}
	return res
}

// 求所有数与v取min的和.
func (a *IncreasingArray) SumWithUpClamp(v int) int {
	pos := sort.SearchInts(a.arr, v)
	return a.presum[pos] + (len(a.arr)-pos)*v
}

// 求所有数与v取max的和.
func (a *IncreasingArray) SumWithLowClamp(v int) int {
	pos := sort.SearchInts(a.arr, v)
	return pos*v + a.presum[len(a.arr)] - a.presum[pos]
}

// 求所有数与v的绝对值差的和.
func (a *IncreasingArray) DiffSum(v int) int {
	pos := sort.SearchInts(a.arr, v)
	n := len(a.arr)
	leftSum := v*pos - a.presum[pos]
	rightSum := a.presum[n] - a.presum[pos] - v*(n-pos)
	return leftSum + rightSum
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
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
