package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	abc294f()
}

// https://atcoder.jp/contests/abc294/tasks/abc294_f
// 两排糖水，糖成份x和水成分 y。
// !从两排中各取一个糖水混合，得糖浓度 (x1+x2)/(x1+y1+x2+y2)。
// 问所有方案中，糖浓度第 k大是多少(1<=k<=m*n)。
// m,n<=5e4
// !01分数规划=> 二分答案+公式变形
// (xi+xj)<=mid*(xi+yi+xj+yj)
// !将i和j分离得到
// (mid-1)*xi +mid*yi + (mid-1)*xj +mid*yj >=0
// !把(mid-1)*x+mid*y作为新的数组,也就是求两个数组中和大于等于0的对数 (排序+双指针或者二分都以)
func abc294f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	pairs1 := make([][2]int, n)
	for i := range pairs1 {
		fmt.Fscan(in, &pairs1[i][0], &pairs1[i][1])
	}
	pairs2 := make([][2]int, m)
	for i := range pairs2 {
		fmt.Fscan(in, &pairs2[i][0], &pairs2[i][1])
	}

	kthMin := m*n + 1 - k
	countNgt := func(mid float64) int {
		A := make([]float64, len(pairs1))
		for i, pair := range pairs1 {
			a, b := float64(pair[0]), float64(pair[1])
			A[i] = (mid-1)*a + mid*b
		}
		sort.Float64Slice(A).Sort()

		B := make([]float64, len(pairs2))
		for i, pair := range pairs2 {
			a, b := float64(pair[0]), float64(pair[1])
			B[i] = (mid-1)*a + mid*b
		}
		sort.Float64Slice(B).Sort()

		res, right := 0, len(B)-1
		for _, a := range A {
			for right >= 0 && a+B[right] >= 0 {
				right--
			}
			res += len(B) - 1 - right
		}
		return res
	}

	percent := GetKth1Float64(0, 1, countNgt, kthMin) * 100
	fmt.Fprintln(out, percent)
}

// 378. 有序矩阵中第 K 小的元素
// https://leetcode.cn/problems/kth-smallest-element-in-a-sorted-matrix/description/
func kthSmallest(matrix [][]int, k int) int {
	countNgt := func(mid int) int {
		res := 0
		for _, row := range matrix {
			res += sort.Search(len(row), func(i int) bool { return row[i] > mid })
		}
		return res
	}

	left, right := matrix[0][0], matrix[len(matrix)-1][len(matrix[0])-1]
	return GetKth0(left, right, countNgt, k-1)
}

// 668. 乘法表中第k小的数
// https://leetcode.cn/problems/kth-smallest-number-in-multiplication-table/description/
func findKthNumber(m int, n int, k int) int {
	countNgt := func(mid int) int {
		res := 0
		for i := 1; i <= m; i++ {
			res += min(mid/i, n)
		}
		return res
	}

	left, right := 1, m*n
	return GetKth1(left, right, countNgt, k)
}

// 719. 找出第 K 小的距离对
// https://leetcode.cn/problems/find-k-th-smallest-pair-distance/
func smallestDistancePair(nums []int, k int) int {
	sort.Ints(nums)

	// 有多少个数对的差值不大于mid.
	countNgt := func(mid int) int {
		res, left := 0, 0
		for right, num := range nums {
			for left <= right && num-nums[left] > mid {
				left++
			}
			res += right - left
		}
		return res
	}

	left, right := 0, nums[len(nums)-1]-nums[0]
	// return GetKth0(left, right, countNgt, k-1)
	return GetKth1(left, right, countNgt, k)
}

// 给定二分答案的区间[left,right]，求第kth小的答案.
// countNgt: 答案不超过mid时，满足条件的个数.
// kth从0开始.
func GetKth0(left int, right int, countNgt func(mid int) int, kth int) int {
	for left <= right {
		mid := left + (right-left)>>1
		if countNgt(mid) <= kth {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return right + 1
}

// 给定二分答案的区间[left,right]，求第kth小的答案.
// countNgt: 答案不超过mid时，满足条件的个数.
// kth从1开始.
func GetKth1(left int, right int, countNgt func(mid int) int, kth int) int {
	for left <= right {
		mid := left + (right-left)>>1
		if countNgt(mid) < kth {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

const EPS = 1e-12

// 给定二分答案的区间[left,right]，求第kth小的答案.
// countNgt: 答案不超过mid时，满足条件的个数.
// kth从0开始.
func GetKth0Float64(left float64, right float64, countNgt func(mid float64) int, kth int) float64 {
	for left <= right {
		mid := left + (right-left)/2
		if countNgt(mid) <= kth {
			left = mid + EPS
		} else {
			right = mid - EPS
		}
	}
	return right + EPS
}

// 给定二分答案的区间[left,right]，求第kth小的答案.
// countNgt: 答案不超过mid时，满足条件的个数.
// kth从1开始.
func GetKth1Float64(left float64, right float64, countNgt func(mid float64) int, kth int) float64 {
	for left <= right {
		mid := left + (right-left)/2
		if countNgt(mid) < kth {
			left = mid + EPS
		} else {
			right = mid - EPS
		}
	}
	return left
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
