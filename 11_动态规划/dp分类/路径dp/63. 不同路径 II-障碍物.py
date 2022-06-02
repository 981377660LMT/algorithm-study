from functools import lru_cache
from typing import List


DIR2 = ((1, 0), (0, 1))


class Solution:
    def uniquePathsWithObstacles(self, obstacleGrid: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(row: int, col: int) -> int:
            if obstacleGrid[row][col] == 1:
                return 0
            if (row, col) == (ROW - 1, COL - 1):
                return 1

            res = 0
            for dr, dc in DIR2:
                nr, nc = row + dr, col + dc
                if 0 <= nr < ROW and 0 <= nc < COL and obstacleGrid[nr][nc] == 0:
                    res += dfs(nr, nc)
            return res

        ROW, COL = len(obstacleGrid), len(obstacleGrid[0])
        res = dfs(0, 0)
        dfs.cache_clear()
        return res

