from functools import lru_cache
from typing import List

# 有 m 个房子排成一排，你需要给每个房子涂上 n 种颜色之一
# 请你返回房子涂色方案的最小总花费，使得每个房子都被涂色后，恰好组成 target 个街区。如果没有可用的涂色方案，请返回 -1 。
# 1 <= m <= 100
# 1 <= n <= 20
# 1 <= target <= m
class Solution:
    def minCost(self, houses: List[int], cost: List[List[int]], m: int, n: int, target: int) -> int:
        @lru_cache(None)
        def dfs(index: int, remain: int, pre: int) -> int:
            if remain < 0:
                return int(1e20)
            if index == m:
                return 0 if remain == 0 else int(1e20)

            res = int(1e20)
            if houses[index] == 0:
                for i in range(n):
                    color = i + 1
                    res = min(
                        res, dfs(index + 1, remain - int(color != pre), color) + cost[index][i]
                    )
            else:
                res = dfs(index + 1, remain - int(houses[index] != pre), houses[index])
            return res

        res = dfs(0, target, -1)
        dfs.cache_clear()
        return res if res < int(1e19) else -1

