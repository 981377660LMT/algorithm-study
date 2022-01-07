from typing import List
from bisect import bisect_right


class Solution:
    def minOperations(self, nums: List[int]) -> int:
        n = len(nums)
        nums = sorted(set(nums))

        res = n
        for i, start in enumerate(nums):
            end = start + n - 1
            upper = bisect_right(nums, end) - 1
            cand = upper - i + 1
            res = min(res, n - cand)

        return res

