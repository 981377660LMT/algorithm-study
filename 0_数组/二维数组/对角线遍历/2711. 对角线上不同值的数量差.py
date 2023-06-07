# !按照对角线进行遍历 O(n*m)
# 给你一个下标从 0 开始、大小为 m x n 的二维矩阵 grid ，请你求解大小同样为 m x n 的答案矩阵 answer 。
# 矩阵 answer 中每个单元格 (r, c) 的值可以按下述方式进行计算：
# 令 topLeft[r][c] 为矩阵 grid 中单元格 (r, c) 左上角对角线上 不同值 的数量。
# 令 bottomRight[r][c] 为矩阵 grid 中单元格 (r, c) 右下角对角线上 不同值 的数量。
# 然后 answer[r][c] = |topLeft[r][c] - bottomRight[r][c]| 。


from typing import List
from enumerateDiagnal import enumerateDiagnal


class Solution:
    def differenceOfDistinctValues(self, grid: List[List[int]]) -> List[List[int]]:
        ROW, COL = len(grid), len(grid[0])
        topLeft = [[0] * COL for _ in range(ROW)]
        bottomRight = [[0] * COL for _ in range(ROW)]

        for group in enumerateDiagnal(grid, direction=0):
            visited = set()
            for r, c in group:
                topLeft[r][c] = len(visited)
                visited.add(grid[r][c])
        for group in enumerateDiagnal(grid, direction=1):
            visited = set()
            for r, c in group:
                bottomRight[r][c] = len(visited)
                visited.add(grid[r][c])

        return [[abs(topLeft[r][c] - bottomRight[r][c]) for c in range(COL)] for r in range(ROW)]
