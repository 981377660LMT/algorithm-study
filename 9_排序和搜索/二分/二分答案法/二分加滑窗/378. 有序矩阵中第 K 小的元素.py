# https://leetcode.cn/problems/kth-smallest-element-in-a-sorted-matrix/
# 378. 有序矩阵中第 K 小的元素(1<=k<=n^2)
# -1e9<=nums[i]<=1e9

from bisect import bisect_right
from typing import List


class Solution:
    def kthSmallest(self, matrix: List[List[int]], k: int) -> int:
        def countNGT(mid: int) -> int:
            """有多少个不超过mid的候选"""
            return sum(bisect_right(row, mid) for row in matrix)

        left, right = -int(1e9), int(1e9)
        while left <= right:
            mid = (left + right) // 2
            if countNGT(mid) < k:
                left = mid + 1
            else:
                right = mid - 1
        return left
