from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxScore(self, grid: List[List[int]]) -> int:
        mp = defaultdict(set)
        for r, row in enumerate(grid):
            for _, val in enumerate(row):
                mp[val].add(r)

        ROW = len(grid)
        MASK = (1 << ROW) - 1

        @lru_cache(None)
        def dfs(index: int, visited: int) -> int:
            if visited == MASK:
                return 0
            if index == m:
                return 0
            res = dfs(index + 1, visited)
            for row in mp[values[index]]:
                if visited & (1 << row):
                    continue
                res = max2(res, values[index] + dfs(index + 1, visited | (1 << row)))
            return res

        values = sorted(mp.keys(), reverse=True)
        m = len(values)
        res = dfs(0, 0)
        dfs.cache_clear()
        return res
