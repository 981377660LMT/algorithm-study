// // 区间和的第 k 小（数组元素均为非负）
// // 每个区间和可以视作一个有序上三角矩阵中的元素，在数组元素均为非负时，该矩阵从左往右和从下往上均为非降序列
// // 1508 https://leetcode-cn.com/problems/range-sum-of-sorted-subarray-sums/
// kthSmallestRangeSum := func(a []int, k int) int {
// 	// 1 <= k <= n*(n+1)/2
// 	n := len(a)
// 	sum := make([]int, n+1)
// 	for i, v := range a {
// 		sum[i+1] = sum[i] + v
// 	}
// 	ans := sort.Search(sum[n], func(v int) bool {
// 		cnt := 0
// 		for l, r := 0, 1; r <= n; {
// 			if v < sum[r]-sum[l] {
// 				l++
// 			} else {
// 				cnt += r - l
// 				r++
// 			}
// 		}
// 		return cnt >= k
// 	})
// 	return ans
// }

package main

import "sort"

// https://github.com/EndlessCheng/codeforces-go/blob/04c3e9c75964d0a5dd701db92e5c3fdf0bc353bc/copypasta/sort.go#L381
// 区间和的第 k 小（数组元素均为非负,1<=k<=n*(n+1)/2
// !每个区间和可以视作一个有序上三角矩阵中的元素，在数组元素均为非负时，该矩阵从左往右和从下往上均为非降序列
// 1508 https://leetcode-cn.com/problems/range-sum-of-sorted-subarray-sums/

func KthSubarraySum(nums []int, k int) int {
	n := len(nums)
	preSum := make([]int, n+1)
	for i, v := range nums {
		preSum[i+1] = preSum[i] + v
	}
	res := sort.Search(preSum[n], func(s int) bool {
		count := 0
		for l, r := 0, 1; r <= n; {
			if s < preSum[r]-preSum[l] {
				l++
			} else {
				count += r - l
				r++
			}
		}
		return count >= k
	})
	return res
}
