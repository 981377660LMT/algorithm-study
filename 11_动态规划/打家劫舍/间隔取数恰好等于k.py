from functools import lru_cache

# n * k ≤ 100,000


class Solution:
    def solve(self, nums, k):
        @lru_cache(None)
        def dfs(curSum, index):
            if curSum > k:
                return False
            if curSum == k:
                return True
            if index >= len(nums):
                return False

            return dfs(curSum + nums[index], index + 2) or dfs(curSum, index + 1)

        return dfs(0, 0)
