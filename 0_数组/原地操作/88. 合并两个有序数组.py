# 88. 合并两个有序数组
# https://leetcode.cn/problems/merge-sorted-array/


from typing import List


class Solution:
    def merge(self, nums1: List[int], m: int, nums2: List[int], n: int) -> None:
        """
        Do not return anything, modify nums1 in-place instead.
        使用三个指针从尾部向前合并，避免覆盖 nums1 中还未处理的元素。
        """
        # p1 指向 nums1 中有效区的末尾，p2 指向 nums2 末尾，p 指向合并后数组的末尾
        p1, p2, p = m - 1, n - 1, m + n - 1

        # 只要 nums2 还有没放入的，就继续比较、写入
        while p2 >= 0:
            if p1 >= 0 and nums1[p1] > nums2[p2]:
                nums1[p] = nums1[p1]
                p1 -= 1
            else:
                nums1[p] = nums2[p2]
                p2 -= 1
            p -= 1
