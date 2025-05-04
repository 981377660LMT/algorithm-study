# https://leetcode.cn/problems/median-of-two-sorted-arrays/description/
# 给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的 中位数 。
# 算法的时间复杂度应该为 O(log (m+n)) 。

from typing import List


INF = int(1e18)


class Solution:
    def findMedianSortedArrays(self, nums1: List[int], nums2: List[int]) -> float:
        if len(nums1) > len(nums2):
            nums1, nums2 = nums2, nums1

        len1, len2 = len(nums1), len(nums2)
        left, right = 0, len1
        max1, min2 = 0, 0  # 前一部分的最大值和后一部分的最小值

        while left <= right:
            # 前一部分包含 nums1[0 .. i-1] 和 nums2[0 .. j-1]
            # 后一部分包含 nums1[i .. m-1] 和 nums2[j .. n-1]
            i = (left + right) // 2
            j = (len1 + len2 + 1) // 2 - i

            # a1, b1, a2, b2 分别表示 nums1[i-1], nums1[i], nums2[j-1], nums2[j]
            a1 = -INF if i == 0 else nums1[i - 1]
            b1 = INF if i == len1 else nums1[i]
            a2 = -INF if j == 0 else nums2[j - 1]
            b2 = INF if j == len2 else nums2[j]

            if a1 <= b2:
                max1 = a1 if a1 > a2 else a2
                min2 = b1 if b1 < b2 else b2
                left = i + 1
            else:
                right = i - 1

        return max1 if (len1 + len2) & 1 else (max1 + min2) / 2
