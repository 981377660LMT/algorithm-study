from functools import lru_cache
from typing import List

MOD = int(1e9 + 7)

# time complexity: o(kn^2)


class Solution:
    def solve(self, nums: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(index: int, pre: int, count: int) -> int:
            if count > k:
                return 0
            if index == n:
                return int(count == k)

            res = dfs(index + 1, pre, count) % MOD
            if nums[index] > pre:
                res += dfs(index + 1, nums[index], count + 1)
                res %= MOD

            return res % MOD

        n = len(nums)
        res = dfs(0, -int(1e20), 0)
        dfs.cache_clear()
        return res % MOD


print(Solution().solve(nums=[1, 2, 3, 0], k=2))
