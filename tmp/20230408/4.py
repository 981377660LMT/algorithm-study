from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始的 m x n 整数矩阵 grid 。你一开始的位置在 左上角 格子 (0, 0) 。

# 当你在格子 (i, j) 的时候，你可以移动到以下格子之一：

# 满足 j < k <= grid[i][j] + j 的格子 (i, k) （向右移动），或者
# 满足 i < k <= grid[i][j] + i 的格子 (k, j) （向下移动）。
# 请你返回到达 右下角 格子 (m - 1, n - 1) 需要经过的最少移动格子数，如果无法到达右下角格子，请你返回 -1 。


class Solution:
    def minimumVisitedCells(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
