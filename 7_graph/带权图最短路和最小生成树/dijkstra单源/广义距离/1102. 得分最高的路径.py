from typing import List
from heapq import heappush, heappop

# !路径的得分是`该路径上的 最小 值`。例如，路径 8 →  4 →  5 →  9 的值为 4
# !找出所有路径中得分 `最高` 的那条路径，返回其 得分。
# 边权:最小值
# 还可以排序+并查集获取更多信息

INF = int(1e18)
DIR4 = ((-1, 0), (1, 0), (0, -1), (0, 1))


class Solution:
    def maximumMinimumPath(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        dist = [[-INF] * COL for _ in range(ROW)]
        dist[0][0] = grid[0][0]
        pq = [(-grid[0][0], 0, 0)]  # (cost,row,col) 大根堆

        while pq:
            curCost, curRow, curCol = heappop(pq)
            curCost = -curCost
            if curCost < dist[curRow][curCol]:
                continue
            if (curRow, curCol) == (ROW - 1, COL - 1):
                return curCost
            for dr, dc in DIR4:
                nr, nc = curRow + dr, curCol + dc
                if 0 <= nr < ROW and 0 <= nc < COL:
                    cand = min(curCost, grid[nr][nc])  # !边权:最小值
                    if cand > dist[nr][nc]:
                        dist[nr][nc] = cand
                        heappush(pq, (-cand, nr, nc))

        raise Exception("No path found")


print(Solution().maximumMinimumPath([[5, 4, 5], [1, 2, 6], [7, 4, 6]]))
