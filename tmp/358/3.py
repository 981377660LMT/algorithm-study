from heapq import heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


# 给你一个下标从 0 开始、大小为 n x n 的二维矩阵 grid ，其中 (r, c) 表示：

# 如果 grid[r][c] = 1 ，则表示一个存在小偷的单元格
# 如果 grid[r][c] = 0 ，则表示一个空单元格
# 你最开始位于单元格 (0, 0) 。在一步移动中，你可以移动到矩阵中的任一相邻单元格，包括存在小偷的单元格。

# 矩阵中路径的 安全系数 定义为：从路径中任一单元格到矩阵中任一小偷所在单元格的 最小 曼哈顿距离。

# 返回所有通向单元格 (n - 1, n - 1) 的路径中的 最大安全系数 。

# 单元格 (r, c) 的某个 相邻 单元格，是指在矩阵中存在的 (r, c + 1)、(r, c - 1)、(r + 1, c) 和 (r - 1, c) 之一。


# 两个单元格 (a, b) 和 (x, y) 之间的 曼哈顿距离 等于 | a - x | + | b - y | ，其中 |val| 表示 val 的绝对值。


# 预处理出每个格子的安全系数
DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


# !网格图多源bfs
class Solution:
    def maximumSafenessFactor(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        values = [[0] * COL for _ in range(ROW)]
        queue = deque()
        visited = [[False] * COL for _ in range(ROW)]
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == 1:
                    queue.append((r, c))
                    visited[r][c] = True

        while queue:
            len_ = len(queue)
            for _ in range(len_):
                curRow, curCol = queue.popleft()
                for dr, dc in DIR4:
                    newRow, newCol = curRow + dr, curCol + dc
                    if 0 <= newRow < ROW and 0 <= newCol < COL and not visited[newRow][newCol]:
                        values[newRow][newCol] = values[curRow][curCol] + 1
                        visited[newRow][newCol] = True
                        queue.append((newRow, newCol))

        pq = [(-values[0][0], 0, 0)]  # (safety, row, col)
        dist = [[0] * COL for _ in range(ROW)]
        dist[0][0] = values[0][0]
        while pq:
            safety, curRow, curCol = heappop(pq)
            safety = -safety
            if curRow == ROW - 1 and curCol == COL - 1:
                return safety
            for dr, dc in DIR4:
                newRow, newCol = curRow + dr, curCol + dc
                if 0 <= newRow < ROW and 0 <= newCol < COL:
                    newSafety = min(safety, values[newRow][newCol])
                    if newSafety > dist[newRow][newCol]:
                        dist[newRow][newCol] = newSafety
                        heappush(pq, (-newSafety, newRow, newCol))
        return 0


# grid = [[1,0,0],[0,0,0],[0,0,1]]
print(Solution().maximumSafenessFactor([[1, 0, 0], [0, 0, 0], [0, 0, 1]]))
#  [[0,0,1],[0,0,0],[0,0,0]]
print(Solution().maximumSafenessFactor([[0, 0, 1], [0, 0, 0], [0, 0, 0]]))
