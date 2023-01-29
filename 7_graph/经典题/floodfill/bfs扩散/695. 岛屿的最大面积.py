from collections import deque
from typing import List

# bfs写floodfill更加方便


class Solution:
    def maxAreaOfIsland(self, grid: List[List[int]]) -> int:
        m, n = len(grid), len(grid[0])

        def bfs(i, j):
            res = 1
            queue = deque([(i, j)])
            grid[i][j] = 0
            while queue:
                x, y = queue.popleft()
                for dx, dy in ((1, 0), (-1, 0), (0, 1), (0, -1)):
                    nx, ny = x + dx, y + dy
                    if 0 <= nx < m and 0 <= ny < n and grid[nx][ny]:
                        grid[nx][ny] = 0
                        res += 1
                        queue.append((nx, ny))
            return res

        res = 0
        for i in range(m):
            for j in range(n):
                if grid[i][j]:
                    res = max(res, bfs(i, j))
        return res

