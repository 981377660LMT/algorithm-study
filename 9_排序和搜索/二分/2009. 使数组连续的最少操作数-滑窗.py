from typing import List


class Solution:
    def minOperations(self, nums: List[int]) -> int:
        n = len(nums)
        nums = sorted(set(nums))

        left = 0
        res = n
        for right, num in enumerate(nums):
            while num - nums[left] >= n:
                left += 1
            cand = right - left + 1
            res = min(res, n - cand)

        return res

