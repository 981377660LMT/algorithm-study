from functools import lru_cache
from typing import List

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumCost(self, m: int, n: int, horizontalCut: List[int], verticalCut: List[int]) -> int:
        @lru_cache(None)
        def dfs(row1: int, col1: int, row2: int, col2: int) -> int:
            if row1 == row2 and col1 == col2:
                return 0
            res = INF
            # 横着切
            for i in range(row1, row2):
                cur = dfs(row1, col1, i, col2) + dfs(i + 1, col1, row2, col2) + horizontalCut[i]
                res = min2(res, cur)
            # 竖着切
            for i in range(col1, col2):
                cur = dfs(row1, col1, row2, i) + dfs(row1, i + 1, row2, col2) + verticalCut[i]
                res = min2(res, cur)
            return res

        res = dfs(0, 0, m - 1, n - 1)
        dfs.cache_clear()
        return res
