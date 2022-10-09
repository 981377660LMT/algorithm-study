from functools import cache, lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# m == grid.length
# n == grid[i].length
# 1 <= m, n <= 5 * 104
# 1 <= m * n <= 5 * 104
# 0 <= grid[i][j] <= 100
# 1 <= k <= 50
# 你从起点 (0, 0) 出发，每一步只能往 下 或者往 右 ，你想要到达终点 (m - 1, n - 1) 。

# 请你返回路径和能被 k 整除的路径数目，由于答案可能很大，返回答案对 109 + 7 取余 的结果。


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
