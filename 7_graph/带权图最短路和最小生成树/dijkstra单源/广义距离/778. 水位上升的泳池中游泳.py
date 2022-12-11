# 当开始下雨时，在时间为 t 时，水池中的水位为 t 。
# 你可以从一个平台游向四周相邻的任意一个平台，
# 但是前提是此时水位必须同时淹没这两个平台。
# 假定你可以瞬间移动无限距离，也就是默认在方格内部游动是不耗时的。
# 当然，在你游泳的时候你必须待在坐标方格里面。


# 你从坐标方格的左上平台 (0，0) 出发。
# 返回 你到达坐标方格的右下平台 (n-1, n-1) 所需的最少时间 。
# !边权:两个点的最大值
# !还可以二分/并查集 使用并查集的话可以获取更多信息

from heapq import heappop, heappush
from typing import List


INF = int(1e18)
DIR4 = ((0, 1), (0, -1), (1, 0), (-1, 0))


class Solution:
    def swimInWater(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        dist = [[INF] * COL for _ in range(ROW)]
        dist[0][0] = grid[0][0]
        pq = [(grid[0][0], 0, 0)]
        while pq:
            curDist, curRow, curCol = heappop(pq)
            if curDist > dist[curRow][curCol]:
                continue
            if (curRow, curCol) == (ROW - 1, COL - 1):
                return curDist
            for dr, dc in DIR4:
                nextRow, nextCol = curRow + dr, curCol + dc
                if 0 <= nextRow < ROW and 0 <= nextCol < COL:
                    cand = max(curDist, grid[nextRow][nextCol])  # !边权:两个点的最大值
                    if cand < dist[nextRow][nextCol]:
                        dist[nextRow][nextCol] = cand
                        heappush(pq, (cand, nextRow, nextCol))
        raise Exception("No path found")


assert Solution().swimInWater([[0, 2], [1, 3]]) == 3
