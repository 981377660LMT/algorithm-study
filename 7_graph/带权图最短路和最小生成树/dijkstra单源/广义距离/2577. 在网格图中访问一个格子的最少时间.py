# 给你一个 m x n 的矩阵 grid ，每个元素都为 非负 整数，
# 其中 grid[row][col] 表示可以访问格子 (row, col) 的 最早 时间。
# 也就是说当你访问格子 (row, col) 时，最少已经经过的时间为 grid[row][col] 。
# 你从 最左上角 出发，出发时刻为 0 ，
# !你必须一直移动到上下左右相邻四个格子中的 任意 一个格子（即不能停留在格子上）。
# 每次移动都需要花费 1 单位时间。
# 请你返回 最早 到达右下角格子的时间，如果你无法到达右下角的格子，请你返回 -1 。
# 2<=m,n<=1000 grid[0][0]==0
# m*n<=1e5


# !边权是动态的dijk

from heapq import heappop, heappush
from typing import List

DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]
INF = int(1e18)


class Solution:
    def minimumTime(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        if grid[0][1] > 1 and grid[1][0] > 1:
            return -1
        pq = [(0, 0, 0)]  # (dist, row, col)
        dist = [[INF] * COL for _ in range(ROW)]
        dist[0][0] = 0
        while pq:
            curDist, curRow, curCol = heappop(pq)
            if curRow == ROW - 1 and curCol == COL - 1:
                return curDist
            if curDist > dist[curRow][curCol]:
                continue
            for dr, dc in DIR4:
                newRow, newCol = curRow + dr, curCol + dc
                if 0 <= newRow < ROW and 0 <= newCol < COL:
                    cand = curDist + 1
                    if grid[newRow][newCol] > cand:
                        diff = grid[newRow][newCol] - cand
                        cand += diff if diff % 2 == 0 else diff + 1
                    if cand < dist[newRow][newCol]:
                        dist[newRow][newCol] = cand
                        heappush(pq, (cand, newRow, newCol))
        return -1
