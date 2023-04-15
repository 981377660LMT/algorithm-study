from bisect import bisect_left
from typing import List
from heapq import merge
from itertools import islice

# 时间复杂度: O(klogn)
# 因为n<=1e4  k最大1e8 所以超时了


class Solution:
    def smallestDistancePair(self, nums: List[int], k: int) -> int:
        n, nums = len(nums), sorted(nums)
        gen = lambda i: (nums[j] - nums[i] for j in range(i + 1, n))
        allGen = [gen(i) for i in range(n)]
        iterable = merge(*allGen)
        return next(islice(iterable, k - 1, None))

    def smallestDistancePair2(self, nums: List[int], k: int) -> int:
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
