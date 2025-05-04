package main

import "fmt"

func main() {
	fmt.Println(KthSmallestOfTwoSortedArrays([]int{1, 3, 5, 7, 9}, []int{2, 4, 6, 8, 10}, 5)) // 5
}

// 寻找两个正序数组的中位数
// https://leetcode.cn/problems/median-of-two-sorted-arrays/
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	n1, n2 := len(nums1), len(nums2)
	if (n1+n2)&1 == 1 {
		return float64(KthSmallestOfTwoSortedArrays(nums1, nums2, (n1+n2)>>1))
	} else {
		res1 := KthSmallestOfTwoSortedArrays(nums1, nums2, (n1+n2)>>1-1)
		res2 := KthSmallestOfTwoSortedArrays(nums1, nums2, (n1+n2)>>1)
		return float64(res1+res2) / 2
	}
}

// 两个有序数组第k小的数.
// !0 <= k < len(nums1) + len(nums2)
func KthSmallestOfTwoSortedArrays(nums1 []int, nums2 []int, k int) int {
	i1, i2 := 0, 0
	j1, j2 := len(nums1), len(nums2)
	for {
		if i1 == j1 {
			return nums2[i2+k]
		}
		if i2 == j2 {
			return nums1[i1+k]
		}
		m1, m2 := (j1-i1)>>1, (j2-i2)>>1
		if m1+m2 < k {
			if nums1[i1+m1] < nums2[i2+m2] {
				i1 += m1 + 1
				k -= m1 + 1
			} else {
				i2 += m2 + 1
				k -= m2 + 1
			}
		} else {
			if nums1[i1+m1] < nums2[i2+m2] {
				j2 = i2 + m2
			} else {
				j1 = i1 + m1
			}
		}
	}
}
