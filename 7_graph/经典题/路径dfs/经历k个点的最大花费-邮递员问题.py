from collections import defaultdict
from functools import lru_cache
from typing import List

# 2 <= n <= 15
# 1 <= k <= 50

# 经过k个点的最大花费

# Traveling Salesman 问题


class Solution:
    def maximumCost(self, n: int, highways: List[List[int]], k: int) -> int:
        @lru_cache(None)
        def dfs(cur: int, visited: int) -> int:
            if bin(visited).count('1') == k + 1:  # visited.bit_count()== k+1
                return 0

            res = -int(1e20)
            for next, weight in adjMap[cur].items():
                if not (visited >> next) & 1:
                    res = max(res, dfs(next, visited | (1 << next)) + weight)
            return res

        adjMap = defaultdict(lambda: defaultdict(lambda: -int(1e20)))
        for u, v, w in highways:
            adjMap[u][v] = w
            adjMap[v][u] = w

        res = -1
        for i in range(n):
            res = max(res, dfs(i, 1 << i))
        print(dfs.cache_info())
        dfs.cache_clear()
        return res


print(
    Solution().maximumCost(
        n=5, highways=[[0, 1, 4], [2, 1, 3], [1, 4, 11], [3, 2, 3], [3, 4, 2]], k=3
    )
)
