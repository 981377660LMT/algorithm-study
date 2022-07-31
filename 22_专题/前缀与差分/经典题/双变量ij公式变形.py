from typing import List

# n ≤ 100,000

INF = int(1e20)


class Solution:
    def solve(self, nums: List[int]) -> int:
        """找到nums[i] + nums[j] + (i - j)的最大值(i<j)"""
        preMax = -INF
        res = -INF
        for i, num in enumerate(nums):
            res = max(res, num - i + preMax)
            preMax = max(preMax, num + i)
        return res


print(Solution().solve(nums=[5, 5, 1, 1, 1, 7]))
