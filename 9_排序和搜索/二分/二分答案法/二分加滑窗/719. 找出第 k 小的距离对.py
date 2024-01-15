# https://leetcode.cn/problems/find-k-th-smallest-pair-distance/
# 719. 找出第 K 小的距离对(一维)
# 两个数组第k大的和

from typing import List


class Solution:
    def smallestDistancePair(self, nums: List[int], k: int) -> int:
        """nlog1e9"""

        def countNGT(mid: int) -> int:
            """有多少个不超过mid的候选"""
            res, left, n = 0, 0, len(nums)
            for right in range(n):
                while left <= right and nums[right] - nums[left] > mid:
                    left += 1
                res += right - left
            return res

        nums = sorted(nums)
        left, right = 0, int(1e18)
        while left <= right:
            mid = (left + right) // 2
            if countNGT(mid) < k:
                left = mid + 1
            else:
                right = mid - 1
        return left
