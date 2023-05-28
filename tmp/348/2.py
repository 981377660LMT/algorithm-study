from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始、大小为 m x n 的二维矩阵 grid ，请你求解大小同样为 m x n 的答案矩阵 answer 。

# 矩阵 answer 中每个单元格 (r, c) 的值可以按下述方式进行计算：

# 令 topLeft[r][c] 为矩阵 grid 中单元格 (r, c) 左上角对角线上 不同值 的数量。
# 令 bottomRight[r][c] 为矩阵 grid 中单元格 (r, c) 右下角对角线上 不同值 的数量。
# 然后 answer[r][c] = |topLeft[r][c] - bottomRight[r][c]| 。

# 返回矩阵 answer 。

# 矩阵对角线 是从最顶行或最左列的某个单元格开始，向右下方向走到矩阵末尾的对角线。

# 如果单元格 (r1, c1) 和单元格 (r, c) 属于同一条对角线且 r1 < r ，则单元格 (r1, c1) 属于单元格 (r, c) 的左上对角线。类似地，可以定义右下对角线。


class Solution:
    def differenceOfDistinctValues(self, grid: List[List[int]]) -> List[List[int]]:
        ROW, COL = len(grid), len(grid[0])
        topleft = [[0] * COL for _ in range(ROW)]
        bottomright = [[0] * COL for _ in range(ROW)]

        for r in range(ROW):
            for c in range(COL):
                curR, curC = r - 1, c - 1
                s = set()
                while curR >= 0 and curC >= 0:
                    s.add(grid[curR][curC])
                    curR -= 1
                    curC -= 1
                topleft[r][c] = len(s)
                curR, curC = r + 1, c + 1
                s = set()
                while curR < ROW and curC < COL:
                    s.add(grid[curR][curC])
                    curR += 1
                    curC += 1
                bottomright[r][c] = len(s)
        return [[abs(topleft[r][c] - bottomright[r][c]) for c in range(COL)] for r in range(ROW)]
