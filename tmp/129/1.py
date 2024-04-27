from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二维 3 x 3 的矩阵 grid ，每个格子都是一个字符，要么是 'B' ，要么是 'W' 。字符 'W' 表示白色，字符 'B' 表示黑色。


# 你的任务是改变 至多一个 格子的颜色，使得矩阵中存在一个 2 x 2 颜色完全相同的正方形。
# 如果可以得到一个相同颜色的 2 x 2 正方形，那么返回 true ，否则返回 false 。
class Solution:
    def canMakeSquare(self, grid: List[List[str]]) -> bool:
        def check(mat: List[List[[str]]]) -> bool:
            for i in range(2):
                for j in range(2):
                    if mat[i][j] == mat[i + 1][j] == mat[i][j + 1] == mat[i + 1][j + 1]:
                        return True
            return False

        if check(grid):
            return True
        for i in range(3):
            for j in range(3):
                if grid[i][j] == "B":
                    grid[i][j] = "W"
                    if check(grid):
                        return True
                    grid[i][j] = "B"
                else:
                    grid[i][j] = "B"
                    if check(grid):
                        return True
                    grid[i][j] = "W"
        return False
