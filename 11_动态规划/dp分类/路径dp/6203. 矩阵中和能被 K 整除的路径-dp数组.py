from functools import lru_cache
from typing import List

MOD = int(1e9 + 7)

# !1. cache不用cache_clear() 超时 内存占用过大拖慢了gc的速度
# !也可以用gc.collect


class Solution:
    def numberOfPaths(self, grid: List[List[int]], k: int) -> int:
        @lru_cache(None)
        def dfs(row: int, col: int, mod: int) -> int:
            if row == ROW - 1 and col == COL - 1:
                return 1 if mod == 0 else 0
            res = 0
            if row + 1 < ROW:
                res += dfs(row + 1, col, (mod + grid[row + 1][col]) % k)
                res %= MOD
            if col + 1 < COL:
                res += dfs(row, col + 1, (mod + grid[row][col + 1]) % k)
                res %= MOD
            return res

        ROW, COL = len(grid), len(grid[0])
        res = dfs(0, 0, grid[0][0] % k)
        dfs.cache_clear()
        return res

    def numberOfPaths2(self, grid: List[List[int]], k: int) -> int:
        """
        dp[r][c][mod] = dp[r-1][c][mod-grid[r][c]] + dp[r][c-1][mod-grid[r][c]]
        """

        ROW, COL = len(grid), len(grid[0])
        dp = [[[0] * k for _ in range(COL)] for _ in range(ROW)]
        dp[0][0][grid[0][0] % k] = 1
        for r in range(ROW):
            for c in range(COL):
                for mod in range(k):
                    if r - 1 >= 0:
                        dp[r][c][mod] += dp[r - 1][c][(mod - grid[r][c]) % k]
                        dp[r][c][mod] %= MOD
                    if c - 1 >= 0:
                        dp[r][c][mod] += dp[r][c - 1][(mod - grid[r][c]) % k]
                        dp[r][c][mod] %= MOD
        return dp[ROW - 1][COL - 1][0]


print(Solution().numberOfPaths(grid=[[5, 2, 4], [3, 0, 5], [0, 7, 2]], k=3))
print(Solution().numberOfPaths2(grid=[[5, 2, 4], [3, 0, 5], [0, 7, 2]], k=3))
