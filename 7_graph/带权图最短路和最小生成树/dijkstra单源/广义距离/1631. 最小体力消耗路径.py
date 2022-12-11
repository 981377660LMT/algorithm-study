from heapq import heappop, heappush
from typing import List

# 一条路径耗费的 体力值 是路径上相邻格子之间 高度差绝对值 的 最大值 决定的。
# 请你返回从左上角走到右下角的最小 体力消耗值 。
# 边权:绝对值的差

INF = int(1e18)
DIR4 = ((-1, 0), (1, 0), (0, -1), (0, 1))


class Solution:
    def minimumEffortPath(self, heights: List[List[int]]) -> int:
        ROW, COL = len(heights), len(heights[0])
        dist = [[INF] * COL for _ in range(ROW)]
        dist[0][0] = 0
        pq = [(0, 0, 0)]  # (cost,row,col) 小根堆
        while pq:
            curCost, curRow, curCol = heappop(pq)
            if curCost > dist[curRow][curCol]:
                continue
            if (curRow, curCol) == (ROW - 1, COL - 1):
                return curCost
            for dr, dc in DIR4:
                nr, nc = curRow + dr, curCol + dc
                if 0 <= nr < ROW and 0 <= nc < COL:
                    cand = max(curCost, abs(heights[nr][nc] - heights[curRow][curCol]))  # !边权
                    if cand < dist[nr][nc]:
                        dist[nr][nc] = cand
                        heappush(pq, (cand, nr, nc))
        raise Exception("No path found")


assert Solution().minimumEffortPath([[1, 2, 2], [3, 8, 2], [5, 3, 5]]) == 2
