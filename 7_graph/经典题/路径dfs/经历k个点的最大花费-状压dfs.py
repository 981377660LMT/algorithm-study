from collections import defaultdict
from functools import lru_cache
from typing import List

# 2 <= n <= 15
# 1 <= k <= 50

# 经过k个点的最大花费


class Solution:
    def maximumCost(self, n: int, highways: List[List[int]], k: int) -> int:
        @lru_cache(None)
        def dfs(cur: int, visited: int, remain: int) -> int:
            if remain == 0:
                return 0

            res = -int(1e20)
            for next, weight in adjMap[cur].items():
                if not (visited >> next) & 1:
                    res = max(res, dfs(next, visited | (1 << next), remain - 1) + weight)
            return res

        adjMap = defaultdict(lambda: defaultdict(lambda: -int(1e20)))
        for u, v, w in highways:
            adjMap[u][v] = w
            adjMap[v][u] = w

        res = -1
        for i in range(n):
            res = max(res, dfs(i, 1 << i, k))
        dfs.cache_clear()
        return res

