from bisect import bisect_left
from typing import List


class Solution:
    def smallestDistancePair(self, nums: List[int], k: int) -> int:
        """nlog1e9"""

        def countNGT(mid: int) -> int:
            """距离小于等于mid的个数"""
            res, left = 0, 0
            for right in range(len(nums)):
                while nums[right] - nums[left] > mid:
                    left += 1
                res += right - left
            return res

        nums = sorted(nums)
        return bisect_left(range(int(1e9)), k, key=countNGT)  # countNGT等于时k往左移
