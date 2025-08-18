# 3651. 带传送的最小路径成本_分层图最短路
# https://leetcode.cn/problems/minimum-cost-path-with-teleportations/description/
#
# 给你一个 m x n 的二维整数数组 grid 和一个整数 k。
# 你从左上角的单元格 (0, 0) 出发，目标是到达右下角的单元格 (m - 1, n - 1)。
# 有两种移动方式可用：
# 普通移动：你可以从当前单元格 (i, j) 向右或向下移动，即移动到 (i, j + 1)（右）或 (i + 1, j)（下）。
# 成本为目标单元格的值。
# 传送：你可以从任意单元格 (i, j) 传送到任意满足 grid[x][y] <= grid[i][j] 的单元格 (x, y)；
# 此移动的成本为 0。你最多可以传送 k 次。
# 返回从 (0, 0) 到达单元格 (m - 1, n - 1) 的 最小 总成本。
#
# !关键：第 k 层 会转移到 k + 1层, 但是第k + 1层的所有状态(x, y, k + 1)只会被第k层的状态(i, j, k)更新一次

from typing import List
from heapq import heappop, heappush


INF = int(1e18)


class Solution:
    def minCost(self, grid: List[List[int]], k: int) -> int:
        m, n = len(grid), len(grid[0])
        dist = [[[INF] * (k + 1) for _ in range(n)] for _ in range(m)]
        dist[0][0][0] = 0

        pq = [(0, 0, 0, 0)]
        cells = sorted((grid[i][j], i, j) for i in range(m) for j in range(n))
        ptrs = [0] * (k + 1)  # !每个 k 对应的 cells 中的索引，防止重复访问

        while pq:
            curDist, curK, x, y = heappop(pq)
            if x == m - 1 and y == n - 1:
                return curDist
            if curDist > dist[x][y][curK]:
                continue

            for nx, ny in ((x + 1, y), (x, y + 1)):
                if nx < m and ny < n:
                    nextDist = curDist + grid[nx][ny]
                    if nextDist < dist[nx][ny][curK]:
                        dist[nx][ny][curK] = nextDist
                        heappush(pq, (nextDist, curK, nx, ny))

            if curK < k:
                v = grid[x][y]
                p = ptrs[curK]
                while p < len(cells) and cells[p][0] <= v:
                    _, nx, ny = cells[p]
                    if curDist < dist[nx][ny][curK + 1]:
                        dist[nx][ny][curK + 1] = curDist
                        heappush(pq, (curDist, curK + 1, nx, ny))
                    p += 1
                ptrs[curK] = p

        return -1
