# n ≤ 30 where n is the length of nums
# 0 ≤ nums[i] ≤ 100

# 是否可以选出两个子集和相等，使他们的和最大
# 注意到nums[i]<=100 表示值域范围最多6000

# jump group1 group2

from functools import lru_cache


class Solution:
    def solve(self, nums):
        """时间复杂度O(n*sum(nums))"""

        @lru_cache(None)
        def dfs(index: int, diff: int) -> int:
            if index == n:
                return 0 if diff == 0 else -int(1e20)
            res = -int(1e20)
            res = max(res, dfs(index + 1, diff))
            res = max(res, dfs(index + 1, diff - nums[index]))
            res = max(res, dfs(index + 1, diff + nums[index]) + nums[index])
            return res

        n = len(nums)
        return dfs(0, 0)


print(Solution().solve(nums=[1, 4, 3, 5]))

