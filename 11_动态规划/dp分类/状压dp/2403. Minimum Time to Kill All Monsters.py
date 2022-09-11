# 1 <= power.length <= 17
# 1 <= power[i] <= 1e9

from functools import lru_cache
from math import ceil
from typing import List

INF = int(1e18)


class Solution:
    def minimumTime(self, power: List[int]) -> int:
        @lru_cache(None)
        def dfs(visited: int, gain: int) -> int:
            if visited == target:
                return 0

            res = INF
            for i in range(n):
                if visited & (1 << i):
                    continue
                cost = ceil(power[i] / gain)
                cand = cost + dfs(visited | (1 << i), gain + 1)
                res = cand if cand < res else res

            return res

        n = len(power)
        target = (1 << n) - 1
        res = dfs(0, 1)
        dfs.cache_clear()
        return res

    def minimumTime2(self, power: List[int]) -> int:
        n = len(power)
        target = (1 << n) - 1
        dp = [INF] * (target + 1)
        dp[0] = 0

        for state in range(target + 1):
            gain = state.bit_count() + 1
            for i in range(n):
                if state & (1 << i):
                    continue
                cost = ceil(power[i] / gain)
                dp[state | (1 << i)] = min(dp[state | (1 << i)], cost + dp[state])

        return dp[target]


print(Solution().minimumTime2(power=[3, 1, 4]))
print(Solution().minimumTime2(power=[1, 1, 4]))
print(Solution().minimumTime2(power=[1, 2, 4, 9]))
