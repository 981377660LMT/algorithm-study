from typing import List, Tuple
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def checkXMatrix(self, grid: List[List[int]]) -> bool:
        #  矩阵对角线上的所有元素都 不是 0
        # 矩阵中所有其他元素都是 0
        ROW, COL = len(grid), len(grid[0])
        for i in range(ROW):
            for j in range(COL):
                # 对角线
                if i == j or i + j == ROW - 1:
                    if grid[i][j] == 0:
                        return False
                # 其他元素
                else:
                    if grid[i][j] != 0:
                        return False
        return True

