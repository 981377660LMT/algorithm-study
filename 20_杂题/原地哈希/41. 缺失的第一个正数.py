from typing import List

INF = int(1e18)


class Solution:
    def firstMissingPositive(self, nums: List[int]) -> int:
        n = len(nums)
        for i, v in enumerate(nums):
            if not (1 <= v <= n):
                nums[i] = INF
        for v in nums:
            pos = abs(v) - 1
            if 0 <= pos < n:
                nums[pos] = -abs(nums[pos])
        for i, v in enumerate(nums):
            if v > 0:
                return i + 1
        return n + 1
