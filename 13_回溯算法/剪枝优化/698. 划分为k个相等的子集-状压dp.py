from functools import lru_cache
from typing import List


class Solution:
    def canPartitionKSubsets(self, nums: List[int], k: int) -> bool:
        @lru_cache(None)
        def dfs(visited: int, curSum: int) -> bool:
            if visited == (1 << n) - 1:
                return True
            for i in range(n):
                if (visited >> i) & 1:
                    continue
                if curSum + nums[i] <= div:
                    if dfs(visited | (1 << i), (curSum + nums[i]) % div):
                        return True
            return False

        div, mod = divmod(sum(nums), k)
        if mod != 0:
            return False

        n = len(nums)
        return dfs(0, 0)

