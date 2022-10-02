from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 沙漏的最大总和


class Solution:
    def maxSum(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        res = 0
        for r in range(ROW - 2):
            for c in range(COL - 2):
                curSum = (
                    grid[r][c]
                    + grid[r][c + 1]
                    + grid[r][c + 2]
                    + grid[r + 1][c + 1]
                    + grid[r + 2][c]
                    + grid[r + 2][c + 1]
                    + grid[r + 2][c + 2]
                )
                res = max(res, curSum)
        return res
