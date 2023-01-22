from typing import List

# 给你两个整数数组 nums1 和 nums2 ，它们已经按非降序排序，
# 请你返回两个数组的 最小公共整数 。
# 如果两个数组 nums1 和 nums2 没有公共整数，请你返回 -1 。


class Solution:
    def getCommon(self, nums1: List[int], nums2: List[int]) -> int:
        """数组已经排序"""
        i, j = 0, 0
        while i < len(nums1) and j < len(nums2):
            if nums1[i] == nums2[j]:
                return nums1[i]
            elif nums1[i] < nums2[j]:
                i += 1
            else:
                j += 1
        return -1

    def getCommon2(self, nums1: List[int], nums2: List[int]) -> int:
        """数组未排序"""
        s1, s2 = set(nums1), set(nums2)
        res = s1 & s2
        if res:
            return min(res)
        return -1
