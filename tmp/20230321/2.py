from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 骑士在一张 n x n 的棋盘上巡视。在有效的巡视方案中，骑士会从棋盘的 左上角 出发，并且访问棋盘上的每个格子 恰好一次 。

# 给你一个 n x n 的整数矩阵 grid ，由范围 [0, n * n - 1] 内的不同整数组成，其中 grid[row][col] 表示单元格 (row, col) 是骑士访问的第 grid[row][col] 个单元格。骑士的行动是从下标 0 开始的。

# 如果 grid 表示了骑士的有效巡视方案，返回 true；否则返回 false。

DIR8 = set([(1, 2), (2, 1), (-1, 2), (2, -1), (1, -2), (-2, 1), (-1, -2), (-2, -1)])


class Solution:
    def checkValidGrid(self, grid: List[List[int]]) -> bool:
        if grid[0][0] != 0:
            return False
        pos = {}
        for i in range(len(grid)):
            for j in range(len(grid[0])):
                pos[grid[i][j]] = (i, j)
        for i in range(len(grid) * len(grid[0]) - 1):
            x, y = pos[i]
            x1, y1 = pos[i + 1]
            if (x1 - x, y1 - y) not in DIR8:
                return False
        return True


# [[24,11,22,17,4],[21,16,5,12,9],[6,23,10,3,18],[15,20,1,8,13],[0,7,14,19,2]]
print(
    Solution().checkValidGrid(
        [
            [24, 11, 22, 17, 4],
            [21, 16, 5, 12, 9],
            [6, 23, 10, 3, 18],
            [15, 20, 1, 8, 13],
            [0, 7, 14, 19, 2],
        ]
    )
)
