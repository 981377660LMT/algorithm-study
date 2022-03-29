# k ≤ n ≤ 100
from functools import lru_cache

# O(n^2*k)


class Solution:
    def solve(self, nums, k) -> int:
        """
        return the length of the longest increasing subsequence with at least k odd numbers.
        LIS:每个数取还是不取
        """

        @lru_cache(None)
        def dfs(index: int, need: int, pre: int) -> int:
            if index == len(nums):
                return 0 if need == 0 else -int(1e20)
            if nums[index] > pre:
                nextNeed = need - 1 if nums[index] & 1 else need
                nextNeed = max(nextNeed, 0)
                return max(dfs(index + 1, nextNeed, nums[index]) + 1, dfs(index + 1, need, pre))
            else:
                return dfs(index + 1, need, pre)

        return max(0, dfs(0, k, -int(1e20)))

