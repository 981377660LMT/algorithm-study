from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个大小为 m x n 的二维矩阵 grid 。你需要判断每一个格子 grid[i][j] 是否满足：


# 如果它下面的格子存在，那么它需要等于它下面的格子，也就是 grid[i][j] == grid[i + 1][j] 。
# 如果它右边的格子存在，那么它需要不等于它右边的格子，也就是 grid[i][j] != grid[i][j + 1] 。
# 如果 所有 格子都满足以上条件，那么返回 true ，否则返回 false 。
class Solution:
    def satisfiesConditions(self, grid: List[List[int]]) -> bool:
        ROW, COL = len(grid), len(grid[0])
        for i in range(ROW):
            for j in range(COL):
                if i + 1 < ROW and grid[i][j] != grid[i + 1][j]:
                    return False
                if j + 1 < COL and grid[i][j] == grid[i][j + 1]:
                    return False
        return True
