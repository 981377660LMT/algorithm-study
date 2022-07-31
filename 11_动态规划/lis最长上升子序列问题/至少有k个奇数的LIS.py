# k ≤ n ≤ 100
from functools import lru_cache

INF = int(1e20)
# O(n^2*k)

# 至少有k个奇数的LIS
class Solution:
    def solve(self, nums, k) -> int:
        """
        return the length of the longest increasing subsequence with at least k odd numbers.
        LIS:每个数取还是不取
        """

        @lru_cache(None)
        def dfs(index: int, need: int, pre: int) -> int:
            if index == len(nums):
                return 0 if need == 0 else -INF
            if nums[index] > pre:
                nextNeed = need - 1 if nums[index] & 1 else need
                nextNeed = max(nextNeed, 0)
                return max(dfs(index + 1, nextNeed, nums[index]) + 1, dfs(index + 1, need, pre))
            else:
                return dfs(index + 1, need, pre)

        return max(0, dfs(0, k, -INF))
