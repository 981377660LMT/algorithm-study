from functools import lru_cache
from typing import List, Tuple
from collections import defaultdict, Counter


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minRemainingSpace(self, N: List[int], V: int) -> int:
        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain < 0:
                return int(1e20)
            if index == n:
                return remain

            res = dfs(index + 1, remain)
            if remain - N[index] >= 0:
                res = min(res, dfs(index + 1, remain - N[index]))
            return res

        n = len(N)
        res = dfs(0, V)
        dfs.cache_clear()
        return res

