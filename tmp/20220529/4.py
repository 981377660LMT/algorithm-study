from heapq import heappop, heappush
from typing import List


MOD = int(1e9 + 7)
INF = int(1e20)

DIR4 = ((0, 1), (0, -1), (1, 0), (-1, 0))


class Solution:
    def minimumObstacles(self, grid: List[List[int]]) -> int:
        #  返回需要移除的障碍物的 最小 数目
        ROW, COL = len(grid), len(grid[0])
        pq = [(0, 0, 0)]
        dist = [[INF] * COL for _ in range(ROW)]
        dist[0][0] = 0
        while pq:
            remove, curRow, curCol = heappop(pq)
            if dist[curRow][curCol] < remove:
                continue
            if (curRow, curCol) == (ROW - 1, COL - 1):
                return remove
            for dRow, dCol in DIR4:
                nRow, nCol = curRow + dRow, curCol + dCol
                if 0 <= nRow < ROW and 0 <= nCol < COL:
                    cost = grid[nRow][nCol]
                    if dist[curRow][curCol] + cost < dist[nRow][nCol]:
                        dist[nRow][nCol] = dist[curRow][curCol] + cost
                        heappush(pq, (dist[nRow][nCol], nRow, nCol))
        return -1


print('1e3'.isnumeric())
