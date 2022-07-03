from functools import cache
from itertools import product
from typing import List

MOD = int(1e9 + 7)
DIR4 = [[0, 1], [1, 0], [0, -1], [-1, 0]]


class Solution:
    def countPaths(self, grid: List[List[int]]) -> int:
        @cache
        def dfs(row: int, col: int) -> int:
            res = 1
            for dr, dc in DIR4:
                nr, nc = row + dr, col + dc
                if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] > grid[row][col]:
                    res += dfs(nr, nc)
                    res %= MOD
            return res

        res = 0
        ROW, COL = len(grid), len(grid[0])
        for r, c in product(range(ROW), range(COL)):
            res += dfs(r, c)
            res %= MOD
        return res
