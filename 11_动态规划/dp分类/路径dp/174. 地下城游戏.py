# 计算确保骑士能够拯救到公主所需的最低初始健康点数
from functools import lru_cache
from typing import List


DIR2 = ((0, 1), (1, 0))


class Solution:
    def calculateMinimumHP(self, dungeon: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(row: int, col: int) -> int:
            if (row, col) == (ROW - 1, COL - 1):
                return max(1, 1 - dungeon[row][col])

            nextMin = int(1e20)
            for dr, dc in DIR2:
                nr, nc = row + dr, col + dc
                if 0 <= nr < ROW and 0 <= nc < COL:
                    nextMin = min(nextMin, dfs(nr, nc))
            return max(1, nextMin - dungeon[row][col])

        ROW, COL = len(dungeon), len(dungeon[0])
        res = dfs(0, 0)
        dfs.cache_clear()
        return res

