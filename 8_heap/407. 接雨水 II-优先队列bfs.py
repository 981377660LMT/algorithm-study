# 407. 接雨水 II-优先队列bfs
# 给你一个 m x n 的矩阵，其中的值均为非负整数，代表二维高度图每个单元的高度
# 请计算图中形状最多能接多少体积的雨水。
# ROW<=200 COL<=200
# !每次从最低的点开始搜索

from heapq import heappop, heappush
from typing import List

DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


class Solution:
    def trapRainWater(self, heightMap: List[List[int]]) -> int:
        ROW, COL = len(heightMap), len(heightMap[0])
        if ROW < 3 or COL < 3:
            return 0

        visited = [[False] * COL for _ in range(ROW)]
        pq = []  # 边界开始搜索
        for r in range(ROW):
            for c in range(COL):
                if r == 0 or r == ROW - 1 or c == 0 or c == COL - 1:
                    visited[r][c] = True
                    heappush(pq, (heightMap[r][c], r, c))

        res = 0
        while pq:
            curHeight, curRow, curCol = heappop(pq)
            for dr, dc in DIR4:
                nr, nc = curRow + dr, curCol + dc
                if 0 <= nr < ROW and 0 <= nc < COL and not visited[nr][nc]:
                    visited[nr][nc] = True
                    res += max(0, curHeight - heightMap[nr][nc])  # !这个点能接的雨水
                    heappush(pq, (max(curHeight, heightMap[nr][nc]), nr, nc))

        return res


assert (
    Solution().trapRainWater(heightMap=[[1, 4, 3, 1, 3, 2], [3, 2, 1, 3, 2, 4], [2, 3, 3, 2, 3, 1]])
    == 4
)
