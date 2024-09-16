// Api:
//  GetMedian(sortedNums []int, start, end int) int // 求有序数组区间中位数.
//  DistSum(sortedNums []int) func(k int) int // 有序数组中所有点到`x=k`的距离之和.
//  DistSumRange(sortedNums []int) func(k, start, end int) int // 有序数组区间所有点到`x=k`的距离之和.
//  DistSumOfAllPairs(sortedNums []int) int // 有序数组中所有点对两两距离之和.
//  DistSumOfAllPairsRange(sortedNums []int, start, end int) int // 有序数组区间中所有点对两两距离之和.
//  KthNearest(sortedNums []int, pos int, k int) int // 距离x=pos第k(k>=1)近的数到pos的距离.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// CF371E()
	abc364D()
}

// Subway Innovation
// https://www.luogu.com.cn/problem/CF371E
// 在直线上给定 n 个点，定义其距离为横坐标的绝对值，请你保留 k 个点，使这些点两两之间的距离和最小，输出这 k 个点的坐标。
func CF371E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	var k int
	fmt.Fscan(in, &k)

	points := make([][2]int, 0, k)
	for i := 0; i < n; i++ {
		points = append(points, [2]int{nums[i], i})
	}
	sort.Slice(points, func(i, j int) bool { return points[i][0] < points[j][0] })
	sortedNums := make([]int, n)
	for i := 0; i < n; i++ {
		sortedNums[i] = points[i][0]
	}

	D := DistSumRange(sortedNums)
	curSum := DistSumOfAllPairsRange(sortedNums, 0, k)
	bestLeft, bestSum := 0, curSum
	for left := 1; left+k-1 < n; left++ {
		right := left + k - 1
		// reove left-1, add right
		curSum -= D(sortedNums[left-1], left-1, right)
		curSum += D(sortedNums[right], left, right+1)
		if curSum < bestSum {
			bestLeft = left
			bestSum = curSum
		}
	}

	for i := bestLeft; i < bestLeft+k; i++ {
		fmt.Fprint(out, points[i][1]+1, " ")
	}
}

// D - K-th Nearest
// https://atcoder.jp/contests/abc364/tasks/abc364_d
func abc364D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	sort.Ints(nums)

	for i := 0; i < q; i++ {
		var pos, k int
		fmt.Fscan(in, &pos, &k)
		fmt.Fprintln(out, KthNearest(nums, pos, k))
	}
}

// 100123. 执行操作使频率分数最大
func maxFrequencyScore(nums []int, k int64) int {
	sort.Ints(nums)
	D := DistSumRange(nums)
	res, left := 0, 0
	for right := 0; right < len(nums); right++ {
		for left <= right {
			median := GetMedian(nums, left, right+1)
			if D(median, left, right+1) <= int(k) {
				break
			}
			left++
		}
		res = max(res, right-left+1)
	}
	return res
}

// 求有序数组中位数(向下取整).
func GetMedian(sortedNums []int, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > len(sortedNums) {
		end = len(sortedNums)
	}
	if start >= end {
		return 0
	}
	if (end-start)&1 == 0 {
		return (sortedNums[(end+start)/2-1] + sortedNums[(end+start)/2]) / 2
	}
	return sortedNums[(end+start)/2]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 有序数组中所有点到`x=k`的距离之和.
func DistSum(sortedNums []int) func(k int) int {
	n := len(sortedNums)
	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + sortedNums[i]
	}

	return func(k int) int {
		pos := sort.SearchInts(sortedNums, k+1)
		leftSum := k*pos - preSum[pos]
		rightSum := preSum[n] - preSum[pos] - k*(n-pos)
		return leftSum + rightSum
	}
}

// 有序数组切片[start:end)中所有点到`x=k`的距离之和.
func DistSumRange(sortedNums []int) func(k int, start, end int) int {
	n := len(sortedNums)
	preSum := make([]int, n+1)
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + sortedNums[i]
	}

	return func(k, start, end int) int {
		if start < 0 {
			start = 0
		}
		if end > n {
			end = n
		}
		if start >= end {
			return 0
		}
		pos := sort.SearchInts(sortedNums, k)
		if pos <= start {
			return (preSum[end] - preSum[start]) - k*(end-start)
		}
		if pos >= end {
			return k*(end-start) - (preSum[end] - preSum[start])
		}
		leftSum := k*(pos-start) - (preSum[pos] - preSum[start])
		rightSum := preSum[end] - preSum[pos] - k*(end-pos)
		return leftSum + rightSum
	}
}

// 有序数组中所有点对两两距离之和.一共有`n*(n-1)//2`对点对.
func DistSumOfAllPairs(sortedNums []int) int {
	n := len(sortedNums)
	res := 0
	preSum := 0
	for i := 0; i < n; i++ {
		res += sortedNums[i]*i - preSum
		preSum += sortedNums[i]
	}
	return res
}

func DistSumOfAllPairsRange(sortedNums []int, start, end int) int {
	res := 0
	preSum := 0
	for i := start; i < end; i++ {
		res += sortedNums[i]*i - preSum
		preSum += sortedNums[i]
	}
	return res
}

// 距离x=pos第k(k>=1)近的数到pos的距离.
// 二分答案.
func KthNearest(sortedNums []int, pos int, k int) int {
	// 离pos距离为d以内的数的个数>=k
	check := func(d int) bool {
		left, right := pos-d, pos+d
		count := sort.SearchInts(sortedNums, right+1) - sort.SearchInts(sortedNums, left)
		return count >= k
	}

	upper := max(abs(sortedNums[0]-pos), abs(sortedNums[len(sortedNums)-1]-pos))
	return MinLeft(upper+1, check, 0)
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

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
