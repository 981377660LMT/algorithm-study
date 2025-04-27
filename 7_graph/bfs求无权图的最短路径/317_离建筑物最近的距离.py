# https://leetcode.cn/problems/shortest-distance-from-all-buildings/solutions/2416642/chi-jian-zhu-wu-zui-jin-de-ju-chi-by-lee-u22n/
# 给你一个由 0、1 和 2 组成的二维网格，其中：
#
# 0 代表你可以自由通过和选择建造的空地
# 1 代表你无法通行的建筑物
# 2 代表你无法通行的障碍物
# 通过调研，你希望从它出发能在 最短的距离和 内抵达周边全部的建筑物。请你计算出这个最佳的选址到周边全部建筑物的 最短距离和。
# 请你计算出这个最佳的选址到周边全部建筑物的 最短距离和。

from typing import List
from collections import deque

INF = int(1e18)


class Solution:
    def shortestDistance(self, grid: List[List[int]]) -> int:
        """
        For each building, do a BFS to compute distance to all reachable empty lands.
        Accumulate distances and reach counts at each empty cell.
        Finally, among cells reachable from all buildings, pick the one with minimal total distance.
        Time: O(K * MN), K = number of buildings.
        Space: O(MN).
        """
        if not grid or not grid[0]:
            return -1
        m, n = len(grid), len(grid[0])
        dist = [[0] * n for _ in range(m)]
        reach = [[0] * n for _ in range(m)]
        total_buildings = 0

        dirs = [(1, 0), (-1, 0), (0, 1), (0, -1)]

        def bfs(bi: int, bj: int) -> None:
            visited = [[False] * n for _ in range(m)]
            q = deque()
            q.append((bi, bj, 0))
            visited[bi][bj] = True

            while q:
                x, y, d = q.popleft()
                for dx, dy in dirs:
                    nx, ny = x + dx, y + dy
                    if 0 <= nx < m and 0 <= ny < n and not visited[nx][ny] and grid[nx][ny] == 0:
                        visited[nx][ny] = True
                        dist[nx][ny] += d + 1
                        reach[nx][ny] += 1
                        q.append((nx, ny, d + 1))

        for i in range(m):
            for j in range(n):
                if grid[i][j] == 1:
                    total_buildings += 1
                    bfs(i, j)

        res = INF
        for i in range(m):
            for j in range(n):
                if grid[i][j] == 0 and reach[i][j] == total_buildings:
                    res = min(res, dist[i][j])

        return res if res < INF else -1


if __name__ == "__main__":
    sol = Solution()
    grid1 = [[1, 0, 2, 0, 1], [0, 0, 0, 0, 0], [0, 0, 1, 0, 0]]
    # Expected output: 7
    print(sol.shortestDistance(grid1))
